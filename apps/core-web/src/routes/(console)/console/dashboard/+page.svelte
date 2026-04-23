<script lang="ts">
  import type {
    KpiData,
    BookingStatusCount,
    MonthlyRevenue,
    UpcomingDeparture,
    RecentLead,
    PageData
  } from './+page.server';

  let { data }: { data: PageData } = $props();

  // ---- local state — synced via $effect (Svelte 5 rule) ----
  let kpi = $state<KpiData>({
    total_bookings_month: 0,
    total_revenue_month: 0,
    seats_sold: 0,
    leads_new_month: 0
  });
  let bookingByStatus = $state<BookingStatusCount[]>([]);
  let revenueChart = $state<MonthlyRevenue[]>([]);
  let upcomingDepartures = $state<UpcomingDeparture[]>([]);
  let recentLeads = $state<RecentLead[]>([]);

  $effect(() => {
    kpi = data.kpi ?? kpi;
    bookingByStatus = data.bookingByStatus ?? [];
    revenueChart = data.revenueChart ?? [];
    upcomingDepartures = data.upcomingDepartures ?? [];
    recentLeads = data.recentLeads ?? [];
  });

  // ---- helpers ----
  function formatIDR(amount: number): string {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(amount);
  }

  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });
  }

  function relativeTime(iso: string): string {
    const diff = Date.now() - new Date(iso).getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);
    if (minutes < 1) return 'Baru saja';
    if (minutes < 60) return `${minutes} menit lalu`;
    if (hours < 24) return `${hours} jam lalu`;
    if (days === 1) return 'Kemarin';
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }

  // Lead status labels & colors (consistent with leads page)
  const LEAD_STATUS_LABELS: Record<string, string> = {
    new: 'Baru',
    contacted: 'Dihubungi',
    qualified: 'Qualified',
    converted: 'Konversi',
    lost: 'Tidak Jadi'
  };

  // Max revenue for bar chart scaling
  const maxRevenue = $derived(Math.max(...revenueChart.map((r) => r.revenue), 1));
</script>

<main class="page-shell">
  <!-- Topbar -->
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">dashboard</span>
      <span class="topbar-current">Dashboard</span>
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
        <h2>Executive Dashboard</h2>
        <p>Ringkasan performa bisnis bulan ini</p>
      </div>
    </div>

    {#if data.error}
      <div class="error-banner" role="alert">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- ================================================================
         Section 1 — KPI Cards
    ================================================================= -->
    <div class="kpi-grid">
      <!-- Total Booking -->
      <div class="kpi-card">
        <div class="kpi-icon-wrap kpi-blue">
          <span class="material-symbols-outlined">book_online</span>
        </div>
        <div class="kpi-body">
          <span class="kpi-value">{kpi.total_bookings_month}</span>
          <span class="kpi-label">Total Booking Bulan Ini</span>
        </div>
      </div>

      <!-- Total Pendapatan -->
      <div class="kpi-card">
        <div class="kpi-icon-wrap kpi-green">
          <span class="material-symbols-outlined">payments</span>
        </div>
        <div class="kpi-body">
          <span class="kpi-value kpi-idr">{formatIDR(kpi.total_revenue_month)}</span>
          <span class="kpi-label">Total Pendapatan</span>
        </div>
      </div>

      <!-- Kursi Terjual -->
      <div class="kpi-card">
        <div class="kpi-icon-wrap kpi-purple">
          <span class="material-symbols-outlined">airline_seat_recline_normal</span>
        </div>
        <div class="kpi-body">
          <span class="kpi-value">{kpi.seats_sold}</span>
          <span class="kpi-label">Kursi Terjual</span>
        </div>
      </div>

      <!-- Lead Baru -->
      <div class="kpi-card">
        <div class="kpi-icon-wrap kpi-amber">
          <span class="material-symbols-outlined">person_add</span>
        </div>
        <div class="kpi-body">
          <span class="kpi-value">{kpi.leads_new_month}</span>
          <span class="kpi-label">Lead Baru</span>
        </div>
      </div>
    </div>

    <!-- ================================================================
         Section 2 — Revenue Chart + Booking by Status (side by side)
    ================================================================= -->
    <div class="charts-row">
      <!-- Revenue Trend -->
      <section class="section-block chart-block">
        <div class="section-header">
          <span class="material-symbols-outlined section-icon">bar_chart</span>
          <h3>Tren Pendapatan</h3>
        </div>
        <div class="panel chart-panel">
          {#if revenueChart.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">bar_chart</span>
              <p>Belum ada data pendapatan.</p>
            </div>
          {:else}
            <div class="bar-chart">
              {#each revenueChart as point (point.month)}
                {@const pct = Math.round((point.revenue / maxRevenue) * 100)}
                <div class="bar-item">
                  <span class="bar-value">{formatIDR(point.revenue)}</span>
                  <div class="bar-track">
                    <div class="bar-fill" style="height:{pct}%"></div>
                  </div>
                  <span class="bar-label">{point.month}</span>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </section>

      <!-- Booking by Status -->
      <section class="section-block chart-block">
        <div class="section-header">
          <span class="material-symbols-outlined section-icon">donut_large</span>
          <h3>Booking per Status</h3>
        </div>
        <div class="panel status-panel">
          {#if bookingByStatus.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">donut_large</span>
              <p>Belum ada data booking.</p>
            </div>
          {:else}
            <div class="status-list">
              {#each bookingByStatus as item (item.status)}
                {@const total = bookingByStatus.reduce((a, b) => a + b.count, 0)}
                {@const pct = total > 0 ? Math.round((item.count / total) * 100) : 0}
                <div class="status-row">
                  <span class="status-dot" style="background:{item.color}"></span>
                  <span class="status-row-label">{item.label}</span>
                  <div class="status-bar-track">
                    <div class="status-bar-fill" style="width:{pct}%; background:{item.color}"></div>
                  </div>
                  <span class="status-count">{item.count}</span>
                  <span class="status-pct">{pct}%</span>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </section>
    </div>

    <!-- ================================================================
         Section 3 — Keberangkatan Mendatang
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">flight_takeoff</span>
        <h3>Keberangkatan Mendatang</h3>
      </div>

      {#if upcomingDepartures.length === 0}
        <div class="empty-state panel-empty">
          <span class="material-symbols-outlined">flight_takeoff</span>
          <p>Tidak ada keberangkatan mendatang.</p>
        </div>
      {:else}
        <div class="panel">
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Paket</th>
                  <th>Tanggal Berangkat</th>
                  <th class="align-right">Kursi Tersisa</th>
                  <th class="align-right">Total Kursi</th>
                </tr>
              </thead>
              <tbody>
                {#each upcomingDepartures.slice(0, 5) as dep (dep.id)}
                  <tr>
                    <td>
                      <span class="dep-name">{dep.package_name}</span>
                    </td>
                    <td>
                      <span class="date-cell">{formatDate(dep.departure_date)}</span>
                    </td>
                    <td class="align-right">
                      <span class="seats-remaining" class:seats-low={dep.seats_remaining <= 5}>
                        {dep.seats_remaining}
                      </span>
                    </td>
                    <td class="align-right">
                      <span class="seats-total">{dep.total_seats}</span>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          <div class="table-footer">
            {upcomingDepartures.length} keberangkatan
          </div>
        </div>
      {/if}
    </section>

    <!-- ================================================================
         Section 4 — Lead Terbaru
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">group_add</span>
        <h3>Lead Terbaru</h3>
      </div>

      {#if recentLeads.length === 0}
        <div class="empty-state panel-empty">
          <span class="material-symbols-outlined">person_add</span>
          <p>Belum ada lead terbaru.</p>
        </div>
      {:else}
        <div class="panel">
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Nama</th>
                  <th>Status</th>
                  <th>CS Assigned</th>
                  <th>Waktu</th>
                </tr>
              </thead>
              <tbody>
                {#each recentLeads.slice(0, 5) as lead (lead.id)}
                  <tr>
                    <td>
                      <span class="lead-name">{lead.name}</span>
                    </td>
                    <td>
                      <span class="status-badge status-badge--{lead.status}">
                        {LEAD_STATUS_LABELS[lead.status] ?? lead.status}
                      </span>
                    </td>
                    <td>
                      {#if lead.cs_name}
                        <span class="cs-name">
                          <span class="cs-avatar">{lead.cs_name.charAt(0)}</span>
                          {lead.cs_name}
                        </span>
                      {:else}
                        <span class="unassigned">Belum assign</span>
                      {/if}
                    </td>
                    <td>
                      <span class="time-cell" title={new Date(lead.created_at).toLocaleString('id-ID')}>
                        {relativeTime(lead.created_at)}
                      </span>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          <div class="table-footer">
            <a href="/console/leads" class="view-all-link">
              Lihat semua lead
              <span class="material-symbols-outlined">arrow_forward</span>
            </a>
          </div>
        </div>
      {/if}
    </section>
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
    gap: 0.5rem;
    color: #434655;
  }

  .breadcrumb-icon {
    font-size: 1.1rem;
    color: #004ac6;
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
    margin-bottom: 1.5rem;
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

  /* ---- KPI cards ---- */
  .kpi-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  @media (max-width: 900px) {
    .kpi-grid { grid-template-columns: repeat(2, 1fr); }
  }

  .kpi-card {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1.1rem;
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .kpi-icon-wrap {
    width: 2.75rem;
    height: 2.75rem;
    border-radius: 0.35rem;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }

  .kpi-icon-wrap .material-symbols-outlined {
    font-size: 1.4rem;
  }

  .kpi-blue { background: #dbeafe; color: #1e3a8a; }
  .kpi-green { background: #d1fae5; color: #065f46; }
  .kpi-purple { background: #ede9fe; color: #4c1d95; }
  .kpi-amber { background: #fef3c7; color: #7d5f00; }

  .kpi-body {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
    min-width: 0;
  }

  .kpi-value {
    font-size: 1.65rem;
    font-weight: 800;
    color: #191c1e;
    line-height: 1;
    font-variant-numeric: tabular-nums;
  }

  .kpi-idr {
    font-size: 1.15rem;
  }

  .kpi-label {
    font-size: 0.72rem;
    color: #434655;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  /* ---- charts row ---- */
  .charts-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
    margin-bottom: 0;
  }

  @media (max-width: 860px) {
    .charts-row { grid-template-columns: 1fr; }
  }

  .chart-block { margin-bottom: 1.5rem; }

  /* ---- section blocks ---- */
  .section-block {
    margin-bottom: 1.5rem;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.85rem;
  }

  .section-header h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: #191c1e;
  }

  .section-icon {
    font-size: 1.1rem;
    color: #004ac6;
  }

  /* ---- panel ---- */
  .panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  /* ---- bar chart ---- */
  .chart-panel {
    padding: 1.25rem 1rem 1rem;
  }

  .bar-chart {
    display: flex;
    align-items: flex-end;
    gap: 0.75rem;
    height: 8rem;
  }

  .bar-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 100%;
    gap: 0.3rem;
  }

  .bar-value {
    font-size: 0.55rem;
    color: #434655;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
    text-align: center;
  }

  .bar-track {
    flex: 1;
    width: 100%;
    display: flex;
    align-items: flex-end;
    background: #f2f4f6;
    border-radius: 0.15rem 0.15rem 0 0;
    overflow: hidden;
  }

  .bar-fill {
    width: 100%;
    background: linear-gradient(180deg, #2563eb, #004ac6);
    border-radius: 0.15rem 0.15rem 0 0;
    min-height: 2px;
    transition: height 0.3s ease;
  }

  .bar-label {
    font-size: 0.65rem;
    color: #434655;
    font-weight: 600;
  }

  /* ---- status panel ---- */
  .status-panel {
    padding: 1rem;
  }

  .status-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .status-row {
    display: flex;
    align-items: center;
    gap: 0.6rem;
  }

  .status-dot {
    width: 0.6rem;
    height: 0.6rem;
    border-radius: 999px;
    flex-shrink: 0;
  }

  .status-row-label {
    font-size: 0.78rem;
    color: #191c1e;
    min-width: 6rem;
    flex-shrink: 0;
  }

  .status-bar-track {
    flex: 1;
    height: 0.5rem;
    background: #f2f4f6;
    border-radius: 999px;
    overflow: hidden;
  }

  .status-bar-fill {
    height: 100%;
    border-radius: 999px;
    min-width: 2px;
    transition: width 0.3s ease;
  }

  .status-count {
    font-size: 0.78rem;
    font-weight: 700;
    color: #191c1e;
    min-width: 1.5rem;
    text-align: right;
    font-variant-numeric: tabular-nums;
  }

  .status-pct {
    font-size: 0.68rem;
    color: #737686;
    min-width: 2.5rem;
    text-align: right;
  }

  /* ---- table ---- */
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
    font-weight: 600;
    color: #191c1e;
    font-size: 0.82rem;
  }

  .date-cell {
    color: #434655;
  }

  .seats-remaining {
    font-weight: 700;
    color: #065f46;
    font-variant-numeric: tabular-nums;
  }

  .seats-remaining.seats-low { color: #ba1a1a; }

  .seats-total {
    color: #434655;
    font-variant-numeric: tabular-nums;
  }

  /* ---- lead row ---- */
  .lead-name {
    font-weight: 700;
    color: #191c1e;
    font-size: 0.82rem;
  }

  /* ---- status badge (consistent with leads page) ---- */
  .status-badge {
    display: inline-flex;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
  }

  .status-badge--new { background: #dbeafe; color: #1e3a8a; }
  .status-badge--contacted { background: #fef9c3; color: #7d5f00; }
  .status-badge--qualified { background: #ede9fe; color: #4c1d95; }
  .status-badge--converted { background: #d1fae5; color: #065f46; }
  .status-badge--lost { background: #fee2e2; color: #991b1b; }

  /* ---- cs cell ---- */
  .cs-name {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.75rem;
    color: #191c1e;
  }

  .cs-avatar {
    width: 1.5rem;
    height: 1.5rem;
    border-radius: 999px;
    background: #b4c5ff;
    color: #00174b;
    font-size: 0.62rem;
    font-weight: 700;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }

  .unassigned {
    font-size: 0.72rem;
    color: #b0b3c1;
    font-style: italic;
  }

  .time-cell {
    font-size: 0.72rem;
    color: #434655;
  }

  /* ---- table footer ---- */
  .table-footer {
    padding: 0.55rem 0.85rem;
    font-size: 0.68rem;
    color: #434655;
    border-top: 1px solid rgb(195 198 215 / 0.35);
    background: #f7f9fb;
    display: flex;
    align-items: center;
  }

  .view-all-link {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    color: #004ac6;
    font-weight: 600;
    font-size: 0.72rem;
    text-decoration: none;
  }

  .view-all-link:hover { text-decoration: underline; }
  .view-all-link .material-symbols-outlined { font-size: 0.85rem; }

  /* ---- empty states ---- */
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
