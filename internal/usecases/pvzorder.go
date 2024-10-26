package usecases

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/metrics"
	"slices"
	"time"

	"homework/internal/abstractions"
	"homework/internal/domain"

	_ "github.com/gojuno/minimock/v3"
)

const (
	// TimeForReturn is a time for return
	TimeForReturn = 2 * 24 * time.Hour
)

var _ abstractions.IPVZOrderUseCase = &PVZOrderUseCase{}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i PVZOrderRepository -s _mock.go -o ./mocks
//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i OrderPackagerInterface -s _mock.go -o ./mocks

// PVZOrderRepository is an interface for order repository
type PVZOrderRepository interface {
	CreateOrder(ctx context.Context, order domain.PVZOrder) error
	DeleteOrder(ctx context.Context, orderID string) error
	SetOrderIssued(ctx context.Context, orderID string) error
	SetOrderReturned(ctx context.Context, orderID string) error
	GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error)
	GetOrder(ctx context.Context, orderID string) (domain.PVZOrder, error)
	GetReturns(ctx context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error)
}

type OrderPackagerInterface interface {
	PackageOrder(order domain.PVZOrder, packagingType domain.PackagingType) (domain.PVZOrder, error)
}

type PVZOrderCache interface {
	GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error, bool)
	GetReturns(ctx context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error, bool)
	SetGetOrders(ctx context.Context, userID string, orders []domain.PVZOrder, options ...abstractions.GetOrdersOptFunc) error
	SetGetReturns(ctx context.Context, orders []domain.PVZOrder, options ...abstractions.PagePaginationOptFunc) error
	GetOrder(ctx context.Context, orderID string) (domain.PVZOrder, error, bool)
	SetOrder(ctx context.Context, order domain.PVZOrder) (domain.PVZOrder, error)
}

// PVZOrderUseCase is a use case for order operations
type PVZOrderUseCase struct {
	repo         PVZOrderRepository
	packager     OrderPackagerInterface
	currentPVZID string
	cache        PVZOrderCache
}

// NewPVZOrderUseCase creates a new order use case
func NewPVZOrderUseCase(repo PVZOrderRepository, packager OrderPackagerInterface, currentPVZID string) *PVZOrderUseCase {
	return &PVZOrderUseCase{
		repo:         repo,
		packager:     packager,
		currentPVZID: currentPVZID,
	}
}

func (P *PVZOrderUseCase) checkOrderID(ctx context.Context, orderID string) error {
	_, err := P.repo.GetOrder(ctx, orderID)
	if err == nil {
		return fmt.Errorf("%w: order already exists", domain.ErrAlreadyExists)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	return nil
}

func (P *PVZOrderUseCase) packageOrder(order domain.PVZOrder, packaging domain.PackagingType, additionalFilm bool) (domain.PVZOrder, error) {
	order, err := P.packager.PackageOrder(order, packaging)
	if err != nil {
		return domain.PVZOrder{}, err
	}

	if additionalFilm {
		order, err = P.packager.PackageOrder(order, domain.PackagingTypeFilm)
		if err != nil {
			return domain.PVZOrder{}, err
		}
	}

	return order, nil
}

// AcceptOrderDelivery accepts order delivery
func (P *PVZOrderUseCase) AcceptOrderDelivery(ctx context.Context, orderID, recipientID string, storageTime time.Duration, cost, weight int, packaging domain.PackagingType, additionalFilm bool) error {
	if err := P.checkOrderID(ctx, orderID); err != nil {
		return err
	}

	if packaging == domain.PackagingTypeFilm && additionalFilm {
		return fmt.Errorf("%w: additional film is not allowed for film packaging", domain.ErrInvalidArgument)
	}

	order := domain.NewPVZOrder(
		orderID,
		P.currentPVZID,
		recipientID,
		cost,
		weight,
		storageTime,
		packaging,
		additionalFilm,
	)

	order, err := P.packageOrder(order, packaging, additionalFilm)
	if err != nil {
		return err
	}

	return P.repo.CreateOrder(ctx, order)
}

// ReturnOrderDelivery returns order delivery
func (P *PVZOrderUseCase) ReturnOrderDelivery(ctx context.Context, orderID string) error {
	order, err := P.repo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.PVZID != P.currentPVZID {
		return fmt.Errorf("%w: order does not belong to this PVZ", domain.ErrInvalidArgument)
	}

	if order.ReceivedAt.Add(order.StorageTime).After(time.Now()) {
		return fmt.Errorf("%w: storage time has not expired", domain.ErrInvalidArgument)
	}

	if !order.IssuedAt.IsZero() {
		return fmt.Errorf("%w: order is already issued", domain.ErrInvalidArgument)
	}

	return P.repo.DeleteOrder(ctx, orderID)
}

// GiveOrderToClient gives order to client
func (P *PVZOrderUseCase) GiveOrderToClient(ctx context.Context, orderIDs []string) error {
	if len(orderIDs) == 0 {
		return fmt.Errorf("%w: orderIDs is empty", domain.ErrInvalidArgument)
	}

	orders := make([]domain.PVZOrder, len(orderIDs))

	var err error
	for i, orderID := range orderIDs {
		orders[i], err = P.repo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}
	}

	return P.processOrders(ctx, orders)
}

func (P *PVZOrderUseCase) processOrders(ctx context.Context, orders []domain.PVZOrder) error {
	userID := orders[0].RecipientID

	if err := validateGiveOrdersToClient(orders, P.currentPVZID, userID); err != nil {
		return err
	}

	return P.setOrdersIssued(ctx, orders)
}

func validateGiveOrdersToClient(orders []domain.PVZOrder, currentPVZID string, userID string) error {
	for _, order := range orders {
		if err := validateGiveOrderToClient(order, currentPVZID); err != nil {
			return err
		}

		if order.RecipientID != userID {
			return fmt.Errorf("%w: orders do not belong to the same user", domain.ErrInvalidArgument)
		}
	}

	return nil
}

func validateGiveOrderToClient(order domain.PVZOrder, currentPVZID string) error {
	if order.PVZID != currentPVZID {
		return fmt.Errorf("%w: order does not belong to this PVZ", domain.ErrInvalidArgument)
	}

	if !order.IssuedAt.IsZero() {
		return fmt.Errorf("%w: order is already issued", domain.ErrInvalidArgument)
	}

	if order.ReceivedAt.Add(order.StorageTime).Before(time.Now()) {
		return fmt.Errorf("%w: orders storage time has expired", domain.ErrInvalidArgument)
	}

	return nil
}

func (P *PVZOrderUseCase) setOrdersIssued(ctx context.Context, orders []domain.PVZOrder) error {
	for _, order := range orders {
		err := P.repo.SetOrderIssued(ctx, order.OrderID)
		if err != nil {
			return err
		}
		metrics.IncOrdersIssued(P.currentPVZID)
	}
	return nil
}

// GetOrders gets orders
func (P *PVZOrderUseCase) GetOrders(ctx context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error) {
	if slices.ContainsFunc(options, func(optFunc abstractions.GetOrdersOptFunc) bool {
		opts := &abstractions.GetOrdersOptions{}
		_ = optFunc(opts)
		return opts.SamePVZ
	}) {
		options = append(options, abstractions.WithPVZID(P.currentPVZID))
	}

	orders, err, ok := P.cache.GetOrders(ctx, userID, options...)
	if err != nil {
		return nil, err
	}

	if ok {
		return orders, nil
	}

	orders, err = P.repo.GetOrders(ctx, userID, options...)
	if err != nil {
		return nil, err
	}

	err = P.cache.SetGetOrders(ctx, userID, orders, options...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func validateAcceptReturn(userID string, order domain.PVZOrder) error {
	if order.RecipientID != userID {
		return fmt.Errorf("%w: user is not recipient", domain.ErrInvalidArgument)
	}

	if !order.ReturnedAt.IsZero() {
		return fmt.Errorf("%w: order is already returned", domain.ErrInvalidArgument)
	}

	if order.IssuedAt.IsZero() {
		return fmt.Errorf("%w: order is not issued", domain.ErrInvalidArgument)
	}

	if order.IssuedAt.Add(TimeForReturn).Before(time.Now()) {
		return fmt.Errorf("%w: time for return has expired", domain.ErrInvalidArgument)
	}

	return nil
}

// AcceptReturn accepts return
func (P *PVZOrderUseCase) AcceptReturn(ctx context.Context, userID, orderID string) error {
	order, err := P.repo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if err := validateAcceptReturn(userID, order); err != nil {
		return err
	}

	return P.repo.SetOrderReturned(ctx, orderID)
}

// GetReturns gets returns
func (P *PVZOrderUseCase) GetReturns(ctx context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error) {
	orders, err, ok := P.cache.GetReturns(ctx, options...)
	if err != nil {
		return nil, err
	}

	if ok {
		return orders, nil
	}

	orders, err = P.repo.GetReturns(ctx, options...)
	if err != nil {
		return nil, err
	}

	err = P.cache.SetGetReturns(ctx, orders, options...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
