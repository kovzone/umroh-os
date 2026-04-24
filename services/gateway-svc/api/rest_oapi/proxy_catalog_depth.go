// proxy_catalog_depth.go — gateway REST handlers for catalog Wave 3 depth RPCs.
// BL-CAT-010: BulkImportPackages  POST /v1/admin/catalog/packages/bulk-import
// BL-CAT-011: BulkUpdatePackages  POST /v1/admin/catalog/packages/bulk-update
// BL-CAT-013: GetPackageVersion   GET  /v1/admin/catalog/packages/:id/version
// BL-BOOK-007: GetSeatsByChannel  GET  /v1/bookings/departures/:id/seats-by-channel

package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Request body types
// ---------------------------------------------------------------------------

// BulkImportRowBody is one row in the BulkImportPackages request body.
type BulkImportRowBody struct {
	Name          string   `json:"name"`
	Kind          string   `json:"kind"`
	Description   string   `json:"description,omitempty"`
	CoverPhotoURL string   `json:"cover_photo_url,omitempty"`
	Highlights    []string `json:"highlights,omitempty"`
	AddonIDs      []string `json:"addon_ids,omitempty"`
	HotelIDs      []string `json:"hotel_ids,omitempty"`
	ItineraryID   string   `json:"itinerary_id,omitempty"`
	AirlineID     string   `json:"airline_id,omitempty"`
	MuthawwifID   string   `json:"muthawwif_id,omitempty"`
	Status        string   `json:"status,omitempty"`
}

// BulkImportPackagesBody is the request body for BulkImportPackages.
type BulkImportPackagesBody struct {
	UserID   string              `json:"user_id,omitempty"`
	BranchID string              `json:"branch_id,omitempty"`
	Rows     []BulkImportRowBody `json:"rows"`
}

// BulkUpdateRowBody is one row in the BulkUpdatePackages request body.
type BulkUpdateRowBody struct {
	ID          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status,omitempty"`
	Highlights  []string `json:"highlights,omitempty"`
	AddonIDs    []string `json:"addon_ids,omitempty"`
	HotelIDs    []string `json:"hotel_ids,omitempty"`
}

// BulkUpdatePackagesBody is the request body for BulkUpdatePackages.
type BulkUpdatePackagesBody struct {
	UserID   string              `json:"user_id,omitempty"`
	BranchID string              `json:"branch_id,omitempty"`
	Rows     []BulkUpdateRowBody `json:"rows"`
}

// ---------------------------------------------------------------------------
// BL-CAT-010: BulkImportPackages
// ---------------------------------------------------------------------------

func (s *Server) BulkImportPackages(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.BulkImportPackages"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body BulkImportPackagesBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	rows := make([]catalog_grpc_adapter.BulkImportRowInput, 0, len(body.Rows))
	for _, r := range body.Rows {
		rows = append(rows, catalog_grpc_adapter.BulkImportRowInput{
			Name:          r.Name,
			Kind:          r.Kind,
			Description:   r.Description,
			CoverPhotoURL: r.CoverPhotoURL,
			Highlights:    r.Highlights,
			AddonIDs:      r.AddonIDs,
			HotelIDs:      r.HotelIDs,
			ItineraryID:   r.ItineraryID,
			AirlineID:     r.AirlineID,
			MuthawwifID:   r.MuthawwifID,
			Status:        r.Status,
		})
	}

	result, err := s.svc.BulkImportPackages(ctx, &catalog_grpc_adapter.BulkImportPackagesParams{
		UserID:   body.UserID,
		BranchID: body.BranchID,
		Rows:     rows,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type rowJSON struct {
		Index     int    `json:"index"`
		PackageID string `json:"package_id,omitempty"`
		Error     string `json:"error,omitempty"`
	}
	data := make([]rowJSON, 0, len(result.Results))
	for _, r := range result.Results {
		data = append(data, rowJSON{Index: r.Index, PackageID: r.PackageID, Error: r.Error})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"results":    data,
			"total_rows": result.TotalRows,
			"successful": result.Successful,
			"failed":     result.Failed,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-CAT-011: BulkUpdatePackages
// ---------------------------------------------------------------------------

func (s *Server) BulkUpdatePackages(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.BulkUpdatePackages"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body BulkUpdatePackagesBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	rows := make([]catalog_grpc_adapter.BulkUpdateRowInput, 0, len(body.Rows))
	for _, r := range body.Rows {
		rows = append(rows, catalog_grpc_adapter.BulkUpdateRowInput{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Status:      r.Status,
			Highlights:  r.Highlights,
			AddonIDs:    r.AddonIDs,
			HotelIDs:    r.HotelIDs,
		})
	}

	result, err := s.svc.BulkUpdatePackages(ctx, &catalog_grpc_adapter.BulkUpdatePackagesParams{
		UserID:   body.UserID,
		BranchID: body.BranchID,
		Rows:     rows,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type rowJSON struct {
		Index     int    `json:"index"`
		PackageID string `json:"package_id,omitempty"`
		Error     string `json:"error,omitempty"`
	}
	data := make([]rowJSON, 0, len(result.Results))
	for _, r := range result.Results {
		data = append(data, rowJSON{Index: r.Index, PackageID: r.PackageID, Error: r.Error})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"results":    data,
			"total_rows": result.TotalRows,
			"successful": result.Successful,
			"failed":     result.Failed,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-CAT-013: GetPackageVersion
// ---------------------------------------------------------------------------

func (s *Server) GetPackageVersion(c *fiber.Ctx, packageID string) error {
	const op = "rest_oapi.Server.GetPackageVersion"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("package_id", packageID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("package_id", packageID).Msg("")

	result, err := s.svc.GetPackageVersion(ctx, packageID)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"package_id": result.PackageID,
			"version":    result.Version,
			"updated_at": result.UpdatedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-BOOK-007: GetSeatsByChannel
// ---------------------------------------------------------------------------

func (s *Server) GetSeatsByChannel(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetSeatsByChannel"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	result, err := s.svc.GetSeatsByChannel(ctx, departureID)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type channelJSON struct {
		Channel      string `json:"channel"`
		Seats        int    `json:"seats"`
		BookingCount int    `json:"booking_count"`
	}
	data := make([]channelJSON, 0, len(result.ByChannel))
	for _, r := range result.ByChannel {
		data = append(data, channelJSON{
			Channel:      r.Channel,
			Seats:        r.Seats,
			BookingCount: r.BookingCount,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"departure_id": result.DepartureID,
			"by_channel":   data,
		},
	})
}
