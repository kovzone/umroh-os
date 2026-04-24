// crm_depth_messages.go — hand-written proto message types for CRM depth RPCs
// (Wave 7 / BL-CRM-010..012, BL-CRM-017..066).
//
// All 53 request/response types with nil-safe getters.

package pb

// ---------------------------------------------------------------------------
// Agency Registration
// ---------------------------------------------------------------------------

type AgentRegisterRequest struct {
	AgencyName string `json:"agency_name"`
	OwnerName  string `json:"owner_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	City       string `json:"city"`
	Province   string `json:"province"`
}

func (x *AgentRegisterRequest) Reset()         {}
func (x *AgentRegisterRequest) String() string { return x.AgencyName }
func (x *AgentRegisterRequest) ProtoMessage()  {}
func (x *AgentRegisterRequest) GetAgencyName() string {
	if x == nil {
		return ""
	}
	return x.AgencyName
}
func (x *AgentRegisterRequest) GetOwnerName() string {
	if x == nil {
		return ""
	}
	return x.OwnerName
}
func (x *AgentRegisterRequest) GetPhone() string {
	if x == nil {
		return ""
	}
	return x.Phone
}
func (x *AgentRegisterRequest) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *AgentRegisterRequest) GetCity() string {
	if x == nil {
		return ""
	}
	return x.City
}
func (x *AgentRegisterRequest) GetProvince() string {
	if x == nil {
		return ""
	}
	return x.Province
}

type AgentResult struct {
	AgentId    string `json:"agent_id"`
	AgencyName string `json:"agency_name"`
	OwnerName  string `json:"owner_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Status     string `json:"status"`
	Tier       string `json:"tier"`
	CreatedAt  string `json:"created_at"`
}

func (x *AgentResult) Reset()         {}
func (x *AgentResult) String() string { return x.AgentId }
func (x *AgentResult) ProtoMessage()  {}
func (x *AgentResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AgentResult) GetAgencyName() string {
	if x == nil {
		return ""
	}
	return x.AgencyName
}
func (x *AgentResult) GetOwnerName() string {
	if x == nil {
		return ""
	}
	return x.OwnerName
}
func (x *AgentResult) GetPhone() string {
	if x == nil {
		return ""
	}
	return x.Phone
}
func (x *AgentResult) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *AgentResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *AgentResult) GetTier() string {
	if x == nil {
		return ""
	}
	return x.Tier
}
func (x *AgentResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type AgentKycRequest struct {
	AgentId  string `json:"agent_id"`
	KtpUrl   string `json:"ktp_url"`
	NpwpUrl  string `json:"npwp_url"`
	SiupUrl  string `json:"siup_url"`
	SelfieUrl string `json:"selfie_url"`
}

func (x *AgentKycRequest) Reset()         {}
func (x *AgentKycRequest) String() string { return x.AgentId }
func (x *AgentKycRequest) ProtoMessage()  {}
func (x *AgentKycRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AgentKycRequest) GetKtpUrl() string {
	if x == nil {
		return ""
	}
	return x.KtpUrl
}
func (x *AgentKycRequest) GetNpwpUrl() string {
	if x == nil {
		return ""
	}
	return x.NpwpUrl
}
func (x *AgentKycRequest) GetSiupUrl() string {
	if x == nil {
		return ""
	}
	return x.SiupUrl
}
func (x *AgentKycRequest) GetSelfieUrl() string {
	if x == nil {
		return ""
	}
	return x.SelfieUrl
}

type AgentMoURequest struct {
	AgentId     string `json:"agent_id"`
	SignatureUrl string `json:"signature_url"`
	SignedAt    string `json:"signed_at"`
}

func (x *AgentMoURequest) Reset()         {}
func (x *AgentMoURequest) String() string { return x.AgentId }
func (x *AgentMoURequest) ProtoMessage()  {}
func (x *AgentMoURequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AgentMoURequest) GetSignatureUrl() string {
	if x == nil {
		return ""
	}
	return x.SignatureUrl
}
func (x *AgentMoURequest) GetSignedAt() string {
	if x == nil {
		return ""
	}
	return x.SignedAt
}

type AgentIdRequest struct {
	AgentId string `json:"agent_id"`
}

func (x *AgentIdRequest) Reset()         {}
func (x *AgentIdRequest) String() string { return x.AgentId }
func (x *AgentIdRequest) ProtoMessage()  {}
func (x *AgentIdRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}

type ReplicaSiteResult struct {
	AgentId  string `json:"agent_id"`
	SiteUrl  string `json:"site_url"`
	Theme    string `json:"theme"`
	LogoUrl  string `json:"logo_url"`
	UpdatedAt string `json:"updated_at"`
}

func (x *ReplicaSiteResult) Reset()         {}
func (x *ReplicaSiteResult) String() string { return x.SiteUrl }
func (x *ReplicaSiteResult) ProtoMessage()  {}
func (x *ReplicaSiteResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ReplicaSiteResult) GetSiteUrl() string {
	if x == nil {
		return ""
	}
	return x.SiteUrl
}
func (x *ReplicaSiteResult) GetTheme() string {
	if x == nil {
		return ""
	}
	return x.Theme
}
func (x *ReplicaSiteResult) GetLogoUrl() string {
	if x == nil {
		return ""
	}
	return x.LogoUrl
}
func (x *ReplicaSiteResult) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

type UpdateReplicaSiteRequest struct {
	AgentId string `json:"agent_id"`
	Theme   string `json:"theme"`
	LogoUrl string `json:"logo_url"`
	CustomDomain string `json:"custom_domain"`
}

func (x *UpdateReplicaSiteRequest) Reset()         {}
func (x *UpdateReplicaSiteRequest) String() string { return x.AgentId }
func (x *UpdateReplicaSiteRequest) ProtoMessage()  {}
func (x *UpdateReplicaSiteRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *UpdateReplicaSiteRequest) GetTheme() string {
	if x == nil {
		return ""
	}
	return x.Theme
}
func (x *UpdateReplicaSiteRequest) GetLogoUrl() string {
	if x == nil {
		return ""
	}
	return x.LogoUrl
}
func (x *UpdateReplicaSiteRequest) GetCustomDomain() string {
	if x == nil {
		return ""
	}
	return x.CustomDomain
}

// ---------------------------------------------------------------------------
// Content & Marketing
// ---------------------------------------------------------------------------

type SocialShareRequest struct {
	AgentId   string `json:"agent_id"`
	PackageId string `json:"package_id"`
	Platform  string `json:"platform"`
}

func (x *SocialShareRequest) Reset()         {}
func (x *SocialShareRequest) String() string { return "" }
func (x *SocialShareRequest) ProtoMessage()  {}
func (x *SocialShareRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *SocialShareRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *SocialShareRequest) GetPlatform() string {
	if x == nil {
		return ""
	}
	return x.Platform
}

type SocialShareResult struct {
	ShareUrl  string `json:"share_url"`
	ShortCode string `json:"short_code"`
}

func (x *SocialShareResult) Reset()         {}
func (x *SocialShareResult) String() string { return x.ShareUrl }
func (x *SocialShareResult) ProtoMessage()  {}
func (x *SocialShareResult) GetShareUrl() string {
	if x == nil {
		return ""
	}
	return x.ShareUrl
}
func (x *SocialShareResult) GetShortCode() string {
	if x == nil {
		return ""
	}
	return x.ShortCode
}

type BusinessCardRequest struct {
	AgentId  string `json:"agent_id"`
	Template string `json:"template"`
}

func (x *BusinessCardRequest) Reset()         {}
func (x *BusinessCardRequest) String() string { return x.AgentId }
func (x *BusinessCardRequest) ProtoMessage()  {}
func (x *BusinessCardRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *BusinessCardRequest) GetTemplate() string {
	if x == nil {
		return ""
	}
	return x.Template
}

type BusinessCardResult struct {
	CardUrl   string `json:"card_url"`
	CreatedAt string `json:"created_at"`
}

func (x *BusinessCardResult) Reset()         {}
func (x *BusinessCardResult) String() string { return x.CardUrl }
func (x *BusinessCardResult) ProtoMessage()  {}
func (x *BusinessCardResult) GetCardUrl() string {
	if x == nil {
		return ""
	}
	return x.CardUrl
}
func (x *BusinessCardResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ContentBankListRequest struct {
	AgentId  string `json:"agent_id"`
	Category string `json:"category"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ContentBankListRequest) Reset()         {}
func (x *ContentBankListRequest) String() string { return "" }
func (x *ContentBankListRequest) ProtoMessage()  {}
func (x *ContentBankListRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ContentBankListRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *ContentBankListRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ContentBankListRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type ContentAssetItem struct {
	AssetId   string `json:"asset_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ContentBankListResult struct {
	Assets   []*ContentAssetItem `json:"assets"`
	Total    int64               `json:"total"`
	Page     int32               `json:"page"`
	PageSize int32               `json:"page_size"`
}

func (x *ContentBankListResult) Reset()         {}
func (x *ContentBankListResult) String() string { return "" }
func (x *ContentBankListResult) ProtoMessage()  {}
func (x *ContentBankListResult) GetAssets() []*ContentAssetItem {
	if x == nil {
		return nil
	}
	return x.Assets
}
func (x *ContentBankListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *ContentBankListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ContentBankListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type CreateContentAssetRequest struct {
	AgentId  string `json:"agent_id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Url      string `json:"url"`
	Category string `json:"category"`
}

func (x *CreateContentAssetRequest) Reset()         {}
func (x *CreateContentAssetRequest) String() string { return x.Title }
func (x *CreateContentAssetRequest) ProtoMessage()  {}
func (x *CreateContentAssetRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateContentAssetRequest) GetTitle() string {
	if x == nil {
		return ""
	}
	return x.Title
}
func (x *CreateContentAssetRequest) GetType() string {
	if x == nil {
		return ""
	}
	return x.Type
}
func (x *CreateContentAssetRequest) GetUrl() string {
	if x == nil {
		return ""
	}
	return x.Url
}
func (x *CreateContentAssetRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}

type ContentAssetResult struct {
	AssetId   string `json:"asset_id"`
	AgentId   string `json:"agent_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

func (x *ContentAssetResult) Reset()         {}
func (x *ContentAssetResult) String() string { return x.AssetId }
func (x *ContentAssetResult) ProtoMessage()  {}
func (x *ContentAssetResult) GetAssetId() string {
	if x == nil {
		return ""
	}
	return x.AssetId
}
func (x *ContentAssetResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ContentAssetResult) GetTitle() string {
	if x == nil {
		return ""
	}
	return x.Title
}
func (x *ContentAssetResult) GetType() string {
	if x == nil {
		return ""
	}
	return x.Type
}
func (x *ContentAssetResult) GetUrl() string {
	if x == nil {
		return ""
	}
	return x.Url
}
func (x *ContentAssetResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type WatermarkFlyerRequest struct {
	AgentId  string `json:"agent_id"`
	FlyerUrl string `json:"flyer_url"`
	Text     string `json:"text"`
}

func (x *WatermarkFlyerRequest) Reset()         {}
func (x *WatermarkFlyerRequest) String() string { return x.AgentId }
func (x *WatermarkFlyerRequest) ProtoMessage()  {}
func (x *WatermarkFlyerRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *WatermarkFlyerRequest) GetFlyerUrl() string {
	if x == nil {
		return ""
	}
	return x.FlyerUrl
}
func (x *WatermarkFlyerRequest) GetText() string {
	if x == nil {
		return ""
	}
	return x.Text
}

type WatermarkFlyerResult struct {
	ResultUrl string `json:"result_url"`
}

func (x *WatermarkFlyerResult) Reset()         {}
func (x *WatermarkFlyerResult) String() string { return x.ResultUrl }
func (x *WatermarkFlyerResult) ProtoMessage()  {}
func (x *WatermarkFlyerResult) GetResultUrl() string {
	if x == nil {
		return ""
	}
	return x.ResultUrl
}

type ProgramGalleryRequest struct {
	PackageId string `json:"package_id"`
	Page      int32  `json:"page"`
	PageSize  int32  `json:"page_size"`
}

func (x *ProgramGalleryRequest) Reset()         {}
func (x *ProgramGalleryRequest) String() string { return x.PackageId }
func (x *ProgramGalleryRequest) ProtoMessage()  {}
func (x *ProgramGalleryRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *ProgramGalleryRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ProgramGalleryRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type GalleryItem struct {
	ImageUrl string `json:"image_url"`
	Caption  string `json:"caption"`
}

type ProgramGalleryResult struct {
	PackageId string         `json:"package_id"`
	Items     []*GalleryItem `json:"items"`
	Total     int64          `json:"total"`
}

func (x *ProgramGalleryResult) Reset()         {}
func (x *ProgramGalleryResult) String() string { return x.PackageId }
func (x *ProgramGalleryResult) ProtoMessage()  {}
func (x *ProgramGalleryResult) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *ProgramGalleryResult) GetItems() []*GalleryItem {
	if x == nil {
		return nil
	}
	return x.Items
}
func (x *ProgramGalleryResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}

type TrackingCodeRequest struct {
	AgentId string `json:"agent_id"`
	Code    string `json:"code"`
	Type    string `json:"type"`
}

func (x *TrackingCodeRequest) Reset()         {}
func (x *TrackingCodeRequest) String() string { return x.AgentId }
func (x *TrackingCodeRequest) ProtoMessage()  {}
func (x *TrackingCodeRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *TrackingCodeRequest) GetCode() string {
	if x == nil {
		return ""
	}
	return x.Code
}
func (x *TrackingCodeRequest) GetType() string {
	if x == nil {
		return ""
	}
	return x.Type
}

type TrackingCodeResult struct {
	AgentId   string `json:"agent_id"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	UpdatedAt string `json:"updated_at"`
}

func (x *TrackingCodeResult) Reset()         {}
func (x *TrackingCodeResult) String() string { return x.AgentId }
func (x *TrackingCodeResult) ProtoMessage()  {}
func (x *TrackingCodeResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *TrackingCodeResult) GetCode() string {
	if x == nil {
		return ""
	}
	return x.Code
}
func (x *TrackingCodeResult) GetType() string {
	if x == nil {
		return ""
	}
	return x.Type
}
func (x *TrackingCodeResult) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

type AdsManagerRequest struct {
	AgentId   string `json:"agent_id"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
}

func (x *AdsManagerRequest) Reset()         {}
func (x *AdsManagerRequest) String() string { return x.AgentId }
func (x *AdsManagerRequest) ProtoMessage()  {}
func (x *AdsManagerRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AdsManagerRequest) GetDateFrom() string {
	if x == nil {
		return ""
	}
	return x.DateFrom
}
func (x *AdsManagerRequest) GetDateTo() string {
	if x == nil {
		return ""
	}
	return x.DateTo
}

type AdsManagerResult struct {
	AgentId     string  `json:"agent_id"`
	Impressions int64   `json:"impressions"`
	Clicks      int64   `json:"clicks"`
	Conversions int64   `json:"conversions"`
	Spend       float64 `json:"spend"`
}

func (x *AdsManagerResult) Reset()         {}
func (x *AdsManagerResult) String() string { return x.AgentId }
func (x *AdsManagerResult) ProtoMessage()  {}
func (x *AdsManagerResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AdsManagerResult) GetImpressions() int64 {
	if x == nil {
		return 0
	}
	return x.Impressions
}
func (x *AdsManagerResult) GetClicks() int64 {
	if x == nil {
		return 0
	}
	return x.Clicks
}
func (x *AdsManagerResult) GetConversions() int64 {
	if x == nil {
		return 0
	}
	return x.Conversions
}
func (x *AdsManagerResult) GetSpend() float64 {
	if x == nil {
		return 0
	}
	return x.Spend
}

type CreateUtmLinkRequest struct {
	AgentId    string `json:"agent_id"`
	BaseUrl    string `json:"base_url"`
	UtmSource  string `json:"utm_source"`
	UtmMedium  string `json:"utm_medium"`
	UtmCampaign string `json:"utm_campaign"`
}

func (x *CreateUtmLinkRequest) Reset()         {}
func (x *CreateUtmLinkRequest) String() string { return x.AgentId }
func (x *CreateUtmLinkRequest) ProtoMessage()  {}
func (x *CreateUtmLinkRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateUtmLinkRequest) GetBaseUrl() string {
	if x == nil {
		return ""
	}
	return x.BaseUrl
}
func (x *CreateUtmLinkRequest) GetUtmSource() string {
	if x == nil {
		return ""
	}
	return x.UtmSource
}
func (x *CreateUtmLinkRequest) GetUtmMedium() string {
	if x == nil {
		return ""
	}
	return x.UtmMedium
}
func (x *CreateUtmLinkRequest) GetUtmCampaign() string {
	if x == nil {
		return ""
	}
	return x.UtmCampaign
}

type UtmLinkResult struct {
	LinkId    string `json:"link_id"`
	AgentId   string `json:"agent_id"`
	FullUrl   string `json:"full_url"`
	ShortUrl  string `json:"short_url"`
	CreatedAt string `json:"created_at"`
}

func (x *UtmLinkResult) Reset()         {}
func (x *UtmLinkResult) String() string { return x.LinkId }
func (x *UtmLinkResult) ProtoMessage()  {}
func (x *UtmLinkResult) GetLinkId() string {
	if x == nil {
		return ""
	}
	return x.LinkId
}
func (x *UtmLinkResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *UtmLinkResult) GetFullUrl() string {
	if x == nil {
		return ""
	}
	return x.FullUrl
}
func (x *UtmLinkResult) GetShortUrl() string {
	if x == nil {
		return ""
	}
	return x.ShortUrl
}
func (x *UtmLinkResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ListUtmLinksRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ListUtmLinksRequest) Reset()         {}
func (x *ListUtmLinksRequest) String() string { return x.AgentId }
func (x *ListUtmLinksRequest) ProtoMessage()  {}
func (x *ListUtmLinksRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ListUtmLinksRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListUtmLinksRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type UtmLinkListResult struct {
	Links    []*UtmLinkResult `json:"links"`
	Total    int64            `json:"total"`
	Page     int32            `json:"page"`
	PageSize int32            `json:"page_size"`
}

func (x *UtmLinkListResult) Reset()         {}
func (x *UtmLinkListResult) String() string { return "" }
func (x *UtmLinkListResult) ProtoMessage()  {}
func (x *UtmLinkListResult) GetLinks() []*UtmLinkResult {
	if x == nil {
		return nil
	}
	return x.Links
}
func (x *UtmLinkListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *UtmLinkListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *UtmLinkListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type CreateLandingPageRequest struct {
	AgentId   string `json:"agent_id"`
	Title     string `json:"title"`
	PackageId string `json:"package_id"`
	Template  string `json:"template"`
}

func (x *CreateLandingPageRequest) Reset()         {}
func (x *CreateLandingPageRequest) String() string { return x.Title }
func (x *CreateLandingPageRequest) ProtoMessage()  {}
func (x *CreateLandingPageRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateLandingPageRequest) GetTitle() string {
	if x == nil {
		return ""
	}
	return x.Title
}
func (x *CreateLandingPageRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *CreateLandingPageRequest) GetTemplate() string {
	if x == nil {
		return ""
	}
	return x.Template
}

type LandingPageResult struct {
	PageId    string `json:"page_id"`
	AgentId   string `json:"agent_id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (x *LandingPageResult) Reset()         {}
func (x *LandingPageResult) String() string { return x.PageId }
func (x *LandingPageResult) ProtoMessage()  {}
func (x *LandingPageResult) GetPageId() string {
	if x == nil {
		return ""
	}
	return x.PageId
}
func (x *LandingPageResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *LandingPageResult) GetTitle() string {
	if x == nil {
		return ""
	}
	return x.Title
}
func (x *LandingPageResult) GetUrl() string {
	if x == nil {
		return ""
	}
	return x.Url
}
func (x *LandingPageResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *LandingPageResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ListLandingPagesRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ListLandingPagesRequest) Reset()         {}
func (x *ListLandingPagesRequest) String() string { return x.AgentId }
func (x *ListLandingPagesRequest) ProtoMessage()  {}
func (x *ListLandingPagesRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ListLandingPagesRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListLandingPagesRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type LandingPageListResult struct {
	Pages    []*LandingPageResult `json:"pages"`
	Total    int64                `json:"total"`
	Page     int32                `json:"page"`
	PageSize int32                `json:"page_size"`
}

func (x *LandingPageListResult) Reset()         {}
func (x *LandingPageListResult) String() string { return "" }
func (x *LandingPageListResult) ProtoMessage()  {}
func (x *LandingPageListResult) GetPages() []*LandingPageResult {
	if x == nil {
		return nil
	}
	return x.Pages
}
func (x *LandingPageListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *LandingPageListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *LandingPageListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type ScheduleContentRequest struct {
	AgentId   string `json:"agent_id"`
	AssetId   string `json:"asset_id"`
	Platform  string `json:"platform"`
	ScheduledAt string `json:"scheduled_at"`
	Caption   string `json:"caption"`
}

func (x *ScheduleContentRequest) Reset()         {}
func (x *ScheduleContentRequest) String() string { return x.AgentId }
func (x *ScheduleContentRequest) ProtoMessage()  {}
func (x *ScheduleContentRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ScheduleContentRequest) GetAssetId() string {
	if x == nil {
		return ""
	}
	return x.AssetId
}
func (x *ScheduleContentRequest) GetPlatform() string {
	if x == nil {
		return ""
	}
	return x.Platform
}
func (x *ScheduleContentRequest) GetScheduledAt() string {
	if x == nil {
		return ""
	}
	return x.ScheduledAt
}
func (x *ScheduleContentRequest) GetCaption() string {
	if x == nil {
		return ""
	}
	return x.Caption
}

type ScheduledContentResult struct {
	ScheduleId  string `json:"schedule_id"`
	AgentId     string `json:"agent_id"`
	AssetId     string `json:"asset_id"`
	Platform    string `json:"platform"`
	ScheduledAt string `json:"scheduled_at"`
	Status      string `json:"status"`
}

func (x *ScheduledContentResult) Reset()         {}
func (x *ScheduledContentResult) String() string { return x.ScheduleId }
func (x *ScheduledContentResult) ProtoMessage()  {}
func (x *ScheduledContentResult) GetScheduleId() string {
	if x == nil {
		return ""
	}
	return x.ScheduleId
}
func (x *ScheduledContentResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ScheduledContentResult) GetAssetId() string {
	if x == nil {
		return ""
	}
	return x.AssetId
}
func (x *ScheduledContentResult) GetPlatform() string {
	if x == nil {
		return ""
	}
	return x.Platform
}
func (x *ScheduledContentResult) GetScheduledAt() string {
	if x == nil {
		return ""
	}
	return x.ScheduledAt
}
func (x *ScheduledContentResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type ListScheduledContentRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ListScheduledContentRequest) Reset()         {}
func (x *ListScheduledContentRequest) String() string { return x.AgentId }
func (x *ListScheduledContentRequest) ProtoMessage()  {}
func (x *ListScheduledContentRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ListScheduledContentRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListScheduledContentRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type ScheduledContentListResult struct {
	Items    []*ScheduledContentResult `json:"items"`
	Total    int64                     `json:"total"`
	Page     int32                     `json:"page"`
	PageSize int32                     `json:"page_size"`
}

func (x *ScheduledContentListResult) Reset()         {}
func (x *ScheduledContentListResult) String() string { return "" }
func (x *ScheduledContentListResult) ProtoMessage()  {}
func (x *ScheduledContentListResult) GetItems() []*ScheduledContentResult {
	if x == nil {
		return nil
	}
	return x.Items
}
func (x *ScheduledContentListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *ScheduledContentListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ScheduledContentListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type ContentAnalyticsRequest struct {
	AgentId  string `json:"agent_id"`
	AssetId  string `json:"asset_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (x *ContentAnalyticsRequest) Reset()         {}
func (x *ContentAnalyticsRequest) String() string { return x.AgentId }
func (x *ContentAnalyticsRequest) ProtoMessage()  {}
func (x *ContentAnalyticsRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ContentAnalyticsRequest) GetAssetId() string {
	if x == nil {
		return ""
	}
	return x.AssetId
}
func (x *ContentAnalyticsRequest) GetDateFrom() string {
	if x == nil {
		return ""
	}
	return x.DateFrom
}
func (x *ContentAnalyticsRequest) GetDateTo() string {
	if x == nil {
		return ""
	}
	return x.DateTo
}

type ContentAnalyticsResult struct {
	AgentId  string  `json:"agent_id"`
	Views    int64   `json:"views"`
	Likes    int64   `json:"likes"`
	Shares   int64   `json:"shares"`
	Leads    int64   `json:"leads"`
	Ctr      float64 `json:"ctr"`
}

func (x *ContentAnalyticsResult) Reset()         {}
func (x *ContentAnalyticsResult) String() string { return x.AgentId }
func (x *ContentAnalyticsResult) ProtoMessage()  {}
func (x *ContentAnalyticsResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ContentAnalyticsResult) GetViews() int64 {
	if x == nil {
		return 0
	}
	return x.Views
}
func (x *ContentAnalyticsResult) GetLikes() int64 {
	if x == nil {
		return 0
	}
	return x.Likes
}
func (x *ContentAnalyticsResult) GetShares() int64 {
	if x == nil {
		return 0
	}
	return x.Shares
}
func (x *ContentAnalyticsResult) GetLeads() int64 {
	if x == nil {
		return 0
	}
	return x.Leads
}
func (x *ContentAnalyticsResult) GetCtr() float64 {
	if x == nil {
		return 0
	}
	return x.Ctr
}

// ---------------------------------------------------------------------------
// CRM / Leads
// ---------------------------------------------------------------------------

type CreateAgentLeadRequest struct {
	AgentId   string `json:"agent_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	PackageId string `json:"package_id"`
	Source    string `json:"source"`
	Notes     string `json:"notes"`
}

func (x *CreateAgentLeadRequest) Reset()         {}
func (x *CreateAgentLeadRequest) String() string { return x.AgentId }
func (x *CreateAgentLeadRequest) ProtoMessage()  {}
func (x *CreateAgentLeadRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateAgentLeadRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateAgentLeadRequest) GetPhone() string {
	if x == nil {
		return ""
	}
	return x.Phone
}
func (x *CreateAgentLeadRequest) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *CreateAgentLeadRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *CreateAgentLeadRequest) GetSource() string {
	if x == nil {
		return ""
	}
	return x.Source
}
func (x *CreateAgentLeadRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type AgentLeadResult struct {
	LeadId    string `json:"lead_id"`
	AgentId   string `json:"agent_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	PackageId string `json:"package_id"`
	CreatedAt string `json:"created_at"`
}

func (x *AgentLeadResult) Reset()         {}
func (x *AgentLeadResult) String() string { return x.LeadId }
func (x *AgentLeadResult) ProtoMessage()  {}
func (x *AgentLeadResult) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *AgentLeadResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AgentLeadResult) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *AgentLeadResult) GetPhone() string {
	if x == nil {
		return ""
	}
	return x.Phone
}
func (x *AgentLeadResult) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *AgentLeadResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *AgentLeadResult) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *AgentLeadResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ListAgentLeadsRequest struct {
	AgentId  string `json:"agent_id"`
	Status   string `json:"status"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ListAgentLeadsRequest) Reset()         {}
func (x *ListAgentLeadsRequest) String() string { return x.AgentId }
func (x *ListAgentLeadsRequest) ProtoMessage()  {}
func (x *ListAgentLeadsRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ListAgentLeadsRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *ListAgentLeadsRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListAgentLeadsRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type AgentLeadListResult struct {
	Leads    []*AgentLeadResult `json:"leads"`
	Total    int64              `json:"total"`
	Page     int32              `json:"page"`
	PageSize int32              `json:"page_size"`
}

func (x *AgentLeadListResult) Reset()         {}
func (x *AgentLeadListResult) String() string { return "" }
func (x *AgentLeadListResult) ProtoMessage()  {}
func (x *AgentLeadListResult) GetLeads() []*AgentLeadResult {
	if x == nil {
		return nil
	}
	return x.Leads
}
func (x *AgentLeadListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *AgentLeadListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *AgentLeadListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type SetLeadReminderRequest struct {
	LeadId     string `json:"lead_id"`
	AgentId    string `json:"agent_id"`
	RemindAt   string `json:"remind_at"`
	Note       string `json:"note"`
}

func (x *SetLeadReminderRequest) Reset()         {}
func (x *SetLeadReminderRequest) String() string { return x.LeadId }
func (x *SetLeadReminderRequest) ProtoMessage()  {}
func (x *SetLeadReminderRequest) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *SetLeadReminderRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *SetLeadReminderRequest) GetRemindAt() string {
	if x == nil {
		return ""
	}
	return x.RemindAt
}
func (x *SetLeadReminderRequest) GetNote() string {
	if x == nil {
		return ""
	}
	return x.Note
}

type LeadReminderResult struct {
	ReminderId string `json:"reminder_id"`
	LeadId     string `json:"lead_id"`
	AgentId    string `json:"agent_id"`
	RemindAt   string `json:"remind_at"`
	Note       string `json:"note"`
	Status     string `json:"status"`
}

func (x *LeadReminderResult) Reset()         {}
func (x *LeadReminderResult) String() string { return x.ReminderId }
func (x *LeadReminderResult) ProtoMessage()  {}
func (x *LeadReminderResult) GetReminderId() string {
	if x == nil {
		return ""
	}
	return x.ReminderId
}
func (x *LeadReminderResult) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *LeadReminderResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *LeadReminderResult) GetRemindAt() string {
	if x == nil {
		return ""
	}
	return x.RemindAt
}
func (x *LeadReminderResult) GetNote() string {
	if x == nil {
		return ""
	}
	return x.Note
}
func (x *LeadReminderResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type ListLeadRemindersRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ListLeadRemindersRequest) Reset()         {}
func (x *ListLeadRemindersRequest) String() string { return x.AgentId }
func (x *ListLeadRemindersRequest) ProtoMessage()  {}
func (x *ListLeadRemindersRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ListLeadRemindersRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListLeadRemindersRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type LeadReminderListResult struct {
	Reminders []*LeadReminderResult `json:"reminders"`
	Total     int64                 `json:"total"`
	Page      int32                 `json:"page"`
	PageSize  int32                 `json:"page_size"`
}

func (x *LeadReminderListResult) Reset()         {}
func (x *LeadReminderListResult) String() string { return "" }
func (x *LeadReminderListResult) ProtoMessage()  {}
func (x *LeadReminderListResult) GetReminders() []*LeadReminderResult {
	if x == nil {
		return nil
	}
	return x.Reminders
}
func (x *LeadReminderListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *LeadReminderListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *LeadReminderListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type BotFilterRequest struct {
	AgentId string `json:"agent_id"`
	LeadId  string `json:"lead_id"`
}

func (x *BotFilterRequest) Reset()         {}
func (x *BotFilterRequest) String() string { return x.LeadId }
func (x *BotFilterRequest) ProtoMessage()  {}
func (x *BotFilterRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *BotFilterRequest) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}

type BotFilterResult struct {
	LeadId  string  `json:"lead_id"`
	IsBot   bool    `json:"is_bot"`
	Score   float64 `json:"score"`
	Reason  string  `json:"reason"`
}

func (x *BotFilterResult) Reset()         {}
func (x *BotFilterResult) String() string { return x.LeadId }
func (x *BotFilterResult) ProtoMessage()  {}
func (x *BotFilterResult) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *BotFilterResult) GetIsBot() bool {
	if x == nil {
		return false
	}
	return x.IsBot
}
func (x *BotFilterResult) GetScore() float64 {
	if x == nil {
		return 0
	}
	return x.Score
}
func (x *BotFilterResult) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}

type CreateDripSequenceRequest struct {
	AgentId  string   `json:"agent_id"`
	Name     string   `json:"name"`
	Triggers []string `json:"triggers"`
	Steps    int32    `json:"steps"`
}

func (x *CreateDripSequenceRequest) Reset()         {}
func (x *CreateDripSequenceRequest) String() string { return x.Name }
func (x *CreateDripSequenceRequest) ProtoMessage()  {}
func (x *CreateDripSequenceRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateDripSequenceRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateDripSequenceRequest) GetTriggers() []string {
	if x == nil {
		return nil
	}
	return x.Triggers
}
func (x *CreateDripSequenceRequest) GetSteps() int32 {
	if x == nil {
		return 0
	}
	return x.Steps
}

type DripSequenceResult struct {
	SequenceId string `json:"sequence_id"`
	AgentId    string `json:"agent_id"`
	Name       string `json:"name"`
	Steps      int32  `json:"steps"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

func (x *DripSequenceResult) Reset()         {}
func (x *DripSequenceResult) String() string { return x.SequenceId }
func (x *DripSequenceResult) ProtoMessage()  {}
func (x *DripSequenceResult) GetSequenceId() string {
	if x == nil {
		return ""
	}
	return x.SequenceId
}
func (x *DripSequenceResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *DripSequenceResult) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *DripSequenceResult) GetSteps() int32 {
	if x == nil {
		return 0
	}
	return x.Steps
}
func (x *DripSequenceResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *DripSequenceResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ListDripSequencesRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ListDripSequencesRequest) Reset()         {}
func (x *ListDripSequencesRequest) String() string { return x.AgentId }
func (x *ListDripSequencesRequest) ProtoMessage()  {}
func (x *ListDripSequencesRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ListDripSequencesRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListDripSequencesRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type DripSequenceListResult struct {
	Sequences []*DripSequenceResult `json:"sequences"`
	Total     int64                 `json:"total"`
	Page      int32                 `json:"page"`
	PageSize  int32                 `json:"page_size"`
}

func (x *DripSequenceListResult) Reset()         {}
func (x *DripSequenceListResult) String() string { return "" }
func (x *DripSequenceListResult) ProtoMessage()  {}
func (x *DripSequenceListResult) GetSequences() []*DripSequenceResult {
	if x == nil {
		return nil
	}
	return x.Sequences
}
func (x *DripSequenceListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *DripSequenceListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *DripSequenceListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type CreateMomentTriggerRequest struct {
	AgentId   string `json:"agent_id"`
	EventType string `json:"event_type"`
	Action    string `json:"action"`
	Template  string `json:"template"`
}

func (x *CreateMomentTriggerRequest) Reset()         {}
func (x *CreateMomentTriggerRequest) String() string { return x.AgentId }
func (x *CreateMomentTriggerRequest) ProtoMessage()  {}
func (x *CreateMomentTriggerRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateMomentTriggerRequest) GetEventType() string {
	if x == nil {
		return ""
	}
	return x.EventType
}
func (x *CreateMomentTriggerRequest) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *CreateMomentTriggerRequest) GetTemplate() string {
	if x == nil {
		return ""
	}
	return x.Template
}

type MomentTriggerResult struct {
	TriggerId string `json:"trigger_id"`
	AgentId   string `json:"agent_id"`
	EventType string `json:"event_type"`
	Action    string `json:"action"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (x *MomentTriggerResult) Reset()         {}
func (x *MomentTriggerResult) String() string { return x.TriggerId }
func (x *MomentTriggerResult) ProtoMessage()  {}
func (x *MomentTriggerResult) GetTriggerId() string {
	if x == nil {
		return ""
	}
	return x.TriggerId
}
func (x *MomentTriggerResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *MomentTriggerResult) GetEventType() string {
	if x == nil {
		return ""
	}
	return x.EventType
}
func (x *MomentTriggerResult) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *MomentTriggerResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *MomentTriggerResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type CreateSegmentRequest struct {
	AgentId  string   `json:"agent_id"`
	Name     string   `json:"name"`
	Criteria []string `json:"criteria"`
}

func (x *CreateSegmentRequest) Reset()         {}
func (x *CreateSegmentRequest) String() string { return x.Name }
func (x *CreateSegmentRequest) ProtoMessage()  {}
func (x *CreateSegmentRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateSegmentRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateSegmentRequest) GetCriteria() []string {
	if x == nil {
		return nil
	}
	return x.Criteria
}

type SegmentResult struct {
	SegmentId  string `json:"segment_id"`
	AgentId    string `json:"agent_id"`
	Name       string `json:"name"`
	LeadCount  int64  `json:"lead_count"`
	CreatedAt  string `json:"created_at"`
}

func (x *SegmentResult) Reset()         {}
func (x *SegmentResult) String() string { return x.SegmentId }
func (x *SegmentResult) ProtoMessage()  {}
func (x *SegmentResult) GetSegmentId() string {
	if x == nil {
		return ""
	}
	return x.SegmentId
}
func (x *SegmentResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *SegmentResult) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *SegmentResult) GetLeadCount() int64 {
	if x == nil {
		return 0
	}
	return x.LeadCount
}
func (x *SegmentResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ListSegmentsRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *ListSegmentsRequest) Reset()         {}
func (x *ListSegmentsRequest) String() string { return x.AgentId }
func (x *ListSegmentsRequest) ProtoMessage()  {}
func (x *ListSegmentsRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ListSegmentsRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListSegmentsRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type SegmentListResult struct {
	Segments []*SegmentResult `json:"segments"`
	Total    int64            `json:"total"`
	Page     int32            `json:"page"`
	PageSize int32            `json:"page_size"`
}

func (x *SegmentListResult) Reset()         {}
func (x *SegmentListResult) String() string { return "" }
func (x *SegmentListResult) ProtoMessage()  {}
func (x *SegmentListResult) GetSegments() []*SegmentResult {
	if x == nil {
		return nil
	}
	return x.Segments
}
func (x *SegmentListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *SegmentListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *SegmentListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type SendBroadcastRequest struct {
	AgentId   string `json:"agent_id"`
	SegmentId string `json:"segment_id"`
	Channel   string `json:"channel"`
	Message   string `json:"message"`
}

func (x *SendBroadcastRequest) Reset()         {}
func (x *SendBroadcastRequest) String() string { return x.AgentId }
func (x *SendBroadcastRequest) ProtoMessage()  {}
func (x *SendBroadcastRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *SendBroadcastRequest) GetSegmentId() string {
	if x == nil {
		return ""
	}
	return x.SegmentId
}
func (x *SendBroadcastRequest) GetChannel() string {
	if x == nil {
		return ""
	}
	return x.Channel
}
func (x *SendBroadcastRequest) GetMessage() string {
	if x == nil {
		return ""
	}
	return x.Message
}

type BroadcastResult struct {
	BroadcastId string `json:"broadcast_id"`
	AgentId     string `json:"agent_id"`
	Sent        int64  `json:"sent"`
	Failed      int64  `json:"failed"`
	Status      string `json:"status"`
}

func (x *BroadcastResult) Reset()         {}
func (x *BroadcastResult) String() string { return x.BroadcastId }
func (x *BroadcastResult) ProtoMessage()  {}
func (x *BroadcastResult) GetBroadcastId() string {
	if x == nil {
		return ""
	}
	return x.BroadcastId
}
func (x *BroadcastResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *BroadcastResult) GetSent() int64 {
	if x == nil {
		return 0
	}
	return x.Sent
}
func (x *BroadcastResult) GetFailed() int64 {
	if x == nil {
		return 0
	}
	return x.Failed
}
func (x *BroadcastResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type AssignLeadFairRequest struct {
	AgentId string   `json:"agent_id"`
	LeadIds []string `json:"lead_ids"`
}

func (x *AssignLeadFairRequest) Reset()         {}
func (x *AssignLeadFairRequest) String() string { return x.AgentId }
func (x *AssignLeadFairRequest) ProtoMessage()  {}
func (x *AssignLeadFairRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AssignLeadFairRequest) GetLeadIds() []string {
	if x == nil {
		return nil
	}
	return x.LeadIds
}

type AssignLeadFairResult struct {
	Assigned int32    `json:"assigned"`
	LeadIds  []string `json:"lead_ids"`
}

func (x *AssignLeadFairResult) Reset()         {}
func (x *AssignLeadFairResult) String() string { return "" }
func (x *AssignLeadFairResult) ProtoMessage()  {}
func (x *AssignLeadFairResult) GetAssigned() int32 {
	if x == nil {
		return 0
	}
	return x.Assigned
}
func (x *AssignLeadFairResult) GetLeadIds() []string {
	if x == nil {
		return nil
	}
	return x.LeadIds
}

type CreateSlaRuleRequest struct {
	AgentId         string `json:"agent_id"`
	Name            string `json:"name"`
	ResponseMinutes int32  `json:"response_minutes"`
	EscalateAfter   int32  `json:"escalate_after"`
}

func (x *CreateSlaRuleRequest) Reset()         {}
func (x *CreateSlaRuleRequest) String() string { return x.Name }
func (x *CreateSlaRuleRequest) ProtoMessage()  {}
func (x *CreateSlaRuleRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CreateSlaRuleRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateSlaRuleRequest) GetResponseMinutes() int32 {
	if x == nil {
		return 0
	}
	return x.ResponseMinutes
}
func (x *CreateSlaRuleRequest) GetEscalateAfter() int32 {
	if x == nil {
		return 0
	}
	return x.EscalateAfter
}

type SlaRuleResult struct {
	RuleId          string `json:"rule_id"`
	AgentId         string `json:"agent_id"`
	Name            string `json:"name"`
	ResponseMinutes int32  `json:"response_minutes"`
	EscalateAfter   int32  `json:"escalate_after"`
	CreatedAt       string `json:"created_at"`
}

func (x *SlaRuleResult) Reset()         {}
func (x *SlaRuleResult) String() string { return x.RuleId }
func (x *SlaRuleResult) ProtoMessage()  {}
func (x *SlaRuleResult) GetRuleId() string {
	if x == nil {
		return ""
	}
	return x.RuleId
}
func (x *SlaRuleResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *SlaRuleResult) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *SlaRuleResult) GetResponseMinutes() int32 {
	if x == nil {
		return 0
	}
	return x.ResponseMinutes
}
func (x *SlaRuleResult) GetEscalateAfter() int32 {
	if x == nil {
		return 0
	}
	return x.EscalateAfter
}
func (x *SlaRuleResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type GetLeadTrailRequest struct {
	LeadId  string `json:"lead_id"`
	AgentId string `json:"agent_id"`
}

func (x *GetLeadTrailRequest) Reset()         {}
func (x *GetLeadTrailRequest) String() string { return x.LeadId }
func (x *GetLeadTrailRequest) ProtoMessage()  {}
func (x *GetLeadTrailRequest) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *GetLeadTrailRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}

type LeadTrailEvent struct {
	EventType string `json:"event_type"`
	Actor     string `json:"actor"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
}

type LeadTrailResult struct {
	LeadId string           `json:"lead_id"`
	Events []*LeadTrailEvent `json:"events"`
}

func (x *LeadTrailResult) Reset()         {}
func (x *LeadTrailResult) String() string { return x.LeadId }
func (x *LeadTrailResult) ProtoMessage()  {}
func (x *LeadTrailResult) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *LeadTrailResult) GetEvents() []*LeadTrailEvent {
	if x == nil {
		return nil
	}
	return x.Events
}

type TagLeadRequest struct {
	LeadId  string   `json:"lead_id"`
	AgentId string   `json:"agent_id"`
	Tags    []string `json:"tags"`
}

func (x *TagLeadRequest) Reset()         {}
func (x *TagLeadRequest) String() string { return x.LeadId }
func (x *TagLeadRequest) ProtoMessage()  {}
func (x *TagLeadRequest) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *TagLeadRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *TagLeadRequest) GetTags() []string {
	if x == nil {
		return nil
	}
	return x.Tags
}

type TagLeadResult struct {
	LeadId string   `json:"lead_id"`
	Tags   []string `json:"tags"`
}

func (x *TagLeadResult) Reset()         {}
func (x *TagLeadResult) String() string { return x.LeadId }
func (x *TagLeadResult) ProtoMessage()  {}
func (x *TagLeadResult) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *TagLeadResult) GetTags() []string {
	if x == nil {
		return nil
	}
	return x.Tags
}

type GenerateQuoteRequest struct {
	LeadId    string `json:"lead_id"`
	AgentId   string `json:"agent_id"`
	PackageId string `json:"package_id"`
	Pax       int32  `json:"pax"`
	Notes     string `json:"notes"`
}

func (x *GenerateQuoteRequest) Reset()         {}
func (x *GenerateQuoteRequest) String() string { return x.LeadId }
func (x *GenerateQuoteRequest) ProtoMessage()  {}
func (x *GenerateQuoteRequest) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *GenerateQuoteRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *GenerateQuoteRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *GenerateQuoteRequest) GetPax() int32 {
	if x == nil {
		return 0
	}
	return x.Pax
}
func (x *GenerateQuoteRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type QuoteResult struct {
	QuoteId    string  `json:"quote_id"`
	LeadId     string  `json:"lead_id"`
	AgentId    string  `json:"agent_id"`
	PackageId  string  `json:"package_id"`
	Pax        int32   `json:"pax"`
	TotalPrice float64 `json:"total_price"`
	PdfUrl     string  `json:"pdf_url"`
	ExpiresAt  string  `json:"expires_at"`
	CreatedAt  string  `json:"created_at"`
}

func (x *QuoteResult) Reset()         {}
func (x *QuoteResult) String() string { return x.QuoteId }
func (x *QuoteResult) ProtoMessage()  {}
func (x *QuoteResult) GetQuoteId() string {
	if x == nil {
		return ""
	}
	return x.QuoteId
}
func (x *QuoteResult) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *QuoteResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *QuoteResult) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *QuoteResult) GetPax() int32 {
	if x == nil {
		return 0
	}
	return x.Pax
}
func (x *QuoteResult) GetTotalPrice() float64 {
	if x == nil {
		return 0
	}
	return x.TotalPrice
}
func (x *QuoteResult) GetPdfUrl() string {
	if x == nil {
		return ""
	}
	return x.PdfUrl
}
func (x *QuoteResult) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}
func (x *QuoteResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type GetQuoteIdRequest struct {
	QuoteId string `json:"quote_id"`
}

func (x *GetQuoteIdRequest) Reset()         {}
func (x *GetQuoteIdRequest) String() string { return x.QuoteId }
func (x *GetQuoteIdRequest) ProtoMessage()  {}
func (x *GetQuoteIdRequest) GetQuoteId() string {
	if x == nil {
		return ""
	}
	return x.QuoteId
}

type BuildPaymentLinkRequest struct {
	QuoteId   string  `json:"quote_id"`
	AgentId   string  `json:"agent_id"`
	Amount    float64 `json:"amount"`
	ExpiresAt string  `json:"expires_at"`
}

func (x *BuildPaymentLinkRequest) Reset()         {}
func (x *BuildPaymentLinkRequest) String() string { return x.QuoteId }
func (x *BuildPaymentLinkRequest) ProtoMessage()  {}
func (x *BuildPaymentLinkRequest) GetQuoteId() string {
	if x == nil {
		return ""
	}
	return x.QuoteId
}
func (x *BuildPaymentLinkRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *BuildPaymentLinkRequest) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *BuildPaymentLinkRequest) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}

type PaymentLinkResult struct {
	LinkId    string  `json:"link_id"`
	QuoteId   string  `json:"quote_id"`
	Url       string  `json:"url"`
	Amount    float64 `json:"amount"`
	ExpiresAt string  `json:"expires_at"`
}

func (x *PaymentLinkResult) Reset()         {}
func (x *PaymentLinkResult) String() string { return x.LinkId }
func (x *PaymentLinkResult) ProtoMessage()  {}
func (x *PaymentLinkResult) GetLinkId() string {
	if x == nil {
		return ""
	}
	return x.LinkId
}
func (x *PaymentLinkResult) GetQuoteId() string {
	if x == nil {
		return ""
	}
	return x.QuoteId
}
func (x *PaymentLinkResult) GetUrl() string {
	if x == nil {
		return ""
	}
	return x.Url
}
func (x *PaymentLinkResult) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *PaymentLinkResult) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}

type RequestDiscountRequest struct {
	QuoteId    string  `json:"quote_id"`
	AgentId    string  `json:"agent_id"`
	Amount     float64 `json:"amount"`
	Reason     string  `json:"reason"`
}

func (x *RequestDiscountRequest) Reset()         {}
func (x *RequestDiscountRequest) String() string { return x.QuoteId }
func (x *RequestDiscountRequest) ProtoMessage()  {}
func (x *RequestDiscountRequest) GetQuoteId() string {
	if x == nil {
		return ""
	}
	return x.QuoteId
}
func (x *RequestDiscountRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *RequestDiscountRequest) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *RequestDiscountRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}

type ApproveDiscountRequest struct {
	DiscountId string `json:"discount_id"`
	ApproverId string `json:"approver_id"`
	Approved   bool   `json:"approved"`
	Note       string `json:"note"`
}

func (x *ApproveDiscountRequest) Reset()         {}
func (x *ApproveDiscountRequest) String() string { return x.DiscountId }
func (x *ApproveDiscountRequest) ProtoMessage()  {}
func (x *ApproveDiscountRequest) GetDiscountId() string {
	if x == nil {
		return ""
	}
	return x.DiscountId
}
func (x *ApproveDiscountRequest) GetApproverId() string {
	if x == nil {
		return ""
	}
	return x.ApproverId
}
func (x *ApproveDiscountRequest) GetApproved() bool {
	if x == nil {
		return false
	}
	return x.Approved
}
func (x *ApproveDiscountRequest) GetNote() string {
	if x == nil {
		return ""
	}
	return x.Note
}

type DiscountApprovalResult struct {
	DiscountId string  `json:"discount_id"`
	QuoteId    string  `json:"quote_id"`
	AgentId    string  `json:"agent_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	UpdatedAt  string  `json:"updated_at"`
}

func (x *DiscountApprovalResult) Reset()         {}
func (x *DiscountApprovalResult) String() string { return x.DiscountId }
func (x *DiscountApprovalResult) ProtoMessage()  {}
func (x *DiscountApprovalResult) GetDiscountId() string {
	if x == nil {
		return ""
	}
	return x.DiscountId
}
func (x *DiscountApprovalResult) GetQuoteId() string {
	if x == nil {
		return ""
	}
	return x.QuoteId
}
func (x *DiscountApprovalResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *DiscountApprovalResult) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *DiscountApprovalResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *DiscountApprovalResult) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

type StaleProspectRequest struct {
	AgentId     string `json:"agent_id"`
	StaleDays   int32  `json:"stale_days"`
	Page        int32  `json:"page"`
	PageSize    int32  `json:"page_size"`
}

func (x *StaleProspectRequest) Reset()         {}
func (x *StaleProspectRequest) String() string { return x.AgentId }
func (x *StaleProspectRequest) ProtoMessage()  {}
func (x *StaleProspectRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *StaleProspectRequest) GetStaleDays() int32 {
	if x == nil {
		return 0
	}
	return x.StaleDays
}
func (x *StaleProspectRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *StaleProspectRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type StaleProspectItem struct {
	LeadId      string `json:"lead_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	LastContact string `json:"last_contact"`
	DaysSince   int32  `json:"days_since"`
}

type StaleProspectListResult struct {
	Prospects []*StaleProspectItem `json:"prospects"`
	Total     int64                `json:"total"`
	Page      int32                `json:"page"`
	PageSize  int32                `json:"page_size"`
}

func (x *StaleProspectListResult) Reset()         {}
func (x *StaleProspectListResult) String() string { return "" }
func (x *StaleProspectListResult) ProtoMessage()  {}
func (x *StaleProspectListResult) GetProspects() []*StaleProspectItem {
	if x == nil {
		return nil
	}
	return x.Prospects
}
func (x *StaleProspectListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *StaleProspectListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *StaleProspectListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

// ---------------------------------------------------------------------------
// Commission / Finance
// ---------------------------------------------------------------------------

type CommissionBalanceResult struct {
	AgentId   string  `json:"agent_id"`
	Balance   float64 `json:"balance"`
	Pending   float64 `json:"pending"`
	Withdrawn float64 `json:"withdrawn"`
	UpdatedAt string  `json:"updated_at"`
}

func (x *CommissionBalanceResult) Reset()         {}
func (x *CommissionBalanceResult) String() string { return x.AgentId }
func (x *CommissionBalanceResult) ProtoMessage()  {}
func (x *CommissionBalanceResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CommissionBalanceResult) GetBalance() float64 {
	if x == nil {
		return 0
	}
	return x.Balance
}
func (x *CommissionBalanceResult) GetPending() float64 {
	if x == nil {
		return 0
	}
	return x.Pending
}
func (x *CommissionBalanceResult) GetWithdrawn() float64 {
	if x == nil {
		return 0
	}
	return x.Withdrawn
}
func (x *CommissionBalanceResult) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

type CommissionEventListRequest struct {
	AgentId  string `json:"agent_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *CommissionEventListRequest) Reset()         {}
func (x *CommissionEventListRequest) String() string { return x.AgentId }
func (x *CommissionEventListRequest) ProtoMessage()  {}
func (x *CommissionEventListRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CommissionEventListRequest) GetDateFrom() string {
	if x == nil {
		return ""
	}
	return x.DateFrom
}
func (x *CommissionEventListRequest) GetDateTo() string {
	if x == nil {
		return ""
	}
	return x.DateTo
}
func (x *CommissionEventListRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *CommissionEventListRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type CommissionEvent struct {
	EventId   string  `json:"event_id"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	BookingId string  `json:"booking_id"`
	CreatedAt string  `json:"created_at"`
}

type CommissionEventListResult struct {
	Events   []*CommissionEvent `json:"events"`
	Total    int64              `json:"total"`
	Page     int32              `json:"page"`
	PageSize int32              `json:"page_size"`
}

func (x *CommissionEventListResult) Reset()         {}
func (x *CommissionEventListResult) String() string { return "" }
func (x *CommissionEventListResult) ProtoMessage()  {}
func (x *CommissionEventListResult) GetEvents() []*CommissionEvent {
	if x == nil {
		return nil
	}
	return x.Events
}
func (x *CommissionEventListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *CommissionEventListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *CommissionEventListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type PayoutRequest struct {
	AgentId       string  `json:"agent_id"`
	Amount        float64 `json:"amount"`
	BankAccount   string  `json:"bank_account"`
	BankCode      string  `json:"bank_code"`
}

func (x *PayoutRequest) Reset()         {}
func (x *PayoutRequest) String() string { return x.AgentId }
func (x *PayoutRequest) ProtoMessage()  {}
func (x *PayoutRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *PayoutRequest) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *PayoutRequest) GetBankAccount() string {
	if x == nil {
		return ""
	}
	return x.BankAccount
}
func (x *PayoutRequest) GetBankCode() string {
	if x == nil {
		return ""
	}
	return x.BankCode
}

type PayoutResult struct {
	PayoutId  string  `json:"payout_id"`
	AgentId   string  `json:"agent_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

func (x *PayoutResult) Reset()         {}
func (x *PayoutResult) String() string { return x.PayoutId }
func (x *PayoutResult) ProtoMessage()  {}
func (x *PayoutResult) GetPayoutId() string {
	if x == nil {
		return ""
	}
	return x.PayoutId
}
func (x *PayoutResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *PayoutResult) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *PayoutResult) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *PayoutResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type PayoutHistoryResult struct {
	AgentId  string          `json:"agent_id"`
	Payouts  []*PayoutResult `json:"payouts"`
	Total    int64           `json:"total"`
}

func (x *PayoutHistoryResult) Reset()         {}
func (x *PayoutHistoryResult) String() string { return x.AgentId }
func (x *PayoutHistoryResult) ProtoMessage()  {}
func (x *PayoutHistoryResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *PayoutHistoryResult) GetPayouts() []*PayoutResult {
	if x == nil {
		return nil
	}
	return x.Payouts
}
func (x *PayoutHistoryResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}

type OverrideCommissionRequest struct {
	AgentId    string  `json:"agent_id"`
	SubAgentId string  `json:"sub_agent_id"`
	BookingId  string  `json:"booking_id"`
	Rate       float64 `json:"rate"`
}

func (x *OverrideCommissionRequest) Reset()         {}
func (x *OverrideCommissionRequest) String() string { return x.AgentId }
func (x *OverrideCommissionRequest) ProtoMessage()  {}
func (x *OverrideCommissionRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *OverrideCommissionRequest) GetSubAgentId() string {
	if x == nil {
		return ""
	}
	return x.SubAgentId
}
func (x *OverrideCommissionRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *OverrideCommissionRequest) GetRate() float64 {
	if x == nil {
		return 0
	}
	return x.Rate
}

type OverrideCommissionResult struct {
	AgentId    string  `json:"agent_id"`
	SubAgentId string  `json:"sub_agent_id"`
	BookingId  string  `json:"booking_id"`
	Amount     float64 `json:"amount"`
	Rate       float64 `json:"rate"`
}

func (x *OverrideCommissionResult) Reset()         {}
func (x *OverrideCommissionResult) String() string { return x.AgentId }
func (x *OverrideCommissionResult) ProtoMessage()  {}
func (x *OverrideCommissionResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *OverrideCommissionResult) GetSubAgentId() string {
	if x == nil {
		return ""
	}
	return x.SubAgentId
}
func (x *OverrideCommissionResult) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *OverrideCommissionResult) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *OverrideCommissionResult) GetRate() float64 {
	if x == nil {
		return 0
	}
	return x.Rate
}

type RoasReportRequest struct {
	AgentId  string `json:"agent_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (x *RoasReportRequest) Reset()         {}
func (x *RoasReportRequest) String() string { return x.AgentId }
func (x *RoasReportRequest) ProtoMessage()  {}
func (x *RoasReportRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *RoasReportRequest) GetDateFrom() string {
	if x == nil {
		return ""
	}
	return x.DateFrom
}
func (x *RoasReportRequest) GetDateTo() string {
	if x == nil {
		return ""
	}
	return x.DateTo
}

type RoasReportResult struct {
	AgentId     string  `json:"agent_id"`
	AdSpend     float64 `json:"ad_spend"`
	Revenue     float64 `json:"revenue"`
	Roas        float64 `json:"roas"`
	Conversions int64   `json:"conversions"`
}

func (x *RoasReportResult) Reset()         {}
func (x *RoasReportResult) String() string { return x.AgentId }
func (x *RoasReportResult) ProtoMessage()  {}
func (x *RoasReportResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *RoasReportResult) GetAdSpend() float64 {
	if x == nil {
		return 0
	}
	return x.AdSpend
}
func (x *RoasReportResult) GetRevenue() float64 {
	if x == nil {
		return 0
	}
	return x.Revenue
}
func (x *RoasReportResult) GetRoas() float64 {
	if x == nil {
		return 0
	}
	return x.Roas
}
func (x *RoasReportResult) GetConversions() int64 {
	if x == nil {
		return 0
	}
	return x.Conversions
}

// ---------------------------------------------------------------------------
// Academy / Gamification
// ---------------------------------------------------------------------------

type AcademyListRequest struct {
	AgentId  string `json:"agent_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *AcademyListRequest) Reset()         {}
func (x *AcademyListRequest) String() string { return x.AgentId }
func (x *AcademyListRequest) ProtoMessage()  {}
func (x *AcademyListRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AcademyListRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *AcademyListRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type AcademyCourseItem struct {
	CourseId    string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int32  `json:"duration"`
	Level       string `json:"level"`
}

type AcademyCourseListResult struct {
	Courses  []*AcademyCourseItem `json:"courses"`
	Total    int64                `json:"total"`
	Page     int32                `json:"page"`
	PageSize int32                `json:"page_size"`
}

func (x *AcademyCourseListResult) Reset()         {}
func (x *AcademyCourseListResult) String() string { return "" }
func (x *AcademyCourseListResult) ProtoMessage()  {}
func (x *AcademyCourseListResult) GetCourses() []*AcademyCourseItem {
	if x == nil {
		return nil
	}
	return x.Courses
}
func (x *AcademyCourseListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *AcademyCourseListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *AcademyCourseListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type CourseProgressRequest struct {
	AgentId  string `json:"agent_id"`
	CourseId string `json:"course_id"`
}

func (x *CourseProgressRequest) Reset()         {}
func (x *CourseProgressRequest) String() string { return x.CourseId }
func (x *CourseProgressRequest) ProtoMessage()  {}
func (x *CourseProgressRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CourseProgressRequest) GetCourseId() string {
	if x == nil {
		return ""
	}
	return x.CourseId
}

type CourseProgressResult struct {
	AgentId     string  `json:"agent_id"`
	CourseId    string  `json:"course_id"`
	Progress    float64 `json:"progress"`
	Completed   bool    `json:"completed"`
	LastAccessAt string `json:"last_access_at"`
}

func (x *CourseProgressResult) Reset()         {}
func (x *CourseProgressResult) String() string { return x.CourseId }
func (x *CourseProgressResult) ProtoMessage()  {}
func (x *CourseProgressResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CourseProgressResult) GetCourseId() string {
	if x == nil {
		return ""
	}
	return x.CourseId
}
func (x *CourseProgressResult) GetProgress() float64 {
	if x == nil {
		return 0
	}
	return x.Progress
}
func (x *CourseProgressResult) GetCompleted() bool {
	if x == nil {
		return false
	}
	return x.Completed
}
func (x *CourseProgressResult) GetLastAccessAt() string {
	if x == nil {
		return ""
	}
	return x.LastAccessAt
}

type SubmitQuizRequest struct {
	AgentId  string            `json:"agent_id"`
	CourseId string            `json:"course_id"`
	Answers  map[string]string `json:"answers"`
}

func (x *SubmitQuizRequest) Reset()         {}
func (x *SubmitQuizRequest) String() string { return x.CourseId }
func (x *SubmitQuizRequest) ProtoMessage()  {}
func (x *SubmitQuizRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *SubmitQuizRequest) GetCourseId() string {
	if x == nil {
		return ""
	}
	return x.CourseId
}
func (x *SubmitQuizRequest) GetAnswers() map[string]string {
	if x == nil {
		return nil
	}
	return x.Answers
}

type QuizResult struct {
	AgentId   string  `json:"agent_id"`
	CourseId  string  `json:"course_id"`
	Score     float64 `json:"score"`
	Passed    bool    `json:"passed"`
	Badge     string  `json:"badge"`
	TakenAt   string  `json:"taken_at"`
}

func (x *QuizResult) Reset()         {}
func (x *QuizResult) String() string { return x.CourseId }
func (x *QuizResult) ProtoMessage()  {}
func (x *QuizResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *QuizResult) GetCourseId() string {
	if x == nil {
		return ""
	}
	return x.CourseId
}
func (x *QuizResult) GetScore() float64 {
	if x == nil {
		return 0
	}
	return x.Score
}
func (x *QuizResult) GetPassed() bool {
	if x == nil {
		return false
	}
	return x.Passed
}
func (x *QuizResult) GetBadge() string {
	if x == nil {
		return ""
	}
	return x.Badge
}
func (x *QuizResult) GetTakenAt() string {
	if x == nil {
		return ""
	}
	return x.TakenAt
}

type SalesScriptListRequest struct {
	AgentId  string `json:"agent_id"`
	Category string `json:"category"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (x *SalesScriptListRequest) Reset()         {}
func (x *SalesScriptListRequest) String() string { return x.AgentId }
func (x *SalesScriptListRequest) ProtoMessage()  {}
func (x *SalesScriptListRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *SalesScriptListRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *SalesScriptListRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *SalesScriptListRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type SalesScriptItem struct {
	ScriptId string `json:"script_id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

type SalesScriptListResult struct {
	Scripts  []*SalesScriptItem `json:"scripts"`
	Total    int64              `json:"total"`
	Page     int32              `json:"page"`
	PageSize int32              `json:"page_size"`
}

func (x *SalesScriptListResult) Reset()         {}
func (x *SalesScriptListResult) String() string { return "" }
func (x *SalesScriptListResult) ProtoMessage()  {}
func (x *SalesScriptListResult) GetScripts() []*SalesScriptItem {
	if x == nil {
		return nil
	}
	return x.Scripts
}
func (x *SalesScriptListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *SalesScriptListResult) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *SalesScriptListResult) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type LeaderboardRequest struct {
	AgentId string `json:"agent_id"`
	Period  string `json:"period"`
	Limit   int32  `json:"limit"`
}

func (x *LeaderboardRequest) Reset()         {}
func (x *LeaderboardRequest) String() string { return x.AgentId }
func (x *LeaderboardRequest) ProtoMessage()  {}
func (x *LeaderboardRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *LeaderboardRequest) GetPeriod() string {
	if x == nil {
		return ""
	}
	return x.Period
}
func (x *LeaderboardRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}

type LeaderboardEntry struct {
	Rank      int32   `json:"rank"`
	AgentId   string  `json:"agent_id"`
	Name      string  `json:"name"`
	Score     float64 `json:"score"`
	Bookings  int64   `json:"bookings"`
}

type LeaderboardResult struct {
	Period  string              `json:"period"`
	Entries []*LeaderboardEntry `json:"entries"`
	MyRank  int32               `json:"my_rank"`
}

func (x *LeaderboardResult) Reset()         {}
func (x *LeaderboardResult) String() string { return x.Period }
func (x *LeaderboardResult) ProtoMessage()  {}
func (x *LeaderboardResult) GetPeriod() string {
	if x == nil {
		return ""
	}
	return x.Period
}
func (x *LeaderboardResult) GetEntries() []*LeaderboardEntry {
	if x == nil {
		return nil
	}
	return x.Entries
}
func (x *LeaderboardResult) GetMyRank() int32 {
	if x == nil {
		return 0
	}
	return x.MyRank
}

type AgentTierResult struct {
	AgentId       string  `json:"agent_id"`
	Tier          string  `json:"tier"`
	Points        float64 `json:"points"`
	NextTier      string  `json:"next_tier"`
	PointsToNext  float64 `json:"points_to_next"`
}

func (x *AgentTierResult) Reset()         {}
func (x *AgentTierResult) String() string { return x.AgentId }
func (x *AgentTierResult) ProtoMessage()  {}
func (x *AgentTierResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AgentTierResult) GetTier() string {
	if x == nil {
		return ""
	}
	return x.Tier
}
func (x *AgentTierResult) GetPoints() float64 {
	if x == nil {
		return 0
	}
	return x.Points
}
func (x *AgentTierResult) GetNextTier() string {
	if x == nil {
		return ""
	}
	return x.NextTier
}
func (x *AgentTierResult) GetPointsToNext() float64 {
	if x == nil {
		return 0
	}
	return x.PointsToNext
}

// ---------------------------------------------------------------------------
// Alumni / Other
// ---------------------------------------------------------------------------

type AlumniReferralItem struct {
	ReferralId string `json:"referral_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type AlumniReferralListResult struct {
	AgentId   string                `json:"agent_id"`
	Referrals []*AlumniReferralItem `json:"referrals"`
	Total     int64                 `json:"total"`
}

func (x *AlumniReferralListResult) Reset()         {}
func (x *AlumniReferralListResult) String() string { return x.AgentId }
func (x *AlumniReferralListResult) ProtoMessage()  {}
func (x *AlumniReferralListResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *AlumniReferralListResult) GetReferrals() []*AlumniReferralItem {
	if x == nil {
		return nil
	}
	return x.Referrals
}
func (x *AlumniReferralListResult) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}

type ReturnIntentResult struct {
	AgentId        string  `json:"agent_id"`
	SavingsBalance float64 `json:"savings_balance"`
	TargetAmount   float64 `json:"target_amount"`
	Progress       float64 `json:"progress"`
	EstimatedDate  string  `json:"estimated_date"`
}

func (x *ReturnIntentResult) Reset()         {}
func (x *ReturnIntentResult) String() string { return x.AgentId }
func (x *ReturnIntentResult) ProtoMessage()  {}
func (x *ReturnIntentResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ReturnIntentResult) GetSavingsBalance() float64 {
	if x == nil {
		return 0
	}
	return x.SavingsBalance
}
func (x *ReturnIntentResult) GetTargetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.TargetAmount
}
func (x *ReturnIntentResult) GetProgress() float64 {
	if x == nil {
		return 0
	}
	return x.Progress
}
func (x *ReturnIntentResult) GetEstimatedDate() string {
	if x == nil {
		return ""
	}
	return x.EstimatedDate
}

type ZakatCalculatorRequest struct {
	AgentId        string  `json:"agent_id"`
	GoldGrams      float64 `json:"gold_grams"`
	SilverGrams    float64 `json:"silver_grams"`
	CashAmount     float64 `json:"cash_amount"`
	BusinessAssets float64 `json:"business_assets"`
	Debts          float64 `json:"debts"`
}

func (x *ZakatCalculatorRequest) Reset()         {}
func (x *ZakatCalculatorRequest) String() string { return x.AgentId }
func (x *ZakatCalculatorRequest) ProtoMessage()  {}
func (x *ZakatCalculatorRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ZakatCalculatorRequest) GetGoldGrams() float64 {
	if x == nil {
		return 0
	}
	return x.GoldGrams
}
func (x *ZakatCalculatorRequest) GetSilverGrams() float64 {
	if x == nil {
		return 0
	}
	return x.SilverGrams
}
func (x *ZakatCalculatorRequest) GetCashAmount() float64 {
	if x == nil {
		return 0
	}
	return x.CashAmount
}
func (x *ZakatCalculatorRequest) GetBusinessAssets() float64 {
	if x == nil {
		return 0
	}
	return x.BusinessAssets
}
func (x *ZakatCalculatorRequest) GetDebts() float64 {
	if x == nil {
		return 0
	}
	return x.Debts
}

type ZakatResult struct {
	AgentId      string  `json:"agent_id"`
	TotalAssets  float64 `json:"total_assets"`
	NetAssets    float64 `json:"net_assets"`
	ZakatDue     float64 `json:"zakat_due"`
	NisabMet     bool    `json:"nisab_met"`
}

func (x *ZakatResult) Reset()         {}
func (x *ZakatResult) String() string { return x.AgentId }
func (x *ZakatResult) ProtoMessage()  {}
func (x *ZakatResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *ZakatResult) GetTotalAssets() float64 {
	if x == nil {
		return 0
	}
	return x.TotalAssets
}
func (x *ZakatResult) GetNetAssets() float64 {
	if x == nil {
		return 0
	}
	return x.NetAssets
}
func (x *ZakatResult) GetZakatDue() float64 {
	if x == nil {
		return 0
	}
	return x.ZakatDue
}
func (x *ZakatResult) GetNisabMet() bool {
	if x == nil {
		return false
	}
	return x.NisabMet
}

type CharityRequest struct {
	AgentId   string  `json:"agent_id"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	Notes     string  `json:"notes"`
}

func (x *CharityRequest) Reset()         {}
func (x *CharityRequest) String() string { return x.AgentId }
func (x *CharityRequest) ProtoMessage()  {}
func (x *CharityRequest) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CharityRequest) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *CharityRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *CharityRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type CharityResult struct {
	CharityId string  `json:"charity_id"`
	AgentId   string  `json:"agent_id"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	CreatedAt string  `json:"created_at"`
}

func (x *CharityResult) Reset()         {}
func (x *CharityResult) String() string { return x.CharityId }
func (x *CharityResult) ProtoMessage()  {}
func (x *CharityResult) GetCharityId() string {
	if x == nil {
		return ""
	}
	return x.CharityId
}
func (x *CharityResult) GetAgentId() string {
	if x == nil {
		return ""
	}
	return x.AgentId
}
func (x *CharityResult) GetAmount() float64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *CharityResult) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *CharityResult) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
