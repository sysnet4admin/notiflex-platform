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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var version = "v0.5.0"

func initTracer() func() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return func() {}
	}
	ctx := context.Background()
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Printf("OTel exporter error: %v", err)
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
	return func() { tp.Shutdown(ctx) }
}

func main() {
	hostname, _ := os.Hostname()
	shutdown := initTracer()
	defer shutdown()
	tracer := otel.Tracer("notiflex-api")

	// Read Valkey password
	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = strings.TrimSpace(string(data))
		}
	}

	// Connect to Valkey with retry
	addr := os.Getenv("VALKEY_ADDR")
	var valkeyClient valkey.Client
	if addr != "" {
		for i := 0; i < 10; i++ {
			var err error
			valkeyClient, err = valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{addr},
				Password:    password,
			})
			if err == nil {
				log.Printf("Valkey connected to %s", addr)
				break
			}
			log.Printf("Valkey retry %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if valkeyClient == nil {
			log.Fatal("Failed to connect to Valkey after 10 retries")
		}
		defer valkeyClient.Close()
	}

	// Connect to Kafka with retry
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	var kafkaProducer sarama.SyncProducer
	if kafkaAddr != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		for i := 0; i < 10; i++ {
			var err error
			kafkaProducer, err = sarama.NewSyncProducer([]string{kafkaAddr}, config)
			if err == nil {
				log.Printf("Kafka connected to %s", kafkaAddr)
				break
			}
			log.Printf("Kafka retry %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if kafkaProducer == nil {
			log.Fatal("Failed to connect to Kafka after 10 retries")
		}
		defer kafkaProducer.Close()
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

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "notify")
		defer span.End()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		tenant := req["tenant"]
		message := req["message"]

		// Valkey INCR
		var id int64
		if valkeyClient != nil {
			_, vSpan := tracer.Start(ctx, "valkey-incr")
			resp := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:notify:"+tenant).Build())
			val, err := resp.AsInt64()
			if err == nil {
				id = val
			}
			vSpan.End()
		}

		// Kafka produce
		if kafkaProducer != nil {
			_, kSpan := tracer.Start(ctx, "kafka-produce")
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"tenant":"%s","message":"%s","id":%d}`, tenant, message, id)),
			}
			_, _, err := kafkaProducer.SendMessage(msg)
			if err != nil {
				log.Printf("[KAFKA] send error: %v", err)
			}
			kSpan.End()
		}

		log.Printf("[NOTIFY] tenant=%s message=%s id=%d", tenant, message, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "accepted",
			"tenant":  tenant,
			"message": message,
			"id":      id,
		})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		var id int64
		if valkeyClient != nil {
			resp := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build())
			val, err := resp.AsInt64()
			if err != nil {
				http.Error(w, "Valkey error", http.StatusInternalServerError)
				return
			}
			id = val
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": hostname,
		})
	})

	fmt.Printf("Notiflex API %s starting on :8080\n", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
