// tests/k6/smoke.js
// ── Smoke Test ── minimal sanity check before heavier tests.
//
// Verifies that every demo-svc endpoint is reachable and responds correctly.
// Run with 1 VU, 10 iterations — fast feedback loop.
//
// Usage:
//   k6 run tests/k6/smoke.js
//   k6 run --env BASE_URL=http://localhost:4000 tests/k6/smoke.js

import http from "k6/http";
import { check, group } from "k6";

export const options = {
  vus: 1,
  iterations: 10,
  thresholds: {
    http_req_failed: ["rate<0.01"], // < 1% errors
    http_req_duration: ["p(95)<500"], // p95 < 500 ms
  },
};

const BASE = __ENV.BASE_URL || "http://localhost:4000";
const DIAG_KEY = __ENV.DIAGNOSTIC_KEY || "8ES3CR3T";

export default function () {
  // ── 1. Liveness ─────────────────────────────────────────
  group("liveness", () => {
    const res = http.get(`${BASE}/system/live`);
    check(res, {
      "live: status 200": (r) => r.status === 200,
      "live: responds in < 200 ms": (r) => r.timings.duration < 200,
    });
  });

  // ── 2. Readiness ────────────────────────────────────────
  group("readiness", () => {
    const res = http.get(`${BASE}/system/ready`);
    check(res, {
      "ready: status 200": (r) => r.status === 200,
      "ready: responds in < 300 ms": (r) => r.timings.duration < 300,
    });
  });

  // ── 3. DB-TX Diagnostic ─────────────────────────────────
  group("db-tx diagnostic", () => {
    const res = http.get(`${BASE}/system/diagnostics/db-tx`, {
      headers: { "X-Diagnostic-Key": DIAG_KEY },
    });
    check(res, {
      "db-tx: status 200": (r) => r.status === 200,
      "db-tx: has data": (r) => JSON.parse(r.body).data !== undefined,
    });
  });

  // ── 4. Issue token ──────────────────────────────────────
  group("issue token", () => {
    const res = http.post(`${BASE}/auth/token`, null, {
      headers: { "Content-Type": "application/json" },
    });
    check(res, {
      "token: status 200": (r) => r.status === 200,
      "token: has token": (r) => JSON.parse(r.body).data.token !== undefined,
      "token: has expires_at": (r) =>
        JSON.parse(r.body).data.expires_at !== undefined,
    });
  });

  // ── 5. Auth me (full round-trip: issue → validate) ─────
  group("auth me", () => {
    // Issue a token first
    const issueRes = http.post(`${BASE}/auth/token`, null, {
      headers: { "Content-Type": "application/json" },
    });
    const token = JSON.parse(issueRes.body).data.token;

    // Validate via /auth/me
    const meRes = http.get(`${BASE}/auth/me`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    check(meRes, {
      "me: status 200": (r) => r.status === 200,
      "me: has user_id": (r) => JSON.parse(r.body).data.user_id !== undefined,
      "me: has roles": (r) => JSON.parse(r.body).data.roles !== undefined,
    });
  });
}
