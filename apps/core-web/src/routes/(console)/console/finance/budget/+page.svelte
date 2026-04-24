<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let fromDate = $state('');
  let toDate = $state('');
  let departureId = $state('');
  let rows = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  async function load(e?: SubmitEvent) {
    if (e) e.preventDefault();
    loading = true; error = '';
    try {
      const params = new URLSearchParams();
      if (fromDate) params.set('from', fromDate);
      if (toDate) params.set('to', toDate);
      if (departureId) params.set('departure_id', departureId);
      const res = await fetch(`${GATEWAY}/v1/finance/budget-vs-actual?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      rows = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { error = 'Gagal memuat data.'; rows = []; }
    loading = false;
  }

  $effect(() => { load(); });

  const totals = $derived({
    budget: rows.reduce((a: number, r: any) => a + (r.budget ?? 0), 0),
    actual: rows.reduce((a: number, r: any) => a + (r.actual ?? 0), 0),
    selisih: rows.reduce((a: number, r: any) => a + ((r.actual ?? 0) - (r.budget ?? 0)), 0),
  });

  function pct(actual: number, budget: number): string {
    if (!budget) return '-';
    return ((actual / budget) * 100).toFixed(1) + '%';
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Budget vs Aktual</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Budget vs Aktual</h2>
      <p>Pantau realisasi anggaran per akun</p>
    </div>

    <form class="filter-bar" onsubmit={load}>
      <div class="field">
        <label>Dari Tanggal</label>
        <input type="date" bind:value={fromDate} />
      </div>
      <div class="field">
        <label>Sampai Tanggal</label>
        <input type="date" bind:value={toDate} />
      </div>
      <div class="field">
        <label>ID Keberangkatan</label>
        <input type="text" bind:value={departureId} placeholder="Opsional" />
      </div>
      <button type="submit" class="btn-primary" disabled={loading} style="align-self:flex-end">
        {loading ? 'Memuat...' : 'Terapkan'}
      </button>
    </form>

    {#if error}<div class="alert-err">{error}</div>{/if}

    <div class="panel">
      {#if loading}
        <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span> Memuat...</div>
      {:else if rows.length === 0}
        <div class="empty">Data tidak tersedia. Coba ubah filter.</div>
      {:else}
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>Kode Akun</th><th>Nama Akun</th><th class="ar">Budget</th><th class="ar">Aktual</th><th class="ar">Selisih</th><th class="ar">%</th></tr>
            </thead>
            <tbody>
              {#each rows as r}
                {@const selisih = (r.actual ?? 0) - (r.budget ?? 0)}
                {@const p = r.budget ? (r.actual / r.budget) * 100 : 0}
                <tr>
                  <td class="mono">{r.account_code ?? '-'}</td>
                  <td>{r.account_name ?? '-'}</td>
                  <td class="ar">{formatIDR(r.budget ?? 0)}</td>
                  <td class="ar">{formatIDR(r.actual ?? 0)}</td>
                  <td class="ar" class:pos={selisih <= 0} class:neg={selisih > 0}>{formatIDR(Math.abs(selisih))}{selisih > 0 ? ' ↑' : ' ↓'}</td>
                  <td class="ar" class:warn={p > 100}>{pct(r.actual, r.budget)}</td>
                </tr>
              {/each}
            </tbody>
            <tfoot>
              {@const ts = totals.selisih}
              <tr class="totals-row">
                <td colspan="2"><strong>Total</strong></td>
                <td class="ar"><strong>{formatIDR(totals.budget)}</strong></td>
                <td class="ar"><strong>{formatIDR(totals.actual)}</strong></td>
                <td class="ar" class:pos={ts <= 0} class:neg={ts > 0}><strong>{formatIDR(Math.abs(ts))}{ts > 0 ? ' ↑' : ' ↓'}</strong></td>
                <td class="ar"><strong>{pct(totals.actual, totals.budget)}</strong></td>
              </tr>
            </tfoot>
          </table>
        </div>
      {/if}
    </div>
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
  .filter-bar { display: flex; gap: 0.75rem; align-items: flex-start; flex-wrap: wrap; margin-bottom: 1.25rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 0.85rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tfoot td { background: #f7f9fb; border-top: 2px solid rgb(195 198 215 / 0.55); }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .mono { font-family: monospace; font-size: 0.72rem; }
  .pos { color: #059669; }
  .neg { color: #dc2626; }
  .warn { color: #d97706; font-weight: 700; }
  .empty { text-align: center; color: #b0b3c1; padding: 3rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
