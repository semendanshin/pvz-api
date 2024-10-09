package pvz_service

import (
	"context"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework/internal/abstractions"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func domainPackagingTypeToDesc(packagingType domain.PackagingType) desc.PackagingType {
	switch packagingType {
	case domain.PackagingTypeBag:
		return desc.PackagingType_BAG
	case domain.PackagingTypeBox:
		return desc.PackagingType_BOX
	case domain.PackagingTypeFilm:
		return desc.PackagingType_FILM
	default:
		return desc.PackagingType_UNKNOWN
	}
}

func domainToDescOrder(order *domain.PVZOrder) *desc.PVZOrder {
	return &desc.PVZOrder{
		OrderId:     order.OrderID,
		PvzId:       order.PVZID,
		RecipientId: order.RecipientID,

		Cost:   int32(order.Cost),
		Weight: int32(order.Weight),

		StorageTime: durationpb.New(order.StorageTime),
		ReceivedAt:  timestamppb.New(order.ReceivedAt),

		Packaging:      domainPackagingTypeToDesc(order.Packaging),
		AdditionalFilm: order.AdditionalFilm,

		IssuedAt:   timestamppb.New(order.IssuedAt),
		ReturnedAt: timestamppb.New(order.ReturnedAt),
	}
}

func (p *PVZService) GetOrders(ctx context.Context, req *desc.GetOrdersRequest) (*desc.GetOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	var options []abstractions.GetOrdersOptFunc
	if req.LastN != nil {
		options = append(options, abstractions.WithLastNOrders(int(req.GetLastN())))
	}
	if req.GetSamePVZ() {
		options = append(options, abstractions.WithSamePVZ())
	}
	if req.Cursor != nil {
		options = append(options, abstractions.WithCursorID(req.GetCursor()))
	}
	if req.Limit != nil {
		options = append(options, abstractions.WithLimit(int(req.GetLimit())))
	}

	orders, err := p.useCase.GetOrders(
		ctx,
		req.GetUserId(),
		options...,
	)
	if err != nil {
		return nil, err
	}

	result := make([]*desc.PVZOrder, 0, len(orders))
	for _, order := range orders {
		result = append(result, domainToDescOrder(&order))
	}

	return &desc.GetOrdersResponse{
		Orders: result,
	}, nil
}
