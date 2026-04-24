package health_check_adapter

// Status mirrors grpc_health_v1.HealthCheckResponse_ServingStatus as a plain
// string so the proto type doesn't leak past the adapter boundary. Values
// are the canonical names from the standard gRPC health protocol; downstream
// handlers render them directly into the REST response.
type Status string

const (
	StatusServing    Status = "SERVING"
	StatusNotServing Status = "NOT_SERVING"
	StatusUnknown    Status = "UNKNOWN"
)

// BackendStatus is one row in the /v1/system/backends response envelope.
// Error is populated only when the Health.Check RPC itself failed (backend
// unreachable, timeout); when the RPC succeeded, Status carries the
// protocol-level verdict and Error is empty.
type BackendStatus struct {
	Name   string
	Status Status
	Error  string
}
