<script lang="ts">
  import type { PageData } from './+page.server';
  let { data }: { data: PageData } = $props();

  const groups = [
    {
      title: 'Tagihan & Pembayaran',
      icon: 'receipt',
      color: '#2563eb',
      items: [
        { label: 'Tagihan Otomatis', desc: 'Jadwalkan & kelola tagihan massal', href: '/console/finance/billing', icon: 'schedule_send' },
        { label: 'Pembayaran Manual', desc: 'Catat pembayaran & DP manual', href: '/console/finance/manual-payment', icon: 'payments' },
        { label: 'Kwitansi Digital', desc: 'Terbitkan & lihat kwitansi', href: '/console/finance/receipts', icon: 'receipt_long' },
      ]
    },
    {
      title: 'Bank & Rekonsiliasi',
      icon: 'account_balance',
      color: '#0891b2',
      items: [
        { label: 'Bank & Rekonsiliasi', desc: 'Catat & rekonsiliasi transaksi bank', href: '/console/finance/bank', icon: 'sync_alt' },
        { label: 'Nilai Tukar', desc: 'Kelola kurs mata uang asing', href: '/console/finance/exchange-rates', icon: 'currency_exchange' },
      ]
    },
    {
      title: 'Buku Besar',
      icon: 'menu_book',
      color: '#7c3aed',
      items: [
        { label: 'AR & AP Subledger', desc: 'Buku besar piutang dan utang', href: '/console/finance/subledger', icon: 'book' },
        { label: 'Master Vendor', desc: 'Kelola data vendor & supplier', href: '/console/finance/vendors', icon: 'store' },
      ]
    },
    {
      title: 'Anggaran & Biaya',
      icon: 'bar_chart',
      color: '#059669',
      items: [
        { label: 'Budget vs Aktual', desc: 'Pantau realisasi anggaran', href: '/console/finance/budget', icon: 'analytics' },
        { label: 'Biaya Proyek & P&L', desc: 'Rincian biaya & laba rugi keberangkatan', href: '/console/finance/project-costs', icon: 'trending_up' },
      ]
    },
    {
      title: 'Aset & Pajak',
      icon: 'domain',
      color: '#d97706',
      items: [
        { label: 'Aset Tetap', desc: 'Inventaris & penyusutan aset', href: '/console/finance/fixed-assets', icon: 'inventory_2' },
        { label: 'Perpajakan', desc: 'PPN, PPh 21 & PPh 23', href: '/console/finance/tax', icon: 'account_balance_wallet' },
      ]
    },
    {
      title: 'Komisi & Persetujuan',
      icon: 'handshake',
      color: '#db2777',
      items: [
        { label: 'Komisi Agen', desc: 'Hitung & setujui komisi agen', href: '/console/finance/commissions', icon: 'percent' },
        { label: 'Otorisasi Pembayaran', desc: 'Setujui atau tolak batch bayar', href: '/console/finance/authorizations', icon: 'verified' },
        { label: 'Kas Kecil', desc: 'Pengeluaran & saldo kas kecil', href: '/console/finance/petty-cash', icon: 'wallet' },
      ]
    },
    {
      title: 'Monitoring & Audit',
      icon: 'monitor_heart',
      color: '#dc2626',
      items: [
        { label: 'Dashboard Real-time', desc: 'KPI keuangan & arus kas live', href: '/console/finance/realtime', icon: 'dashboard' },
        { label: 'Peringatan Jatuh Tempo', desc: 'Piutang & utang yang akan jatuh tempo', href: '/console/finance/aging-alerts', icon: 'warning' },
        { label: 'Jejak Audit', desc: 'Log semua perubahan data keuangan', href: '/console/finance/audit-log', icon: 'history' },
      ]
    },
  ];
</script>

<main class="finance-hub">
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">account_balance</span>
      <span class="breadcrumb-text">Keuangan</span>
    </nav>
    <div class="top-actions">
      <a href="/console/finance/realtime" class="realtime-link">
        <span class="material-symbols-outlined">dashboard</span>
        Dashboard Real-time
      </a>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Modul Keuangan</h2>
      <p>Kelola seluruh aspek keuangan perusahaan travel Umroh & Haji</p>
    </div>

    <!-- Legacy finance report link -->
    <div class="legacy-banner">
      <span class="material-symbols-outlined">info</span>
      <span>Laporan Jurnal & Ringkasan Akun tersedia di <a href="/console/finance/ledger">halaman Buku Besar</a></span>
    </div>

    {#each groups as group}
      <section class="group-section">
        <div class="group-header">
          <span class="material-symbols-outlined group-icon" style="color:{group.color}">{group.icon}</span>
          <h3>{group.title}</h3>
        </div>
        <div class="card-grid">
          {#each group.items as item}
            <a href={item.href} class="nav-card">
              <div class="card-icon-wrap" style="background:{group.color}11">
                <span class="material-symbols-outlined" style="color:{group.color}">{item.icon}</span>
              </div>
              <div class="card-body">
                <div class="card-title">{item.label}</div>
                <div class="card-desc">{item.desc}</div>
              </div>
              <span class="material-symbols-outlined card-arrow">chevron_right</span>
            </a>
          {/each}
        </div>
      </section>
    {/each}
  </section>
</main>

<style>
  .finance-hub { min-height: 100vh; background: #f7f9fb; }

  .topbar {
    position: sticky; top: 0; z-index: 30; height: 4rem;
    background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center; justify-content: space-between;
    gap: 1rem; backdrop-filter: blur(8px);
  }
  .breadcrumb { display: flex; align-items: center; gap: 0.5rem; color: #434655; }
  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .breadcrumb-text { font-size: 0.88rem; font-weight: 600; color: #191c1e; }
  .top-actions { display: flex; align-items: center; gap: 0.5rem; }
  .realtime-link {
    display: inline-flex; align-items: center; gap: 0.35rem;
    font-size: 0.8rem; font-weight: 600; color: #2563eb;
    padding: 0.4rem 0.75rem; border: 1px solid #2563eb;
    border-radius: 0.25rem; text-decoration: none;
  }
  .realtime-link:hover { background: #eff6ff; }
  .realtime-link .material-symbols-outlined { font-size: 1rem; }

  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.5rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .legacy-banner {
    display: flex; align-items: center; gap: 0.5rem;
    background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 0.25rem;
    padding: 0.65rem 0.85rem; font-size: 0.8rem; color: #1e40af;
    margin-bottom: 1.5rem;
  }
  .legacy-banner .material-symbols-outlined { font-size: 1rem; }
  .legacy-banner a { color: #1d4ed8; font-weight: 600; }

  .group-section { margin-bottom: 2rem; }
  .group-header {
    display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.85rem;
  }
  .group-header h3 { margin: 0; font-size: 0.95rem; font-weight: 700; color: #191c1e; }
  .group-icon { font-size: 1.15rem; }

  .card-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 0.75rem;
  }

  .nav-card {
    display: flex; align-items: center; gap: 0.85rem;
    background: #fff; border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.5rem; padding: 0.9rem 1rem;
    text-decoration: none; color: inherit; transition: box-shadow 0.15s, border-color 0.15s;
  }
  .nav-card:hover { box-shadow: 0 2px 12px rgb(0 0 0 / 0.06); border-color: rgb(195 198 215 / 0.8); }

  .card-icon-wrap {
    width: 2.4rem; height: 2.4rem; border-radius: 0.4rem;
    display: grid; place-items: center; flex-shrink: 0;
  }
  .card-icon-wrap .material-symbols-outlined { font-size: 1.2rem; }

  .card-body { flex: 1; min-width: 0; }
  .card-title { font-size: 0.85rem; font-weight: 600; color: #191c1e; }
  .card-desc { font-size: 0.72rem; color: #737686; margin-top: 0.1rem; }

  .card-arrow { font-size: 1.1rem; color: #b0b3c1; flex-shrink: 0; }
</style>
