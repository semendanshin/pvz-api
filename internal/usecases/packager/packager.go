package packager

import (
	"fmt"
	"homework/internal/abstractions"
	"homework/internal/domain"
)

var _ abstractions.OrderPackagerInterface = &OrderPackager{}

// OrderPackager is a packager for orders
type OrderPackager struct {
	strategies map[domain.PackagingType]abstractions.OrderPackagerStrategy
}

// NewOrderPackager creates a new order packager
func NewOrderPackager(strategies map[domain.PackagingType]abstractions.OrderPackagerStrategy) *OrderPackager {
	return &OrderPackager{
		strategies: strategies,
	}
}

// PackageOrder packages an order
func (o OrderPackager) PackageOrder(
	order domain.PVZOrder, packaging domain.PackagingType, additionalFilm bool,
) (domain.PVZOrder, error) {
	strategy, ok := o.strategies[packaging]
	if !ok {
		return domain.PVZOrder{}, fmt.Errorf("unknown packaging type: %s", domain.ErrInvalidArgument)
	}

	if packaging == domain.PackagingTypeFilm && additionalFilm {
		return domain.PVZOrder{}, fmt.Errorf("%w: additional film is not allowed for film packaging", domain.ErrInvalidArgument)
	}

	order, err := strategy.PackageOrder(order)
	if err != nil {
		return domain.PVZOrder{}, err
	}

	return order, nil
}
