<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  // ---- state ----
  let summary = $state<any>(null);
  let cashflow = $state<any>(null);
  let loading = $state(false);
  let error = $state('');
  let tab = $state<'week' | 'month' | 'quarter'>('month');
  let lastRefresh = $state('');

  function getDateRange(t: string) {
    const now = new Date();
    const end = now.toISOString().slice(0, 10);
    let start: string;
    if (t === 'week') {
      const d = new Date(now); d.setDate(d.getDate() - 7);
      start = d.toISOString().slice(0, 10);
    } else if (t === 'month') {
      const d = new Date(now); d.setDate(1);
      start = d.toISOString().slice(0, 10);
    } else {
      const d = new Date(now); d.setMonth(d.getMonth() - 3);
      start = d.toISOString().slice(0, 10);
    }
    return { start, end };
  }

  async function fetchData() {
    loading = true; error = '';
    const { start, end } = getDateRange(tab);
    try {
      const [s, c] = await Promise.all([
        fetch(`${GATEWAY}/v1/finance/realtime-summary`),
        fetch(`${GATEWAY}/v1/finance/cashflow?start_date=${start}&end_date=${end}`)
      ]);
      summary = s.ok ? await s.json() : null;
      cashflow = c.ok ? await c.json() : null;
      lastRefresh = new Date().toLocaleTimeString('id-ID');
    } catch {
      error = 'Gagal memuat data. Menampilkan data kosong.';
      summary = null; cashflow = null;
    }
    loading = false;
  }

  $effect(() => {
    fetchData();
    const timer = setInterval(fetchData, 30000);
    return () => clearInterval(timer);
  });

  $effect(() => {
    // re-fetch when tab changes
    tab;
    fetchData();
  });

  const kpiCards = $derived([
    { label: 'Total Pendapatan', value: summary?.total_revenue ?? 0, icon: 'trending_up', color: '#059669' },
    { label: 'Total Beban', value: summary?.total_expense ?? 0, icon: 'trending_down', color: '#dc2626' },
    { label: 'Laba Bersih', value: summary?.net_profit ?? 0, icon: 'savings', color: '#2563eb' },
    { label: 'Saldo Kas', value: summary?.cash_balance ?? 0, icon: 'account_balance_wallet', color: '#7c3aed' },
    { label: 'Piutang (AR)', value: summary?.ar_balance ?? 0, icon: 'arrow_circle_up', color: '#d97706' },
    { label: 'Utang (AP)', value: summary?.ap_balance ?? 0, icon: 'arrow_circle_down', color: '#db2777' },
  ]);

  const cfCards = $derived([
    { label: 'Operasional', value: cashflow?.operational_net ?? 0 },
    { label: 'Investasi', value: cashflow?.investment_net ?? 0 },
    { label: 'Pendanaan', value: cashflow?.financing_net ?? 0 },
  ]);

  const cfLines = $derived(Array.isArray(cashflow?.lines) ? cashflow.lines : []);
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Dashboard Real-time</span>
    </nav>
    <div class="top-actions">
      {#if lastRefresh}
        <span class="refresh-info">Refresh: {lastRefresh}</span>
      {/if}
      <button class="btn-ghost" onclick={fetchData} disabled={loading}>
        <span class="material-symbols-outlined" class:spin={loading}>refresh</span>
        Muat Ulang
      </button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Dashboard Keuangan Real-time</h2>
      <p>Auto-refresh setiap 30 detik</p>
    </div>

    {#if error}
      <div class="alert-err">{error}</div>
    {/if}

    <!-- KPI Cards -->
    {#if loading && !summary}
      <div class="kpi-grid">
        {#each {length: 6} as _}
          <div class="kpi-card skeleton"></div>
        {/each}
      </div>
    {:else}
      <div class="kpi-grid">
        {#each kpiCards as card}
          <div class="kpi-card">
            <div class="kpi-top">
              <span class="kpi-label">{card.label}</span>
              <span class="material-symbols-outlined kpi-icon" style="color:{card.color}">{card.icon}</span>
            </div>
            <div class="kpi-value" style="color:{card.color}">{formatIDR(card.value)}</div>
          </div>
        {/each}
      </div>
    {/if}

    <!-- Cash Flow Section -->
    <div class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">waterfall_chart</span>
        <h3>Arus Kas</h3>
        <div class="tab-group" style="margin-left:auto">
          {#each (['week','month','quarter'] as const) as t}
            <button class="tab-btn" class:active={tab===t} onclick={() => tab = t}>
              {t === 'week' ? 'Minggu' : t === 'month' ? 'Bulan' : 'Kuartal'}
            </button>
          {/each}
        </div>
      </div>

      <div class="cf-cards">
        {#each cfCards as c}
          <div class="cf-card">
            <div class="cf-label">{c.label}</div>
            <div class="cf-val" class:pos={c.value >= 0} class:neg={c.value < 0}>{formatIDR(c.value)}</div>
          </div>
        {/each}
      </div>

      {#if cfLines.length > 0}
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>Tanggal</th><th>Deskripsi</th><th>Kategori</th><th class="ar">Masuk</th><th class="ar">Keluar</th><th class="ar">Saldo</th></tr>
            </thead>
            <tbody>
              {#each cfLines as line}
                <tr>
                  <td>{line.date ?? '-'}</td>
                  <td>{line.description ?? '-'}</td>
                  <td><span class="chip">{line.category ?? '-'}</span></td>
                  <td class="ar amount-pos">{line.inflow > 0 ? formatIDR(line.inflow) : '—'}</td>
                  <td class="ar amount-neg">{line.outflow > 0 ? formatIDR(line.outflow) : '—'}</td>
                  <td class="ar">{formatIDR(line.balance ?? 0)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {:else if !loading}
        <div class="empty">Data arus kas tidak tersedia</div>
      {/if}
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar {
    position: sticky; top: 0; z-index: 30; height: 4rem;
    background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center; justify-content: space-between;
    backdrop-filter: blur(8px);
  }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .top-actions { display: flex; align-items: center; gap: 0.75rem; }
  .refresh-info { font-size: 0.72rem; color: #737686; }
  .btn-ghost {
    display: inline-flex; align-items: center; gap: 0.3rem;
    border: 1px solid rgb(195 198 215 / 0.55); background: #fff;
    border-radius: 0.25rem; padding: 0.4rem 0.7rem; font-size: 0.78rem;
    font-weight: 600; cursor: pointer; font-family: inherit; color: #191c1e;
  }
  .btn-ghost:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-ghost .material-symbols-outlined { font-size: 1rem; }

  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }

  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }

  .kpi-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; margin-bottom: 1.5rem; }
  .kpi-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1rem 1.1rem; }
  .kpi-card.skeleton { height: 5rem; background: linear-gradient(90deg, #f2f4f6 25%, #e8eaec 50%, #f2f4f6 75%); background-size: 200%; animation: shimmer 1.2s infinite; border-radius: 0.5rem; }
  @keyframes shimmer { from { background-position: 200% 0; } to { background-position: -200% 0; } }
  .kpi-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.5rem; }
  .kpi-label { font-size: 0.72rem; color: #737686; font-weight: 600; }
  .kpi-icon { font-size: 1.1rem; }
  .kpi-value { font-size: 1.1rem; font-weight: 700; font-variant-numeric: tabular-nums; }

  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; margin-bottom: 1.5rem; }
  .section-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 1rem; }
  .section-header h3 { margin: 0; font-size: 0.95rem; font-weight: 700; }
  .section-icon { font-size: 1.1rem; color: #004ac6; }

  .tab-group { display: flex; gap: 0.35rem; }
  .tab-btn { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.3rem 0.65rem; font-size: 0.75rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #737686; }
  .tab-btn.active { background: #2563eb; color: #fff; border-color: #2563eb; }

  .cf-cards { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0.75rem; margin-bottom: 1rem; }
  .cf-card { background: #f7f9fb; border: 1px solid rgb(195 198 215 / 0.35); border-radius: 0.35rem; padding: 0.85rem 1rem; }
  .cf-label { font-size: 0.72rem; color: #737686; font-weight: 600; margin-bottom: 0.3rem; }
  .cf-val { font-size: 1rem; font-weight: 700; font-variant-numeric: tabular-nums; }
  .cf-val.pos { color: #059669; }
  .cf-val.neg { color: #dc2626; }

  .table-wrap { overflow-x: auto; max-height: 320px; overflow-y: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; position: sticky; top: 0; }
  tbody tr:hover { background: #f7f9fb; }
  .ar { text-align: right; }
  .amount-pos { color: #059669; font-variant-numeric: tabular-nums; }
  .amount-neg { color: #dc2626; font-variant-numeric: tabular-nums; }
  .chip { display: inline-flex; padding: 0.1rem 0.4rem; background: #e0f2fe; color: #075985; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .empty { text-align: center; color: #b0b3c1; padding: 2rem; font-size: 0.82rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; }
</style>
