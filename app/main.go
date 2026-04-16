package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/valkey-io/valkey-go"
)

const version = "0.3.0"

func main() {
	hostname, _ := os.Hostname()

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = "valkey-primary.notiflex.svc.cluster.local:6379"
	}

	valkeyPassword := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			valkeyPassword = string(data)
		}
	}

	var client valkey.Client
	var err error
	for i := 0; i < 10; i++ {
		client, err = valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{valkeyAddr},
			Password:    valkeyPassword,
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
	log.Println("Valkey 연결 성공")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		cmd := client.B().Incr().Key("notiflex:id").Build()
		result := client.Do(ctx, cmd)
		id, err := result.AsInt64()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version":    version,
			"go_version": runtime.Version(),
			"pod":        hostname,
		})
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
