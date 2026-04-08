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
	version = "v0.4.0"
	client  valkey.Client
)

func main() {
	hostname, _ := os.Hostname()

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	valkeyPass := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			valkeyPass = string(data)
			log.Printf("Valkey password loaded from file: %s", pwFile)
		}
	}

	if valkeyAddr != "" {
		var err error
		for i := 0; i < 10; i++ {
			client, err = valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{valkeyAddr},
				Password:    valkeyPass,
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

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		var id int64
		if client != nil {
			ctx := context.Background()
			resp := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build())
			val, err := resp.AsInt64()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			id = val
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
		})
	})

	fmt.Println("Notiflex API server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
