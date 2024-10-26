package pgx

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/events/pgx"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
	"homework/internal/usecases"
)

var _ usecases.PVZOrderRepository = &PvzOrderFacade{}

type PvzOrderFacade struct {
	manager    *txmanager.PGXTXManager
	repo       *PostgresRepository
	eventsRepo *pgx.EventsRepository
}

func NewPgxPvzOrderFacade(manager *txmanager.PGXTXManager) *PvzOrderFacade {
	return &PvzOrderFacade{
		manager:    manager,
		repo:       NewPostgresRepository(manager),
		eventsRepo: pgx.NewEventsRepository(manager),
	}
}

func (p *PvzOrderFacade) CreateOrder(ctx context.Context, order domain.PVZOrder) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PvzOrderFacade.CreateOrder")
	defer span.Finish()

	return p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewOrderDeliveryAcceptedEvent(
			order.OrderID,
			order.PVZID,
			order.RecipientID,
			order.Cost,
			order.Weight,
			order.Packaging,
			order.AdditionalFilm,
			order.ReceivedAt,
			order.StorageTime,
		)
		if err := p.repo.CreateOrder(ctx, order); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) DeleteOrder(ctx context.Context, orderID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PvzOrderFacade.DeleteOrder")
	defer span.Finish()

	return p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewOrderDeliveryReturnedEvent(orderID)
		if err := p.repo.DeleteOrder(ctx, orderID); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) SetOrderIssued(ctx context.Context, orderID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PvzOrderFacade.SetOrderIssued")
	defer span.Finish()

	return p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewOrderIssuedEvent(orderID)
		if err := p.repo.SetOrderIssued(ctx, orderID); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) SetOrderReturned(ctx context.Context, orderID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PvzOrderFacade.SetOrderReturned")
	defer span.Finish()

	return p.manager.RunSerializableTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewOrderReturnedEvent(orderID)
		if err := p.repo.SetOrderReturned(ctx, orderID); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PvzOrderFacade.GetOrders")
	defer span.Finish()

	var result []domain.PVZOrder
	var err error
	err = p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		var innerErr error
		result, innerErr = p.repo.GetOrders(ctx, userID, options...)
		return innerErr
	})

	return result, err
}

func (p *PvzOrderFacade) GetOrder(ctx context.Context, orderID string) (domain.PVZOrder, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PvzOrderFacade.GetOrder")
	defer span.Finish()

	var result domain.PVZOrder
	var err error
	err = p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		var innerErr error
		result, innerErr = p.repo.GetOrder(ctx, orderID)
		return innerErr
	})

	return result, err
}

func (p *PvzOrderFacade) GetReturns(ctx context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PvzOrderFacade.GetReturns")
	defer span.Finish()

	var result []domain.PVZOrder
	var err error
	err = p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		var innerErr error
		result, innerErr = p.repo.GetReturns(ctx, options...)
		return innerErr
	})

	return result, err
}
