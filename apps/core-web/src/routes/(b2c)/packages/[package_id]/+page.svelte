<script lang="ts">
  import {
    DeparturePicker,
    FaqAccordion,
    MarketingPageLayout,
    StickyBookingBar
  } from '$lib/components/marketing';
  import type { PackageDetail } from '$lib/features/s1-catalog/types';
  import { formatDepartureDayId } from '$lib/utils/format-departure';

  let { data } = $props();

  function defaultDepartureId(p: PackageDetail): string | null {
    const list = p.departures;
    const open = list.find((d) => d.status === 'open');
    return open?.id ?? list[0]?.id ?? null;
  }

  const pkg = $derived(data.package);

  const gallery = $derived.by(() => {
    const g = pkg.galleryPhotoUrls;
    if (g && g.length > 0) {
      return g;
    }
    return [pkg.coverPhotoUrl];
  });

  const mainGallery = $derived(gallery[0] ?? pkg.coverPhotoUrl);
  const thumbs = $derived(gallery.slice(1, 4));
  const extraPhotoCount = $derived(Math.max(0, gallery.length - 4));

  const secondaryBadges = $derived(pkg.secondaryBadges ?? []);
  const inclusions = $derived(
    pkg.inclusions ??
      pkg.highlights.map((h) => ({
        icon: 'check_circle',
        title: h,
        description: ''
      }))
  );
  const importantNotes = $derived(pkg.importantNotes ?? []);
  const itineraryDays = $derived(pkg.itineraryDays ?? []);
  const faqs = $derived(pkg.faqs ?? []);
  const whatsappHref = $derived(pkg.whatsappHref ?? 'https://wa.me/6281200000000');

  const firstOpenDeparture = $derived(
    pkg.departures.find((d) => d.status === 'open') ?? pkg.departures[0] ?? null
  );

  let selectedDepartureId = $state<string | null>(null);
  let lastSyncedPackageId = '';

  $effect(() => {
    const p = data.package;
    if (p.id === lastSyncedPackageId) {
      return;
    }
    lastSyncedPackageId = p.id;
    selectedDepartureId = defaultDepartureId(p);
  });

  const selectedDeparture = $derived(pkg.departures.find((d) => d.id === selectedDepartureId) ?? firstOpenDeparture);

  const canStartBooking = $derived(
    Boolean(selectedDepartureId && selectedDeparture && selectedDeparture.status === 'open')
  );

  const bookingHref = $derived(
    selectedDepartureId
      ? `/booking/${pkg.id}?departure=${encodeURIComponent(selectedDepartureId)}`
      : `/booking/${pkg.id}`
  );

  const heroPrice = $derived(pkg.displayPriceShort ?? pkg.startingPriceLabel);
  const trustPpiu = $derived(pkg.trustPpiu ?? 'Izin PPIU resmi Kemenag RI');
  const ratingScore = $derived(pkg.ratingScore ?? '—');
  const ratingReviews = $derived(pkg.ratingReviewsLabel ?? '');
</script>

<svelte:head>
  <title>{pkg.name} — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref={bookingHref} ctaTestId="nav-booking-cta">
  <div class="detail-root" data-testid="s1-package-detail">
    <main class="main-shell shell">
      <nav class="breadcrumb" aria-label="Breadcrumb">
        <a href="/packages">Paket Umrah</a>
        <span class="material-symbols-outlined crumb-icon">chevron_right</span>
        <span class="current">{pkg.name}</span>
      </nav>

      <section class="hero-grid">
        <div class="gallery">
          <div class="hero-shot">
            <img src={mainGallery} alt={pkg.name} loading="eager" />
          </div>
          <div class="thumbs">
            {#each thumbs as src (src)}
              <div class="thumb">
                <img {src} alt="" loading="lazy" />
              </div>
            {/each}
            {#if extraPhotoCount > 0}
              <div class="thumb thumb-more">
                <img src={gallery[3] ?? mainGallery} alt="" class="dim" loading="lazy" />
                <span class="more-label">+{extraPhotoCount} foto</span>
              </div>
            {/if}
          </div>
        </div>

        <div class="hero-aside">
          <div class="badges">
            {#if pkg.primaryBadge}
              <span class="badge badge-pop">{pkg.primaryBadge}</span>
            {/if}
            {#each secondaryBadges as b (b)}
              <span class="badge badge-muted">{b}</span>
            {/each}
          </div>
          <h1>{pkg.name}</h1>
          <p class="lede">{pkg.description}</p>

          <div class="price-card">
            <p class="price-kicker">Mulai dari</p>
            <p class="price-row">
              <span class="price-big">{heroPrice}</span>
              <span class="price-unit">/pax</span>
            </p>
            {#if pkg.priceFinePrint}
              <p class="price-note">{pkg.priceFinePrint}</p>
            {/if}
          </div>

          <p class="helper-flow">Pilih jadwal keberangkatan di bawah, lalu lanjut ke form booking.</p>

          <div class="hero-actions">
            <a
              class="btn-gradient {!canStartBooking ? 'is-disabled' : ''}"
              href={bookingHref}
              data-testid="s1-start-booking"
              aria-disabled={!canStartBooking}
            >
              <span>Lanjut booking</span>
              <span class="material-symbols-outlined">arrow_forward</span>
            </a>
            <a class="btn-outline" href={whatsappHref} target="_blank" rel="noreferrer">
              <span class="material-symbols-outlined gold">chat</span>
              Chat WhatsApp
            </a>
          </div>

          <div class="trust-row">
            <p class="trust">
              <span class="material-symbols-outlined gold fill">verified</span>
              <span>{trustPpiu}</span>
            </p>
            <p class="trust">
              <span class="material-symbols-outlined gold fill">star</span>
              <span class="rating">{ratingScore}</span>
              {#if ratingReviews}
                <span class="reviews">{ratingReviews}</span>
              {/if}
            </p>
          </div>
        </div>
      </section>

      <nav class="section-tabs" aria-label="Bagian halaman">
        <a href="#ringkasan">Ringkasan</a>
        <a href="#itinerari">Itinerari</a>
        <a href="#fasilitas">Fasilitas</a>
        <a href="#snk">Syarat & Ketentuan</a>
      </nav>

      <section id="ringkasan" class="section-block two-col">
        <div>
          <h2>Yang sudah termasuk</h2>
          <div class="inclusion-grid">
            {#each inclusions as item (item.title)}
              <article class="inclusion-card">
                <div class="inclusion-icon">
                  <span class="material-symbols-outlined">{item.icon}</span>
                </div>
                <div>
                  <h3>{item.title}</h3>
                  {#if item.description}
                    <p>{item.description}</p>
                  {/if}
                </div>
              </article>
            {/each}
          </div>
        </div>
        <aside class="notes-card">
          <h2>Poin penting perjalanan</h2>
          <ul>
            {#each importantNotes as note (note)}
              <li>
                <span class="material-symbols-outlined note-ic">check_circle</span>
                {note}
              </li>
            {/each}
          </ul>
        </aside>
      </section>

      <section class="section-block">
        <h2>Jadwal keberangkatan</h2>
        <p class="section-lede">
          Tap baris untuk memilih tanggal rombongan. Harga final mengikuti kamar dan promosi saat checkout.
        </p>
        <DeparturePicker
          departures={pkg.departures}
          bind:selectedId={selectedDepartureId}
          groupName={`dep-${pkg.id}`}
        />
      </section>

      <section id="itinerari" class="section-block">
        <h2>Rencana perjalanan</h2>
        <div class="timeline">
          {#each itineraryDays as day, idx (day.dayLabel + idx)}
            <div class="timeline-item">
              <div class="timeline-badge" class:timeline-badge--first={idx === 0}>{day.dayLabel}</div>
              <div class="timeline-body">
                <h3>{day.title}</h3>
                <p>{day.body}</p>
              </div>
            </div>
          {/each}
        </div>
      </section>

      <section id="fasilitas" class="section-block">
        <h2>Fasilitas & layanan</h2>
        <div class="prose-card">
          {#if pkg.facilitiesBody}
            <p>{pkg.facilitiesBody}</p>
          {:else}
            <p>Detail fasilitas mengikuti kontrak pemesanan dan brosur resmi travel untuk paket ini.</p>
          {/if}
        </div>
      </section>

      <section id="snk" class="section-block">
        <h2>Syarat & ketentuan</h2>
        <div class="prose-card">
          {#if pkg.termsSummary}
            <p>{pkg.termsSummary}</p>
          {:else}
            <p>Untuk teks legal lengkap, silakan unduh dari halaman konfirmasi booking atau hubungi tim kami.</p>
          {/if}
        </div>
      </section>

      <section class="section-block faq-block">
        <h2>Pertanyaan umum</h2>
        <FaqAccordion items={faqs} />
      </section>
    </main>

    <StickyBookingBar
      priceLabel={pkg.startingPriceLabel}
      dateLabel={selectedDeparture ? formatDepartureDayId(selectedDeparture) : null}
      {bookingHref}
      {whatsappHref}
      canBook={canStartBooking}
      bookingTestId="s1-start-booking-sticky"
    />
  </div>
</MarketingPageLayout>

<style>
  .detail-root {
    padding-bottom: 6rem;
  }

  .material-symbols-outlined.gold {
    color: #775a19;
  }
  .crumb-icon {
    font-size: 1rem;
    color: #6f7a72;
  }

  .main-shell {
    /* Top offset: MarketingPageLayout .top-nav-spacer (5.2rem) + breathing room */
    padding-top: 1.5rem;
  }

  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.875rem;
    color: #3f4943;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }
  .breadcrumb a {
    color: #3f4943;
    text-decoration: none;
    font-weight: 500;
  }
  .breadcrumb a:hover {
    color: #004d34;
  }
  .breadcrumb .current {
    color: #775a19;
    font-weight: 600;
  }

  .hero-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2.5rem;
    margin-bottom: 4rem;
  }
  @media (min-width: 1024px) {
    .hero-grid {
      grid-template-columns: minmax(0, 7fr) minmax(0, 5fr);
      gap: 3rem;
      align-items: start;
    }
  }

  .gallery {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .hero-shot {
    aspect-ratio: 16 / 10;
    border-radius: 1.5rem;
    overflow: hidden;
    background: #efeded;
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
  }
  .hero-shot img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
  .thumbs {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
  }
  .thumb {
    aspect-ratio: 1;
    border-radius: 1rem;
    overflow: hidden;
    background: #efeded;
  }
  .thumb img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
  .thumb-more {
    position: relative;
  }
  .thumb-more .dim {
    opacity: 0.55;
  }
  .more-label {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 0.95rem;
    color: #1b1c1c;
    text-shadow: 0 0 8px #fbf9f8;
  }

  .hero-aside h1 {
    margin: 0 0 0.75rem;
    font-size: clamp(2rem, 4vw, 3rem);
    font-weight: 800;
    color: #004d34;
    letter-spacing: -0.03em;
    line-height: 1.1;
  }
  .lede {
    margin: 0 0 1.5rem;
    font-size: 1.05rem;
    color: #3f4943;
    line-height: 1.6;
    font-weight: 500;
  }
  .badges {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  .badge {
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    padding: 0.35rem 0.75rem;
    border-radius: 999px;
  }
  .badge-pop {
    background: #fed488;
    color: #785a1a;
  }
  .badge-muted {
    background: #efeded;
    color: #3f4943;
  }

  .price-card {
    padding: 1.5rem;
    border-radius: 1.5rem;
    background: #f5f3f3;
    border: 1px solid rgba(190, 201, 193, 0.25);
    margin-bottom: 1rem;
  }
  .price-kicker {
    margin: 0;
    font-size: 0.8rem;
    font-weight: 600;
    color: #3f4943;
  }
  .price-row {
    margin: 0.25rem 0 0;
    display: flex;
    align-items: baseline;
    gap: 0.35rem;
  }
  .price-big {
    font-size: 2.25rem;
    font-weight: 800;
    color: #1b1c1c;
  }
  .price-unit {
    color: #3f4943;
    font-weight: 500;
  }
  .price-note {
    margin: 0.5rem 0 0;
    font-size: 0.75rem;
    color: rgba(63, 73, 67, 0.75);
    font-style: italic;
  }

  .helper-flow {
    font-size: 0.85rem;
    color: #3f4943;
    margin: 0 0 1rem;
  }

  .hero-actions {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  .btn-gradient {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 1rem 1.5rem;
    border-radius: 1rem;
    font-weight: 700;
    font-size: 1.05rem;
    text-decoration: none;
    color: #fff;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
    box-shadow: 0 12px 28px rgba(0, 77, 52, 0.22);
    border: none;
    cursor: pointer;
  }
  .btn-gradient.is-disabled {
    opacity: 0.45;
    pointer-events: none;
    cursor: not-allowed;
  }
  .btn-outline {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 1rem 1.5rem;
    border-radius: 1rem;
    font-weight: 700;
    font-size: 1.05rem;
    text-decoration: none;
    color: #1b1c1c;
    border: 2px solid rgba(190, 201, 193, 0.45);
    background: transparent;
  }
  .btn-outline:hover {
    background: #f5f3f3;
  }

  .trust-row {
    margin-top: 1.75rem;
    padding-top: 1.25rem;
    border-top: 1px solid rgba(190, 201, 193, 0.35);
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    gap: 1rem;
  }
  .trust {
    margin: 0;
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.75rem;
    font-weight: 700;
    color: #3f4943;
  }
  .rating {
    font-size: 0.9rem;
    color: #1b1c1c;
  }
  .reviews {
    font-weight: 500;
    color: #3f4943;
    font-size: 0.75rem;
  }

  .section-tabs {
    position: sticky;
    top: 5rem;
    z-index: 40;
    display: flex;
    gap: 1.5rem;
    flex-wrap: wrap;
    padding: 1rem 0 0.75rem;
    margin-bottom: 2rem;
    border-bottom: 1px solid rgba(190, 201, 193, 0.35);
    background: rgba(251, 249, 248, 0.92);
    backdrop-filter: blur(8px);
  }
  .section-tabs a {
    color: #3f4943;
    text-decoration: none;
    font-family: 'Plus Jakarta Sans', sans-serif;
    font-weight: 600;
    font-size: 0.95rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid transparent;
    margin-bottom: -2px;
    white-space: nowrap;
  }
  .section-tabs a:hover {
    color: #004d34;
  }
  .section-tabs a:first-child {
    color: #004d34;
    border-bottom-color: #775a19;
  }

  .section-block {
    scroll-margin-top: 7rem;
    margin-bottom: 4rem;
  }
  .section-block h2 {
    margin: 0 0 1rem;
    font-size: 1.65rem;
    font-weight: 700;
    color: #1b1c1c;
  }
  .section-lede {
    margin: 0 0 1.25rem;
    color: #3f4943;
    font-size: 0.95rem;
    max-width: 40rem;
  }

  .two-col {
    display: grid;
    gap: 2.5rem;
  }
  @media (min-width: 1024px) {
    .two-col {
      grid-template-columns: minmax(0, 7fr) minmax(0, 5fr);
      align-items: start;
    }
  }

  .inclusion-grid {
    display: grid;
    gap: 1rem;
    grid-template-columns: 1fr;
  }
  @media (min-width: 640px) {
    .inclusion-grid {
      grid-template-columns: 1fr 1fr;
    }
  }
  .inclusion-card {
    display: flex;
    gap: 1rem;
    padding: 1rem;
    border-radius: 1rem;
    background: #fff;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  }
  .inclusion-icon {
    width: 3rem;
    height: 3rem;
    border-radius: 0.75rem;
    background: rgba(0, 77, 52, 0.08);
    color: #004d34;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }
  .inclusion-card h3 {
    margin: 0 0 0.25rem;
    font-size: 1rem;
    font-weight: 700;
  }
  .inclusion-card p {
    margin: 0;
    font-size: 0.85rem;
    color: #3f4943;
    line-height: 1.45;
  }

  .notes-card {
    padding: 2rem;
    border-radius: 1.5rem;
    background: #26486f;
    color: #fff;
    position: relative;
    overflow: hidden;
  }
  .notes-card h2 {
    color: #fff;
    margin-bottom: 1.25rem;
  }
  .notes-card ul {
    margin: 0;
    padding: 0;
    list-style: none;
  }
  .notes-card li {
    display: flex;
    gap: 0.5rem;
    font-size: 0.875rem;
    line-height: 1.55;
    color: rgba(255, 255, 255, 0.88);
    margin-bottom: 1rem;
  }
  .note-ic {
    font-size: 1rem;
    color: #d2e4ff;
    flex-shrink: 0;
    margin-top: 0.1rem;
  }

  .timeline {
    position: relative;
    max-width: 48rem;
    padding-left: 0.25rem;
  }
  .timeline::before {
    content: '';
    position: absolute;
    left: 1.35rem;
    top: 0.5rem;
    bottom: 0.5rem;
    width: 2px;
    background: rgba(190, 201, 193, 0.45);
  }
  .timeline-item {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 2.5rem;
    position: relative;
  }
  .timeline-badge {
    z-index: 1;
    width: 2.75rem;
    height: 2.75rem;
    border-radius: 999px;
    background: #fed488;
    color: #785a1a;
    font-weight: 800;
    font-size: 0.85rem;
    display: grid;
    place-items: center;
    flex-shrink: 0;
    box-shadow: 0 4px 14px rgba(119, 90, 25, 0.2);
  }
  .timeline-badge--first {
    background: #775a19;
    color: #fff;
  }
  .timeline-body h3 {
    margin: 0 0 0.35rem;
    font-size: 1.1rem;
  }
  .timeline-body p {
    margin: 0;
    color: #3f4943;
    line-height: 1.6;
    font-size: 0.95rem;
  }

  .prose-card {
    padding: 1.5rem;
    border-radius: 1rem;
    background: #fff;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
    max-width: 48rem;
  }
  .prose-card p {
    margin: 0;
    color: #3f4943;
    line-height: 1.65;
    font-size: 0.95rem;
  }

  .faq-block h2 {
    text-align: center;
    margin-bottom: 1.5rem;
  }
</style>
