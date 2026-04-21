<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  let { data } = $props();

  const navCtaHref = $derived(
    data.packages[0] ? `/booking/${data.packages[0].id}` : '/packages'
  );

  const trustItems = [
    { title: 'Izin PPIU No. 123/2024', subtitle: 'Resmi Kemenag RI', icon: 'verified_user' },
    { title: 'Akreditasi A', subtitle: 'Standar Internasional', icon: 'workspace_premium' },
    { title: '25.000+ Jamaah', subtitle: 'Telah Berangkat', icon: 'groups' }
  ];

  const packageImages = [
    '/images/packages/pkg-1.png',
    '/images/packages/pkg-2.png',
    '/images/packages/pkg-3.png',
    '/images/packages/pkg-4.png',
    '/images/packages/pkg-5.png',
    '/images/packages/pkg-6.png'
  ];

  const cardTemplates = [
    {
      title: 'Paket Gold',
      price: 'Rp 28.5jt',
      tag: 'Paling Populer',
      departure: '15 Jan 2025',
      left1: 'Bintang 5',
      right1: 'Garuda Indonesia',
      left2: 'City Tour',
      right2: 'Sisa 5 Kursi',
      right2Critical: true,
      right2Icon: 'event_seat',
      left1Icon: 'hotel',
      right1Icon: 'flight_takeoff',
      left2Icon: 'directions_bus',
      tagStyle: 'secondary'
    },
    {
      title: 'Umroh Plus Turki',
      price: 'Rp 42.0jt',
      tag: 'Plus Turki',
      departure: '10 Feb 2025',
      left1: 'Bintang 5',
      right1: 'Turkish Airlines',
      left2: '3 Hari Turki',
      right2: 'Sisa 12 Kursi',
      right2Critical: false,
      right2Icon: 'event_seat',
      left1Icon: 'hotel',
      right1Icon: 'flight_takeoff',
      left2Icon: 'explore',
      tagStyle: 'primary'
    },
    {
      title: 'Paket Silver',
      price: 'Rp 24.9jt',
      tag: 'Hemat',
      departure: '22 Jan 2025',
      left1: 'Bintang 4',
      right1: 'Saudia Air',
      left2: 'Full Board',
      right2: 'Tersedia',
      right2Critical: false,
      right2Icon: 'event_seat',
      left1Icon: 'hotel',
      right1Icon: 'flight_takeoff',
      left2Icon: 'fastfood',
      tagStyle: 'neutral'
    },
    {
      title: 'Platinum Exclusive',
      price: 'Rp 55.0jt',
      tag: 'Mewah',
      departure: '05 Mar 2025',
      left1: 'Fairmont/Setaraf',
      right1: 'Business Class',
      left2: 'Private Mutawif',
      right2: 'Sisa 2 Kursi',
      right2Critical: true,
      right2Icon: 'event_seat',
      left1Icon: 'hotel',
      right1Icon: 'flight_takeoff',
      left2Icon: 'shield',
      tagStyle: 'secondary'
    },
    {
      title: 'Ramadhan Mubarak',
      price: 'Rp 38.5jt',
      tag: 'Best Value',
      departure: '28 Mar 2025',
      left1: 'Bintang 5',
      right1: 'Sahur & Iftar',
      left2: "12 Hari I'tikaf",
      right2: 'Terbatas',
      right2Critical: false,
      right2Icon: 'event_seat',
      left1Icon: 'hotel',
      right1Icon: 'restaurant',
      left2Icon: 'mosque',
      tagStyle: 'primary'
    },
    {
      title: 'Plus Dubai & Oman',
      price: 'Rp 48.9jt',
      tag: 'New Route',
      departure: '12 Apr 2025',
      left1: 'Premium Hotel',
      right1: 'Emirates',
      left2: '15 Hari Trip',
      right2: 'Buka Kuota',
      right2Critical: false,
      right2Icon: 'event_seat',
      left1Icon: 'hotel',
      right1Icon: 'flight_takeoff',
      left2Icon: 'location_on',
      tagStyle: 'amber'
    }
  ];

  const cards = Array.from({ length: 6 }, (_, idx) => {
    const live = data.packages[idx % Math.max(1, data.packages.length)];
    const template = cardTemplates[idx % cardTemplates.length];
    const sourcePackageId = live?.id ?? data.packages[0]?.id ?? 'demo-package';
    return {
      id: `${sourcePackageId}-${idx + 1}`,
      sourcePackageId,
      image: packageImages[idx % packageImages.length],
      ...template
    };
  });
</script>

<svelte:head>
  <title>Katalog Paket Umroh - UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref={navCtaHref}>
  <main class="main-shell">
    <section class="hero shell">
      <div class="hero-grid">
        <div class="hero-copy">
          <p class="kicker">Katalog Terkurasi 2024/2025</p>
          <h1>Temukan Perjalanan Spiritual Terbaik Anda</h1>
          <p>
            Pilihan paket Umroh premium dengan layanan bintang lima, pendampingan manasik komprehensif, dan kepastian jadwal keberangkatan untuk kekhusyukan ibadah Anda.
          </p>
          <div class="hero-actions">
            <a class="primary-btn hero-btn" href="#katalog">Lihat Paket Populer</a>
            <a class="ghost-btn hero-btn" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
              <span class="material-symbols-outlined fill">chat</span>
              Konsultasi WhatsApp
            </a>
          </div>
        </div>
        <div class="hero-media">
          <div class="hero-image-wrap">
            <img src="/images/kaabah-hero.png" alt="Ka'bah" loading="eager" />
          </div>
          <div class="rating-card">
            <strong>4.9/5</strong>
            <p>Rating Kepuasan Jamaah Indonesia</p>
          </div>
        </div>
      </div>
    </section>

    <section class="trust-strip">
      <div class="shell trust-row">
        {#each trustItems as item, idx (item.title)}
          <article class="trust-item">
            <div class="icon">
              <span class="material-symbols-outlined" class:fill={idx === 1}>{item.icon}</span>
            </div>
            <div>
              <h3>{item.title}</h3>
              <p>{item.subtitle}</p>
            </div>
          </article>
          {#if idx < trustItems.length - 1}
            <span class="divider"></span>
          {/if}
        {/each}
      </div>
    </section>

    <section class="shell filters">
      <div class="filter-wrap">
        <div class="filter-col">
          <p class="filter-label">Bulan Keberangkatan</p>
          <select>
            <option>Semua Bulan</option>
            <option>Januari 2025</option>
            <option>Februari 2025</option>
            <option>Maret 2025</option>
            <option>Ramadhan</option>
          </select>
        </div>
        <div class="filter-col">
          <p class="filter-label">Rentang Budget</p>
          <select>
            <option>Semua Budget</option>
            <option>Hemat (Rp 25jt - 30jt)</option>
            <option>Menengah (Rp 30jt - 40jt)</option>
            <option>Premium (&gt; Rp 40jt)</option>
          </select>
        </div>
        <div class="filter-col">
          <p class="filter-label">Durasi Perjalanan</p>
          <div class="durasi-row">
            <button class="is-active" type="button">9 Hari</button>
            <button type="button">12 Hari</button>
          </div>
        </div>
        <div class="filter-col">
          <p class="filter-label">Jenis Paket</p>
          <div class="chip-row">
            <span class="chip is-active">Reguler</span>
            <span class="chip">Plus Turki</span>
          </div>
        </div>
      </div>
    </section>

    <section id="katalog" class="shell catalog" data-testid="s1-package-catalog">
      <ul class="cards-grid">
        {#each cards as card (card.id)}
          <li class="card">
            <div class="cover-wrap">
              <img src={card.image} alt={card.title} loading="lazy" />
              <span class="tag" class:tag-primary={card.tagStyle === 'primary'} class:tag-secondary={card.tagStyle === 'secondary'} class:tag-neutral={card.tagStyle === 'neutral'} class:tag-amber={card.tagStyle === 'amber'}>{card.tag}</span>
            </div>
            <div class="content">
              <div class="topline">
                <div>
                  <h3>{card.title}</h3>
                  <p class="date"><span class="material-symbols-outlined">calendar_today</span>{card.departure}</p>
                </div>
                <p class="price-wrap">
                  <span>Mulai dari</span>
                  <strong>{card.price}</strong>
                </p>
              </div>
              <div class="meta-grid">
                <p><span class="material-symbols-outlined">{card.left1Icon}</span>{card.left1}</p>
                <p><span class="material-symbols-outlined">{card.right1Icon}</span>{card.right1}</p>
                <p><span class="material-symbols-outlined">{card.left2Icon}</span>{card.left2}</p>
                <p class:critical={card.right2Critical}>
                  <span class="material-symbols-outlined">{card.right2Icon}</span>
                  {card.right2}
                </p>
              </div>
              <div class="actions">
                <a
                  class="btn-primary"
                  href={`/packages/${card.sourcePackageId}`}
                  data-testid={`package-link-${card.sourcePackageId}`}>Lihat Detail</a>
                <a class="btn-secondary" href={`/booking/${card.sourcePackageId}`}>Booking Cepat</a>
              </div>
            </div>
          </li>
        {/each}
      </ul>
      <div class="more-row">
        <button type="button" class="more-btn">
          Tampilkan Lebih Banyak Paket
          <span class="material-symbols-outlined">expand_more</span>
        </button>
      </div>
    </section>
  </main>
</MarketingPageLayout>

<style>
  .main-shell {
    /* Previously 8rem from viewport top; spacer now covers fixed nav (5.2rem) */
    padding-top: calc(8rem - 5.2rem);
  }
  .hero {
    margin-bottom: 4rem;
  }
  .hero-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 3rem;
    align-items: center;
  }
  .kicker {
    margin: 0;
    display: inline-flex;
    border-radius: 999px;
    padding: 0.5rem 1rem;
    text-transform: uppercase;
    font-size: 0.74rem;
    letter-spacing: 0.08em;
    font-weight: 700;
    color: #775a19;
    background: rgba(254, 212, 136, 0.3);
  }
  .hero-copy h1 {
    margin: 1.2rem 0 1rem;
    font-size: clamp(2.25rem, 4.5vw, 3.9rem);
    line-height: 1.1;
    color: #004d34;
    letter-spacing: -0.03em;
  }
  .hero-copy p {
    margin: 0;
    color: #3f4943;
    font-size: 1.1rem;
    line-height: 1.7;
    max-width: 37rem;
  }
  .hero-actions {
    margin-top: 1.7rem;
    display: flex;
    flex-wrap: wrap;
    gap: 0.9rem;
  }
  .hero-btn {
    border-radius: 999px;
    min-height: 3.35rem;
    padding: 0.9rem 1.6rem;
    font-size: 1.02rem;
    font-weight: 700;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.45rem;
  }
  .ghost-btn {
    background: #fff;
    color: #004d34;
    border: 1px solid rgba(190, 201, 193, 0.7);
  }
  .hero-media {
    position: relative;
  }
  .hero-image-wrap {
    border-radius: 2.5rem;
    overflow: hidden;
    box-shadow: 0 26px 42px rgba(17, 24, 39, 0.2);
    transform: rotate(2deg);
    transition: transform 0.6s ease;
  }
  .hero-image-wrap:hover {
    transform: rotate(0deg);
  }
  .hero-image-wrap img {
    width: 100%;
    height: 23rem;
    object-fit: cover;
    display: block;
  }
  .rating-card {
    position: absolute;
    left: -2rem;
    bottom: -2rem;
    border-radius: 1.5rem;
    background: #fff;
    border: 1px solid rgba(190, 201, 193, 0.3);
    box-shadow: 0 16px 28px rgba(0, 0, 0, 0.12);
    padding: 1rem 1rem;
    max-width: 13rem;
  }
  .rating-card strong {
    font-size: 2rem;
    line-height: 1;
    color: #775a19;
  }
  .rating-card p {
    margin: 0.4rem 0 0;
    color: #3f4943;
    font-size: 0.84rem;
    line-height: 1.35;
    font-weight: 600;
  }
  .trust-strip {
    width: 100%;
    background: #f5f3f3;
    margin-bottom: 4rem;
    padding: 2.2rem 0;
  }
  .trust-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1.2rem;
    flex-wrap: wrap;
  }
  .trust-item {
    display: flex;
    align-items: center;
    gap: 0.8rem;
  }
  .trust-item .icon {
    width: 3rem;
    height: 3rem;
    border-radius: 999px;
    display: grid;
    place-items: center;
    background: rgba(0, 103, 71, 0.1);
    color: #006747;
  }
  .trust-item:nth-child(3) .icon {
    background: rgba(119, 90, 25, 0.1);
    color: #775a19;
  }
  .trust-item h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
  }
  .trust-item p {
    margin: 0.15rem 0 0;
    font-size: 0.7rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #3f4943;
  }
  .divider {
    width: 1px;
    height: 2rem;
    background: rgba(190, 201, 193, 0.5);
  }
  .filters {
    margin-bottom: 2.8rem;
  }
  .filter-wrap {
    border-radius: 2rem;
    border: 1px solid rgba(190, 201, 193, 0.24);
    background: #fff;
    padding: 1rem;
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 1rem;
  }
  .filter-label {
    display: block;
    margin: 0 0 0.35rem 0.7rem;
    color: #6f7a72;
    font-size: 0.62rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    font-weight: 700;
  }
  .filter-col select {
    width: 100%;
    border: none;
    border-radius: 1rem;
    background: rgba(228, 226, 226, 0.3);
    padding: 0.8rem 0.9rem;
    font-size: 0.9rem;
    color: #1b1c1c;
    font-family: inherit;
  }
  .durasi-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.4rem;
  }
  .durasi-row button {
    border: none;
    border-radius: 1rem;
    background: rgba(228, 226, 226, 0.3);
    color: #1b1c1c;
    font-weight: 700;
    min-height: 2.9rem;
    cursor: pointer;
  }
  .durasi-row .is-active {
    background: #006747;
    color: #fff;
  }
  .chip-row {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
  .chip {
    border-radius: 999px;
    padding: 0.45rem 0.8rem;
    background: rgba(228, 226, 226, 0.3);
    color: #3f4943;
    font-size: 0.74rem;
    font-weight: 700;
  }
  .chip.is-active {
    background: rgba(254, 212, 136, 0.5);
    color: #775a19;
  }
  .catalog {
    margin-bottom: 5rem;
  }
  .cards-grid {
    margin: 0;
    padding: 0;
    list-style: none;
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 2rem;
  }
  .card {
    border-radius: 2rem;
    overflow: hidden;
    border: 1px solid rgba(190, 201, 193, 0.2);
    background: #fff;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.04);
    transition: box-shadow 0.3s ease;
  }
  .card:hover {
    box-shadow: 0 20px 34px rgba(0, 0, 0, 0.1);
  }
  .cover-wrap {
    height: 16rem;
    position: relative;
    overflow: hidden;
  }
  .cover-wrap img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.7s ease;
  }
  .card:hover .cover-wrap img {
    transform: scale(1.08);
  }
  .tag {
    position: absolute;
    top: 1rem;
    left: 1rem;
    border-radius: 999px;
    padding: 0.37rem 0.85rem;
    font-size: 0.72rem;
    font-weight: 700;
    color: #fff;
  }
  .tag-primary {
    background: #006747;
  }
  .tag-secondary {
    background: #775a19;
  }
  .tag-neutral {
    background: #e4e2e2;
    color: #1b1c1c;
  }
  .tag-amber {
    background: #fed488;
    color: #775a19;
  }
  .content {
    padding: 1.7rem;
    display: grid;
    gap: 0.9rem;
  }
  .topline {
    display: flex;
    justify-content: space-between;
    gap: 0.7rem;
    align-items: flex-start;
  }
  .topline h3 {
    margin: 0;
    color: #004d34;
    font-size: 1.62rem;
    line-height: 1.1;
    letter-spacing: -0.02em;
  }
  .date {
    margin: 0.45rem 0 0;
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
    color: #3f4943;
    font-size: 0.87rem;
  }
  .date .material-symbols-outlined {
    font-size: 1rem;
  }
  .price-wrap {
    margin: 0;
    text-align: right;
  }
  .price-wrap span {
    display: block;
    color: #3f4943;
    font-size: 0.75rem;
    font-weight: 500;
  }
  .price-wrap strong {
    color: #004d34;
    font-size: 1.55rem;
    font-weight: 800;
    letter-spacing: -0.02em;
  }
  .meta-grid {
    padding: 1rem 0;
    border-top: 1px solid rgba(190, 201, 193, 0.35);
    border-bottom: 1px solid rgba(190, 201, 193, 0.35);
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.45rem 0.75rem;
  }
  .meta-grid p {
    margin: 0;
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    color: #3f4943;
    font-size: 0.88rem;
  }
  .meta-grid .material-symbols-outlined {
    color: #775a19;
    font-size: 1.03rem;
  }
  .meta-grid p.critical {
    color: #ba1a1a;
    font-style: italic;
    font-weight: 700;
  }
  .meta-grid p.critical .material-symbols-outlined {
    color: #ba1a1a;
  }
  .actions {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.75rem;
    padding-top: 0.1rem;
  }
  .actions a {
    border-radius: 1rem;
    min-height: 3rem;
    text-decoration: none;
    font-weight: 700;
    font-size: 0.86rem;
    display: inline-flex;
    justify-content: center;
    align-items: center;
  }
  .actions .btn-primary {
    background: #004d34;
    color: #fff;
  }
  .actions .btn-secondary {
    background: #efeded;
    color: #004d34;
  }
  .more-row {
    margin-top: 2.7rem;
    text-align: center;
  }
  .more-btn {
    border: none;
    border-radius: 999px;
    background: #e9e8e7;
    color: #004d34;
    font-weight: 700;
    min-height: 3.2rem;
    padding: 0.65rem 1.7rem;
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    cursor: pointer;
  }
  @media (max-width: 1100px) {
    .hero-grid,
    .cards-grid {
      grid-template-columns: 1fr;
    }
    .hero-image-wrap {
      transform: none;
    }
    .rating-card {
      position: static;
      margin-top: 1rem;
    }
    .filter-wrap {
      grid-template-columns: 1fr 1fr;
    }
    .divider {
      display: none;
    }
  }
  @media (max-width: 760px) {
    .hero-copy h1 {
      font-size: 2.35rem;
    }
    .trust-row {
      flex-direction: column;
      align-items: start;
    }
    .filter-wrap {
      grid-template-columns: 1fr;
    }
    .actions,
    .meta-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
