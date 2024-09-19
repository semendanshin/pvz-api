package usecases

import (
	"errors"
	"fmt"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"slices"
	"time"
)

const (
	// TimeForReturn is a time for return
	TimeForReturn = 2 * 24 * time.Hour
)

var _ abstractions.IPVZOrderUseCase = &PVZOrderUseCase{}

// PVZOrderRepository is an interface for order repository
type PVZOrderRepository interface {
	CreateOrder(order domain.PVZOrder) error
	DeleteOrder(orderID string) error
	SetOrderIssued(orderID string) error
	SetOrderReturned(orderID string) error
	GetOrders(userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error)
	GetOrder(orderID string) (domain.PVZOrder, error)
	GetReturns(options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error)
}

type OrderPackagerInterface interface {
	PackageOrder(order domain.PVZOrder, packagingType domain.PackagingType) (domain.PVZOrder, error)
}

type OrderPackagerStrategy interface {
	PackageOrder(order domain.PVZOrder) (domain.PVZOrder, error)
}

// PVZOrderUseCase is a use case for order operations
type PVZOrderUseCase struct {
	repo         PVZOrderRepository
	packager     OrderPackagerInterface
	currentPVZID string
}

// NewPVZOrderUseCase creates a new order use case
func NewPVZOrderUseCase(repo PVZOrderRepository, packager OrderPackagerInterface, currentPVZID string) *PVZOrderUseCase {
	return &PVZOrderUseCase{
		repo:         repo,
		packager:     packager,
		currentPVZID: currentPVZID,
	}
}

func (P PVZOrderUseCase) checkOrderID(orderID string) error {
	_, err := P.repo.GetOrder(orderID)
	if err == nil {
		return fmt.Errorf("%w: order already exists", domain.ErrAlreadyExists)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	return nil
}

func (P PVZOrderUseCase) packageOrder(order domain.PVZOrder, packaging domain.PackagingType, additionalFilm bool) (domain.PVZOrder, error) {
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
func (P PVZOrderUseCase) AcceptOrderDelivery(orderID, recipientID string, storageTime time.Duration, cost, weight int, packaging domain.PackagingType, additionalFilm bool) error {
	if err := P.checkOrderID(orderID); err != nil {
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

	return P.repo.CreateOrder(order)
}

// ReturnOrderDelivery returns order delivery
func (P PVZOrderUseCase) ReturnOrderDelivery(orderID string) error {
	order, err := P.repo.GetOrder(orderID)
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

	return P.repo.DeleteOrder(orderID)
}

// GiveOrderToClient gives order to client
func (P PVZOrderUseCase) GiveOrderToClient(orderIDs []string) error {
	if len(orderIDs) == 0 {
		return fmt.Errorf("%w: orderIDs is empty", domain.ErrInvalidArgument)
	}

	orders := make([]domain.PVZOrder, len(orderIDs))

	var err error
	for i, orderID := range orderIDs {
		orders[i], err = P.repo.GetOrder(orderID)
		if err != nil {
			return err
		}
	}

	return P.setOrdersIssued(orders)
}

func (P PVZOrderUseCase) processOrders(orders []domain.PVZOrder) error {
	userID := orders[0].RecipientID

	if err := validateGiveOrdersToClient(orders, P.currentPVZID, userID); err != nil {
		return err
	}

	return P.setOrdersIssued(orders)
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

func (P PVZOrderUseCase) setOrdersIssued(orders []domain.PVZOrder) error {
	for _, order := range orders {
		err := P.repo.SetOrderIssued(order.OrderID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetOrders gets orders
func (P PVZOrderUseCase) GetOrders(userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error) {
	if slices.ContainsFunc(options, func(optFunc abstractions.GetOrdersOptFunc) bool {
		opts := &abstractions.GetOrdersOptions{}
		_ = optFunc(opts)
		return opts.SamePVZ
	}) {
		options = append(options, abstractions.WithPVZID(P.currentPVZID))
	}
	return P.repo.GetOrders(userID, options...)
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
func (P PVZOrderUseCase) AcceptReturn(userID, orderID string) error {
	order, err := P.repo.GetOrder(orderID)
	if err != nil {
		return err
	}

	if err := validateAcceptReturn(userID, order); err != nil {
		return err
	}

	return P.repo.SetOrderReturned(orderID)
}

// GetReturns gets returns
func (P PVZOrderUseCase) GetReturns(options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error) {
	return P.repo.GetReturns(options...)
}
