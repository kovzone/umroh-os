<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'collective' | 'ocr' | 'progress'>('collective');

  // Tab 1 — Collective Storage
  let col_departureId = $state('');
  let col_docType = $state('passport');
  let col_uploadUrl = $state('');
  let col_pilgrimId = $state('');
  let col_loading = $state(false);
  let col_error = $state('');
  let col_docs = $state<any[]>([]);

  const DOC_TYPES = [
    { value: 'passport', label: 'Paspor' },
    { value: 'ktp', label: 'KTP' },
    { value: 'photo', label: 'Foto' },
    { value: 'birth_certificate', label: 'Akta Lahir' },
    { value: 'marriage_certificate', label: 'Surat Nikah' },
    { value: 'mahram_letter', label: 'Surat Mahram' },
  ];

  async function submitCollective(e: Event) {
    e.preventDefault();
    col_loading = true; col_error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/collective-docs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          departure_id: col_departureId,
          document_type: col_docType,
          upload_url: col_uploadUrl,
          pilgrim_id: col_pilgrimId,
        }),
      });
      if (!res.ok) throw new Error(`Gagal upload dokumen (${res.status})`);
      const data = await res.json();
      col_docs = [data, ...col_docs];
      col_uploadUrl = ''; col_pilgrimId = '';
    } catch (err) {
      col_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    col_loading = false;
  }

  // Tab 2 — OCR
  let ocr_pilgrimId = $state('');
  let ocr_imageUrl = $state('');
  let ocr_loading = $state(false);
  let ocr_error = $state('');
  let ocr_result = $state<any>(null);

  async function submitOcr(e: Event) {
    e.preventDefault();
    ocr_loading = true; ocr_error = ''; ocr_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/passport-ocr`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pilgrim_id: ocr_pilgrimId, image_url: ocr_imageUrl }),
      });
      if (!res.ok) throw new Error(`Gagal OCR (${res.status})`);
      ocr_result = await res.json();
    } catch (err) {
      ocr_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ocr_loading = false;
  }

  // Tab 3 — Progress & Expiry
  let prog_departureId = $state('');
  let prog_loading = $state(false);
  let prog_error = $state('');
  let prog_data = $state<any>(null);
  let expiry_alerts = $state<any[]>([]);

  async function loadProgress() {
    if (!prog_departureId.trim()) return;
    prog_loading = true; prog_error = '';
    try {
      const [pr, ea] = await Promise.all([
        fetch(`${GATEWAY}/v1/ops/document-progress/${prog_departureId}`),
        fetch(`${GATEWAY}/v1/ops/expiry-alerts?departure_id=${prog_departureId}`),
      ]);
      prog_data = pr.ok ? await pr.json() : null;
      const eaBody = ea.ok ? await ea.json() : { alerts: [] };
      expiry_alerts = eaBody.alerts ?? [];
      if (!pr.ok) throw new Error(`Gagal memuat progress (${pr.status})`);
    } catch (err) {
      prog_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    prog_loading = false;
  }

  const progressItems = $derived(
    prog_data && Array.isArray(prog_data.items) ? prog_data.items : []
  );
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Manajemen Dokumen Jamaah</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Manajemen Dokumen Jamaah</h2>
      <p>BL-OPS-021/022/023 — Upload, OCR, dan tracking kedaluwarsa dokumen</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'collective'} onclick={() => tab = 'collective'}>Penyimpanan Kolektif</button>
      <button class="tab-btn" class:active={tab === 'ocr'} onclick={() => tab = 'ocr'}>OCR Paspor</button>
      <button class="tab-btn" class:active={tab === 'progress'} onclick={() => tab = 'progress'}>Progress &amp; Kedaluwarsa</button>
    </div>

    <!-- TAB 1 -->
    {#if tab === 'collective'}
      <div class="section-block">
        <h3 class="section-title">Upload Dokumen Kolektif</h3>
        <form class="form-grid" onsubmit={submitCollective}>
          <div class="field">
            <label for="col-dep">ID Keberangkatan</label>
            <input id="col-dep" type="text" placeholder="dep-001" bind:value={col_departureId} required />
          </div>
          <div class="field">
            <label for="col-type">Jenis Dokumen</label>
            <select id="col-type" bind:value={col_docType}>
              {#each DOC_TYPES as dt}
                <option value={dt.value}>{dt.label}</option>
              {/each}
            </select>
          </div>
          <div class="field">
            <label for="col-url">URL Upload</label>
            <input id="col-url" type="url" placeholder="https://storage.example.com/doc.pdf" bind:value={col_uploadUrl} required />
          </div>
          <div class="field">
            <label for="col-pilgrim">ID Jamaah</label>
            <input id="col-pilgrim" type="text" placeholder="pilg-001" bind:value={col_pilgrimId} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={col_loading}>
              {#if col_loading}<span class="spinner"></span>{/if}
              Upload Dokumen
            </button>
          </div>
        </form>
        {#if col_error}
          <div class="alert-err">{col_error}</div>
        {/if}
      </div>

      {#if col_docs.length > 0}
        <div class="section-block">
          <h3 class="section-title">Dokumen Terupload</h3>
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>ID Dokumen</th>
                  <th>Jamaah</th>
                  <th>Jenis</th>
                  <th>URL</th>
                  <th>Waktu</th>
                </tr>
              </thead>
              <tbody>
                {#each col_docs as doc}
                  <tr>
                    <td class="mono">{doc.id ?? '-'}</td>
                    <td>{doc.pilgrim_id ?? '-'}</td>
                    <td><span class="chip">{doc.document_type ?? '-'}</span></td>
                    <td>
                      {#if doc.upload_url}
                        <a href={doc.upload_url} target="_blank" class="link">Lihat</a>
                      {:else}—{/if}
                    </td>
                    <td>{doc.created_at ? new Date(doc.created_at).toLocaleString('id-ID') : '-'}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}
    {/if}

    <!-- TAB 2 -->
    {#if tab === 'ocr'}
      <div class="section-block">
        <h3 class="section-title">OCR Paspor</h3>
        <form class="form-grid" onsubmit={submitOcr}>
          <div class="field">
            <label for="ocr-pilgrim">ID Jamaah</label>
            <input id="ocr-pilgrim" type="text" placeholder="pilg-001" bind:value={ocr_pilgrimId} required />
          </div>
          <div class="field">
            <label for="ocr-url">URL Gambar Paspor</label>
            <input id="ocr-url" type="url" placeholder="https://..." bind:value={ocr_imageUrl} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={ocr_loading}>
              {#if ocr_loading}<span class="spinner"></span>{/if}
              Proses OCR
            </button>
          </div>
        </form>
        {#if ocr_error}
          <div class="alert-err">{ocr_error}</div>
        {/if}
      </div>

      {#if ocr_result}
        <div class="section-block">
          <h3 class="section-title">Hasil Ekstraksi OCR</h3>
          <div class="ocr-card">
            <div class="ocr-row"><span class="ocr-label">Nama</span><span>{ocr_result.full_name ?? '-'}</span></div>
            <div class="ocr-row"><span class="ocr-label">No. Paspor</span><span class="mono">{ocr_result.passport_number ?? '-'}</span></div>
            <div class="ocr-row"><span class="ocr-label">Kewarganegaraan</span><span>{ocr_result.nationality ?? '-'}</span></div>
            <div class="ocr-row"><span class="ocr-label">Tanggal Lahir</span><span>{ocr_result.date_of_birth ?? '-'}</span></div>
            <div class="ocr-row"><span class="ocr-label">Berlaku Hingga</span><span>{ocr_result.expiry_date ?? '-'}</span></div>
            <div class="ocr-row"><span class="ocr-label">Jenis Kelamin</span><span>{ocr_result.gender ?? '-'}</span></div>
            <div class="ocr-row"><span class="ocr-label">Tempat Lahir</span><span>{ocr_result.birth_place ?? '-'}</span></div>
            {#if ocr_result.mrz}
              <div class="ocr-row"><span class="ocr-label">MRZ</span><span class="mono text-xs">{ocr_result.mrz}</span></div>
            {/if}
          </div>
        </div>
      {/if}
    {/if}

    <!-- TAB 3 -->
    {#if tab === 'progress'}
      <div class="section-block">
        <h3 class="section-title">Filter Keberangkatan</h3>
        <div class="inline-filter">
          <input type="text" placeholder="ID Keberangkatan" bind:value={prog_departureId} />
          <button class="btn-primary" onclick={loadProgress} disabled={prog_loading}>
            {#if prog_loading}<span class="spinner"></span>{/if}
            Muat Progress
          </button>
        </div>
        {#if prog_error}
          <div class="alert-err">{prog_error}</div>
        {/if}
      </div>

      {#if progressItems.length > 0}
        <div class="section-block">
          <h3 class="section-title">Progress Dokumen per Jenis</h3>
          <div class="progress-list">
            {#each progressItems as item}
              <div class="progress-item">
                <div class="progress-meta">
                  <span class="progress-name">{item.document_type ?? item.type ?? '-'}</span>
                  <span class="progress-count">{item.completed ?? 0}/{item.total ?? 0}</span>
                </div>
                <div class="progress-bar-wrap">
                  <div
                    class="progress-bar-fill"
                    style="width:{item.total > 0 ? Math.round((item.completed / item.total) * 100) : 0}%"
                  ></div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if expiry_alerts.length > 0}
        <div class="section-block">
          <h3 class="section-title">Peringatan Kedaluwarsa</h3>
          <div class="alert-list">
            {#each expiry_alerts as alert}
              <div class="expiry-alert" class:urgent={alert.days_remaining <= 30}>
                <span class="material-symbols-outlined expiry-icon">
                  {alert.days_remaining <= 30 ? 'warning' : 'schedule'}
                </span>
                <div class="expiry-info">
                  <strong>{alert.pilgrim_name ?? alert.pilgrim_id}</strong>
                  <span>{alert.document_type} — kedaluwarsa {alert.expiry_date}</span>
                </div>
                <span class="expiry-days" class:red={alert.days_remaining <= 30}>
                  {alert.days_remaining} hari lagi
                </span>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar {
    position: sticky; top: 0; z-index: 30; height: 4rem;
    background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center;
    backdrop-filter: blur(8px);
  }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 72rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .tab-bar { display: flex; gap: 0.35rem; margin-bottom: 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); padding-bottom: 0; }
  .tab-btn { border: 0; background: transparent; padding: 0.55rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; color: #737686; border-bottom: 2px solid transparent; margin-bottom: -1px; font-family: inherit; }
  .tab-btn.active { color: #2563eb; border-bottom-color: #2563eb; }
  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; margin-bottom: 1.25rem; }
  .section-title { margin: 0 0 1rem; font-size: 0.9rem; font-weight: 700; color: #191c1e; }
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea {
    border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem;
    padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e;
    background: #fff; font-family: inherit;
  }
  .field-actions { align-self: flex-end; }
  .btn-primary {
    display: inline-flex; align-items: center; gap: 0.35rem;
    background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff;
    border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem;
    font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit;
  }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .chip { display: inline-flex; padding: 0.12rem 0.45rem; background: #e0f2fe; color: #075985; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .link { color: #2563eb; text-decoration: none; font-weight: 600; font-size: 0.75rem; }
  .ocr-card { display: grid; grid-template-columns: 1fr 1fr; gap: 0.6rem; }
  .ocr-row { display: flex; flex-direction: column; gap: 0.15rem; padding: 0.65rem 0.85rem; background: #f7f9fb; border-radius: 0.3rem; border: 1px solid rgb(195 198 215 / 0.35); }
  .ocr-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #737686; }
  .text-xs { font-size: 0.65rem; }
  .inline-filter { display: flex; gap: 0.65rem; align-items: center; }
  .inline-filter input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; min-width: 200px; font-family: inherit; }
  .progress-list { display: flex; flex-direction: column; gap: 0.85rem; }
  .progress-item { }
  .progress-meta { display: flex; justify-content: space-between; margin-bottom: 0.35rem; }
  .progress-name { font-size: 0.8rem; font-weight: 600; color: #191c1e; text-transform: capitalize; }
  .progress-count { font-size: 0.75rem; color: #737686; }
  .progress-bar-wrap { height: 0.45rem; background: #e8eaec; border-radius: 999px; overflow: hidden; }
  .progress-bar-fill { height: 100%; background: linear-gradient(90deg,#004ac6,#2563eb); border-radius: 999px; transition: width 0.4s ease; }
  .alert-list { display: flex; flex-direction: column; gap: 0.55rem; }
  .expiry-alert { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: #f7f9fb; border: 1px solid rgb(195 198 215 / 0.35); border-radius: 0.35rem; }
  .expiry-alert.urgent { background: #fff7ed; border-color: #fed7aa; }
  .expiry-icon { font-size: 1.2rem; color: #d97706; }
  .expiry-info { flex: 1; display: flex; flex-direction: column; gap: 0.1rem; font-size: 0.78rem; }
  .expiry-days { font-size: 0.75rem; font-weight: 700; color: #737686; }
  .expiry-days.red { color: #dc2626; }
</style>
