// crm_depth_adapter.go — gateway CRM adapter methods for Wave 7 depth RPCs.
// BL-CRM-010..012, 017..066: agency registration, content, leads/CRM,
// commission, academy, alumni/other.

package crm_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/crm_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Error mapper
// ---------------------------------------------------------------------------

func mapCrmDepthError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("crm-depth call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return fmt.Errorf("%w: %s", apperrors.ErrNotFound, st.Message())
	case grpcCodes.InvalidArgument:
		return fmt.Errorf("%w: %s", apperrors.ErrValidation, st.Message())
	case grpcCodes.PermissionDenied:
		return fmt.Errorf("%w: %s", apperrors.ErrForbidden, st.Message())
	default:
		return fmt.Errorf("%w: crm-depth %s: %s", apperrors.ErrInternal, st.Code(), st.Message())
	}
}

// ---------------------------------------------------------------------------
// Agency Registration
// ---------------------------------------------------------------------------

type RegisterAgentParams struct {
	AgencyName string
	OwnerName  string
	Phone      string
	Email      string
}

type AgentAdapterResult struct {
	ID         string
	AgencyName string
	Status     string
	TierLevel  string
	CreatedAt  string
}

func (a *Adapter) RegisterAgent(ctx context.Context, params *RegisterAgentParams) (*AgentAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.RegisterAgent"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.RegisterAgent(ctx, &pb.CrmDepthAgentRegisterRequest{
		AgencyName: params.AgencyName,
		OwnerName:  params.OwnerName,
		Phone:      params.Phone,
		Email:      params.Email,
	})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AgentAdapterResult{ID: resp.GetId(), AgencyName: resp.GetAgencyName(), Status: resp.GetStatus(), TierLevel: resp.GetTierLevel(), CreatedAt: resp.GetCreatedAt()}, nil
}

type SubmitAgentKycParams struct {
	AgentId string
	KycType string
	DocUrl  string
}

type KycAdapterResult struct {
	AgentId   string
	KycStatus string
	UpdatedAt string
}

func (a *Adapter) SubmitAgentKyc(ctx context.Context, params *SubmitAgentKycParams) (*KycAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.SubmitAgentKyc"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.SubmitAgentKyc(ctx, &pb.CrmDepthSubmitKycRequest{AgentId: params.AgentId, KycType: params.KycType, DocUrl: params.DocUrl})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &KycAdapterResult{AgentId: resp.GetAgentId(), KycStatus: resp.GetKycStatus(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

type SignAgentMoUParams struct {
	AgentId     string
	SignedDocUrl string
}

type MoUAdapterResult struct {
	AgentId   string
	MouStatus string
	SignedAt  string
}

func (a *Adapter) SignAgentMoU(ctx context.Context, params *SignAgentMoUParams) (*MoUAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.SignAgentMoU"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.SignAgentMoU(ctx, &pb.CrmDepthSignMoURequest{AgentId: params.AgentId, SignedDocUrl: params.SignedDocUrl})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &MoUAdapterResult{AgentId: resp.GetAgentId(), MouStatus: resp.GetMouStatus(), SignedAt: resp.GetSignedAt()}, nil
}

func (a *Adapter) GetAgentProfile(ctx context.Context, agentId string) (*AgentAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetAgentProfile"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetAgentProfile(ctx, &pb.CrmDepthGetAgentRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AgentAdapterResult{ID: resp.GetId(), AgencyName: resp.GetAgencyName(), Status: resp.GetStatus(), TierLevel: resp.GetTierLevel(), CreatedAt: resp.GetCreatedAt()}, nil
}

type ReplicaSiteAdapterResult struct {
	AgentId   string
	SiteUrl   string
	UpdatedAt string
}

func (a *Adapter) GetReplicaSite(ctx context.Context, agentId string) (*ReplicaSiteAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetReplicaSite"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetReplicaSite(ctx, &pb.CrmDepthGetAgentRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ReplicaSiteAdapterResult{AgentId: resp.GetAgentId(), SiteUrl: resp.GetSiteUrl(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

type UpdateReplicaSiteParams struct {
	AgentId    string
	CustomSlug string
	ThemeColor string
}

func (a *Adapter) UpdateReplicaSite(ctx context.Context, params *UpdateReplicaSiteParams) (*ReplicaSiteAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.UpdateReplicaSite"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.UpdateReplicaSite(ctx, &pb.CrmDepthUpdateReplicaSiteRequest{AgentId: params.AgentId, CustomSlug: params.CustomSlug, ThemeColor: params.ThemeColor})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ReplicaSiteAdapterResult{AgentId: resp.GetAgentId(), SiteUrl: resp.GetSiteUrl(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

// ---------------------------------------------------------------------------
// Content
// ---------------------------------------------------------------------------

type GetSocialShareLinkParams struct {
	AgentId   string
	Platform  string
	ContentId string
}

type ShareLinkAdapterResult struct {
	ShareUrl  string
	ExpiresAt string
}

func (a *Adapter) GetSocialShareLink(ctx context.Context, params *GetSocialShareLinkParams) (*ShareLinkAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetSocialShareLink"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetSocialShareLink(ctx, &pb.CrmDepthShareLinkRequest{AgentId: params.AgentId, Platform: params.Platform, ContentId: params.ContentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ShareLinkAdapterResult{ShareUrl: resp.GetShareUrl(), ExpiresAt: resp.GetExpiresAt()}, nil
}

type GenerateBusinessCardParams struct {
	AgentId    string
	TemplateId string
}

type BusinessCardAdapterResult struct {
	CardUrl   string
	CreatedAt string
}

func (a *Adapter) GenerateBusinessCard(ctx context.Context, params *GenerateBusinessCardParams) (*BusinessCardAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GenerateBusinessCard"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GenerateBusinessCard(ctx, &pb.CrmDepthBusinessCardRequest{AgentId: params.AgentId, TemplateId: params.TemplateId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &BusinessCardAdapterResult{CardUrl: resp.GetCardUrl(), CreatedAt: resp.GetCreatedAt()}, nil
}

type ListContentBankParams struct {
	AgentId  string
	Category string
	Page     int32
	PageSize int32
}

type ContentAssetRowResult struct {
	Id        string
	Title     string
	Category  string
	AssetUrl  string
	CreatedAt string
}

type ListContentBankResult struct {
	Items    []ContentAssetRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListContentBank(ctx context.Context, params *ListContentBankParams) (*ListContentBankResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListContentBank"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListContentBank(ctx, &pb.CrmDepthListContentBankRequest{AgentId: params.AgentId, Category: params.Category, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]ContentAssetRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, ContentAssetRowResult{Id: r.GetId(), Title: r.GetTitle(), Category: r.GetCategory(), AssetUrl: r.GetAssetUrl(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListContentBankResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type CreateContentAssetParams struct {
	AgentId  string
	Title    string
	Category string
	AssetUrl string
}

func (a *Adapter) CreateContentAsset(ctx context.Context, params *CreateContentAssetParams) (*ContentAssetRowResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateContentAsset"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateContentAsset(ctx, &pb.CrmDepthCreateContentAssetRequest{AgentId: params.AgentId, Title: params.Title, Category: params.Category, AssetUrl: params.AssetUrl})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ContentAssetRowResult{Id: resp.GetId(), Title: resp.GetTitle(), Category: resp.GetCategory(), AssetUrl: resp.GetAssetUrl(), CreatedAt: resp.GetCreatedAt()}, nil
}

type WatermarkFlyerParams struct {
	AgentId       string
	FlyerUrl      string
	WatermarkText string
}

type WatermarkFlyerResult struct {
	WatermarkedUrl string
}

func (a *Adapter) WatermarkFlyer(ctx context.Context, params *WatermarkFlyerParams) (*WatermarkFlyerResult, error) {
	const op = "crm_grpc_adapter.Adapter.WatermarkFlyer"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.WatermarkFlyer(ctx, &pb.CrmDepthWatermarkRequest{AgentId: params.AgentId, FlyerUrl: params.FlyerUrl, WatermarkText: params.WatermarkText})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &WatermarkFlyerResult{WatermarkedUrl: resp.GetWatermarkedUrl()}, nil
}

type ListProgramGalleryParams struct {
	PackageId string
	Page      int32
	PageSize  int32
}

type GalleryRowResult struct {
	Id        string
	ImageUrl  string
	Caption   string
	CreatedAt string
}

type ListProgramGalleryResult struct {
	Items    []GalleryRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListProgramGallery(ctx context.Context, params *ListProgramGalleryParams) (*ListProgramGalleryResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListProgramGallery"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListProgramGallery(ctx, &pb.CrmDepthListGalleryRequest{PackageId: params.PackageId, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]GalleryRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, GalleryRowResult{Id: r.GetId(), ImageUrl: r.GetImageUrl(), Caption: r.GetCaption(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListProgramGalleryResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type SetTrackingCodeParams struct {
	AgentId      string
	TrackingCode string
	Platform     string
}

type TrackingCodeAdapterResult struct {
	AgentId      string
	TrackingCode string
	UpdatedAt    string
}

func (a *Adapter) SetTrackingCode(ctx context.Context, params *SetTrackingCodeParams) (*TrackingCodeAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.SetTrackingCode"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.SetTrackingCode(ctx, &pb.CrmDepthSetTrackingCodeRequest{AgentId: params.AgentId, TrackingCode: params.TrackingCode, Platform: params.Platform})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &TrackingCodeAdapterResult{AgentId: resp.GetAgentId(), TrackingCode: resp.GetTrackingCode(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

type GetAdsManagerStatsParams struct {
	AgentId  string
	DateFrom string
	DateTo   string
}

type AdsStatsAdapterResult struct {
	Impressions int64
	Clicks      int64
	Spend       float64
	Conversions int64
	Roas        float64
}

func (a *Adapter) GetAdsManagerStats(ctx context.Context, params *GetAdsManagerStatsParams) (*AdsStatsAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetAdsManagerStats"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetAdsManagerStats(ctx, &pb.CrmDepthAdsStatsRequest{AgentId: params.AgentId, DateFrom: params.DateFrom, DateTo: params.DateTo})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AdsStatsAdapterResult{Impressions: resp.GetImpressions(), Clicks: resp.GetClicks(), Spend: resp.GetSpend(), Conversions: resp.GetConversions(), Roas: resp.GetRoas()}, nil
}

type CreateUtmLinkParams struct {
	AgentId     string
	DestUrl     string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string
}

type UtmLinkAdapterResult struct {
	Id          string
	AgentId     string
	DestUrl     string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string
	ShortUrl    string
	CreatedAt   string
}

func (a *Adapter) CreateUtmLink(ctx context.Context, params *CreateUtmLinkParams) (*UtmLinkAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateUtmLink"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateUtmLink(ctx, &pb.CrmDepthCreateUtmLinkRequest{AgentId: params.AgentId, DestUrl: params.DestUrl, UtmSource: params.UtmSource, UtmMedium: params.UtmMedium, UtmCampaign: params.UtmCampaign})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &UtmLinkAdapterResult{Id: resp.GetId(), AgentId: resp.GetAgentId(), DestUrl: resp.GetDestUrl(), UtmSource: resp.GetUtmSource(), UtmMedium: resp.GetUtmMedium(), UtmCampaign: resp.GetUtmCampaign(), ShortUrl: resp.GetShortUrl(), CreatedAt: resp.GetCreatedAt()}, nil
}

type ListUtmLinksParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type ListUtmLinksResult struct {
	Items    []UtmLinkAdapterResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListUtmLinks(ctx context.Context, params *ListUtmLinksParams) (*ListUtmLinksResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListUtmLinks"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListUtmLinks(ctx, &pb.CrmDepthListUtmLinksRequest{AgentId: params.AgentId, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]UtmLinkAdapterResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, UtmLinkAdapterResult{Id: r.GetId(), AgentId: r.GetAgentId(), DestUrl: r.GetDestUrl(), UtmSource: r.GetUtmSource(), UtmMedium: r.GetUtmMedium(), UtmCampaign: r.GetUtmCampaign(), ShortUrl: r.GetShortUrl(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListUtmLinksResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type CreateLandingPageParams struct {
	AgentId    string
	Title      string
	TemplateId string
	PackageId  string
}

type LandingPageAdapterResult struct {
	Id        string
	AgentId   string
	Title     string
	PageUrl   string
	CreatedAt string
}

func (a *Adapter) CreateLandingPage(ctx context.Context, params *CreateLandingPageParams) (*LandingPageAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateLandingPage"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateLandingPage(ctx, &pb.CrmDepthCreateLandingPageRequest{AgentId: params.AgentId, Title: params.Title, TemplateId: params.TemplateId, PackageId: params.PackageId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &LandingPageAdapterResult{Id: resp.GetId(), AgentId: resp.GetAgentId(), Title: resp.GetTitle(), PageUrl: resp.GetPageUrl(), CreatedAt: resp.GetCreatedAt()}, nil
}

type ListLandingPagesParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type ListLandingPagesResult struct {
	Items    []LandingPageAdapterResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListLandingPages(ctx context.Context, params *ListLandingPagesParams) (*ListLandingPagesResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListLandingPages"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListLandingPages(ctx, &pb.CrmDepthListLandingPagesRequest{AgentId: params.AgentId, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]LandingPageAdapterResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, LandingPageAdapterResult{Id: r.GetId(), AgentId: r.GetAgentId(), Title: r.GetTitle(), PageUrl: r.GetPageUrl(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListLandingPagesResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type ScheduleContentParams struct {
	AgentId     string
	ContentId   string
	Platform    string
	ScheduledAt string
}

type ScheduleContentAdapterResult struct {
	ScheduleId string
	Status     string
}

func (a *Adapter) ScheduleContent(ctx context.Context, params *ScheduleContentParams) (*ScheduleContentAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.ScheduleContent"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ScheduleContent(ctx, &pb.CrmDepthScheduleContentRequest{AgentId: params.AgentId, ContentId: params.ContentId, Platform: params.Platform, ScheduledAt: params.ScheduledAt})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ScheduleContentAdapterResult{ScheduleId: resp.GetScheduleId(), Status: resp.GetStatus()}, nil
}

type ListScheduledContentParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type ScheduledContentRowResult struct {
	Id          string
	ContentId   string
	Platform    string
	ScheduledAt string
	Status      string
}

type ListScheduledContentResult struct {
	Items    []ScheduledContentRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListScheduledContent(ctx context.Context, params *ListScheduledContentParams) (*ListScheduledContentResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListScheduledContent"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListScheduledContent(ctx, &pb.CrmDepthListScheduledContentRequest{AgentId: params.AgentId, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]ScheduledContentRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, ScheduledContentRowResult{Id: r.GetId(), ContentId: r.GetContentId(), Platform: r.GetPlatform(), ScheduledAt: r.GetScheduledAt(), Status: r.GetStatus()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListScheduledContentResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type GetContentAnalyticsParams struct {
	AgentId  string
	DateFrom string
	DateTo   string
}

type ContentAnalyticsAdapterResult struct {
	TotalViews     int64
	TotalShares    int64
	TotalClicks    int64
	EngagementRate float64
}

func (a *Adapter) GetContentAnalytics(ctx context.Context, params *GetContentAnalyticsParams) (*ContentAnalyticsAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetContentAnalytics"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetContentAnalytics(ctx, &pb.CrmDepthContentAnalyticsRequest{AgentId: params.AgentId, DateFrom: params.DateFrom, DateTo: params.DateTo})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ContentAnalyticsAdapterResult{TotalViews: resp.GetTotalViews(), TotalShares: resp.GetTotalShares(), TotalClicks: resp.GetTotalClicks(), EngagementRate: resp.GetEngagementRate()}, nil
}

// ---------------------------------------------------------------------------
// CRM/Leads
// ---------------------------------------------------------------------------

type CreateAgentLeadParams struct {
	AgentId   string
	Name      string
	Phone     string
	Email     string
	Source    string
	PackageId string
	Notes     string
}

type AgentLeadAdapterResult struct {
	Id        string
	AgentId   string
	Name      string
	Phone     string
	Email     string
	Status    string
	Source    string
	CreatedAt string
	UpdatedAt string
}

func (a *Adapter) CreateAgentLead(ctx context.Context, params *CreateAgentLeadParams) (*AgentLeadAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateAgentLead"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateAgentLead(ctx, &pb.CrmDepthCreateAgentLeadRequest{AgentId: params.AgentId, Name: params.Name, Phone: params.Phone, Email: params.Email, Source: params.Source, PackageId: params.PackageId, Notes: params.Notes})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AgentLeadAdapterResult{Id: resp.GetId(), AgentId: resp.GetAgentId(), Name: resp.GetName(), Phone: resp.GetPhone(), Email: resp.GetEmail(), Status: resp.GetStatus(), Source: resp.GetSource(), CreatedAt: resp.GetCreatedAt(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

type ListAgentLeadsParams struct {
	AgentId      string
	StatusFilter string
	Page         int32
	PageSize     int32
}

type ListAgentLeadsResult struct {
	Items    []AgentLeadAdapterResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListAgentLeads(ctx context.Context, params *ListAgentLeadsParams) (*ListAgentLeadsResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListAgentLeads"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListAgentLeads(ctx, &pb.CrmDepthListAgentLeadsRequest{AgentId: params.AgentId, StatusFilter: params.StatusFilter, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]AgentLeadAdapterResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, AgentLeadAdapterResult{Id: r.GetId(), AgentId: r.GetAgentId(), Name: r.GetName(), Phone: r.GetPhone(), Email: r.GetEmail(), Status: r.GetStatus(), Source: r.GetSource(), CreatedAt: r.GetCreatedAt(), UpdatedAt: r.GetUpdatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListAgentLeadsResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type SetLeadReminderParams struct {
	LeadId   string
	AgentId  string
	RemindAt string
	Message  string
}

type SetLeadReminderAdapterResult struct {
	ReminderId string
	Status     string
}

func (a *Adapter) SetLeadReminder(ctx context.Context, params *SetLeadReminderParams) (*SetLeadReminderAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.SetLeadReminder"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.SetLeadReminder(ctx, &pb.CrmDepthSetLeadReminderRequest{LeadId: params.LeadId, AgentId: params.AgentId, RemindAt: params.RemindAt, Message: params.Message})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SetLeadReminderAdapterResult{ReminderId: resp.GetReminderId(), Status: resp.GetStatus()}, nil
}

type ReminderRowResult struct {
	Id       string
	LeadId   string
	RemindAt string
	Message  string
	Status   string
}

type ListLeadRemindersResult struct {
	Items []ReminderRowResult
}

func (a *Adapter) ListLeadReminders(ctx context.Context, leadId string) (*ListLeadRemindersResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListLeadReminders"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListLeadReminders(ctx, &pb.CrmDepthListLeadRemindersRequest{LeadId: leadId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]ReminderRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, ReminderRowResult{Id: r.GetId(), LeadId: r.GetLeadId(), RemindAt: r.GetRemindAt(), Message: r.GetMessage(), Status: r.GetStatus()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListLeadRemindersResult{Items: items}, nil
}

type FilterBotLeadsAdapterResult struct {
	BotLeadIds    []string
	FilteredCount int32
}

func (a *Adapter) FilterBotLeads(ctx context.Context, agentId string) (*FilterBotLeadsAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.FilterBotLeads"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.FilterBotLeads(ctx, &pb.CrmDepthFilterBotLeadsRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &FilterBotLeadsAdapterResult{BotLeadIds: resp.GetBotLeadIds(), FilteredCount: resp.GetFilteredCount()}, nil
}

type CreateDripSequenceParams struct {
	AgentId string
	Name    string
	Steps   []string
}

type DripSequenceAdapterResult struct {
	SequenceId string
	Status     string
}

func (a *Adapter) CreateDripSequence(ctx context.Context, params *CreateDripSequenceParams) (*DripSequenceAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateDripSequence"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateDripSequence(ctx, &pb.CrmDepthCreateDripSequenceRequest{AgentId: params.AgentId, Name: params.Name, Steps: params.Steps})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &DripSequenceAdapterResult{SequenceId: resp.GetSequenceId(), Status: resp.GetStatus()}, nil
}

type ListDripSequencesParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type DripSequenceRowResult struct {
	Id        string
	AgentId   string
	Name      string
	StepCount int32
	CreatedAt string
}

type ListDripSequencesResult struct {
	Items    []DripSequenceRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListDripSequences(ctx context.Context, params *ListDripSequencesParams) (*ListDripSequencesResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListDripSequences"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListDripSequences(ctx, &pb.CrmDepthListDripSequencesRequest{AgentId: params.AgentId, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]DripSequenceRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, DripSequenceRowResult{Id: r.GetId(), AgentId: r.GetAgentId(), Name: r.GetName(), StepCount: r.GetStepCount(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListDripSequencesResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type CreateMomentTriggerParams struct {
	AgentId   string
	EventType string
	Action    string
	Condition string
}

type MomentTriggerAdapterResult struct {
	TriggerId string
	Status    string
}

func (a *Adapter) CreateMomentTrigger(ctx context.Context, params *CreateMomentTriggerParams) (*MomentTriggerAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateMomentTrigger"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateMomentTrigger(ctx, &pb.CrmDepthCreateMomentTriggerRequest{AgentId: params.AgentId, EventType: params.EventType, Action: params.Action, Condition: params.Condition})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &MomentTriggerAdapterResult{TriggerId: resp.GetTriggerId(), Status: resp.GetStatus()}, nil
}

type CreateSegmentParams struct {
	AgentId  string
	Name     string
	Criteria []string
}

type SegmentAdapterResult struct {
	SegmentId string
	Status    string
}

func (a *Adapter) CreateSegment(ctx context.Context, params *CreateSegmentParams) (*SegmentAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateSegment"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateSegment(ctx, &pb.CrmDepthCreateSegmentRequest{AgentId: params.AgentId, Name: params.Name, Criteria: params.Criteria})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SegmentAdapterResult{SegmentId: resp.GetSegmentId(), Status: resp.GetStatus()}, nil
}

type ListSegmentsParams struct {
	AgentId  string
	Page     int32
	PageSize int32
}

type SegmentRowResult struct {
	Id        string
	AgentId   string
	Name      string
	LeadCount int32
	CreatedAt string
}

type ListSegmentsResult struct {
	Items    []SegmentRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListSegments(ctx context.Context, params *ListSegmentsParams) (*ListSegmentsResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListSegments"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListSegments(ctx, &pb.CrmDepthListSegmentsRequest{AgentId: params.AgentId, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]SegmentRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, SegmentRowResult{Id: r.GetId(), AgentId: r.GetAgentId(), Name: r.GetName(), LeadCount: r.GetLeadCount(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListSegmentsResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type SendBroadcastParams struct {
	AgentId    string
	SegmentIds []string
	Message    string
	Channel    string
}

type SendBroadcastAdapterResult struct {
	BroadcastId    string
	RecipientCount int32
	Status         string
}

func (a *Adapter) SendBroadcast(ctx context.Context, params *SendBroadcastParams) (*SendBroadcastAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.SendBroadcast"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.SendBroadcast(ctx, &pb.CrmDepthSendBroadcastRequest{AgentId: params.AgentId, SegmentIds: params.SegmentIds, Message: params.Message, Channel: params.Channel})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SendBroadcastAdapterResult{BroadcastId: resp.GetBroadcastId(), RecipientCount: resp.GetRecipientCount(), Status: resp.GetStatus()}, nil
}

type AssignLeadFairParams struct {
	AgentId string
	LeadIds []string
}

type AssignLeadFairAdapterResult struct {
	AssignedCount int32
	Status        string
}

func (a *Adapter) AssignLeadFair(ctx context.Context, params *AssignLeadFairParams) (*AssignLeadFairAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.AssignLeadFair"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.AssignLeadFair(ctx, &pb.CrmDepthAssignLeadFairRequest{AgentId: params.AgentId, LeadIds: params.LeadIds})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AssignLeadFairAdapterResult{AssignedCount: resp.GetAssignedCount(), Status: resp.GetStatus()}, nil
}

type CreateSlaRuleParams struct {
	AgentId          string
	RuleName         string
	ResponseTimeMins int32
	EscalationLevel  int32
}

type SlaRuleAdapterResult struct {
	RuleId    string
	Status    string
	CreatedAt string
}

func (a *Adapter) CreateSlaRule(ctx context.Context, params *CreateSlaRuleParams) (*SlaRuleAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateSlaRule"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CreateSlaRule(ctx, &pb.CrmDepthCreateSlaRuleRequest{AgentId: params.AgentId, RuleName: params.RuleName, ResponseTimeMins: params.ResponseTimeMins, EscalationLevel: params.EscalationLevel})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SlaRuleAdapterResult{RuleId: resp.GetRuleId(), Status: resp.GetStatus(), CreatedAt: resp.GetCreatedAt()}, nil
}

type LeadTrailRowResult struct {
	EventType string
	Notes     string
	ActorId   string
	CreatedAt string
}

type GetLeadTrailAdapterResult struct {
	LeadId string
	Events []LeadTrailRowResult
}

func (a *Adapter) GetLeadTrail(ctx context.Context, leadId string) (*GetLeadTrailAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetLeadTrail"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetLeadTrail(ctx, &pb.CrmDepthGetLeadTrailRequest{LeadId: leadId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	events := make([]LeadTrailRowResult, 0, len(resp.GetEvents()))
	for _, r := range resp.GetEvents() {
		events = append(events, LeadTrailRowResult{EventType: r.GetEventType(), Notes: r.GetNotes(), ActorId: r.GetActorId(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetLeadTrailAdapterResult{LeadId: resp.GetLeadId(), Events: events}, nil
}

type TagLeadParams struct {
	LeadId string
	Tags   []string
}

type TagLeadAdapterResult struct {
	LeadId    string
	Tags      []string
	UpdatedAt string
}

func (a *Adapter) TagLead(ctx context.Context, params *TagLeadParams) (*TagLeadAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.TagLead"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.TagLead(ctx, &pb.CrmDepthTagLeadRequest{LeadId: params.LeadId, Tags: params.Tags})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &TagLeadAdapterResult{LeadId: resp.GetLeadId(), Tags: resp.GetTags(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

type GenerateQuoteParams struct {
	LeadId    string
	PackageId string
	Pax       int32
	Discount  float64
}

type QuoteAdapterResult struct {
	QuoteId    string
	LeadId     string
	PackageId  string
	TotalPrice int64
	Status     string
	CreatedAt  string
	ExpiresAt  string
}

func (a *Adapter) GenerateQuote(ctx context.Context, params *GenerateQuoteParams) (*QuoteAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GenerateQuote"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GenerateQuote(ctx, &pb.CrmDepthGenerateQuoteRequest{LeadId: params.LeadId, PackageId: params.PackageId, Pax: params.Pax, Discount: params.Discount})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &QuoteAdapterResult{QuoteId: resp.GetQuoteId(), LeadId: resp.GetLeadId(), PackageId: resp.GetPackageId(), TotalPrice: resp.GetTotalPrice(), Status: resp.GetStatus(), CreatedAt: resp.GetCreatedAt(), ExpiresAt: resp.GetExpiresAt()}, nil
}

func (a *Adapter) GetQuote(ctx context.Context, quoteId string) (*QuoteAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetQuote"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetQuote(ctx, &pb.CrmDepthGetQuoteRequest{QuoteId: quoteId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &QuoteAdapterResult{QuoteId: resp.GetQuoteId(), LeadId: resp.GetLeadId(), PackageId: resp.GetPackageId(), TotalPrice: resp.GetTotalPrice(), Status: resp.GetStatus(), CreatedAt: resp.GetCreatedAt(), ExpiresAt: resp.GetExpiresAt()}, nil
}

type BuildPaymentLinkParams struct {
	QuoteId string
	LeadId  string
	Notes   string
}

type PaymentLinkAdapterResult struct {
	LinkId     string
	PaymentUrl string
	ExpiresAt  string
}

func (a *Adapter) BuildPaymentLink(ctx context.Context, params *BuildPaymentLinkParams) (*PaymentLinkAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.BuildPaymentLink"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.BuildPaymentLink(ctx, &pb.CrmDepthBuildPaymentLinkRequest{QuoteId: params.QuoteId, LeadId: params.LeadId, Notes: params.Notes})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &PaymentLinkAdapterResult{LinkId: resp.GetLinkId(), PaymentUrl: resp.GetPaymentUrl(), ExpiresAt: resp.GetExpiresAt()}, nil
}

type RequestDiscountParams struct {
	QuoteId string
	AgentId string
	Amount  float64
	Reason  string
}

type DiscountAdapterResult struct {
	DiscountId string
	QuoteId    string
	Amount     float64
	Status     string
	CreatedAt  string
}

func (a *Adapter) RequestDiscount(ctx context.Context, params *RequestDiscountParams) (*DiscountAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.RequestDiscount"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.RequestDiscount(ctx, &pb.CrmDepthRequestDiscountRequest{QuoteId: params.QuoteId, AgentId: params.AgentId, Amount: params.Amount, Reason: params.Reason})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &DiscountAdapterResult{DiscountId: resp.GetDiscountId(), QuoteId: resp.GetQuoteId(), Amount: resp.GetAmount(), Status: resp.GetStatus(), CreatedAt: resp.GetCreatedAt()}, nil
}

type ApproveDiscountParams struct {
	DiscountId string
	Decision   string
	Notes      string
}

func (a *Adapter) ApproveDiscount(ctx context.Context, params *ApproveDiscountParams) (*DiscountAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.ApproveDiscount"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ApproveDiscount(ctx, &pb.CrmDepthApproveDiscountRequest{DiscountId: params.DiscountId, Decision: params.Decision, Notes: params.Notes})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &DiscountAdapterResult{DiscountId: resp.GetDiscountId(), QuoteId: resp.GetQuoteId(), Amount: resp.GetAmount(), Status: resp.GetStatus(), CreatedAt: resp.GetCreatedAt()}, nil
}

type GetStaleProspectsParams struct {
	AgentId   string
	StaleDays int32
	Page      int32
	PageSize  int32
}

type GetStaleProspectsResult struct {
	Items    []AgentLeadAdapterResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) GetStaleProspects(ctx context.Context, params *GetStaleProspectsParams) (*GetStaleProspectsResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetStaleProspects"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetStaleProspects(ctx, &pb.CrmDepthStaleProspectsRequest{AgentId: params.AgentId, StaleDays: params.StaleDays, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]AgentLeadAdapterResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, AgentLeadAdapterResult{Id: r.GetId(), AgentId: r.GetAgentId(), Name: r.GetName(), Phone: r.GetPhone(), Email: r.GetEmail(), Status: r.GetStatus(), Source: r.GetSource(), CreatedAt: r.GetCreatedAt(), UpdatedAt: r.GetUpdatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetStaleProspectsResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

// ---------------------------------------------------------------------------
// Commission
// ---------------------------------------------------------------------------

type CommissionBalanceAdapterResult struct {
	AgentId        string
	BalanceIdr     int64
	PendingIdr     int64
	TotalEarnedIdr int64
}

func (a *Adapter) GetCommissionBalance(ctx context.Context, agentId string) (*CommissionBalanceAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetCommissionBalance"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetCommissionBalance(ctx, &pb.CrmDepthGetCommissionRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CommissionBalanceAdapterResult{AgentId: resp.GetAgentId(), BalanceIdr: resp.GetBalanceIdr(), PendingIdr: resp.GetPendingIdr(), TotalEarnedIdr: resp.GetTotalEarnedIdr()}, nil
}

type CommissionEventRowResult struct {
	EventId     string
	EventType   string
	AmountIdr   int64
	ReferenceId string
	CreatedAt   string
}

type GetCommissionEventsAdapterResult struct {
	AgentId string
	Events  []CommissionEventRowResult
	Total   int64
}

func (a *Adapter) GetCommissionEvents(ctx context.Context, agentId string) (*GetCommissionEventsAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetCommissionEvents"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetCommissionEvents(ctx, &pb.CrmDepthGetCommissionRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	events := make([]CommissionEventRowResult, 0, len(resp.GetEvents()))
	for _, r := range resp.GetEvents() {
		events = append(events, CommissionEventRowResult{EventId: r.GetEventId(), EventType: r.GetEventType(), AmountIdr: r.GetAmountIdr(), ReferenceId: r.GetReferenceId(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetCommissionEventsAdapterResult{AgentId: resp.GetAgentId(), Events: events, Total: resp.GetTotal()}, nil
}

type RequestPayoutParams struct {
	AgentId   string
	AmountIdr int64
	BankCode  string
	AccountNo string
}

type PayoutAdapterResult struct {
	PayoutId  string
	Status    string
	CreatedAt string
}

func (a *Adapter) RequestPayout(ctx context.Context, params *RequestPayoutParams) (*PayoutAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.RequestPayout"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.RequestPayout(ctx, &pb.CrmDepthRequestPayoutRequest{AgentId: params.AgentId, AmountIdr: params.AmountIdr, BankCode: params.BankCode, AccountNo: params.AccountNo})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &PayoutAdapterResult{PayoutId: resp.GetPayoutId(), Status: resp.GetStatus(), CreatedAt: resp.GetCreatedAt()}, nil
}

type PayoutHistoryRowResult struct {
	PayoutId  string
	AmountIdr int64
	Status    string
	CreatedAt string
}

type GetPayoutHistoryAdapterResult struct {
	AgentId string
	Items   []PayoutHistoryRowResult
	Total   int64
}

func (a *Adapter) GetPayoutHistory(ctx context.Context, agentId string) (*GetPayoutHistoryAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetPayoutHistory"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetPayoutHistory(ctx, &pb.CrmDepthGetCommissionRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]PayoutHistoryRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, PayoutHistoryRowResult{PayoutId: r.GetPayoutId(), AmountIdr: r.GetAmountIdr(), Status: r.GetStatus(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetPayoutHistoryAdapterResult{AgentId: resp.GetAgentId(), Items: items, Total: resp.GetTotal()}, nil
}

type ComputeOverrideParams struct {
	AgentId    string
	DownlineId string
	SaleAmount int64
}

type OverrideCommissionAdapterResult struct {
	OverrideAmount int64
	Rate           float64
	Tier           string
}

func (a *Adapter) ComputeOverrideCommission(ctx context.Context, params *ComputeOverrideParams) (*OverrideCommissionAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.ComputeOverrideCommission"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ComputeOverrideCommission(ctx, &pb.CrmDepthComputeOverrideRequest{AgentId: params.AgentId, DownlineId: params.DownlineId, SaleAmount: params.SaleAmount})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &OverrideCommissionAdapterResult{OverrideAmount: resp.GetOverrideAmount(), Rate: resp.GetRate(), Tier: resp.GetTier()}, nil
}

type GetRoasReportParams struct {
	AgentId  string
	DateFrom string
	DateTo   string
}

type RoasReportAdapterResult struct {
	AgentId        string
	TotalAdSpend   float64
	TotalRevenue   int64
	Roas           float64
	ConversionRate float64
}

func (a *Adapter) GetRoasReport(ctx context.Context, params *GetRoasReportParams) (*RoasReportAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetRoasReport"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetRoasReport(ctx, &pb.CrmDepthRoasReportRequest{AgentId: params.AgentId, DateFrom: params.DateFrom, DateTo: params.DateTo})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RoasReportAdapterResult{AgentId: resp.GetAgentId(), TotalAdSpend: resp.GetTotalAdSpend(), TotalRevenue: resp.GetTotalRevenue(), Roas: resp.GetRoas(), ConversionRate: resp.GetConversionRate()}, nil
}

// ---------------------------------------------------------------------------
// Academy
// ---------------------------------------------------------------------------

type ListAcademyCoursesParams struct {
	Page     int32
	PageSize int32
}

type CourseRowResult struct {
	Id          string
	Title       string
	Description string
	Category    string
	DurationMin int32
}

type ListAcademyCoursesResult struct {
	Items    []CourseRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListAcademyCourses(ctx context.Context, params *ListAcademyCoursesParams) (*ListAcademyCoursesResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListAcademyCourses"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListAcademyCourses(ctx, &pb.CrmDepthListCoursesRequest{Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]CourseRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, CourseRowResult{Id: r.GetId(), Title: r.GetTitle(), Description: r.GetDescription(), Category: r.GetCategory(), DurationMin: r.GetDurationMin()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListAcademyCoursesResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type GetCourseProgressParams struct {
	AgentId  string
	CourseId string
}

type CourseProgressAdapterResult struct {
	AgentId         string
	CourseId        string
	ProgressPercent float64
	CompletedAt     string
	Status          string
}

func (a *Adapter) GetCourseProgress(ctx context.Context, params *GetCourseProgressParams) (*CourseProgressAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetCourseProgress"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetCourseProgress(ctx, &pb.CrmDepthGetCourseProgressRequest{AgentId: params.AgentId, CourseId: params.CourseId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CourseProgressAdapterResult{AgentId: resp.GetAgentId(), CourseId: resp.GetCourseId(), ProgressPercent: resp.GetProgressPercent(), CompletedAt: resp.GetCompletedAt(), Status: resp.GetStatus()}, nil
}

type SubmitQuizParams struct {
	AgentId  string
	CourseId string
	Answers  []string
}

type QuizAdapterResult struct {
	AgentId  string
	CourseId string
	Score    float64
	Passed   bool
	GradedAt string
}

func (a *Adapter) SubmitQuiz(ctx context.Context, params *SubmitQuizParams) (*QuizAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.SubmitQuiz"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.SubmitQuiz(ctx, &pb.CrmDepthSubmitQuizRequest{AgentId: params.AgentId, CourseId: params.CourseId, Answers: params.Answers})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &QuizAdapterResult{AgentId: resp.GetAgentId(), CourseId: resp.GetCourseId(), Score: resp.GetScore(), Passed: resp.GetPassed(), GradedAt: resp.GetGradedAt()}, nil
}

type ListSalesScriptsParams struct {
	Category string
	Page     int32
	PageSize int32
}

type SalesScriptRowResult struct {
	Id       string
	Title    string
	Category string
	Content  string
}

type ListSalesScriptsResult struct {
	Items    []SalesScriptRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListSalesScripts(ctx context.Context, params *ListSalesScriptsParams) (*ListSalesScriptsResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListSalesScripts"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.ListSalesScripts(ctx, &pb.CrmDepthListSalesScriptsRequest{Category: params.Category, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]SalesScriptRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, SalesScriptRowResult{Id: r.GetId(), Title: r.GetTitle(), Category: r.GetCategory(), Content: r.GetContent()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListSalesScriptsResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type GetLeaderboardParams struct {
	Period   string
	Page     int32
	PageSize int32
}

type LeaderboardRowResult struct {
	Rank      int32
	AgentId   string
	AgentName string
	Sales     int64
	Points    int64
}

type LeaderboardAdapterResult struct {
	Items    []LeaderboardRowResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) GetLeaderboard(ctx context.Context, params *GetLeaderboardParams) (*LeaderboardAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetLeaderboard"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetLeaderboard(ctx, &pb.CrmDepthLeaderboardRequest{Period: params.Period, Page: params.Page, PageSize: params.PageSize})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]LeaderboardRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, LeaderboardRowResult{Rank: r.GetRank(), AgentId: r.GetAgentId(), AgentName: r.GetAgentName(), Sales: r.GetSales(), Points: r.GetPoints()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &LeaderboardAdapterResult{Items: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

type AgentTierAdapterResult struct {
	AgentId   string
	TierLevel string
	TierName  string
	Points    int64
	NextTier  string
}

func (a *Adapter) GetAgentTier(ctx context.Context, agentId string) (*AgentTierAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetAgentTier"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetAgentTier(ctx, &pb.CrmDepthGetAgentRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AgentTierAdapterResult{AgentId: resp.GetAgentId(), TierLevel: resp.GetTierLevel(), TierName: resp.GetTierName(), Points: resp.GetPoints(), NextTier: resp.GetNextTier()}, nil
}

// ---------------------------------------------------------------------------
// Alumni/Other
// ---------------------------------------------------------------------------

type ReferralRowResult struct {
	ReferralId  string
	RefereeName string
	Status      string
	Reward      int64
	CreatedAt   string
}

type GetAlumniReferralsAdapterResult struct {
	AgentId string
	Items   []ReferralRowResult
	Total   int64
}

func (a *Adapter) GetAlumniReferrals(ctx context.Context, agentId string) (*GetAlumniReferralsAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetAlumniReferrals"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetAlumniReferrals(ctx, &pb.CrmDepthGetAgentRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]ReferralRowResult, 0, len(resp.GetItems()))
	for _, r := range resp.GetItems() {
		items = append(items, ReferralRowResult{ReferralId: r.GetReferralId(), RefereeName: r.GetRefereeName(), Status: r.GetStatus(), Reward: r.GetReward(), CreatedAt: r.GetCreatedAt()})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetAlumniReferralsAdapterResult{AgentId: resp.GetAgentId(), Items: items, Total: resp.GetTotal()}, nil
}

type ReturnIntentAdapterResult struct {
	AgentId          string
	EligibleLeads    int32
	PotentialSavings int64
}

func (a *Adapter) GetReturnIntentSavings(ctx context.Context, agentId string) (*ReturnIntentAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetReturnIntentSavings"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.GetReturnIntentSavings(ctx, &pb.CrmDepthGetAgentRequest{AgentId: agentId})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ReturnIntentAdapterResult{AgentId: resp.GetAgentId(), EligibleLeads: resp.GetEligibleLeads(), PotentialSavings: resp.GetPotentialSavings()}, nil
}

type CalculateZakatParams struct {
	IncomeSources    []string
	TotalAssets      int64
	TotalLiabilities int64
}

type ZakatAdapterResult struct {
	ZakatAmount int64
	NisabValue  int64
	Rate        float64
	Eligible    bool
}

func (a *Adapter) CalculateZakat(ctx context.Context, params *CalculateZakatParams) (*ZakatAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.CalculateZakat"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.CalculateZakat(ctx, &pb.CrmDepthCalculateZakatRequest{IncomeSources: params.IncomeSources, TotalAssets: params.TotalAssets, TotalLiabilities: params.TotalLiabilities})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ZakatAdapterResult{ZakatAmount: resp.GetZakatAmount(), NisabValue: resp.GetNisabValue(), Rate: resp.GetRate(), Eligible: resp.GetEligible()}, nil
}

type RecordCharityParams struct {
	AgentId    string
	AmountIdr  int64
	CharityOrg string
	Notes      string
}

type CharityAdapterResult struct {
	CharityId  string
	Status     string
	RecordedAt string
}

func (a *Adapter) RecordCharity(ctx context.Context, params *RecordCharityParams) (*CharityAdapterResult, error) {
	const op = "crm_grpc_adapter.Adapter.RecordCharity"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")
	resp, err := a.depthClient.RecordCharity(ctx, &pb.CrmDepthRecordCharityRequest{AgentId: params.AgentId, AmountIdr: params.AmountIdr, CharityOrg: params.CharityOrg, Notes: params.Notes})
	if err != nil {
		wrapped := mapCrmDepthError(err)
		span.RecordError(wrapped); span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CharityAdapterResult{CharityId: resp.GetCharityId(), Status: resp.GetStatus(), RecordedAt: resp.GetRecordedAt()}, nil
}
