package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const version = "0.6.0"

func initTracer() func() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		log.Println("OTEL_EXPORTER_OTLP_ENDPOINT not set, tracing disabled")
		return func() {}
	}

	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("트레이서 초기화 실패: %v", err)
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
	log.Println("OpenTelemetry 트레이서 초기화 완료")

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tp.Shutdown(ctx)
	}
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	hostname, _ := os.Hostname()

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = "valkey-primary.notiflex.svc.cluster.local:6379"
	}

	valkeyPassword := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			valkeyPassword = string(data)
		}
	}

	var client valkey.Client
	var err error
	for i := 0; i < 10; i++ {
		client, err = valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{valkeyAddr},
			Password:    valkeyPassword,
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

	// Kafka Producer setup
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	var producer sarama.SyncProducer
	if kafkaBroker != "" {
		brokers := strings.Split(kafkaBroker, ",")
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Retry.Max = 5
		for i := 0; i < 10; i++ {
			producer, err = sarama.NewSyncProducer(brokers, config)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 진행): %v", err)
		} else {
			log.Println("Kafka Producer 연결 성공")
			defer producer.Close()

			// Start Consumer
			go func() {
				consumer, err := sarama.NewConsumer(brokers, config)
				if err != nil {
					log.Printf("Kafka Consumer 생성 실패: %v", err)
					return
				}
				defer consumer.Close()
				pc, err := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
				if err != nil {
					log.Printf("Kafka Consumer 파티션 구독 실패: %v", err)
					return
				}
				defer pc.Close()
				log.Println("Kafka Consumer 시작")
				for msg := range pc.Messages() {
					log.Printf("Kafka consumer: received message on %s: %s", msg.Topic, string(msg.Value))
				}
			}()
		}
	}

	tracer := otel.Tracer("notiflex-api")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "GET /health")
		defer span.End()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "GET /id")
		defer span.End()
		cmd := client.B().Incr().Key("notiflex:id").Build()
		result := client.Do(ctx, cmd)
		id, err := result.AsInt64()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idStr := fmt.Sprintf("%d", id)

		// Kafka produce
		if producer != nil {
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":"%s","timestamp":"%s","pod":"%s"}`, idStr, time.Now().Format(time.RFC3339), hostname)),
			}
			if _, _, err := producer.SendMessage(msg); err != nil {
				log.Printf("Kafka 메시지 전송 실패: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           idStr,
			"generated_by": hostname,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "GET /version")
		defer span.End()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version":    version,
			"go_version": runtime.Version(),
			"pod":        hostname,
		})
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
