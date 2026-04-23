package service

// proxy_iam.go — service layer methods for IAM proxy routes.
//
// S1-E-12 / BL-IAM-018: moves all client-facing IAM auth routes from iam-svc
// REST to gateway-svc REST → iam-svc gRPC (ADR 0009).
//
// GetIamSystemLive and GetIamSystemDbTxDiagnostic stay until the iam-svc REST
// surface is fully retired. Once S1-E-12 removes iam-svc HTTP, those two methods
// will also be removed and the iam_rest_adapter import can be dropped.

import (
	"context"
	"fmt"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/adapter/iam_rest_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CheckPermission delegates to iam-svc.CheckPermission over gRPC.
// Used by gateway handlers (e.g. SuspendUser) to gate on a specific permission
// before forwarding the request to iam-svc.
func (s *Service) CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error) {
	const op = "service.Service.CheckPermission"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("resource", params.Resource),
		attribute.String("action", params.Action),
		attribute.String("scope", params.Scope),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.iamGrpc.CheckPermission(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// GetIamSystemLive proxies iam-svc's liveness probe through the REST adapter.
// Retired once iam-svc drops its REST port (BL-IAM-018 / S1-E-12).
func (s *Service) GetIamSystemLive(ctx context.Context) (*iam_rest_adapter.LivenessResult, error) {
	const op = "service.Service.GetIamSystemLive"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op))
	logger.Info().Str("op", op).Msg("")

	result, err := s.adapters.iamRest.GetSystemLive(ctx)
	if err != nil {
		err = fmt.Errorf("call iam-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// GetIamSystemDbTxDiagnostic proxies the iam-svc DB transaction diagnostic.
// Its purpose is verification: gateway → iam produces one trace spanning both
// services, with matching trace_id log lines in each container. See S0-J-05.
// Retired once iam-svc drops its REST port.
func (s *Service) GetIamSystemDbTxDiagnostic(ctx context.Context, message string) (*iam_rest_adapter.DbTxDiagnosticResult, error) {
	const op = "service.Service.GetIamSystemDbTxDiagnostic"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("message", message),
	)
	logger.Info().Str("op", op).Str("message", message).Msg("")

	result, err := s.adapters.iamRest.GetSystemDbTxDiagnostic(ctx, message)
	if err != nil {
		err = fmt.Errorf("call iam-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// Login proxies POST /v1/auth/login to iam-svc.Login over gRPC.
func (s *Service) Login(ctx context.Context, params *iam_grpc_adapter.LoginParams) (*iam_grpc_adapter.LoginResult, error) {
	const op = "service.Service.Login"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

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

// RefreshSession proxies POST /v1/auth/refresh to iam-svc.RefreshSession over gRPC.
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

// Logout proxies DELETE /v1/auth/logout to iam-svc.Logout over gRPC.
func (s *Service) Logout(ctx context.Context, params *iam_grpc_adapter.LogoutParams) error {
	const op = "service.Service.Logout"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	err := s.adapters.iamGrpc.Logout(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

// GetMe proxies GET /v1/me to iam-svc.GetMe over gRPC.
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

// EnrollTOTP proxies POST /v1/me/2fa/enroll to iam-svc.EnrollTOTP over gRPC.
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

// VerifyTOTP proxies POST /v1/me/2fa/verify to iam-svc.VerifyTOTP over gRPC.
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

// SuspendUser proxies POST /v1/users/{id}/suspend to iam-svc.SuspendUser over gRPC.
// Permission-gating (iam.users/suspend/global) is done in the REST handler layer.
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
