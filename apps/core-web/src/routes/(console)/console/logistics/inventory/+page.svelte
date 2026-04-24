<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'stock' | 'warehouse' | 'barcode'>('stock');

  // Tab 1 — Stock & Alerts
  let sa_warehouseId = $state('');
  let sa_loading = $state(false);
  let sa_error = $state('');
  let sa_alerts = $state<any[]>([]);

  async function loadAlerts() {
    sa_loading = true; sa_error = '';
    try {
      const params = sa_warehouseId.trim() ? `?warehouse_id=${sa_warehouseId}` : '';
      const res = await fetch(`${GATEWAY}/v1/logistics/stock-alerts${params}`);
      if (!res.ok) throw new Error(`Gagal memuat stock alerts (${res.status})`);
      const body = await res.json();
      sa_alerts = body.alerts ?? body ?? [];
    } catch (err) {
      sa_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    sa_loading = false;
  }

  $effect(() => { loadAlerts(); });

  // Reorder levels
  let rl_sku = $state('');
  let rl_level = $state('');
  let rl_loading = $state(false);
  let rl_error = $state('');
  let rl_success = $state('');

  async function saveReorderLevel(e: Event) {
    e.preventDefault();
    rl_loading = true; rl_error = ''; rl_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/reorder-levels`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sku: rl_sku, reorder_level: parseInt(rl_level) }),
      });
      if (!res.ok) throw new Error(`Gagal simpan (${res.status})`);
      rl_success = 'Reorder level berhasil diperbarui.';
      rl_sku = ''; rl_level = '';
    } catch (err) {
      rl_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    rl_loading = false;
  }

  // Tab 2 — Multi-Warehouse
  let wh_name = $state('');
  let wh_location = $state('');
  let wh_loading = $state(false);
  let wh_error = $state('');
  let wh_success = $state('');

  async function createWarehouse(e: Event) {
    e.preventDefault();
    wh_loading = true; wh_error = ''; wh_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/warehouses`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: wh_name, location: wh_location }),
      });
      if (!res.ok) throw new Error(`Gagal buat gudang (${res.status})`);
      wh_success = 'Gudang berhasil dibuat.';
      wh_name = ''; wh_location = '';
    } catch (err) {
      wh_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    wh_loading = false;
  }

  let st_fromWarehouse = $state('');
  let st_toWarehouse = $state('');
  let st_sku = $state('');
  let st_qty = $state('');
  let st_loading = $state(false);
  let st_error = $state('');
  let st_success = $state('');

  async function transferStock(e: Event) {
    e.preventDefault();
    st_loading = true; st_error = ''; st_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/stock-transfer`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          from_warehouse_id: st_fromWarehouse,
          to_warehouse_id: st_toWarehouse,
          sku: st_sku,
          quantity: parseInt(st_qty),
        }),
      });
      if (!res.ok) throw new Error(`Gagal transfer stok (${res.status})`);
      st_success = 'Stok berhasil ditransfer.';
      st_sku = ''; st_qty = '';
    } catch (err) {
      st_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    st_loading = false;
  }

  // Tab 3 — Barcode & Labels
  let bc_sku = $state('');
  let bc_quantity = $state('1');
  let bc_format = $state('QR');
  let bc_loading = $state(false);
  let bc_error = $state('');
  let bc_result = $state<any>(null);

  async function generateBarcode(e: Event) {
    e.preventDefault();
    bc_loading = true; bc_error = ''; bc_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/barcode`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sku: bc_sku, quantity: parseInt(bc_quantity), format: bc_format }),
      });
      if (!res.ok) throw new Error(`Gagal generate barcode (${res.status})`);
      bc_result = await res.json();
    } catch (err) {
      bc_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    bc_loading = false;
  }

  let lb_sku = $state('');
  let lb_copies = $state('1');
  let lb_loading = $state(false);
  let lb_error = $state('');
  let lb_result = $state<any>(null);

  async function printLabels(e: Event) {
    e.preventDefault();
    lb_loading = true; lb_error = ''; lb_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/sku-labels`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sku: lb_sku, copies: parseInt(lb_copies) }),
      });
      if (!res.ok) throw new Error(`Gagal cetak label (${res.status})`);
      lb_result = await res.json();
    } catch (err) {
      lb_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    lb_loading = false;
  }

  const ALERT_LEVEL: Record<string, { cls: string; label: string }> = {
    critical: { cls: 'alert-critical', label: 'Kritis' },
    warning: { cls: 'alert-warning', label: 'Peringatan' },
    info: { cls: 'alert-info', label: 'Info' },
  };
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/logistics" class="bc-link">Logistik</a>
      <span class="bc-sep">/</span>
      <span>Inventori &amp; Gudang</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Inventori &amp; Gudang</h2>
      <p>BL-LOG-020/021/022 — Stock alerts, multi-gudang, dan barcode/label</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'stock'} onclick={() => tab = 'stock'}>Stok &amp; Alert</button>
      <button class="tab-btn" class:active={tab === 'warehouse'} onclick={() => tab = 'warehouse'}>Multi-Gudang</button>
      <button class="tab-btn" class:active={tab === 'barcode'} onclick={() => tab = 'barcode'}>Barcode &amp; Label</button>
    </div>

    {#if tab === 'stock'}
      <div class="section-block">
        <h3 class="section-title">Filter Gudang</h3>
        <div class="inline-filter">
          <input type="text" placeholder="ID Gudang (opsional)" bind:value={sa_warehouseId}
            oninput={loadAlerts} />
          <button class="btn-primary" onclick={loadAlerts} disabled={sa_loading}>
            {#if sa_loading}<span class="spinner"></span>{/if}
            Muat Alert
          </button>
        </div>
        {#if sa_error}<div class="alert-err">{sa_error}</div>{/if}
      </div>

      {#if sa_alerts.length > 0}
        <div class="section-block">
          <h3 class="section-title">Stock Alerts ({sa_alerts.length})</h3>
          <div class="alert-cards">
            {#each sa_alerts as alert}
              {@const lvl = ALERT_LEVEL[alert.level ?? alert.severity ?? 'info'] ?? ALERT_LEVEL.info}
              <div class="stock-alert-card {lvl.cls}">
                <div class="stock-alert-top">
                  <span class="stock-sku">{alert.sku ?? '-'}</span>
                  <span class="stock-level-badge">{lvl.label}</span>
                </div>
                <div class="stock-alert-name">{alert.product_name ?? alert.name ?? '-'}</div>
                <div class="stock-alert-qty">
                  Stok: <strong>{alert.current_qty ?? alert.stock ?? 0}</strong>
                  &nbsp;/&nbsp;Min: <strong>{alert.reorder_level ?? alert.min_qty ?? 0}</strong>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else if !sa_loading}
        <div class="empty-state-card">Tidak ada stock alerts aktif.</div>
      {/if}

      <div class="section-block">
        <h3 class="section-title">Set Reorder Level</h3>
        <form class="form-row" onsubmit={saveReorderLevel}>
          <div class="field">
            <label for="rl-sku">SKU</label>
            <input id="rl-sku" type="text" placeholder="SKU-001" bind:value={rl_sku} required />
          </div>
          <div class="field">
            <label for="rl-level">Reorder Level</label>
            <input id="rl-level" type="number" min="0" placeholder="10" bind:value={rl_level} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={rl_loading}>
              {#if rl_loading}<span class="spinner"></span>{/if}
              Simpan
            </button>
          </div>
        </form>
        {#if rl_error}<div class="alert-err">{rl_error}</div>{/if}
        {#if rl_success}<div class="alert-ok">{rl_success}</div>{/if}
      </div>
    {/if}

    {#if tab === 'warehouse'}
      <div class="section-block">
        <h3 class="section-title">Buat Gudang Baru</h3>
        <form class="form-row" onsubmit={createWarehouse}>
          <div class="field">
            <label for="wh-name">Nama Gudang</label>
            <input id="wh-name" type="text" placeholder="Gudang Jakarta Selatan" bind:value={wh_name} required />
          </div>
          <div class="field">
            <label for="wh-loc">Lokasi</label>
            <input id="wh-loc" type="text" placeholder="Jl. Sudirman No.1, Jakarta" bind:value={wh_location} />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={wh_loading}>
              {#if wh_loading}<span class="spinner"></span>{/if}
              Buat Gudang
            </button>
          </div>
        </form>
        {#if wh_error}<div class="alert-err">{wh_error}</div>{/if}
        {#if wh_success}<div class="alert-ok">{wh_success}</div>{/if}
      </div>

      <div class="section-block">
        <h3 class="section-title">Transfer Stok Antar Gudang</h3>
        <form class="form-grid" onsubmit={transferStock}>
          <div class="field">
            <label for="st-from">Dari Gudang</label>
            <input id="st-from" type="text" placeholder="wh-001" bind:value={st_fromWarehouse} required />
          </div>
          <div class="field">
            <label for="st-to">Ke Gudang</label>
            <input id="st-to" type="text" placeholder="wh-002" bind:value={st_toWarehouse} required />
          </div>
          <div class="field">
            <label for="st-sku">SKU</label>
            <input id="st-sku" type="text" placeholder="SKU-001" bind:value={st_sku} required />
          </div>
          <div class="field">
            <label for="st-qty">Jumlah</label>
            <input id="st-qty" type="number" min="1" placeholder="50" bind:value={st_qty} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={st_loading}>
              {#if st_loading}<span class="spinner"></span>{/if}
              Transfer
            </button>
          </div>
        </form>
        {#if st_error}<div class="alert-err">{st_error}</div>{/if}
        {#if st_success}<div class="alert-ok">{st_success}</div>{/if}
      </div>
    {/if}

    {#if tab === 'barcode'}
      <div class="section-block">
        <h3 class="section-title">Generate Barcode</h3>
        <form class="form-row" onsubmit={generateBarcode}>
          <div class="field">
            <label for="bc-sku">SKU</label>
            <input id="bc-sku" type="text" placeholder="SKU-001" bind:value={bc_sku} required />
          </div>
          <div class="field">
            <label for="bc-format">Format</label>
            <select id="bc-format" bind:value={bc_format}>
              <option>QR</option>
              <option>Code128</option>
              <option>EAN13</option>
            </select>
          </div>
          <div class="field">
            <label for="bc-qty">Jumlah</label>
            <input id="bc-qty" type="number" min="1" placeholder="1" bind:value={bc_quantity} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={bc_loading}>
              {#if bc_loading}<span class="spinner"></span>{/if}
              Generate
            </button>
          </div>
        </form>
        {#if bc_error}<div class="alert-err">{bc_error}</div>{/if}
        {#if bc_result}
          <div class="result-card">
            <span class="material-symbols-outlined result-icon">qr_code</span>
            <div class="result-info">
              <div class="result-title">Barcode dibuat</div>
              {#if bc_result.barcode_url ?? bc_result.url}
                <a href={bc_result.barcode_url ?? bc_result.url} target="_blank" class="btn-download">
                  <span class="material-symbols-outlined">download</span>Unduh
                </a>
              {/if}
            </div>
          </div>
        {/if}
      </div>

      <div class="section-block">
        <h3 class="section-title">Cetak Label SKU</h3>
        <form class="form-row" onsubmit={printLabels}>
          <div class="field">
            <label for="lb-sku">SKU</label>
            <input id="lb-sku" type="text" placeholder="SKU-001" bind:value={lb_sku} required />
          </div>
          <div class="field">
            <label for="lb-copies">Jumlah Salinan</label>
            <input id="lb-copies" type="number" min="1" placeholder="1" bind:value={lb_copies} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={lb_loading}>
              {#if lb_loading}<span class="spinner"></span>{/if}
              Cetak Label
            </button>
          </div>
        </form>
        {#if lb_error}<div class="alert-err">{lb_error}</div>{/if}
        {#if lb_result}
          <div class="result-card">
            <span class="material-symbols-outlined result-icon">print</span>
            <div class="result-info">
              <div class="result-title">Label siap dicetak</div>
              {#if lb_result.label_url ?? lb_result.url}
                <a href={lb_result.label_url ?? lb_result.url} target="_blank" class="btn-download">
                  <span class="material-symbols-outlined">download</span>Unduh PDF
                </a>
              {/if}
            </div>
          </div>
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
  .canvas { padding: 1.5rem; max-width: 72rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .tab-bar { display: flex; gap: 0.35rem; margin-bottom: 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); }
  .tab-btn { border: 0; background: transparent; padding: 0.55rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; color: #737686; border-bottom: 2px solid transparent; margin-bottom: -1px; font-family: inherit; }
  .tab-btn.active { color: #2563eb; border-bottom-color: #2563eb; }
  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; margin-bottom: 1.25rem; }
  .section-title { margin: 0 0 1rem; font-size: 0.9rem; font-weight: 700; }
  .inline-filter { display: flex; gap: 0.65rem; align-items: center; }
  .inline-filter input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; min-width: 200px; font-family: inherit; }
  .form-row { display: flex; gap: 0.75rem; align-items: flex-end; flex-wrap: wrap; }
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .alert-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 0.75rem; }
  .stock-alert-card { border-radius: 0.4rem; padding: 0.85rem 1rem; border: 1px solid; }
  .alert-critical { background: #fff1f2; border-color: #fecdd3; }
  .alert-warning { background: #fffbeb; border-color: #fed7aa; }
  .alert-info { background: #eff6ff; border-color: #bfdbfe; }
  .stock-alert-top { display: flex; justify-content: space-between; margin-bottom: 0.25rem; }
  .stock-sku { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; font-weight: 700; }
  .stock-level-badge { font-size: 0.62rem; font-weight: 700; padding: 0.1rem 0.4rem; border-radius: 0.2rem; background: rgba(0,0,0,0.07); }
  .stock-alert-name { font-size: 0.78rem; font-weight: 600; margin-bottom: 0.35rem; }
  .stock-alert-qty { font-size: 0.72rem; color: #434655; }
  .empty-state-card { text-align: center; color: #b0b3c1; padding: 2rem; font-size: 0.82rem; background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; margin-bottom: 1.25rem; }
  .result-card { display: flex; align-items: center; gap: 0.85rem; padding: 0.85rem 1rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 0.35rem; margin-top: 0.85rem; }
  .result-icon { font-size: 1.5rem; color: #16a34a; }
  .result-info { flex: 1; display: flex; align-items: center; gap: 0.75rem; }
  .result-title { font-size: 0.82rem; font-weight: 700; color: #166534; }
  .btn-download { display: inline-flex; align-items: center; gap: 0.3rem; padding: 0.3rem 0.65rem; background: #2563eb; color: #fff; border-radius: 0.2rem; text-decoration: none; font-size: 0.72rem; font-weight: 600; }
  .btn-download .material-symbols-outlined { font-size: 0.85rem; }
</style>
