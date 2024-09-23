package main

import (
	"context"
	"fmt"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
	"log"
	"os"

	"homework/cmd/hw1/cmds"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/pvzorder/pgx"
	"homework/internal/usecases"
	"homework/internal/usecases/packager"
	"homework/internal/usecases/packager/strategies"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() error {
	// Можно получить значение через стандартный пакет flag, но тогда это не будет отображаться в help
	// pvzID := flag.String("pvz", "1", "pvz id")

	pvzID := "1"

	// Наверное не очень круто будет передавать это через флаг. Я подумаю в сторону переменных окружения и .env файла
	postgresURL := "host=localhost port=5430 user=test password=test dbname=test sslmode=disable"

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pool, err := pgxpool.New(ctx, postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	pvzOrderUseCase := initUseCase(pvzID, pool)

	return cmds.Execute(ctx, pvzOrderUseCase)
}

func initUseCase(pvzID string, pool *pgxpool.Pool) abstractions.IPVZOrderUseCase {
	txManager := txmanager.NewPGXTXManager(pool)
	pvzOrderRepoFacade := pgx.NewPgxPvzOrderFacade(txManager)

	orderPackager := packager.NewOrderPackager(
		map[domain.PackagingType]packager.OrderPackagerStrategy{
			domain.PackagingTypeBox:  strategies.NewBoxPackager(),
			domain.PackagingTypeFilm: strategies.NewFilmPackager(),
			domain.PackagingTypeBag:  strategies.NewBagPackager(),
		},
	)

	return usecases.NewPVZOrderUseCase(
		pvzOrderRepoFacade,
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
