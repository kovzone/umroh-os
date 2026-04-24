<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }
  function formatDate(iso: string) {
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }

  let filterLevel = $state('');
  let items = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  // Decision modal
  let decisionItem = $state<any>(null);
  let decisionType = $state<'approve' | 'reject'>('approve');
  let decisionNotes = $state('');
  let decisionLoading = $state(false);
  let decisionError = $state('');

  async function load() {
    loading = true; error = '';
    try {
      const params = new URLSearchParams();
      if (filterLevel) params.set('level', filterLevel);
      const res = await fetch(`${GATEWAY}/v1/finance/payment-authorizations?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      items = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { error = 'Gagal memuat data.'; items = []; }
    loading = false;
  }

  $effect(() => { load(); });
  $effect(() => { filterLevel; load(); });

  const pendingCount = $derived(items.filter((i: any) => i.status === 'pending' || !i.status).length);

  function openDecision(item: any, type: 'approve' | 'reject') {
    decisionItem = item; decisionType = type; decisionNotes = ''; decisionError = '';
  }

  async function submitDecision(e: SubmitEvent) {
    e.preventDefault(); decisionLoading = true; decisionError = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/payment-authorizations/${decisionItem.id}/decision`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ decision: decisionType, notes: decisionNotes })
      });
      if (!res.ok) throw new Error(`${res.status}`);
      decisionItem = null;
      await load();
    } catch (err: any) { decisionError = `Gagal: ${err.message}`; }
    decisionLoading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Otorisasi Pembayaran</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>
          Otorisasi Pembayaran
          {#if pendingCount > 0}
            <span class="badge">{pendingCount}</span>
          {/if}
        </h2>
        <p>Setujui atau tolak permintaan pembayaran yang menunggu otorisasi</p>
      </div>
    </div>

    {#if error}<div class="alert-err">{error}</div>{/if}

    <!-- Filter -->
    <div class="filter-bar">
      <div class="field">
        <label>Filter Level</label>
        <select bind:value={filterLevel}>
          <option value="">Semua Level</option>
          <option value="1">Level 1</option>
          <option value="2">Level 2</option>
          <option value="3">Level 3</option>
        </select>
      </div>
      <button class="btn-ghost" onclick={load}>
        <span class="material-symbols-outlined">refresh</span>
      </button>
    </div>

    <div class="panel">
      {#if loading}
        <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span> Memuat...</div>
      {:else if items.length === 0}
        <div class="empty">Tidak ada permintaan otorisasi yang tertunda.</div>
      {:else}
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>Batch ID</th><th class="ar">Jumlah</th><th>Diminta Oleh</th><th>Level</th><th>Tanggal</th><th>Status</th><th>Aksi</th></tr>
            </thead>
            <tbody>
              {#each items as item}
                <tr>
                  <td class="mono">{item.batch_id ?? item.id}</td>
                  <td class="ar">{formatIDR(item.amount ?? 0)}</td>
                  <td>{item.requested_by ?? '-'}</td>
                  <td><span class="chip">L{item.level ?? '?'}</span></td>
                  <td>{item.created_at ? formatDate(item.created_at) : '-'}</td>
                  <td>
                    {#if item.status === 'approved'}
                      <span class="chip chip-green">Disetujui</span>
                    {:else if item.status === 'rejected'}
                      <span class="chip chip-red">Ditolak</span>
                    {:else}
                      <span class="chip chip-yellow">Menunggu</span>
                    {/if}
                  </td>
                  <td>
                    {#if !item.status || item.status === 'pending'}
                      <div class="row-actions">
                        <button class="btn-approve" onclick={() => openDecision(item, 'approve')}>Setujui</button>
                        <button class="btn-reject" onclick={() => openDecision(item, 'reject')}>Tolak</button>
                      </div>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  </section>
</main>

<!-- Decision Modal -->
{#if decisionItem}
  <div class="overlay" role="presentation" onclick={() => decisionItem = null}></div>
  <div class="modal">
    <div class="modal-header">
      <h3>{decisionType === 'approve' ? 'Setujui' : 'Tolak'} Pembayaran</h3>
      <button class="close-btn" onclick={() => decisionItem = null}>
        <span class="material-symbols-outlined">close</span>
      </button>
    </div>
    <form class="modal-body" onsubmit={submitDecision}>
      {#if decisionError}<div class="alert-err">{decisionError}</div>{/if}
      <p class="modal-info">Batch: <code>{decisionItem.batch_id ?? decisionItem.id}</code> — {formatIDR(decisionItem.amount ?? 0)}</p>
      <div class="field">
        <label>Catatan {decisionType === 'reject' ? '(wajib untuk penolakan)' : '(opsional)'}</label>
        <textarea bind:value={decisionNotes} rows="3" placeholder="Alasan keputusan..." required={decisionType === 'reject'}></textarea>
      </div>
      <div class="modal-actions">
        <button type="button" class="btn-ghost" onclick={() => decisionItem = null}>Batal</button>
        <button type="submit" class={decisionType === 'approve' ? 'btn-approve' : 'btn-reject'} disabled={decisionLoading}>
          {decisionLoading ? 'Memproses...' : decisionType === 'approve' ? 'Ya, Setujui' : 'Ya, Tolak'}
        </button>
      </div>
    </form>
  </div>
{/if}

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; display: flex; align-items: center; gap: 0.5rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .badge { background: #dc2626; color: #fff; border-radius: 99px; padding: 0.1rem 0.5rem; font-size: 0.72rem; font-weight: 700; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .filter-bar { display: flex; align-items: flex-end; gap: 0.75rem; margin-bottom: 1rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; resize: vertical; }
  .btn-ghost { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.8rem; cursor: pointer; font-family: inherit; display: inline-flex; align-items: center; }
  .btn-ghost .material-symbols-outlined { font-size: 1rem; color: #737686; }
  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 0.85rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .mono { font-family: monospace; font-size: 0.72rem; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .chip-green { background: #dcfce7; color: #15803d; }
  .chip-red { background: #fef2f2; color: #dc2626; }
  .chip-yellow { background: #fef9c3; color: #a16207; }
  .row-actions { display: flex; gap: 0.35rem; }
  .btn-approve { background: #059669; color: #fff; border: none; border-radius: 0.2rem; padding: 0.28rem 0.65rem; font-size: 0.72rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-approve:hover { background: #047857; }
  .btn-approve:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-reject { background: #dc2626; color: #fff; border: none; border-radius: 0.2rem; padding: 0.28rem 0.65rem; font-size: 0.72rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-reject:hover { background: #b91c1c; }
  .btn-reject:disabled { opacity: 0.6; cursor: not-allowed; }
  .empty { text-align: center; color: #b0b3c1; padding: 3rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }
  .overlay { position: fixed; inset: 0; background: rgb(0 0 0 / 0.4); z-index: 40; }
  .modal { position: fixed; top: 50%; left: 50%; transform: translate(-50%,-50%); z-index: 50; background: #fff; border-radius: 0.5rem; width: 420px; max-width: 95vw; box-shadow: 0 8px 32px rgb(0 0 0 / 0.18); overflow: hidden; }
  .modal-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); }
  .modal-header h3 { margin: 0; font-size: 1rem; font-weight: 700; }
  .close-btn { border: none; background: transparent; cursor: pointer; color: #737686; display: grid; place-items: center; }
  .close-btn .material-symbols-outlined { font-size: 1.2rem; }
  .modal-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; }
  .modal-info { margin: 0; font-size: 0.82rem; color: #434655; }
  .modal-info code { font-family: monospace; background: #f2f4f6; padding: 0.1rem 0.3rem; border-radius: 0.15rem; }
  .modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
