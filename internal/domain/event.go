package domain

import (
	"github.com/google/uuid"
	"time"
)

type EventType string

const (
	EventTypeUnknown               EventType = "unknown"
	EventTypeOrderDeliveryAccepted EventType = "order_delivery_accepted"
	EventTypeOrderIssued           EventType = "order_issued"
	EventTypeOrderDeliveryReturned EventType = "order_delivery_returned"
	EventTypeOrderReturned         EventType = "order_returned"
)

type Event struct {
	ID        uuid.UUID
	EventType EventType
	Payload   map[string]interface{}
	CreatedAt time.Time
	SentAt    time.Time
}

func NewEvent(eventType EventType, payload map[string]interface{}) Event {
	return Event{
		ID:        uuid.New(),
		EventType: eventType,
		Payload:   payload,
		CreatedAt: time.Now(),
		SentAt:    time.Time{},
	}
}
