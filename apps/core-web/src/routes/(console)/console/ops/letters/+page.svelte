<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  // Section 1 — Official Letters
  let lt_template = $state('Surat Keterangan');
  let lt_departureId = $state('');
  let lt_pilgrimId = $state('');
  let lt_issuedTo = $state('');
  let lt_notes = $state('');
  let lt_loading = $state(false);
  let lt_error = $state('');
  let lt_result = $state<any>(null);

  const TEMPLATE_OPTIONS = ['Surat Keterangan', 'Surat Izin', 'Surat Rekomendasi'];

  async function submitLetter(e: Event) {
    e.preventDefault();
    lt_loading = true; lt_error = ''; lt_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/official-letters`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          template_name: lt_template,
          departure_id: lt_departureId,
          pilgrim_id: lt_pilgrimId,
          issued_to: lt_issuedTo,
          notes: lt_notes,
        }),
      });
      if (!res.ok) throw new Error(`Gagal membuat surat (${res.status})`);
      lt_result = await res.json();
    } catch (err) {
      lt_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    lt_loading = false;
  }

  // Section 2 — Immigration Manifest
  let mf_departureId = $state('');
  let mf_format = $state('PDF');
  let mf_loading = $state(false);
  let mf_error = $state('');
  let mf_result = $state<any>(null);

  async function submitManifest(e: Event) {
    e.preventDefault();
    mf_loading = true; mf_error = ''; mf_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/immigration-manifest/${mf_departureId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ format: mf_format }),
      });
      if (!res.ok) throw new Error(`Gagal buat manifest (${res.status})`);
      mf_result = await res.json();
    } catch (err) {
      mf_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    mf_loading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Surat Resmi &amp; Manifest Imigrasi</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Surat Resmi &amp; Manifest Imigrasi</h2>
      <p>BL-OPS-024/025 — Generate surat resmi dan manifest keberangkatan</p>
    </div>

    <!-- Official Letters -->
    <div class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">description</span>
        <h3>Surat Resmi</h3>
      </div>
      <form class="form-grid" onsubmit={submitLetter}>
        <div class="field">
          <label for="lt-template">Template Surat</label>
          <select id="lt-template" bind:value={lt_template}>
            {#each TEMPLATE_OPTIONS as t}
              <option>{t}</option>
            {/each}
          </select>
        </div>
        <div class="field">
          <label for="lt-dep">ID Keberangkatan</label>
          <input id="lt-dep" type="text" placeholder="dep-001" bind:value={lt_departureId} required />
        </div>
        <div class="field">
          <label for="lt-pilgrim">ID Jamaah</label>
          <input id="lt-pilgrim" type="text" placeholder="pilg-001" bind:value={lt_pilgrimId} required />
        </div>
        <div class="field">
          <label for="lt-issued">Ditujukan Kepada</label>
          <input id="lt-issued" type="text" placeholder="Kepala Dinas..." bind:value={lt_issuedTo} required />
        </div>
        <div class="field field-wide">
          <label for="lt-notes">Catatan</label>
          <textarea id="lt-notes" rows="2" placeholder="Catatan tambahan..." bind:value={lt_notes}></textarea>
        </div>
        <div class="field field-actions">
          <button type="submit" class="btn-primary" disabled={lt_loading}>
            {#if lt_loading}<span class="spinner"></span>{/if}
            Generate Surat
          </button>
        </div>
      </form>
      {#if lt_error}
        <div class="alert-err">{lt_error}</div>
      {/if}
      {#if lt_result}
        <div class="result-card">
          <span class="material-symbols-outlined result-icon">check_circle</span>
          <div class="result-info">
            <div class="result-title">Surat berhasil dibuat</div>
            <div class="result-sub">{lt_result.letter_id ?? lt_result.id ?? ''}</div>
          </div>
          {#if lt_result.letter_url ?? lt_result.url}
            <a href={lt_result.letter_url ?? lt_result.url} target="_blank" class="btn-download">
              <span class="material-symbols-outlined">download</span>
              Unduh
            </a>
          {/if}
        </div>
      {/if}
    </div>

    <!-- Immigration Manifest -->
    <div class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">article</span>
        <h3>Manifest Imigrasi</h3>
      </div>
      <form class="form-grid" onsubmit={submitManifest}>
        <div class="field">
          <label for="mf-dep">ID Keberangkatan</label>
          <input id="mf-dep" type="text" placeholder="dep-001" bind:value={mf_departureId} required />
        </div>
        <div class="field">
          <label for="mf-fmt">Format</label>
          <select id="mf-fmt" bind:value={mf_format}>
            <option>PDF</option>
            <option>CSV</option>
          </select>
        </div>
        <div class="field field-actions">
          <button type="submit" class="btn-primary" disabled={mf_loading}>
            {#if mf_loading}<span class="spinner"></span>{/if}
            Buat Manifest
          </button>
        </div>
      </form>
      {#if mf_error}
        <div class="alert-err">{mf_error}</div>
      {/if}
      {#if mf_result}
        <div class="result-card">
          <span class="material-symbols-outlined result-icon">check_circle</span>
          <div class="result-info">
            <div class="result-title">Manifest berhasil dibuat ({mf_format})</div>
            <div class="result-sub">{mf_result.manifest_id ?? mf_result.id ?? ''}</div>
          </div>
          {#if mf_result.manifest_url ?? mf_result.url}
            <a href={mf_result.manifest_url ?? mf_result.url} target="_blank" class="btn-download">
              <span class="material-symbols-outlined">download</span>
              Unduh
            </a>
          {/if}
        </div>
      {/if}
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 72rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; margin-bottom: 1.25rem; }
  .section-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 1rem; }
  .section-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .section-icon { font-size: 1.15rem; color: #004ac6; }
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { grid-column: span 2; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .result-card { display: flex; align-items: center; gap: 0.85rem; padding: 0.85rem 1rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 0.35rem; margin-top: 0.85rem; }
  .result-icon { font-size: 1.5rem; color: #16a34a; }
  .result-info { flex: 1; }
  .result-title { font-size: 0.82rem; font-weight: 700; color: #166534; }
  .result-sub { font-size: 0.72rem; color: #737686; margin-top: 0.1rem; }
  .btn-download { display: inline-flex; align-items: center; gap: 0.3rem; padding: 0.45rem 0.85rem; background: #2563eb; color: #fff; border-radius: 0.25rem; text-decoration: none; font-size: 0.78rem; font-weight: 600; }
  .btn-download .material-symbols-outlined { font-size: 0.95rem; }
</style>
