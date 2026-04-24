<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  const categories = ['Semua', 'Panduan Ibadah', 'Persiapan Umrah', 'Tips Perjalanan', 'Kisah Jamaah', 'Info Terbaru'];
  let activeCategory = $state('Semua');

  const articles = [
    {
      id: 'persiapan-fisik-umrah',
      category: 'Persiapan Umrah',
      title: 'Persiapan Fisik yang Perlu Dilakukan 3 Bulan Sebelum Berangkat Umrah',
      excerpt: 'Kondisi fisik prima adalah kunci kelancaran ibadah umrah. Berikut panduan lengkap persiapan fisik yang disarankan para ahli kesehatan dan muthawwif berpengalaman.',
      author: 'Ustadz Wahyu Hidayat',
      date: '15 Januari 2025',
      readTime: '7 menit',
      cover: 'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=800&q=80',
      featured: true
    },
    {
      id: 'doa-tawaf',
      category: 'Panduan Ibadah',
      title: 'Kumpulan Doa dan Zikir Saat Tawaf yang Perlu Dihafal',
      excerpt: 'Tawaf merupakan salah satu rukun umrah yang paling agung. Pelajari doa-doa yang dianjurkan untuk dibaca selama tujuh putaran mengelilingi Ka\'bah.',
      author: 'Ustadz Ahmad Fauzi',
      date: '10 Januari 2025',
      readTime: '10 menit',
      cover: 'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&q=80',
      featured: true
    },
    {
      id: 'pilih-paket-umrah',
      category: 'Persiapan Umrah',
      title: 'Cara Memilih Paket Umrah yang Tepat Sesuai Budget dan Kebutuhan',
      excerpt: 'Bingung memilih antara paket umrah reguler, plus, atau platinum? Artikel ini membantu Anda membuat keputusan berdasarkan prioritas dan anggaran.',
      author: 'Tim UmrohOS',
      date: '5 Januari 2025',
      readTime: '5 menit',
      cover: 'https://images.unsplash.com/photo-1469474968028-56623f02e42e?w=800&q=80',
      featured: false
    },
    {
      id: 'dokumen-umrah',
      category: 'Info Terbaru',
      title: 'Update Persyaratan Dokumen Umrah 2025: Apa yang Berubah?',
      excerpt: 'Pemerintah Arab Saudi memperbarui beberapa persyaratan dokumen untuk calon jamaah umrah tahun 2025. Simak perubahan terbaru dan cara mempersiapkannya.',
      author: 'Tim UmrohOS',
      date: '1 Januari 2025',
      readTime: '4 menit',
      cover: 'https://images.unsplash.com/photo-1450101499163-c8848c66ca85?w=800&q=80',
      featured: false
    },
    {
      id: 'hotel-mekah',
      category: 'Tips Perjalanan',
      title: 'Panduan Memilih Hotel di Makkah: Pertimbangkan Jarak dan Fasilitas',
      excerpt: 'Jarak hotel dari Masjidil Haram sangat mempengaruhi kualitas ibadah. Pelajari kriteria pemilihan hotel yang tepat dan apa bedanya bintang 3, 4, dan 5.',
      author: 'Dewi Santika',
      date: '28 Desember 2024',
      readTime: '6 menit',
      cover: 'https://images.unsplash.com/photo-1564501049412-61c2a3083791?w=800&q=80',
      featured: false
    },
    {
      id: 'kisah-jamaah-pak-bambang',
      category: 'Kisah Jamaah',
      title: 'Kisah Pak Bambang: Umrah di Usia 68 Tahun, Impian yang Akhirnya Terwujud',
      excerpt: 'Setelah menabung selama 12 tahun dan sempat diragukan kondisi fisiknya, Pak Bambang akhirnya dapat menunaikan ibadah umrah dengan nyaman bersama UmrohOS.',
      author: 'Tim Redaksi',
      date: '20 Desember 2024',
      readTime: '8 menit',
      cover: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&q=80',
      featured: false
    }
  ];

  const filteredArticles = $derived(
    activeCategory === 'Semua' ? articles : articles.filter(a => a.category === activeCategory)
  );

  const featuredArticles = $derived(articles.filter(a => a.featured).slice(0, 2));
</script>

<svelte:head>
  <title>Blog & Artikel Umrah — UmrohOS</title>
  <meta name="description" content="Artikel, panduan, dan tips terbaru seputar perjalanan umrah dari tim UmrohOS." />
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="blog-root">

    <!-- Hero -->
    <section class="blog-hero">
      <div class="shell">
        <p class="kicker">Blog & Artikel</p>
        <h1>Inspirasi & Panduan Perjalanan Ibadah</h1>
        <p class="hero-sub">Artikel terkurasi tentang persiapan umrah, panduan ibadah, dan kisah inspiratif jamaah kami.</p>
      </div>
    </section>

    <!-- Featured -->
    {#if featuredArticles.length > 0}
    <section class="featured">
      <div class="shell">
        <h2 class="section-label">Artikel Pilihan</h2>
        <div class="featured-grid">
          {#each featuredArticles as article (article.id)}
            <article class="featured-card">
              <div class="featured-img">
                <img src={article.cover} alt={article.title} loading="eager" />
                <span class="category-tag">{article.category}</span>
              </div>
              <div class="featured-content">
                <h2><a href="/blog/{article.id}">{article.title}</a></h2>
                <p>{article.excerpt}</p>
                <div class="meta">
                  <span class="material-symbols-outlined">person</span>
                  <span>{article.author}</span>
                  <span class="dot">·</span>
                  <span>{article.date}</span>
                  <span class="dot">·</span>
                  <span>{article.readTime} baca</span>
                </div>
              </div>
            </article>
          {/each}
        </div>
      </div>
    </section>
    {/if}

    <!-- All Articles -->
    <section class="articles">
      <div class="shell">
        <div class="articles-header">
          <h2 class="section-label">Semua Artikel</h2>
          <div class="categories" role="group" aria-label="Filter kategori">
            {#each categories as cat (cat)}
              <button
                class="cat-btn"
                class:active={activeCategory === cat}
                onclick={() => activeCategory = cat}
                type="button"
              >{cat}</button>
            {/each}
          </div>
        </div>

        <div class="articles-grid">
          {#each filteredArticles as article (article.id)}
            <article class="article-card">
              <a href="/blog/{article.id}" class="card-img-link">
                <img src={article.cover} alt={article.title} loading="lazy" />
                <span class="category-tag">{article.category}</span>
              </a>
              <div class="card-content">
                <h3><a href="/blog/{article.id}">{article.title}</a></h3>
                <p>{article.excerpt}</p>
                <div class="meta">
                  <span>{article.author}</span>
                  <span class="dot">·</span>
                  <span>{article.date}</span>
                  <span class="dot">·</span>
                  <span>{article.readTime} baca</span>
                </div>
              </div>
            </article>
          {/each}
        </div>

        {#if filteredArticles.length === 0}
          <div class="empty">
            <span class="material-symbols-outlined">article</span>
            <p>Belum ada artikel dalam kategori ini.</p>
          </div>
        {/if}

        <div class="load-more-row">
          <button type="button" class="load-more-btn">
            Muat Artikel Lainnya
            <span class="material-symbols-outlined">expand_more</span>
          </button>
        </div>
      </div>
    </section>

    <!-- Newsletter -->
    <section class="newsletter">
      <div class="shell newsletter-inner">
        <div class="nl-copy">
          <h2>Dapatkan Artikel Terbaru</h2>
          <p>Daftarkan email Anda dan terima panduan umrah, tips perjalanan, serta informasi terbaru langsung di kotak masuk Anda.</p>
        </div>
        <form class="nl-form" onsubmit={(e) => e.preventDefault()}>
          <input type="email" placeholder="Masukkan alamat email Anda" autocomplete="email" required />
          <button type="submit">Daftar Sekarang</button>
        </form>
      </div>
    </section>

  </div>
</MarketingPageLayout>

<style>
  .blog-root {
    padding-top: 5.2rem;
    background: #fbf9f8;
  }
  .shell {
    max-width: 80rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  /* Hero */
  .blog-hero {
    padding: 4.5rem 0 3rem;
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
  .blog-hero h1 {
    margin: 0;
    font-size: clamp(2rem, 4vw, 3.5rem);
    font-weight: 800;
    color: #004d34;
    letter-spacing: -0.02em;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .hero-sub {
    margin: 1rem auto 0;
    max-width: 40rem;
    color: #57534e;
    font-size: 1.1rem;
  }
  /* Featured */
  .featured {
    padding: 4rem 0;
    background: #fff;
  }
  .section-label {
    margin: 0 0 2rem;
    font-size: 1.4rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .featured-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }
  .featured-card {
    border-radius: 1.5rem;
    overflow: hidden;
    border: 1px solid rgba(190,201,193,0.3);
    background: #fbf9f8;
    transition: box-shadow 0.2s;
  }
  .featured-card:hover {
    box-shadow: 0 12px 28px rgba(0,0,0,0.08);
  }
  .featured-img {
    position: relative;
    height: 16rem;
    overflow: hidden;
  }
  .featured-img img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.5s ease;
  }
  .featured-card:hover .featured-img img {
    transform: scale(1.05);
  }
  .category-tag {
    position: absolute;
    top: 1rem;
    left: 1rem;
    background: #006747;
    color: #fff;
    border-radius: 999px;
    padding: 0.3rem 0.75rem;
    font-size: 0.7rem;
    font-weight: 700;
  }
  .featured-content {
    padding: 1.6rem;
  }
  .featured-content h2 {
    margin: 0 0 0.7rem;
    font-size: 1.3rem;
    font-weight: 700;
    line-height: 1.35;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .featured-content h2 a {
    color: #004d34;
    text-decoration: none;
  }
  .featured-content h2 a:hover {
    text-decoration: underline;
  }
  .featured-content p {
    margin: 0 0 1rem;
    color: #57534e;
    line-height: 1.65;
    font-size: 0.95rem;
  }
  /* Articles */
  .articles {
    padding: 4rem 0 5rem;
  }
  .articles-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
    margin-bottom: 2rem;
  }
  .categories {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  .cat-btn {
    border: 1px solid rgba(190,201,193,0.5);
    border-radius: 999px;
    padding: 0.45rem 1rem;
    background: #fff;
    color: #57534e;
    font-size: 0.84rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
  }
  .cat-btn.active {
    background: #006747;
    color: #fff;
    border-color: transparent;
  }
  .articles-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
  }
  .article-card {
    border-radius: 1.5rem;
    overflow: hidden;
    border: 1px solid rgba(190,201,193,0.2);
    background: #fff;
    transition: box-shadow 0.2s;
  }
  .article-card:hover {
    box-shadow: 0 10px 24px rgba(0,0,0,0.08);
  }
  .card-img-link {
    display: block;
    position: relative;
    height: 12rem;
    overflow: hidden;
  }
  .card-img-link img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.5s ease;
  }
  .article-card:hover .card-img-link img {
    transform: scale(1.05);
  }
  .card-content {
    padding: 1.4rem;
  }
  .card-content h3 {
    margin: 0 0 0.6rem;
    font-size: 1rem;
    font-weight: 700;
    line-height: 1.4;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .card-content h3 a {
    color: #004d34;
    text-decoration: none;
  }
  .card-content h3 a:hover {
    text-decoration: underline;
  }
  .card-content p {
    margin: 0 0 0.9rem;
    color: #57534e;
    font-size: 0.88rem;
    line-height: 1.6;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  .meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.78rem;
    color: #9ca3af;
  }
  .meta .material-symbols-outlined {
    font-size: 0.9rem;
  }
  .dot {
    color: #d1d5db;
  }
  .empty {
    text-align: center;
    padding: 4rem 0;
    color: #9ca3af;
  }
  .empty .material-symbols-outlined {
    font-size: 3rem;
    display: block;
    margin-bottom: 1rem;
  }
  .load-more-row {
    text-align: center;
    margin-top: 3rem;
  }
  .load-more-btn {
    border: 1px solid rgba(190,201,193,0.5);
    border-radius: 999px;
    background: #fff;
    color: #004d34;
    font-weight: 700;
    padding: 0.8rem 2rem;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    cursor: pointer;
    transition: background 0.15s;
  }
  .load-more-btn:hover {
    background: #f0f9f4;
  }
  /* Newsletter */
  .newsletter {
    padding: 4rem 0;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
  }
  .newsletter-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 2rem;
    flex-wrap: wrap;
  }
  .nl-copy h2 {
    margin: 0;
    font-size: clamp(1.5rem, 3vw, 2rem);
    font-weight: 800;
    color: #fff;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .nl-copy p {
    margin: 0.7rem 0 0;
    color: rgba(255,255,255,0.8);
    max-width: 36rem;
  }
  .nl-form {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  .nl-form input {
    border: none;
    border-radius: 999px;
    padding: 0.8rem 1.4rem;
    font-size: 0.95rem;
    min-width: 18rem;
    outline: none;
  }
  .nl-form button {
    border: none;
    border-radius: 999px;
    background: #fed488;
    color: #5d4201;
    font-weight: 700;
    padding: 0.8rem 1.8rem;
    cursor: pointer;
    transition: opacity 0.15s;
  }
  .nl-form button:hover {
    opacity: 0.9;
  }
  @media (max-width: 1100px) {
    .featured-grid, .articles-grid {
      grid-template-columns: 1fr 1fr;
    }
  }
  @media (max-width: 640px) {
    .featured-grid, .articles-grid {
      grid-template-columns: 1fr;
    }
    .articles-header {
      flex-direction: column;
      align-items: flex-start;
    }
  }
</style>
