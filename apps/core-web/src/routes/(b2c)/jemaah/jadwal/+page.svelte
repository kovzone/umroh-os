<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';
  import { onMount } from 'svelte';

  interface ScheduleItem {
    time: string;
    title: string;
    location: string;
    type: 'ibadah' | 'transport' | 'makan' | 'istirahat';
    done?: boolean;
  }

  const schedule: ScheduleItem[] = [
    { time: '05:30', title: 'Shalat Subuh Berjamaah', location: 'Masjidil Haram', type: 'ibadah' },
    { time: '07:00', title: 'Sarapan Pagi', location: 'Hotel Grand Makkah — Lantai 3', type: 'makan' },
    { time: '08:30', title: 'Thawaf Sunnah', location: 'Masjidil Haram — Lantai 1', type: 'ibadah' },
    { time: '11:00', title: 'Kembali ke Hotel', location: 'Hotel Grand Makkah', type: 'transport' },
    { time: '12:00', title: 'Shalat Dzuhur', location: 'Masjidil Haram', type: 'ibadah' },
    { time: '13:00', title: 'Makan Siang', location: 'Hotel Grand Makkah — Lantai 3', type: 'makan' },
    { time: '14:00', title: 'Istirahat Siang', location: 'Hotel Grand Makkah — Kamar', type: 'istirahat' },
    { time: '15:45', title: 'Shalat Ashar', location: 'Masjidil Haram', type: 'ibadah' },
    { time: '17:00', title: 'Sa\'i Sunnah', location: 'Masjidil Haram — Mas\'a', type: 'ibadah' },
    { time: '18:30', title: 'Shalat Maghrib', location: 'Masjidil Haram', type: 'ibadah' },
    { time: '19:30', title: 'Makan Malam', location: 'Hotel Grand Makkah — Lantai 3', type: 'makan' },
    { time: '19:45', title: 'Shalat Isya', location: 'Masjidil Haram', type: 'ibadah' },
    { time: '22:00', title: 'Istirahat Malam', location: 'Hotel Grand Makkah — Kamar', type: 'istirahat' },
  ];

  const announcements = [
    { id: 1, time: '07:45', title: 'Perubahan Jadwal Sa\'i', message: 'Sa\'i sore dipindah ke pukul 17:00. Harap berkumpul di lobi pukul 16:45.', urgent: true },
    { id: 2, time: '06:30', title: 'Cuaca Makkah Hari Ini', message: 'Suhu 38°C, sangat panas. Harap bawa payung dan minum air cukup.', urgent: false },
    { id: 3, time: '05:00', title: 'Reminder Paspor', message: 'Selalu bawa kartu tanda jamaah saat keluar hotel. Jangan pernah tinggalkan.', urgent: false },
  ];

  const prayerTimes = [
    { name: 'Subuh', time: '05:30', passed: true },
    { name: 'Dzuhur', time: '12:30', passed: true },
    { name: 'Ashar', time: '15:45', passed: false, next: true },
    { name: 'Maghrib', time: '18:30', passed: false },
    { name: 'Isya', time: '19:45', passed: false },
  ];

  let countdown = $state('');
  let nowMinutes = $state(0);

  function timeToMinutes(t: string): number {
    const [h, m] = t.split(':').map(Number);
    return h * 60 + m;
  }

  function updateCountdown() {
    const now = new Date();
    nowMinutes = now.getHours() * 60 + now.getMinutes();
    const next = prayerTimes.find(p => timeToMinutes(p.time) > nowMinutes);
    if (next) {
      const diff = timeToMinutes(next.time) - nowMinutes;
      const h = Math.floor(diff / 60);
      const m = diff % 60;
      countdown = `${h > 0 ? h + ' jam ' : ''}${m} menit lagi`;
    } else {
      countdown = 'Isya telah berlalu';
    }
  }

  onMount(() => {
    updateCountdown();
    const interval = setInterval(updateCountdown, 60000);
    return () => clearInterval(interval);
  });

  const typeColor: Record<string, string> = {
    ibadah: '#006747',
    transport: '#1565c0',
    makan: '#775a19',
    istirahat: '#6b7280',
  };

  const typeLabel: Record<string, string> = {
    ibadah: 'Ibadah',
    transport: 'Transportasi',
    makan: 'Makan',
    istirahat: 'Istirahat',
  };
</script>

<svelte:head>
  <title>Jadwal Hari Ini — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="jadwal-root">
    <div class="shell">
      <div class="page-header">
        <div>
          <a href="/jemaah" class="back-link">
            <span class="material-symbols-outlined">arrow_back</span>
            Portal Jamaah
          </a>
          <h1>Jadwal Hari Ini</h1>
          <p>Jumat, 24 Januari 2025 · Makkah Al-Mukarramah</p>
        </div>
      </div>

      <!-- Prayer countdown -->
      <div class="prayer-card">
        <div class="prayer-header">
          <div>
            <div class="prayer-title">Waktu Shalat — Makkah</div>
            <div class="prayer-next">
              <span class="material-symbols-outlined">schedule</span>
              Shalat Ashar {countdown}
            </div>
          </div>
          <a href="/ibadah/shalat" class="prayer-detail-link">Detail</a>
        </div>
        <div class="prayer-times">
          {#each prayerTimes as p}
            <div class="prayer-time-item" class:next={p.next} class:passed={p.passed}>
              <div class="pt-name">{p.name}</div>
              <div class="pt-time">{p.time}</div>
            </div>
          {/each}
        </div>
      </div>

      <div class="two-col">
        <!-- Timeline -->
        <div class="timeline-section">
          <h2>Timeline Kegiatan</h2>
          <div class="timeline">
            {#each schedule as item}
              <div class="tl-item">
                <div class="tl-time">{item.time}</div>
                <div class="tl-dot" style="background: {typeColor[item.type]}"></div>
                <div class="tl-content">
                  <div class="tl-title">{item.title}</div>
                  <div class="tl-location">
                    <span class="material-symbols-outlined">location_on</span>
                    {item.location}
                  </div>
                  <span class="tl-badge" style="background: color-mix(in srgb, {typeColor[item.type]} 12%, transparent); color: {typeColor[item.type]}">
                    {typeLabel[item.type]}
                  </span>
                </div>
              </div>
            {/each}
          </div>
        </div>

        <!-- Announcements -->
        <div class="ann-section">
          <h2>Pengumuman Langsung</h2>
          <div class="ann-list">
            {#each announcements as ann (ann.id)}
              <div class="ann-card" class:urgent={ann.urgent}>
                {#if ann.urgent}
                  <div class="ann-urgent-badge">
                    <span class="material-symbols-outlined">priority_high</span>
                    Penting
                  </div>
                {/if}
                <div class="ann-time">{ann.time}</div>
                <div class="ann-title">{ann.title}</div>
                <div class="ann-msg">{ann.message}</div>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .jadwal-root {
    padding-top: calc(5.2rem + 2rem);
    padding-bottom: 5rem;
    background: #fbf9f8;
    min-height: 100vh;
  }
  .shell { max-width: 80rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    color: #006747;
    font-weight: 600;
    font-size: 0.85rem;
    text-decoration: none;
    margin-bottom: 0.75rem;
  }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 1.5rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  /* Prayer card */
  .prayer-card {
    background: linear-gradient(135deg, #004d34, #006747);
    border-radius: 1.5rem;
    padding: 1.5rem;
    color: #fff;
    margin-bottom: 2rem;
  }
  .prayer-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 1.25rem; }
  .prayer-title { font-size: 0.78rem; opacity: 0.8; text-transform: uppercase; letter-spacing: 0.07em; margin-bottom: 0.3rem; }
  .prayer-next { display: flex; align-items: center; gap: 0.35rem; font-size: 1.1rem; font-weight: 700; font-family: 'Plus Jakarta Sans', sans-serif; }
  .prayer-next .material-symbols-outlined { font-size: 1.1rem; }
  .prayer-detail-link { color: rgba(255,255,255,0.7); font-size: 0.82rem; text-decoration: none; border: 1px solid rgba(255,255,255,0.3); border-radius: 999px; padding: 0.3rem 0.8rem; flex-shrink: 0; }
  .prayer-detail-link:hover { background: rgba(255,255,255,0.1); }
  .prayer-times { display: flex; gap: 0.5rem; flex-wrap: wrap; }
  .prayer-time-item {
    flex: 1;
    min-width: 80px;
    background: rgba(255,255,255,0.1);
    border-radius: 0.75rem;
    padding: 0.65rem 0.85rem;
    text-align: center;
  }
  .prayer-time-item.next { background: rgba(255,255,255,0.25); }
  .prayer-time-item.passed { opacity: 0.5; }
  .pt-name { font-size: 0.72rem; opacity: 0.8; margin-bottom: 0.2rem; }
  .pt-time { font-size: 1rem; font-weight: 700; font-family: 'Plus Jakarta Sans', sans-serif; }
  /* Two col layout */
  .two-col { display: grid; grid-template-columns: 1fr 360px; gap: 1.5rem; align-items: start; }
  /* Timeline */
  .timeline-section h2, .ann-section h2 { margin: 0 0 1.25rem; font-size: 1.1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .timeline { position: relative; }
  .timeline::before { content: ''; position: absolute; left: 3.8rem; top: 0; bottom: 0; width: 2px; background: rgba(190,201,193,0.3); }
  .tl-item { display: flex; gap: 0; align-items: flex-start; margin-bottom: 0; position: relative; }
  .tl-time { width: 3.5rem; font-size: 0.78rem; font-weight: 700; color: #6b7280; padding-top: 0.9rem; flex-shrink: 0; }
  .tl-dot { width: 0.7rem; height: 0.7rem; border-radius: 50%; flex-shrink: 0; margin-top: 1.15rem; margin-right: 0.85rem; position: relative; z-index: 1; border: 2px solid #fbf9f8; }
  .tl-content {
    background: #fff;
    border-radius: 1rem;
    padding: 0.85rem 1rem;
    border: 1px solid rgba(190,201,193,0.2);
    flex: 1;
    margin-bottom: 0.75rem;
  }
  .tl-title { font-size: 0.92rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; margin-bottom: 0.3rem; }
  .tl-location { display: flex; align-items: center; gap: 0.25rem; font-size: 0.78rem; color: #9ca3af; margin-bottom: 0.4rem; }
  .tl-location .material-symbols-outlined { font-size: 0.9rem; }
  .tl-badge { font-size: 0.68rem; font-weight: 700; border-radius: 999px; padding: 0.2rem 0.6rem; }
  /* Announcements */
  .ann-section { position: sticky; top: 6rem; }
  .ann-list { display: flex; flex-direction: column; gap: 0.85rem; }
  .ann-card {
    background: #fff;
    border-radius: 1.2rem;
    padding: 1rem 1.1rem;
    border: 1px solid rgba(190,201,193,0.2);
  }
  .ann-card.urgent { border-color: rgba(186,26,26,0.2); background: rgba(186,26,26,0.02); }
  .ann-urgent-badge { display: inline-flex; align-items: center; gap: 0.25rem; font-size: 0.7rem; font-weight: 700; color: #ba1a1a; background: rgba(186,26,26,0.1); border-radius: 999px; padding: 0.2rem 0.65rem; margin-bottom: 0.4rem; }
  .ann-urgent-badge .material-symbols-outlined { font-size: 0.85rem; }
  .ann-time { font-size: 0.72rem; color: #9ca3af; margin-bottom: 0.2rem; }
  .ann-title { font-size: 0.9rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; margin-bottom: 0.3rem; }
  .ann-msg { font-size: 0.82rem; color: #57534e; line-height: 1.5; }
  @media (max-width: 900px) {
    .two-col { grid-template-columns: 1fr; }
    .ann-section { position: static; }
  }
</style>
