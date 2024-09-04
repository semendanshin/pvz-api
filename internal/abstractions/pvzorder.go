package abstractions

import (
	"homework/internal/domain"

	"time"
)

// PaginationOptions is a struct for pagination options
type PaginationOptions struct {
	Page     int
	PageSize int
}

// PaginationOptFunc is a type for pagination options
type PaginationOptFunc func(*PaginationOptions) error

// WithPage is an option to get orders by page
func WithPage(page int) PaginationOptFunc {
	return func(o *PaginationOptions) error {
		o.Page = page
		return nil
	}
}

// WithPageSize is an option to get orders by page size
func WithPageSize(pageSize int) PaginationOptFunc {
	return func(o *PaginationOptions) error {
		o.PageSize = pageSize
		return nil
	}
}

// NewPaginationOptions creates new pagination options
func NewPaginationOptions(options ...PaginationOptFunc) (*PaginationOptions, error) {
	opts := &PaginationOptions{
		Page:     0,
		PageSize: 10,
	}
	for _, opt := range options {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}
	return opts, nil
}

// GetOrdersOptions is a struct for get orders options
type GetOrdersOptions struct {
	LastNOrders int
	PVZID       string
	*PaginationOptions
}

// GetOrdersOptFunc is a type for order options
type GetOrdersOptFunc func(*GetOrdersOptions) error

// WithLastNOrders is an option to get last N orders
func WithLastNOrders(lastNOrders int) GetOrdersOptFunc {
	return func(o *GetOrdersOptions) error {
		o.LastNOrders = lastNOrders
		return nil
	}
}

// WithPVZID is an option to get orders by PVZ ID
func WithPVZID(pvzID string) GetOrdersOptFunc {
	return func(o *GetOrdersOptions) error {
		o.PVZID = pvzID
		return nil
	}
}

// WithPaginationOptions is an option to get orders with pagination options
func WithPaginationOptions(opts *PaginationOptions) GetOrdersOptFunc {
	return func(o *GetOrdersOptions) error {
		o.PaginationOptions = opts
		return nil
	}
}

// NewGetOrdersOptions creates new get orders options
func NewGetOrdersOptions(options ...GetOrdersOptFunc) (*GetOrdersOptions, error) {
	opts := GetOrdersOptions{
		PaginationOptions: &PaginationOptions{
			Page:     0,
			PageSize: 10,
		},
	}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}
	return &opts, nil
}

// PVZOrderRepository is an interface for order repository
type PVZOrderRepository interface {
	CreateOrder(order domain.PVZOrder) error
	DeleteOrder(orderID string) error
	SetOrderIssued(orderID string) error
	SetOrderReturned(orderID string) error
	GetOrders(userID string, options ...GetOrdersOptFunc) ([]domain.PVZOrder, error)
	GetOrder(orderID string) (domain.PVZOrder, error)
	GetReturns(options ...PaginationOptFunc) ([]domain.PVZOrder, error)
}

// IPVZOrderUseCase is an interface for order use cases
type IPVZOrderUseCase interface {
	AcceptOrderDelivery(orderID, recipientID string, storageTime time.Duration) error
	ReturnOrderDelivery(orderID string) error
	GiveOrderToClient(orderIDs []string) error
	GetOrders(userID string, options ...GetOrdersOptFunc) ([]domain.PVZOrder, error)
	AcceptReturn(userID, orderID string) error
	GetReturns(options ...PaginationOptFunc) ([]domain.PVZOrder, error)
}
