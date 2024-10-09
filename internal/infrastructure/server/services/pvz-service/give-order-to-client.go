package pvz_service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) GiveOrdersToClient(ctx context.Context, req *desc.GiveOrderToClientRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
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
