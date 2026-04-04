package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var version = "v0.1.0"

func main() {
	hostname, _ := os.Hostname()

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

	fmt.Printf("Notiflex API %s starting on :8080\n", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
