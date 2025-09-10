package kafka

import (
	"encoding/json"
	"homework/internal/infrastructure/saramawrapper/consumer_group"
	"log"
)

func NewEventsHandler() func(msg *consumer_group.Msg) error {
	return handleEvent
}

func handleEvent(msg *consumer_group.Msg) error {
	var eventEntity Event
	if err := json.Unmarshal([]byte(msg.Payload), &eventEntity); err != nil {
		log.Printf("[QueueEventsHandler] Failed to unmarshal message: %v\n", err)
		return nil
	}

	log.Printf("[handler] Received event: %v\n", eventEntity)

	return nil
}
