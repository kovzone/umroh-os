<script lang="ts">
  const tools = [
    {
      href: '/console/crm/flyer-generator',
      icon: 'image',
      title: 'Flyer Generator',
      desc: 'Buat flyer promosi paket Umroh dengan template siap pakai. Export ke PNG.',
      badge: 'BL-CRM-013'
    },
    {
      href: '/console/crm/omni-flyer',
      icon: 'photo_library',
      title: 'Omni-Flyer',
      desc: 'Buat versi flyer multi-channel: WhatsApp Story, Instagram Post, Twitter Banner.',
      badge: 'BL-CRM-014'
    },
    {
      href: '/console/crm/itinerary',
      icon: 'route',
      title: 'Itinerary Viewer',
      desc: 'Tampilan itinerary interaktif per paket. Bagikan link ke calon jamaah.',
      badge: 'BL-CRM-015'
    },
    {
      href: '/console/crm/copywriting',
      icon: 'auto_awesome',
      title: 'Copywriting Automation',
      desc: 'Generate 3 variasi teks marketing dari detail paket secara otomatis.',
      badge: 'BL-CRM-016'
    }
  ];

  const agencyPortal = [
    {
      href: '/console/crm/agency',
      icon: 'storefront',
      title: 'Agen & Onboarding',
      desc: 'Daftarkan & kelola agen, proses KYC dan MoU.',
      badge: 'BL-CRM-010'
    },
    {
      href: '/console/crm/agency/content',
      icon: 'perm_media',
      title: 'Konten & Pemasaran',
      desc: 'Bank konten, watermark flyer, UTM, kalender, dan analitik.',
      badge: 'BL-CRM-021'
    },
    {
      href: '/console/crm/agency/leads',
      icon: 'person_search',
      title: 'CRM Leads Agen',
      desc: 'Kelola leads agen, drip sequence, segmen, broadcast.',
      badge: 'BL-CRM-027'
    },
    {
      href: '/console/crm/agency/commission',
      icon: 'account_balance_wallet',
      title: 'Komisi & Wallet',
      desc: 'Saldo komisi, riwayat transaksi, dan permintaan payout.',
      badge: 'BL-CRM-012'
    },
    {
      href: '/console/crm/agency/academy',
      icon: 'school',
      title: 'Akademi Agen',
      desc: 'Modul pelatihan, leaderboard, skrip penjualan, dan badge.',
      badge: 'BL-CRM-032'
    },
    {
      href: '/console/crm/agency/team',
      icon: 'group',
      title: 'Tim & Downline',
      desc: 'Struktur downline, tier leveling, kalkulator zakat.',
      badge: 'BL-CRM-037'
    }
  ];
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">hub</span>
      <span class="sep">/</span>
      <span class="topbar-current">CRM Tools</span>
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
        <h2>CRM Tools</h2>
        <p>Alat pemasaran dan komunikasi untuk tim CRM & CS</p>
      </div>
    </div>

    <div class="tools-grid">
      {#each tools as tool}
        <a href={tool.href} class="tool-card">
          <div class="tool-icon">
            <span class="material-symbols-outlined">{tool.icon}</span>
          </div>
          <div class="tool-body">
            <div class="tool-header">
              <span class="tool-title">{tool.title}</span>
              <span class="tool-badge">{tool.badge}</span>
            </div>
            <p class="tool-desc">{tool.desc}</p>
          </div>
          <span class="material-symbols-outlined tool-arrow">arrow_forward</span>
        </a>
      {/each}
    </div>

    <div class="section-divider">
      <span class="material-symbols-outlined">storefront</span>
      <span>Portal Agen</span>
    </div>

    <div class="tools-grid">
      {#each agencyPortal as tool}
        <a href={tool.href} class="tool-card tool-card--agency">
          <div class="tool-icon tool-icon--agency">
            <span class="material-symbols-outlined">{tool.icon}</span>
          </div>
          <div class="tool-body">
            <div class="tool-header">
              <span class="tool-title">{tool.title}</span>
              <span class="tool-badge tool-badge--agency">{tool.badge}</span>
            </div>
            <p class="tool-desc">{tool.desc}</p>
          </div>
          <span class="material-symbols-outlined tool-arrow">arrow_forward</span>
        </a>
      {/each}
    </div>

    <div class="leads-shortcut">
      <span class="material-symbols-outlined">person_add</span>
      <span>Kelola leads?</span>
      <a href="/console/leads">Buka Leads Manager</a>
    </div>
  </section>
</main>

<style>
  .page-shell {
    min-height: 100vh;
    background: #f7f9fb;
  }

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
    font-size: 0.88rem;
  }

  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .sep { color: #b0b3c1; }
  .topbar-current { font-weight: 600; color: #191c1e; }

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

  .canvas {
    padding: 1.5rem;
    max-width: 72rem;
  }

  .page-head {
    margin-bottom: 1.75rem;
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

  .section-divider {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.82rem;
    font-weight: 700;
    color: #434655;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    margin: 1.5rem 0 1rem;
  }

  .section-divider .material-symbols-outlined { font-size: 1rem; color: #059669; }

  .tools-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  @media (max-width: 700px) {
    .tools-grid { grid-template-columns: 1fr; }
  }

  .tool-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1.1rem;
    text-decoration: none;
    color: inherit;
    transition: box-shadow 0.15s, border-color 0.15s;
  }

  .tool-card:hover {
    border-color: #2563eb;
    box-shadow: 0 2px 8px rgb(37 99 235 / 0.1);
  }

  .tool-card--agency:hover {
    border-color: #059669;
    box-shadow: 0 2px 8px rgb(5 150 105 / 0.1);
  }

  .tool-icon {
    width: 3rem;
    height: 3rem;
    border-radius: 0.35rem;
    background: #dbeafe;
    color: #1e3a8a;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }

  .tool-icon--agency {
    background: #d1fae5;
    color: #065f46;
  }

  .tool-icon .material-symbols-outlined { font-size: 1.5rem; }

  .tool-body {
    flex: 1;
    min-width: 0;
  }

  .tool-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.3rem;
  }

  .tool-title {
    font-size: 0.95rem;
    font-weight: 700;
    color: #191c1e;
  }

  .tool-badge {
    font-size: 0.6rem;
    font-weight: 700;
    background: #ede9fe;
    color: #4c1d95;
    padding: 0.1rem 0.35rem;
    border-radius: 0.2rem;
  }

  .tool-badge--agency {
    background: #d1fae5;
    color: #065f46;
  }

  .tool-desc {
    margin: 0;
    font-size: 0.78rem;
    color: #434655;
    line-height: 1.4;
  }

  .tool-arrow {
    font-size: 1.1rem;
    color: #b0b3c1;
    flex-shrink: 0;
    transition: color 0.15s;
  }

  .tool-card:hover .tool-arrow { color: #2563eb; }
  .tool-card--agency:hover .tool-arrow { color: #059669; }

  .leads-shortcut {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    background: #f2f4f6;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    font-size: 0.82rem;
    color: #434655;
  }

  .leads-shortcut .material-symbols-outlined { font-size: 1rem; color: #2563eb; }
  .leads-shortcut a { color: #004ac6; font-weight: 600; text-decoration: none; }
  .leads-shortcut a:hover { text-decoration: underline; }
</style>
