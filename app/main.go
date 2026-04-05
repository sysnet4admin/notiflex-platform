package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/valkey-io/valkey-go"
)

var version = "v0.3.0"
var valkeyClient valkey.Client

func main() {
	initValkey()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	http.HandleFunc("/notify", notifyHandler)

	port := "8080"
	fmt.Printf("Notiflex API %s starting on :%s\n", version, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version,
	})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
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
	ctx := context.Background()
	cmd := valkeyClient.B().Incr().Key("notiflex:notification").Build()
	result := valkeyClient.Do(ctx, cmd)
	id, _ := result.AsInt64()

	pod := os.Getenv("HOSTNAME")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notification_id": id,
		"pod":             pod,
		"version":         version,
		"status":          "queued",
	})
}
