package txmanager

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	key   string
	inner func(context.Context) error
)

const engineKey key = "engine"

type QueryEngine interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PGXTXManager struct {
	pool *pgxpool.Pool
}

func NewPGXTXManager(pool *pgxpool.Pool) *PGXTXManager {
	return &PGXTXManager{pool: pool}
}

func (p *PGXTXManager) RunReadUncommittedTransaction(ctx context.Context, f inner) error {
	return p.runTransaction(ctx, pgx.ReadUncommitted, f)
}

func (p *PGXTXManager) RunReadCommittedTransaction(ctx context.Context, f inner) error {
	return p.runTransaction(ctx, pgx.ReadCommitted, f)
}

func (p *PGXTXManager) RunRepeatableReadTransaction(ctx context.Context, f inner) error {
	return p.runTransaction(ctx, pgx.RepeatableRead, f)
}

func (p *PGXTXManager) RunSerializableTransaction(ctx context.Context, f inner) error {
	return p.runTransaction(ctx, pgx.Serializable, f)
}

func (p *PGXTXManager) runTransaction(ctx context.Context, level pgx.TxIsoLevel, f inner) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: level})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ctx = context.WithValue(ctx, engineKey, tx)

	err = f(ctx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PGXTXManager) GetQueryEngine(ctx context.Context) QueryEngine {
	engine, ok := ctx.Value(engineKey).(QueryEngine)
	if !ok {
		return p.pool
	}

	return engine
}
