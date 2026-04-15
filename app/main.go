package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/valkey-io/valkey-go"
)

var valkeyClient valkey.Client
var kafkaProducer sarama.SyncProducer

func main() {
	hostname, _ := os.Hostname()
	addr := os.Getenv("VALKEY_ADDR")
	password := os.Getenv("VALKEY_PASSWORD")

	// 파일 기반 비밀번호 (CSI Secret Manager 마운트)
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
	}

	// Valkey 연결 (10회 재시도, 3초 간격)
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

	// Kafka Producer 초기화
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.RequiredAcks = sarama.WaitForAll

		brokers := strings.Split(kafkaBroker, ",")
		for i := 0; i < 10; i++ {
			kafkaProducer, err = sarama.NewSyncProducer(brokers, config)
			if err == nil {
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (Producer 비활성화): %v", err)
		} else {
			defer kafkaProducer.Close()
			log.Println("Kafka Producer 연결 성공")

			// Consumer goroutine
			go startConsumer(brokers, hostname)
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"version":"v0.6.0","service":"notiflex-api"}`)
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		cmd := valkeyClient.B().Incr().Key("notiflex:id").Build()
		result := valkeyClient.Do(ctx, cmd)
		id, err := result.AsInt64()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}

		// Kafka produce
		if kafkaProducer != nil {
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":"%d","timestamp":"%s","pod":"%s"}`, id, time.Now().Format(time.RFC3339), hostname)),
			}
			_, _, err := kafkaProducer.SendMessage(msg)
			if err != nil {
				log.Printf("Kafka send error: %v", err)
			}
		}

		fmt.Fprintf(w, `{"id":"%d","generated_by":"%s"}`, id, hostname)
	})

	fmt.Println("Notiflex API server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

func startConsumer(brokers []string, hostname string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	for i := 0; i < 10; i++ {
		group, err := sarama.NewConsumerGroup(brokers, "notiflex-consumer", config)
		if err != nil {
			log.Printf("Kafka Consumer 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}
		defer group.Close()
		log.Println("Kafka Consumer 연결 성공")

		handler := &consumerHandler{hostname: hostname}
		for {
			if err := group.Consume(context.Background(), []string{"notifications"}, handler); err != nil {
				log.Printf("Kafka consumer error: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

type consumerHandler struct {
	hostname string
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Kafka consumer [%s]: received on %s: %s", h.hostname, msg.Topic, string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
