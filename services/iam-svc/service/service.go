package service

import (
	"context"
	"time"

	"iam-svc/store/postgres_store"
	"iam-svc/util/token"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for iam-svc.
//
// S1-E-04 (BL-IAM-001) adds the first real IAM endpoints:
// internal login / refresh / logout / current-user, plus the
// TOTP enrollment + verify half-flow. Login-time TOTP enforcement
// is deferred to S1-E-06.
//
// Deferred (sibling S1-E-04 cards + S1-E-06 depth card):
//   - ValidateToken / CheckPermission gRPC (BL-IAM-002)
//   - Suspend / revoke-all (BL-IAM-003)
//   - Audit log writes (BL-IAM-004)
//   - Admin user / role / branch CRUD (S1-E-06)
type IService interface {
	// System
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Auth — BL-IAM-001 (implemented in service/auth.go).
	Login(ctx context.Context, params *LoginParams) (*LoginResult, error)
	RefreshSession(ctx context.Context, params *RefreshSessionParams) (*RefreshSessionResult, error)
	Logout(ctx context.Context, params *LogoutParams) (*LogoutResult, error)

	// Me + TOTP — BL-IAM-001 (implemented in service/me.go).
	GetMe(ctx context.Context, params *GetMeParams) (*GetMeResult, error)
	EnrollTOTP(ctx context.Context, params *EnrollTOTPParams) (*EnrollTOTPResult, error)
	VerifyTOTP(ctx context.Context, params *VerifyTOTPParams) (*VerifyTOTPResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore

	// Auth dependencies (BL-IAM-001).
	tokenMaker         token.Maker
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	totpIssuer         string
	totpEncryptionKey  []byte
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
	tokenMaker token.Maker,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	totpIssuer string,
	totpEncryptionKey []byte,
) IService {
	return &Service{
		logger:            logger,
		tracer:            tracer,
		appName:           appName,
		store:             store,
		tokenMaker:        tokenMaker,
		accessTokenTTL:    accessTokenTTL,
		refreshTokenTTL:   refreshTokenTTL,
		totpIssuer:        totpIssuer,
		totpEncryptionKey: totpEncryptionKey,
	}
}
