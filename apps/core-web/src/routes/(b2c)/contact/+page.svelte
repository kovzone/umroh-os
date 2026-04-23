<script lang="ts">
  import { onMount } from 'svelte';
  import { createLead, listPackages } from '$lib/features/s4-crm/crm-api';
  import type { PackageOption } from '$lib/features/s4-crm/types';

  // ---- form state (Svelte 5 runes) ----
  let name = $state('');
  let phone = $state('');
  let email = $state('');
  let interestPackageId = $state('');
  let notes = $state('');

  // UTM / source (populated in onMount from URL)
  let utmSource = $state('');
  let utmMedium = $state('');
  let utmCampaign = $state('');
  let source = $state('');

  // packages dropdown
  let packages = $state<PackageOption[]>([]);

  // UI state
  type FormState = 'idle' | 'submitting' | 'success' | 'error';
  let formState = $state<FormState>('idle');
  let errorMessage = $state('');

  // ---- validation ----
  const nameError = $derived(name.length > 0 && name.trim().length < 2 ? 'Nama minimal 2 karakter' : '');
  const phoneError = $derived(
    phone.length > 0 && !/^(\+62|08)\d{7,12}$/.test(phone.trim())
      ? 'Format: 08xxx atau +62xxx (minimal 9 digit)'
      : ''
  );
  const canSubmit = $derived(
    name.trim().length >= 2 &&
    /^(\+62|08)\d{7,12}$/.test(phone.trim()) &&
    formState !== 'submitting'
  );

  // ---- lifecycle ----
  onMount(async () => {
    // Capture UTM params from URL
    const params = new URLSearchParams(window.location.search);
    utmSource = params.get('utm_source') ?? '';
    utmMedium = params.get('utm_medium') ?? '';
    utmCampaign = params.get('utm_campaign') ?? '';
    // Derive source from utm_source or referrer
    if (utmSource) {
      source = utmSource;
    } else if (document.referrer.includes('instagram')) {
      source = 'instagram';
    } else if (document.referrer.includes('facebook')) {
      source = 'facebook';
    } else if (document.referrer.includes('tiktok')) {
      source = 'tiktok';
    } else {
      source = 'organic';
    }

    // Load packages for dropdown
    try {
      packages = await listPackages();
    } catch {
      // non-critical — dropdown just stays empty
    }
  });

  // ---- submit ----
  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!canSubmit) return;

    formState = 'submitting';
    errorMessage = '';

    try {
      await createLead({
        name: name.trim(),
        phone: phone.trim(),
        email: email.trim() || undefined,
        source: source || undefined,
        utm_source: utmSource || undefined,
        utm_medium: utmMedium || undefined,
        utm_campaign: utmCampaign || undefined,
        interest_package_id: interestPackageId || undefined,
        notes: notes.trim() || undefined
      });
      formState = 'success';
    } catch (err) {
      formState = 'error';
      errorMessage = err instanceof Error ? err.message : 'Terjadi kesalahan. Silakan coba lagi.';
    }
  }

  function handleRetry() {
    formState = 'idle';
    errorMessage = '';
  }
</script>

<svelte:head>
  <title>Hubungi Kami — UmrohOS</title>
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@400;500;600;700&display=swap"
  />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
  />
</svelte:head>

<main class="contact-shell">
  <div class="hero-band">
    <div class="hero-inner">
      <span class="kicker">
        <span class="material-symbols-outlined">support_agent</span>
        Tim Konsultan Siap Membantu
      </span>
      <h1>Hubungi Kami</h1>
      <p>Isi formulir di bawah dan tim kami akan menghubungi Anda dalam waktu singkat untuk membantu perjalanan ibadah Anda.</p>
    </div>
  </div>

  <div class="content-grid">
    <!-- Left: info panel -->
    <aside class="info-panel">
      <div class="info-card">
        <div class="info-icon">
          <span class="material-symbols-outlined">chat</span>
        </div>
        <div>
          <h3>WhatsApp</h3>
          <p>Respon cepat via WhatsApp Business</p>
          <a class="info-link" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
            +62 812-0000-0000
          </a>
        </div>
      </div>
      <div class="info-card">
        <div class="info-icon info-icon--gold">
          <span class="material-symbols-outlined">schedule</span>
        </div>
        <div>
          <h3>Jam Operasional</h3>
          <p>Senin – Sabtu</p>
          <strong>08.00 – 20.00 WIB</strong>
        </div>
      </div>
      <div class="info-card">
        <div class="info-icon info-icon--teal">
          <span class="material-symbols-outlined">verified_user</span>
        </div>
        <div>
          <h3>Izin Resmi</h3>
          <p>PPIU No. 123/Kemenag/2024</p>
          <strong>Terdaftar & Terpercaya</strong>
        </div>
      </div>
    </aside>

    <!-- Right: form card -->
    <div class="form-card">
      {#if formState === 'success'}
        <div class="success-state" data-testid="lead-success">
          <div class="success-icon">
            <span class="material-symbols-outlined">check_circle</span>
          </div>
          <h2>Terima kasih!</h2>
          <p>Tim kami akan menghubungi Anda segera.</p>
          <p class="success-sub">
            Pesan Anda telah diterima. CS kami akan menghubungi nomor <strong>{phone}</strong> dalam waktu dekat.
          </p>
          <a href="/packages" class="primary-btn">
            <span class="material-symbols-outlined">inventory_2</span>
            Lihat Katalog Paket
          </a>
        </div>

      {:else}
        <div class="form-header">
          <h2>Formulir Konsultasi</h2>
          <p>Semua informasi Anda dijaga kerahasiaannya.</p>
        </div>

        <form onsubmit={handleSubmit} novalidate>
          <!-- Nama -->
          <div class="field-group" class:has-error={!!nameError}>
            <label for="lead-name">
              Nama Lengkap
              <span class="required" aria-hidden="true">*</span>
            </label>
            <input
              id="lead-name"
              type="text"
              placeholder="Contoh: Ahmad Fauzi"
              bind:value={name}
              autocomplete="name"
              required
              data-testid="lead-name"
            />
            {#if nameError}
              <span class="field-error">{nameError}</span>
            {/if}
          </div>

          <!-- Nomor HP -->
          <div class="field-group" class:has-error={!!phoneError}>
            <label for="lead-phone">
              Nomor HP
              <span class="required" aria-hidden="true">*</span>
            </label>
            <div class="input-hint-wrap">
              <input
                id="lead-phone"
                type="tel"
                placeholder="08xxx atau +62xxx"
                bind:value={phone}
                autocomplete="tel"
                required
                data-testid="lead-phone"
              />
              <span class="field-hint">Format: 08123xxxxxxx atau +62812xxxxxxx</span>
            </div>
            {#if phoneError}
              <span class="field-error">{phoneError}</span>
            {/if}
          </div>

          <!-- Email (optional) -->
          <div class="field-group">
            <label for="lead-email">
              Email
              <span class="optional">(opsional)</span>
            </label>
            <input
              id="lead-email"
              type="email"
              placeholder="email@contoh.com"
              bind:value={email}
              autocomplete="email"
              data-testid="lead-email"
            />
          </div>

          <!-- Minat Paket (optional) -->
          <div class="field-group">
            <label for="lead-package">
              Minat Paket
              <span class="optional">(opsional)</span>
            </label>
            <div class="select-wrap">
              <select
                id="lead-package"
                bind:value={interestPackageId}
                data-testid="lead-package"
              >
                <option value="">Pilih paket (jika ada)</option>
                {#each packages as pkg (pkg.id)}
                  <option value={pkg.id}>{pkg.name}</option>
                {/each}
              </select>
              <span class="select-arrow material-symbols-outlined">expand_more</span>
            </div>
          </div>

          <!-- Pesan (optional) -->
          <div class="field-group">
            <label for="lead-notes">
              Pesan / Pertanyaan
              <span class="optional">(opsional)</span>
            </label>
            <textarea
              id="lead-notes"
              rows="4"
              placeholder="Tuliskan pertanyaan atau kebutuhan Anda di sini..."
              bind:value={notes}
              data-testid="lead-notes"
            ></textarea>
          </div>

          {#if formState === 'error'}
            <div class="error-banner" data-testid="lead-error">
              <span class="material-symbols-outlined">error</span>
              <span>{errorMessage}</span>
            </div>
          {/if}

          <div class="form-actions">
            {#if formState === 'error'}
              <button
                type="button"
                class="ghost-btn"
                onclick={handleRetry}
              >
                <span class="material-symbols-outlined">refresh</span>
                Coba Lagi
              </button>
            {/if}
            <button
              type="submit"
              class="primary-btn"
              disabled={!canSubmit}
              data-testid="lead-submit"
            >
              {#if formState === 'submitting'}
                <span class="material-symbols-outlined spin">progress_activity</span>
                Mengirim...
              {:else}
                <span class="material-symbols-outlined">send</span>
                Kirim Pesan
              {/if}
            </button>
          </div>

          <p class="privacy-note">
            <span class="material-symbols-outlined">lock</span>
            Data Anda aman dan tidak akan dibagikan kepada pihak ketiga.
          </p>
        </form>
      {/if}
    </div>
  </div>
</main>

<style>
  :global(.material-symbols-outlined) {
    font-family: 'Material Symbols Outlined', sans-serif;
    font-variation-settings: 'FILL' 0, 'wght' 450, 'GRAD' 0, 'opsz' 24;
  }

  .contact-shell {
    min-height: 100vh;
    background: #f7f9fb;
    font-family: 'IBM Plex Sans', ui-sans-serif, system-ui, sans-serif;
    color: #191c1e;
  }

  /* ---- hero band ---- */
  .hero-band {
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
    padding: 4rem 1.5rem 3rem;
    text-align: center;
  }

  .hero-inner {
    max-width: 42rem;
    margin: 0 auto;
  }

  .kicker {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    background: rgb(255 255 255 / 0.15);
    color: #d4f5e5;
    border-radius: 999px;
    padding: 0.35rem 0.85rem;
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    margin-bottom: 1.2rem;
  }

  .kicker .material-symbols-outlined {
    font-size: 0.95rem;
  }

  .hero-band h1 {
    margin: 0 0 0.85rem;
    font-size: clamp(2rem, 5vw, 3rem);
    font-weight: 800;
    color: #fff;
    letter-spacing: -0.03em;
  }

  .hero-band p {
    margin: 0;
    color: #b5dfce;
    font-size: 1rem;
    line-height: 1.7;
  }

  /* ---- content grid ---- */
  .content-grid {
    max-width: 72rem;
    margin: 0 auto;
    padding: 2.5rem 1.5rem 4rem;
    display: grid;
    grid-template-columns: 18rem 1fr;
    gap: 2rem;
    align-items: start;
  }

  /* ---- info panel ---- */
  .info-panel {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .info-card {
    display: flex;
    align-items: flex-start;
    gap: 0.85rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.5rem;
    padding: 1rem 1.1rem;
  }

  .info-icon {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 0.4rem;
    background: rgb(0 103 71 / 0.1);
    color: #006747;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }

  .info-icon--gold {
    background: rgb(119 90 25 / 0.1);
    color: #775a19;
  }

  .info-icon--teal {
    background: rgb(0 74 198 / 0.1);
    color: #004ac6;
  }

  .info-icon .material-symbols-outlined { font-size: 1.2rem; }

  .info-card h3 {
    margin: 0 0 0.15rem;
    font-size: 0.88rem;
    font-weight: 700;
    color: #191c1e;
  }

  .info-card p {
    margin: 0 0 0.2rem;
    font-size: 0.75rem;
    color: #434655;
  }

  .info-card strong {
    font-size: 0.78rem;
    color: #191c1e;
  }

  .info-link {
    font-size: 0.82rem;
    font-weight: 600;
    color: #006747;
    text-decoration: none;
  }

  .info-link:hover { text-decoration: underline; }

  /* ---- form card ---- */
  .form-card {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.5rem;
    padding: 2rem 2.25rem;
  }

  .form-header {
    margin-bottom: 1.75rem;
  }

  .form-header h2 {
    margin: 0 0 0.35rem;
    font-size: 1.35rem;
    font-weight: 800;
    color: #004d34;
    letter-spacing: -0.02em;
  }

  .form-header p {
    margin: 0;
    font-size: 0.82rem;
    color: #434655;
  }

  /* ---- fields ---- */
  .field-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    margin-bottom: 1.2rem;
  }

  .field-group label {
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }

  .required {
    color: #ba1a1a;
  }

  .optional {
    font-weight: 400;
    text-transform: none;
    letter-spacing: 0;
    color: #737686;
  }

  .field-group input,
  .field-group textarea {
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.3rem;
    padding: 0.6rem 0.75rem;
    font-size: 0.88rem;
    color: #191c1e;
    background: #fff;
    font-family: inherit;
    transition: border-color 0.15s;
    outline: none;
  }

  .field-group input:focus,
  .field-group textarea:focus {
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgb(37 99 235 / 0.12);
  }

  .field-group textarea { resize: vertical; min-height: 7rem; }

  .has-error input,
  .has-error textarea {
    border-color: #ba1a1a;
  }

  .field-error {
    font-size: 0.72rem;
    color: #ba1a1a;
  }

  .input-hint-wrap { display: flex; flex-direction: column; gap: 0.2rem; }

  .field-hint {
    font-size: 0.68rem;
    color: #737686;
  }

  /* select wrapper */
  .select-wrap {
    position: relative;
  }

  .select-wrap select {
    width: 100%;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.3rem;
    padding: 0.6rem 2rem 0.6rem 0.75rem;
    font-size: 0.88rem;
    color: #191c1e;
    background: #fff;
    font-family: inherit;
    appearance: none;
    cursor: pointer;
    outline: none;
    transition: border-color 0.15s;
  }

  .select-wrap select:focus {
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgb(37 99 235 / 0.12);
  }

  .select-arrow {
    position: absolute;
    right: 0.55rem;
    top: 50%;
    transform: translateY(-50%);
    font-size: 1.05rem;
    color: #737686;
    pointer-events: none;
  }

  /* ---- error banner ---- */
  .error-banner {
    display: flex;
    align-items: center;
    gap: 0.55rem;
    background: #ffdad6;
    color: #93000a;
    border-radius: 0.3rem;
    padding: 0.65rem 0.85rem;
    font-size: 0.82rem;
    margin-bottom: 1rem;
  }

  .error-banner .material-symbols-outlined { font-size: 1.1rem; flex-shrink: 0; }

  /* ---- actions ---- */
  .form-actions {
    display: flex;
    gap: 0.65rem;
    justify-content: flex-end;
    margin-bottom: 1rem;
  }

  .primary-btn,
  .ghost-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    border-radius: 0.3rem;
    padding: 0.65rem 1.25rem;
    font-size: 0.88rem;
    font-weight: 700;
    cursor: pointer;
    border: none;
    text-decoration: none;
    font-family: inherit;
    transition: opacity 0.15s, background 0.15s;
  }

  .primary-btn {
    background: linear-gradient(90deg, #004d34, #006747);
    color: #fff;
  }

  .primary-btn:hover { opacity: 0.9; }
  .primary-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  .ghost-btn {
    background: #fff;
    color: #191c1e;
    border: 1px solid rgb(195 198 215 / 0.55);
  }

  .ghost-btn:hover { background: #f2f4f6; }
  .ghost-btn .material-symbols-outlined { font-size: 1rem; }

  /* ---- privacy note ---- */
  .privacy-note {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    margin: 0;
    font-size: 0.68rem;
    color: #737686;
  }

  .privacy-note .material-symbols-outlined { font-size: 0.85rem; }

  /* ---- success state ---- */
  .success-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 2rem 1rem;
    gap: 0.75rem;
  }

  .success-icon {
    width: 4rem;
    height: 4rem;
    border-radius: 999px;
    background: rgb(0 103 71 / 0.1);
    display: grid;
    place-items: center;
    color: #006747;
  }

  .success-icon .material-symbols-outlined {
    font-size: 2rem;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }

  .success-state h2 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 800;
    color: #004d34;
  }

  .success-state p {
    margin: 0;
    color: #434655;
    font-size: 0.9rem;
  }

  .success-sub {
    font-size: 0.82rem !important;
    color: #737686 !important;
    max-width: 28rem;
    line-height: 1.6;
  }

  /* ---- spinner ---- */
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }

  /* ---- responsive ---- */
  @media (max-width: 900px) {
    .content-grid {
      grid-template-columns: 1fr;
    }

    .info-panel {
      flex-direction: row;
      flex-wrap: wrap;
    }

    .info-card {
      flex: 1 1 12rem;
    }
  }

  @media (max-width: 600px) {
    .form-card { padding: 1.25rem 1rem; }
    .hero-band { padding: 3rem 1rem 2.5rem; }
    .form-actions { flex-direction: column; }
    .primary-btn, .ghost-btn { justify-content: center; }
  }
</style>
