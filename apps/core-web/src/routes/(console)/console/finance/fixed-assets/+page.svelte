<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let assets = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  // Add modal
  let showAddModal = $state(false);
  let fName = $state('');
  let fCategory = $state('kendaraan');
  let fPurchaseDate = $state('');
  let fPurchaseCost = $state('');
  let fUsefulLife = $state('');
  let fResidualValue = $state('');
  let addLoading = $state(false);
  let addError = $state('');
  let addSuccess = $state('');

  // Depreciation
  let depAsOf = $state('');
  let depLoading = $state(false);
  let depResult = $state<any>(null);
  let depError = $state('');

  async function load() {
    loading = true; error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/fixed-assets`);
      const body = res.ok ? await res.json() : null;
      assets = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { error = 'Gagal memuat data aset.'; assets = []; }
    loading = false;
  }

  $effect(() => { load(); });

  async function submitAsset(e: SubmitEvent) {
    e.preventDefault(); addLoading = true; addError = ''; addSuccess = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/fixed-assets`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: fName, category: fCategory, purchase_date: fPurchaseDate, purchase_cost: parseFloat(fPurchaseCost), useful_life_months: parseInt(fUsefulLife), residual_value: parseFloat(fResidualValue || '0') })
      });
      if (!res.ok) throw new Error(`${res.status}`);
      addSuccess = 'Aset berhasil ditambahkan.';
      await load();
      setTimeout(() => { showAddModal = false; addSuccess = ''; fName = ''; fPurchaseCost = ''; fUsefulLife = ''; fResidualValue = ''; }, 1200);
    } catch (err: any) { addError = `Gagal: ${err.message}`; }
    addLoading = false;
  }

  async function calcDepreciation(e: SubmitEvent) {
    e.preventDefault(); depLoading = true; depError = ''; depResult = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/depreciation`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ as_of: depAsOf })
      });
      const body = res.ok ? await res.json() : null;
      if (!res.ok) throw new Error(`${res.status}`);
      depResult = body?.data ?? body;
      await load();
    } catch (err: any) { depError = `Gagal: ${err.message}`; }
    depLoading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Aset Tetap</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Aset Tetap</h2>
        <p>Inventaris dan perhitungan penyusutan aset</p>
      </div>
      <button class="btn-primary" onclick={() => showAddModal = true}>
        <span class="material-symbols-outlined">add</span>
        Tambah Aset
      </button>
    </div>

    {#if error}<div class="alert-err">{error}</div>{/if}

    <!-- Depreciation calc -->
    <div class="dep-card">
      <form class="dep-form" onsubmit={calcDepreciation}>
        <span class="material-symbols-outlined" style="color:#d97706">calculate</span>
        <strong>Hitung Penyusutan per:</strong>
        <input type="date" bind:value={depAsOf} required />
        <button type="submit" class="btn-warn" disabled={depLoading}>
          {depLoading ? 'Menghitung...' : 'Hitung'}
        </button>
      </form>
      {#if depError}<div class="alert-err" style="margin-top:0.5rem">{depError}</div>{/if}
      {#if depResult}
        <div class="dep-result">
          <span class="material-symbols-outlined" style="color:#059669">check_circle</span>
          <strong>{depResult.asset_count ?? depResult.count ?? '?'} aset</strong> diproses —
          Total penyusutan: <strong>{formatIDR(depResult.total_depreciation ?? depResult.total ?? 0)}</strong>
        </div>
      {/if}
    </div>

    <div class="panel">
      {#if loading}
        <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span> Memuat...</div>
      {:else if assets.length === 0}
        <div class="empty">Belum ada aset. Klik "Tambah Aset".</div>
      {:else}
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>Nama Aset</th><th>Kategori</th><th>Tgl Beli</th><th class="ar">Biaya Perolehan</th><th class="ar">Akum. Penyusutan</th><th class="ar">Nilai Buku</th></tr>
            </thead>
            <tbody>
              {#each assets as a}
                <tr>
                  <td class="fw">{a.name}</td>
                  <td><span class="chip">{a.category}</span></td>
                  <td>{a.purchase_date ?? '-'}</td>
                  <td class="ar">{formatIDR(a.purchase_cost ?? 0)}</td>
                  <td class="ar amount-neg">{formatIDR(a.accumulated_depreciation ?? 0)}</td>
                  <td class="ar fw">{formatIDR((a.purchase_cost ?? 0) - (a.accumulated_depreciation ?? 0))}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  </section>
</main>

<!-- Add Asset Modal -->
{#if showAddModal}
  <div class="overlay" role="presentation" onclick={() => showAddModal = false}></div>
  <div class="modal">
    <div class="modal-header">
      <h3>Tambah Aset Tetap</h3>
      <button class="close-btn" onclick={() => showAddModal = false}><span class="material-symbols-outlined">close</span></button>
    </div>
    <form class="modal-body" onsubmit={submitAsset}>
      {#if addError}<div class="alert-err">{addError}</div>{/if}
      {#if addSuccess}<div class="alert-ok">{addSuccess}</div>{/if}
      <div class="field">
        <label>Nama Aset *</label>
        <input type="text" bind:value={fName} required placeholder="Mobil Toyota Hiace" />
      </div>
      <div class="field">
        <label>Kategori *</label>
        <select bind:value={fCategory}>
          <option value="kendaraan">Kendaraan</option>
          <option value="peralatan">Peralatan</option>
          <option value="bangunan">Bangunan</option>
          <option value="lainnya">Lainnya</option>
        </select>
      </div>
      <div class="field">
        <label>Tanggal Beli *</label>
        <input type="date" bind:value={fPurchaseDate} required />
      </div>
      <div class="field">
        <label>Biaya Perolehan (Rp) *</label>
        <input type="number" bind:value={fPurchaseCost} required min="0" placeholder="350000000" />
      </div>
      <div class="field">
        <label>Umur Manfaat (bulan) *</label>
        <input type="number" bind:value={fUsefulLife} required min="1" placeholder="60" />
      </div>
      <div class="field">
        <label>Nilai Sisa (Rp)</label>
        <input type="number" bind:value={fResidualValue} min="0" placeholder="0" />
      </div>
      <div class="modal-actions">
        <button type="button" class="btn-ghost" onclick={() => showAddModal = false}>Batal</button>
        <button type="submit" class="btn-primary" disabled={addLoading}>{addLoading ? 'Menyimpan...' : 'Simpan'}</button>
      </div>
    </form>
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
  .dep-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1rem 1.25rem; margin-bottom: 1.25rem; }
  .dep-form { display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; font-size: 0.82rem; color: #434655; }
  .dep-form .material-symbols-outlined { font-size: 1.2rem; }
  .dep-form input { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.8rem; color: #191c1e; font-family: inherit; outline: none; }
  .btn-warn { background: #d97706; color: #fff; border: none; border-radius: 0.25rem; padding: 0.42rem 0.85rem; font-size: 0.78rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-warn:disabled { opacity: 0.6; cursor: not-allowed; }
  .dep-result { display: flex; align-items: center; gap: 0.5rem; margin-top: 0.65rem; font-size: 0.82rem; color: #434655; }
  .dep-result .material-symbols-outlined { font-size: 1.1rem; }
  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 0.85rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .fw { font-weight: 600; }
  .amount-neg { color: #dc2626; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .empty { text-align: center; color: #b0b3c1; padding: 3rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary .material-symbols-outlined { font-size: 1rem; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-ghost { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.45rem 0.85rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #191c1e; }
  .overlay { position: fixed; inset: 0; background: rgb(0 0 0 / 0.4); z-index: 40; }
  .modal { position: fixed; top: 50%; left: 50%; transform: translate(-50%,-50%); z-index: 50; background: #fff; border-radius: 0.5rem; width: 440px; max-width: 95vw; box-shadow: 0 8px 32px rgb(0 0 0 / 0.18); overflow: hidden; }
  .modal-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); }
  .modal-header h3 { margin: 0; font-size: 1rem; font-weight: 700; }
  .close-btn { border: none; background: transparent; cursor: pointer; color: #737686; display: grid; place-items: center; }
  .close-btn .material-symbols-outlined { font-size: 1.2rem; }
  .modal-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; max-height: 70vh; overflow-y: auto; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
