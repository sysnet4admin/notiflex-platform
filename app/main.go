package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	valkey "github.com/valkey-io/valkey-go"
)

const appVersion = "v0.2.0"
const defaultValkeyAddr = "valkey-primary.notiflex.svc.cluster.local:6379"
const idKey = "notiflex:id"

type server struct {
	podName string
	client  valkey.Client
}

func newServer() (*server, error) {
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "unknown-pod"
	}

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = defaultValkeyAddr
	}
	valkeyPassword := os.Getenv("VALKEY_PASSWORD")
	if valkeyPasswordFile := os.Getenv("VALKEY_PASSWORD_FILE"); valkeyPasswordFile != "" {
		if data, err := os.ReadFile(valkeyPasswordFile); err == nil {
			valkeyPassword = strings.TrimSpace(string(data))
		} else {
			log.Printf("VALKEY_PASSWORD_FILE 읽기 실패, 환경변수 VALKEY_PASSWORD 사용: %v", err)
		}
	}

	var (
		client valkey.Client
		err    error
	)
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
		return nil, fmt.Errorf("valkey 연결 실패: %w", err)
	}

	return &server{
		podName: podName,
		client:  client,
	}, nil
}

func (s *server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) idHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	id, err := s.client.Do(ctx, s.client.B().Incr().Key(idKey).Build()).AsInt64()
	if err != nil {
		log.Printf("failed to increment id in valkey: %v", err)
		http.Error(w, "failed to generate id", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":           strconv.FormatInt(id, 10),
		"generated_by": s.podName,
	})
}

func (s *server) versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"version":      appVersion,
		"generated_by": s.podName,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func main() {
	srv, err := newServer()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.client.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", srv.healthHandler)
	mux.HandleFunc("/id", srv.idHandler)
	mux.HandleFunc("/version", srv.versionHandler)

	addr := ":8080"
	log.Printf("notiflex-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(fmt.Errorf("server stopped: %w", err))
	}
}
