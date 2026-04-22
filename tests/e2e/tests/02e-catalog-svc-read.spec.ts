import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { backendServices } from "../lib/services";

// S1-E-02 / BL-CAT-001 — catalog-svc public read smoke.
//
// Exercises the § Catalog contract in `docs/contracts/slice-S1.md`:
//
//   GET /v1/packages              — active-only list + filters + cursor
//   GET /v1/packages/{id}         — active-only detail with 404 on draft/archived
//
// Depends on the dev fixtures seeded by migration
// `000009_seed_catalog_dev_fixtures` — one active, one draft, one archived.
// Reviewer runs `make dev-bootstrap` (or a targeted `make migrate-up`)
// before this spec.

const catalog = backendServices.find((s) => s.name === "catalog-svc")!;

const ACTIVE_ID = "pkg_01JCDE00000000000000000001";
const DRAFT_ID = "pkg_01JCDE00000000000000000002";
const ARCHIVED_ID = "pkg_01JCDE00000000000000000003";
const UNKNOWN_ID = "pkg_01JCDE00000000000000000099";

test.describe.serial("catalog-svc — public read (S1-E-02 / BL-CAT-001)", () => {
  test("GET /v1/packages returns only active packages", async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get("/v1/packages");
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(Array.isArray(body.packages)).toBe(true);

    const ids = body.packages.map((p: { id: string }) => p.id);
    expect(ids).toContain(ACTIVE_ID);
    expect(ids).not.toContain(DRAFT_ID);
    expect(ids).not.toContain(ARCHIVED_ID);

    const active = body.packages.find((p: { id: string }) => p.id === ACTIVE_ID);
    expect(active.kind).toBe("umrah_reguler");
    expect(active.starting_price.list_currency).toBe("IDR");
    expect(active.starting_price.settlement_currency).toBe("IDR");
    expect(active.starting_price.list_amount).toBeGreaterThan(0);
    expect(active.next_departure).toBeTruthy();
    expect(active.next_departure.remaining_seats).toBeGreaterThanOrEqual(0);

    expect(body.page.has_more).toBe(false);
  });

  test("GET /v1/packages?kind=umrah_reguler filters", async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get("/v1/packages?kind=umrah_reguler");
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.packages.length).toBeGreaterThanOrEqual(1);
    for (const p of body.packages) {
      expect(p.kind).toBe("umrah_reguler");
    }
  });

  test("GET /v1/packages?kind=banana returns invalid_query_param 400", async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get("/v1/packages?kind=banana");
    expect(res.status()).toBe(400);
    const body = await res.json();
    expect(body.error.code).toBe("invalid_query_param");
    expect(typeof body.error.message).toBe("string");
    expect(typeof body.error.trace_id).toBe("string");
  });

  test("GET /v1/packages?limit=0 returns invalid_query_param 400", async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get("/v1/packages?limit=0");
    expect(res.status()).toBe(400);
    const body = await res.json();
    expect(body.error.code).toBe("invalid_query_param");
  });

  test("GET /v1/packages?cursor=garbage returns invalid_cursor 400", async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get("/v1/packages?cursor=garbage");
    expect(res.status()).toBe(400);
    const body = await res.json();
    expect(body.error.code).toBe("invalid_cursor");
  });

  test(`GET /v1/packages/{active-id} returns full detail`, async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get(`/v1/packages/${ACTIVE_ID}`);
    expect(res.status()).toBe(200);
    const body = await res.json();

    expect(body.package.id).toBe(ACTIVE_ID);
    expect(body.package.kind).toBe("umrah_reguler");
    expect(Array.isArray(body.package.highlights)).toBe(true);
    expect(Array.isArray(body.package.hotels)).toBe(true);
    expect(Array.isArray(body.package.add_ons)).toBe(true);
    expect(Array.isArray(body.package.departures)).toBe(true);

    expect(body.package.airline).toBeTruthy();
    expect(body.package.muthawwif).toBeTruthy();
    expect(body.package.itinerary).toBeTruthy();
    expect(body.package.itinerary.days.length).toBeGreaterThan(0);

    // Every surfaced departure is open or closed (never departed/completed/cancelled).
    for (const dep of body.package.departures) {
      expect(["open", "closed"]).toContain(dep.status);
      expect(dep.remaining_seats).toBeGreaterThanOrEqual(0);
    }
  });

  test(`GET /v1/packages/{draft-id} returns 404 package_not_found`, async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get(`/v1/packages/${DRAFT_ID}`);
    expect(res.status()).toBe(404);
    const body = await res.json();
    expect(body.error.code).toBe("package_not_found");
    expect(typeof body.error.trace_id).toBe("string");
  });

  test(`GET /v1/packages/{archived-id} returns 404 package_not_found`, async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get(`/v1/packages/${ARCHIVED_ID}`);
    expect(res.status()).toBe(404);
    const body = await res.json();
    expect(body.error.code).toBe("package_not_found");
  });

  test(`GET /v1/packages/{unknown-id} returns 404 package_not_found`, async () => {
    const api = await createApiClient(catalog.baseURL);
    const res = await api.get(`/v1/packages/${UNKNOWN_ID}`);
    expect(res.status()).toBe(404);
    const body = await res.json();
    expect(body.error.code).toBe("package_not_found");
  });
});
