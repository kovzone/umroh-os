<script lang="ts">
  interface CopyVariation {
    id: string;
    label: string;
    tone: string;
    color: string;
    headline: string;
    body: string;
    cta: string;
  }

  // Form state
  let packageName = $state('');
  let durationDays = $state('');
  let destinations = $state('');
  let highlights = $state('');
  let departureDate = $state('');
  let contactInfo = $state('0812-3456-7890');

  // Generation state
  let loading = $state(false);
  let variations = $state<CopyVariation[]>([]);
  let generated = $state(false);
  let copiedId = $state<string | null>(null);

  const isFormValid = $derived(
    packageName.trim().length > 2 &&
    durationDays.trim().length > 0 &&
    destinations.trim().length > 2
  );

  function generateVariations() {
    if (!isFormValid) return;
    loading = true;
    generated = false;
    variations = [];

    // Simulate AI thinking with a timeout
    setTimeout(() => {
      const name = packageName.trim();
      const days = durationDays.trim();
      const dest = destinations.trim();
      const highs = highlights.trim() || 'hotel bintang 5, muthawif berpengalaman, visa garansi';
      const date = departureDate
        ? new Date(departureDate).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
        : 'segera';
      const contact = contactInfo.trim();

      variations = [
        {
          id: 'v1',
          label: 'Variasi 1',
          tone: 'Professional',
          color: '#1e40af',
          headline: `Wujudkan Ibadah Umroh Impian Bersama ${name}`,
          body: `Kami menghadirkan paket Umroh ${days} hari ke ${dest} yang dirancang untuk kenyamanan dan kekhusyukan ibadah Anda. Dengan fasilitas ${highs}, setiap momen perjalanan suci Anda akan terasa istimewa. Tim profesional kami siap mendampingi dari keberangkatan hingga kepulangan.`,
          cta: `Konsultasikan perjalanan Umroh Anda sekarang. Hubungi: ${contact}`
        },
        {
          id: 'v2',
          label: 'Variasi 2',
          tone: 'Emosional',
          color: '#065f46',
          headline: `Langkah Pertama Menuju Tanah Suci Dimulai dari Sini`,
          body: `Setiap Muslim mendambakan momen bersujud di Masjidil Haram. Dengan paket ${name}, impian itu kini dalam jangkauan tangan Anda. Selama ${days} hari di ${dest}, rasakan kedamaian yang tiada duanya — jauh dari hiruk-pikuk dunia, dekat dengan Sang Pencipta. ${highs.charAt(0).toUpperCase() + highs.slice(1)} — semua untuk memastikan ibadah Anda sempurna.`,
          cta: `Jangan tunda lagi. Hubungi kami di ${contact} dan jadwalkan keberangkatan Anda.`
        },
        {
          id: 'v3',
          label: 'Variasi 3',
          tone: 'Urgensi',
          color: '#9f1239',
          headline: `Kursi Terbatas! Paket ${name} Berangkat ${date}`,
          body: `Hanya tersisa beberapa kursi untuk paket ${name} — ${days} hari ${dest} dengan ${highs}. Jangan sampai menyesal! Ribuan jamaah sudah mempercayakan perjalanan ibadah mereka kepada kami. Daftar sekarang sebelum terlambat dan amankan seat Anda sebelum penuh.`,
          cta: `DAFTAR SEKARANG sebelum kehabisan! WA/Telp: ${contact}`
        }
      ];

      loading = false;
      generated = true;
    }, 1800);
  }

  function reset() {
    variations = [];
    generated = false;
    packageName = '';
    durationDays = '';
    destinations = '';
    highlights = '';
    departureDate = '';
  }

  function copyVariation(v: CopyVariation) {
    const text = `${v.headline}\n\n${v.body}\n\n${v.cta}`;
    navigator.clipboard.writeText(text).then(() => {
      copiedId = v.id;
      setTimeout(() => { copiedId = null; }, 2000);
    });
  }

  function copyAll() {
    const text = variations.map((v, i) =>
      `--- ${v.label} (${v.tone}) ---\n\nHeadline:\n${v.headline}\n\nBody:\n${v.body}\n\nCTA:\n${v.cta}`
    ).join('\n\n');
    navigator.clipboard.writeText(text).then(() => {
      copiedId = 'all';
      setTimeout(() => { copiedId = null; }, 2000);
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
      <span class="topbar-current">Copywriting Automation</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn" title="Notifikasi">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar" aria-label="Profile">AD</button>
    </div>
  </header>

  <section class="canvas-page">
    <div class="page-head">
      <div>
        <h2>Copywriting Automation</h2>
        <p>BL-CRM-016 — Generate 3 variasi teks marketing dari detail paket Umroh</p>
      </div>
    </div>

    <!-- AI Stub Notice -->
    <div class="ai-notice">
      <span class="material-symbols-outlined">auto_awesome</span>
      <div>
        <strong>AI Integration Pending</strong>
        <span>— Teks dihasilkan dari template cerdas. AI integration pending API key.</span>
      </div>
    </div>

    <div class="layout">
      <!-- Input form -->
      <div class="form-panel">
        <div class="panel-title">
          <span class="material-symbols-outlined">edit_note</span>
          Detail Paket
        </div>

        <div class="form-grid">
          <div class="field-group span-2">
            <label class="field-label" for="pkg-name">
              Nama Paket <span class="required">*</span>
            </label>
            <input
              id="pkg-name"
              type="text"
              class="field-input"
              bind:value={packageName}
              placeholder="contoh: Umroh Ramadhan Premium 2026"
            />
          </div>

          <div class="field-group">
            <label class="field-label" for="duration">
              Durasi (hari) <span class="required">*</span>
            </label>
            <input
              id="duration"
              type="number"
              class="field-input"
              bind:value={durationDays}
              placeholder="15"
              min="1"
              max="90"
            />
          </div>

          <div class="field-group">
            <label class="field-label" for="departure">Tanggal Keberangkatan</label>
            <input
              id="departure"
              type="date"
              class="field-input"
              bind:value={departureDate}
            />
          </div>

          <div class="field-group span-2">
            <label class="field-label" for="destinations">
              Destinasi <span class="required">*</span>
            </label>
            <input
              id="destinations"
              type="text"
              class="field-input"
              bind:value={destinations}
              placeholder="contoh: Makkah, Madinah, Turki"
            />
          </div>

          <div class="field-group span-2">
            <label class="field-label" for="highlights">Highlights & Fasilitas</label>
            <textarea
              id="highlights"
              class="field-textarea"
              rows="3"
              bind:value={highlights}
              placeholder="contoh: Hotel bintang 5, makan 3x sehari, muthawif profesional, ziarah full"
            ></textarea>
          </div>

          <div class="field-group span-2">
            <label class="field-label" for="contact">Info Kontak</label>
            <input
              id="contact"
              type="text"
              class="field-input"
              bind:value={contactInfo}
              placeholder="Nomor HP / media sosial"
            />
          </div>
        </div>

        <div class="form-actions">
          {#if generated}
            <button type="button" class="ghost-btn" onclick={reset}>
              <span class="material-symbols-outlined">refresh</span>
              Reset
            </button>
          {/if}
          <button
            type="button"
            class="generate-btn"
            onclick={generateVariations}
            disabled={!isFormValid || loading}
          >
            {#if loading}
              <span class="material-symbols-outlined spin">progress_activity</span>
              Generating...
            {:else}
              <span class="material-symbols-outlined">auto_awesome</span>
              Generate Copywriting
            {/if}
          </button>
        </div>
      </div>

      <!-- Results -->
      <div class="results-panel">
        {#if loading}
          <div class="loading-state">
            <div class="loading-spinner">
              <span class="material-symbols-outlined spin-large">auto_awesome</span>
            </div>
            <div class="loading-text">
              <strong>Sedang generate variasi teks...</strong>
              <span>Menganalisis detail paket dan menyesuaikan tone</span>
            </div>
          </div>

        {:else if generated && variations.length > 0}
          <div class="results-header">
            <div class="results-title">
              <span class="material-symbols-outlined">check_circle</span>
              3 Variasi Berhasil Dibuat
            </div>
            <button
              type="button"
              class="copy-all-btn"
              onclick={copyAll}
            >
              <span class="material-symbols-outlined">
                {copiedId === 'all' ? 'check' : 'content_copy'}
              </span>
              {copiedId === 'all' ? 'Tersalin!' : 'Copy Semua'}
            </button>
          </div>

          <div class="variations-list">
            {#each variations as v}
              <div class="variation-card" style="--v-color: {v.color}">
                <div class="var-header">
                  <div class="var-meta">
                    <span class="var-label" style="color: {v.color}">{v.label}</span>
                    <span class="var-tone" style="background: color-mix(in srgb, {v.color} 12%, white); color: {v.color}">
                      {v.tone}
                    </span>
                  </div>
                  <button
                    type="button"
                    class="copy-btn"
                    onclick={() => copyVariation(v)}
                    title="Copy ke clipboard"
                  >
                    <span class="material-symbols-outlined">
                      {copiedId === v.id ? 'check' : 'content_copy'}
                    </span>
                    {copiedId === v.id ? 'Tersalin!' : 'Copy'}
                  </button>
                </div>

                <div class="var-section">
                  <div class="var-section-label">Headline</div>
                  <div class="var-headline" style="color: {v.color}">{v.headline}</div>
                </div>

                <div class="var-section">
                  <div class="var-section-label">Body</div>
                  <div class="var-body">{v.body}</div>
                </div>

                <div class="var-section">
                  <div class="var-section-label">CTA</div>
                  <div class="var-cta">{v.cta}</div>
                </div>
              </div>
            {/each}
          </div>

        {:else}
          <div class="empty-results">
            <span class="material-symbols-outlined">auto_awesome</span>
            <strong>Siap Generate</strong>
            <p>Isi detail paket di sebelah kiri, lalu klik "Generate Copywriting" untuk mendapatkan 3 variasi teks marketing.</p>
            <div class="tips">
              <div class="tip-item">
                <span class="material-symbols-outlined">lightbulb</span>
                <span>Semakin lengkap detail paket, semakin relevan hasilnya</span>
              </div>
              <div class="tip-item">
                <span class="material-symbols-outlined">lightbulb</span>
                <span>3 tone berbeda: Professional, Emosional, dan Urgensi</span>
              </div>
              <div class="tip-item">
                <span class="material-symbols-outlined">lightbulb</span>
                <span>Setiap variasi langsung bisa di-copy ke WhatsApp / Instagram</span>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
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
  .top-actions { display: flex; align-items: center; gap: 0.35rem; }
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas-page { padding: 1.5rem; max-width: 100rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  /* AI notice */
  .ai-notice {
    display: flex; align-items: flex-start; gap: 0.6rem;
    padding: 0.75rem 1rem;
    background: #fefce8;
    border: 1px solid #fde047;
    border-radius: 0.25rem;
    font-size: 0.8rem;
    color: #713f12;
    margin-bottom: 1.25rem;
  }
  .ai-notice .material-symbols-outlined { font-size: 1rem; color: #ca8a04; flex-shrink: 0; margin-top: 1px; }
  .ai-notice strong { font-weight: 700; }

  /* Layout */
  .layout {
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 1.5rem;
    align-items: start;
  }

  @media (max-width: 900px) {
    .layout { grid-template-columns: 1fr; }
  }

  /* Form panel */
  .form-panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1.25rem;
  }

  .panel-title {
    display: flex; align-items: center; gap: 0.4rem;
    font-size: 0.88rem; font-weight: 700; color: #191c1e;
    margin-bottom: 1rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
  }
  .panel-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }

  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.85rem;
    margin-bottom: 1rem;
  }

  .field-group { display: flex; flex-direction: column; gap: 0.3rem; }
  .span-2 { grid-column: span 2; }

  .field-label {
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }
  .required { color: #ba1a1a; }

  .field-input, .field-textarea {
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    padding: 0.5rem 0.65rem;
    font-size: 0.82rem;
    color: #191c1e;
    background: #fff;
    font-family: inherit;
    outline: none;
    width: 100%;
    box-sizing: border-box;
  }
  .field-input:focus, .field-textarea:focus { border-color: #2563eb; }
  .field-textarea { resize: vertical; min-height: 5rem; }

  .form-actions { display: flex; gap: 0.5rem; justify-content: flex-end; }

  .ghost-btn {
    display: inline-flex; align-items: center; gap: 0.3rem;
    padding: 0.55rem 0.85rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    background: #fff;
    font-size: 0.82rem; font-weight: 600; color: #434655;
    cursor: pointer; font-family: inherit;
  }
  .ghost-btn:hover { background: #f2f4f6; }
  .ghost-btn .material-symbols-outlined { font-size: 0.9rem; }

  .generate-btn {
    display: inline-flex; align-items: center; gap: 0.4rem;
    padding: 0.55rem 1rem;
    background: linear-gradient(90deg, #7c3aed, #6d28d9);
    color: #fff; border: 0; border-radius: 0.25rem;
    font-size: 0.85rem; font-weight: 700;
    cursor: pointer; font-family: inherit;
  }
  .generate-btn:hover { opacity: 0.9; }
  .generate-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .generate-btn .material-symbols-outlined { font-size: 1rem; }

  /* Results panel */
  .results-panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    min-height: 24rem;
    overflow: hidden;
  }

  /* Loading */
  .loading-state {
    display: flex; flex-direction: column; align-items: center; gap: 1.25rem;
    padding: 4rem 2rem;
  }

  .loading-spinner {
    width: 4.5rem; height: 4.5rem;
    background: #f3e8ff;
    border-radius: 50%;
    display: grid; place-items: center;
  }

  .spin-large { font-size: 2rem; color: #7c3aed; animation: spin 1.2s linear infinite; }

  .loading-text {
    display: flex; flex-direction: column; align-items: center; gap: 0.3rem;
    text-align: center;
  }
  .loading-text strong { font-size: 0.95rem; color: #191c1e; }
  .loading-text span { font-size: 0.78rem; color: #737686; }

  /* Results header */
  .results-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 0.9rem 1.1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
    background: #f0fdf4;
  }

  .results-title {
    display: flex; align-items: center; gap: 0.4rem;
    font-size: 0.88rem; font-weight: 700; color: #065f46;
  }
  .results-title .material-symbols-outlined { font-size: 1rem; }

  .copy-all-btn {
    display: inline-flex; align-items: center; gap: 0.3rem;
    padding: 0.38rem 0.75rem;
    border: 1px solid #065f46;
    border-radius: 0.25rem;
    background: #fff;
    font-size: 0.75rem; font-weight: 600; color: #065f46;
    cursor: pointer; font-family: inherit;
  }
  .copy-all-btn:hover { background: #f0fdf4; }
  .copy-all-btn .material-symbols-outlined { font-size: 0.85rem; }

  /* Variations list */
  .variations-list {
    display: flex; flex-direction: column;
    padding: 1rem;
    gap: 1rem;
  }

  .variation-card {
    border: 1.5px solid color-mix(in srgb, var(--v-color) 25%, #e5e7eb);
    border-radius: 0.35rem;
    overflow: hidden;
  }

  .var-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 0.65rem 0.85rem;
    background: color-mix(in srgb, var(--v-color) 6%, white);
    border-bottom: 1px solid color-mix(in srgb, var(--v-color) 15%, #e5e7eb);
  }

  .var-meta { display: flex; align-items: center; gap: 0.5rem; }
  .var-label { font-size: 0.85rem; font-weight: 700; }
  .var-tone { padding: 0.15rem 0.5rem; border-radius: 999px; font-size: 0.65rem; font-weight: 700; }

  .copy-btn {
    display: inline-flex; align-items: center; gap: 0.25rem;
    padding: 0.3rem 0.6rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.2rem;
    background: #fff; font-size: 0.72rem; font-weight: 600;
    color: #434655; cursor: pointer; font-family: inherit;
  }
  .copy-btn:hover { background: #f2f4f6; }
  .copy-btn .material-symbols-outlined { font-size: 0.82rem; }

  .var-section { padding: 0.65rem 0.85rem; border-bottom: 1px solid rgb(195 198 215 / 0.25); }
  .var-section:last-child { border-bottom: 0; }

  .var-section-label {
    font-size: 0.6rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #737686;
    margin-bottom: 0.35rem;
  }

  .var-headline { font-size: 0.92rem; font-weight: 700; line-height: 1.35; }
  .var-body { font-size: 0.8rem; color: #374151; line-height: 1.55; }
  .var-cta { font-size: 0.8rem; font-weight: 600; color: #374151; font-style: italic; }

  /* Empty results */
  .empty-results {
    display: flex; flex-direction: column; align-items: center; gap: 0.65rem;
    padding: 3.5rem 2rem;
    color: #b0b3c1;
    text-align: center;
  }
  .empty-results .material-symbols-outlined { font-size: 2.75rem; color: #d1d5db; }
  .empty-results strong { font-size: 0.95rem; color: #434655; }
  .empty-results p { margin: 0; font-size: 0.82rem; color: #737686; max-width: 28rem; line-height: 1.5; }

  .tips {
    display: flex; flex-direction: column; gap: 0.4rem;
    margin-top: 0.5rem;
    text-align: left;
    width: 100%;
    max-width: 24rem;
  }

  .tip-item {
    display: flex; align-items: flex-start; gap: 0.4rem;
    font-size: 0.75rem; color: #737686;
  }
  .tip-item .material-symbols-outlined { font-size: 0.85rem; color: #fbbf24; flex-shrink: 0; margin-top: 1px; }

  /* Animations */
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
