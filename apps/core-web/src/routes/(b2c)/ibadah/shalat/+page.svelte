<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';
  import { onMount } from 'svelte';

  interface Prayer {
    name: string;
    time: string;
    timeMinutes: number;
    icon: string;
  }

  const prayers: Prayer[] = [
    { name: 'Subuh', time: '05:30', timeMinutes: 5 * 60 + 30, icon: 'wb_twilight' },
    { name: 'Dzuhur', time: '12:30', timeMinutes: 12 * 60 + 30, icon: 'wb_sunny' },
    { name: 'Ashar', time: '15:45', timeMinutes: 15 * 60 + 45, icon: 'sunny' },
    { name: 'Maghrib', time: '18:30', timeMinutes: 18 * 60 + 30, icon: 'wb_twilight' },
    { name: 'Isya', time: '19:45', timeMinutes: 19 * 60 + 45, icon: 'nights_stay' },
  ];

  let nowMinutes = $state(0);
  let adhanEnabled = $state(true);
  let countdown = $state('');
  let nextPrayer = $state('');

  function updateTime() {
    const now = new Date();
    nowMinutes = now.getHours() * 60 + now.getMinutes();
    const next = prayers.find(p => p.timeMinutes > nowMinutes);
    if (next) {
      const diff = next.timeMinutes - nowMinutes;
      const h = Math.floor(diff / 60);
      const m = diff % 60;
      countdown = h > 0 ? `${h} jam ${m} menit` : `${m} menit`;
      nextPrayer = next.name;
    } else {
      countdown = 'Isya telah berlalu';
      nextPrayer = 'Subuh besok';
    }
  }

  onMount(() => {
    updateTime();
    const t = setInterval(updateTime, 60000);
    return () => clearInterval(t);
  });

  function isPassed(p: Prayer): boolean { return p.timeMinutes < nowMinutes; }
  function isNext(p: Prayer): boolean {
    const next = prayers.find(pr => pr.timeMinutes > nowMinutes);
    return next?.name === p.name;
  }
</script>

<svelte:head>
  <title>Waktu Shalat — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="shalat-root">
    <div class="shell">
      <a href="/ibadah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Panduan Ibadah
      </a>
      <div class="page-header">
        <h1>Waktu Shalat</h1>
        <div class="city-badge">
          <span class="material-symbols-outlined">location_on</span>
          Makkah Al-Mukarramah
        </div>
      </div>

      <!-- Next prayer countdown -->
      <div class="next-prayer-card">
        <div class="np-label">Shalat Berikutnya</div>
        <div class="np-name">{nextPrayer}</div>
        <div class="np-countdown">
          <span class="material-symbols-outlined">schedule</span>
          {countdown}
        </div>
      </div>

      <!-- Prayer times list -->
      <div class="prayers-list">
        {#each prayers as p}
          <div class="prayer-row" class:passed={isPassed(p)} class:next={isNext(p)}>
            <div class="pr-icon-wrap">
              <span class="material-symbols-outlined pr-icon">{p.icon}</span>
            </div>
            <div class="pr-info">
              <div class="pr-name">{p.name}</div>
              {#if isPassed(p)}<div class="pr-status passed">Sudah berlalu</div>{/if}
              {#if isNext(p)}<div class="pr-status next">Selanjutnya</div>{/if}
            </div>
            <div class="pr-time">{p.time}</div>
          </div>
        {/each}
      </div>

      <!-- Adhan toggle -->
      <div class="adhan-card">
        <div class="adhan-info">
          <span class="material-symbols-outlined adhan-icon">volume_up</span>
          <div>
            <div class="adhan-title">Notifikasi Adzan</div>
            <div class="adhan-desc">Terima pengingat 10 menit sebelum waktu shalat</div>
          </div>
        </div>
        <button
          class="toggle-btn"
          class:on={adhanEnabled}
          onclick={() => adhanEnabled = !adhanEnabled}
          role="switch"
          aria-checked={adhanEnabled}
        >
          <span class="toggle-knob"></span>
        </button>
      </div>

      <div class="info-box">
        <span class="material-symbols-outlined">info</span>
        <p>Waktu shalat berdasarkan metode Umm Al-Qura (Arab Saudi). Selalu perhatikan pengumuman resmi dari Masjidil Haram.</p>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .shalat-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 56rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 2rem; flex-wrap: wrap; gap: 1rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .city-badge { display: flex; align-items: center; gap: 0.35rem; font-size: 0.85rem; color: #006747; font-weight: 600; background: rgba(0,103,71,0.08); border-radius: 999px; padding: 0.4rem 0.9rem; }
  .city-badge .material-symbols-outlined { font-size: 1rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .next-prayer-card {
    background: linear-gradient(135deg, #004d34, #006747);
    border-radius: 1.5rem;
    padding: 2.5rem;
    color: #fff;
    text-align: center;
    margin-bottom: 2rem;
  }
  .np-label { font-size: 0.8rem; opacity: 0.75; text-transform: uppercase; letter-spacing: 0.1em; margin-bottom: 0.5rem; }
  .np-name { font-size: 2.5rem; font-weight: 800; font-family: 'Plus Jakarta Sans', sans-serif; margin-bottom: 0.75rem; }
  .np-countdown { display: inline-flex; align-items: center; gap: 0.5rem; font-size: 1rem; background: rgba(255,255,255,0.15); border-radius: 999px; padding: 0.5rem 1.25rem; }
  .np-countdown .material-symbols-outlined { font-size: 1.1rem; }
  /* Prayers list */
  .prayers-list { display: flex; flex-direction: column; gap: 0.75rem; margin-bottom: 1.5rem; }
  .prayer-row {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.2rem 1.5rem;
    background: #fff;
    border-radius: 1.2rem;
    border: 1px solid rgba(190,201,193,0.2);
    transition: all 0.2s;
  }
  .prayer-row.passed { opacity: 0.5; }
  .prayer-row.next { border-color: #006747; background: rgba(0,103,71,0.03); box-shadow: 0 4px 12px rgba(0,103,71,0.1); }
  .pr-icon-wrap {
    width: 3rem;
    height: 3rem;
    border-radius: 0.85rem;
    background: rgba(0,103,71,0.08);
    display: grid;
    place-items: center;
    flex-shrink: 0;
    color: #006747;
  }
  .prayer-row.next .pr-icon-wrap { background: #006747; color: #fff; }
  .pr-icon { font-size: 1.3rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .pr-info { flex: 1; }
  .pr-name { font-size: 1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .pr-status { font-size: 0.72rem; font-weight: 600; margin-top: 0.15rem; }
  .pr-status.passed { color: #9ca3af; }
  .pr-status.next { color: #006747; }
  .pr-time { font-size: 1.4rem; font-weight: 800; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; flex-shrink: 0; }
  .prayer-row.next .pr-time { color: #006747; }
  /* Adhan toggle */
  .adhan-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #fff;
    border-radius: 1.2rem;
    border: 1px solid rgba(190,201,193,0.2);
    padding: 1.2rem 1.5rem;
    margin-bottom: 1.5rem;
    gap: 1rem;
  }
  .adhan-info { display: flex; align-items: center; gap: 0.85rem; }
  .adhan-icon { font-size: 1.5rem; color: #006747; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .adhan-title { font-size: 0.95rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .adhan-desc { font-size: 0.78rem; color: #9ca3af; margin-top: 0.1rem; }
  .toggle-btn {
    position: relative;
    width: 3.2rem;
    height: 1.8rem;
    border-radius: 999px;
    border: none;
    background: #d1d5db;
    cursor: pointer;
    transition: background 0.2s;
    flex-shrink: 0;
  }
  .toggle-btn.on { background: #006747; }
  .toggle-knob {
    position: absolute;
    top: 0.2rem;
    left: 0.2rem;
    width: 1.4rem;
    height: 1.4rem;
    border-radius: 50%;
    background: #fff;
    transition: transform 0.2s;
    box-shadow: 0 1px 3px rgba(0,0,0,0.2);
  }
  .toggle-btn.on .toggle-knob { transform: translateX(1.4rem); }
  .info-box { display: flex; gap: 0.75rem; align-items: flex-start; background: rgba(0,103,71,0.06); border-radius: 1rem; padding: 1rem 1.4rem; }
  .info-box .material-symbols-outlined { color: #006747; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .info-box p { margin: 0; font-size: 0.85rem; color: #57534e; line-height: 1.6; }
</style>
