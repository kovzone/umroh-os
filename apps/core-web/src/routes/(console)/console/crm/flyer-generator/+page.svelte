<script lang="ts">
  import { onMount } from 'svelte';

  // ---- types ----
  interface PackageOption {
    id: string;
    name: string;
    price: number;
    departure_date: string;
    duration_days: number;
  }

  interface FlyerConfig {
    packageId: string;
    template: 'elegant' | 'bold' | 'minimal';
    headline: string;
    description: string;
    priceCallout: string;
    contactInfo: string;
  }

  // ---- mock packages ----
  const MOCK_PACKAGES: PackageOption[] = [
    { id: 'pkg1', name: 'Umroh Ramadhan Premium', price: 32_500_000, departure_date: '2026-03-01', duration_days: 15 },
    { id: 'pkg2', name: 'Umroh Reguler April', price: 24_500_000, departure_date: '2026-04-10', duration_days: 12 },
    { id: 'pkg3', name: 'Umroh Plus Turki', price: 38_000_000, departure_date: '2026-05-05', duration_days: 18 },
    { id: 'pkg4', name: 'Umroh Hemat Juni', price: 19_900_000, departure_date: '2026-06-15', duration_days: 10 },
    { id: 'pkg5', name: 'Umroh VIP Juli', price: 55_000_000, departure_date: '2026-07-01', duration_days: 14 }
  ];

  const TEMPLATES = [
    { id: 'elegant', label: 'Elegant', desc: 'Emas & gelap — kesan mewah' },
    { id: 'bold', label: 'Bold', desc: 'Biru cerah — energik & modern' },
    { id: 'minimal', label: 'Minimal', desc: 'Putih bersih — profesional' }
  ] as const;

  // ---- state ----
  let packages = $state<PackageOption[]>(MOCK_PACKAGES);
  let config = $state<FlyerConfig>({
    packageId: MOCK_PACKAGES[0].id,
    template: 'elegant',
    headline: 'Wujudkan Ibadah Umroh Impianmu',
    description: 'Paket lengkap dengan hotel bintang 5, pembimbing berpengalaman, dan layanan premium.',
    priceCallout: 'Mulai dari',
    contactInfo: 'Hubungi: 0812-3456-7890 | @umrohos'
  });
  let canvasEl = $state<HTMLCanvasElement | null>(null);
  let generating = $state(false);

  const selectedPackage = $derived(packages.find(p => p.id === config.packageId) ?? packages[0]);

  function formatIDR(n: number): string {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' });
  }

  // ---- canvas drawing ----
  function drawFlyer() {
    if (!canvasEl) return;
    const ctx = canvasEl.getContext('2d');
    if (!ctx) return;

    const W = 540; // display at half of 1080 for canvas
    const H = 960; // display at half of 1920
    canvasEl.width = W;
    canvasEl.height = H;

    const pkg = selectedPackage;

    if (config.template === 'elegant') {
      // Dark gold gradient background
      const grad = ctx.createLinearGradient(0, 0, 0, H);
      grad.addColorStop(0, '#1a0a00');
      grad.addColorStop(0.5, '#2d1500');
      grad.addColorStop(1, '#0d0600');
      ctx.fillStyle = grad;
      ctx.fillRect(0, 0, W, H);

      // Gold accent bar at top
      const goldGrad = ctx.createLinearGradient(0, 0, W, 0);
      goldGrad.addColorStop(0, '#c8860a');
      goldGrad.addColorStop(0.5, '#f0c040');
      goldGrad.addColorStop(1, '#c8860a');
      ctx.fillStyle = goldGrad;
      ctx.fillRect(0, 0, W, 6);

      // Logo circle
      ctx.beginPath();
      ctx.arc(W / 2, 90, 44, 0, Math.PI * 2);
      ctx.fillStyle = '#c8860a';
      ctx.fill();
      ctx.fillStyle = '#fff';
      ctx.font = 'bold 28px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText('O', W / 2, 90);

      // Agency name
      ctx.fillStyle = '#f0c040';
      ctx.font = '600 15px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText('UmrohOS Travel', W / 2, 160);

      // Divider line
      ctx.fillStyle = '#c8860a';
      ctx.fillRect(W / 2 - 60, 172, 120, 1.5);

      // Headline
      ctx.fillStyle = '#ffffff';
      ctx.font = 'bold 30px sans-serif';
      wrapText(ctx, config.headline, W / 2, 220, W - 60, 38, 'center');

      // Description
      ctx.fillStyle = '#e8c87a';
      ctx.font = '400 15px sans-serif';
      wrapText(ctx, config.description, W / 2, 310, W - 80, 22, 'center');

      // Package name band
      ctx.fillStyle = '#c8860a';
      roundRect(ctx, 30, 400, W - 60, 70, 8);
      ctx.fill();
      ctx.fillStyle = '#fff';
      ctx.font = 'bold 19px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(pkg.name, W / 2, 435);

      // Departure info
      ctx.fillStyle = '#f0c040';
      ctx.font = '500 14px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(`Berangkat: ${formatDate(pkg.departure_date)}  •  ${pkg.duration_days} hari`, W / 2, 500);

      // Price callout box
      const priceGrad = ctx.createLinearGradient(0, 540, 0, 640);
      priceGrad.addColorStop(0, '#c8860a');
      priceGrad.addColorStop(1, '#8b5e00');
      ctx.fillStyle = priceGrad;
      roundRect(ctx, 60, 540, W - 120, 90, 10);
      ctx.fill();
      ctx.fillStyle = '#fff8e1';
      ctx.font = '400 13px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(config.priceCallout, W / 2, 572);
      ctx.fillStyle = '#ffffff';
      ctx.font = 'bold 30px sans-serif';
      ctx.fillText(formatIDR(pkg.price), W / 2, 614);

      // Features icons row
      const features = ['Visa Umroh', 'Hotel Bintang 5', 'Full AC Bus', 'Muthawif'];
      const featureX = [W * 0.15, W * 0.38, W * 0.62, W * 0.85];
      ctx.fillStyle = '#e8c87a';
      ctx.font = '12px sans-serif';
      ctx.textAlign = 'center';
      features.forEach((f, i) => {
        ctx.fillStyle = '#c8860a';
        ctx.beginPath();
        ctx.arc(featureX[i], 688, 18, 0, Math.PI * 2);
        ctx.fill();
        ctx.fillStyle = '#fff';
        ctx.font = 'bold 10px sans-serif';
        const initials = f.split(' ').map(w => w[0]).join('').substring(0, 2);
        ctx.fillText(initials, featureX[i], 692);
        ctx.fillStyle = '#e8c87a';
        ctx.font = '11px sans-serif';
        ctx.fillText(f, featureX[i], 716);
      });

      // CTA button
      const ctaGrad = ctx.createLinearGradient(0, 760, 0, 820);
      ctaGrad.addColorStop(0, '#f0c040');
      ctaGrad.addColorStop(1, '#c8860a');
      ctx.fillStyle = ctaGrad;
      roundRect(ctx, W / 2 - 110, 760, 220, 56, 28);
      ctx.fill();
      ctx.fillStyle = '#1a0a00';
      ctx.font = 'bold 18px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText('DAFTAR SEKARANG', W / 2, 788);

      // Contact
      ctx.fillStyle = '#e8c87a';
      ctx.font = '13px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(config.contactInfo, W / 2, 870);

      // Bottom gold bar
      ctx.fillStyle = goldGrad;
      ctx.fillRect(0, H - 6, W, 6);

    } else if (config.template === 'bold') {
      // Bold blue
      const grad = ctx.createLinearGradient(0, 0, W, H);
      grad.addColorStop(0, '#1e40af');
      grad.addColorStop(1, '#0369a1');
      ctx.fillStyle = grad;
      ctx.fillRect(0, 0, W, H);

      // White geometric accent
      ctx.fillStyle = 'rgba(255,255,255,0.06)';
      ctx.beginPath();
      ctx.arc(W + 60, -60, 260, 0, Math.PI * 2);
      ctx.fill();

      // Top accent bar
      ctx.fillStyle = '#38bdf8';
      ctx.fillRect(0, 0, W, 5);

      // Logo
      ctx.beginPath();
      ctx.arc(W / 2, 90, 44, 0, Math.PI * 2);
      ctx.fillStyle = '#fff';
      ctx.fill();
      ctx.fillStyle = '#1e40af';
      ctx.font = 'bold 28px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText('O', W / 2, 90);

      ctx.fillStyle = '#bfdbfe';
      ctx.font = '600 15px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText('UmrohOS Travel', W / 2, 162);

      ctx.fillStyle = '#93c5fd';
      ctx.fillRect(W / 2 - 50, 172, 100, 2);

      ctx.fillStyle = '#ffffff';
      ctx.font = 'bold 32px sans-serif';
      wrapText(ctx, config.headline, W / 2, 222, W - 60, 40, 'center');

      ctx.fillStyle = '#bfdbfe';
      ctx.font = '15px sans-serif';
      wrapText(ctx, config.description, W / 2, 316, W - 80, 22, 'center');

      ctx.fillStyle = 'rgba(255,255,255,0.15)';
      roundRect(ctx, 30, 400, W - 60, 70, 8);
      ctx.fill();
      ctx.fillStyle = '#fff';
      ctx.font = 'bold 18px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(pkg.name, W / 2, 435);

      ctx.fillStyle = '#93c5fd';
      ctx.font = '14px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(`${formatDate(pkg.departure_date)}  •  ${pkg.duration_days} hari`, W / 2, 500);

      ctx.fillStyle = '#fff';
      roundRect(ctx, 60, 535, W - 120, 90, 10);
      ctx.fill();
      ctx.fillStyle = '#1e40af';
      ctx.font = '13px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(config.priceCallout, W / 2, 564);
      ctx.fillStyle = '#0369a1';
      ctx.font = 'bold 30px sans-serif';
      ctx.fillText(formatIDR(pkg.price), W / 2, 608);

      ctx.fillStyle = '#38bdf8';
      roundRect(ctx, W / 2 - 110, 760, 220, 56, 28);
      ctx.fill();
      ctx.fillStyle = '#fff';
      ctx.font = 'bold 18px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText('DAFTAR SEKARANG', W / 2, 788);

      ctx.fillStyle = '#bfdbfe';
      ctx.font = '13px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(config.contactInfo, W / 2, 872);

      ctx.fillStyle = '#38bdf8';
      ctx.fillRect(0, H - 5, W, 5);

    } else {
      // Minimal white
      ctx.fillStyle = '#ffffff';
      ctx.fillRect(0, 0, W, H);

      ctx.fillStyle = '#2563eb';
      ctx.fillRect(0, 0, W, 5);

      ctx.beginPath();
      ctx.arc(W / 2, 90, 44, 0, Math.PI * 2);
      ctx.fillStyle = '#dbeafe';
      ctx.fill();
      ctx.fillStyle = '#1e40af';
      ctx.font = 'bold 28px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText('O', W / 2, 90);

      ctx.fillStyle = '#434655';
      ctx.font = '600 14px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText('UmrohOS Travel', W / 2, 162);

      ctx.fillStyle = '#e5e7eb';
      ctx.fillRect(30, 178, W - 60, 1);

      ctx.fillStyle = '#111827';
      ctx.font = 'bold 28px sans-serif';
      wrapText(ctx, config.headline, W / 2, 218, W - 60, 36, 'center');

      ctx.fillStyle = '#6b7280';
      ctx.font = '14px sans-serif';
      wrapText(ctx, config.description, W / 2, 308, W - 80, 21, 'center');

      ctx.fillStyle = '#f9fafb';
      ctx.strokeStyle = '#e5e7eb';
      ctx.lineWidth = 1;
      roundRect(ctx, 30, 400, W - 60, 70, 8);
      ctx.fill();
      ctx.stroke();
      ctx.fillStyle = '#111827';
      ctx.font = 'bold 18px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(pkg.name, W / 2, 435);

      ctx.fillStyle = '#6b7280';
      ctx.font = '13px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(`${formatDate(pkg.departure_date)}  •  ${pkg.duration_days} hari`, W / 2, 500);

      ctx.fillStyle = '#f0f9ff';
      ctx.strokeStyle = '#bae6fd';
      roundRect(ctx, 60, 535, W - 120, 90, 10);
      ctx.fill();
      ctx.stroke();
      ctx.fillStyle = '#475569';
      ctx.font = '12px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(config.priceCallout, W / 2, 564);
      ctx.fillStyle = '#0c4a6e';
      ctx.font = 'bold 28px sans-serif';
      ctx.fillText(formatIDR(pkg.price), W / 2, 606);

      ctx.fillStyle = '#2563eb';
      roundRect(ctx, W / 2 - 110, 760, 220, 56, 28);
      ctx.fill();
      ctx.fillStyle = '#fff';
      ctx.font = 'bold 17px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText('DAFTAR SEKARANG', W / 2, 788);

      ctx.fillStyle = '#6b7280';
      ctx.font = '12px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'alphabetic';
      ctx.fillText(config.contactInfo, W / 2, 870);

      ctx.fillStyle = '#2563eb';
      ctx.fillRect(0, H - 5, W, 5);
    }
  }

  function roundRect(ctx: CanvasRenderingContext2D, x: number, y: number, w: number, h: number, r: number) {
    ctx.beginPath();
    ctx.moveTo(x + r, y);
    ctx.lineTo(x + w - r, y);
    ctx.quadraticCurveTo(x + w, y, x + w, y + r);
    ctx.lineTo(x + w, y + h - r);
    ctx.quadraticCurveTo(x + w, y + h, x + w - r, y + h);
    ctx.lineTo(x + r, y + h);
    ctx.quadraticCurveTo(x, y + h, x, y + h - r);
    ctx.lineTo(x, y + r);
    ctx.quadraticCurveTo(x, y, x + r, y);
    ctx.closePath();
  }

  function wrapText(
    ctx: CanvasRenderingContext2D,
    text: string,
    x: number,
    y: number,
    maxWidth: number,
    lineHeight: number,
    align: 'left' | 'center' | 'right' = 'left'
  ) {
    ctx.textAlign = align;
    ctx.textBaseline = 'alphabetic';
    const words = text.split(' ');
    let line = '';
    let currentY = y;
    for (const word of words) {
      const testLine = line ? `${line} ${word}` : word;
      if (ctx.measureText(testLine).width > maxWidth && line) {
        ctx.fillText(line, x, currentY);
        line = word;
        currentY += lineHeight;
      } else {
        line = testLine;
      }
    }
    if (line) ctx.fillText(line, x, currentY);
  }

  $effect(() => {
    // re-draw whenever config changes
    void config.template;
    void config.headline;
    void config.description;
    void config.priceCallout;
    void config.contactInfo;
    void config.packageId;
    drawFlyer();
  });

  onMount(() => {
    drawFlyer();
  });

  function downloadPNG() {
    if (!canvasEl) return;
    generating = true;
    // Create a high-res canvas
    const exportCanvas = document.createElement('canvas');
    exportCanvas.width = 1080;
    exportCanvas.height = 1920;
    const exportCtx = exportCanvas.getContext('2d');
    if (!exportCtx) { generating = false; return; }
    exportCtx.scale(2, 2);
    // Draw onto export canvas by drawing the preview canvas scaled up
    exportCtx.drawImage(canvasEl, 0, 0, 1080, 1920);
    const link = document.createElement('a');
    link.download = `flyer-${config.packageId}-${config.template}.png`;
    link.href = exportCanvas.toDataURL('image/png');
    link.click();
    generating = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">hub</span>
      <span class="sep">/</span>
      <a href="/console/crm" class="crumb-link">CRM Tools</a>
      <span class="sep">/</span>
      <span class="topbar-current">Flyer Generator</span>
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
        <h2>Flyer Generator</h2>
        <p>BL-CRM-013 — Buat flyer promosi paket Umroh, preview langsung, export PNG</p>
      </div>
      <button class="dl-btn" onclick={downloadPNG} disabled={generating}>
        <span class="material-symbols-outlined">{generating ? 'hourglass_top' : 'download'}</span>
        Download PNG
      </button>
    </div>

    <div class="generator-layout">
      <!-- Controls panel -->
      <div class="controls-panel">
        <div class="field-group">
          <label class="field-label" for="pkg-select">Paket Keberangkatan</label>
          <select id="pkg-select" class="field-input" bind:value={config.packageId}>
            {#each packages as p}
              <option value={p.id}>{p.name} — {formatDate(p.departure_date)}</option>
            {/each}
          </select>
        </div>

        <div class="field-group">
          <label class="field-label">Template Desain</label>
          <div class="template-options">
            {#each TEMPLATES as t}
              <button
                type="button"
                class="template-opt"
                class:selected={config.template === t.id}
                onclick={() => { config = { ...config, template: t.id }; }}
              >
                <span class="template-name">{t.label}</span>
                <span class="template-desc">{t.desc}</span>
              </button>
            {/each}
          </div>
        </div>

        <div class="field-group">
          <label class="field-label" for="headline">Headline</label>
          <input
            id="headline"
            type="text"
            class="field-input"
            bind:value={config.headline}
            placeholder="Judul utama flyer"
          />
        </div>

        <div class="field-group">
          <label class="field-label" for="description">Deskripsi</label>
          <textarea
            id="description"
            class="field-textarea"
            rows="3"
            bind:value={config.description}
            placeholder="Deskripsi paket singkat"
          ></textarea>
        </div>

        <div class="field-group">
          <label class="field-label" for="price-label">Label Harga</label>
          <input
            id="price-label"
            type="text"
            class="field-input"
            bind:value={config.priceCallout}
            placeholder="contoh: Mulai dari"
          />
        </div>

        <div class="field-group">
          <label class="field-label" for="contact">Info Kontak</label>
          <input
            id="contact"
            type="text"
            class="field-input"
            bind:value={config.contactInfo}
            placeholder="Nomor HP / media sosial"
          />
        </div>

        {#if selectedPackage}
          <div class="pkg-summary">
            <span class="material-symbols-outlined">info</span>
            <div>
              <strong>{selectedPackage.name}</strong><br/>
              <span>{formatDate(selectedPackage.departure_date)} · {selectedPackage.duration_days} hari · {formatIDR(selectedPackage.price)}</span>
            </div>
          </div>
        {/if}

        <button class="dl-btn-full" onclick={downloadPNG} disabled={generating}>
          <span class="material-symbols-outlined">{generating ? 'hourglass_top' : 'download'}</span>
          Download PNG (1080×1920)
        </button>
      </div>

      <!-- Preview -->
      <div class="preview-wrap">
        <div class="preview-label">
          <span class="material-symbols-outlined">preview</span>
          Preview (skala 50%)
        </div>
        <div class="canvas-frame">
          <canvas bind:this={canvasEl} class="flyer-canvas"></canvas>
        </div>
      </div>
    </div>
  </section>
</main>

<style>
  .page-shell {
    min-height: 100vh;
    background: #f7f9fb;
  }

  .topbar {
    position: sticky;
    top: 0;
    z-index: 30;
    height: 4rem;
    background: rgb(255 255 255 / 0.9);
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    backdrop-filter: blur(8px);
  }

  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.88rem;
    color: #434655;
  }

  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .sep { color: #b0b3c1; }
  .crumb-link { color: #434655; text-decoration: none; }
  .crumb-link:hover { color: #004ac6; }
  .topbar-current { font-weight: 600; color: #191c1e; }

  .top-actions { display: flex; align-items: center; gap: 0.35rem; }

  .icon-btn {
    border: 0; background: transparent; color: #434655;
    width: 2rem; height: 2rem; border-radius: 0.25rem;
    cursor: pointer; display: grid; place-items: center;
  }
  .icon-btn:hover { background: #eceef0; }

  .avatar {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #b4c5ff; color: #00174b;
    width: 2rem; height: 2rem; border-radius: 0.25rem;
    font-weight: 700; font-size: 0.65rem; cursor: pointer;
  }

  .canvas-page {
    padding: 1.5rem;
    max-width: 100rem;
  }

  .page-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 1.5rem;
    gap: 1rem;
  }

  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .dl-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.55rem 1rem;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
    border: 0;
    border-radius: 0.25rem;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    white-space: nowrap;
    font-family: inherit;
  }

  .dl-btn:hover { opacity: 0.9; }
  .dl-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .dl-btn .material-symbols-outlined { font-size: 1rem; }

  .generator-layout {
    display: grid;
    grid-template-columns: 380px 1fr;
    gap: 1.5rem;
    align-items: start;
  }

  @media (max-width: 900px) {
    .generator-layout { grid-template-columns: 1fr; }
  }

  /* Controls */
  .controls-panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .field-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .field-label {
    font-size: 0.68rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }

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

  .field-input:focus, .field-textarea:focus {
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgb(37 99 235 / 0.12);
  }

  .field-textarea { resize: vertical; min-height: 4.5rem; }

  .template-options {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.5rem;
  }

  .template-opt {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
    padding: 0.55rem 0.4rem;
    border: 1.5px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    background: #fff;
    cursor: pointer;
    text-align: center;
    font-family: inherit;
    transition: border-color 0.1s;
  }

  .template-opt:hover { border-color: #93c5fd; }
  .template-opt.selected { border-color: #2563eb; background: #eff6ff; }

  .template-name {
    font-size: 0.78rem;
    font-weight: 700;
    color: #191c1e;
  }

  .template-desc {
    font-size: 0.6rem;
    color: #434655;
    line-height: 1.3;
  }

  .pkg-summary {
    display: flex;
    gap: 0.6rem;
    padding: 0.65rem;
    background: #f0f9ff;
    border: 1px solid #bae6fd;
    border-radius: 0.25rem;
    font-size: 0.75rem;
    color: #0c4a6e;
  }

  .pkg-summary .material-symbols-outlined {
    font-size: 1rem;
    flex-shrink: 0;
    color: #0284c7;
    margin-top: 1px;
  }

  .pkg-summary strong { font-weight: 700; }
  .pkg-summary span { color: #0369a1; }

  .dl-btn-full {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.35rem;
    padding: 0.65rem;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
    border: 0;
    border-radius: 0.25rem;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    width: 100%;
  }

  .dl-btn-full:hover { opacity: 0.9; }
  .dl-btn-full:disabled { opacity: 0.5; cursor: not-allowed; }
  .dl-btn-full .material-symbols-outlined { font-size: 1rem; }

  /* Preview */
  .preview-wrap {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
  }

  .preview-label {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.75rem;
    color: #434655;
    font-weight: 600;
  }

  .preview-label .material-symbols-outlined { font-size: 1rem; }

  .canvas-frame {
    background: #e6e8ea;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    overflow: hidden;
    display: inline-block;
    line-height: 0;
  }

  .flyer-canvas {
    display: block;
    max-width: 100%;
    height: auto;
  }
</style>
