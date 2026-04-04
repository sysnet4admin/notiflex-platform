package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/valkey-io/valkey-go"
)

var version = "v0.4.0"

func main() {
	hostname, _ := os.Hostname()

	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
	}

	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "valkey-primary.notiflex.svc.cluster.local:6379"
	}

	var client valkey.Client
	var err error
	for i := 0; i < 10; i++ {
		client, err = valkey.NewClient(valkey.ClientOption{
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
	defer client.Close()

	// Kafka producer setup
	broker := os.Getenv("KAFKA_BROKER")
	var producer sarama.SyncProducer
	if broker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		for i := 0; i < 10; i++ {
			producer, err = sarama.NewSyncProducer([]string{broker}, config)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 실행): %v", err)
		} else {
			defer producer.Close()
			log.Printf("Kafka producer 연결 성공: %s", broker)

			// Start consumer goroutine
			go func() {
				consumer, err := sarama.NewConsumer([]string{broker}, sarama.NewConfig())
				if err != nil {
					log.Printf("Kafka consumer 생성 실패: %v", err)
					return
				}
				defer consumer.Close()
				pc, err := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
				if err != nil {
					log.Printf("Kafka partition consumer 생성 실패: %v", err)
					return
				}
				defer pc.Close()
				for msg := range pc.Messages() {
					log.Printf("[Kafka] 수신: %s", string(msg.Value))
				}
			}()
		}
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": "notification queued",
		})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		result := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build())
		id, err := result.AsInt64()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send to Kafka
		if producer != nil {
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s"}`, id, hostname)),
			}
			if _, _, err := producer.SendMessage(msg); err != nil {
				log.Printf("Kafka 전송 실패: %v", err)
			}
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
