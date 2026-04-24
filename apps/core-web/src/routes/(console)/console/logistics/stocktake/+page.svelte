<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  // Start stocktake
  let st_warehouseId = $state('');
  let st_notes = $state('');
  let st_loading = $state(false);
  let st_error = $state('');

  // Active stocktake
  let activeStocktake = $state<any>(null);

  // Count form
  let ct_sku = $state('');
  let ct_qty = $state('');
  let ct_loading = $state(false);
  let ct_error = $state('');
  let ct_success = $state('');
  let ct_entries = $state<any[]>([]);

  // Finalize
  let fn_loading = $state(false);
  let fn_error = $state('');
  let fn_variance = $state<any>(null);

  async function startStocktake(e: Event) {
    e.preventDefault();
    st_loading = true; st_error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/stocktake`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ warehouse_id: st_warehouseId, notes: st_notes }),
      });
      if (!res.ok) throw new Error(`Gagal mulai stok opname (${res.status})`);
      activeStocktake = await res.json();
      st_warehouseId = ''; st_notes = '';
    } catch (err) {
      st_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    st_loading = false;
  }

  async function recordCount(e: Event) {
    e.preventDefault();
    if (!activeStocktake?.id) return;
    ct_loading = true; ct_error = ''; ct_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/stocktake/${activeStocktake.id}/count`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sku: ct_sku, counted_qty: parseFloat(ct_qty) }),
      });
      if (!res.ok) throw new Error(`Gagal catat hitungan (${res.status})`);
      const entry = await res.json();
      ct_entries = [...ct_entries, entry];
      ct_success = `SKU ${ct_sku} — ${ct_qty} unit dicatat.`;
      ct_sku = ''; ct_qty = '';
    } catch (err) {
      ct_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ct_loading = false;
  }

  async function finalizeStocktake() {
    if (!activeStocktake?.id) return;
    fn_loading = true; fn_error = ''; fn_variance = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/stocktake/${activeStocktake.id}/finalize`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
      });
      if (!res.ok) throw new Error(`Gagal finalisasi (${res.status})`);
      fn_variance = await res.json();
      activeStocktake = { ...activeStocktake, status: 'finalized' };
    } catch (err) {
      fn_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    fn_loading = false;
  }

  const varianceItems = $derived(
    fn_variance && Array.isArray(fn_variance.variance_items) ? fn_variance.variance_items : []
  );
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/logistics" class="bc-link">Logistik</a>
      <span class="bc-sep">/</span>
      <span>Stok Opname Digital</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Stok Opname Digital</h2>
      <p>BL-LOG-024 — Mulai, catat hitungan, dan finalisasi stok opname</p>
    </div>

    {#if !activeStocktake}
      <!-- Start stocktake -->
      <div class="section-block">
        <h3 class="section-title">Mulai Stok Opname</h3>
        <form class="form-row" onsubmit={startStocktake}>
          <div class="field">
            <label for="st-wh">ID Gudang</label>
            <input id="st-wh" type="text" placeholder="wh-001" bind:value={st_warehouseId} required />
          </div>
          <div class="field field-wide">
            <label for="st-notes">Catatan</label>
            <input id="st-notes" type="text" placeholder="Stok opname bulanan..." bind:value={st_notes} />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={st_loading}>
              {#if st_loading}<span class="spinner"></span>{/if}
              Mulai Stok Opname
            </button>
          </div>
        </form>
        {#if st_error}<div class="alert-err">{st_error}</div>{/if}
      </div>
    {:else}
      <!-- Active stocktake -->
      <div class="section-block active-banner">
        <div class="active-info">
          <span class="material-symbols-outlined active-icon">inventory</span>
          <div>
            <div class="active-title">Stok Opname Aktif</div>
            <div class="active-sub">ID: <strong class="mono">{activeStocktake.id}</strong>
              &nbsp;·&nbsp; Gudang: <strong>{activeStocktake.warehouse_id ?? '-'}</strong>
              &nbsp;·&nbsp; Status: <span class="chip {activeStocktake.status === 'finalized' ? 'chip-green' : 'chip-blue'}">{activeStocktake.status ?? 'active'}</span>
            </div>
          </div>
        </div>
      </div>

      {#if activeStocktake.status !== 'finalized'}
        <!-- Count form -->
        <div class="section-block">
          <h3 class="section-title">Catat Hitungan Fisik</h3>
          <form class="form-row" onsubmit={recordCount}>
            <div class="field">
              <label for="ct-sku">SKU</label>
              <input id="ct-sku" type="text" placeholder="SKU-001" bind:value={ct_sku} required />
            </div>
            <div class="field">
              <label for="ct-qty">Jumlah Terhitung</label>
              <input id="ct-qty" type="number" min="0" step="0.01" placeholder="100" bind:value={ct_qty} required />
            </div>
            <div class="field field-actions">
              <button type="submit" class="btn-primary" disabled={ct_loading}>
                {#if ct_loading}<span class="spinner"></span>{/if}
                Catat
              </button>
            </div>
          </form>
          {#if ct_error}<div class="alert-err">{ct_error}</div>{/if}
          {#if ct_success}<div class="alert-ok">{ct_success}</div>{/if}
        </div>

        {#if ct_entries.length > 0}
          <div class="section-block">
            <h3 class="section-title">Hitungan Sesi Ini ({ct_entries.length} SKU)</h3>
            <div class="table-wrap">
              <table>
                <thead>
                  <tr><th>SKU</th><th class="ar">Qty Terhitung</th><th class="ar">Qty Sistem</th><th class="ar">Variance</th></tr>
                </thead>
                <tbody>
                  {#each ct_entries as entry}
                    {@const variance = (entry.counted_qty ?? 0) - (entry.system_qty ?? 0)}
                    <tr>
                      <td class="mono">{entry.sku ?? '-'}</td>
                      <td class="ar">{entry.counted_qty ?? '-'}</td>
                      <td class="ar">{entry.system_qty ?? '-'}</td>
                      <td class="ar {variance > 0 ? 'pos' : variance < 0 ? 'neg' : ''}">
                        {variance > 0 ? '+' : ''}{variance}
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          </div>
        {/if}

        <!-- Finalize -->
        <div class="section-block finalize-block">
          <h3 class="section-title">Finalisasi</h3>
          <p class="finalize-desc">Setelah finalisasi, stok opname tidak dapat diubah dan variance report akan dibuat.</p>
          {#if fn_error}<div class="alert-err">{fn_error}</div>{/if}
          <button class="btn-finalize" onclick={finalizeStocktake} disabled={fn_loading}>
            {#if fn_loading}<span class="spinner"></span>{/if}
            Finalisasi Stok Opname
          </button>
        </div>
      {/if}

      {#if fn_variance}
        <!-- Variance report -->
        <div class="section-block">
          <h3 class="section-title">Laporan Variance</h3>
          <div class="variance-summary">
            <div class="var-card var-total">
              <div class="var-label">Total SKU</div>
              <div class="var-val">{fn_variance.total_skus ?? varianceItems.length}</div>
            </div>
            <div class="var-card var-ok">
              <div class="var-label">Cocok</div>
              <div class="var-val">{fn_variance.matched ?? '-'}</div>
            </div>
            <div class="var-card var-diff">
              <div class="var-label">Selisih</div>
              <div class="var-val">{fn_variance.variance_count ?? varianceItems.filter((i: any) => i.variance !== 0).length}</div>
            </div>
          </div>
          {#if varianceItems.length > 0}
            <div class="table-wrap mt">
              <table>
                <thead>
                  <tr><th>SKU</th><th>Nama</th><th class="ar">Sistem</th><th class="ar">Terhitung</th><th class="ar">Variance</th></tr>
                </thead>
                <tbody>
                  {#each varianceItems as item}
                    {@const v = (item.counted_qty ?? 0) - (item.system_qty ?? 0)}
                    <tr class={v !== 0 ? 'row-diff' : ''}>
                      <td class="mono">{item.sku}</td>
                      <td>{item.product_name ?? '-'}</td>
                      <td class="ar">{item.system_qty}</td>
                      <td class="ar">{item.counted_qty}</td>
                      <td class="ar {v > 0 ? 'pos' : v < 0 ? 'neg' : ''}">{v > 0 ? '+' : ''}{v}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>
      {/if}
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 72rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; margin-bottom: 1.25rem; }
  .section-title { margin: 0 0 1rem; font-size: 0.9rem; font-weight: 700; }
  .form-row { display: flex; gap: 0.75rem; align-items: flex-end; flex-wrap: wrap; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { flex: 1; min-width: 180px; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .active-banner { background: #eff6ff; border-color: #bfdbfe; }
  .active-info { display: flex; align-items: center; gap: 0.85rem; }
  .active-icon { font-size: 1.75rem; color: #2563eb; }
  .active-title { font-size: 0.9rem; font-weight: 700; color: #1e40af; margin-bottom: 0.25rem; }
  .active-sub { font-size: 0.78rem; color: #434655; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .chip { display: inline-flex; padding: 0.1rem 0.4rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .chip-blue { background: #dbeafe; color: #1d4ed8; }
  .chip-green { background: #dcfce7; color: #166534; }
  .finalize-block { border-color: #fef08a; background: #fefce8; }
  .finalize-desc { font-size: 0.78rem; color: #737686; margin: 0 0 1rem; }
  .btn-finalize { display: inline-flex; align-items: center; gap: 0.35rem; background: #ca8a04; color: #fff; border: 0; border-radius: 0.25rem; padding: 0.55rem 1.1rem; font-size: 0.85rem; font-weight: 700; cursor: pointer; font-family: inherit; }
  .btn-finalize:disabled { opacity: 0.6; cursor: not-allowed; }
  .table-wrap { overflow-x: auto; }
  .mt { margin-top: 1rem; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .row-diff { background: #fff7ed !important; }
  .ar { text-align: right; }
  .pos { color: #059669; font-weight: 600; }
  .neg { color: #dc2626; font-weight: 600; }
  .variance-summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0.75rem; margin-bottom: 0.75rem; }
  .var-card { padding: 0.85rem 1rem; border-radius: 0.35rem; text-align: center; }
  .var-total { background: #f2f4f6; }
  .var-ok { background: #dcfce7; }
  .var-diff { background: #fef9c3; }
  .var-label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #737686; margin-bottom: 0.35rem; }
  .var-val { font-size: 1.2rem; font-weight: 700; color: #191c1e; }
</style>
