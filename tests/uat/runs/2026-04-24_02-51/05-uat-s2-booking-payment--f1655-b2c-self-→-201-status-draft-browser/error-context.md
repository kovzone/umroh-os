# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: 05-uat-s2-booking-payment.spec.ts >> S2 Booking — Draft Creation (BL-BOOK-001) >> S2-BOOK-01: POST /v1/bookings (b2c_self) → 201, status draft
- Location: tests/05-uat-s2-booking-payment.spec.ts:46:7

# Error details

```
Error: Login admin gagal: HTTP 404 — <!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="../../favicon.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>UmrohOS — core-web</title>
    
		<link href="../../_app/immutable/assets/0.D2qjOvIQ.css" rel="stylesheet">
  </head>
  <body data-sveltekit-preload-data="hover">
    <div style="display: contents"><!--[--><!--[0--><!--[--><div class="app svelte-12qhfyh"><!--[0--><header class="topbar svelte-1elxaub"><div class="brand svelte-1elxaub"><a href="/" class="brand-link svelte-1elxaub"><strong>UmrohOS</strong></a> <span class="muted svelte-1elxaub">core-web</span></div> <nav class="nav svelte-1elxaub"><a href="/packages" class="nav-link svelte-1elxaub" data-testid="nav-packages">Packages</a> <button type="button" disabled="" title="Sign-in lands with F1.5 — iam-svc auth" class="btn-signin svelte-1elxaub" data-testid="signin-button">Sign in</button></nav></header><!--]--> <main class="container svelte-12qhfyh"><!--[-1--><!--[--><h1>404</h1> <p>Not Found</p><!--]--><!--]--><!----></main> <!--[0--><footer class="footer svelte-jz8lnl"><div class="inner svelte-jz8lnl"><nav class="svelte-jz8lnl"><a href="/system/status" data-testid="footer-status-link" class="svelte-jz8lnl">Service status</a></nav> <p class="meta svelte-jz8lnl">UmrohOS · core-web · scaffold v0.1</p></div></footer><!--]--></div><!--]--><!--]--> <!--[-1--><!--]--><!--]-->
			
			<script>
				{
					__sveltekit_ni0pth = {
						base: new URL("../..", location).pathname.slice(0, -1)
					};

					const element = document.currentScript.parentElement;

					Promise.all([
						import("../../_app/immutable/entry/start.Dxu-Iad9.js"),
						import("../../_app/immutable/entry/app.B5rq3UOE.js")
					]).then(([kit, app]) => {
						kit.start(app, element, {
							node_ids: [0, 1],
							data: [null],
							form: null,
							error: {message:"Not Found"},
							status: 404
						});
					});
				}
			</script>
		</div>
  </body>
</html>

Cek apakah server up di http://216.176.238.161 dan migration seed sudah jalan.
```

```
error: relation "crm.lead_status_history" does not exist
```

# Test source

```ts
  206 |         full_name: `${UAT_PREFIX} Test Jamaah`,
  207 |         email,
  208 |         whatsapp: "+628112345678",
  209 |         domicile: "Jakarta",
  210 |         is_lead: true,
  211 |       },
  212 |     ],
  213 |     notes: `${UAT_PREFIX} automated test booking`,
  214 |     ...overrides,
  215 |   });
  216 | 
  217 |   if (res.status() !== 201 && res.status() !== 200) {
  218 |     const body = await res.text();
  219 |     throw new Error(`createUatBooking gagal: HTTP ${res.status()} — ${body}`);
  220 |   }
  221 | 
  222 |   const body = await res.json();
  223 |   return {
  224 |     id: body.data?.id || body.id,
  225 |     packageId,
  226 |     departureId,
  227 |   };
  228 | }
  229 | 
  230 | /**
  231 |  * Create invoice untuk booking (staff call).
  232 |  */
  233 | export async function createUatInvoice(
  234 |   api: ApiClient,
  235 |   bookingId: string
  236 | ): Promise<UatInvoice> {
  237 |   const res = await api.post("/v1/invoices", {
  238 |     booking_id: bookingId,
  239 |     gateway: "mock",
  240 |   });
  241 | 
  242 |   if (res.status() !== 201 && res.status() !== 200) {
  243 |     const body = await res.text();
  244 |     throw new Error(`createUatInvoice gagal: HTTP ${res.status()} — ${body}`);
  245 |   }
  246 | 
  247 |   const body = await res.json();
  248 |   return { id: body.data?.id || body.id, bookingId };
  249 | }
  250 | 
  251 | /**
  252 |  * Buat lead UAT via public endpoint.
  253 |  */
  254 | export async function createUatLead(
  255 |   overrides: Record<string, unknown> = {}
  256 | ): Promise<UatLead> {
  257 |   const api = await createApiClient(UAT_ENV.gatewayUrl);
  258 |   const email = uatEmail("lead");
  259 |   const res = await api.post("/v1/leads", {
  260 |     name: `${UAT_PREFIX} Test Lead`,
  261 |     email,
  262 |     phone: "+628119876543",
  263 |     message: "Test lead dari UAT automated testing",
  264 |     source_note: "uat-test",
  265 |     ...overrides,
  266 |   });
  267 | 
  268 |   if (res.status() !== 201 && res.status() !== 200) {
  269 |     const body = await res.text();
  270 |     throw new Error(`createUatLead gagal: HTTP ${res.status()} — ${body}`);
  271 |   }
  272 | 
  273 |   const body = await res.json();
  274 |   return { id: body.data?.id || body.id, email };
  275 | }
  276 | 
  277 | // ─── Cleanup ─────────────────────────────────────────────────────────────────
  278 | 
  279 | /**
  280 |  * Hapus package UAT via API (soft-delete / archive).
  281 |  * Gagal silent — sudah terhapus atau tidak ada = OK.
  282 |  */
  283 | export async function deleteUatPackage(
  284 |   api: ApiClient,
  285 |   packageId: string
  286 | ): Promise<void> {
  287 |   try {
  288 |     await api.delete(`/v1/packages/${packageId}`);
  289 |   } catch {
  290 |     // Silent — already deleted or not found
  291 |   }
  292 | }
  293 | 
  294 | /**
  295 |  * Cleanup via DB langsung: hapus semua data [UAT] sekaligus.
  296 |  * Lebih efisien daripada call API satu-satu.
  297 |  * Panggil ini di afterAll test suite.
  298 |  */
  299 | export async function cleanupUatData(): Promise<void> {
  300 |   const client = createDbClient();
  301 |   await client.connect();
  302 |   try {
  303 |     await client.query("BEGIN");
  304 | 
  305 |     // CRM
> 306 |     await client.query(
      |     ^ error: relation "crm.lead_status_history" does not exist
  307 |       `DELETE FROM crm.lead_status_history WHERE lead_id IN (
  308 |         SELECT id FROM crm.leads WHERE email ILIKE 'uat.%@%' OR name ILIKE '[UAT]%'
  309 |       )`
  310 |     );
  311 |     await client.query(
  312 |       `DELETE FROM crm.leads WHERE email ILIKE 'uat.%@%' OR name ILIKE '[UAT]%'`
  313 |     );
  314 | 
  315 |     // Payment
  316 |     await client.query(
  317 |       `DELETE FROM payment.payment_events WHERE invoice_id IN (
  318 |         SELECT pi.id FROM payment.invoices pi
  319 |         JOIN booking.bookings b ON b.id = pi.booking_id
  320 |         WHERE b.notes ILIKE '%[UAT]%'
  321 |       )`
  322 |     );
  323 |     await client.query(
  324 |       `DELETE FROM payment.virtual_accounts WHERE invoice_id IN (
  325 |         SELECT pi.id FROM payment.invoices pi
  326 |         JOIN booking.bookings b ON b.id = pi.booking_id
  327 |         WHERE b.notes ILIKE '%[UAT]%'
  328 |       )`
  329 |     );
  330 |     await client.query(
  331 |       `DELETE FROM payment.invoices WHERE booking_id IN (
  332 |         SELECT id FROM booking.bookings WHERE notes ILIKE '%[UAT]%'
  333 |       )`
  334 |     );
  335 | 
  336 |     // Finance
  337 |     await client.query(
  338 |       `DELETE FROM finance.journal_lines WHERE entry_id IN (
  339 |         SELECT je.id FROM finance.journal_entries je
  340 |         WHERE je.source_id::text IN (
  341 |           SELECT pi.id::text FROM payment.invoices pi
  342 |           JOIN booking.bookings b ON b.id = pi.booking_id
  343 |           WHERE b.notes ILIKE '%[UAT]%'
  344 |         )
  345 |       )`
  346 |     );
  347 |     await client.query(
  348 |       `DELETE FROM finance.journal_entries WHERE source_id::text IN (
  349 |         SELECT pi.id::text FROM payment.invoices pi
  350 |         JOIN booking.bookings b ON b.id = pi.booking_id
  351 |         WHERE b.notes ILIKE '%[UAT]%'
  352 |       )`
  353 |     );
  354 | 
  355 |     // Logistics
  356 |     await client.query(
  357 |       `DELETE FROM logistics.fulfillment_tasks WHERE booking_id IN (
  358 |         SELECT id FROM booking.bookings WHERE notes ILIKE '%[UAT]%'
  359 |       )`
  360 |     );
  361 | 
  362 |     // Booking
  363 |     await client.query(
  364 |       `DELETE FROM booking.pilgrim_documents WHERE jamaah_id IN (
  365 |         SELECT bj.id FROM booking.jamaah bj
  366 |         JOIN booking.bookings b ON b.id = bj.booking_id
  367 |         WHERE b.notes ILIKE '%[UAT]%'
  368 |       )`
  369 |     );
  370 |     await client.query(
  371 |       `DELETE FROM booking.jamaah WHERE booking_id IN (
  372 |         SELECT id FROM booking.bookings WHERE notes ILIKE '%[UAT]%'
  373 |       )`
  374 |     );
  375 |     await client.query(
  376 |       `DELETE FROM booking.bookings WHERE notes ILIKE '%[UAT]%'`
  377 |     );
  378 | 
  379 |     // Catalog
  380 |     await client.query(
  381 |       `DELETE FROM catalog.departures WHERE package_id IN (
  382 |         SELECT id FROM catalog.packages WHERE name ILIKE '[UAT]%'
  383 |       )`
  384 |     );
  385 |     await client.query(
  386 |       `DELETE FROM catalog.packages WHERE name ILIKE '[UAT]%'`
  387 |     );
  388 | 
  389 |     await client.query("COMMIT");
  390 |   } catch (err) {
  391 |     await client.query("ROLLBACK");
  392 |     console.error("UAT cleanup gagal, manual cleanup diperlukan:", err);
  393 |     throw err;
  394 |   } finally {
  395 |     await client.end();
  396 |   }
  397 | }
  398 | 
```