# catalog-svc — Status

## Implementation checklist

- [x] Scaffolded from baseline template (S0)
- [x] Wired into `docker-compose.dev.yml` (S0)
- [x] Initial DDL (`catalog` schema + 10 tables + 5 enums) — migration `000008` (S1-E-02 / BL-CAT-001, 2026-04-22)
- [x] sqlc queries for packages, departures, hotels, airlines, muthawwif, itineraries, addons — read-side only (BL-CAT-001)
- [x] OpenAPI spec for public catalog read (list + detail) (BL-CAT-001)
- [x] CRUD handlers for packages — **read only** (list + active detail); admin write endpoints land in BL-CAT-005..011
- [ ] CRUD handlers for hotels / airlines / muthawwif — admin CRUD deferred
- [ ] `GET /v1/package-departures/{id}` with live `remaining_seats` — BL-CAT-002
- [ ] Atomic seat reservation (gRPC `ReserveSeats` / `ReleaseSeats`) — S1-E-03 territory
- [ ] Bulk import / export — BL-CAT-010 / BL-CAT-011
- [x] Unit tests (cursor helpers) + integration/e2e (Playwright `02e-catalog-svc-read.spec.ts`)
- [ ] Verified by reviewer

## Current status

**S1-E-02 / BL-CAT-001 in progress** on branch `feat/s1-e-02-catalog-read-packages` (2026-04-22). Public catalog read (list + active detail) compiles, runs in dev, and passes the nine e2e assertions plus curl smoke. Awaiting reviewer sign-off, then merge to `dev`. Departure-detail (`BL-CAT-002`) and the B2C integration wrap-up (`BL-CAT-004`) are the next cards under the same task code.
