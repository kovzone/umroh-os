<script lang="ts">
  const trustItems = [
    { title: "Izin PPIU No. 123/2024", subtitle: "Terdaftar di Kemenag RI" },
    { title: "Akreditasi A", subtitle: "Kualitas Layanan Terbaik" },
    { title: "25,000+", subtitle: "Jamaah Terlayani Sejak 2018" },
  ];

  const packages = [
    {
      tier: "Pilihan Hemat",
      name: "Paket Silver",
      price: "Rp 28.5jt",
      highlights: [
        "Hotel Bintang 3 (Dekat Masjid)",
        "Penerbangan Ekonomi Direct",
        "Makan 3x Sehari (Prasmanan)",
      ],
      cta: "Detail Paket",
      popular: false,
    },
    {
      tier: "Keseimbangan Sempurna",
      name: "Paket Gold",
      price: "Rp 34.2jt",
      highlights: [
        "Hotel Bintang 4 (Depan Masjid)",
        "Penerbangan Garuda/Saudia",
        "City Tour Full Access",
        "Perlengkapan Premium",
      ],
      cta: "Pesan Sekarang",
      popular: true,
    },
    {
      tier: "Kemewahan Total",
      name: "Paket Platinum",
      price: "Rp 48.9jt",
      highlights: [
        "Hotel Bintang 5 (VVIP Access)",
        "Business Class Flights",
        "Private Transport & Guide",
      ],
      cta: "Detail Paket",
      popular: false,
    },
  ];

  const bookingSteps = [
    {
      title: "Pilih Pakett",
      desc: "Pilih paket umrah yang sesuai dengan budget dan jadwal Anda. Konsultasikan dengan admin kami secara gratis.",
    },
    {
      title: "Lengkapi Dokumen",
      desc: "Unggah Paspor, Kartu Keluarga, dan dokumen lainnya melalui sistem dashboard kami yang aman dan mudah.",
    },
    {
      title: "Mulai Perjalanan",
      desc: "Selesaikan pembayaran dan ikuti manasik umrah. Anda siap berangkat menuju tanah suci dengan tenang.",
    },
  ];

  const testimonials = [
    {
      quote:
        "Alhamdulillah, perjalanan umrah bersama UmrohOS sangat tertata. Bimbingan ustadz sangat mendalam dan fasilitas hotelnya luar biasa dekat.",
      author: "Hj. Ahmad Fauzi",
      package: "Jamaah Paket Platinum",
      avatar:
        "https://lh3.googleusercontent.com/aida-public/AB6AXuDPROeTP-SYQfBebj5bK_9BdmVmLrsHmGFjsqwzip4P9T_xuPVe-fb44h_6p9qR-26STeolkw-ZC_THAqaJrDsympxROv6wlgNfD6Px2LccWGQz4_IASL-x0ff67-D3NGzhQBu-kZ93VrJWfBAd_NNc8NS910gOJ0EM5arHq1UfssWmRM422VhAnrh-4SKXy9YWugFNG56HNkYmOSp-ExYg56S9imKVFzAn3njexlOmDdvHr3GBs8UHB2SzLDq1g5-cb62IJuiI3Q",
    },
    {
      quote:
        "Sangat terbantu dengan sistem cicilannya. Admin responsif membantu dokumen dari awal sampai akhir. Sangat merekomendasikan!",
      author: "Siti Aminah",
      package: "Jamaah Paket Gold",
      avatar:
        "https://lh3.googleusercontent.com/aida-public/AB6AXuDlltemjg_5MMn8lq499CiBs8P9_5uR47NKbZUOXkaS1NctPGX-rePQu-6UaDkc-hUbfdAhvORUBFLcqBVKhPfs-mbLYETcXhz70lU4L-FsjW-ZKxhuDeZRVnCGxJhwuF3ha4xwQj0vNdSZHH5rFC5mL3lyF2q6agKM_ZUTmWAjk0h96LbougfEFM5kpAR1aW_RelQyYyebg4KE29paHImjgHd_RU0pMKXlmFkRg6rgSgKeKd_3cuKGBuTarnypykwHueqUoY5OKw",
    },
    {
      quote:
        "Pelayanan kelas satu. Dari manasik sampai kepulangan, semua tim sangat sigap. UmrohOS benar-benar membuat ibadah jadi tenang.",
      author: "Bambang S.",
      package: "Jamaah Paket Gold",
      avatar:
        "https://lh3.googleusercontent.com/aida-public/AB6AXuCt_WzsFhKzeOvuRxcQP0kDDw619cDS0HsUI7s2uAYuZT4Rkjc6LMQEoLp7mUyQhXYoi5of1Sk7Fkf1uK7ooAYWo71J9Qdpl5tnYlFLaMt_pLSato-z2U7MbNmbNIu7GQa4AQ_uAyP9b4Cxpebz5CdxysvsNHFFMwhQVxy6OhYZOdIDWS771eslvMybGGbn6fIPQVkIPiaicjiuZGx4m5sEZDTCzxQl48pLv3cKzSI0nuTPmmt0-f896ciPzTH_cbFdvhD7dKEzuQ",
    },
  ];

  const faqs = [
    {
      q: "Apakah ada cicilan?",
      a: "Ya, kami bekerja sama dengan lembaga keuangan syariah terpercaya untuk memberikan fasilitas cicilan hingga 3 tahun dengan proses yang mudah dan tanpa riba.",
    },
    {
      q: "Apa saja dokumen yang diperlukan?",
      a: "Dokumen utama yang diperlukan adalah Paspor (asli), Kartu Keluarga, Buku Nikah (jika suami istri), dan Sertifikat Vaksinasi. Tim kami akan memandu proses pengumpulannya.",
    },
    {
      q: "Kapan jadwal keberangkatan terdekat?",
      a: "Kami memiliki jadwal keberangkatan setiap bulan. Jadwal terdekat yang masih tersedia adalah di bulan depan. Silakan hubungi admin untuk kuota detail.",
    },
  ];
</script>

<svelte:head>
  <title>UmrohOS — Perjalanan Umrah Nyaman dan Terpercaya</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link
    rel="preconnect"
    href="https://fonts.gstatic.com"
    crossorigin="anonymous"
  />
  <link
    href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Manrope:wght@400;500;600;700&display=swap"
    rel="stylesheet"
  />
  <link
    href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
    rel="stylesheet"
  />
</svelte:head>

<div class="home">
  <nav class="top-nav">
    <div class="max shell nav-inner">
      <a class="brand" href="/">UmrohOS</a>
      <div class="nav-links">
        <a class="active" href="/">Beranda</a>
        <a href="/packages">Paket Umrah</a>
        <a href="#proses-booking">Proses Booking</a>
        <a href="#testimoni">Testimoni</a>
        <a href="#faq">FAQ</a>
      </div>
      <div class="nav-right">
        <a
          class="wa-btn"
          href="https://wa.me/6281200000000"
          target="_blank"
          rel="noreferrer">WhatsApp</a
        >
        <a class="primary-btn" href="/packages">Daftar Sekarang</a>
      </div>
    </div>
  </nav>

  <section class="hero">
    <div class="max shell hero-grid">
      <div class="hero-copy">
        <p class="kicker">Premium Umrah Experience</p>
        <h1>
          Wujudkan Perjalanan Ibadah yang <span>Nyaman</span> dan Terpercaya
        </h1>
        <p>
          Nikmati pengalaman spiritual yang dikurasi secara profesional dengan
          layanan concierge 24/7, akomodasi bintang 5, dan bimbingan ibadah yang
          mutawatir.
        </p>
        <div class="hero-cta">
          <a class="hero-btn primary-btn" href="/packages">Lihat Paket Umrah</a>
          <a
            class="hero-btn ghost-btn"
            href="https://wa.me/6281200000000"
            target="_blank"
            rel="noreferrer"
          >
            <span class="material-symbols-outlined">chat</span>
            Konsultasi WhatsApp
          </a>
        </div>
      </div>
      <div class="hero-media">
        <div class="hero-img">
          <img
            alt="Ka'bah"
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuA_fuk10FByeGtF1-q1WcIPMx5SRMx_vFxAfrNdPJENi_rn9B6XExsOBot9NQuQBOamcu1nBSwrdg7XqdF0y9C40aetZ1vu14LXG7_nsOsd34x3xjO9yKctFWY3tt-EVXwwnCrXg7pg7TZy38wjm0C8lMsNw8eNLOGQrykuJvt3PTDl9eJl_WazFbCKTL3eSznVaJvVgvcd8nzSYy2NmdKAeI1tc1DjRCRL07tdfETV4bVDKJdNZYoasVTAXLqIJcSx8IRv6c_EpQ"
            loading="eager"
          />
        </div>
        <div class="license-card">
          <div class="icon-wrap">
            <span class="material-symbols-outlined">verified_user</span>
          </div>
          <div>
            <p>Izin Resmi PPIU</p>
            <strong>No. 123/2024</strong>
          </div>
        </div>
      </div>
    </div>
    <div class="hero-gradient"></div>
  </section>

  <section class="trust-strip">
    <div class="max shell trust-grid">
      {#each trustItems as item, idx (item.title)}
        <article class="trust-item" class:with-border={idx === 1}>
          <h3>{item.title}</h3>
          <p>{item.subtitle}</p>
        </article>
      {/each}
    </div>
  </section>

  <section class="featured">
    <div class="max shell">
      <div class="section-head">
        <h2>Paket Umrah Pilihan</h2>
        <p>
          Pilih paket yang sesuai dengan kebutuhan perjalanan spiritual Anda
          dengan transparansi harga dan fasilitas.
        </p>
      </div>
      <div class="packages-grid">
        {#each packages as item (item.name)}
          <article class="pkg-card" class:popular={item.popular}>
            {#if item.popular}
              <p class="popular-tag">Paling Populer</p>
            {/if}
            <p class="tier">{item.tier}</p>
            <h3>{item.name}</h3>
            <p class="price">{item.price}<span>/pax</span></p>
            <ul>
              {#each item.highlights as highlight (highlight)}
                <li>
                  <span class="material-symbols-outlined">check_circle</span>
                  {highlight}
                </li>
              {/each}
            </ul>
            <a class="pkg-btn" class:solid={item.popular} href="/packages"
              >{item.cta}</a
            >
          </article>
        {/each}
      </div>
    </div>
  </section>

  <section id="proses-booking" class="process">
    <div class="max shell">
      <div class="process-head">
        <div>
          <h2>Proses Booking Mudah</h2>
          <p>Hanya butuh 3 langkah sederhana untuk mengamankan seat Anda.</p>
        </div>
        <a href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
          Bantuan 24/7
          <span class="material-symbols-outlined">arrow_forward</span>
        </a>
      </div>
      <div class="process-grid">
        {#each bookingSteps as step, idx (step.title)}
          <article>
            <div class="step-num">0{idx + 1}</div>
            <h4>{step.title}</h4>
            <p>{step.desc}</p>
          </article>
        {/each}
      </div>
    </div>
  </section>

  <section id="testimoni" class="testimonials">
    <div class="max shell">
      <h2>Apa Kata Jamaah Kami</h2>
      <div class="testi-grid">
        {#each testimonials as item (item.author)}
          <blockquote>
            <span class="quote-mark">"</span>
            <p>{item.quote}</p>
            <div class="author">
              <img alt={item.author} src={item.avatar} loading="lazy" />
              <div>
                <strong>{item.author}</strong>
                <small>{item.package}</small>
              </div>
            </div>
          </blockquote>
        {/each}
      </div>
    </div>
  </section>

  <section id="faq" class="faq">
    <div class="max shell faq-shell">
      <h2>Pertanyaan Umum</h2>
      <div class="faq-list">
        {#each faqs as item (item.q)}
          <details>
            <summary>
              <h4>{item.q}</h4>
              <span class="material-symbols-outlined">expand_more</span>
            </summary>
            <p>{item.a}</p>
          </details>
        {/each}
      </div>
    </div>
  </section>

  <footer class="footer">
    <div class="max shell footer-grid">
      <div class="brand-col">
        <h4>UmrohOS</h4>
        <p>
          Penyedia layanan Umrah digital yang mengutamakan transparansi,
          keamanan, dan kenyamanan ibadah Anda.
        </p>
      </div>
      <div>
        <h5>Layanan</h5>
        <a href="/packages">Paket Umrah 2024</a>
        <a href="/packages">Umrah Plus Turki</a>
        <a href="/packages">Tabungan Umrah</a>
        <a href="/packages">Corporate Travel</a>
      </div>
      <div>
        <h5>Perusahaan</h5>
        <a href="/">Tentang Kami</a>
        <a href="/">Karir</a>
        <a href="/">Kebijakan Privasi</a>
        <a href="/">Syarat & Ketentuan</a>
      </div>
      <div>
        <h5>Kontak Kami</h5>
        <p>info@umrohos.com</p>
        <p>+62 812 3456 7890</p>
        <p>Jl. Sudirman No. 123, Jakarta Selatan</p>
      </div>
    </div>
  </footer>
</div>

<style>
  .home {
    background: #fbf9f8;
    color: #1b1c1c;
    font-family: "Manrope", sans-serif;
  }
  .home h1,
  .home h2,
  .home h3,
  .home h4,
  .home h5,
  .brand {
    font-family: "Plus Jakarta Sans", sans-serif;
  }
  .material-symbols-outlined {
    font-variation-settings:
      "FILL" 0,
      "wght" 400,
      "GRAD" 0,
      "opsz" 24;
    font-size: 1.2rem;
    line-height: 1;
  }
  .max {
    max-width: 80rem;
    margin: 0 auto;
  }
  .shell {
    padding: 0 1.5rem;
  }
  .top-nav {
    position: fixed;
    top: 0;
    z-index: 50;
    width: 100%;
    border-bottom: 1px solid rgba(220, 220, 220, 0.5);
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(20px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
  }
  .nav-inner {
    height: 5rem;
    display: grid;
    grid-template-columns: auto 1fr auto;
    align-items: center;
    gap: 2rem;
  }
  .brand {
    position: static;
    transform: none;
    font-size: 1.6rem;
    font-weight: 800;
    color: #0f3f30;
    text-decoration: none;
  }
  .nav-links,
  .nav-right {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }
  .nav-links {
    justify-self: center;
    white-space: nowrap;
  }
  .nav-links a {
    color: #6b7280;
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s ease;
  }
  .nav-links a:hover {
    color: #065f46;
  }
  .nav-links .active {
    color: #064e3b;
    font-weight: 700;
    border-bottom: 2px solid #b45309;
    padding-bottom: 0.15rem;
  }
  .wa-btn {
    color: #775a19;
    font-weight: 700;
    text-decoration: none;
  }
  .nav-right {
    justify-self: end;
    white-space: nowrap;
  }
  .primary-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 0.75rem;
    padding: 0.68rem 1.4rem;
    font-weight: 700;
    text-decoration: none;
    color: #fff;
    background: linear-gradient(90deg, #004d34, #006747);
  }
  .hero {
    position: relative;
    overflow: hidden;
    padding: 8rem 0 5.5rem;
  }
  .hero-grid {
    position: relative;
    z-index: 2;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 3rem;
    align-items: center;
  }
  .kicker {
    display: inline-block;
    margin: 0;
    border-radius: 999px;
    background: #fed488;
    color: #5d4201;
    font-size: 0.78rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    padding: 0.4rem 0.9rem;
  }
  .hero-copy h1 {
    margin: 1rem 0;
    font-size: clamp(2.4rem, 4.7vw, 4.6rem);
    line-height: 1.1;
    letter-spacing: -0.03em;
  }
  .hero-copy h1 span {
    color: #004d34;
  }
  .hero-copy p {
    margin: 0;
    max-width: 38rem;
    color: #3f4943;
    font-size: 1.14rem;
    line-height: 1.75;
  }
  .hero-cta {
    margin-top: 2rem;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }
  .hero-btn {
    min-height: 3.45rem;
    padding: 0.9rem 1.6rem;
    border-radius: 0.75rem;
    text-decoration: none;
    font-size: 1.08rem;
    font-weight: 700;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.45rem;
  }
  .ghost-btn {
    border: 1px solid #bec9c1;
    color: #1b1c1c;
    background: #fff;
  }
  .hero-media {
    position: relative;
  }
  .hero-img {
    border-radius: 2rem;
    overflow: hidden;
    box-shadow: 0 26px 42px rgba(17, 24, 39, 0.22);
    transform: rotate(2deg);
    transition: transform 0.5s ease;
  }
  .hero-img:hover {
    transform: rotate(0deg);
  }
  .hero-img img {
    width: 100%;
    height: 40rem;
    object-fit: cover;
    display: block;
  }
  .license-card {
    position: absolute;
    left: -1.8rem;
    bottom: -2rem;
    border: 1px solid #e4e2e2;
    border-radius: 1rem;
    background: #fff;
    padding: 1rem 1.1rem;
    box-shadow: 0 14px 26px rgba(0, 0, 0, 0.12);
    display: flex;
    align-items: center;
    gap: 0.8rem;
  }
  .icon-wrap {
    width: 3rem;
    height: 3rem;
    border-radius: 999px;
    background: #775a1919;
    display: grid;
    place-items: center;
    color: #775a19;
  }
  .license-card p {
    margin: 0;
    font-size: 0.78rem;
    color: #3f4943;
  }
  .license-card strong {
    font-size: 0.92rem;
  }
  .hero-gradient {
    position: absolute;
    top: 0;
    right: 0;
    width: 35%;
    height: 100%;
    background: linear-gradient(270deg, #a0f4ca33, transparent);
    pointer-events: none;
  }
  .trust-strip {
    background: #f5f3f3;
    padding: 3rem 0;
  }
  .trust-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
    text-align: center;
    align-items: center;
  }
  .trust-item h3 {
    margin: 0;
    color: #004d34;
    font-size: 2rem;
    font-weight: 800;
  }
  .trust-item p {
    margin: 0.35rem 0 0;
    color: #3f4943;
    font-weight: 500;
  }
  .trust-item.with-border {
    border-left: 1px solid #bec9c1;
    border-right: 1px solid #bec9c1;
  }
  .featured {
    padding: 6rem 0;
  }
  .section-head {
    text-align: center;
    margin-bottom: 3.8rem;
  }
  .section-head h2 {
    margin: 0;
    font-size: clamp(2rem, 4vw, 3.2rem);
    font-weight: 800;
  }
  .section-head p {
    margin: 0.8rem auto 0;
    max-width: 40rem;
    color: #3f4943;
    font-size: 1.12rem;
  }
  .packages-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
  }
  .pkg-card {
    border: 1px solid #e4e2e2;
    border-radius: 2rem;
    background: #fff;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    box-shadow: 0 4px 10px rgba(17, 24, 39, 0.05);
    transition: box-shadow 0.2s ease;
  }
  .pkg-card:hover {
    box-shadow: 0 16px 28px rgba(17, 24, 39, 0.1);
  }
  .pkg-card.popular {
    position: relative;
    border: 2px solid #004d34;
    box-shadow: 0 18px 36px rgba(17, 24, 39, 0.13);
    transform: scale(1.05);
    z-index: 2;
  }
  .popular-tag {
    position: absolute;
    top: -1rem;
    left: 50%;
    transform: translateX(-50%);
    margin: 0;
    border-radius: 999px;
    padding: 0.4rem 1.1rem;
    background: #004d34;
    color: #fff;
    font-size: 0.82rem;
    font-weight: 700;
  }
  .tier {
    margin: 0;
    color: #6b7280;
    text-transform: uppercase;
    letter-spacing: 0.12em;
    font-size: 0.76rem;
    font-weight: 700;
  }
  .pkg-card h3 {
    margin: 0.8rem 0 0;
    font-size: 2rem;
    font-weight: 800;
  }
  .price {
    margin: 1.2rem 0 0;
    color: #004d34;
    font-size: 2.5rem;
    font-weight: 800;
  }
  .price span {
    color: #6b7280;
    font-size: 1rem;
    margin-left: 0.25rem;
    font-weight: 500;
  }
  .pkg-card ul {
    list-style: none;
    margin: 1.5rem 0 2rem;
    padding: 0;
    display: grid;
    gap: 0.9rem;
    flex: 1;
  }
  .pkg-card li {
    display: flex;
    gap: 0.55rem;
    align-items: flex-start;
    color: #3f4943;
  }
  .pkg-card li .material-symbols-outlined {
    color: #004d34;
    font-size: 1rem;
    margin-top: 0.15rem;
  }
  .pkg-btn {
    min-height: 3.25rem;
    border-radius: 0.75rem;
    border: 1px solid #004d34;
    color: #004d34;
    font-weight: 700;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
  .pkg-btn.solid {
    color: #fff;
    background: #004d34;
  }
  .process {
    padding: 6rem 0;
    background: #efeded;
  }
  .process-head {
    display: flex;
    align-items: end;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 4rem;
  }
  .process-head h2 {
    margin: 0;
    font-size: clamp(2rem, 4vw, 3rem);
    font-weight: 800;
  }
  .process-head p {
    margin: 0.8rem 0 0;
    color: #3f4943;
    font-size: 1.1rem;
  }
  .process-head a {
    color: #004d34;
    font-weight: 700;
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    text-decoration: none;
  }
  .process-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 3rem;
  }
  .step-num {
    width: 4rem;
    height: 4rem;
    border-radius: 1rem;
    background: #004d34;
    color: #fff;
    display: grid;
    place-items: center;
    font-size: 1.4rem;
    font-weight: 700;
    margin-bottom: 1.5rem;
  }
  .process-grid h4 {
    margin: 0 0 0.9rem;
    font-size: 1.65rem;
    font-weight: 700;
  }
  .process-grid p {
    margin: 0;
    color: #3f4943;
    line-height: 1.7;
  }
  .testimonials {
    padding: 6rem 0;
  }
  .testimonials h2 {
    margin: 0 0 4rem;
    text-align: center;
    font-size: clamp(2rem, 4vw, 3rem);
    font-weight: 800;
  }
  .testi-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
  }
  blockquote {
    position: relative;
    margin: 0;
    border-radius: 2rem;
    background: #f5f3f3;
    padding: 2.4rem;
    color: #3f4943;
    font-size: 1.1rem;
    line-height: 1.8;
    font-style: italic;
  }
  .quote-mark {
    position: absolute;
    top: 0.8rem;
    left: 1.3rem;
    color: #004d341a;
    font-size: 5rem;
    font-style: normal;
  }
  blockquote p {
    position: relative;
    z-index: 1;
    margin: 0;
  }
  .author {
    margin-top: 2rem;
    display: flex;
    gap: 0.75rem;
    align-items: center;
    font-style: normal;
  }
  .author img {
    width: 3rem;
    height: 3rem;
    border-radius: 999px;
    object-fit: cover;
  }
  .author strong {
    display: block;
    color: #1b1c1c;
    font-size: 1rem;
  }
  .author small {
    color: #775a19;
    font-size: 0.86rem;
  }
  .faq {
    padding: 6rem 0;
    background: #fff;
  }
  .faq-shell {
    max-width: 56rem;
  }
  .faq h2 {
    margin: 0 0 4rem;
    text-align: center;
    font-size: clamp(2rem, 4vw, 3rem);
    font-weight: 800;
  }
  .faq-list {
    display: grid;
    gap: 1.5rem;
  }
  details {
    border-bottom: 1px solid #bec9c1;
    padding-bottom: 1.5rem;
  }
  summary {
    list-style: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
  }
  summary::-webkit-details-marker {
    display: none;
  }
  summary h4 {
    margin: 0;
    font-size: 1.3rem;
    font-weight: 700;
  }
  summary .material-symbols-outlined {
    color: #004d34;
    transition: transform 0.2s ease;
  }
  details[open] summary .material-symbols-outlined {
    transform: rotate(180deg);
  }
  details p {
    margin: 1rem 0 0;
    color: #3f4943;
    line-height: 1.8;
  }
  .footer {
    border-top: 1px solid #e4e2e2;
    background: #fafafa;
    padding: 5rem 0 2.2rem;
  }
  .footer-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 2.2rem;
  }
  .brand-col h4,
  .footer h5 {
    margin: 0 0 1rem;
    font-size: 1.1rem;
    font-weight: 700;
    color: #0f3f30;
  }
  .footer p,
  .footer a {
    margin: 0 0 0.65rem;
    display: block;
    color: #6b7280;
    text-decoration: none;
  }
  .footer a:hover {
    color: #775a19;
  }
  @media (max-width: 1100px) {
    .nav-links {
      display: none;
    }
    .nav-inner {
      grid-template-columns: auto auto;
    }
    .hero-grid,
    .packages-grid,
    .process-grid,
    .testi-grid,
    .footer-grid {
      grid-template-columns: 1fr;
    }
    .hero-img {
      transform: none;
    }
    .hero-img img {
      height: 28rem;
    }
    .license-card {
      position: static;
      margin-top: 1rem;
      width: fit-content;
    }
    .pkg-card.popular {
      transform: none;
    }
    .trust-grid {
      grid-template-columns: 1fr;
      gap: 1.2rem;
    }
    .trust-item.with-border {
      border: none;
      border-top: 1px solid #bec9c1;
      border-bottom: 1px solid #bec9c1;
      padding: 1.2rem 0;
    }
  }
</style>
