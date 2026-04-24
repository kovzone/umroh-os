<script lang="ts">
  type AcadSection = 'courses' | 'leaderboard' | 'scripts' | 'badges';

  const SECTIONS: Array<{ id: AcadSection; label: string; icon: string }> = [
    { id: 'courses', label: 'Kursus & Modul', icon: 'school' },
    { id: 'leaderboard', label: 'Leaderboard', icon: 'leaderboard' },
    { id: 'scripts', label: 'Skrip Penjualan', icon: 'chat' },
    { id: 'badges', label: 'Badge & Sertifikasi', icon: 'workspace_premium' }
  ];

  let activeSection = $state<AcadSection>('courses');

  // --- Courses ---
  const courses = $state([
    { id: 'co1', title: 'Dasar-Dasar Menjual Paket Umroh', category: 'Sales', duration: '2.5 jam', lessons: 8, progress: 100, enrolled: 112, badge: true },
    { id: 'co2', title: 'Teknik Follow-Up yang Efektif', category: 'Sales', duration: '1.5 jam', lessons: 5, progress: 75, enrolled: 88, badge: false },
    { id: 'co3', title: 'Digital Marketing untuk Agen Travel', category: 'Marketing', duration: '3 jam', lessons: 10, progress: 40, enrolled: 65, badge: true },
    { id: 'co4', title: 'Pengelolaan Downline & Tim', category: 'Manajemen', duration: '2 jam', lessons: 6, progress: 0, enrolled: 41, badge: false },
    { id: 'co5', title: 'Syarat & Fiqih Umroh untuk Agen', category: 'Pengetahuan', duration: '1 jam', lessons: 4, progress: 100, enrolled: 135, badge: true },
    { id: 'co6', title: 'Penanganan Keberatan Calon Jamaah', category: 'Sales', duration: '1 jam', lessons: 3, progress: 20, enrolled: 54, badge: false }
  ]);

  let courseFilter = $state('semua');
  const CATEGORIES = ['semua', 'Sales', 'Marketing', 'Manajemen', 'Pengetahuan'];
  const filteredCourses = $derived(
    courseFilter === 'semua' ? courses : courses.filter(c => c.category === courseFilter)
  );

  // --- Leaderboard ---
  const leaderboard = $state([
    { rank: 1, agent: 'PT Mitra Barokah Tours', score: 9_840, courses_done: 6, tier: 'Gold', trend: 'up' },
    { rank: 2, agent: 'Karima Travel Group', score: 9_120, courses_done: 5, tier: 'Gold', trend: 'stable' },
    { rank: 3, agent: 'CV Cahaya Umroh', score: 7_650, courses_done: 4, tier: 'Silver', trend: 'up' },
    { rank: 4, agent: 'Ustadz Fahmi Network', score: 6_400, courses_done: 4, tier: 'Silver', trend: 'down' },
    { rank: 5, agent: 'Naufal Berkah Tour', score: 5_100, courses_done: 3, tier: 'Bronze', trend: 'up' },
    { rank: 6, agent: 'Rina Hadiati Agency', score: 3_800, courses_done: 2, tier: 'Bronze', trend: 'stable' }
  ]);

  // --- Scripts ---
  let scriptSearch = $state('');
  const scripts = $state([
    { id: 'sc1', title: 'Pembuka WA — Calon Jamaah Baru', category: 'WhatsApp', tags: ['pembuka', 'wa'], preview: 'Assalamualaikum Bapak/Ibu [Nama], perkenalkan saya [Nama Agen] dari [Travel]. Ada kabar baik untuk Anda...' },
    { id: 'sc2', title: 'Penanganan Keberatan Harga', category: 'Sales', tags: ['keberatan', 'harga'], preview: 'Saya mengerti kekhawatiran Bapak/Ibu soal harga. Mari saya jelaskan apa saja yang sudah termasuk dalam paket ini...' },
    { id: 'sc3', title: 'Follow-up 3 Hari Tidak Balas', category: 'WhatsApp', tags: ['follow-up', 'stale'], preview: 'Selamat siang Bapak/Ibu [Nama], saya ingin memastikan apakah informasi paket Umroh yang saya kirimkan sudah diterima...' },
    { id: 'sc4', title: 'Closing — Ajakan Booking Hari Ini', category: 'Sales', tags: ['closing', 'booking'], preview: 'Keberangkatan bulan [Bulan] tinggal [N] seat lagi. Untuk mengamankan tempat Bapak/Ibu, kita perlu proses DP hari ini...' },
    { id: 'sc5', title: 'Permintaan Testimoni Jamaah', category: 'Post-Service', tags: ['testimoni', 'alumni'], preview: 'Alhamdulillah, perjalanan Umroh Bapak/Ibu telah selesai. Boleh kami minta cerita pengalaman singkat untuk membantu jamaah lain?' }
  ]);

  const filteredScripts = $derived(
    !scriptSearch.trim()
      ? scripts
      : scripts.filter(s =>
          s.title.toLowerCase().includes(scriptSearch.toLowerCase()) ||
          s.tags.some(t => t.includes(scriptSearch.toLowerCase()))
        )
  );

  let expandedScript = $state<string | null>(null);

  // --- Badges ---
  const badges = $state([
    { id: 'bg1', name: 'Agen Terverifikasi', icon: 'verified', color: '#2563eb', bg: '#dbeafe', earned: true, date: '2026-01-15' },
    { id: 'bg2', name: 'Sales Champion', icon: 'emoji_events', color: '#d97706', bg: '#fef3c7', earned: true, date: '2026-02-10' },
    { id: 'bg3', name: 'Digital Marketer', icon: 'campaign', color: '#7c3aed', bg: '#ede9fe', earned: true, date: '2026-03-05' },
    { id: 'bg4', name: 'Mentor Downline', icon: 'group', color: '#059669', bg: '#d1fae5', earned: false, date: '' },
    { id: 'bg5', name: 'Top Recruiter', icon: 'person_add', color: '#dc2626', bg: '#fee2e2', earned: false, date: '' },
    { id: 'bg6', name: 'Fiqih Umroh Expert', icon: 'menu_book', color: '#0891b2', bg: '#cffafe', earned: true, date: '2026-03-20' }
  ]);

  function fmtDate(iso: string): string {
    if (!iso) return '';
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">hub</span>
      <span class="sep">/</span>
      <a href="/console/crm" class="bc-link">CRM</a>
      <span class="sep">/</span>
      <a href="/console/crm/agency" class="bc-link">Portal Agen</a>
      <span class="sep">/</span>
      <span class="topbar-current">Akademi Agen</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Akademi Agen</h2>
        <p>Tingkatkan kemampuan agen melalui kursus, leaderboard, skrip jual, dan sertifikasi</p>
      </div>
    </div>

    <div class="section-tabs">
      {#each SECTIONS as sec}
        <button
          class="sec-tab"
          class:active={activeSection === sec.id}
          onclick={() => { activeSection = sec.id; }}
        >
          <span class="material-symbols-outlined">{sec.icon}</span>
          {sec.label}
        </button>
      {/each}
    </div>

    <!-- Courses -->
    {#if activeSection === 'courses'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title"><span class="material-symbols-outlined">school</span>Kursus & Modul</div>
          <div class="cat-filters">
            {#each CATEGORIES as cat}
              <button class="cat-btn" class:active={courseFilter === cat} onclick={() => { courseFilter = cat; }}>
                {cat}
              </button>
            {/each}
          </div>
        </div>
        <div class="course-grid">
          {#each filteredCourses as course (course.id)}
            <div class="course-card">
              <div class="course-head">
                <span class="course-cat">{course.category}</span>
                {#if course.badge}
                  <span class="badge-chip">
                    <span class="material-symbols-outlined">workspace_premium</span>
                    Sertifikat
                  </span>
                {/if}
              </div>
              <div class="course-title">{course.title}</div>
              <div class="course-meta">
                <span><span class="material-symbols-outlined icon-xs">schedule</span>{course.duration}</span>
                <span><span class="material-symbols-outlined icon-xs">menu_book</span>{course.lessons} pelajaran</span>
                <span><span class="material-symbols-outlined icon-xs">group</span>{course.enrolled} terdaftar</span>
              </div>
              <div class="progress-wrap">
                <div class="progress-bar">
                  <div class="progress-fill" style="width:{course.progress}%"></div>
                </div>
                <span class="progress-label">{course.progress}%</span>
              </div>
              <button class="course-btn">
                {course.progress === 0 ? 'Mulai Kursus' : course.progress === 100 ? 'Lihat Ulang' : 'Lanjutkan'}
              </button>
            </div>
          {/each}
        </div>
      </div>

    <!-- Leaderboard -->
    {:else if activeSection === 'leaderboard'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">leaderboard</span>Leaderboard Akademi</div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Peringkat</th>
                <th>Agen</th>
                <th>Skor</th>
                <th>Kursus Selesai</th>
                <th>Tier</th>
                <th>Tren</th>
              </tr>
            </thead>
            <tbody>
              {#each leaderboard as row}
                <tr class:top-three={row.rank <= 3}>
                  <td>
                    {#if row.rank === 1}
                      <span class="rank-medal gold">1</span>
                    {:else if row.rank === 2}
                      <span class="rank-medal silver">2</span>
                    {:else if row.rank === 3}
                      <span class="rank-medal bronze">3</span>
                    {:else}
                      <span class="rank-num">{row.rank}</span>
                    {/if}
                  </td>
                  <td class="font-semibold">{row.agent}</td>
                  <td><span class="score-num">{row.score.toLocaleString('id-ID')}</span></td>
                  <td>{row.courses_done} kursus</td>
                  <td><span class="tier-badge tier-{row.tier.toLowerCase()}">{row.tier}</span></td>
                  <td>
                    {#if row.trend === 'up'}
                      <span class="trend-up"><span class="material-symbols-outlined">arrow_upward</span></span>
                    {:else if row.trend === 'down'}
                      <span class="trend-down"><span class="material-symbols-outlined">arrow_downward</span></span>
                    {:else}
                      <span class="trend-stable">—</span>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>

    <!-- Scripts -->
    {:else if activeSection === 'scripts'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title"><span class="material-symbols-outlined">chat</span>Skrip Penjualan</div>
          <div class="search-wrap">
            <span class="material-symbols-outlined search-icon">search</span>
            <input type="text" placeholder="Cari skrip..." bind:value={scriptSearch} />
          </div>
        </div>
        <div class="script-list">
          {#each filteredScripts as script (script.id)}
            <div class="script-card">
              <div class="script-card-head" role="button" tabindex="0"
                onclick={() => { expandedScript = expandedScript === script.id ? null : script.id; }}
                onkeydown={(e) => { if (e.key === 'Enter') { expandedScript = expandedScript === script.id ? null : script.id; } }}
              >
                <div class="script-info">
                  <span class="script-title">{script.title}</span>
                  <span class="script-cat">{script.category}</span>
                </div>
                <div class="script-tags">
                  {#each script.tags as tag}
                    <span class="tag">#{tag}</span>
                  {/each}
                </div>
                <span class="material-symbols-outlined expand-icon">
                  {expandedScript === script.id ? 'expand_less' : 'expand_more'}
                </span>
              </div>
              {#if expandedScript === script.id}
                <div class="script-preview">{script.preview}</div>
                <div class="script-actions">
                  <button class="sm-action-btn">
                    <span class="material-symbols-outlined">content_copy</span>
                    Salin
                  </button>
                </div>
              {/if}
            </div>
          {/each}
        </div>
      </div>

    <!-- Badges -->
    {:else if activeSection === 'badges'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">workspace_premium</span>Badge & Sertifikasi</div>
        <div class="badge-grid">
          {#each badges as badge (badge.id)}
            <div class="badge-card" class:earned={badge.earned} class:locked={!badge.earned}>
              <div class="badge-icon" style="background:{badge.bg};color:{badge.color}">
                <span class="material-symbols-outlined">{badge.icon}</span>
              </div>
              <div class="badge-name">{badge.name}</div>
              {#if badge.earned}
                <div class="badge-date">Diperoleh {fmtDate(badge.date)}</div>
                <span class="badge-earned-chip">
                  <span class="material-symbols-outlined">check_circle</span>
                  Diperoleh
                </span>
              {:else}
                <div class="badge-locked-label">Belum Diperoleh</div>
                <span class="badge-locked-chip">
                  <span class="material-symbols-outlined">lock</span>
                  Terkunci
                </span>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; justify-content: space-between; gap: 1rem; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.88rem; color: #434655; }
  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .sep { color: #b0b3c1; }
  .bc-link { color: #004ac6; text-decoration: none; font-weight: 500; }
  .bc-link:hover { text-decoration: underline; }
  .topbar-current { font-weight: 600; color: #191c1e; }
  .top-actions { display: flex; align-items: center; gap: 0.35rem; }
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .section-tabs { display: flex; gap: 0.25rem; flex-wrap: wrap; margin-bottom: 1.25rem; }
  .sec-tab { display: inline-flex; align-items: center; gap: 0.35rem; padding: 0.45rem 0.85rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; background: #fff; font-size: 0.78rem; color: #434655; cursor: pointer; }
  .sec-tab .material-symbols-outlined { font-size: 1rem; }
  .sec-tab:hover { background: #f2f4f6; }
  .sec-tab.active { border-color: #2563eb; color: #004ac6; background: #eff6ff; font-weight: 700; }

  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; overflow: hidden; margin-bottom: 1rem; }
  .section-header { display: flex; align-items: center; justify-content: space-between; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); flex-wrap: wrap; gap: 0.5rem; }
  .section-title { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; font-weight: 700; color: #191c1e; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .section-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }
  .section-header .section-title { padding: 0; border: 0; }

  .cat-filters { display: flex; gap: 0.25rem; flex-wrap: wrap; }
  .cat-btn { padding: 0.3rem 0.65rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.72rem; color: #434655; cursor: pointer; }
  .cat-btn:hover { background: #f2f4f6; }
  .cat-btn.active { border-color: #2563eb; color: #004ac6; background: #eff6ff; font-weight: 700; }

  /* Courses */
  .course-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 0.85rem; padding: 1rem; }
  .course-card { border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; padding: 1rem; display: flex; flex-direction: column; gap: 0.55rem; }
  .course-head { display: flex; align-items: center; justify-content: space-between; }
  .course-cat { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #2563eb; }
  .badge-chip { display: inline-flex; align-items: center; gap: 0.2rem; font-size: 0.6rem; font-weight: 700; color: #92400e; background: #fef3c7; padding: 0.1rem 0.35rem; border-radius: 0.2rem; }
  .badge-chip .material-symbols-outlined { font-size: 0.8rem; }
  .course-title { font-weight: 700; font-size: 0.85rem; color: #191c1e; line-height: 1.3; }
  .course-meta { display: flex; gap: 0.75rem; flex-wrap: wrap; }
  .course-meta span { display: inline-flex; align-items: center; gap: 0.2rem; font-size: 0.68rem; color: #737686; }
  .icon-xs { font-size: 0.85rem; }
  .progress-wrap { display: flex; align-items: center; gap: 0.5rem; }
  .progress-bar { flex: 1; height: 6px; background: #e2e8f0; border-radius: 999px; }
  .progress-fill { height: 100%; background: #2563eb; border-radius: 999px; transition: width 0.3s; }
  .progress-label { font-size: 0.62rem; color: #437655; font-weight: 700; min-width: 2.5rem; text-align: right; }
  .course-btn { padding: 0.45rem 0.75rem; border: 1px solid #2563eb; border-radius: 0.25rem; background: #fff; color: #2563eb; font-size: 0.75rem; font-weight: 600; cursor: pointer; margin-top: auto; }
  .course-btn:hover { background: #eff6ff; }

  /* Leaderboard */
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .top-three { background: #fffbeb; }
  .font-semibold { font-weight: 600; color: #191c1e; }

  .rank-medal { display: inline-flex; width: 1.8rem; height: 1.8rem; border-radius: 999px; align-items: center; justify-content: center; font-size: 0.75rem; font-weight: 700; }
  .rank-medal.gold { background: #fef3c7; color: #92400e; }
  .rank-medal.silver { background: #f1f5f9; color: #334155; }
  .rank-medal.bronze { background: #fde8d8; color: #7c2d12; }
  .rank-num { font-weight: 700; color: #434655; }

  .score-num { font-weight: 700; color: #2563eb; }

  .tier-badge { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .tier-gold { background: #fef3c7; color: #92400e; }
  .tier-silver { background: #f1f5f9; color: #334155; }
  .tier-bronze { background: #fde8d8; color: #7c2d12; }

  .trend-up .material-symbols-outlined { font-size: 1rem; color: #059669; }
  .trend-down .material-symbols-outlined { font-size: 1rem; color: #dc2626; }
  .trend-stable { color: #b0b3c1; }

  /* Scripts */
  .search-wrap { position: relative; }
  .search-icon { position: absolute; left: 0.65rem; top: 50%; transform: translateY(-50%); font-size: 1rem; color: #737686; }
  .search-wrap input { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.4rem 0.7rem 0.4rem 2.1rem; font-size: 0.82rem; color: #191c1e; min-width: 14rem; }

  .script-list { display: flex; flex-direction: column; }
  .script-card { border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .script-card:last-child { border-bottom: 0; }
  .script-card-head { display: flex; align-items: center; gap: 0.85rem; padding: 0.85rem 1rem; cursor: pointer; }
  .script-card-head:hover { background: #f7f9fb; }
  .script-info { flex: 1; }
  .script-title { font-weight: 700; font-size: 0.82rem; color: #191c1e; display: block; }
  .script-cat { font-size: 0.62rem; color: #2563eb; margin-top: 0.1rem; display: block; }
  .script-tags { display: flex; gap: 0.25rem; flex-wrap: wrap; }
  .tag { font-size: 0.6rem; color: #737686; background: #f1f5f9; padding: 0.1rem 0.3rem; border-radius: 0.15rem; }
  .expand-icon { font-size: 1.2rem; color: #b0b3c1; flex-shrink: 0; }
  .script-preview { padding: 0.75rem 1rem; font-size: 0.78rem; color: #434655; background: #f7f9fb; border-top: 1px solid rgb(195 198 215 / 0.35); line-height: 1.5; }
  .script-actions { padding: 0.5rem 1rem; border-top: 1px solid rgb(195 198 215 / 0.35); }
  .sm-action-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.3rem 0.6rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.72rem; font-weight: 600; color: #191c1e; cursor: pointer; }
  .sm-action-btn:hover { background: #f2f4f6; }
  .sm-action-btn .material-symbols-outlined { font-size: 0.9rem; }

  /* Badges */
  .badge-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 0.85rem; padding: 1rem; }
  .badge-card { border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; padding: 1rem; display: flex; flex-direction: column; align-items: center; gap: 0.4rem; text-align: center; }
  .badge-card.locked { opacity: 0.5; }
  .badge-icon { width: 3rem; height: 3rem; border-radius: 999px; display: grid; place-items: center; }
  .badge-icon .material-symbols-outlined { font-size: 1.5rem; }
  .badge-name { font-weight: 700; font-size: 0.8rem; color: #191c1e; }
  .badge-date { font-size: 0.62rem; color: #737686; }
  .badge-locked-label { font-size: 0.62rem; color: #b0b3c1; }
  .badge-earned-chip { display: inline-flex; align-items: center; gap: 0.2rem; font-size: 0.6rem; background: #d1fae5; color: #065f46; padding: 0.1rem 0.35rem; border-radius: 999px; font-weight: 700; }
  .badge-earned-chip .material-symbols-outlined { font-size: 0.75rem; }
  .badge-locked-chip { display: inline-flex; align-items: center; gap: 0.2rem; font-size: 0.6rem; background: #f1f5f9; color: #64748b; padding: 0.1rem 0.35rem; border-radius: 999px; font-weight: 700; }
  .badge-locked-chip .material-symbols-outlined { font-size: 0.75rem; }
</style>
