<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let bookingId = $state('');
  let loading = $state(false);
  let error = $state('');
  let relations = $state<any[]>([]);

  async function loadRelations() {
    if (!bookingId.trim()) return;
    loading = true; error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/mahram-relations/${bookingId}`);
      if (!res.ok) throw new Error(`Gagal memuat data (${res.status})`);
      const body = await res.json();
      relations = body.relations ?? body ?? [];
    } catch (err) {
      error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    loading = false;
  }

  // Add relation form
  let form_pilgrimId = $state('');
  let form_mahramId = $state('');
  let form_relation = $state('Suami');
  let form_loading = $state(false);
  let form_error = $state('');
  let form_success = $state('');

  const RELATION_OPTIONS = ['Suami', 'Ayah', 'Saudara Laki-laki', 'Lainnya'];

  async function addRelation(e: Event) {
    e.preventDefault();
    form_loading = true; form_error = ''; form_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/mahram-relations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          booking_id: bookingId,
          pilgrim_id: form_pilgrimId,
          mahram_pilgrim_id: form_mahramId,
          relation: form_relation,
        }),
      });
      if (!res.ok) throw new Error(`Gagal menyimpan (${res.status})`);
      const newRel = await res.json();
      relations = [...relations, newRel];
      form_pilgrimId = ''; form_mahramId = '';
      form_success = 'Hubungan mahram berhasil ditambahkan.';
    } catch (err) {
      form_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    form_loading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Hubungan Mahram</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Hubungan Mahram</h2>
      <p>BL-OPS-022 — Manajemen relasi mahram antar jamaah</p>
    </div>

    <!-- Filter -->
    <div class="section-block">
      <h3 class="section-title">Cari Berdasarkan Booking</h3>
      <div class="inline-filter">
        <input type="text" placeholder="Booking ID" bind:value={bookingId} />
        <button class="btn-primary" onclick={loadRelations} disabled={loading}>
          {#if loading}<span class="spinner"></span>{/if}
          Cari
        </button>
      </div>
      {#if error}
        <div class="alert-err">{error}</div>
      {/if}
    </div>

    <!-- Relations table -->
    {#if relations.length > 0}
      <div class="section-block">
        <h3 class="section-title">Daftar Hubungan Mahram</h3>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Jamaah</th>
                <th>Mahram</th>
                <th>Hubungan</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {#each relations as rel}
                <tr>
                  <td class="mono">{rel.pilgrim_id ?? rel.pilgrim_name ?? '-'}</td>
                  <td class="mono">{rel.mahram_pilgrim_id ?? rel.mahram_name ?? '-'}</td>
                  <td><span class="chip chip-blue">{rel.relation ?? '-'}</span></td>
                  <td>
                    <span class="chip {rel.verified ? 'chip-green' : 'chip-yellow'}">
                      {rel.verified ? 'Terverifikasi' : 'Menunggu'}
                    </span>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {:else if !loading && bookingId}
      <div class="empty-state">Tidak ada data mahram untuk booking ini.</div>
    {/if}

    <!-- Add relation form -->
    <div class="section-block">
      <h3 class="section-title">Tambah Hubungan Mahram</h3>
      <form class="form-grid" onsubmit={addRelation}>
        <div class="field">
          <label for="f-pilgrim">ID Jamaah</label>
          <input id="f-pilgrim" type="text" placeholder="pilg-001" bind:value={form_pilgrimId} required />
        </div>
        <div class="field">
          <label for="f-mahram">ID Mahram</label>
          <input id="f-mahram" type="text" placeholder="pilg-002" bind:value={form_mahramId} required />
        </div>
        <div class="field">
          <label for="f-rel">Hubungan</label>
          <select id="f-rel" bind:value={form_relation}>
            {#each RELATION_OPTIONS as opt}
              <option>{opt}</option>
            {/each}
          </select>
        </div>
        <div class="field field-actions">
          <button type="submit" class="btn-primary" disabled={form_loading}>
            {#if form_loading}<span class="spinner"></span>{/if}
            Tambah Hubungan
          </button>
        </div>
      </form>
      {#if form_error}
        <div class="alert-err">{form_error}</div>
      {/if}
      {#if form_success}
        <div class="alert-ok">{form_success}</div>
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
  .section-title { margin: 0 0 1rem; font-size: 0.9rem; font-weight: 700; color: #191c1e; }
  .inline-filter { display: flex; gap: 0.65rem; align-items: center; }
  .inline-filter input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; min-width: 200px; font-family: inherit; }
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; }
  .field-actions { align-self: flex-end; }
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
  .chip-yellow { background: #fef9c3; color: #854d0e; }
  .empty-state { text-align: center; color: #b0b3c1; padding: 2rem; font-size: 0.82rem; background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; margin-bottom: 1.25rem; }
</style>
