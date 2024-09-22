package main

import (
	"fmt"
	"homework/cmd/hw1/cmds"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
	"homework/internal/usecases/packager"
	"homework/internal/usecases/packager/strategies"
	"os"
)

func Run() error {
	// Можно получить значение через стандартный пакет flag, но тогда это не будет отображаться в help
	// ordersFile := flag.String("orders", "orders.json", "orders file")
	// pvzID := flag.String("pvz", "1", "pvz id")

	ordersFile := "orders.json"
	pvzID := "1"

	pvzOrderUseCase := initUseCase(ordersFile, pvzID)

	return cmds.Execute(pvzOrderUseCase)
}

func initUseCase(ordersFile string, pvzID string) abstractions.IPVZOrderUseCase {
	pvzOrderRepository := pvzorder.NewJSONRepository(ordersFile)

	orderPackager := packager.NewOrderPackager(
		map[domain.PackagingType]abstractions.OrderPackagerStrategy{
			domain.PackagingTypeBox:  strategies.NewBoxPackager(),
			domain.PackagingTypeFilm: strategies.NewFilmPackager(),
			domain.PackagingTypeBag:  strategies.NewBagPackager(),
		},
	)

	return usecases.NewPVZOrderUseCase(
		pvzOrderRepository,
		orderPackager,
		pvzID,
	)
}

func main() {
	if err := Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
