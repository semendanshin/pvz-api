package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type EventType string

func NewEventType(eventType string) (EventType, error) {
	switch eventType {
	case EventTypeOrderDeliveryAccepted.String():
		return EventTypeOrderDeliveryAccepted, nil
	case EventTypeOrderIssued.String():
		return EventTypeOrderIssued, nil
	case EventTypeOrderDeliveryReturned.String():
		return EventTypeOrderDeliveryReturned, nil
	case EventTypeOrderReturned.String():
		return EventTypeOrderReturned, nil
	default:
		return EventTypeUnknown, fmt.Errorf("unknown event type %s: %w", eventType, ErrInvalidArgument)
	}
}

func (e EventType) String() string {
	return string(e)
}

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

func NewOrderDeliveryAcceptedEvent(orderID, pvzID, recipientID string, cost, weight int, packaging PackagingType, additionalFilm bool, receivedAt time.Time, storageTime time.Duration) Event {
	return NewEvent(EventTypeOrderDeliveryAccepted, map[string]interface{}{
		"order_id":        orderID,
		"pvz_id":          pvzID,
		"recipient_id":    recipientID,
		"cost":            cost,
		"weight":          weight,
		"packaging":       packaging,
		"additional_film": additionalFilm,
		"received_at":     receivedAt,
		"storage_time":    storageTime,
	})
}

func NewOrderIssuedEvent(orderID string) Event {
	return NewEvent(EventTypeOrderIssued, map[string]interface{}{
		"order_id": orderID,
	})
}

func NewOrderDeliveryReturnedEvent(orderID string) Event {
	return NewEvent(EventTypeOrderDeliveryReturned, map[string]interface{}{
		"order_id": orderID,
	})
}

func NewOrderReturnedEvent(orderID string) Event {
	return NewEvent(EventTypeOrderReturned, map[string]interface{}{
		"order_id": orderID,
	})
}
