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
	"go.opentelemetry.io/otel/sdk/trace"
)

const version = "v0.6.0"

var (
	client   valkey.Client
	producer sarama.SyncProducer
	tracer   = otel.Tracer("notiflex")
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

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	log.Printf("OpenTelemetry 트레이싱 활성화: %s", endpoint)

	return func() {
		tp.Shutdown(ctx)
	}
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	addr := os.Getenv("VALKEY_ADDR")
	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
	}

	var err error
	for i := 0; i < 10; i++ {
		client, err = valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{addr},
			Password:    password,
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

	// Kafka Producer 초기화
	broker := os.Getenv("KAFKA_BROKER")
	if broker != "" {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		for i := 0; i < 10; i++ {
			producer, err = sarama.NewSyncProducer([]string{broker}, cfg)
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
			log.Printf("Kafka Producer 연결 성공: %s", broker)
			go startConsumer(broker)
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "GET /health")
		defer span.End()
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "GET /id")
		defer span.End()

		_, valkeySpan := tracer.Start(ctx, "valkey.INCR")
		cmd := client.B().Incr().Key("notiflex:id").Build()
		id, err := client.Do(ctx, cmd).AsInt64()
		valkeySpan.End()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		pod := os.Getenv("HOSTNAME")

		// Kafka로 이벤트 전송
		if producer != nil {
			_, kafkaSpan := tracer.Start(ctx, "kafka.produce")
			msg := fmt.Sprintf(`{"id":%d,"pod":"%s","time":"%s"}`, id, pod, time.Now().Format(time.RFC3339))
			_, _, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(msg),
			})
			if err != nil {
				log.Printf("Kafka 전송 실패: %v", err)
			}
			kafkaSpan.End()
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": pod,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "GET /version")
		defer span.End()
		json.NewEncoder(w).Encode(map[string]string{"version": version})
	})

	fmt.Printf("Notiflex API %s listening on :8080\n", version)
	http.ListenAndServe(":8080", nil)
}

func startConsumer(broker string) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup([]string{broker}, "notiflex-consumer", cfg)
	if err != nil {
		log.Printf("Kafka Consumer 생성 실패: %v", err)
		return
	}
	defer group.Close()

	handler := &consumerHandler{}
	for {
		if err := group.Consume(context.Background(), []string{"notifications"}, handler); err != nil {
			log.Printf("Kafka Consumer 에러: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

type consumerHandler struct{}

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
