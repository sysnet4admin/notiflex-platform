package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	version       = "v0.7.0"
	valkeyClient  valkey.Client
	kafkaProducer sarama.SyncProducer
	tracer        trace.Tracer
)

func initTracer() func() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return func() {}
	}

	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("OTel exporter 생성 실패: %v", err)
		return func() {}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("notiflex-api"),
			semconv.ServiceVersionKey.String(version),
		)),
	)
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("notiflex-api")
	log.Printf("OTel 트레이싱 초기화 완료: %s", endpoint)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tp.Shutdown(ctx)
	}
}

func main() {
	hostname, _ := os.Hostname()

	shutdown := initTracer()
	defer shutdown()

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	valkeyPassword := os.Getenv("VALKEY_PASSWORD")

	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			valkeyPassword = string(data)
			log.Printf("Valkey 비밀번호를 파일에서 로드: %s", pwFile)
		}
	}

	if valkeyAddr != "" {
		var err error
		for i := 0; i < 10; i++ {
			valkeyClient, err = valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{valkeyAddr},
				Password:    valkeyPassword,
			})
			if err == nil {
				log.Printf("Valkey 연결 성공: %s", valkeyAddr)
				break
			}
			log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Fatalf("Valkey 연결 실패: %v", err)
		}
		defer valkeyClient.Close()
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Retry.Max = 5

		var err error
		for i := 0; i < 10; i++ {
			kafkaProducer, err = sarama.NewSyncProducer([]string{kafkaBroker}, config)
			if err == nil {
				log.Printf("Kafka Producer 연결 성공: %s", kafkaBroker)
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 진행): %v", err)
		} else {
			defer kafkaProducer.Close()
			go startConsumer(kafkaBroker, hostname)
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": version,
			"pod":     hostname,
		})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "generate-id")
			defer span.End()
		}

		var id int64
		if valkeyClient != nil {
			if tracer != nil {
				_, vSpan := tracer.Start(ctx, "valkey-incr")
				defer vSpan.End()
			}
			result := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build())
			val, err := result.AsInt64()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			id = val
		}

		if kafkaProducer != nil {
			if tracer != nil {
				_, kSpan := tracer.Start(ctx, "kafka-produce")
				defer kSpan.End()
			}
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s","time":"%s"}`, id, hostname, time.Now().Format(time.RFC3339))),
			}
			_, _, err := kafkaProducer.SendMessage(msg)
			if err != nil {
				log.Printf("Kafka 전송 실패: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
		})
	})

	fmt.Println("Notiflex API server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func startConsumer(broker, hostname string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup([]string{broker}, "notiflex-consumer", config)
	if err != nil {
		log.Printf("Kafka Consumer Group 생성 실패: %v", err)
		return
	}
	defer group.Close()

	handler := &consumerHandler{hostname: hostname}
	for {
		if err := group.Consume(context.Background(), []string{"notifications"}, handler); err != nil {
			log.Printf("Kafka Consumer 에러: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

type consumerHandler struct {
	hostname string
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("[Kafka Consumer] topic=%s partition=%d offset=%d value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
