package main

import (
	"context"
	"fmt"
	"homework/internal/infrastructure/clients/queue/dummy"
	"homework/internal/infrastructure/repositories/events/pgx"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homework/internal/usecases"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
		return fmt.Errorf("error loading .env file %w", err)
	}

	postgresURL := loadPostgresURL()

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

	txm := txmanager.NewPGXTXManager(pool)
	repo := pgx.NewEventsRepository(txm)

	client := dummy.NewDummyEventsSender()

	worker := usecases.NewEventsSender(repo, client, 10)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		<-stop

		worker.Stop()
	}()

	if err := worker.Run(ctx, 5*time.Second); err != nil {
		return fmt.Errorf("error running worker: %w", err)
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
