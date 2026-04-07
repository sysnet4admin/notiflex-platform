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
	version     = "v0.5.0"
	valkeyClient valkey.Client
)

func main() {
	hostname, _ := os.Hostname()

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	valkeyPassword := os.Getenv("VALKEY_PASSWORD")

	// 파일 기반 Secret (CSI Driver)이 있으면 우선 사용
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			valkeyPassword = string(data)
			log.Printf("Valkey 비밀번호를 파일에서 로드: %s", pwFile)
		}
	}

	if valkeyAddr != "" {
		var err error
		for i := 0; i < 10; i++ {
			valkeyClient, err = valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{valkeyAddr},
				Password:    valkeyPassword,
			})
			if err == nil {
				log.Printf("Valkey 연결 성공: %s", valkeyAddr)
				break
			}
			log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Fatalf("Valkey 연결 실패: %v", err)
		}
		defer valkeyClient.Close()
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
		ctx := context.Background()
		var id int64

		if valkeyClient != nil {
			result := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build())
			val, err := result.AsInt64()
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
