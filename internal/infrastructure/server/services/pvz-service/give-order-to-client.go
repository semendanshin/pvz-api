package pvz_service

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) GiveOrderToClient(ctx context.Context, req *desc.GiveOrderToClientRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PVZService.GiveOrderToClient")
	defer span.Finish()

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
