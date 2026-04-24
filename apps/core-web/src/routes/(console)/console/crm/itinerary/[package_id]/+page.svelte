<script lang="ts">
  import type { PageData } from './+page.server';

  let { data }: { data: PageData } = $props();

  const itin = $derived(data.itinerary);

  // Accordion state: set of open day indices
  let openDays = $state<Set<number>>(new Set([0]));

  function toggleDay(idx: number) {
    const next = new Set(openDays);
    if (next.has(idx)) {
      next.delete(idx);
    } else {
      next.add(idx);
    }
    openDays = next;
  }

  function expandAll() {
    if (!itin) return;
    openDays = new Set(itin.days.map((_, i) => i));
  }

  function collapseAll() {
    openDays = new Set();
  }

  const ACTIVITY_ICON_MAP: Record<string, string> = {
    flight: 'flight',
    hotel: 'hotel',
    ibadah: 'mosque',
    meal: 'restaurant',
    tour: 'camera_outdoor',
    ziarah: 'explore',
    transport: 'directions_bus'
  };

  const ACTIVITY_COLOR_MAP: Record<string, string> = {
    flight: '#0369a1',
    hotel: '#4c1d95',
    ibadah: '#065f46',
    meal: '#b45309',
    tour: '#7c3aed',
    ziarah: '#9f1239',
    transport: '#374151'
  };

  const ACTIVITY_BG_MAP: Record<string, string> = {
    flight: '#e0f2fe',
    hotel: '#ede9fe',
    ibadah: '#d1fae5',
    meal: '#fef3c7',
    tour: '#f3e8ff',
    ziarah: '#ffe4e6',
    transport: '#f3f4f6'
  };

  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
  }

  function formatIDR(n: number): string {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let copied = $state(false);

  function copyLink() {
    navigator.clipboard.writeText(data.shareUrl).then(() => {
      copied = true;
      setTimeout(() => { copied = false; }, 2000);
    }).catch(() => {
      prompt('Copy link berikut:', data.shareUrl);
    });
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">hub</span>
      <span class="sep">/</span>
      <a href="/console/crm" class="crumb-link">CRM Tools</a>
      <span class="sep">/</span>
      <span class="topbar-current">Itinerary</span>
    </nav>
    <div class="top-actions">
      <button class="share-btn" onclick={copyLink}>
        <span class="material-symbols-outlined">{copied ? 'check' : 'link'}</span>
        {copied ? 'Tersalin!' : 'Salin Link'}
      </button>
      <button class="icon-btn" title="Notifikasi">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar" aria-label="Profile">AD</button>
    </div>
  </header>

  <section class="canvas-page">
    {#if !itin}
      <div class="empty-state">
        <span class="material-symbols-outlined">route</span>
        <p>Itinerary tidak ditemukan.</p>
        <a href="/console/crm" class="back-link">Kembali ke CRM Tools</a>
      </div>
    {:else}
      <!-- Package Hero -->
      <div class="itin-hero">
        <div class="hero-main">
          <div class="hero-badges">
            {#each itin.destinations as dest}
              <span class="dest-badge">
                <span class="material-symbols-outlined">location_on</span>
                {dest}
              </span>
            {/each}
          </div>
          <h1 class="hero-title">{itin.name}</h1>
          <div class="hero-meta">
            <span class="meta-item">
              <span class="material-symbols-outlined">calendar_today</span>
              {formatDate(itin.departure_date)} — {formatDate(itin.return_date)}
            </span>
            <span class="meta-item">
              <span class="material-symbols-outlined">schedule</span>
              {itin.duration_days} hari
            </span>
            <span class="meta-item">
              <span class="material-symbols-outlined">flight</span>
              {itin.airline}
            </span>
            <span class="meta-item price-meta">
              <span class="material-symbols-outlined">payments</span>
              {formatIDR(itin.price)} / pax
            </span>
          </div>
        </div>

        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-value">{itin.duration_days}</span>
            <span class="stat-label">Hari</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-value">{itin.quota}</span>
            <span class="stat-label">Kuota</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-value">{itin.destinations.length}</span>
            <span class="stat-label">Kota</span>
          </div>
        </div>
      </div>

      <!-- Info cards -->
      <div class="info-grid">
        <div class="info-card">
          <span class="material-symbols-outlined info-icon">person</span>
          <div>
            <div class="info-label">Pembimbing</div>
            <div class="info-value">{itin.guide_name}</div>
          </div>
        </div>
        <div class="info-card">
          <span class="material-symbols-outlined info-icon">hotel</span>
          <div>
            <div class="info-label">Hotel Makkah</div>
            <div class="info-value">{itin.hotel_makkah}</div>
          </div>
        </div>
        <div class="info-card">
          <span class="material-symbols-outlined info-icon">hotel</span>
          <div>
            <div class="info-label">Hotel Madinah</div>
            <div class="info-value">{itin.hotel_madinah}</div>
          </div>
        </div>
      </div>

      <!-- Highlights -->
      <div class="highlights-section">
        <h3 class="section-title">
          <span class="material-symbols-outlined">star</span>
          Fasilitas & Highlights
        </h3>
        <div class="highlights-grid">
          {#each itin.highlights as h}
            <div class="highlight-item">
              <span class="material-symbols-outlined highlight-check">check_circle</span>
              {h}
            </div>
          {/each}
        </div>
      </div>

      <!-- Day-by-day itinerary -->
      <div class="days-section">
        <div class="days-header">
          <h3 class="section-title">
            <span class="material-symbols-outlined">route</span>
            Itinerary Hari per Hari
          </h3>
          <div class="days-controls">
            <button class="ctrl-btn" onclick={expandAll}>Buka Semua</button>
            <button class="ctrl-btn" onclick={collapseAll}>Tutup Semua</button>
          </div>
        </div>

        <div class="days-list">
          {#each itin.days as day, i (day.day)}
            <div class="day-item">
              <!-- Day header (accordion trigger) -->
              <button
                type="button"
                class="day-header"
                class:open={openDays.has(i)}
                onclick={() => toggleDay(i)}
                aria-expanded={openDays.has(i)}
              >
                <div class="day-num">
                  <span class="day-n">Hari {day.day}</span>
                </div>
                <div class="day-title-wrap">
                  <span class="day-title">{day.title}</span>
                  <span class="day-date">{formatDate(day.date)}</span>
                </div>
                {#if day.hotel}
                  <span class="day-hotel">
                    <span class="material-symbols-outlined">hotel</span>
                    {day.hotel.name}
                  </span>
                {/if}
                <span class="material-symbols-outlined chevron">
                  {openDays.has(i) ? 'keyboard_arrow_up' : 'keyboard_arrow_down'}
                </span>
              </button>

              <!-- Day content -->
              {#if openDays.has(i)}
                <div class="day-content">
                  <div class="timeline">
                    {#each day.activities as act, j}
                      <div class="timeline-item">
                        <div class="tl-time">{act.time}</div>
                        <div class="tl-dot-wrap">
                          <div
                            class="tl-dot"
                            style="background: {ACTIVITY_BG_MAP[act.type] ?? '#f3f4f6'}; color: {ACTIVITY_COLOR_MAP[act.type] ?? '#374151'}"
                          >
                            <span class="material-symbols-outlined">
                              {ACTIVITY_ICON_MAP[act.type] ?? act.icon}
                            </span>
                          </div>
                          {#if j < day.activities.length - 1}
                            <div class="tl-line"></div>
                          {/if}
                        </div>
                        <div class="tl-body">
                          <div class="tl-title">{act.title}</div>
                          {#if act.desc}
                            <div class="tl-desc">{act.desc}</div>
                          {/if}
                          <span
                            class="tl-type-badge"
                            style="background: {ACTIVITY_BG_MAP[act.type] ?? '#f3f4f6'}; color: {ACTIVITY_COLOR_MAP[act.type] ?? '#374151'}"
                          >
                            {act.type}
                          </span>
                        </div>
                      </div>
                    {/each}
                  </div>

                  {#if day.hotel}
                    <div class="hotel-card">
                      <span class="material-symbols-outlined">hotel</span>
                      <div class="hotel-info">
                        <div class="hotel-name">{day.hotel.name}</div>
                        <div class="hotel-meta">
                          {'★'.repeat(day.hotel.stars)} · {day.hotel.location}
                        </div>
                      </div>
                    </div>
                  {/if}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      </div>

      <!-- Share bar -->
      <div class="share-bar">
        <span class="material-symbols-outlined">share</span>
        <span>Bagikan itinerary ini ke calon jamaah:</span>
        <code class="share-url">{data.shareUrl}</code>
        <button class="share-copy-btn" onclick={copyLink}>
          <span class="material-symbols-outlined">{copied ? 'check' : 'content_copy'}</span>
          {copied ? 'Tersalin!' : 'Copy Link'}
        </button>
      </div>
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }

  .topbar {
    position: sticky; top: 0; z-index: 30;
    height: 4rem;
    background: rgb(255 255 255 / 0.9);
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem;
    display: flex; align-items: center; justify-content: space-between; gap: 1rem;
    backdrop-filter: blur(8px);
  }

  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.88rem; color: #434655; }
  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .sep { color: #b0b3c1; }
  .crumb-link { color: #434655; text-decoration: none; }
  .crumb-link:hover { color: #004ac6; }
  .topbar-current { font-weight: 600; color: #191c1e; }

  .top-actions { display: flex; align-items: center; gap: 0.5rem; }

  .share-btn {
    display: inline-flex; align-items: center; gap: 0.3rem;
    padding: 0.42rem 0.75rem;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff; border: 0; border-radius: 0.25rem;
    font-size: 0.78rem; font-weight: 600; cursor: pointer; font-family: inherit;
  }
  .share-btn:hover { opacity: 0.9; }
  .share-btn .material-symbols-outlined { font-size: 0.9rem; }

  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas-page { padding: 1.5rem; max-width: 72rem; }

  /* Hero */
  .itin-hero {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1.5rem;
    background: linear-gradient(135deg, #1e40af, #0369a1);
    border-radius: 0.5rem;
    padding: 1.75rem;
    margin-bottom: 1.25rem;
    color: #fff;
  }

  .hero-main { flex: 1; }

  .hero-badges { display: flex; gap: 0.4rem; flex-wrap: wrap; margin-bottom: 0.75rem; }

  .dest-badge {
    display: inline-flex; align-items: center; gap: 0.25rem;
    padding: 0.2rem 0.55rem;
    background: rgba(255,255,255,0.18);
    border-radius: 999px;
    font-size: 0.72rem;
    font-weight: 600;
  }
  .dest-badge .material-symbols-outlined { font-size: 0.75rem; }

  .hero-title { margin: 0 0 0.85rem; font-size: 1.6rem; line-height: 1.2; }

  .hero-meta { display: flex; flex-wrap: wrap; gap: 0.85rem; }

  .meta-item {
    display: inline-flex; align-items: center; gap: 0.3rem;
    font-size: 0.8rem; color: rgba(255,255,255,0.85);
  }
  .meta-item .material-symbols-outlined { font-size: 0.9rem; }
  .price-meta { font-weight: 700; color: #fde68a; }

  .hero-stats {
    display: flex;
    align-items: center;
    gap: 1rem;
    background: rgba(255,255,255,0.12);
    border-radius: 0.35rem;
    padding: 1rem 1.25rem;
    flex-shrink: 0;
  }

  .stat-item { display: flex; flex-direction: column; align-items: center; gap: 0.1rem; }
  .stat-value { font-size: 1.6rem; font-weight: 800; line-height: 1; }
  .stat-label { font-size: 0.65rem; color: rgba(255,255,255,0.7); }
  .stat-divider { width: 1px; height: 2rem; background: rgba(255,255,255,0.25); }

  /* Info grid */
  .info-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.85rem;
    margin-bottom: 1.25rem;
  }

  @media (max-width: 700px) { .info-grid { grid-template-columns: 1fr; } }

  .info-card {
    display: flex; align-items: flex-start; gap: 0.65rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 0.85rem;
  }

  .info-icon { font-size: 1.1rem; color: #2563eb; margin-top: 1px; }
  .info-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #737686; margin-bottom: 0.2rem; }
  .info-value { font-size: 0.82rem; font-weight: 600; color: #191c1e; }

  /* Highlights */
  .highlights-section {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1.1rem;
    margin-bottom: 1.25rem;
  }

  .section-title {
    display: flex; align-items: center; gap: 0.5rem;
    margin: 0 0 0.85rem;
    font-size: 0.95rem; font-weight: 700; color: #191c1e;
  }
  .section-title .material-symbols-outlined { font-size: 1.05rem; color: #2563eb; }

  .highlights-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 0.5rem;
  }

  .highlight-item {
    display: flex; align-items: center; gap: 0.4rem;
    font-size: 0.82rem; color: #191c1e;
  }
  .highlight-check { font-size: 1rem; color: #065f46; }

  /* Days */
  .days-section {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    overflow: hidden;
    margin-bottom: 1.25rem;
  }

  .days-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 1rem 1.1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
  }

  .days-header .section-title { margin: 0; }

  .days-controls { display: flex; gap: 0.4rem; }

  .ctrl-btn {
    padding: 0.3rem 0.6rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.2rem;
    background: #fff;
    font-size: 0.72rem;
    font-weight: 600;
    color: #434655;
    cursor: pointer;
    font-family: inherit;
  }
  .ctrl-btn:hover { background: #f2f4f6; }

  .days-list { display: flex; flex-direction: column; }

  .day-item { border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .day-item:last-child { border-bottom: 0; }

  .day-header {
    display: flex;
    align-items: center;
    gap: 0.85rem;
    width: 100%;
    padding: 0.9rem 1.1rem;
    border: 0;
    background: transparent;
    cursor: pointer;
    text-align: left;
    font-family: inherit;
    transition: background 0.1s;
  }

  .day-header:hover { background: #f7f9fb; }
  .day-header.open { background: #f0f9ff; }

  .day-num {
    display: flex; align-items: center; justify-content: center;
    width: 2.2rem; height: 2.2rem;
    border-radius: 0.25rem;
    background: #1e40af;
    flex-shrink: 0;
  }
  .day-n { font-size: 0.62rem; font-weight: 700; color: #fff; line-height: 1.2; text-align: center; }

  .day-title-wrap { flex: 1; min-width: 0; }
  .day-title { display: block; font-size: 0.88rem; font-weight: 700; color: #191c1e; }
  .day-date { display: block; font-size: 0.72rem; color: #737686; margin-top: 0.1rem; }

  .day-hotel {
    display: inline-flex; align-items: center; gap: 0.25rem;
    font-size: 0.72rem; color: #4c1d95;
    background: #ede9fe; padding: 0.15rem 0.45rem; border-radius: 0.2rem;
    white-space: nowrap;
    flex-shrink: 0;
  }
  .day-hotel .material-symbols-outlined { font-size: 0.78rem; }

  .chevron { font-size: 1.2rem; color: #737686; flex-shrink: 0; }

  /* Day content */
  .day-content { padding: 0.75rem 1.1rem 1.1rem 1.1rem; }

  .timeline { display: flex; flex-direction: column; gap: 0; margin-bottom: 0.75rem; }

  .timeline-item { display: flex; gap: 0.75rem; align-items: flex-start; }

  .tl-time {
    width: 3.2rem;
    font-size: 0.72rem;
    font-weight: 700;
    color: #434655;
    padding-top: 0.55rem;
    flex-shrink: 0;
    font-variant-numeric: tabular-nums;
  }

  .tl-dot-wrap {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex-shrink: 0;
  }

  .tl-dot {
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }
  .tl-dot .material-symbols-outlined { font-size: 0.9rem; }

  .tl-line { width: 2px; flex: 1; min-height: 0.85rem; background: #e5e7eb; margin: 0.15rem auto; }

  .tl-body { flex: 1; padding-top: 0.3rem; padding-bottom: 0.75rem; }

  .tl-title { font-size: 0.85rem; font-weight: 700; color: #191c1e; margin-bottom: 0.2rem; }

  .tl-desc { font-size: 0.78rem; color: #434655; margin-bottom: 0.3rem; line-height: 1.4; }

  .tl-type-badge {
    display: inline-flex;
    padding: 0.1rem 0.35rem;
    border-radius: 0.15rem;
    font-size: 0.6rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .hotel-card {
    display: flex; align-items: center; gap: 0.6rem;
    padding: 0.65rem 0.85rem;
    background: #f5f3ff;
    border: 1px solid #ddd6fe;
    border-radius: 0.25rem;
    margin-top: 0.25rem;
  }

  .hotel-card .material-symbols-outlined { font-size: 1.1rem; color: #7c3aed; flex-shrink: 0; }
  .hotel-name { font-size: 0.82rem; font-weight: 700; color: #4c1d95; }
  .hotel-meta { font-size: 0.72rem; color: #6d28d9; }

  /* Share bar */
  .share-bar {
    display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap;
    padding: 0.85rem 1rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    font-size: 0.82rem; color: #434655;
    margin-bottom: 1.5rem;
  }

  .share-bar .material-symbols-outlined { font-size: 1rem; color: #2563eb; }

  .share-url {
    flex: 1;
    font-size: 0.75rem;
    background: #f2f4f6;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.2rem;
    padding: 0.3rem 0.55rem;
    font-family: 'IBM Plex Mono', monospace;
    color: #374151;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
  }

  .share-copy-btn {
    display: inline-flex; align-items: center; gap: 0.3rem;
    padding: 0.4rem 0.75rem;
    background: #2563eb; color: #fff; border: 0;
    border-radius: 0.25rem; font-size: 0.75rem; font-weight: 600;
    cursor: pointer; font-family: inherit; white-space: nowrap;
  }
  .share-copy-btn:hover { opacity: 0.9; }
  .share-copy-btn .material-symbols-outlined { font-size: 0.85rem; }

  /* Empty */
  .empty-state {
    display: flex; flex-direction: column; align-items: center; gap: 0.75rem;
    padding: 5rem 1rem; color: #b0b3c1;
  }
  .empty-state .material-symbols-outlined { font-size: 3rem; }
  .empty-state p { margin: 0; font-size: 0.9rem; }
  .back-link { color: #004ac6; font-weight: 600; text-decoration: none; font-size: 0.82rem; }
</style>
