
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"


	"github.com/IBM/sarama"
)

var (
	idCounter   int64
	podName     = "unknown"
	kafkaBroker = "localhost:9092"
)

func main() {
	if name, err := os.Hostname(); err == nil {
		podName = name
	}

	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		kafkaBroker = broker
	}

	producer, err := newProducer()
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	go startConsumer()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		newID := atomic.AddInt64(&idCounter, 1)
		resp := map[string]interface{}{
			"id":        newID,
			"pod_name":  podName,
			"version":   "v0.3.0",
		}

		msg := &sarama.ProducerMessage{
			Topic: "notifications",
			Value: sarama.StringEncoder(fmt.Sprintf("New ID generated: %d", newID)),
		}
		_, _, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %v", err)
		} else {
			log.Printf("Sent message to Kafka: %s", msg.Value)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"version": "v0.3.0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Starting Notiflex API server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V4_1_0_0 // Specify Kafka version
	return sarama.NewSyncProducer(strings.Split(kafkaBroker, ","), config)
}

func startConsumer() {
	config := sarama.NewConfig()
	config.Version = sarama.V4_1_0_0 // Specify Kafka version
	consumer, err := sarama.NewConsumer(strings.Split(kafkaBroker, ","), config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()

	log.Println("Kafka consumer started. Waiting for messages...")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received Kafka message: %s", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Kafka consumer error: %v", err)
		}
	}
}
