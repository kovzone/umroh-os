// proxy_dispatch_phase6.go — gateway service dispatch for Phase 6 new routes.
//
// Covers:
//   - Catalog master data CRUD (Wave 1A): hotels, airlines, muthawwif, addons,
//     departure pricing
//   - Finance depth RPCs (Wave 1B): RecognizeRevenue, GetPLReport, GetBalanceSheet
//   - IAM admin (Wave 1C): user management, role management, permissions,
//     role assignment
//   - Jamaah manifest (Wave 1A): GetDepartureManifest
//   - Finance disbursement + aging (BL-FIN-010/011)
//   - Logistics procurement + GRN + kit (BL-LOG-010..012)
//   - Ops field scanning + bus boarding (BL-OPS-010/011)
//   - IAM security features (BL-IAM-007/011/014/016)
//   - Visa pipeline (BL-VISA-001..003)
//
// Each method is a thin delegation to the appropriate adapter.
// No business logic lives here; all logic lives in the backend services.
package service

import (
	"context"

	"gateway-svc/adapter/booking_grpc_adapter"
	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/adapter/jamaah_grpc_adapter"
	"gateway-svc/adapter/logistics_grpc_adapter"
	"gateway-svc/adapter/ops_grpc_adapter"
	"gateway-svc/adapter/visa_grpc_adapter"
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

// ---------------------------------------------------------------------------
// Finance — disbursement + aging (BL-FIN-010/011)
// ---------------------------------------------------------------------------

func (s *Service) CreateDisbursementBatch(ctx context.Context, params *finance_grpc_adapter.CreateDisbursementBatchParams) (*finance_grpc_adapter.CreateDisbursementBatchResult, error) {
	return s.adapters.financeGrpc.CreateDisbursementBatch(ctx, params)
}

func (s *Service) ApproveDisbursement(ctx context.Context, params *finance_grpc_adapter.ApproveDisbursementParams) (*finance_grpc_adapter.ApproveDisbursementResult, error) {
	return s.adapters.financeGrpc.ApproveDisbursement(ctx, params)
}

func (s *Service) GetARAPAging(ctx context.Context, params *finance_grpc_adapter.GetARAPAgingParams) (*finance_grpc_adapter.GetARAPAgingResult, error) {
	return s.adapters.financeGrpc.GetARAPAging(ctx, params)
}

// ---------------------------------------------------------------------------
// Logistics — procurement + GRN + kit (BL-LOG-010..012)
// ---------------------------------------------------------------------------

func (s *Service) CreatePurchaseRequest(ctx context.Context, params *logistics_grpc_adapter.CreatePurchaseRequestParams) (*logistics_grpc_adapter.CreatePurchaseRequestResult, error) {
	return s.adapters.logisticsGrpc.CreatePurchaseRequest(ctx, params)
}

func (s *Service) ApprovePurchaseRequest(ctx context.Context, params *logistics_grpc_adapter.ApprovePurchaseRequestParams) (*logistics_grpc_adapter.ApprovePurchaseRequestResult, error) {
	return s.adapters.logisticsGrpc.ApprovePurchaseRequest(ctx, params)
}

func (s *Service) RecordGRNWithQC(ctx context.Context, params *logistics_grpc_adapter.RecordGRNWithQCParams) (*logistics_grpc_adapter.RecordGRNWithQCResult, error) {
	return s.adapters.logisticsGrpc.RecordGRNWithQC(ctx, params)
}

func (s *Service) CreateKitAssembly(ctx context.Context, params *logistics_grpc_adapter.CreateKitAssemblyParams) (*logistics_grpc_adapter.CreateKitAssemblyResult, error) {
	return s.adapters.logisticsGrpc.CreateKitAssembly(ctx, params)
}

// ---------------------------------------------------------------------------
// Ops — field scanning + bus boarding (BL-OPS-010/011)
// ---------------------------------------------------------------------------

func (s *Service) RecordScan(ctx context.Context, params *ops_grpc_adapter.RecordScanParams) (*ops_grpc_adapter.RecordScanResult, error) {
	return s.adapters.opsGrpc.RecordScan(ctx, params)
}

func (s *Service) RecordBusBoarding(ctx context.Context, params *ops_grpc_adapter.RecordBusBoardingParams) (*ops_grpc_adapter.RecordBusBoardingResult, error) {
	return s.adapters.opsGrpc.RecordBusBoarding(ctx, params)
}

func (s *Service) GetBoardingRoster(ctx context.Context, departureID, busNumber string) (*ops_grpc_adapter.GetBoardingRosterResult, error) {
	return s.adapters.opsGrpc.GetBoardingRoster(ctx, departureID, busNumber)
}

// ---------------------------------------------------------------------------
// IAM — Phase 6 security features (BL-IAM-007/011/014/016)
// ---------------------------------------------------------------------------

func (s *Service) SetDataScope(ctx context.Context, params *iam_grpc_adapter.SetDataScopeParams) (*iam_grpc_adapter.SetDataScopeResult, error) {
	return s.adapters.iamGrpc.SetDataScope(ctx, params)
}

func (s *Service) CreateAPIKey(ctx context.Context, params *iam_grpc_adapter.CreateAPIKeyParams) (*iam_grpc_adapter.CreateAPIKeyResult, error) {
	return s.adapters.iamGrpc.CreateAPIKey(ctx, params)
}

func (s *Service) RevokeAPIKey(ctx context.Context, keyID string) (*iam_grpc_adapter.RevokeAPIKeyResult, error) {
	return s.adapters.iamGrpc.RevokeAPIKey(ctx, keyID)
}

func (s *Service) GetGlobalConfig(ctx context.Context, keys []string) (*iam_grpc_adapter.GetGlobalConfigResult, error) {
	return s.adapters.iamGrpc.GetGlobalConfig(ctx, keys)
}

func (s *Service) SetGlobalConfig(ctx context.Context, params *iam_grpc_adapter.SetGlobalConfigParams) (*iam_grpc_adapter.SetGlobalConfigResult, error) {
	return s.adapters.iamGrpc.SetGlobalConfig(ctx, params)
}

func (s *Service) SearchActivityLog(ctx context.Context, params *iam_grpc_adapter.SearchActivityLogParams) (*iam_grpc_adapter.SearchActivityLogResult, error) {
	return s.adapters.iamGrpc.SearchActivityLog(ctx, params)
}

// ---------------------------------------------------------------------------
// IAM — security depth (BL-IAM-010/012/013/015/017)
// ---------------------------------------------------------------------------

func (s *Service) GetPasswordPolicy(ctx context.Context) (*iam_grpc_adapter.PasswordPolicyResult, error) {
	return s.adapters.iamGrpc.GetPasswordPolicy(ctx)
}

func (s *Service) SetPasswordPolicy(ctx context.Context, params *iam_grpc_adapter.SetPasswordPolicyParams) (*iam_grpc_adapter.PasswordPolicyResult, error) {
	return s.adapters.iamGrpc.SetPasswordPolicy(ctx, params)
}

func (s *Service) RecordLoginAnomaly(ctx context.Context, params *iam_grpc_adapter.RecordLoginAnomalyParams) (*iam_grpc_adapter.RecordLoginAnomalyResult, error) {
	return s.adapters.iamGrpc.RecordLoginAnomaly(ctx, params)
}

func (s *Service) ListSessions(ctx context.Context, params *iam_grpc_adapter.ListSessionsParams) (*iam_grpc_adapter.ListSessionsResult, error) {
	return s.adapters.iamGrpc.ListSessions(ctx, params)
}

func (s *Service) RevokeSession(ctx context.Context, params *iam_grpc_adapter.RevokeSessionParams) (*iam_grpc_adapter.RevokeSessionResult, error) {
	return s.adapters.iamGrpc.RevokeSession(ctx, params)
}

func (s *Service) UpsertCommTemplate(ctx context.Context, params *iam_grpc_adapter.UpsertCommTemplateParams) (*iam_grpc_adapter.UpsertCommTemplateResult, error) {
	return s.adapters.iamGrpc.UpsertCommTemplate(ctx, params)
}

func (s *Service) ListCommTemplates(ctx context.Context, params *iam_grpc_adapter.ListCommTemplatesParams) (*iam_grpc_adapter.ListCommTemplatesResult, error) {
	return s.adapters.iamGrpc.ListCommTemplates(ctx, params)
}

func (s *Service) TriggerBackup(ctx context.Context, params *iam_grpc_adapter.TriggerBackupParams) (*iam_grpc_adapter.TriggerBackupResult, error) {
	return s.adapters.iamGrpc.TriggerBackup(ctx, params)
}

func (s *Service) GetBackupHistory(ctx context.Context, limit int32) (*iam_grpc_adapter.GetBackupHistoryResult, error) {
	return s.adapters.iamGrpc.GetBackupHistory(ctx, limit)
}

// ---------------------------------------------------------------------------
// Visa — pipeline (BL-VISA-001..003)
// ---------------------------------------------------------------------------

func (s *Service) TransitionVisaStatus(ctx context.Context, params *visa_grpc_adapter.TransitionStatusParams) (*visa_grpc_adapter.TransitionStatusResult, error) {
	return s.adapters.visaGrpc.TransitionStatus(ctx, params)
}

func (s *Service) BulkSubmitVisa(ctx context.Context, params *visa_grpc_adapter.BulkSubmitParams) (*visa_grpc_adapter.BulkSubmitResult, error) {
	return s.adapters.visaGrpc.BulkSubmit(ctx, params)
}

func (s *Service) GetVisaApplications(ctx context.Context, departureID, statusFilter string) (*visa_grpc_adapter.GetApplicationsResult, error) {
	return s.adapters.visaGrpc.GetApplications(ctx, departureID, statusFilter)
}

// ---------------------------------------------------------------------------
// Vendor readiness (BL-OPS-020)
// ---------------------------------------------------------------------------

func (s *Service) UpdateVendorReadiness(ctx context.Context, params *catalog_grpc_adapter.UpdateVendorReadinessParams) (*catalog_grpc_adapter.ReadinessResult, error) {
	return s.adapters.catalogGrpc.UpdateVendorReadiness(ctx, params)
}

func (s *Service) GetDepartureReadiness(ctx context.Context, params *catalog_grpc_adapter.GetDepartureReadinessParams) (*catalog_grpc_adapter.ReadinessResult, error) {
	return s.adapters.catalogGrpc.GetDepartureReadiness(ctx, params)
}

// ---------------------------------------------------------------------------
// Catalog depth — Wave 3 (BL-CAT-010/011/013)
// ---------------------------------------------------------------------------

func (s *Service) BulkImportPackages(ctx context.Context, params *catalog_grpc_adapter.BulkImportPackagesParams) (*catalog_grpc_adapter.BulkImportPackagesResult, error) {
	return s.adapters.catalogGrpc.BulkImportPackages(ctx, params)
}

func (s *Service) BulkUpdatePackages(ctx context.Context, params *catalog_grpc_adapter.BulkUpdatePackagesParams) (*catalog_grpc_adapter.BulkUpdatePackagesResult, error) {
	return s.adapters.catalogGrpc.BulkUpdatePackages(ctx, params)
}

func (s *Service) GetPackageVersion(ctx context.Context, packageID string) (*catalog_grpc_adapter.PackageVersionResult, error) {
	return s.adapters.catalogGrpc.GetPackageVersion(ctx, packageID)
}

// ---------------------------------------------------------------------------
// Booking depth — Wave 3 (BL-BOOK-007)
// ---------------------------------------------------------------------------

func (s *Service) GetSeatsByChannel(ctx context.Context, departureID string) (*booking_grpc_adapter.GetSeatsByChannelResult, error) {
	return s.adapters.bookingGrpc.GetSeatsByChannel(ctx, departureID)
}
