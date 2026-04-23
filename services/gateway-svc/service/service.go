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
		},
	}
}
