<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let payouts = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  // Create form
  let fAgentId = $state('');
  let fDepartureId = $state('');
  let fBasisAmount = $state('');
  let fRatePct = $state('');
  let fNotes = $state('');
  let formLoading = $state(false);
  let formError = $state('');
  let formSuccess = $state('');

  // Decision modal
  let decisionItem = $state<any>(null);
  let decisionType = $state<'approve' | 'reject'>('approve');
  let decisionNotes = $state('');
  let decisionLoading = $state(false);

  async function load() {
    loading = true; error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/commission-payouts`);
      const body = res.ok ? await res.json() : null;
      payouts = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { error = 'Gagal memuat data komisi.'; payouts = []; }
    loading = false;
  }

  $effect(() => { load(); });

  async function submitForm(e: SubmitEvent) {
    e.preventDefault(); formLoading = true; formError = ''; formSuccess = '';
    try {
      const commission_amount = (parseFloat(fBasisAmount) * parseFloat(fRatePct)) / 100;
      const res = await fetch(`${GATEWAY}/v1/finance/commission-payouts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ agent_id: fAgentId, departure_id: fDepartureId, basis_amount: parseFloat(fBasisAmount), rate_percent: parseFloat(fRatePct), commission_amount, notes: fNotes })
      });
      if (!res.ok) throw new Error(`${res.status}`);
      formSuccess = `Komisi berhasil dibuat — ${formatIDR(commission_amount)}`;
      fAgentId = ''; fDepartureId = ''; fBasisAmount = ''; fRatePct = ''; fNotes = '';
      await load();
    } catch (err: any) { formError = `Gagal: ${err.message}`; }
    formLoading = false;
  }

  async function submitDecision(e: SubmitEvent) {
    e.preventDefault(); decisionLoading = true;
    try {
      await fetch(`${GATEWAY}/v1/finance/commission-payouts/${decisionItem.id}/decision`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ decision: decisionType, notes: decisionNotes })
      });
      decisionItem = null;
      await load();
    } catch { alert('Gagal memproses keputusan.'); }
    decisionLoading = false;
  }

  const commission = $derived(fBasisAmount && fRatePct ? (parseFloat(fBasisAmount) * parseFloat(fRatePct)) / 100 : 0);
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Komisi Agen</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Komisi Agen</h2>
      <p>Buat dan kelola komisi agen penjualan</p>
    </div>

    <div class="grid-2">
      <!-- Create form -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">percent</span>
          <h3>Buat Komisi</h3>
        </div>
        <form class="form-body" onsubmit={submitForm}>
          {#if formError}<div class="alert-err">{formError}</div>{/if}
          {#if formSuccess}<div class="alert-ok">{formSuccess}</div>{/if}
          <div class="field">
            <label>ID Agen *</label>
            <input type="text" bind:value={fAgentId} required placeholder="AGT-001" />
          </div>
          <div class="field">
            <label>ID Keberangkatan *</label>
            <input type="text" bind:value={fDepartureId} required placeholder="DEP-2024-03" />
          </div>
          <div class="field">
            <label>Basis (Rp) *</label>
            <input type="number" bind:value={fBasisAmount} required min="0" placeholder="50000000" />
          </div>
          <div class="field">
            <label>Tarif (%) *</label>
            <input type="number" bind:value={fRatePct} required min="0" max="100" step="0.1" placeholder="2.5" />
          </div>
          {#if commission > 0}
            <div class="calc-preview">
              Komisi: <strong>{formatIDR(commission)}</strong>
            </div>
          {/if}
          <div class="field">
            <label>Catatan</label>
            <textarea bind:value={fNotes} rows="2" placeholder="Opsional..."></textarea>
          </div>
          <button type="submit" class="btn-primary" disabled={formLoading}>
            {formLoading ? 'Menyimpan...' : 'Buat Komisi'}
          </button>
        </form>
      </div>

      <!-- Payouts list -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">pending_actions</span>
          <h3>Komisi Tertunda</h3>
          <button class="btn-ghost" onclick={load} style="margin-left:auto">
            <span class="material-symbols-outlined">refresh</span>
          </button>
        </div>
        {#if loading}
          <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span> Memuat...</div>
        {:else if payouts.length === 0}
          <div class="empty">Tidak ada komisi yang menunggu persetujuan.</div>
        {:else}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Agen</th><th>Keberangkatan</th><th class="ar">Komisi</th><th>Status</th><th>Aksi</th></tr>
              </thead>
              <tbody>
                {#each payouts as p}
                  <tr>
                    <td>{p.agent_id}</td>
                    <td>{p.departure_id ?? '-'}</td>
                    <td class="ar fw">{formatIDR(p.commission_amount ?? 0)}</td>
                    <td>
                      {#if p.status === 'approved'}
                        <span class="chip chip-green">Disetujui</span>
                      {:else if p.status === 'rejected'}
                        <span class="chip chip-red">Ditolak</span>
                      {:else}
                        <span class="chip chip-yellow">Menunggu</span>
                      {/if}
                    </td>
                    <td>
                      {#if !p.status || p.status === 'pending'}
                        <div class="row-actions">
                          <button class="btn-approve" onclick={() => { decisionItem = p; decisionType = 'approve'; decisionNotes = ''; }}>Setujui</button>
                          <button class="btn-reject" onclick={() => { decisionItem = p; decisionType = 'reject'; decisionNotes = ''; }}>Tolak</button>
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
    </div>
  </section>
</main>

<!-- Decision Modal -->
{#if decisionItem}
  <div class="overlay" role="presentation" onclick={() => decisionItem = null}></div>
  <div class="modal">
    <div class="modal-header">
      <h3>{decisionType === 'approve' ? 'Setujui' : 'Tolak'} Komisi</h3>
      <button class="close-btn" onclick={() => decisionItem = null}><span class="material-symbols-outlined">close</span></button>
    </div>
    <form class="modal-body" onsubmit={submitDecision}>
      <p>Agen: <strong>{decisionItem.agent_id}</strong> — {formatIDR(decisionItem.commission_amount ?? 0)}</p>
      <div class="field">
        <label>Catatan</label>
        <textarea bind:value={decisionNotes} rows="3" placeholder="Catatan keputusan..."></textarea>
      </div>
      <div class="modal-actions">
        <button type="button" class="btn-ghost" onclick={() => decisionItem = null}>Batal</button>
        <button type="submit" class={decisionType === 'approve' ? 'btn-approve' : 'btn-reject'} disabled={decisionLoading}>
          {decisionLoading ? 'Memproses...' : decisionType === 'approve' ? 'Setujui' : 'Tolak'}
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
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  @media (max-width: 768px) { .grid-2 { grid-template-columns: 1fr; } }
  .card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .card-header { display: flex; align-items: center; gap: 0.5rem; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .card-header .material-symbols-outlined { color: #004ac6; font-size: 1.1rem; }
  .card-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .form-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #15803d; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.82rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field textarea, .field select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; resize: vertical; }
  .calc-preview { background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; color: #1e40af; }
  .btn-primary { display: inline-flex; align-items: center; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-ghost { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.35rem 0.55rem; cursor: pointer; font-family: inherit; display: inline-flex; align-items: center; }
  .btn-ghost .material-symbols-outlined { font-size: 1rem; color: #737686; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .ar { text-align: right; font-variant-numeric: tabular-nums; }
  .fw { font-weight: 600; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; background: #e0f2fe; color: #075985; }
  .chip-green { background: #dcfce7; color: #15803d; }
  .chip-red { background: #fef2f2; color: #dc2626; }
  .chip-yellow { background: #fef9c3; color: #a16207; }
  .row-actions { display: flex; gap: 0.3rem; }
  .btn-approve { background: #059669; color: #fff; border: none; border-radius: 0.2rem; padding: 0.25rem 0.55rem; font-size: 0.72rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-approve:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-reject { background: #dc2626; color: #fff; border: none; border-radius: 0.2rem; padding: 0.25rem 0.55rem; font-size: 0.72rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-reject:disabled { opacity: 0.6; cursor: not-allowed; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }
  .overlay { position: fixed; inset: 0; background: rgb(0 0 0 / 0.4); z-index: 40; }
  .modal { position: fixed; top: 50%; left: 50%; transform: translate(-50%,-50%); z-index: 50; background: #fff; border-radius: 0.5rem; width: 400px; max-width: 95vw; box-shadow: 0 8px 32px rgb(0 0 0 / 0.18); overflow: hidden; }
  .modal-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); }
  .modal-header h3 { margin: 0; font-size: 1rem; font-weight: 700; }
  .close-btn { border: none; background: transparent; cursor: pointer; color: #737686; display: grid; place-items: center; }
  .close-btn .material-symbols-outlined { font-size: 1.2rem; }
  .modal-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; }
  .modal-body p { margin: 0; font-size: 0.82rem; color: #434655; }
  .modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
