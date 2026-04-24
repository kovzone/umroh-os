import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { clearState } from "../lib/state";
import { gateway } from "../lib/services";

// Post ADR 0009 + BL-IAM-019 / S1-E-14 every downstream backend is gRPC-only,
// so system probes are not REST-reachable on individual backends. The only
// REST `/system/*` surface is gateway-svc's own probes. Per-backend liveness
// is served by the standard grpc.health.v1.Health protocol and consumed by
// docker-compose healthchecks via `grpc_health_probe` (BL-MON-001).

test.describe.serial("gateway-svc system probes", () => {
  test("clear state from previous run", () => {
    clearState();
  });

  test("GET gateway-svc /system/live returns ok", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/system/live");
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.ok).toBe(true);
  });

  test("GET gateway-svc /system/ready returns ok", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/system/ready");
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.ok).toBe(true);
  });
});
