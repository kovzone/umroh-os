// jamaah_manifest_stub.go — gateway-side gRPC client stub for jamaah-svc
// manifest RPC (Wave 1A / Phase 6).
//
// Mirrors services/jamaah-svc/api/grpc_api/pb/manifest_messages.go and
// the GetDepartureManifest RPC in jamaah_grpc_ext.go.
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	JamaahService_GetDepartureManifest_FullMethodName = "/pb.JamaahService/GetDepartureManifest"
)

// ---------------------------------------------------------------------------
// Domain objects
// ---------------------------------------------------------------------------

// ManifestJamaah represents one pilgrim row in the departure manifest.
type ManifestJamaah struct {
	BookingID     string
	Name          string
	NIK           string
	Phone         string
	RoomType      string
	BookingStatus string
	DocStatus     string
}

func (x *ManifestJamaah) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *ManifestJamaah) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *ManifestJamaah) GetNIK() string {
	if x == nil {
		return ""
	}
	return x.NIK
}
func (x *ManifestJamaah) GetPhone() string {
	if x == nil {
		return ""
	}
	return x.Phone
}
func (x *ManifestJamaah) GetRoomType() string {
	if x == nil {
		return ""
	}
	return x.RoomType
}
func (x *ManifestJamaah) GetBookingStatus() string {
	if x == nil {
		return ""
	}
	return x.BookingStatus
}
func (x *ManifestJamaah) GetDocStatus() string {
	if x == nil {
		return ""
	}
	return x.DocStatus
}

// ---------------------------------------------------------------------------
// Request / response types
// ---------------------------------------------------------------------------

type GetDepartureManifestRequest struct {
	DepartureID string
}

func (x *GetDepartureManifestRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type GetDepartureManifestResponse struct {
	DepartureID string
	TotalJamaah int32
	LunasPaid   int32
	DocComplete int32
	JamaahList  []*ManifestJamaah
}

func (x *GetDepartureManifestResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetDepartureManifestResponse) GetTotalJamaah() int32 {
	if x == nil {
		return 0
	}
	return x.TotalJamaah
}
func (x *GetDepartureManifestResponse) GetLunasPaid() int32 {
	if x == nil {
		return 0
	}
	return x.LunasPaid
}
func (x *GetDepartureManifestResponse) GetDocComplete() int32 {
	if x == nil {
		return 0
	}
	return x.DocComplete
}
func (x *GetDepartureManifestResponse) GetJamaahList() []*ManifestJamaah {
	if x == nil {
		return nil
	}
	return x.JamaahList
}

// ---------------------------------------------------------------------------
// ManifestClient interface + implementation
// ---------------------------------------------------------------------------

// ManifestClient is the consumer-side interface for jamaah-svc manifest RPCs.
type ManifestClient interface {
	GetDepartureManifest(ctx context.Context, in *GetDepartureManifestRequest, opts ...grpc.CallOption) (*GetDepartureManifestResponse, error)
}

type manifestClient struct {
	cc grpc.ClientConnInterface
}

// NewManifestClient wraps a conn so gateway-svc can call jamaah manifest RPCs.
func NewManifestClient(cc grpc.ClientConnInterface) ManifestClient {
	return &manifestClient{cc}
}

func (c *manifestClient) GetDepartureManifest(ctx context.Context, in *GetDepartureManifestRequest, opts ...grpc.CallOption) (*GetDepartureManifestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDepartureManifestResponse)
	err := c.cc.Invoke(ctx, JamaahService_GetDepartureManifest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
