package pgx

import (
	"context"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
)

type EventsRepository struct {
	manager *txmanager.PGXTXManager
}

func NewEventsRepository(manager *txmanager.PGXTXManager) *EventsRepository {
	return &EventsRepository{
		manager: manager,
	}
}

func (r *EventsRepository) Create(ctx context.Context, event domain.Event) error {
	const query = `
		INSERT INTO events (id, event_type, payload, created_at, sent_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	engine := r.manager.GetQueryEngine(ctx)

	entity := NewEvent(event)

	_, err := engine.Exec(ctx, query,
		entity.ID,
		entity.EventType,
		entity.Payload,
		entity.CreatedAt,
		entity.SentAt,
	)
	if err != nil {
		return err
	}

	return nil
}
