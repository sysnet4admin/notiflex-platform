package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"
)

var version = "v0.3.0"

func main() {
	hostname, _ := os.Hostname()

	// Read Valkey password
	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = strings.TrimSpace(string(data))
		}
	}

	// Connect to Valkey with retry
	addr := os.Getenv("VALKEY_ADDR")
	var valkeyClient valkey.Client
	if addr != "" {
		for i := 0; i < 10; i++ {
			var err error
			valkeyClient, err = valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{addr},
				Password:    password,
			})
			if err == nil {
				log.Printf("Valkey connected to %s", addr)
				break
			}
			log.Printf("Valkey retry %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if valkeyClient == nil {
			log.Fatal("Failed to connect to Valkey after 10 retries")
		}
		defer valkeyClient.Close()
	}

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
		ctx := context.Background()
		var id int64
		if valkeyClient != nil {
			resp := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build())
			val, err := resp.AsInt64()
			if err != nil {
				http.Error(w, "Valkey error", http.StatusInternalServerError)
				return
			}
			id = val
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": hostname,
		})
	})

	fmt.Printf("Notiflex API %s starting on :8080\n", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
