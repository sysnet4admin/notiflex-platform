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
)

var version = "v0.5.0"
var valkeyClient valkey.Client
var kafkaProducer sarama.SyncProducer
var tracer = otel.Tracer("notiflex-api")

func main() {
	shutdown := initTracer()
	defer shutdown()

	initValkey()
	initKafka()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	http.HandleFunc("/notify", notifyHandler)
	http.HandleFunc("/version", versionHandler)

	port := "8080"
	fmt.Printf("Notiflex API %s starting on :%s\n", version, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}

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
		log.Printf("OTel exporter init failed: %v", err)
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
	log.Printf("OTel tracer connected to %s", endpoint)

	return func() {
		tp.Shutdown(ctx)
	}
}

func initValkey() {
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
			log.Printf("Valkey connected to %s", addr)
			return
		}
		log.Printf("Valkey retry %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	log.Fatalf("Valkey connection failed after 10 retries: %v", err)
}

func initKafka() {
	broker := os.Getenv("KAFKA_ADDR")
	if broker == "" {
		return
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	for i := 0; i < 10; i++ {
		kafkaProducer, err = sarama.NewSyncProducer([]string{broker}, config)
		if err == nil {
			log.Printf("Kafka connected to %s", broker)
			return
		}
		log.Printf("Kafka retry %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	log.Printf("Kafka connection failed, continuing without Kafka: %v", err)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version,
	})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	pod := os.Getenv("HOSTNAME")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"pod":     pod,
		"version": version,
	})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "idHandler")
	defer span.End()

	ctx := context.Background()
	cmd := valkeyClient.B().Incr().Key("notiflex:id").Build()
	result := valkeyClient.Do(ctx, cmd)
	id, _ := result.AsInt64()

	pod := os.Getenv("HOSTNAME")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     pod,
		"version": version,
	})
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "notifyHandler")
	defer span.End()

	ctx := context.Background()
	cmd := valkeyClient.B().Incr().Key("notiflex:notification").Build()
	result := valkeyClient.Do(ctx, cmd)
	id, _ := result.AsInt64()

	pod := os.Getenv("HOSTNAME")

	if kafkaProducer != nil {
		msg := &sarama.ProducerMessage{
			Topic: "notifications",
			Value: sarama.StringEncoder(fmt.Sprintf(`{"notification_id":%d,"pod":"%s"}`, id, pod)),
		}
		kafkaProducer.SendMessage(msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notification_id": id,
		"pod":             pod,
		"version":         version,
		"status":          "queued",
	})
}
