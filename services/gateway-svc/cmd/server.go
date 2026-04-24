package main

import (
	"fmt"
	"log"
	"os"

	"gateway-svc/api/rest_oapi"
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/monitoring"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// runRestServer runs the REST server using the OpenAPI-generated routes and handler.
//
// Route topology (per ADR 0009 / BL-GTW-001):
//
//   - /system/live, /system/ready      — public (own probes)
//   - /v1/<svc>/system/live            — public probe proxies (remaining REST backends)
//   - /v1/packages*, /v1/package-departures/* — public catalog read (security: [])
//   - /v1/auth/*                       — IAM auth routes (Login/Refresh public; Logout bearer)
//   - /v1/me, /v1/me/2fa/*             — bearer-auth IAM user routes
//   - /v1/users/:id/suspend            — bearer + permission (super_admin)
//   - POST /v1/packages                — staff create package (bearer + catalog.package.manage)
//   - PUT /v1/packages/:id             — staff update package (bearer + catalog.package.manage)
//   - DELETE /v1/packages/:id          — staff delete package (bearer + catalog.package.manage)
//   - POST /v1/packages/:id/departures — staff create departure (bearer + catalog.package.manage)
//   - PUT /v1/departures/:id           — staff update departure (bearer + catalog.package.manage)
//
// S1-E-11: /v1/catalog/system/live removed (catalog-svc is gRPC-only).
// S1-E-12: /v1/iam/system/live + /v1/iam/system/diagnostics/db-tx removed (iam-svc is gRPC-only).
//          IAM client-facing auth routes moved here from iam-svc REST.
// S1-E-03: POST /v1/bookings added (public in S1; proxied to booking-svc gRPC).
// S4-E-02: POST /v1/leads (public), GET+PUT /v1/leads[/:id] (bearer) added.
//
// iamValidator is the *iam_grpc_adapter.Adapter produced in start.go; it is
// passed as the interface type so unit tests can substitute a stub.
func runRestServer(port int, api rest_oapi.ServerInterface, iamValidator middleware.IamValidator, metricsEnabled bool, serviceName string) {
	app := fiber.New()

	// CORS — gateway is the edge, accept cross-origin from any browser client.
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}
	app.Use(cors.New(corsConfig))

	// OpenTelemetry — start an inbound span for every request so cross-service
	// traces initiated by an upstream caller (via W3C traceparent) continue in
	// this process's spans instead of starting a new trace. gateway-svc is the
	// trace origin for edge-initiated requests and a trace continuer when
	// something (e.g. a frontend with OTel instrumentation) already propagates.
	// Span names are prefixed with the service name so a multi-service trace
	// in Tempo is easy to scan (e.g. "gateway-svc GET /v1/iam/system/live").
	app.Use(otelfiber.Middleware(
		otelfiber.WithSpanNameFormatter(func(c *fiber.Ctx) string {
			return serviceName + " " + c.Method() + " " + c.Route().Path
		}),
	))

	if metricsEnabled {
		app.Use(monitoring.RecoveryMiddleware())
		app.Use(monitoring.Middleware())
		app.Get("/metrics", adaptor.HTTPHandler(monitoring.Handler()))
	}

	app.Use(middleware.ErrorHandler())

	wrapper := rest_oapi.ServerInterfaceWrapper{Handler: api}

	// ── Public routes ────────────────────────────────────────────────────────

	// System routes (own probes only — gateway has no DB so no /diagnostics/db-tx).
	system := app.Group("/system")
	{
		system.Get("/live", wrapper.Liveness)
		system.Get("/ready", wrapper.Readiness)
	}

	// /v1 proxy routes — one per backend, dispatched via its adapter.
	v1 := app.Group("/v1")
	{
		// System-probe proxies (REST adapters; retire with each BL-REFACTOR-* card).
		// Note: /v1/catalog/system/live removed in S1-E-11 (catalog-svc is gRPC-only).
		// Note: /v1/iam/system/live + /v1/iam/system/diagnostics/db-tx removed in S1-E-12
		//       (iam-svc is gRPC-only). The seven pure-scaffold backends (booking/crm/
		//       jamaah/logistics/ops/payment/visa) retired their REST surfaces in S1-E-13.
		v1.Get("/finance/system/live", wrapper.GetFinanceSystemLive)

		// Public catalog read (BL-GTW-002 / S1-E-10) — gRPC adapter.
		// security: [] in openapi.yaml — no bearer required.
		v1.Get("/packages", wrapper.ListPackages)
		v1.Get("/packages/:id", wrapper.GetPackageById)
		v1.Get("/package-departures/:id", wrapper.GetPackageDepartureById)

		// Booking draft (BL-GTW-003 / S1-E-03) — public in S1; auth arrives with F4.
		// Proxied to booking-svc.BookingService/CreateDraftBooking via gRPC.
		v1.Post("/bookings", wrapper.CreateDraftBooking)

		// CRM lead capture (S4-E-02 / BL-CRM-001) — public.
		// POST /v1/leads — submitting a lead from a landing page or B2C form.
		v1.Post("/leads", wrapper.CreateLead)

		// S2 webhook routes (BL-PAY-003/004 / ISSUE-007/008) — public, signature-protected.
		// payment-svc performs signature verification; no bearer required at gateway.
		v1.Post("/webhooks/midtrans", wrapper.WebhookMidtrans)
		v1.Post("/webhooks/xendit", wrapper.WebhookXendit)
		// Mock trigger — dev only (always registered; payment-svc skips signature for "mock").
		v1.Post("/webhooks/mock/trigger", wrapper.WebhookMockTrigger)
	}

	// IAM auth public routes (BL-IAM-018 / S1-E-12) — no bearer required.
	// Login + refresh are public; they issue / rotate the bearer.
	auth := v1.Group("/auth")
	{
		auth.Post("/login", wrapper.Login)
		auth.Post("/refresh", wrapper.RefreshSession)
	}

	// ── Protected routes (BL-GTW-001 / S1-E-09) ─────────────────────────────
	//
	// Every route in this group requires a valid Bearer token. The
	// RequireBearerToken middleware calls iam-svc.ValidateToken (gRPC) and
	// injects *middleware.Identity into c.Locals(middleware.IdentityKey) for
	// downstream handlers to consume. Failure modes per F1-W7:
	//
	//   - Missing / malformed Authorization header → 401 UNAUTHORIZED
	//   - iam-svc rejects (expired, revoked)       → 401 UNAUTHORIZED
	//   - iam-svc unreachable / timeout             → 502 SERVICE_UNAVAILABLE
	v1Protected := app.Group("/v1", middleware.RequireBearerToken(iamValidator))
	{
		// IAM auth bearer routes (BL-IAM-018 / S1-E-12).
		v1Protected.Delete("/auth/logout", wrapper.Logout)
		v1Protected.Get("/me", wrapper.GetMe)
		v1Protected.Post("/me/2fa/enroll", wrapper.EnrollTOTP)
		v1Protected.Post("/me/2fa/verify", wrapper.VerifyTOTP)
		v1Protected.Post("/users/:id/suspend", wrapper.SuspendUser)

		// Staff catalog write routes (BL-CAT-014 / S1-E-07).
		// All require bearer + catalog.package.manage permission.
		// Response shapes align with public read models (PackageDetail, DepartureDetail).
		v1Protected.Post("/packages", wrapper.CreatePackage)
		v1Protected.Put("/packages/:id", wrapper.UpdatePackageById)
		v1Protected.Delete("/packages/:id", wrapper.DeletePackageById)
		v1Protected.Post("/packages/:id/departures", wrapper.CreateDeparture)
		v1Protected.Put("/departures/:id", wrapper.UpdateDepartureById)

		// S2 booking submit (BL-BOOK-005 / ISSUE-006) — bearer required.
		// CS/admin action: transitions draft → pending_payment.
		v1Protected.Post("/bookings/:id/submit", wrapper.SubmitBooking)

		// CRM lead management (S4-E-02 / BL-CRM-002..003) — bearer required.
		// GET  /v1/leads     — list leads (cs/admin)
		// GET  /v1/leads/:id — get single lead
		// PUT  /v1/leads/:id — update lead (status, notes, assigned_cs_id)
		v1Protected.Get("/leads", wrapper.ListLeads)
		v1Protected.Get("/leads/:id", wrapper.GetLeadByID)
		v1Protected.Put("/leads/:id", wrapper.UpdateLeadByID)

		// Finance report routes (S5-E-01 / BL-FIN-004..005) — bearer required.
		// GET /v1/finance/summary  — aggregate per-account balances
		// GET /v1/finance/journals — paginated journal entries + lines
		v1Protected.Get("/finance/summary", wrapper.GetFinanceSummary)
		v1Protected.Get("/finance/journals", wrapper.ListJournals)

		// Finance depth routes (Phase 6 / Wave 1B) — bearer required.
		v1Protected.Post("/finance/recognize-revenue", wrapper.RecognizeRevenue)
		v1Protected.Get("/finance/pl", wrapper.GetPLReport)
		v1Protected.Get("/finance/balance-sheet", wrapper.GetBalanceSheet)

		// Catalog master data (Phase 6 / Wave 1A) — bearer required.
		v1Protected.Get("/masters/hotels", wrapper.ListHotels)
		v1Protected.Post("/masters/hotels", wrapper.CreateHotel)
		v1Protected.Put("/masters/hotels/:id", wrapper.UpdateHotel)
		v1Protected.Delete("/masters/hotels/:id", wrapper.DeleteHotel)
		v1Protected.Get("/masters/airlines", wrapper.ListAirlines)
		v1Protected.Post("/masters/airlines", wrapper.CreateAirline)
		v1Protected.Put("/masters/airlines/:id", wrapper.UpdateAirline)
		v1Protected.Delete("/masters/airlines/:id", wrapper.DeleteAirline)
		v1Protected.Get("/masters/muthawwif", wrapper.ListMuthawwif)
		v1Protected.Post("/masters/muthawwif", wrapper.CreateMuthawwif)
		v1Protected.Put("/masters/muthawwif/:id", wrapper.UpdateMuthawwif)
		v1Protected.Delete("/masters/muthawwif/:id", wrapper.DeleteMuthawwif)
		v1Protected.Get("/masters/addons", wrapper.ListAddons)
		v1Protected.Post("/masters/addons", wrapper.CreateAddon)
		v1Protected.Put("/masters/addons/:id", wrapper.UpdateAddon)
		v1Protected.Delete("/masters/addons/:id", wrapper.DeleteAddon)
		v1Protected.Put("/departures/:id/pricing", wrapper.SetDeparturePricing)
		v1Protected.Get("/departures/:id/pricing", wrapper.GetDeparturePricing)

		// IAM admin routes (Phase 6 / Wave 1C) — bearer required.
		v1Protected.Get("/admin/users", wrapper.ListUsers)
		v1Protected.Post("/admin/users", wrapper.CreateUser)
		v1Protected.Get("/admin/users/:id", wrapper.GetUser)
		v1Protected.Put("/admin/users/:id", wrapper.UpdateUser)
		v1Protected.Post("/admin/users/:id/reset-password", wrapper.ResetUserPassword)
		v1Protected.Get("/admin/roles", wrapper.ListRoles)
		v1Protected.Post("/admin/roles", wrapper.CreateRole)
		v1Protected.Put("/admin/roles/:id", wrapper.UpdateRole)
		v1Protected.Delete("/admin/roles/:id", wrapper.DeleteRole)
		v1Protected.Get("/admin/permissions", wrapper.ListPermissions)
		v1Protected.Post("/admin/users/:id/roles/:role_id", wrapper.AssignRoleToUser)
		v1Protected.Delete("/admin/users/:id/roles/:role_id", wrapper.RevokeRoleFromUser)

		// Jamaah manifest (Phase 6 / Wave 1A) — bearer required.
		v1Protected.Get("/manifest/:departure_id", wrapper.GetDepartureManifest)

		// S3 ops routes (S3 Wave 2) — bearer required.
		v1Protected.Post("/ops/room-allocation", wrapper.RunRoomAllocation)
		v1Protected.Get("/ops/room-allocation/:departure_id", wrapper.GetRoomAllocation)
		v1Protected.Post("/ops/id-cards", wrapper.GenerateIDCard)
		v1Protected.Post("/ops/id-cards/verify", wrapper.VerifyIDCard)
		v1Protected.Get("/ops/manifest/:departure_id/export", wrapper.ExportManifest)

		// S3 logistics routes (S3 Wave 2 / ISSUE-018) — bearer required.
		// GET /v1/fulfillment-tasks — list all tasks, optional ?status=&departure_id= filters.
		v1Protected.Get("/fulfillment-tasks", wrapper.ListFulfillmentTasks)
		v1Protected.Post("/logistics/ship", wrapper.ShipFulfillmentTask)
		v1Protected.Post("/logistics/pickup-qr", wrapper.GeneratePickupQR)
		v1Protected.Post("/logistics/pickup-qr/redeem", wrapper.RedeemPickupQR)

		// S3 finance GRN (S3 Wave 2) — bearer required.
		v1Protected.Post("/finance/grn", wrapper.OnGRNReceived)

		// S3 jamaah OCR (S3 Wave 2) — bearer required.
		v1Protected.Post("/documents/:id/ocr", wrapper.TriggerDocumentOCR)
		v1Protected.Get("/documents/:id/ocr", wrapper.GetDocumentOCRStatus)

		// S2 invoice routes (BL-PAY-001 / ISSUE-005) — bearer required.
		v1Protected.Post("/invoices", wrapper.CreateInvoice)
		v1Protected.Get("/invoices/:id", wrapper.GetInvoiceByID)
		v1Protected.Post("/invoices/:id/virtual-accounts", wrapper.IssueVirtualAccountForInvoice)

		// S2 payment link (BL-PAY-020) — bearer required.
		// CS closing tool: re-issue VA for an existing booking.
		v1Protected.Post("/payments/link", wrapper.ReissuePaymentLink)

		// S5 finance correction (BL-FIN-006) — bearer required.
		// POST /v1/finance/journals/:id/correct — reverse a journal entry.
		v1Protected.Post("/finance/journals/:id/correct", wrapper.CorrectJournal)

		// Phase 6 IAM security routes (BL-IAM-007/011/014/016) — bearer required.
		v1Protected.Put("/admin/users/:id/data-scope", wrapper.SetDataScope)
		v1Protected.Post("/admin/api-keys", wrapper.CreateAPIKey)
		v1Protected.Delete("/admin/api-keys/:id", wrapper.RevokeAPIKey)
		v1Protected.Get("/admin/config", wrapper.GetGlobalConfig)
		v1Protected.Put("/admin/config/:key", wrapper.SetGlobalConfig)
		v1Protected.Get("/admin/activity-log", wrapper.SearchActivityLog)

		// IAM security depth (BL-IAM-010/012/013/015/017) — bearer required.
		v1Protected.Get("/admin/security/password-policy", wrapper.GetPasswordPolicy)
		v1Protected.Put("/admin/security/password-policy", wrapper.SetPasswordPolicy)
		v1Protected.Post("/admin/security/anomalies", wrapper.RecordLoginAnomaly)
		v1Protected.Get("/admin/security/sessions", wrapper.ListSessions)
		v1Protected.Delete("/admin/security/sessions/:id", wrapper.RevokeSession)
		v1Protected.Get("/admin/comm-templates", wrapper.ListCommTemplates)
		v1Protected.Put("/admin/comm-templates/:channel/:name", wrapper.UpsertCommTemplate)
		v1Protected.Get("/admin/backups", wrapper.GetBackupHistory)
		v1Protected.Post("/admin/backups", wrapper.TriggerBackup)

		// Phase 6 Finance disbursement + aging (BL-FIN-010/011) — bearer required.
		v1Protected.Post("/finance/disbursements", wrapper.CreateDisbursementBatch)
		v1Protected.Put("/finance/disbursements/:id/decision", wrapper.ApproveDisbursement)
		v1Protected.Get("/finance/aging", wrapper.GetARAPAging)

		// Phase 6 Ops scanning + bus boarding (BL-OPS-010/011) — bearer required.
		v1Protected.Post("/ops/scans", wrapper.RecordScan)
		v1Protected.Post("/ops/bus-boarding", wrapper.RecordBusBoarding)
		v1Protected.Get("/ops/bus-boarding/:departure_id", wrapper.GetBoardingRoster)

		// Phase 6 Logistics procurement + GRN + kit (BL-LOG-010..012) — bearer required.
		v1Protected.Post("/logistics/purchase-requests", wrapper.CreatePurchaseRequest)
		v1Protected.Put("/logistics/purchase-requests/:id/decision", wrapper.ApprovePurchaseRequest)
		v1Protected.Post("/logistics/grn-qc", wrapper.RecordGRNWithQC)
		v1Protected.Post("/logistics/kit-assembly", wrapper.CreateKitAssembly)

		// Phase 6 Visa pipeline (BL-VISA-001..003) — bearer required.
		v1Protected.Put("/visas/:id/status", wrapper.TransitionVisaStatus)
		v1Protected.Post("/visas/bulk-submit", wrapper.BulkSubmitVisa)
		v1Protected.Get("/visas", wrapper.GetVisaApplications)

		// Vendor readiness (BL-OPS-020) — bearer required.
		v1Protected.Get("/departures/:id/readiness", wrapper.GetDepartureReadiness)
		v1Protected.Put("/departures/:id/readiness", wrapper.UpdateDepartureReadiness)

		// Catalog depth — Wave 3 (BL-CAT-010/011/013) — bearer required.
		v1Protected.Post("/admin/catalog/packages/bulk-import", wrapper.BulkImportPackages)
		v1Protected.Post("/admin/catalog/packages/bulk-update", wrapper.BulkUpdatePackages)
		v1Protected.Get("/admin/catalog/packages/:id/version", wrapper.GetPackageVersion)

		// Booking depth — Wave 3 (BL-BOOK-007) — bearer required.
		v1Protected.Get("/bookings/departures/:id/seats-by-channel", wrapper.GetSeatsByChannel)

		// Wave 4 Finance depth (BL-FIN-020..041) — bearer required.
		v1Protected.Post("/finance/billing/schedule", wrapper.ScheduleBilling)
		v1Protected.Post("/finance/bank-transactions", wrapper.RecordBankTransaction)
		v1Protected.Get("/finance/bank-reconciliation", wrapper.GetBankReconciliation)
		v1Protected.Get("/finance/ar-subledger", wrapper.GetARSubledger)
		v1Protected.Post("/finance/receipts", wrapper.IssueDigitalReceipt)
		v1Protected.Get("/finance/receipts/:id", wrapper.GetDigitalReceipt)
		v1Protected.Post("/finance/manual-payments", wrapper.RecordManualPayment)
		v1Protected.Post("/finance/vendors", wrapper.CreateFinanceVendor)
		v1Protected.Put("/finance/vendors/:id", wrapper.UpdateFinanceVendor)
		v1Protected.Get("/finance/vendors", wrapper.ListFinanceVendors)
		v1Protected.Delete("/finance/vendors/:id", wrapper.DeleteFinanceVendor)
		v1Protected.Get("/finance/ap-subledger", wrapper.GetAPSubledger)
		v1Protected.Get("/finance/payment-authorizations", wrapper.ListPendingAuthorizations)
		v1Protected.Put("/finance/payment-authorizations/:id/decision", wrapper.DecidePaymentAuthorization)
		v1Protected.Post("/finance/petty-cash", wrapper.RecordPettyCash)
		v1Protected.Post("/finance/petty-cash/close-period", wrapper.ClosePettyCashPeriod)
		v1Protected.Get("/finance/project-costs/:departure_id", wrapper.GetProjectCosts)
		v1Protected.Get("/finance/departure-pl/:departure_id", wrapper.GetDeparturePL)
		v1Protected.Get("/finance/budget-vs-actual", wrapper.GetBudgetVsActual)
		v1Protected.Post("/finance/auto-journal", wrapper.TriggerAutoJournal)
		v1Protected.Get("/finance/revenue-recognition-policy", wrapper.GetRevenueRecognitionPolicy)
		v1Protected.Put("/finance/revenue-recognition-policy", wrapper.SetRevenueRecognitionPolicy)
		v1Protected.Post("/finance/exchange-rates", wrapper.SetExchangeRate)
		v1Protected.Get("/finance/exchange-rates", wrapper.GetExchangeRate)
		v1Protected.Post("/finance/fixed-assets", wrapper.CreateFixedAsset)
		v1Protected.Get("/finance/fixed-assets", wrapper.ListFixedAssets)
		v1Protected.Post("/finance/depreciation", wrapper.RunDepreciation)
		v1Protected.Post("/finance/tax/calculate", wrapper.CalculateTax)
		v1Protected.Get("/finance/tax/report", wrapper.GetTaxReport)
		v1Protected.Post("/finance/commission-payouts", wrapper.CreateCommissionPayout)
		v1Protected.Put("/finance/commission-payouts/:id/decision", wrapper.DecideCommissionPayout)
		v1Protected.Get("/finance/realtime-summary", wrapper.GetRealtimeFinancialSummary)
		v1Protected.Get("/finance/cashflow", wrapper.GetCashFlowDashboard)
		v1Protected.Get("/finance/aging-alerts", wrapper.GetAgingAlerts)
		v1Protected.Get("/finance/audit-log", wrapper.SearchFinanceAuditLog)

		// Wave 5 Ops depth (BL-OPS-021..042) — bearer required.
		v1Protected.Post("/ops/collective-docs", wrapper.StoreCollectiveDocument)
		v1Protected.Get("/ops/collective-docs", wrapper.GetCollectiveDocuments)
		v1Protected.Put("/ops/collective-docs/:id/acl", wrapper.SetDocumentACL)
		v1Protected.Post("/ops/passport-ocr", wrapper.ExtractPassportOCR)
		v1Protected.Post("/ops/mahram", wrapper.SetMahramRelation)
		v1Protected.Get("/ops/mahram", wrapper.GetMahramRelations)
		v1Protected.Get("/ops/document-progress", wrapper.GetDocumentProgress)
		v1Protected.Get("/ops/expiry-alerts", wrapper.GetExpiryAlerts)
		v1Protected.Post("/ops/official-letters", wrapper.GenerateOfficialLetter)
		v1Protected.Post("/ops/immigration-manifest", wrapper.GenerateImmigrationManifest)
		v1Protected.Post("/ops/transport-assignments", wrapper.AssignTransport)
		v1Protected.Get("/ops/transport-assignments", wrapper.GetTransportAssignments)
		v1Protected.Post("/ops/manifest-delta", wrapper.PublishManifestDelta)
		v1Protected.Post("/ops/staff-assignments", wrapper.AssignStaff)
		v1Protected.Post("/ops/passport-handover", wrapper.RecordPassportHandover)
		v1Protected.Get("/ops/passport-log", wrapper.GetPassportLog)
		v1Protected.Get("/ops/visa-progress", wrapper.GetVisaProgress)
		v1Protected.Post("/ops/e-visa", wrapper.StoreEVisa)
		v1Protected.Get("/ops/e-visa/:pilgrim_id", wrapper.GetEVisa)
		v1Protected.Post("/ops/external-provider", wrapper.TriggerExternalProvider)
		v1Protected.Post("/ops/refunds", wrapper.CreateRefund)
		v1Protected.Put("/ops/refunds/:id/approve", wrapper.ApproveRefund)
		v1Protected.Post("/ops/penalties", wrapper.RecordPenalty)
		v1Protected.Post("/ops/luggage-scan", wrapper.RecordLuggageScan)
		v1Protected.Get("/ops/luggage-count", wrapper.GetLuggageCount)
		v1Protected.Post("/ops/broadcast", wrapper.BroadcastSchedule)
		v1Protected.Post("/ops/tasreh", wrapper.IssueDigitalTasreh)
		v1Protected.Post("/ops/raudhah", wrapper.RecordRaudhahEntry)
		v1Protected.Post("/ops/audio-devices", wrapper.RegisterAudioDevice)
		v1Protected.Put("/ops/audio-devices/:id/status", wrapper.UpdateAudioDeviceStatus)
		v1Protected.Post("/ops/zamzam", wrapper.RecordZamzamDistribution)
		v1Protected.Post("/ops/room-checkin", wrapper.RecordRoomCheckIn)

		// Wave 5 Logistics depth (BL-LOG-013..029) — bearer required.
		v1Protected.Get("/logistics/purchase-requests", wrapper.ListPurchaseRequests)
		v1Protected.Get("/logistics/budget-sync/:departure_id", wrapper.GetBudgetSyncStatus)
		v1Protected.Get("/logistics/tiered-approvals", wrapper.GetTieredApprovals)
		v1Protected.Put("/logistics/tiered-approvals/:id/decision", wrapper.DecideTieredApproval)
		v1Protected.Post("/logistics/auto-vendor", wrapper.AutoSelectVendor)
		v1Protected.Post("/logistics/partial-grn", wrapper.RecordPartialGRN)
		v1Protected.Post("/logistics/reverse-grn", wrapper.ReverseGRN)
		v1Protected.Post("/logistics/barcode", wrapper.GenerateBarcode)
		v1Protected.Post("/logistics/sku-labels", wrapper.PrintSKULabel)
		v1Protected.Post("/logistics/warehouses", wrapper.CreateWarehouse)
		v1Protected.Post("/logistics/stock-transfer", wrapper.TransferStock)
		v1Protected.Get("/logistics/stock-alerts", wrapper.GetStockAlerts)
		v1Protected.Put("/logistics/reorder-levels", wrapper.SetReorderLevel)
		v1Protected.Post("/logistics/stocktake", wrapper.StartStocktake)
		v1Protected.Post("/logistics/stocktake/:id/count", wrapper.RecordStocktakeCount)
		v1Protected.Put("/logistics/stocktake/:id/finalize", wrapper.FinalizeStocktake)
		v1Protected.Post("/logistics/size-sync", wrapper.SyncFulfillmentSizes)
		v1Protected.Post("/logistics/courier-tracking", wrapper.RecordCourierTracking)
		v1Protected.Post("/logistics/returns", wrapper.RecordReturn)
		v1Protected.Post("/logistics/exchanges", wrapper.ProcessExchange)

		// Wave 7 CRM depth — agency platform (BL-CRM-010..012, 017..066) — bearer required.
		// Agency registration
		v1Protected.Post("/crm/agents", wrapper.RegisterAgent)
		v1Protected.Post("/crm/agents/kyc", wrapper.SubmitAgentKyc)
		v1Protected.Post("/crm/agents/mou", wrapper.SignAgentMoU)
		v1Protected.Get("/crm/agents/:agent_id", wrapper.GetAgentProfile)
		v1Protected.Get("/crm/agents/:agent_id/replica-site", wrapper.GetReplicaSite)
		v1Protected.Put("/crm/agents/:agent_id/replica-site", wrapper.UpdateReplicaSite)
		// Content
		v1Protected.Post("/crm/content/share-link", wrapper.GetSocialShareLink)
		v1Protected.Post("/crm/content/business-card", wrapper.GenerateBusinessCard)
		v1Protected.Get("/crm/content/bank", wrapper.ListContentBank)
		v1Protected.Post("/crm/content/bank", wrapper.CreateContentAsset)
		v1Protected.Post("/crm/content/watermark", wrapper.WatermarkFlyer)
		v1Protected.Get("/crm/content/gallery", wrapper.ListProgramGallery)
		v1Protected.Post("/crm/content/tracking-code", wrapper.SetTrackingCode)
		v1Protected.Get("/crm/agents/:agent_id/ads-stats", wrapper.GetAdsManagerStats)
		v1Protected.Post("/crm/content/utm-links", wrapper.CreateUtmLink)
		v1Protected.Get("/crm/agents/:agent_id/utm-links", wrapper.ListUtmLinks)
		v1Protected.Post("/crm/content/landing-pages", wrapper.CreateLandingPage)
		v1Protected.Get("/crm/agents/:agent_id/landing-pages", wrapper.ListLandingPages)
		v1Protected.Post("/crm/content/schedule", wrapper.ScheduleContent)
		v1Protected.Get("/crm/agents/:agent_id/scheduled-content", wrapper.ListScheduledContent)
		v1Protected.Get("/crm/agents/:agent_id/content-analytics", wrapper.GetContentAnalytics)
		// CRM / Leads
		v1Protected.Post("/crm/agent-leads", wrapper.CreateAgentLead)
		v1Protected.Get("/crm/agents/:agent_id/leads", wrapper.ListAgentLeads)
		v1Protected.Post("/crm/leads/reminders", wrapper.SetLeadReminder)
		v1Protected.Get("/crm/leads/:lead_id/reminders", wrapper.ListLeadReminders)
		v1Protected.Get("/crm/agents/:agent_id/bot-leads", wrapper.FilterBotLeads)
		v1Protected.Post("/crm/drip-sequences", wrapper.CreateDripSequence)
		v1Protected.Get("/crm/agents/:agent_id/drip-sequences", wrapper.ListDripSequences)
		v1Protected.Post("/crm/moment-triggers", wrapper.CreateMomentTrigger)
		v1Protected.Post("/crm/segments", wrapper.CreateSegment)
		v1Protected.Get("/crm/agents/:agent_id/segments", wrapper.ListSegments)
		v1Protected.Post("/crm/broadcasts", wrapper.SendBroadcast)
		v1Protected.Post("/crm/leads/assign-fair", wrapper.AssignLeadFair)
		v1Protected.Post("/crm/sla-rules", wrapper.CreateSlaRule)
		v1Protected.Get("/crm/leads/:lead_id/trail", wrapper.GetLeadTrail)
		v1Protected.Post("/crm/leads/tag", wrapper.TagLead)
		v1Protected.Post("/crm/quotes", wrapper.GenerateQuote)
		v1Protected.Get("/crm/quotes/:quote_id", wrapper.GetQuote)
		v1Protected.Post("/crm/payment-links", wrapper.BuildPaymentLink)
		v1Protected.Post("/crm/discounts", wrapper.RequestDiscount)
		v1Protected.Put("/crm/discounts/:discount_id/decision", wrapper.ApproveDiscount)
		v1Protected.Get("/crm/agents/:agent_id/stale-prospects", wrapper.GetStaleProspects)
		// Commission
		v1Protected.Get("/crm/agents/:agent_id/commission-balance", wrapper.GetCommissionBalance)
		v1Protected.Get("/crm/agents/:agent_id/commission-events", wrapper.GetCommissionEvents)
		v1Protected.Post("/crm/payouts", wrapper.RequestPayout)
		v1Protected.Get("/crm/agents/:agent_id/payout-history", wrapper.GetPayoutHistory)
		v1Protected.Post("/crm/override-commission", wrapper.ComputeOverrideCommission)
		v1Protected.Get("/crm/agents/:agent_id/roas", wrapper.GetRoasReport)
		// Academy
		v1Protected.Get("/crm/academy/courses", wrapper.ListAcademyCourses)
		v1Protected.Get("/crm/agents/:agent_id/course-progress", wrapper.GetCourseProgress)
		v1Protected.Post("/crm/academy/quiz", wrapper.SubmitQuiz)
		v1Protected.Get("/crm/academy/sales-scripts", wrapper.ListSalesScripts)
		v1Protected.Get("/crm/leaderboard", wrapper.GetLeaderboard)
		v1Protected.Get("/crm/agents/:agent_id/tier", wrapper.GetAgentTier)
		// Alumni / Other
		v1Protected.Get("/crm/agents/:agent_id/referrals", wrapper.GetAlumniReferrals)
		v1Protected.Get("/crm/agents/:agent_id/return-intent", wrapper.GetReturnIntentSavings)
		v1Protected.Post("/crm/zakat/calculate", wrapper.CalculateZakat)
		v1Protected.Post("/crm/charity", wrapper.RecordCharity)
	}

	go func() {
		log.Printf("rest server started successfully on :%d", port)

		err := app.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Printf("failed to listen at port: %v!\n", port)
			log.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}()
}
