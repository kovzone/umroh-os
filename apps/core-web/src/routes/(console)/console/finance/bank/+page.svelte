<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let tab = $state<'record' | 'reconcile'>('record');

  // Tab 1: Record Transaction
  let rAccountId = $state('');
  let rRefNo = $state('');
  let rAmount = $state('');
  let rTxDate = $state('');
  let rDescription = $state('');
  let rDirection = $state<'credit' | 'debit'>('credit');
  let rLoading = $state(false);
  let rError = $state('');
  let rSuccess = $state('');

  async function submitRecord(e: SubmitEvent) {
    e.preventDefault(); rLoading = true; rError = ''; rSuccess = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/bank-transactions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ account_id: rAccountId, ref_no: rRefNo, amount: parseFloat(rAmount), tx_date: rTxDate, description: rDescription, direction: rDirection })
      });
      if (!res.ok) throw new Error(`${res.status}`);
      rSuccess = 'Transaksi berhasil dicatat.';
      rRefNo = ''; rAmount = ''; rTxDate = ''; rDescription = '';
    } catch (err: any) { rError = `Gagal mencatat: ${err.message}`; }
    rLoading = false;
  }

  // Tab 2: Reconciliation
  let recAccountId = $state('');
  let recFrom = $state('');
  let recTo = $state('');
  let recLoading = $state(false);
  let recError = $state('');
  let recRows = $state<any[]>([]);

  async function loadRec(e: SubmitEvent) {
    e.preventDefault(); recLoading = true; recError = ''; recRows = [];
    try {
      const params = new URLSearchParams();
      if (recAccountId) params.set('account_id', recAccountId);
      if (recFrom) params.set('from', recFrom);
      if (recTo) params.set('to', recTo);
      const res = await fetch(`${GATEWAY}/v1/finance/bank-reconciliation?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      recRows = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch (err: any) { recError = `Gagal memuat: ${err.message}`; }
    recLoading = false;
  }

  const runningBalance = $derived(
    recRows.reduce((acc: number, r: any) => acc + (r.credit ?? 0) - (r.debit ?? 0), 0)
  );
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Bank & Rekonsiliasi</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Bank & Rekonsiliasi</h2>
      <p>Catat transaksi bank dan lakukan rekonsiliasi</p>
    </div>

    <div class="tabs">
      <button class="tab-btn" class:active={tab === 'record'} onclick={() => tab = 'record'}>Catat Transaksi</button>
      <button class="tab-btn" class:active={tab === 'reconcile'} onclick={() => tab = 'reconcile'}>Rekonsiliasi</button>
    </div>

    {#if tab === 'record'}
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">add_card</span>
          <h3>Catat Transaksi Bank</h3>
        </div>
        <form class="form-body" onsubmit={submitRecord}>
          {#if rError}<div class="alert-err">{rError}</div>{/if}
          {#if rSuccess}<div class="alert-ok">{rSuccess}</div>{/if}
          <div class="form-grid">
            <div class="field">
              <label>ID Rekening *</label>
              <input type="text" bind:value={rAccountId} required placeholder="ACC-001" />
            </div>
            <div class="field">
              <label>No. Referensi *</label>
              <input type="text" bind:value={rRefNo} required placeholder="REF/2024/001" />
            </div>
            <div class="field">
              <label>Jumlah (Rp) *</label>
              <input type="number" bind:value={rAmount} required min="0" placeholder="10000000" />
            </div>
            <div class="field">
              <label>Tanggal Transaksi *</label>
              <input type="date" bind:value={rTxDate} required />
            </div>
            <div class="field full">
              <label>Deskripsi</label>
              <input type="text" bind:value={rDescription} placeholder="Transfer ke vendor hotel" />
            </div>
            <div class="field full">
              <label>Arah *</label>
              <div class="radio-group">
                <label class="radio-label">
                  <input type="radio" bind:group={rDirection} value="credit" /> Kredit (masuk)
                </label>
                <label class="radio-label">
                  <input type="radio" bind:group={rDirection} value="debit" /> Debit (keluar)
                </label>
              </div>
            </div>
          </div>
          <button type="submit" class="btn-primary" disabled={rLoading}>
            {rLoading ? 'Menyimpan...' : 'Catat Transaksi'}
          </button>
        </form>
      </div>

    {:else}
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">sync_alt</span>
          <h3>Rekonsiliasi Bank</h3>
        </div>
        <form class="form-body" onsubmit={loadRec}>
          {#if recError}<div class="alert-err">{recError}</div>{/if}
          <div class="filter-row">
            <div class="field">
              <label>ID Rekening</label>
              <input type="text" bind:value={recAccountId} placeholder="ACC-001" />
            </div>
            <div class="field">
              <label>Dari Tanggal</label>
              <input type="date" bind:value={recFrom} />
            </div>
            <div class="field">
              <label>Sampai Tanggal</label>
              <input type="date" bind:value={recTo} />
            </div>
            <button type="submit" class="btn-primary" disabled={recLoading} style="align-self:flex-end">
              {recLoading ? 'Memuat...' : 'Tampilkan'}
            </button>
          </div>
        </form>

        {#if recRows.length > 0}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Tanggal</th><th>Ref</th><th>Deskripsi</th><th class="ar">Debit</th><th class="ar">Kredit</th><th class="ar">Saldo</th><th>Status</th></tr>
              </thead>
              <tbody>
                {#each recRows as r}
                  <tr>
                    <td>{r.date ?? '-'}</td>
                    <td class="mono">{r.ref_no ?? '-'}</td>
                    <td>{r.description ?? '-'}</td>
                    <td class="ar amount-neg">{r.debit > 0 ? formatIDR(r.debit) : '—'}</td>
                    <td class="ar amount-pos">{r.credit > 0 ? formatIDR(r.credit) : '—'}</td>
                    <td class="ar">{formatIDR(r.balance ?? 0)}</td>
                    <td>
                      {#if r.reconciled}
                        <span class="chip chip-green">Rekonsiliasi</span>
                      {:else}
                        <span class="chip chip-gray">Belum</span>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          <div class="balance-bar">
            Saldo Berjalan: <strong class={runningBalance >= 0 ? 'pos' : 'neg'}>{formatIDR(runningBalance)}</strong>
          </div>
        {:else if !recLoading}
          <div class="empty">Terapkan filter untuk melihat data rekonsiliasi.</div>
        {/if}
      </div>
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
  .tabs { display: flex; gap: 0; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; width: fit-content; margin-bottom: 1rem; overflow: hidden; }
  .tab-btn { border: none; background: #fff; padding: 0.45rem 1rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #737686; border-right: 1px solid rgb(195 198 215 / 0.55); }
  .tab-btn:last-child { border-right: none; }
  .tab-btn.active { background: #2563eb; color: #fff; }
  .card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .card-header { display: flex; align-items: center; gap: 0.5rem; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .card-header .material-symbols-outlined { color: #004ac6; font-size: 1.1rem; }
  .card-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .form-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 1rem; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0.85rem; }
  .full { grid-column: 1 / -1; }
  .filter-row { display: flex; gap: 0.75rem; flex-wrap: wrap; align-items: flex-start; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #15803d; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.82rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .field input:focus { border-color: #2563eb; }
  .radio-group { display: flex; gap: 1.25rem; }
  .radio-label { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; cursor: pointer; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .ar { text-align: right; }
  .mono { font-family: monospace; font-size: 0.72rem; }
  .amount-pos { color: #059669; font-variant-numeric: tabular-nums; }
  .amount-neg { color: #dc2626; font-variant-numeric: tabular-nums; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .chip-green { background: #dcfce7; color: #15803d; }
  .chip-gray { background: #f2f4f6; color: #737686; }
  .balance-bar { padding: 0.75rem 1rem; background: #f7f9fb; border-top: 1px solid rgb(195 198 215 / 0.35); font-size: 0.82rem; color: #434655; }
  .balance-bar .pos { color: #059669; }
  .balance-bar .neg { color: #dc2626; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
</style>
