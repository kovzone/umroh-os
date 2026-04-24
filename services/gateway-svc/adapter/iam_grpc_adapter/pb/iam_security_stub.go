// iam_security_stub.go — gateway-side gRPC client stubs for iam-svc
// security depth RPCs (BL-IAM-010/012/013/015/017).

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	IamSecurityService_GetPasswordPolicy_FullMethodName  = "/pb.IamSecurityService/GetPasswordPolicy"
	IamSecurityService_SetPasswordPolicy_FullMethodName  = "/pb.IamSecurityService/SetPasswordPolicy"
	IamSecurityService_RecordLoginAnomaly_FullMethodName = "/pb.IamSecurityService/RecordLoginAnomaly"
	IamSecurityService_ListSessions_FullMethodName       = "/pb.IamSecurityService/ListSessions"
	IamSecurityService_RevokeSession_FullMethodName      = "/pb.IamSecurityService/RevokeSession"
	IamSecurityService_UpsertCommTemplate_FullMethodName = "/pb.IamSecurityService/UpsertCommTemplate"
	IamSecurityService_ListCommTemplates_FullMethodName  = "/pb.IamSecurityService/ListCommTemplates"
	IamSecurityService_TriggerBackup_FullMethodName      = "/pb.IamSecurityService/TriggerBackup"
	IamSecurityService_GetBackupHistory_FullMethodName   = "/pb.IamSecurityService/GetBackupHistory"
)

// ---------------------------------------------------------------------------
// Message types — mirror iam-svc pb/security_messages.go
// ---------------------------------------------------------------------------

// BL-IAM-010: Password policy

type GwGetPasswordPolicyRequest struct{}

type GwGetPasswordPolicyResponse struct {
	MinLength      int32
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMfa     bool
	UpdatedAt      string
}

func (x *GwGetPasswordPolicyResponse) GetMinLength() int32 {
	if x == nil {
		return 0
	}
	return x.MinLength
}
func (x *GwGetPasswordPolicyResponse) GetRequireUpper() bool {
	if x == nil {
		return false
	}
	return x.RequireUpper
}
func (x *GwGetPasswordPolicyResponse) GetRequireDigit() bool {
	if x == nil {
		return false
	}
	return x.RequireDigit
}
func (x *GwGetPasswordPolicyResponse) GetRequireSpecial() bool {
	if x == nil {
		return false
	}
	return x.RequireSpecial
}
func (x *GwGetPasswordPolicyResponse) GetRequireMfa() bool {
	if x == nil {
		return false
	}
	return x.RequireMfa
}
func (x *GwGetPasswordPolicyResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

type GwSetPasswordPolicyRequest struct {
	MinLength      int32
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMfa     bool
	UpdatedBy      string
}

// BL-IAM-012: Login anomaly

type GwRecordLoginAnomalyRequest struct {
	UserId      string
	Ip          string
	UserAgent   string
	AnomalyKind string
	Details     string
}

type GwRecordLoginAnomalyResponse struct {
	AlertId   string
	CreatedAt string
}

func (x *GwRecordLoginAnomalyResponse) GetAlertId() string {
	if x == nil {
		return ""
	}
	return x.AlertId
}
func (x *GwRecordLoginAnomalyResponse) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

// BL-IAM-013: Sessions

type GwListSessionsRequest struct {
	UserId     string
	IncludeAll bool
}

type GwSessionEntry struct {
	SessionId string
	UserAgent string
	Ip        string
	IssuedAt  string
	ExpiresAt string
	RevokedAt string
	IsActive  bool
}

func (x *GwSessionEntry) GetSessionId() string { if x == nil { return "" }; return x.SessionId }
func (x *GwSessionEntry) GetUserAgent() string  { if x == nil { return "" }; return x.UserAgent }
func (x *GwSessionEntry) GetIp() string         { if x == nil { return "" }; return x.Ip }
func (x *GwSessionEntry) GetIssuedAt() string   { if x == nil { return "" }; return x.IssuedAt }
func (x *GwSessionEntry) GetExpiresAt() string  { if x == nil { return "" }; return x.ExpiresAt }
func (x *GwSessionEntry) GetRevokedAt() string  { if x == nil { return "" }; return x.RevokedAt }
func (x *GwSessionEntry) GetIsActive() bool     { if x == nil { return false }; return x.IsActive }

type GwListSessionsResponse struct {
	Sessions []*GwSessionEntry
}

func (x *GwListSessionsResponse) GetSessions() []*GwSessionEntry {
	if x == nil {
		return nil
	}
	return x.Sessions
}

type GwRevokeSessionRequest struct {
	SessionId   string
	RequestorId string
}

type GwRevokeSessionResponse struct {
	SessionId string
	RevokedAt string
}

func (x *GwRevokeSessionResponse) GetSessionId() string { if x == nil { return "" }; return x.SessionId }
func (x *GwRevokeSessionResponse) GetRevokedAt() string { if x == nil { return "" }; return x.RevokedAt }

// BL-IAM-015: Comm templates

type GwUpsertCommTemplateRequest struct {
	Channel   string
	Name      string
	Subject   string
	Body      string
	Variables []string
	UpdatedBy string
}

type GwUpsertCommTemplateResponse struct {
	Key       string
	UpdatedAt string
}

func (x *GwUpsertCommTemplateResponse) GetKey() string       { if x == nil { return "" }; return x.Key }
func (x *GwUpsertCommTemplateResponse) GetUpdatedAt() string { if x == nil { return "" }; return x.UpdatedAt }

type GwListCommTemplatesRequest struct {
	Channel string
}

type GwCommTemplate struct {
	Key       string
	Channel   string
	Name      string
	Subject   string
	Body      string
	Variables []string
	UpdatedAt string
}

func (x *GwCommTemplate) GetKey() string         { if x == nil { return "" }; return x.Key }
func (x *GwCommTemplate) GetChannel() string     { if x == nil { return "" }; return x.Channel }
func (x *GwCommTemplate) GetName() string        { if x == nil { return "" }; return x.Name }
func (x *GwCommTemplate) GetSubject() string     { if x == nil { return "" }; return x.Subject }
func (x *GwCommTemplate) GetBody() string        { if x == nil { return "" }; return x.Body }
func (x *GwCommTemplate) GetVariables() []string { if x == nil { return nil }; return x.Variables }
func (x *GwCommTemplate) GetUpdatedAt() string   { if x == nil { return "" }; return x.UpdatedAt }

type GwListCommTemplatesResponse struct {
	Templates []*GwCommTemplate
}

func (x *GwListCommTemplatesResponse) GetTemplates() []*GwCommTemplate {
	if x == nil {
		return nil
	}
	return x.Templates
}

// BL-IAM-017: Backup

type GwTriggerBackupRequest struct {
	TriggeredBy string
	Label       string
}

type GwTriggerBackupResponse struct {
	BackupId    string
	Status      string
	ScheduledAt string
}

func (x *GwTriggerBackupResponse) GetBackupId() string    { if x == nil { return "" }; return x.BackupId }
func (x *GwTriggerBackupResponse) GetStatus() string      { if x == nil { return "" }; return x.Status }
func (x *GwTriggerBackupResponse) GetScheduledAt() string { if x == nil { return "" }; return x.ScheduledAt }

type GwGetBackupHistoryRequest struct {
	Limit int32
}

type GwBackupEntry struct {
	BackupId    string
	Status      string
	TriggeredBy string
	ScheduledAt string
}

func (x *GwBackupEntry) GetBackupId() string    { if x == nil { return "" }; return x.BackupId }
func (x *GwBackupEntry) GetStatus() string      { if x == nil { return "" }; return x.Status }
func (x *GwBackupEntry) GetTriggeredBy() string { if x == nil { return "" }; return x.TriggeredBy }
func (x *GwBackupEntry) GetScheduledAt() string { if x == nil { return "" }; return x.ScheduledAt }

type GwGetBackupHistoryResponse struct {
	Backups []*GwBackupEntry
}

func (x *GwGetBackupHistoryResponse) GetBackups() []*GwBackupEntry {
	if x == nil {
		return nil
	}
	return x.Backups
}

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

// IamSecurityClient is the gateway client interface for iam-svc security RPCs.
type IamSecurityClient interface {
	GetPasswordPolicy(ctx context.Context, in *GwGetPasswordPolicyRequest, opts ...grpc.CallOption) (*GwGetPasswordPolicyResponse, error)
	SetPasswordPolicy(ctx context.Context, in *GwSetPasswordPolicyRequest, opts ...grpc.CallOption) (*GwGetPasswordPolicyResponse, error)
	RecordLoginAnomaly(ctx context.Context, in *GwRecordLoginAnomalyRequest, opts ...grpc.CallOption) (*GwRecordLoginAnomalyResponse, error)
	ListSessions(ctx context.Context, in *GwListSessionsRequest, opts ...grpc.CallOption) (*GwListSessionsResponse, error)
	RevokeSession(ctx context.Context, in *GwRevokeSessionRequest, opts ...grpc.CallOption) (*GwRevokeSessionResponse, error)
	UpsertCommTemplate(ctx context.Context, in *GwUpsertCommTemplateRequest, opts ...grpc.CallOption) (*GwUpsertCommTemplateResponse, error)
	ListCommTemplates(ctx context.Context, in *GwListCommTemplatesRequest, opts ...grpc.CallOption) (*GwListCommTemplatesResponse, error)
	TriggerBackup(ctx context.Context, in *GwTriggerBackupRequest, opts ...grpc.CallOption) (*GwTriggerBackupResponse, error)
	GetBackupHistory(ctx context.Context, in *GwGetBackupHistoryRequest, opts ...grpc.CallOption) (*GwGetBackupHistoryResponse, error)
}

type iamSecurityClient struct {
	cc grpc.ClientConnInterface
}

func NewIamSecurityClient(cc grpc.ClientConnInterface) IamSecurityClient {
	return &iamSecurityClient{cc}
}

func (c *iamSecurityClient) GetPasswordPolicy(ctx context.Context, in *GwGetPasswordPolicyRequest, opts ...grpc.CallOption) (*GwGetPasswordPolicyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwGetPasswordPolicyResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_GetPasswordPolicy_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) SetPasswordPolicy(ctx context.Context, in *GwSetPasswordPolicyRequest, opts ...grpc.CallOption) (*GwGetPasswordPolicyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwGetPasswordPolicyResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_SetPasswordPolicy_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) RecordLoginAnomaly(ctx context.Context, in *GwRecordLoginAnomalyRequest, opts ...grpc.CallOption) (*GwRecordLoginAnomalyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwRecordLoginAnomalyResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_RecordLoginAnomaly_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) ListSessions(ctx context.Context, in *GwListSessionsRequest, opts ...grpc.CallOption) (*GwListSessionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwListSessionsResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_ListSessions_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) RevokeSession(ctx context.Context, in *GwRevokeSessionRequest, opts ...grpc.CallOption) (*GwRevokeSessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwRevokeSessionResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_RevokeSession_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) UpsertCommTemplate(ctx context.Context, in *GwUpsertCommTemplateRequest, opts ...grpc.CallOption) (*GwUpsertCommTemplateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwUpsertCommTemplateResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_UpsertCommTemplate_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) ListCommTemplates(ctx context.Context, in *GwListCommTemplatesRequest, opts ...grpc.CallOption) (*GwListCommTemplatesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwListCommTemplatesResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_ListCommTemplates_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) TriggerBackup(ctx context.Context, in *GwTriggerBackupRequest, opts ...grpc.CallOption) (*GwTriggerBackupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwTriggerBackupResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_TriggerBackup_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamSecurityClient) GetBackupHistory(ctx context.Context, in *GwGetBackupHistoryRequest, opts ...grpc.CallOption) (*GwGetBackupHistoryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GwGetBackupHistoryResponse)
	if err := c.cc.Invoke(ctx, IamSecurityService_GetBackupHistory_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
