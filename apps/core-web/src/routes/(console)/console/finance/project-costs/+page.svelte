<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let departureId = $state('');
  let costs = $state<any[]>([]);
  let pl = $state<any>(null);
  let loading = $state(false);
  let error = $state('');

  async function search(e: SubmitEvent) {
    e.preventDefault();
    if (!departureId.trim()) return;
    loading = true; error = ''; costs = []; pl = null;
    try {
      const [cRes, pRes] = await Promise.all([
        fetch(`${GATEWAY}/v1/finance/project-costs/${departureId}`),
        fetch(`${GATEWAY}/v1/finance/departure-pl/${departureId}`)
      ]);
      const cBody = cRes.ok ? await cRes.json() : null;
      const pBody = pRes.ok ? await pRes.json() : null;
      costs = Array.isArray(cBody?.data) ? cBody.data : (Array.isArray(cBody) ? cBody : []);
      pl = pBody?.data ?? pBody ?? null;
      if (!cRes.ok && !pRes.ok) error = 'Data tidak ditemukan untuk keberangkatan ini.';
    } catch { error = 'Gagal memuat data.'; }
    loading = false;
  }

  const variance = $derived(pl ? (pl.revenue_actual ?? 0) - (pl.cost_actual ?? 0) : 0);
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Biaya Proyek & P&L</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Biaya Proyek & P&L Keberangkatan</h2>
      <p>Analisis biaya dan laba rugi per keberangkatan</p>
    </div>

    <form class="search-bar" onsubmit={search}>
      <div class="field">
        <label>ID Keberangkatan</label>
        <input type="text" bind:value={departureId} placeholder="DEP-2024-03" />
      </div>
      <button type="submit" class="btn-primary" disabled={loading}>
        {#if loading}
          <span class="material-symbols-outlined spin">progress_activity</span>
        {:else}
          <span class="material-symbols-outlined">search</span>
        {/if}
        Cari
      </button>
    </form>

    {#if error}<div class="alert-err">{error}</div>{/if}

    {#if costs.length > 0 || pl}
      <div class="grid-2">
        <!-- Costs table -->
        <div class="card">
          <div class="card-header">
            <span class="material-symbols-outlined">table_chart</span>
            <h3>Rincian Biaya</h3>
          </div>
          {#if costs.length === 0}
            <div class="empty">Tidak ada data biaya.</div>
          {:else}
            <div class="table-wrap">
              <table>
                <thead>
                  <tr><th>Kategori</th><th>Deskripsi</th><th class="ar">Budget</th><th class="ar">Aktual</th></tr>
                </thead>
                <tbody>
                  {#each costs as c}
                    <tr>
                      <td><span class="chip">{c.category ?? '-'}</span></td>
                      <td>{c.description ?? '-'}</td>
                      <td class="ar">{formatIDR(c.budget ?? 0)}</td>
                      <td class="ar">{formatIDR(c.actual ?? 0)}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>

        <!-- P&L summary -->
        <div class="card">
          <div class="card-header">
            <span class="material-symbols-outlined">trending_up</span>
            <h3>P&L Keberangkatan</h3>
          </div>
          {#if !pl}
            <div class="empty">Data P&L tidak tersedia.</div>
          {:else}
            <div class="pl-body">
              <div class="pl-row">
                <span>Pendapatan Aktual</span>
                <span class="pos">{formatIDR(pl.revenue_actual ?? 0)}</span>
              </div>
              <div class="pl-row sub">
                <span>Budget Pendapatan</span>
                <span>{formatIDR(pl.revenue_budget ?? 0)}</span>
              </div>
              <div class="pl-divider"></div>
              <div class="pl-row">
                <span>Biaya Aktual</span>
                <span class="neg">{formatIDR(pl.cost_actual ?? 0)}</span>
              </div>
              <div class="pl-row sub">
                <span>Budget Biaya</span>
                <span>{formatIDR(pl.cost_budget ?? 0)}</span>
              </div>
              <div class="pl-divider"></div>
              <div class="pl-row highlight">
                <span>Laba Kotor</span>
                <span class={variance >= 0 ? 'pos' : 'neg'}>{formatIDR(variance)}</span>
              </div>
              {#if pl.variance_pct !== undefined}
                <div class="pl-row">
                  <span>Varians (%)</span>
                  <span class={pl.variance_pct >= 0 ? 'pos' : 'neg'}>{pl.variance_pct?.toFixed(1)}%</span>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </div>
    {:else if !loading}
      <div class="hint">Masukkan ID keberangkatan dan klik Cari untuk melihat data.</div>
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .search-bar { display: flex; gap: 0.75rem; align-items: flex-end; margin-bottom: 1.25rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; width: 220px; }
  .field input:focus { border-color: #2563eb; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-primary .material-symbols-outlined { font-size: 1rem; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  @media (max-width: 768px) { .grid-2 { grid-template-columns: 1fr; } }
  .card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .card-header { display: flex; align-items: center; gap: 0.5rem; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .card-header .material-symbols-outlined { color: #004ac6; font-size: 1.1rem; }
  .card-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:last-child td { border-bottom: 0; }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .pl-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.5rem; }
  .pl-row { display: flex; justify-content: space-between; font-size: 0.82rem; color: #434655; }
  .pl-row.sub { padding-left: 1rem; font-size: 0.75rem; color: #737686; }
  .pl-row.highlight { font-size: 0.95rem; font-weight: 700; color: #191c1e; }
  .pl-row span:last-child { font-variant-numeric: tabular-nums; }
  .pl-divider { height: 1px; background: rgb(195 198 215 / 0.45); margin: 0.25rem 0; }
  .pos { color: #059669; font-weight: 600; }
  .neg { color: #dc2626; font-weight: 600; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
  .hint { text-align: center; color: #b0b3c1; padding: 3rem; font-size: 0.85rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
