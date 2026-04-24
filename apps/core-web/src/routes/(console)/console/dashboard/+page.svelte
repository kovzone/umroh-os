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

  // ---- BL-DASH-005: Dual View ----
  let dualView = $state(false);

  // Mock data for dual view panels
  const salesMetrics = {
    bookingsMonth: 47,
    revenue: 1_342_500_000,
    conversionRate: 38,
    avgTicket: 28_563_830,
    topPackages: [
      { name: 'Umroh Ramadhan Premium', bookings: 18, revenue: 585_000_000 },
      { name: 'Umroh Reguler April', bookings: 14, revenue: 343_000_000 },
      { name: 'Umroh Plus Turki', bookings: 8, revenue: 304_000_000 },
      { name: 'Umroh Hemat Juni', bookings: 7, revenue: 139_300_000 }
    ],
    weeklyTrend: [
      { week: 'Mg 1', bookings: 9 },
      { week: 'Mg 2', bookings: 14 },
      { week: 'Mg 3', bookings: 11 },
      { week: 'Mg 4', bookings: 13 }
    ]
  };

  const opsMetrics = {
    departuresThisWeek: 3,
    paxConfirmed: 127,
    pendingTasks: 8,
    logisticsStatus: [
      { label: 'Dokumen Lengkap', count: 103, total: 127, color: '#065f46' },
      { label: 'Visa Selesai', count: 98, total: 127, color: '#0369a1' },
      { label: 'Kit Terpacking', count: 89, total: 127, color: '#7c3aed' },
      { label: 'Boarding Pass', count: 67, total: 127, color: '#b45309' }
    ],
    upcomingDepartsThisWeek: [
      { name: 'Umroh Ramadhan Premium', date: '2026-03-01', pax: 45, status: 'ready' },
      { name: 'Umroh Reguler April', date: '2026-04-10', pax: 42, status: 'prep' },
      { name: 'Umroh Plus Turki', date: '2026-05-05', pax: 40, status: 'pending' }
    ]
  };

  const maxSalesBookings = $derived(Math.max(...salesMetrics.weeklyTrend.map(w => w.bookings), 1));
</script>

<main class="page-shell">
  <!-- Topbar -->
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">dashboard</span>
      <span class="topbar-current">Dashboard</span>
    </nav>
    <div class="top-actions">
      <!-- BL-DASH-005: Dual view toggle -->
      <div class="view-toggle" role="group" aria-label="Pilih tampilan">
        <button
          type="button"
          class="view-toggle-btn"
          class:active={!dualView}
          onclick={() => { dualView = false; }}
          title="Tampilan Standar"
        >
          <span class="material-symbols-outlined">dashboard</span>
          Standar
        </button>
        <button
          type="button"
          class="view-toggle-btn"
          class:active={dualView}
          onclick={() => { dualView = true; }}
          title="Dual View"
        >
          <span class="material-symbols-outlined">view_column</span>
          Dual View
        </button>
      </div>
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
      {#if dualView}
        <span class="dual-badge">
          <span class="material-symbols-outlined">view_column</span>
          Dual View Aktif
        </span>
      {/if}
    </div>

    <!-- ================================================================
         BL-DASH-005: Dual View — Sales & Ops side by side
    ================================================================= -->
    {#if dualView}
      <div class="dual-layout">
        <!-- LEFT: Sales & Bookings Panel -->
        <div class="dual-panel dual-sales">
          <div class="dual-panel-header">
            <span class="material-symbols-outlined">trending_up</span>
            <h3>Sales & Booking</h3>
          </div>

          <!-- Sales KPIs -->
          <div class="dual-kpi-row">
            <div class="dual-kpi">
              <span class="dkpi-val">{salesMetrics.bookingsMonth}</span>
              <span class="dkpi-lbl">Booking Bulan Ini</span>
            </div>
            <div class="dual-kpi">
              <span class="dkpi-val dkpi-idr">{formatIDR(salesMetrics.revenue)}</span>
              <span class="dkpi-lbl">Revenue</span>
            </div>
            <div class="dual-kpi">
              <span class="dkpi-val">{salesMetrics.conversionRate}%</span>
              <span class="dkpi-lbl">Konversi</span>
            </div>
          </div>

          <!-- Weekly trend -->
          <div class="dual-section">
            <div class="dual-section-label">Booking Mingguan</div>
            <div class="mini-bar-chart">
              {#each salesMetrics.weeklyTrend as w}
                {@const pct = Math.round((w.bookings / maxSalesBookings) * 100)}
                <div class="mini-bar-item">
                  <span class="mini-bar-val">{w.bookings}</span>
                  <div class="mini-bar-track">
                    <div class="mini-bar-fill" style="height:{pct}%; background: linear-gradient(180deg, #2563eb, #004ac6)"></div>
                  </div>
                  <span class="mini-bar-lbl">{w.week}</span>
                </div>
              {/each}
            </div>
          </div>

          <!-- Top packages -->
          <div class="dual-section">
            <div class="dual-section-label">Top Paket</div>
            <div class="top-pkg-list">
              {#each salesMetrics.topPackages as pkg, i}
                <div class="top-pkg-row">
                  <span class="top-pkg-rank">#{i + 1}</span>
                  <div class="top-pkg-info">
                    <span class="top-pkg-name">{pkg.name}</span>
                    <span class="top-pkg-rev">{formatIDR(pkg.revenue)}</span>
                  </div>
                  <span class="top-pkg-count">{pkg.bookings}</span>
                </div>
              {/each}
            </div>
          </div>

          <a href="/console/bookings" class="dual-view-link">
            Lihat semua booking
            <span class="material-symbols-outlined">arrow_forward</span>
          </a>
        </div>

        <!-- RIGHT: Operations & Logistics Panel -->
        <div class="dual-panel dual-ops">
          <div class="dual-panel-header">
            <span class="material-symbols-outlined">local_shipping</span>
            <h3>Operasional & Logistik</h3>
          </div>

          <!-- Ops KPIs -->
          <div class="dual-kpi-row">
            <div class="dual-kpi">
              <span class="dkpi-val">{opsMetrics.departuresThisWeek}</span>
              <span class="dkpi-lbl">Keberangkatan Minggu Ini</span>
            </div>
            <div class="dual-kpi">
              <span class="dkpi-val">{opsMetrics.paxConfirmed}</span>
              <span class="dkpi-lbl">Pax Confirmed</span>
            </div>
            <div class="dual-kpi">
              <span class="dkpi-val ops-pending">{opsMetrics.pendingTasks}</span>
              <span class="dkpi-lbl">Tugas Pending</span>
            </div>
          </div>

          <!-- Logistics status progress -->
          <div class="dual-section">
            <div class="dual-section-label">Status Persiapan Jamaah</div>
            <div class="logistics-status-list">
              {#each opsMetrics.logisticsStatus as ls}
                {@const pct = Math.round((ls.count / ls.total) * 100)}
                <div class="ls-row">
                  <span class="ls-label">{ls.label}</span>
                  <div class="ls-bar-track">
                    <div class="ls-bar-fill" style="width:{pct}%; background:{ls.color}"></div>
                  </div>
                  <span class="ls-count" style="color:{ls.color}">{ls.count}/{ls.total}</span>
                  <span class="ls-pct">{pct}%</span>
                </div>
              {/each}
            </div>
          </div>

          <!-- Upcoming departures this week -->
          <div class="dual-section">
            <div class="dual-section-label">Keberangkatan Mendatang</div>
            <div class="ops-dep-list">
              {#each opsMetrics.upcomingDepartsThisWeek as dep}
                <div class="ops-dep-row">
                  <div class="ops-dep-info">
                    <span class="ops-dep-name">{dep.name}</span>
                    <span class="ops-dep-date">{formatDate(dep.date)} · {dep.pax} pax</span>
                  </div>
                  <span
                    class="ops-dep-status"
                    class:status-ready={dep.status === 'ready'}
                    class:status-prep={dep.status === 'prep'}
                    class:status-pending={dep.status === 'pending'}
                  >
                    {dep.status === 'ready' ? 'Siap' : dep.status === 'prep' ? 'Persiapan' : 'Pending'}
                  </span>
                </div>
              {/each}
            </div>
          </div>

          <a href="/console/ops" class="dual-view-link ops-link">
            Lihat Ops Board
            <span class="material-symbols-outlined">arrow_forward</span>
          </a>
        </div>
      </div>
    {/if}

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
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
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

  /* ---- BL-DASH-005: View toggle ---- */
  .view-toggle {
    display: inline-flex;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.3rem;
    overflow: hidden;
    background: #f2f4f6;
  }

  .view-toggle-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.35rem 0.65rem;
    border: 0;
    background: transparent;
    font-size: 0.72rem;
    font-weight: 600;
    color: #434655;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.1s, color 0.1s;
  }

  .view-toggle-btn .material-symbols-outlined { font-size: 0.85rem; }

  .view-toggle-btn:hover { background: #e6e8ea; }

  .view-toggle-btn.active {
    background: #2563eb;
    color: #fff;
  }

  /* ---- dual badge ---- */
  .dual-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.25rem 0.55rem;
    background: #ede9fe;
    color: #4c1d95;
    border-radius: 0.2rem;
    font-size: 0.68rem;
    font-weight: 700;
  }
  .dual-badge .material-symbols-outlined { font-size: 0.8rem; }

  /* ---- dual view layout ---- */
  .dual-layout {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
    margin-bottom: 1.5rem;
  }

  @media (max-width: 900px) {
    .dual-layout { grid-template-columns: 1fr; }
  }

  .dual-panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.4rem;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .dual-panel-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.85rem 1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
  }

  .dual-sales .dual-panel-header {
    background: #eff6ff;
  }

  .dual-ops .dual-panel-header {
    background: #f0fdf4;
  }

  .dual-sales .dual-panel-header .material-symbols-outlined {
    font-size: 1.1rem;
    color: #2563eb;
  }

  .dual-ops .dual-panel-header .material-symbols-outlined {
    font-size: 1.1rem;
    color: #065f46;
  }

  .dual-panel-header h3 {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 700;
    color: #191c1e;
  }

  /* dual KPIs row */
  .dual-kpi-row {
    display: flex;
    gap: 0;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
  }

  .dual-kpi {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.1rem;
    padding: 0.85rem 0.5rem;
    border-right: 1px solid rgb(195 198 215 / 0.35);
  }

  .dual-kpi:last-child { border-right: 0; }

  .dkpi-val {
    font-size: 1.25rem;
    font-weight: 800;
    color: #191c1e;
    font-variant-numeric: tabular-nums;
    line-height: 1;
  }

  .dkpi-idr { font-size: 0.78rem; }

  .ops-pending { color: #b45309; }

  .dkpi-lbl {
    font-size: 0.6rem;
    color: #737686;
    text-align: center;
    line-height: 1.3;
    max-width: 6rem;
  }

  /* sections within dual panel */
  .dual-section {
    padding: 0.85rem 1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.25);
  }

  .dual-section:last-of-type { border-bottom: 0; }

  .dual-section-label {
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #737686;
    margin-bottom: 0.65rem;
  }

  /* mini bar chart */
  .mini-bar-chart {
    display: flex;
    align-items: flex-end;
    gap: 0.5rem;
    height: 5rem;
  }

  .mini-bar-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 100%;
    gap: 0.2rem;
  }

  .mini-bar-val { font-size: 0.6rem; color: #434655; font-weight: 700; }

  .mini-bar-track {
    flex: 1;
    width: 100%;
    display: flex;
    align-items: flex-end;
    background: #f2f4f6;
    border-radius: 0.15rem 0.15rem 0 0;
    overflow: hidden;
  }

  .mini-bar-fill {
    width: 100%;
    border-radius: 0.15rem 0.15rem 0 0;
    min-height: 2px;
  }

  .mini-bar-lbl { font-size: 0.6rem; color: #737686; font-weight: 600; }

  /* top packages */
  .top-pkg-list { display: flex; flex-direction: column; gap: 0.5rem; }

  .top-pkg-row {
    display: flex;
    align-items: center;
    gap: 0.55rem;
  }

  .top-pkg-rank {
    font-size: 0.68rem;
    font-weight: 700;
    color: #b0b3c1;
    width: 1.4rem;
    flex-shrink: 0;
  }

  .top-pkg-info { flex: 1; min-width: 0; }

  .top-pkg-name {
    display: block;
    font-size: 0.75rem;
    font-weight: 600;
    color: #191c1e;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .top-pkg-rev { font-size: 0.65rem; color: #065f46; font-weight: 600; }

  .top-pkg-count {
    font-size: 0.82rem;
    font-weight: 800;
    color: #2563eb;
    font-variant-numeric: tabular-nums;
    flex-shrink: 0;
  }

  /* logistics status */
  .logistics-status-list { display: flex; flex-direction: column; gap: 0.55rem; }

  .ls-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .ls-label { font-size: 0.72rem; color: #191c1e; min-width: 8rem; flex-shrink: 0; }

  .ls-bar-track {
    flex: 1;
    height: 0.45rem;
    background: #f2f4f6;
    border-radius: 999px;
    overflow: hidden;
  }

  .ls-bar-fill { height: 100%; border-radius: 999px; min-width: 2px; transition: width 0.3s; }

  .ls-count { font-size: 0.68rem; font-weight: 700; min-width: 3rem; text-align: right; font-variant-numeric: tabular-nums; }

  .ls-pct { font-size: 0.62rem; color: #737686; min-width: 2rem; text-align: right; }

  /* ops departures */
  .ops-dep-list { display: flex; flex-direction: column; gap: 0.5rem; }

  .ops-dep-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .ops-dep-info { flex: 1; min-width: 0; }

  .ops-dep-name {
    display: block;
    font-size: 0.75rem;
    font-weight: 600;
    color: #191c1e;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .ops-dep-date { font-size: 0.65rem; color: #737686; }

  .ops-dep-status {
    font-size: 0.65rem;
    font-weight: 700;
    padding: 0.12rem 0.4rem;
    border-radius: 0.2rem;
    flex-shrink: 0;
  }

  .status-ready { background: #d1fae5; color: #065f46; }
  .status-prep { background: #fef3c7; color: #b45309; }
  .status-pending { background: #fee2e2; color: #991b1b; }

  /* dual view link */
  .dual-view-link {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.6rem 1rem;
    font-size: 0.72rem;
    font-weight: 600;
    color: #2563eb;
    text-decoration: none;
    border-top: 1px solid rgb(195 198 215 / 0.35);
    background: #f7f9fb;
    margin-top: auto;
  }

  .dual-view-link:hover { text-decoration: underline; }
  .dual-view-link .material-symbols-outlined { font-size: 0.85rem; }

  .ops-link { color: #065f46; }
</style>
