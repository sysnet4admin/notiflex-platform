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

var version = "v0.4.0"
var valkeyClient valkey.Client

func main() {
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "valkey-primary.notiflex.svc.cluster.local:6379"
	}

	var password string
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
			log.Printf("Valkey 비밀번호를 파일에서 읽음: %s", pwFile)
		}
	}
	if password == "" {
		password = os.Getenv("VALKEY_PASSWORD")
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

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	http.HandleFunc("/notify", notifyHandler)

	log.Printf("Notiflex API %s starting on :8080", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version,
	})
}

func nextID(ctx context.Context) (int64, error) {
	cmd := valkeyClient.B().Incr().Key("notiflex:id").Build()
	return valkeyClient.Do(ctx, cmd).AsInt64()
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := nextID(r.Context())
	if err != nil {
		http.Error(w, "ID generation failed", http.StatusInternalServerError)
		return
	}
	pod := os.Getenv("HOSTNAME")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     pod,
		"message": "notification sent",
		"version": version,
	})
	fmt.Printf("Notification %d sent from pod %s\n", id, pod)
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	id, err := nextID(r.Context())
	if err != nil {
		http.Error(w, "ID generation failed", http.StatusInternalServerError)
		return
	}
	pod := os.Getenv("HOSTNAME")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     pod,
		"version": version,
	})
	fmt.Printf("Generated ID: %d on pod %s\n", id, pod)
}
