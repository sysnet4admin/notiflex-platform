// notiflex-api HTTP server
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
)

const appVersion = "v0.1.1"

var counter int64

func main() {
	hostname, _ := os.Hostname()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddInt64(&counter, 1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
			"version":      appVersion,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": appVersion,
			"go":      runtime.Version(),
			"pod":     hostname,
		})
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
