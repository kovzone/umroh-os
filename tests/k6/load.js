// tests/k6/load.js
// ── Load Test ── simulates realistic sustained traffic against demo-svc.
//
// Ramps up to a target number of virtual users, holds steady, then ramps
// down.  Exercises all demo-svc endpoints with a mix of read-heavy (probes,
// auth/me) and write (issue token, db-tx diagnostic) operations.
//
// Usage:
//   k6 run tests/k6/load.js
//   k6 run --env BASE_URL=http://localhost:4000 tests/k6/load.js

import http from "k6/http";
import { check, sleep, group } from "k6";
import { Counter, Trend } from "k6/metrics";

// ── Custom metrics ────────────────────────────────────────
const tokensIssued = new Counter("tokens_issued_total");
const tokensFailed = new Counter("tokens_failed_total");
const tokenLatency = new Trend("token_issue_duration_ms");
const meLatency = new Trend("auth_me_duration_ms");
const diagLatency = new Trend("diag_dbtx_duration_ms");

// ── SLO Thresholds ────────────────────────────────────────
export const options = {
  scenarios: {
    average_load: {
      executor: "ramping-vus",
      stages: [
        { duration: "30s", target: 50 }, // warm-up
        { duration: "2m", target: 200 }, // sustained load
        { duration: "30s", target: 0 }, // cool-down
      ],
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"], // < 1% errors
    http_req_duration: ["p(95)<300", "p(99)<800"],
    token_issue_duration_ms: ["p(95)<400", "p(99)<1000"],
    auth_me_duration_ms: ["p(95)<200"],
    diag_dbtx_duration_ms: ["p(95)<500"],
    tokens_issued_total: ["count>500"], // must complete ≥ 500 issues
  },
};

const BASE = __ENV.BASE_URL || "http://localhost:4000";
const DIAG_KEY = __ENV.DIAGNOSTIC_KEY || "8ES3CR3T";

// ── VU logic ──────────────────────────────────────────────
export default function () {
  // 1. Health probes (lightweight, read-only)
  group("health probes", () => {
    const live = http.get(`${BASE}/system/live`);
    check(live, { "live 200": (r) => r.status === 200 });

    const ready = http.get(`${BASE}/system/ready`);
    check(ready, { "ready 200": (r) => r.status === 200 });
  });

  // 2. Issue a token
  let token;
  group("issue token", () => {
    const start = Date.now();
    const res = http.post(
      `${BASE}/auth/token`,
      JSON.stringify({
        user_id: `user-${__VU}-${__ITER}`,
        branch_id: "branch-load-test",
        roles: ["tester"],
      }),
      { headers: { "Content-Type": "application/json" } },
    );
    const dur = Date.now() - start;
    tokenLatency.add(dur);

    const ok = check(res, {
      "token 200": (r) => r.status === 200,
      "has token": (r) => JSON.parse(r.body).data.token !== undefined,
    });
    ok ? tokensIssued.add(1) : tokensFailed.add(1);

    if (ok) {
      token = JSON.parse(res.body).data.token;
    }
  });

  // 3. Validate token via /auth/me
  if (token) {
    group("auth me", () => {
      const start = Date.now();
      const res = http.get(`${BASE}/auth/me`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      meLatency.add(Date.now() - start);
      check(res, {
        "me 200": (r) => r.status === 200,
        "me has user_id": (r) =>
          JSON.parse(r.body).data.user_id !== undefined,
      });
    });
  }

  // 4. DB-TX diagnostic (write path — insert + read round-trip)
  group("db-tx diagnostic", () => {
    const start = Date.now();
    const res = http.get(`${BASE}/system/diagnostics/db-tx`, {
      headers: { "X-Diagnostic-Key": DIAG_KEY },
    });
    diagLatency.add(Date.now() - start);
    check(res, {
      "diag 200": (r) => r.status === 200,
    });
  });

  sleep(0.5); // think-time between iterations
}
