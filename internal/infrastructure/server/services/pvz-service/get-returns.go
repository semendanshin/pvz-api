package pvz_service

import (
	"context"
	"fmt"
	"homework/internal/abstractions"
	"homework/internal/domain"
	desc "homework/pkg/pvz-service/v1"
)

func (p *PVZService) GetReturns(ctx context.Context, req *desc.GetReturnsRequest) (*desc.GetReturnsResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidArgument, err)
	}

	var options []abstractions.PagePaginationOptFunc
	if req.Page != nil {
		options = append(options, abstractions.WithPage(int(req.GetPage())))
	}
	if req.PageSize != nil {
		options = append(options, abstractions.WithPageSize(int(req.GetPageSize())))
	}

	returns, err := p.useCase.GetReturns(ctx, options...)
	if err != nil {
		return nil, err
	}

	var descReturns []*desc.PVZOrder
	for _, ret := range returns {
		descReturns = append(descReturns, domainToDescOrder(&ret))
	}

	return &desc.GetReturnsResponse{
		Returns: descReturns,
	}, nil
}
