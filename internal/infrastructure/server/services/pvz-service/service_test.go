package pvz_service

import (
	"context"
	"net"
	"testing"
	"time"

	"homework/internal/abstractions"
	"homework/internal/abstractions/mocks"
	"homework/internal/domain"
	"homework/internal/infrastructure/server/middleware"
	desc "homework/pkg/pvz-service/v1"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	buffer = 1024 * 1024
)

func setupSuite(useCase abstractions.IPVZOrderUseCase) (desc.PvzServiceClient, func()) {
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.NewErrorMiddleware(),
		),
	)

	desc.RegisterPvzServiceServer(baseServer, NewPVZService(useCase))

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	return desc.NewPvzServiceClient(conn), func() {
		defer conn.Close()
		baseServer.Stop()
	}
}

func TestPVZService_AcceptOrderDelivery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctrl := minimock.NewController(t)
	useCase := mocks.NewIPVZOrderUseCaseMock(ctrl)

	client, teardown := setupSuite(useCase)
	defer teardown()

	type args struct {
		body *desc.AcceptOrderDeliveryRequest
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				body: &desc.AcceptOrderDeliveryRequest{
					OrderId:        "orderID",
					RecipientId:    "recipientID",
					StorageTime:    durationpb.New(2 * 24 * 60 * 60 * 1000000000),
					Cost:           10000,
					Weight:         1000,
					Packaging:      desc.PackagingType_BOX,
					AdditionalFilm: false,
				},
			},
			setup: func() {
				useCase.AcceptOrderDeliveryMock.Expect(
					minimock.AnyContext,
					"orderID",
					"recipientID",
					time.Duration(2*24*60*60*1000000000),
					10000,
					1000,
					domain.PackagingTypeBox,
					false,
				).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "order already exists",
			args: args{
				body: &desc.AcceptOrderDeliveryRequest{
					OrderId:        "orderID",
					RecipientId:    "recipientID",
					StorageTime:    durationpb.New(2 * 24 * 60 * 60 * 1000000000),
					Cost:           10000,
					Weight:         1000,
					Packaging:      desc.PackagingType_BOX,
					AdditionalFilm: false,
				},
			},
			setup: func() {
				useCase.AcceptOrderDeliveryMock.Expect(
					minimock.AnyContext,
					"orderID",
					"recipientID",
					time.Duration(2*24*60*60*1000000000),
					10000,
					1000,
					domain.PackagingTypeBox,
					false,
				).Return(domain.ErrAlreadyExists)
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Error(t, err)
				code, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.AlreadyExists, code.Code())
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			_, err := client.AcceptOrderDelivery(
				ctx,
				tt.args.body,
			)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZService_AcceptReturn(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctrl := minimock.NewController(t)
	useCase := mocks.NewIPVZOrderUseCaseMock(ctrl)

	client, teardown := setupSuite(useCase)
	defer teardown()

	type args struct {
		body *desc.AcceptReturnRequest
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				body: &desc.AcceptReturnRequest{
					OrderId: "orderID",
					UserId:  "userID",
				},
			},
			setup: func() {
				useCase.AcceptReturnMock.Expect(
					minimock.AnyContext,
					"userID",
					"orderID",
				).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "not found",
			args: args{
				body: &desc.AcceptReturnRequest{
					OrderId: "orderID",
					UserId:  "userID",
				},
			},
			setup: func() {
				useCase.AcceptReturnMock.Expect(
					minimock.AnyContext,
					"userID",
					"orderID",
				).Return(domain.ErrNotFound)
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Error(t, err)
				code, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, code.Code())
				return true
			},
		},
		{
			name: "invalid argument",
			args: args{
				body: &desc.AcceptReturnRequest{
					OrderId: "orderID",
					UserId:  "userID",
				},
			},
			setup: func() {
				useCase.AcceptReturnMock.Expect(
					minimock.AnyContext,
					"userID",
					"orderID",
				).Return(domain.ErrInvalidArgument)
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Error(t, err)
				code, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, code.Code())
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			_, err := client.AcceptReturn(
				ctx,
				tt.args.body,
			)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZService_GiveOrderToClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctrl := minimock.NewController(t)
	useCase := mocks.NewIPVZOrderUseCaseMock(ctrl)

	client, teardown := setupSuite(useCase)
	defer teardown()

	type args struct {
		body *desc.GiveOrderToClientRequest
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				body: &desc.GiveOrderToClientRequest{
					OrderIds: []string{"orderID"},
				},
			},
			setup: func() {
				useCase.GiveOrderToClientMock.Expect(
					minimock.AnyContext,
					[]string{"orderID"},
				).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "empty orderIDs",
			args: args{
				body: &desc.GiveOrderToClientRequest{
					OrderIds: []string{},
				},
			},
			setup: func() {},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Error(t, err)
				code, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, code.Code())
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			_, err := client.GiveOrderToClient(
				ctx,
				tt.args.body,
			)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZService_GetOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctrl := minimock.NewController(t)
	useCase := mocks.NewIPVZOrderUseCaseMock(ctrl)

	client, teardown := setupSuite(useCase)
	defer teardown()

	type args struct {
		body *desc.GetOrdersRequest
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				body: &desc.GetOrdersRequest{
					UserId: "userID",
				},
			},
			setup: func() {
				useCase.GetOrdersMock.Expect(
					minimock.AnyContext,
					"userID",
				).Return([]domain.PVZOrder{}, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "empty userID",
			args: args{
				body: &desc.GetOrdersRequest{
					UserId: "",
				},
			},
			setup: func() {},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Error(t, err)
				code, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, code.Code())
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			_, err := client.GetOrders(
				ctx,
				tt.args.body,
			)
			tt.wantErr(t, err)
		})
	}
}

func TestPVZService_ReturnOrderDelivery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctrl := minimock.NewController(t)
	useCase := mocks.NewIPVZOrderUseCaseMock(ctrl)

	client, teardown := setupSuite(useCase)
	defer teardown()

	type args struct {
		body *desc.ReturnOrderDeliveryRequest
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				body: &desc.ReturnOrderDeliveryRequest{
					OrderId: "orderID",
				},
			},
			setup: func() {
				useCase.ReturnOrderDeliveryMock.Expect(
					minimock.AnyContext,
					"orderID",
				).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "not found",
			args: args{
				body: &desc.ReturnOrderDeliveryRequest{
					OrderId: "orderID",
				},
			},
			setup: func() {
				useCase.ReturnOrderDeliveryMock.Expect(
					minimock.AnyContext,
					"orderID",
				).Return(domain.ErrNotFound)
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Error(t, err)
				code, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, code.Code())
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			_, err := client.ReturnOrderDelivery(
				ctx,
				tt.args.body,
			)
			tt.wantErr(t, err)
		})
	}
}
