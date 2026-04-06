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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		endpoint = "tempo.monitoring.svc.cluster.local:4317"
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
	tracer = tp.Tracer("notiflex")

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
	if kafkaBroker == "" {
		kafkaBroker = "notiflex-kafka-kafka-bootstrap.kafka.svc.cluster.local:9092"
	}

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
		log.Printf("Kafka 연결 실패 (Producer 비활성): %v", err)
	} else {
		defer producer.Close()
	}

	// Kafka Consumer (background)
	go func() {
		consumer, err := sarama.NewConsumer([]string{kafkaBroker}, nil)
		if err != nil {
			log.Printf("Kafka Consumer 시작 실패: %v", err)
			return
		}
		defer consumer.Close()
		pc, err := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Kafka Partition Consumer 실패: %v", err)
			return
		}
		defer pc.Close()
		for msg := range pc.Messages() {
			log.Printf("[Kafka Consumer] %s", string(msg.Value))
		}
	}()

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

		resp := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build())
		id, _ := resp.ToInt64()

		if producer != nil {
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s","version":"%s"}`, id, hostname, version)),
			}
			producer.SendMessage(msg)
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
