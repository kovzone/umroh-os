<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  // Savings simulation parameters
  let targetPackage = $state<'reguler' | 'gold' | 'platinum'>('gold');
  let monthlySavings = $state(1500000);
  let currentSavings = $state(0);
  let additionalPax = $state(0); // how many extra people traveling

  const packagePrices: Record<string, { label: string; price: number; desc: string }> = {
    reguler: { label: 'Paket Reguler', price: 28500000, desc: 'Hotel bintang 3, maskapai ekonomi' },
    gold: { label: 'Paket Gold', price: 34200000, desc: 'Hotel bintang 4, maskapai premium' },
    platinum: { label: 'Paket Platinum', price: 48900000, desc: 'Hotel bintang 5, business class' }
  };

  const selectedPkg = $derived(packagePrices[targetPackage]);
  const totalTarget = $derived(selectedPkg.price * (1 + additionalPax));
  const remaining = $derived(Math.max(0, totalTarget - currentSavings));
  const monthsNeeded = $derived(
    monthlySavings > 0 ? Math.ceil(remaining / monthlySavings) : null
  );
  const yearsNeeded = $derived(monthsNeeded != null ? (monthsNeeded / 12).toFixed(1) : null);
  const targetDate = $derived.by(() => {
    if (!monthsNeeded) return null;
    const d = new Date();
    d.setMonth(d.getMonth() + monthsNeeded);
    return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' });
  });
  const progressPercent = $derived(
    totalTarget > 0 ? Math.min(100, (currentSavings / totalTarget) * 100) : 0
  );

  function formatRp(val: number): string {
    return 'Rp ' + new Intl.NumberFormat('id-ID').format(Math.round(val));
  }

  function handleSlider(e: Event, field: 'monthly' | 'current') {
    const v = parseInt((e.target as HTMLInputElement).value);
    if (field === 'monthly') monthlySavings = v;
    else currentSavings = v;
  }

  const savingsMilestones = $derived.by(() => {
    if (!monthsNeeded) return [];
    return [3, 6, 12, 24].map(m => ({
      label: m < 12 ? `${m} bulan` : `${m/12} tahun`,
      amount: currentSavings + monthlySavings * m,
      percent: Math.min(100, ((currentSavings + monthlySavings * m) / totalTarget) * 100)
    }));
  });
</script>

<svelte:head>
  <title>Kalkulator Tabungan Umrah — UmrohOS</title>
  <meta name="description" content="Hitung berapa lama Anda perlu menabung untuk mewujudkan impian umrah. Simulasi tabungan umrah gratis dari UmrohOS." />
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="calc-root">

    <section class="calc-hero">
      <div class="shell">
        <p class="kicker">Kalkulator Tabungan</p>
        <h1>Rencanakan Perjalanan Umrah Impian Anda</h1>
        <p class="hero-sub">Simulasikan berapa lama Anda perlu menabung dan kapan bisa berangkat. Gratis, tanpa daftar.</p>
        <p class="disclaimer">* Kalkulasi ini bersifat estimasi dan tidak termasuk fluktuasi harga. Harga final mengikuti paket aktual saat pemesanan.</p>
      </div>
    </section>

    <section class="calc-main">
      <div class="shell calc-grid">

        <!-- Left: Controls -->
        <div class="controls">
          <div class="control-group">
            <label class="group-title">1. Pilih Paket Umrah</label>
            <div class="package-options">
              {#each Object.entries(packagePrices) as [key, pkg] (key)}
                <label class="pkg-option" class:selected={targetPackage === key}>
                  <input type="radio" name="package" value={key} bind:group={targetPackage} />
                  <div class="pkg-option-content">
                    <span class="pkg-name">{pkg.label}</span>
                    <span class="pkg-price">{formatRp(pkg.price)}</span>
                    <span class="pkg-desc">{pkg.desc}</span>
                  </div>
                  {#if targetPackage === key}
                    <span class="check material-symbols-outlined">check_circle</span>
                  {/if}
                </label>
              {/each}
            </div>
          </div>

          <div class="control-group">
            <label class="group-title">2. Jumlah Jamaah Tambahan</label>
            <p class="group-hint">Berapa anggota keluarga lain yang akan ikut? (di luar Anda sendiri)</p>
            <div class="number-input">
              <button type="button" onclick={() => additionalPax = Math.max(0, additionalPax - 1)}>−</button>
              <span>{additionalPax}</span>
              <button type="button" onclick={() => additionalPax = Math.min(10, additionalPax + 1)}>+</button>
            </div>
          </div>

          <div class="control-group">
            <div class="slider-header">
              <label class="group-title">3. Tabungan Saat Ini</label>
              <span class="slider-value">{formatRp(currentSavings)}</span>
            </div>
            <input
              type="range" min="0" max={totalTarget} step="500000"
              value={currentSavings}
              oninput={(e) => handleSlider(e, 'current')}
              class="slider"
            />
            <div class="slider-labels">
              <span>Rp 0</span>
              <span>{formatRp(totalTarget)}</span>
            </div>
          </div>

          <div class="control-group">
            <div class="slider-header">
              <label class="group-title">4. Tabungan per Bulan</label>
              <span class="slider-value">{formatRp(monthlySavings)}</span>
            </div>
            <input
              type="range" min="100000" max="10000000" step="100000"
              value={monthlySavings}
              oninput={(e) => handleSlider(e, 'monthly')}
              class="slider"
            />
            <div class="slider-labels">
              <span>Rp 100.000</span>
              <span>Rp 10.000.000</span>
            </div>
          </div>
        </div>

        <!-- Right: Results -->
        <div class="results">
          <div class="result-card main-result">
            <div class="result-icon">
              <span class="material-symbols-outlined">mosque</span>
            </div>
            {#if monthsNeeded && monthsNeeded <= 0}
              <div class="ready-now">
                <h2>Anda Sudah Siap! 🎉</h2>
                <p>Tabungan Anda sudah mencukupi. Segera booking paket umrah Anda!</p>
                <a href="/packages" class="btn-go">Pilih Paket Sekarang</a>
              </div>
            {:else if monthsNeeded}
              <p class="result-label">Estimasi waktu menabung</p>
              <h2 class="result-main">
                {monthsNeeded < 12
                  ? `${monthsNeeded} bulan`
                  : yearsNeeded + ' tahun'}
              </h2>
              {#if targetDate}
                <p class="result-date">
                  <span class="material-symbols-outlined">calendar_month</span>
                  Siap berangkat sekitar <strong>{targetDate}</strong>
                </p>
              {/if}
            {:else}
              <p class="result-label">Masukkan nilai tabungan bulanan</p>
            {/if}
          </div>

          <!-- Progress -->
          <div class="result-card">
            <h3>Progres Tabungan</h3>
            <div class="progress-info">
              <span>Terkumpul: {formatRp(currentSavings)}</span>
              <span>{progressPercent.toFixed(1)}%</span>
            </div>
            <div class="progress-bar-wrap">
              <div class="progress-bar" style="width: {progressPercent}%"></div>
            </div>
            <div class="progress-info">
              <span>Target: {formatRp(totalTarget)}</span>
              <span>Kurang: {formatRp(remaining)}</span>
            </div>
          </div>

          <!-- Milestones -->
          {#if savingsMilestones.length > 0}
          <div class="result-card">
            <h3>Estimasi Pencapaian</h3>
            <div class="milestones-list">
              {#each savingsMilestones as ms (ms.label)}
                <div class="milestone-row">
                  <span class="ms-label">{ms.label}</span>
                  <div class="ms-bar-wrap">
                    <div class="ms-bar" style="width: {ms.percent}%"></div>
                  </div>
                  <span class="ms-amount">{formatRp(ms.amount)}</span>
                  <span class="ms-pct">{ms.percent.toFixed(0)}%</span>
                </div>
              {/each}
            </div>
          </div>
          {/if}

          <div class="result-cta">
            <p>Ingin konsultasi lebih lanjut tentang paket dan cicilan?</p>
            <div class="result-cta-btns">
              <a href="https://wa.me/6281200000000" class="btn-wa" target="_blank" rel="noreferrer">
                <span class="material-symbols-outlined">chat</span>
                WhatsApp Kami
              </a>
              <a href="/packages" class="btn-packages">Lihat Paket</a>
            </div>
          </div>
        </div>

      </div>
    </section>

  </div>
</MarketingPageLayout>

<style>
  .calc-root {
    padding-top: 5.2rem;
    background: #fbf9f8;
  }
  .shell {
    max-width: 80rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  /* Hero */
  .calc-hero {
    padding: 4rem 0 3rem;
    background: linear-gradient(180deg, #f0f9f4 0%, #fbf9f8 100%);
    text-align: center;
  }
  .kicker {
    display: inline-block;
    margin: 0 0 1rem;
    background: rgba(254,212,136,0.3);
    color: #775a19;
    border-radius: 999px;
    padding: 0.4rem 1rem;
    font-size: 0.76rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .calc-hero h1 {
    margin: 0 0 0.7rem;
    font-size: clamp(1.8rem, 3.5vw, 3rem);
    font-weight: 800;
    color: #004d34;
    letter-spacing: -0.02em;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .hero-sub {
    margin: 0 auto 0.8rem;
    max-width: 40rem;
    color: #57534e;
    font-size: 1.05rem;
  }
  .disclaimer {
    margin: 0 auto;
    max-width: 40rem;
    color: #9ca3af;
    font-size: 0.8rem;
    font-style: italic;
  }
  /* Calculator grid */
  .calc-main {
    padding: 3rem 0 6rem;
  }
  .calc-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2.5rem;
    align-items: start;
  }
  /* Controls */
  .controls {
    display: grid;
    gap: 2rem;
  }
  .control-group {
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.6rem;
    border: 1px solid rgba(190,201,193,0.3);
  }
  .group-title {
    display: block;
    font-size: 0.88rem;
    font-weight: 700;
    color: #1b1c1c;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    margin-bottom: 1rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .group-hint {
    margin: -0.5rem 0 0.8rem;
    font-size: 0.82rem;
    color: #9ca3af;
  }
  .package-options {
    display: grid;
    gap: 0.7rem;
  }
  .pkg-option {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.9rem 1rem;
    border-radius: 1rem;
    border: 1.5px solid rgba(190,201,193,0.3);
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
  }
  .pkg-option input[type=radio] { display: none; }
  .pkg-option.selected {
    border-color: #006747;
    background: rgba(0,103,71,0.04);
  }
  .pkg-option-content {
    flex: 1;
  }
  .pkg-name {
    display: block;
    font-weight: 700;
    font-size: 0.9rem;
    color: #1b1c1c;
  }
  .pkg-price {
    display: block;
    font-size: 1.05rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .pkg-desc {
    display: block;
    font-size: 0.78rem;
    color: #9ca3af;
  }
  .check {
    color: #006747;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .number-input {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }
  .number-input button {
    width: 2.8rem;
    height: 2.8rem;
    border-radius: 50%;
    border: 1.5px solid rgba(190,201,193,0.5);
    background: #fbf9f8;
    font-size: 1.4rem;
    font-weight: 700;
    color: #004d34;
    cursor: pointer;
    display: grid;
    place-items: center;
    transition: background 0.15s;
  }
  .number-input button:hover { background: #f0f9f4; }
  .number-input span {
    font-size: 1.6rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
    min-width: 2rem;
    text-align: center;
  }
  .slider-header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 0.8rem;
  }
  .slider-value {
    font-size: 1.05rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .slider {
    width: 100%;
    accent-color: #006747;
    height: 6px;
    cursor: pointer;
  }
  .slider-labels {
    display: flex;
    justify-content: space-between;
    margin-top: 0.4rem;
    font-size: 0.76rem;
    color: #9ca3af;
  }
  /* Results */
  .results {
    display: grid;
    gap: 1.5rem;
    position: sticky;
    top: 7rem;
  }
  .result-card {
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.6rem;
    border: 1px solid rgba(190,201,193,0.3);
  }
  .result-card.main-result {
    text-align: center;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
    color: #fff;
    border: none;
  }
  .result-icon .material-symbols-outlined {
    font-size: 3rem;
    color: rgba(255,255,255,0.6);
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
    margin-bottom: 0.5rem;
    display: block;
  }
  .result-label {
    margin: 0 0 0.4rem;
    font-size: 0.85rem;
    color: rgba(255,255,255,0.7);
  }
  .result-main {
    margin: 0;
    font-size: clamp(2rem, 4vw, 3rem);
    font-weight: 800;
    font-family: 'Plus Jakarta Sans', sans-serif;
    color: #fff;
  }
  .result-date {
    margin: 0.8rem 0 0;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.9rem;
    color: rgba(255,255,255,0.8);
  }
  .result-date .material-symbols-outlined {
    font-size: 1.1rem;
  }
  .ready-now { color: #fff; }
  .ready-now h2 { margin: 0; font-size: 1.6rem; font-family: 'Plus Jakarta Sans', sans-serif; }
  .ready-now p { margin: 0.5rem 0 1.2rem; opacity: 0.85; }
  .btn-go {
    display: inline-flex;
    text-decoration: none;
    background: #fed488;
    color: #5d4201;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.7rem 1.6rem;
  }
  .result-card h3 {
    margin: 0 0 1rem;
    font-size: 0.95rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .progress-info {
    display: flex;
    justify-content: space-between;
    font-size: 0.82rem;
    color: #6b7280;
    margin-bottom: 0.5rem;
  }
  .progress-info:last-child { margin-top: 0.4rem; margin-bottom: 0; }
  .progress-bar-wrap {
    height: 10px;
    background: #e4e2e2;
    border-radius: 999px;
    overflow: hidden;
  }
  .progress-bar {
    height: 100%;
    background: linear-gradient(90deg, #004d34, #22c55e);
    border-radius: 999px;
    transition: width 0.5s ease;
  }
  .milestones-list {
    display: grid;
    gap: 0.85rem;
  }
  .milestone-row {
    display: grid;
    grid-template-columns: 4rem 1fr 6rem 2.5rem;
    gap: 0.6rem;
    align-items: center;
  }
  .ms-label {
    font-size: 0.8rem;
    font-weight: 700;
    color: #6b7280;
  }
  .ms-bar-wrap {
    height: 8px;
    background: #e4e2e2;
    border-radius: 999px;
    overflow: hidden;
  }
  .ms-bar {
    height: 100%;
    background: #006747;
    border-radius: 999px;
    transition: width 0.5s ease;
  }
  .ms-amount {
    font-size: 0.78rem;
    font-weight: 600;
    color: #1b1c1c;
    text-align: right;
  }
  .ms-pct {
    font-size: 0.78rem;
    color: #9ca3af;
    text-align: right;
  }
  .result-cta {
    text-align: center;
    padding: 1.5rem;
    background: #f0f9f4;
    border-radius: 1.5rem;
    border: 1px solid rgba(0,103,71,0.1);
  }
  .result-cta p {
    margin: 0 0 1rem;
    color: #57534e;
    font-size: 0.9rem;
  }
  .result-cta-btns {
    display: flex;
    gap: 0.7rem;
    justify-content: center;
    flex-wrap: wrap;
  }
  .btn-wa {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    background: #006747;
    color: #fff;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.7rem 1.4rem;
    font-size: 0.88rem;
  }
  .btn-wa .material-symbols-outlined { font-size: 1rem; }
  .btn-packages {
    text-decoration: none;
    border: 1.5px solid #006747;
    color: #006747;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.7rem 1.4rem;
    font-size: 0.88rem;
  }
  @media (max-width: 900px) {
    .calc-grid {
      grid-template-columns: 1fr;
    }
    .results { position: static; }
  }
</style>
