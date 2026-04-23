package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	appVersion      = "v0.1.3"
	idCounterKey    = "notiflex:id:counter"
	valkeyRetries   = 10
	valkeyRetryWait = 3 * time.Second
)

var valkeyClient valkey.Client
var tracer = otel.Tracer("notiflex-api")

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GET /health")
	defer span.End()

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "GET /id")
	defer span.End()

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	nextID, err := valkeyClient.Do(timeoutCtx, valkeyClient.B().Incr().Key(idCounterKey).Build()).ToInt64()
	if err != nil {
		span.RecordError(err)
		log.Printf("failed to get next id from valkey: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to generate id",
		})
		return
	}
	span.SetAttributes(attribute.Int64("notiflex.id.value", nextID))

	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "unknown"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":           strconv.FormatInt(nextID, 10),
		"generated_by": podName,
	})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GET /version")
	defer span.End()

	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "unknown"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"service": "notiflex-api",
		"version": appVersion,
		"pod":     podName,
	})
}

func initTracerProvider(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "tempo.monitoring.svc.cluster.local:4317"
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("otlp exporter init failed: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			attribute.String("service.name", "notiflex-api"),
			attribute.String("service.version", appVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("otel resource init failed: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func newValkeyClientFromEnv() (valkey.Client, error) {
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "valkey-primary.notiflex.svc.cluster.local:6379"
	}
	password := os.Getenv("VALKEY_PASSWORD")

	for i := 0; i < valkeyRetries; i++ {
		client, err := valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{addr},
			Password:    password,
		})
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, pingErr := client.Do(ctx, client.B().Ping().Build()).ToString()
			cancel()
			if pingErr == nil {
				return client, nil
			}
			client.Close()
			err = pingErr
		}

		log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(valkeyRetryWait)
	}

	return nil, fmt.Errorf("valkey 연결 실패: %s", addr)
}

func main() {
	shutdown, err := initTracerProvider(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if shutdownErr := shutdown(ctx); shutdownErr != nil {
			log.Printf("failed to shutdown tracer provider: %v", shutdownErr)
		}
	}()

	client, err := newValkeyClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	valkeyClient = client
	defer valkeyClient.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/id", idHandler)
	mux.HandleFunc("/version", versionHandler)

	addr := ":8080"
	log.Printf("notiflex api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
