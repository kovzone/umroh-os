// admin_grpc_ext.go — IamAdminHandler interface + grpc ServiceDesc registration.
//
// The admin surface is an extension of IamService: it reuses the same gRPC
// server struct but registers additional methods under the same service name.
// Because we cannot re-run protoc here we hand-write the ServiceDesc additions
// and wire them via RegisterIamAdminHandlerExtension which must be called AFTER
// RegisterIamServiceServer.
//
// S1-E-06 depth card.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

// IamAdminHandler is the interface the gRPC server layer must implement for
// the admin surface. Each method maps 1:1 to a service/IService admin call.
type IamAdminHandler interface {
	// User management
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error)
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	ResetUserPassword(ctx context.Context, req *ResetUserPasswordRequest) (*ResetUserPasswordResponse, error)

	// Role management
	ListRoles(ctx context.Context, req *ListRolesRequest) (*ListRolesResponse, error)
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error)
	UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*UpdateRoleResponse, error)
	DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error)
	ListPermissions(ctx context.Context, req *ListPermissionsRequest) (*ListPermissionsResponse, error)
	AssignRoleToUser(ctx context.Context, req *AssignRoleToUserRequest) (*AssignRoleToUserResponse, error)
	RevokeRoleFromUser(ctx context.Context, req *RevokeRoleFromUserRequest) (*RevokeRoleFromUserResponse, error)
}

// ── Handler adapters ──────────────────────────────────────────────────────────

func _IamAdmin_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/ListUsers"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).ListUsers(ctx, req.(*ListUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/CreateUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/UpdateUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/GetUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_ResetUserPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetUserPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).ResetUserPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/ResetUserPassword"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).ResetUserPassword(ctx, req.(*ResetUserPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_ListRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).ListRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/ListRoles"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).ListRoles(ctx, req.(*ListRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_CreateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).CreateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/CreateRole"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).CreateRole(ctx, req.(*CreateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_UpdateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).UpdateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/UpdateRole"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).UpdateRole(ctx, req.(*UpdateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_DeleteRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).DeleteRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/DeleteRole"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).DeleteRole(ctx, req.(*DeleteRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_ListPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).ListPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/ListPermissions"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).ListPermissions(ctx, req.(*ListPermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_AssignRoleToUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignRoleToUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).AssignRoleToUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/AdminAssignRoleToUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).AssignRoleToUser(ctx, req.(*AssignRoleToUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamAdmin_RevokeRoleFromUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeRoleFromUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamAdminHandler).RevokeRoleFromUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamService/AdminRevokeRoleFromUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamAdminHandler).RevokeRoleFromUser(ctx, req.(*RevokeRoleFromUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ── Registration ──────────────────────────────────────────────────────────────

// IamAdminServiceDesc extends IamService with admin RPCs.
// Call grpc.RegisterService(s, &IamAdminServiceDesc) AFTER
// RegisterIamServiceServer so the methods append to the same service descriptor.
//
// NOTE: In practice the go-grpc runtime does not support splitting a single
// proto service across two ServiceDescs. The recommended approach is therefore
// to register a *separate* service name. We use "pb.IamAdminService" here.
var IamAdminServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.IamAdminService",
	HandlerType: (*IamAdminHandler)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "ListUsers", Handler: _IamAdmin_ListUsers_Handler},
		{MethodName: "CreateUser", Handler: _IamAdmin_CreateUser_Handler},
		{MethodName: "UpdateUser", Handler: _IamAdmin_UpdateUser_Handler},
		{MethodName: "GetUser", Handler: _IamAdmin_GetUser_Handler},
		{MethodName: "ResetUserPassword", Handler: _IamAdmin_ResetUserPassword_Handler},
		{MethodName: "ListRoles", Handler: _IamAdmin_ListRoles_Handler},
		{MethodName: "CreateRole", Handler: _IamAdmin_CreateRole_Handler},
		{MethodName: "UpdateRole", Handler: _IamAdmin_UpdateRole_Handler},
		{MethodName: "DeleteRole", Handler: _IamAdmin_DeleteRole_Handler},
		{MethodName: "ListPermissions", Handler: _IamAdmin_ListPermissions_Handler},
		{MethodName: "AssignRoleToUser", Handler: _IamAdmin_AssignRoleToUser_Handler},
		{MethodName: "RevokeRoleFromUser", Handler: _IamAdmin_RevokeRoleFromUser_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "iam_admin.proto",
}

// RegisterIamAdminHandler registers the admin extension on a gRPC server.
func RegisterIamAdminHandler(s grpc.ServiceRegistrar, srv IamAdminHandler) {
	s.RegisterService(&IamAdminServiceDesc, srv)
}
