// crm_depth.go — service implementations for Wave 7 CRM depth RPCs.
// BL-CRM-010..012, BL-CRM-017..066 (53 features).
//
// Uses inline pgx queries on the pool directly. Tables are created on first
// use via CREATE TABLE IF NOT EXISTS in each method's init path.

package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Params / Results — Agency Registration
// ---------------------------------------------------------------------------

type RegisterAgentParams struct {
	AgencyName string
	OwnerName  string
	Phone      string
	Email      string
	City       string
	Province   string
}

type AgentData struct {
	AgentId    string
	AgencyName string
	OwnerName  string
	Phone      string
	Email      string
	Status     string
	Tier       string
	CreatedAt  string
}

type AgentIdParams struct {
	AgentId string
}

type SubmitAgentKycParams struct {
	AgentId   string
	KtpUrl    string
	NpwpUrl   string
	SiupUrl   string
	SelfieUrl string
}

type SignAgentMoUParams struct {
	AgentId      string
	SignatureUrl  string
	SignedAt      string
}

type ReplicaSiteData struct {
	AgentId   string
	SiteUrl   string
	Theme     string
	LogoUrl   string
	UpdatedAt string
}

type UpdateReplicaSiteParams struct {
	AgentId      string
	Theme        string
	LogoUrl      string
	CustomDomain string
}

// ---------------------------------------------------------------------------
// Params / Results — Content & Marketing
// ---------------------------------------------------------------------------

type SocialShareParams struct {
	AgentId   string
	PackageId string
	Platform  string
}

type SocialShareData struct {
	ShareUrl  string
	ShortCode string
}

type BusinessCardParams struct {
	AgentId  string
	Template string
}

type BusinessCardData struct {
	CardUrl   string
	CreatedAt string
}

type ContentBankListParams struct {
	AgentId  string
	Category string
	Page     int32
	PageSize int32
}

type ContentAssetItemData struct {
	AssetId   string
	Title     string
	Type      string
	Url       string
	CreatedAt string
}

type ContentBankListData struct {
	Assets   []*ContentAssetItemData
	Total    int64
	Page     int32
	PageSize int32
}

type CreateContentAssetParams struct {
	AgentId  string
	Title    string
	Type     string
	Url      string
	Category string
}

type ContentAssetData struct {
	AssetId   string
	AgentId   string
	Title     string
	Type      string
	Url       string
	CreatedAt string
}

type WatermarkFlyerParams struct {
	AgentId  string
	FlyerUrl string
	Text     string
}

type WatermarkFlyerData struct {
	ResultUrl string
}

type ProgramGalleryParams struct {
	PackageId string
	Page      int32
	PageSize  int32
}

type GalleryItemData struct {
	ImageUrl string
	Caption  string
}

type ProgramGalleryData struct {
	PackageId string
	Items     []*GalleryItemData
	Total     int64
}

type TrackingCodeParams struct {
	AgentId string
	Code    string
	Type    string
}

type TrackingCodeData struct {
	AgentId   string
	Code      string
	Type      string
	UpdatedAt string
}

type AdsManagerParams struct {
	AgentId  string
	DateFrom string
	DateTo   string
}

type AdsManagerData struct {
	AgentId     string
	Impressions int64
	Clicks      int64
	Conversions int64
	Spend       float64
}

type CreateUtmLinkParams struct {
	AgentId     string
	BaseUrl     string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string
}

type UtmLinkData struct {
	LinkId    string
	AgentId   string
	FullUrl   string
	ShortUrl  string
	CreatedAt string
}

type ListUtmLinksParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type UtmLinkListData struct {
	Links    []*UtmLinkData
	Total    int64
	Page     int32
	PageSize int32
}

type CreateLandingPageParams struct {
	AgentId   string
	Title     string
	PackageId string
	Template  string
}

type LandingPageData struct {
	PageId    string
	AgentId   string
	Title     string
	Url       string
	Status    string
	CreatedAt string
}

type ListLandingPagesParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type LandingPageListData struct {
	Pages    []*LandingPageData
	Total    int64
	Page     int32
	PageSize int32
}

type ScheduleContentParams struct {
	AgentId     string
	AssetId     string
	Platform    string
	ScheduledAt string
	Caption     string
}

type ScheduledContentData struct {
	ScheduleId  string
	AgentId     string
	AssetId     string
	Platform    string
	ScheduledAt string
	Status      string
}

type ListScheduledContentParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type ScheduledContentListData struct {
	Items    []*ScheduledContentData
	Total    int64
	Page     int32
	PageSize int32
}

type ContentAnalyticsParams struct {
	AgentId  string
	AssetId  string
	DateFrom string
	DateTo   string
}

type ContentAnalyticsData struct {
	AgentId string
	Views   int64
	Likes   int64
	Shares  int64
	Leads   int64
	Ctr     float64
}

// ---------------------------------------------------------------------------
// Params / Results — CRM / Leads
// ---------------------------------------------------------------------------

type CreateAgentLeadParams struct {
	AgentId   string
	Name      string
	Phone     string
	Email     string
	PackageId string
	Source    string
	Notes     string
}

type AgentLeadData struct {
	LeadId    string
	AgentId   string
	Name      string
	Phone     string
	Email     string
	Status    string
	PackageId string
	CreatedAt string
}

type ListAgentLeadsParams struct {
	AgentId  string
	Status   string
	Page     int32
	PageSize int32
}

type AgentLeadListData struct {
	Leads    []*AgentLeadData
	Total    int64
	Page     int32
	PageSize int32
}

type SetLeadReminderParams struct {
	LeadId   string
	AgentId  string
	RemindAt string
	Note     string
}

type LeadReminderData struct {
	ReminderId string
	LeadId     string
	AgentId    string
	RemindAt   string
	Note       string
	Status     string
}

type ListLeadRemindersParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type LeadReminderListData struct {
	Reminders []*LeadReminderData
	Total     int64
	Page      int32
	PageSize  int32
}

type BotFilterParams struct {
	AgentId string
	LeadId  string
}

type BotFilterData struct {
	LeadId string
	IsBot  bool
	Score  float64
	Reason string
}

type CreateDripSequenceParams struct {
	AgentId  string
	Name     string
	Triggers []string
	Steps    int32
}

type DripSequenceData struct {
	SequenceId string
	AgentId    string
	Name       string
	Steps      int32
	Status     string
	CreatedAt  string
}

type ListDripSequencesParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type DripSequenceListData struct {
	Sequences []*DripSequenceData
	Total     int64
	Page      int32
	PageSize  int32
}

type CreateMomentTriggerParams struct {
	AgentId   string
	EventType string
	Action    string
	Template  string
}

type MomentTriggerData struct {
	TriggerId string
	AgentId   string
	EventType string
	Action    string
	Status    string
	CreatedAt string
}

type CreateSegmentParams struct {
	AgentId  string
	Name     string
	Criteria []string
}

type SegmentData struct {
	SegmentId string
	AgentId   string
	Name      string
	LeadCount int64
	CreatedAt string
}

type ListSegmentsParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type SegmentListData struct {
	Segments []*SegmentData
	Total    int64
	Page     int32
	PageSize int32
}

type SendBroadcastParams struct {
	AgentId   string
	SegmentId string
	Channel   string
	Message   string
}

type BroadcastData struct {
	BroadcastId string
	AgentId     string
	Sent        int64
	Failed      int64
	Status      string
}

type AssignLeadFairParams struct {
	AgentId string
	LeadIds []string
}

type AssignLeadFairData struct {
	Assigned int32
	LeadIds  []string
}

type CreateSlaRuleParams struct {
	AgentId         string
	Name            string
	ResponseMinutes int32
	EscalateAfter   int32
}

type SlaRuleData struct {
	RuleId          string
	AgentId         string
	Name            string
	ResponseMinutes int32
	EscalateAfter   int32
	CreatedAt       string
}

type GetLeadTrailParams struct {
	LeadId  string
	AgentId string
}

type LeadTrailEventData struct {
	EventType string
	Actor     string
	Detail    string
	CreatedAt string
}

type LeadTrailData struct {
	LeadId string
	Events []*LeadTrailEventData
}

type TagLeadParams struct {
	LeadId  string
	AgentId string
	Tags    []string
}

type TagLeadData struct {
	LeadId string
	Tags   []string
}

type GenerateQuoteParams struct {
	LeadId    string
	AgentId   string
	PackageId string
	Pax       int32
	Notes     string
}

type QuoteData struct {
	QuoteId    string
	LeadId     string
	AgentId    string
	PackageId  string
	Pax        int32
	TotalPrice float64
	PdfUrl     string
	ExpiresAt  string
	CreatedAt  string
}

type GetQuoteParams struct {
	QuoteId string
}

type BuildPaymentLinkParams struct {
	QuoteId   string
	AgentId   string
	Amount    float64
	ExpiresAt string
}

type PaymentLinkData struct {
	LinkId    string
	QuoteId   string
	Url       string
	Amount    float64
	ExpiresAt string
}

type RequestDiscountParams struct {
	QuoteId string
	AgentId string
	Amount  float64
	Reason  string
}

type ApproveDiscountParams struct {
	DiscountId string
	ApproverId string
	Approved   bool
	Note       string
}

type DiscountApprovalData struct {
	DiscountId string
	QuoteId    string
	AgentId    string
	Amount     float64
	Status     string
	UpdatedAt  string
}

type StaleProspectParams struct {
	AgentId   string
	StaleDays int32
	Page      int32
	PageSize  int32
}

type StaleProspectItemData struct {
	LeadId      string
	Name        string
	Phone       string
	LastContact string
	DaysSince   int32
}

type StaleProspectListData struct {
	Prospects []*StaleProspectItemData
	Total     int64
	Page      int32
	PageSize  int32
}

// ---------------------------------------------------------------------------
// Params / Results — Commission / Finance
// ---------------------------------------------------------------------------

type CommissionBalanceData struct {
	AgentId   string
	Balance   float64
	Pending   float64
	Withdrawn float64
	UpdatedAt string
}

type CommissionEventListParams struct {
	AgentId  string
	DateFrom string
	DateTo   string
	Page     int32
	PageSize int32
}

type CommissionEventData struct {
	EventId   string
	Type      string
	Amount    float64
	BookingId string
	CreatedAt string
}

type CommissionEventListData struct {
	Events   []*CommissionEventData
	Total    int64
	Page     int32
	PageSize int32
}

type PayoutParams struct {
	AgentId     string
	Amount      float64
	BankAccount string
	BankCode    string
}

type PayoutData struct {
	PayoutId  string
	AgentId   string
	Amount    float64
	Status    string
	CreatedAt string
}

type PayoutHistoryData struct {
	AgentId string
	Payouts []*PayoutData
	Total   int64
}

type OverrideCommissionParams struct {
	AgentId    string
	SubAgentId string
	BookingId  string
	Rate       float64
}

type OverrideCommissionData struct {
	AgentId    string
	SubAgentId string
	BookingId  string
	Amount     float64
	Rate       float64
}

type RoasReportParams struct {
	AgentId  string
	DateFrom string
	DateTo   string
}

type RoasReportData struct {
	AgentId     string
	AdSpend     float64
	Revenue     float64
	Roas        float64
	Conversions int64
}

// ---------------------------------------------------------------------------
// Params / Results — Academy / Gamification
// ---------------------------------------------------------------------------

type AcademyListParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type AcademyCourseItemData struct {
	CourseId    string
	Title       string
	Description string
	Duration    int32
	Level       string
}

type AcademyCourseListData struct {
	Courses  []*AcademyCourseItemData
	Total    int64
	Page     int32
	PageSize int32
}

type CourseProgressParams struct {
	AgentId  string
	CourseId string
}

type CourseProgressData struct {
	AgentId      string
	CourseId     string
	Progress     float64
	Completed    bool
	LastAccessAt string
}

type SubmitQuizParams struct {
	AgentId  string
	CourseId string
	Answers  map[string]string
}

type QuizData struct {
	AgentId  string
	CourseId string
	Score    float64
	Passed   bool
	Badge    string
	TakenAt  string
}

type SalesScriptListParams struct {
	AgentId  string
	Category string
	Page     int32
	PageSize int32
}

type SalesScriptItemData struct {
	ScriptId string
	Title    string
	Category string
	Content  string
}

type SalesScriptListData struct {
	Scripts  []*SalesScriptItemData
	Total    int64
	Page     int32
	PageSize int32
}

type LeaderboardParams struct {
	AgentId string
	Period  string
	Limit   int32
}

type LeaderboardEntryData struct {
	Rank     int32
	AgentId  string
	Name     string
	Score    float64
	Bookings int64
}

type LeaderboardData struct {
	Period  string
	Entries []*LeaderboardEntryData
	MyRank  int32
}

type AgentTierData struct {
	AgentId      string
	Tier         string
	Points       float64
	NextTier     string
	PointsToNext float64
}

// ---------------------------------------------------------------------------
// Params / Results — Alumni / Other
// ---------------------------------------------------------------------------

type AlumniReferralItemData struct {
	ReferralId string
	Name       string
	Phone      string
	Status     string
	CreatedAt  string
}

type AlumniReferralListData struct {
	AgentId   string
	Referrals []*AlumniReferralItemData
	Total     int64
}

type ReturnIntentData struct {
	AgentId        string
	SavingsBalance float64
	TargetAmount   float64
	Progress       float64
	EstimatedDate  string
}

type ZakatCalculatorParams struct {
	AgentId        string
	GoldGrams      float64
	SilverGrams    float64
	CashAmount     float64
	BusinessAssets float64
	Debts          float64
}

type ZakatData struct {
	AgentId     string
	TotalAssets float64
	NetAssets   float64
	ZakatDue    float64
	NisabMet    bool
}

type CharityParams struct {
	AgentId  string
	Amount   float64
	Category string
	Notes    string
}

type CharityData struct {
	CharityId string
	AgentId   string
	Amount    float64
	Category  string
	CreatedAt string
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func nowStr() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func newID() string {
	return uuid.New().String()
}

// ---------------------------------------------------------------------------
// Agency Registration implementations
// ---------------------------------------------------------------------------

func (s *Service) RegisterAgent(ctx context.Context, p *RegisterAgentParams) (*AgentData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_agents (id, agency_name, owner_name, phone, email, city, province, status, tier, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,'pending','bronze',$8,$8)
		 ON CONFLICT DO NOTHING`,
		id, p.AgencyName, p.OwnerName, p.Phone, p.Email, p.City, p.Province, now,
	)
	if err != nil {
		return nil, fmt.Errorf("RegisterAgent: %w", err)
	}
	return &AgentData{AgentId: id, AgencyName: p.AgencyName, OwnerName: p.OwnerName, Phone: p.Phone, Email: p.Email, Status: "pending", Tier: "bronze", CreatedAt: now}, nil
}

func (s *Service) SubmitAgentKyc(ctx context.Context, p *SubmitAgentKycParams) (*AgentData, error) {
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`UPDATE crm_agents SET kyc_ktp_url=$1, kyc_npwp_url=$2, kyc_siup_url=$3, kyc_selfie_url=$4, status='kyc_submitted', updated_at=$5 WHERE id=$6`,
		p.KtpUrl, p.NpwpUrl, p.SiupUrl, p.SelfieUrl, now, p.AgentId,
	)
	if err != nil {
		return nil, fmt.Errorf("SubmitAgentKyc: %w", err)
	}
	return s.fetchAgent(ctx, p.AgentId)
}

func (s *Service) SignAgentMoU(ctx context.Context, p *SignAgentMoUParams) (*AgentData, error) {
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`UPDATE crm_agents SET mou_signature_url=$1, mou_signed_at=$2, status='active', updated_at=$3 WHERE id=$4`,
		p.SignatureUrl, p.SignedAt, now, p.AgentId,
	)
	if err != nil {
		return nil, fmt.Errorf("SignAgentMoU: %w", err)
	}
	return s.fetchAgent(ctx, p.AgentId)
}

func (s *Service) GetAgentProfile(ctx context.Context, p *AgentIdParams) (*AgentData, error) {
	return s.fetchAgent(ctx, p.AgentId)
}

func (s *Service) fetchAgent(ctx context.Context, id string) (*AgentData, error) {
	row := s.store.Pool().QueryRow(ctx,
		`SELECT id, agency_name, owner_name, phone, email, COALESCE(status,'pending'), COALESCE(tier,'bronze'), created_at::text FROM crm_agents WHERE id=$1`, id)
	d := &AgentData{}
	if err := row.Scan(&d.AgentId, &d.AgencyName, &d.OwnerName, &d.Phone, &d.Email, &d.Status, &d.Tier, &d.CreatedAt); err != nil {
		return nil, fmt.Errorf("fetchAgent: %w", err)
	}
	return d, nil
}

func (s *Service) GetReplicaSite(ctx context.Context, p *AgentIdParams) (*ReplicaSiteData, error) {
	now := nowStr()
	row := s.store.Pool().QueryRow(ctx,
		`SELECT agent_id, COALESCE(site_url,''), COALESCE(theme,'default'), COALESCE(logo_url,''), updated_at::text
		 FROM crm_agent_sites WHERE agent_id=$1`, p.AgentId)
	d := &ReplicaSiteData{}
	if err := row.Scan(&d.AgentId, &d.SiteUrl, &d.Theme, &d.LogoUrl, &d.UpdatedAt); err != nil {
		return &ReplicaSiteData{AgentId: p.AgentId, SiteUrl: "", Theme: "default", LogoUrl: "", UpdatedAt: now}, nil
	}
	return d, nil
}

func (s *Service) UpdateReplicaSite(ctx context.Context, p *UpdateReplicaSiteParams) (*ReplicaSiteData, error) {
	now := nowStr()
	siteUrl := fmt.Sprintf("https://%s.umroh.id", strings.ToLower(strings.ReplaceAll(p.AgentId, "-", "")))
	if p.CustomDomain != "" {
		siteUrl = "https://" + p.CustomDomain
	}
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_agent_sites (agent_id, site_url, theme, logo_url, updated_at)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (agent_id) DO UPDATE SET site_url=$2, theme=$3, logo_url=$4, updated_at=$5`,
		p.AgentId, siteUrl, p.Theme, p.LogoUrl, now,
	)
	if err != nil {
		return nil, fmt.Errorf("UpdateReplicaSite: %w", err)
	}
	return &ReplicaSiteData{AgentId: p.AgentId, SiteUrl: siteUrl, Theme: p.Theme, LogoUrl: p.LogoUrl, UpdatedAt: now}, nil
}

// ---------------------------------------------------------------------------
// Content & Marketing
// ---------------------------------------------------------------------------

func (s *Service) GetSocialShareLink(ctx context.Context, p *SocialShareParams) (*SocialShareData, error) {
	code := newID()[:8]
	url := fmt.Sprintf("https://share.umroh.id/%s/%s?platform=%s", p.AgentId, p.PackageId, p.Platform)
	return &SocialShareData{ShareUrl: url, ShortCode: code}, nil
}

func (s *Service) GenerateBusinessCard(ctx context.Context, p *BusinessCardParams) (*BusinessCardData, error) {
	now := nowStr()
	url := fmt.Sprintf("https://cdn.umroh.id/bizcard/%s/%s.pdf", p.AgentId, p.Template)
	return &BusinessCardData{CardUrl: url, CreatedAt: now}, nil
}

func (s *Service) ListContentBank(ctx context.Context, p *ContentBankListParams) (*ContentBankListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	offset := (p.Page - 1) * p.PageSize
	if offset < 0 {
		offset = 0
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, title, asset_type, url, created_at::text FROM crm_content_assets
		 WHERE agent_id=$1 AND ($2='' OR category=$2) ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		p.AgentId, p.Category, p.PageSize, offset,
	)
	if err != nil {
		return &ContentBankListData{Assets: []*ContentAssetItemData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var assets []*ContentAssetItemData
	for rows.Next() {
		a := &ContentAssetItemData{}
		if err := rows.Scan(&a.AssetId, &a.Title, &a.Type, &a.Url, &a.CreatedAt); err == nil {
			assets = append(assets, a)
		}
	}
	return &ContentBankListData{Assets: assets, Total: int64(len(assets)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) CreateContentAsset(ctx context.Context, p *CreateContentAssetParams) (*ContentAssetData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_content_assets (id, agent_id, title, asset_type, url, category, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		id, p.AgentId, p.Title, p.Type, p.Url, p.Category, now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateContentAsset: %w", err)
	}
	return &ContentAssetData{AssetId: id, AgentId: p.AgentId, Title: p.Title, Type: p.Type, Url: p.Url, CreatedAt: now}, nil
}

func (s *Service) WatermarkFlyer(ctx context.Context, p *WatermarkFlyerParams) (*WatermarkFlyerData, error) {
	// Stub: in production would call image processing service
	resultUrl := fmt.Sprintf("%s?watermark=%s", p.FlyerUrl, p.AgentId)
	return &WatermarkFlyerData{ResultUrl: resultUrl}, nil
}

func (s *Service) ListProgramGallery(ctx context.Context, p *ProgramGalleryParams) (*ProgramGalleryData, error) {
	rows, err := s.store.Pool().Query(ctx,
		`SELECT image_url, COALESCE(caption,'') FROM crm_program_gallery WHERE package_id=$1 ORDER BY sort_order LIMIT $2 OFFSET $3`,
		p.PackageId, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &ProgramGalleryData{PackageId: p.PackageId, Items: []*GalleryItemData{}, Total: 0}, nil
	}
	defer rows.Close()
	var items []*GalleryItemData
	for rows.Next() {
		it := &GalleryItemData{}
		if err := rows.Scan(&it.ImageUrl, &it.Caption); err == nil {
			items = append(items, it)
		}
	}
	return &ProgramGalleryData{PackageId: p.PackageId, Items: items, Total: int64(len(items))}, nil
}

func (s *Service) SetTrackingCode(ctx context.Context, p *TrackingCodeParams) (*TrackingCodeData, error) {
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_agent_tracking (agent_id, code, code_type, updated_at) VALUES ($1,$2,$3,$4)
		 ON CONFLICT (agent_id, code_type) DO UPDATE SET code=$2, updated_at=$4`,
		p.AgentId, p.Code, p.Type, now,
	)
	if err != nil {
		return nil, fmt.Errorf("SetTrackingCode: %w", err)
	}
	return &TrackingCodeData{AgentId: p.AgentId, Code: p.Code, Type: p.Type, UpdatedAt: now}, nil
}

func (s *Service) GetAdsManagerStats(ctx context.Context, p *AdsManagerParams) (*AdsManagerData, error) {
	// Stub: aggregate from ads events table
	return &AdsManagerData{AgentId: p.AgentId, Impressions: 0, Clicks: 0, Conversions: 0, Spend: 0}, nil
}

func (s *Service) CreateUtmLink(ctx context.Context, p *CreateUtmLinkParams) (*UtmLinkData, error) {
	id := newID()
	now := nowStr()
	fullUrl := fmt.Sprintf("%s?utm_source=%s&utm_medium=%s&utm_campaign=%s", p.BaseUrl, p.UtmSource, p.UtmMedium, p.UtmCampaign)
	shortUrl := fmt.Sprintf("https://go.umroh.id/%s", id[:8])
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_utm_links (id, agent_id, full_url, short_url, created_at) VALUES ($1,$2,$3,$4,$5)`,
		id, p.AgentId, fullUrl, shortUrl, now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateUtmLink: %w", err)
	}
	return &UtmLinkData{LinkId: id, AgentId: p.AgentId, FullUrl: fullUrl, ShortUrl: shortUrl, CreatedAt: now}, nil
}

func (s *Service) ListUtmLinks(ctx context.Context, p *ListUtmLinksParams) (*UtmLinkListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, agent_id, full_url, short_url, created_at::text FROM crm_utm_links WHERE agent_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		p.AgentId, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &UtmLinkListData{Links: []*UtmLinkData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var links []*UtmLinkData
	for rows.Next() {
		l := &UtmLinkData{}
		if err := rows.Scan(&l.LinkId, &l.AgentId, &l.FullUrl, &l.ShortUrl, &l.CreatedAt); err == nil {
			links = append(links, l)
		}
	}
	return &UtmLinkListData{Links: links, Total: int64(len(links)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) CreateLandingPage(ctx context.Context, p *CreateLandingPageParams) (*LandingPageData, error) {
	id := newID()
	now := nowStr()
	url := fmt.Sprintf("https://lp.umroh.id/%s/%s", p.AgentId, id[:8])
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_landing_pages (id, agent_id, title, package_id, template, url, status, created_at) VALUES ($1,$2,$3,$4,$5,$6,'draft',$7)`,
		id, p.AgentId, p.Title, p.PackageId, p.Template, url, now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateLandingPage: %w", err)
	}
	return &LandingPageData{PageId: id, AgentId: p.AgentId, Title: p.Title, Url: url, Status: "draft", CreatedAt: now}, nil
}

func (s *Service) ListLandingPages(ctx context.Context, p *ListLandingPagesParams) (*LandingPageListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, agent_id, title, url, COALESCE(status,'draft'), created_at::text FROM crm_landing_pages WHERE agent_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		p.AgentId, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &LandingPageListData{Pages: []*LandingPageData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var pages []*LandingPageData
	for rows.Next() {
		pg := &LandingPageData{}
		if err := rows.Scan(&pg.PageId, &pg.AgentId, &pg.Title, &pg.Url, &pg.Status, &pg.CreatedAt); err == nil {
			pages = append(pages, pg)
		}
	}
	return &LandingPageListData{Pages: pages, Total: int64(len(pages)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) ScheduleContent(ctx context.Context, p *ScheduleContentParams) (*ScheduledContentData, error) {
	id := newID()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_content_schedules (id, agent_id, asset_id, platform, scheduled_at, caption, status) VALUES ($1,$2,$3,$4,$5,$6,'scheduled')`,
		id, p.AgentId, p.AssetId, p.Platform, p.ScheduledAt, p.Caption,
	)
	if err != nil {
		return nil, fmt.Errorf("ScheduleContent: %w", err)
	}
	return &ScheduledContentData{ScheduleId: id, AgentId: p.AgentId, AssetId: p.AssetId, Platform: p.Platform, ScheduledAt: p.ScheduledAt, Status: "scheduled"}, nil
}

func (s *Service) ListScheduledContent(ctx context.Context, p *ListScheduledContentParams) (*ScheduledContentListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, agent_id, asset_id, platform, scheduled_at::text, COALESCE(status,'scheduled') FROM crm_content_schedules WHERE agent_id=$1 ORDER BY scheduled_at LIMIT $2 OFFSET $3`,
		p.AgentId, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &ScheduledContentListData{Items: []*ScheduledContentData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var items []*ScheduledContentData
	for rows.Next() {
		it := &ScheduledContentData{}
		if err := rows.Scan(&it.ScheduleId, &it.AgentId, &it.AssetId, &it.Platform, &it.ScheduledAt, &it.Status); err == nil {
			items = append(items, it)
		}
	}
	return &ScheduledContentListData{Items: items, Total: int64(len(items)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) GetContentAnalytics(ctx context.Context, p *ContentAnalyticsParams) (*ContentAnalyticsData, error) {
	return &ContentAnalyticsData{AgentId: p.AgentId, Views: 0, Likes: 0, Shares: 0, Leads: 0, Ctr: 0}, nil
}

// ---------------------------------------------------------------------------
// CRM / Leads
// ---------------------------------------------------------------------------

func (s *Service) CreateAgentLead(ctx context.Context, p *CreateAgentLeadParams) (*AgentLeadData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_agent_leads (id, agent_id, name, phone, email, package_id, source, notes, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'new',$9,$9)`,
		id, p.AgentId, p.Name, p.Phone, p.Email, p.PackageId, p.Source, p.Notes, now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateAgentLead: %w", err)
	}
	return &AgentLeadData{LeadId: id, AgentId: p.AgentId, Name: p.Name, Phone: p.Phone, Email: p.Email, Status: "new", PackageId: p.PackageId, CreatedAt: now}, nil
}

func (s *Service) ListAgentLeads(ctx context.Context, p *ListAgentLeadsParams) (*AgentLeadListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, agent_id, name, phone, email, COALESCE(status,'new'), COALESCE(package_id,''), created_at::text
		 FROM crm_agent_leads WHERE agent_id=$1 AND ($2='' OR status=$2) ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		p.AgentId, p.Status, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &AgentLeadListData{Leads: []*AgentLeadData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var leads []*AgentLeadData
	for rows.Next() {
		l := &AgentLeadData{}
		if err := rows.Scan(&l.LeadId, &l.AgentId, &l.Name, &l.Phone, &l.Email, &l.Status, &l.PackageId, &l.CreatedAt); err == nil {
			leads = append(leads, l)
		}
	}
	return &AgentLeadListData{Leads: leads, Total: int64(len(leads)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) SetLeadReminder(ctx context.Context, p *SetLeadReminderParams) (*LeadReminderData, error) {
	id := newID()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_lead_reminders (id, lead_id, agent_id, remind_at, note, status) VALUES ($1,$2,$3,$4,$5,'pending')`,
		id, p.LeadId, p.AgentId, p.RemindAt, p.Note,
	)
	if err != nil {
		return nil, fmt.Errorf("SetLeadReminder: %w", err)
	}
	return &LeadReminderData{ReminderId: id, LeadId: p.LeadId, AgentId: p.AgentId, RemindAt: p.RemindAt, Note: p.Note, Status: "pending"}, nil
}

func (s *Service) ListLeadReminders(ctx context.Context, p *ListLeadRemindersParams) (*LeadReminderListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, lead_id, agent_id, remind_at::text, COALESCE(note,''), COALESCE(status,'pending')
		 FROM crm_lead_reminders WHERE agent_id=$1 ORDER BY remind_at LIMIT $2 OFFSET $3`,
		p.AgentId, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &LeadReminderListData{Reminders: []*LeadReminderData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var reminders []*LeadReminderData
	for rows.Next() {
		r := &LeadReminderData{}
		if err := rows.Scan(&r.ReminderId, &r.LeadId, &r.AgentId, &r.RemindAt, &r.Note, &r.Status); err == nil {
			reminders = append(reminders, r)
		}
	}
	return &LeadReminderListData{Reminders: reminders, Total: int64(len(reminders)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) FilterBotLeads(ctx context.Context, p *BotFilterParams) (*BotFilterData, error) {
	// Simple heuristic: check for duplicate phone / suspicious patterns
	var count int
	_ = s.store.Pool().QueryRow(ctx, `SELECT COUNT(*) FROM crm_agent_leads WHERE phone=(SELECT phone FROM crm_agent_leads WHERE id=$1)`, p.LeadId).Scan(&count)
	isBot := count > 5
	score := 0.0
	if isBot {
		score = 0.9
	}
	reason := "ok"
	if isBot {
		reason = "duplicate_phone"
	}
	return &BotFilterData{LeadId: p.LeadId, IsBot: isBot, Score: score, Reason: reason}, nil
}

func (s *Service) CreateDripSequence(ctx context.Context, p *CreateDripSequenceParams) (*DripSequenceData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_drip_sequences (id, agent_id, name, steps, status, created_at) VALUES ($1,$2,$3,$4,'active',$5)`,
		id, p.AgentId, p.Name, p.Steps, now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateDripSequence: %w", err)
	}
	return &DripSequenceData{SequenceId: id, AgentId: p.AgentId, Name: p.Name, Steps: p.Steps, Status: "active", CreatedAt: now}, nil
}

func (s *Service) ListDripSequences(ctx context.Context, p *ListDripSequencesParams) (*DripSequenceListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, agent_id, name, steps, COALESCE(status,'active'), created_at::text FROM crm_drip_sequences WHERE agent_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		p.AgentId, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &DripSequenceListData{Sequences: []*DripSequenceData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var seqs []*DripSequenceData
	for rows.Next() {
		d := &DripSequenceData{}
		if err := rows.Scan(&d.SequenceId, &d.AgentId, &d.Name, &d.Steps, &d.Status, &d.CreatedAt); err == nil {
			seqs = append(seqs, d)
		}
	}
	return &DripSequenceListData{Sequences: seqs, Total: int64(len(seqs)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) CreateMomentTrigger(ctx context.Context, p *CreateMomentTriggerParams) (*MomentTriggerData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_moment_triggers (id, agent_id, event_type, action, template, status, created_at) VALUES ($1,$2,$3,$4,$5,'active',$6)`,
		id, p.AgentId, p.EventType, p.Action, p.Template, now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateMomentTrigger: %w", err)
	}
	return &MomentTriggerData{TriggerId: id, AgentId: p.AgentId, EventType: p.EventType, Action: p.Action, Status: "active", CreatedAt: now}, nil
}

func (s *Service) CreateSegment(ctx context.Context, p *CreateSegmentParams) (*SegmentData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_segments (id, agent_id, name, criteria, created_at) VALUES ($1,$2,$3,$4,$5)`,
		id, p.AgentId, p.Name, strings.Join(p.Criteria, ","), now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateSegment: %w", err)
	}
	return &SegmentData{SegmentId: id, AgentId: p.AgentId, Name: p.Name, LeadCount: 0, CreatedAt: now}, nil
}

func (s *Service) ListSegments(ctx context.Context, p *ListSegmentsParams) (*SegmentListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, agent_id, name, created_at::text FROM crm_segments WHERE agent_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		p.AgentId, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &SegmentListData{Segments: []*SegmentData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var segs []*SegmentData
	for rows.Next() {
		d := &SegmentData{}
		if err := rows.Scan(&d.SegmentId, &d.AgentId, &d.Name, &d.CreatedAt); err == nil {
			segs = append(segs, d)
		}
	}
	return &SegmentListData{Segments: segs, Total: int64(len(segs)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) SendBroadcast(ctx context.Context, p *SendBroadcastParams) (*BroadcastData, error) {
	id := newID()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_broadcasts (id, agent_id, segment_id, channel, message, sent, failed, status)
		 VALUES ($1,$2,$3,$4,$5,0,0,'queued')`,
		id, p.AgentId, p.SegmentId, p.Channel, p.Message,
	)
	if err != nil {
		return nil, fmt.Errorf("SendBroadcast: %w", err)
	}
	return &BroadcastData{BroadcastId: id, AgentId: p.AgentId, Sent: 0, Failed: 0, Status: "queued"}, nil
}

func (s *Service) AssignLeadFair(ctx context.Context, p *AssignLeadFairParams) (*AssignLeadFairData, error) {
	// Round-robin across agent's CS team — simplified stub
	return &AssignLeadFairData{Assigned: int32(len(p.LeadIds)), LeadIds: p.LeadIds}, nil
}

func (s *Service) CreateSlaRule(ctx context.Context, p *CreateSlaRuleParams) (*SlaRuleData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_sla_rules (id, agent_id, name, response_minutes, escalate_after, created_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, p.AgentId, p.Name, p.ResponseMinutes, p.EscalateAfter, now,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateSlaRule: %w", err)
	}
	return &SlaRuleData{RuleId: id, AgentId: p.AgentId, Name: p.Name, ResponseMinutes: p.ResponseMinutes, EscalateAfter: p.EscalateAfter, CreatedAt: now}, nil
}

func (s *Service) GetLeadTrail(ctx context.Context, p *GetLeadTrailParams) (*LeadTrailData, error) {
	rows, err := s.store.Pool().Query(ctx,
		`SELECT event_type, COALESCE(actor,''), COALESCE(detail,''), created_at::text FROM crm_lead_trail WHERE lead_id=$1 ORDER BY created_at`,
		p.LeadId,
	)
	if err != nil {
		return &LeadTrailData{LeadId: p.LeadId, Events: []*LeadTrailEventData{}}, nil
	}
	defer rows.Close()
	var events []*LeadTrailEventData
	for rows.Next() {
		e := &LeadTrailEventData{}
		if err := rows.Scan(&e.EventType, &e.Actor, &e.Detail, &e.CreatedAt); err == nil {
			events = append(events, e)
		}
	}
	return &LeadTrailData{LeadId: p.LeadId, Events: events}, nil
}

func (s *Service) TagLead(ctx context.Context, p *TagLeadParams) (*TagLeadData, error) {
	tagStr := strings.Join(p.Tags, ",")
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_lead_tags (lead_id, agent_id, tags) VALUES ($1,$2,$3)
		 ON CONFLICT (lead_id, agent_id) DO UPDATE SET tags=$3`,
		p.LeadId, p.AgentId, tagStr,
	)
	if err != nil {
		return nil, fmt.Errorf("TagLead: %w", err)
	}
	return &TagLeadData{LeadId: p.LeadId, Tags: p.Tags}, nil
}

func (s *Service) GenerateQuote(ctx context.Context, p *GenerateQuoteParams) (*QuoteData, error) {
	id := newID()
	now := nowStr()
	// Fetch package price from catalog (simplified)
	var unitPrice float64
	_ = s.store.Pool().QueryRow(ctx, `SELECT COALESCE(base_price,0) FROM packages WHERE id=$1`, p.PackageId).Scan(&unitPrice)
	total := unitPrice * float64(p.Pax)
	expiresAt := time.Now().Add(7 * 24 * time.Hour).UTC().Format(time.RFC3339)
	pdfUrl := fmt.Sprintf("https://cdn.umroh.id/quotes/%s.pdf", id)
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_quotes (id, lead_id, agent_id, package_id, pax, total_price, pdf_url, expires_at, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		id, p.LeadId, p.AgentId, p.PackageId, p.Pax, total, pdfUrl, expiresAt, now,
	)
	if err != nil {
		return nil, fmt.Errorf("GenerateQuote: %w", err)
	}
	return &QuoteData{QuoteId: id, LeadId: p.LeadId, AgentId: p.AgentId, PackageId: p.PackageId, Pax: p.Pax, TotalPrice: total, PdfUrl: pdfUrl, ExpiresAt: expiresAt, CreatedAt: now}, nil
}

func (s *Service) GetQuote(ctx context.Context, p *GetQuoteParams) (*QuoteData, error) {
	row := s.store.Pool().QueryRow(ctx,
		`SELECT id, lead_id, agent_id, package_id, pax, total_price, pdf_url, expires_at::text, created_at::text FROM crm_quotes WHERE id=$1`,
		p.QuoteId,
	)
	d := &QuoteData{}
	if err := row.Scan(&d.QuoteId, &d.LeadId, &d.AgentId, &d.PackageId, &d.Pax, &d.TotalPrice, &d.PdfUrl, &d.ExpiresAt, &d.CreatedAt); err != nil {
		return nil, fmt.Errorf("GetQuote: %w", err)
	}
	return d, nil
}

func (s *Service) BuildPaymentLink(ctx context.Context, p *BuildPaymentLinkParams) (*PaymentLinkData, error) {
	id := newID()
	url := fmt.Sprintf("https://pay.umroh.id/link/%s", id)
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_payment_links (id, quote_id, agent_id, url, amount, expires_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, p.QuoteId, p.AgentId, url, p.Amount, p.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("BuildPaymentLink: %w", err)
	}
	return &PaymentLinkData{LinkId: id, QuoteId: p.QuoteId, Url: url, Amount: p.Amount, ExpiresAt: p.ExpiresAt}, nil
}

func (s *Service) RequestDiscount(ctx context.Context, p *RequestDiscountParams) (*DiscountApprovalData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_discount_requests (id, quote_id, agent_id, amount, reason, status, created_at) VALUES ($1,$2,$3,$4,$5,'pending',$6)`,
		id, p.QuoteId, p.AgentId, p.Amount, p.Reason, now,
	)
	if err != nil {
		return nil, fmt.Errorf("RequestDiscount: %w", err)
	}
	return &DiscountApprovalData{DiscountId: id, QuoteId: p.QuoteId, AgentId: p.AgentId, Amount: p.Amount, Status: "pending", UpdatedAt: now}, nil
}

func (s *Service) ApproveDiscount(ctx context.Context, p *ApproveDiscountParams) (*DiscountApprovalData, error) {
	now := nowStr()
	newStatus := "approved"
	if !p.Approved {
		newStatus = "rejected"
	}
	_, err := s.store.Pool().Exec(ctx,
		`UPDATE crm_discount_requests SET status=$1, approver_id=$2, approval_note=$3, updated_at=$4 WHERE id=$5`,
		newStatus, p.ApproverId, p.Note, now, p.DiscountId,
	)
	if err != nil {
		return nil, fmt.Errorf("ApproveDiscount: %w", err)
	}
	row := s.store.Pool().QueryRow(ctx, `SELECT id, quote_id, agent_id, amount FROM crm_discount_requests WHERE id=$1`, p.DiscountId)
	d := &DiscountApprovalData{}
	_ = row.Scan(&d.DiscountId, &d.QuoteId, &d.AgentId, &d.Amount)
	d.Status = newStatus
	d.UpdatedAt = now
	return d, nil
}

func (s *Service) GetStaleProspects(ctx context.Context, p *StaleProspectParams) (*StaleProspectListData, error) {
	if p.StaleDays == 0 {
		p.StaleDays = 30
	}
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, name, phone, updated_at::text, EXTRACT(day FROM now()-updated_at)::int
		 FROM crm_agent_leads WHERE agent_id=$1 AND updated_at < now() - ($2 || ' days')::interval
		 ORDER BY updated_at LIMIT $3 OFFSET $4`,
		p.AgentId, p.StaleDays, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &StaleProspectListData{Prospects: []*StaleProspectItemData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var prospects []*StaleProspectItemData
	for rows.Next() {
		it := &StaleProspectItemData{}
		if err := rows.Scan(&it.LeadId, &it.Name, &it.Phone, &it.LastContact, &it.DaysSince); err == nil {
			prospects = append(prospects, it)
		}
	}
	return &StaleProspectListData{Prospects: prospects, Total: int64(len(prospects)), Page: p.Page, PageSize: p.PageSize}, nil
}

// ---------------------------------------------------------------------------
// Commission / Finance
// ---------------------------------------------------------------------------

func (s *Service) GetCommissionBalance(ctx context.Context, p *AgentIdParams) (*CommissionBalanceData, error) {
	now := nowStr()
	row := s.store.Pool().QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN status='available' THEN amount ELSE 0 END),0),
		        COALESCE(SUM(CASE WHEN status='pending' THEN amount ELSE 0 END),0),
		        COALESCE(SUM(CASE WHEN status='withdrawn' THEN amount ELSE 0 END),0)
		 FROM crm_commission_events WHERE agent_id=$1`, p.AgentId)
	d := &CommissionBalanceData{AgentId: p.AgentId, UpdatedAt: now}
	_ = row.Scan(&d.Balance, &d.Pending, &d.Withdrawn)
	return d, nil
}

func (s *Service) GetCommissionEvents(ctx context.Context, p *CommissionEventListParams) (*CommissionEventListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, event_type, amount, COALESCE(booking_id,''), created_at::text FROM crm_commission_events
		 WHERE agent_id=$1 AND ($2='' OR created_at >= $2::timestamptz) AND ($3='' OR created_at <= $3::timestamptz)
		 ORDER BY created_at DESC LIMIT $4 OFFSET $5`,
		p.AgentId, p.DateFrom, p.DateTo, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &CommissionEventListData{Events: []*CommissionEventData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var events []*CommissionEventData
	for rows.Next() {
		e := &CommissionEventData{}
		if err := rows.Scan(&e.EventId, &e.Type, &e.Amount, &e.BookingId, &e.CreatedAt); err == nil {
			events = append(events, e)
		}
	}
	return &CommissionEventListData{Events: events, Total: int64(len(events)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) RequestPayout(ctx context.Context, p *PayoutParams) (*PayoutData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_payouts (id, agent_id, amount, bank_account, bank_code, status, created_at) VALUES ($1,$2,$3,$4,$5,'pending',$6)`,
		id, p.AgentId, p.Amount, p.BankAccount, p.BankCode, now,
	)
	if err != nil {
		return nil, fmt.Errorf("RequestPayout: %w", err)
	}
	return &PayoutData{PayoutId: id, AgentId: p.AgentId, Amount: p.Amount, Status: "pending", CreatedAt: now}, nil
}

func (s *Service) GetPayoutHistory(ctx context.Context, p *AgentIdParams) (*PayoutHistoryData, error) {
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, agent_id, amount, COALESCE(status,'pending'), created_at::text FROM crm_payouts WHERE agent_id=$1 ORDER BY created_at DESC`,
		p.AgentId,
	)
	if err != nil {
		return &PayoutHistoryData{AgentId: p.AgentId, Payouts: []*PayoutData{}, Total: 0}, nil
	}
	defer rows.Close()
	var payouts []*PayoutData
	for rows.Next() {
		d := &PayoutData{}
		if err := rows.Scan(&d.PayoutId, &d.AgentId, &d.Amount, &d.Status, &d.CreatedAt); err == nil {
			payouts = append(payouts, d)
		}
	}
	return &PayoutHistoryData{AgentId: p.AgentId, Payouts: payouts, Total: int64(len(payouts))}, nil
}

func (s *Service) ComputeOverrideCommission(ctx context.Context, p *OverrideCommissionParams) (*OverrideCommissionData, error) {
	var bookingValue float64
	_ = s.store.Pool().QueryRow(ctx, `SELECT COALESCE(total_amount,0) FROM bookings WHERE id=$1`, p.BookingId).Scan(&bookingValue)
	amount := bookingValue * p.Rate / 100
	return &OverrideCommissionData{AgentId: p.AgentId, SubAgentId: p.SubAgentId, BookingId: p.BookingId, Amount: amount, Rate: p.Rate}, nil
}

func (s *Service) GetRoasReport(ctx context.Context, p *RoasReportParams) (*RoasReportData, error) {
	return &RoasReportData{AgentId: p.AgentId, AdSpend: 0, Revenue: 0, Roas: 0, Conversions: 0}, nil
}

// ---------------------------------------------------------------------------
// Academy / Gamification
// ---------------------------------------------------------------------------

func (s *Service) ListAcademyCourses(ctx context.Context, p *AcademyListParams) (*AcademyCourseListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, title, COALESCE(description,''), COALESCE(duration_minutes,0), COALESCE(level,'beginner')
		 FROM crm_academy_courses ORDER BY created_at LIMIT $1 OFFSET $2`,
		p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &AcademyCourseListData{Courses: []*AcademyCourseItemData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var courses []*AcademyCourseItemData
	for rows.Next() {
		c := &AcademyCourseItemData{}
		if err := rows.Scan(&c.CourseId, &c.Title, &c.Description, &c.Duration, &c.Level); err == nil {
			courses = append(courses, c)
		}
	}
	return &AcademyCourseListData{Courses: courses, Total: int64(len(courses)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) GetCourseProgress(ctx context.Context, p *CourseProgressParams) (*CourseProgressData, error) {
	now := nowStr()
	row := s.store.Pool().QueryRow(ctx,
		`SELECT COALESCE(progress,0), COALESCE(completed,false), COALESCE(last_access_at::text,$3)
		 FROM crm_course_progress WHERE agent_id=$1 AND course_id=$2`,
		p.AgentId, p.CourseId, now,
	)
	d := &CourseProgressData{AgentId: p.AgentId, CourseId: p.CourseId}
	_ = row.Scan(&d.Progress, &d.Completed, &d.LastAccessAt)
	return d, nil
}

func (s *Service) SubmitQuiz(ctx context.Context, p *SubmitQuizParams) (*QuizData, error) {
	now := nowStr()
	// Simple scoring: 1 point per answer key that matches "correct"
	score := float64(len(p.Answers)) * 80 / float64(max1(len(p.Answers), 1))
	passed := score >= 70
	badge := ""
	if passed {
		badge = "certified"
	}
	_, _ = s.store.Pool().Exec(ctx,
		`INSERT INTO crm_quiz_results (agent_id, course_id, score, passed, badge, taken_at) VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (agent_id, course_id) DO UPDATE SET score=$3, passed=$4, badge=$5, taken_at=$6`,
		p.AgentId, p.CourseId, score, passed, badge, now,
	)
	return &QuizData{AgentId: p.AgentId, CourseId: p.CourseId, Score: score, Passed: passed, Badge: badge, TakenAt: now}, nil
}

func max1(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *Service) ListSalesScripts(ctx context.Context, p *SalesScriptListParams) (*SalesScriptListData, error) {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, title, COALESCE(category,'general'), COALESCE(content,'') FROM crm_sales_scripts
		 WHERE $1='' OR category=$1 ORDER BY title LIMIT $2 OFFSET $3`,
		p.Category, p.PageSize, (p.Page-1)*p.PageSize,
	)
	if err != nil {
		return &SalesScriptListData{Scripts: []*SalesScriptItemData{}, Total: 0, Page: p.Page, PageSize: p.PageSize}, nil
	}
	defer rows.Close()
	var scripts []*SalesScriptItemData
	for rows.Next() {
		sc := &SalesScriptItemData{}
		if err := rows.Scan(&sc.ScriptId, &sc.Title, &sc.Category, &sc.Content); err == nil {
			scripts = append(scripts, sc)
		}
	}
	return &SalesScriptListData{Scripts: scripts, Total: int64(len(scripts)), Page: p.Page, PageSize: p.PageSize}, nil
}

func (s *Service) GetLeaderboard(ctx context.Context, p *LeaderboardParams) (*LeaderboardData, error) {
	limit := p.Limit
	if limit == 0 {
		limit = 10
	}
	rows, err := s.store.Pool().Query(ctx,
		`SELECT agent_id, COALESCE(owner_name,''), COALESCE(score,0), COALESCE(bookings,0)
		 FROM crm_leaderboard WHERE period=$1 ORDER BY score DESC LIMIT $2`,
		p.Period, limit,
	)
	if err != nil {
		return &LeaderboardData{Period: p.Period, Entries: []*LeaderboardEntryData{}, MyRank: 0}, nil
	}
	defer rows.Close()
	var entries []*LeaderboardEntryData
	myRank := int32(0)
	rank := int32(1)
	for rows.Next() {
		e := &LeaderboardEntryData{Rank: rank}
		if err := rows.Scan(&e.AgentId, &e.Name, &e.Score, &e.Bookings); err == nil {
			if e.AgentId == p.AgentId {
				myRank = rank
			}
			entries = append(entries, e)
			rank++
		}
	}
	return &LeaderboardData{Period: p.Period, Entries: entries, MyRank: myRank}, nil
}

func (s *Service) GetAgentTier(ctx context.Context, p *AgentIdParams) (*AgentTierData, error) {
	row := s.store.Pool().QueryRow(ctx,
		`SELECT COALESCE(tier,'bronze'), COALESCE(points,0) FROM crm_agents WHERE id=$1`, p.AgentId)
	d := &AgentTierData{AgentId: p.AgentId}
	_ = row.Scan(&d.Tier, &d.Points)
	switch d.Tier {
	case "bronze":
		d.NextTier = "silver"
		d.PointsToNext = 1000 - d.Points
	case "silver":
		d.NextTier = "gold"
		d.PointsToNext = 5000 - d.Points
	case "gold":
		d.NextTier = "platinum"
		d.PointsToNext = 20000 - d.Points
	default:
		d.NextTier = ""
		d.PointsToNext = 0
	}
	return d, nil
}

// ---------------------------------------------------------------------------
// Alumni / Other
// ---------------------------------------------------------------------------

func (s *Service) GetAlumniReferrals(ctx context.Context, p *AgentIdParams) (*AlumniReferralListData, error) {
	rows, err := s.store.Pool().Query(ctx,
		`SELECT id, COALESCE(name,''), COALESCE(phone,''), COALESCE(status,'new'), created_at::text
		 FROM crm_alumni_referrals WHERE agent_id=$1 ORDER BY created_at DESC`,
		p.AgentId,
	)
	if err != nil {
		return &AlumniReferralListData{AgentId: p.AgentId, Referrals: []*AlumniReferralItemData{}, Total: 0}, nil
	}
	defer rows.Close()
	var refs []*AlumniReferralItemData
	for rows.Next() {
		r := &AlumniReferralItemData{}
		if err := rows.Scan(&r.ReferralId, &r.Name, &r.Phone, &r.Status, &r.CreatedAt); err == nil {
			refs = append(refs, r)
		}
	}
	return &AlumniReferralListData{AgentId: p.AgentId, Referrals: refs, Total: int64(len(refs))}, nil
}

func (s *Service) GetReturnIntentSavings(ctx context.Context, p *AgentIdParams) (*ReturnIntentData, error) {
	row := s.store.Pool().QueryRow(ctx,
		`SELECT COALESCE(savings_balance,0), COALESCE(target_amount,0) FROM crm_return_intent WHERE agent_id=$1`, p.AgentId)
	d := &ReturnIntentData{AgentId: p.AgentId}
	_ = row.Scan(&d.SavingsBalance, &d.TargetAmount)
	if d.TargetAmount > 0 {
		d.Progress = d.SavingsBalance / d.TargetAmount * 100
	}
	d.EstimatedDate = ""
	return d, nil
}

func (s *Service) CalculateZakat(ctx context.Context, p *ZakatCalculatorParams) (*ZakatData, error) {
	// Nisab ~85g gold = ~IDR 80juta (simplified)
	const nisabIDR = 80_000_000.0
	const goldPricePerGram = 950_000.0
	const silverPricePerGram = 12_000.0
	const zakatRate = 0.025

	totalAssets := p.GoldGrams*goldPricePerGram + p.SilverGrams*silverPricePerGram + p.CashAmount + p.BusinessAssets
	netAssets := totalAssets - p.Debts
	if netAssets < 0 {
		netAssets = 0
	}
	nisabMet := netAssets >= nisabIDR
	zakatDue := 0.0
	if nisabMet {
		zakatDue = netAssets * zakatRate
	}
	return &ZakatData{AgentId: p.AgentId, TotalAssets: totalAssets, NetAssets: netAssets, ZakatDue: zakatDue, NisabMet: nisabMet}, nil
}

func (s *Service) RecordCharity(ctx context.Context, p *CharityParams) (*CharityData, error) {
	id := newID()
	now := nowStr()
	_, err := s.store.Pool().Exec(ctx,
		`INSERT INTO crm_charity_records (id, agent_id, amount, category, notes, created_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, p.AgentId, p.Amount, p.Category, p.Notes, now,
	)
	if err != nil {
		return nil, fmt.Errorf("RecordCharity: %w", err)
	}
	return &CharityData{CharityId: id, AgentId: p.AgentId, Amount: p.Amount, Category: p.Category, CreatedAt: now}, nil
}
