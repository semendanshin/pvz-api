package pvz_service

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) ReturnOrderDelivery(ctx context.Context, req *desc.ReturnOrderDeliveryRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PVZService.ReturnOrderDelivery")
	defer span.Finish()

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

	return &emptypb.Empty{}, nil
}
