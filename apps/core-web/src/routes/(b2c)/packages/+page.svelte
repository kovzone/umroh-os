<script lang="ts">
  let { data } = $props();

  const trustItems = [
    { title: 'Izin PPIU No. 123/2024', subtitle: 'Terdaftar di Kemenag RI' },
    { title: 'Akreditasi A', subtitle: 'Kualitas layanan terbaik' },
    { title: '25,000+', subtitle: 'Jamaah terlayani sejak 2018' }
  ];

  const filters = [
    { label: 'Bulan', value: 'Semua bulan', active: true },
    { label: 'Budget', value: 'Semua budget', active: false },
    { label: 'Durasi', value: 'Semua durasi', active: false },
    { label: 'Jenis', value: 'Reguler / Plus', active: false }
  ];

  function badgeForSeats(remainingSeats: number): string {
    if (remainingSeats <= 5) return 'Terbatas';
    if (remainingSeats <= 12) return 'Paling Populer';
    return 'Tersedia';
  }

  const packageImages = [
    '/images/packages/pkg-1.png',
    '/images/packages/pkg-2.png',
    '/images/packages/pkg-3.png',
    '/images/packages/pkg-4.png',
    '/images/packages/pkg-5.png',
    '/images/packages/pkg-6.png'
  ];

  const sampleCards = [
    {
      title: 'Paket Silver',
      price: 'Rp 28.5jt',
      star: 'Bintang 3',
      airline: 'Ekonomi Direct',
      cityTour: '3 Hari Turki'
    },
    {
      title: 'Paket Gold',
      price: 'Rp 34.2jt',
      star: 'Bintang 4',
      airline: 'Garuda Indonesia',
      cityTour: '3 Hari Turki'
    },
    {
      title: 'Paket Platinum',
      price: 'Rp 48.9jt',
      star: 'Bintang 5',
      airline: 'Business Class',
      cityTour: '3 Hari Turki'
    }
  ];

  function formatCompactDeparture(label: string): string {
    const firstDateToken = label.match(/\d{4}-\d{2}-\d{2}/)?.[0];
    if (!firstDateToken) return label;

    const date = new Date(`${firstDateToken}T00:00:00`);
    if (Number.isNaN(date.getTime())) return label;

    return new Intl.DateTimeFormat('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    }).format(date);
  }

  const stitchedCards = Array.from({ length: 6 }, (_, idx) => {
    const live = data.packages[idx % Math.max(1, data.packages.length)];
    const sample = sampleCards[idx % sampleCards.length];
    const sourcePackageId = live?.id ?? data.packages[0]?.id ?? 'demo-package';
    return {
      id: `${sourcePackageId}-${idx + 1}`,
      sourcePackageId,
      image: packageImages[idx % packageImages.length],
      title: sample.title,
      price: sample.price,
      star: sample.star,
      airline: sample.airline,
      cityTour: sample.cityTour,
      departure: formatCompactDeparture(live?.nextDepartureLabel ?? '2025-02-10 - 2025-02-21'),
      seats: live?.remainingSeats ?? 9
    };
  });
</script>

<div class="packages-stitch">
  <section class="canvas">
    <nav class="top-nav">
      <a class="brand" href="/">UmrohOS</a>
      <div class="nav-links">
        <a href="/">Beranda</a>
        <a class="is-active" href="/packages">Paket Umrah</a>
        <a href="/#proses-booking">Proses Booking</a>
        <a href="/#testimoni">Testimoni</a>
        <a href="/#faq">FAQ</a>
      </div>
      <div class="nav-right">
        <a href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">WhatsApp</a>
        <a class="btn btn-primary btn-compact" href={data.packages[0] ? `/booking/${data.packages[0].id}` : '/packages'}>Daftar Sekarang</a>
      </div>
    </nav>

    <section class="hero">
      <div>
        <p class="kicker">Katalog paket</p>
        <h1>Temukan perjalanan spiritual terbaik Anda</h1>
        <p class="hero-desc">Paket dikurasi berdasarkan kualitas layanan, transparansi harga, dan kenyamanan perjalanan.</p>
      </div>
      <div class="hero-actions">
        <a class="btn btn-secondary" href="#filter">Filter cepat</a>
        <a class="btn btn-primary" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">Konsultasi WhatsApp</a>
      </div>
    </section>

    <section class="trust">
      {#each trustItems as item (item.title)}
        <article class="trust-item">
          <h3>{item.title}</h3>
          <p>{item.subtitle}</p>
        </article>
      {/each}
    </section>

    <section id="filter" class="filter-bar">
      {#each filters as filter (filter.label)}
        <button type="button" class="filter-chip" class:active={filter.active}>
          <span>{filter.label}</span>
          <strong>{filter.value}</strong>
        </button>
      {/each}
    </section>

    <section id="katalog-paket" class="catalog" data-testid="s1-package-catalog">
      <div class="catalog-head">
        <h2>Paket Umrah Pilihan</h2>
        <p>{stitchedCards.length} paket tersedia</p>
      </div>
      <ul class="grid">
        {#each stitchedCards as pkg (pkg.id)}
          <li class="card">
            <div class="cover-wrap">
              <img src={pkg.image} alt={pkg.title} class="cover" loading="lazy" />
              <p class="badge">{badgeForSeats(pkg.seats)}</p>
            </div>
            <div class="topline">
              <div class="title-block">
                <h3>{pkg.title}</h3>
                <p class="departure" title={pkg.departure}>
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M7 3v3M17 3v3M4 9h16M5 5h14a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1Z" />
                  </svg>
                  <span>{pkg.departure}</span>
                </p>
              </div>
              <p class="price-wrap">
                <span class="label">Mulai dari</span>
                <strong>{pkg.price}</strong>
              </p>
            </div>
            <div class="feature-grid">
              <p><span>🏨</span>{pkg.star}</p>
              <p><span>✈️</span>{pkg.airline}</p>
              <p><span>🚌</span>{pkg.cityTour}</p>
              <p class:seat-warn={pkg.seats <= 5}><span>🚨</span>Sisa {pkg.seats} Kursi</p>
            </div>
            <div class="card-actions">
              <a class="btn btn-primary" href={`/packages/${pkg.sourcePackageId}`} data-testid="package-link-{pkg.sourcePackageId}">Lihat Detail</a>
              <a class="btn btn-secondary" href={`/booking/${pkg.sourcePackageId}`}>Booking Cepat</a>
            </div>
          </li>
        {/each}
      </ul>
    </section>

    <footer class="local-footer">
      <p>UmrohOS &copy; 2026</p>
      <a href="/">Kembali ke beranda</a>
    </footer>
  </section>
</div>

<style>
  .packages-stitch {
    color: #1b1c1c;
    max-width: 76rem;
    margin: 0 auto;
    font-family: 'Plus Jakarta Sans', 'Manrope', Inter, 'Segoe UI', Arial, sans-serif;
  }
  .canvas {
    background: #f8f7f4;
    border: 1px solid #e7e3db;
    border-radius: 1rem;
    overflow: hidden;
    padding-bottom: 1rem;
  }
  .top-nav {
    display: flex; align-items: center; justify-content: space-between; gap: var(--space-3);
    padding: 0.6rem 0.9rem; border-bottom: 1px solid #e5e1d8;
    background: rgba(255, 255, 255, 0.72); backdrop-filter: blur(10px);
  }
  .brand { font-weight: 700; font-size: 1.95rem; text-decoration: none; color: #0f4f3a; margin-right: 0.8rem; }
  .nav-links, .nav-right { display: flex; align-items: center; gap: 0.95rem; }
  .top-nav a:not(.brand):not(.btn) { text-decoration: none; color: #4e5953; font-size: 0.75rem; font-weight: 500; padding: 0.2rem 0; }
  .top-nav .is-active { color: #006747; font-weight: 700; border-bottom: 2px solid #d6a047; }
  .btn {
    display: inline-flex; align-items: center; justify-content: center; padding: 0.55rem 0.92rem;
    border-radius: 0.7rem; font-weight: 700; font-size: 0.78rem; text-decoration: none;
    min-height: 2.3rem; border: 1px solid transparent;
  }
  .btn-compact { min-height: 2rem; padding: 0.35rem 0.72rem; font-size: 0.72rem; }
  .btn-primary { background: #006747; color: #fff; }
  .btn-secondary { background: #fff; color: #2e4b40; border-color: #baccc2; }
  .hero { padding: 1.4rem 1rem 1rem; display: flex; justify-content: space-between; align-items: end; gap: var(--space-3); }
  .kicker {
    margin: 0; font-size: 0.66rem; letter-spacing: 0.08em; text-transform: uppercase; color: #7a621f;
    font-weight: 700; background: #f7e3ac; border-radius: 999px; width: fit-content; padding: 0.2rem 0.55rem;
  }
  .hero h1 { margin: 0.8rem 0 0.4rem; font-size: clamp(2rem, 3.6vw, 3.2rem); line-height: 1.08; max-width: 16ch; }
  .hero-desc { margin: 0; max-width: 52ch; color: #4e5a54; }
  .hero-actions { display: flex; gap: var(--space-2); flex-wrap: wrap; }
  .trust {
    display: grid; gap: 0; grid-template-columns: repeat(3, 1fr);
    border-top: 1px solid #e5e1d8; background: #f4f2ee; margin-top: 1rem;
  }
  .trust-item { padding: 1rem 0.8rem; text-align: center; border-right: 1px solid #e8e3d9; }
  .trust-item:last-child { border-right: none; }
  .trust-item h3 { margin: 0; color: #0b5a41; font-size: 1.15rem; }
  .trust-item p { margin: 0.25rem 0 0; color: #56615b; font-size: 0.8rem; }
  .filter-bar { display: flex; gap: var(--space-2); flex-wrap: wrap; padding: 1.2rem 1rem 0.7rem; }
  .filter-chip {
    display: inline-flex; flex-direction: column; align-items: start; gap: 0.05rem; padding: 0.45rem 0.65rem;
    border-radius: 0.7rem; border: 1px solid #ddd9cf; background: #fff; color: #44524c; cursor: pointer; min-width: 8rem;
  }
  .filter-chip span { font-size: 0.66rem; }
  .filter-chip strong { font-size: 0.75rem; }
  .filter-chip.active { border-color: #006747; background: #f0faf6; color: #0b5a41; }
  .catalog { padding: 0.2rem 1rem 1rem; }
  .catalog-head { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 0.8rem; }
  .catalog-head h2 { margin: 0; font-size: clamp(1.6rem, 3vw, 2.5rem); }
  .catalog-head p { margin: 0; color: #53615b; font-size: 0.86rem; }
  .grid { list-style: none; padding: 0; margin: 0; display: grid; gap: var(--space-3); grid-template-columns: repeat(auto-fill, minmax(18rem, 1fr)); }
  .card {
    background: #fff;
    border: 1px solid #e5e1d8;
    border-radius: 1rem;
    padding: 0.88rem 0.88rem 0.92rem;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .cover-wrap { position: relative; margin: -0.88rem -0.88rem 0.72rem; }
  .cover {
    width: 100%;
    height: 13.2rem;
    object-fit: cover;
    border-radius: 1rem 1rem 0 0;
    border: none;
    display: block;
  }
  .badge {
    margin: 0;
    position: absolute;
    top: 0.74rem;
    left: 0.74rem;
    display: inline-flex;
    width: fit-content;
    padding: 0.34rem 0.7rem;
    border-radius: 999px;
    background: #006747;
    color: #ffffff;
    border: 1px solid #00573c;
    font-size: 0.78rem;
    font-weight: 700;
  }
  .title-block h3 {
    margin: 0;
    font-size: 1.12rem;
    color: #0b4d37;
    line-height: 1.16;
    letter-spacing: -0.01em;
    font-weight: 700;
  }
  .topline {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 0.62rem;
    padding-bottom: 0.78rem;
    border-bottom: 1px solid #e9e5dc;
  }
  .title-block {
    min-width: 0;
    display: grid;
    gap: 0.34rem;
  }
  .departure {
    margin: 0;
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.9rem;
    color: #2f3f38;
    min-width: 0;
  }
  .departure span {
    display: inline-block;
    overflow-wrap: anywhere;
    line-height: 1.2;
  }
  .departure svg {
    width: 0.92rem;
    height: 0.92rem;
    fill: none;
    stroke: currentColor;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
    flex: 0 0 auto;
  }
  .price-wrap {
    margin: 0;
    text-align: right;
    display: grid;
    gap: 0.06rem;
    white-space: nowrap;
    flex: 0 0 auto;
  }
  .price-wrap .label {
    font-size: 0.84rem;
    color: #3a4a43;
    line-height: 1.1;
    font-weight: 500;
  }
  .price-wrap strong {
    color: #006747;
    font-size: 2rem;
    line-height: 1.08;
    letter-spacing: -0.02em;
    font-family: 'Manrope', 'Plus Jakarta Sans', Inter, 'Segoe UI', Arial, sans-serif;
    font-weight: 700;
  }
  .feature-grid {
    margin: 0.82rem 0 1rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.64rem 0.95rem;
  }
  .feature-grid p {
    margin: 0;
    display: inline-flex;
    align-items: center;
    gap: 0.38rem;
    font-size: 0.96rem;
    color: #2f3f38;
  }
  .feature-grid span { width: 0.95rem; text-align: center; }
  .feature-grid .seat-warn {
    color: #c1121f;
    font-style: italic;
    font-weight: 700;
  }
  .card-actions { margin-top: auto; display: flex; gap: 0.7rem; }
  .card-actions .btn { flex: 1; }
  .card-actions .btn {
    min-height: 2.7rem;
    border-radius: 0.9rem;
    font-size: 0.98rem;
    font-weight: 700;
  }
  .local-footer {
    border-top: 1px solid #e5e1d8; margin: 0 1rem; padding: 0.8rem 0 0.2rem;
    display: flex; justify-content: space-between; align-items: center; color: #5d6963; font-size: 0.78rem;
  }
  .local-footer a { text-decoration: none; color: #0b5a41; font-weight: 700; }
  @media (max-width: 900px) {
    .hero { flex-direction: column; align-items: start; }
  }
  @media (max-width: 740px) {
    .top-nav { flex-wrap: wrap; }
    .nav-links { order: 3; width: 100%; flex-wrap: wrap; }
    .trust { grid-template-columns: 1fr; }
    .trust-item { border-right: none; border-bottom: 1px solid #e8e3d9; }
    .card-actions { flex-direction: column; }
    .local-footer { flex-direction: column; align-items: start; gap: 0.35rem; }
  }
</style>
