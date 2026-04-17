---
id: Q076
title: Daily App form factor — native vs PWA vs responsive web
asked_by: session 2026-04-17 F12 draft
asked_date: 2026-04-17
blocks: F12
status: open
---

# Q076 — Daily App form factor

## Context

PRD Section J names this "Daily App" which linguistically implies a standalone mobile app, but the sitemap (lines 707–713) lists all alumni routes under `/jamaah/*` suggesting an embedded surface inside the jamaah web portal. The phrasing mismatch needs resolution before implementation.

Three viable form factors:

1. **Native iOS + Android** (Flutter / React Native / SwiftUI + Jetpack) — standalone mobile app in app stores.
2. **Svelte PWA** (Progressive Web App) — single codebase, installable via browser prompt, offline-capable, device-sensor access.
3. **Responsive web inside jamaah portal** — no install, just a responsive portal.

Each has cost + capability tradeoffs. Daily Worship features (qibla compass, push notifications for adzan, offline Quran during Saudi trip) benefit from deeper device integration that responsive web can't provide fully.

## The question

1. **Form factor choice** — native / PWA / responsive web?
2. **Single codebase or per-platform** — if native, Flutter + web or separate builds?
3. **App store presence** — if native, is app-store approval overhead in scope?
4. **Install flow** — how does user discover and install (WA link, QR, browser prompt)?
5. **Offline behavior** — which features are offline-first?
6. **Push notifications** — Web Push, Apple Push Notification Service (APNs), Firebase Cloud Messaging (FCM)?
7. **Device sensors** — qibla compass needs gyroscope + magnetometer; camera for check-in QR; GPS for prayer times.
8. **Update cadence** — app-store review vs web deploy?

## Options considered

- **Option A — Svelte PWA installable from jamaah portal.** Single codebase with main ERP frontend; manifest + service worker enable install + offline. Web Push for notifications.
  - Pros: one codebase; no app-store; fast iteration; good-enough sensor access via modern web APIs.
  - Cons: iOS Web Push support is newer (iOS 16.4+); no App Store discovery.
- **Option B — Responsive web only, no PWA / install.** Simplest.
  - Pros: simplest; no service worker.
  - Cons: no offline; no push; no install affordance — misses the "Daily App" expectation.
- **Option C — Native iOS + Android via Flutter or React Native.** Full native capability.
  - Pros: full sensor access; app-store discovery; best notification support.
  - Cons: two additional build pipelines; app-store review (weeks lag); separate auth flows.
- **Option D — Responsive web now, native later (Phase 2).** Start lean; upgrade when usage justifies.
  - Pros: fastest to ship the responsive experience; reserve native for when demand validates.
  - Cons: lose the Daily App / offline promise upfront.

## Recommendation

**Option A — Svelte PWA installable from jamaah portal; Phase 2 option to add native wrapper (Capacitor or similar) for deeper app-store presence if usage warrants.**

Option B misses too much of the Daily App value proposition (offline manasik during Saudi trip, adzan notifications, qibla) — responsive-web-only is what makes F12 less valuable than it should be. Option C's dual-native approach is the right answer for a high-volume consumer app but over-engineered for an alumni utility where many users already have stronger specialized apps (e.g. Muslim Pro) installed — competing head-to-head isn't MVP-viable. Option D pushes the decision forward but the interim experience is weaker than a PWA would be, so we give up polish for no speed gain (PWA is comparable in dev time to responsive web).

Option A threads the needle: single Svelte codebase shares auth + styling + API adapters with the main ERP / jamaah portal; PWA manifest + service worker make it installable + offline-capable; modern Web APIs give adequate sensor access (gyroscope via DeviceOrientation, geolocation, Web Push on iOS 16.4+ and all Android); no app-store overhead; can wrap in Capacitor later for native distribution.

Defaults to propose: Svelte PWA manifest + service worker; mobile-first responsive layout; install prompt on first alumni visit (deferred 3 sessions to avoid annoyance); offline-first for prayer times + qibla + manasik + daily dzikir; Web Push for adzan + kajian reminders (user opt-in); Capacitor wrapper deferred to Phase 2 if app-store presence becomes strategic.

Reversibility: adding Capacitor native wrapper later is additive (same codebase). Abandoning PWA for pure responsive is a downgrade (drops install + push + offline) but schema unaffected.

## Answer

TBD — awaiting stakeholder input. Deciders: CTO / frontend lead (PWA viability vs native), agency owner (app-store presence preference), design lead (UX expectation for "Daily App" phrasing).
