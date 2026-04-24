// proxy_crm_depth.go — gateway REST handlers for CRM agency Wave 7 depth RPCs.
// BL-CRM-010..012, 017..066

package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/crm_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Agency Registration
// ---------------------------------------------------------------------------

func (s *Server) RegisterAgent(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RegisterAgent"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgencyName string `json:"agency_name"`
		OwnerName  string `json:"owner_name"`
		Phone      string `json:"phone"`
		Email      string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RegisterAgent(ctx, &crm_grpc_adapter.RegisterAgentParams{
		AgencyName: body.AgencyName,
		OwnerName:  body.OwnerName,
		Phone:      body.Phone,
		Email:      body.Email,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) SubmitAgentKyc(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SubmitAgentKyc"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId string `json:"agent_id"`
		KycType string `json:"kyc_type"`
		DocUrl  string `json:"doc_url"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SubmitAgentKyc(ctx, &crm_grpc_adapter.SubmitAgentKycParams{
		AgentId: body.AgentId,
		KycType: body.KycType,
		DocUrl:  body.DocUrl,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) SignAgentMoU(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SignAgentMoU"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId      string `json:"agent_id"`
		SignedDocUrl string `json:"signed_doc_url"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SignAgentMoU(ctx, &crm_grpc_adapter.SignAgentMoUParams{
		AgentId:      body.AgentId,
		SignedDocUrl: body.SignedDocUrl,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetAgentProfile(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetAgentProfile"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetAgentProfile(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetReplicaSite(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetReplicaSite"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetReplicaSite(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) UpdateReplicaSite(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.UpdateReplicaSite"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId    string `json:"agent_id"`
		CustomSlug string `json:"custom_slug"`
		ThemeColor string `json:"theme_color"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.UpdateReplicaSite(ctx, &crm_grpc_adapter.UpdateReplicaSiteParams{
		AgentId:    body.AgentId,
		CustomSlug: body.CustomSlug,
		ThemeColor: body.ThemeColor,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// Content
// ---------------------------------------------------------------------------

func (s *Server) GetSocialShareLink(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetSocialShareLink"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId   string `json:"agent_id"`
		Platform  string `json:"platform"`
		ContentId string `json:"content_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.GetSocialShareLink(ctx, &crm_grpc_adapter.GetSocialShareLinkParams{
		AgentId:   body.AgentId,
		Platform:  body.Platform,
		ContentId: body.ContentId,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GenerateBusinessCard(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GenerateBusinessCard"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId    string `json:"agent_id"`
		TemplateId string `json:"template_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.GenerateBusinessCard(ctx, &crm_grpc_adapter.GenerateBusinessCardParams{
		AgentId:    body.AgentId,
		TemplateId: body.TemplateId,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) ListContentBank(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListContentBank"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	agentId := c.Query("agent_id")
	category := c.Query("category")
	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListContentBank(ctx, &crm_grpc_adapter.ListContentBankParams{
		AgentId:  agentId,
		Category: category,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) CreateContentAsset(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateContentAsset"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId  string `json:"agent_id"`
		Title    string `json:"title"`
		Category string `json:"category"`
		AssetUrl string `json:"asset_url"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateContentAsset(ctx, &crm_grpc_adapter.CreateContentAssetParams{
		AgentId:  body.AgentId,
		Title:    body.Title,
		Category: body.Category,
		AssetUrl: body.AssetUrl,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) WatermarkFlyer(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.WatermarkFlyer"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId       string `json:"agent_id"`
		FlyerUrl      string `json:"flyer_url"`
		WatermarkText string `json:"watermark_text"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.WatermarkFlyer(ctx, &crm_grpc_adapter.WatermarkFlyerParams{
		AgentId:       body.AgentId,
		FlyerUrl:      body.FlyerUrl,
		WatermarkText: body.WatermarkText,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) ListProgramGallery(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListProgramGallery"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	packageId := c.Query("package_id")
	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListProgramGallery(ctx, &crm_grpc_adapter.ListProgramGalleryParams{
		PackageId: packageId,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) SetTrackingCode(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SetTrackingCode"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId      string `json:"agent_id"`
		TrackingCode string `json:"tracking_code"`
		Platform     string `json:"platform"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SetTrackingCode(ctx, &crm_grpc_adapter.SetTrackingCodeParams{
		AgentId:      body.AgentId,
		TrackingCode: body.TrackingCode,
		Platform:     body.Platform,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetAdsManagerStats(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetAdsManagerStats"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	result, err := s.svc.GetAdsManagerStats(ctx, &crm_grpc_adapter.GetAdsManagerStatsParams{
		AgentId:  agentId,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) CreateUtmLink(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateUtmLink"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId     string `json:"agent_id"`
		DestUrl     string `json:"dest_url"`
		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateUtmLink(ctx, &crm_grpc_adapter.CreateUtmLinkParams{
		AgentId:     body.AgentId,
		DestUrl:     body.DestUrl,
		UtmSource:   body.UtmSource,
		UtmMedium:   body.UtmMedium,
		UtmCampaign: body.UtmCampaign,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ListUtmLinks(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.ListUtmLinks"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListUtmLinks(ctx, &crm_grpc_adapter.ListUtmLinksParams{
		AgentId:  agentId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) CreateLandingPage(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateLandingPage"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId    string `json:"agent_id"`
		Title      string `json:"title"`
		TemplateId string `json:"template_id"`
		PackageId  string `json:"package_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateLandingPage(ctx, &crm_grpc_adapter.CreateLandingPageParams{
		AgentId:    body.AgentId,
		Title:      body.Title,
		TemplateId: body.TemplateId,
		PackageId:  body.PackageId,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ListLandingPages(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.ListLandingPages"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListLandingPages(ctx, &crm_grpc_adapter.ListLandingPagesParams{
		AgentId:  agentId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) ScheduleContent(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ScheduleContent"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId     string `json:"agent_id"`
		ContentId   string `json:"content_id"`
		Platform    string `json:"platform"`
		ScheduledAt string `json:"scheduled_at"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.ScheduleContent(ctx, &crm_grpc_adapter.ScheduleContentParams{
		AgentId:     body.AgentId,
		ContentId:   body.ContentId,
		Platform:    body.Platform,
		ScheduledAt: body.ScheduledAt,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ListScheduledContent(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.ListScheduledContent"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListScheduledContent(ctx, &crm_grpc_adapter.ListScheduledContentParams{
		AgentId:  agentId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetContentAnalytics(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetContentAnalytics"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	result, err := s.svc.GetContentAnalytics(ctx, &crm_grpc_adapter.GetContentAnalyticsParams{
		AgentId:  agentId,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// CRM / Leads
// ---------------------------------------------------------------------------

func (s *Server) CreateAgentLead(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateAgentLead"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId   string `json:"agent_id"`
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		Source    string `json:"source"`
		PackageId string `json:"package_id"`
		Notes     string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateAgentLead(ctx, &crm_grpc_adapter.CreateAgentLeadParams{
		AgentId:   body.AgentId,
		Name:      body.Name,
		Phone:     body.Phone,
		Email:     body.Email,
		Source:    body.Source,
		PackageId: body.PackageId,
		Notes:     body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ListAgentLeads(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.ListAgentLeads"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	statusFilter := c.Query("status")
	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListAgentLeads(ctx, &crm_grpc_adapter.ListAgentLeadsParams{
		AgentId:      agentId,
		StatusFilter: statusFilter,
		Page:         page,
		PageSize:     pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) SetLeadReminder(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SetLeadReminder"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		LeadId   string `json:"lead_id"`
		AgentId  string `json:"agent_id"`
		RemindAt string `json:"remind_at"`
		Message  string `json:"message"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SetLeadReminder(ctx, &crm_grpc_adapter.SetLeadReminderParams{
		LeadId:   body.LeadId,
		AgentId:  body.AgentId,
		RemindAt: body.RemindAt,
		Message:  body.Message,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ListLeadReminders(c *fiber.Ctx, leadId string) error {
	const op = "rest_oapi.Server.ListLeadReminders"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.ListLeadReminders(ctx, leadId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) FilterBotLeads(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.FilterBotLeads"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.FilterBotLeads(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) CreateDripSequence(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateDripSequence"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId string   `json:"agent_id"`
		Name    string   `json:"name"`
		Steps   []string `json:"steps"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateDripSequence(ctx, &crm_grpc_adapter.CreateDripSequenceParams{
		AgentId: body.AgentId,
		Name:    body.Name,
		Steps:   body.Steps,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ListDripSequences(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.ListDripSequences"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListDripSequences(ctx, &crm_grpc_adapter.ListDripSequencesParams{
		AgentId:  agentId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) CreateMomentTrigger(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateMomentTrigger"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId   string `json:"agent_id"`
		EventType string `json:"event_type"`
		Action    string `json:"action"`
		Condition string `json:"condition"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateMomentTrigger(ctx, &crm_grpc_adapter.CreateMomentTriggerParams{
		AgentId:   body.AgentId,
		EventType: body.EventType,
		Action:    body.Action,
		Condition: body.Condition,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) CreateSegment(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateSegment"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId  string   `json:"agent_id"`
		Name     string   `json:"name"`
		Criteria []string `json:"criteria"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateSegment(ctx, &crm_grpc_adapter.CreateSegmentParams{
		AgentId:  body.AgentId,
		Name:     body.Name,
		Criteria: body.Criteria,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ListSegments(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.ListSegments"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListSegments(ctx, &crm_grpc_adapter.ListSegmentsParams{
		AgentId:  agentId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) SendBroadcast(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SendBroadcast"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId    string   `json:"agent_id"`
		SegmentIds []string `json:"segment_ids"`
		Message    string   `json:"message"`
		Channel    string   `json:"channel"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SendBroadcast(ctx, &crm_grpc_adapter.SendBroadcastParams{
		AgentId:    body.AgentId,
		SegmentIds: body.SegmentIds,
		Message:    body.Message,
		Channel:    body.Channel,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) AssignLeadFair(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.AssignLeadFair"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId string   `json:"agent_id"`
		LeadIds []string `json:"lead_ids"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.AssignLeadFair(ctx, &crm_grpc_adapter.AssignLeadFairParams{
		AgentId: body.AgentId,
		LeadIds: body.LeadIds,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) CreateSlaRule(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateSlaRule"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId          string `json:"agent_id"`
		RuleName         string `json:"rule_name"`
		ResponseTimeMins int32  `json:"response_time_mins"`
		EscalationLevel  int32  `json:"escalation_level"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateSlaRule(ctx, &crm_grpc_adapter.CreateSlaRuleParams{
		AgentId:          body.AgentId,
		RuleName:         body.RuleName,
		ResponseTimeMins: body.ResponseTimeMins,
		EscalationLevel:  body.EscalationLevel,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) GetLeadTrail(c *fiber.Ctx, leadId string) error {
	const op = "rest_oapi.Server.GetLeadTrail"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetLeadTrail(ctx, leadId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) TagLead(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.TagLead"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		LeadId string   `json:"lead_id"`
		Tags   []string `json:"tags"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.TagLead(ctx, &crm_grpc_adapter.TagLeadParams{
		LeadId: body.LeadId,
		Tags:   body.Tags,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GenerateQuote(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GenerateQuote"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		LeadId    string  `json:"lead_id"`
		PackageId string  `json:"package_id"`
		Pax       int32   `json:"pax"`
		Discount  float64 `json:"discount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.GenerateQuote(ctx, &crm_grpc_adapter.GenerateQuoteParams{
		LeadId:    body.LeadId,
		PackageId: body.PackageId,
		Pax:       body.Pax,
		Discount:  body.Discount,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) GetQuote(c *fiber.Ctx, quoteId string) error {
	const op = "rest_oapi.Server.GetQuote"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetQuote(ctx, quoteId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) BuildPaymentLink(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.BuildPaymentLink"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		QuoteId string `json:"quote_id"`
		LeadId  string `json:"lead_id"`
		Notes   string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.BuildPaymentLink(ctx, &crm_grpc_adapter.BuildPaymentLinkParams{
		QuoteId: body.QuoteId,
		LeadId:  body.LeadId,
		Notes:   body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) RequestDiscount(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RequestDiscount"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		QuoteId string  `json:"quote_id"`
		AgentId string  `json:"agent_id"`
		Amount  float64 `json:"amount"`
		Reason  string  `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RequestDiscount(ctx, &crm_grpc_adapter.RequestDiscountParams{
		QuoteId: body.QuoteId,
		AgentId: body.AgentId,
		Amount:  body.Amount,
		Reason:  body.Reason,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) ApproveDiscount(c *fiber.Ctx, discountId string) error {
	const op = "rest_oapi.Server.ApproveDiscount"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		Decision string `json:"decision"`
		Notes    string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.ApproveDiscount(ctx, &crm_grpc_adapter.ApproveDiscountParams{
		DiscountId: discountId,
		Decision:   body.Decision,
		Notes:      body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetStaleProspects(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetStaleProspects"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	staleDays := int32(c.QueryInt("stale_days", 30))
	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.GetStaleProspects(ctx, &crm_grpc_adapter.GetStaleProspectsParams{
		AgentId:   agentId,
		StaleDays: staleDays,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// Commission
// ---------------------------------------------------------------------------

func (s *Server) GetCommissionBalance(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetCommissionBalance"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetCommissionBalance(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetCommissionEvents(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetCommissionEvents"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetCommissionEvents(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) RequestPayout(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RequestPayout"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId   string `json:"agent_id"`
		AmountIdr int64  `json:"amount_idr"`
		BankCode  string `json:"bank_code"`
		AccountNo string `json:"account_no"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RequestPayout(ctx, &crm_grpc_adapter.RequestPayoutParams{
		AgentId:   body.AgentId,
		AmountIdr: body.AmountIdr,
		BankCode:  body.BankCode,
		AccountNo: body.AccountNo,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (s *Server) GetPayoutHistory(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetPayoutHistory"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetPayoutHistory(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) ComputeOverrideCommission(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ComputeOverrideCommission"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId    string `json:"agent_id"`
		DownlineId string `json:"downline_id"`
		SaleAmount int64  `json:"sale_amount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.ComputeOverrideCommission(ctx, &crm_grpc_adapter.ComputeOverrideParams{
		AgentId:    body.AgentId,
		DownlineId: body.DownlineId,
		SaleAmount: body.SaleAmount,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetRoasReport(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetRoasReport"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	result, err := s.svc.GetRoasReport(ctx, &crm_grpc_adapter.GetRoasReportParams{
		AgentId:  agentId,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// Academy
// ---------------------------------------------------------------------------

func (s *Server) ListAcademyCourses(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListAcademyCourses"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListAcademyCourses(ctx, &crm_grpc_adapter.ListAcademyCoursesParams{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetCourseProgress(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetCourseProgress"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	courseId := c.Query("course_id")

	result, err := s.svc.GetCourseProgress(ctx, &crm_grpc_adapter.GetCourseProgressParams{
		AgentId:  agentId,
		CourseId: courseId,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) SubmitQuiz(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SubmitQuiz"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId  string   `json:"agent_id"`
		CourseId string   `json:"course_id"`
		Answers  []string `json:"answers"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SubmitQuiz(ctx, &crm_grpc_adapter.SubmitQuizParams{
		AgentId:  body.AgentId,
		CourseId: body.CourseId,
		Answers:  body.Answers,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) ListSalesScripts(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListSalesScripts"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	category := c.Query("category")
	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.ListSalesScripts(ctx, &crm_grpc_adapter.ListSalesScriptsParams{
		Category: category,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetLeaderboard(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetLeaderboard"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	period := c.Query("period")
	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))

	result, err := s.svc.GetLeaderboard(ctx, &crm_grpc_adapter.GetLeaderboardParams{
		Period:   period,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetAgentTier(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetAgentTier"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetAgentTier(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// Alumni / Other
// ---------------------------------------------------------------------------

func (s *Server) GetAlumniReferrals(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetAlumniReferrals"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetAlumniReferrals(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) GetReturnIntentSavings(c *fiber.Ctx, agentId string) error {
	const op = "rest_oapi.Server.GetReturnIntentSavings"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetReturnIntentSavings(ctx, agentId)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) CalculateZakat(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CalculateZakat"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		IncomeSources    []string `json:"income_sources"`
		TotalAssets      int64    `json:"total_assets"`
		TotalLiabilities int64    `json:"total_liabilities"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CalculateZakat(ctx, &crm_grpc_adapter.CalculateZakatParams{
		IncomeSources:    body.IncomeSources,
		TotalAssets:      body.TotalAssets,
		TotalLiabilities: body.TotalLiabilities,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (s *Server) RecordCharity(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordCharity"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AgentId    string `json:"agent_id"`
		AmountIdr  int64  `json:"amount_idr"`
		CharityOrg string `json:"charity_org"`
		Notes      string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordCharity(ctx, &crm_grpc_adapter.RecordCharityParams{
		AgentId:    body.AgentId,
		AmountIdr:  body.AmountIdr,
		CharityOrg: body.CharityOrg,
		Notes:      body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}
