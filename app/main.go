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
	version      = "v0.6.0"
	valkeyClient valkey.Client
	kafkaProducer sarama.SyncProducer
)

func main() {
	hostname, _ := os.Hostname()

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	valkeyPassword := os.Getenv("VALKEY_PASSWORD")

	// 파일 기반 Secret (CSI Driver)이 있으면 우선 사용
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			valkeyPassword = string(data)
			log.Printf("Valkey 비밀번호를 파일에서 로드: %s", pwFile)
		}
	}

	if valkeyAddr != "" {
		var err error
		for i := 0; i < 10; i++ {
			valkeyClient, err = valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{valkeyAddr},
				Password:    valkeyPassword,
			})
			if err == nil {
				log.Printf("Valkey 연결 성공: %s", valkeyAddr)
				break
			}
			log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Fatalf("Valkey 연결 실패: %v", err)
		}
		defer valkeyClient.Close()
	}

	// Kafka Producer 초기화
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Retry.Max = 5

		var err error
		for i := 0; i < 10; i++ {
			kafkaProducer, err = sarama.NewSyncProducer([]string{kafkaBroker}, config)
			if err == nil {
				log.Printf("Kafka Producer 연결 성공: %s", kafkaBroker)
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속 진행): %v", err)
		} else {
			defer kafkaProducer.Close()
			// Consumer 시작 (백그라운드)
			go startConsumer(kafkaBroker, hostname)
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

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		var id int64

		if valkeyClient != nil {
			result := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build())
			val, err := result.AsInt64()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			id = val
		}

		// Kafka에 이벤트 발행
		if kafkaProducer != nil {
			msg := &sarama.ProducerMessage{
				Topic: "notifications",
				Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s","time":"%s"}`, id, hostname, time.Now().Format(time.RFC3339))),
			}
			_, _, err := kafkaProducer.SendMessage(msg)
			if err != nil {
				log.Printf("Kafka 전송 실패: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"generated_by": hostname,
		})
	})

	fmt.Println("Notiflex API server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func startConsumer(broker, hostname string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup([]string{broker}, "notiflex-consumer", config)
	if err != nil {
		log.Printf("Kafka Consumer Group 생성 실패: %v", err)
		return
	}
	defer group.Close()

	handler := &consumerHandler{hostname: hostname}
	for {
		if err := group.Consume(context.Background(), []string{"notifications"}, handler); err != nil {
			log.Printf("Kafka Consumer 에러: %v", err)
			time.Sleep(5 * time.Second)
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
		log.Printf("[Kafka Consumer] topic=%s partition=%d offset=%d value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
