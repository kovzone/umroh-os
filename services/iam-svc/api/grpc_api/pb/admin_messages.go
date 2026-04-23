// admin_messages.go — hand-written message types for the IAM admin gRPC surface.
//
// These are plain Go structs that mirror what would be generated from .proto
// definitions. We avoid regenerating the .proto / .pb.go here because the
// admin surface is delivered as a hand-rolled extension (admin_grpc_ext.go)
// registered alongside the main IamService.
//
// S1-E-06 depth card (BL-IAM-005..017).

package pb

// ── Shared value types ───────────────────────────────────────────────────────

// AdminPermission represents one iam.permissions row.
type AdminPermission struct {
	Id       string `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}

// AdminRole represents one iam.roles row with its permission set.
type AdminRole struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Permissions []*AdminPermission `json:"permissions"`
	CreatedAt   int64              `json:"created_at"` // Unix seconds
}

// AdminUserSummary is a lightweight user record used in list responses.
type AdminUserSummary struct {
	Id          string   `json:"id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	BranchId    string   `json:"branch_id"`
	Status      string   `json:"status"`
	Roles       []string `json:"roles"`
	LastLoginAt int64    `json:"last_login_at"` // Unix seconds; 0 = never
	CreatedAt   int64    `json:"created_at"`
}

// AdminUser is the full user record returned from GetUser / CreateUser /
// UpdateUser. Roles are strings (names) for wire simplicity.
type AdminUser struct {
	Id        string   `json:"id"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	BranchId  string   `json:"branch_id"`
	Status    string   `json:"status"`
	Roles     []string `json:"roles"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
}

// ── User Management Requests / Responses ─────────────────────────────────────

type ListUsersRequest struct {
	Status   string `json:"status"`    // optional: active|suspended|pending
	BranchId string `json:"branch_id"` // optional UUID
	Cursor   string `json:"cursor"`    // opaque: "createdAt_unix:uuid"
	Limit    int32  `json:"limit"`     // default 20, max 100
}

type ListUsersResponse struct {
	Users      []*AdminUserSummary `json:"users"`
	NextCursor string              `json:"next_cursor"` // empty = no more pages
}

type CreateUserRequest struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	BranchId string   `json:"branch_id"`
	RoleIds  []string `json:"role_ids"`
}

type CreateUserResponse struct {
	User *AdminUser `json:"user"`
}

type UpdateUserRequest struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`    // optional; empty = keep current
	Status  string   `json:"status"`  // optional; empty = keep current
	RoleIds []string `json:"role_ids"` // optional; nil = keep current, [] = remove all
}

type UpdateUserResponse struct {
	User *AdminUser `json:"user"`
}

type GetUserRequest struct {
	Id string `json:"id"`
}

type GetUserResponse struct {
	User *AdminUser `json:"user"`
}

type ResetUserPasswordRequest struct {
	Id          string `json:"id"`
	NewPassword string `json:"new_password"`
}

type ResetUserPasswordResponse struct {
	Ok bool `json:"ok"`
}

// ── Role Management Requests / Responses ─────────────────────────────────────

type ListRolesRequest struct {
	Cursor string `json:"cursor"` // opaque: "createdAt_unix:uuid"
	Limit  int32  `json:"limit"`
}

type ListRolesResponse struct {
	Roles      []*AdminRole `json:"roles"`
	NextCursor string       `json:"next_cursor"`
}

type CreateRoleRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIds []string `json:"permission_ids"`
}

type CreateRoleResponse struct {
	Role *AdminRole `json:"role"`
}

type UpdateRoleRequest struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`        // optional
	Description   string   `json:"description"` // optional
	PermissionIds []string `json:"permission_ids"` // optional; nil = keep current
}

type UpdateRoleResponse struct {
	Role *AdminRole `json:"role"`
}

type DeleteRoleRequest struct {
	Id string `json:"id"`
}

type DeleteRoleResponse struct {
	Ok bool `json:"ok"`
}

type ListPermissionsRequest struct{}

type ListPermissionsResponse struct {
	Permissions []*AdminPermission `json:"permissions"`
}

type AssignRoleToUserRequest struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}

type AssignRoleToUserResponse struct {
	Ok bool `json:"ok"`
}

type RevokeRoleFromUserRequest struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}

type RevokeRoleFromUserResponse struct {
	Ok bool `json:"ok"`
}
