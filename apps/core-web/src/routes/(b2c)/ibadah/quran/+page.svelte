<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  interface Surah {
    number: number;
    name: string;
    arabic: string;
    ayahs: number;
    type: 'Makkiyah' | 'Madaniyah';
  }

  interface Ayah {
    number: number;
    arabic: string;
    latin: string;
    translation: string;
  }

  const surahs: Surah[] = [
    { number: 1, name: 'Al-Fatihah', arabic: 'الفاتحة', ayahs: 7, type: 'Makkiyah' },
    { number: 2, name: 'Al-Baqarah', arabic: 'البقرة', ayahs: 286, type: 'Madaniyah' },
    { number: 3, name: 'Ali Imran', arabic: 'آل عمران', ayahs: 200, type: 'Madaniyah' },
    { number: 4, name: 'An-Nisa', arabic: 'النساء', ayahs: 176, type: 'Madaniyah' },
    { number: 5, name: 'Al-Maidah', arabic: 'المائدة', ayahs: 120, type: 'Madaniyah' },
    { number: 6, name: 'Al-Anam', arabic: 'الأنعام', ayahs: 165, type: 'Makkiyah' },
    { number: 7, name: 'Al-Araf', arabic: 'الأعراف', ayahs: 206, type: 'Makkiyah' },
    { number: 8, name: 'Al-Anfal', arabic: 'الأنفال', ayahs: 75, type: 'Madaniyah' },
    { number: 9, name: 'At-Taubah', arabic: 'التوبة', ayahs: 129, type: 'Madaniyah' },
    { number: 10, name: 'Yunus', arabic: 'يونس', ayahs: 109, type: 'Makkiyah' },
    { number: 11, name: 'Hud', arabic: 'هود', ayahs: 123, type: 'Makkiyah' },
    { number: 12, name: 'Yusuf', arabic: 'يوسف', ayahs: 111, type: 'Makkiyah' },
    { number: 13, name: 'Ar-Rad', arabic: 'الرعد', ayahs: 43, type: 'Madaniyah' },
    { number: 14, name: 'Ibrahim', arabic: 'إبراهيم', ayahs: 52, type: 'Makkiyah' },
    { number: 15, name: 'Al-Hijr', arabic: 'الحجر', ayahs: 99, type: 'Makkiyah' },
    { number: 16, name: 'An-Nahl', arabic: 'النحل', ayahs: 128, type: 'Makkiyah' },
    { number: 17, name: 'Al-Isra', arabic: 'الإسراء', ayahs: 111, type: 'Makkiyah' },
    { number: 18, name: 'Al-Kahf', arabic: 'الكهف', ayahs: 110, type: 'Makkiyah' },
    { number: 19, name: 'Maryam', arabic: 'مريم', ayahs: 98, type: 'Makkiyah' },
    { number: 20, name: 'Ta-Ha', arabic: 'طه', ayahs: 135, type: 'Makkiyah' },
  ];

  // Mock ayahs for Al-Baqarah (first 3)
  const mockAyahs: Record<number, Ayah[]> = {
    1: [
      { number: 1, arabic: 'بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ', latin: 'Bismillahir rahmanir rahim', translation: 'Dengan nama Allah Yang Maha Pengasih, Maha Penyayang.' },
      { number: 2, arabic: 'الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ', latin: 'Alhamdu lillahi rabbil alamin', translation: 'Segala puji bagi Allah, Tuhan seluruh alam,' },
      { number: 3, arabic: 'الرَّحْمَٰنِ الرَّحِيمِ', latin: 'Ar rahmanir rahim', translation: 'Yang Maha Pengasih, Maha Penyayang,' },
    ],
    2: [
      { number: 1, arabic: 'الم', latin: 'Alif Lam Mim', translation: 'Alif Lam Mim.' },
      { number: 2, arabic: 'ذَٰلِكَ الْكِتَابُ لَا رَيْبَ ۛ فِيهِ ۛ هُدًى لِّلْمُتَّقِينَ', latin: 'Zaalikal kitaabu laa rayba feeh, hudal lil muttaqeen', translation: 'Kitab (Al-Qur\'an) ini tidak ada keraguan padanya; petunjuk bagi mereka yang bertakwa,' },
      { number: 255, arabic: 'اللَّهُ لَا إِلَٰهَ إِلَّا هُوَ الْحَيُّ الْقَيُّومُ', latin: 'Allahu la ilaha illa huwal hayyul qayyum', translation: 'Allah, tidak ada tuhan selain Dia. Yang Mahahidup, Yang terus menerus mengurus (makhluk-Nya).' },
    ],
  };

  let searchQuery = $state('');
  let selectedSurah = $state<Surah | null>(null);

  const bookmark = { surahNumber: 2, surahName: 'Al-Baqarah', ayah: 255 };

  const filteredSurahs = $derived(
    searchQuery
      ? surahs.filter(s => s.name.toLowerCase().includes(searchQuery.toLowerCase()) || s.arabic.includes(searchQuery))
      : surahs
  );

  function selectSurah(s: Surah) {
    selectedSurah = s;
  }

  function backToList() {
    selectedSurah = null;
  }

  const currentAyahs = $derived(
    selectedSurah ? (mockAyahs[selectedSurah.number] ?? []) : []
  );
</script>

<svelte:head>
  <title>Al-Quran Digital — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="quran-root">
    <div class="shell">
      <a href="/ibadah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Panduan Ibadah
      </a>

      {#if selectedSurah}
        <!-- Surah view -->
        <div class="surah-header">
          <button class="back-surah-btn" onclick={backToList}>
            <span class="material-symbols-outlined">arrow_back</span>
          </button>
          <div>
            <h1>{selectedSurah.name}</h1>
            <p class="surah-meta">{selectedSurah.arabic} · {selectedSurah.ayahs} ayat · {selectedSurah.type}</p>
          </div>
          <div class="surah-number-badge">{selectedSurah.number}</div>
        </div>

        <div class="ayahs-list">
          {#if currentAyahs.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">menu_book</span>
              <p>Konten surah ini belum tersedia dalam demo.</p>
            </div>
          {:else}
            {#each currentAyahs as ayah (ayah.number)}
              <div class="ayah-card">
                <div class="ayah-number">{ayah.number}</div>
                <div class="ayah-content">
                  <p class="ayah-arabic">{ayah.arabic}</p>
                  <p class="ayah-latin">{ayah.latin}</p>
                  <p class="ayah-translation">{ayah.translation}</p>
                </div>
              </div>
            {/each}
            <div class="ayah-note">Menampilkan {currentAyahs.length} ayat dari {selectedSurah.ayahs} ayat (mode demo)</div>
          {/if}
        </div>

      {:else}
        <!-- Surah list -->
        <div class="page-header">
          <h1>Al-Quran Digital</h1>
          <p>Baca dan cari ayat Al-Quran</p>
        </div>

        <!-- Bookmark -->
        <button class="bookmark-card" onclick={() => selectSurah(surahs.find(s => s.number === bookmark.surahNumber)!)}>
          <span class="material-symbols-outlined bm-icon">bookmark</span>
          <div>
            <div class="bm-label">Terakhir Dibaca</div>
            <div class="bm-value">{bookmark.surahName}: Ayat {bookmark.ayah}</div>
          </div>
          <span class="material-symbols-outlined bm-arrow">arrow_forward</span>
        </button>

        <!-- Search -->
        <div class="search-bar">
          <span class="material-symbols-outlined">search</span>
          <input type="search" placeholder="Cari nama surah..." bind:value={searchQuery} />
          {#if searchQuery}
            <button class="clear-search" onclick={() => searchQuery = ''}>
              <span class="material-symbols-outlined">close</span>
            </button>
          {/if}
        </div>

        <!-- Surah list -->
        <div class="surah-list">
          {#each filteredSurahs as s (s.number)}
            <button class="surah-item" onclick={() => selectSurah(s)}>
              <div class="surah-num-badge">{s.number}</div>
              <div class="surah-info">
                <div class="surah-name">{s.name}</div>
                <div class="surah-details">{s.arabic} · {s.ayahs} ayat · {s.type}</div>
              </div>
              <span class="material-symbols-outlined surah-arrow">chevron_right</span>
            </button>
          {/each}
          {#if filteredSurahs.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">search_off</span>
              <p>Tidak ditemukan surah dengan kata kunci tersebut</p>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
</MarketingPageLayout>

<style>
  .quran-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 64rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 1.5rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  /* Bookmark */
  .bookmark-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    width: 100%;
    background: linear-gradient(135deg, rgba(119,90,25,0.08), rgba(254,212,136,0.1));
    border: 1px solid rgba(119,90,25,0.15);
    border-radius: 1.2rem;
    padding: 1rem 1.2rem;
    margin-bottom: 1.25rem;
    cursor: pointer;
    font-family: inherit;
    text-align: left;
    transition: box-shadow 0.15s;
  }
  .bookmark-card:hover { box-shadow: 0 4px 12px rgba(119,90,25,0.1); }
  .bm-icon { font-size: 1.5rem; color: #775a19; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .bm-label { font-size: 0.72rem; color: #9ca3af; text-transform: uppercase; letter-spacing: 0.07em; }
  .bm-value { font-size: 0.92rem; font-weight: 700; color: #775a19; font-family: 'Plus Jakarta Sans', sans-serif; margin-top: 0.1rem; }
  .bm-arrow { color: #775a19; margin-left: auto; }
  /* Search */
  .search-bar {
    display: flex;
    align-items: center;
    background: #fff;
    border: 1.5px solid rgba(190,201,193,0.4);
    border-radius: 999px;
    padding: 0.65rem 1.2rem;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
    box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  }
  .search-bar .material-symbols-outlined { color: #9ca3af; font-size: 1.1rem; flex-shrink: 0; }
  .search-bar input { border: none; outline: none; flex: 1; font-size: 0.92rem; background: transparent; color: #1b1c1c; }
  .clear-search { border: none; background: none; cursor: pointer; color: #9ca3af; padding: 0; display: flex; }
  /* Surah list */
  .surah-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .surah-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem 1.2rem;
    background: #fff;
    border-radius: 1rem;
    border: 1px solid rgba(190,201,193,0.2);
    cursor: pointer;
    font-family: inherit;
    text-align: left;
    transition: box-shadow 0.15s;
    width: 100%;
  }
  .surah-item:hover { box-shadow: 0 3px 10px rgba(0,0,0,0.06); }
  .surah-num-badge { width: 2.2rem; height: 2.2rem; border-radius: 0.5rem; background: rgba(0,103,71,0.08); display: grid; place-items: center; font-size: 0.78rem; font-weight: 800; color: #006747; flex-shrink: 0; font-family: 'Plus Jakarta Sans', sans-serif; }
  .surah-info { flex: 1; }
  .surah-name { font-size: 0.92rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .surah-details { font-size: 0.75rem; color: #9ca3af; margin-top: 0.1rem; }
  .surah-arrow { color: #d1d5db; font-size: 1.2rem; }
  .surah-item:hover .surah-arrow { color: #006747; }
  /* Surah header */
  .surah-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
  .back-surah-btn { width: 2.5rem; height: 2.5rem; border-radius: 50%; border: 1.5px solid rgba(190,201,193,0.4); background: #fff; cursor: pointer; display: grid; place-items: center; color: #6b7280; flex-shrink: 0; }
  .surah-header h1 { margin: 0; font-size: 1.5rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .surah-meta { margin: 0.2rem 0 0; font-size: 0.82rem; color: #9ca3af; }
  .surah-number-badge { width: 3rem; height: 3rem; border-radius: 0.85rem; background: rgba(0,103,71,0.1); display: grid; place-items: center; font-size: 1.1rem; font-weight: 800; color: #006747; margin-left: auto; flex-shrink: 0; font-family: 'Plus Jakarta Sans', sans-serif; }
  /* Ayahs */
  .ayahs-list { display: flex; flex-direction: column; gap: 1.25rem; }
  .ayah-card {
    display: flex;
    gap: 1rem;
    background: #fff;
    border-radius: 1.2rem;
    padding: 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
  }
  .ayah-number { width: 2rem; height: 2rem; border-radius: 50%; background: rgba(0,103,71,0.08); display: grid; place-items: center; font-size: 0.72rem; font-weight: 800; color: #006747; flex-shrink: 0; margin-top: 0.25rem; font-family: 'Plus Jakarta Sans', sans-serif; }
  .ayah-content { flex: 1; }
  .ayah-arabic { margin: 0 0 0.75rem; font-size: 1.5rem; line-height: 2.2; color: #1b1c1c; direction: rtl; text-align: right; font-family: 'Amiri', 'Scheherazade', serif; }
  .ayah-latin { margin: 0 0 0.4rem; font-size: 0.85rem; color: #6b7280; font-style: italic; }
  .ayah-translation { margin: 0; font-size: 0.88rem; color: #57534e; line-height: 1.6; }
  .ayah-note { text-align: center; font-size: 0.78rem; color: #9ca3af; padding: 1rem; }
  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; padding: 3rem 1rem; color: #9ca3af; }
  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.85rem; }
</style>
