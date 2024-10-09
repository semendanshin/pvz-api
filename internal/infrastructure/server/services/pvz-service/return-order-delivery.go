package pvz_service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) ReturnOrderDelivery(ctx context.Context, req *desc.ReturnOrderDeliveryRequest) (*emptypb.Empty, error) {
	return nil, nil
}
