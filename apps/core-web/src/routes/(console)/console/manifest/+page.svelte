<script lang="ts">
  import type { DepartureSummary, DepartureStatus, PageData } from './+page.server';

  let { data }: { data: PageData } = $props();

  // ---- local state — synced via $effect ----
  let departures = $state<DepartureSummary[]>([]);

  $effect(() => {
    departures = data.departures ?? [];
  });

  // ---- filter state ----
  let statusFilter = $state<DepartureStatus | 'all'>('all');

  // ---- derived ----
  const filtered = $derived(
    statusFilter === 'all'
      ? departures
      : departures.filter((d) => d.status === statusFilter)
  );

  // ---- helpers ----
  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });
  }

  const STATUS_LABELS: Record<DepartureStatus, string> = {
    upcoming: 'Mendatang',
    ongoing: 'Berlangsung',
    completed: 'Selesai',
    cancelled: 'Dibatalkan'
  };

  const STATUS_FILTERS: Array<{ value: DepartureStatus | 'all'; label: string }> = [
    { value: 'all', label: 'Semua' },
    { value: 'upcoming', label: 'Mendatang' },
    { value: 'ongoing', label: 'Berlangsung' },
    { value: 'completed', label: 'Selesai' },
    { value: 'cancelled', label: 'Dibatalkan' }
  ];
</script>

<main class="page-shell">
  <!-- Topbar -->
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <a href="/console/dashboard" class="back-link">
        <span class="material-symbols-outlined">chevron_left</span>
        Dashboard
      </a>
      <span class="breadcrumb-sep">/</span>
      <span class="topbar-current">Manifest Jamaah</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn" title="Notifikasi">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar" aria-label="Profile">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Manifest Jamaah</h2>
        <p>Pilih keberangkatan untuk melihat daftar jamaah</p>
      </div>
    </div>

    {#if data.error}
      <div class="error-banner" role="alert">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- Filter tabs -->
    <div class="filter-tabs">
      {#each STATUS_FILTERS as f}
        <button
          type="button"
          class="filter-tab"
          class:active={statusFilter === f.value}
          onclick={() => { statusFilter = f.value as DepartureStatus | 'all'; }}
        >
          {f.label}
          {#if f.value !== 'all'}
            <span class="tab-count">{departures.filter(d => d.status === f.value).length}</span>
          {:else}
            <span class="tab-count">{departures.length}</span>
          {/if}
        </button>
      {/each}
    </div>

    <!-- Table -->
    {#if filtered.length === 0}
      <div class="empty-state panel-empty">
        <span class="material-symbols-outlined">flight_takeoff</span>
        <p>Tidak ada keberangkatan{statusFilter !== 'all' ? ` dengan status "${STATUS_LABELS[statusFilter as DepartureStatus]}"` : ''}.</p>
      </div>
    {:else}
      <div class="panel">
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Nama Paket</th>
                <th>Tanggal Berangkat</th>
                <th>Tanggal Kembali</th>
                <th class="align-right">Total Jamaah</th>
                <th class="align-right">Kursi</th>
                <th>Status</th>
                <th class="align-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each filtered as dep (dep.id)}
                {@const remaining = dep.total_seats - dep.booked_seats}
                <tr>
                  <td>
                    <span class="dep-name">{dep.package_name}</span>
                  </td>
                  <td>
                    <span class="date-cell">{formatDate(dep.departure_date)}</span>
                  </td>
                  <td>
                    <span class="date-cell">{formatDate(dep.return_date)}</span>
                  </td>
                  <td class="align-right">
                    <span class="num-cell">{dep.booked_seats}</span>
                  </td>
                  <td class="align-right">
                    <span class="seats-cell">
                      <span class="seats-booked">{dep.booked_seats}</span>
                      <span class="seats-sep">/</span>
                      <span class="seats-total">{dep.total_seats}</span>
                      {#if remaining > 0}
                        <span class="seats-remaining">({remaining} sisa)</span>
                      {/if}
                    </span>
                  </td>
                  <td>
                    <span class="status-badge status-badge--{dep.status}">
                      {STATUS_LABELS[dep.status]}
                    </span>
                  </td>
                  <td class="align-right">
                    <a href="/console/manifest/{dep.id}" class="lihat-btn">
                      <span class="material-symbols-outlined">list_alt</span>
                      Lihat Manifest
                    </a>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="table-footer">
          {filtered.length} keberangkatan ditampilkan
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
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.1rem;
    font-size: 0.82rem;
    color: #434655;
    text-decoration: none;
    font-weight: 500;
  }

  .back-link:hover { color: #004ac6; }
  .back-link .material-symbols-outlined { font-size: 1rem; }

  .breadcrumb-sep {
    color: #b0b3c1;
    font-size: 0.78rem;
  }

  .topbar-current {
    font-size: 0.88rem;
    font-weight: 600;
    color: #191c1e;
  }

  .top-actions {
    display: flex;
    align-items: center;
    gap: 0.35rem;
  }

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

  /* ---- filter tabs ---- */
  .filter-tabs {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .filter-tab {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.4rem 0.75rem;
    border-radius: 0.25rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    font-size: 0.78rem;
    font-weight: 600;
    color: #434655;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
  }

  .filter-tab:hover { background: #f2f4f6; }
  .filter-tab.active {
    border-color: #2563eb;
    color: #004ac6;
    background: rgb(37 99 235 / 0.07);
  }

  .tab-count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: #eceef0;
    color: #434655;
    border-radius: 999px;
    font-size: 0.62rem;
    font-weight: 700;
    min-width: 1.2rem;
    height: 1.2rem;
    padding: 0 0.25rem;
  }

  .filter-tab.active .tab-count {
    background: #dbeafe;
    color: #1e3a8a;
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

  .align-right { text-align: right; }

  /* ---- cells ---- */
  .dep-name {
    font-weight: 700;
    color: #191c1e;
    font-size: 0.82rem;
  }

  .date-cell { color: #434655; }

  .num-cell {
    font-weight: 600;
    color: #191c1e;
    font-variant-numeric: tabular-nums;
  }

  .seats-cell {
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
    font-variant-numeric: tabular-nums;
  }

  .seats-booked { font-weight: 700; color: #191c1e; }
  .seats-sep { color: #b0b3c1; }
  .seats-total { color: #434655; }
  .seats-remaining { font-size: 0.68rem; color: #065f46; }

  /* ---- status badge ---- */
  .status-badge {
    display: inline-flex;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
  }

  .status-badge--upcoming { background: #dbeafe; color: #1e3a8a; }
  .status-badge--ongoing { background: #d1fae5; color: #065f46; }
  .status-badge--completed { background: #f1f5f9; color: #475569; }
  .status-badge--cancelled { background: #fee2e2; color: #991b1b; }

  /* ---- lihat manifest button ---- */
  .lihat-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.3rem 0.65rem;
    border-radius: 0.2rem;
    border: 1px solid #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    font-size: 0.72rem;
    font-weight: 600;
    color: #fff;
    text-decoration: none;
    transition: opacity 0.1s;
  }

  .lihat-btn:hover { opacity: 0.85; }
  .lihat-btn .material-symbols-outlined { font-size: 0.9rem; }

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
