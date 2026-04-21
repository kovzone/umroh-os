<script lang="ts">
  let heroImageFailed = $state(false);

  const trustItems = [
    { title: 'Izin PPIU No. 123/2024', subtitle: 'Terdaftar di Kemenag RI' },
    { title: 'Akreditasi A', subtitle: 'Kualitas layanan terbaik' },
    { title: '25,000+', subtitle: 'Jamaah terlayani sejak 2018' }
  ];

  const packages = [
    {
      tier: 'Pilihan Hemat',
      name: 'Paket Silver',
      price: 'Rp 28.5jt',
      priceSuffix: '/pax',
      highlights: ['Hotel Bintang 3 (Dekat Masjid)', 'Penerbangan Ekonomi Direct', 'Makan 3x Sehari (Prasmanan)'],
      cta: 'Detail Paket',
      popular: false
    },
    {
      tier: 'Keseimbangan Sempurna',
      name: 'Paket Gold',
      price: 'Rp 34.2jt',
      priceSuffix: '/pax',
      highlights: ['Hotel Bintang 4 (Depan Masjid)', 'Penerbangan Garuda/Saudia', 'City Tour Full Access', 'Perlengkapan Premium'],
      cta: 'Pesan Sekarang',
      popular: true
    },
    {
      tier: 'Kemewahan Total',
      name: 'Paket Platinum',
      price: 'Rp 48.9jt',
      priceSuffix: '/pax',
      highlights: ['Hotel Bintang 5 (VVIP Access)', 'Business Class Flights', 'Private Transport & Guide'],
      cta: 'Detail Paket',
      popular: false
    }
  ];

  const bookingSteps = [
    {
      title: 'Pilih Paket',
      desc: 'Pilih paket umrah yang sesuai dengan budget dan jadwal Anda. Konsultasikan dengan admin kami secara gratis.'
    },
    {
      title: 'Lengkapi Dokumen',
      desc: 'Unggah Paspor, Kartu Keluarga, dan dokumen lainnya melalui sistem dashboard kami yang aman dan mudah.'
    },
    {
      title: 'Mulai Perjalanan',
      desc: 'Selesaikan pembayaran dan ikuti manasik umrah. Anda siap berangkat menuju tanah suci dengan tenang.'
    }
  ];

  const testimonials = [
    {
      quote: 'Alurnya jelas dan tim sangat responsif. Kami jadi lebih tenang dari awal sampai berangkat.',
      author: 'Ibu Siti, Jakarta'
    },
    {
      quote: 'Paket dan harga transparan. Konsultasi via WhatsApp cepat dan membantu.',
      author: 'Bapak Fahmi, Bandung'
    }
  ];

  const faqs = [
    {
      q: 'Apakah tersedia cicilan untuk paket umrah?',
      a: 'Ya, tersedia opsi cicilan bertahap. Simulasi pembiayaan dijelaskan saat konsultasi.'
    },
    {
      q: 'Dokumen apa saja yang dibutuhkan?',
      a: 'Minimal paspor aktif, KTP, dan dokumen pendukung lainnya sesuai ketentuan keberangkatan.'
    },
    {
      q: 'Kapan jadwal keberangkatan terdekat?',
      a: 'Jadwal terdekat bisa dilihat langsung di halaman paket atau ditanyakan via WhatsApp.'
    }
  ];
</script>

<svelte:head>
  <title>UmrohOS — Perjalanan Umrah Nyaman dan Terpercaya</title>
</svelte:head>

<div class="home-stitch">
  <section class="canvas">
    <nav class="top-nav">
      <a class="brand" href="/">UmrohOS</a>
      <div class="nav-links">
        <a class="is-active" href="/">Beranda</a>
        <a href="/packages">Paket Umrah</a>
        <a href="#proses-booking">Proses Booking</a>
        <a href="#testimoni">Testimoni</a>
        <a href="#faq">FAQ</a>
      </div>
      <div class="nav-right">
        <a href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">WhatsApp</a>
        <a class="btn btn-primary btn-compact" href="/packages">Daftar Sekarang</a>
      </div>
    </nav>

    <section class="hero" data-testid="landing-hero">
      <div class="hero-left">
        <p class="kicker">Premium Umrah Experience</p>
        <h1>
          Wujudkan Perjalanan Ibadah yang <span>Nyaman</span> dan Terpercaya
        </h1>
        <p class="desc">
          Nikmati pengalaman spiritual yang dikurasi secara profesional dengan layanan concierge 24/7,
          akomodasi terbaik, dan bimbingan ibadah yang mutawatir.
        </p>
        <div class="cta-row">
          <a class="btn btn-primary" href="/packages" data-testid="browse-packages-cta">Lihat Paket Umrah</a>
          <a class="btn btn-secondary" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer"
            >Konsultasi WhatsApp</a
          >
        </div>
      </div>
      <div class="hero-right">
        {#if !heroImageFailed}
          <img
            src="/images/kaabah-hero.png"
            alt="Ka'bah"
            loading="eager"
            onerror={() => {
              heroImageFailed = true;
            }}
          />
        {:else}
          <div class="image-fallback" role="img" aria-label="Ilustrasi Ka'bah">Ka'bah</div>
        {/if}
        <div class="floating-card">
          <p>Izin Resmi PPIU</p>
          <strong>No. 123/2024</strong>
        </div>
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
    <section class="section package-showcase">
      <div class="showcase-head">
        <h2>Paket Umrah Pilihan</h2>
        <p>Pilih paket yang sesuai dengan kebutuhan perjalanan spiritual Anda dengan transparansi harga dan fasilitas.</p>
      </div>
      <div class="package-grid">
        {#each packages as pkg (pkg.name)}
          <article class="package-card" class:package-card--popular={pkg.popular}>
            {#if pkg.popular}
              <p class="popular-badge">Paling Populer</p>
            {/if}
            <p class="tier">{pkg.tier}</p>
            <h3>{pkg.name}</h3>
            <p class="price">{pkg.price}<span>{pkg.priceSuffix}</span></p>
            <ul>
              {#each pkg.highlights as highlight (highlight)}
                <li>{highlight}</li>
              {/each}
            </ul>
            <a class="pkg-cta" class:pkg-cta--solid={pkg.popular} href="/packages">{pkg.cta}</a>
          </article>
        {/each}
      </div>
    </section>

    <section id="proses-booking" class="section process">
      <div class="process-head">
        <div>
          <h2>Proses Booking Mudah</h2>
          <p>Hanya butuh 3 langkah sederhana untuk mengamankan seat Anda.</p>
        </div>
        <a class="help-link" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">Bantuan 24/7 →</a>
      </div>
      <div class="process-grid">
        {#each bookingSteps as step, idx (step.title)}
          <article class="step-card">
            <span class="step-index">0{idx + 1}</span>
            <h3>{step.title}</h3>
            <p>{step.desc}</p>
          </article>
        {/each}
      </div>
    </section>

    <section id="testimoni" class="section testimonials">
      <h2>Testimoni jamaah</h2>
      <div class="testimonial-grid">
        {#each testimonials as item (item.author)}
          <blockquote>
            <p>"{item.quote}"</p>
            <cite>{item.author}</cite>
          </blockquote>
        {/each}
      </div>
    </section>

    <section id="faq" class="section faq">
      <h2>Pertanyaan yang sering diajukan</h2>
      <div class="faq-list">
        {#each faqs as item (item.q)}
          <details>
            <summary>{item.q}</summary>
            <p>{item.a}</p>
          </details>
        {/each}
      </div>
    </section>
  </section>
</div>

<style>
  .home-stitch {
    color: #1b1c1c;
  }
  .canvas {
    background: #f8f7f4;
    border: 1px solid #e7e3db;
    border-radius: 1rem;
    overflow: hidden;
    margin-bottom: var(--space-4);
    padding-bottom: 1.15rem;
  }
  .top-nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: 0.6rem 0.9rem;
    border-bottom: 1px solid #e5e1d8;
    background: rgba(255, 255, 255, 0.72);
    backdrop-filter: blur(10px);
  }
  .brand {
    font-weight: 700;
    font-size: 1.95rem;
    text-decoration: none;
    color: #0f4f3a;
    margin-right: 0.8rem;
  }
  .nav-links,
  .nav-right {
    display: flex;
    align-items: center;
    gap: 0.95rem;
  }
  .top-nav a:not(.brand):not(.btn) {
    text-decoration: none;
    color: #4e5953;
    font-size: 0.75rem;
    font-weight: 500;
    padding: 0.2rem 0;
  }
  .top-nav a:not(.brand):not(.btn):hover {
    color: #0f4f3a;
  }
  .top-nav .is-active {
    color: #006747;
    font-weight: 700;
    border-bottom: 2px solid #d6a047;
  }
  .hero {
    padding: 2rem 1.2rem 1.6rem;
    display: grid;
    grid-template-columns: 1fr 1.05fr;
    gap: 1.6rem;
  }
  .kicker {
    margin: 0;
    font-size: 0.72rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #7a621f;
    font-weight: 700;
    background: #f7e3ac;
    border-radius: 999px;
    width: fit-content;
    padding: 0.2rem 0.55rem;
  }
  .hero-left h1 {
    font-size: clamp(2rem, 4vw, 4rem);
    margin: var(--space-2) 0;
    line-height: 1.06;
    color: #191f1c;
    max-width: 14ch;
    letter-spacing: -0.02em;
  }
  .hero-left h1 span {
    color: #006747;
  }
  .desc {
    color: #3b544a;
    max-width: 48ch;
    margin: 0;
    line-height: 1.55;
  }
  .hero-right {
    position: relative;
  }
  .hero-right img {
    width: 100%;
    height: 100%;
    min-height: 30rem;
    object-fit: cover;
    border-radius: 1.3rem;
    border: 1px solid #d9d5cb;
  }
  .image-fallback {
    width: 100%;
    min-height: 30rem;
    border-radius: 1.3rem;
    border: 1px solid #d9d5cb;
    background: linear-gradient(135deg, #d4d7d3, #b8c0ba);
    display: flex;
    align-items: center;
    justify-content: center;
    color: #3d4742;
    font-weight: 700;
  }
  .floating-card {
    position: absolute;
    left: 1.2rem;
    bottom: -1.05rem;
    background: #fff;
    border: 1px solid #e3dfd5;
    border-radius: 0.8rem;
    padding: 0.6rem 0.85rem;
    box-shadow: 0 14px 24px rgba(21, 33, 27, 0.14);
  }
  .floating-card p {
    margin: 0;
    font-size: 0.72rem;
    color: #6a726d;
  }
  .floating-card strong {
    font-size: 0.88rem;
    color: #1f332c;
  }
  .cta-row {
    margin: var(--space-4) 0 0 0;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-3);
  }
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.62rem 1.05rem;
    border-radius: var(--radius-card);
    font-weight: 600;
    font-size: 0.9rem;
    text-decoration: none;
    min-height: 2.6rem;
  }
  .btn-compact {
    min-height: 2.2rem;
    padding: 0.45rem 0.8rem;
    font-size: 0.82rem;
  }
  .btn-primary {
    background: #006747;
    color: #fff;
  }
  .btn-primary:hover {
    background: #00533a;
  }
  .btn-secondary {
    border: 1px solid #baccc2;
    color: #2e4b40;
    background: #fff;
  }
  .btn-secondary:hover {
    border-color: #90a89a;
  }
  .btn:focus-visible {
    outline: 2px solid #006747;
    outline-offset: 2px;
  }
  .trust {
    display: grid;
    gap: 0;
    grid-template-columns: repeat(3, 1fr);
    border-top: 1px solid #e5e1d8;
    background: #f4f2ee;
    margin-top: 1.35rem;
  }
  .trust-item {
    padding: 1.5rem 1rem 1.35rem;
    text-align: center;
    border-right: 1px solid #e8e3d9;
  }
  .trust-item:last-child {
    border-right: none;
  }
  .trust-item h3 {
    margin: 0;
    color: #0b5a41;
    font-size: 2rem;
    line-height: 1.1;
    letter-spacing: -0.01em;
    font-weight: 700;
  }
  .trust-item p {
    margin: 0.45rem 0 0;
    color: #56615b;
    font-size: 1rem;
    line-height: 1.2;
  }
  .section {
    margin-bottom: calc(var(--space-4) * 1.2);
    padding: 0 1rem;
  }
  .section h2 {
    margin: 0 0 var(--space-3);
    color: #1e2b26;
  }
  .package-grid,
  .process-grid,
  .testimonial-grid {
    display: grid;
    gap: var(--space-3);
    grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr));
  }
  .package-card,
  .step-card,
  blockquote {
    margin: 0;
    background: #ffffff;
    border: 1px solid #e4e7e2;
    border-radius: 1.25rem;
    padding: var(--space-3);
  }
  .package-showcase {
    background: #fcfbf8;
    border-top: 1px solid #e8e3d9;
    margin-top: 0;
    margin-bottom: var(--space-4);
    padding-top: 1.8rem;
    padding-bottom: 2rem;
  }
  .showcase-head {
    text-align: center;
    max-width: 44rem;
    margin: 0 auto 1.2rem;
  }
  .showcase-head h2 {
    font-size: clamp(2rem, 4vw, 3rem);
    line-height: 1.1;
    margin-bottom: 0.7rem;
  }
  .showcase-head p {
    margin: 0;
    color: #4f5e57;
    font-size: 1rem;
  }
  .package-card {
    position: relative;
    padding: 1.15rem 1.15rem 1.3rem;
    border-color: #e8e6e1;
    box-shadow: 0 6px 20px rgba(27, 36, 31, 0.06);
  }
  .package-card--popular {
    border: 2px solid #006747;
    box-shadow: 0 18px 32px rgba(27, 36, 31, 0.14);
    transform: translateY(-0.6rem);
  }
  .popular-badge {
    position: absolute;
    top: -0.8rem;
    left: 50%;
    transform: translateX(-50%);
    margin: 0;
    padding: 0.28rem 0.8rem;
    border-radius: 999px;
    background: #006747;
    color: #fff;
    font-size: 0.74rem;
    font-weight: 700;
  }
  .tier {
    margin: 0 0 var(--space-2) 0;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    font-size: 0.76rem;
    color: #6e6451;
    font-weight: 600;
  }
  .package-card h3,
  .step-card h3 {
    margin: 0;
    color: #203c31;
  }
  .price {
    margin: var(--space-2) 0 0.4rem;
    color: #006747;
    font-weight: 700;
    font-size: 2.1rem;
    letter-spacing: -0.02em;
    line-height: 1;
  }
  .price span {
    color: #5c6a64;
    font-size: 1.1rem;
    font-weight: 500;
    margin-left: 0.1rem;
  }
  .package-card ul {
    margin: 0.9rem 0 1.1rem;
    padding-left: 1.1rem;
    color: #4f5e57;
    display: grid;
    gap: 0.65rem;
  }
  .pkg-cta {
    margin-top: auto;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    min-height: 3rem;
    border: 1px solid #006747;
    border-radius: 0.7rem;
    text-decoration: none;
    color: #006747;
    font-weight: 700;
    background: #fff;
    width: 100%;
  }
  .pkg-cta--solid {
    background: #006747;
    color: #fff;
  }
  .step-index {
    display: inline-flex;
    width: 3.85rem;
    height: 3.4rem;
    align-items: center;
    justify-content: center;
    border-radius: 0.95rem;
    background: #006747;
    color: #ffffff;
    font-weight: 700;
    font-size: 1.7rem;
    margin-bottom: var(--space-2);
  }
  .step-card p,
  blockquote p {
    margin: var(--space-2) 0 0;
    color: #4b5e54;
    font-size: 1.24rem;
    line-height: 1.5;
  }
  .process {
    background: #f2f1ef;
    border-top: 1px solid #e8e3d9;
    padding-top: 2.4rem;
    padding-bottom: 2.3rem;
  }
  .process-head {
    display: flex;
    justify-content: space-between;
    align-items: end;
    gap: var(--space-3);
    margin-bottom: 1.4rem;
  }
  .process-head h2 {
    margin: 0;
    font-size: clamp(2rem, 3.4vw, 3rem);
    line-height: 1.1;
  }
  .process-head p {
    margin: 0.45rem 0 0;
    color: #4f5e57;
    font-size: 0.98rem;
  }
  .help-link {
    text-decoration: none;
    color: #0b5a41;
    font-weight: 700;
    white-space: nowrap;
  }
  .help-link:hover {
    text-decoration: underline;
  }
  .step-card {
    background: transparent;
    border: none;
    border-radius: 0;
    padding: 0;
  }
  .step-card h3 {
    margin-top: 1rem;
    font-size: 2rem;
    line-height: 1.2;
  }
  .step-index {
    box-shadow: none;
  }
  .process-grid {
    gap: 2.1rem;
  }
  blockquote cite {
    display: block;
    margin-top: var(--space-2);
    color: #2e4b40;
    font-style: normal;
    font-weight: 600;
  }
  .faq-list {
    display: grid;
    gap: var(--space-2);
  }
  .testimonials h2 {
    text-align: center;
    font-size: clamp(2.2rem, 3.6vw, 3.3rem);
    margin-bottom: 1.4rem;
  }
  details {
    background: #fff;
    border: 1px solid #e1e5df;
    border-radius: var(--radius-card);
    padding: var(--space-3);
  }
  summary {
    cursor: pointer;
    font-weight: 600;
    color: #25392f;
  }
  details p {
    margin: var(--space-2) 0 0;
    color: #51625a;
  }
  @media (max-width: 1080px) {
    .hero {
      grid-template-columns: 1fr;
    }
    .hero-right img {
      min-height: 24rem;
    }
    .image-fallback {
      min-height: 24rem;
    }
  }
  @media (max-width: 740px) {
    .top-nav {
      flex-wrap: wrap;
    }
    .nav-links {
      order: 3;
      width: 100%;
      flex-wrap: wrap;
      gap: var(--space-2);
    }
    .nav-right {
      margin-left: auto;
    }
    .trust {
      grid-template-columns: 1fr;
    }
    .package-card--popular {
      transform: none;
    }
    .showcase-head h2 {
      font-size: 2.2rem;
    }
    .price {
      font-size: 1.7rem;
    }
    .process-head {
      flex-direction: column;
      align-items: start;
    }
    .trust-item {
      border-right: none;
      border-bottom: 1px solid #e8e3d9;
      padding: 1rem;
    }
    .trust-item:last-child {
      border-bottom: none;
    }
    .trust-item h3 {
      font-size: 1.35rem;
    }
    .trust-item p {
      font-size: 0.88rem;
    }
    .hero,
    .trust-item,
    .package-card,
    .step-card,
    blockquote,
    details {
      padding: var(--space-3);
    }
    .step-card h3 {
      font-size: 1.6rem;
    }
    .step-card p {
      font-size: 1rem;
    }
    .step-index {
      width: 3rem;
      height: 2.7rem;
      font-size: 1.35rem;
    }
    .cta-row {
      display: grid;
      grid-template-columns: 1fr;
    }
    .btn {
      width: 100%;
    }
  }
</style>
