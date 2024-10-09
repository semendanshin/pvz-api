package pvz_service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) AcceptOrderDelivery(ctx context.Context, req *desc.AcceptOrderDeliveryRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	packagingType, err := domain.NewPackagingType(req.GetPackaging().String())
	if err != nil {
		return nil, err
	}

	err = p.useCase.AcceptOrderDelivery(
		ctx,
		req.GetOrderId(),
		req.GetRecipientId(),
		req.GetStorageTime().AsDuration(),
		int(req.GetCost()),
		int(req.GetWeight()),
		packagingType,
		req.GetAdditionalFilm(),
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
