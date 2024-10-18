package pvz_service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func packagingTypeFromProto(packaging desc.PackagingType) domain.PackagingType {
	switch packaging {
	case desc.PackagingType_BOX:
		return domain.PackagingTypeBox
	case desc.PackagingType_BAG:
		return domain.PackagingTypeBag
	case desc.PackagingType_FILM:
		return domain.PackagingTypeFilm
	default:
		return domain.PackagingTypeUnknown
	}
}

func (p *PVZService) AcceptOrderDelivery(ctx context.Context, req *desc.AcceptOrderDeliveryRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidArgument, err)
	}

	err := p.useCase.AcceptOrderDelivery(
		ctx,
		req.GetOrderId(),
		req.GetRecipientId(),
		req.GetStorageTime().AsDuration(),
		int(req.GetCost()),
		int(req.GetWeight()),
		packagingTypeFromProto(req.Packaging),
		req.GetAdditionalFilm(),
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
