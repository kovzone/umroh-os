// tests/k6/stress.js
// ── Stress / Break-point Test ── find the system's limit.
//
// Ramps VUs aggressively to discover at what concurrency level the service
// starts to degrade (latency > threshold or error rate > threshold).
// Thresholds use abortOnFail — k6 will stop early when the system breaks.
//
// Usage:
//   k6 run tests/k6/stress.js
//   k6 run --env BASE_URL=http://localhost:4000 tests/k6/stress.js

import http from "k6/http";
import { check, sleep, group } from "k6";

export const options = {
  scenarios: {
    breakpoint: {
      executor: "ramping-vus",
      stages: [
        { duration: "1m", target: 100 },
        { duration: "1m", target: 300 },
        { duration: "1m", target: 600 },
        { duration: "1m", target: 1000 },
        { duration: "30s", target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_failed: [
      { threshold: "rate<0.05", abortOnFail: true }, // abort if > 5% errors
    ],
    http_req_duration: [
      { threshold: "p(99)<2000", abortOnFail: true }, // abort if p99 > 2s
    ],
  },
};

const BASE = __ENV.BASE_URL || "http://localhost:4000";
const DIAG_KEY = __ENV.DIAGNOSTIC_KEY || "8ES3CR3T";

export default function () {
  // Mix of endpoints — heaviest first (db-tx), then token round-trip.

  group("stress: db-tx diagnostic", () => {
    const res = http.get(`${BASE}/system/diagnostics/db-tx`, {
      headers: { "X-Diagnostic-Key": DIAG_KEY },
    });
    check(res, {
      "diag: status < 500": (r) => r.status < 500,
    });
  });

  group("stress: token round-trip", () => {
    // Issue
    const issueRes = http.post(`${BASE}/auth/token`, null, {
      headers: { "Content-Type": "application/json" },
    });
    check(issueRes, {
      "token: status < 500": (r) => r.status < 500,
    });

    // Validate (only if issue succeeded)
    if (issueRes.status === 200) {
      const token = JSON.parse(issueRes.body).data.token;
      const meRes = http.get(`${BASE}/auth/me`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      check(meRes, {
        "me: status < 500": (r) => r.status < 500,
      });
    }
  });

  group("stress: probes", () => {
    const live = http.get(`${BASE}/system/live`);
    check(live, { "live: status < 500": (r) => r.status < 500 });

    const ready = http.get(`${BASE}/system/ready`);
    check(ready, { "ready: status < 500": (r) => r.status < 500 });
  });

  sleep(0.1); // minimal think-time — stress test
}
