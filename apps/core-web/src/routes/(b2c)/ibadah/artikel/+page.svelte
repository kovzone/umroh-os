<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  type ArticleCategory = 'kajian' | 'kisah' | 'tips' | 'sejarah' | 'doa';

  interface Article {
    id: number;
    title: string;
    category: ArticleCategory;
    date: string;
    excerpt: string;
    readTime: string;
    icon: string;
  }

  const articles: Article[] = [
    { id: 1, title: 'Keutamaan Shalat di Masjidil Haram: Pahala 100.000 Kali', category: 'kajian', date: '24 Jan 2025', excerpt: 'Shalat satu kali di Masjidil Haram bernilai 100.000 kali shalat di masjid lain. Inilah yang membuat setiap langkah menuju masjid suci ini begitu berharga.', readTime: '5 menit', icon: 'mosque' },
    { id: 2, title: 'Kisah Siti Hajar dan Lahirnya Tradisi Sa\'i', category: 'kisah', date: '23 Jan 2025', excerpt: 'Sa\'i yang kita lakukan hari ini adalah warisan perjuangan Siti Hajar, ibu yang penuh tawakal, yang berlari-lari mencari air untuk putranya Ismail.', readTime: '7 menit', icon: 'history_edu' },
    { id: 3, title: '10 Tips Menjaga Stamina Selama Umroh di Makkah', category: 'tips', date: '22 Jan 2025', excerpt: 'Cuaca panas Makkah bisa menguras stamina. Berikut 10 tips praktis dari dokter dan pembimbing perjalanan untuk menjaga kesehatan Anda selama ibadah.', readTime: '4 menit', icon: 'health_and_safety' },
    { id: 4, title: 'Sejarah Ka\'bah: Dari Masa Nabi Ibrahim hingga Kini', category: 'sejarah', date: '21 Jan 2025', excerpt: 'Ka\'bah telah berdiri ribuan tahun dan telah mengalami berbagai renovasi. Berikut perjalanan sejarah bangunan paling suci di muka bumi ini.', readTime: '8 menit', icon: 'account_balance' },
    { id: 5, title: 'Doa-doa Mustajab di Tempat-tempat Suci Makkah', category: 'doa', date: '20 Jan 2025', excerpt: 'Ada beberapa tempat di Makkah yang terkenal sebagai tempat dikabulkannya doa: Multazam, sumur Zamzam, Hajar Aswad, dan Maqam Ibrahim.', readTime: '6 menit', icon: 'auto_awesome' },
    { id: 6, title: 'Makna Spiritual di Balik Setiap Gerakan Thawaf', category: 'kajian', date: '19 Jan 2025', excerpt: 'Setiap putaran dalam thawaf menyimpan makna mendalam. Bukan sekedar ritual, thawaf adalah ekspresi penghambaan manusia kepada Allah yang Maha Esa.', readTime: '6 menit', icon: 'rotate_right' },
  ];

  const filters: Array<{ val: ArticleCategory | 'all'; label: string }> = [
    { val: 'all', label: 'Semua' },
    { val: 'kajian', label: 'Kajian' },
    { val: 'kisah', label: 'Kisah' },
    { val: 'tips', label: 'Tips' },
    { val: 'sejarah', label: 'Sejarah' },
    { val: 'doa', label: 'Doa' },
  ];

  let activeFilter = $state<ArticleCategory | 'all'>('all');

  const filtered = $derived(
    activeFilter === 'all' ? articles : articles.filter(a => a.category === activeFilter)
  );

  const catColor: Record<ArticleCategory, string> = {
    kajian: '#006747',
    kisah: '#775a19',
    tips: '#1565c0',
    sejarah: '#6a1c6a',
    doa: '#c62828',
  };
</script>

<svelte:head>
  <title>Artikel Islami — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="artikel-root">
    <div class="shell">
      <a href="/ibadah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Panduan Ibadah
      </a>
      <div class="page-header">
        <h1>Artikel & Kajian Islami</h1>
        <p>Bacaan islami untuk memperdalam spiritualitas perjalanan umroh Anda</p>
      </div>

      <!-- Filters -->
      <div class="filter-row">
        {#each filters as f}
          <button
            class="filter-btn"
            class:active={activeFilter === f.val}
            onclick={() => activeFilter = f.val}
          >
            {f.label}
          </button>
        {/each}
      </div>

      <!-- Articles grid -->
      <div class="articles-grid">
        {#each filtered as article (article.id)}
          <div class="article-card">
            <div class="article-icon" style="background: color-mix(in srgb, {catColor[article.category]} 12%, transparent); color: {catColor[article.category]}">
              <span class="material-symbols-outlined">{article.icon}</span>
            </div>
            <div class="article-content">
              <div class="article-meta">
                <span class="cat-badge" style="background: color-mix(in srgb, {catColor[article.category]} 12%, transparent); color: {catColor[article.category]}">
                  {filters.find(f => f.val === article.category)?.label}
                </span>
                <span class="article-date">{article.date}</span>
              </div>
              <h3>{article.title}</h3>
              <p>{article.excerpt}</p>
              <div class="article-footer">
                <span class="read-time">
                  <span class="material-symbols-outlined">schedule</span>
                  {article.readTime} baca
                </span>
                <button class="read-more-btn">
                  Baca Selengkapnya
                  <span class="material-symbols-outlined">arrow_forward</span>
                </button>
              </div>
            </div>
          </div>
        {/each}
      </div>

      {#if filtered.length === 0}
        <div class="empty-state">
          <span class="material-symbols-outlined">article</span>
          <p>Tidak ada artikel untuk kategori ini</p>
        </div>
      {/if}
    </div>
  </div>
</MarketingPageLayout>

<style>
  .artikel-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 72rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .filter-row { display: flex; gap: 0.5rem; margin-bottom: 2rem; flex-wrap: wrap; }
  .filter-btn {
    padding: 0.5rem 1rem;
    border-radius: 999px;
    border: 1.5px solid rgba(190,201,193,0.3);
    background: #fff;
    font-size: 0.85rem;
    font-weight: 600;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.15s;
    font-family: inherit;
  }
  .filter-btn.active { background: #006747; color: #fff; border-color: #006747; }
  .articles-grid { display: grid; gap: 1.25rem; }
  .article-card {
    display: flex;
    gap: 1.25rem;
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    transition: box-shadow 0.2s;
    align-items: flex-start;
  }
  .article-card:hover { box-shadow: 0 6px 16px rgba(0,0,0,0.06); }
  .article-icon {
    width: 3.2rem;
    height: 3.2rem;
    border-radius: 1rem;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }
  .article-icon .material-symbols-outlined { font-size: 1.4rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .article-content { flex: 1; }
  .article-meta { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.5rem; }
  .cat-badge { font-size: 0.7rem; font-weight: 700; border-radius: 999px; padding: 0.2rem 0.65rem; }
  .article-date { font-size: 0.75rem; color: #9ca3af; }
  .article-content h3 { margin: 0 0 0.5rem; font-size: 1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; line-height: 1.4; }
  .article-content p { margin: 0 0 1rem; font-size: 0.85rem; color: #57534e; line-height: 1.6; }
  .article-footer { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 0.5rem; }
  .read-time { display: inline-flex; align-items: center; gap: 0.25rem; font-size: 0.78rem; color: #9ca3af; }
  .read-time .material-symbols-outlined { font-size: 0.88rem; }
  .read-more-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.82rem;
    font-weight: 700;
    color: #006747;
    background: rgba(0,103,71,0.08);
    border: none;
    border-radius: 999px;
    padding: 0.4rem 0.9rem;
    cursor: pointer;
    font-family: inherit;
  }
  .read-more-btn .material-symbols-outlined { font-size: 0.9rem; }
  .read-more-btn:hover { background: rgba(0,103,71,0.14); }
  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; padding: 3rem 1rem; color: #9ca3af; }
  .empty-state .material-symbols-outlined { font-size: 3rem; }
  .empty-state p { margin: 0; font-size: 0.9rem; }
</style>
