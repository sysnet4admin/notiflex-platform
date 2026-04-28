package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/valkey-io/valkey-go"
)

var vk valkey.Client

func main() {
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("VALKEY_PASSWORD")

	var err error
	for i := 0; i < 10; i++ {
		vk, err = valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{addr},
			Password:    password,
		})
		if err == nil {
			break
		}
		log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("Valkey 연결 실패: %v", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"version":"v0.2.0"}`)
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		id, err := vk.Do(ctx, vk.B().Incr().Key("notiflex:id").Build()).AsInt64()
		if err != nil {
			http.Error(w, "Valkey error: "+err.Error(), 500)
			return
		}
		pod := os.Getenv("POD_NAME")
		if pod == "" {
			pod = "local"
		}
		fmt.Fprintf(w, `{"id":%d,"pod":"%s"}`+"\n", id, pod)
	})

	http.ListenAndServe(":8080", nil)
}
