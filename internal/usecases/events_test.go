package usecases

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"homework/internal/domain"
	"homework/internal/usecases/mocks"
	"testing"
	"time"
)

func TestEventsProcessor_RunOnce(t *testing.T) {
	t.Parallel()

	const limit = 10

	ctrl := minimock.NewController(t)

	repo := mocks.NewEventsRepositoryMock(ctrl)
	queueProducer := mocks.NewQueueProducerMock(ctrl)

	processor := NewEventsProcessor(repo, queueProducer, limit)

	id1 := uuid.New()
	id2 := uuid.New()

	events := []domain.Event{
		{
			ID:        id1,
			EventType: domain.EventTypeOrderIssued,
			Payload: map[string]interface{}{
				"order_id": 1,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:        id2,
			EventType: domain.EventTypeOrderReturned,
			Payload: map[string]interface{}{
				"order_id": 2,
			},
			CreatedAt: time.Now(),
		},
	}

	repo.GetPendingEventsMock.Expect(minimock.AnyContext, limit).Return(events, nil)

	queueProducer.SendEventMock.Times(2)
	repo.MarkAsSentMock.Times(2)

	queueProducer.SendEventMock.Inspect(func(ctx context.Context, event domain.Event) {
		assert.Contains(t, events, event)
	})
	repo.MarkAsSentMock.Inspect(func(ctx context.Context, id uuid.UUID) {
		assert.Contains(t, []uuid.UUID{id1, id2}, id)
	})

	queueProducer.SendEventMock.Return(nil)
	repo.MarkAsSentMock.Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := processor.RunOnce(ctx)

	assert.NoError(t, err)
}

func TestEventsProcessor_RunOnce_Error(t *testing.T) {
	t.Parallel()

	const limit = 10

	ctrl := minimock.NewController(t)

	repo := mocks.NewEventsRepositoryMock(ctrl)
	queueProducer := mocks.NewQueueProducerMock(ctrl)

	processor := NewEventsProcessor(repo, queueProducer, limit)

	id1 := uuid.New()
	id2 := uuid.New()

	events := []domain.Event{
		{
			ID:        id1,
			EventType: domain.EventTypeOrderIssued,
			Payload: map[string]interface{}{
				"order_id": 1,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:        id2,
			EventType: domain.EventTypeOrderReturned,
			Payload: map[string]interface{}{
				"order_id": 2,
			},
			CreatedAt: time.Now(),
		},
	}

	repo.GetPendingEventsMock.Expect(minimock.AnyContext, limit).Return(events, nil)

	queueProducer.SendEventMock.Times(1)
	repo.MarkAsSentMock.Times(1)

	queueProducer.SendEventMock.Return(nil)
	repo.MarkAsSentMock.Return(domain.ErrNotFound)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := processor.RunOnce(ctx)

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestEventsProcessor_RunOnce_Empty(t *testing.T) {
	t.Parallel()

	const limit = 10

	ctrl := minimock.NewController(t)

	repo := mocks.NewEventsRepositoryMock(ctrl)
	queueProducer := mocks.NewQueueProducerMock(ctrl)

	processor := NewEventsProcessor(repo, queueProducer, limit)

	repo.GetPendingEventsMock.Expect(minimock.AnyContext, limit).Return(nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := processor.RunOnce(ctx)

	assert.NoError(t, err)
}
