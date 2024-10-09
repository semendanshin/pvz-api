package pvz_service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) AcceptReturn(ctx context.Context, req *desc.AcceptReturnRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	err := p.useCase.AcceptReturn(
		ctx,
		req.GetUserId(),
		req.GetOrderId(),
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
