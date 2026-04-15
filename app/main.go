package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

var (
	version       = "v0.6.0"
	valkeyClient  valkey.Client
	kafkaProducer sarama.AsyncProducer
	tracer        = otel.Tracer("notiflex-api")
)

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
		log.Printf("Failed to create trace exporter: %v", err)
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
	log.Println("OpenTelemetry tracing 초기화 완료")

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("TracerProvider shutdown error: %v", err)
		}
	}
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	hostname, _ := os.Hostname()

	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "valkey-primary.notiflex.svc.cluster.local:6379"
	}
	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
	}

	var err error
	for i := 0; i < 10; i++ {
		valkeyClient, err = valkey.NewClient(valkey.ClientOption{
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
	defer valkeyClient.Close()
	log.Println("Valkey 연결 성공")

	// Kafka Producer 설정
	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = false
		config.Producer.Return.Errors = true

		brokers := strings.Split(broker, ",")
		for i := 0; i < 10; i++ {
			kafkaProducer, err = sarama.NewAsyncProducer(brokers, config)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 실행): %v", err)
		} else {
			defer kafkaProducer.Close()
			log.Println("Kafka Producer 연결 성공")
			go func() {
				for e := range kafkaProducer.Errors() {
					log.Printf("Kafka 전송 에러: %v", e.Err)
				}
			}()
		}
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "GET /version")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": version,
			"pod":     hostname,
		})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "GET /id")
		defer span.End()

		_, valkeySpan := tracer.Start(ctx, "valkey.INCR")
		result := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build())
		id, err := result.AsInt64()
		valkeySpan.End()
		if err != nil {
			http.Error(w, fmt.Sprintf("Valkey error: %v", err), http.StatusInternalServerError)
			return
		}

		idStr := fmt.Sprintf("%d", id)

		// Kafka에 메시지 전송
		if kafkaProducer != nil {
			_, kafkaSpan := tracer.Start(ctx, "kafka.produce")
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":"%s","pod":"%s","timestamp":"%s"}`, idStr, hostname, time.Now().Format(time.RFC3339))),
			}
			kafkaProducer.Input() <- msg
			kafkaSpan.End()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           idStr,
			"generated_by": hostname,
		})
	})

	log.Println("Notiflex API v0.6.0 starting on :8080")
	http.ListenAndServe(":8080", nil)
}
