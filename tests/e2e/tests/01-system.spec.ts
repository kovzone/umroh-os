import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { clearState } from "../lib/state";
import { services } from "../lib/services";

test.describe.serial("System probes", () => {
  test("clear state from previous run", () => {
    clearState();
  });

  for (const svc of services) {
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
  }
});
