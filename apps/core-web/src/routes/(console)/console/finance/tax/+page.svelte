<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  // Calculator
  let cBase = $state('');
  let cType = $state('PPN');
  let cRate = $state('11');
  let calcLoading = $state(false);
  let calcResult = $state<any>(null);
  let calcError = $state('');

  async function calculate(e: SubmitEvent) {
    e.preventDefault(); calcLoading = true; calcError = ''; calcResult = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/tax/calculate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ base_amount: parseFloat(cBase), tax_type: cType, rate: parseFloat(cRate) })
      });
      const body = res.ok ? await res.json() : null;
      if (!res.ok) throw new Error(`${res.status}`);
      calcResult = body?.data ?? body;
    } catch (err: any) { calcError = `Gagal: ${err.message}`; }
    calcLoading = false;
  }

  // Report
  let rFrom = $state('');
  let rTo = $state('');
  let rType = $state('');
  let reportLoading = $state(false);
  let reportRows = $state<any[]>([]);
  let reportError = $state('');

  async function loadReport(e: SubmitEvent) {
    e.preventDefault(); reportLoading = true; reportError = ''; reportRows = [];
    try {
      const params = new URLSearchParams();
      if (rFrom) params.set('from', rFrom);
      if (rTo) params.set('to', rTo);
      if (rType) params.set('tax_type', rType);
      const res = await fetch(`${GATEWAY}/v1/finance/tax/report?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      reportRows = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch (err: any) { reportError = `Gagal: ${err.message}`; }
    reportLoading = false;
  }

  // Auto-fill rate based on tax type
  $effect(() => {
    if (cType === 'PPN') cRate = '11';
    else if (cType === 'PPh 21') cRate = '5';
    else if (cType === 'PPh 23') cRate = '2';
  });

  const taxAmount = $derived(calcResult ? (calcResult.tax_amount ?? (parseFloat(cBase || '0') * parseFloat(cRate || '0') / 100)) : null);
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Perpajakan</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Perpajakan (PPN/PPh)</h2>
      <p>Kalkulator pajak dan laporan perpajakan</p>
    </div>

    <div class="grid-2">
      <!-- Calculator -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">calculate</span>
          <h3>Kalkulator Pajak</h3>
        </div>
        <form class="form-body" onsubmit={calculate}>
          {#if calcError}<div class="alert-err">{calcError}</div>{/if}
          <div class="field">
            <label>Dasar Pengenaan Pajak (Rp) *</label>
            <input type="number" bind:value={cBase} required min="0" placeholder="100000000" />
          </div>
          <div class="field">
            <label>Jenis Pajak *</label>
            <select bind:value={cType}>
              <option value="PPN">PPN</option>
              <option value="PPh 21">PPh 21</option>
              <option value="PPh 23">PPh 23</option>
            </select>
          </div>
          <div class="field">
            <label>Tarif (%)</label>
            <input type="number" bind:value={cRate} step="0.1" placeholder="11" />
          </div>
          <button type="submit" class="btn-primary" disabled={calcLoading}>
            {calcLoading ? 'Menghitung...' : 'Hitung Pajak'}
          </button>

          {#if calcResult || (cBase && cRate)}
            <div class="tax-result">
              <div class="tax-row">
                <span>DPP</span>
                <span>{formatIDR(parseFloat(cBase || '0'))}</span>
              </div>
              <div class="tax-row">
                <span>Tarif</span>
                <span>{cRate}%</span>
              </div>
              <div class="tax-divider"></div>
              <div class="tax-row tax-total">
                <span>{cType}</span>
                <span class="tax-val">{formatIDR(parseFloat(cBase || '0') * parseFloat(cRate || '0') / 100)}</span>
              </div>
              <div class="tax-row tax-gross">
                <span>Total Bruto</span>
                <span>{formatIDR(parseFloat(cBase || '0') * (1 + parseFloat(cRate || '0') / 100))}</span>
              </div>
            </div>
          {/if}
        </form>
      </div>

      <!-- Report -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">summarize</span>
          <h3>Laporan Pajak</h3>
        </div>
        <form class="form-body" onsubmit={loadReport}>
          {#if reportError}<div class="alert-err">{reportError}</div>{/if}
          <div class="field">
            <label>Dari Tanggal</label>
            <input type="date" bind:value={rFrom} />
          </div>
          <div class="field">
            <label>Sampai Tanggal</label>
            <input type="date" bind:value={rTo} />
          </div>
          <div class="field">
            <label>Jenis Pajak</label>
            <select bind:value={rType}>
              <option value="">Semua</option>
              <option value="PPN">PPN</option>
              <option value="PPh 21">PPh 21</option>
              <option value="PPh 23">PPh 23</option>
            </select>
          </div>
          <button type="submit" class="btn-primary" disabled={reportLoading}>
            {reportLoading ? 'Memuat...' : 'Tampilkan Laporan'}
          </button>
        </form>

        {#if reportRows.length > 0}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Tanggal</th><th>Referensi</th><th>Jenis</th><th class="ar">DPP</th><th class="ar">Pajak</th></tr>
              </thead>
              <tbody>
                {#each reportRows as r}
                  <tr>
                    <td>{r.date ?? '-'}</td>
                    <td class="mono">{r.reference ?? '-'}</td>
                    <td><span class="chip">{r.tax_type ?? '-'}</span></td>
                    <td class="ar">{formatIDR(r.base_amount ?? 0)}</td>
                    <td class="ar fw">{formatIDR(r.tax_amount ?? 0)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else if !reportLoading}
          <div class="empty">Terapkan filter untuk melihat laporan.</div>
        {/if}
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
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  @media (max-width: 768px) { .grid-2 { grid-template-columns: 1fr; } }
  .card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .card-header { display: flex; align-items: center; gap: 0.5rem; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .card-header .material-symbols-outlined { color: #004ac6; font-size: 1.1rem; }
  .card-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .form-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .btn-primary { display: inline-flex; align-items: center; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .tax-result { background: #f7f9fb; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.35rem; padding: 1rem; display: flex; flex-direction: column; gap: 0.45rem; }
  .tax-row { display: flex; justify-content: space-between; font-size: 0.8rem; color: #434655; }
  .tax-row span:last-child { font-variant-numeric: tabular-nums; }
  .tax-total { font-weight: 600; color: #191c1e; }
  .tax-gross { font-weight: 700; color: #004ac6; }
  .tax-val { font-size: 0.95rem; color: #dc2626; }
  .tax-divider { height: 1px; background: rgb(195 198 215 / 0.45); margin: 0.1rem 0; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .mono { font-family: monospace; font-size: 0.72rem; }
  .fw { font-weight: 600; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
</style>
