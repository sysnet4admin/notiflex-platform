package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	valkey "github.com/valkey-io/valkey-go"
)

const appVersion = "v0.2.0"
const defaultValkeyAddr = "valkey-primary.notiflex.svc.cluster.local:6379"
const defaultKafkaBroker = "notiflex-kafka-kafka-bootstrap.kafka.svc.cluster.local:9092"
const idKey = "notiflex:id"
const kafkaTopic = "notifications"
const kafkaConsumerGroup = "notiflex-api"

type server struct {
	podName string
	client  valkey.Client
	producer sarama.SyncProducer
}

func newServer() (*server, error) {
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "unknown-pod"
	}

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = defaultValkeyAddr
	}
	valkeyPassword := os.Getenv("VALKEY_PASSWORD")
	if valkeyPasswordFile := os.Getenv("VALKEY_PASSWORD_FILE"); valkeyPasswordFile != "" {
		if data, err := os.ReadFile(valkeyPasswordFile); err == nil {
			valkeyPassword = strings.TrimSpace(string(data))
		} else {
			log.Printf("VALKEY_PASSWORD_FILE 읽기 실패, 환경변수 VALKEY_PASSWORD 사용: %v", err)
		}
	}

	var (
		client valkey.Client
		err    error
	)
	for i := 0; i < 10; i++ {
		client, err = valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{valkeyAddr},
			Password:    valkeyPassword,
		})
		if err == nil {
			break
		}
		log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("valkey 연결 실패: %w", err)
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = defaultKafkaBroker
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V4_1_0_0
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Retry.Max = 5
	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("kafka producer 연결 실패: %w", err)
	}

	return &server{
		podName:  podName,
		client:   client,
		producer: producer,
	}, nil
}

func (s *server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) idHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	id, err := s.client.Do(ctx, s.client.B().Incr().Key(idKey).Build()).AsInt64()
	if err != nil {
		log.Printf("failed to increment id in valkey: %v", err)
		http.Error(w, "failed to generate id", http.StatusInternalServerError)
		return
	}

	message := map[string]string{
		"id":           strconv.FormatInt(id, 10),
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"generated_by": s.podName,
	}
	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("failed to marshal kafka message: %v", err)
		http.Error(w, "failed to generate id", http.StatusInternalServerError)
		return
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.ByteEncoder(body),
	})
	if err != nil {
		log.Printf("failed to publish message to kafka: %v", err)
		http.Error(w, "failed to publish event", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":           strconv.FormatInt(id, 10),
		"generated_by": s.podName,
	})
}

func (s *server) versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"version":      appVersion,
		"generated_by": s.podName,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func main() {
	srv, err := newServer()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.client.Close()
	defer func() {
		if err := srv.producer.Close(); err != nil {
			log.Printf("kafka producer close failed: %v", err)
		}
	}()

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = defaultKafkaBroker
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go startKafkaConsumer(ctx, kafkaBroker)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", srv.healthHandler)
	mux.HandleFunc("/id", srv.idHandler)
	mux.HandleFunc("/version", srv.versionHandler)

	addr := ":8080"
	log.Printf("notiflex-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(fmt.Errorf("server stopped: %w", err))
	}
}

type consumerHandler struct{}

func (consumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (consumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Kafka consumer: received message on %s: %s", msg.Topic, string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

func startKafkaConsumer(ctx context.Context, kafkaBroker string) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0
	cfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup([]string{kafkaBroker}, kafkaConsumerGroup, cfg)
	if err != nil {
		log.Printf("Kafka consumer group init failed: %v", err)
		return
	}
	defer func() {
		if closeErr := consumerGroup.Close(); closeErr != nil {
			log.Printf("Kafka consumer group close failed: %v", closeErr)
		}
	}()

	handler := consumerHandler{}
	for {
		if err := consumerGroup.Consume(ctx, []string{kafkaTopic}, handler); err != nil {
			log.Printf("Kafka consumer loop error: %v", err)
			time.Sleep(3 * time.Second)
		}
		if ctx.Err() != nil {
			return
		}
	}
}
