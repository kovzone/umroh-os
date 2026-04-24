// proxy_dispatch_crm_depth.go — thin dispatch methods for Wave 7 CRM depth RPCs.
// Each method delegates directly to the crm_grpc_adapter.Adapter.
package service

import (
	"context"

	"gateway-svc/adapter/crm_grpc_adapter"
)

func (s *Service) RegisterAgent(ctx context.Context, params *crm_grpc_adapter.RegisterAgentParams) (*crm_grpc_adapter.AgentAdapterResult, error) {
	return s.adapters.crmGrpc.RegisterAgent(ctx, params)
}

func (s *Service) SubmitAgentKyc(ctx context.Context, params *crm_grpc_adapter.SubmitAgentKycParams) (*crm_grpc_adapter.KycAdapterResult, error) {
	return s.adapters.crmGrpc.SubmitAgentKyc(ctx, params)
}

func (s *Service) SignAgentMoU(ctx context.Context, params *crm_grpc_adapter.SignAgentMoUParams) (*crm_grpc_adapter.MoUAdapterResult, error) {
	return s.adapters.crmGrpc.SignAgentMoU(ctx, params)
}

func (s *Service) GetAgentProfile(ctx context.Context, agentId string) (*crm_grpc_adapter.AgentAdapterResult, error) {
	return s.adapters.crmGrpc.GetAgentProfile(ctx, agentId)
}

func (s *Service) GetReplicaSite(ctx context.Context, agentId string) (*crm_grpc_adapter.ReplicaSiteAdapterResult, error) {
	return s.adapters.crmGrpc.GetReplicaSite(ctx, agentId)
}

func (s *Service) UpdateReplicaSite(ctx context.Context, params *crm_grpc_adapter.UpdateReplicaSiteParams) (*crm_grpc_adapter.ReplicaSiteAdapterResult, error) {
	return s.adapters.crmGrpc.UpdateReplicaSite(ctx, params)
}

func (s *Service) GetSocialShareLink(ctx context.Context, params *crm_grpc_adapter.GetSocialShareLinkParams) (*crm_grpc_adapter.ShareLinkAdapterResult, error) {
	return s.adapters.crmGrpc.GetSocialShareLink(ctx, params)
}

func (s *Service) GenerateBusinessCard(ctx context.Context, params *crm_grpc_adapter.GenerateBusinessCardParams) (*crm_grpc_adapter.BusinessCardAdapterResult, error) {
	return s.adapters.crmGrpc.GenerateBusinessCard(ctx, params)
}

func (s *Service) ListContentBank(ctx context.Context, params *crm_grpc_adapter.ListContentBankParams) (*crm_grpc_adapter.ListContentBankResult, error) {
	return s.adapters.crmGrpc.ListContentBank(ctx, params)
}

func (s *Service) CreateContentAsset(ctx context.Context, params *crm_grpc_adapter.CreateContentAssetParams) (*crm_grpc_adapter.ContentAssetRowResult, error) {
	return s.adapters.crmGrpc.CreateContentAsset(ctx, params)
}

func (s *Service) WatermarkFlyer(ctx context.Context, params *crm_grpc_adapter.WatermarkFlyerParams) (*crm_grpc_adapter.WatermarkFlyerResult, error) {
	return s.adapters.crmGrpc.WatermarkFlyer(ctx, params)
}

func (s *Service) ListProgramGallery(ctx context.Context, params *crm_grpc_adapter.ListProgramGalleryParams) (*crm_grpc_adapter.ListProgramGalleryResult, error) {
	return s.adapters.crmGrpc.ListProgramGallery(ctx, params)
}

func (s *Service) SetTrackingCode(ctx context.Context, params *crm_grpc_adapter.SetTrackingCodeParams) (*crm_grpc_adapter.TrackingCodeAdapterResult, error) {
	return s.adapters.crmGrpc.SetTrackingCode(ctx, params)
}

func (s *Service) GetAdsManagerStats(ctx context.Context, params *crm_grpc_adapter.GetAdsManagerStatsParams) (*crm_grpc_adapter.AdsStatsAdapterResult, error) {
	return s.adapters.crmGrpc.GetAdsManagerStats(ctx, params)
}

func (s *Service) CreateUtmLink(ctx context.Context, params *crm_grpc_adapter.CreateUtmLinkParams) (*crm_grpc_adapter.UtmLinkAdapterResult, error) {
	return s.adapters.crmGrpc.CreateUtmLink(ctx, params)
}

func (s *Service) ListUtmLinks(ctx context.Context, params *crm_grpc_adapter.ListUtmLinksParams) (*crm_grpc_adapter.ListUtmLinksResult, error) {
	return s.adapters.crmGrpc.ListUtmLinks(ctx, params)
}

func (s *Service) CreateLandingPage(ctx context.Context, params *crm_grpc_adapter.CreateLandingPageParams) (*crm_grpc_adapter.LandingPageAdapterResult, error) {
	return s.adapters.crmGrpc.CreateLandingPage(ctx, params)
}

func (s *Service) ListLandingPages(ctx context.Context, params *crm_grpc_adapter.ListLandingPagesParams) (*crm_grpc_adapter.ListLandingPagesResult, error) {
	return s.adapters.crmGrpc.ListLandingPages(ctx, params)
}

func (s *Service) ScheduleContent(ctx context.Context, params *crm_grpc_adapter.ScheduleContentParams) (*crm_grpc_adapter.ScheduleContentAdapterResult, error) {
	return s.adapters.crmGrpc.ScheduleContent(ctx, params)
}

func (s *Service) ListScheduledContent(ctx context.Context, params *crm_grpc_adapter.ListScheduledContentParams) (*crm_grpc_adapter.ListScheduledContentResult, error) {
	return s.adapters.crmGrpc.ListScheduledContent(ctx, params)
}

func (s *Service) GetContentAnalytics(ctx context.Context, params *crm_grpc_adapter.GetContentAnalyticsParams) (*crm_grpc_adapter.ContentAnalyticsAdapterResult, error) {
	return s.adapters.crmGrpc.GetContentAnalytics(ctx, params)
}

func (s *Service) CreateAgentLead(ctx context.Context, params *crm_grpc_adapter.CreateAgentLeadParams) (*crm_grpc_adapter.AgentLeadAdapterResult, error) {
	return s.adapters.crmGrpc.CreateAgentLead(ctx, params)
}

func (s *Service) ListAgentLeads(ctx context.Context, params *crm_grpc_adapter.ListAgentLeadsParams) (*crm_grpc_adapter.ListAgentLeadsResult, error) {
	return s.adapters.crmGrpc.ListAgentLeads(ctx, params)
}

func (s *Service) SetLeadReminder(ctx context.Context, params *crm_grpc_adapter.SetLeadReminderParams) (*crm_grpc_adapter.SetLeadReminderAdapterResult, error) {
	return s.adapters.crmGrpc.SetLeadReminder(ctx, params)
}

func (s *Service) ListLeadReminders(ctx context.Context, leadId string) (*crm_grpc_adapter.ListLeadRemindersResult, error) {
	return s.adapters.crmGrpc.ListLeadReminders(ctx, leadId)
}

func (s *Service) FilterBotLeads(ctx context.Context, agentId string) (*crm_grpc_adapter.FilterBotLeadsAdapterResult, error) {
	return s.adapters.crmGrpc.FilterBotLeads(ctx, agentId)
}

func (s *Service) CreateDripSequence(ctx context.Context, params *crm_grpc_adapter.CreateDripSequenceParams) (*crm_grpc_adapter.DripSequenceAdapterResult, error) {
	return s.adapters.crmGrpc.CreateDripSequence(ctx, params)
}

func (s *Service) ListDripSequences(ctx context.Context, params *crm_grpc_adapter.ListDripSequencesParams) (*crm_grpc_adapter.ListDripSequencesResult, error) {
	return s.adapters.crmGrpc.ListDripSequences(ctx, params)
}

func (s *Service) CreateMomentTrigger(ctx context.Context, params *crm_grpc_adapter.CreateMomentTriggerParams) (*crm_grpc_adapter.MomentTriggerAdapterResult, error) {
	return s.adapters.crmGrpc.CreateMomentTrigger(ctx, params)
}

func (s *Service) CreateSegment(ctx context.Context, params *crm_grpc_adapter.CreateSegmentParams) (*crm_grpc_adapter.SegmentAdapterResult, error) {
	return s.adapters.crmGrpc.CreateSegment(ctx, params)
}

func (s *Service) ListSegments(ctx context.Context, params *crm_grpc_adapter.ListSegmentsParams) (*crm_grpc_adapter.ListSegmentsResult, error) {
	return s.adapters.crmGrpc.ListSegments(ctx, params)
}

func (s *Service) SendBroadcast(ctx context.Context, params *crm_grpc_adapter.SendBroadcastParams) (*crm_grpc_adapter.SendBroadcastAdapterResult, error) {
	return s.adapters.crmGrpc.SendBroadcast(ctx, params)
}

func (s *Service) AssignLeadFair(ctx context.Context, params *crm_grpc_adapter.AssignLeadFairParams) (*crm_grpc_adapter.AssignLeadFairAdapterResult, error) {
	return s.adapters.crmGrpc.AssignLeadFair(ctx, params)
}

func (s *Service) CreateSlaRule(ctx context.Context, params *crm_grpc_adapter.CreateSlaRuleParams) (*crm_grpc_adapter.SlaRuleAdapterResult, error) {
	return s.adapters.crmGrpc.CreateSlaRule(ctx, params)
}

func (s *Service) GetLeadTrail(ctx context.Context, leadId string) (*crm_grpc_adapter.GetLeadTrailAdapterResult, error) {
	return s.adapters.crmGrpc.GetLeadTrail(ctx, leadId)
}

func (s *Service) TagLead(ctx context.Context, params *crm_grpc_adapter.TagLeadParams) (*crm_grpc_adapter.TagLeadAdapterResult, error) {
	return s.adapters.crmGrpc.TagLead(ctx, params)
}

func (s *Service) GenerateQuote(ctx context.Context, params *crm_grpc_adapter.GenerateQuoteParams) (*crm_grpc_adapter.QuoteAdapterResult, error) {
	return s.adapters.crmGrpc.GenerateQuote(ctx, params)
}

func (s *Service) GetQuote(ctx context.Context, quoteId string) (*crm_grpc_adapter.QuoteAdapterResult, error) {
	return s.adapters.crmGrpc.GetQuote(ctx, quoteId)
}

func (s *Service) BuildPaymentLink(ctx context.Context, params *crm_grpc_adapter.BuildPaymentLinkParams) (*crm_grpc_adapter.PaymentLinkAdapterResult, error) {
	return s.adapters.crmGrpc.BuildPaymentLink(ctx, params)
}

func (s *Service) RequestDiscount(ctx context.Context, params *crm_grpc_adapter.RequestDiscountParams) (*crm_grpc_adapter.DiscountAdapterResult, error) {
	return s.adapters.crmGrpc.RequestDiscount(ctx, params)
}

func (s *Service) ApproveDiscount(ctx context.Context, params *crm_grpc_adapter.ApproveDiscountParams) (*crm_grpc_adapter.DiscountAdapterResult, error) {
	return s.adapters.crmGrpc.ApproveDiscount(ctx, params)
}

func (s *Service) GetStaleProspects(ctx context.Context, params *crm_grpc_adapter.GetStaleProspectsParams) (*crm_grpc_adapter.GetStaleProspectsResult, error) {
	return s.adapters.crmGrpc.GetStaleProspects(ctx, params)
}

func (s *Service) GetCommissionBalance(ctx context.Context, agentId string) (*crm_grpc_adapter.CommissionBalanceAdapterResult, error) {
	return s.adapters.crmGrpc.GetCommissionBalance(ctx, agentId)
}

func (s *Service) GetCommissionEvents(ctx context.Context, agentId string) (*crm_grpc_adapter.GetCommissionEventsAdapterResult, error) {
	return s.adapters.crmGrpc.GetCommissionEvents(ctx, agentId)
}

func (s *Service) RequestPayout(ctx context.Context, params *crm_grpc_adapter.RequestPayoutParams) (*crm_grpc_adapter.PayoutAdapterResult, error) {
	return s.adapters.crmGrpc.RequestPayout(ctx, params)
}

func (s *Service) GetPayoutHistory(ctx context.Context, agentId string) (*crm_grpc_adapter.GetPayoutHistoryAdapterResult, error) {
	return s.adapters.crmGrpc.GetPayoutHistory(ctx, agentId)
}

func (s *Service) ComputeOverrideCommission(ctx context.Context, params *crm_grpc_adapter.ComputeOverrideParams) (*crm_grpc_adapter.OverrideCommissionAdapterResult, error) {
	return s.adapters.crmGrpc.ComputeOverrideCommission(ctx, params)
}

func (s *Service) GetRoasReport(ctx context.Context, params *crm_grpc_adapter.GetRoasReportParams) (*crm_grpc_adapter.RoasReportAdapterResult, error) {
	return s.adapters.crmGrpc.GetRoasReport(ctx, params)
}

func (s *Service) ListAcademyCourses(ctx context.Context, params *crm_grpc_adapter.ListAcademyCoursesParams) (*crm_grpc_adapter.ListAcademyCoursesResult, error) {
	return s.adapters.crmGrpc.ListAcademyCourses(ctx, params)
}

func (s *Service) GetCourseProgress(ctx context.Context, params *crm_grpc_adapter.GetCourseProgressParams) (*crm_grpc_adapter.CourseProgressAdapterResult, error) {
	return s.adapters.crmGrpc.GetCourseProgress(ctx, params)
}

func (s *Service) SubmitQuiz(ctx context.Context, params *crm_grpc_adapter.SubmitQuizParams) (*crm_grpc_adapter.QuizAdapterResult, error) {
	return s.adapters.crmGrpc.SubmitQuiz(ctx, params)
}

func (s *Service) ListSalesScripts(ctx context.Context, params *crm_grpc_adapter.ListSalesScriptsParams) (*crm_grpc_adapter.ListSalesScriptsResult, error) {
	return s.adapters.crmGrpc.ListSalesScripts(ctx, params)
}

func (s *Service) GetLeaderboard(ctx context.Context, params *crm_grpc_adapter.GetLeaderboardParams) (*crm_grpc_adapter.LeaderboardAdapterResult, error) {
	return s.adapters.crmGrpc.GetLeaderboard(ctx, params)
}

func (s *Service) GetAgentTier(ctx context.Context, agentId string) (*crm_grpc_adapter.AgentTierAdapterResult, error) {
	return s.adapters.crmGrpc.GetAgentTier(ctx, agentId)
}

func (s *Service) GetAlumniReferrals(ctx context.Context, agentId string) (*crm_grpc_adapter.GetAlumniReferralsAdapterResult, error) {
	return s.adapters.crmGrpc.GetAlumniReferrals(ctx, agentId)
}

func (s *Service) GetReturnIntentSavings(ctx context.Context, agentId string) (*crm_grpc_adapter.ReturnIntentAdapterResult, error) {
	return s.adapters.crmGrpc.GetReturnIntentSavings(ctx, agentId)
}

func (s *Service) CalculateZakat(ctx context.Context, params *crm_grpc_adapter.CalculateZakatParams) (*crm_grpc_adapter.ZakatAdapterResult, error) {
	return s.adapters.crmGrpc.CalculateZakat(ctx, params)
}

func (s *Service) RecordCharity(ctx context.Context, params *crm_grpc_adapter.RecordCharityParams) (*crm_grpc_adapter.CharityAdapterResult, error) {
	return s.adapters.crmGrpc.RecordCharity(ctx, params)
}
