package pgx

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"homework/internal/domain"
)

type Event struct {
	ID        uuid.UUID          `db:"id"`
	EventType string             `db:"event_type"`
	Payload   map[string]any     `db:"payload"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	SentAt    pgtype.Timestamptz `db:"sent_at"`
}

func NewEvent(model domain.Event) Event {
	return Event{
		ID:        model.ID,
		EventType: string(model.EventType),
		Payload:   model.Payload,
		CreatedAt: pgtype.Timestamptz{Time: model.CreatedAt, Valid: !model.CreatedAt.IsZero()},
		SentAt:    pgtype.Timestamptz{Time: model.SentAt, Valid: !model.SentAt.IsZero()},
	}
}

func (e Event) ToDomain() domain.Event {
	return domain.Event{
		ID:        e.ID,
		EventType: domain.EventType(e.EventType),
		Payload:   e.Payload,
		CreatedAt: e.CreatedAt.Time,
		SentAt:    e.SentAt.Time,
	}
}
