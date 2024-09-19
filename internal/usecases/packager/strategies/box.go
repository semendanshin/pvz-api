package strategies

import (
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecases"
)

var _ usecases.OrderPackagerStrategy = &BoxPackager{}

const (
	BoxPackagingCost        = 20_00
	BoxPackagingWeightLimit = 30_000
)

// BoxPackager is a packager for orders
type BoxPackager struct{}

// NewBoxPackager creates a new box packager
func NewBoxPackager() *BoxPackager {
	return &BoxPackager{}
}

// PackageOrder packages an order
func (b BoxPackager) PackageOrder(order domain.PVZOrder) (domain.PVZOrder, error) {
	if order.Weight > BoxPackagingWeightLimit {
		return domain.PVZOrder{}, fmt.Errorf("%w: weight limit exceeded", domain.ErrInvalidArgument)
	}

	order.Cost += BoxPackagingCost
	return order, nil
}
