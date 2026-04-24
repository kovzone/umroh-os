<script lang="ts">
  import type {
    KpiData,
    BookingStatusCount,
    MonthlyRevenue,
    UpcomingDeparture,
    RecentLead,
    VendorReadiness,
    SeatAvailability,
    CashFlowSummary,
    FinancialReport,
    AdCampaign,
    CsPerformance,
    BusStatus,
    RaudhahStatus,
    LuggageStatus,
    IncidentReport,
    WarehouseHealth,
    LogisticsExecution,
    ArApAging,
    InventoryHealth,
    FulfillmentMonitor,
    DamageReport,
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

  // Wave 6 widget state
  let vendorReadiness = $state<VendorReadiness[]>([]);
  let seatAvailability = $state<SeatAvailability[]>([]);
  let cashFlow = $state<CashFlowSummary>({ period: 'weekly', data: [], net_today: 0, net_week: 0 });
  let financialReport = $state<FinancialReport>({ revenue_month: 0, cogs_month: 0, gross_profit: 0, gross_margin_pct: 0, opex_month: 0, net_profit: 0, net_margin_pct: 0, total_assets: 0, total_liabilities: 0, equity: 0 });
  let adCampaigns = $state<AdCampaign[]>([]);
  let csPerformance = $state<CsPerformance[]>([]);
  let busRadar = $state<BusStatus[]>([]);
  let raudhahStatus = $state<RaudhahStatus[]>([]);
  let luggageTracking = $state<LuggageStatus[]>([]);
  let incidentFeed = $state<IncidentReport[]>([]);
  let warehouseHealth = $state<WarehouseHealth>({ total_stock_value: 0, items_critical: 0, items_reorder: 0, items_ok: 0, categories: [] });
  let logisticsExecution = $state<LogisticsExecution>({ paid_unshipped_count: 0, paid_unshipped_aging_avg_days: 0, grn_backlog: 0, po_backlog: 0, aging_buckets: [] });
  let arApAging = $state<ArApAging>({ ar_current: 0, ar_30: 0, ar_60: 0, ar_90plus: 0, ap_current: 0, ap_30: 0, ap_60: 0, ap_90plus: 0, alerts: [] });
  let inventoryHealth = $state<InventoryHealth>({ total_skus: 0, healthy_pct: 0, low_stock_pct: 0, out_of_stock_pct: 0, top_items: [] });
  let fulfillmentMonitor = $state<FulfillmentMonitor>({ open_pos: 0, overdue_pos: 0, pending_fulfillments: 0, overdue_fulfillments: 0, backlog_items: [] });
  let damageReports = $state<DamageReport[]>([]);

  // BL-DASH-005: dashboard mode
  let dashMode = $state<'operasional' | 'eksekutif'>('operasional');

  // Incident severity filter
  let incidentSeverityFilter = $state<string>('all');

  $effect(() => {
    kpi = data.kpi ?? kpi;
    bookingByStatus = data.bookingByStatus ?? [];
    revenueChart = data.revenueChart ?? [];
    upcomingDepartures = data.upcomingDepartures ?? [];
    recentLeads = data.recentLeads ?? [];
    vendorReadiness = data.vendorReadiness ?? [];
    seatAvailability = data.seatAvailability ?? [];
    cashFlow = data.cashFlow ?? cashFlow;
    financialReport = data.financialReport ?? financialReport;
    adCampaigns = data.adCampaigns ?? [];
    csPerformance = data.csPerformance ?? [];
    busRadar = data.busRadar ?? [];
    raudhahStatus = data.raudhahStatus ?? [];
    luggageTracking = data.luggageTracking ?? [];
    incidentFeed = data.incidentFeed ?? [];
    warehouseHealth = data.warehouseHealth ?? warehouseHealth;
    logisticsExecution = data.logisticsExecution ?? logisticsExecution;
    arApAging = data.arApAging ?? arApAging;
    inventoryHealth = data.inventoryHealth ?? inventoryHealth;
    fulfillmentMonitor = data.fulfillmentMonitor ?? fulfillmentMonitor;
    damageReports = data.damageReports ?? [];
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

  // BL-DASH-005 handled above via dashMode state

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
      <!-- BL-DASH-005: Mode toggle Operasional / Eksekutif -->
      <div class="view-toggle" role="group" aria-label="Pilih mode dashboard">
        <button
          type="button"
          class="view-toggle-btn"
          class:active={dashMode === 'operasional'}
          onclick={() => { dashMode = 'operasional'; }}
          title="Mode Operasional"
        >
          <span class="material-symbols-outlined">settings</span>
          Operasional
        </button>
        <button
          type="button"
          class="view-toggle-btn"
          class:active={dashMode === 'eksekutif'}
          onclick={() => { dashMode = 'eksekutif'; }}
          title="Mode Eksekutif"
        >
          <span class="material-symbols-outlined">leaderboard</span>
          Eksekutif
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
        <h2>{dashMode === 'eksekutif' ? 'Dashboard Eksekutif' : 'Dashboard Operasional'}</h2>
        <p>Ringkasan performa bisnis bulan ini</p>
      </div>
      <span class="dual-badge" class:badge-ops={dashMode === 'operasional'} class:badge-exec={dashMode === 'eksekutif'}>
        <span class="material-symbols-outlined">{dashMode === 'operasional' ? 'settings' : 'leaderboard'}</span>
        {dashMode === 'operasional' ? 'Mode Operasional' : 'Mode Eksekutif'}
      </span>
    </div>

    <!-- ================================================================
         BL-DASH-005: Dual Panel — Sales & Ops side by side (always shown)
    ================================================================= -->
    {#if true}
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

    <!-- ================================================================
         BL-DASH-001: Vendor Execution Readiness (Operasional)
    ================================================================= -->
    {#if dashMode === 'operasional'}
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">task_alt</span>
        <h3>Kesiapan Vendor Keberangkatan</h3>
      </div>
      <div class="panel">
        {#each vendorReadiness as vr}
          {@const pct = vr.checklist_total > 0 ? Math.round((vr.checklist_done / vr.checklist_total) * 100) : 0}
          <div class="w6-vendor-row">
            <div class="w6-vendor-info">
              <span class="w6-vendor-name">{vr.vendor_name}</span>
              <span class="w6-vendor-dep">{vr.departure} · {formatDate(vr.departure_date)}</span>
            </div>
            <div class="w6-vendor-bar-wrap">
              <div class="w6-bar-track">
                <div class="w6-bar-fill" style="width:{pct}%; background:{pct >= 80 ? '#059669' : pct >= 50 ? '#d97706' : '#dc2626'}"></div>
              </div>
              <span class="w6-bar-pct" style="color:{pct >= 80 ? '#059669' : pct >= 50 ? '#d97706' : '#dc2626'}">{vr.checklist_done}/{vr.checklist_total} ({pct}%)</span>
            </div>
          </div>
        {/each}
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-002: Seat Availability (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">airline_seat_recline_normal</span>
        <h3>Ketersediaan Kursi Real-Time</h3>
      </div>
      <div class="w6-seat-grid">
        {#each seatAvailability as sa}
          {@const soldPct = sa.seats_total > 0 ? Math.round((sa.seats_sold / sa.seats_total) * 100) : 0}
          <div class="panel w6-seat-card">
            <div class="w6-seat-name">{sa.package_name}</div>
            <div class="w6-seat-stats">
              <span class="w6-seat-num sold">{sa.seats_sold}<span class="w6-seat-lbl">Terjual</span></span>
              <span class="w6-seat-num avail">{sa.seats_available}<span class="w6-seat-lbl">Tersisa</span></span>
              <span class="w6-seat-num total">{sa.seats_total}<span class="w6-seat-lbl">Total</span></span>
            </div>
            <div class="w6-bar-track" style="height:0.5rem">
              <div class="w6-bar-fill" style="width:{soldPct}%; background:{soldPct >= 90 ? '#dc2626' : soldPct >= 70 ? '#d97706' : '#059669'}"></div>
            </div>
            <div class="w6-seat-foot">{soldPct}% terisi</div>
          </div>
        {/each}
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-008: Live Bus Radar (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">directions_bus</span>
        <h3>Live Bus Radar</h3>
      </div>
      <div class="panel">
        <div class="table-wrap">
          <table>
            <thead><tr><th>Bus</th><th>Rute</th><th>Lokasi Terakhir</th><th>Penumpang</th><th>Status</th></tr></thead>
            <tbody>
              {#each busRadar as bus}
                <tr>
                  <td><b>{bus.plate}</b><br><small style="color:#737686">{bus.driver}</small></td>
                  <td>{bus.route}</td>
                  <td>{bus.last_location}</td>
                  <td>{bus.pax}</td>
                  <td>
                    <span class="w6-badge" class:w6-badge-green={bus.status === 'on_route'} class:w6-badge-amber={bus.status === 'idle'} class:w6-badge-red={bus.status === 'breakdown'} class:w6-badge-gray={bus.status === 'maintenance'}>
                      {bus.status === 'on_route' ? 'Dalam Rute' : bus.status === 'idle' ? 'Standby' : bus.status === 'maintenance' ? 'Maintenance' : 'Rusak'}
                    </span>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-009: Raudhah Status (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">mosque</span>
        <h3>Status Raudhah per Keberangkatan</h3>
      </div>
      <div class="panel">
        {#each raudhahStatus as rs}
          {@const pct = rs.total_jemaah > 0 ? Math.round((rs.entered_raudhah / rs.total_jemaah) * 100) : 0}
          <div class="w6-vendor-row">
            <div class="w6-vendor-info">
              <span class="w6-vendor-name">{rs.departure_name}</span>
              <span class="w6-vendor-dep">{rs.entered_raudhah} dari {rs.total_jemaah} jemaah masuk Raudhah</span>
            </div>
            <div class="w6-vendor-bar-wrap">
              <div class="w6-bar-track"><div class="w6-bar-fill" style="width:{pct}%; background:#059669"></div></div>
              <span class="w6-bar-pct" style="color:#059669">{pct}%</span>
            </div>
          </div>
        {/each}
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-010: Luggage Tracking (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">luggage</span>
        <h3>Tracking Koper Jemaah</h3>
      </div>
      <div class="panel">
        <div class="table-wrap">
          <table>
            <thead><tr><th>Keberangkatan</th><th>Total Koper</th><th>Di Asal</th><th>Transit</th><th>Terkirim</th></tr></thead>
            <tbody>
              {#each luggageTracking as lt}
                {@const delivPct = lt.total_bags > 0 ? Math.round((lt.delivered / lt.total_bags) * 100) : 0}
                <tr>
                  <td>{lt.departure_name}</td>
                  <td>{lt.total_bags}</td>
                  <td>{lt.at_origin}</td>
                  <td>{lt.in_transit}</td>
                  <td><span style="color:{delivPct >= 90 ? '#059669' : '#d97706'}; font-weight:700">{lt.delivered} ({delivPct}%)</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-011: Incident Report Feed (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">warning</span>
        <h3>Feed Insiden</h3>
        <div class="w6-filter-group" role="group">
          {#each ['all','critical','high','medium','low'] as sv}
            <button type="button" class="w6-filter-btn" class:active={incidentSeverityFilter === sv} onclick={() => { incidentSeverityFilter = sv; }}>
              {sv === 'all' ? 'Semua' : sv}
            </button>
          {/each}
        </div>
      </div>
      <div class="panel">
        {#each incidentFeed.filter(i => incidentSeverityFilter === 'all' || i.severity === incidentSeverityFilter) as inc}
          <div class="w6-incident-row">
            <span class="w6-severity w6-sev-{inc.severity}">{inc.severity.toUpperCase()}</span>
            <div class="w6-incident-info">
              <span class="w6-incident-title">{inc.title}</span>
              <span class="w6-incident-meta">{inc.category} · {inc.reported_by} · {relativeTime(inc.reported_at)}</span>
            </div>
            <span class="w6-badge" class:w6-badge-red={inc.status === 'open'} class:w6-badge-amber={inc.status === 'investigating'} class:w6-badge-green={inc.status === 'resolved'}>
              {inc.status === 'open' ? 'Terbuka' : inc.status === 'investigating' ? 'Investigasi' : 'Selesai'}
            </span>
          </div>
        {/each}
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-012: Warehouse Health (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">warehouse</span>
        <h3>Kesehatan Gudang</h3>
      </div>
      <div class="charts-row">
        <div class="panel chart-panel">
          <div class="w6-wh-kpis">
            <div class="w6-wh-kpi"><span class="w6-wh-val">{formatIDR(warehouseHealth.total_stock_value)}</span><span class="w6-wh-lbl">Nilai Stok</span></div>
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:#dc2626">{warehouseHealth.items_critical}</span><span class="w6-wh-lbl">Kritis</span></div>
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:#d97706">{warehouseHealth.items_reorder}</span><span class="w6-wh-lbl">Reorder</span></div>
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:#059669">{warehouseHealth.items_ok}</span><span class="w6-wh-lbl">OK</span></div>
          </div>
        </div>
        <div class="panel chart-panel">
          {#each warehouseHealth.categories as cat}
            {@const maxVal = Math.max(...warehouseHealth.categories.map(c => c.value), 1)}
            {@const catPct = Math.round((cat.value / maxVal) * 100)}
            <div class="w6-wh-cat-row">
              <span class="w6-wh-cat-name">{cat.name}</span>
              <div class="w6-bar-track"><div class="w6-bar-fill" style="width:{catPct}%; background:{cat.status === 'critical' ? '#dc2626' : cat.status === 'reorder' ? '#d97706' : '#059669'}"></div></div>
              <span class="w6-wh-cat-val">{formatIDR(cat.value)}</span>
            </div>
          {/each}
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-013: Logistics Execution Monitor (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">local_shipping</span>
        <h3>Monitor Eksekusi Logistik</h3>
      </div>
      <div class="charts-row">
        <div class="panel chart-panel">
          <div class="w6-wh-kpis">
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:#dc2626">{logisticsExecution.paid_unshipped_count}</span><span class="w6-wh-lbl">Dibayar Belum Kirim</span></div>
            <div class="w6-wh-kpi"><span class="w6-wh-val">{logisticsExecution.paid_unshipped_aging_avg_days}h</span><span class="w6-wh-lbl">Rata-rata Aging</span></div>
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:#d97706">{logisticsExecution.grn_backlog}</span><span class="w6-wh-lbl">Backlog GRN</span></div>
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:#d97706">{logisticsExecution.po_backlog}</span><span class="w6-wh-lbl">Backlog PO</span></div>
          </div>
        </div>
        <div class="panel chart-panel">
          <div class="dual-section-label" style="font-size:0.7rem;font-weight:700;color:#737686;margin-bottom:0.7rem">Aging Kiriman Tertunda</div>
          <div class="bar-chart" style="height:6rem">
            {#each logisticsExecution.aging_buckets as bucket}
              {@const maxCount = Math.max(...logisticsExecution.aging_buckets.map(b => b.count), 1)}
              {@const bPct = Math.round((bucket.count / maxCount) * 100)}
              <div class="bar-item">
                <span class="bar-value" style="font-size:0.72rem">{bucket.count}</span>
                <div class="bar-track"><div class="bar-fill" style="height:{bPct}%; background:{bucket.color}"></div></div>
                <span class="bar-label" style="font-size:0.6rem">{bucket.label}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-015: Inventory Health (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">inventory_2</span>
        <h3>Kesehatan Inventori</h3>
      </div>
      <div class="charts-row">
        <div class="panel chart-panel">
          <div class="w6-inv-donut">
            <div class="w6-inv-stat"><span class="w6-inv-num" style="color:#059669">{inventoryHealth.healthy_pct}%</span><span class="w6-inv-lbl">Sehat</span></div>
            <div class="w6-inv-stat"><span class="w6-inv-num" style="color:#d97706">{inventoryHealth.low_stock_pct}%</span><span class="w6-inv-lbl">Stok Rendah</span></div>
            <div class="w6-inv-stat"><span class="w6-inv-num" style="color:#dc2626">{inventoryHealth.out_of_stock_pct}%</span><span class="w6-inv-lbl">Habis</span></div>
            <div class="w6-inv-stat"><span class="w6-inv-num">{inventoryHealth.total_skus}</span><span class="w6-inv-lbl">Total SKU</span></div>
          </div>
        </div>
        <div class="panel chart-panel">
          <div class="dual-section-label" style="font-size:0.7rem;font-weight:700;color:#737686;margin-bottom:0.7rem">Item Perlu Perhatian</div>
          {#each inventoryHealth.top_items as item}
            <div class="w6-inv-item-row">
              <span class="w6-inv-item-name">{item.name}</span>
              <span class="w6-badge" class:w6-badge-green={item.status === 'healthy'} class:w6-badge-amber={item.status === 'low'} class:w6-badge-red={item.status === 'out'}>{item.status === 'healthy' ? 'OK' : item.status === 'low' ? 'Rendah' : 'Habis'}</span>
              <span style="font-size:0.72rem;font-weight:700;min-width:2.5rem;text-align:right">{item.qty}</span>
            </div>
          {/each}
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-016: Fulfillment & PO Monitor (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">assignment</span>
        <h3>Monitor Fulfillment & PO</h3>
      </div>
      <div class="w6-fulfil-kpis">
        <div class="panel w6-fulfil-kpi"><span class="w6-wh-val">{fulfillmentMonitor.open_pos}</span><span class="w6-wh-lbl">PO Terbuka</span></div>
        <div class="panel w6-fulfil-kpi"><span class="w6-wh-val" style="color:#dc2626">{fulfillmentMonitor.overdue_pos}</span><span class="w6-wh-lbl">PO Jatuh Tempo</span></div>
        <div class="panel w6-fulfil-kpi"><span class="w6-wh-val">{fulfillmentMonitor.pending_fulfillments}</span><span class="w6-wh-lbl">Fulfillment Pending</span></div>
        <div class="panel w6-fulfil-kpi"><span class="w6-wh-val" style="color:#d97706">{fulfillmentMonitor.overdue_fulfillments}</span><span class="w6-wh-lbl">Fulfillment Terlambat</span></div>
      </div>
      <div class="panel" style="margin-top:1rem">
        <div class="table-wrap">
          <table>
            <thead><tr><th>No. PO</th><th>Vendor</th><th>Due Date</th><th>Status</th><th class="align-right">Nilai</th></tr></thead>
            <tbody>
              {#each fulfillmentMonitor.backlog_items as item}
                <tr>
                  <td><b>{item.po_number}</b></td>
                  <td>{item.vendor}</td>
                  <td>{formatDate(item.due_date)}</td>
                  <td><span class="w6-badge" class:w6-badge-red={item.status === 'overdue'} class:w6-badge-amber={item.status === 'pending'}>{item.status === 'overdue' ? 'Terlambat' : 'Pending'}</span></td>
                  <td class="align-right">{formatIDR(item.amount)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-017: Damage Report (Operasional)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">broken_image</span>
        <h3>Laporan Kerusakan Barang</h3>
      </div>
      <div class="panel">
        <div class="table-wrap">
          <table>
            <thead><tr><th>Item</th><th>Jml Rusak</th><th>Lokasi</th><th>Dilaporkan Oleh</th><th>Waktu</th><th class="align-right">Est. Kerugian</th><th>Status</th></tr></thead>
            <tbody>
              {#each damageReports as dmg}
                <tr>
                  <td>{dmg.item_name}</td>
                  <td>{dmg.qty_damaged}</td>
                  <td>{dmg.location}</td>
                  <td>{dmg.reported_by}</td>
                  <td>{relativeTime(dmg.reported_at)}</td>
                  <td class="align-right">{formatIDR(dmg.estimated_loss)}</td>
                  <td><span class="w6-badge" class:w6-badge-amber={dmg.status === 'pending'} class:w6-badge-gray={dmg.status === 'reviewed'} class:w6-badge-red={dmg.status === 'written_off'}>{dmg.status === 'pending' ? 'Pending' : dmg.status === 'reviewed' ? 'Ditinjau' : 'Dihapusbuku'}</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>
    {/if}

    <!-- ================================================================
         BL-DASH-003: Cash Flow Widget (Eksekutif)
    ================================================================= -->
    {#if dashMode === 'eksekutif'}
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">account_balance_wallet</span>
        <h3>Arus Kas Instan</h3>
      </div>
      <div class="charts-row">
        <div class="panel chart-panel">
          <div class="w6-wh-kpis">
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:{cashFlow.net_today >= 0 ? '#059669' : '#dc2626'}">{formatIDR(cashFlow.net_today)}</span><span class="w6-wh-lbl">Net Hari Ini</span></div>
            <div class="w6-wh-kpi"><span class="w6-wh-val" style="color:{cashFlow.net_week >= 0 ? '#059669' : '#dc2626'}">{formatIDR(cashFlow.net_week)}</span><span class="w6-wh-lbl">Net Minggu Ini</span></div>
          </div>
        </div>
        <div class="panel chart-panel">
          <div class="dual-section-label" style="font-size:0.7rem;font-weight:700;color:#737686;margin-bottom:0.7rem">Cash In vs Out (7 hari)</div>
          {@const maxCF = Math.max(...cashFlow.data.map(d => Math.max(d.cash_in, d.cash_out)), 1)}
          <div class="bar-chart" style="height:7rem; gap:0.35rem">
            {#each cashFlow.data as day}
              {@const inPct = Math.round((day.cash_in / maxCF) * 100)}
              {@const outPct = Math.round((day.cash_out / maxCF) * 100)}
              <div class="bar-item" style="gap:0.1rem">
                <div class="mini-bar-track" style="height:100%; background:transparent; gap:2px; flex-direction:row; align-items:flex-end">
                  <div style="width:50%; height:{inPct}%; background:#059669; border-radius:2px 2px 0 0; min-height:2px"></div>
                  <div style="width:50%; height:{outPct}%; background:#dc2626; border-radius:2px 2px 0 0; min-height:2px"></div>
                </div>
                <span class="bar-label" style="font-size:0.58rem">{day.label}</span>
              </div>
            {/each}
          </div>
          <div style="display:flex;gap:1rem;margin-top:0.5rem">
            <span style="font-size:0.65rem;color:#059669">▌ Cash In</span>
            <span style="font-size:0.65rem;color:#dc2626">▌ Cash Out</span>
          </div>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-004: Executive Financial Report (Eksekutif)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">analytics</span>
        <h3>Laporan Keuangan Eksekutif</h3>
      </div>
      <div class="charts-row">
        <!-- P&L Summary -->
        <div class="panel chart-panel">
          <div class="dual-section-label" style="font-size:0.7rem;font-weight:700;color:#737686;margin-bottom:0.75rem">Laba Rugi Bulan Ini</div>
          <div class="w6-fin-row"><span>Pendapatan</span><span style="color:#059669;font-weight:700">{formatIDR(financialReport.revenue_month)}</span></div>
          <div class="w6-fin-row"><span>HPP</span><span style="color:#dc2626">({formatIDR(financialReport.cogs_month)})</span></div>
          <div class="w6-fin-divider"></div>
          <div class="w6-fin-row"><span>Laba Kotor</span><span style="font-weight:800">{formatIDR(financialReport.gross_profit)} <small style="color:#737686">({financialReport.gross_margin_pct}%)</small></span></div>
          <div class="w6-fin-row"><span>Beban Operasional</span><span style="color:#dc2626">({formatIDR(financialReport.opex_month)})</span></div>
          <div class="w6-fin-divider"></div>
          <div class="w6-fin-row w6-fin-net"><span>Laba Bersih</span><span style="color:{financialReport.net_profit >= 0 ? '#059669' : '#dc2626'};font-size:1.1rem;font-weight:800">{formatIDR(financialReport.net_profit)} <small>({financialReport.net_margin_pct}%)</small></span></div>
        </div>
        <!-- Balance Sheet Highlight -->
        <div class="panel chart-panel">
          <div class="dual-section-label" style="font-size:0.7rem;font-weight:700;color:#737686;margin-bottom:0.75rem">Neraca Ringkas</div>
          <div class="w6-fin-row"><span>Total Aset</span><span style="font-weight:700">{formatIDR(financialReport.total_assets)}</span></div>
          <div class="w6-fin-row"><span>Total Kewajiban</span><span style="color:#dc2626">{formatIDR(financialReport.total_liabilities)}</span></div>
          <div class="w6-fin-divider"></div>
          <div class="w6-fin-row w6-fin-net"><span>Ekuitas</span><span style="color:#059669;font-weight:800">{formatIDR(financialReport.equity)}</span></div>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-006: Ad Budget Monitor (Eksekutif)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">campaign</span>
        <h3>Monitor Anggaran Iklan</h3>
      </div>
      <div class="panel">
        <div class="table-wrap">
          <table>
            <thead><tr><th>Kampanye</th><th>Channel</th><th>Spend / Budget</th><th>Closing</th><th class="align-right">CPL</th><th class="align-right">CPA</th></tr></thead>
            <tbody>
              {#each adCampaigns as camp}
                {@const spendPct = camp.budget > 0 ? Math.round((camp.spend / camp.budget) * 100) : 0}
                <tr>
                  <td>{camp.campaign}</td>
                  <td>{camp.channel}</td>
                  <td>
                    <div style="display:flex;align-items:center;gap:0.4rem">
                      <div class="w6-bar-track" style="width:5rem;flex-shrink:0"><div class="w6-bar-fill" style="width:{Math.min(spendPct,100)}%; background:{spendPct >= 90 ? '#dc2626' : '#2563eb'}"></div></div>
                      <span style="font-size:0.68rem">{formatIDR(camp.spend)} / {formatIDR(camp.budget)}</span>
                    </div>
                  </td>
                  <td>{camp.closings}</td>
                  <td class="align-right">{formatIDR(camp.cpl)}</td>
                  <td class="align-right">{formatIDR(camp.cpa)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-007: CS Performance Board (Eksekutif)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">support_agent</span>
        <h3>Papan Performa CS</h3>
      </div>
      <div class="panel">
        <div class="table-wrap">
          <table>
            <thead><tr><th>CS</th><th>Lead Dihandle</th><th>Konversi</th><th>Tingkat Konversi</th><th class="align-right">Avg Respon (mnt)</th><th class="align-right">Breach SLA</th></tr></thead>
            <tbody>
              {#each csPerformance.sort((a, b) => b.converted - a.converted) as cs, i}
                <tr>
                  <td><b>{i + 1}. {cs.cs_name}</b></td>
                  <td>{cs.leads_handled}</td>
                  <td>{cs.converted}</td>
                  <td>
                    <div style="display:flex;align-items:center;gap:0.4rem">
                      <div class="w6-bar-track" style="width:4rem"><div class="w6-bar-fill" style="width:{cs.conversion_pct}%; background:#059669"></div></div>
                      <span style="font-size:0.68rem;font-weight:700;color:#059669">{cs.conversion_pct}%</span>
                    </div>
                  </td>
                  <td class="align-right" style="color:{cs.avg_response_min > 30 ? '#dc2626' : '#059669'}">{cs.avg_response_min} mnt</td>
                  <td class="align-right" style="color:{cs.sla_breach > 0 ? '#dc2626' : '#059669'}">{cs.sla_breach}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- ================================================================
         BL-DASH-014: AR/AP Aging Snapshot (Eksekutif)
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">account_balance</span>
        <h3>Snapshot Likuiditas AR/AP</h3>
      </div>
      <div class="charts-row">
        <div class="panel chart-panel">
          <div class="dual-section-label" style="font-size:0.7rem;font-weight:700;color:#737686;margin-bottom:0.75rem">Piutang Usaha (AR)</div>
          {#each [['Lancar', arApAging.ar_current, '#059669'], ['1–30 hari', arApAging.ar_30, '#d97706'], ['31–60 hari', arApAging.ar_60, '#f97316'], ['60+ hari', arApAging.ar_90plus, '#dc2626']] as [label, val, color]}
            <div class="w6-aging-row"><span class="w6-aging-lbl">{label}</span><span class="w6-aging-val" style="color:{color}">{formatIDR(val as number)}</span></div>
          {/each}
        </div>
        <div class="panel chart-panel">
          <div class="dual-section-label" style="font-size:0.7rem;font-weight:700;color:#737686;margin-bottom:0.75rem">Hutang Usaha (AP)</div>
          {#each [['Lancar', arApAging.ap_current, '#059669'], ['1–30 hari', arApAging.ap_30, '#d97706'], ['31–60 hari', arApAging.ap_60, '#f97316'], ['60+ hari', arApAging.ap_90plus, '#dc2626']] as [label, val, color]}
            <div class="w6-aging-row"><span class="w6-aging-lbl">{label}</span><span class="w6-aging-val" style="color:{color}">{formatIDR(val as number)}</span></div>
          {/each}
          {#if arApAging.alerts.length > 0}
            <div style="margin-top:0.85rem;display:flex;flex-direction:column;gap:0.4rem">
              {#each arApAging.alerts as alert}
                <div class="w6-alert-row">
                  <span class="w6-badge" class:w6-badge-amber={alert.type === 'ap'} class:w6-badge-red={alert.type === 'ar'}>{alert.type.toUpperCase()}</span>
                  <span style="font-size:0.72rem;flex:1">{alert.message}</span>
                  <span style="font-size:0.7rem;font-weight:700;color:#dc2626">{formatIDR(alert.amount)}</span>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    </section>
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

  /* ---- BL-DASH-005: mode badge ---- */
  .badge-ops { background: #dbeafe; color: #1e3a8a; }
  .badge-exec { background: #ede9fe; color: #4c1d95; }

  /* ---- Wave 6 widget shared ---- */
  .w6-bar-track {
    flex: 1;
    height: 0.4rem;
    background: #f2f4f6;
    border-radius: 999px;
    overflow: hidden;
  }

  .w6-bar-fill {
    height: 100%;
    border-radius: 999px;
    min-width: 2px;
    transition: width 0.3s;
  }

  .w6-badge {
    display: inline-flex;
    align-items: center;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .w6-badge-green { background: #d1fae5; color: #065f46; }
  .w6-badge-amber { background: #fef3c7; color: #b45309; }
  .w6-badge-red { background: #fee2e2; color: #991b1b; }
  .w6-badge-gray { background: #f3f4f6; color: #374151; }

  /* BL-DASH-001: vendor readiness */
  .w6-vendor-row {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
  }

  .w6-vendor-row:last-child { border-bottom: 0; }

  .w6-vendor-info { flex: 1; min-width: 0; }

  .w6-vendor-name {
    display: block;
    font-size: 0.82rem;
    font-weight: 700;
    color: #191c1e;
  }

  .w6-vendor-dep { font-size: 0.68rem; color: #737686; }

  .w6-vendor-bar-wrap {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    min-width: 14rem;
  }

  .w6-bar-pct { font-size: 0.7rem; font-weight: 700; white-space: nowrap; }

  /* BL-DASH-002: seat availability */
  .w6-seat-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(14rem, 1fr));
    gap: 1rem;
    margin-bottom: 0;
  }

  .w6-seat-card { padding: 1rem; display: flex; flex-direction: column; gap: 0.65rem; }

  .w6-seat-name {
    font-size: 0.8rem;
    font-weight: 700;
    color: #191c1e;
    line-height: 1.3;
  }

  .w6-seat-stats {
    display: flex;
    gap: 1rem;
  }

  .w6-seat-num {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: 1.25rem;
    font-weight: 800;
    font-variant-numeric: tabular-nums;
  }

  .w6-seat-num.sold { color: #2563eb; }
  .w6-seat-num.avail { color: #059669; }
  .w6-seat-num.total { color: #737686; }

  .w6-seat-lbl { font-size: 0.6rem; font-weight: 500; color: #737686; }

  .w6-seat-foot { font-size: 0.65rem; color: #737686; }

  /* BL-DASH-011: incident feed */
  .w6-filter-group {
    display: flex;
    gap: 0.25rem;
    margin-left: auto;
  }

  .w6-filter-btn {
    border: 1px solid rgb(195 198 215 / 0.7);
    background: #fff;
    border-radius: 0.2rem;
    padding: 0.2rem 0.55rem;
    font-size: 0.65rem;
    font-weight: 600;
    color: #434655;
    cursor: pointer;
    text-transform: capitalize;
  }

  .w6-filter-btn.active {
    background: #2563eb;
    border-color: #2563eb;
    color: #fff;
  }

  .w6-incident-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.65rem 1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
  }

  .w6-incident-row:last-child { border-bottom: 0; }

  .w6-severity {
    font-size: 0.6rem;
    font-weight: 800;
    padding: 0.18rem 0.4rem;
    border-radius: 0.2rem;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .w6-sev-critical { background: #fce7f3; color: #9d174d; }
  .w6-sev-high { background: #fee2e2; color: #991b1b; }
  .w6-sev-medium { background: #fef3c7; color: #b45309; }
  .w6-sev-low { background: #d1fae5; color: #065f46; }

  .w6-incident-info { flex: 1; min-width: 0; }

  .w6-incident-title {
    display: block;
    font-size: 0.8rem;
    font-weight: 600;
    color: #191c1e;
  }

  .w6-incident-meta { font-size: 0.65rem; color: #737686; }

  /* BL-DASH-012: warehouse health */
  .w6-wh-kpis {
    display: flex;
    gap: 0;
    flex-wrap: wrap;
  }

  .w6-wh-kpi {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0.85rem 0.5rem;
    border-right: 1px solid rgb(195 198 215 / 0.35);
    min-width: 6rem;
  }

  .w6-wh-kpi:last-child { border-right: 0; }

  .w6-wh-val {
    font-size: 1.15rem;
    font-weight: 800;
    color: #191c1e;
    font-variant-numeric: tabular-nums;
    line-height: 1;
  }

  .w6-wh-lbl {
    font-size: 0.62rem;
    color: #737686;
    text-align: center;
    margin-top: 0.2rem;
  }

  .w6-wh-cat-row {
    display: flex;
    align-items: center;
    gap: 0.65rem;
    margin-bottom: 0.55rem;
  }

  .w6-wh-cat-name { font-size: 0.72rem; min-width: 10rem; color: #434655; }

  .w6-wh-cat-val { font-size: 0.68rem; font-weight: 700; min-width: 6rem; text-align: right; color: #191c1e; }

  /* BL-DASH-015: inventory health */
  .w6-inv-donut {
    display: flex;
    gap: 0;
    flex-wrap: wrap;
  }

  .w6-inv-stat {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0.85rem 0.25rem;
    border-right: 1px solid rgb(195 198 215 / 0.35);
    min-width: 5rem;
  }

  .w6-inv-stat:last-child { border-right: 0; }

  .w6-inv-num { font-size: 1.35rem; font-weight: 800; line-height: 1; }

  .w6-inv-lbl { font-size: 0.62rem; color: #737686; margin-top: 0.2rem; }

  .w6-inv-item-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.4rem 0;
    border-bottom: 1px solid rgb(195 198 215 / 0.25);
  }

  .w6-inv-item-row:last-child { border-bottom: 0; }

  .w6-inv-item-name { flex: 1; font-size: 0.75rem; color: #191c1e; }

  /* BL-DASH-016: fulfillment monitor */
  .w6-fulfil-kpis {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: 0;
  }

  .w6-fulfil-kpi {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 1rem 0.5rem;
  }

  @media (max-width: 860px) {
    .w6-fulfil-kpis { grid-template-columns: repeat(2, 1fr); }
    .w6-seat-grid { grid-template-columns: 1fr; }
  }

  /* BL-DASH-003/004: finance rows */
  .w6-fin-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.4rem 0;
    font-size: 0.8rem;
    color: #191c1e;
    border-bottom: 1px solid rgb(195 198 215 / 0.2);
  }

  .w6-fin-divider { height: 1px; background: rgb(195 198 215 / 0.6); margin: 0.25rem 0; }

  .w6-fin-net { font-size: 0.9rem; font-weight: 700; border-bottom: 0; }

  /* BL-DASH-014: aging */
  .w6-aging-row {
    display: flex;
    justify-content: space-between;
    padding: 0.4rem 0;
    font-size: 0.8rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.2);
  }

  .w6-aging-lbl { color: #434655; }

  .w6-aging-val { font-weight: 700; }

  .w6-alert-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.35rem 0.5rem;
    background: #fff7ed;
    border-radius: 0.2rem;
    border: 1px solid #fed7aa;
  }
</style>
