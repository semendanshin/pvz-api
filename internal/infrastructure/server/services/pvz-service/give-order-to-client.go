package pvz_service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) GiveOrdersToClient(ctx context.Context, req *desc.GiveOrderToClientRequest) (*emptypb.Empty, error) {
	return nil, nil
}
