package json

import (
	"context"
	"encoding/json"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/usecases"
	"os"
	"slices"
	"time"
)

var _ usecases.PVZOrderRepository = &JSONRepository{}

type pvzOrder struct {
	OrderID     string `json:"order_id"`
	PVZID       string `json:"pvz_id"`
	RecipientID string `json:"recipient_id"`

	Weight         int    `json:"weight"`
	Cost           int    `json:"cost"`
	Packaging      string `json:"packaging"`
	AdditionalFilm bool   `json:"additional_film"`

	ReceivedAt  time.Time     `json:"received_at"`
	StorageTime time.Duration `json:"storage_time"`

	IssuedAt   time.Time `json:"issued_at"`
	ReturnedAt time.Time `json:"returned_at"`

	DeletedAt time.Time `json:"deleted_at" `
}

type fileStruct struct {
	Orders []pvzOrder `json:"orders"`
}

func ensureExists(pathToFile string) error {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		err = writeFile(pathToFile, fileStruct{Orders: make([]pvzOrder, 0)})
		if err != nil {
			return err
		}
	}

	return nil
}

func readFile(pathToFile string) (fileStruct, error) {
	err := ensureExists(pathToFile)
	if err != nil {
		return fileStruct{}, err
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		return fileStruct{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	var result fileStruct
	err = decoder.Decode(&result)
	if err != nil {
		return fileStruct{}, err
	}

	return result, nil
}

func writeFile(pathToFile string, fileStruct fileStruct) error {
	file, err := os.Create(pathToFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(fileStruct)
	if err != nil {
		return err
	}

	return nil
}

func convertToDomain(order pvzOrder) domain.PVZOrder {
	return domain.PVZOrder{
		OrderID:        order.OrderID,
		PVZID:          order.PVZID,
		RecipientID:    order.RecipientID,
		ReceivedAt:     order.ReceivedAt,
		StorageTime:    order.StorageTime,
		IssuedAt:       order.IssuedAt,
		ReturnedAt:     order.ReturnedAt,
		Weight:         order.Weight,
		Cost:           order.Cost,
		AdditionalFilm: order.AdditionalFilm,
		Packaging:      domain.PackagingType(order.Packaging),
	}
}

func convertToRepo(order domain.PVZOrder) pvzOrder {
	return pvzOrder{
		OrderID:        order.OrderID,
		PVZID:          order.PVZID,
		RecipientID:    order.RecipientID,
		ReceivedAt:     order.ReceivedAt,
		StorageTime:    order.StorageTime,
		IssuedAt:       order.IssuedAt,
		ReturnedAt:     order.ReturnedAt,
		DeletedAt:      time.Time{},
		Weight:         order.Weight,
		Cost:           order.Cost,
		AdditionalFilm: order.AdditionalFilm,
		Packaging:      string(order.Packaging),
	}
}

// JSONRepository is a JSON repository for PVZ orders
type JSONRepository struct {
	pathToFile string
}

// NewJSONRepository creates a new JSON repository
func NewJSONRepository(pathToFile string) *JSONRepository {
	return &JSONRepository{pathToFile: pathToFile}
}

// CreateOrder creates a new order
func (J *JSONRepository) CreateOrder(_ context.Context, order domain.PVZOrder) error {
	fileStruct, err := readFile(J.pathToFile)
	if err != nil {
		return err
	}

	fileStruct.Orders = append(fileStruct.Orders, convertToRepo(order))

	err = writeFile(J.pathToFile, fileStruct)

	return err
}

// DeleteOrder deletes an order
func (J *JSONRepository) DeleteOrder(_ context.Context, orderID string) error {
	fileStruct, err := readFile(J.pathToFile)
	if err != nil {
		return err
	}

	for i, order := range fileStruct.Orders {
		if order.OrderID == orderID {
			fileStruct.Orders[i].DeletedAt = time.Now()
			break
		}
	}

	err = writeFile(J.pathToFile, fileStruct)

	return err
}

// SetOrderIssued sets the order as issued
func (J *JSONRepository) SetOrderIssued(_ context.Context, orderID string) error {
	fileStruct, err := readFile(J.pathToFile)
	if err != nil {
		return err
	}

	for i, order := range fileStruct.Orders {
		if order.OrderID == orderID {
			fileStruct.Orders[i].IssuedAt = time.Now()
			break
		}
	}

	err = writeFile(J.pathToFile, fileStruct)

	return err
}

// SetOrderReturned sets the order as returned
func (J *JSONRepository) SetOrderReturned(_ context.Context, orderID string) error {
	fileStruct, err := readFile(J.pathToFile)
	if err != nil {
		return err
	}

	for i, order := range fileStruct.Orders {
		if order.OrderID == orderID {
			fileStruct.Orders[i].ReturnedAt = time.Now()
			break
		}
	}

	err = writeFile(J.pathToFile, fileStruct)

	return err
}

func filter(order pvzOrder, userID string, options abstractions.GetOrdersOptions) bool {
	return order.RecipientID == userID && (options.PVZID == "" || order.PVZID == options.PVZID) && order.IssuedAt.IsZero()
}

func getOrders(fileStruct fileStruct, userID string, options abstractions.GetOrdersOptions) []domain.PVZOrder {
	orders := make([]domain.PVZOrder, 0)
	for _, order := range fileStruct.Orders {
		if !filter(order, userID, options) {
			continue
		}
		orders = append(orders, convertToDomain(order))
	}

	return orders
}

// GetOrders gets orders
func (J *JSONRepository) GetOrders(_ context.Context, userID string, options ...abstractions.GetOrdersOptFunc) ([]domain.PVZOrder, error) {
	fileStruct, err := readFile(J.pathToFile)
	if err != nil {
		return nil, err
	}

	getOrdersOptions, err := abstractions.NewGetOrdersOptions(options...)
	if err != nil {
		return nil, err
	}

	orders := getOrders(fileStruct, userID, *getOrdersOptions)
	sortOrdersByReceivedAt(orders)

	orders = applyLastNOrdersFilter(orders, getOrdersOptions.LastNOrders)
	orders = applyCursorIDFilter(orders, getOrdersOptions.CursorID)
	orders = applyLimitFilter(orders, getOrdersOptions.Limit)

	return orders, nil
}

func sortOrdersByReceivedAt(orders []domain.PVZOrder) {
	slices.SortFunc(orders, func(i, j domain.PVZOrder) int {
		return int(j.ReceivedAt.Sub(i.ReceivedAt))
	})
}

func applyLastNOrdersFilter(orders []domain.PVZOrder, lastNOrders int) []domain.PVZOrder {
	if lastNOrders != 0 && len(orders) > lastNOrders {
		return orders[:lastNOrders]
	}
	return orders
}

func applyCursorIDFilter(orders []domain.PVZOrder, cursorID string) []domain.PVZOrder {
	if cursorID == "" {
		return orders
	}
	for i, order := range orders {
		if order.OrderID == cursorID {
			return orders[i:]
		}
	}
	return orders
}

func applyLimitFilter(orders []domain.PVZOrder, limit int) []domain.PVZOrder {
	if limit != 0 && len(orders) > limit {
		return orders[:limit]
	}
	return orders
}

// GetOrder gets an order
func (J *JSONRepository) GetOrder(_ context.Context, orderID string) (domain.PVZOrder, error) {
	fileStruct, err := readFile(J.pathToFile)
	if err != nil {
		return domain.PVZOrder{}, err
	}

	for _, order := range fileStruct.Orders {
		if order.OrderID == orderID && order.DeletedAt.IsZero() {
			return convertToDomain(order), nil
		}
	}

	return domain.PVZOrder{}, domain.ErrNotFound
}

func getReturns(fileStruct fileStruct) []domain.PVZOrder {
	returns := make([]domain.PVZOrder, 0)
	for _, order := range fileStruct.Orders {
		if !order.ReturnedAt.IsZero() || !order.DeletedAt.IsZero() {
			returns = append(returns, convertToDomain(order))
		}
	}

	return returns
}

// GetReturns gets returns
func (J *JSONRepository) GetReturns(_ context.Context, options ...abstractions.PagePaginationOptFunc) ([]domain.PVZOrder, error) {
	fileStruct, err := readFile(J.pathToFile)
	if err != nil {
		return nil, err
	}

	paginationOptions, err := abstractions.NewPaginationOptions(options...)
	if err != nil {
		return nil, err
	}

	returns := getReturns(fileStruct)

	returns = applyPagination(returns, *paginationOptions)

	return returns, nil
}

func applyPagination(orders []domain.PVZOrder, options abstractions.PagePaginationOptions) []domain.PVZOrder {
	if options.Page*options.PageSize >= len(orders) {
		return []domain.PVZOrder{}
	}

	orders = orders[options.Page*options.PageSize:]

	if len(orders) > options.PageSize {
		orders = orders[:options.PageSize]
	}
	return orders
}
