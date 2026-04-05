package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/valkey-io/valkey-go"
)

var version = "v0.6.0"
var valkeyClient valkey.Client
var kafkaProducer sarama.SyncProducer

func main() {
	// Valkey 연결
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "valkey-primary.notiflex.svc.cluster.local:6379"
	}

	var password string
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
			log.Printf("Valkey 비밀번호를 파일에서 읽음: %s", pwFile)
		}
	}
	if password == "" {
		password = os.Getenv("VALKEY_PASSWORD")
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

	// Kafka Producer 연결
	broker := os.Getenv("KAFKA_BROKER")
	if broker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Retry.Max = 5
		for i := 0; i < 10; i++ {
			kafkaProducer, err = sarama.NewSyncProducer([]string{broker}, config)
			if err == nil {
				log.Printf("Kafka 연결 성공: %s", broker)
				break
			}
			log.Printf("Kafka 연결 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Printf("Kafka 연결 실패 (Producer 비활성): %v", err)
		} else {
			defer kafkaProducer.Close()
			go consumeNotifications(broker)
		}
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	http.HandleFunc("/notify", notifyHandler)

	log.Printf("Notiflex API %s starting on :8080", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version,
	})
}

func nextID(ctx context.Context) (int64, error) {
	cmd := valkeyClient.B().Incr().Key("notiflex:id").Build()
	return valkeyClient.Do(ctx, cmd).AsInt64()
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := nextID(r.Context())
	if err != nil {
		http.Error(w, "ID generation failed", http.StatusInternalServerError)
		return
	}
	pod := os.Getenv("HOSTNAME")

	if kafkaProducer != nil {
		msg := &sarama.ProducerMessage{
			Topic: "notifications",
			Key:   sarama.StringEncoder(strconv.FormatInt(id, 10)),
			Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s","type":"notify"}`, id, pod)),
		}
		_, _, err := kafkaProducer.SendMessage(msg)
		if err != nil {
			log.Printf("Kafka send failed: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     pod,
		"message": "notification sent",
		"version": version,
	})
	fmt.Printf("Notification %d sent from pod %s\n", id, pod)
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	id, err := nextID(r.Context())
	if err != nil {
		http.Error(w, "ID generation failed", http.StatusInternalServerError)
		return
	}
	pod := os.Getenv("HOSTNAME")

	if kafkaProducer != nil {
		msg := &sarama.ProducerMessage{
			Topic: "notifications",
			Key:   sarama.StringEncoder(strconv.FormatInt(id, 10)),
			Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s","type":"id"}`, id, pod)),
		}
		_, _, err := kafkaProducer.SendMessage(msg)
		if err != nil {
			log.Printf("Kafka send failed: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     pod,
		"version": version,
	})
	fmt.Printf("Generated ID: %d on pod %s\n", id, pod)
}

func consumeNotifications(broker string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	for i := 0; i < 10; i++ {
		group, err := sarama.NewConsumerGroup([]string{broker}, "notiflex-consumer", config)
		if err != nil {
			log.Printf("Kafka Consumer 재시도 %d/10: %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}
		defer group.Close()
		log.Println("Kafka Consumer 시작: notifications 토픽")
		handler := &consumerHandler{}
		for {
			if err := group.Consume(context.Background(), []string{"notifications"}, handler); err != nil {
				log.Printf("Kafka Consumer 에러: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

type consumerHandler struct{}

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
