<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  type DzikirTab = 'pagi' | 'sore' | 'setelah_shalat';

  interface DzikirItem {
    id: string;
    arabic: string;
    latin: string;
    meaning: string;
    target: number;
  }

  const dzikirCards: DzikirItem[] = [
    { id: 'subhanallah', arabic: 'سُبْحَانَ اللهِ', latin: 'Subhanallah', meaning: 'Maha Suci Allah', target: 33 },
    { id: 'alhamdulillah', arabic: 'الْحَمْدُ لِلَّهِ', latin: 'Alhamdulillah', meaning: 'Segala puji bagi Allah', target: 33 },
    { id: 'allahuakbar', arabic: 'اللهُ أَكْبَرُ', latin: 'Allahu Akbar', meaning: 'Allah Maha Besar', target: 34 },
    { id: 'istighfar', arabic: 'أَسْتَغْفِرُ اللهَ', latin: 'Astaghfirullah', meaning: 'Aku memohon ampun kepada Allah', target: 100 },
    { id: 'shalawat', arabic: 'اللَّهُمَّ صَلِّ عَلَى مُحَمَّدٍ', latin: 'Allahumma shalli ala Muhammad', meaning: 'Ya Allah, limpahkan shalawat atas Muhammad', target: 100 },
  ];

  const dzikirCollections: Record<DzikirTab, { arabic: string; latin: string; source: string }[]> = {
    pagi: [
      { arabic: 'أَصْبَحْنَا وَأَصْبَحَ الْمُلْكُ لِلَّهِ', latin: 'Asbahna wa asbahal mulku lillah', source: 'Dibaca 1x setelah subuh' },
      { arabic: 'اللَّهُمَّ بِكَ أَصْبَحْنَا', latin: 'Allahumma bika asbahna', source: 'Dibaca 1x setelah subuh' },
    ],
    sore: [
      { arabic: 'أَمْسَيْنَا وَأَمْسَى الْمُلْكُ لِلَّهِ', latin: 'Amsayna wa amsal mulku lillah', source: 'Dibaca 1x setelah ashar' },
      { arabic: 'اللَّهُمَّ بِكَ أَمْسَيْنَا', latin: 'Allahumma bika amsayna', source: 'Dibaca 1x setelah ashar' },
    ],
    setelah_shalat: [
      { arabic: 'سُبْحَانَ اللهِ × ٣٣', latin: 'Subhanallah × 33', source: 'Setelah shalat fardhu' },
      { arabic: 'الْحَمْدُ لِلَّهِ × ٣٣', latin: 'Alhamdulillah × 33', source: 'Setelah shalat fardhu' },
      { arabic: 'اللهُ أَكْبَرُ × ٣٤', latin: 'Allahu Akbar × 34', source: 'Setelah shalat fardhu' },
    ],
  };

  let counters = $state<Record<string, number>>(
    Object.fromEntries(dzikirCards.map(d => [d.id, 0]))
  );

  let activeTab = $state<DzikirTab>('pagi');

  function increment(id: string) {
    counters[id] = (counters[id] ?? 0) + 1;
  }

  function resetCounter(id: string) {
    counters[id] = 0;
  }

  function resetAll() {
    counters = Object.fromEntries(dzikirCards.map(d => [d.id, 0]));
  }

  const tabs: Array<{ key: DzikirTab; label: string }> = [
    { key: 'pagi', label: 'Dzikir Pagi' },
    { key: 'sore', label: 'Dzikir Sore' },
    { key: 'setelah_shalat', label: 'Setelah Shalat' },
  ];
</script>

<svelte:head>
  <title>Dzikir — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="dzikir-root">
    <div class="shell">
      <a href="/ibadah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Panduan Ibadah
      </a>
      <div class="page-header">
        <div>
          <h1>Dzikir & Doa</h1>
          <p>Kumpulan dzikir harian dengan penghitung tasbih</p>
        </div>
        <button class="reset-all-btn" onclick={resetAll}>
          <span class="material-symbols-outlined">restart_alt</span>
          Reset Semua
        </button>
      </div>

      <!-- Counter cards -->
      <h2 class="section-title">Tasbih Digital</h2>
      <div class="dzikir-grid">
        {#each dzikirCards as d (d.id)}
          {@const count = counters[d.id] ?? 0}
          {@const done = count >= d.target}
          <div class="dzikir-card" class:done>
            <div class="dz-arabic">{d.arabic}</div>
            <div class="dz-latin">{d.latin}</div>
            <div class="dz-meaning">{d.meaning}</div>
            <div class="dz-counter-row">
              <div class="dz-count" class:done>{count}<span class="dz-target">/{d.target}</span></div>
              <div class="dz-actions">
                <button class="dz-reset-btn" onclick={() => resetCounter(d.id)} title="Reset">
                  <span class="material-symbols-outlined">restart_alt</span>
                </button>
                <button
                  class="dz-tap-btn"
                  class:done
                  onclick={() => increment(d.id)}
                  disabled={done}
                >
                  {#if done}
                    <span class="material-symbols-outlined">check</span>
                  {:else}
                    +1
                  {/if}
                </button>
              </div>
            </div>
            {#if !done}
              <div class="dz-progress-bar-wrap">
                <div class="dz-progress-bar" style="width: {Math.min(100, (count / d.target) * 100)}%"></div>
              </div>
            {/if}
          </div>
        {/each}
      </div>

      <!-- Dzikir collections tabs -->
      <h2 class="section-title">Kumpulan Dzikir</h2>
      <div class="tabs">
        {#each tabs as tab}
          <button
            class="tab-btn"
            class:active={activeTab === tab.key}
            onclick={() => activeTab = tab.key}
          >
            {tab.label}
          </button>
        {/each}
      </div>
      <div class="collection-list">
        {#each dzikirCollections[activeTab] as item, idx}
          <div class="coll-item">
            <div class="coll-arabic">{item.arabic}</div>
            <div class="coll-latin">{item.latin}</div>
            <div class="coll-source">{item.source}</div>
          </div>
        {/each}
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .dzikir-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 72rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; margin-bottom: 2rem; flex-wrap: wrap; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .reset-all-btn { display: inline-flex; align-items: center; gap: 0.4rem; border: 1.5px solid rgba(190,201,193,0.4); background: #fff; color: #6b7280; border-radius: 999px; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; flex-shrink: 0; }
  .reset-all-btn .material-symbols-outlined { font-size: 1rem; }
  .section-title { margin: 0 0 1.2rem; font-size: 1.1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  /* Dzikir grid */
  .dzikir-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 1.25rem; margin-bottom: 3rem; }
  .dzikir-card {
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    transition: box-shadow 0.2s;
  }
  .dzikir-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.06); }
  .dzikir-card.done { border-color: rgba(0,103,71,0.2); background: rgba(0,103,71,0.02); }
  .dz-arabic { font-size: 1.6rem; text-align: right; direction: rtl; color: #1b1c1c; margin-bottom: 0.5rem; font-family: 'Amiri', 'Scheherazade', serif; line-height: 1.8; }
  .dz-latin { font-size: 0.88rem; font-weight: 700; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; margin-bottom: 0.15rem; }
  .dz-meaning { font-size: 0.78rem; color: #9ca3af; margin-bottom: 1rem; }
  .dz-counter-row { display: flex; align-items: center; justify-content: space-between; gap: 0.75rem; margin-bottom: 0.75rem; }
  .dz-count { font-size: 2rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .dz-count.done { color: #006747; }
  .dz-target { font-size: 1rem; color: #9ca3af; font-weight: 400; }
  .dz-actions { display: flex; align-items: center; gap: 0.5rem; }
  .dz-reset-btn { width: 2rem; height: 2rem; border-radius: 50%; border: 1.5px solid rgba(190,201,193,0.4); background: #fff; color: #9ca3af; cursor: pointer; display: grid; place-items: center; }
  .dz-reset-btn .material-symbols-outlined { font-size: 0.9rem; }
  .dz-tap-btn {
    width: 3rem;
    height: 3rem;
    border-radius: 50%;
    background: #006747;
    color: #fff;
    border: none;
    font-size: 0.95rem;
    font-weight: 800;
    cursor: pointer;
    display: grid;
    place-items: center;
    transition: background 0.15s, transform 0.1s;
    font-family: inherit;
  }
  .dz-tap-btn:active { transform: scale(0.94); }
  .dz-tap-btn.done { background: rgba(0,103,71,0.15); color: #006747; cursor: not-allowed; }
  .dz-tap-btn .material-symbols-outlined { font-size: 1.2rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .dz-progress-bar-wrap { height: 4px; background: rgba(190,201,193,0.2); border-radius: 999px; overflow: hidden; }
  .dz-progress-bar { height: 100%; background: linear-gradient(90deg, #004d34, #22c55e); border-radius: 999px; transition: width 0.3s ease; }
  /* Tabs */
  .tabs { display: flex; gap: 0.5rem; margin-bottom: 1.5rem; flex-wrap: wrap; }
  .tab-btn {
    padding: 0.55rem 1.1rem;
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
  .tab-btn.active { background: #006747; color: #fff; border-color: #006747; }
  /* Collection */
  .collection-list { display: flex; flex-direction: column; gap: 0.85rem; }
  .coll-item {
    background: #fff;
    border-radius: 1.2rem;
    padding: 1.25rem 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
  }
  .coll-arabic { font-size: 1.4rem; text-align: right; direction: rtl; color: #1b1c1c; margin-bottom: 0.5rem; font-family: 'Amiri', 'Scheherazade', serif; line-height: 2; }
  .coll-latin { font-size: 0.88rem; color: #57534e; font-style: italic; margin-bottom: 0.25rem; }
  .coll-source { font-size: 0.72rem; color: #9ca3af; }
</style>
