package pvzorder

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/pvzorder"
	"os"
	"testing"
	"time"
)

var _ = []domain.PVZOrder{
	{
		OrderID:        "1",
		PVZID:          "1",
		RecipientID:    "1",
		ReceivedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		StorageTime:    24 * time.Hour,
		IssuedAt:       time.Time{},
		ReturnedAt:     time.Time{},
		Weight:         1000,
		Cost:           1000,
		AdditionalFilm: false,
		Packaging:      domain.PackagingTypeBox,
	},
	{
		OrderID:        "2",
		PVZID:          "1",
		RecipientID:    "1",
		ReceivedAt:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		StorageTime:    24 * time.Hour,
		IssuedAt:       time.Time{},
		ReturnedAt:     time.Time{},
		Weight:         1000,
		Cost:           1000,
		AdditionalFilm: false,
		Packaging:      domain.PackagingTypeBag,
	},
	{
		OrderID:        "3",
		PVZID:          "2",
		RecipientID:    "1",
		ReceivedAt:     time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		StorageTime:    24 * time.Hour,
		IssuedAt:       time.Time{},
		ReturnedAt:     time.Time{},
		Weight:         1000,
		Cost:           1000,
		AdditionalFilm: false,
		Packaging:      domain.PackagingTypeFilm,
	},
	{
		OrderID:        "4",
		PVZID:          "1",
		RecipientID:    "2",
		ReceivedAt:     time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
		StorageTime:    24 * time.Hour,
		IssuedAt:       time.Time{},
		ReturnedAt:     time.Time{},
		Weight:         1000,
		Cost:           1000,
		AdditionalFilm: false,
		Packaging:      domain.PackagingTypeBox,
	},
	{
		OrderID:        "5",
		PVZID:          "1",
		RecipientID:    "2",
		ReceivedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		StorageTime:    24 * time.Hour,
		IssuedAt:       time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		ReturnedAt:     time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Weight:         1000,
		Cost:           1000,
		AdditionalFilm: false,
		Packaging:      domain.PackagingTypeBox,
	},
	{
		OrderID:        "6",
		PVZID:          "1",
		RecipientID:    "2",
		ReceivedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		StorageTime:    24 * time.Hour,
		IssuedAt:       time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		ReturnedAt:     time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Weight:         1000,
		Cost:           1000,
		AdditionalFilm: false,
		Packaging:      domain.PackagingTypeBox,
	},
	{
		OrderID:        "7",
		PVZID:          "1",
		RecipientID:    "2",
		ReceivedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		StorageTime:    24 * time.Hour,
		IssuedAt:       time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		ReturnedAt:     time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Weight:         1000,
		Cost:           1000,
		AdditionalFilm: false,
		Packaging:      domain.PackagingTypeBox,
	},
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = source.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = destination.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = source.WriteTo(destination)
	if err != nil {
		return err
	}

	return nil
}

func setupTest(t *testing.T) (*pvzorder.JSONRepository, func()) {
	t.Helper()

	fileName := fmt.Sprintf("test_%v.json", time.Now().Unix())

	err := copyFile("sample_data.json", fileName)
	if err != nil {
		t.Fatal(err)
	}

	repo := pvzorder.NewJSONRepository(fileName)

	return repo, func() {
		err := os.Remove(fileName)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestJSONRepository_CreateOrder(t *testing.T) {
	repo, tearDown := setupTest(t)
	defer tearDown()

	order := domain.NewPVZOrder(
		"100",
		"1",
		"1",
		1000,
		1000,
		24*time.Hour,
		domain.PackagingTypeBox,
		false,
	)

	err := repo.CreateOrder(order)
	assert.NoError(t, err)

	// Check if the order was created
	actual, err := repo.GetOrder("100")
	assert.NoError(t, err)
	assert.Equal(t, order.ReceivedAt.UnixMilli(), actual.ReceivedAt.UnixMilli())
	// при маршалинге и демаршалинге время из наносекунд?? превращается в миллисекунды??, поэтому сравнивать их не получится
	// если убрать следующую строку, то тест будет ломаться, хотя фактически время правильное(
	order.ReceivedAt = actual.ReceivedAt
	assert.Equal(t, order, actual)
}

func TestJSONRepository_DeleteOrder(t *testing.T) {
	repo, tearDown := setupTest(t)
	defer tearDown()

	err := repo.DeleteOrder("1")
	assert.NoError(t, err)

	// Check if the order was deleted
	_, err = repo.GetOrder("1")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestJSONRepository_SetOrderIssued(t *testing.T) {
	repo, tearDown := setupTest(t)
	defer tearDown()

	err := repo.SetOrderIssued("1")
	assert.NoError(t, err)

	// Check if the order was issued
	order, err := repo.GetOrder("1")
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, order.IssuedAt)
}

func TestJSONRepository_SetOrderReturned(t *testing.T) {
	repo, tearDown := setupTest(t)
	defer tearDown()

	err := repo.SetOrderReturned("1")
	assert.NoError(t, err)

	// Check if the order was returned
	order, err := repo.GetOrder("1")
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, order.ReturnedAt)
}

func TestJSONRepository_GetOrders(t *testing.T) {
	repo, tearDown := setupTest(t)
	defer tearDown()

	type args struct {
		userID string
		opts   []abstractions.GetOrdersOptFunc
	}

	tests := []struct {
		name    string
		args    args
		want    func([]domain.PVZOrder) bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID: "1",
				opts:   make([]abstractions.GetOrdersOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 3)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success with PVZID",
			args: args{
				userID: "1",
				opts: []abstractions.GetOrdersOptFunc{
					abstractions.WithPVZID("2"),
				},
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 1) && assert.Equal(t, "3", orders[0].OrderID)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success with pagination",
			args: args{
				userID: "1",
				opts: []abstractions.GetOrdersOptFunc{
					abstractions.WithCursorID("3"),
					abstractions.WithLimit(2),
				},
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 2) && assert.Equal(t, "3", orders[0].OrderID) && assert.Equal(t, "2", orders[1].OrderID)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success empty",
			args: args{
				userID: "3",
				opts:   make([]abstractions.GetOrdersOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 0)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := repo.GetOrders(tt.args.userID, tt.args.opts...)
			tt.want(orders)
			tt.wantErr(t, err)
		})
	}
}

func TestJSONRepository_GetOrder(t *testing.T) {
	repo, tearDown := setupTest(t)
	defer tearDown()

	type args struct {
		orderID string
	}

	tests := []struct {
		name    string
		args    args
		want    domain.PVZOrder
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderID: "1",
			},
			want: domain.PVZOrder{
				OrderID:        "1",
				PVZID:          "1",
				RecipientID:    "1",
				ReceivedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				StorageTime:    24 * time.Hour,
				IssuedAt:       time.Time{},
				ReturnedAt:     time.Time{},
				Weight:         1000,
				Cost:           1000,
				AdditionalFilm: false,
				Packaging:      domain.PackagingTypeBox,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				orderID: "100",
			},
			want: domain.PVZOrder{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && errors.Is(err, domain.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetOrder(tt.args.orderID)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestJSONRepository_GetReturns(t *testing.T) {
	repo, tearDown := setupTest(t)
	defer tearDown()

	type args struct {
		opts []abstractions.PagePaginationOptFunc
	}

	tests := []struct {
		name    string
		args    args
		want    func([]domain.PVZOrder) bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				opts: make([]abstractions.PagePaginationOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 3)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success with pagination",
			args: args{
				opts: []abstractions.PagePaginationOptFunc{
					abstractions.WithPage(1),
					abstractions.WithPageSize(2),
				},
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 1) && assert.Equal(t, "7", orders[0].OrderID)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := repo.GetReturns(tt.args.opts...)
			tt.want(orders)
			tt.wantErr(t, err)
		})
	}
}
