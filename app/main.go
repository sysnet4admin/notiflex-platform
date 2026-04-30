package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "sync/atomic"
)

var idCounter int64

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "OK")
    })

    http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintln(w, `{"version": "v0.1.1"}`)
    })

    http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
        podName := os.Getenv("HOSTNAME")
        id := atomic.AddInt64(&idCounter, 1)
        fmt.Fprintf(w, "ID: %d, Pod: %s", id, podName)
        fmt.Fprintln(w)
    })

    log.Println("Starting server on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
