// Package postgres_store provides database access layer with transaction support.
package postgres_store

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"catalog-svc/util/apperrors"
)

// WrapDBError maps PostgreSQL/pgx errors to domain sentinels so the API can return
// the correct HTTP status. Use this inside WithTx/WithTxOptions callbacks when
// a query returns an error.
//
// Mapping:
//   - pgx.ErrNoRows → ErrNotFound
//   - 23505 (unique_violation) → ErrConflict
//   - 23503 (foreign_key_violation) → ErrConflict
//   - Other PG errors → ErrInternal
//   - Unknown → ErrInternal
func WrapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.Join(apperrors.ErrNotFound, err)
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return errors.Join(apperrors.ErrConflict, err)
		case "23503": // foreign_key_violation
			return errors.Join(apperrors.ErrConflict, err)
		default:
			return errors.Join(apperrors.ErrInternal, err)
		}
	}
	return errors.Join(apperrors.ErrInternal, err)
}
