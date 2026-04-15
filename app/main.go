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

var (
	version     = "v0.4.0"
	valkeyClient valkey.Client
)

func main() {
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

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
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
		ctx := context.Background()
		result := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build())
		id, err := result.AsInt64()
		if err != nil {
			http.Error(w, fmt.Sprintf("Valkey error: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
		})
	})

	log.Println("Notiflex API v0.4.0 starting on :8080")
	http.ListenAndServe(":8080", nil)
}
