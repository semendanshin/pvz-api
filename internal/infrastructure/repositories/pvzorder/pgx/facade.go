package pgx

import (
	"context"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
	"homework/internal/usecases"
)

var _ usecases.PVZOrderRepository = &PvzOrderFacade{}

type PvzOrderFacade struct {
	manager *txmanager.PGXTXManager
	repo    *PostgresRepository
}

func NewPgxPvzOrderFacade(manager *txmanager.PGXTXManager) *PvzOrderFacade {
	return &PvzOrderFacade{
		manager: manager,
		repo:    NewPostgresRepository(manager),
	}
}

func (p *PvzOrderFacade) CreateOrder(ctx context.Context, order domain.PVZOrder) error {
	return p.manager.RunSerializableTransaction(ctx, func(ctx context.Context) error {
		return p.repo.CreateOrder(ctx, order)
	})
}

func (p *PvzOrderFacade) DeleteOrder(ctx context.Context, orderID string) error {
	return p.manager.RunSerializableTransaction(ctx, func(ctx context.Context) error {
		return p.repo.DeleteOrder(ctx, orderID)
	})
}

func (p *PvzOrderFacade) SetOrderIssued(ctx context.Context, orderID string) error {
	return p.manager.RunSerializableTransaction(ctx, func(ctx context.Context) error {
		return p.repo.SetOrderIssued(ctx, orderID)
	})
}

func (p *PvzOrderFacade) SetOrderReturned(ctx context.Context, orderID string) error {
	return p.manager.RunSerializableTransaction(ctx, func(ctx context.Context) error {
		return p.repo.SetOrderReturned(ctx, orderID)
	})
}

func (p *PvzOrderFacade) GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error) {
	return p.repo.GetOrders(ctx, userID, options...)
}

func (p *PvzOrderFacade) GetOrder(ctx context.Context, orderID string) (domain.PVZOrder, error) {
	return p.repo.GetOrder(ctx, orderID)
}

func (p *PvzOrderFacade) GetReturns(ctx context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error) {
	return p.repo.GetReturns(ctx, options...)
}
