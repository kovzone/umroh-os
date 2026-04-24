<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  const categories = [
    {
      id: 'persiapan',
      icon: 'checklist',
      title: 'Persiapan Umrah',
      desc: 'Semua yang perlu Anda ketahui sebelum berangkat',
      count: 12,
      color: '#006747'
    },
    {
      id: 'dokumen',
      icon: 'description',
      title: 'Dokumen & Visa',
      desc: 'Panduan lengkap pengurusan dokumen dan visa',
      count: 8,
      color: '#775a19'
    },
    {
      id: 'ibadah',
      icon: 'mosque',
      title: 'Tata Cara Ibadah',
      desc: 'Panduan tata cara umrah sesuai sunnah',
      count: 15,
      color: '#004d34'
    },
    {
      id: 'kesehatan',
      icon: 'health_and_safety',
      title: 'Kesehatan Perjalanan',
      desc: 'Tips menjaga kesehatan selama perjalanan',
      count: 9,
      color: '#1565c0'
    },
    {
      id: 'hotel',
      icon: 'hotel',
      title: 'Akomodasi & Hotel',
      desc: 'Panduan memilih dan memaksimalkan akomodasi',
      count: 6,
      color: '#6a1c6a'
    },
    {
      id: 'faq',
      icon: 'help_outline',
      title: 'FAQ & Tips Praktis',
      desc: 'Jawaban atas pertanyaan yang paling sering ditanyakan',
      count: 20,
      color: '#c62828'
    }
  ];

  const popularArticles = [
    {
      id: 'niat-umrah',
      category: 'Tata Cara Ibadah',
      title: 'Bacaan Niat dan Doa Umrah Lengkap (Arab, Latin, Terjemahan)',
      icon: 'mosque',
      readTime: '5 menit'
    },
    {
      id: 'wajib-umrah',
      category: 'Tata Cara Ibadah',
      title: 'Rukun dan Wajib Umrah: Perbedaan dan Konsekuensinya',
      icon: 'mosque',
      readTime: '8 menit'
    },
    {
      id: 'paspor-umrah',
      category: 'Dokumen & Visa',
      title: 'Cara Membuat dan Memperpanjang Paspor untuk Umrah',
      icon: 'description',
      readTime: '6 menit'
    },
    {
      id: 'vaksin-umrah',
      category: 'Kesehatan Perjalanan',
      title: 'Vaksin Wajib dan Dianjurkan untuk Jamaah Umrah 2025',
      icon: 'health_and_safety',
      readTime: '5 menit'
    },
    {
      id: 'ihram',
      category: 'Tata Cara Ibadah',
      title: 'Panduan Lengkap Memakai Pakaian Ihram dan Larangannya',
      icon: 'mosque',
      readTime: '7 menit'
    },
    {
      id: 'miqat',
      category: 'Tata Cara Ibadah',
      title: 'Mengenal Miqat: Batas Tempat Mulai Ihram dari Indonesia',
      icon: 'location_on',
      readTime: '4 menit'
    }
  ];

  const faqs = [
    {
      q: 'Berapa lama proses pengurusan visa umrah?',
      a: 'Proses pengurusan visa umrah biasanya memakan waktu 7–14 hari kerja setelah seluruh dokumen lengkap dan diserahkan. Tim kami memantau status visa secara proaktif.'
    },
    {
      q: 'Apakah ada batasan usia untuk berangkat umrah?',
      a: 'Tidak ada batasan usia minimum untuk umrah. Anak-anak di bawah 12 tahun perlu didampingi mahram. Lansia di atas 70 tahun mungkin memerlukan surat keterangan sehat dari dokter.'
    },
    {
      q: 'Berapa lama durasi perjalanan umrah rata-rata?',
      a: 'Paket umrah reguler UmrohOS berlangsung 9–14 hari, termasuk perjalanan. Program umrah plus (Turki, Mesir, dll.) dapat berlangsung hingga 17 hari.'
    },
    {
      q: 'Apakah bisa umrah sendiri tanpa mahram bagi wanita?',
      a: 'Peraturan Saudi Arabia tahun 2021 membolehkan wanita berumur 45 tahun ke atas untuk berangkat umrah dalam kelompok resmi tanpa mahram, sepanjang bepergian bersama rombongan terdaftar.'
    },
    {
      q: 'Bagaimana jika paspor saya kurang dari 6 bulan masa berlakunya?',
      a: 'Paspor harus memiliki masa berlaku minimal 6 bulan sejak tanggal keberangkatan. Pastikan memperbarui paspor Anda sebelum mendaftar umrah.'
    }
  ];

  let openFaq = $state<number | null>(null);
  let searchQuery = $state('');

  function toggleFaq(idx: number) {
    openFaq = openFaq === idx ? null : idx;
  }
</script>

<svelte:head>
  <title>Panduan Umrah — UmrohOS</title>
  <meta name="description" content="Panduan lengkap persiapan dan tata cara ibadah umrah: dokumen, visa, kesehatan, doa, dan tips perjalanan dari tim UmrohOS." />
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="kb-root">

    <!-- Hero -->
    <section class="kb-hero">
      <div class="shell">
        <p class="kicker">Pusat Panduan</p>
        <h1>Panduan Lengkap Umrah</h1>
        <p class="hero-sub">Semua informasi yang Anda butuhkan untuk perjalanan ibadah yang khusyuk dan lancar, dari awal hingga kembali ke tanah air.</p>
        <div class="search-bar">
          <span class="material-symbols-outlined">search</span>
          <input
            type="search"
            placeholder="Cari panduan, contoh: 'cara ihram' atau 'persyaratan visa'..."
            bind:value={searchQuery}
          />
          {#if searchQuery}
            <button class="clear-btn" type="button" onclick={() => searchQuery = ''}>
              <span class="material-symbols-outlined">close</span>
            </button>
          {/if}
        </div>
      </div>
    </section>

    <!-- Categories -->
    <section class="categories-section">
      <div class="shell">
        <h2 class="section-title">Pilih Topik Panduan</h2>
        <div class="categories-grid">
          {#each categories as cat (cat.id)}
            <a href="/panduan/{cat.id}" class="cat-card" style="--cat-color: {cat.color}">
              <div class="cat-icon">
                <span class="material-symbols-outlined">{cat.icon}</span>
              </div>
              <div class="cat-info">
                <h3>{cat.title}</h3>
                <p>{cat.desc}</p>
                <span class="article-count">{cat.count} artikel</span>
              </div>
              <span class="cat-arrow material-symbols-outlined">arrow_forward</span>
            </a>
          {/each}
        </div>
      </div>
    </section>

    <!-- Popular Articles -->
    <section class="popular-section">
      <div class="shell">
        <h2 class="section-title">Panduan Paling Banyak Dibaca</h2>
        <div class="popular-grid">
          {#each popularArticles as article (article.id)}
            <a href="/panduan/{article.id}" class="popular-card">
              <div class="popular-icon">
                <span class="material-symbols-outlined">{article.icon}</span>
              </div>
              <div class="popular-info">
                <span class="popular-category">{article.category}</span>
                <h3>{article.title}</h3>
                <span class="read-time">
                  <span class="material-symbols-outlined">schedule</span>
                  {article.readTime} baca
                </span>
              </div>
              <span class="popular-arrow material-symbols-outlined">chevron_right</span>
            </a>
          {/each}
        </div>
      </div>
    </section>

    <!-- FAQ -->
    <section class="faq-section">
      <div class="shell">
        <div class="faq-inner">
          <div class="faq-head">
            <h2>Pertanyaan yang Sering Ditanyakan</h2>
            <p>Jawaban cepat untuk pertanyaan umum seputar perjalanan umrah</p>
          </div>
          <div class="faq-list">
            {#each faqs as item, idx (item.q)}
              <details class="faq-item" open={openFaq === idx}>
                <summary onclick={(e) => { e.preventDefault(); toggleFaq(idx); }}>
                  <span>{item.q}</span>
                  <span class="faq-arrow material-symbols-outlined" class:open={openFaq === idx}>expand_more</span>
                </summary>
                {#if openFaq === idx}
                  <div class="faq-answer">
                    <p>{item.a}</p>
                  </div>
                {/if}
              </details>
            {/each}
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="cta-section">
      <div class="shell cta-inner">
        <div>
          <h2>Masih Ada Pertanyaan?</h2>
          <p>Tim konsultan kami siap membantu Anda 24/7 melalui WhatsApp.</p>
        </div>
        <div class="cta-btns">
          <a class="btn-wa" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
            <span class="material-symbols-outlined">chat</span>
            Chat WhatsApp
          </a>
          <a class="btn-packages" href="/packages">Lihat Paket Umrah</a>
        </div>
      </div>
    </section>

  </div>
</MarketingPageLayout>

<style>
  .kb-root {
    padding-top: 5.2rem;
    background: #fbf9f8;
  }
  .shell {
    max-width: 80rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  /* Hero */
  .kb-hero {
    padding: 4.5rem 0 3.5rem;
    background: linear-gradient(180deg, #f0f9f4 0%, #fbf9f8 100%);
    text-align: center;
  }
  .kicker {
    display: inline-block;
    margin: 0 0 1rem;
    background: rgba(254,212,136,0.3);
    color: #775a19;
    border-radius: 999px;
    padding: 0.4rem 1rem;
    font-size: 0.76rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .kb-hero h1 {
    margin: 0;
    font-size: clamp(2rem, 4vw, 3.5rem);
    font-weight: 800;
    color: #004d34;
    letter-spacing: -0.02em;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .hero-sub {
    margin: 1rem auto 2rem;
    max-width: 42rem;
    color: #57534e;
    font-size: 1.05rem;
    line-height: 1.7;
  }
  .search-bar {
    position: relative;
    max-width: 44rem;
    margin: 0 auto;
    display: flex;
    align-items: center;
    background: #fff;
    border: 1.5px solid rgba(190,201,193,0.5);
    border-radius: 999px;
    padding: 0.75rem 1.2rem;
    gap: 0.75rem;
    box-shadow: 0 4px 16px rgba(0,0,0,0.06);
  }
  .search-bar .material-symbols-outlined {
    color: #9ca3af;
    font-size: 1.2rem;
    flex-shrink: 0;
  }
  .search-bar input {
    border: none;
    outline: none;
    flex: 1;
    font-size: 0.95rem;
    color: #1b1c1c;
    background: transparent;
  }
  .clear-btn {
    border: none;
    background: none;
    cursor: pointer;
    color: #9ca3af;
    padding: 0;
    display: flex;
    align-items: center;
  }
  /* Categories */
  .categories-section {
    padding: 4rem 0;
    background: #fff;
  }
  .section-title {
    margin: 0 0 2.5rem;
    font-size: 1.6rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .categories-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1.5rem;
  }
  .cat-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.5rem;
    border-radius: 1.5rem;
    border: 1px solid rgba(190,201,193,0.3);
    background: #fbf9f8;
    text-decoration: none;
    transition: box-shadow 0.2s, transform 0.2s;
  }
  .cat-card:hover {
    box-shadow: 0 10px 24px rgba(0,0,0,0.08);
    transform: translateY(-2px);
  }
  .cat-icon {
    width: 3.2rem;
    height: 3.2rem;
    border-radius: 1rem;
    background: color-mix(in srgb, var(--cat-color) 12%, transparent);
    display: grid;
    place-items: center;
    flex-shrink: 0;
    color: var(--cat-color);
  }
  .cat-icon .material-symbols-outlined {
    font-size: 1.5rem;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .cat-info {
    flex: 1;
  }
  .cat-info h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .cat-info p {
    margin: 0.2rem 0 0.4rem;
    font-size: 0.82rem;
    color: #6b7280;
    line-height: 1.4;
  }
  .article-count {
    font-size: 0.76rem;
    font-weight: 600;
    color: var(--cat-color);
  }
  .cat-arrow {
    color: #d1d5db;
    font-size: 1.2rem;
    transition: color 0.2s;
  }
  .cat-card:hover .cat-arrow {
    color: var(--cat-color);
  }
  /* Popular */
  .popular-section {
    padding: 4rem 0;
  }
  .popular-grid {
    display: grid;
    gap: 1rem;
  }
  .popular-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.2rem 1.5rem;
    border-radius: 1.2rem;
    border: 1px solid rgba(190,201,193,0.2);
    background: #fff;
    text-decoration: none;
    transition: box-shadow 0.2s;
  }
  .popular-card:hover {
    box-shadow: 0 6px 16px rgba(0,0,0,0.06);
  }
  .popular-icon {
    width: 2.8rem;
    height: 2.8rem;
    border-radius: 0.75rem;
    background: rgba(0,103,71,0.1);
    display: grid;
    place-items: center;
    flex-shrink: 0;
    color: #006747;
  }
  .popular-icon .material-symbols-outlined {
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .popular-info {
    flex: 1;
  }
  .popular-category {
    font-size: 0.7rem;
    font-weight: 700;
    color: #775a19;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .popular-info h3 {
    margin: 0.2rem 0 0.4rem;
    font-size: 0.95rem;
    font-weight: 600;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
    line-height: 1.35;
  }
  .read-time {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.78rem;
    color: #9ca3af;
  }
  .read-time .material-symbols-outlined {
    font-size: 0.9rem;
  }
  .popular-arrow {
    color: #d1d5db;
    transition: color 0.2s;
  }
  .popular-card:hover .popular-arrow {
    color: #006747;
  }
  /* FAQ */
  .faq-section {
    padding: 4rem 0;
    background: #fff;
  }
  .faq-inner {
    display: grid;
    grid-template-columns: 1fr 1.8fr;
    gap: 4rem;
    align-items: start;
  }
  .faq-head h2 {
    margin: 0;
    font-size: 1.6rem;
    font-weight: 700;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
    line-height: 1.3;
  }
  .faq-head p {
    margin: 0.8rem 0 0;
    color: #6b7280;
    line-height: 1.65;
  }
  .faq-list {
    display: grid;
    gap: 0.75rem;
  }
  .faq-item {
    border-radius: 1rem;
    border: 1px solid rgba(190,201,193,0.3);
    background: #fbf9f8;
    overflow: hidden;
  }
  .faq-item summary {
    list-style: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 1.1rem 1.3rem;
    font-weight: 600;
    color: #1b1c1c;
  }
  .faq-item summary::-webkit-details-marker { display: none; }
  .faq-arrow {
    flex-shrink: 0;
    color: #9ca3af;
    transition: transform 0.2s;
  }
  .faq-arrow.open {
    transform: rotate(180deg);
    color: #006747;
  }
  .faq-answer {
    padding: 0 1.3rem 1.1rem;
  }
  .faq-answer p {
    margin: 0;
    color: #57534e;
    line-height: 1.7;
  }
  /* CTA */
  .cta-section {
    padding: 4rem 0;
    background: #f0f9f4;
  }
  .cta-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 2rem;
    flex-wrap: wrap;
  }
  .cta-inner h2 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .cta-inner p {
    margin: 0.5rem 0 0;
    color: #57534e;
  }
  .cta-btns {
    display: flex;
    gap: 0.9rem;
    flex-shrink: 0;
    flex-wrap: wrap;
  }
  .btn-wa {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    background: #006747;
    color: #fff;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.8rem 1.8rem;
  }
  .btn-wa .material-symbols-outlined {
    font-size: 1rem;
  }
  .btn-packages {
    text-decoration: none;
    border: 1.5px solid #006747;
    color: #006747;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.8rem 1.8rem;
  }
  @media (max-width: 1100px) {
    .categories-grid {
      grid-template-columns: 1fr 1fr;
    }
    .faq-inner {
      grid-template-columns: 1fr;
      gap: 2rem;
    }
  }
  @media (max-width: 640px) {
    .categories-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
