<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let departure = $state('');
  let dueDate = $state('');
  let notes = $state('');
  let loading = $state(false);
  let error = $state('');
  let result = $state<any>(null);

  const recentRuns = $state([
    { id: 'br_001', departure: 'DEP-2024-03', count: 12, total: 240_000_000, date: '2024-03-01', status: 'selesai' },
    { id: 'br_002', departure: 'DEP-2024-04', count: 8, total: 160_000_000, date: '2024-04-01', status: 'selesai' },
    { id: 'br_003', departure: 'DEP-2024-05', count: 15, total: 300_000_000, date: '2024-05-01', status: 'proses' },
  ]);

  async function submit(e: SubmitEvent) {
    e.preventDefault(); loading = true; error = ''; result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/billing/schedule`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ departure_id: departure, due_date: dueDate, notes })
      });
      const body = res.ok ? await res.json() : null;
      if (!res.ok) throw new Error(body?.message ?? `${res.status}`);
      result = body;
    } catch (err: any) { error = err.message ?? 'Gagal membuat tagihan.'; }
    loading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Tagihan Otomatis</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Tagihan Otomatis</h2>
      <p>Jadwalkan pembuatan invoice massal untuk keberangkatan</p>
    </div>

    <div class="grid-2">
      <!-- Form -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">schedule_send</span>
          <h3>Jadwalkan Tagihan</h3>
        </div>
        <form onsubmit={submit} class="form-body">
          {#if error}<div class="alert-err">{error}</div>{/if}
          {#if result}
            <div class="alert-ok">
              <span class="material-symbols-outlined">check_circle</span>
              <strong>{result.invoice_count ?? result.count ?? '?'} invoice</strong> dibuat — Total {formatIDR(result.total_amount ?? result.total ?? 0)}
            </div>
          {/if}
          <div class="field">
            <label>ID Keberangkatan *</label>
            <input type="text" bind:value={departure} required placeholder="DEP-2024-06" />
          </div>
          <div class="field">
            <label>Tanggal Jatuh Tempo *</label>
            <input type="date" bind:value={dueDate} required />
          </div>
          <div class="field">
            <label>Catatan</label>
            <textarea bind:value={notes} rows="3" placeholder="Catatan opsional..."></textarea>
          </div>
          <button type="submit" class="btn-primary" disabled={loading}>
            {#if loading}
              <span class="material-symbols-outlined spin">progress_activity</span> Memproses...
            {:else}
              <span class="material-symbols-outlined">send</span> Buat Tagihan
            {/if}
          </button>
        </form>
      </div>

      <!-- Recent runs -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">history</span>
          <h3>Riwayat Pembuatan Tagihan</h3>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>Keberangkatan</th><th>Invoice</th><th>Total</th><th>Tanggal</th><th>Status</th></tr>
            </thead>
            <tbody>
              {#each recentRuns as r}
                <tr>
                  <td class="fw">{r.departure}</td>
                  <td>{r.count}</td>
                  <td class="mono">{formatIDR(r.total)}</td>
                  <td>{r.date}</td>
                  <td><span class="chip" class:chip-green={r.status === 'selesai'} class:chip-yellow={r.status === 'proses'}>{r.status}</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
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
  .form-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 1rem; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; }
  .alert-ok { display: flex; align-items: center; gap: 0.5rem; background: #f0fdf4; border: 1px solid #bbf7d0; color: #15803d; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.82rem; }
  .alert-ok .material-symbols-outlined { font-size: 1.1rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; resize: vertical; }
  .field input:focus, .field textarea:focus { border-color: #2563eb; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.55rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-primary .material-symbols-outlined { font-size: 1rem; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 0.85rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:last-child td { border-bottom: 0; }
  .fw { font-weight: 600; }
  .mono { font-variant-numeric: tabular-nums; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .chip-green { background: #dcfce7; color: #15803d; }
  .chip-yellow { background: #fef9c3; color: #a16207; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
