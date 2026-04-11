# UmrohOS — Vision

## What it is

UmrohOS is an end-to-end ERP for Umrah and Hajj travel agencies. One platform that consolidates everything an agency needs to operate at scale: B2C package booking, B2B agent networks, marketing/CRM, master product/inventory, document and visa processing, warehouse logistics, PSAK-compliant finance and accounting, jamaah field operations in Saudi Arabia, alumni hub, and executive dashboards.

## Who it's for

Mid-to-large pilgrimage travel operators in Indonesia (initially), serving 1,000–10,000+ jamaah per year. The platform is single-tenant per company — each customer runs their own instance — but multi-branch within a company. It must support a hierarchy of users: prospective pilgrims, active jamaah, alumni, customer service agents, travel agents/resellers, branch managers, ops/warehouse/finance staff, tour leaders in the field, branch directors, and the company owner.

## The problem

Pilgrimage agencies today juggle a fragmented stack — accounting in one tool, design in another, CRM in a third, separate spreadsheets for warehouse, separate manual workflows for visa applications, separate WhatsApp threads for field updates. Data is re-keyed across every handoff. Errors surface in critical moments: a visa is misfiled, a payment is misallocated, a manifest goes out wrong. Scale is impossible without proportional headcount.

## The bet

A single source of truth — one platform where every entity (jamaah, package, hotel, flight, payment, visa, document, journal entry) lives once and is referenced everywhere — eliminates re-keying and broken handoffs. Combined with intelligent automation (OCR, smart room/bus allocation, virtual account reconciliation, Temporal-orchestrated visa pipelines), it lets an agency 10x throughput without 10x headcount.

## Non-goals

- **Multi-tenant SaaS.** UmrohOS is single-tenant per customer. No tenant isolation in the data layer beyond branch scoping.
- **A general travel ERP.** Built specifically for Umrah and Hajj. Domain logic (mahram validation, tasreh management, Raudhah Shield, PSAK accounting) is first-class, not pluggable.
- **A consumer mobile app marketplace.** The jamaah and field apps are companions to the agency platform, not standalone products.

## Source of truth

The full PRD is at `docs/UmrohOS - Product Requirements Document.docx.md` (~1,620 lines, mostly Indonesian terminology). When this vision file disagrees with the PRD, the PRD wins — but flag the conflict in your session note so it can be reconciled.
