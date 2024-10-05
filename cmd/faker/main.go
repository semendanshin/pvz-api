package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
	"homework/internal/usecases"
	"log"
	"os"
	"time"

	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/pvzorder/pgx"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"

	"github.com/jackc/pgx/v5/pgxpool"
)

func loadPostgresURL() string {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUsername := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUsername, postgresPassword, postgresDatabase)
}

func main() {
	count := flag.Int("count", 100, "count of orders")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	postgresURL := loadPostgresURL()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pool, err := pgxpool.New(ctx, postgresURL)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create pool: %w", err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(fmt.Errorf("failed to ping database: %w", err))
	}

	txManager := txmanager.NewPGXTXManager(pool)
	pvzOrderRepoFacade := pgx.NewPgxPvzOrderFacade(txManager)

	for i := 0; i < *count; i++ {
		err := createRandomOrder(ctx, pvzOrderRepoFacade)
		if err != nil {
			log.Printf("failed to create order: %v", err)
		}
	}
}

func createRandomOrder(ctx context.Context, repo *pgx.PvzOrderFacade) error {
	var order domain.PVZOrder
	{
		order.OrderID = gofakeit.LetterN(10)
		order.PVZID = gofakeit.LetterN(10)
		order.RecipientID = gofakeit.LetterN(10)

		order.Cost = gofakeit.Number(1000, 1000000)
		order.Weight = gofakeit.Number(10, 10000)

		order.Packaging = domain.PackagingType(gofakeit.RandomString([]string{"box", "film", "bag"}))
		if order.Packaging != domain.PackagingTypeFilm {
			order.AdditionalFilm = gofakeit.Bool()
		}

		order.ReceivedAt = gofakeit.PastDate()
		order.StorageTime = time.Duration(gofakeit.Number(1, 14)) * 24 * time.Hour

		if gofakeit.Bool() {
			order.IssuedAt = gofakeit.DateRange(order.ReceivedAt, order.ReceivedAt.Add(order.StorageTime))
		}
		if !order.IssuedAt.IsZero() && gofakeit.Bool() {
			order.ReturnedAt = gofakeit.DateRange(order.IssuedAt, order.IssuedAt.Add(usecases.TimeForReturn))
		}
	}

	fmt.Printf("Creating order: %+v\n", order)

	err := repo.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	fmt.Println("Order created")

	return nil
}
