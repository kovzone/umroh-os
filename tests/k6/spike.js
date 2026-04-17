// tests/k6/spike.js
// ── Spike Test ── simulates a sudden traffic burst (e.g. payday, promo).
//
// Holds a steady baseline, then spikes VUs sharply, returns to baseline,
// holds again, and ramps down.  Tests whether the system recovers gracefully
// after the burst.
//
// Usage:
//   k6 run tests/k6/spike.js
//   k6 run --env BASE_URL=http://localhost:4000 tests/k6/spike.js

import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    spike: {
      executor: "ramping-vus",
      stages: [
        { duration: "30s", target: 50 }, // baseline
        { duration: "10s", target: 1000 }, // sudden spike
        { duration: "10s", target: 50 }, // back to baseline
        { duration: "30s", target: 50 }, // hold baseline
        { duration: "10s", target: 0 }, // ramp down
      ],
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.02"], // allow up to 2% during spike
    http_req_duration: ["p(95)<1000"], // p95 < 1s
  },
};

const BASE = __ENV.BASE_URL || "http://localhost:4000";

export default function () {
  // Read-heavy mix — simulates many users hitting probes + validating tokens
  // during a traffic burst.

  // Lightweight probe (majority of spike traffic)
  const live = http.get(`${BASE}/system/live`);
  check(live, { "live: status 200": (r) => r.status === 200 });

  const ready = http.get(`${BASE}/system/ready`);
  check(ready, { "ready: status 200": (r) => r.status === 200 });

  // Token round-trip (issue → validate)
  const issueRes = http.post(`${BASE}/auth/token`, null, {
    headers: { "Content-Type": "application/json" },
  });
  check(issueRes, { "token: status 200": (r) => r.status === 200 });

  if (issueRes.status === 200) {
    const token = JSON.parse(issueRes.body).data.token;
    const meRes = http.get(`${BASE}/auth/me`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    check(meRes, { "me: status 200": (r) => r.status === 200 });
  }

  sleep(0.2); // short think-time
}
