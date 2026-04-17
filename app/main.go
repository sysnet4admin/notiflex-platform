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

	"github.com/IBM/sarama"
	"github.com/valkey-io/valkey-go"
)

const version = "v0.6.0"

var (
	client   valkey.Client
	producer sarama.SyncProducer
	hostname string
)

func readSecret(path string) string {
	if path == "" {
		return ""
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func initKafka(broker string) error {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0 // ← run-50 피드백 반영: Kafka 4.x 상수
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	var err error
	for i := 0; i < 10; i++ {
		producer, err = sarama.NewSyncProducer([]string{broker}, cfg)
		if err == nil {
			break
		}
		log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return err
}

func startConsumer(broker, topic string) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer([]string{broker}, cfg)
	if err != nil {
		log.Printf("Kafka consumer 생성 실패: %v", err)
		return
	}
	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Consumer 파티션 구독 실패: %v", err)
		return
	}

	log.Printf("Kafka consumer 시작: topic=%s", topic)
	for msg := range pc.Messages() {
		log.Printf("Kafka consumer: received on %s: %s", topic, string(msg.Value))
	}
}

func main() {
	hostname, _ = os.Hostname()

	addr := os.Getenv("VALKEY_ADDR")
	password := readSecret(os.Getenv("VALKEY_PASSWORD_FILE"))
	if password == "" {
		password = os.Getenv("VALKEY_PASSWORD")
	}

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

	// Kafka
	broker := os.Getenv("KAFKA_BROKER")
	if broker != "" {
		if err := initKafka(broker); err != nil {
			log.Printf("Kafka 초기화 실패, 계속 진행: %v", err)
		}
		go startConsumer(broker, "notifications")
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		resp := client.Do(ctx, client.B().Incr().Key("notiflex:id").Build())
		n, err := resp.AsInt64()
		if err != nil {
			http.Error(w, fmt.Sprintf("valkey error: %v", err), 500)
			return
		}

		kafkaStatus := "disabled"
		if producer != nil {
			payload := fmt.Sprintf(`{"id":%d,"timestamp":"%s"}`, n, time.Now().UTC().Format(time.RFC3339))
			_, _, kerr := producer.SendMessage(&sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(payload),
			})
			if kerr != nil {
				kafkaStatus = fmt.Sprintf("error:%v", kerr)
			} else {
				kafkaStatus = "sent"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", n),
			"generated_by": hostname,
			"source":       "valkey+csi",
			"kafka":        kafkaStatus,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"version":      version,
			"generated_by": hostname,
		})
	})

	log.Printf("Starting server %s on :8080", version)
	http.ListenAndServe(":8080", nil)
}
