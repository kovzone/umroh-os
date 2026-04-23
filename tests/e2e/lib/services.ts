// Service registry — base URLs for every backend that still exposes a REST
// surface (the ADR 0009 sweep has moved most backends to gRPC-only).
//
// Two registries:
//   - backendServices: backends whose /system/live, /system/ready,
//     /system/diagnostics/db-tx probes are still REST-accessible. The seven
//     pure-scaffold backends (booking/crm/jamaah/logistics/ops/payment/visa)
//     retired their REST surfaces in BL-REFACTOR-002..008 / S1-E-13; probe
//     them via `grpc_health_probe -addr=<svc>:<grpc-port>` instead. finance-svc
//     follows in BL-IAM-019 / S1-E-14. catalog-svc moved earlier in
//     BL-REFACTOR-001 (G7).
//   - gateway: the edge proxy. Exposes /system/live + /system/ready locally
//     plus the few /v1/<shortName>/system/live proxies that are still wired
//     (currently iam + finance).

export interface ServiceEntry {
  name: string;       // "iam-svc"
  shortName: string;  // "iam" — used for gateway proxy paths /v1/<shortName>/...
  baseURL: string;
}

// Default to 127.0.0.1 (not "localhost") so Node/Playwright on Windows resolves to IPv4.
// Docker Desktop publishes host ports on IPv4; ::1 often gets ECONNREFUSED for the same port.
export const backendServices: ServiceEntry[] = [
  { name: "iam-svc",     shortName: "iam",     baseURL: process.env.IAM_SVC_URL     || "http://127.0.0.1:4001" },
  { name: "finance-svc", shortName: "finance", baseURL: process.env.FINANCE_SVC_URL || "http://127.0.0.1:4009" },
];

export interface GatewayEntry {
  name: string;
  baseURL: string;
}

export const gateway: GatewayEntry = {
  name: "gateway-svc",
  baseURL: process.env.GATEWAY_SVC_URL || "http://127.0.0.1:4000",
};

// Backwards-compatible alias for callers that still import `services`.
export const services = backendServices;
