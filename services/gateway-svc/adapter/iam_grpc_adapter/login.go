package iam_grpc_adapter

import (
	"context"
	"time"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// LoginParams is the adapter-level input carried from the gateway REST handler
// through to iam-svc over gRPC.
type LoginParams struct {
	Email     string
	Password  string
	TOTPCode  string // optional; accepted-but-unused until S1-E-06
	UserAgent string // forwarded from Fiber c.Get("User-Agent")
	IP        string // forwarded from the Fiber X-Forwarded-For/c.IP() derivation
}

// UserProfile is the adapter-level user envelope. Keeps proto types out of
// gateway handlers.
type UserProfile struct {
	UserID   string
	Email    string
	Name     string
	BranchID string
	Status   string
}

// LoginResult is the adapter-level output.
type LoginResult struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	User             UserProfile
}

// Login delegates to iam-svc.Login over gRPC. Transport / auth errors flow
// through mapIamError — invalid credentials arrive as ErrUnauthorized, a
// suspended-account gate as ErrForbidden, and iam-svc being unreachable as
// ErrServiceUnavailable (502).
func (a *Adapter) Login(ctx context.Context, params *LoginParams) (*LoginResult, error) {
	const op = "iam_grpc_adapter.Adapter.Login"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "Login"),
		attribute.String("email", params.Email),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	// Do NOT log params.Password.
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
		AccessExpiresAt:  time.Unix(resp.GetAccessExpiresAtUnix(), 0),
		RefreshExpiresAt: time.Unix(resp.GetRefreshExpiresAtUnix(), 0),
		User:             userProfileFromProto(resp.GetUser()),
	}, nil
}

// userProfileFromProto strips the proto wrapper from the user envelope.
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
