<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let balance = $state(0);
  let entries = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  // Form
  let fAmount = $state('');
  let fDirection = $state<'masuk' | 'keluar'>('keluar');
  let fDescription = $state('');
  let fCategory = $state('ATK');
  let fDate = $state('');
  let fEvidenceUrl = $state('');
  let formLoading = $state(false);
  let formError = $state('');
  let formSuccess = $state('');

  // Close period dialog
  let showCloseDialog = $state(false);
  let closeLoading = $state(false);

  async function load() {
    loading = true; error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/petty-cash`);
      const body = res.ok ? await res.json() : null;
      balance = body?.balance ?? body?.current_balance ?? 0;
      entries = Array.isArray(body?.data) ? body.data : (Array.isArray(body?.entries) ? body.entries : []);
    } catch { error = 'Gagal memuat data kas kecil.'; }
    loading = false;
  }

  $effect(() => { load(); });

  async function submitEntry(e: SubmitEvent) {
    e.preventDefault(); formLoading = true; formError = ''; formSuccess = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/petty-cash`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount: parseFloat(fAmount), direction: fDirection, description: fDescription, category: fCategory, date: fDate, evidence_url: fEvidenceUrl })
      });
      if (!res.ok) throw new Error(`${res.status}`);
      formSuccess = 'Entri berhasil dicatat.';
      fAmount = ''; fDescription = ''; fDate = ''; fEvidenceUrl = '';
      await load();
    } catch (err: any) { formError = `Gagal: ${err.message}`; }
    formLoading = false;
  }

  async function closePeriod() {
    closeLoading = true;
    try {
      await fetch(`${GATEWAY}/v1/finance/petty-cash/close-period`, { method: 'POST' });
      showCloseDialog = false;
      await load();
    } catch { alert('Gagal menutup periode.'); }
    closeLoading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Kas Kecil</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Kas Kecil</h2>
        <p>Kelola pengeluaran kas kecil harian</p>
      </div>
      <button class="btn-warning" onclick={() => showCloseDialog = true}>
        <span class="material-symbols-outlined">lock</span>
        Tutup Periode
      </button>
    </div>

    {#if error}<div class="alert-err">{error}</div>{/if}

    <!-- Balance card -->
    <div class="balance-card">
      <div class="balance-label">Saldo Kas Kecil</div>
      <div class="balance-val">{formatIDR(balance)}</div>
    </div>

    <div class="grid-2">
      <!-- Entry form -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">add_circle</span>
          <h3>Catat Entri</h3>
        </div>
        <form class="form-body" onsubmit={submitEntry}>
          {#if formError}<div class="alert-err">{formError}</div>{/if}
          {#if formSuccess}<div class="alert-ok">{formSuccess}</div>{/if}
          <div class="field">
            <label>Jumlah (Rp) *</label>
            <input type="number" bind:value={fAmount} required min="0" placeholder="50000" />
          </div>
          <div class="field">
            <label>Arah *</label>
            <div class="radio-group">
              <label class="radio-label"><input type="radio" bind:group={fDirection} value="masuk" /> Masuk</label>
              <label class="radio-label"><input type="radio" bind:group={fDirection} value="keluar" /> Keluar</label>
            </div>
          </div>
          <div class="field">
            <label>Kategori</label>
            <select bind:value={fCategory}>
              <option value="ATK">ATK</option>
              <option value="Transport">Transport</option>
              <option value="Makan">Makan</option>
              <option value="Lainnya">Lainnya</option>
            </select>
          </div>
          <div class="field">
            <label>Deskripsi *</label>
            <input type="text" bind:value={fDescription} required placeholder="Beli alat tulis kantor" />
          </div>
          <div class="field">
            <label>Tanggal *</label>
            <input type="date" bind:value={fDate} required />
          </div>
          <div class="field">
            <label>URL Bukti</label>
            <input type="url" bind:value={fEvidenceUrl} placeholder="https://..." />
          </div>
          <button type="submit" class="btn-primary" disabled={formLoading}>
            {formLoading ? 'Menyimpan...' : 'Catat Entri'}
          </button>
        </form>
      </div>

      <!-- Recent entries -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">history</span>
          <h3>Riwayat Entri</h3>
        </div>
        {#if loading}
          <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span> Memuat...</div>
        {:else if entries.length === 0}
          <div class="empty">Belum ada entri kas kecil.</div>
        {:else}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Tanggal</th><th>Deskripsi</th><th>Kategori</th><th>Arah</th><th class="ar">Jumlah</th></tr>
              </thead>
              <tbody>
                {#each entries as e}
                  <tr>
                    <td>{e.date ?? '-'}</td>
                    <td>{e.description ?? '-'}</td>
                    <td><span class="chip">{e.category ?? '-'}</span></td>
                    <td>
                      {#if e.direction === 'masuk'}
                        <span class="chip chip-green">Masuk</span>
                      {:else}
                        <span class="chip chip-red">Keluar</span>
                      {/if}
                    </td>
                    <td class="ar" class:pos={e.direction === 'masuk'} class:neg={e.direction === 'keluar'}>{formatIDR(e.amount ?? 0)}</td>
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

<!-- Close period dialog -->
{#if showCloseDialog}
  <div class="overlay" role="presentation" onclick={() => showCloseDialog = false}></div>
  <div class="modal">
    <div class="modal-header">
      <h3>Tutup Periode Kas Kecil</h3>
    </div>
    <div class="modal-body">
      <p>Apakah Anda yakin ingin menutup periode kas kecil? Semua entri periode ini akan dikunci.</p>
      <p><strong>Saldo saat ini: {formatIDR(balance)}</strong></p>
      <div class="modal-actions">
        <button class="btn-ghost" onclick={() => showCloseDialog = false}>Batal</button>
        <button class="btn-warning" onclick={closePeriod} disabled={closeLoading}>
          {closeLoading ? 'Menutup...' : 'Ya, Tutup Periode'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 1.25rem; gap: 1rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #15803d; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.82rem; }
  .balance-card { background: linear-gradient(135deg, #004ac6, #2563eb); color: #fff; border-radius: 0.75rem; padding: 1.5rem 2rem; margin-bottom: 1.25rem; display: inline-block; min-width: 280px; }
  .balance-label { font-size: 0.75rem; font-weight: 600; opacity: 0.8; letter-spacing: 0.06em; text-transform: uppercase; }
  .balance-val { font-size: 1.8rem; font-weight: 800; font-variant-numeric: tabular-nums; margin-top: 0.25rem; }
  .grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  @media (max-width: 768px) { .grid-2 { grid-template-columns: 1fr; } }
  .card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .card-header { display: flex; align-items: center; gap: 0.5rem; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .card-header .material-symbols-outlined { color: #004ac6; font-size: 1.1rem; }
  .card-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .form-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .radio-group { display: flex; gap: 1.25rem; }
  .radio-label { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; cursor: pointer; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-warning { display: inline-flex; align-items: center; gap: 0.35rem; background: #d97706; color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-warning:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-warning .material-symbols-outlined { font-size: 1rem; }
  .btn-ghost { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.45rem 0.85rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #191c1e; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .pos { color: #059669; }
  .neg { color: #dc2626; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .chip-green { background: #dcfce7; color: #15803d; }
  .chip-red { background: #fef2f2; color: #dc2626; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }
  .overlay { position: fixed; inset: 0; background: rgb(0 0 0 / 0.4); z-index: 40; }
  .modal { position: fixed; top: 50%; left: 50%; transform: translate(-50%,-50%); z-index: 50; background: #fff; border-radius: 0.5rem; width: 420px; max-width: 95vw; box-shadow: 0 8px 32px rgb(0 0 0 / 0.18); overflow: hidden; }
  .modal-header { padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); }
  .modal-header h3 { margin: 0; font-size: 1rem; font-weight: 700; }
  .modal-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; }
  .modal-body p { margin: 0; font-size: 0.85rem; color: #434655; }
  .modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
