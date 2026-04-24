<script lang="ts">
  type DeviceType = 'handphone' | 'walkie_talkie' | 'tablet' | 'power_bank';
  type DeviceStatus = 'available' | 'assigned' | 'maintenance' | 'lost';

  interface Device {
    id: string;
    imei: string;
    type: DeviceType;
    label: string;
    assignedTo: string | null;
    assignedDate: string | null;
    status: DeviceStatus;
  }

  let devices = $state<Device[]>([
    { id: 'd1', imei: '356938035643809', type: 'handphone', label: 'HP Samsung A54 #1', assignedTo: 'Ustaz Ahmad Fauzi', assignedDate: '15 Jan 2025', status: 'assigned' },
    { id: 'd2', imei: '490154203237518', type: 'walkie_talkie', label: 'Walkie-Talkie #1', assignedTo: 'Petugas Koper Grup A', assignedDate: '15 Jan 2025', status: 'assigned' },
    { id: 'd3', imei: '358240051111110', type: 'tablet', label: 'Tablet iPad Mini #1', assignedTo: null, assignedDate: null, status: 'available' },
    { id: 'd4', imei: '012345678901234', type: 'power_bank', label: 'Power Bank 20K #1', assignedTo: null, assignedDate: null, status: 'available' },
    { id: 'd5', imei: '999000000000001', type: 'handphone', label: 'HP Xiaomi Note 12 #2', assignedTo: 'Petugas Medis', assignedDate: '15 Jan 2025', status: 'assigned' },
    { id: 'd6', imei: '111222333444555', type: 'walkie_talkie', label: 'Walkie-Talkie #2', assignedTo: null, assignedDate: null, status: 'maintenance' },
  ]);

  let assignForm = $state<{ deviceId: string; assignTo: string }>({ deviceId: '', assignTo: '' });
  let formVisible = $state(false);
  let searchQuery = $state('');

  const filtered = $derived(
    devices.filter(d =>
      !searchQuery ||
      d.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
      d.imei.includes(searchQuery) ||
      (d.assignedTo ?? '').toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  const stats = $derived({
    total: devices.length,
    assigned: devices.filter(d => d.status === 'assigned').length,
    available: devices.filter(d => d.status === 'available').length,
    maintenance: devices.filter(d => d.status === 'maintenance').length,
  });

  function handleAssign(e: Event) {
    e.preventDefault();
    if (!assignForm.deviceId || !assignForm.assignTo.trim()) return;
    const now = new Date().toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
    devices = devices.map(d =>
      d.id === assignForm.deviceId
        ? { ...d, assignedTo: assignForm.assignTo, assignedDate: now, status: 'assigned' }
        : d
    );
    assignForm = { deviceId: '', assignTo: '' };
    formVisible = false;
  }

  function returnDevice(id: string) {
    devices = devices.map(d =>
      d.id === id ? { ...d, assignedTo: null, assignedDate: null, status: 'available' } : d
    );
  }

  const typeLabel: Record<DeviceType, string> = {
    handphone: 'Handphone',
    walkie_talkie: 'Walkie-Talkie',
    tablet: 'Tablet',
    power_bank: 'Power Bank',
  };

  const typeIcon: Record<DeviceType, string> = {
    handphone: 'smartphone',
    walkie_talkie: 'radio',
    tablet: 'tablet',
    power_bank: 'battery_charging_full',
  };

  const statusLabel: Record<DeviceStatus, string> = {
    available: 'Tersedia',
    assigned: 'Dipinjam',
    maintenance: 'Servis',
    lost: 'Hilang',
  };

  const statusColor: Record<DeviceStatus, string> = {
    available: '#065f46',
    assigned: '#004ac6',
    maintenance: '#7d4f00',
    lost: '#93000a',
  };

  const statusBg: Record<DeviceStatus, string> = {
    available: '#d1fae5',
    assigned: 'rgba(37,99,235,0.1)',
    maintenance: '#fff3cd',
    lost: '#ffdad6',
  };
</script>

<main class="page-shell">
  <header class="topbar">
    <div class="search-wrap">
      <span class="material-symbols-outlined search-icon">search</span>
      <input type="text" placeholder="Cari perangkat, IMEI, atau assignee..." bind:value={searchQuery} />
    </div>
    <div class="top-right">
      <button class="add-btn" onclick={() => formVisible = !formVisible}>
        <span class="material-symbols-outlined">add</span>
        Assign Perangkat
      </button>
    </div>
  </header>

  <section class="canvas">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Perangkat Komunikasi</span>
    </nav>

    <div class="page-head">
      <h2>Manajemen Perangkat Komunikasi</h2>
      <p>BL-JMJ-010 — Inventaris dan serah terima perangkat lapangan</p>
    </div>

    <!-- Stats -->
    <div class="stats-row">
      <div class="stat-card"><div class="stat-val">{stats.total}</div><div class="stat-label">Total Perangkat</div></div>
      <div class="stat-card"><div class="stat-val" style="color:#004ac6">{stats.assigned}</div><div class="stat-label">Dipinjam</div></div>
      <div class="stat-card"><div class="stat-val" style="color:#065f46">{stats.available}</div><div class="stat-label">Tersedia</div></div>
      <div class="stat-card"><div class="stat-val" style="color:#7d4f00">{stats.maintenance}</div><div class="stat-label">Servis</div></div>
    </div>

    <!-- Assign form -->
    {#if formVisible}
      <div class="assign-form-card">
        <h3>Assign Perangkat ke Staf</h3>
        <form class="assign-form" onsubmit={handleAssign}>
          <div class="field">
            <label>Pilih Perangkat</label>
            <select bind:value={assignForm.deviceId} required>
              <option value="">— Pilih perangkat —</option>
              {#each devices.filter(d => d.status === 'available') as d}
                <option value={d.id}>{d.label} ({typeLabel[d.type]})</option>
              {/each}
            </select>
          </div>
          <div class="field">
            <label>Assign Kepada</label>
            <input type="text" placeholder="Nama staf atau fungsi..." bind:value={assignForm.assignTo} required />
          </div>
          <div class="form-btns">
            <button type="submit" class="btn-primary">
              <span class="material-symbols-outlined">check</span>
              Simpan Penugasan
            </button>
            <button type="button" class="btn-ghost" onclick={() => formVisible = false}>Batal</button>
          </div>
        </form>
      </div>
    {/if}

    <!-- Table -->
    <div class="panel">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Perangkat</th>
              <th>IMEI / Serial</th>
              <th>Jenis</th>
              <th>Dipinjam Oleh</th>
              <th>Tanggal</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each filtered as d (d.id)}
              <tr>
                <td class="device-label-cell">
                  <span class="device-icon">
                    <span class="material-symbols-outlined">{typeIcon[d.type]}</span>
                  </span>
                  {d.label}
                </td>
                <td><span class="imei-code">{d.imei}</span></td>
                <td>{typeLabel[d.type]}</td>
                <td>{d.assignedTo ?? <span style="color:#b0b3c1">—</span>}</td>
                <td>
                  {#if d.assignedDate}
                    <span class="timestamp">{d.assignedDate}</span>
                  {:else}
                    <span style="color:#b0b3c1">—</span>
                  {/if}
                </td>
                <td>
                  <span class="status-badge" style="background: {statusBg[d.status]}; color: {statusColor[d.status]}">
                    {statusLabel[d.status]}
                  </span>
                </td>
                <td>
                  {#if d.status === 'assigned'}
                    <button class="action-btn" onclick={() => returnDevice(d.id)}>
                      <span class="material-symbols-outlined">undo</span>
                      Kembalikan
                    </button>
                  {:else if d.status === 'available'}
                    <button class="action-btn assign" onclick={() => { formVisible = true; assignForm.deviceId = d.id; }}>
                      <span class="material-symbols-outlined">person_add</span>
                      Assign
                    </button>
                  {:else}
                    <span style="color:#b0b3c1; font-size:0.72rem">—</span>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <div class="table-footer">
        Menampilkan {filtered.length} dari {devices.length} perangkat
      </div>
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; justify-content: space-between; gap: 1rem; backdrop-filter: blur(8px); }
  .search-wrap { flex: 1; max-width: 32rem; position: relative; }
  .search-icon { position: absolute; left: 0.65rem; top: 50%; transform: translateY(-50%); font-size: 1.05rem; color: #737686; }
  .search-wrap input { width: 100%; border: 1px solid rgb(195 198 215 / 0.55); background: #f2f4f6; border-radius: 0.25rem; padding: 0.48rem 0.7rem 0.48rem 2.1rem; font-size: 0.85rem; color: #191c1e; }
  .top-right { display: flex; align-items: center; gap: 0.5rem; }
  .add-btn { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.45rem 0.85rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .add-btn .material-symbols-outlined { font-size: 1rem; }
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; margin-bottom: 1rem; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .page-head { margin-bottom: 1.5rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .stats-row { display: flex; gap: 1rem; flex-wrap: wrap; margin-bottom: 1.5rem; }
  .stat-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1rem 1.5rem; flex: 1; min-width: 120px; }
  .stat-val { font-size: 1.8rem; font-weight: 800; color: #004ac6; font-family: 'Plus Jakarta Sans', sans-serif; }
  .stat-label { font-size: 0.72rem; color: #737686; margin-top: 0.15rem; }
  /* Assign form */
  .assign-form-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; margin-bottom: 1.5rem; }
  .assign-form-card h3 { margin: 0 0 1rem; font-size: 0.92rem; font-weight: 700; }
  .assign-form { display: flex; gap: 1rem; flex-wrap: wrap; align-items: flex-end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; min-width: 200px; flex: 1; }
  .field label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field select, .field input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; }
  .form-btns { display: flex; gap: 0.5rem; align-items: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.3rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.42rem 0.85rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; white-space: nowrap; }
  .btn-ghost { background: #fff; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.85rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; color: #434655; font-family: inherit; white-space: nowrap; }
  /* Table */
  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.25rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .device-label-cell { display: flex; align-items: center; gap: 0.5rem; }
  .device-icon { width: 1.8rem; height: 1.8rem; border-radius: 0.3rem; background: rgba(37,99,235,0.08); display: grid; place-items: center; color: #2563eb; flex-shrink: 0; }
  .device-icon .material-symbols-outlined { font-size: 0.95rem; }
  .imei-code { font-family: 'IBM Plex Mono', monospace; font-size: 0.68rem; color: #434655; }
  .status-badge { font-size: 0.68rem; font-weight: 700; padding: 0.2rem 0.55rem; border-radius: 0.2rem; }
  .timestamp { font-size: 0.68rem; color: #737686; }
  .action-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.28rem 0.55rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.7rem; font-weight: 600; cursor: pointer; color: #191c1e; font-family: inherit; }
  .action-btn:hover { background: #f2f4f6; }
  .action-btn .material-symbols-outlined { font-size: 0.85rem; }
  .action-btn.assign .material-symbols-outlined { color: #2563eb; }
  .table-footer { padding: 0.55rem 0.85rem; font-size: 0.68rem; color: #434655; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; }
</style>
