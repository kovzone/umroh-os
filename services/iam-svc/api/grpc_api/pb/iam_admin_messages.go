// iam_admin_messages.go — hand-written proto message types for IAM Phase 6
// admin RPCs (BL-IAM-007/010/011/014/016).

package pb

// ---------------------------------------------------------------------------
// SetDataScope (BL-IAM-007)
// ---------------------------------------------------------------------------

type SetDataScopeRequest struct {
	UserId    string
	ScopeType string // "global" | "branch" | "own_only"
	BranchId  string // required when ScopeType = "branch"
}

func (x *SetDataScopeRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *SetDataScopeRequest) GetScopeType() string {
	if x == nil {
		return ""
	}
	return x.ScopeType
}
func (x *SetDataScopeRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}

type SetDataScopeResponse struct {
	UserId    string
	ScopeType string
}

func (x *SetDataScopeResponse) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *SetDataScopeResponse) GetScopeType() string {
	if x == nil {
		return ""
	}
	return x.ScopeType
}

// ---------------------------------------------------------------------------
// CreateAPIKey (BL-IAM-014)
// ---------------------------------------------------------------------------

type CreateAPIKeyRequest struct {
	Name      string
	Scopes    []string
	ExpiresAt string // RFC3339; "" = no expiry
	CreatedBy string
}

func (x *CreateAPIKeyRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateAPIKeyRequest) GetScopes() []string {
	if x == nil {
		return nil
	}
	return x.Scopes
}
func (x *CreateAPIKeyRequest) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}
func (x *CreateAPIKeyRequest) GetCreatedBy() string {
	if x == nil {
		return ""
	}
	return x.CreatedBy
}

type CreateAPIKeyResponse struct {
	KeyId        string
	PlaintextKey string // returned exactly once
	KeyPrefix    string
	ExpiresAt    string // RFC3339 or ""
}

func (x *CreateAPIKeyResponse) GetKeyId() string {
	if x == nil {
		return ""
	}
	return x.KeyId
}
func (x *CreateAPIKeyResponse) GetPlaintextKey() string {
	if x == nil {
		return ""
	}
	return x.PlaintextKey
}
func (x *CreateAPIKeyResponse) GetKeyPrefix() string {
	if x == nil {
		return ""
	}
	return x.KeyPrefix
}
func (x *CreateAPIKeyResponse) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}

// ---------------------------------------------------------------------------
// RevokeAPIKey (BL-IAM-014)
// ---------------------------------------------------------------------------

type RevokeAPIKeyRequest struct {
	KeyId string
}

func (x *RevokeAPIKeyRequest) GetKeyId() string {
	if x == nil {
		return ""
	}
	return x.KeyId
}

type RevokeAPIKeyResponse struct {
	KeyId     string
	RevokedAt string // RFC3339
}

func (x *RevokeAPIKeyResponse) GetKeyId() string {
	if x == nil {
		return ""
	}
	return x.KeyId
}
func (x *RevokeAPIKeyResponse) GetRevokedAt() string {
	if x == nil {
		return ""
	}
	return x.RevokedAt
}

// ---------------------------------------------------------------------------
// GetGlobalConfig (BL-IAM-016)
// ---------------------------------------------------------------------------

type GetGlobalConfigRequest struct {
	Keys []string // empty = return all
}

func (x *GetGlobalConfigRequest) GetKeys() []string {
	if x == nil {
		return nil
	}
	return x.Keys
}

type ConfigEntryPb struct {
	Key         string
	Value       string
	Description string
	UpdatedAt   string // RFC3339
}

type GetGlobalConfigResponse struct {
	Configs []*ConfigEntryPb
}

func (x *GetGlobalConfigResponse) GetConfigs() []*ConfigEntryPb {
	if x == nil {
		return nil
	}
	return x.Configs
}

// ---------------------------------------------------------------------------
// SetGlobalConfig (BL-IAM-016)
// ---------------------------------------------------------------------------

type SetGlobalConfigRequest struct {
	Key         string
	Value       string
	Description string // optional
	UpdatedBy   string
}

func (x *SetGlobalConfigRequest) GetKey() string {
	if x == nil {
		return ""
	}
	return x.Key
}
func (x *SetGlobalConfigRequest) GetValue() string {
	if x == nil {
		return ""
	}
	return x.Value
}
func (x *SetGlobalConfigRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *SetGlobalConfigRequest) GetUpdatedBy() string {
	if x == nil {
		return ""
	}
	return x.UpdatedBy
}

type SetGlobalConfigResponse struct {
	Key       string
	Value     string
	UpdatedAt string // RFC3339
}

func (x *SetGlobalConfigResponse) GetKey() string {
	if x == nil {
		return ""
	}
	return x.Key
}
func (x *SetGlobalConfigResponse) GetValue() string {
	if x == nil {
		return ""
	}
	return x.Value
}
func (x *SetGlobalConfigResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// SearchActivityLog (BL-IAM-011)
// ---------------------------------------------------------------------------

type SearchActivityLogRequest struct {
	UserId   string
	Resource string
	Action   string
	From     string // RFC3339; "" = no lower bound
	To       string // RFC3339; "" = no upper bound
	Limit    int32  // default 50; max 200
	Cursor   string // opaque; "" = first page
}

func (x *SearchActivityLogRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *SearchActivityLogRequest) GetResource() string {
	if x == nil {
		return ""
	}
	return x.Resource
}
func (x *SearchActivityLogRequest) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *SearchActivityLogRequest) GetFrom() string {
	if x == nil {
		return ""
	}
	return x.From
}
func (x *SearchActivityLogRequest) GetTo() string {
	if x == nil {
		return ""
	}
	return x.To
}
func (x *SearchActivityLogRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}
func (x *SearchActivityLogRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}

type ActivityLogEntryPb struct {
	Id         string
	UserId     string
	Resource   string
	Action     string
	ResourceId string
	CreatedAt  string // RFC3339
}

type SearchActivityLogResponse struct {
	Logs       []*ActivityLogEntryPb
	NextCursor string
}

func (x *SearchActivityLogResponse) GetLogs() []*ActivityLogEntryPb {
	if x == nil {
		return nil
	}
	return x.Logs
}
func (x *SearchActivityLogResponse) GetNextCursor() string {
	if x == nil {
		return ""
	}
	return x.NextCursor
}
