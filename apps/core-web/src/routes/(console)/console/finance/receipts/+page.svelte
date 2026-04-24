<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }
  function formatDate(iso: string) {
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' });
  }

  // Form state
  let fBookingId = $state('');
  let fPaymentId = $state('');
  let fAmount = $state('');
  let fNotes = $state('');
  let formLoading = $state(false);
  let formError = $state('');
  let formSuccess = $state('');

  // List state
  let receipts = $state<any[]>([]);
  let listLoading = $state(false);

  // Modal state
  let modalReceipt = $state<any>(null);
  let modalLoading = $state(false);

  async function loadReceipts() {
    listLoading = true;
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/receipts`);
      const body = res.ok ? await res.json() : null;
      receipts = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { receipts = []; }
    listLoading = false;
  }

  $effect(() => { loadReceipts(); });

  async function submitForm(e: SubmitEvent) {
    e.preventDefault(); formLoading = true; formError = ''; formSuccess = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/receipts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ booking_id: fBookingId, payment_id: fPaymentId, amount: parseFloat(fAmount), notes: fNotes })
      });
      if (!res.ok) throw new Error(`${res.status}`);
      formSuccess = 'Kwitansi berhasil diterbitkan.';
      fBookingId = ''; fPaymentId = ''; fAmount = ''; fNotes = '';
      await loadReceipts();
    } catch (err: any) { formError = `Gagal menerbitkan: ${err.message}`; }
    formLoading = false;
  }

  async function viewReceipt(id: string) {
    modalLoading = true; modalReceipt = { id, _loading: true };
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/receipts/${id}`);
      const body = res.ok ? await res.json() : null;
      modalReceipt = body?.data ?? body ?? { id, error: 'Tidak ditemukan' };
    } catch { modalReceipt = { id, error: 'Gagal memuat.' }; }
    modalLoading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Kwitansi Digital</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Kwitansi Digital</h2>
      <p>Terbitkan dan kelola kwitansi pembayaran</p>
    </div>

    <div class="grid-2">
      <!-- Issue form -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">receipt_long</span>
          <h3>Terbitkan Kwitansi</h3>
        </div>
        <form class="form-body" onsubmit={submitForm}>
          {#if formError}<div class="alert-err">{formError}</div>{/if}
          {#if formSuccess}<div class="alert-ok">{formSuccess}</div>{/if}
          <div class="field">
            <label>ID Booking *</label>
            <input type="text" bind:value={fBookingId} required placeholder="BK-2024-001" />
          </div>
          <div class="field">
            <label>ID Pembayaran *</label>
            <input type="text" bind:value={fPaymentId} required placeholder="PAY-001" />
          </div>
          <div class="field">
            <label>Jumlah (Rp) *</label>
            <input type="number" bind:value={fAmount} required min="0" placeholder="25000000" />
          </div>
          <div class="field">
            <label>Catatan</label>
            <textarea bind:value={fNotes} rows="2" placeholder="Opsional..."></textarea>
          </div>
          <button type="submit" class="btn-primary" disabled={formLoading}>
            {formLoading ? 'Menerbitkan...' : 'Terbitkan Kwitansi'}
          </button>
        </form>
      </div>

      <!-- List -->
      <div class="card">
        <div class="card-header">
          <span class="material-symbols-outlined">list_alt</span>
          <h3>Kwitansi Terbaru</h3>
        </div>
        {#if listLoading}
          <div class="loading-row"><span class="material-symbols-outlined spin">progress_activity</span> Memuat...</div>
        {:else if receipts.length === 0}
          <div class="empty">Belum ada kwitansi diterbitkan.</div>
        {:else}
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>No. Kwitansi</th><th>Booking</th><th>Jumlah</th><th>Tanggal</th><th></th></tr>
              </thead>
              <tbody>
                {#each receipts as r}
                  <tr>
                    <td class="mono">{r.receipt_number ?? r.id}</td>
                    <td>{r.booking_id}</td>
                    <td>{formatIDR(r.amount)}</td>
                    <td>{r.issued_at ? formatDate(r.issued_at) : '-'}</td>
                    <td>
                      <button class="btn-sm" onclick={() => viewReceipt(r.id ?? r.receipt_number)}>Lihat</button>
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

<!-- Receipt Modal -->
{#if modalReceipt}
  <div class="overlay" role="presentation" onclick={() => modalReceipt = null}></div>
  <div class="modal">
    {#if modalLoading || modalReceipt._loading}
      <div class="modal-loading"><span class="material-symbols-outlined spin">progress_activity</span></div>
    {:else if modalReceipt.error}
      <div class="alert-err">{modalReceipt.error}</div>
      <button class="btn-ghost" onclick={() => modalReceipt = null}>Tutup</button>
    {:else}
      <div class="receipt">
        <div class="receipt-header">
          <div class="receipt-logo">UmrohOS</div>
          <h2 class="receipt-title">KWITANSI</h2>
          <div class="receipt-no">No. {modalReceipt.receipt_number ?? modalReceipt.id}</div>
        </div>
        <div class="receipt-body">
          <div class="receipt-row">
            <span>Tanggal</span><span>{modalReceipt.issued_at ? formatDate(modalReceipt.issued_at) : '-'}</span>
          </div>
          <div class="receipt-row">
            <span>Booking</span><span>{modalReceipt.booking_id ?? '-'}</span>
          </div>
          <div class="receipt-row">
            <span>Pembayaran</span><span>{modalReceipt.payment_id ?? '-'}</span>
          </div>
          {#if modalReceipt.notes}
            <div class="receipt-row">
              <span>Catatan</span><span>{modalReceipt.notes}</span>
            </div>
          {/if}
          <div class="receipt-amount">
            <span>JUMLAH DITERIMA</span>
            <span class="receipt-amount-val">{formatIDR(modalReceipt.amount ?? 0)}</span>
          </div>
        </div>
        <div class="receipt-footer">Dokumen ini adalah kwitansi resmi yang diterbitkan secara digital.</div>
      </div>
      <div style="display:flex;justify-content:center;margin-top:1rem">
        <button class="btn-ghost" onclick={() => modalReceipt = null}>Tutup</button>
      </div>
    {/if}
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
  .field input, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; resize: vertical; }
  .btn-primary { display: inline-flex; align-items: center; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-ghost { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.45rem 0.85rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #191c1e; }
  .btn-sm { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.2rem; padding: 0.25rem 0.55rem; font-size: 0.72rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #2563eb; }
  .btn-sm:hover { background: #eff6ff; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.55rem 0.75rem; font-size: 0.76rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:last-child td { border-bottom: 0; }
  .mono { font-family: monospace; font-size: 0.72rem; }
  .empty { text-align: center; color: #b0b3c1; padding: 2.5rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }

  /* Modal */
  .overlay { position: fixed; inset: 0; background: rgb(0 0 0 / 0.4); z-index: 40; }
  .modal { position: fixed; top: 50%; left: 50%; transform: translate(-50%,-50%); z-index: 50; background: #fff; border-radius: 0.5rem; padding: 1.5rem; width: 420px; max-width: 95vw; box-shadow: 0 8px 32px rgb(0 0 0 / 0.18); }
  .modal-loading { display: flex; justify-content: center; padding: 2rem; }
  .receipt { border: 2px dashed #b0b3c1; border-radius: 0.35rem; padding: 1.5rem; }
  .receipt-header { text-align: center; margin-bottom: 1.25rem; }
  .receipt-logo { font-size: 0.75rem; font-weight: 800; color: #004ac6; letter-spacing: 0.08em; text-transform: uppercase; }
  .receipt-title { margin: 0.25rem 0; font-size: 1.4rem; font-weight: 900; color: #191c1e; letter-spacing: 0.05em; }
  .receipt-no { font-size: 0.75rem; color: #737686; font-family: monospace; }
  .receipt-body { display: flex; flex-direction: column; gap: 0.55rem; }
  .receipt-row { display: flex; justify-content: space-between; font-size: 0.8rem; color: #434655; }
  .receipt-row span:last-child { font-weight: 600; color: #191c1e; }
  .receipt-amount { display: flex; justify-content: space-between; align-items: center; margin-top: 0.75rem; padding-top: 0.75rem; border-top: 2px solid #191c1e; }
  .receipt-amount span:first-child { font-size: 0.72rem; font-weight: 800; letter-spacing: 0.06em; text-transform: uppercase; }
  .receipt-amount-val { font-size: 1.1rem; font-weight: 800; color: #004ac6; font-variant-numeric: tabular-nums; }
  .receipt-footer { margin-top: 1rem; font-size: 0.65rem; color: #b0b3c1; text-align: center; font-style: italic; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1.4rem; }
</style>
