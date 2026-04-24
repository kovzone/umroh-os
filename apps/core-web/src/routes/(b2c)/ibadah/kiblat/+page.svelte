<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  // Mock: Makkah direction from Jakarta ≈ 295° (Northwest)
  let compassDegrees = $state(295);
  let animating = $state(false);

  function calibrate() {
    animating = true;
    // Simulate compass rotation animation
    const target = 295 + (Math.random() - 0.5) * 10;
    compassDegrees = Math.round(target);
    setTimeout(() => { animating = false; }, 1200);
  }

  const directionLabel = $derived(() => {
    const d = ((compassDegrees % 360) + 360) % 360;
    if (d >= 337.5 || d < 22.5) return 'Utara';
    if (d < 67.5) return 'Timur Laut';
    if (d < 112.5) return 'Timur';
    if (d < 157.5) return 'Tenggara';
    if (d < 202.5) return 'Selatan';
    if (d < 247.5) return 'Barat Daya';
    if (d < 292.5) return 'Barat';
    return 'Barat Laut';
  });
</script>

<svelte:head>
  <title>Kiblat — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="kiblat-root">
    <div class="shell">
      <a href="/ibadah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Panduan Ibadah
      </a>
      <div class="page-header">
        <h1>Kompas Kiblat</h1>
        <p>Arahkan perangkat ke utara untuk kalibrasi</p>
      </div>

      <!-- Compass -->
      <div class="compass-section">
        <div class="compass-card">
          <div class="compass-outer">
            <!-- Cardinal directions (fixed) -->
            <div class="cardinal n">U</div>
            <div class="cardinal e">T</div>
            <div class="cardinal s">S</div>
            <div class="cardinal w">B</div>

            <!-- Compass rose (rotates based on degrees) -->
            <div class="compass-inner">
              <!-- Tick marks -->
              {#each Array(36) as _, i}
                <div class="tick" style="transform: rotate({i * 10}deg) translateX(-50%)">
                  <div class="tick-line" class:major={i % 9 === 0}></div>
                </div>
              {/each}
            </div>

            <!-- Qibla arrow (points toward Makkah) -->
            <div
              class="qibla-arrow-wrap"
              class:animate={animating}
              style="transform: rotate({compassDegrees}deg)"
            >
              <div class="qibla-arrow">
                <div class="arrow-head"></div>
                <div class="arrow-body"></div>
                <div class="arrow-tail"></div>
              </div>
            </div>

            <!-- Center dot -->
            <div class="compass-center">
              <span class="material-symbols-outlined">mosque</span>
            </div>
          </div>

          <div class="compass-info">
            <div class="deg-display">
              <span class="deg-value">{compassDegrees}°</span>
              <span class="deg-dir">{directionLabel()}</span>
            </div>
            <div class="compass-label">Arah Kiblat dari Lokasi Anda</div>
          </div>
        </div>

        <div class="calibrate-hint">
          <span class="material-symbols-outlined">info</span>
          <p>Arahkan perangkat ke utara untuk kalibrasi. Jauhkan dari benda logam atau elektromagnetik.</p>
        </div>

        <button class="calibrate-btn" onclick={calibrate} disabled={animating}>
          <span class="material-symbols-outlined" class:spin={animating}>my_location</span>
          {animating ? 'Mengkalibrasi...' : 'Kalibrasi Ulang'}
        </button>
      </div>

      <!-- Info -->
      <div class="makkah-info">
        <div class="makkah-icon">
          <span class="material-symbols-outlined">mosque</span>
        </div>
        <div>
          <div class="makkah-title">Ka'bah — Masjidil Haram</div>
          <div class="makkah-coords">21.4225° N, 39.8262° E · Makkah Al-Mukarramah, Saudi Arabia</div>
        </div>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .kiblat-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 56rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .compass-section { display: flex; flex-direction: column; align-items: center; gap: 1.5rem; margin-bottom: 2rem; }
  .compass-card {
    background: #fff;
    border-radius: 2rem;
    padding: 2.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    width: 100%;
    max-width: 380px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.06);
  }
  .compass-outer {
    position: relative;
    width: 240px;
    height: 240px;
    border-radius: 50%;
    background: radial-gradient(circle at center, #f0f9f4 0%, #fbf9f8 100%);
    border: 2px solid rgba(0,103,71,0.15);
    box-shadow: inset 0 2px 8px rgba(0,0,0,0.05), 0 4px 16px rgba(0,0,0,0.08);
  }
  .cardinal {
    position: absolute;
    font-size: 0.75rem;
    font-weight: 800;
    color: #6b7280;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .cardinal.n { top: 8px; left: 50%; transform: translateX(-50%); color: #006747; }
  .cardinal.e { right: 8px; top: 50%; transform: translateY(-50%); }
  .cardinal.s { bottom: 8px; left: 50%; transform: translateX(-50%); }
  .cardinal.w { left: 8px; top: 50%; transform: translateY(-50%); }
  .compass-inner { position: absolute; inset: 0; }
  .tick {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    transform-origin: 0 0;
  }
  .tick-line {
    position: absolute;
    width: 1px;
    height: 8px;
    background: rgba(190,201,193,0.5);
    top: -118px;
    left: 0;
  }
  .tick-line.major { height: 14px; background: rgba(0,103,71,0.3); width: 2px; top: -118px; }
  /* Qibla arrow */
  .qibla-arrow-wrap {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: transform 1.2s cubic-bezier(0.22, 0.61, 0.36, 1);
  }
  .qibla-arrow-wrap.animate { transition: transform 1.2s cubic-bezier(0.22, 0.61, 0.36, 1); }
  .qibla-arrow { display: flex; flex-direction: column; align-items: center; }
  .arrow-head {
    width: 0;
    height: 0;
    border-left: 10px solid transparent;
    border-right: 10px solid transparent;
    border-bottom: 20px solid #006747;
  }
  .arrow-body {
    width: 6px;
    height: 60px;
    background: linear-gradient(to bottom, #006747, rgba(0,103,71,0.3));
    border-radius: 0 0 3px 3px;
  }
  .arrow-tail {
    width: 0;
    height: 0;
    border-left: 8px solid transparent;
    border-right: 8px solid transparent;
    border-top: 14px solid rgba(0,103,71,0.3);
  }
  .compass-center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    background: #fff;
    border: 2px solid rgba(0,103,71,0.2);
    display: grid;
    place-items: center;
    color: #006747;
    z-index: 2;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }
  .compass-center .material-symbols-outlined { font-size: 1.1rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .compass-info { text-align: center; }
  .deg-display { display: flex; align-items: baseline; gap: 0.5rem; justify-content: center; margin-bottom: 0.25rem; }
  .deg-value { font-size: 2.5rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .deg-dir { font-size: 1rem; color: #006747; font-weight: 600; }
  .compass-label { font-size: 0.8rem; color: #9ca3af; }
  .calibrate-hint { display: flex; align-items: flex-start; gap: 0.5rem; background: rgba(0,103,71,0.06); border-radius: 1rem; padding: 0.9rem 1.2rem; max-width: 380px; width: 100%; }
  .calibrate-hint .material-symbols-outlined { color: #006747; font-size: 1rem; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .calibrate-hint p { margin: 0; font-size: 0.82rem; color: #57534e; line-height: 1.5; }
  .calibrate-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: #006747;
    color: #fff;
    border: none;
    border-radius: 999px;
    padding: 0.8rem 2rem;
    font-size: 0.9rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
  }
  .calibrate-btn:disabled { opacity: 0.7; cursor: not-allowed; }
  .makkah-info {
    display: flex;
    align-items: center;
    gap: 1rem;
    background: #fff;
    border: 1px solid rgba(190,201,193,0.2);
    border-radius: 1.2rem;
    padding: 1.2rem 1.5rem;
  }
  .makkah-icon { width: 3rem; height: 3rem; border-radius: 0.85rem; background: rgba(0,103,71,0.08); display: grid; place-items: center; color: #006747; flex-shrink: 0; }
  .makkah-icon .material-symbols-outlined { font-size: 1.5rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .makkah-title { font-size: 0.92rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .makkah-coords { font-size: 0.78rem; color: #9ca3af; margin-top: 0.1rem; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; }
</style>
