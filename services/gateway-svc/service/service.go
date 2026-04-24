package service

import (
	"context"

	"gateway-svc/adapter/booking_grpc_adapter"
	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/crm_grpc_adapter"
	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/adapter/finance_rest_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/adapter/iam_rest_adapter"
	"gateway-svc/adapter/jamaah_grpc_adapter"
	"gateway-svc/adapter/logistics_grpc_adapter"
	"gateway-svc/adapter/ops_grpc_adapter"
	"gateway-svc/adapter/payment_grpc_adapter"
	"gateway-svc/adapter/visa_grpc_adapter"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for gateway-svc.
//
// gateway-svc fronts every backend. Methods either return process-local state
// (Liveness, Readiness) or dispatch to a backend through the corresponding
// adapter. Under ADR 0009 the target shape is REST-in, gRPC-out:
//   - catalog-svc: gRPC-only from S1-E-11 (catalog_grpc_adapter)
//   - iam-svc: gRPC-only from S1-E-12 (iam_grpc_adapter); iam_rest_adapter
//     retained only for interim probes until HTTP surface is fully removed.
//   - finance_rest_adapter: retires with BL-IAM-019 / S1-E-14.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)

	// Per-backend liveness proxies — IAM still serves via REST adapter until
	// S1-E-12 removes the IAM REST surface.
	GetIamSystemLive(ctx context.Context) (*iam_rest_adapter.LivenessResult, error)

	// Traced cross-service path — the S0-J-05 observability acceptance check
	// flows through gateway-svc → iam-svc here.
	GetIamSystemDbTxDiagnostic(ctx context.Context, message string) (*iam_rest_adapter.DbTxDiagnosticResult, error)

	// catalog-svc is gRPC-only from S1-E-11; no REST liveness proxy.
	// Public catalog read (BL-GTW-002 / S1-E-10) — proxies to catalog-svc
	// gRPC via catalog_grpc_adapter. Mirrors GET /v1/packages,
	// /v1/packages/{id}, /v1/package-departures/{id}.
	ListPackages(ctx context.Context, params *catalog_grpc_adapter.ListPackagesParams) (*catalog_grpc_adapter.ListPackagesResult, error)
	GetPackage(ctx context.Context, params *catalog_grpc_adapter.GetPackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	GetPackageDeparture(ctx context.Context, params *catalog_grpc_adapter.GetPackageDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)

	// Staff catalog write (BL-CAT-014 / S1-E-07) — bearer-protected; gateway
	// handler calls CheckPermission(catalog.package.manage) before delegating.
	CreatePackage(ctx context.Context, params *catalog_grpc_adapter.CreatePackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	UpdatePackage(ctx context.Context, params *catalog_grpc_adapter.UpdatePackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	DeletePackage(ctx context.Context, params *catalog_grpc_adapter.DeletePackageParams) (*catalog_grpc_adapter.DeletePackageResult, error)
	CreateDeparture(ctx context.Context, params *catalog_grpc_adapter.CreateDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)
	UpdateDeparture(ctx context.Context, params *catalog_grpc_adapter.UpdateDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)

	// IAM auth routes (BL-IAM-018 / S1-E-12) — proxied via iam_grpc_adapter.
	// Public (no bearer): Login, RefreshSession.
	// Bearer-required (via middleware): Logout, GetMe, EnrollTOTP, VerifyTOTP.
	// Bearer + permission: SuspendUser.
	CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error)
	Login(ctx context.Context, params *iam_grpc_adapter.LoginParams) (*iam_grpc_adapter.LoginResult, error)
	RefreshSession(ctx context.Context, params *iam_grpc_adapter.RefreshSessionParams) (*iam_grpc_adapter.RefreshSessionResult, error)
	Logout(ctx context.Context, params *iam_grpc_adapter.LogoutParams) error
	GetMe(ctx context.Context, params *iam_grpc_adapter.GetMeParams) (*iam_grpc_adapter.GetMeResult, error)
	EnrollTOTP(ctx context.Context, params *iam_grpc_adapter.EnrollTOTPParams) (*iam_grpc_adapter.EnrollTOTPResult, error)
	VerifyTOTP(ctx context.Context, params *iam_grpc_adapter.VerifyTOTPParams) (*iam_grpc_adapter.VerifyTOTPResult, error)
	SuspendUser(ctx context.Context, params *iam_grpc_adapter.SuspendUserParams) (*iam_grpc_adapter.SuspendUserResult, error)

	// Finance liveness proxy — interim REST adapter; retires with
	// BL-IAM-019 / S1-E-14 when /v1/finance/* moves to gRPC.
	GetFinanceSystemLive(ctx context.Context) (*finance_rest_adapter.LivenessResult, error)

	// Finance report routes (S5-E-01 / BL-FIN-004..005) — proxied via finance_grpc_adapter.
	// GET /v1/finance/summary  — bearer; aggregate per-account balances.
	GetFinanceSummary(ctx context.Context) (*finance_grpc_adapter.GetFinanceSummaryResult, error)
	// GET /v1/finance/journals — bearer; paginated journal entries + lines.
	ListJournalEntries(ctx context.Context, params *finance_grpc_adapter.ListJournalEntriesParams) (*finance_grpc_adapter.ListJournalEntriesResult, error)

	// Finance depth routes (Phase 6 / Wave 1B) — bearer required.
	RecognizeRevenue(ctx context.Context, departureID string, totalAmountIdr int64) (*finance_grpc_adapter.RecognizeRevenueResult, error)
	GetPLReport(ctx context.Context, from, to string) (*finance_grpc_adapter.PLReportResult, error)
	GetBalanceSheet(ctx context.Context, asOf string) (*finance_grpc_adapter.BalanceSheetResult, error)

	// Booking draft (BL-BOOK-001..006 / S1-E-03 / BL-GTW-003) — proxied via
	// booking_grpc_adapter. Public in S1 (auth arrives with F4).
	CreateDraftBooking(ctx context.Context, params *booking_grpc_adapter.CreateDraftBookingParams) (*booking_grpc_adapter.CreateDraftBookingResult, error)

	// SubmitBooking transitions draft → pending_payment (S2 / BL-BOOK-005).
	// Bearer required (CS/admin action).
	SubmitBooking(ctx context.Context, params *booking_grpc_adapter.SubmitBookingParams) (*booking_grpc_adapter.SubmitBookingResult, error)

	// CRM lead management (S4-E-02 / BL-CRM-001..003) — proxied via crm_grpc_adapter.
	// POST /v1/leads — public (lead capture from landing pages, B2C forms).
	CreateLead(ctx context.Context, params *crm_grpc_adapter.CreateLeadParams) (*crm_grpc_adapter.LeadResult, error)
	// GET /v1/leads — bearer (cs/admin only).
	ListLeads(ctx context.Context, params *crm_grpc_adapter.ListLeadsParams) (*crm_grpc_adapter.ListLeadsResult, error)
	// GET /v1/leads/:id — bearer.
	GetLead(ctx context.Context, id string) (*crm_grpc_adapter.LeadResult, error)
	// PUT /v1/leads/:id — bearer.
	UpdateLead(ctx context.Context, params *crm_grpc_adapter.UpdateLeadParams) (*crm_grpc_adapter.LeadResult, error)

	// Catalog master data (Phase 6 / Wave 1A) — bearer required.
	ListHotels(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListHotelsResult, error)
	CreateHotel(ctx context.Context, params *catalog_grpc_adapter.CreateHotelParams) (*catalog_grpc_adapter.HotelResult, error)
	UpdateHotel(ctx context.Context, params *catalog_grpc_adapter.UpdateHotelParams) (*catalog_grpc_adapter.HotelResult, error)
	DeleteHotel(ctx context.Context, id string) error
	ListAirlines(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListAirlinesResult, error)
	CreateAirline(ctx context.Context, params *catalog_grpc_adapter.CreateAirlineParams) (*catalog_grpc_adapter.AirlineResult, error)
	UpdateAirline(ctx context.Context, params *catalog_grpc_adapter.UpdateAirlineParams) (*catalog_grpc_adapter.AirlineResult, error)
	DeleteAirline(ctx context.Context, id string) error
	ListMuthawwif(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListMuthawwifResult, error)
	CreateMuthawwif(ctx context.Context, params *catalog_grpc_adapter.CreateMuthawwifParams) (*catalog_grpc_adapter.MuthawwifResult, error)
	UpdateMuthawwif(ctx context.Context, params *catalog_grpc_adapter.UpdateMuthawwifParams) (*catalog_grpc_adapter.MuthawwifResult, error)
	DeleteMuthawwif(ctx context.Context, id string) error
	ListAddons(ctx context.Context, params *catalog_grpc_adapter.ListMastersParams) (*catalog_grpc_adapter.ListAddonsResult, error)
	CreateAddon(ctx context.Context, params *catalog_grpc_adapter.CreateAddonParams) (*catalog_grpc_adapter.AddonResult, error)
	UpdateAddon(ctx context.Context, params *catalog_grpc_adapter.UpdateAddonParams) (*catalog_grpc_adapter.AddonResult, error)
	DeleteAddon(ctx context.Context, id string) error
	SetDeparturePricing(ctx context.Context, departureID string, pricings []*catalog_grpc_adapter.PricingUpsertInput) (*catalog_grpc_adapter.PricingResult, error)
	GetDeparturePricing(ctx context.Context, departureID string) (*catalog_grpc_adapter.PricingResult, error)

	// IAM admin routes (Phase 6 / Wave 1C) — bearer required.
	ListUsers(ctx context.Context, params *iam_grpc_adapter.ListUsersParams) (*iam_grpc_adapter.ListUsersResult, error)
	CreateUser(ctx context.Context, params *iam_grpc_adapter.CreateUserParams) (*iam_grpc_adapter.AdminUserResult, error)
	UpdateUser(ctx context.Context, id string, params *iam_grpc_adapter.UpdateUserParams) (*iam_grpc_adapter.AdminUserResult, error)
	GetUser(ctx context.Context, id string) (*iam_grpc_adapter.AdminUserResult, error)
	ResetUserPassword(ctx context.Context, id, newPassword string) error
	ListRoles(ctx context.Context) (*iam_grpc_adapter.ListRolesResult, error)
	CreateRole(ctx context.Context, params *iam_grpc_adapter.CreateRoleParams) (*iam_grpc_adapter.AdminRoleResult, error)
	UpdateRole(ctx context.Context, id string, params *iam_grpc_adapter.UpdateRoleParams) (*iam_grpc_adapter.AdminRoleResult, error)
	DeleteRole(ctx context.Context, id string) error
	ListPermissions(ctx context.Context) (*iam_grpc_adapter.ListPermissionsResult, error)
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RevokeRoleFromUser(ctx context.Context, userID, roleID string) error

	// Jamaah manifest (Phase 6 / Wave 1A) — bearer required.
	GetDepartureManifest(ctx context.Context, departureID string) (*jamaah_grpc_adapter.GetDepartureManifestResult, error)

	// S3 ops routes (S3 Wave 2) — bearer required.
	RunRoomAllocation(ctx context.Context, departureID string, jamaahIDs []string) (*ops_grpc_adapter.RunRoomAllocationResult, error)
	GetRoomAllocation(ctx context.Context, departureID string) (*ops_grpc_adapter.GetRoomAllocationResult, error)
	GenerateIDCard(ctx context.Context, jamaahID, departureID, cardType, jamaahName, departureName string) (*ops_grpc_adapter.GenerateIDCardResult, error)
	VerifyIDCard(ctx context.Context, token string) (*ops_grpc_adapter.VerifyIDCardResult, error)
	ExportManifest(ctx context.Context, departureID string) (*ops_grpc_adapter.ExportManifestResult, error)

	// S3 logistics routes (S3 Wave 2) — bearer required.
	ShipFulfillmentTask(ctx context.Context, bookingID, carrier, notes string) (*logistics_grpc_adapter.ShipFulfillmentTaskResult, error)
	GeneratePickupQR(ctx context.Context, bookingID string) (*logistics_grpc_adapter.GeneratePickupQRResult, error)
	RedeemPickupQR(ctx context.Context, token string) (*logistics_grpc_adapter.RedeemPickupQRResult, error)
	// ListFulfillmentTasks returns a paginated list of fulfillment tasks (ISSUE-018).
	ListFulfillmentTasks(ctx context.Context, params *logistics_grpc_adapter.ListFulfillmentTasksParams) (*logistics_grpc_adapter.ListFulfillmentTasksResult, error)

	// S3 finance GRN (S3 Wave 2) — bearer required.
	OnGRNReceived(ctx context.Context, grnID, departureID string, amountIdr int64) (*finance_grpc_adapter.OnGRNReceivedResult, error)

	// S5 finance correction (BL-FIN-006) — bearer required.
	// POST /v1/finance/journals/:id/correct — post reversing counter-entry.
	CorrectJournal(ctx context.Context, params *finance_grpc_adapter.CorrectJournalParams) (*finance_grpc_adapter.CorrectJournalResult, error)

	// S3 jamaah OCR (S3 Wave 2) — bearer required.
	TriggerOCR(ctx context.Context, documentID string) (*jamaah_grpc_adapter.TriggerOCRResult, error)
	GetOCRStatus(ctx context.Context, documentID string) (*jamaah_grpc_adapter.GetOCRStatusResult, error)

	// S2 payment link (BL-PAY-020) — bearer required.
	// POST /v1/payments/link — CS re-issues VA for existing booking.
	ReissuePaymentLink(ctx context.Context, params *payment_grpc_adapter.ReissuePaymentLinkParams) (*payment_grpc_adapter.ReissuePaymentLinkResult, error)

	// S2 invoice routes (BL-PAY-001 / ISSUE-005) — bearer required.
	// POST /v1/invoices — create invoice + VA for a booking.
	IssueVirtualAccount(ctx context.Context, params *payment_grpc_adapter.IssueVirtualAccountParams) (*payment_grpc_adapter.IssueVirtualAccountResult, error)
	// GET /v1/invoices/:id — fetch invoice details by UUID.
	GetInvoiceByID(ctx context.Context, params *payment_grpc_adapter.GetInvoiceByIDParams) (*payment_grpc_adapter.GetInvoiceByIDResult, error)

	// S2 webhook routes (BL-PAY-003/004 / ISSUE-007/008) — public, signature-protected.
	// POST /v1/webhooks/:gateway — forward raw webhook payload to payment-svc.
	ProcessWebhook(ctx context.Context, params *payment_grpc_adapter.ProcessWebhookParams) (*payment_grpc_adapter.ProcessWebhookResult, error)

	// Phase 6 finance disbursement + aging (BL-FIN-010/011) — bearer required.
	CreateDisbursementBatch(ctx context.Context, params *finance_grpc_adapter.CreateDisbursementBatchParams) (*finance_grpc_adapter.CreateDisbursementBatchResult, error)
	ApproveDisbursement(ctx context.Context, params *finance_grpc_adapter.ApproveDisbursementParams) (*finance_grpc_adapter.ApproveDisbursementResult, error)
	GetARAPAging(ctx context.Context, params *finance_grpc_adapter.GetARAPAgingParams) (*finance_grpc_adapter.GetARAPAgingResult, error)

	// Phase 6 logistics procurement + GRN + kit (BL-LOG-010..012) — bearer required.
	CreatePurchaseRequest(ctx context.Context, params *logistics_grpc_adapter.CreatePurchaseRequestParams) (*logistics_grpc_adapter.CreatePurchaseRequestResult, error)
	ApprovePurchaseRequest(ctx context.Context, params *logistics_grpc_adapter.ApprovePurchaseRequestParams) (*logistics_grpc_adapter.ApprovePurchaseRequestResult, error)
	RecordGRNWithQC(ctx context.Context, params *logistics_grpc_adapter.RecordGRNWithQCParams) (*logistics_grpc_adapter.RecordGRNWithQCResult, error)
	CreateKitAssembly(ctx context.Context, params *logistics_grpc_adapter.CreateKitAssemblyParams) (*logistics_grpc_adapter.CreateKitAssemblyResult, error)

	// Phase 6 ops field scanning + bus boarding (BL-OPS-010/011) — bearer required.
	RecordScan(ctx context.Context, params *ops_grpc_adapter.RecordScanParams) (*ops_grpc_adapter.RecordScanResult, error)
	RecordBusBoarding(ctx context.Context, params *ops_grpc_adapter.RecordBusBoardingParams) (*ops_grpc_adapter.RecordBusBoardingResult, error)
	GetBoardingRoster(ctx context.Context, departureID, busNumber string) (*ops_grpc_adapter.GetBoardingRosterResult, error)

	// Phase 6 IAM security features (BL-IAM-007/011/014/016) — bearer required.
	SetDataScope(ctx context.Context, params *iam_grpc_adapter.SetDataScopeParams) (*iam_grpc_adapter.SetDataScopeResult, error)
	CreateAPIKey(ctx context.Context, params *iam_grpc_adapter.CreateAPIKeyParams) (*iam_grpc_adapter.CreateAPIKeyResult, error)
	RevokeAPIKey(ctx context.Context, keyID string) (*iam_grpc_adapter.RevokeAPIKeyResult, error)
	GetGlobalConfig(ctx context.Context, keys []string) (*iam_grpc_adapter.GetGlobalConfigResult, error)
	SetGlobalConfig(ctx context.Context, params *iam_grpc_adapter.SetGlobalConfigParams) (*iam_grpc_adapter.SetGlobalConfigResult, error)
	SearchActivityLog(ctx context.Context, params *iam_grpc_adapter.SearchActivityLogParams) (*iam_grpc_adapter.SearchActivityLogResult, error)

	// IAM security depth (BL-IAM-010/012/013/015/017) — Wave 2.
	GetPasswordPolicy(ctx context.Context) (*iam_grpc_adapter.PasswordPolicyResult, error)
	SetPasswordPolicy(ctx context.Context, params *iam_grpc_adapter.SetPasswordPolicyParams) (*iam_grpc_adapter.PasswordPolicyResult, error)
	RecordLoginAnomaly(ctx context.Context, params *iam_grpc_adapter.RecordLoginAnomalyParams) (*iam_grpc_adapter.RecordLoginAnomalyResult, error)
	ListSessions(ctx context.Context, params *iam_grpc_adapter.ListSessionsParams) (*iam_grpc_adapter.ListSessionsResult, error)
	RevokeSession(ctx context.Context, params *iam_grpc_adapter.RevokeSessionParams) (*iam_grpc_adapter.RevokeSessionResult, error)
	UpsertCommTemplate(ctx context.Context, params *iam_grpc_adapter.UpsertCommTemplateParams) (*iam_grpc_adapter.UpsertCommTemplateResult, error)
	ListCommTemplates(ctx context.Context, params *iam_grpc_adapter.ListCommTemplatesParams) (*iam_grpc_adapter.ListCommTemplatesResult, error)
	TriggerBackup(ctx context.Context, params *iam_grpc_adapter.TriggerBackupParams) (*iam_grpc_adapter.TriggerBackupResult, error)
	GetBackupHistory(ctx context.Context, limit int32) (*iam_grpc_adapter.GetBackupHistoryResult, error)

	// Phase 6 visa pipeline (BL-VISA-001..003) — bearer required.
	TransitionVisaStatus(ctx context.Context, params *visa_grpc_adapter.TransitionStatusParams) (*visa_grpc_adapter.TransitionStatusResult, error)
	BulkSubmitVisa(ctx context.Context, params *visa_grpc_adapter.BulkSubmitParams) (*visa_grpc_adapter.BulkSubmitResult, error)
	GetVisaApplications(ctx context.Context, departureID, statusFilter string) (*visa_grpc_adapter.GetApplicationsResult, error)

	// Vendor readiness (BL-OPS-020) — bearer required.
	UpdateVendorReadiness(ctx context.Context, params *catalog_grpc_adapter.UpdateVendorReadinessParams) (*catalog_grpc_adapter.ReadinessResult, error)
	GetDepartureReadiness(ctx context.Context, params *catalog_grpc_adapter.GetDepartureReadinessParams) (*catalog_grpc_adapter.ReadinessResult, error)

	// Catalog depth — Wave 3 (BL-CAT-010/011/013) — bearer required.
	BulkImportPackages(ctx context.Context, params *catalog_grpc_adapter.BulkImportPackagesParams) (*catalog_grpc_adapter.BulkImportPackagesResult, error)
	BulkUpdatePackages(ctx context.Context, params *catalog_grpc_adapter.BulkUpdatePackagesParams) (*catalog_grpc_adapter.BulkUpdatePackagesResult, error)
	GetPackageVersion(ctx context.Context, packageID string) (*catalog_grpc_adapter.PackageVersionResult, error)

	// Booking depth — Wave 3 (BL-BOOK-007) — bearer required.
	GetSeatsByChannel(ctx context.Context, departureID string) (*booking_grpc_adapter.GetSeatsByChannelResult, error)

	// Ops depth — Wave 5 (BL-OPS-021..042) — bearer required.
	StoreCollectiveDocument(ctx context.Context, params *ops_grpc_adapter.StoreCollectiveDocumentParams) (*ops_grpc_adapter.StoreCollectiveDocumentResult, error)
	GetCollectiveDocuments(ctx context.Context, params *ops_grpc_adapter.GetCollectiveDocumentsParams) (*ops_grpc_adapter.GetCollectiveDocumentsResult, error)
	SetDocumentACL(ctx context.Context, params *ops_grpc_adapter.SetDocumentACLParams) (*ops_grpc_adapter.SetDocumentACLResult, error)
	ExtractPassportOCR(ctx context.Context, params *ops_grpc_adapter.ExtractPassportOCRParams) (*ops_grpc_adapter.ExtractPassportOCRResult, error)
	SetMahramRelation(ctx context.Context, params *ops_grpc_adapter.SetMahramRelationParams) (*ops_grpc_adapter.SetMahramRelationResult, error)
	GetMahramRelations(ctx context.Context, params *ops_grpc_adapter.GetMahramRelationsParams) (*ops_grpc_adapter.GetMahramRelationsResult, error)
	GetDocumentProgress(ctx context.Context, params *ops_grpc_adapter.GetDocumentProgressParams) (*ops_grpc_adapter.GetDocumentProgressResult, error)
	GetExpiryAlerts(ctx context.Context, params *ops_grpc_adapter.GetExpiryAlertsParams) (*ops_grpc_adapter.GetExpiryAlertsResult, error)
	GenerateOfficialLetter(ctx context.Context, params *ops_grpc_adapter.GenerateOfficialLetterParams) (*ops_grpc_adapter.GenerateOfficialLetterResult, error)
	GenerateImmigrationManifest(ctx context.Context, params *ops_grpc_adapter.GenerateImmigrationManifestParams) (*ops_grpc_adapter.GenerateImmigrationManifestResult, error)
	AssignTransport(ctx context.Context, params *ops_grpc_adapter.AssignTransportParams) (*ops_grpc_adapter.AssignTransportResult, error)
	GetTransportAssignments(ctx context.Context, params *ops_grpc_adapter.GetTransportAssignmentsParams) (*ops_grpc_adapter.GetTransportAssignmentsResult, error)
	PublishManifestDelta(ctx context.Context, params *ops_grpc_adapter.PublishManifestDeltaParams) (*ops_grpc_adapter.PublishManifestDeltaResult, error)
	AssignStaff(ctx context.Context, params *ops_grpc_adapter.AssignStaffParams) (*ops_grpc_adapter.AssignStaffResult, error)
	RecordPassportHandover(ctx context.Context, params *ops_grpc_adapter.RecordPassportHandoverParams) (*ops_grpc_adapter.RecordPassportHandoverResult, error)
	GetPassportLog(ctx context.Context, params *ops_grpc_adapter.GetPassportLogParams) (*ops_grpc_adapter.GetPassportLogResult, error)
	GetVisaProgress(ctx context.Context, params *ops_grpc_adapter.GetVisaProgressParams) (*ops_grpc_adapter.GetVisaProgressResult, error)
	StoreEVisa(ctx context.Context, params *ops_grpc_adapter.StoreEVisaParams) (*ops_grpc_adapter.StoreEVisaResult, error)
	GetEVisa(ctx context.Context, params *ops_grpc_adapter.GetEVisaParams) (*ops_grpc_adapter.GetEVisaResult, error)
	TriggerExternalProvider(ctx context.Context, params *ops_grpc_adapter.TriggerExternalProviderParams) (*ops_grpc_adapter.TriggerExternalProviderResult, error)
	CreateRefund(ctx context.Context, params *ops_grpc_adapter.CreateRefundParams) (*ops_grpc_adapter.CreateRefundResult, error)
	ApproveRefund(ctx context.Context, params *ops_grpc_adapter.ApproveRefundParams) (*ops_grpc_adapter.ApproveRefundResult, error)
	RecordPenalty(ctx context.Context, params *ops_grpc_adapter.RecordPenaltyParams) (*ops_grpc_adapter.RecordPenaltyResult, error)
	RecordLuggageScan(ctx context.Context, params *ops_grpc_adapter.RecordLuggageScanParams) (*ops_grpc_adapter.RecordLuggageScanResult, error)
	GetLuggageCount(ctx context.Context, params *ops_grpc_adapter.GetLuggageCountParams) (*ops_grpc_adapter.GetLuggageCountResult, error)
	BroadcastSchedule(ctx context.Context, params *ops_grpc_adapter.BroadcastScheduleParams) (*ops_grpc_adapter.BroadcastScheduleResult, error)
	IssueDigitalTasreh(ctx context.Context, params *ops_grpc_adapter.IssueDigitalTasrehParams) (*ops_grpc_adapter.IssueDigitalTasrehResult, error)
	RecordRaudhahEntry(ctx context.Context, params *ops_grpc_adapter.RecordRaudhahEntryParams) (*ops_grpc_adapter.RecordRaudhahEntryResult, error)
	RegisterAudioDevice(ctx context.Context, params *ops_grpc_adapter.RegisterAudioDeviceParams) (*ops_grpc_adapter.RegisterAudioDeviceResult, error)
	UpdateAudioDeviceStatus(ctx context.Context, params *ops_grpc_adapter.UpdateAudioDeviceStatusParams) (*ops_grpc_adapter.UpdateAudioDeviceStatusResult, error)
	RecordZamzamDistribution(ctx context.Context, params *ops_grpc_adapter.RecordZamzamDistributionParams) (*ops_grpc_adapter.RecordZamzamDistributionResult, error)
	RecordRoomCheckIn(ctx context.Context, params *ops_grpc_adapter.RecordRoomCheckInParams) (*ops_grpc_adapter.RecordRoomCheckInResult, error)

	// Logistics depth — Wave 5 (BL-LOG-013..029) — bearer required.
	ListPurchaseRequests(ctx context.Context, params *logistics_grpc_adapter.ListPurchaseRequestsParams) (*logistics_grpc_adapter.ListPurchaseRequestsResult, error)
	GetBudgetSyncStatus(ctx context.Context, params *logistics_grpc_adapter.GetBudgetSyncStatusParams) (*logistics_grpc_adapter.GetBudgetSyncStatusResult, error)
	GetTieredApprovals(ctx context.Context, params *logistics_grpc_adapter.GetTieredApprovalsParams) (*logistics_grpc_adapter.GetTieredApprovalsResult, error)
	DecideTieredApproval(ctx context.Context, params *logistics_grpc_adapter.DecideTieredApprovalParams) (*logistics_grpc_adapter.DecideTieredApprovalResult, error)
	AutoSelectVendor(ctx context.Context, params *logistics_grpc_adapter.AutoSelectVendorParams) (*logistics_grpc_adapter.AutoSelectVendorResult, error)
	RecordPartialGRN(ctx context.Context, params *logistics_grpc_adapter.RecordPartialGRNParams) (*logistics_grpc_adapter.RecordPartialGRNResult, error)
	ReverseGRN(ctx context.Context, params *logistics_grpc_adapter.ReverseGRNParams) (*logistics_grpc_adapter.ReverseGRNResult, error)
	GenerateBarcode(ctx context.Context, params *logistics_grpc_adapter.GenerateBarcodeParams) (*logistics_grpc_adapter.GenerateBarcodeResult, error)
	PrintSKULabel(ctx context.Context, params *logistics_grpc_adapter.PrintSKULabelParams) (*logistics_grpc_adapter.PrintSKULabelResult, error)
	CreateWarehouse(ctx context.Context, params *logistics_grpc_adapter.CreateWarehouseParams) (*logistics_grpc_adapter.CreateWarehouseResult, error)
	TransferStock(ctx context.Context, params *logistics_grpc_adapter.TransferStockParams) (*logistics_grpc_adapter.TransferStockResult, error)
	GetStockAlerts(ctx context.Context, params *logistics_grpc_adapter.GetStockAlertsParams) (*logistics_grpc_adapter.GetStockAlertsResult, error)
	SetReorderLevel(ctx context.Context, params *logistics_grpc_adapter.SetReorderLevelParams) (*logistics_grpc_adapter.SetReorderLevelResult, error)
	StartStocktake(ctx context.Context, params *logistics_grpc_adapter.StartStocktakeParams) (*logistics_grpc_adapter.StartStocktakeResult, error)
	RecordStocktakeCount(ctx context.Context, params *logistics_grpc_adapter.RecordStocktakeCountParams) (*logistics_grpc_adapter.RecordStocktakeCountResult, error)
	FinalizeStocktake(ctx context.Context, params *logistics_grpc_adapter.FinalizeStocktakeParams) (*logistics_grpc_adapter.FinalizeStocktakeResult, error)
	SyncFulfillmentSizes(ctx context.Context, params *logistics_grpc_adapter.SyncFulfillmentSizesParams) (*logistics_grpc_adapter.SyncFulfillmentSizesResult, error)
	RecordCourierTracking(ctx context.Context, params *logistics_grpc_adapter.RecordCourierTrackingParams) (*logistics_grpc_adapter.RecordCourierTrackingResult, error)
	RecordReturn(ctx context.Context, params *logistics_grpc_adapter.RecordReturnParams) (*logistics_grpc_adapter.RecordReturnResult, error)
	ProcessExchange(ctx context.Context, params *logistics_grpc_adapter.ProcessExchangeParams) (*logistics_grpc_adapter.ProcessExchangeResult, error)

	// Wave 4 Finance depth (BL-FIN-020..041) — bearer required.
	ScheduleBilling(ctx context.Context, params *finance_grpc_adapter.ScheduleBillingParams) (*finance_grpc_adapter.ScheduleBillingResult, error)
	RecordBankTransaction(ctx context.Context, params *finance_grpc_adapter.RecordBankTransactionParams) (*finance_grpc_adapter.RecordBankTransactionResult, error)
	GetBankReconciliation(ctx context.Context, params *finance_grpc_adapter.GetBankReconciliationParams) (*finance_grpc_adapter.GetBankReconciliationResult, error)
	GetARSubledger(ctx context.Context, params *finance_grpc_adapter.GetARSubledgerParams) (*finance_grpc_adapter.GetARSubledgerResult, error)
	IssueDigitalReceipt(ctx context.Context, params *finance_grpc_adapter.IssueDigitalReceiptParams) (*finance_grpc_adapter.IssueDigitalReceiptResult, error)
	GetDigitalReceipt(ctx context.Context, params *finance_grpc_adapter.GetDigitalReceiptParams) (*finance_grpc_adapter.DigitalReceiptResult, error)
	RecordManualPayment(ctx context.Context, params *finance_grpc_adapter.RecordManualPaymentParams) (*finance_grpc_adapter.RecordManualPaymentResult, error)
	CreateFinanceVendor(ctx context.Context, params *finance_grpc_adapter.CreateVendorParams) (*finance_grpc_adapter.CreateVendorResult, error)
	UpdateFinanceVendor(ctx context.Context, params *finance_grpc_adapter.UpdateVendorParams) (*finance_grpc_adapter.UpdateVendorResult, error)
	ListFinanceVendors(ctx context.Context, params *finance_grpc_adapter.ListVendorsParams) (*finance_grpc_adapter.ListVendorsResult, error)
	DeleteFinanceVendor(ctx context.Context, params *finance_grpc_adapter.DeleteVendorParams) (*finance_grpc_adapter.DeleteVendorResult, error)
	GetAPSubledger(ctx context.Context, params *finance_grpc_adapter.GetAPSubledgerParams) (*finance_grpc_adapter.GetAPSubledgerResult, error)
	ListPendingAuthorizations(ctx context.Context, params *finance_grpc_adapter.ListPendingAuthorizationsParams) (*finance_grpc_adapter.ListPendingAuthorizationsResult, error)
	DecidePaymentAuthorization(ctx context.Context, params *finance_grpc_adapter.DecidePaymentAuthorizationParams) (*finance_grpc_adapter.DecidePaymentAuthorizationResult, error)
	RecordPettyCash(ctx context.Context, params *finance_grpc_adapter.RecordPettyCashParams) (*finance_grpc_adapter.RecordPettyCashResult, error)
	ClosePettyCashPeriod(ctx context.Context, params *finance_grpc_adapter.ClosePettyCashPeriodParams) (*finance_grpc_adapter.ClosePettyCashPeriodResult, error)
	GetProjectCosts(ctx context.Context, params *finance_grpc_adapter.GetProjectCostsParams) (*finance_grpc_adapter.GetProjectCostsResult, error)
	GetDeparturePL(ctx context.Context, params *finance_grpc_adapter.GetDeparturePLParams) (*finance_grpc_adapter.GetDeparturePLResult, error)
	GetBudgetVsActual(ctx context.Context, params *finance_grpc_adapter.GetBudgetVsActualParams) (*finance_grpc_adapter.GetBudgetVsActualResult, error)
	TriggerAutoJournal(ctx context.Context, params *finance_grpc_adapter.TriggerAutoJournalParams) (*finance_grpc_adapter.TriggerAutoJournalResult, error)
	GetRevenueRecognitionPolicy(ctx context.Context) (*finance_grpc_adapter.RevenueRecognitionPolicyResult, error)
	SetRevenueRecognitionPolicy(ctx context.Context, params *finance_grpc_adapter.SetRevenueRecognitionPolicyParams) error
	SetExchangeRate(ctx context.Context, params *finance_grpc_adapter.SetExchangeRateParams) (*finance_grpc_adapter.SetExchangeRateResult, error)
	GetExchangeRate(ctx context.Context, params *finance_grpc_adapter.GetExchangeRateParams) (*finance_grpc_adapter.GetExchangeRateResult, error)
	CreateFixedAsset(ctx context.Context, params *finance_grpc_adapter.CreateFixedAssetParams) (*finance_grpc_adapter.CreateFixedAssetResult, error)
	ListFixedAssets(ctx context.Context, category string) (*finance_grpc_adapter.ListFixedAssetsResult, error)
	RunDepreciation(ctx context.Context, params *finance_grpc_adapter.RunDepreciationParams) (*finance_grpc_adapter.RunDepreciationResult, error)
	CalculateTax(ctx context.Context, params *finance_grpc_adapter.CalculateTaxParams) (*finance_grpc_adapter.CalculateTaxResult, error)
	GetTaxReport(ctx context.Context, params *finance_grpc_adapter.GetTaxReportParams) (*finance_grpc_adapter.GetTaxReportResult, error)
	CreateCommissionPayout(ctx context.Context, params *finance_grpc_adapter.CreateCommissionPayoutParams) (*finance_grpc_adapter.CreateCommissionPayoutResult, error)
	DecideCommissionPayout(ctx context.Context, params *finance_grpc_adapter.DecideCommissionPayoutParams) (*finance_grpc_adapter.DecideCommissionPayoutResult, error)
	GetRealtimeFinancialSummary(ctx context.Context) (*finance_grpc_adapter.GetRealtimeFinancialSummaryResult, error)
	GetCashFlowDashboard(ctx context.Context, params *finance_grpc_adapter.GetCashFlowDashboardParams) (*finance_grpc_adapter.GetCashFlowDashboardResult, error)
	GetAgingAlerts(ctx context.Context, params *finance_grpc_adapter.GetAgingAlertsParams) (*finance_grpc_adapter.GetAgingAlertsResult, error)
	SearchFinanceAuditLog(ctx context.Context, params *finance_grpc_adapter.SearchFinanceAuditLogParams) (*finance_grpc_adapter.SearchFinanceAuditLogResult, error)
}

// Adapters bundles the adapters this service can dispatch through.
// One field per backend; populated at construction time in cmd/start.go.
// Mixed REST + gRPC during the ADR-0009 transition — each backend
// graduates to gRPC-only as its BL-REFACTOR-* card lands.
type Adapters struct {
	iamRest      *iam_rest_adapter.Adapter
	iamGrpc      *iam_grpc_adapter.Adapter
	catalogGrpc  *catalog_grpc_adapter.Adapter
	financeRest  *finance_rest_adapter.Adapter
	financeGrpc  *finance_grpc_adapter.Adapter
	bookingGrpc  *booking_grpc_adapter.Adapter
	crmGrpc      *crm_grpc_adapter.Adapter
	jamaahGrpc   *jamaah_grpc_adapter.Adapter
	opsGrpc       *ops_grpc_adapter.Adapter
	logisticsGrpc *logistics_grpc_adapter.Adapter
	paymentGrpc   *payment_grpc_adapter.Adapter
	visaGrpc      *visa_grpc_adapter.Adapter
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	appName  string
	adapters Adapters
}

// NewServiceParams keeps the constructor readable as the adapter list grows.
type NewServiceParams struct {
	Logger        *zerolog.Logger
	Tracer        trace.Tracer
	AppName       string
	IamRest       *iam_rest_adapter.Adapter
	IamGrpc       *iam_grpc_adapter.Adapter
	CatalogGrpc   *catalog_grpc_adapter.Adapter
	FinanceRest   *finance_rest_adapter.Adapter
	FinanceGrpc   *finance_grpc_adapter.Adapter
	BookingGrpc   *booking_grpc_adapter.Adapter
	CrmGrpc       *crm_grpc_adapter.Adapter
	JamaahGrpc    *jamaah_grpc_adapter.Adapter
	OpsGrpc       *ops_grpc_adapter.Adapter
	LogisticsGrpc *logistics_grpc_adapter.Adapter
	PaymentGrpc   *payment_grpc_adapter.Adapter
	VisaGrpc      *visa_grpc_adapter.Adapter
}

func NewService(p NewServiceParams) IService {
	return &Service{
		logger:  p.Logger,
		tracer:  p.Tracer,
		appName: p.AppName,
		adapters: Adapters{
			iamRest:       p.IamRest,
			iamGrpc:       p.IamGrpc,
			catalogGrpc:   p.CatalogGrpc,
			financeRest:   p.FinanceRest,
			financeGrpc:   p.FinanceGrpc,
			bookingGrpc:   p.BookingGrpc,
			crmGrpc:       p.CrmGrpc,
			jamaahGrpc:    p.JamaahGrpc,
			opsGrpc:       p.OpsGrpc,
			logisticsGrpc: p.LogisticsGrpc,
			paymentGrpc:   p.PaymentGrpc,
			visaGrpc:      p.VisaGrpc,
		},
	}
}
