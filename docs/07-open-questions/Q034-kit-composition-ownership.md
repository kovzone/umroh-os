---
id: Q034
title: Kit composition ownership — catalog-svc vs logistics-svc
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8, F2
status: answered
---

# Q034 — Kit composition ownership

## Context

A "kit" (Indonesian *paket*) is the bundle of physical goods a jamaah receives: koper, ihram, kain seragam, buku doa, tas travel, etc. Kit composition varies per package tier — a Paket Silver might include 5 items, Paket Gold 7, Paket Platinum 9 with premium koper.

There are two different things sharing the word "kit":

- **Kit definition** — the named recipe ("Paket Silver 2026 Kit = 1× koper-24in + 1× ihram-M + 1× buku-doa + 1× travel-bag + 1× seragam-set-L"). Attribute of the package.
- **Kit instance** — an actual assembled bundle in the warehouse, reserved for a specific jamaah.

PRD module #119 *Perakitan Paket* (line 401) lives in Section F (Logistics). But package cost / tier composition is implicitly catalog-svc territory (PRD Section D — Master Product, pricing, HPP). Which service owns what?

## The question

1. **Does `kit_definition` (the recipe) live in `catalog-svc` or `logistics-svc`?**
2. If catalog-svc: how does logistics-svc read it — gRPC call per kit assembly, or event-sync'd cache?
3. If logistics-svc: how does catalog-svc express "Paket Silver includes kit definition X" — by FK to logistics, by embedded list, or by loose coupling (catalog just names the kit, logistics defines components)?
4. **Kit versioning** — when Paket Silver 2026 includes koper-24in and Paket Silver 2027 switches to koper-26in, is this a new kit_definition or a versioned edit of the same one?
5. **Component substitution at assembly time** — if koper-24in is out of stock but koper-25in is available, can warehouse substitute? If yes, is the substitution logged against the kit_instance or against a change of kit_definition?
6. **Cost implications** — kit composition directly affects HPP (which belongs to catalog). If logistics owns the recipe, how does catalog know what a kit costs?

## Options considered

- **Option A — Kit definitions owned by catalog-svc; kit instances owned by logistics-svc.** Catalog defines recipes (SKU list + qty per package tier). Logistics reads the definition at assembly time via gRPC and creates instances. HPP calculation lives naturally with catalog.
  - Pros: clean bounded-context split; HPP stays with packaging; logistics owns execution.
  - Cons: cross-service call on every assembly; catalog-svc must know about SKUs (which are otherwise logistics's domain).
- **Option B — Kit definitions owned by logistics-svc; catalog-svc references by ID.** Logistics owns everything kit-related. Catalog attaches a `kit_definition_id` to each package tier.
  - Pros: all kit logic in one service; logistics owns SKUs and kit_definitions coherently.
  - Cons: HPP calculation needs logistics data (crosses into catalog's pricing concerns).
- **Option C — Kit definitions in catalog-svc, but stored as denormalized SKU list (no FK).** Catalog keeps `package_tiers.kit_components` as `[{sku, qty}]` jsonb. Logistics validates SKU existence at assembly time.
  - Pros: no cross-service call at query time; catalog is self-contained.
  - Cons: SKU list can drift from logistics's truth; validation gaps.

## Recommendation

**Option A — kit definitions in catalog-svc, kit instances in logistics-svc.**

Kits are product attributes first, physical artefacts second. When marketing prints a brochure saying "Paket Silver includes...", that list needs to be in the same service as the package itself (catalog-svc). Option B sounds cleaner for logistics but puts product composition in a warehouse service — the first time catalog-svc needs kit data for the brochure it has to cross-call logistics. Option C's jsonb avoids the FK but loses validation (nothing stops catalog from listing a non-existent SKU).

Defaults to propose: `catalog.kit_definitions { id, package_tier_id, components: [{sku, qty}] }` with SKU validated against `logistics-svc.CheckSku(sku)` at definition creation time. Logistics reads via `logistics-svc.GetKitDefinition(id)` (which reaches into catalog) at assembly time — or more cleanly, catalog-svc pushes kit-definition-changed events that logistics caches. Kit versioning: new year = new kit_definition row linked to the new package tier. Component substitution at assembly: allowed with supervisor approval, logged against the kit_instance (not the definition) — kit_instance.substitutions jsonb captures what was swapped. HPP computed from the components × current average cost per SKU (read from logistics), cached in the package row for quick brochure rendering.

Reversibility: moving the table between services later is expensive but feasible via an extraction migration. Schema shape (components as jsonb vs normalized kit_components table) is lighter to change.

## Answer

**Decided:** **Option A** — **`kit_definition` in catalog-svc**, **`kit_instance` in logistics-svc**; **SKU existence validated** at definition write; **event/cache read** at assembly acceptable; **substitutions on `kit_instance`** with supervisor flag; **new kit_definition row per catalog year/tier change**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
