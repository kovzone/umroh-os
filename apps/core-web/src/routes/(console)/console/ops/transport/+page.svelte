<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'rooming' | 'transport' | 'delta'>('rooming');

  // Tab 1 — Smart Rooming
  let rm_departureId = $state('');
  let rm_maxPerRoom = $state(4);
  let rm_pisahMahram = $state(false);
  let rm_prioritasKeluarga = $state(false);
  let rm_loading = $state(false);
  let rm_error = $state('');
  let rm_result = $state<any>(null);

  async function submitRooming(e: Event) {
    e.preventDefault();
    rm_loading = true; rm_error = ''; rm_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/smart-rooming`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          departure_id: rm_departureId,
          max_per_room: rm_maxPerRoom,
          separate_mahram: rm_pisahMahram,
          prioritize_family: rm_prioritasKeluarga,
        }),
      });
      if (!res.ok) throw new Error(`Gagal penempatan kamar (${res.status})`);
      rm_result = await res.json();
    } catch (err) {
      rm_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    rm_loading = false;
  }

  // Tab 2 — Transport
  let tr_departureId = $state('');
  let tr_vehicleType = $state('Bus');
  let tr_vehicleId = $state('');
  let tr_pilgrimIds = $state('');
  let tr_loading = $state(false);
  let tr_error = $state('');
  let tr_assignments = $state<any[]>([]);

  async function submitTransport(e: Event) {
    e.preventDefault();
    tr_loading = true; tr_error = '';
    try {
      const ids = tr_pilgrimIds.split('\n').map(s => s.trim()).filter(Boolean);
      const res = await fetch(`${GATEWAY}/v1/ops/transport-assignments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          departure_id: tr_departureId,
          vehicle_type: tr_vehicleType,
          vehicle_id: tr_vehicleId,
          pilgrim_ids: ids,
        }),
      });
      if (!res.ok) throw new Error(`Gagal assign transportasi (${res.status})`);
      await loadTransportAssignments();
      tr_vehicleId = ''; tr_pilgrimIds = '';
    } catch (err) {
      tr_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    tr_loading = false;
  }

  async function loadTransportAssignments() {
    if (!tr_departureId.trim()) return;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/transport-assignments/${tr_departureId}`);
      if (res.ok) {
        const body = await res.json();
        tr_assignments = body.assignments ?? body ?? [];
      }
    } catch { /* ignore */ }
  }

  // Tab 3 — Manifest Delta
  let md_departureId = $state('');
  let md_changeType = $state('add');
  let md_entityId = $state('');
  let md_notes = $state('');
  let md_loading = $state(false);
  let md_error = $state('');
  let md_success = $state('');

  const CHANGE_TYPES = [
    { value: 'add', label: 'Tambah' },
    { value: 'remove', label: 'Hapus' },
    { value: 'update', label: 'Perbarui' },
    { value: 'swap', label: 'Tukar' },
  ];

  async function submitDelta(e: Event) {
    e.preventDefault();
    md_loading = true; md_error = ''; md_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/manifest-delta`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          departure_id: md_departureId,
          change_type: md_changeType,
          entity_id: md_entityId,
          notes: md_notes,
        }),
      });
      if (!res.ok) throw new Error(`Gagal catat delta (${res.status})`);
      md_success = 'Perubahan manifest berhasil dicatat.';
      md_entityId = ''; md_notes = '';
    } catch (err) {
      md_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    md_loading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Akomodasi &amp; Transportasi</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Akomodasi &amp; Transportasi</h2>
      <p>BL-OPS-026/027/028 — Penempatan kamar, transportasi, dan delta manifest</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'rooming'} onclick={() => tab = 'rooming'}>Penempatan Kamar</button>
      <button class="tab-btn" class:active={tab === 'transport'} onclick={() => tab = 'transport'}>Transportasi</button>
      <button class="tab-btn" class:active={tab === 'delta'} onclick={() => tab = 'delta'}>Delta Manifest</button>
    </div>

    {#if tab === 'rooming'}
      <div class="section-block">
        <h3 class="section-title">Smart Rooming</h3>
        <form class="form-grid" onsubmit={submitRooming}>
          <div class="field">
            <label for="rm-dep">ID Keberangkatan</label>
            <input id="rm-dep" type="text" placeholder="dep-001" bind:value={rm_departureId} required />
          </div>
          <div class="field">
            <label for="rm-max">Maks per Kamar</label>
            <input id="rm-max" type="number" min="1" max="10" bind:value={rm_maxPerRoom} />
          </div>
          <div class="field checkbox-field">
            <label class="checkbox-label">
              <input type="checkbox" bind:checked={rm_pisahMahram} />
              Pisah Mahram
            </label>
          </div>
          <div class="field checkbox-field">
            <label class="checkbox-label">
              <input type="checkbox" bind:checked={rm_prioritasKeluarga} />
              Prioritas Keluarga
            </label>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={rm_loading}>
              {#if rm_loading}<span class="spinner"></span>{/if}
              Proses Penempatan
            </button>
          </div>
        </form>
        {#if rm_error}
          <div class="alert-err">{rm_error}</div>
        {/if}
        {#if rm_result}
          <div class="result-card">
            <span class="material-symbols-outlined result-icon">check_circle</span>
            <div>
              <div class="result-title">Penempatan selesai</div>
              <div class="result-sub">{rm_result.rooms_allocated ?? rm_result.total_rooms ?? 0} kamar dialokasikan</div>
            </div>
          </div>
        {/if}
      </div>
    {/if}

    {#if tab === 'transport'}
      <div class="section-block">
        <h3 class="section-title">Assign Transportasi</h3>
        <form class="form-grid" onsubmit={submitTransport}>
          <div class="field">
            <label for="tr-dep">ID Keberangkatan</label>
            <input id="tr-dep" type="text" placeholder="dep-001" bind:value={tr_departureId}
              oninput={loadTransportAssignments} required />
          </div>
          <div class="field">
            <label for="tr-vtype">Jenis Kendaraan</label>
            <select id="tr-vtype" bind:value={tr_vehicleType}>
              <option>Bus</option>
              <option>Van</option>
              <option>Hiace</option>
            </select>
          </div>
          <div class="field">
            <label for="tr-vid">ID Kendaraan</label>
            <input id="tr-vid" type="text" placeholder="BUS-001" bind:value={tr_vehicleId} required />
          </div>
          <div class="field field-wide">
            <label for="tr-pilgrims">ID Jamaah (satu per baris)</label>
            <textarea id="tr-pilgrims" rows="4" placeholder="pilg-001&#10;pilg-002&#10;pilg-003" bind:value={tr_pilgrimIds}></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={tr_loading}>
              {#if tr_loading}<span class="spinner"></span>{/if}
              Assign Transportasi
            </button>
          </div>
        </form>
        {#if tr_error}
          <div class="alert-err">{tr_error}</div>
        {/if}
      </div>

      {#if tr_assignments.length > 0}
        <div class="section-block">
          <h3 class="section-title">Penugasan Transportasi Aktif</h3>
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>ID Jamaah</th>
                  <th>Kendaraan</th>
                  <th>Jenis</th>
                  <th>Keberangkatan</th>
                </tr>
              </thead>
              <tbody>
                {#each tr_assignments as a}
                  <tr>
                    <td class="mono">{a.pilgrim_id ?? '-'}</td>
                    <td class="mono">{a.vehicle_id ?? '-'}</td>
                    <td><span class="chip chip-blue">{a.vehicle_type ?? '-'}</span></td>
                    <td class="mono">{a.departure_id ?? '-'}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}
    {/if}

    {#if tab === 'delta'}
      <div class="section-block">
        <h3 class="section-title">Catat Delta Manifest</h3>
        <form class="form-grid" onsubmit={submitDelta}>
          <div class="field">
            <label for="md-dep">ID Keberangkatan</label>
            <input id="md-dep" type="text" placeholder="dep-001" bind:value={md_departureId} required />
          </div>
          <div class="field">
            <label for="md-type">Jenis Perubahan</label>
            <select id="md-type" bind:value={md_changeType}>
              {#each CHANGE_TYPES as ct}
                <option value={ct.value}>{ct.label}</option>
              {/each}
            </select>
          </div>
          <div class="field">
            <label for="md-entity">ID Entitas</label>
            <input id="md-entity" type="text" placeholder="pilg-001" bind:value={md_entityId} required />
          </div>
          <div class="field field-wide">
            <label for="md-notes">Catatan</label>
            <textarea id="md-notes" rows="3" placeholder="Alasan perubahan..." bind:value={md_notes}></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={md_loading}>
              {#if md_loading}<span class="spinner"></span>{/if}
              Simpan Delta
            </button>
          </div>
        </form>
        {#if md_error}
          <div class="alert-err">{md_error}</div>
        {/if}
        {#if md_success}
          <div class="alert-ok">{md_success}</div>
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
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { grid-column: span 2; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .checkbox-field { justify-content: flex-end; }
  .checkbox-label { display: flex; align-items: center; gap: 0.45rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; }
  .checkbox-label input { width: 1rem; height: 1rem; cursor: pointer; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .result-card { display: flex; align-items: center; gap: 0.85rem; padding: 0.85rem 1rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 0.35rem; margin-top: 0.85rem; }
  .result-icon { font-size: 1.5rem; color: #16a34a; }
  .result-title { font-size: 0.82rem; font-weight: 700; color: #166534; }
  .result-sub { font-size: 0.72rem; color: #737686; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .chip { display: inline-flex; padding: 0.12rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .chip-blue { background: #e0f2fe; color: #075985; }
</style>
