# Bounded Contexts

UmrohOS is decomposed into bounded contexts (DDD term: a region of the domain with its own ubiquitous language and data ownership). Each context maps to one service. Cross-context interaction goes via gRPC reads or in-process sagas in the orchestrating service (per ADR 0006). The F6 visa pipeline is the one exception — it will use Temporal when implemented.

## The contexts

### 1. Identity & Access (`iam-svc`)

**Owns:** users, roles, permissions, branches, sessions, audit log.

**Language:** user, role, permission, branch, scope, token, audit.

**Boundaries:** Every other service calls iam-svc to validate tokens and check permissions. iam-svc never knows about packages, bookings, or jamaah.

### 2. Product Catalog (`catalog-svc`)

**Owns:** packages (sellable), hotels, airlines, muthawwif, itineraries, seat inventory, package pricing.

**Language:** package, hotel, airline, muthawwif, itinerary, seat, room type (double/triple/quad), departure date.

**Boundaries:** Catalog is the master record for sellable products. Bookings hold a reference to a package, but never copy package data — they always read fresh.

### 3. Pilgrim (`jamaah-svc`)

**Owns:** jamaah biodata, family graph, mahram relations, documents (KTP, passport, vaccine), OCR results.

**Language:** jamaah, mahram, family, document, MRZ, K-Family Code.

**Boundaries:** The pilgrim's identity and documents live here. Bookings reference jamaah by ID; visa applications read passport data via gRPC. Jamaah-svc owns the family graph and is the authority on mahram validation.

### 4. Booking (`booking-svc`)

**Owns:** bookings, booking items (jamaah-package links), room allocations, bus allocations, manifests (lightweight; full manifest generation is in ops-svc).

**Language:** booking, reservation, room, bus, allocation, status (draft / DP / lunas / cancelled).

**Boundaries:** Booking is the link between a jamaah and a package. It does not store payment data (that's payment-svc) or visa data (visa-svc). It does store room/bus assignments because those are scheduled at booking time.

### 5. Payment (`payment-svc`)

**Owns:** invoices, virtual accounts, payment events (gateway webhooks), refunds.

**Language:** invoice, virtual account, gateway, settlement, refund, DP, lunas, installment.

**Boundaries:** Payment is the only service that talks to Midtrans/Xendit. It signals booking when payment status changes via direct gRPC call (`booking-svc.MarkBookingPaid`) per ADR 0006. It does not write journal entries — finance-svc consumes payment events for that.

### 6. Visa (`visa-svc`)

**Owns:** visa applications, status history, e-visas, tasreh records, Raudhah Shield monitoring state.

**Language:** visa, e-visa, tasreh, MOFA, Sajil, Nusuk, Raudhah, status (waiting_docs / ready / submitted / issued / rejected).

**Boundaries:** Visa-svc is the only service that talks to MOFA/Sajil/Nusuk. It reads passport data from jamaah-svc; it signals booking when a visa is attached.

### 7. Operations (`ops-svc`)

**Owns:** document verification queue, luggage tags, airport handling events, full manifest generation, smart room/bus allocation algorithm.

**Language:** verification, manifest, luggage tag, ALL system, airport handling, room grouping, bus seating.

**Boundaries:** Ops is the workflow surface for back-office operational staff. It reads from jamaah, booking, catalog, and writes its own working state. The smart-grouping algorithm is business logic that lives here.

### 8. Logistics (`logistics-svc`)

**Owns:** stock items, warehouses, purchase orders, GRN (goods received notes), kit assemblies, shipments.

**Language:** stock, warehouse, PO, GRN, kit, SKU, shipment, courier.

**Boundaries:** Logistics is triggered by booking payment status (via direct gRPC call from payment-svc per ADR 0006). It does not know about jamaah identity beyond shipping address.

### 9. Finance & Accounting (`finance-svc`)

**Owns:** chart of accounts, journal entries, AR, AP, tax records, FX rates, job-order costs.

**Language:** journal, debit, credit, COA, AR, AP, PPh, PPN, FX, job order, balance sheet, P&L.

**Boundaries:** Finance is downstream — it consumes events from payment, logistics, and crm to write journal entries. It does not write back to those services. PSAK compliance lives here.

### 10. Marketing & CRM (`crm-svc`)

**Owns:** leads, campaigns, agents, agent network hierarchy, commission ledger, broadcasts, alumni community threads, ZISWAF transactions.

**Language:** lead, campaign, agent, super-agent, branch, commission, override, broadcast, alumni, referral.

**Boundaries:** CRM owns the agent network. When a booking is attributed to an agent, crm-svc calculates commission via gRPC invoked from payment-svc when the booking hits paid_in_full (in-process; per ADR 0006). It does not own the booking itself.

### 11. Workflows (`broker-svc`) — **DEFERRED for MVP (see ADR 0006)**

**Owns:** Temporal workflow definitions and activities. Stateless from a business-data perspective — Temporal owns workflow state.

**Language:** workflow, activity, signal, query, saga, compensation.

**Boundaries:** broker-svc is the only service that imports the Temporal SDK as a workflow author. It calls into other services via gRPC adapters.

## Context map (relationships)

| Upstream | Downstream | Type |
|---|---|---|
| iam-svc | every service | upstream/conformist (token validation) |
| catalog-svc | booking-svc, ops-svc | upstream/customer-supplier |
| jamaah-svc | booking-svc, visa-svc, ops-svc | upstream/customer-supplier |
| booking-svc | payment-svc, logistics-svc, ops-svc, finance-svc | upstream/customer-supplier |
| payment-svc | finance-svc | upstream (event source; direct gRPC per ADR 0006) |
| visa-svc | booking-svc | partner (signals on visa attached) |
| crm-svc | booking-svc | partner (lead → booking attribution) |
| broker-svc *(deferred, F6)* | visa-svc, jamaah-svc, booking-svc | orchestrator for the visa pipeline when F6 is implemented (per ADR 0006) |
