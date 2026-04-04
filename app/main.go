package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var version = "v0.2.0"

func main() {
	hostname, _ := os.Hostname()
	var counter int64

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

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		tenant := req["tenant"]
		message := req["message"]
		log.Printf("[NOTIFY] tenant=%s message=%s", tenant, message)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"tenant":  tenant,
			"message": message,
		})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddInt64(&counter, 1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": hostname,
		})
	})

	fmt.Printf("Notiflex API %s starting on :8080\n", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
