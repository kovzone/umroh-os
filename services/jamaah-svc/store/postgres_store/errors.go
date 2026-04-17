// Package postgres_store provides database access layer with transaction support.
package postgres_store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrRollbackFailed indicates that a transaction rollback failed.
	ErrRollbackFailed = errors.New("rollback failed")

	// ErrPostgresRollbackError indicates a PostgreSQL-specific error during rollback.
	ErrPostgresRollbackError = errors.New("postgres rollback error")

	// ErrPostgresCommitError indicates a PostgreSQL-specific error during commit.
	ErrPostgresCommitError = errors.New("postgres commit error")

	// ErrPostgresBeginTransactionError indicates a PostgreSQL-specific error during begin transaction.
	ErrPostgresBeginTransactionError = errors.New("postgres begin transaction error")

	// ErrFailedToExecuteTransaction indicates that a transaction execution failed.
	ErrFailedToExecuteTransaction = errors.New("failed to execute transaction")
)

// handleRollbackError processes rollback errors and wraps them appropriately.
func handleRollbackError(ctx context.Context, transaction pgx.Tx, originalErr error) error {
	rbErr := transaction.Rollback(ctx)
	if rbErr == nil {
		// Rollback succeeded, return original error
		return originalErr
	}

	if errors.Is(rbErr, pgx.ErrTxClosed) {
		// Transaction already closed, return original error
		return originalErr
	}

	// Handle PostgreSQL-specific rollback errors
	var pgErr *pgconn.PgError
	if errors.As(rbErr, &pgErr) {
		// PostgreSQL error during rollback - this is a serious server issue
		return fmt.Errorf(
			"%w: Code=%s, Message=%s, RollbackError=%w, OriginalError=%w",
			ErrPostgresRollbackError,
			pgErr.Code,
			pgErr.Message,
			rbErr,
			originalErr,
		)
	}

	// Generic rollback error
	return fmt.Errorf("%w: RollbackError=%w, OriginalError=%w", ErrRollbackFailed, rbErr, originalErr)
}

// isRetryableTransactionError returns true if the error is a PostgreSQL serialization_failure (40001)
// or deadlock_detected (40P01), indicating the transaction can be retried.
func isRetryableTransactionError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40001" || pgErr.Code == "40P01"
	}

	return false
}

// handleCommitError processes commit errors and wraps them appropriately.
func handleCommitError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		// Generic commit failure
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Check for concurrency/serialization errors
	switch pgErr.Code {
	case "40001": // serialization_failure
		return fmt.Errorf("transaction serialization failure: %w", err)
	case "40P01": // deadlock_detected
		return fmt.Errorf("transaction deadlock detected: %w", err)
	default:
		// Other PostgreSQL errors during commit
		return fmt.Errorf("%w: Code=%s, Message=%s", ErrPostgresCommitError, pgErr.Code, pgErr.Message)
	}
}
