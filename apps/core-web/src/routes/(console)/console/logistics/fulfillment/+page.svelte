<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'size' | 'courier' | 'returns'>('size');

  // Tab 1 — Size Sync
  let sz_departureId = $state('');
  let sz_loading = $state(false);
  let sz_error = $state('');
  let sz_result = $state<any>(null);

  async function syncSizes(e: Event) {
    e.preventDefault();
    sz_loading = true; sz_error = ''; sz_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/size-sync`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ departure_id: sz_departureId }),
      });
      if (!res.ok) throw new Error(`Gagal sinkronisasi ukuran (${res.status})`);
      sz_result = await res.json();
    } catch (err) {
      sz_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    sz_loading = false;
  }

  const sizePresets = $derived(
    sz_result && Array.isArray(sz_result.presets) ? sz_result.presets : []
  );

  // Tab 2 — Courier Tracking
  let ct_taskId = $state('');
  let ct_courier = $state('');
  let ct_tracking = $state('');
  let ct_status = $state('picked_up');
  let ct_note = $state('');
  let ct_loading = $state(false);
  let ct_error = $state('');
  let ct_success = $state('');

  const COURIER_STATUSES = [
    { value: 'picked_up', label: 'Diambil' },
    { value: 'in_transit', label: 'Dalam Perjalanan' },
    { value: 'out_for_delivery', label: 'Dikirim' },
    { value: 'delivered', label: 'Terkirim' },
    { value: 'failed', label: 'Gagal Kirim' },
    { value: 'returned', label: 'Dikembalikan' },
  ];

  async function submitCourier(e: Event) {
    e.preventDefault();
    ct_loading = true; ct_error = ''; ct_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/courier-tracking`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          fulfillment_task_id: ct_taskId,
          courier_name: ct_courier,
          tracking_number: ct_tracking,
          status: ct_status,
          note: ct_note,
        }),
      });
      if (!res.ok) throw new Error(`Gagal update tracking (${res.status})`);
      ct_success = 'Status courier berhasil diperbarui.';
      ct_note = '';
    } catch (err) {
      ct_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ct_loading = false;
  }

  // Tab 3 — Returns & Exchanges
  let rt_bookingId = $state('');
  let rt_sku = $state('');
  let rt_qty = $state('');
  let rt_reason = $state('');
  let rt_loading = $state(false);
  let rt_error = $state('');
  let rt_success = $state('');

  async function submitReturn(e: Event) {
    e.preventDefault();
    rt_loading = true; rt_error = ''; rt_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/returns`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          booking_id: rt_bookingId,
          sku: rt_sku,
          quantity: parseInt(rt_qty),
          reason: rt_reason,
        }),
      });
      if (!res.ok) throw new Error(`Gagal catat retur (${res.status})`);
      rt_success = 'Retur berhasil dicatat.';
      rt_sku = ''; rt_qty = ''; rt_reason = '';
    } catch (err) {
      rt_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    rt_loading = false;
  }

  let ex_returnId = $state('');
  let ex_newSku = $state('');
  let ex_newQty = $state('');
  let ex_loading = $state(false);
  let ex_error = $state('');
  let ex_success = $state('');

  async function submitExchange(e: Event) {
    e.preventDefault();
    ex_loading = true; ex_error = ''; ex_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/exchanges`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          return_id: ex_returnId,
          new_sku: ex_newSku,
          new_quantity: parseInt(ex_newQty),
        }),
      });
      if (!res.ok) throw new Error(`Gagal proses tukar (${res.status})`);
      ex_success = 'Penukaran berhasil diproses.';
      ex_returnId = ''; ex_newSku = ''; ex_newQty = '';
    } catch (err) {
      ex_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ex_loading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/logistics" class="bc-link">Logistik</a>
      <span class="bc-sep">/</span>
      <span>Fulfillment &amp; Pengiriman</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Fulfillment &amp; Pengiriman</h2>
      <p>BL-LOG-025/026/027/028/029 — Sinkronisasi ukuran, courier tracking, retur &amp; tukar</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'size'} onclick={() => tab = 'size'}>Sinkronisasi Ukuran</button>
      <button class="tab-btn" class:active={tab === 'courier'} onclick={() => tab = 'courier'}>Courier Tracking</button>
      <button class="tab-btn" class:active={tab === 'returns'} onclick={() => tab = 'returns'}>Retur &amp; Tukar</button>
    </div>

    {#if tab === 'size'}
      <div class="section-block">
        <h3 class="section-title">Sinkronisasi Ukuran Jamaah</h3>
        <form class="form-row" onsubmit={syncSizes}>
          <div class="field">
            <label for="sz-dep">ID Keberangkatan</label>
            <input id="sz-dep" type="text" placeholder="dep-001" bind:value={sz_departureId} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={sz_loading}>
              {#if sz_loading}<span class="spinner"></span>{/if}
              Sinkronisasi
            </button>
          </div>
        </form>
        {#if sz_error}<div class="alert-err">{sz_error}</div>{/if}
      </div>

      {#if sizePresets.length > 0}
        <div class="section-block">
          <h3 class="section-title">Preset Ukuran Jamaah ({sizePresets.length})</h3>
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Jamaah</th><th>Baju (Pria)</th><th>Gamis (Wanita)</th><th>Sepatu</th><th>Ihram</th></tr>
              </thead>
              <tbody>
                {#each sizePresets as preset}
                  <tr>
                    <td class="mono">{preset.pilgrim_id ?? '-'}</td>
                    <td>{preset.shirt_size ?? '-'}</td>
                    <td>{preset.dress_size ?? '-'}</td>
                    <td>{preset.shoe_size ?? '-'}</td>
                    <td>{preset.ihram_size ?? '-'}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}
    {/if}

    {#if tab === 'courier'}
      <div class="section-block">
        <h3 class="section-title">Update Status Courier</h3>
        <form class="form-grid" onsubmit={submitCourier}>
          <div class="field">
            <label for="ct-task">Fulfillment Task ID</label>
            <input id="ct-task" type="text" placeholder="task-001" bind:value={ct_taskId} required />
          </div>
          <div class="field">
            <label for="ct-courier">Nama Courier</label>
            <input id="ct-courier" type="text" placeholder="JNE / SiCepat / Tiki" bind:value={ct_courier} required />
          </div>
          <div class="field">
            <label for="ct-tracking">Nomor Resi</label>
            <input id="ct-tracking" type="text" placeholder="JNE-123456" bind:value={ct_tracking} required />
          </div>
          <div class="field">
            <label for="ct-status">Status</label>
            <select id="ct-status" bind:value={ct_status}>
              {#each COURIER_STATUSES as s}
                <option value={s.value}>{s.label}</option>
              {/each}
            </select>
          </div>
          <div class="field field-wide">
            <label for="ct-note">Catatan</label>
            <input id="ct-note" type="text" placeholder="Catatan pengiriman..." bind:value={ct_note} />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={ct_loading}>
              {#if ct_loading}<span class="spinner"></span>{/if}
              Update
            </button>
          </div>
        </form>
        {#if ct_error}<div class="alert-err">{ct_error}</div>{/if}
        {#if ct_success}<div class="alert-ok">{ct_success}</div>{/if}
      </div>
    {/if}

    {#if tab === 'returns'}
      <div class="two-col">
        <!-- Return form -->
        <div class="section-block">
          <h3 class="section-title">Catat Retur Barang</h3>
          <form class="form-col" onsubmit={submitReturn}>
            <div class="field">
              <label for="rt-booking">Booking ID</label>
              <input id="rt-booking" type="text" placeholder="BKG-001" bind:value={rt_bookingId} required />
            </div>
            <div class="field">
              <label for="rt-sku">SKU</label>
              <input id="rt-sku" type="text" placeholder="SKU-001" bind:value={rt_sku} required />
            </div>
            <div class="field">
              <label for="rt-qty">Jumlah</label>
              <input id="rt-qty" type="number" min="1" placeholder="1" bind:value={rt_qty} required />
            </div>
            <div class="field">
              <label for="rt-reason">Alasan Retur</label>
              <textarea id="rt-reason" rows="3" placeholder="Ukuran tidak sesuai..." bind:value={rt_reason} required></textarea>
            </div>
            <button type="submit" class="btn-primary" disabled={rt_loading}>
              {#if rt_loading}<span class="spinner"></span>{/if}
              Catat Retur
            </button>
            {#if rt_error}<div class="alert-err">{rt_error}</div>{/if}
            {#if rt_success}<div class="alert-ok">{rt_success}</div>{/if}
          </form>
        </div>

        <!-- Exchange form -->
        <div class="section-block">
          <h3 class="section-title">Proses Tukar Barang</h3>
          <form class="form-col" onsubmit={submitExchange}>
            <div class="field">
              <label for="ex-return">Return ID</label>
              <input id="ex-return" type="text" placeholder="RET-001" bind:value={ex_returnId} required />
            </div>
            <div class="field">
              <label for="ex-sku">SKU Pengganti</label>
              <input id="ex-sku" type="text" placeholder="SKU-002" bind:value={ex_newSku} required />
            </div>
            <div class="field">
              <label for="ex-qty">Jumlah Pengganti</label>
              <input id="ex-qty" type="number" min="1" placeholder="1" bind:value={ex_newQty} required />
            </div>
            <button type="submit" class="btn-primary" disabled={ex_loading}>
              {#if ex_loading}<span class="spinner"></span>{/if}
              Proses Tukar
            </button>
            {#if ex_error}<div class="alert-err">{ex_error}</div>{/if}
            {#if ex_success}<div class="alert-ok">{ex_success}</div>{/if}
          </form>
        </div>
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
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 0.75rem; align-items: end; }
  .form-col { display: flex; flex-direction: column; gap: 0.7rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { grid-column: span 2; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; align-self: flex-start; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; }
  @media (max-width: 800px) { .two-col { grid-template-columns: 1fr; } }
</style>
