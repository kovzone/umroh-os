<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'staff' | 'log'>('staff');

  // Tab 1 — Staff Assignment
  let sa_departureId = $state('');
  let sa_staffUserId = $state('');
  let sa_role = $state('Pembimbing');
  let sa_pilgrimIds = $state('');
  let sa_loading = $state(false);
  let sa_error = $state('');
  let sa_success = $state('');

  const STAFF_ROLES = ['Pembimbing', 'Ketua Regu', 'Petugas Visa'];

  async function submitStaff(e: Event) {
    e.preventDefault();
    sa_loading = true; sa_error = ''; sa_success = '';
    try {
      const ids = sa_pilgrimIds.split('\n').map(s => s.trim()).filter(Boolean);
      const res = await fetch(`${GATEWAY}/v1/ops/staff-assignments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          departure_id: sa_departureId,
          staff_user_id: sa_staffUserId,
          role: sa_role,
          pilgrim_ids: ids,
        }),
      });
      if (!res.ok) throw new Error(`Gagal assign staf (${res.status})`);
      sa_success = 'Penugasan staf berhasil disimpan.';
      sa_staffUserId = ''; sa_pilgrimIds = '';
    } catch (err) {
      sa_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    sa_loading = false;
  }

  // Tab 2 — Passport Log
  let log_departureId = $state('');
  let log_loading = $state(false);
  let log_error = $state('');
  let log_entries = $state<any[]>([]);

  async function loadLog() {
    if (!log_departureId.trim()) return;
    log_loading = true; log_error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/passport-log/${log_departureId}`);
      if (!res.ok) throw new Error(`Gagal memuat log (${res.status})`);
      const body = await res.json();
      log_entries = body.logs ?? body ?? [];
    } catch (err) {
      log_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    log_loading = false;
  }

  let ht_pilgrimId = $state('');
  let ht_fromUser = $state('');
  let ht_toUser = $state('');
  let ht_notes = $state('');
  let ht_loading = $state(false);
  let ht_error = $state('');
  let ht_success = $state('');

  async function submitHandover(e: Event) {
    e.preventDefault();
    ht_loading = true; ht_error = ''; ht_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/passport-log`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          departure_id: log_departureId,
          pilgrim_id: ht_pilgrimId,
          from_user_id: ht_fromUser,
          to_user_id: ht_toUser,
          notes: ht_notes,
        }),
      });
      if (!res.ok) throw new Error(`Gagal catat serah terima (${res.status})`);
      const newEntry = await res.json();
      log_entries = [newEntry, ...log_entries];
      ht_pilgrimId = ''; ht_fromUser = ''; ht_toUser = ''; ht_notes = '';
      ht_success = 'Serah terima paspor berhasil dicatat.';
    } catch (err) {
      ht_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ht_loading = false;
  }

  function formatDate(iso: string) {
    try { return new Date(iso).toLocaleString('id-ID'); } catch { return iso; }
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Penugasan Staf &amp; Log Paspor</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Penugasan Staf &amp; Log Paspor</h2>
      <p>BL-OPS-029/030 — Assign staf ke keberangkatan dan lacak serah terima paspor</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'staff'} onclick={() => tab = 'staff'}>Penugasan Staf</button>
      <button class="tab-btn" class:active={tab === 'log'} onclick={() => tab = 'log'}>Log Serah Terima Paspor</button>
    </div>

    {#if tab === 'staff'}
      <div class="section-block">
        <h3 class="section-title">Assign Staf ke Keberangkatan</h3>
        <form class="form-grid" onsubmit={submitStaff}>
          <div class="field">
            <label for="sa-dep">ID Keberangkatan</label>
            <input id="sa-dep" type="text" placeholder="dep-001" bind:value={sa_departureId} required />
          </div>
          <div class="field">
            <label for="sa-staff">ID Staf</label>
            <input id="sa-staff" type="text" placeholder="user-001" bind:value={sa_staffUserId} required />
          </div>
          <div class="field">
            <label for="sa-role">Peran</label>
            <select id="sa-role" bind:value={sa_role}>
              {#each STAFF_ROLES as r}
                <option>{r}</option>
              {/each}
            </select>
          </div>
          <div class="field field-wide">
            <label for="sa-pilgrims">ID Jamaah (satu per baris, opsional)</label>
            <textarea id="sa-pilgrims" rows="4" placeholder="pilg-001&#10;pilg-002" bind:value={sa_pilgrimIds}></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={sa_loading}>
              {#if sa_loading}<span class="spinner"></span>{/if}
              Simpan Penugasan
            </button>
          </div>
        </form>
        {#if sa_error}<div class="alert-err">{sa_error}</div>{/if}
        {#if sa_success}<div class="alert-ok">{sa_success}</div>{/if}
      </div>
    {/if}

    {#if tab === 'log'}
      <div class="section-block">
        <h3 class="section-title">Filter Keberangkatan</h3>
        <div class="inline-filter">
          <input type="text" placeholder="ID Keberangkatan" bind:value={log_departureId} />
          <button class="btn-primary" onclick={loadLog} disabled={log_loading}>
            {#if log_loading}<span class="spinner"></span>{/if}
            Muat Log
          </button>
        </div>
        {#if log_error}<div class="alert-err">{log_error}</div>{/if}
      </div>

      {#if log_entries.length > 0}
        <div class="section-block">
          <h3 class="section-title">Riwayat Serah Terima Paspor</h3>
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Jamaah</th>
                  <th>Dari</th>
                  <th>Kepada</th>
                  <th>Catatan</th>
                  <th>Waktu</th>
                </tr>
              </thead>
              <tbody>
                {#each log_entries as entry}
                  <tr>
                    <td class="mono">{entry.pilgrim_id ?? '-'}</td>
                    <td class="mono">{entry.from_user_id ?? '-'}</td>
                    <td class="mono">{entry.to_user_id ?? '-'}</td>
                    <td>{entry.notes ?? '-'}</td>
                    <td>{entry.created_at ? formatDate(entry.created_at) : '-'}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}

      <div class="section-block">
        <h3 class="section-title">Catat Serah Terima Paspor</h3>
        <form class="form-grid" onsubmit={submitHandover}>
          <div class="field">
            <label for="ht-pilgrim">ID Jamaah</label>
            <input id="ht-pilgrim" type="text" placeholder="pilg-001" bind:value={ht_pilgrimId} required />
          </div>
          <div class="field">
            <label for="ht-from">Dari (User ID)</label>
            <input id="ht-from" type="text" placeholder="user-001" bind:value={ht_fromUser} required />
          </div>
          <div class="field">
            <label for="ht-to">Kepada (User ID)</label>
            <input id="ht-to" type="text" placeholder="user-002" bind:value={ht_toUser} required />
          </div>
          <div class="field field-wide">
            <label for="ht-notes">Catatan</label>
            <textarea id="ht-notes" rows="2" placeholder="Catatan serah terima..." bind:value={ht_notes}></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={ht_loading}>
              {#if ht_loading}<span class="spinner"></span>{/if}
              Catat Serah Terima
            </button>
          </div>
        </form>
        {#if ht_error}<div class="alert-err">{ht_error}</div>{/if}
        {#if ht_success}<div class="alert-ok">{ht_success}</div>{/if}
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
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { grid-column: span 2; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
</style>
