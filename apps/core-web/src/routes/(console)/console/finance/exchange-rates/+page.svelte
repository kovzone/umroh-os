<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let fFrom = $state('USD');
  let fTo = $state('IDR');
  let fRate = $state('');
  let fDate = $state('');
  let formLoading = $state(false);
  let formError = $state('');
  let formSuccess = $state('');

  let rates = $state<any[]>([]);
  let ratesLoading = $state(false);
  let qFrom = $state('USD');
  let qTo = $state('IDR');

  async function loadRates() {
    ratesLoading = true;
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/exchange-rates?from=${qFrom}&to=${qTo}`);
      const body = res.ok ? await res.json() : null;
      rates = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { rates = []; }
    ratesLoading = false;
  }

  $effect(() => { loadRates(); });

  async function submitRate(e: SubmitEvent) {
    e.preventDefault(); formLoading = true; formError = ''; formSuccess = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/exchange-rates`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ from_currency: fFrom, to_currency: fTo, rate: parseFloat(fRate), effective_date: fDate })
      });
      if (!res.ok) throw new Error(`${res.status}`);
      formSuccess = `Kurs ${fFrom}/${fTo} berhasil disimpan.`;
      fRate = ''; fDate = '';
      await loadRates();
    } catch (err: any) { formError = `Gagal: ${err.message}`; }
    formLoading = false;
  }

  const currencies = ['IDR', 'USD', 'SAR', 'EUR', 'MYR'];
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Nilai Tukar Mata Uang</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Nilai Tukar Mata Uang</h2>
      <p>Kelola kurs mata uang untuk multi-currency</p>
    </div>

    <div class="grid-2">
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">currency_exchange</span>
          <h3>Set Kurs Baru</h3>
        </div>
        <form class="form-body" onsubmit={submitRate}>
          {#if formError}<div class="alert-err">{formError}</div>{/if}
          {#if formSuccess}<div class="alert-ok">{formSuccess}</div>{/if}
          <div class="field">
            <label>Mata Uang Asal *</label>
            <select bind:value={fFrom}>
              {#each currencies as c}<option value={c}>{c}</option>{/each}
            </select>
          </div>
          <div class="field">
            <label>Mata Uang Tujuan *</label>
            <select bind:value={fTo}>
              {#each currencies as c}<option value={c}>{c}</option>{/each}
            </select>
          </div>
          <div class="field">
            <label>Kurs *</label>
            <input type="text" bind:value={fRate} required placeholder="15750.50" />
          </div>
          <div class="field">
            <label>Berlaku Mulai *</label>
            <input type="date" bind:value={fDate} required />
          </div>
          <button type="submit" class="btn-primary" disabled={formLoading}>
            {formLoading ? 'Menyimpan...' : 'Simpan Kurs'}
          </button>
        </form>
      </div>

      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">table_chart</span>
          <h3>Kurs Saat Ini</h3>
        </div>
        <div class="filter-bar-inner">
          <div class="field-inline">
            <label>Dari:</label>
            <select bind:value={qFrom} onchange={loadRates}>
              {#each currencies as c}<option value={c}>{c}</option>{/each}
            </select>
          </div>
          <div class="field-inline">
            <label>Ke:</label>
            <select bind:value={qTo} onchange={loadRates}>
              {#each currencies as c}<option value={c}>{c}</option>{/each}
            </select>
          </div>
          <button class="btn-ghost" onclick={loadRates}>
            <span class="material-symbols-outlined">refresh</span>
          </button>
        </div>
        {#if ratesLoading}
          <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span></div>
        {:else if rates.length === 0}
          <div class="empty">Belum ada data kurs.</div>
        {:else}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Dari</th><th>Ke</th><th class="ar">Kurs</th><th>Berlaku Mulai</th></tr>
              </thead>
              <tbody>
                {#each rates as r}
                  <tr>
                    <td><span class="currency-chip">{r.from_currency}</span></td>
                    <td><span class="currency-chip">{r.to_currency}</span></td>
                    <td class="ar mono">{r.rate?.toLocaleString('id-ID')}</td>
                    <td>{r.effective_date ?? '-'}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
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
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #15803d; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.82rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-ghost { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.42rem 0.55rem; cursor: pointer; font-family: inherit; display: inline-flex; align-items: center; }
  .btn-ghost .material-symbols-outlined { font-size: 1rem; color: #737686; }
  .filter-bar-inner { display: flex; align-items: center; gap: 0.75rem; padding: 0.85rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.25); flex-wrap: wrap; }
  .field-inline { display: flex; align-items: center; gap: 0.4rem; }
  .field-inline label { font-size: 0.75rem; font-weight: 600; color: #434655; }
  .field-inline select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.35rem 0.55rem; font-size: 0.78rem; color: #191c1e; font-family: inherit; outline: none; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 0.85rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:last-child td { border-bottom: 0; }
  .ar { text-align: right; }
  .mono { font-family: monospace; font-variant-numeric: tabular-nums; }
  .currency-chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; background: #f0fdf4; color: #15803d; font-size: 0.7rem; font-weight: 700; letter-spacing: 0.05em; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; justify-content: center; padding: 1.5rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1.2rem; color: #737686; }
</style>
