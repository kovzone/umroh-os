# crm-svc — Overview

## Purpose

Marketing, CRM, agent network, commissions, broadcasts, alumni community. Owns the entire B2B agent ecosystem and the customer relationship layer.

## Bounded context

Marketing & CRM. See `docs/02-domain/00-bounded-contexts.md` § 10.

## PRD source

PRD sections B (B2B Front-End), C (Marketing & Sales / CRM), J (Daily App & Alumni Hub).

## Owns (data)

- `leads`
- `campaigns`
- `agents` (with hierarchy)
- `commission_ledger`
- `broadcasts`
- `alumni_threads`
- `referral_codes`
- `ziswaf_transactions`

## Boundaries (does NOT own)

- User accounts (`iam-svc`) — agents have a `user_id` reference
- Bookings (`booking-svc`) — crm reads bookings to compute commissions
- Payment / financial vendor records (`finance-svc`)

## Interactions

- **Inbound:** booking-svc events for commission calc; ops/visa events for customer notifications.
- **Outbound:** Meta WhatsApp API, Meta/Google Ads API, email service.

## Notable behaviors

- **Multi-level commission overriding** — super-agents earn the difference between their level and a sub-agent's level.
- **Replicated agent websites** — each agent has a personal landing page; flyers auto-watermarked with their WhatsApp.
- **Lead nurturing** via automated WhatsApp campaigns.
- **Broadcast handling** — message bursts to 5,000+ agents must queue and rate-limit.
- **Alumni community** — threads, Q&A, referral codes, savings goals for return pilgrimage.
