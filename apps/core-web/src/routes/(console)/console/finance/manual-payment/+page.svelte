<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let fBookingId = $state('');
  let fAmount = $state('');
  let fPaymentDate = $state('');
  let fMethod = $state('transfer');
  let fEvidenceUrl = $state('');
  let fNotes = $state('');
  let loading = $state(false);
  let error = $state('');
  let result = $state<any>(null);
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Pembayaran Manual & DP</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Pembayaran Manual & DP</h2>
      <p>Catat pembayaran manual atau uang muka dari jamaah</p>
    </div>

    <div class="card" style="max-width:560px">
      <div class="card-header">
        <span class="material-symbols-outlined">payments</span>
        <h3>Form Pembayaran</h3>
      </div>
      <form class="form-body" onsubmit={async (e) => {
        e.preventDefault(); loading = true; error = ''; result = null;
        try {
          const res = await fetch(`${GATEWAY}/v1/finance/manual-payments`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ booking_id: fBookingId, amount: parseFloat(fAmount), payment_date: fPaymentDate, method: fMethod, evidence_url: fEvidenceUrl, notes: fNotes })
          });
          const body = res.ok ? await res.json() : null;
          if (!res.ok) throw new Error(body?.message ?? `${res.status}`);
          result = body?.data ?? body;
          fBookingId = ''; fAmount = ''; fPaymentDate = ''; fEvidenceUrl = ''; fNotes = '';
        } catch (err: any) { error = err.message ?? 'Gagal menyimpan.'; }
        loading = false;
      }}>
        {#if error}<div class="alert-err">{error}</div>{/if}

        {#if result}
          <div class="result-ok">
            <span class="material-symbols-outlined ok-icon">check_circle</span>
            <div>
              <div class="ok-title">Pembayaran Berhasil Dicatat</div>
              {#if result.entry_id}<div class="ok-detail">Entry ID: <code>{result.entry_id}</code></div>{/if}
              {#if result.journal_id}<div class="ok-detail">Journal ID: <code>{result.journal_id}</code></div>{/if}
              {#if result.amount}<div class="ok-detail">Jumlah: <strong>{formatIDR(result.amount)}</strong></div>{/if}
            </div>
          </div>
        {/if}

        <div class="field">
          <label>ID Booking *</label>
          <input type="text" bind:value={fBookingId} required placeholder="BK-2024-001" />
        </div>
        <div class="field">
          <label>Jumlah (Rp) *</label>
          <input type="number" bind:value={fAmount} required min="0" placeholder="5000000" />
        </div>
        <div class="field">
          <label>Tanggal Pembayaran *</label>
          <input type="date" bind:value={fPaymentDate} required />
        </div>
        <div class="field">
          <label>Metode *</label>
          <select bind:value={fMethod}>
            <option value="transfer">Transfer Bank</option>
            <option value="tunai">Tunai</option>
            <option value="cheque">Cek/Giro</option>
          </select>
        </div>
        <div class="field">
          <label>URL Bukti Pembayaran</label>
          <input type="url" bind:value={fEvidenceUrl} placeholder="https://..." />
        </div>
        <div class="field">
          <label>Catatan</label>
          <textarea bind:value={fNotes} rows="3" placeholder="Opsional..."></textarea>
        </div>
        <button type="submit" class="btn-primary" disabled={loading}>
          {#if loading}
            <span class="material-symbols-outlined spin">progress_activity</span> Menyimpan...
          {:else}
            <span class="material-symbols-outlined">save</span> Simpan Pembayaran
          {/if}
        </button>
      </form>
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
  .card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .card-header { display: flex; align-items: center; gap: 0.5rem; padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .card-header .material-symbols-outlined { color: #004ac6; font-size: 1.1rem; }
  .card-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .form-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; }
  .result-ok { display: flex; align-items: flex-start; gap: 0.75rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 0.25rem; padding: 0.85rem 1rem; }
  .ok-icon { font-size: 1.4rem; color: #15803d; flex-shrink: 0; }
  .ok-title { font-size: 0.85rem; font-weight: 700; color: #15803d; }
  .ok-detail { font-size: 0.78rem; color: #166534; margin-top: 0.2rem; }
  .ok-detail code { font-family: monospace; background: #dcfce7; padding: 0.1rem 0.3rem; border-radius: 0.15rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; resize: vertical; }
  .field input:focus, .field select:focus, .field textarea:focus { border-color: #2563eb; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.55rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-primary .material-symbols-outlined { font-size: 1rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
