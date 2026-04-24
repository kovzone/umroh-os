<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  const navCards = [
    { href: '/jemaah/jadwal', icon: 'event', label: 'Jadwal', desc: 'Jadwal harian & pengumuman langsung', color: '#006747' },
    { href: '/dokumen', icon: 'description', label: 'Dokumen', desc: 'Kelengkapan dokumen perjalanan', color: '#775a19' },
    { href: '/jemaah/darurat', icon: 'emergency', label: 'Darurat', desc: 'Tombol SOS & kontak darurat', color: '#ba1a1a' },
    { href: '/jemaah/koper', icon: 'luggage', label: 'Koper', desc: 'Lacak status bagasi & koper', color: '#004d34' },
    { href: '/jemaah/forum', icon: 'forum', label: 'Forum', desc: 'Grup diskusi rombongan', color: '#1565c0' },
    { href: '/jemaah/sertifikat', icon: 'workspace_premium', label: 'Sertifikat', desc: 'Unduh e-sertifikat umroh', color: '#6a1c6a' },
    { href: '/jemaah/laporan', icon: 'report', label: 'Laporan', desc: 'Laporkan kendala perjalanan', color: '#c62828' },
    { href: '/jemaah/pengumuman', icon: 'campaign', label: 'Pengumuman', desc: 'Papan info & reuni', color: '#2e7d32' },
  ];

  const reminders = [
    { id: 1, time: '06:00', label: 'Shalat Subuh berjamaah di lobi hotel', done: true },
    { id: 2, time: '08:30', label: 'Persiapan thawaf pagi — kumpul di lobi', done: false },
    { id: 3, time: '12:00', label: 'Makan siang bersama di lantai 3', done: false },
    { id: 4, time: '14:30', label: 'Sa\'i bersama — bawa kartu jamaah', done: false },
    { id: 5, time: '18:00', label: 'Shalat Maghrib berjamaah di Masjidil Haram', done: false },
  ];

  let remindersState = $state(reminders.map(r => ({ ...r })));

  function toggleReminder(id: number) {
    remindersState = remindersState.map(r => r.id === id ? { ...r, done: !r.done } : r);
  }

  const doneCount = $derived(remindersState.filter(r => r.done).length);
</script>

<svelte:head>
  <title>Dasbor Jemaah — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="jm-root">
    <div class="shell">
      <div class="page-header">
        <div>
          <p class="kicker">Portal Jamaah</p>
          <h1>Assalamu'alaikum, Bambang</h1>
          <p class="subtitle">Perjalanan Umroh · 15 Jan – 28 Jan 2025 · Makkah</p>
        </div>
        <div class="trip-badge">
          <span class="material-symbols-outlined">flight</span>
          <div>
            <div class="trip-day">Hari ke-8</div>
            <div class="trip-total">dari 14 hari</div>
          </div>
        </div>
      </div>

      <!-- Quick reminders -->
      <div class="section-card">
        <div class="section-header">
          <h2>
            <span class="material-symbols-outlined">notifications_active</span>
            Pengingat Hari Ini
          </h2>
          <span class="badge-count">{doneCount}/{remindersState.length} selesai</span>
        </div>
        <div class="reminder-list">
          {#each remindersState as r (r.id)}
            <button class="reminder-item" class:done={r.done} onclick={() => toggleReminder(r.id)}>
              <span class="check-icon material-symbols-outlined">
                {r.done ? 'check_circle' : 'radio_button_unchecked'}
              </span>
              <span class="reminder-time">{r.time}</span>
              <span class="reminder-label">{r.label}</span>
            </button>
          {/each}
        </div>
      </div>

      <!-- Nav cards -->
      <h2 class="nav-section-title">Menu Jamaah</h2>
      <div class="nav-grid">
        {#each navCards as card (card.href)}
          <a href={card.href} class="nav-card" style="--card-color: {card.color}">
            <div class="nav-icon">
              <span class="material-symbols-outlined">{card.icon}</span>
            </div>
            <div class="nav-text">
              <div class="nav-label">{card.label}</div>
              <div class="nav-desc">{card.desc}</div>
            </div>
            <span class="material-symbols-outlined nav-arrow">arrow_forward</span>
          </a>
        {/each}
      </div>

      <!-- Ibadah section link -->
      <a href="/ibadah" class="ibadah-banner">
        <span class="material-symbols-outlined ib-icon">mosque</span>
        <div>
          <div class="ib-title">Panduan Ibadah & Spiritual</div>
          <div class="ib-desc">Waktu shalat, kiblat, Al-Quran, dzikir, manasik, dan lebih banyak lagi</div>
        </div>
        <span class="material-symbols-outlined ib-arrow">arrow_forward</span>
      </a>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .jm-root {
    padding-top: calc(5.2rem + 2rem);
    padding-bottom: 5rem;
    background: #fbf9f8;
    min-height: 100vh;
  }
  .shell {
    max-width: 72rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  .page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }
  .kicker {
    display: inline-block;
    margin: 0 0 0.5rem;
    background: rgba(254,212,136,0.3);
    color: #775a19;
    border-radius: 999px;
    padding: 0.3rem 0.9rem;
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .page-header h1 {
    margin: 0;
    font-size: 1.9rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .subtitle { margin: 0.4rem 0 0; color: #6b7280; font-size: 0.9rem; }
  .trip-badge {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: linear-gradient(135deg, #004d34, #006747);
    color: #fff;
    border-radius: 1.2rem;
    padding: 1rem 1.5rem;
    flex-shrink: 0;
  }
  .trip-badge .material-symbols-outlined { font-size: 2rem; }
  .trip-day { font-size: 1.5rem; font-weight: 800; font-family: 'Plus Jakarta Sans', sans-serif; }
  .trip-total { font-size: 0.8rem; opacity: 0.8; }
  .section-card {
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.5rem;
    border: 1px solid rgba(190,201,193,0.3);
    margin-bottom: 2rem;
  }
  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
  }
  .section-header h2 {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 0;
    font-size: 1.05rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .section-header h2 .material-symbols-outlined { color: #006747; font-size: 1.2rem; }
  .badge-count {
    font-size: 0.78rem;
    font-weight: 700;
    background: rgba(0,103,71,0.1);
    color: #006747;
    border-radius: 999px;
    padding: 0.25rem 0.75rem;
  }
  .reminder-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .reminder-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.7rem 0.9rem;
    border-radius: 0.75rem;
    border: 1px solid rgba(190,201,193,0.2);
    background: #fbf9f8;
    cursor: pointer;
    text-align: left;
    font-family: inherit;
    transition: background 0.15s;
  }
  .reminder-item:hover { background: #f0f9f4; }
  .reminder-item.done { opacity: 0.55; }
  .check-icon { font-size: 1.2rem; color: #9ca3af; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .reminder-item.done .check-icon { color: #006747; }
  .reminder-time { font-size: 0.8rem; font-weight: 700; color: #775a19; flex-shrink: 0; }
  .reminder-label { font-size: 0.88rem; color: #1b1c1c; }
  .reminder-item.done .reminder-label { text-decoration: line-through; color: #9ca3af; }
  .nav-section-title {
    margin: 0 0 1.2rem;
    font-size: 1.2rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .nav-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
  }
  .nav-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem 1.2rem;
    background: #fff;
    border: 1px solid rgba(190,201,193,0.3);
    border-radius: 1.2rem;
    text-decoration: none;
    color: inherit;
    transition: box-shadow 0.2s, transform 0.2s;
  }
  .nav-card:hover { box-shadow: 0 6px 16px rgba(0,0,0,0.08); transform: translateY(-1px); }
  .nav-icon {
    width: 3rem;
    height: 3rem;
    border-radius: 0.85rem;
    background: color-mix(in srgb, var(--card-color) 12%, transparent);
    display: grid;
    place-items: center;
    flex-shrink: 0;
    color: var(--card-color);
  }
  .nav-icon .material-symbols-outlined { font-size: 1.4rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .nav-text { flex: 1; }
  .nav-label { font-size: 0.95rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .nav-desc { font-size: 0.78rem; color: #9ca3af; margin-top: 0.15rem; }
  .nav-arrow { color: #d1d5db; font-size: 1.1rem; transition: color 0.2s; }
  .nav-card:hover .nav-arrow { color: var(--card-color); }
  .ibadah-banner {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.2rem 1.5rem;
    background: linear-gradient(135deg, rgba(0,77,52,0.08), rgba(0,103,71,0.04));
    border: 1px solid rgba(0,103,71,0.2);
    border-radius: 1.2rem;
    text-decoration: none;
    color: inherit;
    transition: box-shadow 0.2s;
  }
  .ibadah-banner:hover { box-shadow: 0 4px 12px rgba(0,103,71,0.1); }
  .ib-icon { font-size: 2rem; color: #006747; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .ib-title { font-size: 1rem; font-weight: 700; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .ib-desc { font-size: 0.82rem; color: #57534e; margin-top: 0.2rem; }
  .ib-arrow { color: #006747; margin-left: auto; flex-shrink: 0; }
</style>
