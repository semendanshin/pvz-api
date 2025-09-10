package usecases

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecases/mocks"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tests := []struct {
		name    string
		args    args
		setup   func(repoMock *mocks.PVZOrderRepositoryMock, packagerMock *mocks.OrderPackagerInterfaceMock, cacheMock *mocks.PVZOrderCacheMock)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderID:        "orderID",
				recipientID:    "recipientID",
				storageTime:    1 * time.Hour,
				cost:           100,
				weight:         1,
				packaging:      domain.PackagingTypeBox,
				additionalFilm: false,
			},
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, packagerMock *mocks.OrderPackagerInterfaceMock, _ *mocks.PVZOrderCacheMock) {
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
				cost:           100,
				weight:         1,
				packaging:      domain.PackagingTypeBox,
				additionalFilm: false,
			},
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, _ *mocks.OrderPackagerInterfaceMock, _ *mocks.PVZOrderCacheMock) {
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
				cost:           100,
				weight:         1,
				packaging:      domain.PackagingTypeFilm,
				additionalFilm: true,
			},
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, _ *mocks.OrderPackagerInterfaceMock, _ *mocks.PVZOrderCacheMock) {
				repoMock.GetOrderMock.Return(domain.PVZOrder{}, domain.ErrNotFound)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := minimock.NewController(t)
			repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)
			packagerMock := mocks.NewOrderPackagerInterfaceMock(ctrl)
			cacheMock := mocks.NewPVZOrderCacheMock(ctrl)
			uc := NewPVZOrderUseCase(repoMock, packagerMock, pvzID, cacheMock)
			tt.setup(repoMock, packagerMock, cacheMock)
			err := uc.AcceptOrderDelivery(ctx, tt.args.orderID, tt.args.recipientID, tt.args.storageTime, tt.args.cost, tt.args.weight, tt.args.packaging, tt.args.additionalFilm)
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tests := []struct {
		name    string
		args    args
		setup   func(repoMock *mocks.PVZOrderRepositoryMock, cacheMock *mocks.PVZOrderCacheMock)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderID: "orderID",
			},
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, cacheMock *mocks.PVZOrderCacheMock) {
				cacheMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{PVZID: pvzID, ReceivedAt: time.Now().Add(-3 * time.Hour), StorageTime: 2 * time.Hour}
				repoMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cacheMock.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
				repoMock.DeleteOrderMock.Expect(minimock.AnyContext, "orderID").Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				orderID: "orderID",
			},
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, cacheMock *mocks.PVZOrderCacheMock) {
				cacheMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				repoMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, domain.ErrNotFound)
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
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, cacheMock *mocks.PVZOrderCacheMock) {
				cacheMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{PVZID: "anotherPVZID", ReceivedAt: time.Now().Add(-3 * time.Hour), StorageTime: 2 * time.Hour}
				repoMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cacheMock.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
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
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, cacheMock *mocks.PVZOrderCacheMock) {
				cacheMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{ReceivedAt: time.Now(), StorageTime: 2 * time.Hour}
				repoMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cacheMock.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
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
			setup: func(repoMock *mocks.PVZOrderRepositoryMock, cacheMock *mocks.PVZOrderCacheMock) {
				cacheMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{IssuedAt: time.Now()}
				repoMock.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cacheMock.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := minimock.NewController(t)
			repo := mocks.NewPVZOrderRepositoryMock(ctrl)
			cache := mocks.NewPVZOrderCacheMock(ctrl)
			uc := NewPVZOrderUseCase(repo, nil, pvzID, cache)
			tt.setup(repo, cache)
			err := uc.ReturnOrderDelivery(ctx, tt.args.orderID)
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tests := []struct {
		name    string
		args    args
		setup   func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderIDs: []string{"orderID"},
			},
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{
					OrderID:     "orderID",
					RecipientID: "userID",
					PVZID:       pvzID,
					ReceivedAt:  time.Now().Add(-1 * time.Hour),
					StorageTime: 2 * time.Hour,
				}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
				repo.SetOrderIssuedMock.Expect(minimock.AnyContext, "orderID").Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				orderIDs: []string{"orderID"},
			},
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, domain.ErrNotFound)
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
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{
					OrderID:     "orderID",
					RecipientID: "userID",
					PVZID:       "anotherPVZID",
					IssuedAt:    time.Now(),
					ReceivedAt:  time.Now().Add(-1 * time.Hour),
					StorageTime: 2 * time.Hour,
				}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "Order storage time expired",
			args: args{
				orderIDs: []string{"orderID"},
			},
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{
					OrderID:     "orderID",
					RecipientID: "userID",
					PVZID:       pvzID,
					ReceivedAt:  time.Now().Add(-3 * time.Hour),
					StorageTime: 2 * time.Hour,
				}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
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
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				// Allow any order of calls and return data based on orderID
				cache.GetOrderMock.Set(func(_ context.Context, id string) (domain.PVZOrder, error, bool) {
					return domain.PVZOrder{}, nil, false
				})
				repo.GetOrderMock.Set(func(_ context.Context, id string) (domain.PVZOrder, error) {
					switch id {
					case "orderID":
						return domain.PVZOrder{
							OrderID:     "orderID",
							RecipientID: "userID",
							PVZID:       pvzID,
							ReceivedAt:  time.Now().Add(-1 * time.Hour),
							StorageTime: 2 * time.Hour,
						}, nil
					case "anotherOrderID":
						return domain.PVZOrder{
							OrderID:     "anotherOrderID",
							RecipientID: "anotherUserID",
							PVZID:       pvzID,
							ReceivedAt:  time.Now().Add(-1 * time.Hour),
							StorageTime: 2 * time.Hour,
						}, nil
					default:
						return domain.PVZOrder{}, errors.New("unexpected id")
					}
				})
				cache.SetOrderMock.Set(func(_ context.Context, _ domain.PVZOrder) error { return nil })
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := minimock.NewController(t)
			repo := mocks.NewPVZOrderRepositoryMock(ctrl)
			cache := mocks.NewPVZOrderCacheMock(ctrl)
			uc := NewPVZOrderUseCase(repo, nil, pvzID, cache)
			tt.setup(repo, cache)
			err := uc.GiveOrderToClient(ctx, tt.args.orderIDs)
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
	cacheMock := mocks.NewPVZOrderCacheMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, nil, pvzID, cacheMock)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
				cacheMock.GetOrdersMock.Expect(minimock.AnyContext, "userID").Return(nil, nil, false)
				got := []domain.PVZOrder{{RecipientID: "userID", PVZID: pvzID}}
				repoMock.GetOrdersMock.Expect(minimock.AnyContext, "userID").Return(got, nil)
				cacheMock.SetGetOrdersMock.Expect(minimock.AnyContext, "userID", got).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Cache hit",
			args: args{userID: "userID"},
			setup: func() {
				cacheMock.GetOrdersMock.Expect(minimock.AnyContext, "userID").Return([]domain.PVZOrder{{RecipientID: "userID", PVZID: pvzID}}, nil, true)
			},
			wantErr: assert.NoError,
		},
		{
			name: "No orders",
			args: args{
				userID: "userID",
			},
			setup: func() {
				cacheMock.GetOrdersMock.Expect(minimock.AnyContext, "userID").Return(nil, nil, false)
				got := []domain.PVZOrder{}
				repoMock.GetOrdersMock.Expect(minimock.AnyContext, "userID").Return(got, nil)
				cacheMock.SetGetOrdersMock.Expect(minimock.AnyContext, "userID", got).Return(nil)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := useCase.GetOrders(ctx, tt.args.userID)
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tests := []struct {
		name    string
		args    args
		setup   func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{RecipientID: "userID", IssuedAt: time.Now().Add(-(TimeForReturn - time.Hour))}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
				repo.SetOrderReturnedMock.Expect(minimock.AnyContext, "orderID").Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				userID:  "userID",
				orderID: "orderID",
			},
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, domain.ErrNotFound)
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
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{RecipientID: "anotherUserID"}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
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
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{RecipientID: "userID", ReturnedAt: time.Now()}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
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
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{RecipientID: "userID"}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
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
			setup: func(repo *mocks.PVZOrderRepositoryMock, cache *mocks.PVZOrderCacheMock) {
				cache.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(domain.PVZOrder{}, nil, false)
				order := domain.PVZOrder{RecipientID: "userID", IssuedAt: time.Now().Add(-(TimeForReturn + time.Hour))}
				repo.GetOrderMock.Expect(minimock.AnyContext, "orderID").Return(order, nil)
				cache.SetOrderMock.Expect(minimock.AnyContext, order).Return(nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := minimock.NewController(t)
			repo := mocks.NewPVZOrderRepositoryMock(ctrl)
			cache := mocks.NewPVZOrderCacheMock(ctrl)
			uc := NewPVZOrderUseCase(repo, nil, pvzID, cache)
			tt.setup(repo, cache)
			err := uc.AcceptReturn(ctx, tt.args.userID, tt.args.orderID)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZOrderUseCase_GetReturns(t *testing.T) {
	t.Parallel()

	const pvzID = "currentPVZID"

	ctrl := minimock.NewController(t)
	repoMock := mocks.NewPVZOrderRepositoryMock(ctrl)
	cacheMock := mocks.NewPVZOrderCacheMock(ctrl)

	useCase := NewPVZOrderUseCase(repoMock, nil, pvzID, cacheMock)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tests := []struct {
		name    string
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success (cache miss)",
			setup: func() {
				cacheMock.GetReturnsMock.Expect(minimock.AnyContext).Return(nil, nil, false)
				got := []domain.PVZOrder{{PVZID: pvzID}}
				repoMock.GetReturnsMock.Expect(minimock.AnyContext).Return(got, nil)
				cacheMock.SetGetReturnsMock.Expect(minimock.AnyContext, got).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success (cache hit)",
			setup: func() {
				cacheMock.GetReturnsMock.Expect(minimock.AnyContext).Return([]domain.PVZOrder{{PVZID: pvzID}}, nil, true)
			},
			wantErr: assert.NoError,
		},
		{
			name: "No returns",
			setup: func() {
				cacheMock.GetReturnsMock.Expect(minimock.AnyContext).Return(nil, nil, false)
				got := []domain.PVZOrder{}
				repoMock.GetReturnsMock.Expect(minimock.AnyContext).Return(got, nil)
				cacheMock.SetGetReturnsMock.Expect(minimock.AnyContext, got).Return(nil)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := useCase.GetReturns(ctx)
			tt.wantErr(t, err)
		})
	}
}
