package service

import (
	"context"

	"crm-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for crm-svc.
//
// Scaffold methods:
//   - Liveness / Readiness / DbTxDiagnostic
//
// S4-E-02 adds CRM lead management:
//   - CreateLead, GetLead, UpdateLead, ListLeads
//   - OnBookingCreated, OnBookingPaidInFull (event callbacks)
//
// Wave 7 adds agency platform depth (BL-CRM-010..012, BL-CRM-017..066).
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// S4-E-02 — CRM lead management (BL-CRM-001..003)
	CreateLead(ctx context.Context, params *CreateLeadParams) (*LeadResult, error)
	GetLead(ctx context.Context, params *GetLeadParams) (*LeadResult, error)
	UpdateLead(ctx context.Context, params *UpdateLeadParams) (*LeadResult, error)
	ListLeads(ctx context.Context, params *ListLeadsParams) (*ListLeadsResult, error)
	OnBookingCreated(ctx context.Context, params *OnBookingCreatedParams) (*OnBookingCreatedResult, error)
	OnBookingPaidInFull(ctx context.Context, params *OnBookingPaidInFullParams) (*OnBookingPaidInFullResult, error)

	// Wave 7 — Agency platform depth (BL-CRM-010..012, BL-CRM-017..066)
	RegisterAgent(ctx context.Context, p *RegisterAgentParams) (*AgentData, error)
	SubmitAgentKyc(ctx context.Context, p *SubmitAgentKycParams) (*AgentData, error)
	SignAgentMoU(ctx context.Context, p *SignAgentMoUParams) (*AgentData, error)
	GetAgentProfile(ctx context.Context, p *AgentIdParams) (*AgentData, error)
	GetReplicaSite(ctx context.Context, p *AgentIdParams) (*ReplicaSiteData, error)
	UpdateReplicaSite(ctx context.Context, p *UpdateReplicaSiteParams) (*ReplicaSiteData, error)
	GetSocialShareLink(ctx context.Context, p *SocialShareParams) (*SocialShareData, error)
	GenerateBusinessCard(ctx context.Context, p *BusinessCardParams) (*BusinessCardData, error)
	ListContentBank(ctx context.Context, p *ContentBankListParams) (*ContentBankListData, error)
	CreateContentAsset(ctx context.Context, p *CreateContentAssetParams) (*ContentAssetData, error)
	WatermarkFlyer(ctx context.Context, p *WatermarkFlyerParams) (*WatermarkFlyerData, error)
	ListProgramGallery(ctx context.Context, p *ProgramGalleryParams) (*ProgramGalleryData, error)
	SetTrackingCode(ctx context.Context, p *TrackingCodeParams) (*TrackingCodeData, error)
	GetAdsManagerStats(ctx context.Context, p *AdsManagerParams) (*AdsManagerData, error)
	CreateUtmLink(ctx context.Context, p *CreateUtmLinkParams) (*UtmLinkData, error)
	ListUtmLinks(ctx context.Context, p *ListUtmLinksParams) (*UtmLinkListData, error)
	CreateLandingPage(ctx context.Context, p *CreateLandingPageParams) (*LandingPageData, error)
	ListLandingPages(ctx context.Context, p *ListLandingPagesParams) (*LandingPageListData, error)
	ScheduleContent(ctx context.Context, p *ScheduleContentParams) (*ScheduledContentData, error)
	ListScheduledContent(ctx context.Context, p *ListScheduledContentParams) (*ScheduledContentListData, error)
	GetContentAnalytics(ctx context.Context, p *ContentAnalyticsParams) (*ContentAnalyticsData, error)
	CreateAgentLead(ctx context.Context, p *CreateAgentLeadParams) (*AgentLeadData, error)
	ListAgentLeads(ctx context.Context, p *ListAgentLeadsParams) (*AgentLeadListData, error)
	SetLeadReminder(ctx context.Context, p *SetLeadReminderParams) (*LeadReminderData, error)
	ListLeadReminders(ctx context.Context, p *ListLeadRemindersParams) (*LeadReminderListData, error)
	FilterBotLeads(ctx context.Context, p *BotFilterParams) (*BotFilterData, error)
	CreateDripSequence(ctx context.Context, p *CreateDripSequenceParams) (*DripSequenceData, error)
	ListDripSequences(ctx context.Context, p *ListDripSequencesParams) (*DripSequenceListData, error)
	CreateMomentTrigger(ctx context.Context, p *CreateMomentTriggerParams) (*MomentTriggerData, error)
	CreateSegment(ctx context.Context, p *CreateSegmentParams) (*SegmentData, error)
	ListSegments(ctx context.Context, p *ListSegmentsParams) (*SegmentListData, error)
	SendBroadcast(ctx context.Context, p *SendBroadcastParams) (*BroadcastData, error)
	AssignLeadFair(ctx context.Context, p *AssignLeadFairParams) (*AssignLeadFairData, error)
	CreateSlaRule(ctx context.Context, p *CreateSlaRuleParams) (*SlaRuleData, error)
	GetLeadTrail(ctx context.Context, p *GetLeadTrailParams) (*LeadTrailData, error)
	TagLead(ctx context.Context, p *TagLeadParams) (*TagLeadData, error)
	GenerateQuote(ctx context.Context, p *GenerateQuoteParams) (*QuoteData, error)
	GetQuote(ctx context.Context, p *GetQuoteParams) (*QuoteData, error)
	BuildPaymentLink(ctx context.Context, p *BuildPaymentLinkParams) (*PaymentLinkData, error)
	RequestDiscount(ctx context.Context, p *RequestDiscountParams) (*DiscountApprovalData, error)
	ApproveDiscount(ctx context.Context, p *ApproveDiscountParams) (*DiscountApprovalData, error)
	GetStaleProspects(ctx context.Context, p *StaleProspectParams) (*StaleProspectListData, error)
	GetCommissionBalance(ctx context.Context, p *AgentIdParams) (*CommissionBalanceData, error)
	GetCommissionEvents(ctx context.Context, p *CommissionEventListParams) (*CommissionEventListData, error)
	RequestPayout(ctx context.Context, p *PayoutParams) (*PayoutData, error)
	GetPayoutHistory(ctx context.Context, p *AgentIdParams) (*PayoutHistoryData, error)
	ComputeOverrideCommission(ctx context.Context, p *OverrideCommissionParams) (*OverrideCommissionData, error)
	GetRoasReport(ctx context.Context, p *RoasReportParams) (*RoasReportData, error)
	ListAcademyCourses(ctx context.Context, p *AcademyListParams) (*AcademyCourseListData, error)
	GetCourseProgress(ctx context.Context, p *CourseProgressParams) (*CourseProgressData, error)
	SubmitQuiz(ctx context.Context, p *SubmitQuizParams) (*QuizData, error)
	ListSalesScripts(ctx context.Context, p *SalesScriptListParams) (*SalesScriptListData, error)
	GetLeaderboard(ctx context.Context, p *LeaderboardParams) (*LeaderboardData, error)
	GetAgentTier(ctx context.Context, p *AgentIdParams) (*AgentTierData, error)
	GetAlumniReferrals(ctx context.Context, p *AgentIdParams) (*AlumniReferralListData, error)
	GetReturnIntentSavings(ctx context.Context, p *AgentIdParams) (*ReturnIntentData, error)
	CalculateZakat(ctx context.Context, p *ZakatCalculatorParams) (*ZakatData, error)
	RecordCharity(ctx context.Context, p *CharityParams) (*CharityData, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
) IService {
	return &Service{
		logger:  logger,
		tracer:  tracer,
		appName: appName,
		store:   store,
	}
}
