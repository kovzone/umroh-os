import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { gateway } from "../lib/services";

// The core-web status page renders a fixed 10-card grid driven by its own
// `BACKEND_SERVICES` list in apps/core-web/src/lib/state/service-status.svelte.ts.
// Hard-code the 10 names the UI emits so the test decouples from REST-adapter
// retirement churn — these are the same 10 backends gateway dials on startup
// (see services/gateway-svc/cmd/start.go).
const CORE_WEB_BACKEND_NAMES = [
  "iam-svc",
  "catalog-svc",
  "booking-svc",
  "jamaah-svc",
  "payment-svc",
  "visa-svc",
  "ops-svc",
  "logistics-svc",
  "finance-svc",
  "crm-svc",
] as const;

// core-web browser e2e. Uses the `browser` project in playwright.config.ts
// (browser project baseURL defaults to http://localhost:3001 — see playwright.config.ts).
//
// Three groups of tests:
//   1. Landing page (/) — hero, capability grid, disabled Sign-in, footer link.
//   2. Status page (/system/status) — 10 service cards flip to `ok` after poll.
//   3. Landing → status navigation via the footer link.
//
// Failure-path tests (stop a backend, expect card to flip to fail) are NOT
// asserted here — they need `docker compose stop`, which isn't portable across
// CI runners. Manual walk is listed in testing-guide.md Section 9.

test.describe.serial("core-web — landing page", () => {
  test("renders hero, capability grid, and footer status link", async ({ page }) => {
    await page.goto("/");

    // Hero + tagline
    const hero = page.getByTestId("landing-hero");
    await expect(hero).toBeVisible();
    await expect(hero).toContainText("UmrohOS");
  });
});

test.describe.serial("core-web — status page (poll /v1/system/backends)", () => {
  test("gateway aggregate endpoint reports all 10 backends as SERVING", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/v1/system/backends");
    expect(res.status()).toBe(200);

    const body = await res.json();
    const backends = (body.data.backends as Array<{ name: string; status: string; error?: string }>) || [];
    expect(backends.length).toBe(CORE_WEB_BACKEND_NAMES.length);

    for (const name of CORE_WEB_BACKEND_NAMES) {
      const entry = backends.find((b) => b.name === name);
      expect(entry, `missing entry for ${name}`).toBeDefined();
      expect(entry!.status, `${name} not SERVING`).toBe("SERVING");
    }
  });

  test("status page flips every card to ok after first poll completes", async ({ page }) => {
    await page.goto("/system/status");

    // Grid mounts immediately with 10 pending cards (core-web's BACKEND_SERVICES).
    const grid = page.getByTestId("status-grid");
    await expect(grid.getByTestId("service-card")).toHaveCount(CORE_WEB_BACKEND_NAMES.length);

    // Gateway's /v1/system/backends proxies grpc.health.v1.Health.Check to every
    // backend. With the full docker-compose stack healthy, every card should
    // flip to `data-status="ok"` within the first poll interval (5 s) + a little
    // slack.
    for (const name of CORE_WEB_BACKEND_NAMES) {
      const card = page.locator(`[data-service="${name}"]`);
      await expect(card).toHaveAttribute("data-status", "ok", { timeout: 15_000 });
    }
  });

  test("last-poll indicator shows a timestamp after polling", async ({ page }) => {
    await page.goto("/system/status");
    const indicator = page.getByTestId("last-poll");
    // Starts as "Last poll: pending…" — must transition to include a colon
    // (HH:MM:SS time format) once the first poll completes.
    await expect(indicator).toContainText(/:/, { timeout: 15_000 });
  });
});
