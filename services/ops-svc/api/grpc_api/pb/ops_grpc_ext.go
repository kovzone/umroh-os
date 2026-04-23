// ops_grpc_ext.go — hand-written gRPC service extension for ops-svc RPCs.
//
// Extends OpsServiceServer with:
//   - RunRoomAllocation  (BL-OPS-002)
//   - GetRoomAllocation  (BL-OPS-002)
//   - GenerateIDCard     (BL-OPS-003)
//   - VerifyIDCard       (BL-OPS-003)
//   - ExportManifest     (BL-OPS-001)
//
// Run `make genpb` after updating ops.proto to replace with generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Full method name constants for ops RPCs.
const (
	OpsService_RunRoomAllocation_FullMethodName = "/pb.OpsService/RunRoomAllocation"
	OpsService_GetRoomAllocation_FullMethodName = "/pb.OpsService/GetRoomAllocation"
	OpsService_GenerateIDCard_FullMethodName    = "/pb.OpsService/GenerateIDCard"
	OpsService_VerifyIDCard_FullMethodName      = "/pb.OpsService/VerifyIDCard"
	OpsService_ExportManifest_FullMethodName    = "/pb.OpsService/ExportManifest"
)

// OpsServiceHandler is the server-side interface for all ops RPCs.
type OpsServiceHandler interface {
	RunRoomAllocation(context.Context, *RunRoomAllocationRequest) (*RunRoomAllocationResponse, error)
	GetRoomAllocation(context.Context, *GetRoomAllocationRequest) (*GetRoomAllocationResponse, error)
	GenerateIDCard(context.Context, *GenerateIDCardRequest) (*GenerateIDCardResponse, error)
	VerifyIDCard(context.Context, *VerifyIDCardRequest) (*VerifyIDCardResponse, error)
	ExportManifest(context.Context, *ExportManifestRequest) (*ExportManifestResponse, error)
}

// UnimplementedOpsServiceHandler provides safe defaults for OpsServiceHandler.
type UnimplementedOpsServiceHandler struct{}

func (UnimplementedOpsServiceHandler) RunRoomAllocation(context.Context, *RunRoomAllocationRequest) (*RunRoomAllocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunRoomAllocation not implemented")
}

func (UnimplementedOpsServiceHandler) GetRoomAllocation(context.Context, *GetRoomAllocationRequest) (*GetRoomAllocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomAllocation not implemented")
}

func (UnimplementedOpsServiceHandler) GenerateIDCard(context.Context, *GenerateIDCardRequest) (*GenerateIDCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateIDCard not implemented")
}

func (UnimplementedOpsServiceHandler) VerifyIDCard(context.Context, *VerifyIDCardRequest) (*VerifyIDCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyIDCard not implemented")
}

func (UnimplementedOpsServiceHandler) ExportManifest(context.Context, *ExportManifestRequest) (*ExportManifestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportManifest not implemented")
}

// ---------------------------------------------------------------------------
// Handler dispatch helpers
// ---------------------------------------------------------------------------

func _OpsService_RunRoomAllocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunRoomAllocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServiceHandler).RunRoomAllocation(ctx, req.(*RunRoomAllocationRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RunRoomAllocation_FullMethodName}, handler)
}

func _OpsService_GetRoomAllocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomAllocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServiceHandler).GetRoomAllocation(ctx, req.(*GetRoomAllocationRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetRoomAllocation_FullMethodName}, handler)
}

func _OpsService_GenerateIDCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateIDCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServiceHandler).GenerateIDCard(ctx, req.(*GenerateIDCardRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GenerateIDCard_FullMethodName}, handler)
}

func _OpsService_VerifyIDCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyIDCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServiceHandler).VerifyIDCard(ctx, req.(*VerifyIDCardRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_VerifyIDCard_FullMethodName}, handler)
}

func _OpsService_ExportManifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportManifestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServiceHandler).ExportManifest(ctx, req.(*ExportManifestRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_ExportManifest_FullMethodName}, handler)
}

// gRPC method descriptors.
var OpsService_RunRoomAllocation_MethodDesc = grpc.MethodDesc{
	MethodName: "RunRoomAllocation",
	Handler:    _OpsService_RunRoomAllocation_Handler,
}

var OpsService_GetRoomAllocation_MethodDesc = grpc.MethodDesc{
	MethodName: "GetRoomAllocation",
	Handler:    _OpsService_GetRoomAllocation_Handler,
}

var OpsService_GenerateIDCard_MethodDesc = grpc.MethodDesc{
	MethodName: "GenerateIDCard",
	Handler:    _OpsService_GenerateIDCard_Handler,
}

var OpsService_VerifyIDCard_MethodDesc = grpc.MethodDesc{
	MethodName: "VerifyIDCard",
	Handler:    _OpsService_VerifyIDCard_Handler,
}

var OpsService_ExportManifest_MethodDesc = grpc.MethodDesc{
	MethodName: "ExportManifest",
	Handler:    _OpsService_ExportManifest_Handler,
}

// RegisterOpsServiceServerFull registers the combined OpsService:
// generated Healthz RPC + all hand-written ops RPCs.
func RegisterOpsServiceServerFull(s grpc.ServiceRegistrar, srv interface {
	OpsServiceServer
	OpsServiceHandler
}) {
	desc := grpc.ServiceDesc{
		ServiceName: "pb.OpsService",
		HandlerType: (*OpsServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Healthz",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(HealthzRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(OpsServiceServer).Healthz(ctx, req.(*HealthzRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_Healthz_FullMethodName}
					return interceptor(ctx, in, info, handler)
				},
			},
			OpsService_RunRoomAllocation_MethodDesc,
			OpsService_GetRoomAllocation_MethodDesc,
			OpsService_GenerateIDCard_MethodDesc,
			OpsService_VerifyIDCard_MethodDesc,
			OpsService_ExportManifest_MethodDesc,
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "ops.proto",
	}
	s.RegisterService(&desc, srv)
}
