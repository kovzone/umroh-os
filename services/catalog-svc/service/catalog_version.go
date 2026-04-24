package service

// BL-CAT-013 — Package versioning.
//
// GetPackageVersion returns a content-hash of the package's key mutable fields
// so callers can detect changes without fetching the full package. The version
// is computed on-read (not stored) as:
//
//	MD5( package_id + "|" + name + "|" + status + "|" + updated_at_unix )
//
// as a lowercase hex string.

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"

	"catalog-svc/store/postgres_store"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetPackageVersionParams is the input for GetPackageVersion.
type GetPackageVersionParams struct {
	PackageID string
}

// PackageVersionResult carries the computed version.
type PackageVersionResult struct {
	PackageID string
	Version   string
	UpdatedAt string
}

// GetPackageVersion computes and returns the version hash for a package.
func (s *Service) GetPackageVersion(ctx context.Context, params *GetPackageVersionParams) (*PackageVersionResult, error) {
	const op = "service.Service.GetPackageVersion"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.package_id", params.PackageID),
	)
	logger.Info().Str("op", op).Str("package_id", params.PackageID).Msg("")

	if params.PackageID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("package_id is required"))
	}

	row, err := s.store.GetPackageByIDForStaff(ctx, params.PackageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("package %q not found", params.PackageID))
		}
		return nil, fmt.Errorf("get package: %w", postgres_store.WrapDBError(err))
	}

	updatedAtUnix := int64(0)
	updatedAtStr := ""
	if row.UpdatedAt.Valid {
		updatedAtUnix = row.UpdatedAt.Time.UTC().Unix()
		updatedAtStr = row.UpdatedAt.Time.UTC().Format("2006-01-02T15:04:05Z")
	}

	version := computePackageVersion(row.ID, row.Name, string(row.Status), updatedAtUnix)

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.version", version))
	return &PackageVersionResult{
		PackageID: row.ID,
		Version:   version,
		UpdatedAt: updatedAtStr,
	}, nil
}

// computePackageVersion computes MD5(id + "|" + name + "|" + status + "|" + updatedAtUnix)
// as a lowercase hex string.
func computePackageVersion(id, name, status string, updatedAtUnix int64) string {
	raw := id + "|" + name + "|" + status + "|" + strconv.FormatInt(updatedAtUnix, 10)
	sum := md5.Sum([]byte(raw)) //nolint:gosec // not used for security
	return fmt.Sprintf("%x", sum)
}
