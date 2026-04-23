// iam_admin_phase6_stub.go — gateway-side gRPC client stubs for iam-svc
// Phase 6 admin/security RPCs (BL-IAM-007/011/014/016).

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	IamAdminService_SetDataScope_FullMethodName      = "/pb.IamAdminService/SetDataScope"
	IamAdminService_CreateAPIKey_FullMethodName      = "/pb.IamAdminService/CreateAPIKey"
	IamAdminService_RevokeAPIKey_FullMethodName      = "/pb.IamAdminService/RevokeAPIKey"
	IamAdminService_GetGlobalConfig_FullMethodName   = "/pb.IamAdminService/GetGlobalConfig"
	IamAdminService_SetGlobalConfig_FullMethodName   = "/pb.IamAdminService/SetGlobalConfig"
	IamAdminService_SearchActivityLog_FullMethodName = "/pb.IamAdminService/SearchActivityLog"
)

// ---------------------------------------------------------------------------
// Message types — mirror iam-svc pb/iam_admin_messages.go
// ---------------------------------------------------------------------------

type SetDataScopeRequest struct {
	UserId    string
	ScopeType string
	BranchId  string
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

type CreateAPIKeyRequest struct {
	Name      string
	Scopes    []string
	ExpiresAt string
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
	PlaintextKey string
	KeyPrefix    string
	ExpiresAt    string
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
	RevokedAt string
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

type GetGlobalConfigRequest struct {
	Keys []string
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
	UpdatedAt   string
}

func (x *ConfigEntryPb) GetKey() string {
	if x == nil {
		return ""
	}
	return x.Key
}
func (x *ConfigEntryPb) GetValue() string {
	if x == nil {
		return ""
	}
	return x.Value
}
func (x *ConfigEntryPb) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *ConfigEntryPb) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
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

type SetGlobalConfigRequest struct {
	Key         string
	Value       string
	Description string
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
	UpdatedAt string
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

type SearchActivityLogRequest struct {
	UserId   string
	Resource string
	Action   string
	From     string
	To       string
	Limit    int32
	Cursor   string
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
	CreatedAt  string
}

func (x *ActivityLogEntryPb) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *ActivityLogEntryPb) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *ActivityLogEntryPb) GetResource() string {
	if x == nil {
		return ""
	}
	return x.Resource
}
func (x *ActivityLogEntryPb) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *ActivityLogEntryPb) GetResourceId() string {
	if x == nil {
		return ""
	}
	return x.ResourceId
}
func (x *ActivityLogEntryPb) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
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

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

// IamAdminPhase6Client is the Phase 6 security/config IAM admin client interface.
type IamAdminPhase6Client interface {
	SetDataScope(ctx context.Context, in *SetDataScopeRequest, opts ...grpc.CallOption) (*SetDataScopeResponse, error)
	CreateAPIKey(ctx context.Context, in *CreateAPIKeyRequest, opts ...grpc.CallOption) (*CreateAPIKeyResponse, error)
	RevokeAPIKey(ctx context.Context, in *RevokeAPIKeyRequest, opts ...grpc.CallOption) (*RevokeAPIKeyResponse, error)
	GetGlobalConfig(ctx context.Context, in *GetGlobalConfigRequest, opts ...grpc.CallOption) (*GetGlobalConfigResponse, error)
	SetGlobalConfig(ctx context.Context, in *SetGlobalConfigRequest, opts ...grpc.CallOption) (*SetGlobalConfigResponse, error)
	SearchActivityLog(ctx context.Context, in *SearchActivityLogRequest, opts ...grpc.CallOption) (*SearchActivityLogResponse, error)
}

type iamAdminPhase6Client struct {
	cc grpc.ClientConnInterface
}

func NewIamAdminPhase6Client(cc grpc.ClientConnInterface) IamAdminPhase6Client {
	return &iamAdminPhase6Client{cc}
}

func (c *iamAdminPhase6Client) SetDataScope(ctx context.Context, in *SetDataScopeRequest, opts ...grpc.CallOption) (*SetDataScopeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetDataScopeResponse)
	if err := c.cc.Invoke(ctx, IamAdminService_SetDataScope_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminPhase6Client) CreateAPIKey(ctx context.Context, in *CreateAPIKeyRequest, opts ...grpc.CallOption) (*CreateAPIKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAPIKeyResponse)
	if err := c.cc.Invoke(ctx, IamAdminService_CreateAPIKey_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminPhase6Client) RevokeAPIKey(ctx context.Context, in *RevokeAPIKeyRequest, opts ...grpc.CallOption) (*RevokeAPIKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RevokeAPIKeyResponse)
	if err := c.cc.Invoke(ctx, IamAdminService_RevokeAPIKey_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminPhase6Client) GetGlobalConfig(ctx context.Context, in *GetGlobalConfigRequest, opts ...grpc.CallOption) (*GetGlobalConfigResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGlobalConfigResponse)
	if err := c.cc.Invoke(ctx, IamAdminService_GetGlobalConfig_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminPhase6Client) SetGlobalConfig(ctx context.Context, in *SetGlobalConfigRequest, opts ...grpc.CallOption) (*SetGlobalConfigResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetGlobalConfigResponse)
	if err := c.cc.Invoke(ctx, IamAdminService_SetGlobalConfig_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminPhase6Client) SearchActivityLog(ctx context.Context, in *SearchActivityLogRequest, opts ...grpc.CallOption) (*SearchActivityLogResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchActivityLogResponse)
	if err := c.cc.Invoke(ctx, IamAdminService_SearchActivityLog_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
