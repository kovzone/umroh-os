<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'partial' | 'reversal'>('partial');

  // Tab 1 — Partial GRN
  let pg_prId = $state('');
  let pg_notes = $state('');
  let pg_items = $state([{ sku: '', qty: '', unit_cost: '' }]);
  let pg_loading = $state(false);
  let pg_error = $state('');
  let pg_result = $state<any>(null);

  function addItem() {
    pg_items = [...pg_items, { sku: '', qty: '', unit_cost: '' }];
  }

  function removeItem(i: number) {
    pg_items = pg_items.filter((_, idx) => idx !== i);
  }

  function updateItem(i: number, field: string, value: string) {
    pg_items = pg_items.map((item, idx) =>
      idx === i ? { ...item, [field]: value } : item
    );
  }

  async function submitGRN(e: Event) {
    e.preventDefault();
    pg_loading = true; pg_error = ''; pg_result = null;
    try {
      const items = pg_items.map(i => ({
        sku: i.sku,
        quantity: parseFloat(i.qty),
        unit_cost: parseFloat(i.unit_cost),
      }));
      const res = await fetch(`${GATEWAY}/v1/logistics/partial-grn`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pr_id: pg_prId, notes: pg_notes, items }),
      });
      if (!res.ok) throw new Error(`Gagal buat GRN (${res.status})`);
      pg_result = await res.json();
      pg_prId = ''; pg_notes = '';
      pg_items = [{ sku: '', qty: '', unit_cost: '' }];
    } catch (err) {
      pg_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    pg_loading = false;
  }

  // Tab 2 — Reversal
  let rv_grnId = $state('');
  let rv_reason = $state('');
  let rv_loading = $state(false);
  let rv_error = $state('');
  let rv_success = $state('');

  async function submitReversal(e: Event) {
    e.preventDefault();
    rv_loading = true; rv_error = ''; rv_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/reverse-grn`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ grn_id: rv_grnId, reason: rv_reason }),
      });
      if (!res.ok) throw new Error(`Gagal reversal GRN (${res.status})`);
      rv_success = 'GRN berhasil di-reverse.';
      rv_grnId = ''; rv_reason = '';
    } catch (err) {
      rv_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    rv_loading = false;
  }

  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/logistics" class="bc-link">Logistik</a>
      <span class="bc-sep">/</span>
      <span>Penerimaan Barang (GRN)</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Penerimaan Barang (GRN)</h2>
      <p>BL-LOG-017/018 — GRN parsial dan reversal</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'partial'} onclick={() => tab = 'partial'}>GRN Parsial</button>
      <button class="tab-btn" class:active={tab === 'reversal'} onclick={() => tab = 'reversal'}>Reversal GRN</button>
    </div>

    {#if tab === 'partial'}
      <div class="section-block">
        <h3 class="section-title">Buat GRN Parsial</h3>
        <form onsubmit={submitGRN}>
          <div class="form-row mb">
            <div class="field">
              <label for="pg-pr">PR ID</label>
              <input id="pg-pr" type="text" placeholder="PR-001" bind:value={pg_prId} required />
            </div>
            <div class="field field-wide">
              <label for="pg-notes">Catatan</label>
              <input id="pg-notes" type="text" placeholder="Catatan penerimaan..." bind:value={pg_notes} />
            </div>
          </div>

          <div class="items-header">
            <span class="items-label">Item yang Diterima</span>
            <button type="button" class="btn-add" onclick={addItem}>
              <span class="material-symbols-outlined">add</span>
              Tambah Item
            </button>
          </div>

          <div class="items-table">
            <div class="items-head">
              <span>SKU</span>
              <span>Jumlah</span>
              <span>Harga Satuan (Rp)</span>
              <span></span>
            </div>
            {#each pg_items as item, i}
              <div class="items-row">
                <input type="text" placeholder="SKU-001" value={item.sku}
                  oninput={(e) => updateItem(i, 'sku', (e.target as HTMLInputElement).value)} required />
                <input type="number" min="0" step="0.01" placeholder="10" value={item.qty}
                  oninput={(e) => updateItem(i, 'qty', (e.target as HTMLInputElement).value)} required />
                <input type="number" min="0" placeholder="50000" value={item.unit_cost}
                  oninput={(e) => updateItem(i, 'unit_cost', (e.target as HTMLInputElement).value)} />
                <button type="button" class="btn-remove" onclick={() => removeItem(i)}
                  disabled={pg_items.length === 1}>
                  <span class="material-symbols-outlined">delete</span>
                </button>
              </div>
            {/each}
          </div>

          <div class="form-footer">
            <button type="submit" class="btn-primary" disabled={pg_loading}>
              {#if pg_loading}<span class="spinner"></span>{/if}
              Buat GRN
            </button>
          </div>
        </form>
        {#if pg_error}<div class="alert-err">{pg_error}</div>{/if}
      </div>

      {#if pg_result}
        <div class="section-block">
          <h3 class="section-title">GRN Dibuat</h3>
          <div class="result-info">
            <div class="result-row"><span>GRN ID</span><strong class="mono">{pg_result.grn_id ?? pg_result.id ?? '-'}</strong></div>
            <div class="result-row"><span>PR ID</span><strong class="mono">{pg_result.pr_id ?? '-'}</strong></div>
            <div class="result-row"><span>Total Item</span><strong>{pg_result.items_count ?? (pg_result.items ?? []).length ?? '-'}</strong></div>
          </div>
          {#if pg_result.items?.length > 0}
            <div class="table-wrap mt">
              <table>
                <thead>
                  <tr><th>SKU</th><th class="ar">Qty Diterima</th><th class="ar">Harga Satuan</th><th class="ar">Total</th></tr>
                </thead>
                <tbody>
                  {#each pg_result.items as item}
                    <tr>
                      <td class="mono">{item.sku}</td>
                      <td class="ar">{item.quantity_received ?? item.qty}</td>
                      <td class="ar">{formatIDR(item.unit_cost ?? 0)}</td>
                      <td class="ar">{formatIDR((item.quantity_received ?? item.qty ?? 0) * (item.unit_cost ?? 0))}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>
      {/if}
    {/if}

    {#if tab === 'reversal'}
      <div class="section-block">
        <h3 class="section-title">Reversal GRN</h3>
        <form class="form-grid" onsubmit={submitReversal}>
          <div class="field">
            <label for="rv-grn">GRN ID</label>
            <input id="rv-grn" type="text" placeholder="GRN-001" bind:value={rv_grnId} required />
          </div>
          <div class="field field-wide">
            <label for="rv-reason">Alasan Reversal</label>
            <textarea id="rv-reason" rows="3" placeholder="Barang tidak sesuai spesifikasi..." bind:value={rv_reason} required></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-danger" disabled={rv_loading}>
              {#if rv_loading}<span class="spinner"></span>{/if}
              Reverse GRN
            </button>
          </div>
        </form>
        {#if rv_error}<div class="alert-err">{rv_error}</div>{/if}
        {#if rv_success}<div class="alert-ok">{rv_success}</div>{/if}
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
  .form-row { display: flex; gap: 0.75rem; align-items: flex-end; flex-wrap: wrap; }
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; align-items: end; }
  .mb { margin-bottom: 1rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { flex: 1; min-width: 200px; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .field-actions { align-self: flex-end; }
  .items-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.65rem; }
  .items-label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .btn-add { display: inline-flex; align-items: center; gap: 0.2rem; padding: 0.3rem 0.65rem; border: 1px solid #2563eb; background: #eff6ff; color: #2563eb; border-radius: 0.25rem; font-size: 0.75rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-add .material-symbols-outlined { font-size: 0.9rem; }
  .items-table { border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.35rem; overflow: hidden; margin-bottom: 1rem; }
  .items-head { display: grid; grid-template-columns: 1fr 1fr 1fr auto; gap: 0; background: #f2f4f6; padding: 0.45rem 0.75rem; font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .items-row { display: grid; grid-template-columns: 1fr 1fr 1fr auto; gap: 0; align-items: center; border-top: 1px solid rgb(195 198 215 / 0.35); padding: 0.35rem 0.5rem; }
  .items-row input { border: 1px solid rgb(195 198 215 / 0.35); border-radius: 0.2rem; padding: 0.35rem 0.5rem; font-size: 0.8rem; width: 100%; font-family: inherit; background: #fff; margin: 0.2rem; }
  .btn-remove { border: 0; background: transparent; color: #dc2626; cursor: pointer; padding: 0.25rem; display: grid; place-items: center; border-radius: 0.2rem; margin: 0.2rem; }
  .btn-remove:hover { background: #fee2e2; }
  .btn-remove:disabled { opacity: 0.3; cursor: not-allowed; }
  .btn-remove .material-symbols-outlined { font-size: 1rem; }
  .form-footer { display: flex; justify-content: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-danger { display: inline-flex; align-items: center; gap: 0.35rem; background: #dc2626; color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-danger:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .result-info { display: flex; flex-direction: column; gap: 0.5rem; }
  .result-row { display: flex; justify-content: space-between; font-size: 0.82rem; padding: 0.45rem 0; border-bottom: 1px solid rgb(195 198 215 / 0.3); }
  .result-row span { color: #737686; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .table-wrap { overflow-x: auto; }
  .mt { margin-top: 1rem; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .ar { text-align: right; }
  tbody tr:hover { background: #f7f9fb; }
</style>
