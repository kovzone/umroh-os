// iam_admin_stub.go — gateway-side gRPC client stub for iam-svc admin RPCs
// (Wave 1C / Phase 6).
//
// Mirrors services/iam-svc/api/grpc_api/pb/admin_messages.go and
// admin_grpc_ext.go. Run `make genpb` to replace with generated code.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	IamAdminService_ListUsers_FullMethodName        = "/pb.IamAdminService/ListUsers"
	IamAdminService_CreateUser_FullMethodName       = "/pb.IamAdminService/CreateUser"
	IamAdminService_UpdateUser_FullMethodName       = "/pb.IamAdminService/UpdateUser"
	IamAdminService_GetUser_FullMethodName          = "/pb.IamAdminService/GetUser"
	IamAdminService_ResetUserPassword_FullMethodName = "/pb.IamAdminService/ResetUserPassword"
	IamAdminService_ListRoles_FullMethodName        = "/pb.IamAdminService/ListRoles"
	IamAdminService_CreateRole_FullMethodName       = "/pb.IamAdminService/CreateRole"
	IamAdminService_UpdateRole_FullMethodName       = "/pb.IamAdminService/UpdateRole"
	IamAdminService_DeleteRole_FullMethodName       = "/pb.IamAdminService/DeleteRole"
	IamAdminService_ListPermissions_FullMethodName  = "/pb.IamAdminService/ListPermissions"
	IamAdminService_AssignRoleToUser_FullMethodName = "/pb.IamAdminService/AssignRoleToUser"
	IamAdminService_RevokeRoleFromUser_FullMethodName = "/pb.IamAdminService/RevokeRoleFromUser"
)

// ---------------------------------------------------------------------------
// Shared value types
// ---------------------------------------------------------------------------

// AdminPermission represents one iam.permissions row.
type AdminPermission struct {
	Id       string
	Resource string
	Action   string
	Scope    string
}

func (x *AdminPermission) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *AdminPermission) GetResource() string {
	if x == nil {
		return ""
	}
	return x.Resource
}
func (x *AdminPermission) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *AdminPermission) GetScope() string {
	if x == nil {
		return ""
	}
	return x.Scope
}

// AdminRole represents one iam.roles row with its permission set.
type AdminRole struct {
	Id          string
	Name        string
	Description string
	Permissions []*AdminPermission
	CreatedAt   int64
}

func (x *AdminRole) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *AdminRole) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *AdminRole) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *AdminRole) GetPermissions() []*AdminPermission {
	if x == nil {
		return nil
	}
	return x.Permissions
}
func (x *AdminRole) GetCreatedAt() int64 {
	if x == nil {
		return 0
	}
	return x.CreatedAt
}

// AdminUserSummary is a lightweight user record used in list responses.
type AdminUserSummary struct {
	Id          string
	Email       string
	Name        string
	BranchId    string
	Status      string
	Roles       []string
	LastLoginAt int64
	CreatedAt   int64
}

func (x *AdminUserSummary) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *AdminUserSummary) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *AdminUserSummary) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *AdminUserSummary) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *AdminUserSummary) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *AdminUserSummary) GetRoles() []string {
	if x == nil {
		return nil
	}
	return x.Roles
}
func (x *AdminUserSummary) GetLastLoginAt() int64 {
	if x == nil {
		return 0
	}
	return x.LastLoginAt
}
func (x *AdminUserSummary) GetCreatedAt() int64 {
	if x == nil {
		return 0
	}
	return x.CreatedAt
}

// AdminUser is the full user record.
type AdminUser struct {
	Id        string
	Email     string
	Name      string
	BranchId  string
	Status    string
	Roles     []string
	CreatedAt int64
	UpdatedAt int64
}

func (x *AdminUser) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *AdminUser) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *AdminUser) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *AdminUser) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *AdminUser) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *AdminUser) GetRoles() []string {
	if x == nil {
		return nil
	}
	return x.Roles
}
func (x *AdminUser) GetCreatedAt() int64 {
	if x == nil {
		return 0
	}
	return x.CreatedAt
}
func (x *AdminUser) GetUpdatedAt() int64 {
	if x == nil {
		return 0
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// User management request/response types
// ---------------------------------------------------------------------------

type ListUsersRequest struct {
	Status   string
	BranchId string
	Cursor   string
	Limit    int32
}

func (x *ListUsersRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *ListUsersRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *ListUsersRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}
func (x *ListUsersRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}

type ListUsersResponse struct {
	Users      []*AdminUserSummary
	NextCursor string
}

func (x *ListUsersResponse) GetUsers() []*AdminUserSummary {
	if x == nil {
		return nil
	}
	return x.Users
}
func (x *ListUsersResponse) GetNextCursor() string {
	if x == nil {
		return ""
	}
	return x.NextCursor
}

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
	BranchId string
	RoleIds  []string
}

func (x *CreateUserRequest) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *CreateUserRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateUserRequest) GetPassword() string {
	if x == nil {
		return ""
	}
	return x.Password
}
func (x *CreateUserRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *CreateUserRequest) GetRoleIds() []string {
	if x == nil {
		return nil
	}
	return x.RoleIds
}

type CreateUserResponse struct {
	User *AdminUser
}

func (x *CreateUserResponse) GetUser() *AdminUser {
	if x == nil {
		return nil
	}
	return x.User
}

type UpdateUserRequest struct {
	Id      string
	Name    string
	Status  string
	RoleIds []string
}

func (x *UpdateUserRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateUserRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *UpdateUserRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *UpdateUserRequest) GetRoleIds() []string {
	if x == nil {
		return nil
	}
	return x.RoleIds
}

type UpdateUserResponse struct {
	User *AdminUser
}

func (x *UpdateUserResponse) GetUser() *AdminUser {
	if x == nil {
		return nil
	}
	return x.User
}

type GetUserRequest struct {
	Id string
}

func (x *GetUserRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

type GetUserResponse struct {
	User *AdminUser
}

func (x *GetUserResponse) GetUser() *AdminUser {
	if x == nil {
		return nil
	}
	return x.User
}

type ResetUserPasswordRequest struct {
	Id          string
	NewPassword string
}

func (x *ResetUserPasswordRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *ResetUserPasswordRequest) GetNewPassword() string {
	if x == nil {
		return ""
	}
	return x.NewPassword
}

type ResetUserPasswordResponse struct {
	Ok bool
}

func (x *ResetUserPasswordResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

// ---------------------------------------------------------------------------
// Role management request/response types
// ---------------------------------------------------------------------------

type ListRolesRequest struct {
	Cursor string
	Limit  int32
}

func (x *ListRolesRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}
func (x *ListRolesRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}

type ListRolesResponse struct {
	Roles      []*AdminRole
	NextCursor string
}

func (x *ListRolesResponse) GetRoles() []*AdminRole {
	if x == nil {
		return nil
	}
	return x.Roles
}
func (x *ListRolesResponse) GetNextCursor() string {
	if x == nil {
		return ""
	}
	return x.NextCursor
}

type CreateRoleRequest struct {
	Name          string
	Description   string
	PermissionIds []string
}

func (x *CreateRoleRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateRoleRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *CreateRoleRequest) GetPermissionIds() []string {
	if x == nil {
		return nil
	}
	return x.PermissionIds
}

type CreateRoleResponse struct {
	Role *AdminRole
}

func (x *CreateRoleResponse) GetRole() *AdminRole {
	if x == nil {
		return nil
	}
	return x.Role
}

type UpdateRoleRequest struct {
	Id            string
	Name          string
	Description   string
	PermissionIds []string
}

func (x *UpdateRoleRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateRoleRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *UpdateRoleRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *UpdateRoleRequest) GetPermissionIds() []string {
	if x == nil {
		return nil
	}
	return x.PermissionIds
}

type UpdateRoleResponse struct {
	Role *AdminRole
}

func (x *UpdateRoleResponse) GetRole() *AdminRole {
	if x == nil {
		return nil
	}
	return x.Role
}

type DeleteRoleRequest struct {
	Id string
}

func (x *DeleteRoleRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

type DeleteRoleResponse struct {
	Ok bool
}

func (x *DeleteRoleResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

type ListPermissionsRequest struct{}

type ListPermissionsResponse struct {
	Permissions []*AdminPermission
}

func (x *ListPermissionsResponse) GetPermissions() []*AdminPermission {
	if x == nil {
		return nil
	}
	return x.Permissions
}

type AssignRoleToUserRequest struct {
	UserId string
	RoleId string
}

func (x *AssignRoleToUserRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *AssignRoleToUserRequest) GetRoleId() string {
	if x == nil {
		return ""
	}
	return x.RoleId
}

type AssignRoleToUserResponse struct {
	Ok bool
}

func (x *AssignRoleToUserResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

type RevokeRoleFromUserRequest struct {
	UserId string
	RoleId string
}

func (x *RevokeRoleFromUserRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *RevokeRoleFromUserRequest) GetRoleId() string {
	if x == nil {
		return ""
	}
	return x.RoleId
}

type RevokeRoleFromUserResponse struct {
	Ok bool
}

func (x *RevokeRoleFromUserResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

// ---------------------------------------------------------------------------
// IamAdminClient interface + implementation
// ---------------------------------------------------------------------------

// IamAdminClient is the consumer-side interface for iam-svc admin RPCs.
type IamAdminClient interface {
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	ResetUserPassword(ctx context.Context, in *ResetUserPasswordRequest, opts ...grpc.CallOption) (*ResetUserPasswordResponse, error)
	ListRoles(ctx context.Context, in *ListRolesRequest, opts ...grpc.CallOption) (*ListRolesResponse, error)
	CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error)
	UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...grpc.CallOption) (*UpdateRoleResponse, error)
	DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error)
	ListPermissions(ctx context.Context, in *ListPermissionsRequest, opts ...grpc.CallOption) (*ListPermissionsResponse, error)
	AssignRoleToUser(ctx context.Context, in *AssignRoleToUserRequest, opts ...grpc.CallOption) (*AssignRoleToUserResponse, error)
	RevokeRoleFromUser(ctx context.Context, in *RevokeRoleFromUserRequest, opts ...grpc.CallOption) (*RevokeRoleFromUserResponse, error)
}

type iamAdminClient struct {
	cc grpc.ClientConnInterface
}

// NewIamAdminClient wraps a conn so gateway-svc can call iam admin RPCs.
func NewIamAdminClient(cc grpc.ClientConnInterface) IamAdminClient {
	return &iamAdminClient{cc}
}

func (c *iamAdminClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, IamAdminService_ListUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, IamAdminService_CreateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserResponse)
	err := c.cc.Invoke(ctx, IamAdminService_UpdateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, IamAdminService_GetUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) ResetUserPassword(ctx context.Context, in *ResetUserPasswordRequest, opts ...grpc.CallOption) (*ResetUserPasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResetUserPasswordResponse)
	err := c.cc.Invoke(ctx, IamAdminService_ResetUserPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) ListRoles(ctx context.Context, in *ListRolesRequest, opts ...grpc.CallOption) (*ListRolesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRolesResponse)
	err := c.cc.Invoke(ctx, IamAdminService_ListRoles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRoleResponse)
	err := c.cc.Invoke(ctx, IamAdminService_CreateRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...grpc.CallOption) (*UpdateRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateRoleResponse)
	err := c.cc.Invoke(ctx, IamAdminService_UpdateRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteRoleResponse)
	err := c.cc.Invoke(ctx, IamAdminService_DeleteRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) ListPermissions(ctx context.Context, in *ListPermissionsRequest, opts ...grpc.CallOption) (*ListPermissionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPermissionsResponse)
	err := c.cc.Invoke(ctx, IamAdminService_ListPermissions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) AssignRoleToUser(ctx context.Context, in *AssignRoleToUserRequest, opts ...grpc.CallOption) (*AssignRoleToUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AssignRoleToUserResponse)
	err := c.cc.Invoke(ctx, IamAdminService_AssignRoleToUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamAdminClient) RevokeRoleFromUser(ctx context.Context, in *RevokeRoleFromUserRequest, opts ...grpc.CallOption) (*RevokeRoleFromUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RevokeRoleFromUserResponse)
	err := c.cc.Invoke(ctx, IamAdminService_RevokeRoleFromUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
