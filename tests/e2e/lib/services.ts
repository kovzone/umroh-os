// Service registry — base URLs for every REST surface in the stack.
//
// Post BL-IAM-019 / S1-E-14 every downstream backend is gRPC-only per
// ADR 0009; no spec hits a backend directly anymore. The only REST surface
// is gateway-svc, which fronts every backend via its per-service gRPC
// adapter. Probe individual backends via
// `grpc_health_probe -addr=<svc>:<grpc-port>` (docker-compose healthcheck
// uses this too, per BL-MON-001).
//
// `backendServices` is kept as an empty-by-default export purely so existing
// iteration patterns (if any land later) don't break. All current specs use
// `gateway`.

export interface ServiceEntry {
  name: string;
  shortName: string;
  baseURL: string;
}

export const backendServices: ServiceEntry[] = [];

export interface GatewayEntry {
  name: string;
  baseURL: string;
}

// Default to 127.0.0.1 (not "localhost") so Node/Playwright on Windows resolves to IPv4.
// Docker Desktop publishes host ports on IPv4; ::1 often gets ECONNREFUSED for the same port.
export const gateway: GatewayEntry = {
  name: "gateway-svc",
  baseURL: process.env.GATEWAY_SVC_URL || "http://127.0.0.1:4000",
};

// Backwards-compatible alias for callers that still import `services`.
export const services = backendServices;
