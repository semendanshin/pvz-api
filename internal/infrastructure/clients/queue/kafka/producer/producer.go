package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"homework/internal/domain"
	"homework/internal/infrastructure/clients/queue/kafka"
	"homework/internal/usecases"

	"github.com/IBM/sarama"
)

var _ usecases.QueueProducer = &Producer{}

// Producer is a Kafka producer.
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer.
func NewProducer(producer sarama.SyncProducer, topic string) *Producer {
	return &Producer{producer: producer, topic: topic}
}

// SendEvent sends an event to the queue.
func (p *Producer) SendEvent(_ context.Context, event domain.Event) error {
	eventEntity := kafka.NewEvent(event)

	data, err := json.Marshal(eventEntity)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
