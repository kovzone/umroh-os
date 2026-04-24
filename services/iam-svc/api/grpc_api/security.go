// security.go — gRPC handlers for IAM security depth features.
// BL-IAM-010: GetPasswordPolicy, SetPasswordPolicy
// BL-IAM-012: RecordLoginAnomaly
// BL-IAM-013: ListSessions, RevokeSession
// BL-IAM-015: UpsertCommTemplate, ListCommTemplates
// BL-IAM-017: TriggerBackup, GetBackupHistory

package grpc_api

import (
	"context"

	"iam-svc/api/grpc_api/pb"
	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// BL-IAM-010: Password policy
// ---------------------------------------------------------------------------

func (s *Server) GetPasswordPolicy(ctx context.Context, req *pb.GetPasswordPolicyRequest) (*pb.GetPasswordPolicyResponse, error) {
	const op = "grpc_api.Server.GetPasswordPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetPasswordPolicy(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetPasswordPolicyResponse{
		MinLength:      int32(result.MinLength),
		RequireUpper:   result.RequireUpper,
		RequireDigit:   result.RequireDigit,
		RequireSpecial: result.RequireSpecial,
		RequireMfa:     result.RequireMFA,
		UpdatedAt:      result.UpdatedAt,
	}, nil
}

func (s *Server) SetPasswordPolicy(ctx context.Context, req *pb.SetPasswordPolicyRequest) (*pb.GetPasswordPolicyResponse, error) {
	const op = "grpc_api.Server.SetPasswordPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.SetPasswordPolicy(ctx, &service.SetPasswordPolicyParams{
		MinLength:      int(req.GetMinLength()),
		RequireUpper:   req.GetRequireUpper(),
		RequireDigit:   req.GetRequireDigit(),
		RequireSpecial: req.GetRequireSpecial(),
		RequireMFA:     req.GetRequireMfa(),
		UpdatedBy:      req.GetUpdatedBy(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetPasswordPolicyResponse{
		MinLength:      int32(result.MinLength),
		RequireUpper:   result.RequireUpper,
		RequireDigit:   result.RequireDigit,
		RequireSpecial: result.RequireSpecial,
		RequireMfa:     result.RequireMFA,
		UpdatedAt:      result.UpdatedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-012: Anomaly alerts
// ---------------------------------------------------------------------------

func (s *Server) RecordLoginAnomaly(ctx context.Context, req *pb.RecordLoginAnomalyRequest) (*pb.RecordLoginAnomalyResponse, error) {
	const op = "grpc_api.Server.RecordLoginAnomaly"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("user_id", req.GetUserId()),
		attribute.String("anomaly_kind", req.GetAnomalyKind()),
	)
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordLoginAnomaly(ctx, &service.RecordLoginAnomalyParams{
		UserID:      req.GetUserId(),
		IP:          req.GetIp(),
		UserAgent:   req.GetUserAgent(),
		AnomalyKind: req.GetAnomalyKind(),
		Details:     req.GetDetails(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.RecordLoginAnomalyResponse{
		AlertId:   result.AlertID,
		CreatedAt: result.CreatedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-013: Session history + revoke
// ---------------------------------------------------------------------------

func (s *Server) ListSessions(ctx context.Context, req *pb.ListSessionsRequest) (*pb.ListSessionsResponse, error) {
	const op = "grpc_api.Server.ListSessions"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", req.GetUserId()))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListSessions(ctx, &service.ListSessionsParams{
		UserID:     req.GetUserId(),
		IncludeAll: req.GetIncludeAll(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	sessions := make([]*pb.SessionEntry, 0, len(result.Sessions))
	for _, s := range result.Sessions {
		sessions = append(sessions, &pb.SessionEntry{
			SessionId: s.SessionID,
			UserAgent: s.UserAgent,
			Ip:        s.IP,
			IssuedAt:  s.IssuedAt,
			ExpiresAt: s.ExpiresAt,
			RevokedAt: s.RevokedAt,
			IsActive:  s.IsActive,
		})
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListSessionsResponse{Sessions: sessions}, nil
}

func (s *Server) RevokeSession(ctx context.Context, req *pb.RevokeSessionRequest) (*pb.RevokeSessionResponse, error) {
	const op = "grpc_api.Server.RevokeSession"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("session_id", req.GetSessionId()))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RevokeSession(ctx, &service.RevokeSessionParams{
		SessionID:   req.GetSessionId(),
		RequestorID: req.GetRequestorId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.RevokeSessionResponse{
		SessionId: result.SessionID,
		RevokedAt: result.RevokedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-015: Communication templates
// ---------------------------------------------------------------------------

func (s *Server) UpsertCommTemplate(ctx context.Context, req *pb.UpsertCommTemplateRequest) (*pb.UpsertCommTemplateResponse, error) {
	const op = "grpc_api.Server.UpsertCommTemplate"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("channel", req.GetChannel()).Str("name", req.GetName()).Msg("")

	result, err := s.svc.UpsertCommTemplate(ctx, &service.UpsertCommTemplateParams{
		Channel:   service.CommTemplateChannel(req.GetChannel()),
		Name:      req.GetName(),
		Subject:   req.GetSubject(),
		Body:      req.GetBody(),
		Variables: req.GetVariables(),
		UpdatedBy: req.GetUpdatedBy(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpsertCommTemplateResponse{
		Key:       result.Key,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (s *Server) ListCommTemplates(ctx context.Context, req *pb.ListCommTemplatesRequest) (*pb.ListCommTemplatesResponse, error) {
	const op = "grpc_api.Server.ListCommTemplates"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListCommTemplates(ctx, &service.ListCommTemplatesParams{
		Channel: service.CommTemplateChannel(req.GetChannel()),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	templates := make([]*pb.CommTemplate, 0, len(result.Templates))
	for _, t := range result.Templates {
		templates = append(templates, &pb.CommTemplate{
			Key:       t.Key,
			Channel:   string(t.Channel),
			Name:      t.Name,
			Subject:   t.Subject,
			Body:      t.Body,
			Variables: t.Variables,
			UpdatedAt: t.UpdatedAt,
		})
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListCommTemplatesResponse{Templates: templates}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-017: DB backup schedule/status
// ---------------------------------------------------------------------------

func (s *Server) TriggerBackup(ctx context.Context, req *pb.TriggerBackupRequest) (*pb.TriggerBackupResponse, error) {
	const op = "grpc_api.Server.TriggerBackup"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("triggered_by", req.GetTriggeredBy()))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.TriggerBackup(ctx, &service.TriggerBackupParams{
		TriggeredBy: req.GetTriggeredBy(),
		Label:       req.GetLabel(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.TriggerBackupResponse{
		BackupId:    result.BackupID,
		Status:      string(result.Status),
		ScheduledAt: result.ScheduledAt,
	}, nil
}

func (s *Server) GetBackupHistory(ctx context.Context, req *pb.GetBackupHistoryRequest) (*pb.GetBackupHistoryResponse, error) {
	const op = "grpc_api.Server.GetBackupHistory"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetBackupHistory(ctx, req.GetLimit())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	backups := make([]*pb.BackupEntry, 0, len(result.Backups))
	for _, b := range result.Backups {
		backups = append(backups, &pb.BackupEntry{
			BackupId:    b.BackupID,
			Status:      string(b.Status),
			TriggeredBy: b.TriggeredBy,
			ScheduledAt: b.ScheduledAt,
		})
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetBackupHistoryResponse{Backups: backups}, nil
}
