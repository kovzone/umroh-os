import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { backendServices, gateway } from "../lib/services";

test.describe.serial("gateway-svc — own probes + per-backend proxies", () => {
  test(`GET ${gateway.name} /system/live returns ok`, async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/system/live");
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.ok).toBe(true);
  });

  test(`GET ${gateway.name} /system/ready returns ok`, async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/system/ready");
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.ok).toBe(true);
  });

  // One proxy probe per backend — proves the REST adapter is wired and
  // upstream service is reachable from inside the docker network.
  for (const svc of backendServices) {
    test(`GET ${gateway.name} /v1/${svc.shortName}/system/live proxies upstream`, async () => {
      const api = await createApiClient(gateway.baseURL);
      const res = await api.get(`/v1/${svc.shortName}/system/live`);
      expect(res.status()).toBe(200);
      const body = await res.json();
      expect(body.data.ok).toBe(true);
    });
  }
});
