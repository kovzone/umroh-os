<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  let tab = $state<'refund' | 'penalty'>('refund');

  // ---- Refunds ----
  let rf_bookingId = $state('');
  let rf_amount = $state('');
  let rf_reason = $state('');
  let rf_bankAccount = $state('');
  let rf_loading = $state(false);
  let rf_error = $state('');
  let rf_refunds = $state<any[]>([]);
  let rf_listLoading = $state(false);

  async function loadRefunds() {
    rf_listLoading = true;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/refunds`);
      if (res.ok) {
        const body = await res.json();
        rf_refunds = body.refunds ?? body ?? [];
      }
    } catch { /* ignore */ }
    rf_listLoading = false;
  }

  $effect(() => { loadRefunds(); });

  async function submitRefund(e: Event) {
    e.preventDefault();
    rf_loading = true; rf_error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/refunds`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          booking_id: rf_bookingId,
          amount: parseFloat(rf_amount),
          reason: rf_reason,
          bank_account: rf_bankAccount,
        }),
      });
      if (!res.ok) throw new Error(`Gagal buat refund (${res.status})`);
      const newRefund = await res.json();
      rf_refunds = [newRefund, ...rf_refunds];
      rf_bookingId = ''; rf_amount = ''; rf_reason = ''; rf_bankAccount = '';
    } catch (err) {
      rf_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    rf_loading = false;
  }

  async function decideRefund(id: string, decision: 'approve' | 'reject') {
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/refunds/${id}/decision`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ decision }),
      });
      if (res.ok) {
        const updated = await res.json();
        rf_refunds = rf_refunds.map(r => r.id === id ? updated : r);
      }
    } catch { /* ignore */ }
  }

  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  // ---- Penalties ----
  let pn_bookingId = $state('');
  let pn_penaltyType = $state('cancellation');
  let pn_amount = $state('');
  let pn_notes = $state('');
  let pn_loading = $state(false);
  let pn_error = $state('');
  let pn_success = $state('');

  const PENALTY_TYPES = [
    { value: 'cancellation', label: 'Pembatalan' },
    { value: 'late_payment', label: 'Keterlambatan Pembayaran' },
    { value: 'document_issue', label: 'Masalah Dokumen' },
    { value: 'no_show', label: 'Tidak Hadir' },
  ];

  async function submitPenalty(e: Event) {
    e.preventDefault();
    pn_loading = true; pn_error = ''; pn_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/penalties`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          booking_id: pn_bookingId,
          penalty_type: pn_penaltyType,
          amount: parseFloat(pn_amount),
          notes: pn_notes,
        }),
      });
      if (!res.ok) throw new Error(`Gagal buat penalti (${res.status})`);
      pn_success = 'Penalti berhasil dicatat.';
      pn_bookingId = ''; pn_amount = ''; pn_notes = '';
    } catch (err) {
      pn_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    pn_loading = false;
  }

  const STATUS_CLS: Record<string, string> = {
    pending: 'chip-yellow',
    approved: 'chip-green',
    rejected: 'chip-red',
    processing: 'chip-blue',
  };
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Refund &amp; Penalti</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Refund &amp; Penalti</h2>
      <p>BL-OPS-034 — Pengelolaan permintaan refund dan penalti jamaah</p>
    </div>

    <div class="tab-bar">
      <button class="tab-btn" class:active={tab === 'refund'} onclick={() => tab = 'refund'}>Refund</button>
      <button class="tab-btn" class:active={tab === 'penalty'} onclick={() => tab = 'penalty'}>Penalti</button>
    </div>

    {#if tab === 'refund'}
      <div class="section-block">
        <h3 class="section-title">Buat Permintaan Refund</h3>
        <form class="form-grid" onsubmit={submitRefund}>
          <div class="field">
            <label for="rf-booking">Booking ID</label>
            <input id="rf-booking" type="text" placeholder="BKG-001" bind:value={rf_bookingId} required />
          </div>
          <div class="field">
            <label for="rf-amount">Jumlah Refund (Rp)</label>
            <input id="rf-amount" type="number" min="0" placeholder="5000000" bind:value={rf_amount} required />
          </div>
          <div class="field">
            <label for="rf-bank">Rekening Bank</label>
            <input id="rf-bank" type="text" placeholder="BCA 123456789" bind:value={rf_bankAccount} />
          </div>
          <div class="field field-wide">
            <label for="rf-reason">Alasan</label>
            <textarea id="rf-reason" rows="2" placeholder="Alasan permintaan refund..." bind:value={rf_reason} required></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={rf_loading}>
              {#if rf_loading}<span class="spinner"></span>{/if}
              Buat Refund
            </button>
          </div>
        </form>
        {#if rf_error}<div class="alert-err">{rf_error}</div>{/if}
      </div>

      <div class="section-block">
        <h3 class="section-title">Daftar Refund</h3>
        {#if rf_listLoading}
          <div class="loading-row"><span class="spinner-dark"></span> Memuat...</div>
        {:else if rf_refunds.length === 0}
          <div class="empty-state">Belum ada data refund.</div>
        {:else}
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Booking</th>
                  <th>Jumlah</th>
                  <th>Alasan</th>
                  <th>Status</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                {#each rf_refunds as refund}
                  <tr>
                    <td class="mono">{refund.id ?? '-'}</td>
                    <td class="mono">{refund.booking_id ?? '-'}</td>
                    <td>{formatIDR(refund.amount ?? 0)}</td>
                    <td>{refund.reason ?? '-'}</td>
                    <td>
                      <span class="chip {STATUS_CLS[refund.status] ?? 'chip-gray'}">
                        {refund.status ?? '-'}
                      </span>
                    </td>
                    <td>
                      {#if refund.status === 'pending'}
                        <div class="action-group">
                          <button class="btn-approve" onclick={() => decideRefund(refund.id, 'approve')}>Setujui</button>
                          <button class="btn-reject" onclick={() => decideRefund(refund.id, 'reject')}>Tolak</button>
                        </div>
                      {:else}
                        <span class="text-muted">—</span>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    {/if}

    {#if tab === 'penalty'}
      <div class="section-block">
        <h3 class="section-title">Catat Penalti</h3>
        <form class="form-grid" onsubmit={submitPenalty}>
          <div class="field">
            <label for="pn-booking">Booking ID</label>
            <input id="pn-booking" type="text" placeholder="BKG-001" bind:value={pn_bookingId} required />
          </div>
          <div class="field">
            <label for="pn-type">Jenis Penalti</label>
            <select id="pn-type" bind:value={pn_penaltyType}>
              {#each PENALTY_TYPES as pt}
                <option value={pt.value}>{pt.label}</option>
              {/each}
            </select>
          </div>
          <div class="field">
            <label for="pn-amount">Jumlah (Rp)</label>
            <input id="pn-amount" type="number" min="0" placeholder="500000" bind:value={pn_amount} required />
          </div>
          <div class="field field-wide">
            <label for="pn-notes">Catatan</label>
            <textarea id="pn-notes" rows="2" placeholder="Detail penalti..." bind:value={pn_notes}></textarea>
          </div>
          <div class="field field-actions">
            <button type="submit" class="btn-primary" disabled={pn_loading}>
              {#if pn_loading}<span class="spinner"></span>{/if}
              Simpan Penalti
            </button>
          </div>
        </form>
        {#if pn_error}<div class="alert-err">{pn_error}</div>{/if}
        {#if pn_success}<div class="alert-ok">{pn_success}</div>{/if}
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
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { grid-column: span 2; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  .spinner-dark { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(100 100 100 / 0.3); border-top-color: #434655; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
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
  .action-group { display: flex; gap: 0.4rem; }
  .btn-approve { padding: 0.25rem 0.55rem; font-size: 0.7rem; font-weight: 600; border: 0; border-radius: 0.2rem; background: #dcfce7; color: #166534; cursor: pointer; font-family: inherit; }
  .btn-approve:hover { background: #bbf7d0; }
  .btn-reject { padding: 0.25rem 0.55rem; font-size: 0.7rem; font-weight: 600; border: 0; border-radius: 0.2rem; background: #fee2e2; color: #991b1b; cursor: pointer; font-family: inherit; }
  .btn-reject:hover { background: #fecaca; }
  .text-muted { color: #b0b3c1; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; color: #737686; padding: 1rem 0; }
  .empty-state { text-align: center; color: #b0b3c1; padding: 2rem; font-size: 0.82rem; }
</style>
