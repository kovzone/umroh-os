// jamaah_grpc_ext.go — hand-written gRPC interface extension for S3-E-02.
//
// Extends JamaahServiceServer with UploadDocument and ReviewDocument.
// Run `make genpb` after updating jamaah.proto to replace these stubs.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	JamaahService_UploadDocument_FullMethodName      = "/pb.JamaahService/UploadDocument"
	JamaahService_ReviewDocument_FullMethodName      = "/pb.JamaahService/ReviewDocument"
	JamaahService_GetDepartureManifest_FullMethodName = "/pb.JamaahService/GetDepartureManifest"
)

// DocumentHandler is the server-side interface for the document RPCs.
type DocumentHandler interface {
	UploadDocument(context.Context, *UploadDocumentRequest) (*UploadDocumentResponse, error)
	ReviewDocument(context.Context, *ReviewDocumentRequest) (*ReviewDocumentResponse, error)
}

// UnimplementedDocumentHandler provides safe defaults.
type UnimplementedDocumentHandler struct{}

func (UnimplementedDocumentHandler) UploadDocument(context.Context, *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadDocument not implemented")
}
func (UnimplementedDocumentHandler) ReviewDocument(context.Context, *ReviewDocumentRequest) (*ReviewDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReviewDocument not implemented")
}

// ManifestHandler is the server-side interface for manifest RPCs.
type ManifestHandler interface {
	GetDepartureManifest(context.Context, *GetDepartureManifestRequest) (*GetDepartureManifestResponse, error)
}

// UnimplementedManifestHandler provides safe defaults.
type UnimplementedManifestHandler struct{}

func (UnimplementedManifestHandler) GetDepartureManifest(context.Context, *GetDepartureManifestRequest) (*GetDepartureManifestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDepartureManifest not implemented")
}

func _JamaahService_UploadDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadDocumentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DocumentHandler).UploadDocument(ctx, req.(*UploadDocumentRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: JamaahService_UploadDocument_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _JamaahService_ReviewDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReviewDocumentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DocumentHandler).ReviewDocument(ctx, req.(*ReviewDocumentRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: JamaahService_ReviewDocument_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _JamaahService_GetDepartureManifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDepartureManifestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManifestHandler).GetDepartureManifest(ctx, req.(*GetDepartureManifestRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: JamaahService_GetDepartureManifest_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

// RegisterJamaahServiceServerWithExtensions registers the combined JamaahService.
func RegisterJamaahServiceServerWithExtensions(s grpc.ServiceRegistrar, srv interface {
	JamaahServiceServer
	DocumentHandler
	ManifestHandler
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
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "jamaah.proto",
	}
	s.RegisterService(&desc, srv)
}
