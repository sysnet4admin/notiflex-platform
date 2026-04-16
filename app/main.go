package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const version = "0.4.0"

var valkeyAddr = getenv("VALKEY_ADDR", "valkey-primary.notiflex.svc.cluster.local:6379")

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// valkeyIncr calls INCR key over RESP protocol without external deps.
func valkeyIncr(key string) (int64, error) {
	conn, err := net.DialTimeout("tcp", valkeyAddr, 2*time.Second)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	fmt.Fprintf(conn, "*2\r\n$4\r\nINCR\r\n$%d\r\n%s\r\n", len(key), key)
	r := bufio.NewReader(conn)
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, ":") {
		return 0, fmt.Errorf("unexpected response: %q", line)
	}
	return strconv.ParseInt(line[1:], 10, 64)
}

func main() {
	hostname, _ := os.Hostname()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version":      version,
			"generated_by": hostname,
		})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id, err := valkeyIncr("notiflex:id")
		if err != nil {
			http.Error(w, "valkey error: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
			"source":       "valkey",
		})
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
