package kafkat

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"homework/internal/domain"
	"testing"
	"time"

	kafkaClient "homework/internal/infrastructure/clients/queue/kafka"
	"homework/internal/infrastructure/saramawrapper"
	"homework/internal/infrastructure/saramawrapper/producer"

	"github.com/IBM/sarama"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

func startContainer(ctx context.Context) (brokers []string, tearDown func(), err error) {
	kafkaContainer, err := kafka.Run(ctx, "confluentinc/confluent-local:7.5.0",
		testcontainers.WithEnv(map[string]string{
			"KAFKA_AUTO_CREATE_TOPICS_ENABLE": "true",
		}),
	)
	if err != nil {
		return
	}

	brokers, err = kafkaContainer.Brokers(ctx)
	if err != nil {
		return
	}

	tearDown = func() {
		_ = kafkaContainer.Terminate(ctx)
	}

	return
}

func convertMsg(in *sarama.ConsumerMessage) (domain.Event, error) {
	var event struct {
		ID        uuid.UUID              `json:"id"`
		EventType domain.EventType       `json:"event_type"`
		Payload   map[string]interface{} `json:"payload"`
		CreatedAt time.Time              `json:"created_at"`
		SentAt    time.Time              `json:"sent_at"`
	}
	err := json.Unmarshal(in.Value, &event)
	if err != nil {
		return domain.Event{}, err
	}
	return domain.Event{
		ID:        event.ID,
		EventType: event.EventType,
		Payload:   event.Payload,
		CreatedAt: event.CreatedAt,
		SentAt:    event.SentAt,
	}, nil
}

func TestPublisher(t *testing.T) {
	t.Parallel()

	const topic = "test-topic"

	ctx := context.Background()

	t.Logf("Starting Kafka container")
	brokers, tearDown, err := startContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to start Kafka container: %v", err)
	}
	defer tearDown()
	t.Logf("Kafka container started successfully")

	t.Logf("Creating Kafka consumer group")
	config := sarama.NewConfig()
	client, err := sarama.NewConsumerGroup(brokers, "test-group", config)
	if err != nil {
		t.Fatalf("Failed to create consumer group: %v", err)
	}

	consumer, ready, done, cancel := NewTestKafkaConsumer(t)
	go func() {
		if err := client.Consume(context.Background(), []string{topic}, consumer); err != nil {
			cancel()
		}
	}()

	<-ready

	config.Producer.Return.Successes = true

	t.Logf("Creating Kafka producer")

	prod, err := producer.NewSyncSaramaProducer(
		saramawrapper.Config{Brokers: brokers},
		producer.WithIdempotent(),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		producer.WithProducerPartitioner(sarama.NewHashPartitioner),
	)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer func(prod sarama.SyncProducer) {
		err := prod.Close()
		if err != nil {
			t.Fatalf("Failed to close producer: %v", err)
		}
	}(prod)

	meow := kafkaClient.NewProducer(prod, topic)

	t.Logf("Sending event")

	event := domain.NewEvent(
		domain.EventTypeOrderReturned,
		map[string]interface{}{
			"key": "value",
		},
	)

	err = meow.SendEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to send event: %v", err)
	}

	<-done

	rawMessage := consumer.message
	t.Logf("Received raw message: %+v", rawMessage)

	assert.NotNil(t, rawMessage)

	message, err := convertMsg(rawMessage)
	if err != nil {
		t.Fatalf("Failed to convert message: %v", err)
	}

	t.Logf("Received event: %+v", message)

	assert.Equal(t, event.CreatedAt.UTC(), message.CreatedAt.UTC())
	message.CreatedAt = event.CreatedAt
	assert.Equal(t, event, message)
}
