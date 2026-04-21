import { test, expect } from "@playwright/test";
import { backendServices } from "../lib/services";

// core-web browser e2e. Uses the `browser` project in playwright.config.ts
// (baseURL defaults to http://localhost:3001).
//
// Three groups of tests:
//   1. Landing page (/) — hero, capability grid, disabled Sign-in, footer link.
//   2. Status page (/system/status) — 10 service cards flip to `ok` after poll.
//   3. Landing → status navigation via the footer link.
//
// Failure-path tests (stop a backend, expect card to flip to fail) are NOT
// asserted here — they need `docker compose stop`, which isn't portable across
// CI runners. Manual walk is listed in testing-guide.md Section 9.

test.describe.serial("core-web — landing page", () => {
  test("renders hero, capability grid, and footer status link", async ({ page }) => {
    await page.goto("/");

    // Hero + tagline
    const hero = page.getByTestId("landing-hero");
    await expect(hero).toBeVisible();
    await expect(hero).toContainText("UmrohOS");

    // Capability grid has 4 cards (must match the landing page's capabilities list)
    const grid = page.getByTestId("capability-grid");
    await expect(grid.getByTestId("capability-card")).toHaveCount(4);

    // Sign-in button is visible but disabled (auth lands in F1.5)
    const signin = page.getByTestId("signin-button");
    await expect(signin).toBeVisible();
    await expect(signin).toBeDisabled();

    // Footer has the status-page link
    const statusLink = page.getByTestId("footer-status-link");
    await expect(statusLink).toBeVisible();
    await expect(statusLink).toHaveAttribute("href", "/system/status");
  });

  test("footer link navigates to /system/status", async ({ page }) => {
    await page.goto("/");
    await page.getByTestId("footer-status-link").click();
    await expect(page).toHaveURL(/\/system\/status$/);
  });

  test("browse packages CTA navigates to /packages", async ({ page }) => {
    await page.goto("/");
    await page.getByTestId("browse-packages-cta").click();
    await expect(page).toHaveURL(/\/packages$/);
    await expect(page.getByTestId("s1-package-catalog")).toBeVisible();
  });
});

test.describe.serial("core-web — S1 catalog shells (S1-L-01)", () => {
  test("packages list links to package detail", async ({ page }) => {
    await page.goto("/packages");
    await page.getByTestId("package-link-demo-pkg-umrah-12d").click();
    await expect(page).toHaveURL(/\/packages\/demo-pkg-umrah-12d$/);
    await expect(page.getByTestId("s1-package-detail")).toBeVisible();
  });

  test("package detail links to booking shell", async ({ page }) => {
    await page.goto("/packages/demo-pkg-umrah-12d");
    await page.getByTestId("s1-start-booking").click();
    await expect(page).toHaveURL(/\/booking\/demo-pkg-umrah-12d$/);
    await expect(page.getByTestId("s1-booking-draft-shell")).toBeVisible();
  });
});

test.describe.serial("core-web — service status page", () => {
  test("renders 10 service cards, all OK after poll", async ({ page }) => {
    await page.goto("/system/status");

    // Grid mounts immediately with 10 pending cards.
    const grid = page.getByTestId("status-grid");
    await expect(grid.getByTestId("service-card")).toHaveCount(backendServices.length);

    // Each card should reach data-status="ok" within the polling window
    // (first poll fires immediately; allow slack for container startup).
    for (const svc of backendServices) {
      const card = page.locator(`[data-service="${svc.name}"]`);
      await expect(card).toHaveAttribute("data-status", "ok", { timeout: 15_000 });
    }
  });

  test("last-poll indicator shows a timestamp after polling", async ({ page }) => {
    await page.goto("/system/status");
    const indicator = page.getByTestId("last-poll");
    // Starts as "Last poll: pending…" — must transition to include a colon
    // (HH:MM:SS time format) once the first poll completes.
    await expect(indicator).toContainText(/:/, { timeout: 15_000 });
  });
});
