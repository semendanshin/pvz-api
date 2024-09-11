package abstractions

import (
	"homework/internal/domain"

	"time"
)

// PagePaginationOptions is a struct for pagination options
type PagePaginationOptions struct {
	Page     int
	PageSize int
}

// PagePaginationOptFunc is a type for pagination options
type PagePaginationOptFunc func(*PagePaginationOptions) error

// WithPage is an option to get orders by page
func WithPage(page int) PagePaginationOptFunc {
	return func(o *PagePaginationOptions) error {
		o.Page = page
		return nil
	}
}

// WithPageSize is an option to get orders by page size
func WithPageSize(pageSize int) PagePaginationOptFunc {
	return func(o *PagePaginationOptions) error {
		o.PageSize = pageSize
		return nil
	}
}

// NewPaginationOptions creates new pagination options
func NewPaginationOptions(options ...PagePaginationOptFunc) (*PagePaginationOptions, error) {
	opts := &PagePaginationOptions{
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
	LastNOrders     int
	PVZID           string
	CursorCreatedAt time.Time
	Limit           int
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

// WithCursorCreatedAt is an option to get orders by cursor created at
func WithCursorCreatedAt(cursorCreatedAt time.Time) GetOrdersOptFunc {
	return func(o *GetOrdersOptions) error {
		o.CursorCreatedAt = cursorCreatedAt
		return nil
	}
}

// WithLimit is an option to get orders by limit
func WithLimit(limit int) GetOrdersOptFunc {
	return func(o *GetOrdersOptions) error {
		o.Limit = limit
		return nil
	}
}

// NewGetOrdersOptions creates new get orders options
func NewGetOrdersOptions(options ...GetOrdersOptFunc) (*GetOrdersOptions, error) {
	opts := GetOrdersOptions{}
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
	GetReturns(options ...PagePaginationOptFunc) ([]domain.PVZOrder, error)
}

// IPVZOrderUseCase is an interface for order use cases
type IPVZOrderUseCase interface {
	AcceptOrderDelivery(orderID, recipientID string, storageTime time.Duration) error
	ReturnOrderDelivery(orderID string) error
	GiveOrderToClient(orderIDs []string) error
	GetOrders(userID string, options ...GetOrdersOptFunc) ([]domain.PVZOrder, error)
	AcceptReturn(userID, orderID string) error
	GetReturns(options ...PagePaginationOptFunc) ([]domain.PVZOrder, error)
}
