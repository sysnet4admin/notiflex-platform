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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	version  = "v0.7.0"
	client   valkey.Client
	producer sarama.SyncProducer
	tracer   = otel.Tracer("notiflex-api")
)

func initTracer() func() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return func() {}
	}

	ctx := context.Background()
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("OTel gRPC 연결 실패: %v", err)
		return func() {}
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Printf("OTel exporter 생성 실패: %v", err)
		return func() {}
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("notiflex-api"),
			semconv.ServiceVersionKey.String(version),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	log.Println("OpenTelemetry 트레이서 초기화 완료")

	return func() {
		tp.Shutdown(ctx)
	}
}

func main() {
	hostname, _ := os.Hostname()

	shutdown := initTracer()
	defer shutdown()

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	valkeyPass := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			valkeyPass = string(data)
			log.Printf("Valkey password loaded from file: %s", pwFile)
		}
	}

	if valkeyAddr != "" {
		var err error
		for i := 0; i < 10; i++ {
			client, err = valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{valkeyAddr},
				Password:    valkeyPass,
			})
			if err == nil {
				break
			}
			log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Fatalf("Valkey 연결 실패: %v", err)
		}
		defer client.Close()
		log.Println("Valkey 연결 성공")
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.RequiredAcks = sarama.WaitForAll
		var err error
		for i := 0; i < 10; i++ {
			producer, err = sarama.NewSyncProducer([]string{kafkaBroker}, config)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 진행): %v", err)
		} else {
			defer producer.Close()
			log.Println("Kafka Producer 연결 성공")
			go consumeMessages(kafkaBroker, hostname)
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "GET /version")
		defer span.End()
		span.SetAttributes(attribute.String("version", version))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": version,
			"pod":     hostname,
		})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "GET /id")
		defer span.End()

		var id int64
		if client != nil {
			_, valkeySpan := tracer.Start(ctx, "valkey.INCR")
			resp := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build())
			val, err := resp.AsInt64()
			valkeySpan.End()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			id = val
		}

		if producer != nil {
			_, kafkaSpan := tracer.Start(ctx, "kafka.produce")
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s","time":"%s"}`, id, hostname, time.Now().Format(time.RFC3339))),
			}
			_, _, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Kafka 전송 실패: %v", err)
			}
			kafkaSpan.End()
		}

		span.SetAttributes(attribute.Int64("id", id))
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

func consumeMessages(broker, hostname string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup([]string{broker}, "notiflex-consumer", config)
	if err != nil {
		log.Printf("Kafka Consumer 생성 실패: %v", err)
		return
	}
	defer group.Close()
	log.Println("Kafka Consumer 시작")

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
