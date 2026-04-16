package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

const version = "0.5.0"

var (
	valkeyAddr  = getenv("VALKEY_ADDR", "valkey-primary.notiflex.svc.cluster.local:6379")
	kafkaBroker = getenv("KAFKA_BROKER", "notiflex-kafka-kafka-bootstrap.kafka.svc.cluster.local:9092")
	kafkaTopic  = getenv("KAFKA_TOPIC", "notifications")
)

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

var producer sarama.SyncProducer

func initKafka() error {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Return.Successes = true
	cfg.Version = sarama.V4_0_0_0
	p, err := sarama.NewSyncProducer([]string{kafkaBroker}, cfg)
	if err != nil {
		return err
	}
	producer = p
	return nil
}

func publishEvent(id int64) error {
	if producer == nil {
		return fmt.Errorf("producer not ready")
	}
	payload, _ := json.Marshal(map[string]any{
		"id":    id,
		"ts":    time.Now().UTC().Format(time.RFC3339),
		"event": "id.generated",
	})
	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafkaTopic,
		Key:   sarama.StringEncoder(strconv.FormatInt(id, 10)),
		Value: sarama.ByteEncoder(payload),
	})
	return err
}

// consumerLoop runs in a goroutine and logs every message on the topic.
func consumerLoop(ctx context.Context) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Version = sarama.V4_0_0_0
	for {
		c, err := sarama.NewConsumer([]string{kafkaBroker}, cfg)
		if err != nil {
			log.Printf("kafka consumer connect: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		partitions, err := c.Partitions(kafkaTopic)
		if err != nil {
			log.Printf("kafka partitions: %v", err)
			c.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("kafka consumer started on topic=%s partitions=%v", kafkaTopic, partitions)
		for _, part := range partitions {
			pc, err := c.ConsumePartition(kafkaTopic, part, sarama.OffsetNewest)
			if err != nil {
				log.Printf("consume partition %d: %v", part, err)
				continue
			}
			go func(part int32, pc sarama.PartitionConsumer) {
				for msg := range pc.Messages() {
					log.Printf("kafka recv topic=%s partition=%d offset=%d key=%s value=%s",
						msg.Topic, part, msg.Offset, string(msg.Key), string(msg.Value))
				}
			}(part, pc)
		}
		<-ctx.Done()
		c.Close()
		return
	}
}

func main() {
	hostname, _ := os.Hostname()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Kafka 초기화는 실패해도 API는 살아 있게 유지
	if err := initKafka(); err != nil {
		log.Printf("kafka producer init failed: %v (continuing)", err)
	}
	go consumerLoop(ctx)

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
		kafkaStatus := "sent"
		if err := publishEvent(id); err != nil {
			log.Printf("kafka publish id=%d err=%v", id, err)
			kafkaStatus = "skipped:" + err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
			"source":       "valkey",
			"kafka":        kafkaStatus,
		})
	})

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
