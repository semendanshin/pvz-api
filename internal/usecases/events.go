package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"homework/internal/domain"

	_ "github.com/gojuno/minimock/v3"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i EventsRepository -s _mock.go -o ./mocks
//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i KafkaClient -s _mock.go -o ./mocks

type EventsRepository interface {
	// GetPendingEvents returns a list of events that have not been sent yet.
	GetPendingEvents(ctx context.Context, limit int) ([]domain.Event, error)
	// MarkAsSent marks the event as sent.
	MarkAsSent(ctx context.Context, id uuid.UUID) error
}

type KafkaClient interface {
	// SendEvent sends an event to Kafka.
	SendEvent(ctx context.Context, event domain.Event) error
}

type EventsSender struct {
	repo   EventsRepository
	client KafkaClient

	done chan struct{}
}

func NewEventsSender(repo EventsRepository, client KafkaClient) *EventsSender {
	return &EventsSender{
		repo:   repo,
		client: client,
		done:   make(chan struct{}),
	}
}

func (e *EventsSender) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-e.done:
			return nil
		case <-ticker.C:
			if err := e.RunOnce(ctx); err != nil {
				return err
			}
		}
	}
}

func (e *EventsSender) Stop() {
	close(e.done)
}

func (e *EventsSender) RunOnce(ctx context.Context) error {
	events, err := e.repo.GetPendingEvents(ctx)
	if err != nil {
		return fmt.Errorf("error getting pending events: %w", err)
	}

	for _, event := range events {
		fmt.Println(event)
		if err := e.processEvent(ctx, event); err != nil {
			return fmt.Errorf("error processing event: %w", err)
		}
	}

	return nil
}

func (e *EventsSender) processEvent(ctx context.Context, event domain.Event) error {
	if err := e.client.SendEvent(ctx, event); err != nil {
		return fmt.Errorf("error sending event: %w", err)
	}

	if err := e.repo.MarkAsSent(ctx, event.ID); err != nil {
		return fmt.Errorf("error marking event as sent: %w", err)
	}

	return nil
}
