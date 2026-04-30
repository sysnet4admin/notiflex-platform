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

var valkeyClient valkey.Client
var kafkaProducer sarama.SyncProducer

func main() {
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
	}

	var err error
	for i := 0; i < 10; i++ {
		valkeyClient, err = valkey.NewClient(valkey.ClientOption{
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
	defer valkeyClient.Close()

	broker := os.Getenv("KAFKA_BROKER")
	if broker != "" {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		cfg.Version = sarama.V4_1_0_0
		kafkaProducer, err = sarama.NewSyncProducer([]string{broker}, cfg)
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속): %v", err)
		} else {
			defer kafkaProducer.Close()
			go consumeKafka(broker)
		}
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	fmt.Println("Notiflex API server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

func consumeKafka(broker string) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0
	consumer, err := sarama.NewConsumer([]string{broker}, cfg)
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
		log.Printf("[Kafka] 수신: key=%s value=%s", string(msg.Key), string(msg.Value))
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": "v0.3.0"})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	result, err := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build()).AsInt64()
	if err != nil {
		http.Error(w, "Valkey error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	pod := os.Getenv("POD_NAME")
	if pod == "" {
		pod = "unknown"
	}

	if kafkaProducer != nil {
		msg := &sarama.ProducerMessage{
			Topic: "notifications",
			Key:   sarama.StringEncoder(fmt.Sprintf("id-%d", result)),
			Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s"}`, result, pod)),
		}
		if _, _, err := kafkaProducer.SendMessage(msg); err != nil {
			log.Printf("[Kafka] 전송 실패: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": result, "pod": pod})
}
