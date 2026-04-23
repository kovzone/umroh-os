# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: 04-uat-s1-auth-catalog.spec.ts >> S1 Auth — Permissions & Protected Routes (BL-IAM-002) >> S1-PERM-01: route staff tanpa token → 401
- Location: tests/04-uat-s1-auth-catalog.spec.ts:136:7

# Error details

```
Error: Staff route tanpa token harus 401

expect(received).toBe(expected) // Object.is equality

Expected: 401
Received: 404
```

# Test source

```ts
  40  |     const res = await api.post("/v1/auth/login", {
  41  |       email: UAT_ENV.adminEmail,
  42  |       password: UAT_ENV.adminPassword,
  43  |     });
  44  | 
  45  |     expect(res.status(), "Login harus 200").toBe(200);
  46  |     const body = await res.json();
  47  |     expect(body.data.access_token, "access_token harus PASETO v2.local.*").toMatch(/^v2\.local\./);
  48  |     expect(body.data.refresh_token).toHaveLength(64);
  49  |     expect(body.data.user.email).toBe(UAT_ENV.adminEmail);
  50  |     expect(body.data.user.user_id).toBe(UAT_ENV.adminUserId);
  51  |     expect(body.data.user.status).toBe("active");
  52  |     expect(new Date(body.data.access_expires_at).getTime()).toBeGreaterThan(Date.now());
  53  | 
  54  |     adminToken = body.data.access_token;
  55  |   });
  56  | 
  57  |   test("S1-AUTH-02: login password salah → 401, bukan 500", async () => {
  58  |     const api = await createApiClient(gateway.baseURL);
  59  |     const res = await api.post("/v1/auth/login", {
  60  |       email: UAT_ENV.adminEmail,
  61  |       password: "password-yang-salah",
  62  |     });
  63  | 
  64  |     expect(res.status(), "Password salah harus 401").toBe(401);
  65  |     const body = await res.json();
  66  |     expect(body.error.code).toBe("UNAUTHORIZED");
  67  |   });
  68  | 
  69  |   test("S1-AUTH-03: login email tidak dikenal → 401, bukan 404 (tidak leak user existence)", async () => {
  70  |     const api = await createApiClient(gateway.baseURL);
  71  |     const res = await api.post("/v1/auth/login", {
  72  |       email: "ghost@umrohos.dev",
  73  |       password: "apapun",
  74  |     });
  75  | 
  76  |     expect(res.status(), "Email tidak ada harus 401, bukan 404").toBe(401);
  77  |   });
  78  | 
  79  |   test("S1-AUTH-04: GET /v1/me dengan token valid → 200 + user data", async () => {
  80  |     const api = await createApiClient(gateway.baseURL, adminToken);
  81  |     const res = await api.get("/v1/me");
  82  | 
  83  |     expect(res.status(), "/v1/me dengan token valid harus 200").toBe(200);
  84  |     const body = await res.json();
  85  |     expect(body.data.email).toBe(UAT_ENV.adminEmail);
  86  |   });
  87  | 
  88  |   test("S1-AUTH-05: GET /v1/me tanpa token → 401", async () => {
  89  |     const api = await createApiClient(gateway.baseURL);
  90  |     const res = await api.get("/v1/me");
  91  | 
  92  |     expect(res.status(), "Request tanpa token harus 401").toBe(401);
  93  |   });
  94  | 
  95  |   test("S1-AUTH-06: GET /v1/me dengan token garbage → 401", async () => {
  96  |     const api = await createApiClient(gateway.baseURL, "Bearer ini-bukan-token");
  97  |     const res = await api.get("/v1/me");
  98  | 
  99  |     expect(res.status(), "Token garbage harus 401").toBe(401);
  100 |   });
  101 | 
  102 |   test("S1-AUTH-07: POST /v1/auth/refresh → token baru", async () => {
  103 |     // Login dulu untuk dapat refresh token
  104 |     const api = await createApiClient(gateway.baseURL);
  105 |     const loginRes = await api.post("/v1/auth/login", {
  106 |       email: UAT_ENV.adminEmail,
  107 |       password: UAT_ENV.adminPassword,
  108 |     });
  109 |     expect(loginRes.status()).toBe(200);
  110 |     const loginBody = await loginRes.json();
  111 |     const refreshToken = loginBody.data.refresh_token;
  112 | 
  113 |     // Gunakan refresh token
  114 |     const refreshRes = await api.post("/v1/auth/refresh", {
  115 |       refresh_token: refreshToken,
  116 |     });
  117 | 
  118 |     expect(refreshRes.status(), "Refresh token harus 200").toBe(200);
  119 |     const refreshBody = await refreshRes.json();
  120 |     expect(refreshBody.data.access_token).toMatch(/^v2\.local\./);
  121 |     expect(refreshBody.data.access_token).not.toBe(loginBody.data.access_token);
  122 |   });
  123 | 
  124 |   test("S1-AUTH-08: iam-svc down → gateway harus 502, bukan 200 palsu (fail-closed)", async () => {
  125 |     // Test ini hanya bisa dilakukan dengan mematikan iam-svc — skip di UAT normal
  126 |     // Tandai sebagai informational
  127 |     test.skip(true, "Membutuhkan akses ke docker compose untuk stop iam-svc");
  128 |   });
  129 | });
  130 | 
  131 | // ═══════════════════════════════════════════════════════════════════════════════
  132 | // 2. PERMISSIONS & AUDIT (BL-IAM-002..004)
  133 | // ═══════════════════════════════════════════════════════════════════════════════
  134 | 
  135 | test.describe.serial("S1 Auth — Permissions & Protected Routes (BL-IAM-002)", () => {
  136 |   test("S1-PERM-01: route staff tanpa token → 401", async () => {
  137 |     const api = await createApiClient(gateway.baseURL);
  138 |     // POST /v1/packages adalah staff-only route
  139 |     const res = await api.post("/v1/packages", { name: "test" });
> 140 |     expect(res.status(), "Staff route tanpa token harus 401").toBe(401);
      |                                                               ^ Error: Staff route tanpa token harus 401
  141 |   });
  142 | 
  143 |   test("S1-PERM-02: DELETE /v1/auth/logout dengan token valid → 200", async () => {
  144 |     // Login sementara
  145 |     const api = await createApiClient(gateway.baseURL);
  146 |     const loginRes = await api.post("/v1/auth/login", {
  147 |       email: UAT_ENV.adminEmail,
  148 |       password: UAT_ENV.adminPassword,
  149 |     });
  150 |     expect(loginRes.status()).toBe(200);
  151 |     const token = (await loginRes.json()).data.access_token;
  152 | 
  153 |     const authedApi = await createApiClient(gateway.baseURL, token);
  154 |     const logoutRes = await authedApi.delete("/v1/auth/logout");
  155 | 
  156 |     expect(logoutRes.status(), "Logout harus 200 atau 204").toBeOneOf([200, 204]);
  157 |   });
  158 | 
  159 |   test("S1-PERM-03: token yang sudah di-logout tidak bisa digunakan lagi → 401", async () => {
  160 |     // Login dan logout
  161 |     const api = await createApiClient(gateway.baseURL);
  162 |     const loginRes = await api.post("/v1/auth/login", {
  163 |       email: UAT_ENV.adminEmail,
  164 |       password: UAT_ENV.adminPassword,
  165 |     });
  166 |     const token = (await loginRes.json()).data.access_token;
  167 |     const authedApi = await createApiClient(gateway.baseURL, token);
  168 |     await authedApi.delete("/v1/auth/logout");
  169 | 
  170 |     // Coba gunakan token yang sudah di-revoke
  171 |     const res = await authedApi.get("/v1/me");
  172 |     expect(res.status(), "Token revoked harus 401").toBe(401);
  173 |   });
  174 | });
  175 | 
  176 | // ═══════════════════════════════════════════════════════════════════════════════
  177 | // 3. CATALOG READ — Public Endpoints (BL-CAT-001..004)
  178 | // ═══════════════════════════════════════════════════════════════════════════════
  179 | 
  180 | test.describe.serial("S1 Catalog — Public Read (BL-CAT-001..004)", () => {
  181 |   test("S1-CAT-01: GET /v1/packages → hanya package active yang tampil", async () => {
  182 |     const api = await createApiClient(gateway.baseURL);
  183 |     const res = await api.get("/v1/packages");
  184 | 
  185 |     expect(res.status(), "List packages harus 200").toBe(200);
  186 |     const body = await res.json();
  187 |     expect(Array.isArray(body.packages), "packages harus array").toBe(true);
  188 |     expect(body.packages.length, "Harus ada setidaknya 1 package active").toBeGreaterThan(0);
  189 | 
  190 |     // Package active harus ada
  191 |     const ids = body.packages.map((p: { id: string }) => p.id);
  192 |     expect(ids).toContain(UAT_ENV.activePkgId);
  193 | 
  194 |     // Package draft TIDAK boleh muncul
  195 |     expect(ids, "Package draft tidak boleh muncul di list publik").not.toContain(UAT_ENV.draftPkgId);
  196 | 
  197 |     // Shape validasi
  198 |     const active = body.packages.find((p: { id: string }) => p.id === UAT_ENV.activePkgId);
  199 |     expect(active).toBeTruthy();
  200 |     expect(active.starting_price?.list_currency).toBe("IDR");
  201 |     expect(active.starting_price?.list_amount).toBeGreaterThan(0);
  202 |   });
  203 | 
  204 |   test("S1-CAT-02: GET /v1/packages/{active_id} → detail lengkap", async () => {
  205 |     const api = await createApiClient(gateway.baseURL);
  206 |     const res = await api.get(`/v1/packages/${UAT_ENV.activePkgId}`);
  207 | 
  208 |     expect(res.status(), "Package detail harus 200").toBe(200);
  209 |     const body = await res.json();
  210 |     expect(body.data?.id || body.id).toBe(UAT_ENV.activePkgId);
  211 |   });
  212 | 
  213 |   test("S1-CAT-03: GET /v1/packages/{draft_id} → 404 (draft tidak boleh publik)", async () => {
  214 |     const api = await createApiClient(gateway.baseURL);
  215 |     const res = await api.get(`/v1/packages/${UAT_ENV.draftPkgId}`);
  216 | 
  217 |     expect(res.status(), "Package draft harus 404 di endpoint publik").toBe(404);
  218 |   });
  219 | 
  220 |   test("S1-CAT-04: GET /v1/package-departures/{active_dep_id} → detail + remaining_seats", async () => {
  221 |     const api = await createApiClient(gateway.baseURL);
  222 |     const res = await api.get(`/v1/package-departures/${UAT_ENV.activeDepId}`);
  223 | 
  224 |     expect(res.status(), "Departure detail harus 200").toBe(200);
  225 |     const body = await res.json();
  226 |     const dep = body.data || body;
  227 |     expect(dep.id).toBe(UAT_ENV.activeDepId);
  228 |     expect(typeof dep.remaining_seats).toBe("number");
  229 |     expect(dep.remaining_seats).toBeGreaterThanOrEqual(0);
  230 |   });
  231 | 
  232 |   test("S1-CAT-05: GET /v1/package-departures/{cancelled_dep} → 404", async () => {
  233 |     const api = await createApiClient(gateway.baseURL);
  234 |     const cancelledDepId = "dep_01JCDF00000000000000000003";
  235 |     const res = await api.get(`/v1/package-departures/${cancelledDepId}`);
  236 | 
  237 |     expect(res.status(), "Departure cancelled harus 404").toBe(404);
  238 |   });
  239 | 
  240 |   test("S1-CAT-06: GET /v1/packages tanpa auth → 200 (public endpoint)", async () => {
```