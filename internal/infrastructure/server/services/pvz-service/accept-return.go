package pvz_service

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) AcceptReturn(ctx context.Context, req *desc.AcceptReturnRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PVZService.AcceptReturn")
	defer span.Finish()

	if err := req.ValidateAll(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidArgument, err)
	}

	err := p.useCase.AcceptReturn(
		ctx,
		req.GetUserId(),
		req.GetOrderId(),
	)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
