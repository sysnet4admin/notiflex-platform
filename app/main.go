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
)

const (
	appVersion      = "v0.1.2"
	idCounterKey    = "notiflex:id:counter"
	valkeyRetries   = 10
	valkeyRetryWait = 3 * time.Second
)

var valkeyClient valkey.Client

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func idHandler(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	nextID, err := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key(idCounterKey).Build()).ToInt64()
	if err != nil {
		log.Printf("failed to get next id from valkey: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to generate id",
		})
		return
	}

	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "unknown"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":           strconv.FormatInt(nextID, 10),
		"generated_by": podName,
	})
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
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
