---
id: Q021
title: Luggage Tag QR payload scheme + ALL System scan protocol
asked_by: session 2026-04-15 F7 draft
asked_date: 2026-04-15
blocks: F7
status: open
---

# Q021 — Luggage Tag QR + ALL System scan protocol

## Context

PRD modules #95 (ID Card & Luggage Tag with QR, line 325) and #101 (ALL System scan, line 347) both rely on a QR-encoded identity primitive that gets scanned across the whole operational chain: bus boarding (#104), luggage counting (#102), zamzam distribution (#107), tasreh (#105), airport handling (#101).

The PRD doesn't specify what encodes in the QR or how the scanner talks to the backend. Two security concerns drive the question:

1. **Spoofing** — if a QR just encodes `jamaah_id` plainly, anyone with a printer and the ID can impersonate a jamaah (attend on their behalf, claim their zamzam, etc.).
2. **Offline reliability** — airport WiFi is famously bad. Scans need to work offline and reconcile when online, without losing events or double-counting.

## The question

1. **QR payload.** Plain `jamaah_id` / booking_code, or a signed token? If signed, what signing scheme (HMAC with a server secret, JWT, something else)? Is the token scoped to a specific departure (shorter TTL) or the jamaah's lifetime?

2. **Tag reuse after allocation changes.** If a room or bus assignment changes after tags are printed, does the system invalidate the old QR token? A signed short-lived token forces reissue; a stable signed token keeps the physical tag useful through changes.

3. **ALL System scan protocol — online vs offline.**
   - Online-only: app rejects scans without connectivity. Airport reality suggests this is unacceptable.
   - Offline-first with sync: app queues scans locally, syncs on reconnect with idempotency keys. Battery/storage implications; conflict resolution if two tablets scan the same jamaah.
   - Hybrid: try online; fall back to local queue on failure.

4. **Idempotency key** for scan events. `(device_id, scanned_at, event_kind, jamaah_id)` or something else?

5. **Symbology.** Standard QR code, or a specific barcode type (Code 128, Data Matrix) for speed/reliability?

## Options considered

- **Option A — Signed token, stable-per-jamaah, offline-first scanning, standard QR.** QR payload = HMAC-signed `{ jamaah_id, booking_id, tag_seq, iat }` with a server-side key. Tag stays valid across allocation changes (only `room_number` / `bus_number` on the backend record changes; QR stays). Offline-first: app queues scans with local IndexedDB, syncs with `(device_id, scanned_at, event_kind, jamaah_id)` as the idempotency key. Standard QR, wide reader compatibility.
  - Pros: secure (can't spoof without the secret), survives allocation changes (one physical tag per trip), works in offline airport conditions.
  - Cons: requires secret-rotation story; slightly more backend complexity.
- **Option B — Plain JamaahID QR, online-only scanning.** QR encodes a UUID; backend resolves at scan time.
  - Pros: simplest.
  - Cons: spoofable; fails without connectivity.
- **Option C — JWT token per departure, offline-first.** Short-lived JWT issued when tags are printed; expires at H+72 after departure return.
  - Pros: clean expiry semantics.
  - Cons: reissue on every allocation change; tag becomes useless after return (fine for most cases, but loss-claim windows extend weeks).

## Recommendation

**Option A — signed HMAC token, stable per jamaah + booking, offline-first with standard QR.**

Payload:
```json
{
  "j": "<jamaah_id>",
  "b": "<booking_id>",
  "s": "<tag_seq e.g. 1_of_3 or 2_of_3>",
  "iat": 1713187200
}
```
Plus an HMAC signature as the final field. The whole thing base64-encodes into the QR.

Scan flow: reader decodes, verifies signature against the server's current key (with a rotation-overlap window for smooth key rolls), resolves `jamaah_id` + `booking_id`, writes the scan event.

Offline-first: field apps maintain a local IndexedDB queue, sync when online. Idempotency key `(device_id, scanned_at, event_kind, jamaah_id)` — a UNIQUE constraint on the backend dedupes duplicate syncs.

Signing-key rotation: same pattern as F1 PASETO keys. Two keys active at a time; new keys sign, both keys verify; old key retires after the maximum tag lifetime (~90 days post-departure).

Symbology: standard QR at Medium error correction (good balance of density + damage tolerance for tags that get wet or crumpled on a 2-week trip).

Reasoning: Option B is too insecure for a system that touches zamzam quotas and Raudhah access. Option C's reissue-on-change pain outweighs the expiry benefits. Option A gives security, stability, and offline operability at reasonable implementation cost.

Reversibility: token schema is versionable via a leading version byte; signing key rotation is routine; offline queue is per-device local state.

## Answer

TBD — awaiting stakeholder input. Deciders: agency ops lead + any security-minded reviewer.
