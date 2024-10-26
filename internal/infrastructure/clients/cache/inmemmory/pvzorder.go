package inmemmory

import (
	"context"
	"fmt"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/usecases"
	"time"
)

var _ usecases.PVZOrderCache = &PVZOrder{}

// PVZOrder is a cache for PVZ orders
type PVZOrder struct {
	cache *Cache[string, interface{}]
}

// NewPVZOrder creates a new PVZ order cache
func NewPVZOrder(ttl time.Duration, maxItems int, invalidationStrategy InvalidationStrategy[string, interface{}]) *PVZOrder {
	return &PVZOrder{
		cache: NewCache[string, interface{}](
			ttl,
			maxItems,
			invalidationStrategy,
		),
	}
}

func (P PVZOrder) GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error, bool) {
	key := fmt.Sprintf("GetOrders:%s:%v", userID, options)

	v, ok := P.cache.Get(key)

	if !ok {
		return nil, nil, false
	}

	result, ok := v.([]domain.PVZOrder)
	if !ok {
		return nil, fmt.Errorf("invalid cache value type"), false
	}

	return result, nil, true
}

func (P PVZOrder) GetReturns(ctx context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error, bool) {
	key := fmt.Sprintf("GetReturns:%v", options)

	v, ok := P.cache.Get(key)

	if !ok {
		return nil, nil, false
	}

	result, ok := v.([]domain.PVZOrder)
	if !ok {
		return nil, fmt.Errorf("invalid cache value type"), false
	}

	return result, nil, true
}

func (P PVZOrder) SetGetOrders(ctx context.Context, userID string, orders []domain.PVZOrder, options ...abstractions.GetOrdersOptFunc) error {
	key := fmt.Sprintf("GetOrders:%s:%v", userID, options)

	P.cache.Set(key, orders)

	return nil
}

func (P PVZOrder) SetGetReturns(ctx context.Context, orders []domain.PVZOrder, options ...abstractions.PagePaginationOptFunc) error {
	key := fmt.Sprintf("GetReturns:%v", options)

	P.cache.Set(key, orders)

	return nil
}

func (P PVZOrder) GetOrder(ctx context.Context, orderID string) (domain.PVZOrder, error, bool) {
	key := fmt.Sprintf("GetOrder:%s", orderID)

	v, ok := P.cache.Get(key)

	if !ok {
		return domain.PVZOrder{}, nil, false
	}

	result, ok := v.(domain.PVZOrder)
	if !ok {
		return domain.PVZOrder{}, fmt.Errorf("invalid cache value type"), false
	}

	return result, nil, true
}

func (P PVZOrder) SetOrder(ctx context.Context, order domain.PVZOrder) (domain.PVZOrder, error) {
	key := fmt.Sprintf("GetOrder:%s", order.OrderID)

	P.cache.Set(key, order)

	return order, nil
}
