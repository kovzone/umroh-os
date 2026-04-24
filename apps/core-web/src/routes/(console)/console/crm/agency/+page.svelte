<script lang="ts">
  interface Agent {
    id: string;
    name: string;
    tier: string;
    status: 'active' | 'pending_kyc' | 'pending_mou' | 'registered';
    commission_total: number;
    joined: string;
  }

  const kpis = $state([
    { label: 'Agen Aktif', value: '142', icon: 'storefront', color: '#059669', bg: '#d1fae5' },
    { label: 'Menunggu KYC', value: '18', icon: 'badge', color: '#d97706', bg: '#fef3c7' },
    { label: 'Menunggu MoU', value: '9', icon: 'handshake', color: '#7c3aed', bg: '#ede9fe' },
    { label: 'Total Komisi Dibayar', value: 'Rp 847 jt', icon: 'account_balance_wallet', color: '#2563eb', bg: '#dbeafe' }
  ]);

  const navCards = [
    { href: '/console/crm/agency/onboarding', icon: 'how_to_reg', title: 'Onboarding', desc: 'Pipeline registrasi, KYC & MoU agen' },
    { href: '/console/crm/agency/replica-site', icon: 'language', title: 'Replica Site', desc: 'Kelola landing page personal agen' },
    { href: '/console/crm/agency/content', icon: 'perm_media', title: 'Konten & Pemasaran', desc: 'Bank konten, flyer, UTM, kalender' },
    { href: '/console/crm/agency/leads', icon: 'person_search', title: 'CRM Leads', desc: 'Leads, drip sequence, broadcast' },
    { href: '/console/crm/agency/commission', icon: 'account_balance_wallet', title: 'Komisi & Wallet', desc: 'Saldo, payout, override komisi' },
    { href: '/console/crm/agency/academy', icon: 'school', title: 'Akademi Agen', desc: 'Kursus, leaderboard, skrip jual' },
    { href: '/console/crm/agency/team', icon: 'group', title: 'Tim & Downline', desc: 'Struktur downline & tier leveling' }
  ];

  let recentAgents = $state<Agent[]>([
    { id: 'a1', name: 'PT Mitra Barokah Tours', tier: 'Gold', status: 'active', commission_total: 45_000_000, joined: '2026-01-15' },
    { id: 'a2', name: 'CV Cahaya Umroh', tier: 'Silver', status: 'active', commission_total: 22_500_000, joined: '2026-02-03' },
    { id: 'a3', name: 'Ustadz Fahmi Network', tier: 'Bronze', status: 'pending_kyc', commission_total: 0, joined: '2026-04-10' },
    { id: 'a4', name: 'PT Andalus Wisata', tier: 'Silver', status: 'pending_mou', commission_total: 0, joined: '2026-04-12' },
    { id: 'a5', name: 'Rina Hadiati Agency', tier: 'Bronze', status: 'registered', commission_total: 0, joined: '2026-04-20' },
    { id: 'a6', name: 'Karima Travel Group', tier: 'Gold', status: 'active', commission_total: 67_800_000, joined: '2025-11-22' }
  ]);

  const STATUS_LABELS: Record<string, string> = {
    active: 'Aktif',
    pending_kyc: 'Menunggu KYC',
    pending_mou: 'Menunggu MoU',
    registered: 'Terdaftar'
  };

  function fmtCurrency(val: number): string {
    if (val === 0) return '—';
    return 'Rp ' + (val / 1_000_000).toFixed(1) + ' jt';
  }

  function fmtDate(iso: string): string {
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">hub</span>
      <span class="sep">/</span>
      <a href="/console/crm" class="bc-link">CRM Tools</a>
      <span class="sep">/</span>
      <span class="topbar-current">Portal Agen</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Portal Agen</h2>
        <p>Pusat manajemen agen, onboarding, komisi, dan pemasaran</p>
      </div>
      <a href="/console/crm/agency/onboarding" class="primary-btn">
        <span class="material-symbols-outlined">person_add</span>
        Daftarkan Agen Baru
      </a>
    </div>

    <!-- KPI Cards -->
    <div class="kpi-grid">
      {#each kpis as kpi}
        <div class="kpi-card">
          <div class="kpi-icon" style="background:{kpi.bg};color:{kpi.color}">
            <span class="material-symbols-outlined">{kpi.icon}</span>
          </div>
          <div class="kpi-body">
            <div class="kpi-value">{kpi.value}</div>
            <div class="kpi-label">{kpi.label}</div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Navigation Cards -->
    <div class="section-block">
      <div class="section-title">
        <span class="material-symbols-outlined">grid_view</span>
        Sub-Modul
      </div>
      <div class="nav-grid">
        {#each navCards as card}
          <a href={card.href} class="nav-card">
            <span class="material-symbols-outlined nav-icon">{card.icon}</span>
            <div class="nav-text">
              <div class="nav-label">{card.title}</div>
              <div class="nav-desc">{card.desc}</div>
            </div>
            <span class="material-symbols-outlined nav-arrow">arrow_forward</span>
          </a>
        {/each}
      </div>
    </div>

    <!-- Recent Registrations -->
    <div class="section-block">
      <div class="section-title">
        <span class="material-symbols-outlined">history</span>
        Registrasi Agen Terbaru
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Nama Agen</th>
              <th>Tier</th>
              <th>Status</th>
              <th>Total Komisi</th>
              <th>Bergabung</th>
              <th class="align-right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each recentAgents as agent (agent.id)}
              <tr>
                <td><span class="agent-name">{agent.name}</span></td>
                <td><span class="tier-badge tier-{agent.tier.toLowerCase()}">{agent.tier}</span></td>
                <td><span class="status-badge status-{agent.status}">{STATUS_LABELS[agent.status]}</span></td>
                <td>{fmtCurrency(agent.commission_total)}</td>
                <td>{fmtDate(agent.joined)}</td>
                <td class="align-right">
                  <a href="/console/crm/agency/onboarding" class="action-link">Detail</a>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <div class="table-footer">
        <a href="/console/crm/agency/onboarding" class="view-all-link">Lihat semua agen &rarr;</a>
      </div>
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }

  .topbar {
    position: sticky; top: 0; z-index: 30; height: 4rem;
    background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center; justify-content: space-between;
    gap: 1rem; backdrop-filter: blur(8px);
  }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.88rem; color: #434655; }
  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .sep { color: #b0b3c1; }
  .bc-link { color: #004ac6; text-decoration: none; font-weight: 500; }
  .bc-link:hover { text-decoration: underline; }
  .topbar-current { font-weight: 600; color: #191c1e; }
  .top-actions { display: flex; align-items: center; gap: 0.35rem; }
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas { padding: 1.5rem; max-width: 96rem; }

  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.5rem; flex-wrap: wrap; gap: 0.75rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; line-height: 1.2; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .primary-btn {
    display: inline-flex; align-items: center; gap: 0.4rem;
    padding: 0.5rem 1rem; background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff; border: none; border-radius: 0.3rem; font-size: 0.82rem;
    font-weight: 600; cursor: pointer; text-decoration: none;
  }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }
  .primary-btn:hover { opacity: 0.9; }

  .kpi-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 1rem; margin-bottom: 1.5rem; }
  .kpi-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; padding: 1rem; display: flex; align-items: center; gap: 0.85rem; }
  .kpi-icon { width: 2.8rem; height: 2.8rem; border-radius: 0.35rem; display: grid; place-items: center; flex-shrink: 0; }
  .kpi-icon .material-symbols-outlined { font-size: 1.4rem; }
  .kpi-value { font-size: 1.4rem; font-weight: 700; color: #191c1e; line-height: 1.1; }
  .kpi-label { font-size: 0.72rem; color: #434655; margin-top: 0.15rem; }

  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; margin-bottom: 1.25rem; overflow: hidden; }
  .section-title { display: flex; align-items: center; gap: 0.5rem; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); font-size: 0.82rem; font-weight: 700; color: #191c1e; }
  .section-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }

  .nav-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 0.75rem; padding: 1rem; }
  .nav-card { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.35rem; text-decoration: none; color: inherit; transition: border-color 0.15s; }
  .nav-card:hover { border-color: #2563eb; }
  .nav-icon { font-size: 1.3rem; color: #004ac6; flex-shrink: 0; }
  .nav-text { flex: 1; min-width: 0; }
  .nav-label { font-size: 0.82rem; font-weight: 700; color: #191c1e; }
  .nav-desc { font-size: 0.68rem; color: #737686; margin-top: 0.1rem; }
  .nav-arrow { font-size: 0.9rem; color: #b0b3c1; flex-shrink: 0; }
  .nav-card:hover .nav-arrow { color: #2563eb; }

  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .align-right { text-align: right; }

  .agent-name { font-weight: 600; color: #191c1e; }

  .tier-badge { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .tier-gold { background: #fef3c7; color: #92400e; }
  .tier-silver { background: #f1f5f9; color: #334155; }
  .tier-bronze { background: #fde8d8; color: #7c2d12; }

  .status-badge { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .status-active { background: #d1fae5; color: #065f46; }
  .status-pending_kyc { background: #fef3c7; color: #92400e; }
  .status-pending_mou { background: #ede9fe; color: #4c1d95; }
  .status-registered { background: #f1f5f9; color: #334155; }

  .action-link { font-size: 0.72rem; font-weight: 600; color: #2563eb; text-decoration: none; }
  .action-link:hover { text-decoration: underline; }

  .table-footer { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; }
  .view-all-link { font-size: 0.72rem; color: #2563eb; text-decoration: none; font-weight: 600; }
  .view-all-link:hover { text-decoration: underline; }
</style>
