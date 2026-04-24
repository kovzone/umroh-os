package service

// BL-CAT-010 / BL-CAT-011 — BulkImportPackages + BulkUpdatePackages.
//
// Both operations process rows sequentially. A failure in one row never aborts
// the others (per-row transactional isolation). Each row result records either
// a package_id (on success) or an error message (on failure).

import (
	"context"
	"fmt"
)

// ---------------------------------------------------------------------------
// Input / output types for bulk operations
// ---------------------------------------------------------------------------

// BulkImportRowInput is the input for a single row in a bulk import.
type BulkImportRowInput struct {
	Name          string
	Kind          string
	Description   string
	CoverPhotoUrl string
	Highlights    []string
	AddonIDs      []string
	HotelIDs      []string
	ItineraryID   string
	AirlineID     string
	MuthawwifID   string
	Status        string
}

// BulkImportRowResult is the per-row result for a bulk import.
type BulkImportRowResult struct {
	Index     int
	PackageID string // populated on success
	Error     string // populated on failure
}

// BulkImportPackagesParams is the input for BulkImportPackages.
type BulkImportPackagesParams struct {
	UserID   string
	BranchID string
	Rows     []BulkImportRowInput
}

// BulkImportPackagesResult is the response for BulkImportPackages.
type BulkImportPackagesResult struct {
	Results    []BulkImportRowResult
	TotalRows  int
	Successful int
	Failed     int
}

// BulkUpdateRowInput is the input for a single row in a bulk update.
type BulkUpdateRowInput struct {
	ID          string
	Name        string
	Description string
	Status      string
	Highlights  []string
	AddonIDs    []string
	HotelIDs    []string
}

// BulkUpdateRowResult is the per-row result for a bulk update.
type BulkUpdateRowResult struct {
	Index     int
	PackageID string // populated on success (same as input ID)
	Error     string // populated on failure
}

// BulkUpdatePackagesParams is the input for BulkUpdatePackages.
type BulkUpdatePackagesParams struct {
	UserID   string
	BranchID string
	Rows     []BulkUpdateRowInput
}

// BulkUpdatePackagesResult is the response for BulkUpdatePackages.
type BulkUpdatePackagesResult struct {
	Results    []BulkUpdateRowResult
	TotalRows  int
	Successful int
	Failed     int
}

// ---------------------------------------------------------------------------
// BulkImportPackages (BL-CAT-010)
// ---------------------------------------------------------------------------

// BulkImportPackages processes each row by calling CreatePackage.
// One row failure never aborts others — per-row transactional isolation.
func (s *Service) BulkImportPackages(ctx context.Context, params *BulkImportPackagesParams) (*BulkImportPackagesResult, error) {
	result := &BulkImportPackagesResult{
		TotalRows: len(params.Rows),
		Results:   make([]BulkImportRowResult, 0, len(params.Rows)),
	}

	for i, row := range params.Rows {
		detail, err := s.CreatePackage(ctx, &CreatePackageParams{
			UserID:        params.UserID,
			BranchID:      params.BranchID,
			Kind:          row.Kind,
			Name:          row.Name,
			Description:   row.Description,
			CoverPhotoUrl: row.CoverPhotoUrl,
			Highlights:    row.Highlights,
			AddonIDs:      row.AddonIDs,
			HotelIDs:      row.HotelIDs,
			ItineraryID:   row.ItineraryID,
			AirlineID:     row.AirlineID,
			MuthawwifID:   row.MuthawwifID,
			Status:        row.Status,
		})
		if err != nil {
			result.Results = append(result.Results, BulkImportRowResult{
				Index: i,
				Error: fmt.Sprintf("row %d: %s", i, err.Error()),
			})
			result.Failed++
			continue
		}
		result.Results = append(result.Results, BulkImportRowResult{
			Index:     i,
			PackageID: detail.ID,
		})
		result.Successful++
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// BulkUpdatePackages (BL-CAT-011)
// ---------------------------------------------------------------------------

// BulkUpdatePackages processes each row by calling UpdatePackage.
// One row failure never aborts others — per-row transactional isolation.
func (s *Service) BulkUpdatePackages(ctx context.Context, params *BulkUpdatePackagesParams) (*BulkUpdatePackagesResult, error) {
	result := &BulkUpdatePackagesResult{
		TotalRows: len(params.Rows),
		Results:   make([]BulkUpdateRowResult, 0, len(params.Rows)),
	}

	for i, row := range params.Rows {
		updateParams := &UpdatePackageParams{
			UserID:      params.UserID,
			BranchID:    params.BranchID,
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Status:      row.Status,
		}
		// Nil slices mean "no change"; non-nil mean "replace".
		if row.Highlights != nil {
			updateParams.Highlights = row.Highlights
		}
		if row.AddonIDs != nil {
			updateParams.AddonIDs = row.AddonIDs
		}
		if row.HotelIDs != nil {
			updateParams.HotelIDs = row.HotelIDs
		}

		_, err := s.UpdatePackage(ctx, updateParams)
		if err != nil {
			result.Results = append(result.Results, BulkUpdateRowResult{
				Index:     i,
				PackageID: row.ID,
				Error:     fmt.Sprintf("row %d (id=%s): %s", i, row.ID, err.Error()),
			})
			result.Failed++
			continue
		}
		result.Results = append(result.Results, BulkUpdateRowResult{
			Index:     i,
			PackageID: row.ID,
		})
		result.Successful++
	}

	return result, nil
}
