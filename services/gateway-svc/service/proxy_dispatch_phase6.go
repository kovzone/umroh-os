// proxy_dispatch_phase6.go — gateway service dispatch for Phase 6 new routes.
//
// Covers:
//   - Catalog master data CRUD (Wave 1A): hotels, airlines, muthawwif, addons,
//     departure pricing
//   - Finance depth RPCs (Wave 1B): RecognizeRevenue, GetPLReport, GetBalanceSheet
//   - IAM admin (Wave 1C): user management, role management, permissions,
//     role assignment
//   - Jamaah manifest (Wave 1A): GetDepartureManifest
//
// Each method is a thin delegation to the appropriate adapter.
// No business logic lives here; all logic lives in the backend services.
package service

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/adapter/jamaah_grpc_adapter"
)

// ---------------------------------------------------------------------------
// Catalog masters — Hotels
// ---------------------------------------------------------------------------

func (s *Service) ListHotels(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListHotelsResult, error) {
	return s.adapters.catalogGrpc.ListHotels(ctx, params)
}

func (s *Service) CreateHotel(ctx context.Context, params *catalog_grpc_adapter.CreateHotelParams) (*catalog_grpc_adapter.HotelResult, error) {
	return s.adapters.catalogGrpc.CreateHotel(ctx, params)
}

func (s *Service) UpdateHotel(ctx context.Context, params *catalog_grpc_adapter.UpdateHotelParams) (*catalog_grpc_adapter.HotelResult, error) {
	return s.adapters.catalogGrpc.UpdateHotel(ctx, params)
}

func (s *Service) DeleteHotel(ctx context.Context, id string) error {
	return s.adapters.catalogGrpc.DeleteHotel(ctx, "", id)
}

// ---------------------------------------------------------------------------
// Catalog masters — Airlines
// ---------------------------------------------------------------------------

func (s *Service) ListAirlines(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListAirlinesResult, error) {
	return s.adapters.catalogGrpc.ListAirlines(ctx, params)
}

func (s *Service) CreateAirline(ctx context.Context, params *catalog_grpc_adapter.CreateAirlineParams) (*catalog_grpc_adapter.AirlineResult, error) {
	return s.adapters.catalogGrpc.CreateAirline(ctx, params)
}

func (s *Service) UpdateAirline(ctx context.Context, params *catalog_grpc_adapter.UpdateAirlineParams) (*catalog_grpc_adapter.AirlineResult, error) {
	return s.adapters.catalogGrpc.UpdateAirline(ctx, params)
}

func (s *Service) DeleteAirline(ctx context.Context, id string) error {
	return s.adapters.catalogGrpc.DeleteAirline(ctx, "", id)
}

// ---------------------------------------------------------------------------
// Catalog masters — Muthawwif
// ---------------------------------------------------------------------------

func (s *Service) ListMuthawwif(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListMuthawwifResult, error) {
	return s.adapters.catalogGrpc.ListMuthawwif(ctx, params)
}

func (s *Service) CreateMuthawwif(ctx context.Context, params *catalog_grpc_adapter.CreateMuthawwifParams) (*catalog_grpc_adapter.MuthawwifResult, error) {
	return s.adapters.catalogGrpc.CreateMuthawwif(ctx, params)
}

func (s *Service) UpdateMuthawwif(ctx context.Context, params *catalog_grpc_adapter.UpdateMuthawwifParams) (*catalog_grpc_adapter.MuthawwifResult, error) {
	return s.adapters.catalogGrpc.UpdateMuthawwif(ctx, params)
}

func (s *Service) DeleteMuthawwif(ctx context.Context, id string) error {
	return s.adapters.catalogGrpc.DeleteMuthawwif(ctx, "", id)
}

// ---------------------------------------------------------------------------
// Catalog masters — Addons
// ---------------------------------------------------------------------------

func (s *Service) ListAddons(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListAddonsResult, error) {
	return s.adapters.catalogGrpc.ListAddons(ctx, params)
}

func (s *Service) CreateAddon(ctx context.Context, params *catalog_grpc_adapter.CreateAddonParams) (*catalog_grpc_adapter.AddonResult, error) {
	return s.adapters.catalogGrpc.CreateAddon(ctx, params)
}

func (s *Service) UpdateAddon(ctx context.Context, params *catalog_grpc_adapter.UpdateAddonParams) (*catalog_grpc_adapter.AddonResult, error) {
	return s.adapters.catalogGrpc.UpdateAddon(ctx, params)
}

func (s *Service) DeleteAddon(ctx context.Context, id string) error {
	return s.adapters.catalogGrpc.DeleteAddon(ctx, "", id)
}

// ---------------------------------------------------------------------------
// Catalog masters — Departure pricing
// ---------------------------------------------------------------------------

func (s *Service) SetDeparturePricing(ctx context.Context, departureID string, pricings []*catalog_grpc_adapter.PricingUpsertInput) (*catalog_grpc_adapter.PricingResult, error) {
	return s.adapters.catalogGrpc.SetDeparturePricing(ctx, "", departureID, pricings)
}

func (s *Service) GetDeparturePricing(ctx context.Context, departureID string) (*catalog_grpc_adapter.PricingResult, error) {
	return s.adapters.catalogGrpc.GetDeparturePricing(ctx, "", departureID)
}

// ---------------------------------------------------------------------------
// Finance depth
// ---------------------------------------------------------------------------

func (s *Service) RecognizeRevenue(ctx context.Context, departureID string, totalAmountIdr int64) (*finance_grpc_adapter.RecognizeRevenueResult, error) {
	return s.adapters.financeGrpc.RecognizeRevenue(ctx, departureID, totalAmountIdr)
}

func (s *Service) GetPLReport(ctx context.Context, from, to string) (*finance_grpc_adapter.PLReportResult, error) {
	return s.adapters.financeGrpc.GetPLReport(ctx, from, to)
}

func (s *Service) GetBalanceSheet(ctx context.Context, asOf string) (*finance_grpc_adapter.BalanceSheetResult, error) {
	return s.adapters.financeGrpc.GetBalanceSheet(ctx, asOf)
}

// ---------------------------------------------------------------------------
// IAM admin — Users
// ---------------------------------------------------------------------------

func (s *Service) ListUsers(ctx context.Context, params *iam_grpc_adapter.ListUsersParams) (*iam_grpc_adapter.ListUsersResult, error) {
	return s.adapters.iamGrpc.ListUsers(ctx, params)
}

func (s *Service) CreateUser(ctx context.Context, params *iam_grpc_adapter.CreateUserParams) (*iam_grpc_adapter.AdminUserResult, error) {
	return s.adapters.iamGrpc.CreateUser(ctx, params)
}

func (s *Service) UpdateUser(ctx context.Context, id string, params *iam_grpc_adapter.UpdateUserParams) (*iam_grpc_adapter.AdminUserResult, error) {
	return s.adapters.iamGrpc.UpdateUser(ctx, id, params)
}

func (s *Service) GetUser(ctx context.Context, id string) (*iam_grpc_adapter.AdminUserResult, error) {
	return s.adapters.iamGrpc.GetUser(ctx, id)
}

func (s *Service) ResetUserPassword(ctx context.Context, id, newPassword string) error {
	return s.adapters.iamGrpc.ResetUserPassword(ctx, id, newPassword)
}

// ---------------------------------------------------------------------------
// IAM admin — Roles
// ---------------------------------------------------------------------------

func (s *Service) ListRoles(ctx context.Context) (*iam_grpc_adapter.ListRolesResult, error) {
	return s.adapters.iamGrpc.ListRoles(ctx)
}

func (s *Service) CreateRole(ctx context.Context, params *iam_grpc_adapter.CreateRoleParams) (*iam_grpc_adapter.AdminRoleResult, error) {
	return s.adapters.iamGrpc.CreateRole(ctx, params)
}

func (s *Service) UpdateRole(ctx context.Context, id string, params *iam_grpc_adapter.UpdateRoleParams) (*iam_grpc_adapter.AdminRoleResult, error) {
	return s.adapters.iamGrpc.UpdateRole(ctx, id, params)
}

func (s *Service) DeleteRole(ctx context.Context, id string) error {
	return s.adapters.iamGrpc.DeleteRole(ctx, id)
}

func (s *Service) ListPermissions(ctx context.Context) (*iam_grpc_adapter.ListPermissionsResult, error) {
	return s.adapters.iamGrpc.ListPermissions(ctx)
}

func (s *Service) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	return s.adapters.iamGrpc.AssignRoleToUser(ctx, userID, roleID)
}

func (s *Service) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	return s.adapters.iamGrpc.RevokeRoleFromUser(ctx, userID, roleID)
}

// ---------------------------------------------------------------------------
// Jamaah manifest
// ---------------------------------------------------------------------------

func (s *Service) GetDepartureManifest(ctx context.Context, departureID string) (*jamaah_grpc_adapter.GetDepartureManifestResult, error) {
	return s.adapters.jamaahGrpc.GetDepartureManifest(ctx, departureID)
}
