<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'progress' | 'evisa' | 'external'>('progress');

  // Tab 1 — Visa Progress
  let vp_departureId = $state('');
  let vp_loading = $state(false);
  let vp_error = $state('');
  let vp_rows = $state<any[]>([]);

  const VISA_STATUSES: Record<string, { label: string; cls: string }> = {
    Submitted: { label: 'Diajukan', cls: 'chip-blue' },
    Approved: { label: 'Disetujui', cls: 'chip-green' },
    Rejected: { label: 'Ditolak', cls: 'chip-red' },
    Pending: { label: 'Pending', cls: 'chip-yellow' },
  };

  async function loadVisaProgress() {
    if (!vp_departureId.trim()) return;
    vp_loading = true; vp_error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/visa-progress/${vp_departureId}`);
      if (!res.ok) throw new Error(`Gagal memuat progress visa (${res.status})`);
      const body = await res.json();
      vp_rows = body.progress ?? body.visas ?? body ?? [];
    } catch (err) {
      vp_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    vp_loading = false;
  }

  // Tab 2 — E-Visa
  let ev_pilgrimId = $state('');
  let ev_departureId = $state('');
  let ev_loading = $state(false);
  let ev_error = $state('');
  let ev_data = $state<any>(null);

  let up_file_url = $state('');
  let up_loading = $state(false);
  let up_error = $state('');
  let up_success = $state('');

  async function loadEvisa() {
    if (!ev_pilgrimId.trim()) return;
    ev_loading = true; ev_error = ''; ev_data = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/evisa/${ev_pilgrimId}?departure_id=${ev_departureId}`);
      if (!res.ok) throw new Error(`Gagal memuat e-visa (${res.status})`);
      ev_data = await res.json();
    } catch (err) {
      ev_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ev_loading = false;
  }

  async function uploadEvisa(e: Event) {
    e.preventDefault();
    up_loading = true; up_error = ''; up_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/evisa`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          pilgrim_id: ev_pilgrimId,
          departure_id: ev_departureId,
          file_url: up_file_url,
        }),
      });
      if (!res.ok) throw new Error(`Gagal upload e-visa (${res.status})`);
      up_success = 'E-Visa berhasil diupload.';
      up_file_url = '';
      await loadEvisa();
    } catch (err) {
      up_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    up_loading = false;
  }

  // Tab 3 — External Provider
  let ext_provider = $state('iata');
  let ext_action = $state('');
  let ext_refId = $state('');
  let ext_payload = $state('{}');
  let ext_loading = $state(false);
  let ext_error = $state('');
  let ext_response = $state<any>(null);

  async function submitExternal(e: Event) {
    e.preventDefault();
    ext_loading = true; ext_error = ''; ext_response = null;
    try {
      let parsedPayload: any = {};
      try { parsedPayload = JSON.parse(ext_payload); } catch { throw new Error('Payload JSON tidak valid'); }
      const res = await fetch(`${GATEWAY}/v1/ops/external-provider`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          provider: ext_provider,
          action: ext_action,
          reference_id: ext_refId,
          payload: parsedPayload,
        }),
      });
      if (!res.ok) throw new Error(`Gagal panggil provider (${res.status})`);
      ext_response = await res.json();
    } catch (err) {
      ext_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ext_loading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Tracking Visa &amp; E-Visa</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Tracking Visa &amp; Repositori E-Visa</h2>
      <p>BL-OPS-031/032/033 — Monitor progress visa, kelola e-visa, integrasi eksternal</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'progress'} onclick={() => tab = 'progress'}>Progress Visa</button>
      <button class="tab-btn" class:active={tab === 'evisa'} onclick={() => tab = 'evisa'}>E-Visa</button>
      <button class="tab-btn" class:active={tab === 'external'} onclick={() => tab = 'external'}>Integrasi Eksternal</button>
    </div>

    {#if tab === 'progress'}
      <div class="section-block">
        <h3 class="section-title">Filter Keberangkatan</h3>
        <div class="inline-filter">
          <input type="text" placeholder="ID Keberangkatan" bind:value={vp_departureId} />
          <button class="btn-primary" onclick={loadVisaProgress} disabled={vp_loading}>
            {#if vp_loading}<span class="spinner"></span>{/if}
            Muat Progress
          </button>
        </div>
        {#if vp_error}<div class="alert-err">{vp_error}</div>{/if}
      </div>

      {#if vp_rows.length > 0}
        <div class="section-block">
          <h3 class="section-title">Status Visa Jamaah</h3>
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Jamaah</th>
                  <th>Nama</th>
                  <th>No. Paspor</th>
                  <th>Status</th>
                  <th>Diajukan</th>
                  <th>Diperbarui</th>
                </tr>
              </thead>
              <tbody>
                {#each vp_rows as row}
                  {@const st = VISA_STATUSES[row.status] ?? { label: row.status, cls: 'chip-gray' }}
                  <tr>
                    <td class="mono">{row.pilgrim_id ?? '-'}</td>
                    <td>{row.pilgrim_name ?? '-'}</td>
                    <td class="mono">{row.passport_number ?? '-'}</td>
                    <td><span class="chip {st.cls}">{st.label}</span></td>
                    <td>{row.submitted_at ? new Date(row.submitted_at).toLocaleDateString('id-ID') : '-'}</td>
                    <td>{row.updated_at ? new Date(row.updated_at).toLocaleDateString('id-ID') : '-'}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}
    {/if}

    {#if tab === 'evisa'}
      <div class="section-block">
        <h3 class="section-title">Cari E-Visa Jamaah</h3>
        <div class="form-row">
          <div class="field">
            <label for="ev-pilgrim">ID Jamaah</label>
            <input id="ev-pilgrim" type="text" placeholder="pilg-001" bind:value={ev_pilgrimId} />
          </div>
          <div class="field">
            <label for="ev-dep">ID Keberangkatan</label>
            <input id="ev-dep" type="text" placeholder="dep-001" bind:value={ev_departureId} />
          </div>
          <div class="field field-actions">
            <button class="btn-primary" onclick={loadEvisa} disabled={ev_loading}>
              {#if ev_loading}<span class="spinner"></span>{/if}
              Cari E-Visa
            </button>
          </div>
        </div>
        {#if ev_error}<div class="alert-err">{ev_error}</div>{/if}
      </div>

      {#if ev_data}
        <div class="section-block">
          <h3 class="section-title">Data E-Visa</h3>
          <div class="evisa-card">
            <div class="evisa-header">
              <span class="material-symbols-outlined">badge</span>
              <span>E-VISA — {ev_data.visa_number ?? ev_data.id ?? 'N/A'}</span>
              <span class="chip chip-green ml-auto">{ev_data.status ?? 'Active'}</span>
            </div>
            <div class="evisa-body">
              <div class="evisa-row"><span>Jamaah</span><strong>{ev_data.pilgrim_name ?? ev_pilgrimId}</strong></div>
              <div class="evisa-row"><span>No. Paspor</span><strong class="mono">{ev_data.passport_number ?? '-'}</strong></div>
              <div class="evisa-row"><span>Berlaku Hingga</span><strong>{ev_data.expiry_date ?? '-'}</strong></div>
              {#if ev_data.file_url}
                <div class="evisa-row">
                  <span>File</span>
                  <a href={ev_data.file_url} target="_blank" class="btn-download">
                    <span class="material-symbols-outlined">download</span>Unduh
                  </a>
                </div>
              {/if}
            </div>
          </div>
        </div>
      {/if}

      <div class="section-block">
        <h3 class="section-title">Upload E-Visa</h3>
        <form class="form-row" onsubmit={uploadEvisa}>
          <div class="field">
            <label for="up-url">URL File E-Visa</label>
            <input id="up-url" type="url" placeholder="https://..." bind:value={up_file_url} required />
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={up_loading}>
              {#if up_loading}<span class="spinner"></span>{/if}
              Upload
            </button>
          </div>
        </form>
        {#if up_error}<div class="alert-err">{up_error}</div>{/if}
        {#if up_success}<div class="alert-ok">{up_success}</div>{/if}
      </div>
    {/if}

    {#if tab === 'external'}
      <div class="section-block">
        <h3 class="section-title">Integrasi Provider Eksternal</h3>
        <form class="form-grid" onsubmit={submitExternal}>
          <div class="field">
            <label for="ext-prov">Provider</label>
            <select id="ext-prov" bind:value={ext_provider}>
              <option value="iata">IATA</option>
              <option value="e-hajj">e-Hajj</option>
            </select>
          </div>
          <div class="field">
            <label for="ext-action">Action</label>
            <input id="ext-action" type="text" placeholder="check_status" bind:value={ext_action} required />
          </div>
          <div class="field">
            <label for="ext-ref">Reference ID</label>
            <input id="ext-ref" type="text" placeholder="REF-001" bind:value={ext_refId} />
          </div>
          <div class="field field-wide">
            <label for="ext-payload">Payload (JSON)</label>
            <textarea id="ext-payload" rows="4" class="mono-ta" bind:value={ext_payload}></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={ext_loading}>
              {#if ext_loading}<span class="spinner"></span>{/if}
              Kirim
            </button>
          </div>
        </form>
        {#if ext_error}<div class="alert-err">{ext_error}</div>{/if}
        {#if ext_response}
          <div class="response-block">
            <div class="response-label">Response</div>
            <pre class="response-pre">{JSON.stringify(ext_response, null, 2)}</pre>
          </div>
        {/if}
      </div>
    {/if}
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
  .tab-bar { display: flex; gap: 0.35rem; margin-bottom: 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); }
  .tab-btn { border: 0; background: transparent; padding: 0.55rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; color: #737686; border-bottom: 2px solid transparent; margin-bottom: -1px; font-family: inherit; }
  .tab-btn.active { color: #2563eb; border-bottom-color: #2563eb; }
  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; margin-bottom: 1.25rem; }
  .section-title { margin: 0 0 1rem; font-size: 0.9rem; font-weight: 700; }
  .inline-filter { display: flex; gap: 0.65rem; align-items: center; }
  .inline-filter input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; min-width: 200px; font-family: inherit; }
  .form-row { display: flex; gap: 0.75rem; align-items: flex-end; flex-wrap: wrap; }
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { grid-column: span 2; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .field-actions { align-self: flex-end; }
  .mono-ta { font-family: 'IBM Plex Mono', monospace; font-size: 0.76rem; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .chip { display: inline-flex; padding: 0.12rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .chip-blue { background: #e0f2fe; color: #075985; }
  .chip-green { background: #dcfce7; color: #166534; }
  .chip-red { background: #fee2e2; color: #991b1b; }
  .chip-yellow { background: #fef9c3; color: #854d0e; }
  .chip-gray { background: #f2f4f6; color: #434655; }
  .ml-auto { margin-left: auto; }
  .evisa-card { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.4rem; overflow: hidden; }
  .evisa-header { display: flex; align-items: center; gap: 0.65rem; padding: 0.75rem 1rem; background: #f2f4f6; font-size: 0.82rem; font-weight: 700; }
  .evisa-body { padding: 0.75rem 1rem; display: flex; flex-direction: column; gap: 0.5rem; }
  .evisa-row { display: flex; justify-content: space-between; font-size: 0.8rem; }
  .evisa-row span { color: #737686; }
  .btn-download { display: inline-flex; align-items: center; gap: 0.3rem; padding: 0.3rem 0.65rem; background: #2563eb; color: #fff; border-radius: 0.2rem; text-decoration: none; font-size: 0.72rem; font-weight: 600; }
  .btn-download .material-symbols-outlined { font-size: 0.85rem; }
  .response-block { margin-top: 1rem; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.35rem; overflow: hidden; }
  .response-label { background: #f2f4f6; padding: 0.4rem 0.75rem; font-size: 0.65rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .response-pre { margin: 0; padding: 0.85rem; font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; color: #191c1e; overflow-x: auto; max-height: 300px; overflow-y: auto; }
</style>
