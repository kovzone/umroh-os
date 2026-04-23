package iam_grpc_adapter

// auth.go — gateway adapter methods for the client-facing IAM auth RPCs added
// in BL-IAM-018 / S1-E-12. These wrap iam-svc's Login, RefreshSession, Logout,
// GetMe, EnrollTOTP, VerifyTOTP, and SuspendUser gRPC methods.
//
// Each method follows the same span + log + mapIamError pattern as validate_token.go.
// Error mapping is done via mapIamError which is already defined in validate_token.go.

import (
	"context"
	"time"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local types (proto types do not leak past this package).
// ---------------------------------------------------------------------------

// UserProfile is the shared user shape returned from Login, GetMe, SuspendUser.
type UserProfile struct {
	UserID   string
	Email    string
	Name     string
	BranchID string
	Status   string // "active" | "suspended" | "pending"
}

// LoginParams is the adapter-level input for Login.
type LoginParams struct {
	Email     string
	Password  string
	TOTPCode  string
	UserAgent string
	IP        string // raw IP literal; empty → omitted in the RPC
}

// LoginResult is the adapter-level output for Login.
type LoginResult struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	User             UserProfile
}

// RefreshSessionParams is the adapter-level input for RefreshSession.
type RefreshSessionParams struct {
	RefreshToken string
	UserAgent    string
	IP           string
}

// RefreshSessionResult is the adapter-level output for RefreshSession.
type RefreshSessionResult struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

// LogoutParams carries the session_id claim extracted by the bearer middleware.
type LogoutParams struct {
	SessionID string
}

// GetMeParams carries the user_id claim from the bearer token.
type GetMeParams struct {
	UserID string
}

// GetMeResult is the adapter-level output for GetMe.
type GetMeResult struct {
	User         UserProfile
	TOTPEnrolled bool
	TOTPVerified bool
}

// EnrollTOTPParams carries the user_id from the bearer token.
type EnrollTOTPParams struct {
	UserID string
}

// EnrollTOTPResult is the adapter-level output for EnrollTOTP.
type EnrollTOTPResult struct {
	Secret     string
	OtpauthURL string
}

// VerifyTOTPParams carries the user_id and the TOTP code from the request body.
type VerifyTOTPParams struct {
	UserID string
	Code   string
}

// VerifyTOTPResult is the adapter-level output for VerifyTOTP.
type VerifyTOTPResult struct {
	VerifiedAt time.Time
}

// SuspendUserParams carries the actor (gateway-verified) and target user IDs.
type SuspendUserParams struct {
	ActorUserID  string
	TargetUserID string
}

// SuspendUserResult is the adapter-level output for SuspendUser.
type SuspendUserResult struct {
	User UserProfile
}

// ---------------------------------------------------------------------------
// RPC wrappers
// ---------------------------------------------------------------------------

// Login delegates to iam-svc.Login over gRPC.
func (a *Adapter) Login(ctx context.Context, params *LoginParams) (*LoginResult, error) {
	const op = "iam_grpc_adapter.Adapter.Login"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "Login"))
	// Do NOT log password — it is a secret. Email is acceptable for tracing.
	span.SetAttributes(attribute.String("email", params.Email))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.Login(ctx, &pb.LoginRequest{
		Email:     params.Email,
		Password:  params.Password,
		TotpCode:  params.TOTPCode,
		UserAgent: params.UserAgent,
		Ip:        params.IP,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &LoginResult{
		AccessToken:      resp.GetAccessToken(),
		RefreshToken:     resp.GetRefreshToken(),
		AccessExpiresAt:  time.Unix(resp.GetAccessExpiresAt(), 0),
		RefreshExpiresAt: time.Unix(resp.GetRefreshExpiresAt(), 0),
		User:             userProfileFromProto(resp.GetUser()),
	}, nil
}

// RefreshSession delegates to iam-svc.RefreshSession over gRPC.
func (a *Adapter) RefreshSession(ctx context.Context, params *RefreshSessionParams) (*RefreshSessionResult, error) {
	const op = "iam_grpc_adapter.Adapter.RefreshSession"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "RefreshSession"))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.RefreshSession(ctx, &pb.RefreshSessionRequest{
		RefreshToken: params.RefreshToken,
		UserAgent:    params.UserAgent,
		Ip:           params.IP,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &RefreshSessionResult{
		AccessToken:      resp.GetAccessToken(),
		RefreshToken:     resp.GetRefreshToken(),
		AccessExpiresAt:  time.Unix(resp.GetAccessExpiresAt(), 0),
		RefreshExpiresAt: time.Unix(resp.GetRefreshExpiresAt(), 0),
	}, nil
}

// Logout delegates to iam-svc.Logout over gRPC.
func (a *Adapter) Logout(ctx context.Context, params *LogoutParams) error {
	const op = "iam_grpc_adapter.Adapter.Logout"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "Logout"))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.iamClient.Logout(ctx, &pb.LogoutRequest{SessionId: params.SessionID})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

// GetMe delegates to iam-svc.GetMe over gRPC.
func (a *Adapter) GetMe(ctx context.Context, params *GetMeParams) (*GetMeResult, error) {
	const op = "iam_grpc_adapter.Adapter.GetMe"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetMe"),
		attribute.String("user_id", params.UserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.GetMe(ctx, &pb.GetMeRequest{UserId: params.UserID})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &GetMeResult{
		User:         userProfileFromProto(resp.GetUser()),
		TOTPEnrolled: resp.GetTotpEnrolled(),
		TOTPVerified: resp.GetTotpVerified(),
	}, nil
}

// EnrollTOTP delegates to iam-svc.EnrollTOTP over gRPC.
func (a *Adapter) EnrollTOTP(ctx context.Context, params *EnrollTOTPParams) (*EnrollTOTPResult, error) {
	const op = "iam_grpc_adapter.Adapter.EnrollTOTP"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "EnrollTOTP"),
		attribute.String("user_id", params.UserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.EnrollTOTP(ctx, &pb.EnrollTOTPRequest{UserId: params.UserID})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	// Do NOT log resp.Secret — it's a one-time plaintext secret.
	span.SetStatus(codes.Ok, "success")
	return &EnrollTOTPResult{
		Secret:     resp.GetSecret(),
		OtpauthURL: resp.GetOtpauthUrl(),
	}, nil
}

// VerifyTOTP delegates to iam-svc.VerifyTOTP over gRPC.
func (a *Adapter) VerifyTOTP(ctx context.Context, params *VerifyTOTPParams) (*VerifyTOTPResult, error) {
	const op = "iam_grpc_adapter.Adapter.VerifyTOTP"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "VerifyTOTP"),
		attribute.String("user_id", params.UserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.VerifyTOTP(ctx, &pb.VerifyTOTPRequest{
		UserId: params.UserID,
		Code:   params.Code,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &VerifyTOTPResult{
		VerifiedAt: time.Unix(resp.GetVerifiedAtUnix(), 0),
	}, nil
}

// SuspendUser delegates to iam-svc.SuspendUser over gRPC.
// Permission-gating (iam.users/suspend/global) is done by the gateway handler
// before calling this adapter, using CheckPermission.
func (a *Adapter) SuspendUser(ctx context.Context, params *SuspendUserParams) (*SuspendUserResult, error) {
	const op = "iam_grpc_adapter.Adapter.SuspendUser"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "SuspendUser"),
		attribute.String("actor_user_id", params.ActorUserID),
		attribute.String("target_user_id", params.TargetUserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.SuspendUser(ctx, &pb.SuspendUserRequest{
		ActorUserId:  params.ActorUserID,
		TargetUserId: params.TargetUserID,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).
			Str("actor_user_id", params.ActorUserID).
			Str("target_user_id", params.TargetUserID).
			Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &SuspendUserResult{
		User: userProfileFromProto(resp.GetUser()),
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func userProfileFromProto(u *pb.UserProfile) UserProfile {
	if u == nil {
		return UserProfile{}
	}
	return UserProfile{
		UserID:   u.GetUserId(),
		Email:    u.GetEmail(),
		Name:     u.GetName(),
		BranchID: u.GetBranchId(),
		Status:   u.GetStatus(),
	}
}
