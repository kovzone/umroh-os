// jamaah_ocr_grpc_ext.go — hand-written gRPC interface extension for OCR RPCs.
//
// Adds TriggerOCR and GetOCRStatus to the JamaahService gRPC service.
// Run `make genpb` after updating jamaah.proto to replace these stubs.
//
// BL-DOC-002 / S3-E-02.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	JamaahService_TriggerOCR_FullMethodName    = "/pb.JamaahService/TriggerOCR"
	JamaahService_GetOCRStatus_FullMethodName  = "/pb.JamaahService/GetOCRStatus"
)

// OCRHandler is the server-side interface for the OCR RPCs.
type OCRHandler interface {
	TriggerOCR(context.Context, *TriggerOCRRequest) (*TriggerOCRResponse, error)
	GetOCRStatus(context.Context, *GetOCRStatusRequest) (*GetOCRStatusResponse, error)
}

// UnimplementedOCRHandler provides safe defaults.
type UnimplementedOCRHandler struct{}

func (UnimplementedOCRHandler) TriggerOCR(context.Context, *TriggerOCRRequest) (*TriggerOCRResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerOCR not implemented")
}
func (UnimplementedOCRHandler) GetOCRStatus(context.Context, *GetOCRStatusRequest) (*GetOCRStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOCRStatus not implemented")
}

func _JamaahService_TriggerOCR_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerOCRRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OCRHandler).TriggerOCR(ctx, req.(*TriggerOCRRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: JamaahService_TriggerOCR_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _JamaahService_GetOCRStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOCRStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OCRHandler).GetOCRStatus(ctx, req.(*GetOCRStatusRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: JamaahService_GetOCRStatus_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

// RegisterJamaahServiceServerFull registers the full JamaahService including
// all hand-written extensions (Documents, Manifest, OCR).
func RegisterJamaahServiceServerFull(s grpc.ServiceRegistrar, srv interface {
	JamaahServiceServer
	DocumentHandler
	ManifestHandler
	OCRHandler
}) {
	desc := grpc.ServiceDesc{
		ServiceName: "pb.JamaahService",
		HandlerType: (*JamaahServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Healthz",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(HealthzRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(JamaahServiceServer).Healthz(ctx, req.(*HealthzRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: JamaahService_Healthz_FullMethodName}
					return interceptor(ctx, in, info, handler)
				},
			},
			{
				MethodName: "UploadDocument",
				Handler:    _JamaahService_UploadDocument_Handler,
			},
			{
				MethodName: "ReviewDocument",
				Handler:    _JamaahService_ReviewDocument_Handler,
			},
			{
				MethodName: "GetDepartureManifest",
				Handler:    _JamaahService_GetDepartureManifest_Handler,
			},
			{
				MethodName: "TriggerOCR",
				Handler:    _JamaahService_TriggerOCR_Handler,
			},
			{
				MethodName: "GetOCRStatus",
				Handler:    _JamaahService_GetOCRStatus_Handler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "jamaah.proto",
	}
	s.RegisterService(&desc, srv)
}
