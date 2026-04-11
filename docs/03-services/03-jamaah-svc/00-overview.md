# jamaah-svc — Overview

## Purpose

Pilgrim profile, family graph, mahram relations, document storage, and OCR results. The authority on jamaah identity.

## Bounded context

Pilgrim. See `docs/02-domain/00-bounded-contexts.md` § 3.

## PRD source

PRD section E (operational handling — document vault, OCR) and parts of section A (jamaah self-registration).

## Owns (data)

- `jamaah` — biodata, contact, status (calon / active / alumni)
- `family_units` — group of related jamaah
- `mahram_relations` — edges in the family graph
- `documents` — KTP, passport, vaccine, family book scans
- `ocr_results` — structured fields extracted from documents (passport MRZ, etc.)

## Boundaries (does NOT own)

- User accounts (`iam-svc`) — even though jamaah may have a login, the auth credentials live in iam
- Bookings (`booking-svc`)
- Visa applications (`visa-svc`)

## Interactions

- **Inbound:** booking-svc reads jamaah and validates mahram; visa-svc reads passport data; ops-svc reads jamaah for verification queue.
- **Outbound:** GCS (document storage), GCP Vision (OCR), iam-svc (audit log).

## Notable behaviors

- **Mahram validation algorithm.** Walks the family graph (recursive CTE) to determine valid mahram for a female jamaah under 45.
- **OCR pipeline.** Document upload → enqueue OCR job → call GCP Vision → parse MRZ → write `ocr_results`.
- **Document vault.** GCS V4 signed URLs (15-min expiry). No permanent public URLs.
- **K-Family Code.** Human-readable family group identifier referenced from booking.
