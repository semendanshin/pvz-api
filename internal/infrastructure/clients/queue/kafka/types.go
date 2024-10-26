package kafka

import (
	"github.com/google/uuid"
	"homework/internal/domain"
	"time"
)

type Event struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}

func NewEvent(event domain.Event) Event {
	return Event{
		ID:        event.ID.String(),
		EventType: string(event.EventType),
		Payload:   event.Payload,
		CreatedAt: event.CreatedAt,
	}
}

func (e Event) ToDomain() (domain.Event, error) {
	id, err := uuid.Parse(e.ID)
	if err != nil {
		return domain.Event{}, err
	}

	eventType, err := domain.NewEventType(e.EventType)
	if err != nil {
		return domain.Event{}, err
	}

	return domain.Event{
		ID:        id,
		EventType: eventType,
		Payload:   e.Payload,
		CreatedAt: e.CreatedAt,
	}, nil
}
