package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

var counter int64

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	fmt.Println("Notiflex API server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": "v0.1.1"})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	id := atomic.AddInt64(&counter, 1)
	pod := os.Getenv("POD_NAME")
	if pod == "" {
		pod = "unknown"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "pod": pod})
}
