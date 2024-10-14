package pgx

import (
	"context"
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
	return p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewEvent(domain.EventTypeOrderDeliveryAccepted, map[string]interface{}{
			"order_id":        order.OrderID,
			"pvz_id":          order.PVZID,
			"recipient_id":    order.RecipientID,
			"cost":            order.Cost,
			"weight":          order.Weight,
			"packaging":       order.Packaging,
			"additional_film": order.AdditionalFilm,
			"received_at":     order.ReceivedAt.Format("2006-01-02 15:04:05"),
			"storage_time":    order.StorageTime.Milliseconds(),
		})
		if err := p.repo.CreateOrder(ctx, order); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) DeleteOrder(ctx context.Context, orderID string) error {
	return p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewEvent(domain.EventTypeOrderDeliveryReturned, map[string]interface{}{
			"order_id": orderID,
		})
		if err := p.repo.DeleteOrder(ctx, orderID); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) SetOrderIssued(ctx context.Context, orderID string) error {
	return p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewEvent(domain.EventTypeOrderIssued, map[string]interface{}{
			"order_id": orderID,
		})
		if err := p.repo.SetOrderIssued(ctx, orderID); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) SetOrderReturned(ctx context.Context, orderID string) error {
	return p.manager.RunSerializableTransaction(ctx, func(ctx context.Context) error {
		event := domain.NewEvent(domain.EventTypeOrderReturned, map[string]interface{}{
			"order_id": orderID,
		})
		if err := p.repo.SetOrderReturned(ctx, orderID); err != nil {
			return err
		}
		return p.eventsRepo.Create(ctx, event)
	})
}

func (p *PvzOrderFacade) GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error) {
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
	var result []domain.PVZOrder
	var err error
	err = p.manager.RunReadCommittedTransaction(ctx, func(ctx context.Context) error {
		var innerErr error
		result, innerErr = p.repo.GetReturns(ctx, options...)
		return innerErr
	})

	return result, err
}
