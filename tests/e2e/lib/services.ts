// Service registry — base URLs for all scaffolded services.
// Add new services here as they are scaffolded.

export interface ServiceEntry {
  name: string;
  baseURL: string;
}

export const services: ServiceEntry[] = [
  {
    name: "iam-svc",
    baseURL: process.env.IAM_SVC_URL || "http://localhost:4001",
  },
  // Uncomment as services are scaffolded:
  // { name: "catalog-svc", baseURL: process.env.CATALOG_SVC_URL || "http://localhost:4002" },
  // { name: "booking-svc", baseURL: process.env.BOOKING_SVC_URL || "http://localhost:4003" },
  // { name: "jamaah-svc",  baseURL: process.env.JAMAAH_SVC_URL  || "http://localhost:4004" },
  // { name: "payment-svc", baseURL: process.env.PAYMENT_SVC_URL || "http://localhost:4005" },
  // { name: "visa-svc",    baseURL: process.env.VISA_SVC_URL    || "http://localhost:4006" },
  // { name: "ops-svc",     baseURL: process.env.OPS_SVC_URL     || "http://localhost:4007" },
];
