package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

var counter atomic.Int64

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"version":"v0.1.1"}`)
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id := counter.Add(1)
		pod := os.Getenv("POD_NAME")
		if pod == "" {
			pod = "local"
		}
		fmt.Fprintf(w, `{"id":%d,"pod":"%s"}`+"\n", id, pod)
	})

	http.ListenAndServe(":8080", nil)
}
