package json

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/pvzorder/json"
	"os"
	"testing"
	"time"
)

type JSONRepositorySuite struct {
	suite.Suite
	storage *json.JSONRepository
}

func (s *JSONRepositorySuite) SetupTest() {
	fileName := fmt.Sprintf("test_%v.json", time.Now().Unix())

	err := copyFile("sample_data.json", fileName)
	s.Require().NoError(err)

	s.storage = json.NewJSONRepository(fileName)

	s.T().Cleanup(func() {
		err := os.Remove(fileName)
		s.Require().NoError(err)
	})
}

func (s *JSONRepositorySuite) TestCreateOrder() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	err := s.storage.CreateOrder(ctx, order)
	s.Require().NoError(err)

	actualOrder, err := s.storage.GetOrder(ctx, "100")
	s.Require().NoError(err)

	s.Equal(order.ReceivedAt.UnixMilli(), actualOrder.ReceivedAt.UnixMilli())
	order.ReceivedAt = actualOrder.ReceivedAt
	s.Equal(order, actualOrder)
}

func (s *JSONRepositorySuite) TestDeleteOrder() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := s.storage.DeleteOrder(ctx, "1")
	s.Require().NoError(err)

	_, err = s.storage.GetOrder(ctx, "1")
	s.Error(err)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *JSONRepositorySuite) TestSetOrderIssued() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := s.storage.SetOrderIssued(ctx, "1")
	s.Require().NoError(err)

	order, err := s.storage.GetOrder(ctx, "1")
	s.Require().NoError(err)

	s.NotZero(order.IssuedAt)
}

func (s *JSONRepositorySuite) TestSetOrderReturned() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := s.storage.SetOrderReturned(ctx, "1")
	s.Require().NoError(err)

	order, err := s.storage.GetOrder(ctx, "1")
	s.Require().NoError(err)

	s.NotZero(order.ReturnedAt)
}

func (s *JSONRepositorySuite) TestGetOrders() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type args struct {
		userID string
		opts   []abstractions.GetOrdersOptFunc
	}

	tests := []struct {
		name string
		args args
		want func([]domain.PVZOrder) bool
	}{
		{
			name: "Success",
			args: args{
				userID: "1",
				opts:   make([]abstractions.GetOrdersOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return s.Len(orders, 3)
			},
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
				return s.Len(orders, 1) && s.Equal("3", orders[0].OrderID)
			},
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
				return s.Len(orders, 2) && s.Equal("3", orders[0].OrderID) && s.Equal("2", orders[1].OrderID)
			},
		},
		{
			name: "Success empty",
			args: args{
				userID: "3",
				opts:   make([]abstractions.GetOrdersOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return s.Len(orders, 0)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orders, err := s.storage.GetOrders(ctx, tt.args.userID, tt.args.opts...)
			tt.want(orders)
			s.NoError(err)
		})
	}
}

func (s *JSONRepositorySuite) TestGetOrder() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
		s.Run(tt.name, func() {
			order, err := s.storage.GetOrder(ctx, tt.args.orderID)
			tt.wantErr(s.T(), err)
			s.Equal(tt.want, order)
		})
	}
}

func TestJSONRepository_Run(t *testing.T) {
	suite.Run(t, &JSONRepositorySuite{})
}
