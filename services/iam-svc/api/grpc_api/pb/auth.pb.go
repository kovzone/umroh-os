// Hand-written supplement to the protoc-generated pb package.
// Added in S1-E-12 (BL-IAM-018) to carry the auth RPC message types until
// the next `make genpb` regenerates iam.pb.go from the updated iam.proto.
// When protoc is re-run, delete this file and keep the generated output.

package pb

// ---------------------------------------------------------------------------
// Shared
// ---------------------------------------------------------------------------

// UserProfile is the shared user shape returned by Login, GetMe, SuspendUser.
type UserProfile struct {
	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Email    string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Name     string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	BranchId string `protobuf:"bytes,4,opt,name=branch_id,json=branchId,proto3" json:"branch_id,omitempty"`
	Status   string `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"` // "active"|"suspended"|"pending"
}

func (x *UserProfile) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}
func (x *UserProfile) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}
func (x *UserProfile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *UserProfile) GetBranchId() string {
	if x != nil {
		return x.BranchId
	}
	return ""
}
func (x *UserProfile) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *UserProfile) Reset()         {}
func (x *UserProfile) String() string { return x.UserId }
func (x *UserProfile) ProtoMessage()  {}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

type LoginRequest struct {
	Email     string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password  string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	TotpCode  string `protobuf:"bytes,3,opt,name=totp_code,json=totpCode,proto3" json:"totp_code,omitempty"`
	UserAgent string `protobuf:"bytes,4,opt,name=user_agent,json=userAgent,proto3" json:"user_agent,omitempty"`
	Ip        string `protobuf:"bytes,5,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *LoginRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}
func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}
func (x *LoginRequest) GetTotpCode() string {
	if x != nil {
		return x.TotpCode
	}
	return ""
}
func (x *LoginRequest) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}
func (x *LoginRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *LoginRequest) Reset()         {}
func (x *LoginRequest) String() string { return x.Email }
func (x *LoginRequest) ProtoMessage()  {}

type LoginResponse struct {
	AccessToken      string       `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken     string       `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	AccessExpiresAt  int64        `protobuf:"varint,3,opt,name=access_expires_at,json=accessExpiresAt,proto3" json:"access_expires_at,omitempty"`
	RefreshExpiresAt int64        `protobuf:"varint,4,opt,name=refresh_expires_at,json=refreshExpiresAt,proto3" json:"refresh_expires_at,omitempty"`
	User             *UserProfile `protobuf:"bytes,5,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *LoginResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}
func (x *LoginResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}
func (x *LoginResponse) GetAccessExpiresAt() int64 {
	if x != nil {
		return x.AccessExpiresAt
	}
	return 0
}
func (x *LoginResponse) GetRefreshExpiresAt() int64 {
	if x != nil {
		return x.RefreshExpiresAt
	}
	return 0
}
func (x *LoginResponse) GetUser() *UserProfile {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *LoginResponse) Reset()         {}
func (x *LoginResponse) String() string { return "" }
func (x *LoginResponse) ProtoMessage()  {}

// ---------------------------------------------------------------------------
// RefreshSession
// ---------------------------------------------------------------------------

type RefreshSessionRequest struct {
	RefreshToken string `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	UserAgent    string `protobuf:"bytes,2,opt,name=user_agent,json=userAgent,proto3" json:"user_agent,omitempty"`
	Ip           string `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *RefreshSessionRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}
func (x *RefreshSessionRequest) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}
func (x *RefreshSessionRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *RefreshSessionRequest) Reset()         {}
func (x *RefreshSessionRequest) String() string { return "" }
func (x *RefreshSessionRequest) ProtoMessage()  {}

type RefreshSessionResponse struct {
	AccessToken      string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken     string `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	AccessExpiresAt  int64  `protobuf:"varint,3,opt,name=access_expires_at,json=accessExpiresAt,proto3" json:"access_expires_at,omitempty"`
	RefreshExpiresAt int64  `protobuf:"varint,4,opt,name=refresh_expires_at,json=refreshExpiresAt,proto3" json:"refresh_expires_at,omitempty"`
}

func (x *RefreshSessionResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}
func (x *RefreshSessionResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}
func (x *RefreshSessionResponse) GetAccessExpiresAt() int64 {
	if x != nil {
		return x.AccessExpiresAt
	}
	return 0
}
func (x *RefreshSessionResponse) GetRefreshExpiresAt() int64 {
	if x != nil {
		return x.RefreshExpiresAt
	}
	return 0
}

func (x *RefreshSessionResponse) Reset()         {}
func (x *RefreshSessionResponse) String() string { return "" }
func (x *RefreshSessionResponse) ProtoMessage()  {}

// ---------------------------------------------------------------------------
// Logout
// ---------------------------------------------------------------------------

type LogoutRequest struct {
	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *LogoutRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *LogoutRequest) Reset()         {}
func (x *LogoutRequest) String() string { return "" }
func (x *LogoutRequest) ProtoMessage()  {}

type LogoutResponse struct {
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *LogoutResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *LogoutResponse) Reset()         {}
func (x *LogoutResponse) String() string { return "" }
func (x *LogoutResponse) ProtoMessage()  {}

// ---------------------------------------------------------------------------
// GetMe
// ---------------------------------------------------------------------------

type GetMeRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetMeRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetMeRequest) Reset()         {}
func (x *GetMeRequest) String() string { return "" }
func (x *GetMeRequest) ProtoMessage()  {}

type GetMeResponse struct {
	User         *UserProfile `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	TotpEnrolled bool         `protobuf:"varint,2,opt,name=totp_enrolled,json=totpEnrolled,proto3" json:"totp_enrolled,omitempty"`
	TotpVerified bool         `protobuf:"varint,3,opt,name=totp_verified,json=totpVerified,proto3" json:"totp_verified,omitempty"`
}

func (x *GetMeResponse) GetUser() *UserProfile {
	if x != nil {
		return x.User
	}
	return nil
}
func (x *GetMeResponse) GetTotpEnrolled() bool {
	if x != nil {
		return x.TotpEnrolled
	}
	return false
}
func (x *GetMeResponse) GetTotpVerified() bool {
	if x != nil {
		return x.TotpVerified
	}
	return false
}

func (x *GetMeResponse) Reset()         {}
func (x *GetMeResponse) String() string { return "" }
func (x *GetMeResponse) ProtoMessage()  {}

// ---------------------------------------------------------------------------
// EnrollTOTP
// ---------------------------------------------------------------------------

type EnrollTOTPRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *EnrollTOTPRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *EnrollTOTPRequest) Reset()         {}
func (x *EnrollTOTPRequest) String() string { return "" }
func (x *EnrollTOTPRequest) ProtoMessage()  {}

type EnrollTOTPResponse struct {
	Secret     string `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
	OtpauthUrl string `protobuf:"bytes,2,opt,name=otpauth_url,json=otpauthUrl,proto3" json:"otpauth_url,omitempty"`
}

func (x *EnrollTOTPResponse) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}
func (x *EnrollTOTPResponse) GetOtpauthUrl() string {
	if x != nil {
		return x.OtpauthUrl
	}
	return ""
}

func (x *EnrollTOTPResponse) Reset()         {}
func (x *EnrollTOTPResponse) String() string { return "" }
func (x *EnrollTOTPResponse) ProtoMessage()  {}

// ---------------------------------------------------------------------------
// VerifyTOTP
// ---------------------------------------------------------------------------

type VerifyTOTPRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Code   string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *VerifyTOTPRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}
func (x *VerifyTOTPRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *VerifyTOTPRequest) Reset()         {}
func (x *VerifyTOTPRequest) String() string { return "" }
func (x *VerifyTOTPRequest) ProtoMessage()  {}

type VerifyTOTPResponse struct {
	VerifiedAtUnix int64 `protobuf:"varint,1,opt,name=verified_at_unix,json=verifiedAtUnix,proto3" json:"verified_at_unix,omitempty"`
}

func (x *VerifyTOTPResponse) GetVerifiedAtUnix() int64 {
	if x != nil {
		return x.VerifiedAtUnix
	}
	return 0
}

func (x *VerifyTOTPResponse) Reset()         {}
func (x *VerifyTOTPResponse) String() string { return "" }
func (x *VerifyTOTPResponse) ProtoMessage()  {}

// ---------------------------------------------------------------------------
// SuspendUser
// ---------------------------------------------------------------------------

type SuspendUserRequest struct {
	ActorUserId  string `protobuf:"bytes,1,opt,name=actor_user_id,json=actorUserId,proto3" json:"actor_user_id,omitempty"`
	TargetUserId string `protobuf:"bytes,2,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
}

func (x *SuspendUserRequest) GetActorUserId() string {
	if x != nil {
		return x.ActorUserId
	}
	return ""
}
func (x *SuspendUserRequest) GetTargetUserId() string {
	if x != nil {
		return x.TargetUserId
	}
	return ""
}

func (x *SuspendUserRequest) Reset()         {}
func (x *SuspendUserRequest) String() string { return "" }
func (x *SuspendUserRequest) ProtoMessage()  {}

type SuspendUserResponse struct {
	User *UserProfile `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *SuspendUserResponse) GetUser() *UserProfile {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *SuspendUserResponse) Reset()         {}
func (x *SuspendUserResponse) String() string { return "" }
func (x *SuspendUserResponse) ProtoMessage()  {}
