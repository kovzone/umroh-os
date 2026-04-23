package service

import (
	"context"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// BL-IAM-018 (S1-E-12) — client-facing auth proxies. Each method is a thin
// wrapper around iam_grpc_adapter; the adapter does the gRPC call + error
// mapping. Business logic stays in iam-svc.

// Login proxies POST /v1/sessions to iam-svc.Login (gRPC).
func (s *Service) Login(ctx context.Context, params *iam_grpc_adapter.LoginParams) (*iam_grpc_adapter.LoginResult, error) {
	const op = "service.Service.Login"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("email", params.Email))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.Login(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// RefreshSession proxies POST /v1/sessions/refresh to iam-svc.RefreshSession (gRPC).
func (s *Service) RefreshSession(ctx context.Context, params *iam_grpc_adapter.RefreshSessionParams) (*iam_grpc_adapter.RefreshSessionResult, error) {
	const op = "service.Service.RefreshSession"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.RefreshSession(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// Logout proxies DELETE /v1/sessions to iam-svc.Logout (gRPC). The session_id
// comes from the validated bearer's payload (resolved by the gateway bearer
// middleware before this method runs).
func (s *Service) Logout(ctx context.Context, params *iam_grpc_adapter.LogoutParams) (*iam_grpc_adapter.LogoutResult, error) {
	const op = "service.Service.Logout"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("session_id", params.SessionID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.Logout(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// GetMe proxies GET /v1/me to iam-svc.GetMe (gRPC).
func (s *Service) GetMe(ctx context.Context, params *iam_grpc_adapter.GetMeParams) (*iam_grpc_adapter.GetMeResult, error) {
	const op = "service.Service.GetMe"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", params.UserID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.GetMe(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// EnrollTOTP proxies POST /v1/me/2fa/enroll to iam-svc.EnrollTotp (gRPC).
func (s *Service) EnrollTOTP(ctx context.Context, params *iam_grpc_adapter.EnrollTOTPParams) (*iam_grpc_adapter.EnrollTOTPResult, error) {
	const op = "service.Service.EnrollTOTP"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", params.UserID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.EnrollTOTP(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// VerifyTOTP proxies POST /v1/me/2fa/verify to iam-svc.VerifyTotp (gRPC).
func (s *Service) VerifyTOTP(ctx context.Context, params *iam_grpc_adapter.VerifyTOTPParams) (*iam_grpc_adapter.VerifyTOTPResult, error) {
	const op = "service.Service.VerifyTOTP"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", params.UserID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.VerifyTOTP(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// SuspendUser proxies POST /v1/users/{id}/suspend to iam-svc.SuspendUser (gRPC).
// iam-svc enforces the permission gate + self-suspend guard.
func (s *Service) SuspendUser(ctx context.Context, params *iam_grpc_adapter.SuspendUserParams) (*iam_grpc_adapter.SuspendUserResult, error) {
	const op = "service.Service.SuspendUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("actor_user_id", params.ActorUserID),
		attribute.String("target_user_id", params.TargetUserID),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.SuspendUser(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}
