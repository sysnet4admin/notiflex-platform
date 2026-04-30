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

var client valkey.Client

func main() {
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("VALKEY_PASSWORD")

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

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	fmt.Println("Notiflex API server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": "v0.2.0"})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	result, err := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build()).AsInt64()
	if err != nil {
		http.Error(w, "Valkey error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	pod := os.Getenv("POD_NAME")
	if pod == "" {
		pod = "unknown"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": result, "pod": pod})
}
