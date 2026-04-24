<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let fUserId = $state('');
  let fAction = $state('');
  let fEntityType = $state('');
  let fEntityId = $state('');
  let fFrom = $state('');
  let fTo = $state('');
  let cursor = $state('');
  let prevCursors = $state<string[]>([]);

  let rows = $state<any[]>([]);
  let nextCursor = $state<string | null>(null);
  let loading = $state(false);
  let error = $state('');

  async function load(cur = '') {
    loading = true; error = '';
    try {
      const params = new URLSearchParams();
      if (fUserId) params.set('user_id', fUserId);
      if (fAction) params.set('action', fAction);
      if (fEntityType) params.set('entity_type', fEntityType);
      if (fEntityId) params.set('entity_id', fEntityId);
      if (fFrom) params.set('from', fFrom);
      if (fTo) params.set('to', fTo);
      if (cur) params.set('cursor', cur);
      params.set('limit', '20');
      const res = await fetch(`${GATEWAY}/v1/finance/audit-log?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      rows = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
      nextCursor = body?.next_cursor ?? null;
    } catch { error = 'Gagal memuat log audit.'; rows = []; nextCursor = null; }
    loading = false;
  }

  $effect(() => { load(); });

  function applyFilter(e: SubmitEvent) {
    e.preventDefault(); prevCursors = []; cursor = ''; load('');
  }

  function goNext() {
    if (!nextCursor) return;
    prevCursors = [...prevCursors, cursor];
    cursor = nextCursor;
    load(nextCursor);
  }

  function goPrev() {
    const prev = prevCursors[prevCursors.length - 1] ?? '';
    prevCursors = prevCursors.slice(0, -1);
    cursor = prev;
    load(prev);
  }

  function truncate(s: string, n = 60) {
    return s && s.length > n ? s.slice(0, n) + '…' : (s ?? '-');
  }

  function formatDate(iso: string) {
    return new Date(iso).toLocaleString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Jejak Audit Keuangan</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Jejak Audit Keuangan</h2>
      <p>Riwayat lengkap perubahan data keuangan</p>
    </div>

    <form class="filter-bar" onsubmit={applyFilter}>
      <div class="field">
        <label>User ID</label>
        <input type="text" bind:value={fUserId} placeholder="usr_001" />
      </div>
      <div class="field">
        <label>Aksi</label>
        <select bind:value={fAction}>
          <option value="">Semua</option>
          <option value="create">Create</option>
          <option value="update">Update</option>
          <option value="delete">Delete</option>
        </select>
      </div>
      <div class="field">
        <label>Tipe Entitas</label>
        <input type="text" bind:value={fEntityType} placeholder="invoice" />
      </div>
      <div class="field">
        <label>ID Entitas</label>
        <input type="text" bind:value={fEntityId} placeholder="inv_001" />
      </div>
      <div class="field">
        <label>Dari</label>
        <input type="date" bind:value={fFrom} />
      </div>
      <div class="field">
        <label>Sampai</label>
        <input type="date" bind:value={fTo} />
      </div>
      <button type="submit" class="btn-primary" disabled={loading} style="align-self:flex-end">
        {loading ? 'Memuat...' : 'Filter'}
      </button>
    </form>

    {#if error}<div class="alert-err">{error}</div>{/if}

    <div class="panel">
      {#if loading}
        <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span> Memuat...</div>
      {:else if rows.length === 0}
        <div class="empty">Tidak ada entri audit yang cocok dengan filter.</div>
      {:else}
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>Waktu</th><th>Pengguna</th><th>Aksi</th><th>Tipe Entitas</th><th>ID Entitas</th><th>Detail</th></tr>
            </thead>
            <tbody>
              {#each rows as r}
                <tr>
                  <td class="nowrap">{r.created_at ? formatDate(r.created_at) : '-'}</td>
                  <td>{r.user_id ?? r.user ?? '-'}</td>
                  <td>
                    {#if r.action === 'create'}
                      <span class="chip chip-green">create</span>
                    {:else if r.action === 'delete'}
                      <span class="chip chip-red">delete</span>
                    {:else}
                      <span class="chip chip-blue">update</span>
                    {/if}
                  </td>
                  <td>{r.entity_type ?? '-'}</td>
                  <td class="mono">{r.entity_id ?? '-'}</td>
                  <td class="detail" title={r.diff ?? r.detail ?? ''}>{truncate(r.diff ?? r.detail ?? r.message ?? '')}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="pagination">
          <button class="btn-page" onclick={goPrev} disabled={prevCursors.length === 0 || loading}>
            <span class="material-symbols-outlined">chevron_left</span> Sebelumnya
          </button>
          <span class="page-info">{rows.length} entri ditampilkan</span>
          <button class="btn-page" onclick={goNext} disabled={!nextCursor || loading}>
            Berikutnya <span class="material-symbols-outlined">chevron_right</span>
          </button>
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
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .filter-bar { display: flex; gap: 0.65rem; align-items: flex-start; flex-wrap: wrap; margin-bottom: 1.25rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.45rem 0.6rem; font-size: 0.78rem; color: #191c1e; font-family: inherit; outline: none; min-width: 110px; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; white-space: nowrap; }
  tbody tr:hover { background: #f7f9fb; }
  .nowrap { white-space: nowrap; }
  .mono { font-family: monospace; font-size: 0.72rem; }
  .detail { max-width: 240px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #737686; font-size: 0.72rem; font-family: monospace; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .chip-green { background: #dcfce7; color: #15803d; }
  .chip-red { background: #fef2f2; color: #dc2626; }
  .chip-blue { background: #eff6ff; color: #1d4ed8; }
  .pagination { display: flex; align-items: center; justify-content: space-between; padding: 0.65rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; }
  .btn-page { display: inline-flex; align-items: center; gap: 0.2rem; border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.35rem 0.65rem; font-size: 0.75rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #2563eb; }
  .btn-page:disabled { opacity: 0.5; cursor: not-allowed; color: #737686; }
  .btn-page .material-symbols-outlined { font-size: 0.95rem; }
  .page-info { font-size: 0.72rem; color: #737686; }
  .empty { text-align: center; color: #b0b3c1; padding: 3rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
