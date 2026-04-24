<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  type NoticeCategory = 'reuni' | 'acara' | 'info';

  interface Notice {
    id: number;
    title: string;
    date: string;
    category: NoticeCategory;
    description: string;
    icon: string;
  }

  const notices: Notice[] = [
    { id: 1, title: 'Reuni Akbar Alumni Umroh 2025', date: '15 Maret 2025', category: 'reuni', description: 'Reuni akbar seluruh alumni umroh UmrohOS tahun 2025 akan diselenggarakan di Jakarta. Pendaftaran dibuka hingga 1 Maret 2025.', icon: 'groups' },
    { id: 2, title: 'Malam Tasyakuran Kepulangan', date: '28 Januari 2025', category: 'acara', description: 'Malam tasyakuran kepulangan akan diselenggarakan di aula Masjid Al-Huda, Jakarta Selatan. Harap konfirmasi kehadiran kepada pembimbing.', icon: 'celebration' },
    { id: 3, title: 'Pengumuman Keberangkatan Batch Februari', date: '20 Januari 2025', category: 'info', description: 'Batch Februari 2025 masih tersedia beberapa kursi. Ajak keluarga atau tetangga untuk bergabung!', icon: 'campaign' },
    { id: 4, title: 'Halaqah Pasca-Umroh', date: '5 Februari 2025', category: 'acara', description: 'Halaqah ilmu untuk mempertahankan spiritualitas pasca-umroh. Bersama Ustaz Khalid Al-Mustafa. Daring via Zoom.', icon: 'school' },
    { id: 5, title: 'Pertemuan Alumni Wilayah Bandung', date: '10 Februari 2025', category: 'reuni', description: 'Silaturahmi alumni UmrohOS se-Bandung Raya di Masjid Agung Bandung. Gratis untuk semua alumni.', icon: 'handshake' },
    { id: 6, title: 'Update Syarat Visa Umroh 2025', date: '1 Januari 2025', category: 'info', description: 'Pemerintah Saudi Arabia telah memperbarui persyaratan visa umroh. Harap baca informasi lengkap sebelum mendaftar.', icon: 'info' },
  ];

  let activeFilter = $state<NoticeCategory | 'all'>('all');

  const filteredNotices = $derived(
    activeFilter === 'all' ? notices : notices.filter(n => n.category === activeFilter)
  );

  const catLabel: Record<NoticeCategory, string> = { reuni: 'Reuni', acara: 'Acara', info: 'Info' };
  const catColor: Record<NoticeCategory, string> = { reuni: '#1565c0', acara: '#006747', info: '#775a19' };

  const filters: Array<{ val: NoticeCategory | 'all'; label: string }> = [
    { val: 'all', label: 'Semua' },
    { val: 'reuni', label: 'Reuni' },
    { val: 'acara', label: 'Acara' },
    { val: 'info', label: 'Info' },
  ];
</script>

<svelte:head>
  <title>Pengumuman — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="pengumuman-root">
    <div class="shell">
      <a href="/jemaah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Portal Jamaah
      </a>
      <div class="page-header">
        <h1>Papan Pengumuman</h1>
        <p>Informasi terbaru seputar reuni, acara, dan berita alumni</p>
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
            <span class="filter-count">
              {f.val === 'all' ? notices.length : notices.filter(n => n.category === f.val).length}
            </span>
          </button>
        {/each}
      </div>

      <!-- Notice cards -->
      <div class="notices-grid">
        {#each filteredNotices as notice (notice.id)}
          <div class="notice-card">
            <div class="notice-icon" style="background: color-mix(in srgb, {catColor[notice.category]} 12%, transparent); color: {catColor[notice.category]}">
              <span class="material-symbols-outlined">{notice.icon}</span>
            </div>
            <div class="notice-content">
              <div class="notice-meta">
                <span class="cat-badge" style="background: color-mix(in srgb, {catColor[notice.category]} 12%, transparent); color: {catColor[notice.category]}">{catLabel[notice.category]}</span>
                <span class="notice-date">
                  <span class="material-symbols-outlined">event</span>
                  {notice.date}
                </span>
              </div>
              <h3>{notice.title}</h3>
              <p>{notice.description}</p>
            </div>
          </div>
        {/each}
      </div>

      {#if filteredNotices.length === 0}
        <div class="empty-state">
          <span class="material-symbols-outlined">inbox</span>
          <p>Tidak ada pengumuman untuk kategori ini</p>
        </div>
      {/if}
    </div>
  </div>
</MarketingPageLayout>

<style>
  .pengumuman-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 72rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .filter-row { display: flex; gap: 0.5rem; margin-bottom: 2rem; flex-wrap: wrap; }
  .filter-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
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
  .filter-btn.active { background: #006747; color: #fff; border-color: #006747; }
  .filter-btn:hover:not(.active) { border-color: #006747; color: #006747; }
  .filter-count { font-size: 0.72rem; background: rgba(190,201,193,0.2); border-radius: 999px; padding: 0.05rem 0.45rem; }
  .filter-btn.active .filter-count { background: rgba(255,255,255,0.2); }
  .notices-grid { display: grid; gap: 1.25rem; }
  .notice-card {
    display: flex;
    gap: 1.25rem;
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    transition: box-shadow 0.2s;
    align-items: flex-start;
  }
  .notice-card:hover { box-shadow: 0 6px 16px rgba(0,0,0,0.06); }
  .notice-icon {
    width: 3.2rem;
    height: 3.2rem;
    border-radius: 1rem;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }
  .notice-icon .material-symbols-outlined { font-size: 1.4rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .notice-content { flex: 1; }
  .notice-meta { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.5rem; flex-wrap: wrap; }
  .cat-badge { font-size: 0.7rem; font-weight: 700; border-radius: 999px; padding: 0.2rem 0.65rem; }
  .notice-date { display: flex; align-items: center; gap: 0.25rem; font-size: 0.78rem; color: #9ca3af; }
  .notice-date .material-symbols-outlined { font-size: 0.88rem; }
  .notice-content h3 { margin: 0 0 0.5rem; font-size: 1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .notice-content p { margin: 0; font-size: 0.85rem; color: #57534e; line-height: 1.6; }
  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; padding: 4rem 1rem; color: #9ca3af; }
  .empty-state .material-symbols-outlined { font-size: 3rem; }
  .empty-state p { margin: 0; font-size: 0.9rem; }
</style>
