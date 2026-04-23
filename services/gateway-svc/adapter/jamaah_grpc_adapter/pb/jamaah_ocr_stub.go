// jamaah_ocr_stub.go — gateway-side gRPC client stub for jamaah-svc OCR RPCs
// (S3 Wave 2).
//
// Mirrors services/jamaah-svc/api/grpc_api/pb/jamaah_ocr_messages.go.
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	JamaahService_TriggerOCR_FullMethodName   = "/pb.JamaahService/TriggerOCR"
	JamaahService_GetOCRStatus_FullMethodName = "/pb.JamaahService/GetOCRStatus"
)

// ---------------------------------------------------------------------------
// Request / response types
// ---------------------------------------------------------------------------

type TriggerOCRRequest struct {
	DocumentId string
}

func (x *TriggerOCRRequest) GetDocumentId() string {
	if x == nil {
		return ""
	}
	return x.DocumentId
}

type TriggerOCRResponse struct {
	DocumentId string
	Status     string
	Confidence float64
	OcrResult  map[string]string
}

func (x *TriggerOCRResponse) GetDocumentId() string {
	if x == nil {
		return ""
	}
	return x.DocumentId
}
func (x *TriggerOCRResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *TriggerOCRResponse) GetConfidence() float64 {
	if x == nil {
		return 0
	}
	return x.Confidence
}
func (x *TriggerOCRResponse) GetOcrResult() map[string]string {
	if x == nil {
		return nil
	}
	return x.OcrResult
}

type GetOCRStatusRequest struct {
	DocumentId string
}

func (x *GetOCRStatusRequest) GetDocumentId() string {
	if x == nil {
		return ""
	}
	return x.DocumentId
}

type GetOCRStatusResponse struct {
	DocumentId string
	Status     string
	Confidence float64
	OcrResult  map[string]string
}

func (x *GetOCRStatusResponse) GetDocumentId() string {
	if x == nil {
		return ""
	}
	return x.DocumentId
}
func (x *GetOCRStatusResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *GetOCRStatusResponse) GetConfidence() float64 {
	if x == nil {
		return 0
	}
	return x.Confidence
}
func (x *GetOCRStatusResponse) GetOcrResult() map[string]string {
	if x == nil {
		return nil
	}
	return x.OcrResult
}

// ---------------------------------------------------------------------------
// JamaahOCRClient interface + implementation
// ---------------------------------------------------------------------------

// JamaahOCRClient is the consumer-side interface for jamaah-svc OCR RPCs.
type JamaahOCRClient interface {
	TriggerOCR(ctx context.Context, req *TriggerOCRRequest, opts ...grpc.CallOption) (*TriggerOCRResponse, error)
	GetOCRStatus(ctx context.Context, req *GetOCRStatusRequest, opts ...grpc.CallOption) (*GetOCRStatusResponse, error)
}

type jamaahOCRClient struct {
	cc grpc.ClientConnInterface
}

// NewJamaahOCRClient wraps a conn so gateway-svc can call jamaah-svc OCR RPCs.
func NewJamaahOCRClient(cc grpc.ClientConnInterface) JamaahOCRClient {
	return &jamaahOCRClient{cc}
}

func (c *jamaahOCRClient) TriggerOCR(ctx context.Context, in *TriggerOCRRequest, opts ...grpc.CallOption) (*TriggerOCRResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TriggerOCRResponse)
	err := c.cc.Invoke(ctx, JamaahService_TriggerOCR_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jamaahOCRClient) GetOCRStatus(ctx context.Context, in *GetOCRStatusRequest, opts ...grpc.CallOption) (*GetOCRStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetOCRStatusResponse)
	err := c.cc.Invoke(ctx, JamaahService_GetOCRStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
