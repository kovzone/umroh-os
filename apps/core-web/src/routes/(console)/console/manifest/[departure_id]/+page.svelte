<script lang="ts">
  import type {
    ManifestEntry,
    DepartureInfo,
    BookingStatus,
    DocStatus,
    RoomType,
    PageData
  } from './+page.server';

  let { data }: { data: PageData } = $props();

  // ---- local state — synced via $effect ----
  let departure = $state<DepartureInfo>({
    id: '',
    package_name: '',
    departure_date: '',
    return_date: '',
    total_seats: 0,
    booked_seats: 0
  });
  let manifest = $state<ManifestEntry[]>([]);

  $effect(() => {
    departure = data.departure ?? departure;
    manifest = data.manifest ?? [];
  });

  // ---- filter state ----
  let filterBookingStatus = $state<BookingStatus | 'all'>('all');
  let filterRoomType = $state<RoomType | 'all'>('all');

  // ---- derived filtered list ----
  const filteredJamaah = $derived(
    manifest.filter((j) => {
      const matchStatus = filterBookingStatus === 'all' || j.booking_status === filterBookingStatus;
      const matchRoom = filterRoomType === 'all' || j.room_type === filterRoomType;
      return matchStatus && matchRoom;
    })
  );

  // ---- summary counts ----
  const countPaid = $derived(manifest.filter((j) => j.booking_status === 'paid').length);
  const countDocComplete = $derived(
    manifest.filter((j) => j.doc_status === 'complete' || j.doc_status === 'verified').length
  );
  const seatsRemaining = $derived(departure.total_seats - departure.booked_seats);

  // ---- helpers ----
  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });
  }

  const BOOKING_STATUS_LABELS: Record<BookingStatus, string> = {
    registered: 'Terdaftar',
    dp_paid: 'Sudah DP',
    paid: 'Lunas',
    cancelled: 'Batal'
  };

  const DOC_STATUS_LABELS: Record<DocStatus, string> = {
    incomplete: 'Belum',
    complete: 'Sebagian',
    verified: 'Lengkap'
  };

  const ROOM_TYPE_LABELS: Record<RoomType, string> = {
    double: 'Double',
    triple: 'Triple',
    quad: 'Quad'
  };

  // ---- export CSV ----
  function exportCSV() {
    const rows = filteredJamaah.map((j) =>
      [j.name, j.nik, j.phone, j.room_type, j.booking_status].join(',')
    );
    const csv = ['Nama,NIK,Telepon,Room,Status', ...rows].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'manifest.csv';
    a.click();
    URL.revokeObjectURL(url);
  }

  // ---- progress bar width ----
  function pct(count: number, total: number): number {
    if (total === 0) return 0;
    return Math.round((count / total) * 100);
  }
</script>

<main class="page-shell">
  <!-- Topbar -->
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <a href="/console/manifest" class="back-link">
        <span class="material-symbols-outlined">chevron_left</span>
        Manifest
      </a>
      <span class="breadcrumb-sep">/</span>
      <span class="topbar-current">
        {departure.package_name || 'Detail'}
        {#if departure.departure_date}
          — {formatDate(departure.departure_date)}
        {/if}
      </span>
    </nav>
    <div class="top-actions">
      <button class="export-btn" onclick={exportCSV} type="button">
        <span class="material-symbols-outlined">download</span>
        Export CSV
      </button>
      <button class="icon-btn" title="Notifikasi">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar" aria-label="Profile">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>{departure.package_name || 'Manifest Keberangkatan'}</h2>
        {#if departure.departure_date}
          <p>
            Berangkat: {formatDate(departure.departure_date)}
            {#if departure.return_date}
              &mdash; Kembali: {formatDate(departure.return_date)}
            {/if}
          </p>
        {/if}
      </div>
    </div>

    {#if data.error}
      <div class="error-banner" role="alert">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- ================================================================
         Summary Cards
    ================================================================= -->
    <div class="summary-grid">
      <!-- Total Jamaah -->
      <div class="summary-card">
        <div class="summary-icon-wrap summary-blue">
          <span class="material-symbols-outlined">groups</span>
        </div>
        <div class="summary-body">
          <span class="summary-value">{manifest.length}</span>
          <span class="summary-label">Total Jamaah</span>
        </div>
      </div>

      <!-- Sudah Lunas -->
      <div class="summary-card">
        <div class="summary-icon-wrap summary-green">
          <span class="material-symbols-outlined">payments</span>
        </div>
        <div class="summary-body">
          <span class="summary-value">{countPaid}</span>
          <span class="summary-label">Sudah Lunas</span>
          <div class="progress-track">
            <div
              class="progress-fill progress-green"
              style="width:{pct(countPaid, manifest.length)}%"
            ></div>
          </div>
          <span class="summary-sub">{pct(countPaid, manifest.length)}% dari total</span>
        </div>
      </div>

      <!-- Dokumen Lengkap -->
      <div class="summary-card">
        <div class="summary-icon-wrap summary-purple">
          <span class="material-symbols-outlined">folder_check</span>
        </div>
        <div class="summary-body">
          <span class="summary-value">{countDocComplete}</span>
          <span class="summary-label">Dokumen Lengkap</span>
          <div class="progress-track">
            <div
              class="progress-fill progress-purple"
              style="width:{pct(countDocComplete, manifest.length)}%"
            ></div>
          </div>
          <span class="summary-sub">{pct(countDocComplete, manifest.length)}% dari total</span>
        </div>
      </div>

      <!-- Kursi Tersisa -->
      <div class="summary-card">
        <div class="summary-icon-wrap summary-amber">
          <span class="material-symbols-outlined">airline_seat_recline_normal</span>
        </div>
        <div class="summary-body">
          <span class="summary-value" class:seats-low={seatsRemaining <= 5}>
            {seatsRemaining}
          </span>
          <span class="summary-label">Kursi Tersisa</span>
          <span class="summary-sub">dari {departure.total_seats} total</span>
        </div>
      </div>
    </div>

    <!-- ================================================================
         Filter row
    ================================================================= -->
    <div class="filter-row">
      <div class="filter-group">
        <label for="filter-booking">Status Booking</label>
        <select id="filter-booking" bind:value={filterBookingStatus}>
          <option value="all">Semua Status</option>
          <option value="registered">Terdaftar</option>
          <option value="dp_paid">Sudah DP</option>
          <option value="paid">Lunas</option>
          <option value="cancelled">Batal</option>
        </select>
      </div>

      <div class="filter-group">
        <label for="filter-room">Tipe Kamar</label>
        <select id="filter-room" bind:value={filterRoomType}>
          <option value="all">Semua Tipe</option>
          <option value="double">Double</option>
          <option value="triple">Triple</option>
          <option value="quad">Quad</option>
        </select>
      </div>

      {#if filterBookingStatus !== 'all' || filterRoomType !== 'all'}
        <button
          type="button"
          class="clear-btn"
          onclick={() => { filterBookingStatus = 'all'; filterRoomType = 'all'; }}
        >
          <span class="material-symbols-outlined">close</span>
          Reset
        </button>
      {/if}

      <span class="filter-count">
        {filteredJamaah.length} dari {manifest.length} jamaah
      </span>
    </div>

    <!-- ================================================================
         Manifest Table
    ================================================================= -->
    {#if filteredJamaah.length === 0}
      <div class="empty-state panel-empty">
        <span class="material-symbols-outlined">person_search</span>
        <p>Tidak ada jamaah yang sesuai filter.</p>
      </div>
    {:else}
      <div class="panel">
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th class="align-center">No</th>
                <th>Nama</th>
                <th>NIK</th>
                <th>No HP</th>
                <th>Room</th>
                <th>Status Booking</th>
                <th>Dok</th>
              </tr>
            </thead>
            <tbody>
              {#each filteredJamaah as jamaah (jamaah.id)}
                <tr>
                  <td class="align-center">
                    <span class="no-cell">{jamaah.no}</span>
                  </td>
                  <td>
                    <span class="jamaah-name">{jamaah.name}</span>
                  </td>
                  <td>
                    <span class="nik-cell">{jamaah.nik}</span>
                  </td>
                  <td>
                    <span class="phone-cell">{jamaah.phone}</span>
                  </td>
                  <td>
                    <span class="room-badge room-badge--{jamaah.room_type}">
                      {ROOM_TYPE_LABELS[jamaah.room_type]}
                    </span>
                  </td>
                  <td>
                    <span class="booking-badge booking-badge--{jamaah.booking_status}">
                      {BOOKING_STATUS_LABELS[jamaah.booking_status]}
                    </span>
                  </td>
                  <td>
                    <span class="doc-badge doc-badge--{jamaah.doc_status}">
                      {DOC_STATUS_LABELS[jamaah.doc_status]}
                    </span>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="table-footer">
          {filteredJamaah.length} jamaah ditampilkan
        </div>
      </div>
    {/if}
  </section>
</main>

<style>
  .page-shell {
    min-height: 100vh;
    background: #f7f9fb;
  }

  /* ---- topbar ---- */
  .topbar {
    position: sticky;
    top: 0;
    z-index: 30;
    height: 4rem;
    background: rgb(255 255 255 / 0.9);
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    backdrop-filter: blur(8px);
  }

  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: #434655;
    min-width: 0;
    flex: 1;
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.1rem;
    font-size: 0.82rem;
    color: #434655;
    text-decoration: none;
    font-weight: 500;
    flex-shrink: 0;
  }

  .back-link:hover { color: #004ac6; }
  .back-link .material-symbols-outlined { font-size: 1rem; }

  .breadcrumb-sep {
    color: #b0b3c1;
    font-size: 0.78rem;
    flex-shrink: 0;
  }

  .topbar-current {
    font-size: 0.88rem;
    font-weight: 600;
    color: #191c1e;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .top-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  .export-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.4rem 0.75rem;
    border-radius: 0.25rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    font-size: 0.78rem;
    font-weight: 600;
    color: #191c1e;
    cursor: pointer;
    font-family: inherit;
  }

  .export-btn:hover { background: #f2f4f6; }
  .export-btn .material-symbols-outlined { font-size: 0.9rem; }

  .icon-btn {
    border: 0;
    background: transparent;
    color: #434655;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    cursor: pointer;
    display: grid;
    place-items: center;
  }

  .icon-btn:hover { background: #eceef0; }

  .avatar {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #b4c5ff;
    color: #00174b;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    font-weight: 700;
    font-size: 0.65rem;
    cursor: pointer;
  }

  /* ---- canvas ---- */
  .canvas {
    padding: 1.5rem;
    max-width: 96rem;
  }

  .page-head {
    margin-bottom: 1.25rem;
  }

  .page-head h2 {
    margin: 0;
    font-size: 1.5rem;
    line-height: 1.2;
  }

  .page-head p {
    margin: 0.3rem 0 0;
    font-size: 0.82rem;
    color: #434655;
  }

  /* ---- error banner ---- */
  .error-banner {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: #ffdad6;
    color: #93000a;
    border-radius: 0.25rem;
    padding: 0.65rem 0.85rem;
    font-size: 0.82rem;
    margin-bottom: 1.25rem;
  }

  /* ---- summary grid ---- */
  .summary-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  @media (max-width: 900px) {
    .summary-grid { grid-template-columns: repeat(2, 1fr); }
  }

  .summary-card {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1.1rem;
    display: flex;
    align-items: flex-start;
    gap: 0.85rem;
  }

  .summary-icon-wrap {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 0.3rem;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }

  .summary-icon-wrap .material-symbols-outlined { font-size: 1.3rem; }

  .summary-blue { background: #dbeafe; color: #1e3a8a; }
  .summary-green { background: #d1fae5; color: #065f46; }
  .summary-purple { background: #ede9fe; color: #4c1d95; }
  .summary-amber { background: #fef3c7; color: #7d5f00; }

  .summary-body {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
    min-width: 0;
    flex: 1;
  }

  .summary-value {
    font-size: 1.65rem;
    font-weight: 800;
    color: #191c1e;
    line-height: 1;
    font-variant-numeric: tabular-nums;
  }

  .summary-value.seats-low { color: #ba1a1a; }

  .summary-label {
    font-size: 0.72rem;
    color: #434655;
  }

  .summary-sub {
    font-size: 0.65rem;
    color: #737686;
  }

  .progress-track {
    height: 0.35rem;
    background: #f2f4f6;
    border-radius: 999px;
    overflow: hidden;
    margin-top: 0.25rem;
  }

  .progress-fill {
    height: 100%;
    border-radius: 999px;
    min-width: 2px;
    transition: width 0.3s ease;
  }

  .progress-green { background: #10b981; }
  .progress-purple { background: #8b5cf6; }

  /* ---- filter row ---- */
  .filter-row {
    display: flex;
    align-items: flex-end;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .filter-group {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
  }

  .filter-group label {
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }

  .filter-group select {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    border-radius: 0.25rem;
    padding: 0.42rem 0.6rem;
    font-size: 0.82rem;
    color: #191c1e;
    min-width: 10rem;
    outline: none;
    font-family: inherit;
  }

  .clear-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.42rem 0.6rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    background: #fff;
    font-size: 0.78rem;
    color: #434655;
    cursor: pointer;
    align-self: flex-end;
    font-family: inherit;
  }

  .clear-btn:hover { background: #f2f4f6; }
  .clear-btn .material-symbols-outlined { font-size: 0.9rem; }

  .filter-count {
    font-size: 0.72rem;
    color: #737686;
    align-self: flex-end;
    margin-left: auto;
  }

  /* ---- panel ---- */
  .panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .table-wrap { overflow-x: auto; }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    padding: 0.62rem 0.85rem;
    text-align: left;
    font-size: 0.78rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    white-space: nowrap;
  }

  th {
    text-transform: uppercase;
    font-size: 0.62rem;
    letter-spacing: 0.08em;
    color: #434655;
    background: #f2f4f6;
  }

  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }

  .align-center { text-align: center; }

  /* ---- cells ---- */
  .no-cell {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.75rem;
    height: 1.75rem;
    border-radius: 999px;
    background: #f2f4f6;
    color: #434655;
    font-size: 0.72rem;
    font-weight: 700;
    font-variant-numeric: tabular-nums;
  }

  .jamaah-name {
    font-weight: 700;
    color: #191c1e;
    font-size: 0.82rem;
  }

  .nik-cell {
    font-family: 'IBM Plex Mono', 'Courier New', monospace;
    font-size: 0.72rem;
    color: #434655;
    letter-spacing: 0.03em;
  }

  .phone-cell {
    font-size: 0.75rem;
    color: #434655;
  }

  /* ---- room type badge ---- */
  .room-badge {
    display: inline-flex;
    padding: 0.12rem 0.4rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
    background: #eceef0;
    color: #434655;
  }

  .room-badge--double { background: #dbeafe; color: #1e3a8a; }
  .room-badge--triple { background: #ede9fe; color: #4c1d95; }
  .room-badge--quad { background: #fef3c7; color: #7d5f00; }

  /* ---- booking status badge ---- */
  .booking-badge {
    display: inline-flex;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
  }

  .booking-badge--registered { background: #dbeafe; color: #1e3a8a; }
  .booking-badge--dp_paid { background: #fef9c3; color: #7d5f00; }
  .booking-badge--paid { background: #d1fae5; color: #065f46; }
  .booking-badge--cancelled { background: #fee2e2; color: #991b1b; }

  /* ---- doc status badge ---- */
  .doc-badge {
    display: inline-flex;
    padding: 0.12rem 0.4rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
  }

  .doc-badge--incomplete { background: #fee2e2; color: #991b1b; }
  .doc-badge--complete { background: #fef9c3; color: #7d5f00; }
  .doc-badge--verified { background: #d1fae5; color: #065f46; }

  /* ---- table footer ---- */
  .table-footer {
    padding: 0.55rem 0.85rem;
    font-size: 0.68rem;
    color: #434655;
    border-top: 1px solid rgb(195 198 215 / 0.35);
    background: #f7f9fb;
  }

  /* ---- empty state ---- */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.6rem;
    padding: 3rem 1rem;
    color: #b0b3c1;
  }

  .panel-empty {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
  }

  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.82rem; }
</style>
