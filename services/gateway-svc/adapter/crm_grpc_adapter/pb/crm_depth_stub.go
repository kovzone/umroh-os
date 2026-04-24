// crm_depth_stub.go — gateway-side gRPC client stubs for crm-svc Wave 7 depth RPCs.
// BL-CRM-010..012, 017..066: agency registration, content, leads/CRM, commission,
// academy, alumni/other.
//
// Mirrors services/crm-svc/api/grpc_api/pb/crm_depth_messages.go.
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Method name constants
// ---------------------------------------------------------------------------

const (
	CrmService_RegisterAgent_FullMethodName             = "/crm.CrmService/RegisterAgent"
	CrmService_SubmitAgentKyc_FullMethodName            = "/crm.CrmService/SubmitAgentKyc"
	CrmService_SignAgentMoU_FullMethodName              = "/crm.CrmService/SignAgentMoU"
	CrmService_GetAgentProfile_FullMethodName           = "/crm.CrmService/GetAgentProfile"
	CrmService_GetReplicaSite_FullMethodName            = "/crm.CrmService/GetReplicaSite"
	CrmService_UpdateReplicaSite_FullMethodName         = "/crm.CrmService/UpdateReplicaSite"
	CrmService_GetSocialShareLink_FullMethodName        = "/crm.CrmService/GetSocialShareLink"
	CrmService_GenerateBusinessCard_FullMethodName      = "/crm.CrmService/GenerateBusinessCard"
	CrmService_ListContentBank_FullMethodName           = "/crm.CrmService/ListContentBank"
	CrmService_CreateContentAsset_FullMethodName        = "/crm.CrmService/CreateContentAsset"
	CrmService_WatermarkFlyer_FullMethodName            = "/crm.CrmService/WatermarkFlyer"
	CrmService_ListProgramGallery_FullMethodName        = "/crm.CrmService/ListProgramGallery"
	CrmService_SetTrackingCode_FullMethodName           = "/crm.CrmService/SetTrackingCode"
	CrmService_GetAdsManagerStats_FullMethodName        = "/crm.CrmService/GetAdsManagerStats"
	CrmService_CreateUtmLink_FullMethodName             = "/crm.CrmService/CreateUtmLink"
	CrmService_ListUtmLinks_FullMethodName              = "/crm.CrmService/ListUtmLinks"
	CrmService_CreateLandingPage_FullMethodName         = "/crm.CrmService/CreateLandingPage"
	CrmService_ListLandingPages_FullMethodName          = "/crm.CrmService/ListLandingPages"
	CrmService_ScheduleContent_FullMethodName           = "/crm.CrmService/ScheduleContent"
	CrmService_ListScheduledContent_FullMethodName      = "/crm.CrmService/ListScheduledContent"
	CrmService_GetContentAnalytics_FullMethodName       = "/crm.CrmService/GetContentAnalytics"
	CrmService_CreateAgentLead_FullMethodName           = "/crm.CrmService/CreateAgentLead"
	CrmService_ListAgentLeads_FullMethodName            = "/crm.CrmService/ListAgentLeads"
	CrmService_SetLeadReminder_FullMethodName           = "/crm.CrmService/SetLeadReminder"
	CrmService_ListLeadReminders_FullMethodName         = "/crm.CrmService/ListLeadReminders"
	CrmService_FilterBotLeads_FullMethodName            = "/crm.CrmService/FilterBotLeads"
	CrmService_CreateDripSequence_FullMethodName        = "/crm.CrmService/CreateDripSequence"
	CrmService_ListDripSequences_FullMethodName         = "/crm.CrmService/ListDripSequences"
	CrmService_CreateMomentTrigger_FullMethodName       = "/crm.CrmService/CreateMomentTrigger"
	CrmService_CreateSegment_FullMethodName             = "/crm.CrmService/CreateSegment"
	CrmService_ListSegments_FullMethodName              = "/crm.CrmService/ListSegments"
	CrmService_SendBroadcast_FullMethodName             = "/crm.CrmService/SendBroadcast"
	CrmService_AssignLeadFair_FullMethodName            = "/crm.CrmService/AssignLeadFair"
	CrmService_CreateSlaRule_FullMethodName             = "/crm.CrmService/CreateSlaRule"
	CrmService_GetLeadTrail_FullMethodName              = "/crm.CrmService/GetLeadTrail"
	CrmService_TagLead_FullMethodName                   = "/crm.CrmService/TagLead"
	CrmService_GenerateQuote_FullMethodName             = "/crm.CrmService/GenerateQuote"
	CrmService_GetQuote_FullMethodName                  = "/crm.CrmService/GetQuote"
	CrmService_BuildPaymentLink_FullMethodName          = "/crm.CrmService/BuildPaymentLink"
	CrmService_RequestDiscount_FullMethodName           = "/crm.CrmService/RequestDiscount"
	CrmService_ApproveDiscount_FullMethodName           = "/crm.CrmService/ApproveDiscount"
	CrmService_GetStaleProspects_FullMethodName         = "/crm.CrmService/GetStaleProspects"
	CrmService_GetCommissionBalance_FullMethodName      = "/crm.CrmService/GetCommissionBalance"
	CrmService_GetCommissionEvents_FullMethodName       = "/crm.CrmService/GetCommissionEvents"
	CrmService_RequestPayout_FullMethodName             = "/crm.CrmService/RequestPayout"
	CrmService_GetPayoutHistory_FullMethodName          = "/crm.CrmService/GetPayoutHistory"
	CrmService_ComputeOverrideCommission_FullMethodName = "/crm.CrmService/ComputeOverrideCommission"
	CrmService_GetRoasReport_FullMethodName             = "/crm.CrmService/GetRoasReport"
	CrmService_ListAcademyCourses_FullMethodName        = "/crm.CrmService/ListAcademyCourses"
	CrmService_GetCourseProgress_FullMethodName         = "/crm.CrmService/GetCourseProgress"
	CrmService_SubmitQuiz_FullMethodName                = "/crm.CrmService/SubmitQuiz"
	CrmService_ListSalesScripts_FullMethodName          = "/crm.CrmService/ListSalesScripts"
	CrmService_GetLeaderboard_FullMethodName            = "/crm.CrmService/GetLeaderboard"
	CrmService_GetAgentTier_FullMethodName              = "/crm.CrmService/GetAgentTier"
	CrmService_GetAlumniReferrals_FullMethodName        = "/crm.CrmService/GetAlumniReferrals"
	CrmService_GetReturnIntentSavings_FullMethodName    = "/crm.CrmService/GetReturnIntentSavings"
	CrmService_CalculateZakat_FullMethodName            = "/crm.CrmService/CalculateZakat"
	CrmService_RecordCharity_FullMethodName             = "/crm.CrmService/RecordCharity"
)

// ---------------------------------------------------------------------------
// Agency Registration messages
// ---------------------------------------------------------------------------

type CrmDepthAgentRegisterRequest struct {
	AgencyName string `json:"agency_name"`
	OwnerName  string `json:"owner_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type CrmDepthAgentResult struct {
	Id         string `json:"id"`
	AgencyName string `json:"agency_name"`
	Status     string `json:"status"`
	TierLevel  string `json:"tier_level"`
	CreatedAt  string `json:"created_at"`
}

func (x *CrmDepthAgentResult) GetId() string        { if x == nil { return "" }; return x.Id }
func (x *CrmDepthAgentResult) GetAgencyName() string { if x == nil { return "" }; return x.AgencyName }
func (x *CrmDepthAgentResult) GetStatus() string    { if x == nil { return "" }; return x.Status }
func (x *CrmDepthAgentResult) GetTierLevel() string { if x == nil { return "" }; return x.TierLevel }
func (x *CrmDepthAgentResult) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthSubmitKycRequest struct {
	AgentId  string `json:"agent_id"`
	KycType  string `json:"kyc_type"`
	DocUrl   string `json:"doc_url"`
}

type CrmDepthKycResult struct {
	AgentId   string `json:"agent_id"`
	KycStatus string `json:"kyc_status"`
	UpdatedAt string `json:"updated_at"`
}

func (x *CrmDepthKycResult) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthKycResult) GetKycStatus() string { if x == nil { return "" }; return x.KycStatus }
func (x *CrmDepthKycResult) GetUpdatedAt() string { if x == nil { return "" }; return x.UpdatedAt }

type CrmDepthSignMoURequest struct {
	AgentId    string `json:"agent_id"`
	SignedDocUrl string `json:"signed_doc_url"`
}

type CrmDepthMoUResult struct {
	AgentId   string `json:"agent_id"`
	MouStatus string `json:"mou_status"`
	SignedAt  string `json:"signed_at"`
}

func (x *CrmDepthMoUResult) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthMoUResult) GetMouStatus() string { if x == nil { return "" }; return x.MouStatus }
func (x *CrmDepthMoUResult) GetSignedAt() string  { if x == nil { return "" }; return x.SignedAt }

type CrmDepthGetAgentRequest struct {
	AgentId string `json:"agent_id"`
}

type CrmDepthReplicaSiteResult struct {
	AgentId   string `json:"agent_id"`
	SiteUrl   string `json:"site_url"`
	UpdatedAt string `json:"updated_at"`
}

func (x *CrmDepthReplicaSiteResult) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthReplicaSiteResult) GetSiteUrl() string   { if x == nil { return "" }; return x.SiteUrl }
func (x *CrmDepthReplicaSiteResult) GetUpdatedAt() string { if x == nil { return "" }; return x.UpdatedAt }

type CrmDepthUpdateReplicaSiteRequest struct {
	AgentId    string `json:"agent_id"`
	CustomSlug string `json:"custom_slug"`
	ThemeColor string `json:"theme_color"`
}

// ---------------------------------------------------------------------------
// Content messages
// ---------------------------------------------------------------------------

type CrmDepthShareLinkRequest struct {
	AgentId  string `json:"agent_id"`
	Platform string `json:"platform"`
	ContentId string `json:"content_id"`
}

type CrmDepthShareLinkResult struct {
	ShareUrl  string `json:"share_url"`
	ExpiresAt string `json:"expires_at"`
}

func (x *CrmDepthShareLinkResult) GetShareUrl() string  { if x == nil { return "" }; return x.ShareUrl }
func (x *CrmDepthShareLinkResult) GetExpiresAt() string { if x == nil { return "" }; return x.ExpiresAt }

type CrmDepthBusinessCardRequest struct {
	AgentId   string `json:"agent_id"`
	TemplateId string `json:"template_id"`
}

type CrmDepthBusinessCardResult struct {
	CardUrl   string `json:"card_url"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthBusinessCardResult) GetCardUrl() string   { if x == nil { return "" }; return x.CardUrl }
func (x *CrmDepthBusinessCardResult) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthListContentBankRequest struct {
	AgentId  string `json:"agent_id"`
	Category string `json:"category"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthContentAssetRow struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	AssetUrl  string `json:"asset_url"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthContentAssetRow) GetId() string        { if x == nil { return "" }; return x.Id }
func (x *CrmDepthContentAssetRow) GetTitle() string     { if x == nil { return "" }; return x.Title }
func (x *CrmDepthContentAssetRow) GetCategory() string  { if x == nil { return "" }; return x.Category }
func (x *CrmDepthContentAssetRow) GetAssetUrl() string  { if x == nil { return "" }; return x.AssetUrl }
func (x *CrmDepthContentAssetRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthListContentBankResult struct {
	Items    []*CrmDepthContentAssetRow `json:"items"`
	Total    int64                      `json:"total"`
	Page     int32                      `json:"page"`
	PageSize int32                      `json:"page_size"`
}

func (x *CrmDepthListContentBankResult) GetItems() []*CrmDepthContentAssetRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListContentBankResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListContentBankResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListContentBankResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthCreateContentAssetRequest struct {
	AgentId  string `json:"agent_id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	AssetUrl string `json:"asset_url"`
}

type CrmDepthWatermarkRequest struct {
	AgentId    string `json:"agent_id"`
	FlyerUrl   string `json:"flyer_url"`
	WatermarkText string `json:"watermark_text"`
}

type CrmDepthWatermarkResult struct {
	WatermarkedUrl string `json:"watermarked_url"`
}

func (x *CrmDepthWatermarkResult) GetWatermarkedUrl() string { if x == nil { return "" }; return x.WatermarkedUrl }

type CrmDepthListGalleryRequest struct {
	PackageId string `json:"package_id"`
	Page      int32  `json:"page"`
	PageSize  int32  `json:"page_size"`
}

type CrmDepthGalleryRow struct {
	Id        string `json:"id"`
	ImageUrl  string `json:"image_url"`
	Caption   string `json:"caption"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthGalleryRow) GetId() string        { if x == nil { return "" }; return x.Id }
func (x *CrmDepthGalleryRow) GetImageUrl() string  { if x == nil { return "" }; return x.ImageUrl }
func (x *CrmDepthGalleryRow) GetCaption() string   { if x == nil { return "" }; return x.Caption }
func (x *CrmDepthGalleryRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthListGalleryResult struct {
	Items    []*CrmDepthGalleryRow `json:"items"`
	Total    int64                 `json:"total"`
	Page     int32                 `json:"page"`
	PageSize int32                 `json:"page_size"`
}

func (x *CrmDepthListGalleryResult) GetItems() []*CrmDepthGalleryRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListGalleryResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListGalleryResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListGalleryResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthSetTrackingCodeRequest struct {
	AgentId      string `json:"agent_id"`
	TrackingCode string `json:"tracking_code"`
	Platform     string `json:"platform"`
}

type CrmDepthTrackingCodeResult struct {
	AgentId      string `json:"agent_id"`
	TrackingCode string `json:"tracking_code"`
	UpdatedAt    string `json:"updated_at"`
}

func (x *CrmDepthTrackingCodeResult) GetAgentId() string      { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthTrackingCodeResult) GetTrackingCode() string { if x == nil { return "" }; return x.TrackingCode }
func (x *CrmDepthTrackingCodeResult) GetUpdatedAt() string    { if x == nil { return "" }; return x.UpdatedAt }

type CrmDepthAdsStatsRequest struct {
	AgentId   string `json:"agent_id"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
}

type CrmDepthAdsStatsResult struct {
	Impressions int64   `json:"impressions"`
	Clicks      int64   `json:"clicks"`
	Spend       float64 `json:"spend"`
	Conversions int64   `json:"conversions"`
	Roas        float64 `json:"roas"`
}

func (x *CrmDepthAdsStatsResult) GetImpressions() int64   { if x == nil { return 0 }; return x.Impressions }
func (x *CrmDepthAdsStatsResult) GetClicks() int64        { if x == nil { return 0 }; return x.Clicks }
func (x *CrmDepthAdsStatsResult) GetSpend() float64       { if x == nil { return 0 }; return x.Spend }
func (x *CrmDepthAdsStatsResult) GetConversions() int64   { if x == nil { return 0 }; return x.Conversions }
func (x *CrmDepthAdsStatsResult) GetRoas() float64        { if x == nil { return 0 }; return x.Roas }

type CrmDepthCreateUtmLinkRequest struct {
	AgentId    string `json:"agent_id"`
	DestUrl    string `json:"dest_url"`
	UtmSource  string `json:"utm_source"`
	UtmMedium  string `json:"utm_medium"`
	UtmCampaign string `json:"utm_campaign"`
}

type CrmDepthUtmLinkRow struct {
	Id          string `json:"id"`
	AgentId     string `json:"agent_id"`
	DestUrl     string `json:"dest_url"`
	UtmSource   string `json:"utm_source"`
	UtmMedium   string `json:"utm_medium"`
	UtmCampaign string `json:"utm_campaign"`
	ShortUrl    string `json:"short_url"`
	CreatedAt   string `json:"created_at"`
}

func (x *CrmDepthUtmLinkRow) GetId() string          { if x == nil { return "" }; return x.Id }
func (x *CrmDepthUtmLinkRow) GetAgentId() string     { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthUtmLinkRow) GetDestUrl() string     { if x == nil { return "" }; return x.DestUrl }
func (x *CrmDepthUtmLinkRow) GetUtmSource() string   { if x == nil { return "" }; return x.UtmSource }
func (x *CrmDepthUtmLinkRow) GetUtmMedium() string   { if x == nil { return "" }; return x.UtmMedium }
func (x *CrmDepthUtmLinkRow) GetUtmCampaign() string { if x == nil { return "" }; return x.UtmCampaign }
func (x *CrmDepthUtmLinkRow) GetShortUrl() string    { if x == nil { return "" }; return x.ShortUrl }
func (x *CrmDepthUtmLinkRow) GetCreatedAt() string   { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthListUtmLinksRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthListUtmLinksResult struct {
	Items    []*CrmDepthUtmLinkRow `json:"items"`
	Total    int64                 `json:"total"`
	Page     int32                 `json:"page"`
	PageSize int32                 `json:"page_size"`
}

func (x *CrmDepthListUtmLinksResult) GetItems() []*CrmDepthUtmLinkRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListUtmLinksResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListUtmLinksResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListUtmLinksResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthCreateLandingPageRequest struct {
	AgentId     string `json:"agent_id"`
	Title       string `json:"title"`
	TemplateId  string `json:"template_id"`
	PackageId   string `json:"package_id"`
}

type CrmDepthLandingPageRow struct {
	Id        string `json:"id"`
	AgentId   string `json:"agent_id"`
	Title     string `json:"title"`
	PageUrl   string `json:"page_url"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthLandingPageRow) GetId() string        { if x == nil { return "" }; return x.Id }
func (x *CrmDepthLandingPageRow) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthLandingPageRow) GetTitle() string     { if x == nil { return "" }; return x.Title }
func (x *CrmDepthLandingPageRow) GetPageUrl() string   { if x == nil { return "" }; return x.PageUrl }
func (x *CrmDepthLandingPageRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthListLandingPagesRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthListLandingPagesResult struct {
	Items    []*CrmDepthLandingPageRow `json:"items"`
	Total    int64                     `json:"total"`
	Page     int32                     `json:"page"`
	PageSize int32                     `json:"page_size"`
}

func (x *CrmDepthListLandingPagesResult) GetItems() []*CrmDepthLandingPageRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListLandingPagesResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListLandingPagesResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListLandingPagesResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthScheduleContentRequest struct {
	AgentId     string `json:"agent_id"`
	ContentId   string `json:"content_id"`
	Platform    string `json:"platform"`
	ScheduledAt string `json:"scheduled_at"`
}

type CrmDepthScheduledContentRow struct {
	Id          string `json:"id"`
	ContentId   string `json:"content_id"`
	Platform    string `json:"platform"`
	ScheduledAt string `json:"scheduled_at"`
	Status      string `json:"status"`
}

func (x *CrmDepthScheduledContentRow) GetId() string          { if x == nil { return "" }; return x.Id }
func (x *CrmDepthScheduledContentRow) GetContentId() string   { if x == nil { return "" }; return x.ContentId }
func (x *CrmDepthScheduledContentRow) GetPlatform() string    { if x == nil { return "" }; return x.Platform }
func (x *CrmDepthScheduledContentRow) GetScheduledAt() string { if x == nil { return "" }; return x.ScheduledAt }
func (x *CrmDepthScheduledContentRow) GetStatus() string      { if x == nil { return "" }; return x.Status }

type CrmDepthScheduleContentResult struct {
	ScheduleId string `json:"schedule_id"`
	Status     string `json:"status"`
}

func (x *CrmDepthScheduleContentResult) GetScheduleId() string { if x == nil { return "" }; return x.ScheduleId }
func (x *CrmDepthScheduleContentResult) GetStatus() string     { if x == nil { return "" }; return x.Status }

type CrmDepthListScheduledContentRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthListScheduledContentResult struct {
	Items    []*CrmDepthScheduledContentRow `json:"items"`
	Total    int64                          `json:"total"`
	Page     int32                          `json:"page"`
	PageSize int32                          `json:"page_size"`
}

func (x *CrmDepthListScheduledContentResult) GetItems() []*CrmDepthScheduledContentRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListScheduledContentResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListScheduledContentResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListScheduledContentResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthContentAnalyticsRequest struct {
	AgentId  string `json:"agent_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type CrmDepthContentAnalyticsResult struct {
	TotalViews   int64   `json:"total_views"`
	TotalShares  int64   `json:"total_shares"`
	TotalClicks  int64   `json:"total_clicks"`
	EngagementRate float64 `json:"engagement_rate"`
}

func (x *CrmDepthContentAnalyticsResult) GetTotalViews() int64      { if x == nil { return 0 }; return x.TotalViews }
func (x *CrmDepthContentAnalyticsResult) GetTotalShares() int64     { if x == nil { return 0 }; return x.TotalShares }
func (x *CrmDepthContentAnalyticsResult) GetTotalClicks() int64     { if x == nil { return 0 }; return x.TotalClicks }
func (x *CrmDepthContentAnalyticsResult) GetEngagementRate() float64 { if x == nil { return 0 }; return x.EngagementRate }

// ---------------------------------------------------------------------------
// CRM/Leads messages
// ---------------------------------------------------------------------------

type CrmDepthCreateAgentLeadRequest struct {
	AgentId   string `json:"agent_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Source    string `json:"source"`
	PackageId string `json:"package_id"`
	Notes     string `json:"notes"`
}

type CrmDepthAgentLeadRow struct {
	Id        string `json:"id"`
	AgentId   string `json:"agent_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (x *CrmDepthAgentLeadRow) GetId() string        { if x == nil { return "" }; return x.Id }
func (x *CrmDepthAgentLeadRow) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthAgentLeadRow) GetName() string      { if x == nil { return "" }; return x.Name }
func (x *CrmDepthAgentLeadRow) GetPhone() string     { if x == nil { return "" }; return x.Phone }
func (x *CrmDepthAgentLeadRow) GetEmail() string     { if x == nil { return "" }; return x.Email }
func (x *CrmDepthAgentLeadRow) GetStatus() string    { if x == nil { return "" }; return x.Status }
func (x *CrmDepthAgentLeadRow) GetSource() string    { if x == nil { return "" }; return x.Source }
func (x *CrmDepthAgentLeadRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }
func (x *CrmDepthAgentLeadRow) GetUpdatedAt() string { if x == nil { return "" }; return x.UpdatedAt }

type CrmDepthListAgentLeadsRequest struct {
	AgentId      string `json:"agent_id"`
	StatusFilter string `json:"status_filter"`
	Page         int32  `json:"page"`
	PageSize     int32  `json:"page_size"`
}

type CrmDepthListAgentLeadsResult struct {
	Items    []*CrmDepthAgentLeadRow `json:"items"`
	Total    int64                   `json:"total"`
	Page     int32                   `json:"page"`
	PageSize int32                   `json:"page_size"`
}

func (x *CrmDepthListAgentLeadsResult) GetItems() []*CrmDepthAgentLeadRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListAgentLeadsResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListAgentLeadsResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListAgentLeadsResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthSetLeadReminderRequest struct {
	LeadId      string `json:"lead_id"`
	AgentId     string `json:"agent_id"`
	RemindAt    string `json:"remind_at"`
	Message     string `json:"message"`
}

type CrmDepthReminderRow struct {
	Id       string `json:"id"`
	LeadId   string `json:"lead_id"`
	RemindAt string `json:"remind_at"`
	Message  string `json:"message"`
	Status   string `json:"status"`
}

func (x *CrmDepthReminderRow) GetId() string       { if x == nil { return "" }; return x.Id }
func (x *CrmDepthReminderRow) GetLeadId() string   { if x == nil { return "" }; return x.LeadId }
func (x *CrmDepthReminderRow) GetRemindAt() string { if x == nil { return "" }; return x.RemindAt }
func (x *CrmDepthReminderRow) GetMessage() string  { if x == nil { return "" }; return x.Message }
func (x *CrmDepthReminderRow) GetStatus() string   { if x == nil { return "" }; return x.Status }

type CrmDepthSetLeadReminderResult struct {
	ReminderId string `json:"reminder_id"`
	Status     string `json:"status"`
}

func (x *CrmDepthSetLeadReminderResult) GetReminderId() string { if x == nil { return "" }; return x.ReminderId }
func (x *CrmDepthSetLeadReminderResult) GetStatus() string     { if x == nil { return "" }; return x.Status }

type CrmDepthListLeadRemindersRequest struct {
	LeadId string `json:"lead_id"`
}

type CrmDepthListLeadRemindersResult struct {
	Items []*CrmDepthReminderRow `json:"items"`
}

func (x *CrmDepthListLeadRemindersResult) GetItems() []*CrmDepthReminderRow {
	if x == nil { return nil }; return x.Items
}

type CrmDepthFilterBotLeadsRequest struct {
	AgentId string `json:"agent_id"`
}

type CrmDepthFilterBotLeadsResult struct {
	BotLeadIds []string `json:"bot_lead_ids"`
	FilteredCount int32 `json:"filtered_count"`
}

func (x *CrmDepthFilterBotLeadsResult) GetBotLeadIds() []string     { if x == nil { return nil }; return x.BotLeadIds }
func (x *CrmDepthFilterBotLeadsResult) GetFilteredCount() int32     { if x == nil { return 0 }; return x.FilteredCount }

type CrmDepthCreateDripSequenceRequest struct {
	AgentId  string   `json:"agent_id"`
	Name     string   `json:"name"`
	Steps    []string `json:"steps"`
}

type CrmDepthDripSequenceRow struct {
	Id        string `json:"id"`
	AgentId   string `json:"agent_id"`
	Name      string `json:"name"`
	StepCount int32  `json:"step_count"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthDripSequenceRow) GetId() string        { if x == nil { return "" }; return x.Id }
func (x *CrmDepthDripSequenceRow) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthDripSequenceRow) GetName() string      { if x == nil { return "" }; return x.Name }
func (x *CrmDepthDripSequenceRow) GetStepCount() int32  { if x == nil { return 0 }; return x.StepCount }
func (x *CrmDepthDripSequenceRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthCreateDripSequenceResult struct {
	SequenceId string `json:"sequence_id"`
	Status     string `json:"status"`
}

func (x *CrmDepthCreateDripSequenceResult) GetSequenceId() string { if x == nil { return "" }; return x.SequenceId }
func (x *CrmDepthCreateDripSequenceResult) GetStatus() string     { if x == nil { return "" }; return x.Status }

type CrmDepthListDripSequencesRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthListDripSequencesResult struct {
	Items    []*CrmDepthDripSequenceRow `json:"items"`
	Total    int64                      `json:"total"`
	Page     int32                      `json:"page"`
	PageSize int32                      `json:"page_size"`
}

func (x *CrmDepthListDripSequencesResult) GetItems() []*CrmDepthDripSequenceRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListDripSequencesResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListDripSequencesResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListDripSequencesResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthCreateMomentTriggerRequest struct {
	AgentId   string `json:"agent_id"`
	EventType string `json:"event_type"`
	Action    string `json:"action"`
	Condition string `json:"condition"`
}

type CrmDepthMomentTriggerResult struct {
	TriggerId string `json:"trigger_id"`
	Status    string `json:"status"`
}

func (x *CrmDepthMomentTriggerResult) GetTriggerId() string { if x == nil { return "" }; return x.TriggerId }
func (x *CrmDepthMomentTriggerResult) GetStatus() string    { if x == nil { return "" }; return x.Status }

type CrmDepthCreateSegmentRequest struct {
	AgentId  string   `json:"agent_id"`
	Name     string   `json:"name"`
	Criteria []string `json:"criteria"`
}

type CrmDepthSegmentRow struct {
	Id        string `json:"id"`
	AgentId   string `json:"agent_id"`
	Name      string `json:"name"`
	LeadCount int32  `json:"lead_count"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthSegmentRow) GetId() string        { if x == nil { return "" }; return x.Id }
func (x *CrmDepthSegmentRow) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthSegmentRow) GetName() string      { if x == nil { return "" }; return x.Name }
func (x *CrmDepthSegmentRow) GetLeadCount() int32  { if x == nil { return 0 }; return x.LeadCount }
func (x *CrmDepthSegmentRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthCreateSegmentResult struct {
	SegmentId string `json:"segment_id"`
	Status    string `json:"status"`
}

func (x *CrmDepthCreateSegmentResult) GetSegmentId() string { if x == nil { return "" }; return x.SegmentId }
func (x *CrmDepthCreateSegmentResult) GetStatus() string    { if x == nil { return "" }; return x.Status }

type CrmDepthListSegmentsRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthListSegmentsResult struct {
	Items    []*CrmDepthSegmentRow `json:"items"`
	Total    int64                 `json:"total"`
	Page     int32                 `json:"page"`
	PageSize int32                 `json:"page_size"`
}

func (x *CrmDepthListSegmentsResult) GetItems() []*CrmDepthSegmentRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListSegmentsResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListSegmentsResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListSegmentsResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthSendBroadcastRequest struct {
	AgentId    string   `json:"agent_id"`
	SegmentIds []string `json:"segment_ids"`
	Message    string   `json:"message"`
	Channel    string   `json:"channel"`
}

type CrmDepthSendBroadcastResult struct {
	BroadcastId    string `json:"broadcast_id"`
	RecipientCount int32  `json:"recipient_count"`
	Status         string `json:"status"`
}

func (x *CrmDepthSendBroadcastResult) GetBroadcastId() string    { if x == nil { return "" }; return x.BroadcastId }
func (x *CrmDepthSendBroadcastResult) GetRecipientCount() int32  { if x == nil { return 0 }; return x.RecipientCount }
func (x *CrmDepthSendBroadcastResult) GetStatus() string         { if x == nil { return "" }; return x.Status }

type CrmDepthAssignLeadFairRequest struct {
	AgentId string   `json:"agent_id"`
	LeadIds []string `json:"lead_ids"`
}

type CrmDepthAssignLeadFairResult struct {
	AssignedCount int32  `json:"assigned_count"`
	Status        string `json:"status"`
}

func (x *CrmDepthAssignLeadFairResult) GetAssignedCount() int32  { if x == nil { return 0 }; return x.AssignedCount }
func (x *CrmDepthAssignLeadFairResult) GetStatus() string        { if x == nil { return "" }; return x.Status }

type CrmDepthCreateSlaRuleRequest struct {
	AgentId         string `json:"agent_id"`
	RuleName        string `json:"rule_name"`
	ResponseTimeMins int32  `json:"response_time_mins"`
	EscalationLevel  int32  `json:"escalation_level"`
}

type CrmDepthSlaRuleResult struct {
	RuleId    string `json:"rule_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthSlaRuleResult) GetRuleId() string    { if x == nil { return "" }; return x.RuleId }
func (x *CrmDepthSlaRuleResult) GetStatus() string    { if x == nil { return "" }; return x.Status }
func (x *CrmDepthSlaRuleResult) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthGetLeadTrailRequest struct {
	LeadId string `json:"lead_id"`
}

type CrmDepthLeadTrailRow struct {
	EventType string `json:"event_type"`
	Notes     string `json:"notes"`
	ActorId   string `json:"actor_id"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthLeadTrailRow) GetEventType() string { if x == nil { return "" }; return x.EventType }
func (x *CrmDepthLeadTrailRow) GetNotes() string     { if x == nil { return "" }; return x.Notes }
func (x *CrmDepthLeadTrailRow) GetActorId() string   { if x == nil { return "" }; return x.ActorId }
func (x *CrmDepthLeadTrailRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthGetLeadTrailResult struct {
	LeadId string                 `json:"lead_id"`
	Events []*CrmDepthLeadTrailRow `json:"events"`
}

func (x *CrmDepthGetLeadTrailResult) GetLeadId() string { if x == nil { return "" }; return x.LeadId }
func (x *CrmDepthGetLeadTrailResult) GetEvents() []*CrmDepthLeadTrailRow {
	if x == nil { return nil }; return x.Events
}

type CrmDepthTagLeadRequest struct {
	LeadId string   `json:"lead_id"`
	Tags   []string `json:"tags"`
}

type CrmDepthTagLeadResult struct {
	LeadId    string   `json:"lead_id"`
	Tags      []string `json:"tags"`
	UpdatedAt string   `json:"updated_at"`
}

func (x *CrmDepthTagLeadResult) GetLeadId() string    { if x == nil { return "" }; return x.LeadId }
func (x *CrmDepthTagLeadResult) GetTags() []string    { if x == nil { return nil }; return x.Tags }
func (x *CrmDepthTagLeadResult) GetUpdatedAt() string { if x == nil { return "" }; return x.UpdatedAt }

type CrmDepthGenerateQuoteRequest struct {
	LeadId    string  `json:"lead_id"`
	PackageId string  `json:"package_id"`
	Pax       int32   `json:"pax"`
	Discount  float64 `json:"discount"`
}

type CrmDepthQuoteResult struct {
	QuoteId    string  `json:"quote_id"`
	LeadId     string  `json:"lead_id"`
	PackageId  string  `json:"package_id"`
	TotalPrice int64   `json:"total_price"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	ExpiresAt  string  `json:"expires_at"`
}

func (x *CrmDepthQuoteResult) GetQuoteId() string    { if x == nil { return "" }; return x.QuoteId }
func (x *CrmDepthQuoteResult) GetLeadId() string     { if x == nil { return "" }; return x.LeadId }
func (x *CrmDepthQuoteResult) GetPackageId() string  { if x == nil { return "" }; return x.PackageId }
func (x *CrmDepthQuoteResult) GetTotalPrice() int64  { if x == nil { return 0 }; return x.TotalPrice }
func (x *CrmDepthQuoteResult) GetStatus() string     { if x == nil { return "" }; return x.Status }
func (x *CrmDepthQuoteResult) GetCreatedAt() string  { if x == nil { return "" }; return x.CreatedAt }
func (x *CrmDepthQuoteResult) GetExpiresAt() string  { if x == nil { return "" }; return x.ExpiresAt }

type CrmDepthGetQuoteRequest struct {
	QuoteId string `json:"quote_id"`
}

type CrmDepthBuildPaymentLinkRequest struct {
	QuoteId  string `json:"quote_id"`
	LeadId   string `json:"lead_id"`
	Notes    string `json:"notes"`
}

type CrmDepthPaymentLinkResult struct {
	LinkId      string `json:"link_id"`
	PaymentUrl  string `json:"payment_url"`
	ExpiresAt   string `json:"expires_at"`
}

func (x *CrmDepthPaymentLinkResult) GetLinkId() string     { if x == nil { return "" }; return x.LinkId }
func (x *CrmDepthPaymentLinkResult) GetPaymentUrl() string { if x == nil { return "" }; return x.PaymentUrl }
func (x *CrmDepthPaymentLinkResult) GetExpiresAt() string  { if x == nil { return "" }; return x.ExpiresAt }

type CrmDepthRequestDiscountRequest struct {
	QuoteId  string  `json:"quote_id"`
	AgentId  string  `json:"agent_id"`
	Amount   float64 `json:"amount"`
	Reason   string  `json:"reason"`
}

type CrmDepthDiscountResult struct {
	DiscountId string  `json:"discount_id"`
	QuoteId    string  `json:"quote_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
}

func (x *CrmDepthDiscountResult) GetDiscountId() string { if x == nil { return "" }; return x.DiscountId }
func (x *CrmDepthDiscountResult) GetQuoteId() string    { if x == nil { return "" }; return x.QuoteId }
func (x *CrmDepthDiscountResult) GetAmount() float64    { if x == nil { return 0 }; return x.Amount }
func (x *CrmDepthDiscountResult) GetStatus() string     { if x == nil { return "" }; return x.Status }
func (x *CrmDepthDiscountResult) GetCreatedAt() string  { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthApproveDiscountRequest struct {
	DiscountId string `json:"discount_id"`
	Decision   string `json:"decision"`
	Notes      string `json:"notes"`
}

type CrmDepthStaleProspectsRequest struct {
	AgentId      string `json:"agent_id"`
	StaleDays    int32  `json:"stale_days"`
	Page         int32  `json:"page"`
	PageSize     int32  `json:"page_size"`
}

type CrmDepthStaleProspectsResult struct {
	Items    []*CrmDepthAgentLeadRow `json:"items"`
	Total    int64                   `json:"total"`
	Page     int32                   `json:"page"`
	PageSize int32                   `json:"page_size"`
}

func (x *CrmDepthStaleProspectsResult) GetItems() []*CrmDepthAgentLeadRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthStaleProspectsResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthStaleProspectsResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthStaleProspectsResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

// ---------------------------------------------------------------------------
// Commission messages
// ---------------------------------------------------------------------------

type CrmDepthGetCommissionRequest struct {
	AgentId string `json:"agent_id"`
}

type CrmDepthCommissionBalanceResult struct {
	AgentId       string  `json:"agent_id"`
	BalanceIdr    int64   `json:"balance_idr"`
	PendingIdr    int64   `json:"pending_idr"`
	TotalEarnedIdr int64  `json:"total_earned_idr"`
}

func (x *CrmDepthCommissionBalanceResult) GetAgentId() string        { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthCommissionBalanceResult) GetBalanceIdr() int64      { if x == nil { return 0 }; return x.BalanceIdr }
func (x *CrmDepthCommissionBalanceResult) GetPendingIdr() int64      { if x == nil { return 0 }; return x.PendingIdr }
func (x *CrmDepthCommissionBalanceResult) GetTotalEarnedIdr() int64  { if x == nil { return 0 }; return x.TotalEarnedIdr }

type CrmDepthCommissionEventRow struct {
	EventId    string  `json:"event_id"`
	EventType  string  `json:"event_type"`
	AmountIdr  int64   `json:"amount_idr"`
	ReferenceId string `json:"reference_id"`
	CreatedAt  string  `json:"created_at"`
}

func (x *CrmDepthCommissionEventRow) GetEventId() string     { if x == nil { return "" }; return x.EventId }
func (x *CrmDepthCommissionEventRow) GetEventType() string   { if x == nil { return "" }; return x.EventType }
func (x *CrmDepthCommissionEventRow) GetAmountIdr() int64    { if x == nil { return 0 }; return x.AmountIdr }
func (x *CrmDepthCommissionEventRow) GetReferenceId() string { if x == nil { return "" }; return x.ReferenceId }
func (x *CrmDepthCommissionEventRow) GetCreatedAt() string   { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthGetCommissionEventsResult struct {
	AgentId string                        `json:"agent_id"`
	Events  []*CrmDepthCommissionEventRow `json:"events"`
	Total   int64                         `json:"total"`
}

func (x *CrmDepthGetCommissionEventsResult) GetAgentId() string { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthGetCommissionEventsResult) GetEvents() []*CrmDepthCommissionEventRow {
	if x == nil { return nil }; return x.Events
}
func (x *CrmDepthGetCommissionEventsResult) GetTotal() int64 { if x == nil { return 0 }; return x.Total }

type CrmDepthRequestPayoutRequest struct {
	AgentId   string `json:"agent_id"`
	AmountIdr int64  `json:"amount_idr"`
	BankCode  string `json:"bank_code"`
	AccountNo string `json:"account_no"`
}

type CrmDepthPayoutResult struct {
	PayoutId  string `json:"payout_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthPayoutResult) GetPayoutId() string  { if x == nil { return "" }; return x.PayoutId }
func (x *CrmDepthPayoutResult) GetStatus() string    { if x == nil { return "" }; return x.Status }
func (x *CrmDepthPayoutResult) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthPayoutHistoryRow struct {
	PayoutId  string `json:"payout_id"`
	AmountIdr int64  `json:"amount_idr"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (x *CrmDepthPayoutHistoryRow) GetPayoutId() string  { if x == nil { return "" }; return x.PayoutId }
func (x *CrmDepthPayoutHistoryRow) GetAmountIdr() int64  { if x == nil { return 0 }; return x.AmountIdr }
func (x *CrmDepthPayoutHistoryRow) GetStatus() string    { if x == nil { return "" }; return x.Status }
func (x *CrmDepthPayoutHistoryRow) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthGetPayoutHistoryResult struct {
	AgentId string                    `json:"agent_id"`
	Items   []*CrmDepthPayoutHistoryRow `json:"items"`
	Total   int64                     `json:"total"`
}

func (x *CrmDepthGetPayoutHistoryResult) GetAgentId() string { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthGetPayoutHistoryResult) GetItems() []*CrmDepthPayoutHistoryRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthGetPayoutHistoryResult) GetTotal() int64 { if x == nil { return 0 }; return x.Total }

type CrmDepthComputeOverrideRequest struct {
	AgentId    string  `json:"agent_id"`
	DownlineId string  `json:"downline_id"`
	SaleAmount int64   `json:"sale_amount"`
}

type CrmDepthOverrideCommissionResult struct {
	OverrideAmount int64   `json:"override_amount"`
	Rate           float64 `json:"rate"`
	Tier           string  `json:"tier"`
}

func (x *CrmDepthOverrideCommissionResult) GetOverrideAmount() int64   { if x == nil { return 0 }; return x.OverrideAmount }
func (x *CrmDepthOverrideCommissionResult) GetRate() float64           { if x == nil { return 0 }; return x.Rate }
func (x *CrmDepthOverrideCommissionResult) GetTier() string            { if x == nil { return "" }; return x.Tier }

type CrmDepthRoasReportRequest struct {
	AgentId  string `json:"agent_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type CrmDepthRoasReportResult struct {
	AgentId        string  `json:"agent_id"`
	TotalAdSpend   float64 `json:"total_ad_spend"`
	TotalRevenue   int64   `json:"total_revenue"`
	Roas           float64 `json:"roas"`
	ConversionRate float64 `json:"conversion_rate"`
}

func (x *CrmDepthRoasReportResult) GetAgentId() string        { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthRoasReportResult) GetTotalAdSpend() float64  { if x == nil { return 0 }; return x.TotalAdSpend }
func (x *CrmDepthRoasReportResult) GetTotalRevenue() int64    { if x == nil { return 0 }; return x.TotalRevenue }
func (x *CrmDepthRoasReportResult) GetRoas() float64          { if x == nil { return 0 }; return x.Roas }
func (x *CrmDepthRoasReportResult) GetConversionRate() float64 { if x == nil { return 0 }; return x.ConversionRate }

// ---------------------------------------------------------------------------
// Academy messages
// ---------------------------------------------------------------------------

type CrmDepthListCoursesRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthCourseRow struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	DurationMin int32  `json:"duration_min"`
}

func (x *CrmDepthCourseRow) GetId() string          { if x == nil { return "" }; return x.Id }
func (x *CrmDepthCourseRow) GetTitle() string       { if x == nil { return "" }; return x.Title }
func (x *CrmDepthCourseRow) GetDescription() string { if x == nil { return "" }; return x.Description }
func (x *CrmDepthCourseRow) GetCategory() string    { if x == nil { return "" }; return x.Category }
func (x *CrmDepthCourseRow) GetDurationMin() int32  { if x == nil { return 0 }; return x.DurationMin }

type CrmDepthListCoursesResult struct {
	Items    []*CrmDepthCourseRow `json:"items"`
	Total    int64                `json:"total"`
	Page     int32                `json:"page"`
	PageSize int32                `json:"page_size"`
}

func (x *CrmDepthListCoursesResult) GetItems() []*CrmDepthCourseRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListCoursesResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListCoursesResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListCoursesResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthGetCourseProgressRequest struct {
	AgentId  string `json:"agent_id"`
	CourseId string `json:"course_id"`
}

type CrmDepthCourseProgressResult struct {
	AgentId         string  `json:"agent_id"`
	CourseId        string  `json:"course_id"`
	ProgressPercent float64 `json:"progress_percent"`
	CompletedAt     string  `json:"completed_at"`
	Status          string  `json:"status"`
}

func (x *CrmDepthCourseProgressResult) GetAgentId() string         { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthCourseProgressResult) GetCourseId() string        { if x == nil { return "" }; return x.CourseId }
func (x *CrmDepthCourseProgressResult) GetProgressPercent() float64 { if x == nil { return 0 }; return x.ProgressPercent }
func (x *CrmDepthCourseProgressResult) GetCompletedAt() string     { if x == nil { return "" }; return x.CompletedAt }
func (x *CrmDepthCourseProgressResult) GetStatus() string          { if x == nil { return "" }; return x.Status }

type CrmDepthSubmitQuizRequest struct {
	AgentId  string   `json:"agent_id"`
	CourseId string   `json:"course_id"`
	Answers  []string `json:"answers"`
}

type CrmDepthQuizResult struct {
	AgentId   string  `json:"agent_id"`
	CourseId  string  `json:"course_id"`
	Score     float64 `json:"score"`
	Passed    bool    `json:"passed"`
	GradedAt  string  `json:"graded_at"`
}

func (x *CrmDepthQuizResult) GetAgentId() string  { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthQuizResult) GetCourseId() string { if x == nil { return "" }; return x.CourseId }
func (x *CrmDepthQuizResult) GetScore() float64   { if x == nil { return 0 }; return x.Score }
func (x *CrmDepthQuizResult) GetPassed() bool     { if x == nil { return false }; return x.Passed }
func (x *CrmDepthQuizResult) GetGradedAt() string { if x == nil { return "" }; return x.GradedAt }

type CrmDepthListSalesScriptsRequest struct {
	Category string `json:"category"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthSalesScriptRow struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

func (x *CrmDepthSalesScriptRow) GetId() string       { if x == nil { return "" }; return x.Id }
func (x *CrmDepthSalesScriptRow) GetTitle() string    { if x == nil { return "" }; return x.Title }
func (x *CrmDepthSalesScriptRow) GetCategory() string { if x == nil { return "" }; return x.Category }
func (x *CrmDepthSalesScriptRow) GetContent() string  { if x == nil { return "" }; return x.Content }

type CrmDepthListSalesScriptsResult struct {
	Items    []*CrmDepthSalesScriptRow `json:"items"`
	Total    int64                     `json:"total"`
	Page     int32                     `json:"page"`
	PageSize int32                     `json:"page_size"`
}

func (x *CrmDepthListSalesScriptsResult) GetItems() []*CrmDepthSalesScriptRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthListSalesScriptsResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthListSalesScriptsResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthListSalesScriptsResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthLeaderboardRequest struct {
	Period   string `json:"period"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type CrmDepthLeaderboardRow struct {
	Rank      int32  `json:"rank"`
	AgentId   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
	Sales     int64  `json:"sales"`
	Points    int64  `json:"points"`
}

func (x *CrmDepthLeaderboardRow) GetRank() int32      { if x == nil { return 0 }; return x.Rank }
func (x *CrmDepthLeaderboardRow) GetAgentId() string  { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthLeaderboardRow) GetAgentName() string { if x == nil { return "" }; return x.AgentName }
func (x *CrmDepthLeaderboardRow) GetSales() int64     { if x == nil { return 0 }; return x.Sales }
func (x *CrmDepthLeaderboardRow) GetPoints() int64    { if x == nil { return 0 }; return x.Points }

type CrmDepthLeaderboardResult struct {
	Items    []*CrmDepthLeaderboardRow `json:"items"`
	Total    int64                     `json:"total"`
	Page     int32                     `json:"page"`
	PageSize int32                     `json:"page_size"`
}

func (x *CrmDepthLeaderboardResult) GetItems() []*CrmDepthLeaderboardRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthLeaderboardResult) GetTotal() int64    { if x == nil { return 0 }; return x.Total }
func (x *CrmDepthLeaderboardResult) GetPage() int32     { if x == nil { return 0 }; return x.Page }
func (x *CrmDepthLeaderboardResult) GetPageSize() int32 { if x == nil { return 0 }; return x.PageSize }

type CrmDepthAgentTierResult struct {
	AgentId   string `json:"agent_id"`
	TierLevel string `json:"tier_level"`
	TierName  string `json:"tier_name"`
	Points    int64  `json:"points"`
	NextTier  string `json:"next_tier"`
}

func (x *CrmDepthAgentTierResult) GetAgentId() string   { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthAgentTierResult) GetTierLevel() string { if x == nil { return "" }; return x.TierLevel }
func (x *CrmDepthAgentTierResult) GetTierName() string  { if x == nil { return "" }; return x.TierName }
func (x *CrmDepthAgentTierResult) GetPoints() int64     { if x == nil { return 0 }; return x.Points }
func (x *CrmDepthAgentTierResult) GetNextTier() string  { if x == nil { return "" }; return x.NextTier }

// ---------------------------------------------------------------------------
// Alumni/Other messages
// ---------------------------------------------------------------------------

type CrmDepthReferralRow struct {
	ReferralId  string `json:"referral_id"`
	RefereeName string `json:"referee_name"`
	Status      string `json:"status"`
	Reward      int64  `json:"reward"`
	CreatedAt   string `json:"created_at"`
}

func (x *CrmDepthReferralRow) GetReferralId() string  { if x == nil { return "" }; return x.ReferralId }
func (x *CrmDepthReferralRow) GetRefereeName() string { if x == nil { return "" }; return x.RefereeName }
func (x *CrmDepthReferralRow) GetStatus() string      { if x == nil { return "" }; return x.Status }
func (x *CrmDepthReferralRow) GetReward() int64       { if x == nil { return 0 }; return x.Reward }
func (x *CrmDepthReferralRow) GetCreatedAt() string   { if x == nil { return "" }; return x.CreatedAt }

type CrmDepthGetAlumniReferralsResult struct {
	AgentId string                 `json:"agent_id"`
	Items   []*CrmDepthReferralRow `json:"items"`
	Total   int64                  `json:"total"`
}

func (x *CrmDepthGetAlumniReferralsResult) GetAgentId() string { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthGetAlumniReferralsResult) GetItems() []*CrmDepthReferralRow {
	if x == nil { return nil }; return x.Items
}
func (x *CrmDepthGetAlumniReferralsResult) GetTotal() int64 { if x == nil { return 0 }; return x.Total }

type CrmDepthReturnIntentResult struct {
	AgentId          string `json:"agent_id"`
	EligibleLeads    int32  `json:"eligible_leads"`
	PotentialSavings int64  `json:"potential_savings"`
}

func (x *CrmDepthReturnIntentResult) GetAgentId() string          { if x == nil { return "" }; return x.AgentId }
func (x *CrmDepthReturnIntentResult) GetEligibleLeads() int32     { if x == nil { return 0 }; return x.EligibleLeads }
func (x *CrmDepthReturnIntentResult) GetPotentialSavings() int64  { if x == nil { return 0 }; return x.PotentialSavings }

type CrmDepthCalculateZakatRequest struct {
	IncomeSources []string `json:"income_sources"`
	TotalAssets   int64    `json:"total_assets"`
	TotalLiabilities int64 `json:"total_liabilities"`
}

type CrmDepthZakatResult struct {
	ZakatAmount int64   `json:"zakat_amount"`
	NisabValue  int64   `json:"nisab_value"`
	Rate        float64 `json:"rate"`
	Eligible    bool    `json:"eligible"`
}

func (x *CrmDepthZakatResult) GetZakatAmount() int64   { if x == nil { return 0 }; return x.ZakatAmount }
func (x *CrmDepthZakatResult) GetNisabValue() int64    { if x == nil { return 0 }; return x.NisabValue }
func (x *CrmDepthZakatResult) GetRate() float64        { if x == nil { return 0 }; return x.Rate }
func (x *CrmDepthZakatResult) GetEligible() bool       { if x == nil { return false }; return x.Eligible }

type CrmDepthRecordCharityRequest struct {
	AgentId    string `json:"agent_id"`
	AmountIdr  int64  `json:"amount_idr"`
	CharityOrg string `json:"charity_org"`
	Notes      string `json:"notes"`
}

type CrmDepthCharityResult struct {
	CharityId string `json:"charity_id"`
	Status    string `json:"status"`
	RecordedAt string `json:"recorded_at"`
}

func (x *CrmDepthCharityResult) GetCharityId() string  { if x == nil { return "" }; return x.CharityId }
func (x *CrmDepthCharityResult) GetStatus() string     { if x == nil { return "" }; return x.Status }
func (x *CrmDepthCharityResult) GetRecordedAt() string { if x == nil { return "" }; return x.RecordedAt }

// ---------------------------------------------------------------------------
// CrmDepthClient — gRPC client interface + implementation
// ---------------------------------------------------------------------------

// CrmDepthClient is the client API for crm-svc Wave 7 depth RPCs.
type CrmDepthClient interface {
	RegisterAgent(ctx context.Context, in *CrmDepthAgentRegisterRequest, opts ...grpc.CallOption) (*CrmDepthAgentResult, error)
	SubmitAgentKyc(ctx context.Context, in *CrmDepthSubmitKycRequest, opts ...grpc.CallOption) (*CrmDepthKycResult, error)
	SignAgentMoU(ctx context.Context, in *CrmDepthSignMoURequest, opts ...grpc.CallOption) (*CrmDepthMoUResult, error)
	GetAgentProfile(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthAgentResult, error)
	GetReplicaSite(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthReplicaSiteResult, error)
	UpdateReplicaSite(ctx context.Context, in *CrmDepthUpdateReplicaSiteRequest, opts ...grpc.CallOption) (*CrmDepthReplicaSiteResult, error)
	GetSocialShareLink(ctx context.Context, in *CrmDepthShareLinkRequest, opts ...grpc.CallOption) (*CrmDepthShareLinkResult, error)
	GenerateBusinessCard(ctx context.Context, in *CrmDepthBusinessCardRequest, opts ...grpc.CallOption) (*CrmDepthBusinessCardResult, error)
	ListContentBank(ctx context.Context, in *CrmDepthListContentBankRequest, opts ...grpc.CallOption) (*CrmDepthListContentBankResult, error)
	CreateContentAsset(ctx context.Context, in *CrmDepthCreateContentAssetRequest, opts ...grpc.CallOption) (*CrmDepthContentAssetRow, error)
	WatermarkFlyer(ctx context.Context, in *CrmDepthWatermarkRequest, opts ...grpc.CallOption) (*CrmDepthWatermarkResult, error)
	ListProgramGallery(ctx context.Context, in *CrmDepthListGalleryRequest, opts ...grpc.CallOption) (*CrmDepthListGalleryResult, error)
	SetTrackingCode(ctx context.Context, in *CrmDepthSetTrackingCodeRequest, opts ...grpc.CallOption) (*CrmDepthTrackingCodeResult, error)
	GetAdsManagerStats(ctx context.Context, in *CrmDepthAdsStatsRequest, opts ...grpc.CallOption) (*CrmDepthAdsStatsResult, error)
	CreateUtmLink(ctx context.Context, in *CrmDepthCreateUtmLinkRequest, opts ...grpc.CallOption) (*CrmDepthUtmLinkRow, error)
	ListUtmLinks(ctx context.Context, in *CrmDepthListUtmLinksRequest, opts ...grpc.CallOption) (*CrmDepthListUtmLinksResult, error)
	CreateLandingPage(ctx context.Context, in *CrmDepthCreateLandingPageRequest, opts ...grpc.CallOption) (*CrmDepthLandingPageRow, error)
	ListLandingPages(ctx context.Context, in *CrmDepthListLandingPagesRequest, opts ...grpc.CallOption) (*CrmDepthListLandingPagesResult, error)
	ScheduleContent(ctx context.Context, in *CrmDepthScheduleContentRequest, opts ...grpc.CallOption) (*CrmDepthScheduleContentResult, error)
	ListScheduledContent(ctx context.Context, in *CrmDepthListScheduledContentRequest, opts ...grpc.CallOption) (*CrmDepthListScheduledContentResult, error)
	GetContentAnalytics(ctx context.Context, in *CrmDepthContentAnalyticsRequest, opts ...grpc.CallOption) (*CrmDepthContentAnalyticsResult, error)
	CreateAgentLead(ctx context.Context, in *CrmDepthCreateAgentLeadRequest, opts ...grpc.CallOption) (*CrmDepthAgentLeadRow, error)
	ListAgentLeads(ctx context.Context, in *CrmDepthListAgentLeadsRequest, opts ...grpc.CallOption) (*CrmDepthListAgentLeadsResult, error)
	SetLeadReminder(ctx context.Context, in *CrmDepthSetLeadReminderRequest, opts ...grpc.CallOption) (*CrmDepthSetLeadReminderResult, error)
	ListLeadReminders(ctx context.Context, in *CrmDepthListLeadRemindersRequest, opts ...grpc.CallOption) (*CrmDepthListLeadRemindersResult, error)
	FilterBotLeads(ctx context.Context, in *CrmDepthFilterBotLeadsRequest, opts ...grpc.CallOption) (*CrmDepthFilterBotLeadsResult, error)
	CreateDripSequence(ctx context.Context, in *CrmDepthCreateDripSequenceRequest, opts ...grpc.CallOption) (*CrmDepthCreateDripSequenceResult, error)
	ListDripSequences(ctx context.Context, in *CrmDepthListDripSequencesRequest, opts ...grpc.CallOption) (*CrmDepthListDripSequencesResult, error)
	CreateMomentTrigger(ctx context.Context, in *CrmDepthCreateMomentTriggerRequest, opts ...grpc.CallOption) (*CrmDepthMomentTriggerResult, error)
	CreateSegment(ctx context.Context, in *CrmDepthCreateSegmentRequest, opts ...grpc.CallOption) (*CrmDepthCreateSegmentResult, error)
	ListSegments(ctx context.Context, in *CrmDepthListSegmentsRequest, opts ...grpc.CallOption) (*CrmDepthListSegmentsResult, error)
	SendBroadcast(ctx context.Context, in *CrmDepthSendBroadcastRequest, opts ...grpc.CallOption) (*CrmDepthSendBroadcastResult, error)
	AssignLeadFair(ctx context.Context, in *CrmDepthAssignLeadFairRequest, opts ...grpc.CallOption) (*CrmDepthAssignLeadFairResult, error)
	CreateSlaRule(ctx context.Context, in *CrmDepthCreateSlaRuleRequest, opts ...grpc.CallOption) (*CrmDepthSlaRuleResult, error)
	GetLeadTrail(ctx context.Context, in *CrmDepthGetLeadTrailRequest, opts ...grpc.CallOption) (*CrmDepthGetLeadTrailResult, error)
	TagLead(ctx context.Context, in *CrmDepthTagLeadRequest, opts ...grpc.CallOption) (*CrmDepthTagLeadResult, error)
	GenerateQuote(ctx context.Context, in *CrmDepthGenerateQuoteRequest, opts ...grpc.CallOption) (*CrmDepthQuoteResult, error)
	GetQuote(ctx context.Context, in *CrmDepthGetQuoteRequest, opts ...grpc.CallOption) (*CrmDepthQuoteResult, error)
	BuildPaymentLink(ctx context.Context, in *CrmDepthBuildPaymentLinkRequest, opts ...grpc.CallOption) (*CrmDepthPaymentLinkResult, error)
	RequestDiscount(ctx context.Context, in *CrmDepthRequestDiscountRequest, opts ...grpc.CallOption) (*CrmDepthDiscountResult, error)
	ApproveDiscount(ctx context.Context, in *CrmDepthApproveDiscountRequest, opts ...grpc.CallOption) (*CrmDepthDiscountResult, error)
	GetStaleProspects(ctx context.Context, in *CrmDepthStaleProspectsRequest, opts ...grpc.CallOption) (*CrmDepthStaleProspectsResult, error)
	GetCommissionBalance(ctx context.Context, in *CrmDepthGetCommissionRequest, opts ...grpc.CallOption) (*CrmDepthCommissionBalanceResult, error)
	GetCommissionEvents(ctx context.Context, in *CrmDepthGetCommissionRequest, opts ...grpc.CallOption) (*CrmDepthGetCommissionEventsResult, error)
	RequestPayout(ctx context.Context, in *CrmDepthRequestPayoutRequest, opts ...grpc.CallOption) (*CrmDepthPayoutResult, error)
	GetPayoutHistory(ctx context.Context, in *CrmDepthGetCommissionRequest, opts ...grpc.CallOption) (*CrmDepthGetPayoutHistoryResult, error)
	ComputeOverrideCommission(ctx context.Context, in *CrmDepthComputeOverrideRequest, opts ...grpc.CallOption) (*CrmDepthOverrideCommissionResult, error)
	GetRoasReport(ctx context.Context, in *CrmDepthRoasReportRequest, opts ...grpc.CallOption) (*CrmDepthRoasReportResult, error)
	ListAcademyCourses(ctx context.Context, in *CrmDepthListCoursesRequest, opts ...grpc.CallOption) (*CrmDepthListCoursesResult, error)
	GetCourseProgress(ctx context.Context, in *CrmDepthGetCourseProgressRequest, opts ...grpc.CallOption) (*CrmDepthCourseProgressResult, error)
	SubmitQuiz(ctx context.Context, in *CrmDepthSubmitQuizRequest, opts ...grpc.CallOption) (*CrmDepthQuizResult, error)
	ListSalesScripts(ctx context.Context, in *CrmDepthListSalesScriptsRequest, opts ...grpc.CallOption) (*CrmDepthListSalesScriptsResult, error)
	GetLeaderboard(ctx context.Context, in *CrmDepthLeaderboardRequest, opts ...grpc.CallOption) (*CrmDepthLeaderboardResult, error)
	GetAgentTier(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthAgentTierResult, error)
	GetAlumniReferrals(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthGetAlumniReferralsResult, error)
	GetReturnIntentSavings(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthReturnIntentResult, error)
	CalculateZakat(ctx context.Context, in *CrmDepthCalculateZakatRequest, opts ...grpc.CallOption) (*CrmDepthZakatResult, error)
	RecordCharity(ctx context.Context, in *CrmDepthRecordCharityRequest, opts ...grpc.CallOption) (*CrmDepthCharityResult, error)
}

type crmDepthClient struct{ cc grpc.ClientConnInterface }

// NewCrmDepthClient returns a CrmDepthClient backed by the given conn.
func NewCrmDepthClient(cc grpc.ClientConnInterface) CrmDepthClient {
	return &crmDepthClient{cc}
}

func (c *crmDepthClient) RegisterAgent(ctx context.Context, in *CrmDepthAgentRegisterRequest, opts ...grpc.CallOption) (*CrmDepthAgentResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthAgentResult)
	if err := c.cc.Invoke(ctx, CrmService_RegisterAgent_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) SubmitAgentKyc(ctx context.Context, in *CrmDepthSubmitKycRequest, opts ...grpc.CallOption) (*CrmDepthKycResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthKycResult)
	if err := c.cc.Invoke(ctx, CrmService_SubmitAgentKyc_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) SignAgentMoU(ctx context.Context, in *CrmDepthSignMoURequest, opts ...grpc.CallOption) (*CrmDepthMoUResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthMoUResult)
	if err := c.cc.Invoke(ctx, CrmService_SignAgentMoU_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetAgentProfile(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthAgentResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthAgentResult)
	if err := c.cc.Invoke(ctx, CrmService_GetAgentProfile_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetReplicaSite(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthReplicaSiteResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthReplicaSiteResult)
	if err := c.cc.Invoke(ctx, CrmService_GetReplicaSite_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) UpdateReplicaSite(ctx context.Context, in *CrmDepthUpdateReplicaSiteRequest, opts ...grpc.CallOption) (*CrmDepthReplicaSiteResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthReplicaSiteResult)
	if err := c.cc.Invoke(ctx, CrmService_UpdateReplicaSite_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetSocialShareLink(ctx context.Context, in *CrmDepthShareLinkRequest, opts ...grpc.CallOption) (*CrmDepthShareLinkResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthShareLinkResult)
	if err := c.cc.Invoke(ctx, CrmService_GetSocialShareLink_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GenerateBusinessCard(ctx context.Context, in *CrmDepthBusinessCardRequest, opts ...grpc.CallOption) (*CrmDepthBusinessCardResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthBusinessCardResult)
	if err := c.cc.Invoke(ctx, CrmService_GenerateBusinessCard_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListContentBank(ctx context.Context, in *CrmDepthListContentBankRequest, opts ...grpc.CallOption) (*CrmDepthListContentBankResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListContentBankResult)
	if err := c.cc.Invoke(ctx, CrmService_ListContentBank_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateContentAsset(ctx context.Context, in *CrmDepthCreateContentAssetRequest, opts ...grpc.CallOption) (*CrmDepthContentAssetRow, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthContentAssetRow)
	if err := c.cc.Invoke(ctx, CrmService_CreateContentAsset_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) WatermarkFlyer(ctx context.Context, in *CrmDepthWatermarkRequest, opts ...grpc.CallOption) (*CrmDepthWatermarkResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthWatermarkResult)
	if err := c.cc.Invoke(ctx, CrmService_WatermarkFlyer_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListProgramGallery(ctx context.Context, in *CrmDepthListGalleryRequest, opts ...grpc.CallOption) (*CrmDepthListGalleryResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListGalleryResult)
	if err := c.cc.Invoke(ctx, CrmService_ListProgramGallery_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) SetTrackingCode(ctx context.Context, in *CrmDepthSetTrackingCodeRequest, opts ...grpc.CallOption) (*CrmDepthTrackingCodeResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthTrackingCodeResult)
	if err := c.cc.Invoke(ctx, CrmService_SetTrackingCode_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetAdsManagerStats(ctx context.Context, in *CrmDepthAdsStatsRequest, opts ...grpc.CallOption) (*CrmDepthAdsStatsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthAdsStatsResult)
	if err := c.cc.Invoke(ctx, CrmService_GetAdsManagerStats_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateUtmLink(ctx context.Context, in *CrmDepthCreateUtmLinkRequest, opts ...grpc.CallOption) (*CrmDepthUtmLinkRow, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthUtmLinkRow)
	if err := c.cc.Invoke(ctx, CrmService_CreateUtmLink_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListUtmLinks(ctx context.Context, in *CrmDepthListUtmLinksRequest, opts ...grpc.CallOption) (*CrmDepthListUtmLinksResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListUtmLinksResult)
	if err := c.cc.Invoke(ctx, CrmService_ListUtmLinks_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateLandingPage(ctx context.Context, in *CrmDepthCreateLandingPageRequest, opts ...grpc.CallOption) (*CrmDepthLandingPageRow, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthLandingPageRow)
	if err := c.cc.Invoke(ctx, CrmService_CreateLandingPage_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListLandingPages(ctx context.Context, in *CrmDepthListLandingPagesRequest, opts ...grpc.CallOption) (*CrmDepthListLandingPagesResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListLandingPagesResult)
	if err := c.cc.Invoke(ctx, CrmService_ListLandingPages_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ScheduleContent(ctx context.Context, in *CrmDepthScheduleContentRequest, opts ...grpc.CallOption) (*CrmDepthScheduleContentResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthScheduleContentResult)
	if err := c.cc.Invoke(ctx, CrmService_ScheduleContent_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListScheduledContent(ctx context.Context, in *CrmDepthListScheduledContentRequest, opts ...grpc.CallOption) (*CrmDepthListScheduledContentResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListScheduledContentResult)
	if err := c.cc.Invoke(ctx, CrmService_ListScheduledContent_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetContentAnalytics(ctx context.Context, in *CrmDepthContentAnalyticsRequest, opts ...grpc.CallOption) (*CrmDepthContentAnalyticsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthContentAnalyticsResult)
	if err := c.cc.Invoke(ctx, CrmService_GetContentAnalytics_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateAgentLead(ctx context.Context, in *CrmDepthCreateAgentLeadRequest, opts ...grpc.CallOption) (*CrmDepthAgentLeadRow, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthAgentLeadRow)
	if err := c.cc.Invoke(ctx, CrmService_CreateAgentLead_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListAgentLeads(ctx context.Context, in *CrmDepthListAgentLeadsRequest, opts ...grpc.CallOption) (*CrmDepthListAgentLeadsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListAgentLeadsResult)
	if err := c.cc.Invoke(ctx, CrmService_ListAgentLeads_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) SetLeadReminder(ctx context.Context, in *CrmDepthSetLeadReminderRequest, opts ...grpc.CallOption) (*CrmDepthSetLeadReminderResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthSetLeadReminderResult)
	if err := c.cc.Invoke(ctx, CrmService_SetLeadReminder_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListLeadReminders(ctx context.Context, in *CrmDepthListLeadRemindersRequest, opts ...grpc.CallOption) (*CrmDepthListLeadRemindersResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListLeadRemindersResult)
	if err := c.cc.Invoke(ctx, CrmService_ListLeadReminders_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) FilterBotLeads(ctx context.Context, in *CrmDepthFilterBotLeadsRequest, opts ...grpc.CallOption) (*CrmDepthFilterBotLeadsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthFilterBotLeadsResult)
	if err := c.cc.Invoke(ctx, CrmService_FilterBotLeads_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateDripSequence(ctx context.Context, in *CrmDepthCreateDripSequenceRequest, opts ...grpc.CallOption) (*CrmDepthCreateDripSequenceResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthCreateDripSequenceResult)
	if err := c.cc.Invoke(ctx, CrmService_CreateDripSequence_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListDripSequences(ctx context.Context, in *CrmDepthListDripSequencesRequest, opts ...grpc.CallOption) (*CrmDepthListDripSequencesResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListDripSequencesResult)
	if err := c.cc.Invoke(ctx, CrmService_ListDripSequences_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateMomentTrigger(ctx context.Context, in *CrmDepthCreateMomentTriggerRequest, opts ...grpc.CallOption) (*CrmDepthMomentTriggerResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthMomentTriggerResult)
	if err := c.cc.Invoke(ctx, CrmService_CreateMomentTrigger_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateSegment(ctx context.Context, in *CrmDepthCreateSegmentRequest, opts ...grpc.CallOption) (*CrmDepthCreateSegmentResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthCreateSegmentResult)
	if err := c.cc.Invoke(ctx, CrmService_CreateSegment_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListSegments(ctx context.Context, in *CrmDepthListSegmentsRequest, opts ...grpc.CallOption) (*CrmDepthListSegmentsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListSegmentsResult)
	if err := c.cc.Invoke(ctx, CrmService_ListSegments_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) SendBroadcast(ctx context.Context, in *CrmDepthSendBroadcastRequest, opts ...grpc.CallOption) (*CrmDepthSendBroadcastResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthSendBroadcastResult)
	if err := c.cc.Invoke(ctx, CrmService_SendBroadcast_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) AssignLeadFair(ctx context.Context, in *CrmDepthAssignLeadFairRequest, opts ...grpc.CallOption) (*CrmDepthAssignLeadFairResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthAssignLeadFairResult)
	if err := c.cc.Invoke(ctx, CrmService_AssignLeadFair_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CreateSlaRule(ctx context.Context, in *CrmDepthCreateSlaRuleRequest, opts ...grpc.CallOption) (*CrmDepthSlaRuleResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthSlaRuleResult)
	if err := c.cc.Invoke(ctx, CrmService_CreateSlaRule_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetLeadTrail(ctx context.Context, in *CrmDepthGetLeadTrailRequest, opts ...grpc.CallOption) (*CrmDepthGetLeadTrailResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthGetLeadTrailResult)
	if err := c.cc.Invoke(ctx, CrmService_GetLeadTrail_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) TagLead(ctx context.Context, in *CrmDepthTagLeadRequest, opts ...grpc.CallOption) (*CrmDepthTagLeadResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthTagLeadResult)
	if err := c.cc.Invoke(ctx, CrmService_TagLead_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GenerateQuote(ctx context.Context, in *CrmDepthGenerateQuoteRequest, opts ...grpc.CallOption) (*CrmDepthQuoteResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthQuoteResult)
	if err := c.cc.Invoke(ctx, CrmService_GenerateQuote_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetQuote(ctx context.Context, in *CrmDepthGetQuoteRequest, opts ...grpc.CallOption) (*CrmDepthQuoteResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthQuoteResult)
	if err := c.cc.Invoke(ctx, CrmService_GetQuote_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) BuildPaymentLink(ctx context.Context, in *CrmDepthBuildPaymentLinkRequest, opts ...grpc.CallOption) (*CrmDepthPaymentLinkResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthPaymentLinkResult)
	if err := c.cc.Invoke(ctx, CrmService_BuildPaymentLink_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) RequestDiscount(ctx context.Context, in *CrmDepthRequestDiscountRequest, opts ...grpc.CallOption) (*CrmDepthDiscountResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthDiscountResult)
	if err := c.cc.Invoke(ctx, CrmService_RequestDiscount_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ApproveDiscount(ctx context.Context, in *CrmDepthApproveDiscountRequest, opts ...grpc.CallOption) (*CrmDepthDiscountResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthDiscountResult)
	if err := c.cc.Invoke(ctx, CrmService_ApproveDiscount_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetStaleProspects(ctx context.Context, in *CrmDepthStaleProspectsRequest, opts ...grpc.CallOption) (*CrmDepthStaleProspectsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthStaleProspectsResult)
	if err := c.cc.Invoke(ctx, CrmService_GetStaleProspects_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetCommissionBalance(ctx context.Context, in *CrmDepthGetCommissionRequest, opts ...grpc.CallOption) (*CrmDepthCommissionBalanceResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthCommissionBalanceResult)
	if err := c.cc.Invoke(ctx, CrmService_GetCommissionBalance_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetCommissionEvents(ctx context.Context, in *CrmDepthGetCommissionRequest, opts ...grpc.CallOption) (*CrmDepthGetCommissionEventsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthGetCommissionEventsResult)
	if err := c.cc.Invoke(ctx, CrmService_GetCommissionEvents_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) RequestPayout(ctx context.Context, in *CrmDepthRequestPayoutRequest, opts ...grpc.CallOption) (*CrmDepthPayoutResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthPayoutResult)
	if err := c.cc.Invoke(ctx, CrmService_RequestPayout_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetPayoutHistory(ctx context.Context, in *CrmDepthGetCommissionRequest, opts ...grpc.CallOption) (*CrmDepthGetPayoutHistoryResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthGetPayoutHistoryResult)
	if err := c.cc.Invoke(ctx, CrmService_GetPayoutHistory_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ComputeOverrideCommission(ctx context.Context, in *CrmDepthComputeOverrideRequest, opts ...grpc.CallOption) (*CrmDepthOverrideCommissionResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthOverrideCommissionResult)
	if err := c.cc.Invoke(ctx, CrmService_ComputeOverrideCommission_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetRoasReport(ctx context.Context, in *CrmDepthRoasReportRequest, opts ...grpc.CallOption) (*CrmDepthRoasReportResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthRoasReportResult)
	if err := c.cc.Invoke(ctx, CrmService_GetRoasReport_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListAcademyCourses(ctx context.Context, in *CrmDepthListCoursesRequest, opts ...grpc.CallOption) (*CrmDepthListCoursesResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListCoursesResult)
	if err := c.cc.Invoke(ctx, CrmService_ListAcademyCourses_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetCourseProgress(ctx context.Context, in *CrmDepthGetCourseProgressRequest, opts ...grpc.CallOption) (*CrmDepthCourseProgressResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthCourseProgressResult)
	if err := c.cc.Invoke(ctx, CrmService_GetCourseProgress_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) SubmitQuiz(ctx context.Context, in *CrmDepthSubmitQuizRequest, opts ...grpc.CallOption) (*CrmDepthQuizResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthQuizResult)
	if err := c.cc.Invoke(ctx, CrmService_SubmitQuiz_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) ListSalesScripts(ctx context.Context, in *CrmDepthListSalesScriptsRequest, opts ...grpc.CallOption) (*CrmDepthListSalesScriptsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthListSalesScriptsResult)
	if err := c.cc.Invoke(ctx, CrmService_ListSalesScripts_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetLeaderboard(ctx context.Context, in *CrmDepthLeaderboardRequest, opts ...grpc.CallOption) (*CrmDepthLeaderboardResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthLeaderboardResult)
	if err := c.cc.Invoke(ctx, CrmService_GetLeaderboard_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetAgentTier(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthAgentTierResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthAgentTierResult)
	if err := c.cc.Invoke(ctx, CrmService_GetAgentTier_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetAlumniReferrals(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthGetAlumniReferralsResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthGetAlumniReferralsResult)
	if err := c.cc.Invoke(ctx, CrmService_GetAlumniReferrals_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) GetReturnIntentSavings(ctx context.Context, in *CrmDepthGetAgentRequest, opts ...grpc.CallOption) (*CrmDepthReturnIntentResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthReturnIntentResult)
	if err := c.cc.Invoke(ctx, CrmService_GetReturnIntentSavings_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) CalculateZakat(ctx context.Context, in *CrmDepthCalculateZakatRequest, opts ...grpc.CallOption) (*CrmDepthZakatResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthZakatResult)
	if err := c.cc.Invoke(ctx, CrmService_CalculateZakat_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
func (c *crmDepthClient) RecordCharity(ctx context.Context, in *CrmDepthRecordCharityRequest, opts ...grpc.CallOption) (*CrmDepthCharityResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CrmDepthCharityResult)
	if err := c.cc.Invoke(ctx, CrmService_RecordCharity_FullMethodName, in, out, cOpts...); err != nil { return nil, err }
	return out, nil
}
