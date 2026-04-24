<script lang="ts">
  type ContentSection = 'bank' | 'flyer' | 'utm' | 'ads' | 'calendar' | 'analytics';

  const SECTIONS: Array<{ id: ContentSection; label: string; icon: string }> = [
    { id: 'bank', label: 'Bank Konten', icon: 'perm_media' },
    { id: 'flyer', label: 'Watermark Flyer', icon: 'water_drop' },
    { id: 'utm', label: 'UTM Link Builder', icon: 'link' },
    { id: 'ads', label: 'Ads Manager Lite', icon: 'campaign' },
    { id: 'calendar', label: 'Kalender Konten', icon: 'calendar_month' },
    { id: 'analytics', label: 'Analitik Konten', icon: 'analytics' }
  ];

  let activeSection = $state<ContentSection>('bank');

  // --- Bank Konten ---
  type Category = 'semua' | 'foto' | 'video' | 'infografis' | 'template';
  let categoryFilter = $state<Category>('semua');

  const contentItems = $state([
    { id: 'c1', title: 'Foto Masjidil Haram Terbaru', category: 'foto', tags: ['makkah', 'haram'], thumb: '#dbeafe' },
    { id: 'c2', title: 'Video Profil Paket Gold 2026', category: 'video', tags: ['gold', 'paket'], thumb: '#fde8d8' },
    { id: 'c3', title: 'Infografis Syarat Umroh', category: 'infografis', tags: ['edukasi', 'syarat'], thumb: '#d1fae5' },
    { id: 'c4', title: 'Template Flyer Ramadan', category: 'template', tags: ['ramadan', 'flyer'], thumb: '#ede9fe' },
    { id: 'c5', title: 'Foto Raudhah Madinah', category: 'foto', tags: ['madinah', 'raudhah'], thumb: '#fef9c3' },
    { id: 'c6', title: 'Template Testimoni Jamaah', category: 'template', tags: ['testimoni'], thumb: '#fce7f3' },
    { id: 'c7', title: 'Infografis Jadwal Keberangkatan', category: 'infografis', tags: ['jadwal'], thumb: '#dbeafe' },
    { id: 'c8', title: 'Video Review Hotel Bintang 5', category: 'video', tags: ['hotel', 'mewah'], thumb: '#fde8d8' }
  ]);

  const filteredContent = $derived(
    categoryFilter === 'semua' ? contentItems : contentItems.filter(c => c.category === categoryFilter)
  );

  // --- UTM Link Builder ---
  let utmBase = $state('https://umrohos.id/paket/gold-2026');
  let utmSource = $state('whatsapp');
  let utmMedium = $state('sosial');
  let utmCampaign = $state('ramadan2026');
  let utmAgent = $state('mitra-barokah');

  const generatedUtm = $derived(
    `${utmBase}?utm_source=${utmSource}&utm_medium=${utmMedium}&utm_campaign=${utmCampaign}&utm_content=${utmAgent}`
  );

  const savedLinks = $state([
    { id: 'u1', agent: 'PT Mitra Barokah', campaign: 'ramadan2026', source: 'instagram', clicks: 1_240, leads: 38 },
    { id: 'u2', agent: 'CV Cahaya Umroh', campaign: 'harbolnas', source: 'facebook', clicks: 530, leads: 14 },
    { id: 'u3', agent: 'Karima Travel', campaign: 'mei2026', source: 'whatsapp', clicks: 870, leads: 29 }
  ]);

  // --- Ads Manager Lite ---
  const campaigns = $state([
    { id: 'ad1', name: 'Ramadan 2026 — Gold Package', platform: 'Meta', budget: 5_000_000, spend: 3_850_000, impressions: 48_200, leads: 62, status: 'active' },
    { id: 'ad2', name: 'Harbolnas Umroh Hemat', platform: 'Google', budget: 3_000_000, spend: 3_000_000, impressions: 32_100, leads: 41, status: 'finished' },
    { id: 'ad3', name: 'Brand Awareness Q2', platform: 'TikTok', budget: 2_500_000, spend: 1_200_000, impressions: 85_000, leads: 18, status: 'active' }
  ]);

  // --- Kalender Konten (mini weekly view) ---
  const weekDays = ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min'];
  const calendarItems = $state([
    { day: 0, title: 'Post IG: Testimoni', type: 'post' },
    { day: 1, title: 'WA Blast: Promo Gold', type: 'blast' },
    { day: 2, title: 'FB: Infografis Syarat', type: 'post' },
    { day: 4, title: 'Story: Video Hotel', type: 'story' },
    { day: 5, title: 'TikTok: Behind the Scene', type: 'video' },
    { day: 6, title: 'WA Blast: Reminder', type: 'blast' }
  ]);

  // --- Analytics ---
  const analyticsData = $state([
    { channel: 'Instagram', reach: 42_100, engagement: 3_840, ctr: 4.2, leads: 58 },
    { channel: 'WhatsApp', reach: 8_700, engagement: 2_100, ctr: 12.1, leads: 44 },
    { channel: 'Facebook', reach: 21_500, engagement: 1_280, ctr: 2.8, leads: 22 },
    { channel: 'TikTok', reach: 85_400, engagement: 9_200, ctr: 1.9, leads: 18 },
    { channel: 'YouTube', reach: 4_100, engagement: 620, ctr: 3.7, leads: 9 }
  ]);

  function fmtRp(v: number): string {
    return 'Rp ' + (v / 1_000_000).toFixed(1) + ' jt';
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
      <span class="topbar-current">Konten & Pemasaran</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Konten & Pemasaran</h2>
        <p>Bank konten, watermark flyer, UTM, iklan, kalender, dan analitik</p>
      </div>
    </div>

    <!-- Section Tabs -->
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

    <!-- Bank Konten -->
    {#if activeSection === 'bank'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title">
            <span class="material-symbols-outlined">perm_media</span>
            Bank Konten
          </div>
          <div class="cat-filters">
            {#each ['semua','foto','video','infografis','template'] as cat}
              <button
                class="cat-btn"
                class:active={categoryFilter === cat}
                onclick={() => { categoryFilter = cat as Category; }}
              >{cat.charAt(0).toUpperCase() + cat.slice(1)}</button>
            {/each}
          </div>
        </div>
        <div class="content-grid">
          {#each filteredContent as item (item.id)}
            <div class="content-card">
              <div class="content-thumb" style="background:{item.thumb}">
                <span class="material-symbols-outlined thumb-icon">
                  {item.category === 'video' ? 'play_circle' : item.category === 'infografis' ? 'bar_chart' : 'image'}
                </span>
              </div>
              <div class="content-info">
                <div class="content-title">{item.title}</div>
                <div class="content-tags">
                  {#each item.tags as tag}
                    <span class="tag">#{tag}</span>
                  {/each}
                </div>
              </div>
              <div class="content-actions">
                <button class="sm-btn">
                  <span class="material-symbols-outlined">download</span>
                </button>
                <button class="sm-btn">
                  <span class="material-symbols-outlined">share</span>
                </button>
              </div>
            </div>
          {/each}
        </div>
      </div>

    <!-- Watermark Flyer -->
    {:else if activeSection === 'flyer'}
      <div class="section-block">
        <div class="section-title">
          <span class="material-symbols-outlined">water_drop</span>
          Watermark Flyer
        </div>
        <div class="flyer-layout">
          <div class="upload-area">
            <span class="material-symbols-outlined upload-icon">upload_file</span>
            <p class="upload-text">Seret file atau klik untuk unggah</p>
            <p class="upload-sub">PNG, JPG, atau PDF — maks. 10 MB</p>
            <button class="primary-btn">Pilih File</button>
          </div>
          <div class="flyer-preview">
            <div class="preview-placeholder">
              <span class="material-symbols-outlined">image</span>
              <p>Preview watermark akan tampil di sini</p>
            </div>
            <div class="flyer-controls">
              <div class="field-row">
                <label class="field-label">Teks Watermark</label>
                <input type="text" class="field-input" value="PT Mitra Barokah Tours | 0811-2233-4455" />
              </div>
              <div class="field-row">
                <label class="field-label">Posisi</label>
                <select class="field-input">
                  <option>Kanan Bawah</option>
                  <option>Kiri Bawah</option>
                  <option>Tengah</option>
                </select>
              </div>
              <div class="field-row">
                <label class="field-label">Opasitas</label>
                <input type="range" min="10" max="100" value="70" />
              </div>
              <button class="primary-btn" style="margin-top:0.5rem">
                <span class="material-symbols-outlined">download</span>
                Unduh dengan Watermark
              </button>
            </div>
          </div>
        </div>
      </div>

    <!-- UTM Link Builder -->
    {:else if activeSection === 'utm'}
      <div class="section-block">
        <div class="section-title">
          <span class="material-symbols-outlined">link</span>
          UTM Link Builder
        </div>
        <div class="utm-layout">
          <div class="utm-form">
            <div class="field-row">
              <label class="field-label">Base URL</label>
              <input type="url" class="field-input" bind:value={utmBase} />
            </div>
            <div class="field-row">
              <label class="field-label">UTM Source</label>
              <input type="text" class="field-input" bind:value={utmSource} placeholder="whatsapp, instagram, google..." />
            </div>
            <div class="field-row">
              <label class="field-label">UTM Medium</label>
              <input type="text" class="field-input" bind:value={utmMedium} placeholder="sosial, email, cpc..." />
            </div>
            <div class="field-row">
              <label class="field-label">UTM Campaign</label>
              <input type="text" class="field-input" bind:value={utmCampaign} placeholder="ramadan2026, harbolnas..." />
            </div>
            <div class="field-row">
              <label class="field-label">UTM Content (Agen)</label>
              <input type="text" class="field-input" bind:value={utmAgent} placeholder="kode-agen" />
            </div>
            <div class="utm-result">
              <div class="utm-result-label">Link yang Dihasilkan:</div>
              <div class="utm-result-url">{generatedUtm}</div>
              <button class="primary-btn">
                <span class="material-symbols-outlined">content_copy</span>
                Salin Link
              </button>
            </div>
          </div>
          <div class="utm-table-wrap">
            <div class="sub-section-title">Link Tersimpan</div>
            <div class="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Agen</th>
                    <th>Campaign</th>
                    <th>Source</th>
                    <th>Klik</th>
                    <th>Leads</th>
                  </tr>
                </thead>
                <tbody>
                  {#each savedLinks as link (link.id)}
                    <tr>
                      <td>{link.agent}</td>
                      <td>{link.campaign}</td>
                      <td>{link.source}</td>
                      <td>{link.clicks.toLocaleString('id-ID')}</td>
                      <td>{link.leads}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

    <!-- Ads Manager Lite -->
    {:else if activeSection === 'ads'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title">
            <span class="material-symbols-outlined">campaign</span>
            Ads Manager Lite
          </div>
          <button class="primary-btn">
            <span class="material-symbols-outlined">add</span>
            Buat Kampanye
          </button>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Nama Kampanye</th>
                <th>Platform</th>
                <th>Budget</th>
                <th>Terpakai</th>
                <th>Impresi</th>
                <th>Leads</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {#each campaigns as camp (camp.id)}
                <tr>
                  <td><span class="font-semibold">{camp.name}</span></td>
                  <td>
                    <span class="platform-badge platform-{camp.platform.toLowerCase()}">{camp.platform}</span>
                  </td>
                  <td>{fmtRp(camp.budget)}</td>
                  <td>
                    <div class="spend-bar">
                      <div class="spend-fill" style="width:{Math.min(camp.spend/camp.budget*100,100)}%"></div>
                    </div>
                    <span class="spend-text">{fmtRp(camp.spend)}</span>
                  </td>
                  <td>{camp.impressions.toLocaleString('id-ID')}</td>
                  <td>{camp.leads}</td>
                  <td>
                    <span class="ad-status ad-status--{camp.status}">
                      {camp.status === 'active' ? 'Aktif' : 'Selesai'}
                    </span>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>

    <!-- Kalender Konten -->
    {:else if activeSection === 'calendar'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title">
            <span class="material-symbols-outlined">calendar_month</span>
            Kalender Konten — Minggu Ini
          </div>
          <button class="primary-btn">
            <span class="material-symbols-outlined">add</span>
            Tambah Jadwal
          </button>
        </div>
        <div class="calendar-grid">
          {#each weekDays as day, i}
            <div class="cal-day">
              <div class="cal-day-header">{day}</div>
              <div class="cal-day-body">
                {#each calendarItems.filter(c => c.day === i) as item}
                  <div class="cal-item cal-item--{item.type}">
                    <span class="cal-item-title">{item.title}</span>
                  </div>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      </div>

    <!-- Analitik Konten -->
    {:else if activeSection === 'analytics'}
      <div class="section-block">
        <div class="section-title">
          <span class="material-symbols-outlined">analytics</span>
          Analitik Konten — Bulan Ini
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Channel</th>
                <th>Reach</th>
                <th>Engagement</th>
                <th>CTR</th>
                <th>Leads</th>
              </tr>
            </thead>
            <tbody>
              {#each analyticsData as row}
                <tr>
                  <td><span class="font-semibold">{row.channel}</span></td>
                  <td>{row.reach.toLocaleString('id-ID')}</td>
                  <td>{row.engagement.toLocaleString('id-ID')}</td>
                  <td>
                    <span class="ctr-badge">{row.ctr}%</span>
                  </td>
                  <td><span class="leads-num">{row.leads}</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="analytics-total">
          Total Reach: {analyticsData.reduce((s,r)=>s+r.reach,0).toLocaleString('id-ID')} &nbsp;|&nbsp;
          Total Leads: {analyticsData.reduce((s,r)=>s+r.leads,0)}
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
  .section-header { display: flex; align-items: center; justify-content: space-between; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .section-title { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; font-weight: 700; color: #191c1e; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .section-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }
  .section-header .section-title { padding: 0; border-bottom: 0; }

  .primary-btn { display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.5rem 1rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.3rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; text-decoration: none; }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  /* Bank Konten */
  .cat-filters { display: flex; gap: 0.25rem; flex-wrap: wrap; }
  .cat-btn { padding: 0.3rem 0.65rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.72rem; color: #434655; cursor: pointer; }
  .cat-btn:hover { background: #f2f4f6; }
  .cat-btn.active { border-color: #2563eb; color: #004ac6; background: #eff6ff; font-weight: 700; }

  .content-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; padding: 1rem; }
  .content-card { border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.35rem; overflow: hidden; }
  .content-thumb { height: 100px; display: grid; place-items: center; }
  .thumb-icon { font-size: 2rem; color: #64748b; }
  .content-info { padding: 0.55rem 0.65rem; }
  .content-title { font-size: 0.75rem; font-weight: 600; color: #191c1e; line-height: 1.3; }
  .content-tags { display: flex; gap: 0.25rem; flex-wrap: wrap; margin-top: 0.3rem; }
  .tag { font-size: 0.6rem; color: #737686; background: #f1f5f9; padding: 0.1rem 0.3rem; border-radius: 0.15rem; }
  .content-actions { display: flex; gap: 0.25rem; padding: 0 0.65rem 0.55rem; }
  .sm-btn { display: grid; place-items: center; width: 1.8rem; height: 1.8rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; cursor: pointer; }
  .sm-btn .material-symbols-outlined { font-size: 1rem; color: #434655; }
  .sm-btn:hover { background: #f2f4f6; }

  /* Flyer */
  .flyer-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .flyer-layout { grid-template-columns: 1fr; } }
  .upload-area { border: 2px dashed rgb(195 198 215 / 0.6); border-radius: 0.4rem; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 0.5rem; padding: 2rem; text-align: center; }
  .upload-icon { font-size: 2.5rem; color: #b0b3c1; }
  .upload-text { margin: 0; font-size: 0.85rem; font-weight: 600; color: #434655; }
  .upload-sub { margin: 0; font-size: 0.72rem; color: #b0b3c1; }
  .flyer-preview { display: flex; flex-direction: column; gap: 0.85rem; }
  .preview-placeholder { height: 160px; background: #f2f4f6; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.3rem; display: grid; place-items: center; color: #b0b3c1; text-align: center; }
  .preview-placeholder .material-symbols-outlined { font-size: 2rem; display: block; margin-bottom: 0.35rem; }
  .preview-placeholder p { margin: 0; font-size: 0.72rem; }
  .flyer-controls { display: flex; flex-direction: column; gap: 0.55rem; }
  .field-row { display: flex; flex-direction: column; gap: 0.2rem; }
  .field-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field-input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; }

  /* UTM */
  .utm-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .utm-layout { grid-template-columns: 1fr; } }
  .utm-form { display: flex; flex-direction: column; gap: 0.65rem; }
  .utm-result { background: #f2f4f6; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.3rem; padding: 0.75rem; margin-top: 0.5rem; }
  .utm-result-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; margin-bottom: 0.3rem; }
  .utm-result-url { font-size: 0.7rem; color: #004ac6; word-break: break-all; margin-bottom: 0.5rem; font-family: monospace; }
  .utm-table-wrap { display: flex; flex-direction: column; gap: 0.5rem; }
  .sub-section-title { font-size: 0.78rem; font-weight: 700; color: #191c1e; }

  /* Ads */
  .font-semibold { font-weight: 600; color: #191c1e; }
  .platform-badge { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .platform-meta { background: #dbeafe; color: #1e40af; }
  .platform-google { background: #fef3c7; color: #92400e; }
  .platform-tiktok { background: #fce7f3; color: #9d174d; }
  .spend-bar { height: 4px; background: #e2e8f0; border-radius: 999px; width: 80px; margin-bottom: 0.2rem; }
  .spend-fill { height: 100%; background: #2563eb; border-radius: 999px; }
  .spend-text { font-size: 0.68rem; color: #434655; }
  .ad-status { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .ad-status--active { background: #d1fae5; color: #065f46; }
  .ad-status--finished { background: #f1f5f9; color: #334155; }

  /* Calendar */
  .calendar-grid { display: grid; grid-template-columns: repeat(7, 1fr); padding: 0.75rem; gap: 0.5rem; }
  @media (max-width: 700px) { .calendar-grid { grid-template-columns: repeat(3, 1fr); } }
  .cal-day { border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.3rem; min-height: 100px; }
  .cal-day-header { padding: 0.4rem 0.5rem; font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: #434655; border-bottom: 1px solid rgb(195 198 215 / 0.35); background: #f2f4f6; }
  .cal-day-body { padding: 0.35rem; display: flex; flex-direction: column; gap: 0.25rem; }
  .cal-item { padding: 0.3rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 600; line-height: 1.3; }
  .cal-item--post { background: #dbeafe; color: #1e40af; }
  .cal-item--blast { background: #d1fae5; color: #065f46; }
  .cal-item--story { background: #fce7f3; color: #9d174d; }
  .cal-item--video { background: #fef3c7; color: #92400e; }
  .cal-item-title { display: block; }

  /* Analytics */
  .ctr-badge { background: #eff6ff; color: #1d4ed8; padding: 0.1rem 0.35rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .leads-num { font-weight: 700; color: #059669; }
  .analytics-total { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; font-size: 0.72rem; color: #434655; }

  /* Shared table */
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
</style>
