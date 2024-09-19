package packager

import (
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecases"
)

var _ usecases.OrderPackagerInterface = &OrderPackager{}

// OrderPackager is a packager for orders
type OrderPackager struct {
	strategies map[domain.PackagingType]usecases.OrderPackagerStrategy
}

// NewOrderPackager creates a new order packager
func NewOrderPackager(strategies map[domain.PackagingType]usecases.OrderPackagerStrategy) *OrderPackager {
	return &OrderPackager{
		strategies: strategies,
	}
}

// PackageOrder packages an order
func (o OrderPackager) PackageOrder(
	order domain.PVZOrder, packaging domain.PackagingType,
) (domain.PVZOrder, error) {
	strategy, ok := o.strategies[packaging]
	if !ok {
		return domain.PVZOrder{}, fmt.Errorf("unknown packaging type: %s", domain.ErrInvalidArgument)
	}

	order, err := strategy.PackageOrder(order)
	if err != nil {
		return domain.PVZOrder{}, err
	}

	return order, nil
}
