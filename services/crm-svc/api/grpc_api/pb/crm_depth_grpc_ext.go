// crm_depth_grpc_ext.go — hand-written gRPC interface extension for Wave 7
// CRM depth RPCs (BL-CRM-010..012, BL-CRM-017..066 = 53 RPCs).
//
// Follows the same pattern as finance_depth_grpc_ext.go.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Server-side handler interface
// ---------------------------------------------------------------------------

// CrmDepthHandler is the server-side interface for Wave 7 CRM depth RPCs.
type CrmDepthHandler interface {
	// Agency Registration
	RegisterAgent(context.Context, *AgentRegisterRequest) (*AgentResult, error)
	SubmitAgentKyc(context.Context, *AgentKycRequest) (*AgentResult, error)
	SignAgentMoU(context.Context, *AgentMoURequest) (*AgentResult, error)
	GetAgentProfile(context.Context, *AgentIdRequest) (*AgentResult, error)
	GetReplicaSite(context.Context, *AgentIdRequest) (*ReplicaSiteResult, error)
	UpdateReplicaSite(context.Context, *UpdateReplicaSiteRequest) (*ReplicaSiteResult, error)
	// Content & Marketing
	GetSocialShareLink(context.Context, *SocialShareRequest) (*SocialShareResult, error)
	GenerateBusinessCard(context.Context, *BusinessCardRequest) (*BusinessCardResult, error)
	ListContentBank(context.Context, *ContentBankListRequest) (*ContentBankListResult, error)
	CreateContentAsset(context.Context, *CreateContentAssetRequest) (*ContentAssetResult, error)
	WatermarkFlyer(context.Context, *WatermarkFlyerRequest) (*WatermarkFlyerResult, error)
	ListProgramGallery(context.Context, *ProgramGalleryRequest) (*ProgramGalleryResult, error)
	SetTrackingCode(context.Context, *TrackingCodeRequest) (*TrackingCodeResult, error)
	GetAdsManagerStats(context.Context, *AdsManagerRequest) (*AdsManagerResult, error)
	CreateUtmLink(context.Context, *CreateUtmLinkRequest) (*UtmLinkResult, error)
	ListUtmLinks(context.Context, *ListUtmLinksRequest) (*UtmLinkListResult, error)
	CreateLandingPage(context.Context, *CreateLandingPageRequest) (*LandingPageResult, error)
	ListLandingPages(context.Context, *ListLandingPagesRequest) (*LandingPageListResult, error)
	ScheduleContent(context.Context, *ScheduleContentRequest) (*ScheduledContentResult, error)
	ListScheduledContent(context.Context, *ListScheduledContentRequest) (*ScheduledContentListResult, error)
	GetContentAnalytics(context.Context, *ContentAnalyticsRequest) (*ContentAnalyticsResult, error)
	// CRM/Leads
	CreateAgentLead(context.Context, *CreateAgentLeadRequest) (*AgentLeadResult, error)
	ListAgentLeads(context.Context, *ListAgentLeadsRequest) (*AgentLeadListResult, error)
	SetLeadReminder(context.Context, *SetLeadReminderRequest) (*LeadReminderResult, error)
	ListLeadReminders(context.Context, *ListLeadRemindersRequest) (*LeadReminderListResult, error)
	FilterBotLeads(context.Context, *BotFilterRequest) (*BotFilterResult, error)
	CreateDripSequence(context.Context, *CreateDripSequenceRequest) (*DripSequenceResult, error)
	ListDripSequences(context.Context, *ListDripSequencesRequest) (*DripSequenceListResult, error)
	CreateMomentTrigger(context.Context, *CreateMomentTriggerRequest) (*MomentTriggerResult, error)
	CreateSegment(context.Context, *CreateSegmentRequest) (*SegmentResult, error)
	ListSegments(context.Context, *ListSegmentsRequest) (*SegmentListResult, error)
	SendBroadcast(context.Context, *SendBroadcastRequest) (*BroadcastResult, error)
	AssignLeadFair(context.Context, *AssignLeadFairRequest) (*AssignLeadFairResult, error)
	CreateSlaRule(context.Context, *CreateSlaRuleRequest) (*SlaRuleResult, error)
	GetLeadTrail(context.Context, *GetLeadTrailRequest) (*LeadTrailResult, error)
	TagLead(context.Context, *TagLeadRequest) (*TagLeadResult, error)
	GenerateQuote(context.Context, *GenerateQuoteRequest) (*QuoteResult, error)
	GetQuote(context.Context, *GetQuoteIdRequest) (*QuoteResult, error)
	BuildPaymentLink(context.Context, *BuildPaymentLinkRequest) (*PaymentLinkResult, error)
	RequestDiscount(context.Context, *RequestDiscountRequest) (*DiscountApprovalResult, error)
	ApproveDiscount(context.Context, *ApproveDiscountRequest) (*DiscountApprovalResult, error)
	GetStaleProspects(context.Context, *StaleProspectRequest) (*StaleProspectListResult, error)
	// Commission/Finance
	GetCommissionBalance(context.Context, *AgentIdRequest) (*CommissionBalanceResult, error)
	GetCommissionEvents(context.Context, *CommissionEventListRequest) (*CommissionEventListResult, error)
	RequestPayout(context.Context, *PayoutRequest) (*PayoutResult, error)
	GetPayoutHistory(context.Context, *AgentIdRequest) (*PayoutHistoryResult, error)
	ComputeOverrideCommission(context.Context, *OverrideCommissionRequest) (*OverrideCommissionResult, error)
	GetRoasReport(context.Context, *RoasReportRequest) (*RoasReportResult, error)
	// Academy/Gamification
	ListAcademyCourses(context.Context, *AcademyListRequest) (*AcademyCourseListResult, error)
	GetCourseProgress(context.Context, *CourseProgressRequest) (*CourseProgressResult, error)
	SubmitQuiz(context.Context, *SubmitQuizRequest) (*QuizResult, error)
	ListSalesScripts(context.Context, *SalesScriptListRequest) (*SalesScriptListResult, error)
	GetLeaderboard(context.Context, *LeaderboardRequest) (*LeaderboardResult, error)
	GetAgentTier(context.Context, *AgentIdRequest) (*AgentTierResult, error)
	// Alumni/Other
	GetAlumniReferrals(context.Context, *AgentIdRequest) (*AlumniReferralListResult, error)
	GetReturnIntentSavings(context.Context, *AgentIdRequest) (*ReturnIntentResult, error)
	CalculateZakat(context.Context, *ZakatCalculatorRequest) (*ZakatResult, error)
	RecordCharity(context.Context, *CharityRequest) (*CharityResult, error)
}

// UnimplementedCrmDepthHandler provides safe defaults.
type UnimplementedCrmDepthHandler struct{}

func (UnimplementedCrmDepthHandler) RegisterAgent(context.Context, *AgentRegisterRequest) (*AgentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAgent not implemented")
}
func (UnimplementedCrmDepthHandler) SubmitAgentKyc(context.Context, *AgentKycRequest) (*AgentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitAgentKyc not implemented")
}
func (UnimplementedCrmDepthHandler) SignAgentMoU(context.Context, *AgentMoURequest) (*AgentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignAgentMoU not implemented")
}
func (UnimplementedCrmDepthHandler) GetAgentProfile(context.Context, *AgentIdRequest) (*AgentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAgentProfile not implemented")
}
func (UnimplementedCrmDepthHandler) GetReplicaSite(context.Context, *AgentIdRequest) (*ReplicaSiteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReplicaSite not implemented")
}
func (UnimplementedCrmDepthHandler) UpdateReplicaSite(context.Context, *UpdateReplicaSiteRequest) (*ReplicaSiteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateReplicaSite not implemented")
}
func (UnimplementedCrmDepthHandler) GetSocialShareLink(context.Context, *SocialShareRequest) (*SocialShareResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSocialShareLink not implemented")
}
func (UnimplementedCrmDepthHandler) GenerateBusinessCard(context.Context, *BusinessCardRequest) (*BusinessCardResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateBusinessCard not implemented")
}
func (UnimplementedCrmDepthHandler) ListContentBank(context.Context, *ContentBankListRequest) (*ContentBankListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListContentBank not implemented")
}
func (UnimplementedCrmDepthHandler) CreateContentAsset(context.Context, *CreateContentAssetRequest) (*ContentAssetResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateContentAsset not implemented")
}
func (UnimplementedCrmDepthHandler) WatermarkFlyer(context.Context, *WatermarkFlyerRequest) (*WatermarkFlyerResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WatermarkFlyer not implemented")
}
func (UnimplementedCrmDepthHandler) ListProgramGallery(context.Context, *ProgramGalleryRequest) (*ProgramGalleryResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProgramGallery not implemented")
}
func (UnimplementedCrmDepthHandler) SetTrackingCode(context.Context, *TrackingCodeRequest) (*TrackingCodeResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetTrackingCode not implemented")
}
func (UnimplementedCrmDepthHandler) GetAdsManagerStats(context.Context, *AdsManagerRequest) (*AdsManagerResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdsManagerStats not implemented")
}
func (UnimplementedCrmDepthHandler) CreateUtmLink(context.Context, *CreateUtmLinkRequest) (*UtmLinkResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUtmLink not implemented")
}
func (UnimplementedCrmDepthHandler) ListUtmLinks(context.Context, *ListUtmLinksRequest) (*UtmLinkListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUtmLinks not implemented")
}
func (UnimplementedCrmDepthHandler) CreateLandingPage(context.Context, *CreateLandingPageRequest) (*LandingPageResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLandingPage not implemented")
}
func (UnimplementedCrmDepthHandler) ListLandingPages(context.Context, *ListLandingPagesRequest) (*LandingPageListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLandingPages not implemented")
}
func (UnimplementedCrmDepthHandler) ScheduleContent(context.Context, *ScheduleContentRequest) (*ScheduledContentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleContent not implemented")
}
func (UnimplementedCrmDepthHandler) ListScheduledContent(context.Context, *ListScheduledContentRequest) (*ScheduledContentListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListScheduledContent not implemented")
}
func (UnimplementedCrmDepthHandler) GetContentAnalytics(context.Context, *ContentAnalyticsRequest) (*ContentAnalyticsResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContentAnalytics not implemented")
}
func (UnimplementedCrmDepthHandler) CreateAgentLead(context.Context, *CreateAgentLeadRequest) (*AgentLeadResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAgentLead not implemented")
}
func (UnimplementedCrmDepthHandler) ListAgentLeads(context.Context, *ListAgentLeadsRequest) (*AgentLeadListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAgentLeads not implemented")
}
func (UnimplementedCrmDepthHandler) SetLeadReminder(context.Context, *SetLeadReminderRequest) (*LeadReminderResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLeadReminder not implemented")
}
func (UnimplementedCrmDepthHandler) ListLeadReminders(context.Context, *ListLeadRemindersRequest) (*LeadReminderListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLeadReminders not implemented")
}
func (UnimplementedCrmDepthHandler) FilterBotLeads(context.Context, *BotFilterRequest) (*BotFilterResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterBotLeads not implemented")
}
func (UnimplementedCrmDepthHandler) CreateDripSequence(context.Context, *CreateDripSequenceRequest) (*DripSequenceResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDripSequence not implemented")
}
func (UnimplementedCrmDepthHandler) ListDripSequences(context.Context, *ListDripSequencesRequest) (*DripSequenceListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDripSequences not implemented")
}
func (UnimplementedCrmDepthHandler) CreateMomentTrigger(context.Context, *CreateMomentTriggerRequest) (*MomentTriggerResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMomentTrigger not implemented")
}
func (UnimplementedCrmDepthHandler) CreateSegment(context.Context, *CreateSegmentRequest) (*SegmentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSegment not implemented")
}
func (UnimplementedCrmDepthHandler) ListSegments(context.Context, *ListSegmentsRequest) (*SegmentListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSegments not implemented")
}
func (UnimplementedCrmDepthHandler) SendBroadcast(context.Context, *SendBroadcastRequest) (*BroadcastResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBroadcast not implemented")
}
func (UnimplementedCrmDepthHandler) AssignLeadFair(context.Context, *AssignLeadFairRequest) (*AssignLeadFairResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignLeadFair not implemented")
}
func (UnimplementedCrmDepthHandler) CreateSlaRule(context.Context, *CreateSlaRuleRequest) (*SlaRuleResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSlaRule not implemented")
}
func (UnimplementedCrmDepthHandler) GetLeadTrail(context.Context, *GetLeadTrailRequest) (*LeadTrailResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeadTrail not implemented")
}
func (UnimplementedCrmDepthHandler) TagLead(context.Context, *TagLeadRequest) (*TagLeadResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TagLead not implemented")
}
func (UnimplementedCrmDepthHandler) GenerateQuote(context.Context, *GenerateQuoteRequest) (*QuoteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateQuote not implemented")
}
func (UnimplementedCrmDepthHandler) GetQuote(context.Context, *GetQuoteIdRequest) (*QuoteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuote not implemented")
}
func (UnimplementedCrmDepthHandler) BuildPaymentLink(context.Context, *BuildPaymentLinkRequest) (*PaymentLinkResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildPaymentLink not implemented")
}
func (UnimplementedCrmDepthHandler) RequestDiscount(context.Context, *RequestDiscountRequest) (*DiscountApprovalResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestDiscount not implemented")
}
func (UnimplementedCrmDepthHandler) ApproveDiscount(context.Context, *ApproveDiscountRequest) (*DiscountApprovalResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveDiscount not implemented")
}
func (UnimplementedCrmDepthHandler) GetStaleProspects(context.Context, *StaleProspectRequest) (*StaleProspectListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStaleProspects not implemented")
}
func (UnimplementedCrmDepthHandler) GetCommissionBalance(context.Context, *AgentIdRequest) (*CommissionBalanceResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommissionBalance not implemented")
}
func (UnimplementedCrmDepthHandler) GetCommissionEvents(context.Context, *CommissionEventListRequest) (*CommissionEventListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommissionEvents not implemented")
}
func (UnimplementedCrmDepthHandler) RequestPayout(context.Context, *PayoutRequest) (*PayoutResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestPayout not implemented")
}
func (UnimplementedCrmDepthHandler) GetPayoutHistory(context.Context, *AgentIdRequest) (*PayoutHistoryResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPayoutHistory not implemented")
}
func (UnimplementedCrmDepthHandler) ComputeOverrideCommission(context.Context, *OverrideCommissionRequest) (*OverrideCommissionResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ComputeOverrideCommission not implemented")
}
func (UnimplementedCrmDepthHandler) GetRoasReport(context.Context, *RoasReportRequest) (*RoasReportResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoasReport not implemented")
}
func (UnimplementedCrmDepthHandler) ListAcademyCourses(context.Context, *AcademyListRequest) (*AcademyCourseListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAcademyCourses not implemented")
}
func (UnimplementedCrmDepthHandler) GetCourseProgress(context.Context, *CourseProgressRequest) (*CourseProgressResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCourseProgress not implemented")
}
func (UnimplementedCrmDepthHandler) SubmitQuiz(context.Context, *SubmitQuizRequest) (*QuizResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitQuiz not implemented")
}
func (UnimplementedCrmDepthHandler) ListSalesScripts(context.Context, *SalesScriptListRequest) (*SalesScriptListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSalesScripts not implemented")
}
func (UnimplementedCrmDepthHandler) GetLeaderboard(context.Context, *LeaderboardRequest) (*LeaderboardResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeaderboard not implemented")
}
func (UnimplementedCrmDepthHandler) GetAgentTier(context.Context, *AgentIdRequest) (*AgentTierResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAgentTier not implemented")
}
func (UnimplementedCrmDepthHandler) GetAlumniReferrals(context.Context, *AgentIdRequest) (*AlumniReferralListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlumniReferrals not implemented")
}
func (UnimplementedCrmDepthHandler) GetReturnIntentSavings(context.Context, *AgentIdRequest) (*ReturnIntentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReturnIntentSavings not implemented")
}
func (UnimplementedCrmDepthHandler) CalculateZakat(context.Context, *ZakatCalculatorRequest) (*ZakatResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateZakat not implemented")
}
func (UnimplementedCrmDepthHandler) RecordCharity(context.Context, *CharityRequest) (*CharityResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordCharity not implemented")
}

// ---------------------------------------------------------------------------
// gRPC handler dispatch functions
// ---------------------------------------------------------------------------

func _CrmService_RegisterAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).RegisterAgent(ctx, req.(*AgentRegisterRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/RegisterAgent"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_SubmitAgentKyc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentKycRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).SubmitAgentKyc(ctx, req.(*AgentKycRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/SubmitAgentKyc"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_SignAgentMoU_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentMoURequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).SignAgentMoU(ctx, req.(*AgentMoURequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/SignAgentMoU"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetAgentProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetAgentProfile(ctx, req.(*AgentIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetAgentProfile"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetReplicaSite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetReplicaSite(ctx, req.(*AgentIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetReplicaSite"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_UpdateReplicaSite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReplicaSiteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).UpdateReplicaSite(ctx, req.(*UpdateReplicaSiteRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/UpdateReplicaSite"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetSocialShareLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SocialShareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetSocialShareLink(ctx, req.(*SocialShareRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetSocialShareLink"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GenerateBusinessCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BusinessCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GenerateBusinessCard(ctx, req.(*BusinessCardRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GenerateBusinessCard"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListContentBank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContentBankListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListContentBank(ctx, req.(*ContentBankListRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListContentBank"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateContentAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateContentAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateContentAsset(ctx, req.(*CreateContentAssetRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateContentAsset"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_WatermarkFlyer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WatermarkFlyerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).WatermarkFlyer(ctx, req.(*WatermarkFlyerRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/WatermarkFlyer"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListProgramGallery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProgramGalleryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListProgramGallery(ctx, req.(*ProgramGalleryRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListProgramGallery"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_SetTrackingCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackingCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).SetTrackingCode(ctx, req.(*TrackingCodeRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/SetTrackingCode"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetAdsManagerStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdsManagerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetAdsManagerStats(ctx, req.(*AdsManagerRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetAdsManagerStats"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateUtmLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUtmLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateUtmLink(ctx, req.(*CreateUtmLinkRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateUtmLink"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListUtmLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUtmLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListUtmLinks(ctx, req.(*ListUtmLinksRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListUtmLinks"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateLandingPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLandingPageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateLandingPage(ctx, req.(*CreateLandingPageRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateLandingPage"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListLandingPages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLandingPagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListLandingPages(ctx, req.(*ListLandingPagesRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListLandingPages"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ScheduleContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScheduleContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ScheduleContent(ctx, req.(*ScheduleContentRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ScheduleContent"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListScheduledContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListScheduledContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListScheduledContent(ctx, req.(*ListScheduledContentRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListScheduledContent"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetContentAnalytics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContentAnalyticsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetContentAnalytics(ctx, req.(*ContentAnalyticsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetContentAnalytics"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateAgentLead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAgentLeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateAgentLead(ctx, req.(*CreateAgentLeadRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateAgentLead"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListAgentLeads_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAgentLeadsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListAgentLeads(ctx, req.(*ListAgentLeadsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListAgentLeads"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_SetLeadReminder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetLeadReminderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).SetLeadReminder(ctx, req.(*SetLeadReminderRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/SetLeadReminder"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListLeadReminders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLeadRemindersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListLeadReminders(ctx, req.(*ListLeadRemindersRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListLeadReminders"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_FilterBotLeads_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BotFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).FilterBotLeads(ctx, req.(*BotFilterRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/FilterBotLeads"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateDripSequence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDripSequenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateDripSequence(ctx, req.(*CreateDripSequenceRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateDripSequence"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListDripSequences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDripSequencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListDripSequences(ctx, req.(*ListDripSequencesRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListDripSequences"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateMomentTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMomentTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateMomentTrigger(ctx, req.(*CreateMomentTriggerRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateMomentTrigger"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateSegment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSegmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateSegment(ctx, req.(*CreateSegmentRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateSegment"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListSegments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSegmentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListSegments(ctx, req.(*ListSegmentsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListSegments"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_SendBroadcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendBroadcastRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).SendBroadcast(ctx, req.(*SendBroadcastRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/SendBroadcast"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_AssignLeadFair_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignLeadFairRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).AssignLeadFair(ctx, req.(*AssignLeadFairRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/AssignLeadFair"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CreateSlaRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSlaRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CreateSlaRule(ctx, req.(*CreateSlaRuleRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CreateSlaRule"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetLeadTrail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLeadTrailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetLeadTrail(ctx, req.(*GetLeadTrailRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetLeadTrail"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_TagLead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TagLeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).TagLead(ctx, req.(*TagLeadRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/TagLead"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GenerateQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateQuoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GenerateQuote(ctx, req.(*GenerateQuoteRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GenerateQuote"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetQuoteIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetQuote(ctx, req.(*GetQuoteIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetQuote"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_BuildPaymentLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuildPaymentLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).BuildPaymentLink(ctx, req.(*BuildPaymentLinkRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/BuildPaymentLink"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_RequestDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).RequestDiscount(ctx, req.(*RequestDiscountRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/RequestDiscount"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ApproveDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ApproveDiscount(ctx, req.(*ApproveDiscountRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ApproveDiscount"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetStaleProspects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StaleProspectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetStaleProspects(ctx, req.(*StaleProspectRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetStaleProspects"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetCommissionBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetCommissionBalance(ctx, req.(*AgentIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetCommissionBalance"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetCommissionEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommissionEventListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetCommissionEvents(ctx, req.(*CommissionEventListRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetCommissionEvents"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_RequestPayout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PayoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).RequestPayout(ctx, req.(*PayoutRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/RequestPayout"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetPayoutHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetPayoutHistory(ctx, req.(*AgentIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetPayoutHistory"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ComputeOverrideCommission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OverrideCommissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ComputeOverrideCommission(ctx, req.(*OverrideCommissionRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ComputeOverrideCommission"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetRoasReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoasReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetRoasReport(ctx, req.(*RoasReportRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetRoasReport"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListAcademyCourses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcademyListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListAcademyCourses(ctx, req.(*AcademyListRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListAcademyCourses"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetCourseProgress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseProgressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetCourseProgress(ctx, req.(*CourseProgressRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetCourseProgress"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_SubmitQuiz_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitQuizRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).SubmitQuiz(ctx, req.(*SubmitQuizRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/SubmitQuiz"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_ListSalesScripts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SalesScriptListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).ListSalesScripts(ctx, req.(*SalesScriptListRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/ListSalesScripts"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetLeaderboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetLeaderboard(ctx, req.(*LeaderboardRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetLeaderboard"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetAgentTier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetAgentTier(ctx, req.(*AgentIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetAgentTier"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetAlumniReferrals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetAlumniReferrals(ctx, req.(*AgentIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetAlumniReferrals"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_GetReturnIntentSavings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).GetReturnIntentSavings(ctx, req.(*AgentIdRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/GetReturnIntentSavings"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_CalculateZakat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ZakatCalculatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).CalculateZakat(ctx, req.(*ZakatCalculatorRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/CalculateZakat"}
	return interceptor(ctx, in, info, handler)
}

func _CrmService_RecordCharity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CharityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrmDepthHandler).RecordCharity(ctx, req.(*CharityRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.CrmService/RecordCharity"}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// Method descriptors
// ---------------------------------------------------------------------------

var _CrmDepth_MethodDescs = []grpc.MethodDesc{
	{MethodName: "RegisterAgent", Handler: _CrmService_RegisterAgent_Handler},
	{MethodName: "SubmitAgentKyc", Handler: _CrmService_SubmitAgentKyc_Handler},
	{MethodName: "SignAgentMoU", Handler: _CrmService_SignAgentMoU_Handler},
	{MethodName: "GetAgentProfile", Handler: _CrmService_GetAgentProfile_Handler},
	{MethodName: "GetReplicaSite", Handler: _CrmService_GetReplicaSite_Handler},
	{MethodName: "UpdateReplicaSite", Handler: _CrmService_UpdateReplicaSite_Handler},
	{MethodName: "GetSocialShareLink", Handler: _CrmService_GetSocialShareLink_Handler},
	{MethodName: "GenerateBusinessCard", Handler: _CrmService_GenerateBusinessCard_Handler},
	{MethodName: "ListContentBank", Handler: _CrmService_ListContentBank_Handler},
	{MethodName: "CreateContentAsset", Handler: _CrmService_CreateContentAsset_Handler},
	{MethodName: "WatermarkFlyer", Handler: _CrmService_WatermarkFlyer_Handler},
	{MethodName: "ListProgramGallery", Handler: _CrmService_ListProgramGallery_Handler},
	{MethodName: "SetTrackingCode", Handler: _CrmService_SetTrackingCode_Handler},
	{MethodName: "GetAdsManagerStats", Handler: _CrmService_GetAdsManagerStats_Handler},
	{MethodName: "CreateUtmLink", Handler: _CrmService_CreateUtmLink_Handler},
	{MethodName: "ListUtmLinks", Handler: _CrmService_ListUtmLinks_Handler},
	{MethodName: "CreateLandingPage", Handler: _CrmService_CreateLandingPage_Handler},
	{MethodName: "ListLandingPages", Handler: _CrmService_ListLandingPages_Handler},
	{MethodName: "ScheduleContent", Handler: _CrmService_ScheduleContent_Handler},
	{MethodName: "ListScheduledContent", Handler: _CrmService_ListScheduledContent_Handler},
	{MethodName: "GetContentAnalytics", Handler: _CrmService_GetContentAnalytics_Handler},
	{MethodName: "CreateAgentLead", Handler: _CrmService_CreateAgentLead_Handler},
	{MethodName: "ListAgentLeads", Handler: _CrmService_ListAgentLeads_Handler},
	{MethodName: "SetLeadReminder", Handler: _CrmService_SetLeadReminder_Handler},
	{MethodName: "ListLeadReminders", Handler: _CrmService_ListLeadReminders_Handler},
	{MethodName: "FilterBotLeads", Handler: _CrmService_FilterBotLeads_Handler},
	{MethodName: "CreateDripSequence", Handler: _CrmService_CreateDripSequence_Handler},
	{MethodName: "ListDripSequences", Handler: _CrmService_ListDripSequences_Handler},
	{MethodName: "CreateMomentTrigger", Handler: _CrmService_CreateMomentTrigger_Handler},
	{MethodName: "CreateSegment", Handler: _CrmService_CreateSegment_Handler},
	{MethodName: "ListSegments", Handler: _CrmService_ListSegments_Handler},
	{MethodName: "SendBroadcast", Handler: _CrmService_SendBroadcast_Handler},
	{MethodName: "AssignLeadFair", Handler: _CrmService_AssignLeadFair_Handler},
	{MethodName: "CreateSlaRule", Handler: _CrmService_CreateSlaRule_Handler},
	{MethodName: "GetLeadTrail", Handler: _CrmService_GetLeadTrail_Handler},
	{MethodName: "TagLead", Handler: _CrmService_TagLead_Handler},
	{MethodName: "GenerateQuote", Handler: _CrmService_GenerateQuote_Handler},
	{MethodName: "GetQuote", Handler: _CrmService_GetQuote_Handler},
	{MethodName: "BuildPaymentLink", Handler: _CrmService_BuildPaymentLink_Handler},
	{MethodName: "RequestDiscount", Handler: _CrmService_RequestDiscount_Handler},
	{MethodName: "ApproveDiscount", Handler: _CrmService_ApproveDiscount_Handler},
	{MethodName: "GetStaleProspects", Handler: _CrmService_GetStaleProspects_Handler},
	{MethodName: "GetCommissionBalance", Handler: _CrmService_GetCommissionBalance_Handler},
	{MethodName: "GetCommissionEvents", Handler: _CrmService_GetCommissionEvents_Handler},
	{MethodName: "RequestPayout", Handler: _CrmService_RequestPayout_Handler},
	{MethodName: "GetPayoutHistory", Handler: _CrmService_GetPayoutHistory_Handler},
	{MethodName: "ComputeOverrideCommission", Handler: _CrmService_ComputeOverrideCommission_Handler},
	{MethodName: "GetRoasReport", Handler: _CrmService_GetRoasReport_Handler},
	{MethodName: "ListAcademyCourses", Handler: _CrmService_ListAcademyCourses_Handler},
	{MethodName: "GetCourseProgress", Handler: _CrmService_GetCourseProgress_Handler},
	{MethodName: "SubmitQuiz", Handler: _CrmService_SubmitQuiz_Handler},
	{MethodName: "ListSalesScripts", Handler: _CrmService_ListSalesScripts_Handler},
	{MethodName: "GetLeaderboard", Handler: _CrmService_GetLeaderboard_Handler},
	{MethodName: "GetAgentTier", Handler: _CrmService_GetAgentTier_Handler},
	{MethodName: "GetAlumniReferrals", Handler: _CrmService_GetAlumniReferrals_Handler},
	{MethodName: "GetReturnIntentSavings", Handler: _CrmService_GetReturnIntentSavings_Handler},
	{MethodName: "CalculateZakat", Handler: _CrmService_CalculateZakat_Handler},
	{MethodName: "RecordCharity", Handler: _CrmService_RecordCharity_Handler},
}

// ---------------------------------------------------------------------------
// RegisterCrmServiceServerWithDepth — registration function for Wave 7.
// Does NOT replace RegisterCrmServiceServer; call both.
// ---------------------------------------------------------------------------

// RegisterCrmServiceServerWithDepth registers the CRM depth RPCs (Wave 7).
// Call this after pb.RegisterCrmServiceServer in runGrpcServer.
func RegisterCrmServiceServerWithDepth(s *grpc.Server, srv CrmDepthHandler) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "pb.CrmService",
		HandlerType: (*CrmDepthHandler)(nil),
		Methods:     _CrmDepth_MethodDescs,
		Streams:     []grpc.StreamDesc{},
		Metadata:    "crm.proto",
	}, srv)
}
