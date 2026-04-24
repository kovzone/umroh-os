// crm_depth.go — gRPC handler methods for Wave 7 CRM depth RPCs.
// Delegates to service.IService for all business logic.

package grpc_api

import (
	"context"

	"crm-svc/api/grpc_api/pb"
	"crm-svc/service"
)

// ---------------------------------------------------------------------------
// Agency Registration
// ---------------------------------------------------------------------------

func (s *Server) RegisterAgent(ctx context.Context, req *pb.AgentRegisterRequest) (*pb.AgentResult, error) {
	res, err := s.svc.RegisterAgent(ctx, &service.RegisterAgentParams{
		AgencyName: req.GetAgencyName(),
		OwnerName:  req.GetOwnerName(),
		Phone:      req.GetPhone(),
		Email:      req.GetEmail(),
		City:       req.GetCity(),
		Province:   req.GetProvince(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AgentResult{
		AgentId:    res.AgentId,
		AgencyName: res.AgencyName,
		OwnerName:  res.OwnerName,
		Phone:      res.Phone,
		Email:      res.Email,
		Status:     res.Status,
		Tier:       res.Tier,
		CreatedAt:  res.CreatedAt,
	}, nil
}

func (s *Server) SubmitAgentKyc(ctx context.Context, req *pb.AgentKycRequest) (*pb.AgentResult, error) {
	res, err := s.svc.SubmitAgentKyc(ctx, &service.SubmitAgentKycParams{
		AgentId:   req.GetAgentId(),
		KtpUrl:    req.GetKtpUrl(),
		NpwpUrl:   req.GetNpwpUrl(),
		SiupUrl:   req.GetSiupUrl(),
		SelfieUrl: req.GetSelfieUrl(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AgentResult{
		AgentId: res.AgentId, AgencyName: res.AgencyName, OwnerName: res.OwnerName,
		Phone: res.Phone, Email: res.Email, Status: res.Status, Tier: res.Tier, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) SignAgentMoU(ctx context.Context, req *pb.AgentMoURequest) (*pb.AgentResult, error) {
	res, err := s.svc.SignAgentMoU(ctx, &service.SignAgentMoUParams{
		AgentId:      req.GetAgentId(),
		SignatureUrl: req.GetSignatureUrl(),
		SignedAt:     req.GetSignedAt(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AgentResult{
		AgentId: res.AgentId, AgencyName: res.AgencyName, OwnerName: res.OwnerName,
		Phone: res.Phone, Email: res.Email, Status: res.Status, Tier: res.Tier, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) GetAgentProfile(ctx context.Context, req *pb.AgentIdRequest) (*pb.AgentResult, error) {
	res, err := s.svc.GetAgentProfile(ctx, &service.AgentIdParams{AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	return &pb.AgentResult{
		AgentId: res.AgentId, AgencyName: res.AgencyName, OwnerName: res.OwnerName,
		Phone: res.Phone, Email: res.Email, Status: res.Status, Tier: res.Tier, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) GetReplicaSite(ctx context.Context, req *pb.AgentIdRequest) (*pb.ReplicaSiteResult, error) {
	res, err := s.svc.GetReplicaSite(ctx, &service.AgentIdParams{AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	return &pb.ReplicaSiteResult{
		AgentId: res.AgentId, SiteUrl: res.SiteUrl, Theme: res.Theme, LogoUrl: res.LogoUrl, UpdatedAt: res.UpdatedAt,
	}, nil
}

func (s *Server) UpdateReplicaSite(ctx context.Context, req *pb.UpdateReplicaSiteRequest) (*pb.ReplicaSiteResult, error) {
	res, err := s.svc.UpdateReplicaSite(ctx, &service.UpdateReplicaSiteParams{
		AgentId:      req.GetAgentId(),
		Theme:        req.GetTheme(),
		LogoUrl:      req.GetLogoUrl(),
		CustomDomain: req.GetCustomDomain(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplicaSiteResult{
		AgentId: res.AgentId, SiteUrl: res.SiteUrl, Theme: res.Theme, LogoUrl: res.LogoUrl, UpdatedAt: res.UpdatedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// Content & Marketing
// ---------------------------------------------------------------------------

func (s *Server) GetSocialShareLink(ctx context.Context, req *pb.SocialShareRequest) (*pb.SocialShareResult, error) {
	res, err := s.svc.GetSocialShareLink(ctx, &service.SocialShareParams{
		AgentId: req.GetAgentId(), PackageId: req.GetPackageId(), Platform: req.GetPlatform(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.SocialShareResult{ShareUrl: res.ShareUrl, ShortCode: res.ShortCode}, nil
}

func (s *Server) GenerateBusinessCard(ctx context.Context, req *pb.BusinessCardRequest) (*pb.BusinessCardResult, error) {
	res, err := s.svc.GenerateBusinessCard(ctx, &service.BusinessCardParams{
		AgentId: req.GetAgentId(), Template: req.GetTemplate(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.BusinessCardResult{CardUrl: res.CardUrl, CreatedAt: res.CreatedAt}, nil
}

func (s *Server) ListContentBank(ctx context.Context, req *pb.ContentBankListRequest) (*pb.ContentBankListResult, error) {
	res, err := s.svc.ListContentBank(ctx, &service.ContentBankListParams{
		AgentId: req.GetAgentId(), Category: req.GetCategory(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*pb.ContentAssetItem, len(res.Assets))
	for i, a := range res.Assets {
		items[i] = &pb.ContentAssetItem{AssetId: a.AssetId, Title: a.Title, Type: a.Type, Url: a.Url, CreatedAt: a.CreatedAt}
	}
	return &pb.ContentBankListResult{Assets: items, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) CreateContentAsset(ctx context.Context, req *pb.CreateContentAssetRequest) (*pb.ContentAssetResult, error) {
	res, err := s.svc.CreateContentAsset(ctx, &service.CreateContentAssetParams{
		AgentId: req.GetAgentId(), Title: req.GetTitle(), Type: req.GetType(), Url: req.GetUrl(), Category: req.GetCategory(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ContentAssetResult{
		AssetId: res.AssetId, AgentId: res.AgentId, Title: res.Title, Type: res.Type, Url: res.Url, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) WatermarkFlyer(ctx context.Context, req *pb.WatermarkFlyerRequest) (*pb.WatermarkFlyerResult, error) {
	res, err := s.svc.WatermarkFlyer(ctx, &service.WatermarkFlyerParams{
		AgentId: req.GetAgentId(), FlyerUrl: req.GetFlyerUrl(), Text: req.GetText(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.WatermarkFlyerResult{ResultUrl: res.ResultUrl}, nil
}

func (s *Server) ListProgramGallery(ctx context.Context, req *pb.ProgramGalleryRequest) (*pb.ProgramGalleryResult, error) {
	res, err := s.svc.ListProgramGallery(ctx, &service.ProgramGalleryParams{
		PackageId: req.GetPackageId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*pb.GalleryItem, len(res.Items))
	for i, it := range res.Items {
		items[i] = &pb.GalleryItem{ImageUrl: it.ImageUrl, Caption: it.Caption}
	}
	return &pb.ProgramGalleryResult{PackageId: res.PackageId, Items: items, Total: res.Total}, nil
}

func (s *Server) SetTrackingCode(ctx context.Context, req *pb.TrackingCodeRequest) (*pb.TrackingCodeResult, error) {
	res, err := s.svc.SetTrackingCode(ctx, &service.TrackingCodeParams{
		AgentId: req.GetAgentId(), Code: req.GetCode(), Type: req.GetType(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.TrackingCodeResult{AgentId: res.AgentId, Code: res.Code, Type: res.Type, UpdatedAt: res.UpdatedAt}, nil
}

func (s *Server) GetAdsManagerStats(ctx context.Context, req *pb.AdsManagerRequest) (*pb.AdsManagerResult, error) {
	res, err := s.svc.GetAdsManagerStats(ctx, &service.AdsManagerParams{
		AgentId: req.GetAgentId(), DateFrom: req.GetDateFrom(), DateTo: req.GetDateTo(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AdsManagerResult{
		AgentId: res.AgentId, Impressions: res.Impressions, Clicks: res.Clicks, Conversions: res.Conversions, Spend: res.Spend,
	}, nil
}

func (s *Server) CreateUtmLink(ctx context.Context, req *pb.CreateUtmLinkRequest) (*pb.UtmLinkResult, error) {
	res, err := s.svc.CreateUtmLink(ctx, &service.CreateUtmLinkParams{
		AgentId: req.GetAgentId(), BaseUrl: req.GetBaseUrl(), UtmSource: req.GetUtmSource(),
		UtmMedium: req.GetUtmMedium(), UtmCampaign: req.GetUtmCampaign(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.UtmLinkResult{
		LinkId: res.LinkId, AgentId: res.AgentId, FullUrl: res.FullUrl, ShortUrl: res.ShortUrl, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) ListUtmLinks(ctx context.Context, req *pb.ListUtmLinksRequest) (*pb.UtmLinkListResult, error) {
	res, err := s.svc.ListUtmLinks(ctx, &service.ListUtmLinksParams{
		AgentId: req.GetAgentId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	links := make([]*pb.UtmLinkResult, len(res.Links))
	for i, l := range res.Links {
		links[i] = &pb.UtmLinkResult{LinkId: l.LinkId, AgentId: l.AgentId, FullUrl: l.FullUrl, ShortUrl: l.ShortUrl, CreatedAt: l.CreatedAt}
	}
	return &pb.UtmLinkListResult{Links: links, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) CreateLandingPage(ctx context.Context, req *pb.CreateLandingPageRequest) (*pb.LandingPageResult, error) {
	res, err := s.svc.CreateLandingPage(ctx, &service.CreateLandingPageParams{
		AgentId: req.GetAgentId(), Title: req.GetTitle(), PackageId: req.GetPackageId(), Template: req.GetTemplate(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.LandingPageResult{
		PageId: res.PageId, AgentId: res.AgentId, Title: res.Title, Url: res.Url, Status: res.Status, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) ListLandingPages(ctx context.Context, req *pb.ListLandingPagesRequest) (*pb.LandingPageListResult, error) {
	res, err := s.svc.ListLandingPages(ctx, &service.ListLandingPagesParams{
		AgentId: req.GetAgentId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	pages := make([]*pb.LandingPageResult, len(res.Pages))
	for i, p := range res.Pages {
		pages[i] = &pb.LandingPageResult{PageId: p.PageId, AgentId: p.AgentId, Title: p.Title, Url: p.Url, Status: p.Status, CreatedAt: p.CreatedAt}
	}
	return &pb.LandingPageListResult{Pages: pages, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) ScheduleContent(ctx context.Context, req *pb.ScheduleContentRequest) (*pb.ScheduledContentResult, error) {
	res, err := s.svc.ScheduleContent(ctx, &service.ScheduleContentParams{
		AgentId: req.GetAgentId(), AssetId: req.GetAssetId(), Platform: req.GetPlatform(),
		ScheduledAt: req.GetScheduledAt(), Caption: req.GetCaption(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ScheduledContentResult{
		ScheduleId: res.ScheduleId, AgentId: res.AgentId, AssetId: res.AssetId,
		Platform: res.Platform, ScheduledAt: res.ScheduledAt, Status: res.Status,
	}, nil
}

func (s *Server) ListScheduledContent(ctx context.Context, req *pb.ListScheduledContentRequest) (*pb.ScheduledContentListResult, error) {
	res, err := s.svc.ListScheduledContent(ctx, &service.ListScheduledContentParams{
		AgentId: req.GetAgentId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*pb.ScheduledContentResult, len(res.Items))
	for i, it := range res.Items {
		items[i] = &pb.ScheduledContentResult{
			ScheduleId: it.ScheduleId, AgentId: it.AgentId, AssetId: it.AssetId,
			Platform: it.Platform, ScheduledAt: it.ScheduledAt, Status: it.Status,
		}
	}
	return &pb.ScheduledContentListResult{Items: items, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) GetContentAnalytics(ctx context.Context, req *pb.ContentAnalyticsRequest) (*pb.ContentAnalyticsResult, error) {
	res, err := s.svc.GetContentAnalytics(ctx, &service.ContentAnalyticsParams{
		AgentId: req.GetAgentId(), AssetId: req.GetAssetId(), DateFrom: req.GetDateFrom(), DateTo: req.GetDateTo(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ContentAnalyticsResult{
		AgentId: res.AgentId, Views: res.Views, Likes: res.Likes, Shares: res.Shares, Leads: res.Leads, Ctr: res.Ctr,
	}, nil
}

// ---------------------------------------------------------------------------
// CRM / Leads
// ---------------------------------------------------------------------------

func (s *Server) CreateAgentLead(ctx context.Context, req *pb.CreateAgentLeadRequest) (*pb.AgentLeadResult, error) {
	res, err := s.svc.CreateAgentLead(ctx, &service.CreateAgentLeadParams{
		AgentId: req.GetAgentId(), Name: req.GetName(), Phone: req.GetPhone(),
		Email: req.GetEmail(), PackageId: req.GetPackageId(), Source: req.GetSource(), Notes: req.GetNotes(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AgentLeadResult{
		LeadId: res.LeadId, AgentId: res.AgentId, Name: res.Name, Phone: res.Phone,
		Email: res.Email, Status: res.Status, PackageId: res.PackageId, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) ListAgentLeads(ctx context.Context, req *pb.ListAgentLeadsRequest) (*pb.AgentLeadListResult, error) {
	res, err := s.svc.ListAgentLeads(ctx, &service.ListAgentLeadsParams{
		AgentId: req.GetAgentId(), Status: req.GetStatus(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	leads := make([]*pb.AgentLeadResult, len(res.Leads))
	for i, l := range res.Leads {
		leads[i] = &pb.AgentLeadResult{
			LeadId: l.LeadId, AgentId: l.AgentId, Name: l.Name, Phone: l.Phone,
			Email: l.Email, Status: l.Status, PackageId: l.PackageId, CreatedAt: l.CreatedAt,
		}
	}
	return &pb.AgentLeadListResult{Leads: leads, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) SetLeadReminder(ctx context.Context, req *pb.SetLeadReminderRequest) (*pb.LeadReminderResult, error) {
	res, err := s.svc.SetLeadReminder(ctx, &service.SetLeadReminderParams{
		LeadId: req.GetLeadId(), AgentId: req.GetAgentId(), RemindAt: req.GetRemindAt(), Note: req.GetNote(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.LeadReminderResult{
		ReminderId: res.ReminderId, LeadId: res.LeadId, AgentId: res.AgentId,
		RemindAt: res.RemindAt, Note: res.Note, Status: res.Status,
	}, nil
}

func (s *Server) ListLeadReminders(ctx context.Context, req *pb.ListLeadRemindersRequest) (*pb.LeadReminderListResult, error) {
	res, err := s.svc.ListLeadReminders(ctx, &service.ListLeadRemindersParams{
		AgentId: req.GetAgentId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	reminders := make([]*pb.LeadReminderResult, len(res.Reminders))
	for i, r := range res.Reminders {
		reminders[i] = &pb.LeadReminderResult{
			ReminderId: r.ReminderId, LeadId: r.LeadId, AgentId: r.AgentId,
			RemindAt: r.RemindAt, Note: r.Note, Status: r.Status,
		}
	}
	return &pb.LeadReminderListResult{Reminders: reminders, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) FilterBotLeads(ctx context.Context, req *pb.BotFilterRequest) (*pb.BotFilterResult, error) {
	res, err := s.svc.FilterBotLeads(ctx, &service.BotFilterParams{AgentId: req.GetAgentId(), LeadId: req.GetLeadId()})
	if err != nil {
		return nil, err
	}
	return &pb.BotFilterResult{LeadId: res.LeadId, IsBot: res.IsBot, Score: res.Score, Reason: res.Reason}, nil
}

func (s *Server) CreateDripSequence(ctx context.Context, req *pb.CreateDripSequenceRequest) (*pb.DripSequenceResult, error) {
	res, err := s.svc.CreateDripSequence(ctx, &service.CreateDripSequenceParams{
		AgentId: req.GetAgentId(), Name: req.GetName(), Triggers: req.GetTriggers(), Steps: req.GetSteps(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DripSequenceResult{
		SequenceId: res.SequenceId, AgentId: res.AgentId, Name: res.Name, Steps: res.Steps, Status: res.Status, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) ListDripSequences(ctx context.Context, req *pb.ListDripSequencesRequest) (*pb.DripSequenceListResult, error) {
	res, err := s.svc.ListDripSequences(ctx, &service.ListDripSequencesParams{
		AgentId: req.GetAgentId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	seqs := make([]*pb.DripSequenceResult, len(res.Sequences))
	for i, seq := range res.Sequences {
		seqs[i] = &pb.DripSequenceResult{
			SequenceId: seq.SequenceId, AgentId: seq.AgentId, Name: seq.Name, Steps: seq.Steps, Status: seq.Status, CreatedAt: seq.CreatedAt,
		}
	}
	return &pb.DripSequenceListResult{Sequences: seqs, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) CreateMomentTrigger(ctx context.Context, req *pb.CreateMomentTriggerRequest) (*pb.MomentTriggerResult, error) {
	res, err := s.svc.CreateMomentTrigger(ctx, &service.CreateMomentTriggerParams{
		AgentId: req.GetAgentId(), EventType: req.GetEventType(), Action: req.GetAction(), Template: req.GetTemplate(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.MomentTriggerResult{
		TriggerId: res.TriggerId, AgentId: res.AgentId, EventType: res.EventType, Action: res.Action, Status: res.Status, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) CreateSegment(ctx context.Context, req *pb.CreateSegmentRequest) (*pb.SegmentResult, error) {
	res, err := s.svc.CreateSegment(ctx, &service.CreateSegmentParams{
		AgentId: req.GetAgentId(), Name: req.GetName(), Criteria: req.GetCriteria(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.SegmentResult{
		SegmentId: res.SegmentId, AgentId: res.AgentId, Name: res.Name, LeadCount: res.LeadCount, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) ListSegments(ctx context.Context, req *pb.ListSegmentsRequest) (*pb.SegmentListResult, error) {
	res, err := s.svc.ListSegments(ctx, &service.ListSegmentsParams{
		AgentId: req.GetAgentId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	segs := make([]*pb.SegmentResult, len(res.Segments))
	for i, seg := range res.Segments {
		segs[i] = &pb.SegmentResult{
			SegmentId: seg.SegmentId, AgentId: seg.AgentId, Name: seg.Name, LeadCount: seg.LeadCount, CreatedAt: seg.CreatedAt,
		}
	}
	return &pb.SegmentListResult{Segments: segs, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) SendBroadcast(ctx context.Context, req *pb.SendBroadcastRequest) (*pb.BroadcastResult, error) {
	res, err := s.svc.SendBroadcast(ctx, &service.SendBroadcastParams{
		AgentId: req.GetAgentId(), SegmentId: req.GetSegmentId(), Channel: req.GetChannel(), Message: req.GetMessage(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.BroadcastResult{
		BroadcastId: res.BroadcastId, AgentId: res.AgentId, Sent: res.Sent, Failed: res.Failed, Status: res.Status,
	}, nil
}

func (s *Server) AssignLeadFair(ctx context.Context, req *pb.AssignLeadFairRequest) (*pb.AssignLeadFairResult, error) {
	res, err := s.svc.AssignLeadFair(ctx, &service.AssignLeadFairParams{
		AgentId: req.GetAgentId(), LeadIds: req.GetLeadIds(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AssignLeadFairResult{Assigned: res.Assigned, LeadIds: res.LeadIds}, nil
}

func (s *Server) CreateSlaRule(ctx context.Context, req *pb.CreateSlaRuleRequest) (*pb.SlaRuleResult, error) {
	res, err := s.svc.CreateSlaRule(ctx, &service.CreateSlaRuleParams{
		AgentId: req.GetAgentId(), Name: req.GetName(), ResponseMinutes: req.GetResponseMinutes(), EscalateAfter: req.GetEscalateAfter(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.SlaRuleResult{
		RuleId: res.RuleId, AgentId: res.AgentId, Name: res.Name,
		ResponseMinutes: res.ResponseMinutes, EscalateAfter: res.EscalateAfter, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) GetLeadTrail(ctx context.Context, req *pb.GetLeadTrailRequest) (*pb.LeadTrailResult, error) {
	res, err := s.svc.GetLeadTrail(ctx, &service.GetLeadTrailParams{LeadId: req.GetLeadId(), AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	events := make([]*pb.LeadTrailEvent, len(res.Events))
	for i, e := range res.Events {
		events[i] = &pb.LeadTrailEvent{EventType: e.EventType, Actor: e.Actor, Detail: e.Detail, CreatedAt: e.CreatedAt}
	}
	return &pb.LeadTrailResult{LeadId: res.LeadId, Events: events}, nil
}

func (s *Server) TagLead(ctx context.Context, req *pb.TagLeadRequest) (*pb.TagLeadResult, error) {
	res, err := s.svc.TagLead(ctx, &service.TagLeadParams{LeadId: req.GetLeadId(), AgentId: req.GetAgentId(), Tags: req.GetTags()})
	if err != nil {
		return nil, err
	}
	return &pb.TagLeadResult{LeadId: res.LeadId, Tags: res.Tags}, nil
}

func (s *Server) GenerateQuote(ctx context.Context, req *pb.GenerateQuoteRequest) (*pb.QuoteResult, error) {
	res, err := s.svc.GenerateQuote(ctx, &service.GenerateQuoteParams{
		LeadId: req.GetLeadId(), AgentId: req.GetAgentId(), PackageId: req.GetPackageId(), Pax: req.GetPax(), Notes: req.GetNotes(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.QuoteResult{
		QuoteId: res.QuoteId, LeadId: res.LeadId, AgentId: res.AgentId, PackageId: res.PackageId,
		Pax: res.Pax, TotalPrice: res.TotalPrice, PdfUrl: res.PdfUrl, ExpiresAt: res.ExpiresAt, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) GetQuote(ctx context.Context, req *pb.GetQuoteIdRequest) (*pb.QuoteResult, error) {
	res, err := s.svc.GetQuote(ctx, &service.GetQuoteParams{QuoteId: req.GetQuoteId()})
	if err != nil {
		return nil, err
	}
	return &pb.QuoteResult{
		QuoteId: res.QuoteId, LeadId: res.LeadId, AgentId: res.AgentId, PackageId: res.PackageId,
		Pax: res.Pax, TotalPrice: res.TotalPrice, PdfUrl: res.PdfUrl, ExpiresAt: res.ExpiresAt, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) BuildPaymentLink(ctx context.Context, req *pb.BuildPaymentLinkRequest) (*pb.PaymentLinkResult, error) {
	res, err := s.svc.BuildPaymentLink(ctx, &service.BuildPaymentLinkParams{
		QuoteId: req.GetQuoteId(), AgentId: req.GetAgentId(), Amount: req.GetAmount(), ExpiresAt: req.GetExpiresAt(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.PaymentLinkResult{
		LinkId: res.LinkId, QuoteId: res.QuoteId, Url: res.Url, Amount: res.Amount, ExpiresAt: res.ExpiresAt,
	}, nil
}

func (s *Server) RequestDiscount(ctx context.Context, req *pb.RequestDiscountRequest) (*pb.DiscountApprovalResult, error) {
	res, err := s.svc.RequestDiscount(ctx, &service.RequestDiscountParams{
		QuoteId: req.GetQuoteId(), AgentId: req.GetAgentId(), Amount: req.GetAmount(), Reason: req.GetReason(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DiscountApprovalResult{
		DiscountId: res.DiscountId, QuoteId: res.QuoteId, AgentId: res.AgentId, Amount: res.Amount, Status: res.Status, UpdatedAt: res.UpdatedAt,
	}, nil
}

func (s *Server) ApproveDiscount(ctx context.Context, req *pb.ApproveDiscountRequest) (*pb.DiscountApprovalResult, error) {
	res, err := s.svc.ApproveDiscount(ctx, &service.ApproveDiscountParams{
		DiscountId: req.GetDiscountId(), ApproverId: req.GetApproverId(), Approved: req.GetApproved(), Note: req.GetNote(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DiscountApprovalResult{
		DiscountId: res.DiscountId, QuoteId: res.QuoteId, AgentId: res.AgentId, Amount: res.Amount, Status: res.Status, UpdatedAt: res.UpdatedAt,
	}, nil
}

func (s *Server) GetStaleProspects(ctx context.Context, req *pb.StaleProspectRequest) (*pb.StaleProspectListResult, error) {
	res, err := s.svc.GetStaleProspects(ctx, &service.StaleProspectParams{
		AgentId: req.GetAgentId(), StaleDays: req.GetStaleDays(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*pb.StaleProspectItem, len(res.Prospects))
	for i, p := range res.Prospects {
		items[i] = &pb.StaleProspectItem{LeadId: p.LeadId, Name: p.Name, Phone: p.Phone, LastContact: p.LastContact, DaysSince: p.DaysSince}
	}
	return &pb.StaleProspectListResult{Prospects: items, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

// ---------------------------------------------------------------------------
// Commission / Finance
// ---------------------------------------------------------------------------

func (s *Server) GetCommissionBalance(ctx context.Context, req *pb.AgentIdRequest) (*pb.CommissionBalanceResult, error) {
	res, err := s.svc.GetCommissionBalance(ctx, &service.AgentIdParams{AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	return &pb.CommissionBalanceResult{
		AgentId: res.AgentId, Balance: res.Balance, Pending: res.Pending, Withdrawn: res.Withdrawn, UpdatedAt: res.UpdatedAt,
	}, nil
}

func (s *Server) GetCommissionEvents(ctx context.Context, req *pb.CommissionEventListRequest) (*pb.CommissionEventListResult, error) {
	res, err := s.svc.GetCommissionEvents(ctx, &service.CommissionEventListParams{
		AgentId: req.GetAgentId(), DateFrom: req.GetDateFrom(), DateTo: req.GetDateTo(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	events := make([]*pb.CommissionEvent, len(res.Events))
	for i, e := range res.Events {
		events[i] = &pb.CommissionEvent{EventId: e.EventId, Type: e.Type, Amount: e.Amount, BookingId: e.BookingId, CreatedAt: e.CreatedAt}
	}
	return &pb.CommissionEventListResult{Events: events, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) RequestPayout(ctx context.Context, req *pb.PayoutRequest) (*pb.PayoutResult, error) {
	res, err := s.svc.RequestPayout(ctx, &service.PayoutParams{
		AgentId: req.GetAgentId(), Amount: req.GetAmount(), BankAccount: req.GetBankAccount(), BankCode: req.GetBankCode(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.PayoutResult{
		PayoutId: res.PayoutId, AgentId: res.AgentId, Amount: res.Amount, Status: res.Status, CreatedAt: res.CreatedAt,
	}, nil
}

func (s *Server) GetPayoutHistory(ctx context.Context, req *pb.AgentIdRequest) (*pb.PayoutHistoryResult, error) {
	res, err := s.svc.GetPayoutHistory(ctx, &service.AgentIdParams{AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	payouts := make([]*pb.PayoutResult, len(res.Payouts))
	for i, p := range res.Payouts {
		payouts[i] = &pb.PayoutResult{PayoutId: p.PayoutId, AgentId: p.AgentId, Amount: p.Amount, Status: p.Status, CreatedAt: p.CreatedAt}
	}
	return &pb.PayoutHistoryResult{AgentId: res.AgentId, Payouts: payouts, Total: res.Total}, nil
}

func (s *Server) ComputeOverrideCommission(ctx context.Context, req *pb.OverrideCommissionRequest) (*pb.OverrideCommissionResult, error) {
	res, err := s.svc.ComputeOverrideCommission(ctx, &service.OverrideCommissionParams{
		AgentId: req.GetAgentId(), SubAgentId: req.GetSubAgentId(), BookingId: req.GetBookingId(), Rate: req.GetRate(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.OverrideCommissionResult{
		AgentId: res.AgentId, SubAgentId: res.SubAgentId, BookingId: res.BookingId, Amount: res.Amount, Rate: res.Rate,
	}, nil
}

func (s *Server) GetRoasReport(ctx context.Context, req *pb.RoasReportRequest) (*pb.RoasReportResult, error) {
	res, err := s.svc.GetRoasReport(ctx, &service.RoasReportParams{
		AgentId: req.GetAgentId(), DateFrom: req.GetDateFrom(), DateTo: req.GetDateTo(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.RoasReportResult{
		AgentId: res.AgentId, AdSpend: res.AdSpend, Revenue: res.Revenue, Roas: res.Roas, Conversions: res.Conversions,
	}, nil
}

// ---------------------------------------------------------------------------
// Academy / Gamification
// ---------------------------------------------------------------------------

func (s *Server) ListAcademyCourses(ctx context.Context, req *pb.AcademyListRequest) (*pb.AcademyCourseListResult, error) {
	res, err := s.svc.ListAcademyCourses(ctx, &service.AcademyListParams{
		AgentId: req.GetAgentId(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	courses := make([]*pb.AcademyCourseItem, len(res.Courses))
	for i, c := range res.Courses {
		courses[i] = &pb.AcademyCourseItem{CourseId: c.CourseId, Title: c.Title, Description: c.Description, Duration: c.Duration, Level: c.Level}
	}
	return &pb.AcademyCourseListResult{Courses: courses, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) GetCourseProgress(ctx context.Context, req *pb.CourseProgressRequest) (*pb.CourseProgressResult, error) {
	res, err := s.svc.GetCourseProgress(ctx, &service.CourseProgressParams{AgentId: req.GetAgentId(), CourseId: req.GetCourseId()})
	if err != nil {
		return nil, err
	}
	return &pb.CourseProgressResult{
		AgentId: res.AgentId, CourseId: res.CourseId, Progress: res.Progress, Completed: res.Completed, LastAccessAt: res.LastAccessAt,
	}, nil
}

func (s *Server) SubmitQuiz(ctx context.Context, req *pb.SubmitQuizRequest) (*pb.QuizResult, error) {
	res, err := s.svc.SubmitQuiz(ctx, &service.SubmitQuizParams{
		AgentId: req.GetAgentId(), CourseId: req.GetCourseId(), Answers: req.GetAnswers(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.QuizResult{
		AgentId: res.AgentId, CourseId: res.CourseId, Score: res.Score, Passed: res.Passed, Badge: res.Badge, TakenAt: res.TakenAt,
	}, nil
}

func (s *Server) ListSalesScripts(ctx context.Context, req *pb.SalesScriptListRequest) (*pb.SalesScriptListResult, error) {
	res, err := s.svc.ListSalesScripts(ctx, &service.SalesScriptListParams{
		AgentId: req.GetAgentId(), Category: req.GetCategory(), Page: req.GetPage(), PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}
	scripts := make([]*pb.SalesScriptItem, len(res.Scripts))
	for i, sc := range res.Scripts {
		scripts[i] = &pb.SalesScriptItem{ScriptId: sc.ScriptId, Title: sc.Title, Category: sc.Category, Content: sc.Content}
	}
	return &pb.SalesScriptListResult{Scripts: scripts, Total: res.Total, Page: res.Page, PageSize: res.PageSize}, nil
}

func (s *Server) GetLeaderboard(ctx context.Context, req *pb.LeaderboardRequest) (*pb.LeaderboardResult, error) {
	res, err := s.svc.GetLeaderboard(ctx, &service.LeaderboardParams{AgentId: req.GetAgentId(), Period: req.GetPeriod(), Limit: req.GetLimit()})
	if err != nil {
		return nil, err
	}
	entries := make([]*pb.LeaderboardEntry, len(res.Entries))
	for i, e := range res.Entries {
		entries[i] = &pb.LeaderboardEntry{Rank: e.Rank, AgentId: e.AgentId, Name: e.Name, Score: e.Score, Bookings: e.Bookings}
	}
	return &pb.LeaderboardResult{Period: res.Period, Entries: entries, MyRank: res.MyRank}, nil
}

func (s *Server) GetAgentTier(ctx context.Context, req *pb.AgentIdRequest) (*pb.AgentTierResult, error) {
	res, err := s.svc.GetAgentTier(ctx, &service.AgentIdParams{AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	return &pb.AgentTierResult{
		AgentId: res.AgentId, Tier: res.Tier, Points: res.Points, NextTier: res.NextTier, PointsToNext: res.PointsToNext,
	}, nil
}

// ---------------------------------------------------------------------------
// Alumni / Other
// ---------------------------------------------------------------------------

func (s *Server) GetAlumniReferrals(ctx context.Context, req *pb.AgentIdRequest) (*pb.AlumniReferralListResult, error) {
	res, err := s.svc.GetAlumniReferrals(ctx, &service.AgentIdParams{AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	refs := make([]*pb.AlumniReferralItem, len(res.Referrals))
	for i, r := range res.Referrals {
		refs[i] = &pb.AlumniReferralItem{ReferralId: r.ReferralId, Name: r.Name, Phone: r.Phone, Status: r.Status, CreatedAt: r.CreatedAt}
	}
	return &pb.AlumniReferralListResult{AgentId: res.AgentId, Referrals: refs, Total: res.Total}, nil
}

func (s *Server) GetReturnIntentSavings(ctx context.Context, req *pb.AgentIdRequest) (*pb.ReturnIntentResult, error) {
	res, err := s.svc.GetReturnIntentSavings(ctx, &service.AgentIdParams{AgentId: req.GetAgentId()})
	if err != nil {
		return nil, err
	}
	return &pb.ReturnIntentResult{
		AgentId: res.AgentId, SavingsBalance: res.SavingsBalance, TargetAmount: res.TargetAmount,
		Progress: res.Progress, EstimatedDate: res.EstimatedDate,
	}, nil
}

func (s *Server) CalculateZakat(ctx context.Context, req *pb.ZakatCalculatorRequest) (*pb.ZakatResult, error) {
	res, err := s.svc.CalculateZakat(ctx, &service.ZakatCalculatorParams{
		AgentId:        req.GetAgentId(),
		GoldGrams:      req.GetGoldGrams(),
		SilverGrams:    req.GetSilverGrams(),
		CashAmount:     req.GetCashAmount(),
		BusinessAssets: req.GetBusinessAssets(),
		Debts:          req.GetDebts(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ZakatResult{
		AgentId: res.AgentId, TotalAssets: res.TotalAssets, NetAssets: res.NetAssets, ZakatDue: res.ZakatDue, NisabMet: res.NisabMet,
	}, nil
}

func (s *Server) RecordCharity(ctx context.Context, req *pb.CharityRequest) (*pb.CharityResult, error) {
	res, err := s.svc.RecordCharity(ctx, &service.CharityParams{
		AgentId: req.GetAgentId(), Amount: req.GetAmount(), Category: req.GetCategory(), Notes: req.GetNotes(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CharityResult{
		CharityId: res.CharityId, AgentId: res.AgentId, Amount: res.Amount, Category: res.Category, CreatedAt: res.CreatedAt,
	}, nil
}
