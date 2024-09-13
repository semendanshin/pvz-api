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

// PVZOrderUseCase is a use case for order operations
type PVZOrderUseCase struct {
	repo         abstractions.PVZOrderRepository
	currentPVZID string
}

// NewPVZOrderUseCase creates a new order use case
func NewPVZOrderUseCase(repo abstractions.PVZOrderRepository, currentPVZID string) *PVZOrderUseCase {
	return &PVZOrderUseCase{
		repo:         repo,
		currentPVZID: currentPVZID,
	}
}

// AcceptOrderDelivery accepts order delivery
func (P PVZOrderUseCase) AcceptOrderDelivery(orderID, recipientID string, storageTime time.Duration) error {
	_, err := P.repo.GetOrder(orderID)
	if err == nil {
		return fmt.Errorf("%w: order already exists", domain.ErrAlreadyExists)
	} else if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	order := domain.PVZOrder{
		OrderID:     orderID,
		PVZID:       P.currentPVZID,
		RecipientID: recipientID,
		ReceivedAt:  time.Now(),
		StorageTime: storageTime,
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
	orders := make([]domain.PVZOrder, 0, len(orderIDs))
	userID := ""

	for _, orderID := range orderIDs {
		order, err := P.repo.GetOrder(orderID)
		if err != nil {
			return err
		}

		if order.PVZID != P.currentPVZID {
			return fmt.Errorf("%w: order %s does not belong to this PVZ", domain.ErrInvalidArgument, order.OrderID)
		}

		if !order.IssuedAt.IsZero() {
			return fmt.Errorf("%w: order %s is already issued", domain.ErrInvalidArgument, order.OrderID)
		}

		if order.ReceivedAt.Add(order.StorageTime).Before(time.Now()) {
			return fmt.Errorf("%w: orders storage time for order %s has expired", domain.ErrInvalidArgument, order.OrderID)
		}

		if userID == "" {
			userID = order.RecipientID
		}

		if order.RecipientID != userID {
			return fmt.Errorf("%w: orders do not belong to the same user", domain.ErrInvalidArgument)
		}

		orders = append(orders, order)
	}

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

// AcceptReturn accepts return
func (P PVZOrderUseCase) AcceptReturn(userID, orderID string) error {
	order, err := P.repo.GetOrder(orderID)
	if err != nil {
		return err
	}

	if order.RecipientID != userID {
		return fmt.Errorf("%w: user is not recipient", domain.ErrInvalidArgument)
	}

	if order.PVZID != P.currentPVZID {
		return fmt.Errorf("%w: order does not belong to this PVZ", domain.ErrInvalidArgument)
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

	return P.repo.SetOrderReturned(orderID)
}

// GetReturns gets returns
func (P PVZOrderUseCase) GetReturns(options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error) {
	return P.repo.GetReturns(options...)
}
