import { defineConfig, devices } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";

dotenv.config({ path: path.resolve(__dirname, ".env") });

// Three projects:
//   - api:     APIRequestContext-only specs (01-*, 02-*, 02a-*, ...). Hit backend REST
//              endpoints directly, no browser required.
//   - uat-api: UAT API-only tests embedded in UAT spec files (04-*, 05-*, 06-*).
//              These run API tests without needing a browser.
//   - browser: real-browser specs (03-*). Load the SvelteKit app through
//              gateway and assert DOM state. Requires `npx playwright install chromium`
//              which `make e2e-install` runs for us.
//
// 2026-04-21 note: the api pattern accepts an optional lowercase letter after
// the digit pair (`02a-...`, `02b-...`) so BL-IAM-001's iam-svc-sessions spec
// can slot in next to the existing gateway tests without re-numbering the
// browser spec.
// 2026-04-24 note: browser pattern now also accepts an optional letter suffix
// (`03b-...`) so the B2C browser E2E spec can coexist with 03-core-web-status.
export default defineConfig({
  testDir: "./tests",
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: "list",
  projects: [
    {
      name: "api",
      testMatch: /0[0-2][a-z]?-.*\.spec\.ts$/,
    },
    {
      name: "uat-api",
      testMatch: /0[4-9]-.*\.spec\.ts$/,
      use: {
        baseURL: process.env.GATEWAY_SVC_URL || "http://216.176.238.161",
      },
    },
    {
      name: "browser",
      testMatch: /0[3-9][a-z]?-.*\.spec\.ts$/,
      use: {
        ...devices["Desktop Chrome"],
        // Prefer "localhost" for the browser baseURL: on Docker Desktop for Windows,
        // 127.0.0.1:<hostPort> can be routed differently than localhost for some ports
        // (observed: 127.0.0.1:3001 hit Grafana while localhost:3001 hit core-web).
        baseURL: process.env.CORE_WEB_URL || "http://localhost:3001",
      },
    },
  ],
});
