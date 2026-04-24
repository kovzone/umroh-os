<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let tab = $state<'ar' | 'ap'>('ar');

  // AR
  let arBookingId = $state('');
  let arPilgrimId = $state('');
  let arFrom = $state('');
  let arTo = $state('');
  let arRows = $state<any[]>([]);
  let arLoading = $state(false);
  let arError = $state('');

  // AP
  let apVendorId = $state('');
  let apFrom = $state('');
  let apTo = $state('');
  let apRows = $state<any[]>([]);
  let apLoading = $state(false);
  let apError = $state('');

  async function loadAR(e: SubmitEvent) {
    e.preventDefault(); arLoading = true; arError = '';
    try {
      const params = new URLSearchParams();
      if (arBookingId) params.set('booking_id', arBookingId);
      if (arPilgrimId) params.set('pilgrim_id', arPilgrimId);
      if (arFrom) params.set('from', arFrom);
      if (arTo) params.set('to', arTo);
      const res = await fetch(`${GATEWAY}/v1/finance/ar-subledger?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      arRows = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
      if (!arRows.length) arError = 'Tidak ada data ditemukan.';
    } catch (err: any) { arError = `Gagal: ${err.message}`; }
    arLoading = false;
  }

  async function loadAP(e: SubmitEvent) {
    e.preventDefault(); apLoading = true; apError = '';
    try {
      const params = new URLSearchParams();
      if (apVendorId) params.set('vendor_id', apVendorId);
      if (apFrom) params.set('from', apFrom);
      if (apTo) params.set('to', apTo);
      const res = await fetch(`${GATEWAY}/v1/finance/ap-subledger?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      apRows = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
      if (!apRows.length) apError = 'Tidak ada data ditemukan.';
    } catch (err: any) { apError = `Gagal: ${err.message}`; }
    apLoading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Buku Besar Pembantu</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Buku Besar Pembantu</h2>
      <p>Laporan AR (Piutang) dan AP (Utang)</p>
    </div>

    <div class="tabs">
      <button class="tab-btn" class:active={tab === 'ar'} onclick={() => tab = 'ar'}>AR (Piutang)</button>
      <button class="tab-btn" class:active={tab === 'ap'} onclick={() => tab = 'ap'}>AP (Utang)</button>
    </div>

    {#if tab === 'ar'}
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">arrow_circle_up</span>
          <h3>Piutang (AR)</h3>
        </div>
        <form class="filter-bar" onsubmit={loadAR}>
          {#if arError}<div class="alert-err">{arError}</div>{/if}
          <div class="field">
            <label>ID Booking</label>
            <input type="text" bind:value={arBookingId} placeholder="BK-2024-001" />
          </div>
          <div class="field">
            <label>ID Jamaah</label>
            <input type="text" bind:value={arPilgrimId} placeholder="PLG-001" />
          </div>
          <div class="field">
            <label>Dari</label>
            <input type="date" bind:value={arFrom} />
          </div>
          <div class="field">
            <label>Sampai</label>
            <input type="date" bind:value={arTo} />
          </div>
          <button type="submit" class="btn-primary" disabled={arLoading} style="align-self:flex-end">
            {arLoading ? 'Memuat...' : 'Tampilkan'}
          </button>
        </form>
        {#if arRows.length > 0}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Tanggal</th><th>Deskripsi</th><th class="ar">Debet</th><th class="ar">Kredit</th><th class="ar">Saldo</th></tr>
              </thead>
              <tbody>
                {#each arRows as r}
                  <tr>
                    <td>{r.date ?? '-'}</td>
                    <td>{r.description ?? '-'}</td>
                    <td class="ar amount-neg">{r.debit > 0 ? formatIDR(r.debit) : '—'}</td>
                    <td class="ar amount-pos">{r.credit > 0 ? formatIDR(r.credit) : '—'}</td>
                    <td class="ar">{formatIDR(r.balance ?? 0)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else if !arLoading}
          <div class="empty">Gunakan filter di atas untuk menampilkan data piutang.</div>
        {/if}
      </div>

    {:else}
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">arrow_circle_down</span>
          <h3>Utang (AP)</h3>
        </div>
        <form class="filter-bar" onsubmit={loadAP}>
          {#if apError}<div class="alert-err">{apError}</div>{/if}
          <div class="field">
            <label>ID Vendor</label>
            <input type="text" bind:value={apVendorId} placeholder="VND-001" />
          </div>
          <div class="field">
            <label>Dari</label>
            <input type="date" bind:value={apFrom} />
          </div>
          <div class="field">
            <label>Sampai</label>
            <input type="date" bind:value={apTo} />
          </div>
          <button type="submit" class="btn-primary" disabled={apLoading} style="align-self:flex-end">
            {apLoading ? 'Memuat...' : 'Tampilkan'}
          </button>
        </form>
        {#if apRows.length > 0}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Tanggal</th><th>Deskripsi</th><th class="ar">Debet</th><th class="ar">Kredit</th><th class="ar">Saldo</th></tr>
              </thead>
              <tbody>
                {#each apRows as r}
                  <tr>
                    <td>{r.date ?? '-'}</td>
                    <td>{r.description ?? '-'}</td>
                    <td class="ar amount-neg">{r.debit > 0 ? formatIDR(r.debit) : '—'}</td>
                    <td class="ar amount-pos">{r.credit > 0 ? formatIDR(r.credit) : '—'}</td>
                    <td class="ar">{formatIDR(r.balance ?? 0)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else if !apLoading}
          <div class="empty">Gunakan filter di atas untuk menampilkan data utang.</div>
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
  .filter-bar { display: flex; gap: 0.75rem; align-items: flex-start; flex-wrap: wrap; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.25); }
  .alert-err { width: 100%; background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.55rem 0.75rem; font-size: 0.78rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 0.85rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:last-child td { border-bottom: 0; }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .amount-pos { color: #059669; }
  .amount-neg { color: #dc2626; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
</style>
