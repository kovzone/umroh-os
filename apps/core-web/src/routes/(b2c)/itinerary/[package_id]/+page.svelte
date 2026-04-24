<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';
  import { page } from '$app/state';

  let { data } = $props();
  const pkg = $derived(data.package);
  const itinerary = $derived(pkg.itineraryDays ?? []);
  const shareUrl = $derived(typeof window !== 'undefined' ? window.location.href : '');

  let copied = $state(false);

  async function copyLink() {
    try {
      await navigator.clipboard.writeText(shareUrl);
      copied = true;
      setTimeout(() => { copied = false; }, 2000);
    } catch {
      // fallback: select text
    }
  }

  function shareWhatsApp() {
    const text = encodeURIComponent(`Cek itinerari perjalanan umrah ${pkg.name}: ${shareUrl}`);
    window.open(`https://wa.me/?text=${text}`, '_blank');
  }
</script>

<svelte:head>
  <title>Itinerari {pkg.name} — UmrohOS</title>
  <meta name="description" content="Program perjalanan harian paket {pkg.name} dari UmrohOS." />
  <meta property="og:title" content="Itinerari {pkg.name} — UmrohOS" />
  <meta property="og:description" content="Lihat program perjalanan harian lengkap paket umrah {pkg.name}." />
  <meta property="og:image" content={pkg.coverPhotoUrl} />
</svelte:head>

<MarketingPageLayout ctaHref="/packages/{pkg.id}" ctaLabel="Daftar Sekarang" packagesLinkActive={false}>
  <div class="itin-root">

    <!-- Hero -->
    <section class="itin-hero">
      <div class="hero-img-wrap">
        <img src={pkg.coverPhotoUrl} alt={pkg.name} loading="eager" />
        <div class="hero-overlay"></div>
      </div>
      <div class="hero-content shell">
        <nav class="breadcrumb">
          <a href="/packages">Paket Umrah</a>
          <span class="material-symbols-outlined">chevron_right</span>
          <a href="/packages/{pkg.id}">{pkg.name}</a>
          <span class="material-symbols-outlined">chevron_right</span>
          <span>Itinerari</span>
        </nav>
        <h1>{pkg.name}</h1>
        <p class="hero-desc">{pkg.description}</p>
        <div class="hero-meta">
          {#if itinerary.length > 0}
            <span>
              <span class="material-symbols-outlined">calendar_month</span>
              {itinerary.length} hari perjalanan
            </span>
          {/if}
          <span>
            <span class="material-symbols-outlined">verified</span>
            Izin PPIU Kemenag RI
          </span>
        </div>
        <div class="share-row">
          <span class="share-label">Bagikan itinerari:</span>
          <button class="share-btn" type="button" onclick={copyLink}>
            <span class="material-symbols-outlined">{copied ? 'check' : 'link'}</span>
            {copied ? 'Tersalin!' : 'Salin Tautan'}
          </button>
          <button class="share-btn wa" type="button" onclick={shareWhatsApp}>
            <span class="material-symbols-outlined">share</span>
            WhatsApp
          </button>
        </div>
      </div>
    </section>

    <div class="itin-body shell">

      {#if itinerary.length > 0}
        <!-- Navigation dots -->
        <nav class="day-nav" aria-label="Hari itinerari">
          {#each itinerary as day (day.dayLabel)}
            <a href="#{day.dayLabel.replace(/\s+/g, '-')}" class="day-dot" title={day.dayLabel}>
              {day.dayLabel.replace('Hari ke-', '')}
            </a>
          {/each}
        </nav>

        <!-- Timeline -->
        <div class="timeline">
          {#each itinerary as day, idx (day.dayLabel)}
            <article class="timeline-item" id="{day.dayLabel.replace(/\s+/g, '-')}">
              <div class="timeline-marker">
                <div class="day-badge">{day.dayLabel.replace('Hari ke-', 'H-')}</div>
                {#if idx < itinerary.length - 1}
                  <div class="timeline-line"></div>
                {/if}
              </div>
              <div class="timeline-content">
                <h2>{day.title}</h2>
                <p>{day.body}</p>
              </div>
            </article>
          {/each}
        </div>

      {:else}
        <div class="no-itinerary">
          <span class="material-symbols-outlined">map</span>
          <h2>Itinerari segera tersedia</h2>
          <p>Program perjalanan lengkap sedang disiapkan. Hubungi kami untuk informasi lebih lanjut.</p>
          <a class="btn-wa" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
            <span class="material-symbols-outlined">chat</span>
            Tanya via WhatsApp
          </a>
        </div>
      {/if}

      <!-- CTA -->
      <div class="bottom-cta">
        <div class="bottom-cta-inner">
          <div>
            <h3>Tertarik dengan paket ini?</h3>
            <p>Booking sekarang sebelum kursi habis. Proses cepat, pembayaran mudah.</p>
          </div>
          <div class="bottom-cta-btns">
            <a class="btn-book" href="/packages/{pkg.id}">Lihat Detail & Harga</a>
            <a class="btn-wa" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
              <span class="material-symbols-outlined">chat</span>
              Konsultasi
            </a>
          </div>
        </div>
      </div>
    </div>

  </div>
</MarketingPageLayout>

<style>
  .itin-root {
    padding-top: 5.2rem;
    background: #fbf9f8;
    min-height: 100vh;
  }
  .shell {
    max-width: 72rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  /* Hero */
  .itin-hero {
    position: relative;
    margin-bottom: 3rem;
  }
  .hero-img-wrap {
    position: relative;
    height: 28rem;
    overflow: hidden;
  }
  .hero-img-wrap img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
  .hero-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(to bottom, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0.7) 100%);
  }
  .hero-content {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding-bottom: 2.5rem;
    color: #fff;
  }
  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.8rem;
    margin-bottom: 1rem;
    opacity: 0.8;
  }
  .breadcrumb a {
    color: inherit;
    text-decoration: none;
  }
  .breadcrumb .material-symbols-outlined {
    font-size: 0.9rem;
  }
  .hero-content h1 {
    margin: 0 0 0.6rem;
    font-size: clamp(1.8rem, 4vw, 3rem);
    font-weight: 800;
    letter-spacing: -0.02em;
    font-family: 'Plus Jakarta Sans', sans-serif;
    line-height: 1.2;
  }
  .hero-desc {
    margin: 0 0 1rem;
    opacity: 0.85;
    max-width: 42rem;
    font-size: 0.95rem;
    line-height: 1.65;
  }
  .hero-meta {
    display: flex;
    gap: 1.5rem;
    flex-wrap: wrap;
    margin-bottom: 1.2rem;
    font-size: 0.85rem;
    opacity: 0.9;
  }
  .hero-meta span {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
  }
  .hero-meta .material-symbols-outlined {
    font-size: 1rem;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .share-row {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
  }
  .share-label {
    font-size: 0.82rem;
    opacity: 0.8;
  }
  .share-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    border: 1.5px solid rgba(255,255,255,0.5);
    border-radius: 999px;
    background: rgba(255,255,255,0.15);
    color: #fff;
    font-size: 0.82rem;
    font-weight: 600;
    padding: 0.4rem 0.9rem;
    cursor: pointer;
    backdrop-filter: blur(8px);
    transition: background 0.15s;
  }
  .share-btn:hover { background: rgba(255,255,255,0.25); }
  .share-btn.wa { border-color: rgba(254,212,136,0.6); color: #fed488; }
  .share-btn .material-symbols-outlined { font-size: 0.9rem; }
  /* Body */
  .itin-body {
    padding-bottom: 5rem;
  }
  /* Day nav */
  .day-nav {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-bottom: 2.5rem;
    padding: 1rem;
    background: #fff;
    border-radius: 1rem;
    border: 1px solid rgba(190,201,193,0.2);
  }
  .day-dot {
    width: 2.4rem;
    height: 2.4rem;
    border-radius: 50%;
    border: 1.5px solid rgba(190,201,193,0.4);
    display: grid;
    place-items: center;
    font-size: 0.78rem;
    font-weight: 700;
    color: #57534e;
    text-decoration: none;
    transition: background 0.15s, color 0.15s;
  }
  .day-dot:hover {
    background: #006747;
    color: #fff;
    border-color: transparent;
  }
  /* Timeline */
  .timeline {
    display: grid;
    gap: 0;
  }
  .timeline-item {
    display: grid;
    grid-template-columns: 5rem 1fr;
    gap: 1.5rem;
    position: relative;
  }
  .timeline-marker {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .day-badge {
    width: 3.2rem;
    height: 3.2rem;
    border-radius: 1rem;
    background: linear-gradient(135deg, #004d34, #006747);
    color: #fff;
    display: grid;
    place-items: center;
    font-size: 0.72rem;
    font-weight: 800;
    font-family: 'Plus Jakarta Sans', sans-serif;
    flex-shrink: 0;
  }
  .timeline-line {
    width: 2px;
    flex: 1;
    min-height: 2rem;
    background: rgba(190,201,193,0.5);
    margin: 0.3rem 0;
  }
  .timeline-content {
    padding-bottom: 2.5rem;
    padding-top: 0.3rem;
  }
  .timeline-content h2 {
    margin: 0 0 0.6rem;
    font-size: 1.1rem;
    font-weight: 700;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .timeline-content p {
    margin: 0;
    color: #57534e;
    line-height: 1.75;
  }
  /* No itinerary */
  .no-itinerary {
    text-align: center;
    padding: 5rem 0;
    color: #9ca3af;
  }
  .no-itinerary .material-symbols-outlined {
    font-size: 3.5rem;
    display: block;
    margin-bottom: 1rem;
    color: #d1d5db;
  }
  .no-itinerary h2 {
    margin: 0 0 0.5rem;
    color: #374151;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .no-itinerary p {
    margin: 0 0 1.5rem;
    max-width: 30rem;
    margin-left: auto;
    margin-right: auto;
  }
  /* Bottom CTA */
  .bottom-cta {
    margin-top: 2rem;
    padding: 2rem;
    background: linear-gradient(135deg, #004d34, #006747);
    border-radius: 2rem;
    color: #fff;
  }
  .bottom-cta-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 2rem;
    flex-wrap: wrap;
  }
  .bottom-cta h3 {
    margin: 0;
    font-size: 1.3rem;
    font-weight: 700;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .bottom-cta p {
    margin: 0.4rem 0 0;
    opacity: 0.8;
    font-size: 0.9rem;
  }
  .bottom-cta-btns {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  .btn-book {
    text-decoration: none;
    background: #fed488;
    color: #5d4201;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.75rem 1.6rem;
    font-size: 0.9rem;
  }
  .btn-wa {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    border: 1.5px solid rgba(255,255,255,0.5);
    color: #fff;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.75rem 1.6rem;
    font-size: 0.9rem;
  }
  .btn-wa .material-symbols-outlined { font-size: 1rem; }
  @media (max-width: 640px) {
    .hero-img-wrap { height: 20rem; }
    .timeline-item { grid-template-columns: 3.5rem 1fr; gap: 1rem; }
    .day-badge { width: 2.8rem; height: 2.8rem; font-size: 0.65rem; }
    .bottom-cta-inner { flex-direction: column; }
  }
</style>
