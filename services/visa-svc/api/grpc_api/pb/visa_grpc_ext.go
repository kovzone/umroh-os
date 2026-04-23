// visa_grpc_ext.go — extension of the generated VisaService gRPC server interface
// to include the Phase 6 RPCs (BL-VISA-001..003).
//
// The actual .proto / generated code is in visa_grpc.pb.go and visa.pb.go.
// The hand-written extension methods declared here satisfy the standard scaffold
// pattern used across all services in this repo.

package pb

import "context"

// VisaHandler is the interface that visa-svc's gRPC server must implement for
// the Phase 6 visa pipeline RPCs.
type VisaHandler interface {
	TransitionStatus(ctx context.Context, req *TransitionStatusRequest) (*TransitionStatusResponse, error)
	BulkSubmit(ctx context.Context, req *BulkSubmitRequest) (*BulkSubmitResponse, error)
	GetApplications(ctx context.Context, req *GetApplicationsRequest) (*GetApplicationsResponse, error)
}
