package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
	"homework/internal/usecases"
)

var _ usecases.PVZOrderRepository = &PostgresRepository{}

type PostgresRepository struct {
	manager *txmanager.PGXTXManager
}

func NewPostgresRepository(manager *txmanager.PGXTXManager) *PostgresRepository {
	return &PostgresRepository{
		manager: manager,
	}
}

func (p *PostgresRepository) CreateOrder(ctx context.Context, order domain.PVZOrder) error {
	const query = `
		INSERT INTO pvz_orders (order_id, pvz_id, recipient_id, cost, weight, packaging, additional_film, received_at, storage_time, issued_at, returned_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	engine := p.manager.GetQueryEngine(ctx)

	entity := newPgxPvzOrder(order)

	_, err := engine.Exec(ctx, query,
		entity.OrderID,
		entity.PVZID,
		entity.RecipientID,
		entity.Cost,
		entity.Weight,
		entity.Packaging,
		entity.AdditionalFilm,
		entity.ReceivedAt,
		entity.StorageTime,
		entity.IssuedAt,
		entity.ReturnedAt,
		entity.DeletedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) DeleteOrder(ctx context.Context, orderID string) error {
	const query = `
		UPDATE pvz_orders
		SET deleted_at = NOW()
		WHERE order_id = $1
	`

	engine := p.manager.GetQueryEngine(ctx)

	_, err := engine.Exec(ctx, query, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: order not found", domain.ErrNotFound)
		}
	}

	return nil
}

func (p *PostgresRepository) SetOrderIssued(ctx context.Context, orderID string) error {
	const query = `
		UPDATE pvz_orders
		SET issued_at = NOW()
		WHERE order_id = $1
	`

	engine := p.manager.GetQueryEngine(ctx)

	_, err := engine.Exec(ctx, query, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: order not found", domain.ErrNotFound)
		}
	}

	return nil
}

func (p *PostgresRepository) SetOrderReturned(ctx context.Context, orderID string) error {
	const query = `
		UPDATE pvz_orders
		SET returned_at = NOW()
		WHERE order_id = $1
	`

	engine := p.manager.GetQueryEngine(ctx)

	_, err := engine.Exec(ctx, query, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: order not found", domain.ErrNotFound)
		}
	}

	return nil
}

func (p *PostgresRepository) GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error) {
	opts, err := abstractions.NewGetOrdersOptions(options...)
	if err != nil {
		return nil, err
	}

	const query = `
		WITH subquery AS (
			SELECT order_id, pvz_id, recipient_id, cost, weight, packaging, additional_film, received_at, storage_time, issued_at, returned_at, deleted_at, 
				   ROW_NUMBER() OVER (ORDER BY received_at DESC) AS rn
			FROM pvz_orders
			WHERE recipient_id = $1 
			  AND (pvz_id = $2 OR $2 = '') 
			  AND deleted_at IS NULL
			ORDER BY received_at DESC
			LIMIT CASE WHEN $3 = 0 THEN NULL ELSE $3 END
		), firstRowN AS (
			SELECT rn 
			FROM subquery 
			WHERE order_id = $4 OR $4 = ''
			ORDER BY rn ASC
			LIMIT 1
		), row_boundary AS (
			SELECT COALESCE((SELECT rn FROM firstRowN), 1) AS start_row
		)
		SELECT order_id, pvz_id, recipient_id, cost, weight, packaging, additional_film, received_at, storage_time, issued_at, returned_at, deleted_at
		FROM subquery, row_boundary
		WHERE subquery.rn >= row_boundary.start_row
		ORDER BY subquery.rn ASC
		LIMIT CASE WHEN $5 = 0 THEN NULL ELSE $5 END
	`

	engine := p.manager.GetQueryEngine(ctx)

	var rows []*pgxPvzOrder

	err = pgxscan.Select(ctx, engine, &rows, query, userID, opts.PVZID, opts.LastNOrders, opts.CursorID, opts.Limit)
	if err != nil {
		return nil, err
	}

	orders := make([]domain.PVZOrder, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, row.ToDomain())
	}

	return orders, nil
}

func (p *PostgresRepository) GetOrder(ctx context.Context, orderID string) (domain.PVZOrder, error) {
	const query = `
		SELECT order_id, pvz_id, recipient_id, cost, weight, packaging, additional_film, received_at, storage_time, issued_at, returned_at, deleted_at
		FROM pvz_orders
		WHERE order_id = $1 AND deleted_at IS NULL
	`

	engine := p.manager.GetQueryEngine(ctx)

	var row pgxPvzOrder

	err := pgxscan.Get(ctx, engine, &row, query, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PVZOrder{}, fmt.Errorf("%w: order not found", domain.ErrNotFound)
		}
		return domain.PVZOrder{}, fmt.Errorf("failed to get order: %w", err)
	}

	return row.ToDomain(), nil
}

func (p *PostgresRepository) GetReturns(ctx context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error) {
	opts, err := abstractions.NewPaginationOptions(options...)
	if err != nil {
		return nil, err
	}

	const query = `
		SELECT order_id, pvz_id, recipient_id, cost, weight, packaging, additional_film, received_at, storage_time, issued_at, returned_at, deleted_at
		FROM pvz_orders
		WHERE returned_at IS NOT NULL AND deleted_at IS NULL
		ORDER BY returned_at DESC
		LIMIT $1 OFFSET $2
	`

	engine := p.manager.GetQueryEngine(ctx)

	var rows []*pgxPvzOrder

	err = pgxscan.Select(ctx, engine, &rows, query, opts.PageSize, opts.Page*opts.PageSize)
	if err != nil {
		return nil, err
	}

	orders := make([]domain.PVZOrder, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, row.ToDomain())
	}

	return orders, nil
}
