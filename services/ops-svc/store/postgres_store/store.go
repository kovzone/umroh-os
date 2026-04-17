// Package postgres_store provides database access layer with transaction support.
package postgres_store

import (
	"context"

	"ops-svc/store/postgres_store/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type IStore interface {
	sqlc.Querier

	WithTx(ctx context.Context, args *WithTxArgs) (*WithTxData, error)
	WithTxOptions(ctx context.Context, args *WithTxOptionsArgs) (*WithTxOptionsData, error)
	ExecRawSQL(ctx context.Context, sql string) error
}

type Store struct {
	*sqlc.Queries

	logger *zerolog.Logger
	tracer trace.Tracer

	pool *pgxpool.Pool
}

func NewStore(logger *zerolog.Logger, tracer trace.Tracer, pool *pgxpool.Pool) IStore {
	// Wrap the pool with logging to log all SQL queries with parameters substituted
	loggingDB := NewLoggingDBTX(pool, logger, tracer)

	return &Store{
		Queries: sqlc.New(loggingDB),

		logger: logger,
		tracer: tracer,

		pool: pool,
	}
}

// ExecRawSQL executes a raw SQL statement (for DDL operations like CREATE TABLE, DROP SEQUENCE, etc.)
func (store *Store) ExecRawSQL(ctx context.Context, sql string) error {
	_, err := store.pool.Exec(ctx, sql)

	return err
}
