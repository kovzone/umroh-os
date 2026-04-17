import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { clearState } from "../lib/state";
import { backendServices } from "../lib/services";

test.describe.serial("Backend system probes", () => {
  test("clear state from previous run", () => {
    clearState();
  });

  for (const svc of backendServices) {
    test(`GET ${svc.name} /system/live returns ok`, async () => {
      const api = await createApiClient(svc.baseURL);
      const res = await api.get("/system/live");
      expect(res.status()).toBe(200);
      const body = await res.json();
      expect(body.data.ok).toBe(true);
    });

    test(`GET ${svc.name} /system/ready returns ok`, async () => {
      const api = await createApiClient(svc.baseURL);
      const res = await api.get("/system/ready");
      expect(res.status()).toBe(200);
      const body = await res.json();
      expect(body.data.ok).toBe(true);
    });

    test(`GET ${svc.name} /system/diagnostics/db-tx commits a row`, async () => {
      const api = await createApiClient(svc.baseURL);
      const message = `e2e-${svc.name}-${Date.now()}`;
      const res = await api.get("/system/diagnostics/db-tx", { message });
      expect(res.status()).toBe(200);
      const body = await res.json();
      expect(body.data.diagnostic_id).toBeGreaterThan(0);
      expect(body.data.message).toContain(message);
    });
  }
});
