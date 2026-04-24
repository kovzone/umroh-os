<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  type ManasikCategory = 'ihram' | 'thawaf' | 'sai' | 'wuquf' | 'mabit' | 'jumrah' | 'tahallul';

  interface ManasikItem {
    id: string;
    category: ManasikCategory;
    title: string;
    description: string;
  }

  const items: ManasikItem[] = [
    { id: '1', category: 'ihram', title: 'Niat dan Pakaian Ihram', description: 'Ihram adalah niat memulai ibadah umroh atau haji dengan mengucapkan talbiyah. Pakaian ihram pria terdiri dari dua helai kain putih tidak berjahit. Wanita boleh memakai pakaian apapun yang menutup aurat kecuali wajah dan telapak tangan.' },
    { id: '2', category: 'ihram', title: 'Larangan Ihram', description: 'Selama berihram dilarang: memotong kuku/rambut, memakai wewangian, berhubungan suami-istri, berburu hewan, dan bagi pria tidak boleh menutup kepala atau memakai pakaian berjahit.' },
    { id: '3', category: 'ihram', title: 'Miqat Jamaah Indonesia', description: 'Jamaah dari Indonesia biasanya berihram di Bandara Jeddah (miqat zamani) atau di pesawat saat melintas miqat makani. Disunahkan mandi dan memakai wewangian sebelum berihram.' },
    { id: '4', category: 'thawaf', title: 'Tata Cara Thawaf', description: 'Thawaf adalah mengelilingi Ka\'bah sebanyak 7 kali putaran berlawanan arah jarum jam. Dimulai dari sudut Hajar Aswad. Saat melewati Hajar Aswad disunnahkan menciumnya atau mengisyaratkan tangan.' },
    { id: '5', category: 'thawaf', title: 'Jenis-jenis Thawaf', description: 'Thawaf Qudum (saat tiba), Thawaf Umroh (rukun umroh), Thawaf Ifadhah (rukun haji), Thawaf Wada\' (saat akan meninggalkan Makkah), dan Thawaf Sunnah (kapan saja).' },
    { id: '6', category: 'sai', title: 'Tata Cara Sa\'i', description: 'Sa\'i adalah berjalan bolak-balik antara Bukit Shafa dan Marwah sebanyak 7 kali. Dimulai dari Shafa dan berakhir di Marwah. Sa\'i mengenang perjuangan Siti Hajar mencari air untuk putranya Ismail.' },
    { id: '7', category: 'wuquf', title: 'Wuquf di Arafah', description: 'Wuquf di Arafah adalah rukun haji terbesar. Dilaksanakan pada 9 Dzulhijjah dari waktu dzuhur hingga terbenam matahari. Jamaah berdiam diri di padang Arafah sambil berdoa dan berdzikir.' },
    { id: '8', category: 'mabit', title: 'Mabit di Muzdalifah', description: 'Setelah wuquf di Arafah, jamaah berangkat ke Muzdalifah untuk bermalam dan mengambil kerikil. Di sini dilaksanakan shalat Maghrib dan Isya dijama\' dan qashar.' },
    { id: '9', category: 'mabit', title: 'Mabit di Mina', description: 'Bermalam di Mina pada malam Tasyriq (11, 12, atau 13 Dzulhijjah). Termasuk wajib haji bagi yang tidak mendapat dispensasi.' },
    { id: '10', category: 'jumrah', title: 'Melempar Jumrah', description: 'Melempar jumrah adalah melempar kerikil ke tiga tiang (Jumrah Ula, Wustha, dan Aqabah) di Mina. Dilakukan pada hari Idul Adha dan hari Tasyriq. Setiap tiang dilempar 7 kali.' },
    { id: '11', category: 'tahallul', title: 'Tahallul Awal dan Akhir', description: 'Tahallul adalah mencukur/memotong rambut sebagai tanda keluar dari ihram. Tahallul awal terjadi setelah melempar jumrah Aqabah. Tahallul akhir setelah thawaf ifadhah dan sa\'i.' },
  ];

  const categories: Array<{ key: ManasikCategory | 'all'; label: string; icon: string }> = [
    { key: 'all', label: 'Semua', icon: 'apps' },
    { key: 'ihram', label: 'Ihram', icon: 'checkroom' },
    { key: 'thawaf', label: 'Thawaf', icon: 'rotate_right' },
    { key: 'sai', label: 'Sa\'i', icon: 'directions_walk' },
    { key: 'wuquf', label: 'Wuquf', icon: 'landscape' },
    { key: 'mabit', label: 'Mabit', icon: 'hotel' },
    { key: 'jumrah', label: 'Jumrah', icon: 'grain' },
    { key: 'tahallul', label: 'Tahallul', icon: 'content_cut' },
  ];

  let activeCategory = $state<ManasikCategory | 'all'>('all');
  let searchQuery = $state('');
  let expandedId = $state<string | null>(null);

  const filtered = $derived(
    items.filter(item => {
      const matchCat = activeCategory === 'all' || item.category === activeCategory;
      const matchSearch = !searchQuery || item.title.toLowerCase().includes(searchQuery.toLowerCase());
      return matchCat && matchSearch;
    })
  );

  const catColor: Record<ManasikCategory, string> = {
    ihram: '#006747',
    thawaf: '#1565c0',
    sai: '#775a19',
    wuquf: '#c62828',
    mabit: '#2e7d32',
    jumrah: '#6a1c6a',
    tahallul: '#004d34',
  };
</script>

<svelte:head>
  <title>Manasik Umroh & Haji — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="manasik-root">
    <div class="shell">
      <a href="/ibadah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Panduan Ibadah
      </a>
      <div class="page-header">
        <h1>Manasik Umroh & Haji</h1>
        <p>Ensiklopedia tata cara ibadah umroh dan haji</p>
      </div>

      <!-- Search -->
      <div class="search-bar">
        <span class="material-symbols-outlined">search</span>
        <input type="search" placeholder="Cari topik manasik..." bind:value={searchQuery} />
      </div>

      <!-- Category filters -->
      <div class="cat-filters">
        {#each categories as cat}
          <button
            class="cat-btn"
            class:active={activeCategory === cat.key}
            onclick={() => activeCategory = cat.key}
          >
            <span class="material-symbols-outlined">{cat.icon}</span>
            {cat.label}
          </button>
        {/each}
      </div>

      <!-- Items list -->
      <div class="items-list">
        {#each filtered as item (item.id)}
          <div class="manasik-item">
            <button class="manasik-header" onclick={() => expandedId = expandedId === item.id ? null : item.id}>
              <div class="mh-left">
                <span
                  class="cat-dot"
                  style="background: {catColor[item.category]}"
                ></span>
                <span class="mh-title">{item.title}</span>
                <span class="mh-cat" style="color: {catColor[item.category]}">
                  {categories.find(c => c.key === item.category)?.label}
                </span>
              </div>
              <span class="material-symbols-outlined expand-icon" class:open={expandedId === item.id}>
                expand_more
              </span>
            </button>
            {#if expandedId === item.id}
              <div class="manasik-body">
                <p>{item.description}</p>
              </div>
            {/if}
          </div>
        {/each}
        {#if filtered.length === 0}
          <div class="empty-state">
            <span class="material-symbols-outlined">search_off</span>
            <p>Tidak ada topik yang sesuai</p>
          </div>
        {/if}
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .manasik-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 72rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 1.5rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .search-bar {
    display: flex;
    align-items: center;
    background: #fff;
    border: 1.5px solid rgba(190,201,193,0.4);
    border-radius: 999px;
    padding: 0.65rem 1.2rem;
    gap: 0.75rem;
    margin-bottom: 1.25rem;
  }
  .search-bar .material-symbols-outlined { color: #9ca3af; font-size: 1.1rem; flex-shrink: 0; }
  .search-bar input { border: none; outline: none; flex: 1; font-size: 0.92rem; background: transparent; color: #1b1c1c; }
  .cat-filters { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1.5rem; }
  .cat-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.5rem 0.9rem;
    border-radius: 999px;
    border: 1.5px solid rgba(190,201,193,0.3);
    background: #fff;
    font-size: 0.82rem;
    font-weight: 600;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.15s;
    font-family: inherit;
  }
  .cat-btn .material-symbols-outlined { font-size: 0.95rem; }
  .cat-btn.active { background: #006747; color: #fff; border-color: #006747; }
  .items-list { display: flex; flex-direction: column; gap: 0.65rem; }
  .manasik-item { background: #fff; border-radius: 1.2rem; border: 1px solid rgba(190,201,193,0.2); overflow: hidden; }
  .manasik-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 1.1rem 1.3rem;
    background: none;
    border: none;
    cursor: pointer;
    width: 100%;
    font-family: inherit;
    text-align: left;
  }
  .manasik-header:hover { background: #fbf9f8; }
  .mh-left { display: flex; align-items: center; gap: 0.75rem; flex: 1; flex-wrap: wrap; }
  .cat-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
  .mh-title { font-size: 0.92rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .mh-cat { font-size: 0.7rem; font-weight: 700; }
  .expand-icon { color: #9ca3af; font-size: 1.2rem; transition: transform 0.2s; flex-shrink: 0; }
  .expand-icon.open { transform: rotate(180deg); color: #006747; }
  .manasik-body { padding: 0 1.3rem 1.1rem; }
  .manasik-body p { margin: 0; font-size: 0.88rem; color: #57534e; line-height: 1.7; }
  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; padding: 3rem 1rem; color: #9ca3af; }
  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.85rem; }
</style>
