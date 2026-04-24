<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  type BagStatus = 'collected' | 'checkin' | 'transit' | 'delivered';

  interface Bag {
    id: string;
    label: string;
    weight: string;
    status: BagStatus;
    lastUpdate: string;
  }

  const bags: Bag[] = [
    { id: 'KPR-001', label: 'Koper Besar (Hitam)', weight: '22 kg', status: 'delivered', lastUpdate: 'Hari ini, 14:30' },
    { id: 'KPR-002', label: 'Koper Kabin (Biru)', weight: '7 kg', status: 'delivered', lastUpdate: 'Hari ini, 14:30' },
    { id: 'KPR-003', label: 'Tas Jinjing (Coklat)', weight: '3 kg', status: 'transit', lastUpdate: 'Hari ini, 12:00' },
  ];

  const timeline = [
    { key: 'collected', label: 'Dikumpulkan', icon: 'inventory', desc: 'Koper dikumpulkan petugas' },
    { key: 'checkin', label: 'Check-in Bandara', icon: 'flight_takeoff', desc: 'Sudah check-in di konter' },
    { key: 'transit', label: 'Dalam Perjalanan', icon: 'local_shipping', desc: 'Sedang dikirim ke hotel' },
    { key: 'delivered', label: 'Terkirim', icon: 'check_circle', desc: 'Telah sampai di kamar' },
  ];

  function getStatusOrder(s: BagStatus): number {
    return { collected: 0, checkin: 1, transit: 2, delivered: 3 }[s];
  }

  const statusLabel: Record<BagStatus, string> = {
    collected: 'Dikumpulkan',
    checkin: 'Check-in Bandara',
    transit: 'Dalam Perjalanan',
    delivered: 'Terkirim',
  };

  const statusColor: Record<BagStatus, string> = {
    collected: '#775a19',
    checkin: '#1565c0',
    transit: '#c62828',
    delivered: '#006747',
  };

  const lastUpdate = 'Hari ini, 14:35 WIB';
</script>

<svelte:head>
  <title>Lacak Koper — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="koper-root">
    <div class="shell">
      <a href="/jemaah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Portal Jamaah
      </a>
      <div class="page-header">
        <div>
          <h1>Pelacakan Koper</h1>
          <p>Status real-time bagasi & koper perjalanan Anda</p>
        </div>
        <div class="last-update">
          <span class="material-symbols-outlined">update</span>
          Terakhir: {lastUpdate}
        </div>
      </div>

      <!-- Overall timeline -->
      <div class="timeline-card">
        <h2>Status Pengiriman</h2>
        <div class="main-timeline">
          {#each timeline as step, idx}
            {@const isActive = bags.some(b => getStatusOrder(b.status) >= idx)}
            {@const isCurrent = bags.some(b => getStatusOrder(b.status) === idx)}
            <div class="tl-step" class:active={isActive} class:current={isCurrent}>
              <div class="tl-icon-wrap">
                <span class="material-symbols-outlined tl-icon">{step.icon}</span>
              </div>
              {#if idx < timeline.length - 1}
                <div class="tl-connector" class:filled={isActive}></div>
              {/if}
              <div class="tl-label">{step.label}</div>
              <div class="tl-desc">{step.desc}</div>
            </div>
          {/each}
        </div>
      </div>

      <!-- Per-bag status -->
      <h2 class="section-title">Daftar Koper ({bags.length} item)</h2>
      <div class="bags-list">
        {#each bags as bag (bag.id)}
          <div class="bag-card">
            <div class="bag-icon-wrap">
              <span class="material-symbols-outlined">luggage</span>
            </div>
            <div class="bag-info">
              <div class="bag-label">{bag.label}</div>
              <div class="bag-meta">
                <span class="bag-id">{bag.id}</span>
                <span class="bag-dot">·</span>
                <span class="bag-weight">{bag.weight}</span>
              </div>
              <div class="bag-update">
                <span class="material-symbols-outlined">schedule</span>
                {bag.lastUpdate}
              </div>
            </div>
            <span class="bag-status" style="background: color-mix(in srgb, {statusColor[bag.status]} 12%, transparent); color: {statusColor[bag.status]}">
              {statusLabel[bag.status]}
            </span>
          </div>
        {/each}
      </div>

      <div class="info-box">
        <span class="material-symbols-outlined">info</span>
        <p>Jika koper Anda tidak ditemukan atau ada kerusakan, segera hubungi pembimbing rombongan atau tekan tombol SOS di halaman Darurat.</p>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .koper-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 64rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; flex-wrap: wrap; margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .last-update { display: flex; align-items: center; gap: 0.35rem; font-size: 0.8rem; color: #9ca3af; flex-shrink: 0; padding-top: 0.5rem; }
  .last-update .material-symbols-outlined { font-size: 0.95rem; }
  /* Timeline card */
  .timeline-card {
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.75rem 2rem;
    border: 1px solid rgba(190,201,193,0.2);
    margin-bottom: 2rem;
  }
  .timeline-card h2 { margin: 0 0 1.5rem; font-size: 1.05rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .main-timeline { display: flex; align-items: flex-start; gap: 0; }
  .tl-step { display: flex; flex-direction: column; align-items: center; flex: 1; position: relative; }
  .tl-icon-wrap {
    width: 3rem;
    height: 3rem;
    border-radius: 50%;
    background: rgba(190,201,193,0.15);
    border: 2px solid rgba(190,201,193,0.3);
    display: grid;
    place-items: center;
    color: #9ca3af;
    z-index: 1;
    transition: all 0.2s;
  }
  .tl-step.active .tl-icon-wrap { background: rgba(0,103,71,0.1); border-color: #006747; color: #006747; }
  .tl-step.current .tl-icon-wrap { background: #006747; border-color: #006747; color: #fff; box-shadow: 0 4px 12px rgba(0,103,71,0.3); }
  .tl-icon { font-size: 1.2rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .tl-connector {
    position: absolute;
    top: 1.5rem;
    left: calc(50% + 1.5rem);
    right: calc(-50% + 1.5rem);
    height: 2px;
    background: rgba(190,201,193,0.3);
  }
  .tl-connector.filled { background: #006747; }
  .tl-label { font-size: 0.78rem; font-weight: 700; color: #9ca3af; margin-top: 0.6rem; text-align: center; font-family: 'Plus Jakarta Sans', sans-serif; }
  .tl-step.active .tl-label { color: #006747; }
  .tl-desc { font-size: 0.68rem; color: #d1d5db; text-align: center; margin-top: 0.15rem; }
  .tl-step.active .tl-desc { color: #9ca3af; }
  /* Bag list */
  .section-title { margin: 0 0 1.2rem; font-size: 1.1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .bags-list { display: flex; flex-direction: column; gap: 0.85rem; margin-bottom: 2rem; }
  .bag-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem 1.2rem;
    background: #fff;
    border-radius: 1.2rem;
    border: 1px solid rgba(190,201,193,0.2);
  }
  .bag-icon-wrap {
    width: 3rem;
    height: 3rem;
    border-radius: 0.85rem;
    background: rgba(0,103,71,0.08);
    display: grid;
    place-items: center;
    color: #006747;
    flex-shrink: 0;
  }
  .bag-icon-wrap .material-symbols-outlined { font-size: 1.4rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .bag-info { flex: 1; }
  .bag-label { font-size: 0.92rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .bag-meta { display: flex; align-items: center; gap: 0.4rem; margin-top: 0.1rem; font-size: 0.78rem; color: #9ca3af; }
  .bag-id { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; color: #006747; }
  .bag-dot { color: #d1d5db; }
  .bag-update { display: flex; align-items: center; gap: 0.25rem; font-size: 0.72rem; color: #9ca3af; margin-top: 0.25rem; }
  .bag-update .material-symbols-outlined { font-size: 0.85rem; }
  .bag-status { font-size: 0.76rem; font-weight: 700; border-radius: 999px; padding: 0.3rem 0.8rem; flex-shrink: 0; }
  .info-box { display: flex; gap: 0.75rem; align-items: flex-start; background: rgba(0,103,71,0.06); border-radius: 1rem; padding: 1rem 1.4rem; }
  .info-box .material-symbols-outlined { color: #006747; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .info-box p { margin: 0; font-size: 0.85rem; color: #57534e; line-height: 1.6; }
</style>
