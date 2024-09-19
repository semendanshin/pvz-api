package strategies

import (
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecases"
)

var _ usecases.OrderPackagerStrategy = &BagPackager{}

const (
	BagPackagingCost        = 5_00
	BagPackagingWeightLimit = 10_000
)

// BagPackager is a packager for orders
type BagPackager struct{}

// NewBagPackager creates a new bag packager
func NewBagPackager() *BagPackager {
	return &BagPackager{}
}

// PackageOrder packages an order
func (b BagPackager) PackageOrder(order domain.PVZOrder) (domain.PVZOrder, error) {
	if order.Weight > BagPackagingWeightLimit {
		return domain.PVZOrder{}, fmt.Errorf("%w: weight limit exceeded", domain.ErrInvalidArgument)
	}

	order.Cost += BagPackagingCost
	return order, nil
}
