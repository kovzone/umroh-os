<script lang="ts">
  import { onMount } from 'svelte';

  interface PackageOption {
    id: string;
    name: string;
    price: number;
    departure_date: string;
    duration_days: number;
  }

  type ChannelId = 'whatsapp' | 'instagram' | 'twitter';

  interface Channel {
    id: ChannelId;
    label: string;
    icon: string;
    width: number;
    height: number;
    desc: string;
    color: string;
  }

  const MOCK_PACKAGES: PackageOption[] = [
    { id: 'pkg1', name: 'Umroh Ramadhan Premium', price: 32_500_000, departure_date: '2026-03-01', duration_days: 15 },
    { id: 'pkg2', name: 'Umroh Reguler April', price: 24_500_000, departure_date: '2026-04-10', duration_days: 12 },
    { id: 'pkg3', name: 'Umroh Plus Turki', price: 38_000_000, departure_date: '2026-05-05', duration_days: 18 },
    { id: 'pkg4', name: 'Umroh Hemat Juni', price: 19_900_000, departure_date: '2026-06-15', duration_days: 10 }
  ];

  const CHANNELS: Channel[] = [
    {
      id: 'whatsapp',
      label: 'WhatsApp Story',
      icon: 'chat',
      width: 1080,
      height: 1920,
      desc: '9:16 Portrait — Story WA',
      color: '#25D366'
    },
    {
      id: 'instagram',
      label: 'Instagram Post',
      icon: 'photo_camera',
      width: 1080,
      height: 1080,
      desc: '1:1 Square — Feed IG',
      color: '#E1306C'
    },
    {
      id: 'twitter',
      label: 'Twitter/X Banner',
      icon: 'alternate_email',
      width: 1500,
      height: 500,
      desc: '3:1 Landscape — X Banner',
      color: '#1DA1F2'
    }
  ];

  // State
  let packageId = $state(MOCK_PACKAGES[0].id);
  let headline = $state('Umroh Impianmu, Berangkat Bersama Kami');
  let contactInfo = $state('0812-3456-7890 | @umrohos');
  let activeChannel = $state<ChannelId>('whatsapp');

  // Canvas refs
  let canvases = $state<Record<ChannelId, HTMLCanvasElement | null>>({
    whatsapp: null,
    instagram: null,
    twitter: null
  });

  const selectedPackage = $derived(MOCK_PACKAGES.find(p => p.id === packageId) ?? MOCK_PACKAGES[0]);

  function formatIDR(n: number): string {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
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

  function wrapText(ctx: CanvasRenderingContext2D, text: string, x: number, y: number, maxWidth: number, lineHeight: number, align: CanvasTextAlign = 'left') {
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
    return currentY;
  }

  function drawWhatsApp(canvas: HTMLCanvasElement, pkg: PackageOption) {
    // Display scale: 270x480 (quarter of 1080x1920)
    canvas.width = 270;
    canvas.height = 480;
    const ctx = canvas.getContext('2d')!;
    const W = 270, H = 480;

    const grad = ctx.createLinearGradient(0, 0, 0, H);
    grad.addColorStop(0, '#064e3b');
    grad.addColorStop(1, '#065f46');
    ctx.fillStyle = grad;
    ctx.fillRect(0, 0, W, H);

    // WA green accent top
    ctx.fillStyle = '#25D366';
    ctx.fillRect(0, 0, W, 4);

    // Logo
    ctx.beginPath();
    ctx.arc(W / 2, 52, 26, 0, Math.PI * 2);
    ctx.fillStyle = '#25D366';
    ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 16px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('O', W / 2, 52);

    ctx.fillStyle = '#6ee7b7';
    ctx.font = '600 10px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText('UmrohOS Travel', W / 2, 92);

    // Headline
    ctx.fillStyle = '#ffffff';
    ctx.font = 'bold 18px sans-serif';
    wrapText(ctx, headline, W / 2, 118, W - 30, 22, 'center');

    // Package band
    ctx.fillStyle = 'rgba(255,255,255,0.12)';
    roundRect(ctx, 12, 168, W - 24, 38, 6);
    ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 11px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(pkg.name, W / 2, 186);

    // Info
    ctx.fillStyle = '#6ee7b7';
    ctx.font = '10px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText(`${formatDate(pkg.departure_date)} · ${pkg.duration_days} hari`, W / 2, 222);

    // Price
    ctx.fillStyle = '#25D366';
    roundRect(ctx, 30, 236, W - 60, 44, 8);
    ctx.fill();
    ctx.fillStyle = '#064e3b';
    ctx.font = '9px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText('Mulai dari', W / 2, 252);
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 16px sans-serif';
    ctx.fillText(formatIDR(pkg.price), W / 2, 271);

    // CTA
    ctx.fillStyle = '#fff';
    roundRect(ctx, W / 2 - 60, 340, 120, 30, 15);
    ctx.fill();
    ctx.fillStyle = '#064e3b';
    ctx.font = 'bold 11px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('DAFTAR SEKARANG', W / 2, 355);

    // Contact
    ctx.fillStyle = '#6ee7b7';
    ctx.font = '9px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText(contactInfo, W / 2, 408);

    ctx.fillStyle = '#25D366';
    ctx.fillRect(0, H - 4, W, 4);
  }

  function drawInstagram(canvas: HTMLCanvasElement, pkg: PackageOption) {
    canvas.width = 320;
    canvas.height = 320;
    const ctx = canvas.getContext('2d')!;
    const W = 320, H = 320;

    const grad = ctx.createLinearGradient(0, 0, W, H);
    grad.addColorStop(0, '#7c3aed');
    grad.addColorStop(0.5, '#db2777');
    grad.addColorStop(1, '#ea580c');
    ctx.fillStyle = grad;
    ctx.fillRect(0, 0, W, H);

    ctx.fillStyle = 'rgba(255,255,255,0.08)';
    ctx.beginPath();
    ctx.arc(-20, H + 20, 200, 0, Math.PI * 2);
    ctx.fill();
    ctx.beginPath();
    ctx.arc(W + 20, -20, 180, 0, Math.PI * 2);
    ctx.fill();

    // Logo
    ctx.beginPath();
    ctx.arc(W / 2, 55, 28, 0, Math.PI * 2);
    ctx.fillStyle = 'rgba(255,255,255,0.25)';
    ctx.fill();
    ctx.strokeStyle = '#fff';
    ctx.lineWidth = 2;
    ctx.stroke();
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 18px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('O', W / 2, 55);

    ctx.fillStyle = 'rgba(255,255,255,0.85)';
    ctx.font = '500 11px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText('UmrohOS Travel', W / 2, 98);

    ctx.fillStyle = '#fff';
    ctx.font = 'bold 22px sans-serif';
    wrapText(ctx, headline, W / 2, 128, W - 40, 28, 'center');

    ctx.fillStyle = 'rgba(255,255,255,0.15)';
    roundRect(ctx, 20, 180, W - 40, 40, 6);
    ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 12px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(pkg.name, W / 2, 200);

    ctx.fillStyle = 'rgba(255,255,255,0.7)';
    ctx.font = '10px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText(`${formatDate(pkg.departure_date)} · ${pkg.duration_days} hari`, W / 2, 234);

    ctx.fillStyle = '#fff';
    ctx.font = 'bold 20px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText(formatIDR(pkg.price), W / 2, 262);

    ctx.fillStyle = '#fff';
    roundRect(ctx, W / 2 - 65, 276, 130, 26, 13);
    ctx.fill();
    const ctagrad = ctx.createLinearGradient(W / 2 - 65, 0, W / 2 + 65, 0);
    ctagrad.addColorStop(0, '#7c3aed');
    ctagrad.addColorStop(1, '#db2777');
    ctx.fillStyle = ctagrad;
    ctx.font = 'bold 11px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('DAFTAR SEKARANG', W / 2, 289);

    ctx.fillStyle = 'rgba(255,255,255,0.6)';
    ctx.font = '9px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText(contactInfo, W / 2, 312);
  }

  function drawTwitter(canvas: HTMLCanvasElement, pkg: PackageOption) {
    canvas.width = 450;
    canvas.height = 150;
    const ctx = canvas.getContext('2d')!;
    const W = 450, H = 150;

    const grad = ctx.createLinearGradient(0, 0, W, 0);
    grad.addColorStop(0, '#0f172a');
    grad.addColorStop(1, '#1e293b');
    ctx.fillStyle = grad;
    ctx.fillRect(0, 0, W, H);

    ctx.fillStyle = '#1DA1F2';
    ctx.fillRect(0, 0, W, 3);

    // Logo
    ctx.beginPath();
    ctx.arc(60, H / 2, 36, 0, Math.PI * 2);
    ctx.fillStyle = '#1DA1F2';
    ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 22px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('O', 60, H / 2);

    // Text content
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 18px sans-serif';
    ctx.textAlign = 'left';
    ctx.textBaseline = 'alphabetic';
    wrapText(ctx, headline, 118, 46, 200, 22, 'left');

    ctx.fillStyle = '#94a3b8';
    ctx.font = '11px sans-serif';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText(`${pkg.name} · ${formatDate(pkg.departure_date)}`, 118, 88);

    ctx.fillStyle = '#64748b';
    ctx.font = '10px sans-serif';
    ctx.fillText(contactInfo, 118, 108);

    // Price block right side
    ctx.fillStyle = '#1DA1F2';
    roundRect(ctx, W - 130, 30, 110, 90, 10);
    ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.font = '10px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'alphabetic';
    ctx.fillText('Mulai dari', W - 75, 58);
    ctx.font = 'bold 14px sans-serif';
    const price = formatIDR(pkg.price);
    // shorten if too long
    ctx.fillText(price.length > 14 ? price.substring(0, 14) + '..' : price, W - 75, 80);
    ctx.fillStyle = '#0f172a';
    roundRect(ctx, W - 125, 92, 100, 20, 10);
    ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 9px sans-serif';
    ctx.fillText('DAFTAR', W - 75, 105);

    ctx.fillStyle = '#1DA1F2';
    ctx.fillRect(0, H - 3, W, 3);
  }

  function drawAll() {
    const pkg = selectedPackage;
    if (canvases.whatsapp) drawWhatsApp(canvases.whatsapp, pkg);
    if (canvases.instagram) drawInstagram(canvases.instagram, pkg);
    if (canvases.twitter) drawTwitter(canvases.twitter, pkg);
  }

  $effect(() => {
    void packageId;
    void headline;
    void contactInfo;
    drawAll();
  });

  onMount(() => drawAll());

  function downloadChannel(channelId: ChannelId) {
    const canvas = canvases[channelId];
    if (!canvas) return;
    const ch = CHANNELS.find(c => c.id === channelId)!;
    const link = document.createElement('a');
    link.download = `omni-flyer-${channelId}-${packageId}.png`;
    link.href = canvas.toDataURL('image/png');
    link.click();
  }

  async function copyToClipboard(channelId: ChannelId) {
    const canvas = canvases[channelId];
    if (!canvas) return;
    canvas.toBlob(async (blob) => {
      if (!blob) return;
      try {
        await navigator.clipboard.write([
          new ClipboardItem({ 'image/png': blob })
        ]);
        alert('Gambar berhasil disalin ke clipboard!');
      } catch {
        alert('Clipboard API tidak didukung browser ini. Gunakan tombol Download.');
      }
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
      <span class="topbar-current">Omni-Flyer</span>
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
        <h2>Omni-Flyer</h2>
        <p>BL-CRM-014 — Generate flyer per channel: WhatsApp Story, Instagram Post, Twitter/X Banner</p>
      </div>
    </div>

    <!-- Controls -->
    <div class="controls-row">
      <div class="field-group">
        <label class="field-label" for="pkg-select">Paket</label>
        <select id="pkg-select" class="field-input" bind:value={packageId}>
          {#each MOCK_PACKAGES as p}
            <option value={p.id}>{p.name}</option>
          {/each}
        </select>
      </div>
      <div class="field-group flex-1">
        <label class="field-label" for="headline-input">Headline</label>
        <input id="headline-input" type="text" class="field-input" bind:value={headline} placeholder="Headline flyer" />
      </div>
      <div class="field-group">
        <label class="field-label" for="contact-input">Kontak</label>
        <input id="contact-input" type="text" class="field-input" bind:value={contactInfo} placeholder="HP / sosmed" />
      </div>
    </div>

    <!-- Channel tabs -->
    <div class="channel-tabs">
      {#each CHANNELS as ch}
        <button
          type="button"
          class="ch-tab"
          class:active={activeChannel === ch.id}
          onclick={() => { activeChannel = ch.id; }}
          style="--ch-color: {ch.color}"
        >
          <span class="material-symbols-outlined">{ch.icon}</span>
          <span class="ch-label">{ch.label}</span>
          <span class="ch-dim">{ch.width}×{ch.height}</span>
        </button>
      {/each}
    </div>

    <!-- Preview + actions for each channel -->
    <div class="channels-grid">
      {#each CHANNELS as ch}
        <div class="channel-card" class:hidden={activeChannel !== ch.id}>
          <div class="ch-card-header" style="--ch-color: {ch.color}">
            <span class="material-symbols-outlined">{ch.icon}</span>
            <strong>{ch.label}</strong>
            <span class="ch-desc">{ch.desc}</span>
            <div class="ch-actions">
              <button class="ch-btn" onclick={() => copyToClipboard(ch.id)} title="Copy ke clipboard">
                <span class="material-symbols-outlined">content_copy</span>
                Copy
              </button>
              <button class="ch-btn primary" onclick={() => downloadChannel(ch.id)}>
                <span class="material-symbols-outlined">download</span>
                Download
              </button>
            </div>
          </div>
          <div class="preview-area" class:landscape={ch.id === 'twitter'}>
            {#if ch.id === 'whatsapp'}
              <canvas bind:this={canvases.whatsapp} class="channel-canvas portrait"></canvas>
            {:else if ch.id === 'instagram'}
              <canvas bind:this={canvases.instagram} class="channel-canvas square"></canvas>
            {:else}
              <canvas bind:this={canvases.twitter} class="channel-canvas wide"></canvas>
            {/if}
          </div>
        </div>
      {/each}
    </div>

    <!-- All channels download bar -->
    <div class="all-dl-bar">
      <span class="material-symbols-outlined">photo_library</span>
      <span>Download semua format sekaligus:</span>
      {#each CHANNELS as ch}
        <button class="mini-dl-btn" onclick={() => downloadChannel(ch.id)} style="--ch-color: {ch.color}">
          <span class="material-symbols-outlined">{ch.icon}</span>
          {ch.label}
        </button>
      {/each}
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

  .canvas-page { padding: 1.5rem; max-width: 80rem; }

  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .controls-row {
    display: flex;
    gap: 1rem;
    align-items: flex-end;
    flex-wrap: wrap;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1rem;
    margin-bottom: 1.25rem;
  }

  .field-group { display: flex; flex-direction: column; gap: 0.3rem; }
  .flex-1 { flex: 1; min-width: 14rem; }

  .field-label {
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }

  .field-input {
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    padding: 0.48rem 0.65rem;
    font-size: 0.82rem;
    color: #191c1e;
    background: #fff;
    font-family: inherit;
    outline: none;
    min-width: 12rem;
  }

  .field-input:focus { border-color: #2563eb; }

  .channel-tabs {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
  }

  .ch-tab {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.5rem 0.85rem;
    border: 1.5px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    background: #fff;
    font-family: inherit;
    cursor: pointer;
    transition: border-color 0.1s, background 0.1s;
  }

  .ch-tab:hover { background: #f2f4f6; }
  .ch-tab.active { border-color: var(--ch-color); background: color-mix(in srgb, var(--ch-color) 8%, white); }
  .ch-tab .material-symbols-outlined { font-size: 1rem; color: var(--ch-color); }
  .ch-label { font-size: 0.82rem; font-weight: 600; color: #191c1e; }
  .ch-dim { font-size: 0.65rem; color: #737686; }

  .channels-grid { display: block; }

  .channel-card { display: flex; flex-direction: column; gap: 0.75rem; }
  .channel-card.hidden { display: none; }

  .ch-card-header {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.85rem 1rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
  }

  .ch-card-header .material-symbols-outlined { font-size: 1.1rem; color: var(--ch-color); }
  .ch-card-header strong { font-size: 0.9rem; color: #191c1e; }
  .ch-desc { font-size: 0.75rem; color: #737686; }

  .ch-actions { margin-left: auto; display: flex; gap: 0.5rem; }

  .ch-btn {
    display: inline-flex; align-items: center; gap: 0.3rem;
    padding: 0.4rem 0.75rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    background: #fff;
    font-size: 0.75rem;
    font-weight: 600;
    color: #191c1e;
    cursor: pointer;
    font-family: inherit;
  }

  .ch-btn:hover { background: #f2f4f6; }
  .ch-btn.primary {
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
    border-color: transparent;
  }
  .ch-btn.primary:hover { opacity: 0.9; }
  .ch-btn .material-symbols-outlined { font-size: 0.9rem; }

  .preview-area {
    display: flex;
    justify-content: flex-start;
    padding: 1.5rem;
    background: #e6e8ea;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
  }

  .preview-area.landscape { justify-content: center; }

  .channel-canvas { display: block; border-radius: 4px; }
  .channel-canvas.portrait { max-height: 480px; width: auto; }
  .channel-canvas.square { max-width: 320px; max-height: 320px; }
  .channel-canvas.wide { max-width: 100%; height: auto; }

  .all-dl-bar {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-top: 1.25rem;
    padding: 0.85rem 1rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    font-size: 0.82rem;
    color: #434655;
  }

  .all-dl-bar .material-symbols-outlined { font-size: 1.05rem; color: #2563eb; }

  .mini-dl-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.35rem 0.65rem;
    border: 1.5px solid var(--ch-color);
    border-radius: 0.2rem;
    background: color-mix(in srgb, var(--ch-color) 8%, white);
    color: var(--ch-color);
    font-size: 0.72rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
  }

  .mini-dl-btn:hover { opacity: 0.85; }
  .mini-dl-btn .material-symbols-outlined { font-size: 0.85rem; }
</style>
