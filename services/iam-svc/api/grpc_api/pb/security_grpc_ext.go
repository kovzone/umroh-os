// security_grpc_ext.go — IamSecurityHandler interface + grpc ServiceDesc.
//
// BL-IAM-010: GetPasswordPolicy, SetPasswordPolicy
// BL-IAM-012: RecordLoginAnomaly
// BL-IAM-013: ListSessions, RevokeSession
// BL-IAM-015: UpsertCommTemplate, ListCommTemplates
// BL-IAM-017: TriggerBackup, GetBackupHistory
//
// Pattern mirrors admin_grpc_ext.go — separate service name "pb.IamSecurityService"
// registered alongside the main IamService. Call RegisterIamSecurityHandler AFTER
// the main service is registered.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

// IamSecurityHandler is implemented by the grpc_api.Server.
type IamSecurityHandler interface {
	GetPasswordPolicy(ctx context.Context, req *GetPasswordPolicyRequest) (*GetPasswordPolicyResponse, error)
	SetPasswordPolicy(ctx context.Context, req *SetPasswordPolicyRequest) (*GetPasswordPolicyResponse, error)
	RecordLoginAnomaly(ctx context.Context, req *RecordLoginAnomalyRequest) (*RecordLoginAnomalyResponse, error)
	ListSessions(ctx context.Context, req *ListSessionsRequest) (*ListSessionsResponse, error)
	RevokeSession(ctx context.Context, req *RevokeSessionRequest) (*RevokeSessionResponse, error)
	UpsertCommTemplate(ctx context.Context, req *UpsertCommTemplateRequest) (*UpsertCommTemplateResponse, error)
	ListCommTemplates(ctx context.Context, req *ListCommTemplatesRequest) (*ListCommTemplatesResponse, error)
	TriggerBackup(ctx context.Context, req *TriggerBackupRequest) (*TriggerBackupResponse, error)
	GetBackupHistory(ctx context.Context, req *GetBackupHistoryRequest) (*GetBackupHistoryResponse, error)
}

// ── Handler adapters ──────────────────────────────────────────────────────────

func _IamSec_GetPasswordPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPasswordPolicyRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).GetPasswordPolicy(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/GetPasswordPolicy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).GetPasswordPolicy(ctx, req.(*GetPasswordPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_SetPasswordPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetPasswordPolicyRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).SetPasswordPolicy(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/SetPasswordPolicy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).SetPasswordPolicy(ctx, req.(*SetPasswordPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_RecordLoginAnomaly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordLoginAnomalyRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).RecordLoginAnomaly(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/RecordLoginAnomaly"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).RecordLoginAnomaly(ctx, req.(*RecordLoginAnomalyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_ListSessions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSessionsRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).ListSessions(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/ListSessions"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).ListSessions(ctx, req.(*ListSessionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_RevokeSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeSessionRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).RevokeSession(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/RevokeSession"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).RevokeSession(ctx, req.(*RevokeSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_UpsertCommTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertCommTemplateRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).UpsertCommTemplate(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/UpsertCommTemplate"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).UpsertCommTemplate(ctx, req.(*UpsertCommTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_ListCommTemplates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCommTemplatesRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).ListCommTemplates(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/ListCommTemplates"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).ListCommTemplates(ctx, req.(*ListCommTemplatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_TriggerBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerBackupRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).TriggerBackup(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/TriggerBackup"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).TriggerBackup(ctx, req.(*TriggerBackupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamSec_GetBackupHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBackupHistoryRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(IamSecurityHandler).GetBackupHistory(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/pb.IamSecurityService/GetBackupHistory"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamSecurityHandler).GetBackupHistory(ctx, req.(*GetBackupHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ── Registration ──────────────────────────────────────────────────────────────

// IamSecurityServiceDesc extends IamService with security RPCs.
var IamSecurityServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.IamSecurityService",
	HandlerType: (*IamSecurityHandler)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetPasswordPolicy", Handler: _IamSec_GetPasswordPolicy_Handler},
		{MethodName: "SetPasswordPolicy", Handler: _IamSec_SetPasswordPolicy_Handler},
		{MethodName: "RecordLoginAnomaly", Handler: _IamSec_RecordLoginAnomaly_Handler},
		{MethodName: "ListSessions", Handler: _IamSec_ListSessions_Handler},
		{MethodName: "RevokeSession", Handler: _IamSec_RevokeSession_Handler},
		{MethodName: "UpsertCommTemplate", Handler: _IamSec_UpsertCommTemplate_Handler},
		{MethodName: "ListCommTemplates", Handler: _IamSec_ListCommTemplates_Handler},
		{MethodName: "TriggerBackup", Handler: _IamSec_TriggerBackup_Handler},
		{MethodName: "GetBackupHistory", Handler: _IamSec_GetBackupHistory_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "iam_security.proto",
}

// RegisterIamSecurityHandler registers the security extension on a gRPC server.
func RegisterIamSecurityHandler(s grpc.ServiceRegistrar, srv IamSecurityHandler) {
	s.RegisterService(&IamSecurityServiceDesc, srv)
}
