package grpc_api

import (
	"context"
	"errors"

	"iam-svc/api/grpc_api/pb"
	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// GetMe — current-user profile. The gateway forwards user_id after bearer
// validation. Mirrors GET /v1/me.
func (s *Server) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	const op = "grpc_api.Server.GetMe"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetMe"),
		attribute.String("user_id", req.GetUserId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("user_id is required"))
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.GetMe(ctx, &service.GetMeParams{UserID: req.GetUserId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetMeResponse{
		User:          userProfileToProto(result.User),
		TotpEnrolled:  result.TOTPEnrolled,
		TotpVerified:  result.TOTPVerified,
	}, nil
}

// EnrollTotp — generate + persist TOTP secret for the caller. Mirrors
// POST /v1/me/2fa/enroll.
func (s *Server) EnrollTotp(ctx context.Context, req *pb.EnrollTotpRequest) (*pb.EnrollTotpResponse, error) {
	const op = "grpc_api.Server.EnrollTotp"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "EnrollTotp"),
		attribute.String("user_id", req.GetUserId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("user_id is required"))
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.EnrollTOTP(ctx, &service.EnrollTOTPParams{UserID: req.GetUserId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.EnrollTotpResponse{
		Secret:     result.Secret,
		OtpauthUrl: result.OtpauthURL,
	}, nil
}

// VerifyTotp — validate a six-digit TOTP code + stamp verified_at. Mirrors
// POST /v1/me/2fa/verify.
func (s *Server) VerifyTotp(ctx context.Context, req *pb.VerifyTotpRequest) (*pb.VerifyTotpResponse, error) {
	const op = "grpc_api.Server.VerifyTotp"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "VerifyTotp"),
		attribute.String("user_id", req.GetUserId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("user_id is required"))
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.VerifyTOTP(ctx, &service.VerifyTOTPParams{
		UserID: req.GetUserId(),
		Code:   req.GetCode(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.VerifyTotpResponse{
		VerifiedAtUnix: result.VerifiedAt.Unix(),
	}, nil
}
