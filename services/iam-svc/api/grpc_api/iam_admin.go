// iam_admin.go — gRPC handlers for IAM Phase 6 admin/security RPCs.
// SetDataScope, CreateAPIKey, RevokeAPIKey, GetGlobalConfig, SetGlobalConfig,
// SearchActivityLog (BL-IAM-007/010/011/014/016).

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
// SetDataScope (BL-IAM-007)
// ---------------------------------------------------------------------------

func (s *Server) SetDataScope(ctx context.Context, req *pb.SetDataScopeRequest) (*pb.SetDataScopeResponse, error) {
	const op = "grpc_api.Server.SetDataScope"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("user_id", req.GetUserId()),
		attribute.String("scope_type", req.GetScopeType()),
	)
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("user_id", req.GetUserId()).Msg("")

	result, err := s.svc.SetDataScope(ctx, &service.SetDataScopeParams{
		UserID:    req.GetUserId(),
		ScopeType: req.GetScopeType(),
		BranchID:  req.GetBranchId(),
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.SetDataScopeResponse{
		UserId:    result.UserID,
		ScopeType: result.ScopeType,
	}, nil
}

// ---------------------------------------------------------------------------
// CreateAPIKey (BL-IAM-014)
// ---------------------------------------------------------------------------

func (s *Server) CreateAPIKey(ctx context.Context, req *pb.CreateAPIKeyRequest) (*pb.CreateAPIKeyResponse, error) {
	const op = "grpc_api.Server.CreateAPIKey"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("name", req.GetName()))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("name", req.GetName()).Msg("")

	result, err := s.svc.CreateAPIKey(ctx, &service.CreateAPIKeyParams{
		Name:      req.GetName(),
		Scopes:    req.GetScopes(),
		ExpiresAt: req.GetExpiresAt(),
		CreatedBy: req.GetCreatedBy(),
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetAttributes(attribute.String("key_id", result.KeyID))
	span.SetStatus(codes.Ok, "success")
	return &pb.CreateAPIKeyResponse{
		KeyId:        result.KeyID,
		PlaintextKey: result.PlaintextKey,
		KeyPrefix:    result.KeyPrefix,
		ExpiresAt:    result.ExpiresAt,
	}, nil
}

// ---------------------------------------------------------------------------
// RevokeAPIKey (BL-IAM-014)
// ---------------------------------------------------------------------------

func (s *Server) RevokeAPIKey(ctx context.Context, req *pb.RevokeAPIKeyRequest) (*pb.RevokeAPIKeyResponse, error) {
	const op = "grpc_api.Server.RevokeAPIKey"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("key_id", req.GetKeyId()))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("key_id", req.GetKeyId()).Msg("")

	result, err := s.svc.RevokeAPIKey(ctx, &service.RevokeAPIKeyParams{
		KeyID: req.GetKeyId(),
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.RevokeAPIKeyResponse{
		KeyId:     result.KeyID,
		RevokedAt: result.RevokedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// GetGlobalConfig (BL-IAM-016)
// ---------------------------------------------------------------------------

func (s *Server) GetGlobalConfig(ctx context.Context, req *pb.GetGlobalConfigRequest) (*pb.GetGlobalConfigResponse, error) {
	const op = "grpc_api.Server.GetGlobalConfig"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Int("key_count", len(req.GetKeys())).Msg("")

	result, err := s.svc.GetGlobalConfig(ctx, &service.GetGlobalConfigParams{
		Keys: req.GetKeys(),
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	resp := &pb.GetGlobalConfigResponse{}
	for _, c := range result.Configs {
		resp.Configs = append(resp.Configs, &pb.ConfigEntryPb{
			Key:         c.Key,
			Value:       c.Value,
			Description: c.Description,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	span.SetStatus(codes.Ok, "success")
	return resp, nil
}

// ---------------------------------------------------------------------------
// SetGlobalConfig (BL-IAM-016)
// ---------------------------------------------------------------------------

func (s *Server) SetGlobalConfig(ctx context.Context, req *pb.SetGlobalConfigRequest) (*pb.SetGlobalConfigResponse, error) {
	const op = "grpc_api.Server.SetGlobalConfig"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("key", req.GetKey()))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("key", req.GetKey()).Msg("")

	result, err := s.svc.SetGlobalConfig(ctx, &service.SetGlobalConfigParams{
		Key:         req.GetKey(),
		Value:       req.GetValue(),
		Description: req.GetDescription(),
		UpdatedBy:   req.GetUpdatedBy(),
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.SetGlobalConfigResponse{
		Key:       result.Key,
		Value:     result.Value,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// SearchActivityLog (BL-IAM-011)
// ---------------------------------------------------------------------------

func (s *Server) SearchActivityLog(ctx context.Context, req *pb.SearchActivityLogRequest) (*pb.SearchActivityLogResponse, error) {
	const op = "grpc_api.Server.SearchActivityLog"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("user_id", req.GetUserId()),
		attribute.String("resource", req.GetResource()),
	)
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("resource", req.GetResource()).Msg("")

	result, err := s.svc.SearchActivityLog(ctx, &service.SearchActivityLogParams{
		UserID:   req.GetUserId(),
		Resource: req.GetResource(),
		Action:   req.GetAction(),
		From:     req.GetFrom(),
		To:       req.GetTo(),
		Limit:    req.GetLimit(),
		Cursor:   req.GetCursor(),
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	resp := &pb.SearchActivityLogResponse{
		NextCursor: result.NextCursor,
	}
	for _, e := range result.Logs {
		resp.Logs = append(resp.Logs, &pb.ActivityLogEntryPb{
			Id:         e.ID,
			UserId:     e.UserID,
			Resource:   e.Resource,
			Action:     e.Action,
			ResourceId: e.ResourceID,
			CreatedAt:  e.CreatedAt,
		})
	}

	span.SetStatus(codes.Ok, "success")
	return resp, nil
}
