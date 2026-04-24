package finance_grpc_adapter

// PingResult is the adapter-layer result produced by Ping. Proto types do not
// leak past the adapter; the gateway's REST service layer (and handlers) see
// these plain Go types only.
//
// The caller identity is NOT in this result — finance-svc is identity-agnostic
// per ADR 0009 (single-point auth at gateway). The gateway's REST handler
// pulls user_id / branch_id / roles from its own c.Locals(IdentityKey) and
// assembles the client-visible response envelope directly.
type PingResult struct {
	Message string
}
