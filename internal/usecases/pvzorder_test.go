package usecases

import (
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"homework/internal/domain"
	"homework/internal/usecases/mocks"
	"testing"
	"time"
)

func TestPVZOrderUseCase_AcceptOrderDelivery(t *testing.T) {
	t.Parallel()

	const pvzID = "currentPVZID"

	type args struct {
		orderID, recipientID string
		storageTime          time.Duration
		cost, weight         int
		packaging            domain.PackagingType
		additionalFilm       bool
	}

	ctrl := minimock.NewController(t)
	repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)
	packagerMock := mocks.NewOrderPackagerInterfaceMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, packagerMock, pvzID)

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderID:        "orderID",
				recipientID:    "recipientID",
				storageTime:    1 * time.Hour,
				cost:           100.00,
				weight:         1.000,
				packaging:      domain.PackagingTypeBox,
				additionalFilm: false,
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{}, domain.ErrNotFound)
				packagerMock.PackageOrderMock.Return(domain.PVZOrder{}, nil)
				repoMock.CreateOrderMock.Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order already exists",
			args: args{
				orderID:        "orderID",
				recipientID:    "recipientID",
				storageTime:    1 * time.Hour,
				cost:           100.00,
				weight:         1.000,
				packaging:      domain.PackagingTypeBox,
				additionalFilm: false,
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{}, nil)

			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrAlreadyExists)
			},
		},
		{
			name: "Additional film is not allowed for film packaging",
			args: args{
				orderID:        "orderID",
				recipientID:    "recipientID",
				storageTime:    1 * time.Hour,
				cost:           100.00,
				weight:         1.000,
				packaging:      domain.PackagingTypeFilm,
				additionalFilm: true,
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{}, domain.ErrNotFound)

			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := useCase.AcceptOrderDelivery(tt.args.orderID, tt.args.recipientID, tt.args.storageTime, tt.args.cost, tt.args.weight, tt.args.packaging, tt.args.additionalFilm)
			tt.wantErr(t, err)
		})
	}

}

func TestPVZOrderUseCase_ReturnOrderDelivery(t *testing.T) {
	t.Parallel()

	const pvzID = "currentPVZID"

	type args struct {
		orderID string
	}

	ctrl := minimock.NewController(t)
	repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, nil, pvzID)

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					PVZID: pvzID,
				}, nil)
				repoMock.DeleteOrderMock.Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{}, domain.ErrNotFound)

			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrNotFound)
			},
		},
		{
			name: "Orders is not for current PVZ",
			args: args{
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					PVZID: "anotherPVZID",
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "Storage time has not expired",
			args: args{
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					IssuedAt: time.Now(),
				}, nil)

			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "Order is already issued",
			args: args{
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					IssuedAt: time.Now(),
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := useCase.ReturnOrderDelivery(tt.args.orderID)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZOrderUseCase_GiveOrderToClient(t *testing.T) {
	t.Parallel()

	const pvzID = "currentPVZID"

	type args struct {
		orderIDs []string
	}

	ctrl := minimock.NewController(t)
	repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, nil, pvzID)

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderIDs: []string{"orderID"},
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
					PVZID:       pvzID,
					ReceivedAt:  time.Now().Add(-1 * time.Hour),
					StorageTime: 2 * time.Hour,
				}, nil)
				repoMock.SetOrderIssuedMock.Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				orderIDs: []string{"orderID"},
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{}, domain.ErrNotFound)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrNotFound)
			},
		},
		{
			name: "Order is not for current PVZ",
			args: args{
				orderIDs: []string{"orderID"},
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
					PVZID:       "anotherPVZID",
					IssuedAt:    time.Now(),
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "Order is not issued",
			args: args{
				orderIDs: []string{"orderID"},
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
					PVZID:       pvzID,
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "Orders do not belong to the same user",
			args: args{
				orderIDs: []string{"orderID", "anotherOrderID"},
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
					PVZID:       pvzID,
					IssuedAt:    time.Now(),
				}, nil)
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "anotherUserID",
					PVZID:       pvzID,
					IssuedAt:    time.Now(),
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := useCase.GiveOrderToClient(tt.args.orderIDs)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZOrderUseCase_GetOrders(t *testing.T) {
	t.Parallel()

	const pvzID = "currentPVZID"

	type args struct {
		userID string
	}

	ctrl := minimock.NewController(t)
	repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, nil, pvzID)

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID: "userID",
			},
			setup: func() {
				repoMock.GetOrdersMock.Return([]domain.PVZOrder{
					{
						RecipientID: "userID",
						PVZID:       pvzID,
					},
				}, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "No orders",
			args: args{
				userID: "userID",
			},
			setup: func() {
				repoMock.GetOrdersMock.Return([]domain.PVZOrder{}, nil)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := useCase.GetOrders(tt.args.userID)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZOrderUseCase_AcceptReturn(t *testing.T) {
	t.Parallel()

	const pvzID = "currentPVZID"

	type args struct {
		userID  string
		orderID string
	}

	ctrl := minimock.NewController(t)
	repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, nil, pvzID)

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
					IssuedAt:    time.Now().Add(-(TimeForReturn - time.Hour)),
				}, nil)
				repoMock.SetOrderReturnedMock.Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{}, domain.ErrNotFound)

			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrNotFound)
			},
		},
		{
			name: "User is not recipient",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "anotherUserID",
				}, nil)

			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)

			},
		},
		{
			name: "Order is already returned",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
					ReturnedAt:  time.Now(),
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "Order is not issued",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "Time for return has expired",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func() {
				repoMock.GetOrderMock.Return(domain.PVZOrder{
					RecipientID: "userID",
					IssuedAt:    time.Now().Add(-(TimeForReturn + time.Hour)),
				}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := useCase.AcceptReturn(tt.args.userID, tt.args.orderID)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZOrderUseCase_GetReturns(t *testing.T) {
	t.Parallel()

	const pvzID = "currentPVZID"

	ctrl := minimock.NewController(t)
	repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, nil, pvzID)

	tests := []struct {
		name    string
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			setup: func() {
				repoMock.GetReturnsMock.Return([]domain.PVZOrder{
					{
						PVZID: pvzID,
					},
				}, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "No returns",
			setup: func() {
				repoMock.GetReturnsMock.Return([]domain.PVZOrder{}, nil)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := useCase.GetReturns()
			tt.wantErr(t, err)
		})
	}
}
