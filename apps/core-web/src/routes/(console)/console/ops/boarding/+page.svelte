<script lang="ts">
  type BoardingStatus = 'checked' | 'pending';

  interface Passenger {
    id: string;
    name: string;
    bookingCode: string;
    seatNo: string;
    status: BoardingStatus;
  }

  interface Bus {
    id: string;
    label: string;
    destination: string;
    departureTime: string;
    capacity: number;
    passengers: Passenger[];
  }

  let buses = $state<Bus[]>([
    {
      id: 'bus-001',
      label: 'Bus A1',
      destination: 'Masjidil Haram',
      departureTime: '08:00',
      capacity: 42,
      passengers: [
        { id: 'p1', name: 'Bambang Suryanto', bookingCode: 'ORD-0012', seatNo: '1A', status: 'checked' },
        { id: 'p2', name: 'Siti Rahayu', bookingCode: 'ORD-0013', seatNo: '1B', status: 'checked' },
        { id: 'p3', name: 'Hendra Wijaya', bookingCode: 'ORD-0014', seatNo: '2A', status: 'pending' },
        { id: 'p4', name: 'Dewi Lestari', bookingCode: 'ORD-0015', seatNo: '2B', status: 'pending' },
        { id: 'p5', name: 'Ahmad Fauzi', bookingCode: 'ORD-0016', seatNo: '3A', status: 'pending' },
      ],
    },
    {
      id: 'bus-002',
      label: 'Bus A2',
      destination: 'Masjid Nabawi',
      departureTime: '08:30',
      capacity: 42,
      passengers: [
        { id: 'p6', name: 'Rini Handayani', bookingCode: 'ORD-0017', seatNo: '1A', status: 'checked' },
        { id: 'p7', name: 'Budi Santoso', bookingCode: 'ORD-0018', seatNo: '1B', status: 'pending' },
        { id: 'p8', name: 'Nur Aisyah', bookingCode: 'ORD-0019', seatNo: '2A', status: 'checked' },
      ],
    },
  ]);

  let expandedBus = $state<string | null>('bus-001');
  let qrInput = $state('');

  const boardedCount = $derived(
    buses.reduce((sum, b) => sum + b.passengers.filter(p => p.status === 'checked').length, 0)
  );
  const totalCount = $derived(
    buses.reduce((sum, b) => sum + b.passengers.length, 0)
  );

  function toggleBoarding(busId: string, passId: string) {
    buses = buses.map(b => {
      if (b.id !== busId) return b;
      return {
        ...b,
        passengers: b.passengers.map(p =>
          p.id === passId
            ? { ...p, status: p.status === 'checked' ? 'pending' : 'checked' }
            : p
        ),
      };
    });
  }

  function handleQrScan(e: Event) {
    e.preventDefault();
    if (!qrInput.trim()) return;
    const code = qrInput.trim().toUpperCase();
    let found = false;
    buses = buses.map(b => ({
      ...b,
      passengers: b.passengers.map(p => {
        if (p.bookingCode === code && p.status === 'pending') {
          found = true;
          return { ...p, status: 'checked' };
        }
        return p;
      }),
    }));
    qrInput = '';
    if (!found) alert(`Kode ${code} tidak ditemukan atau sudah di-check.`);
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Bus Boarding</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Absensi Boarding Bus</h2>
        <p>BL-JMJ-007 — Scan QR & catat kehadiran boarding per bus</p>
      </div>
    </div>

    <!-- Stats -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-val">{boardedCount}</div>
        <div class="stat-label">Sudah Boarding</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{totalCount - boardedCount}</div>
        <div class="stat-label">Belum Boarding</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{totalCount}</div>
        <div class="stat-label">Total Penumpang</div>
      </div>
      <div class="stat-card">
        <div class="stat-val">{buses.length}</div>
        <div class="stat-label">Jumlah Bus</div>
      </div>
    </div>

    <!-- QR Scanner placeholder -->
    <div class="qr-section">
      <div class="qr-placeholder">
        <span class="material-symbols-outlined">qr_code_scanner</span>
        <p>Kamera QR Scanner</p>
        <span class="qr-hint">Placeholder — kamera akan aktif di perangkat lapangan</span>
      </div>
      <form class="manual-input" onsubmit={handleQrScan}>
        <input type="text" placeholder="Masukkan kode booking manual..." bind:value={qrInput} />
        <button type="submit" class="btn-scan">Cari & Check-in</button>
      </form>
    </div>

    <!-- Bus list -->
    <div class="bus-list">
      {#each buses as bus (bus.id)}
        {@const checked = bus.passengers.filter(p => p.status === 'checked').length}
        <div class="bus-card">
          <button class="bus-header" onclick={() => expandedBus = expandedBus === bus.id ? null : bus.id}>
            <div class="bus-icon-wrap">
              <span class="material-symbols-outlined">directions_bus</span>
            </div>
            <div class="bus-info">
              <div class="bus-label">{bus.label} — {bus.destination}</div>
              <div class="bus-meta">Berangkat {bus.departureTime} · Kapasitas {bus.capacity}</div>
            </div>
            <div class="bus-progress">
              <div class="bp-counts">{checked}/{bus.passengers.length}</div>
              <div class="bp-bar-wrap">
                <div class="bp-bar" style="width: {bus.passengers.length > 0 ? (checked / bus.passengers.length * 100) : 0}%"></div>
              </div>
            </div>
            <span class="material-symbols-outlined chevron">
              {expandedBus === bus.id ? 'expand_less' : 'expand_more'}
            </span>
          </button>

          {#if expandedBus === bus.id}
            <div class="bus-passengers">
              <table>
                <thead>
                  <tr>
                    <th>Nama</th>
                    <th>Booking</th>
                    <th>Kursi</th>
                    <th>Status</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {#each bus.passengers as p (p.id)}
                    <tr class:checked={p.status === 'checked'}>
                      <td>{p.name}</td>
                      <td><span class="booking-code">{p.bookingCode}</span></td>
                      <td>{p.seatNo}</td>
                      <td>
                        <span class="status-badge" class:checked={p.status === 'checked'}>
                          {p.status === 'checked' ? 'Checked' : 'Pending'}
                        </span>
                      </td>
                      <td>
                        <button class="toggle-btn" onclick={() => toggleBoarding(bus.id, p.id)}>
                          <span class="material-symbols-outlined">
                            {p.status === 'checked' ? 'undo' : 'check'}
                          </span>
                          {p.status === 'checked' ? 'Batalkan' : 'Check-in'}
                        </button>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.5rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  /* Stats */
  .stats-row { display: flex; gap: 1rem; flex-wrap: wrap; margin-bottom: 1.5rem; }
  .stat-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1rem 1.5rem; flex: 1; min-width: 120px; }
  .stat-val { font-size: 1.8rem; font-weight: 800; color: #004ac6; font-family: 'Plus Jakarta Sans', sans-serif; }
  .stat-label { font-size: 0.72rem; color: #737686; margin-top: 0.15rem; }
  /* QR section */
  .qr-section { display: flex; gap: 1rem; align-items: flex-start; margin-bottom: 1.5rem; flex-wrap: wrap; }
  .qr-placeholder {
    width: 180px;
    height: 120px;
    border: 2px dashed rgb(195 198 215 / 0.6);
    border-radius: 0.5rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.25rem;
    color: #737686;
    flex-shrink: 0;
  }
  .qr-placeholder .material-symbols-outlined { font-size: 2rem; color: #b0b3c1; }
  .qr-placeholder p { margin: 0; font-size: 0.75rem; font-weight: 600; }
  .qr-hint { font-size: 0.62rem; color: #b0b3c1; text-align: center; }
  .manual-input { display: flex; gap: 0.5rem; align-self: center; flex: 1; min-width: 240px; }
  .manual-input input { flex: 1; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.48rem 0.7rem; font-size: 0.85rem; color: #191c1e; }
  .btn-scan { background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.48rem 0.85rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; white-space: nowrap; }
  /* Bus cards */
  .bus-list { display: flex; flex-direction: column; gap: 0.85rem; }
  .bus-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .bus-header { display: flex; align-items: center; gap: 1rem; padding: 1rem 1.2rem; background: none; border: none; width: 100%; cursor: pointer; font-family: inherit; text-align: left; }
  .bus-header:hover { background: #f7f9fb; }
  .bus-icon-wrap { width: 2.5rem; height: 2.5rem; border-radius: 0.4rem; background: rgba(37,99,235,0.08); display: grid; place-items: center; color: #2563eb; flex-shrink: 0; }
  .bus-icon-wrap .material-symbols-outlined { font-size: 1.3rem; }
  .bus-info { flex: 1; }
  .bus-label { font-size: 0.88rem; font-weight: 700; color: #191c1e; }
  .bus-meta { font-size: 0.72rem; color: #737686; margin-top: 0.1rem; }
  .bus-progress { min-width: 100px; }
  .bp-counts { font-size: 0.78rem; font-weight: 700; color: #191c1e; margin-bottom: 0.25rem; }
  .bp-bar-wrap { height: 6px; background: rgb(195 198 215 / 0.3); border-radius: 999px; overflow: hidden; }
  .bp-bar { height: 100%; background: linear-gradient(90deg, #004ac6, #2563eb); border-radius: 999px; transition: width 0.3s ease; }
  .chevron { color: #737686; font-size: 1.1rem; flex-shrink: 0; }
  /* Passengers table */
  .bus-passengers { border-top: 1px solid rgb(195 198 215 / 0.35); overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 1.2rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #737686; background: #f2f4f6; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr.checked { background: rgba(220,252,231,0.3); }
  tbody tr:last-child td { border-bottom: none; }
  .booking-code { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; font-weight: 600; color: #004ac6; }
  .status-badge { font-size: 0.68rem; font-weight: 700; padding: 0.2rem 0.55rem; border-radius: 0.2rem; background: #fff3cd; color: #7d4f00; }
  .status-badge.checked { background: #d1fae5; color: #065f46; }
  .toggle-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.28rem 0.55rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.7rem; font-weight: 600; cursor: pointer; color: #191c1e; font-family: inherit; }
  .toggle-btn:hover { background: #f2f4f6; }
  .toggle-btn .material-symbols-outlined { font-size: 0.85rem; }
</style>
