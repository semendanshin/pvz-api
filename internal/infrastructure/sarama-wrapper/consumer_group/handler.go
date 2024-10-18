package consumer_group

import (
	"log"

	"github.com/IBM/sarama"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type ConsumerGroupHandler struct {
	ready   chan bool
	handler func(msg *Msg) error
}

func NewConsumerGroupHandler(handler func(msg *Msg) error) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		ready:   make(chan bool),
		handler: handler,
	}
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор, пока сессия не завершилась
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := convertMsg(message)
			log.Printf("Message claimed: %+v\n", msg)

			err := h.handler(&msg)
			if err != nil {
				log.Printf("Error handling message: %v\n", err)
				return err
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

type Msg struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Payload   string `json:"payload"`
}

func convertMsg(in *sarama.ConsumerMessage) Msg {
	return Msg{
		Topic:     in.Topic,
		Partition: in.Partition,
		Offset:    in.Offset,
		Key:       string(in.Key),
		Payload:   string(in.Value),
	}
}
