package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

const version = "v0.1.1"

var counter int64

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddInt64(&counter, 1)
		pod := os.Getenv("HOSTNAME")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": pod,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"version": version})
	})

	fmt.Printf("Notiflex API %s listening on :8080\n", version)
	http.ListenAndServe(":8080", nil)
}
