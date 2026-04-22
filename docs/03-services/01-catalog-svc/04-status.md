# catalog-svc — Status

## Implementation checklist

- [x] Scaffolded from baseline template (S0)
- [x] Wired into `docker-compose.dev.yml` (S0)
- [x] Initial DDL (`catalog` schema + 10 tables + 5 enums) — migration `000008` (S1-E-02 / BL-CAT-001, 2026-04-22)
- [x] sqlc queries for packages, departures, hotels, airlines, muthawwif, itineraries, addons — read-side only (BL-CAT-001)
- [x] OpenAPI spec for public catalog read (list + detail) (BL-CAT-001)
- [x] CRUD handlers for packages — **read only** (list + active detail); admin write endpoints land in BL-CAT-005..011
- [ ] CRUD handlers for hotels / airlines / muthawwif — admin CRUD deferred
- [x] `GET /v1/package-departures/{id}` with live `remaining_seats` — BL-CAT-002
- [ ] Atomic seat reservation (gRPC `ReserveSeats` / `ReleaseSeats`) — S1-E-03 territory
- [ ] Bulk import / export — BL-CAT-010 / BL-CAT-011
- [x] Unit tests (cursor helpers) + integration/e2e (Playwright `02e-catalog-svc-read.spec.ts`)
- [x] Verified by reviewer — BL-CAT-001 merged to `dev` on 2026-04-22 (PR #38)

## Current status

**S1-E-02 / BL-CAT-001 + BL-CAT-002 both in flight** (branch `feat/s1-e-02-catalog-departures`, 2026-04-22). Public catalog read (list + detail + departure detail) is fully live in dev. B2C integration wrap-up (`BL-CAT-004`) is the next card under the same task code.

## 2026-04-22 — what landed with BL-CAT-001

- **DDL + enums (migration `000008`).** `catalog` schema with 10 tables and 5 enums: `packages`, `package_departures`, `package_pricing`, `hotels`, `airlines`, `muthawwif`, `itinerary_templates`, `addons`, and the two join tables (`package_hotels`, `package_addons`). ULID `TEXT` PKs with type-prefix CHECKs (`pkg_`, `dep_`, `pkgpr_`, `itn_`, `htl_`, `arl_`, `mtw_`, `addon_`) — deliberate deviation from the repo-default UUID PK, forced by § Catalog's ULID requirement; documented inline in the migration header + data-model doc. Seven-value `package_kind` enum per § Catalog. `CHECK (settlement_currency = 'IDR')` on `package_pricing` + `addons` per Q001; `CHECK (reserved_seats <= total_seats)` as defence-in-depth on departures (atomic guard lands with `S1-E-03`).
- **Dev seed (migration `000009`).** Idempotent seed: 1 active + 1 draft + 1 archived package with full master tree (2 hotels, 1 airline, 1 muthawwif, 1 itinerary with 3 days, 2 addons, 2 upcoming departures, 6 pricing rows). Drives e2e filter / 404 assertions.
- **Public REST surface.** `GET /v1/packages` — active-only list, cursor-paginated, filterable by `kind` / `departure_from/to` / `airline_code` / `hotel_id`, limit 1..100. `GET /v1/packages/{id}` — active-only detail with eager master refs and upcoming open/closed departures. Both match § Catalog wire shapes 1:1. Active-only filter is a top-level unconditional SQL predicate — cannot be bypassed by any filter combination. 404 `package_not_found` is emitted for draft/archived/unknown with identical response shape (no existence oracle).
- **Contract-exact error envelope.** Catalog handlers emit snake_case codes (`package_not_found`, `invalid_query_param`, `invalid_cursor`, `internal_error`), Bahasa messages (`id-ID` per Q003), and OTel `trace_id` pulled from the span context. The shared `util/apperrors` + middleware stays UPPERCASE for system endpoints — catalog handlers write their response directly via `c.Status().JSON()` returning `nil`, bypassing the middleware's raw-`err.Error()` leak path. Cross-service alignment of the shared envelope to the contract shape is follow-up.
- **Handler-level query-param validation.** oapi-codegen does not enforce enum / numeric-range constraints on query parameters by default. Discovered when `?kind=banana` bubbled a Postgres enum cast failure up to 500, and when `?limit=0` silently defaulted to 20. Added `isKnownPackageKind` whitelist + explicit `*params.Limit < 1 || > 100` guard in `api/rest_oapi/packages.go`.
- **Cursor format.** Base64-encoded `{"last_id": "..."}` JSON. Fixed-shape struct with a single `string` field — no custom `UnmarshalJSON`, no gadget-chain surface. Decoded `last_id` reaches SQL strictly as a bound parameter via `WHERE p.id > $cursor`. 4 unit tests lock the wire-format invariants (round-trip, empty, malformed, missing-last-id).
- **sqlc query split.** An initial single-query lateral-join shape (list + next_departure + starting_price in one round trip) was refactored into three simpler queries after sqlc mis-inferred nullability on the lateral-joined columns — would have crashed the row scan when a package had no upcoming departure. `ListActivePackages` + per-package `GetNextDepartureForPackage` + `GetStartingPriceForPackage` is N+2 queries per page but clean and error-safe. `list_amount` cast `numeric::bigint` in SQL for wire-integer conformance per Q001.
- **e2e coverage (`02e-catalog-svc-read.spec.ts`).** 9 assertions under the `api` project (no browser): active-only list, `kind=umrah_reguler` filter, 3 × 400 validation (bad kind, bad limit, bad cursor), full detail happy path (airline + muthawwif + itinerary + 2 departures + 2 hotels + 2 addons + 3 highlights), 3 × 404 `package_not_found` cases. Depends on migration `000009` fixtures.
- **Contract-first gate.** `§ Catalog` was already locked on 2026-04-20 via `S1-J-01`; this card is an exact wire-shape implementation — no contract changes.
- **Mid-session dev rebase.** Lutfi merged PR #36 (booking wizard) and PR #37 (staff catalog CRUD MVP docs) during the session. Rebased 8 commits onto `origin/dev`; git's 3-way merge auto-resolved both overlapping files (`docs/00-overview/06-feature-to-backlog-mapping.md` — my `in progress` flip coexists with Lutfi's new BL-CAT-014 / BL-FE-CAT-001 rows; `docs/03-services/01-catalog-svc/01-api.md` — my "Current status" header coexists with Lutfi's `public` / `staff only` annotations on the planned-endpoints table). No manual conflict resolution needed.
- **Security review.** One round, zero findings at confidence ≥ 8. Sub-agent walked SQL injection on all 7 query params + the path param (all bound via pgx, zero string concatenation), draft/archived leakage (top-level predicate unconditional; identical 404 response for all non-active cases), error-envelope leakage (closed set of static Bahasa messages; raw `err.Error()` never reaches response bodies), cursor deserialization (fixed-shape struct, no gadget-chain surface), secrets (none).

## 2026-04-22 — what landed with BL-CAT-002

- **One new endpoint** — `GET /v1/package-departures/{id}`. Returns departure detail with live `remaining_seats` (computed in SQL as `total_seats - reserved_seats`), pricing per room type (ordered cheapest-first), and a stubbed `vendor_readiness` (`not_started` for ticket/hotel/visa). Real readiness wires from visa-svc / logistics-svc in S3-E-02 / S3-E-06.
- **Hidden-status hiding** — `departed` / `completed` / `cancelled` departures return `404 departure_not_found` (identical shape to unknown-id — no existence oracle). Implemented as an in-query `WHERE status IN ('open', 'closed')` predicate; `pgx.ErrNoRows` maps to `ErrNotFound`.
- **`writeCatalogErrorFor` refactor** — `mapCatalogError` and `writeCatalogError` were parametrised so the `ErrNotFound` branch emits the correct resource-scoped code (`package_not_found` vs `departure_not_found`). No change to caller behaviour for the package endpoints.
- **Migration `000010`** — seed one cancelled departure (`dep_01JCDF00000000000000000003`) for e2e 404 coverage. Immutable post-merge; extended as a new migration not an edit to 000009.
- **e2e coverage** — 3 new cases in `02e-catalog-svc-read.spec.ts`: happy-path 200 (seat math, 3 pricing rows, vendor_readiness shape), cancelled-departure 404, unknown-departure 404. Total now 12 passed.

## § Next

- `BL-CAT-004` (F2-AC) — B2C browse integration wrap-up + gateway route smoke. Exec seq 123; same task code `S1-E-02`.
- After S1-E-02 closes: `BL-CAT-003` (atomic `ReserveSeats` / `ReleaseSeats` gRPC) moves to `S1-E-03`, which also covers the booking draft handler and its submit saga.
