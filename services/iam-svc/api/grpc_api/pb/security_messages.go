// security_messages.go — hand-written proto message types for IAM security depth.
// BL-IAM-010: GetPasswordPolicy, SetPasswordPolicy
// BL-IAM-012: RecordLoginAnomaly
// BL-IAM-013: ListSessions, RevokeSession
// BL-IAM-015: UpsertCommTemplate, ListCommTemplates
// BL-IAM-017: TriggerBackup, GetBackupHistory

package pb

// ---------------------------------------------------------------------------
// BL-IAM-010: Password policy
// ---------------------------------------------------------------------------

type GetPasswordPolicyRequest struct{}

type GetPasswordPolicyResponse struct {
	MinLength      int32
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMfa     bool
	UpdatedAt      string // RFC3339
}

func (x *GetPasswordPolicyResponse) GetMinLength() int32      { if x == nil { return 0 }; return x.MinLength }
func (x *GetPasswordPolicyResponse) GetRequireUpper() bool    { if x == nil { return false }; return x.RequireUpper }
func (x *GetPasswordPolicyResponse) GetRequireDigit() bool    { if x == nil { return false }; return x.RequireDigit }
func (x *GetPasswordPolicyResponse) GetRequireSpecial() bool  { if x == nil { return false }; return x.RequireSpecial }
func (x *GetPasswordPolicyResponse) GetRequireMfa() bool      { if x == nil { return false }; return x.RequireMfa }
func (x *GetPasswordPolicyResponse) GetUpdatedAt() string     { if x == nil { return "" }; return x.UpdatedAt }

type SetPasswordPolicyRequest struct {
	MinLength      int32
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMfa     bool
	UpdatedBy      string
}

func (x *SetPasswordPolicyRequest) GetMinLength() int32      { if x == nil { return 0 }; return x.MinLength }
func (x *SetPasswordPolicyRequest) GetRequireUpper() bool    { if x == nil { return false }; return x.RequireUpper }
func (x *SetPasswordPolicyRequest) GetRequireDigit() bool    { if x == nil { return false }; return x.RequireDigit }
func (x *SetPasswordPolicyRequest) GetRequireSpecial() bool  { if x == nil { return false }; return x.RequireSpecial }
func (x *SetPasswordPolicyRequest) GetRequireMfa() bool      { if x == nil { return false }; return x.RequireMfa }
func (x *SetPasswordPolicyRequest) GetUpdatedBy() string     { if x == nil { return "" }; return x.UpdatedBy }

// ---------------------------------------------------------------------------
// BL-IAM-012: Anomaly alerts
// ---------------------------------------------------------------------------

type RecordLoginAnomalyRequest struct {
	UserId      string
	Ip          string
	UserAgent   string
	AnomalyKind string // "new_ip" | "new_device" | "rapid_login" | "off_hours"
	Details     string
}

func (x *RecordLoginAnomalyRequest) GetUserId() string      { if x == nil { return "" }; return x.UserId }
func (x *RecordLoginAnomalyRequest) GetIp() string          { if x == nil { return "" }; return x.Ip }
func (x *RecordLoginAnomalyRequest) GetUserAgent() string   { if x == nil { return "" }; return x.UserAgent }
func (x *RecordLoginAnomalyRequest) GetAnomalyKind() string { if x == nil { return "" }; return x.AnomalyKind }
func (x *RecordLoginAnomalyRequest) GetDetails() string     { if x == nil { return "" }; return x.Details }

type RecordLoginAnomalyResponse struct {
	AlertId   string
	CreatedAt string // RFC3339
}

func (x *RecordLoginAnomalyResponse) GetAlertId() string   { if x == nil { return "" }; return x.AlertId }
func (x *RecordLoginAnomalyResponse) GetCreatedAt() string { if x == nil { return "" }; return x.CreatedAt }

// ---------------------------------------------------------------------------
// BL-IAM-013: Session history + revoke
// ---------------------------------------------------------------------------

type ListSessionsRequest struct {
	UserId     string
	IncludeAll bool
}

func (x *ListSessionsRequest) GetUserId() string     { if x == nil { return "" }; return x.UserId }
func (x *ListSessionsRequest) GetIncludeAll() bool   { if x == nil { return false }; return x.IncludeAll }

type SessionEntry struct {
	SessionId string
	UserAgent string
	Ip        string
	IssuedAt  string // RFC3339
	ExpiresAt string // RFC3339
	RevokedAt string // RFC3339 or ""
	IsActive  bool
}

func (x *SessionEntry) GetSessionId() string  { if x == nil { return "" }; return x.SessionId }
func (x *SessionEntry) GetUserAgent() string  { if x == nil { return "" }; return x.UserAgent }
func (x *SessionEntry) GetIp() string         { if x == nil { return "" }; return x.Ip }
func (x *SessionEntry) GetIssuedAt() string   { if x == nil { return "" }; return x.IssuedAt }
func (x *SessionEntry) GetExpiresAt() string  { if x == nil { return "" }; return x.ExpiresAt }
func (x *SessionEntry) GetRevokedAt() string  { if x == nil { return "" }; return x.RevokedAt }
func (x *SessionEntry) GetIsActive() bool     { if x == nil { return false }; return x.IsActive }

type ListSessionsResponse struct {
	Sessions []*SessionEntry
}

func (x *ListSessionsResponse) GetSessions() []*SessionEntry {
	if x == nil { return nil }
	return x.Sessions
}

type RevokeSessionRequest struct {
	SessionId   string
	RequestorId string
}

func (x *RevokeSessionRequest) GetSessionId() string   { if x == nil { return "" }; return x.SessionId }
func (x *RevokeSessionRequest) GetRequestorId() string { if x == nil { return "" }; return x.RequestorId }

type RevokeSessionResponse struct {
	SessionId string
	RevokedAt string // RFC3339
}

func (x *RevokeSessionResponse) GetSessionId() string  { if x == nil { return "" }; return x.SessionId }
func (x *RevokeSessionResponse) GetRevokedAt() string  { if x == nil { return "" }; return x.RevokedAt }

// ---------------------------------------------------------------------------
// BL-IAM-015: Communication templates
// ---------------------------------------------------------------------------

type UpsertCommTemplateRequest struct {
	Channel   string
	Name      string
	Subject   string
	Body      string
	Variables []string
	UpdatedBy string
}

func (x *UpsertCommTemplateRequest) GetChannel() string    { if x == nil { return "" }; return x.Channel }
func (x *UpsertCommTemplateRequest) GetName() string       { if x == nil { return "" }; return x.Name }
func (x *UpsertCommTemplateRequest) GetSubject() string    { if x == nil { return "" }; return x.Subject }
func (x *UpsertCommTemplateRequest) GetBody() string       { if x == nil { return "" }; return x.Body }
func (x *UpsertCommTemplateRequest) GetVariables() []string { if x == nil { return nil }; return x.Variables }
func (x *UpsertCommTemplateRequest) GetUpdatedBy() string  { if x == nil { return "" }; return x.UpdatedBy }

type UpsertCommTemplateResponse struct {
	Key       string
	UpdatedAt string // RFC3339
}

func (x *UpsertCommTemplateResponse) GetKey() string       { if x == nil { return "" }; return x.Key }
func (x *UpsertCommTemplateResponse) GetUpdatedAt() string { if x == nil { return "" }; return x.UpdatedAt }

type ListCommTemplatesRequest struct {
	Channel string // "" = all channels
}

func (x *ListCommTemplatesRequest) GetChannel() string { if x == nil { return "" }; return x.Channel }

type CommTemplate struct {
	Key       string
	Channel   string
	Name      string
	Subject   string
	Body      string
	Variables []string
	UpdatedAt string // RFC3339
}

func (x *CommTemplate) GetKey() string        { if x == nil { return "" }; return x.Key }
func (x *CommTemplate) GetChannel() string    { if x == nil { return "" }; return x.Channel }
func (x *CommTemplate) GetName() string       { if x == nil { return "" }; return x.Name }
func (x *CommTemplate) GetSubject() string    { if x == nil { return "" }; return x.Subject }
func (x *CommTemplate) GetBody() string       { if x == nil { return "" }; return x.Body }
func (x *CommTemplate) GetVariables() []string { if x == nil { return nil }; return x.Variables }
func (x *CommTemplate) GetUpdatedAt() string  { if x == nil { return "" }; return x.UpdatedAt }

type ListCommTemplatesResponse struct {
	Templates []*CommTemplate
}

func (x *ListCommTemplatesResponse) GetTemplates() []*CommTemplate {
	if x == nil { return nil }
	return x.Templates
}

// ---------------------------------------------------------------------------
// BL-IAM-017: DB backup
// ---------------------------------------------------------------------------

type TriggerBackupRequest struct {
	TriggeredBy string
	Label       string
}

func (x *TriggerBackupRequest) GetTriggeredBy() string { if x == nil { return "" }; return x.TriggeredBy }
func (x *TriggerBackupRequest) GetLabel() string       { if x == nil { return "" }; return x.Label }

type TriggerBackupResponse struct {
	BackupId    string
	Status      string
	ScheduledAt string // RFC3339
}

func (x *TriggerBackupResponse) GetBackupId() string    { if x == nil { return "" }; return x.BackupId }
func (x *TriggerBackupResponse) GetStatus() string      { if x == nil { return "" }; return x.Status }
func (x *TriggerBackupResponse) GetScheduledAt() string { if x == nil { return "" }; return x.ScheduledAt }

type GetBackupHistoryRequest struct {
	Limit int32
}

func (x *GetBackupHistoryRequest) GetLimit() int32 { if x == nil { return 0 }; return x.Limit }

type BackupEntry struct {
	BackupId    string
	Status      string
	TriggeredBy string
	ScheduledAt string // RFC3339
}

func (x *BackupEntry) GetBackupId() string    { if x == nil { return "" }; return x.BackupId }
func (x *BackupEntry) GetStatus() string      { if x == nil { return "" }; return x.Status }
func (x *BackupEntry) GetTriggeredBy() string { if x == nil { return "" }; return x.TriggeredBy }
func (x *BackupEntry) GetScheduledAt() string { if x == nil { return "" }; return x.ScheduledAt }

type GetBackupHistoryResponse struct {
	Backups []*BackupEntry
}

func (x *GetBackupHistoryResponse) GetBackups() []*BackupEntry {
	if x == nil { return nil }
	return x.Backups
}
