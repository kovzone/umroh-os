// Service registry — base URLs for every scaffolded service.
//
// Two registries:
//   - backendServices: the 10 stateful services. Each exposes the standard
//     scaffold endpoints /system/live, /system/ready, /system/diagnostics/db-tx.
//   - gateway: the edge proxy. Exposes /system/live + /system/ready locally
//     plus one /v1/<shortName>/system/live per backend it fronts.

export interface ServiceEntry {
  name: string;       // "iam-svc"
  shortName: string;  // "iam" — used for gateway proxy paths /v1/<shortName>/...
  baseURL: string;
}

// Default to 127.0.0.1 (not "localhost") so Node/Playwright on Windows resolves to IPv4.
// Docker Desktop publishes host ports on IPv4; ::1 often gets ECONNREFUSED for the same port.
export const backendServices: ServiceEntry[] = [
  { name: "iam-svc",       shortName: "iam",       baseURL: process.env.IAM_SVC_URL       || "http://127.0.0.1:4001" },
  { name: "catalog-svc",   shortName: "catalog",   baseURL: process.env.CATALOG_SVC_URL   || "http://127.0.0.1:4002" },
  { name: "booking-svc",   shortName: "booking",   baseURL: process.env.BOOKING_SVC_URL   || "http://127.0.0.1:4003" },
  { name: "jamaah-svc",    shortName: "jamaah",    baseURL: process.env.JAMAAH_SVC_URL    || "http://127.0.0.1:4004" },
  { name: "payment-svc",   shortName: "payment",   baseURL: process.env.PAYMENT_SVC_URL   || "http://127.0.0.1:4005" },
  { name: "visa-svc",      shortName: "visa",      baseURL: process.env.VISA_SVC_URL      || "http://127.0.0.1:4006" },
  { name: "ops-svc",       shortName: "ops",       baseURL: process.env.OPS_SVC_URL       || "http://127.0.0.1:4007" },
  { name: "logistics-svc", shortName: "logistics", baseURL: process.env.LOGISTICS_SVC_URL || "http://127.0.0.1:4008" },
  { name: "finance-svc",   shortName: "finance",   baseURL: process.env.FINANCE_SVC_URL   || "http://127.0.0.1:4009" },
  { name: "crm-svc",       shortName: "crm",       baseURL: process.env.CRM_SVC_URL       || "http://127.0.0.1:4010" },
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
