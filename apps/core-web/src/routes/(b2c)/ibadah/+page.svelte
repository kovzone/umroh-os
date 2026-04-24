<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  const navCards = [
    { href: '/ibadah/shalat', icon: 'access_time', label: 'Waktu Shalat', desc: 'Jadwal shalat & adzan Makkah', color: '#006747' },
    { href: '/ibadah/kiblat', icon: 'explore', label: 'Kiblat', desc: 'Kompas arah kiblat', color: '#004d34' },
    { href: '/ibadah/quran', icon: 'menu_book', label: 'Al-Quran', desc: 'Baca & cari ayat Al-Quran', color: '#775a19' },
    { href: '/ibadah/dzikir', icon: 'spa', label: 'Dzikir', desc: 'Kumpulan dzikir & counter', color: '#1565c0' },
    { href: '/ibadah/manasik', icon: 'mosque', label: 'Manasik', desc: 'Panduan manasik umroh & haji', color: '#2e7d32' },
    { href: '/ibadah/tanya-jawab', icon: 'quiz', label: 'Tanya Jawab', desc: 'FAQ & kirim pertanyaan agama', color: '#6a1c6a' },
    { href: '/ibadah/artikel', icon: 'article', label: 'Artikel', desc: 'Kajian & bacaan islami', color: '#c62828' },
  ];

  const prayerTimes = [
    { name: 'Subuh', time: '05:30', passed: true },
    { name: 'Dzuhur', time: '12:30', passed: true },
    { name: 'Ashar', time: '15:45', passed: false, next: true },
    { name: 'Maghrib', time: '18:30', passed: false },
    { name: 'Isya', time: '19:45', passed: false },
  ];

  const ayatHarian = {
    arabic: 'وَاسْتَعِينُوا بِالصَّبْرِ وَالصَّلَاةِ ۚ وَإِنَّهَا لَكَبِيرَةٌ إِلَّا عَلَى الْخَاشِعِينَ',
    latin: 'Wasta\'inu biṣ-ṣabri waṣ-ṣalati wainnaha lakabiratunil-alla alal-khashi\'in',
    translation: 'Dan mohonlah pertolongan (kepada Allah) dengan sabar dan shalat. Dan sesungguhnya yang demikian itu sungguh berat, kecuali bagi orang-orang yang khusyuk.',
    source: 'QS. Al-Baqarah: 45',
  };
</script>

<svelte:head>
  <title>Panduan Ibadah — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="ibadah-root">
    <div class="shell">
      <div class="page-header">
        <div>
          <p class="kicker">Panduan Ibadah</p>
          <h1>Sahabat Spiritual</h1>
          <p class="subtitle">Alat bantu ibadah dan spiritual Anda selama perjalanan umroh</p>
        </div>
      </div>

      <!-- Prayer times strip -->
      <div class="prayer-strip">
        <div class="strip-title">
          <span class="material-symbols-outlined">location_on</span>
          Makkah Al-Mukarramah
        </div>
        <div class="prayer-times">
          {#each prayerTimes as p}
            <div class="pt-item" class:next={p.next} class:passed={p.passed}>
              <div class="pt-name">{p.name}</div>
              <div class="pt-time">{p.time}</div>
              {#if p.next}<div class="pt-next-badge">Selanjutnya</div>{/if}
            </div>
          {/each}
        </div>
      </div>

      <!-- Ayat Harian -->
      <div class="ayat-card">
        <div class="ayat-label">Ayat Hari Ini</div>
        <p class="ayat-arabic">{ayatHarian.arabic}</p>
        <p class="ayat-latin">{ayatHarian.latin}</p>
        <p class="ayat-translation">"{ayatHarian.translation}"</p>
        <div class="ayat-source">{ayatHarian.source}</div>
      </div>

      <!-- Nav grid -->
      <h2 class="section-title">Menu Ibadah</h2>
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
    </div>
  </div>
</MarketingPageLayout>

<style>
  .ibadah-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 72rem; margin: 0 auto; padding: 0 1.5rem; }
  .kicker { display: inline-block; margin: 0 0 0.5rem; background: rgba(254,212,136,0.3); color: #775a19; border-radius: 999px; padding: 0.3rem 0.9rem; font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.08em; }
  .page-header { margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .subtitle { margin: 0.4rem 0 0; color: #6b7280; }
  /* Prayer strip */
  .prayer-strip {
    background: linear-gradient(135deg, #004d34, #006747);
    border-radius: 1.5rem;
    padding: 1.5rem;
    color: #fff;
    margin-bottom: 1.5rem;
  }
  .strip-title { display: flex; align-items: center; gap: 0.35rem; font-size: 0.8rem; opacity: 0.8; margin-bottom: 1rem; }
  .strip-title .material-symbols-outlined { font-size: 0.95rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .prayer-times { display: flex; gap: 0.5rem; flex-wrap: wrap; }
  .pt-item { flex: 1; min-width: 80px; background: rgba(255,255,255,0.1); border-radius: 0.75rem; padding: 0.65rem; text-align: center; }
  .pt-item.next { background: rgba(255,255,255,0.25); }
  .pt-item.passed { opacity: 0.5; }
  .pt-name { font-size: 0.7rem; opacity: 0.8; margin-bottom: 0.15rem; }
  .pt-time { font-size: 1rem; font-weight: 700; font-family: 'Plus Jakarta Sans', sans-serif; }
  .pt-next-badge { font-size: 0.62rem; background: rgba(255,255,255,0.2); border-radius: 999px; padding: 0.1rem 0.4rem; margin-top: 0.25rem; display: inline-block; }
  /* Ayat card */
  .ayat-card {
    background: linear-gradient(135deg, rgba(254,212,136,0.15), rgba(119,90,25,0.05));
    border: 1px solid rgba(119,90,25,0.15);
    border-radius: 1.5rem;
    padding: 1.75rem;
    margin-bottom: 2rem;
    text-align: center;
  }
  .ayat-label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #775a19; margin-bottom: 1rem; }
  .ayat-arabic { margin: 0 0 0.75rem; font-size: 1.5rem; line-height: 2; color: #004d34; direction: rtl; font-family: 'Amiri', 'Scheherazade', serif; }
  .ayat-latin { margin: 0 0 0.5rem; font-size: 0.85rem; color: #6b7280; font-style: italic; }
  .ayat-translation { margin: 0 0 0.75rem; font-size: 0.9rem; color: #1b1c1c; line-height: 1.7; max-width: 40rem; margin-left: auto; margin-right: auto; }
  .ayat-source { font-size: 0.78rem; font-weight: 700; color: #775a19; }
  /* Nav grid */
  .section-title { margin: 0 0 1.2rem; font-size: 1.2rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .nav-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem; }
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
  .nav-icon { width: 3rem; height: 3rem; border-radius: 0.85rem; background: color-mix(in srgb, var(--card-color) 12%, transparent); display: grid; place-items: center; flex-shrink: 0; color: var(--card-color); }
  .nav-icon .material-symbols-outlined { font-size: 1.4rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .nav-text { flex: 1; }
  .nav-label { font-size: 0.95rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .nav-desc { font-size: 0.78rem; color: #9ca3af; margin-top: 0.15rem; }
  .nav-arrow { color: #d1d5db; font-size: 1.1rem; transition: color 0.2s; }
  .nav-card:hover .nav-arrow { color: var(--card-color); }
</style>
