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

var (
	version       = "v0.4.0"
	valkeyClient  valkey.Client
	kafkaProducer sarama.SyncProducer
)

func main() {
	hostname, _ := os.Hostname()

	// Valkey 연결 (10회 재시도)
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
	log.Println("Valkey 연결 성공")

	// Kafka Producer 연결 (10회 재시도)
	broker := os.Getenv("KAFKA_BROKER")
	if broker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		for i := 0; i < 10; i++ {
			kafkaProducer, err = sarama.NewSyncProducer([]string{broker}, config)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 진행): %v", err)
		} else {
			defer kafkaProducer.Close()
			log.Println("Kafka Producer 연결 성공")
			go consumeMessages(broker)
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		cmd := valkeyClient.B().Incr().Key("notiflex:id").Build()
		resp := valkeyClient.Do(context.Background(), cmd)
		id, _ := resp.AsInt64()

		if kafkaProducer != nil {
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s"}`, id, hostname)),
			}
			kafkaProducer.SendMessage(msg)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": hostname,
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": version,
			"pod":     hostname,
		})
	})

	fmt.Printf("Notiflex API %s starting on :8080\n", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func consumeMessages(broker string) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Printf("Kafka Consumer 생성 실패: %v", err)
		return
	}
	defer consumer.Close()

	partitions, _ := consumer.Partitions("notifications")
	for _, p := range partitions {
		pc, err := consumer.ConsumePartition("notifications", p, sarama.OffsetNewest)
		if err != nil {
			continue
		}
		go func(pc sarama.PartitionConsumer, partition int32) {
			for msg := range pc.Messages() {
				log.Printf("[Kafka] partition=%d offset=%d value=%s",
					partition, msg.Offset, string(msg.Value))
			}
		}(pc, p)
	}
	select {}
}
