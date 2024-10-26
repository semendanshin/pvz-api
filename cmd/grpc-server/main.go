package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/clients/cache/inmemmory"
	"homework/internal/infrastructure/repositories/pvzorder/pgx"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
	"homework/internal/infrastructure/server"
	"homework/internal/tracer"
	"homework/internal/usecases"
	"homework/internal/usecases/packager"
	"homework/internal/usecases/packager/strategies"
	"log"
	"os"
	"time"
)

func loadPostgresURL() string {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUsername := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUsername, postgresPassword, postgresDatabase)
}

func Run() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	pvzID := os.Getenv("PVZ_ID")
	if pvzID == "" {
		return fmt.Errorf("PVZ_ID must be set")
	}

	postgresURL := loadPostgresURL()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tracer.MustSetup(ctx, "pvz-order-service")

	pool, err := pgxpool.New(ctx, postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	pvzOrderUseCase := initUseCase(pvzID, pool)

	grpcServer := server.NewGRPCServer(pvzOrderUseCase)

	return grpcServer.Run(ctx, "localhost", 8080, 8081)
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

	cache := inmemmory.NewPVZOrder(time.Second, 100, inmemmory.NewLRUInvalidationStrategy[string, interface{}]())

	return usecases.NewPVZOrderUseCase(
		pvzOrderRepoFacade,
		orderPackager,
		pvzID,
		cache,
	)
}

func main() {
	if err := Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
