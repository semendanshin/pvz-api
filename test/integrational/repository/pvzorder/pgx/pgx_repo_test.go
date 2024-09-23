package pgx

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"homework/internal/abstractions"
	"testing"
	"time"

	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/pvzorder/pgx"
	"homework/internal/infrastructure/repositories/utils/pgx/txmanager"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

//go:embed migrations/*.sql
var migrations embed.FS

const (
	dbName = "postgres"
	dbUser
	dbPassword
)

func startContainer(ctx context.Context) (string, string, func(), error) {
	// Define the container request
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:15",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Get the host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to get container port: %w", err)
	}

	return host, port.Port(), func() {
		postgresContainer.Terminate(ctx)
	}, nil
}

func applyMigrations(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func setupDatabase(ctx context.Context, t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	host, port, tearDown, err := startContainer(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to start container: %w", err))
	}

	dbUrl := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=postgres sslmode=disable", host, port)

	db, err := sql.Open(
		"postgres",
		dbUrl,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to open database: %w", err))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(time.Minute * 3)

	err = db.Ping()
	if err != nil {
		t.Fatal(fmt.Errorf("failed to ping database: %w", err))
	}

	err = applyMigrations(ctx, db)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to apply migrations: %w", err))
	}

	_ = db.Close()

	// Connect to the database
	pgxPool, err := pgxpool.New(
		ctx,
		dbUrl,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	return pgxPool, func() {
		defer tearDown()
		defer pgxPool.Close()
	}
}

func setupSuite(ctx context.Context, t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	pgxPool, tearDown := setupDatabase(ctx, t)

	return pgxPool, func() {
		defer tearDown()
	}
}

func TestPGXRepository_CreateOrder(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxPool, tearDown := setupSuite(ctx, t)
	defer tearDown()

	manager := txmanager.NewPGXTXManager(pgxPool)
	repo := pgx.NewPgxPvzOrderFacade(manager)

	order := domain.NewPVZOrder(
		"100",
		"1",
		"1",
		1000,
		1000,
		24*time.Hour,
		domain.PackagingTypeBox,
		false,
	)

	err := repo.CreateOrder(ctx, order)
	assert.NoError(t, err)

	actual, err := repo.GetOrder(ctx, "100")
	assert.NoError(t, err)
	assert.Equal(t, order.ReceivedAt.UnixMilli(), actual.ReceivedAt.UnixMilli())
	assert.Equal(t, order, actual)
}

func TestPGXRepository_DeleteOrder(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxPool, tearDown := setupSuite(ctx, t)
	defer tearDown()

	manager := txmanager.NewPGXTXManager(pgxPool)
	repo := pgx.NewPgxPvzOrderFacade(manager)

	err := repo.DeleteOrder(ctx, "1")
	assert.NoError(t, err)

	_, err = repo.GetOrder(ctx, "1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPGXRepository_SetOrderIssued(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxPool, tearDown := setupSuite(ctx, t)
	defer tearDown()

	manager := txmanager.NewPGXTXManager(pgxPool)
	repo := pgx.NewPgxPvzOrderFacade(manager)

	err := repo.SetOrderIssued(ctx, "1")
	assert.NoError(t, err)

	order, err := repo.GetOrder(ctx, "1")
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, order.IssuedAt)
}

func TestPGXRepository_SetOrderReturned(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxPool, tearDown := setupSuite(ctx, t)
	defer tearDown()

	manager := txmanager.NewPGXTXManager(pgxPool)
	repo := pgx.NewPgxPvzOrderFacade(manager)

	err := repo.SetOrderReturned(ctx, "1")
	assert.NoError(t, err)

	order, err := repo.GetOrder(ctx, "1")
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, order.ReturnedAt)
}

func TestPGXRepository_GetOrders(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxPool, tearDown := setupSuite(ctx, t)
	defer tearDown()

	manager := txmanager.NewPGXTXManager(pgxPool)
	repo := pgx.NewPgxPvzOrderFacade(manager)

	type args struct {
		userID string
		opts   []abstractions.GetOrdersOptFunc
	}

	tests := []struct {
		name    string
		args    args
		want    func([]domain.PVZOrder) bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID: "1",
				opts:   make([]abstractions.GetOrdersOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 3)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success with PVZID",
			args: args{
				userID: "1",
				opts: []abstractions.GetOrdersOptFunc{
					abstractions.WithPVZID("2"),
				},
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 1) && assert.Equal(t, "3", orders[0].OrderID)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success with pagination",
			args: args{
				userID: "1",
				opts: []abstractions.GetOrdersOptFunc{
					abstractions.WithCursorID("3"),
					abstractions.WithLimit(2),
				},
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 2) && assert.Equal(t, "3", orders[0].OrderID) && assert.Equal(t, "2", orders[1].OrderID)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success empty",
			args: args{
				userID: "3",
				opts:   make([]abstractions.GetOrdersOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 0)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := repo.GetOrders(ctx, tt.args.userID, tt.args.opts...)
			tt.want(orders)
			tt.wantErr(t, err)
		})
	}
}

func TestPGXRepository_GetOrder(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxPool, tearDown := setupSuite(ctx, t)
	defer tearDown()

	manager := txmanager.NewPGXTXManager(pgxPool)
	repo := pgx.NewPgxPvzOrderFacade(manager)

	type args struct {
		orderID string
	}

	tests := []struct {
		name    string
		args    args
		want    domain.PVZOrder
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				orderID: "1",
			},
			want: domain.PVZOrder{
				OrderID:        "1",
				PVZID:          "1",
				RecipientID:    "1",
				ReceivedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				StorageTime:    24 * time.Hour,
				IssuedAt:       time.Time{},
				ReturnedAt:     time.Time{},
				Weight:         1000,
				Cost:           1000,
				AdditionalFilm: false,
				Packaging:      domain.PackagingTypeBox,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Order not found",
			args: args{
				orderID: "100",
			},
			want: domain.PVZOrder{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && errors.Is(err, domain.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetOrder(ctx, tt.args.orderID)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestPGXRepository_GetReturns(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxPool, tearDown := setupSuite(ctx, t)
	defer tearDown()

	manager := txmanager.NewPGXTXManager(pgxPool)
	repo := pgx.NewPgxPvzOrderFacade(manager)

	type args struct {
		opts []abstractions.PagePaginationOptFunc
	}

	tests := []struct {
		name    string
		args    args
		want    func([]domain.PVZOrder) bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				opts: make([]abstractions.PagePaginationOptFunc, 0),
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 3)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Success with pagination",
			args: args{
				opts: []abstractions.PagePaginationOptFunc{
					abstractions.WithPage(1),
					abstractions.WithPageSize(2),
				},
			},
			want: func(orders []domain.PVZOrder) bool {
				return assert.Len(t, orders, 1) && assert.Equal(t, "7", orders[0].OrderID)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := repo.GetReturns(ctx, tt.args.opts...)
			tt.want(orders)
			tt.wantErr(t, err)
		})
	}
}
