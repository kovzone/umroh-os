// Package postgres_store provides database access layer with transaction support.
package postgres_store

import (
	"context"
	"fmt"
	"time"

	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Constants for transaction retries
const (
	maxTxRetries   = 3
	retryBaseDelay = time.Millisecond * 10
	retryMaxDelay  = time.Millisecond * 100
)

// TxOptions represents available transaction options.
type TxOptions struct {
	Isolation  pgx.TxIsoLevel
	AccessMode pgx.TxAccessMode
	Deferrable bool
	ReadOnly   bool
}

// DefaultTxOptions returns default transaction options (Serializable, ReadWrite, NotDeferrable).
func DefaultTxOptions() TxOptions {
	return TxOptions{
		Isolation:  pgx.Serializable,
		AccessMode: pgx.ReadWrite,
		Deferrable: false,
		ReadOnly:   false,
	}
}

// ReadOnlyTxOptions returns options optimized for read-only transactions.
func ReadOnlyTxOptions() TxOptions {
	return TxOptions{
		Isolation:  pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
		Deferrable: true,
		ReadOnly:   true,
	}
}

type WithTxOptionsArgs struct {
	Opts TxOptions
	Fn   func(*sqlc.Queries) error
}

type WithTxOptionsData struct{}

// WithTxOptions executes a function within a database transaction with custom options.
// Retries up to maxTxRetries times on serialization_failure (40001) or deadlock_detected (40P01).
func (store *Store) WithTxOptions(ctx context.Context, args *WithTxOptionsArgs) (*WithTxOptionsData, error) {
	const op = "postgres_store.Store.WithTxOptions"

	// Start span
	ctx, span := store.tracer.Start(ctx, op)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.args", fmt.Sprintf("%+v", args)),
	)

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, store.logger)
	logger.Info().
		Str("op", op).
		Str("args", fmt.Sprintf("%+v", args)).
		Msg("")

	// Initialize data
	data := &WithTxOptionsData{}

	deferrable := pgx.NotDeferrable
	if args.Opts.Deferrable {
		deferrable = pgx.Deferrable
	}

	txOptions := &pgx.TxOptions{
		IsoLevel:       args.Opts.Isolation,
		AccessMode:     args.Opts.AccessMode,
		DeferrableMode: deferrable,
		BeginQuery:     "",
		CommitQuery:    "",
	}

	var lastErr error
	for attempt := 0; attempt < maxTxRetries; attempt++ {
		transaction, err := store.pool.BeginTx(ctx, *txOptions)
		if err != nil {
			err = fmt.Errorf("%w: %w", ErrPostgresBeginTransactionError, err)
			logger.Error().Err(err).Msg("")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		// Wrap transaction with logging to log all SQL queries with parameters substituted
		loggingTx := NewLoggingDBTX(transaction, store.logger, store.tracer)
		q := sqlc.New(loggingTx)

		err = args.Fn(q)
		if err != nil {
			return nil, handleRollbackError(ctx, transaction, err)
		}

		err = transaction.Commit(ctx)
		if err != nil {
			_ = transaction.Rollback(ctx) // best-effort; tx may already be aborted
			lastErr = err

			if isRetryableTransactionError(err) && attempt < maxTxRetries-1 {
				backoff := min(retryBaseDelay*time.Duration(1<<attempt), retryMaxDelay)

				logger.Warn().
					Int("attempt", attempt+1).
					Dur("retry_in", backoff).
					Msg("retrying transaction after serialization or deadlock")
				time.Sleep(backoff)

				continue
			}

			return nil, handleCommitError(err)
		}

		// Set span attributes and status
		span.SetAttributes(attribute.String("output.data", fmt.Sprintf("%+v", data)))
		span.SetStatus(codes.Ok, "success")

		return data, nil
	}

	// Should not reach here; return last error if we do
	return nil, handleCommitError(lastErr)
}

type WithTxArgs struct {
	Fn func(*sqlc.Queries) error
}

type WithTxData struct{}

// WithTx executes a function within a database transaction with default options.
func (store *Store) WithTx(ctx context.Context, args *WithTxArgs) (*WithTxData, error) {
	const op = "postgres_store.Store.WithTx"

	// Start span
	ctx, span := store.tracer.Start(ctx, op)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.args", fmt.Sprintf("%+v", args)),
	)

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, store.logger)
	logger.Info().
		Str("op", op).
		Str("args", fmt.Sprintf("%+v", args)).
		Msg("")

	// Initialize data
	data := &WithTxData{}

	_, err := store.WithTxOptions(ctx, &WithTxOptionsArgs{
		Opts: DefaultTxOptions(),
		Fn:   args.Fn,
	})
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrFailedToExecuteTransaction, err)

		logger.Error().Err(err).Msg("")

		// Set span error and status
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.data", fmt.Sprintf("%+v", data)),
	)
	span.SetStatus(codes.Ok, "success")

	return data, nil
}
