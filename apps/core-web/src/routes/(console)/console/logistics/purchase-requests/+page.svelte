<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  // Filters
  let filterStatus = $state('all');
  let filterDeparture = $state('');
  let pr_loading = $state(false);
  let pr_error = $state('');
  let pr_list = $state<any[]>([]);

  async function loadPRs() {
    pr_loading = true; pr_error = '';
    try {
      const params = new URLSearchParams();
      if (filterStatus !== 'all') params.set('status', filterStatus);
      if (filterDeparture.trim()) params.set('departure_id', filterDeparture.trim());
      const res = await fetch(`${GATEWAY}/v1/logistics/purchase-requests?${params}`);
      if (!res.ok) throw new Error(`Gagal memuat PR (${res.status})`);
      const body = await res.json();
      pr_list = body.purchase_requests ?? body.prs ?? body ?? [];
    } catch (err) {
      pr_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    pr_loading = false;
  }

  $effect(() => { loadPRs(); });

  // Budget sidebar
  let bs_departureId = $state('');
  let bs_loading = $state(false);
  let bs_data = $state<any>(null);

  async function loadBudget() {
    if (!bs_departureId.trim()) return;
    bs_loading = true;
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/budget-sync/${bs_departureId}`);
      if (res.ok) bs_data = await res.json();
    } catch { /* ignore */ }
    bs_loading = false;
  }

  // Tiered Approvals
  let ta_loading = $state(false);
  let ta_list = $state<any[]>([]);

  async function loadApprovals() {
    ta_loading = true;
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/tiered-approvals`);
      if (res.ok) {
        const body = await res.json();
        ta_list = body.approvals ?? body ?? [];
      }
    } catch { /* ignore */ }
    ta_loading = false;
  }

  $effect(() => { loadApprovals(); });

  async function decideApproval(id: string, decision: 'approve' | 'reject') {
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/tiered-approvals/${id}/decision`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ decision }),
      });
      if (res.ok) {
        const updated = await res.json();
        ta_list = ta_list.map(a => a.id === id ? updated : a);
      }
    } catch { /* ignore */ }
  }

  const PR_STATUSES = [
    { value: 'all', label: 'Semua Status' },
    { value: 'draft', label: 'Draft' },
    { value: 'pending', label: 'Menunggu' },
    { value: 'approved', label: 'Disetujui' },
    { value: 'rejected', label: 'Ditolak' },
    { value: 'ordered', label: 'Dipesan' },
  ];

  const STATUS_CLS: Record<string, string> = {
    draft: 'chip-gray',
    pending: 'chip-yellow',
    approved: 'chip-green',
    rejected: 'chip-red',
    ordered: 'chip-blue',
  };

  function formatDate(iso: string) {
    try { return new Date(iso).toLocaleDateString('id-ID'); } catch { return iso; }
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/logistics" class="bc-link">Logistik</a>
      <span class="bc-sep">/</span>
      <span>Permintaan Pembelian</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Permintaan Pembelian (PR)</h2>
      <p>BL-LOG-013/014/015 — Kelola PR, tiered approvals, dan budget sync</p>
    </div>

    <div class="layout-cols">
      <!-- Main area -->
      <div class="main-area">
        <!-- Filters -->
        <div class="section-block">
          <div class="filter-row">
            <div class="field">
              <label for="f-status">Status</label>
              <select id="f-status" bind:value={filterStatus} onchange={loadPRs}>
                {#each PR_STATUSES as s}
                  <option value={s.value}>{s.label}</option>
                {/each}
              </select>
            </div>
            <div class="field">
              <label for="f-dep">ID Keberangkatan</label>
              <input id="f-dep" type="text" placeholder="dep-001" bind:value={filterDeparture}
                oninput={loadPRs} />
            </div>
            <div class="field field-actions">
              <button class="btn-ghost" onclick={loadPRs} disabled={pr_loading}>
                {#if pr_loading}<span class="spinner-dark"></span>{/if}
                Muat Ulang
              </button>
            </div>
          </div>
          {#if pr_error}<div class="alert-err">{pr_error}</div>{/if}
        </div>

        <!-- PR Table -->
        <div class="section-block">
          <h3 class="section-title">Daftar Permintaan Pembelian</h3>
          {#if pr_loading && pr_list.length === 0}
            <div class="loading-row"><span class="spinner-dark"></span> Memuat...</div>
          {:else if pr_list.length === 0}
            <div class="empty-state">Belum ada data PR.</div>
          {:else}
            <div class="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>PR ID</th>
                    <th>Keberangkatan</th>
                    <th>Deskripsi</th>
                    <th class="ar">Jumlah</th>
                    <th>Status</th>
                    <th>Tanggal</th>
                  </tr>
                </thead>
                <tbody>
                  {#each pr_list as pr}
                    <tr>
                      <td class="mono">{pr.id ?? '-'}</td>
                      <td class="mono">{pr.departure_id ?? '-'}</td>
                      <td>{pr.description ?? '-'}</td>
                      <td class="ar">{formatIDR(pr.total_amount ?? pr.amount ?? 0)}</td>
                      <td><span class="chip {STATUS_CLS[pr.status] ?? 'chip-gray'}">{pr.status ?? '-'}</span></td>
                      <td>{pr.created_at ? formatDate(pr.created_at) : '-'}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>

        <!-- Tiered Approvals -->
        <div class="section-block">
          <h3 class="section-title">Tiered Approvals</h3>
          {#if ta_loading}
            <div class="loading-row"><span class="spinner-dark"></span> Memuat...</div>
          {:else if ta_list.length === 0}
            <div class="empty-state">Tidak ada approval yang menunggu.</div>
          {:else}
            <div class="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>PR</th>
                    <th>Level</th>
                    <th>Pemohon</th>
                    <th>Status</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {#each ta_list as ap}
                    <tr>
                      <td class="mono">{ap.id ?? '-'}</td>
                      <td class="mono">{ap.pr_id ?? '-'}</td>
                      <td><span class="chip chip-blue">Level {ap.level ?? ap.tier ?? '-'}</span></td>
                      <td class="mono">{ap.requester_id ?? '-'}</td>
                      <td><span class="chip {STATUS_CLS[ap.status] ?? 'chip-gray'}">{ap.status ?? '-'}</span></td>
                      <td>
                        {#if ap.status === 'pending'}
                          <div class="action-group">
                            <button class="btn-approve" onclick={() => decideApproval(ap.id, 'approve')}>Setujui</button>
                            <button class="btn-reject" onclick={() => decideApproval(ap.id, 'reject')}>Tolak</button>
                          </div>
                        {:else}
                          <span class="text-muted">—</span>
                        {/if}
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>
      </div>

      <!-- Budget Sidebar -->
      <div class="sidebar">
        <div class="section-block">
          <h3 class="section-title">Budget vs Aktual</h3>
          <div class="field">
            <label for="bs-dep">ID Keberangkatan</label>
            <input id="bs-dep" type="text" placeholder="dep-001" bind:value={bs_departureId}
              oninput={loadBudget} />
          </div>
          {#if bs_loading}
            <div class="loading-row mt"><span class="spinner-dark"></span></div>
          {:else if bs_data}
            <div class="budget-items">
              <div class="budget-row">
                <span>Total Budget</span>
                <strong>{formatIDR(bs_data.total_budget ?? 0)}</strong>
              </div>
              <div class="budget-row">
                <span>Aktual Terpakai</span>
                <strong class="text-red">{formatIDR(bs_data.actual_spent ?? 0)}</strong>
              </div>
              <div class="budget-row">
                <span>Sisa Anggaran</span>
                <strong class="text-green">{formatIDR((bs_data.total_budget ?? 0) - (bs_data.actual_spent ?? 0))}</strong>
              </div>
              {#if bs_data.utilization_pct !== undefined}
                <div class="progress-bar-wrap mt">
                  <div class="progress-bar-fill" style="width:{Math.min(bs_data.utilization_pct, 100)}%;background:{bs_data.utilization_pct > 90 ? '#dc2626' : '#2563eb'}"></div>
                </div>
                <div class="pct-label">{bs_data.utilization_pct?.toFixed(1)}% terpakai</div>
              {/if}
            </div>
          {:else}
            <div class="empty-state-sm">Masukkan ID keberangkatan</div>
          {/if}
        </div>
      </div>
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 88rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .layout-cols { display: grid; grid-template-columns: 1fr 280px; gap: 1.25rem; align-items: start; }
  @media (max-width: 900px) { .layout-cols { grid-template-columns: 1fr; } }
  .main-area { display: flex; flex-direction: column; gap: 1.25rem; }
  .sidebar { }
  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; }
  .section-title { margin: 0 0 1rem; font-size: 0.9rem; font-weight: 700; }
  .filter-row { display: flex; gap: 0.75rem; flex-wrap: wrap; align-items: flex-end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; }
  .field-actions { align-self: flex-end; }
  .btn-ghost { display: inline-flex; align-items: center; gap: 0.3rem; border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.42rem 0.7rem; font-size: 0.78rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #191c1e; }
  .btn-ghost:disabled { opacity: 0.6; cursor: not-allowed; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .ar { text-align: right; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .chip { display: inline-flex; padding: 0.12rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .chip-blue { background: #e0f2fe; color: #075985; }
  .chip-green { background: #dcfce7; color: #166534; }
  .chip-red { background: #fee2e2; color: #991b1b; }
  .chip-yellow { background: #fef9c3; color: #854d0e; }
  .chip-gray { background: #f2f4f6; color: #434655; }
  .action-group { display: flex; gap: 0.4rem; }
  .btn-approve { padding: 0.25rem 0.55rem; font-size: 0.7rem; font-weight: 600; border: 0; border-radius: 0.2rem; background: #dcfce7; color: #166534; cursor: pointer; font-family: inherit; }
  .btn-reject { padding: 0.25rem 0.55rem; font-size: 0.7rem; font-weight: 600; border: 0; border-radius: 0.2rem; background: #fee2e2; color: #991b1b; cursor: pointer; font-family: inherit; }
  .text-muted { color: #b0b3c1; }
  .spinner-dark { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(100 100 100 / 0.25); border-top-color: #434655; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; color: #737686; padding: 0.75rem 0; }
  .empty-state { text-align: center; color: #b0b3c1; padding: 1.5rem; font-size: 0.82rem; }
  .empty-state-sm { text-align: center; color: #b0b3c1; padding: 1rem 0; font-size: 0.78rem; }
  .budget-items { display: flex; flex-direction: column; gap: 0.6rem; margin-top: 0.75rem; }
  .budget-row { display: flex; justify-content: space-between; align-items: center; font-size: 0.8rem; }
  .budget-row span { color: #737686; }
  .text-red { color: #dc2626; }
  .text-green { color: #059669; }
  .progress-bar-wrap { height: 0.4rem; background: #e8eaec; border-radius: 999px; overflow: hidden; }
  .progress-bar-fill { height: 100%; border-radius: 999px; transition: width 0.4s; }
  .pct-label { font-size: 0.68rem; color: #737686; text-align: right; margin-top: 0.2rem; }
  .mt { margin-top: 0.65rem; }
</style>
