import { defineConfig, devices } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";

dotenv.config({ path: path.resolve(__dirname, ".env") });

// Two projects:
//   - api:     APIRequestContext-only specs (01-*, 02-*). Hit backend REST
//              endpoints directly, no browser required.
//   - browser: real-browser specs (03-*). Load the SvelteKit app through
//              gateway and assert DOM state. Requires `npx playwright install chromium`
//              which `make e2e-install` runs for us.
export default defineConfig({
  testDir: "./tests",
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: "list",
  projects: [
    {
      name: "api",
      testMatch: /0[0-2]-.*\.spec\.ts$/,
    },
    {
      name: "browser",
      testMatch: /0[3-9]-.*\.spec\.ts$/,
      use: {
        ...devices["Desktop Chrome"],
        baseURL: process.env.CORE_WEB_URL || "http://localhost:3001",
      },
    },
  ],
});
