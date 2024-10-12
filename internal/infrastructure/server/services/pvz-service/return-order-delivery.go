package pvz_service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) ReturnOrderDelivery(ctx context.Context, req *desc.ReturnOrderDeliveryRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidArgument, err)
	}

	err := p.useCase.ReturnOrderDelivery(
		ctx,
		req.GetOrderId(),
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
