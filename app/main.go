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

const version = "v0.6.0"

var (
	client   valkey.Client
	producer sarama.SyncProducer
)

func main() {
	addr := os.Getenv("VALKEY_ADDR")
	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
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

	// Kafka Producer 초기화
	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		for i := 0; i < 10; i++ {
			producer, err = sarama.NewSyncProducer([]string{broker}, cfg)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 진행): %v", err)
		} else {
			defer producer.Close()
			log.Printf("Kafka Producer 연결 완료: %s", broker)
			// Consumer goroutine
			go startConsumer(broker)
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		cmd := client.B().Incr().Key("notiflex:id").Build()
		id, err := client.Do(ctx, cmd).AsInt64()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		pod := os.Getenv("HOSTNAME")

		// Kafka에 이벤트 전송
		if producer != nil {
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s","version":"%s"}`, id, pod, version)),
			}
			_, _, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Kafka 전송 실패: %v", err)
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": pod,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"version": version})
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"ready": "true", "version": version})
	})

	fmt.Printf("Notiflex API %s listening on :8080\n", version)
	http.ListenAndServe(":8080", nil)
}

func startConsumer(broker string) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer([]string{broker}, cfg)
	if err != nil {
		log.Printf("Kafka Consumer 생성 실패: %v", err)
		return
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions("notifications")
	if err != nil {
		log.Printf("Kafka 파티션 조회 실패: %v", err)
		return
	}

	log.Printf("Kafka Consumer 시작: notifications (%d 파티션)", len(partitions))
	for _, p := range partitions {
		pc, err := consumer.ConsumePartition("notifications", p, sarama.OffsetNewest)
		if err != nil {
			log.Printf("파티션 %d 소비 실패: %v", p, err)
			continue
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				log.Printf("[Kafka] 수신: partition=%d offset=%d value=%s", msg.Partition, msg.Offset, string(msg.Value))
			}
		}(pc)
	}
	select {} // block forever
}
