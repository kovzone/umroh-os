// admin_phase6.go — gateway-svc adapter methods for iam-svc Phase 6 admin RPCs
// (BL-IAM-007/011/014/016): SetDataScope, CreateAPIKey, RevokeAPIKey,
// GetGlobalConfig, SetGlobalConfig, SearchActivityLog.

package iam_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// SetDataScope
// ---------------------------------------------------------------------------

// SetDataScopeParams holds inputs for PUT /v1/admin/users/:id/data-scope.
type SetDataScopeParams struct {
	UserID    string
	ScopeType string
	BranchID  string
}

// SetDataScopeResult is the response.
type SetDataScopeResult struct {
	UserID    string
	ScopeType string
}

func (a *Adapter) SetDataScope(ctx context.Context, params *SetDataScopeParams) (*SetDataScopeResult, error) {
	const op = "iam_grpc_adapter.Adapter.SetDataScope"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	resp, err := a.adminPhase6Client.SetDataScope(ctx, &pb.SetDataScopeRequest{
		UserId:    params.UserID,
		ScopeType: params.ScopeType,
		BranchId:  params.BranchID,
	})
	if err != nil {
		wrapped := mapIamAdminError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.SetDataScope failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &SetDataScopeResult{
		UserID:    resp.GetUserId(),
		ScopeType: resp.GetScopeType(),
	}, nil
}

// ---------------------------------------------------------------------------
// CreateAPIKey
// ---------------------------------------------------------------------------

// CreateAPIKeyParams holds inputs for POST /v1/admin/api-keys.
type CreateAPIKeyParams struct {
	Name      string
	Scopes    []string
	ExpiresAt string
	CreatedBy string
}

// CreateAPIKeyResult holds the response (plaintext key returned exactly once).
type CreateAPIKeyResult struct {
	KeyID        string
	PlaintextKey string
	KeyPrefix    string
	ExpiresAt    string
}

func (a *Adapter) CreateAPIKey(ctx context.Context, params *CreateAPIKeyParams) (*CreateAPIKeyResult, error) {
	const op = "iam_grpc_adapter.Adapter.CreateAPIKey"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("name", params.Name).Msg("")

	resp, err := a.adminPhase6Client.CreateAPIKey(ctx, &pb.CreateAPIKeyRequest{
		Name:      params.Name,
		Scopes:    params.Scopes,
		ExpiresAt: params.ExpiresAt,
		CreatedBy: params.CreatedBy,
	})
	if err != nil {
		wrapped := mapIamAdminError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.CreateAPIKey failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &CreateAPIKeyResult{
		KeyID:        resp.GetKeyId(),
		PlaintextKey: resp.GetPlaintextKey(),
		KeyPrefix:    resp.GetKeyPrefix(),
		ExpiresAt:    resp.GetExpiresAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// RevokeAPIKey
// ---------------------------------------------------------------------------

// RevokeAPIKeyResult holds the response for DELETE /v1/admin/api-keys/:id.
type RevokeAPIKeyResult struct {
	KeyID     string
	RevokedAt string
}

func (a *Adapter) RevokeAPIKey(ctx context.Context, keyID string) (*RevokeAPIKeyResult, error) {
	const op = "iam_grpc_adapter.Adapter.RevokeAPIKey"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("key_id", keyID).Msg("")

	resp, err := a.adminPhase6Client.RevokeAPIKey(ctx, &pb.RevokeAPIKeyRequest{
		KeyId: keyID,
	})
	if err != nil {
		wrapped := mapIamAdminError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.RevokeAPIKey failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RevokeAPIKeyResult{
		KeyID:     resp.GetKeyId(),
		RevokedAt: resp.GetRevokedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// GetGlobalConfig
// ---------------------------------------------------------------------------

// ConfigEntry holds one global config entry.
type ConfigEntry struct {
	Key         string
	Value       string
	Description string
	UpdatedAt   string
}

// GetGlobalConfigResult holds the response for GET /v1/admin/config.
type GetGlobalConfigResult struct {
	Configs []*ConfigEntry
}

func (a *Adapter) GetGlobalConfig(ctx context.Context, keys []string) (*GetGlobalConfigResult, error) {
	const op = "iam_grpc_adapter.Adapter.GetGlobalConfig"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Int("key_count", len(keys)).Msg("")

	resp, err := a.adminPhase6Client.GetGlobalConfig(ctx, &pb.GetGlobalConfigRequest{
		Keys: keys,
	})
	if err != nil {
		wrapped := mapIamAdminError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.GetGlobalConfig failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	configs := make([]*ConfigEntry, 0, len(resp.GetConfigs()))
	for _, c := range resp.GetConfigs() {
		configs = append(configs, &ConfigEntry{
			Key:         c.GetKey(),
			Value:       c.GetValue(),
			Description: c.GetDescription(),
			UpdatedAt:   c.GetUpdatedAt(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetGlobalConfigResult{Configs: configs}, nil
}

// ---------------------------------------------------------------------------
// SetGlobalConfig
// ---------------------------------------------------------------------------

// SetGlobalConfigParams holds inputs for PUT /v1/admin/config/:key.
type SetGlobalConfigParams struct {
	Key         string
	Value       string
	Description string
	UpdatedBy   string
}

// SetGlobalConfigResult holds the response.
type SetGlobalConfigResult struct {
	Key       string
	Value     string
	UpdatedAt string
}

func (a *Adapter) SetGlobalConfig(ctx context.Context, params *SetGlobalConfigParams) (*SetGlobalConfigResult, error) {
	const op = "iam_grpc_adapter.Adapter.SetGlobalConfig"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("key", params.Key).Msg("")

	resp, err := a.adminPhase6Client.SetGlobalConfig(ctx, &pb.SetGlobalConfigRequest{
		Key:         params.Key,
		Value:       params.Value,
		Description: params.Description,
		UpdatedBy:   params.UpdatedBy,
	})
	if err != nil {
		wrapped := mapIamAdminError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.SetGlobalConfig failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &SetGlobalConfigResult{
		Key:       resp.GetKey(),
		Value:     resp.GetValue(),
		UpdatedAt: resp.GetUpdatedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// SearchActivityLog
// ---------------------------------------------------------------------------

// SearchActivityLogParams holds filters for GET /v1/admin/activity-log.
type SearchActivityLogParams struct {
	UserID   string
	Resource string
	Action   string
	From     string
	To       string
	Limit    int32
	Cursor   string
}

// ActivityLogEntry holds one audit log entry.
type ActivityLogEntry struct {
	ID         string
	UserID     string
	Resource   string
	Action     string
	ResourceID string
	CreatedAt  string
}

// SearchActivityLogResult is the paginated response.
type SearchActivityLogResult struct {
	Logs       []*ActivityLogEntry
	NextCursor string
}

func (a *Adapter) SearchActivityLog(ctx context.Context, params *SearchActivityLogParams) (*SearchActivityLogResult, error) {
	const op = "iam_grpc_adapter.Adapter.SearchActivityLog"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("resource", params.Resource).Msg("")

	resp, err := a.adminPhase6Client.SearchActivityLog(ctx, &pb.SearchActivityLogRequest{
		UserId:   params.UserID,
		Resource: params.Resource,
		Action:   params.Action,
		From:     params.From,
		To:       params.To,
		Limit:    params.Limit,
		Cursor:   params.Cursor,
	})
	if err != nil {
		wrapped := mapIamAdminError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.SearchActivityLog failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	logs := make([]*ActivityLogEntry, 0, len(resp.GetLogs()))
	for _, e := range resp.GetLogs() {
		logs = append(logs, &ActivityLogEntry{
			ID:         e.GetId(),
			UserID:     e.GetUserId(),
			Resource:   e.GetResource(),
			Action:     e.GetAction(),
			ResourceID: e.GetResourceId(),
			CreatedAt:  e.GetCreatedAt(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &SearchActivityLogResult{
		Logs:       logs,
		NextCursor: resp.GetNextCursor(),
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapper for IAM admin Phase 6 calls
// ---------------------------------------------------------------------------

func mapIamAdminError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("iam admin call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	case grpcCodes.Unavailable, grpcCodes.DeadlineExceeded:
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("iam admin unreachable: %s", st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("iam admin call failed: %s", st.Message()))
	}
}
