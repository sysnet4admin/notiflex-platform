package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

var counter int64

func main() {
	hostname, _ := os.Hostname()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"version":"v0.2.0","service":"notiflex-api"}`)
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddInt64(&counter, 1)
		fmt.Fprintf(w, `{"id":%d,"pod":"%s"}`, id, hostname)
	})

	fmt.Println("Notiflex API server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
