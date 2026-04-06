package main

import (
	"context"
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
)

var (
	version  = "v0.6.0"
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
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("OTel exporter 초기화 실패: %v", err)
		return func() {}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("notiflex-api"),
			semconv.ServiceVersionKey.String(version),
		)),
	)
	otel.SetTracerProvider(tp)
	return func() { tp.Shutdown(ctx) }
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	hostname, _ := os.Hostname()
	addr := os.Getenv("VALKEY_ADDR")
	password := os.Getenv("VALKEY_PASSWORD")

	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
	}

	if addr == "" {
		addr = "valkey-primary.notiflex.svc.cluster.local:6379"
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

	// Kafka Producer
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		for i := 0; i < 10; i++ {
			producer, err = sarama.NewSyncProducer([]string{kafkaBroker}, config)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka Producer 연결 실패 (계속 진행): %v", err)
		} else {
			defer producer.Close()
		}

		// Consumer goroutine
		go func() {
			consumer, err := sarama.NewConsumer([]string{kafkaBroker}, nil)
			if err != nil {
				log.Printf("Kafka Consumer 시작 실패: %v", err)
				return
			}
			defer consumer.Close()
			partConsumer, err := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Kafka Partition Consumer 시작 실패: %v", err)
				return
			}
			defer partConsumer.Close()
			for msg := range partConsumer.Messages() {
				log.Printf("[Kafka Consumer] %s", string(msg.Value))
			}
		}()
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "GET /version")
		defer span.End()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"version":"%s"}`, version)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "GET /id")
		defer span.End()

		_, valkeySpan := tracer.Start(ctx, "valkey.incr")
		resp := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build())
		id, err := resp.ToInt64()
		valkeySpan.End()

		if err != nil {
			span.SetAttributes(attribute.String("error", err.Error()))
			http.Error(w, fmt.Sprintf("Valkey error: %v", err), 500)
			return
		}

		span.SetAttributes(attribute.Int64("notiflex.id", id))

		if producer != nil {
			_, kafkaSpan := tracer.Start(ctx, "kafka.produce")
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"event":"id_created","id":%d,"pod":"%s"}`, id, hostname)),
			}
			producer.SendMessage(msg)
			kafkaSpan.End()
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":%d,"pod":"%s"}`, id, hostname)
	})

	fmt.Println("Notiflex API starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
