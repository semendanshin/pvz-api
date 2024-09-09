package strategies

import (
	"homework/internal/abstractions"
	"homework/internal/domain"
)

var _ abstractions.OrderPackagerStrategy = &FilmPackager{}

const (
	FilmPackagingCost = 1_00
)

// FilmPackager is a packager for orders
type FilmPackager struct{}

// NewFilmPackager creates a new film packager
func NewFilmPackager() *FilmPackager {
	return &FilmPackager{}
}

// PackageOrder packages an order
func (f FilmPackager) PackageOrder(order domain.PVZOrder) (domain.PVZOrder, error) {
	order.Cost += FilmPackagingCost
	return order, nil
}
