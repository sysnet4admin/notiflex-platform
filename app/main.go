package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var (
	idCounter int64
	podName   = "unknown"
)

func main() {
	// A change to trigger CI
	if name, err := os.Hostname(); err == nil {
		podName = name
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		newID := atomic.AddInt64(&idCounter, 1)
		resp := map[string]interface{}{
			"id":        newID,
			"pod_name":  podName,
			"version": "v0.2.0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"version": "v0.2.0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Starting Notiflex API server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
