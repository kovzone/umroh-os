---
id: Q041
title: Self-pickup QR security model
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8
status: open
---

# Q041 — Self-pickup QR security model

## Context

PRD module #124 *Pengambilan Mandiri* (line 415) offers self-pickup as an alternative to courier delivery: the system generates a QR-code digital receipt that the jamaah presents at the office counter. F8 W12 needs the security model specified — a naive QR is a bearer token anyone could present.

Concerns: prevent impersonation (someone else collects the kit on behalf without authorization), prevent replay (same QR used twice), handle lost phones (jamaah lost the WhatsApp-sent QR), support delegation (jamaah sends a family member to pick up).

This is analogous to the F7 luggage-tag QR security question (Q021), but the threat model is different: luggage tags are scanned by agency staff at airports; self-pickup QRs are scanned at the agency's office counter by the agency's own staff. Lower-stakes but still a physical handover of valuable goods.

## The question

1. **Signing scheme** — HMAC-signed payload, JWT, opaque random token validated against DB, or something else?
2. **Expiry** — how long is the QR valid? 7 days? 30 days? Until the departure date?
3. **Single-use or multi-use** — can the same QR be scanned once, or does it just attest identity (multi-use)?
4. **Revocation** — can ops invalidate a lost/compromised QR? Generate a replacement?
5. **Verification flow at counter** — does the office staff just scan, or also verify jamaah identity (KTP photo side-by-side)?
6. **Delegation** — if jamaah sends a family member, how is that authorized? Is the QR transferable? Does the delegate need their own KTP match?
7. **Offline fallback** — if the office internet is down, can staff verify the QR locally?

## Options considered

- **Option A — Opaque random token + DB lookup + identity check.** 32-byte random token; DB row per token (dispatch_task_id, expiry, used_at). Office staff scans, system fetches the row, displays jamaah KTP photo, staff confirms identity match, marks token used.
  - Pros: simple; server-side state prevents replay; easy revocation (set `revoked_at`).
  - Cons: requires online DB lookup; no offline verification; state grows linearly with dispatches.
- **Option B — HMAC-signed token (self-contained).** Token encodes `{dispatch_task_id, jamaah_id, expiry}` with HMAC signature. Staff scans, system verifies signature + expiry; separate tracking for used-once.
  - Pros: smaller server-side state (just "used" flag); signatures verifiable offline if key is cached.
  - Cons: revocation requires a blocklist (can't "unsign" a signed token).
- **Option C — JWT-based (industry-standard).** Same as B but uses JWT conventions. Standard library support.
  - Pros: mature tooling, easy audit.
  - Cons: slightly heavier payload; JWT complexity.

## Recommendation

**Option A — opaque token + DB lookup + jamaah identity check, with ops revocation support.**

This is not a scale-sensitive use case (a few hundred pickups/month for most agencies), so the simplicity of DB-backed opaque tokens wins. HMAC/JWT's offline-verifiability isn't worth the revocation complexity — if the office's internet is down, delaying pickups by 15 minutes is a fine failure mode for a non-departure-critical flow. Server-side state makes revocation (lost phone, compromised token) a one-row update. The identity check — KTP photo alongside scan — is the real security control; the QR is just a pointer to "which pickup record are we confirming." That's the pattern that matches how humans verify pickups in the physical world.

Defaults to propose: 32-byte random token, base64url-encoded in the QR; stored in `self_pickup_tokens { token_hash, dispatch_task_id, jamaah_id, valid_until, used_at, revoked_at }`. Token hash (not plaintext) in DB. Validity: 30 days from generation (covers booking-to-departure window for most cases). Single-use (`used_at` set on first successful scan). Revocation: `revoked_at` set by ops; fails verification. Counter staff flow: scan → system displays jamaah KTP photo + booking code + kit summary; staff visually confirms identity; clicks "Confirm pickup." Delegation: jamaah can request a delegation in the jamaah portal (Phase 2 — MVP has jamaah-only pickup); a delegated QR encodes a different `delegate_name` and staff verifies delegate's ID. Offline: flow blocks, "please try again when online" — acceptable failure mode.

Reversibility: switching to HMAC-based tokens later is a migration — only matters if we need offline verification, which isn't on the roadmap.

## Answer

TBD — awaiting stakeholder input. Deciders: ops lead (pickup counter operational experience), security advisor (if retained), CS lead (delegation policy).
