package pvz_service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) GiveOrderToClient(ctx context.Context, req *desc.GiveOrderToClientRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidArgument, err)
	}

	err := p.useCase.GiveOrderToClient(
		ctx,
		req.GetOrderIds(),
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
