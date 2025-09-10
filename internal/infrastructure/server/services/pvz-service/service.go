package pvz_service

import (
	"homework/internal/abstractions"
	desc "homework/pkg/pvz-service/v1"
)

type PVZService struct {
	useCase abstractions.IPVZOrderUseCase

	desc.UnimplementedPvzServiceServer
}

func NewPVZService(useCase abstractions.IPVZOrderUseCase) *PVZService {
	return &PVZService{
		useCase: useCase,
	}
}
