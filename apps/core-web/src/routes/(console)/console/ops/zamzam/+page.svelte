<script lang="ts">
  interface ZamzamRow {
    id: string;
    name: string;
    bookingCode: string;
    allocated: number;
    distributed: number;
    photoProof: boolean;
    distributed_at?: string;
  }

  let rows = $state<ZamzamRow[]>([
    { id: 'z1', name: 'Bambang Suryanto', bookingCode: 'ORD-0012', allocated: 5, distributed: 5, photoProof: true, distributed_at: '24 Jan 2025, 10:30' },
    { id: 'z2', name: 'Siti Rahayu', bookingCode: 'ORD-0013', allocated: 5, distributed: 0, photoProof: false },
    { id: 'z3', name: 'Hendra Wijaya', bookingCode: 'ORD-0014', allocated: 5, distributed: 0, photoProof: false },
    { id: 'z4', name: 'Dewi Lestari', bookingCode: 'ORD-0015', allocated: 5, distributed: 5, photoProof: false, distributed_at: '24 Jan 2025, 10:45' },
    { id: 'z5', name: 'Ahmad Fauzi', bookingCode: 'ORD-0016', allocated: 5, distributed: 0, photoProof: false },
    { id: 'z6', name: 'Rini Handayani', bookingCode: 'ORD-0017', allocated: 5, distributed: 5, photoProof: true, distributed_at: '24 Jan 2025, 11:00' },
    { id: 'z7', name: 'Budi Santoso', bookingCode: 'ORD-0018', allocated: 5, distributed: 0, photoProof: false },
  ]);

  let searchQuery = $state('');

  const filtered = $derived(
    rows.filter(r =>
      !searchQuery ||
      r.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      r.bookingCode.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  const totalAllocated = $derived(rows.reduce((s, r) => s + r.allocated, 0));
  const totalDistributed = $derived(rows.reduce((s, r) => s + r.distributed, 0));
  const doneCount = $derived(rows.filter(r => r.distributed >= r.allocated).length);

  function markDistributed(id: string) {
    const now = new Date();
    const dateStr = now.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
    const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`;
    rows = rows.map(r =>
      r.id === id
        ? { ...r, distributed: r.allocated, distributed_at: `${dateStr}, ${timeStr}` }
        : r
    );
  }

  function togglePhoto(id: string) {
    rows = rows.map(r => r.id === id ? { ...r, photoProof: !r.photoProof } : r);
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <div class="search-wrap">
      <span class="material-symbols-outlined search-icon">search</span>
      <input type="text" placeholder="Cari nama atau booking ID..." bind:value={searchQuery} />
    </div>
    <div class="top-right">
      <button class="icon-btn">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar">A</button>
    </div>
  </header>

  <section class="canvas">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Distribusi Zamzam</span>
    </nav>

    <div class="page-head">
      <h2>Distribusi Air Zamzam</h2>
      <p>BL-JMJ-011 — Manajemen distribusi zamzam per jamaah</p>
    </div>

    <!-- Stats -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-val">{doneCount}</div>
        <div class="stat-label">Selesai Distribusi</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{rows.length - doneCount}</div>
        <div class="stat-label">Belum Didistribusi</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{totalDistributed} L</div>
        <div class="stat-label">Total Terdistribusi</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{totalAllocated} L</div>
        <div class="stat-label">Total Alokasi</div>
      </div>
    </div>

    <!-- Table -->
    <div class="panel">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Nama Jamaah</th>
              <th>Booking</th>
              <th>Alokasi</th>
              <th>Terdistribusi</th>
              <th>Bukti Foto</th>
              <th>Waktu</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each filtered as r (r.id)}
              {@const done = r.distributed >= r.allocated}
              <tr class:done>
                <td>{r.name}</td>
                <td><span class="booking-code">{r.bookingCode}</span></td>
                <td>{r.allocated} L</td>
                <td>
                  <span class="dist-badge" class:done={done} class:pending={!done}>
                    {r.distributed} L
                  </span>
                </td>
                <td>
                  <button class="photo-btn" onclick={() => togglePhoto(r.id)}>
                    <span class="material-symbols-outlined" style="color: {r.photoProof ? '#065f46' : '#b0b3c1'}">
                      {r.photoProof ? 'check_circle' : 'radio_button_unchecked'}
                    </span>
                    {r.photoProof ? 'Ada' : 'Belum'}
                  </button>
                </td>
                <td>
                  {#if r.distributed_at}
                    <span class="timestamp">{r.distributed_at}</span>
                  {:else}
                    <span style="color: #b0b3c1">—</span>
                  {/if}
                </td>
                <td>
                  {#if !done}
                    <button class="action-btn" onclick={() => markDistributed(r.id)}>
                      <span class="material-symbols-outlined">water_drop</span>
                      Tandai Terdistribusi
                    </button>
                  {:else}
                    <span class="done-label">
                      <span class="material-symbols-outlined">check</span>
                      Selesai
                    </span>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <div class="table-footer">
        Menampilkan {filtered.length} dari {rows.length} jamaah
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
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; cursor: pointer; }
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
  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.25rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr.done { background: rgba(220,252,231,0.2); }
  tbody tr:last-child td { border-bottom: 0; }
  .booking-code { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; font-weight: 600; color: #004ac6; }
  .dist-badge { font-size: 0.72rem; font-weight: 700; padding: 0.2rem 0.5rem; border-radius: 0.2rem; }
  .dist-badge.done { background: #d1fae5; color: #065f46; }
  .dist-badge.pending { background: #fff3cd; color: #7d4f00; }
  .photo-btn { display: inline-flex; align-items: center; gap: 0.25rem; background: none; border: none; cursor: pointer; font-size: 0.72rem; color: #434655; font-family: inherit; }
  .photo-btn .material-symbols-outlined { font-size: 1rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .timestamp { font-size: 0.68rem; color: #737686; }
  .action-btn { display: inline-flex; align-items: center; gap: 0.3rem; padding: 0.3rem 0.6rem; border-radius: 0.2rem; border: 1px solid rgb(195 198 215 / 0.55); background: #fff; font-size: 0.7rem; font-weight: 600; color: #191c1e; cursor: pointer; font-family: inherit; }
  .action-btn:hover { background: #f2f4f6; }
  .action-btn .material-symbols-outlined { font-size: 0.85rem; color: #2563eb; }
  .done-label { display: inline-flex; align-items: center; gap: 0.25rem; font-size: 0.72rem; color: #065f46; font-weight: 600; }
  .done-label .material-symbols-outlined { font-size: 0.85rem; }
  .table-footer { padding: 0.55rem 0.85rem; font-size: 0.68rem; color: #434655; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; }
</style>
