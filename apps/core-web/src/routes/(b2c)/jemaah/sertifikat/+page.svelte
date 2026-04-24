<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  const pilgrim = {
    name: 'Bambang Suryanto',
    tripName: 'Paket Umroh Plus Madinah — Januari 2025',
    departureDate: '15 Januari 2025',
    returnDate: '28 Januari 2025',
    umrohCount: 2,
    certNumber: 'UOS/UMR/2025/001234',
    issuedDate: '28 Januari 2025',
    agency: 'UmrohOS Travel — Izin Kemenag No. U.234/2024',
  };

  let downloading = $state(false);

  function handleDownload() {
    downloading = true;
    setTimeout(() => { downloading = false; }, 2000);
  }
</script>

<svelte:head>
  <title>Sertifikat Umroh — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="sertifikat-root">
    <div class="shell">
      <a href="/jemaah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Portal Jamaah
      </a>
      <div class="page-header">
        <h1>E-Sertifikat Umroh</h1>
        <p>Bukti resmi perjalanan ibadah umroh Anda</p>
      </div>

      <!-- Umroh counter badge -->
      <div class="counter-banner">
        <span class="material-symbols-outlined counter-icon">workspace_premium</span>
        <div>
          <div class="counter-title">Umroh ke-{pilgrim.umrohCount}</div>
          <div class="counter-sub">Anda telah melaksanakan {pilgrim.umrohCount} kali umroh bersama UmrohOS. Semoga mabrur.</div>
        </div>
      </div>

      <!-- Certificate preview -->
      <div class="cert-wrapper">
        <div class="cert-card">
          <div class="cert-header">
            <div class="cert-logo">UmrohOS</div>
            <div class="cert-subtitle">SERTIFIKAT UMROH</div>
          </div>

          <div class="cert-ornament top"></div>
          <div class="cert-ornament bottom"></div>

          <div class="cert-body">
            <p class="cert-pre">Dengan ini menerangkan bahwa</p>
            <h2 class="cert-name">{pilgrim.name}</h2>
            <p class="cert-info">telah melaksanakan ibadah Umroh dengan penuh keimanan dan ketaqwaan dalam</p>
            <div class="cert-trip">{pilgrim.tripName}</div>
            <div class="cert-dates">
              <div class="cert-date-item">
                <div class="cert-date-label">Keberangkatan</div>
                <div class="cert-date-val">{pilgrim.departureDate}</div>
              </div>
              <div class="cert-date-divider"></div>
              <div class="cert-date-item">
                <div class="cert-date-label">Kepulangan</div>
                <div class="cert-date-val">{pilgrim.returnDate}</div>
              </div>
            </div>
          </div>

          <div class="cert-footer">
            <div class="cert-number">No. Sertifikat: {pilgrim.certNumber}</div>
            <div class="cert-issued">Diterbitkan: {pilgrim.issuedDate}</div>
            <div class="cert-agency">{pilgrim.agency}</div>
          </div>
        </div>
      </div>

      <!-- Download button -->
      <div class="actions-row">
        <button class="download-btn" onclick={handleDownload} disabled={downloading}>
          {#if downloading}
            <span class="material-symbols-outlined spin">progress_activity</span>
            Menyiapkan PDF...
          {:else}
            <span class="material-symbols-outlined">download</span>
            Unduh Sertifikat PDF
          {/if}
        </button>
        <button class="share-btn">
          <span class="material-symbols-outlined">share</span>
          Bagikan
        </button>
      </div>

      <div class="info-box">
        <span class="material-symbols-outlined">info</span>
        <p>Sertifikat ini sah dan dapat digunakan sebagai bukti ibadah umroh resmi. Nomor sertifikat terdaftar di sistem UmrohOS.</p>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .sertifikat-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 64rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 1.5rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .counter-banner {
    display: flex;
    align-items: center;
    gap: 1rem;
    background: linear-gradient(135deg, rgba(119,90,25,0.1), rgba(254,212,136,0.15));
    border: 1px solid rgba(119,90,25,0.2);
    border-radius: 1.2rem;
    padding: 1.2rem 1.5rem;
    margin-bottom: 2rem;
  }
  .counter-icon { font-size: 2.5rem; color: #775a19; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 48; }
  .counter-title { font-size: 1.2rem; font-weight: 800; color: #775a19; font-family: 'Plus Jakarta Sans', sans-serif; }
  .counter-sub { font-size: 0.85rem; color: #57534e; margin-top: 0.2rem; }
  /* Certificate */
  .cert-wrapper { display: flex; justify-content: center; margin-bottom: 2rem; }
  .cert-card {
    position: relative;
    width: 100%;
    max-width: 600px;
    background: linear-gradient(135deg, #fffbf0, #fff);
    border: 2px solid rgba(119,90,25,0.25);
    border-radius: 1.5rem;
    padding: 2.5rem 3rem;
    text-align: center;
    box-shadow: 0 12px 40px rgba(0,0,0,0.08), inset 0 0 60px rgba(254,212,136,0.1);
    overflow: hidden;
  }
  .cert-ornament {
    position: absolute;
    left: 0; right: 0;
    height: 6px;
    background: linear-gradient(90deg, transparent, #775a19, #fed888, #775a19, transparent);
  }
  .cert-ornament.top { top: 0; }
  .cert-ornament.bottom { bottom: 0; }
  .cert-header { margin-bottom: 1.5rem; }
  .cert-logo { font-size: 1.6rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; letter-spacing: -0.02em; }
  .cert-subtitle { font-size: 0.75rem; font-weight: 700; letter-spacing: 0.2em; color: #775a19; text-transform: uppercase; margin-top: 0.25rem; }
  .cert-body { margin-bottom: 1.5rem; }
  .cert-pre { margin: 0 0 0.5rem; font-size: 0.88rem; color: #6b7280; }
  .cert-name { margin: 0.25rem 0 0.75rem; font-size: 2rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .cert-info { margin: 0 0 0.75rem; font-size: 0.88rem; color: #57534e; line-height: 1.6; }
  .cert-trip { font-size: 1rem; font-weight: 700; color: #775a19; margin-bottom: 1.25rem; font-family: 'Plus Jakarta Sans', sans-serif; }
  .cert-dates { display: flex; align-items: center; justify-content: center; gap: 1.5rem; }
  .cert-date-item { text-align: center; }
  .cert-date-label { font-size: 0.7rem; color: #9ca3af; text-transform: uppercase; letter-spacing: 0.07em; margin-bottom: 0.2rem; }
  .cert-date-val { font-size: 0.92rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .cert-date-divider { width: 2px; height: 2.5rem; background: rgba(119,90,25,0.2); border-radius: 999px; }
  .cert-footer { border-top: 1px dashed rgba(119,90,25,0.2); padding-top: 1rem; font-size: 0.72rem; color: #9ca3af; line-height: 1.8; }
  .cert-number { font-family: 'IBM Plex Mono', monospace; color: #775a19; font-weight: 600; }
  /* Actions */
  .actions-row { display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap; margin-bottom: 1.5rem; }
  .download-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: linear-gradient(135deg, #004d34, #006747);
    color: #fff;
    border: none;
    border-radius: 999px;
    padding: 0.85rem 2rem;
    font-size: 0.95rem;
    font-weight: 700;
    cursor: pointer;
    transition: opacity 0.15s;
    font-family: inherit;
  }
  .download-btn:disabled { opacity: 0.7; cursor: not-allowed; }
  .share-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: #fff;
    color: #006747;
    border: 1.5px solid #006747;
    border-radius: 999px;
    padding: 0.85rem 1.5rem;
    font-size: 0.95rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
  }
  .share-btn:hover { background: rgba(0,103,71,0.06); }
  .info-box {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
    background: rgba(0,103,71,0.06);
    border-radius: 1rem;
    padding: 1rem 1.4rem;
  }
  .info-box .material-symbols-outlined { color: #006747; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .info-box p { margin: 0; font-size: 0.85rem; color: #57534e; line-height: 1.6; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; }
</style>
