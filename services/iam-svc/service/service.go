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
// S1-E-04 (BL-IAM-001) added the first real IAM endpoints:
// internal login / refresh / logout / current-user, plus the
// TOTP enrollment + verify half-flow. Login-time TOTP enforcement
// is deferred to S1-E-06.
//
// S1-E-04 (BL-IAM-002) adds the two internal-gRPC hot-path RPCs
// every consumer service depends on: ValidateToken + CheckPermission.
//
// S1-E-04 (BL-IAM-003) adds the admin-side SuspendUser action that flips
// `iam.users.status` to `suspended` and revokes every active session in one tx.
//
// S1-E-04 (BL-IAM-004) adds the audit-producer RPC (RecordAudit) and wires
// SuspendUser's WithTx to emit its own audit row inside the same transaction.
//
// Deferred (S1-E-06 depth card):
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

	// Permission resolution — BL-IAM-002 (implemented in service/permissions.go).
	ValidateToken(ctx context.Context, params *ValidateTokenParams) (*ValidateTokenResult, error)
	CheckPermission(ctx context.Context, params *CheckPermissionParams) (*CheckPermissionResult, error)

	// Admin actions — BL-IAM-003 (implemented in service/admin.go).
	SuspendUser(ctx context.Context, params *SuspendUserParams) (*SuspendUserResult, error)

	// Audit producer — BL-IAM-004 (implemented in service/audit.go).
	RecordAudit(ctx context.Context, params *RecordAuditParams) (*RecordAuditResult, error)
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
