package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
	"homework/internal/usecases"
)

var _ usecases.EventsRepository = &EventsRepository{}

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

func (r *EventsRepository) GetPendingEvents(ctx context.Context, limit int) ([]domain.Event, error) {
	const query = `SELECT * FROM get_pending_events($1)`

	engine := r.manager.GetQueryEngine(ctx)

	var events []domain.Event

	if err := pgxscan.Select(ctx, engine, &events, query, limit); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no pending events found %w", domain.ErrNotFound)
		}
		return nil, err
	}

	return events, nil
}

func (r *EventsRepository) MarkAsSent(ctx context.Context, id uuid.UUID) error {
	const query = `SELECT mark_event_as_sent($1)`

	engine := r.manager.GetQueryEngine(ctx)

	_, err := engine.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: event not found", domain.ErrNotFound)
		}
	}

	return nil
}
