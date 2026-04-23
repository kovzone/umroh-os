<script lang="ts">
  import { enhance } from '$app/forms';
  import type { SubmitFunction } from '@sveltejs/kit';
  import type { ActionData, PageData } from './$types';

  let { data, form }: { data: PageData; form: ActionData } = $props();

  let submitting = $state(false);

  const formEnhance: SubmitFunction = () => {
    submitting = true;
    return async ({ update }) => {
      await update();
      submitting = false;
    };
  };
</script>

<main class="page-shell">
  <header class="topbar">
    <a href="/console/packages" class="back-link">
      <span class="material-symbols-outlined">inventory_2</span>
      Katalog
    </a>
    <span class="topbar-sep">/</span>
    <a href="/console/packages/{data.packageId}/departures" class="back-link">
      {data.packageName ?? data.packageId}
    </a>
    <span class="topbar-sep">/</span>
    <span class="topbar-current">Tambah Departure</span>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Tambah Keberangkatan</h2>
        {#if data.packageName}
          <p>Paket: <strong>{data.packageName}</strong></p>
        {/if}
      </div>
    </div>

    <div class="form-card">
      {#if form?.serverError}
        <div class="alert-error" role="alert">
          <span class="material-symbols-outlined">error</span>
          <p>{form.serverError}</p>
        </div>
      {/if}

      <form method="POST" use:enhance={formEnhance} novalidate>
        <div class="field-group">
          <div class="field-row">
            <label class="field" for="departure_date">
              <span class="field-label">Tanggal Keberangkatan <span class="required">*</span></span>
              <input
                id="departure_date"
                name="departure_date"
                type="date"
                required
                value={form?.values?.departure_date ?? ''}
                class:input-error={!!form?.errors?.departure_date}
              />
              {#if form?.errors?.departure_date}
                <span class="error-hint">{form.errors.departure_date}</span>
              {/if}
            </label>

            <label class="field" for="return_date">
              <span class="field-label">Tanggal Kembali <span class="required">*</span></span>
              <input
                id="return_date"
                name="return_date"
                type="date"
                required
                value={form?.values?.return_date ?? ''}
                class:input-error={!!form?.errors?.return_date}
              />
              {#if form?.errors?.return_date}
                <span class="error-hint">{form.errors.return_date}</span>
              {/if}
            </label>
          </div>

          <div class="field-row">
            <label class="field" for="total_seats">
              <span class="field-label">Kapasitas Seats <span class="required">*</span></span>
              <input
                id="total_seats"
                name="total_seats"
                type="number"
                min="1"
                step="1"
                required
                placeholder="40"
                value={form?.values?.total_seats ?? ''}
                class:input-error={!!form?.errors?.total_seats}
              />
              {#if form?.errors?.total_seats}
                <span class="error-hint">{form.errors.total_seats}</span>
              {:else}
                <span class="field-hint">Total kursi yang tersedia untuk keberangkatan ini</span>
              {/if}
            </label>

            <label class="field" for="status">
              <span class="field-label">Status <span class="required">*</span></span>
              <select id="status" name="status" value={form?.values?.status ?? 'draft'}>
                <option value="draft">Draft</option>
                <option value="open">Buka (Open)</option>
                <option value="closed">Tutup (Closed)</option>
              </select>
              <span class="field-hint">Draft = belum dipublish ke publik</span>
            </label>
          </div>
        </div>

        <div class="form-actions">
          <a href="/console/packages/{data.packageId}/departures" class="ghost-btn">Batal</a>
          <button type="submit" class="primary-btn" disabled={submitting}>
            {#if submitting}
              <svg class="spinner" viewBox="0 0 24 24" aria-hidden="true">
                <circle class="spinner-track" cx="12" cy="12" r="10"></circle>
                <path class="spinner-fill" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              Menyimpan...
            {:else}
              <span class="material-symbols-outlined">flight_takeoff</span>
              Simpan Departure
            {/if}
          </button>
        </div>
      </form>
    </div>
  </section>
</main>

<style>
  .page-shell {
    min-height: 100vh;
    background: #f7f9fb;
    font-family: 'IBM Plex Sans', ui-sans-serif, system-ui, -apple-system, sans-serif;
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
    gap: 0.5rem;
    backdrop-filter: blur(8px);
    flex-wrap: wrap;
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    color: #004ac6;
    text-decoration: none;
    font-size: 0.82rem;
    font-weight: 500;
    white-space: nowrap;
  }

  .back-link:hover { text-decoration: underline; }
  .back-link .material-symbols-outlined { font-size: 1rem; }

  .topbar-sep { color: #c3c6d7; font-size: 0.9rem; }

  .topbar-current {
    font-size: 0.82rem;
    color: #434655;
    white-space: nowrap;
  }

  .canvas { padding: 1.5rem; max-width: 52rem; }

  .page-head { margin-bottom: 1.25rem; }

  .page-head h2 {
    margin: 0;
    font-size: 1.5rem;
    color: #191c1e;
  }

  .page-head p {
    margin: 0.3rem 0 0;
    font-size: 0.82rem;
    color: #434655;
  }

  .form-card {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    padding: 1.75rem;
  }

  .alert-error {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.85rem 1rem;
    background: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 0.25rem;
    margin-bottom: 1.25rem;
    color: #dc2626;
    font-size: 0.84rem;
  }

  .alert-error .material-symbols-outlined { font-size: 1.1rem; flex-shrink: 0; }
  .alert-error p { margin: 0; }

  .field-group { display: grid; gap: 1.1rem; }
  .field-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  .field { display: grid; gap: 0.4rem; }

  .field-label {
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: #434655;
  }

  .required { color: #dc2626; }

  .field input,
  .field select {
    border: 1px solid #cbd5e1;
    padding: 0.65rem 0.8rem;
    font: inherit;
    font-size: 0.88rem;
    color: #191c1e;
    background: #fff;
    border-radius: 0.2rem;
    transition: border-color 120ms ease;
    width: 100%;
    box-sizing: border-box;
  }

  .field input:focus,
  .field select:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 1px #2563eb;
  }

  .input-error { border-color: #dc2626 !important; }

  .error-hint { font-size: 0.72rem; color: #dc2626; font-weight: 500; }
  .field-hint { font-size: 0.72rem; color: #737686; }

  .form-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    margin-top: 1.5rem;
    padding-top: 1.25rem;
    border-top: 1px solid rgb(195 198 215 / 0.45);
  }

  .ghost-btn {
    display: inline-flex;
    align-items: center;
    padding: 0.62rem 1rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    color: #434655;
    border-radius: 0.25rem;
    font-size: 0.84rem;
    font-weight: 600;
    cursor: pointer;
    text-decoration: none;
    font-family: inherit;
  }

  .ghost-btn:hover { background: #f2f4f6; }

  .primary-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.62rem 1.1rem;
    border: 1px solid #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
    border-radius: 0.25rem;
    font-size: 0.84rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
  }

  .primary-btn:hover:not(:disabled) { background: linear-gradient(90deg, #003a9e, #1d4ed8); }
  .primary-btn:disabled { opacity: 0.7; cursor: not-allowed; }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  .spinner { width: 1rem; height: 1rem; animation: spin 0.8s linear infinite; }
  .spinner-track { fill: none; stroke: currentColor; stroke-width: 4; opacity: 0.25; }
  .spinner-fill { fill: currentColor; opacity: 0.75; }

  @keyframes spin { to { transform: rotate(360deg); } }

  @media (max-width: 600px) {
    .field-row { grid-template-columns: 1fr; }
  }
</style>
