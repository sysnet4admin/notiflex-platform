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

	"github.com/valkey-io/valkey-go"
)

const version = "v0.4.0"

var client valkey.Client

func readSecret(path string) string {
	if path == "" {
		return ""
	}
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("warn: cannot read %s: %v", path, err)
		return ""
	}
	return strings.TrimSpace(string(b))
}

func main() {
	hostname, _ := os.Hostname()
	addr := os.Getenv("VALKEY_ADDR")

	// 파일 기반 → 환경변수 fallback
	password := readSecret(os.Getenv("VALKEY_PASSWORD_FILE"))
	if password == "" {
		password = os.Getenv("VALKEY_PASSWORD")
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
	log.Printf("Valkey connected: %s", addr)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		resp := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build())
		n, err := resp.AsInt64()
		if err != nil {
			http.Error(w, fmt.Sprintf("valkey error: %v", err), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", n),
			"generated_by": hostname,
			"source":       "valkey+csi",
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version":      version,
			"generated_by": hostname,
		})
	})

	log.Printf("Starting server %s on :8080", version)
	http.ListenAndServe(":8080", nil)
}
