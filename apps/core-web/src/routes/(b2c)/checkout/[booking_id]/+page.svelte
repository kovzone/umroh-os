<script lang="ts">
  /**
   * S2-L-02: Checkout page — VA info + countdown + status polling.
   *
   * States:
   *   loading          → issuing VA (or re-issuing after expiry)
   *   waiting_payment  → VA issued, awaiting transfer
   *   partially_paid   → DP received, awaiting remainder
   *   paid             → full payment received
   *   expired          → VA expired, user can request a new one
   *   error            → issue VA failed
   */
  import { onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import MarketingPageLayout from '$lib/components/marketing/MarketingPageLayout.svelte';
  import VAInfoCard from '$lib/components/payment/VAInfoCard.svelte';
  import { issuePaymentInvoice, pollInvoice } from '$lib/features/s2-payment/repository';
  import type { CheckoutStatus, Invoice, InvoiceWithVA } from '$lib/features/s2-payment/types';
  import { formatIdrAmountLabel } from '$lib/utils/format-idr';

  let { data } = $props();

  // ---- State ----

  let invoiceWithVA = $state<InvoiceWithVA | null>(null);
  let invoice = $state<Invoice | null>(null);
  let status = $state<CheckoutStatus>('loading');
  let errorMessage = $state('');

  // ---- Sync from page data (runs once on mount, and if data changes) ----

  $effect(() => {
    invoiceWithVA = data.initialInvoice ?? null;
    errorMessage = data.issueError ?? '';
    status = data.issueError ? 'error' : data.initialInvoice ? 'waiting_payment' : 'loading';
  });

  // ---- Derived ----

  const amountLabel = $derived(
    invoice
      ? formatIdrAmountLabel(invoice.amount_total)
      : invoiceWithVA
        ? formatIdrAmountLabel(invoiceWithVA.va.amount_total)
        : '—'
  );

  const paidAmountLabel = $derived(
    invoice && invoice.paid_amount > 0 ? formatIdrAmountLabel(invoice.paid_amount) : null
  );

  const remainingLabel = $derived(
    invoice && invoice.paid_amount > 0
      ? formatIdrAmountLabel(Math.max(0, invoice.amount_total - invoice.paid_amount))
      : null
  );

  const bookingCodeDisplay = $derived(
    `UMR-${data.bookingId.slice(0, 8).toUpperCase()}`
  );

  // ---- Intervals ----

  let pollingInterval: ReturnType<typeof setInterval> | undefined;

  function startPolling(invoiceId: string) {
    clearInterval(pollingInterval);
    pollingInterval = setInterval(async () => {
      try {
        const result = await pollInvoice(invoiceId);
        invoice = result;

        if (result.status === 'paid') {
          status = 'paid';
          stopPolling();
        } else if (result.status === 'partially_paid') {
          status = 'partially_paid';
        } else if (result.status === 'void') {
          status = 'expired';
          stopPolling();
        }
      } catch {
        // Silently ignore transient poll failures; stop after too many?
        // For MVP: just keep polling until status resolves.
      }
    }, 5000);
  }

  function stopPolling() {
    clearInterval(pollingInterval);
    pollingInterval = undefined;
  }

  // ---- VA expiry handler ----

  function handleVAExpired() {
    if (status !== 'paid' && status !== 'partially_paid') {
      status = 'expired';
      stopPolling();
    }
  }

  // ---- Re-issue VA ----

  async function reissueVA() {
    status = 'loading';
    errorMessage = '';
    invoiceWithVA = null;
    invoice = null;

    try {
      // Backend should create a fresh VA when old one is expired
      const fresh = await issuePaymentInvoice(data.bookingId);
      invoiceWithVA = fresh;
      status = 'waiting_payment';
      startPolling(fresh.invoice_id);
    } catch (err) {
      const e = err as Error;
      errorMessage = e.message ?? 'Gagal menerbitkan VA baru.';
      status = 'error';
    }
  }

  // ---- Lifecycle ----

  $effect(() => {
    if (invoiceWithVA && status === 'waiting_payment') {
      startPolling(invoiceWithVA.invoice_id);
    }

    return () => {
      stopPolling();
    };
  });

  onDestroy(() => {
    stopPolling();
  });
</script>

<MarketingPageLayout ctaHref="/packages" ctaLabel="Masuk" packagesLinkActive={true}>
  <div class="co-shell" data-testid="s2-checkout-shell">
    <main class="shell co-main">

      <!-- Back navigation -->
      <a class="back" href="/packages">
        <span class="material-symbols-outlined">arrow_back</span>
        Kembali ke paket
      </a>

      <div class="co-grid">
        <!-- Main column -->
        <div class="main-col">

          <!-- Page header -->
          <header class="page-h">
            <h1 class="h1">Pembayaran</h1>
            <p class="lead">
              Selesaikan transfer sebelum batas waktu untuk mengamankan reservasi Anda.
            </p>
          </header>

          <!-- ---- LOADING STATE ---- -->
          {#if status === 'loading'}
            <section class="card card-pad state-card" aria-busy="true">
              <div class="spinner" aria-hidden="true"></div>
              <p class="state-title">Menyiapkan instruksi pembayaran…</p>
              <p class="state-sub">Mohon tunggu, kami sedang menerbitkan Virtual Account untuk pesanan Anda.</p>
            </section>

          <!-- ---- ERROR STATE ---- -->
          {:else if status === 'error'}
            <section class="card card-pad state-card error-state" role="alert">
              <span class="material-symbols-outlined state-ico err-ico">error</span>
              <p class="state-title">Gagal menerbitkan VA</p>
              <p class="state-sub">
                {errorMessage || 'Terjadi kesalahan saat menerbitkan Virtual Account. Silakan coba lagi.'}
              </p>
              <button type="button" class="btn-primary" onclick={() => void reissueVA()}>
                Coba lagi
              </button>
            </section>

          <!-- ---- PAID STATE ---- -->
          {:else if status === 'paid'}
            <section class="card card-pad state-card success-state" role="status">
              <span class="material-symbols-outlined state-ico success-ico">check_circle</span>
              <p class="state-title">Pembayaran diterima!</p>
              <p class="state-sub">
                Terima kasih! Pembayaran sebesar <strong>{amountLabel}</strong> telah kami terima.
                Tim kami akan segera memproses pemesanan Anda.
              </p>
              <div class="success-actions">
                <a class="btn-primary" href="/packages">
                  <span class="material-symbols-outlined">home</span>
                  Kembali ke beranda
                </a>
              </div>
              <p class="receipt-hint">
                <span class="material-symbols-outlined">receipt_long</span>
                Kuitansi dikirim via WhatsApp dan email Anda.
              </p>
            </section>

          <!-- ---- EXPIRED STATE ---- -->
          {:else if status === 'expired'}
            <section class="card card-pad state-card expired-state" role="alert">
              <span class="material-symbols-outlined state-ico expired-ico">timer_off</span>
              <p class="state-title">Virtual Account kadaluarsa</p>
              <p class="state-sub">
                Batas waktu pembayaran telah terlewat. Anda dapat meminta VA baru untuk melanjutkan pembayaran.
              </p>
              <button type="button" class="btn-primary" onclick={() => void reissueVA()}>
                <span class="material-symbols-outlined">refresh</span>
                Minta VA Baru
              </button>
            </section>

          <!-- ---- WAITING / PARTIALLY PAID STATES ---- -->
          {:else if (status === 'waiting_payment' || status === 'partially_paid') && invoiceWithVA}

            <!-- Booking code bar -->
            <div class="booking-bar">
              <div class="bbar-left">
                <span class="material-symbols-outlined bbar-ico">confirmation_number</span>
                <div>
                  <p class="bbar-label">Kode booking</p>
                  <p class="bbar-code">{bookingCodeDisplay}</p>
                </div>
              </div>
              <div class="status-chip" class:partial={status === 'partially_paid'}>
                <span class="pulse"></span>
                <span>{status === 'partially_paid' ? 'DP Diterima' : 'Menunggu Pembayaran'}</span>
              </div>
            </div>

            <!-- Partially paid info -->
            {#if status === 'partially_paid' && paidAmountLabel && remainingLabel}
              <div class="partial-alert">
                <span class="material-symbols-outlined">payments</span>
                <div>
                  <p class="pa-title">DP sebesar {paidAmountLabel} telah diterima</p>
                  <p class="pa-sub">Sisa pembayaran yang harus dilunasi: <strong>{remainingLabel}</strong></p>
                </div>
              </div>
            {/if}

            <!-- VA Info Card (main component) -->
            <VAInfoCard va={invoiceWithVA.va} onExpired={handleVAExpired} />

            <!-- Polling indicator -->
            <div class="polling-row" aria-live="polite">
              <span class="poll-dot" aria-hidden="true"></span>
              <p class="poll-text">Memverifikasi pembayaran secara otomatis setiap 5 detik…</p>
            </div>

            <!-- Help & actions -->
            <div class="help-row">
              <a
                class="btn-outline"
                href="https://wa.me/6281200000000"
                target="_blank"
                rel="noopener noreferrer"
              >
                <span class="material-symbols-outlined">support_agent</span>
                Hubungi CS
              </a>
              <button
                type="button"
                class="link-muted"
                onclick={() => void goto('/packages')}
              >
                Bayar nanti (pesanan tersimpan)
              </button>
            </div>
          {/if}

        </div>

        <!-- Aside summary -->
        <aside class="aside-col">
          <div class="summary-card card card-pad">
            <h3 class="h3">Ringkasan pesanan</h3>

            <div class="sum-row">
              <span class="sum-label">Kode booking</span>
              <span class="sum-val mono">{bookingCodeDisplay}</span>
            </div>

            {#if amountLabel !== '—'}
              <div class="sum-row total-row">
                <span class="sum-label">Total pembayaran</span>
                <span class="sum-amt">{amountLabel}</span>
              </div>
            {/if}

            {#if status === 'paid'}
              <div class="sum-status paid-badge">
                <span class="material-symbols-outlined">check_circle</span>
                Lunas
              </div>
            {:else if status === 'partially_paid' && paidAmountLabel}
              <div class="sum-status partial-badge">
                <span class="material-symbols-outlined">timelapse</span>
                DP: {paidAmountLabel}
              </div>
            {:else if status === 'waiting_payment'}
              <div class="sum-status waiting-badge">
                <span class="material-symbols-outlined">pending</span>
                Menunggu transfer
              </div>
            {:else if status === 'expired'}
              <div class="sum-status expired-badge">
                <span class="material-symbols-outlined">timer_off</span>
                VA Kadaluarsa
              </div>
            {/if}

            <div class="sum-note">
              <span class="material-symbols-outlined">verified_user</span>
              <p>Pembayaran Anda diproses secara aman melalui gateway terverifikasi.</p>
            </div>
          </div>

          {#if status === 'waiting_payment' || status === 'partially_paid'}
            <div class="tips-card card card-pad">
              <h4 class="h4">Tips pembayaran</h4>
              <ul class="tips">
                <li>Pastikan nominal transfer tepat sesuai yang tertera.</li>
                <li>Tidak perlu konfirmasi manual — sistem memverifikasi otomatis.</li>
                <li>Simpan bukti transfer sebagai referensi.</li>
              </ul>
            </div>
          {/if}
        </aside>
      </div>
    </main>
  </div>
</MarketingPageLayout>

<style>
  .co-main {
    padding-top: 1.5rem;
    padding-bottom: 4rem;
  }
  .back {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    margin-bottom: 1.75rem;
    font-size: 0.9rem;
    font-weight: 600;
    color: #6f7a72;
    text-decoration: none;
    transition: color 0.15s ease, gap 0.15s ease;
  }
  .back:hover {
    color: #004d34;
    gap: 0.5rem;
  }

  .co-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2.25rem;
    align-items: flex-start;
  }
  @media (min-width: 1024px) {
    .co-grid {
      grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
      gap: 2.5rem;
    }
  }

  .main-col {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .aside-col {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  @media (min-width: 1024px) {
    .aside-col {
      position: sticky;
      top: 5.75rem;
    }
  }

  /* Typography */
  .h1 {
    margin: 0 0 0.35rem;
    font-size: clamp(1.75rem, 3vw, 2.25rem);
    font-weight: 800;
    color: #004d34;
    letter-spacing: -0.02em;
  }
  .h3 {
    margin: 0 0 1.25rem;
    font-size: 1.1rem;
    font-weight: 700;
    color: #1b1c1c;
  }
  .h4 {
    margin: 0 0 0.75rem;
    font-size: 0.95rem;
    font-weight: 700;
    color: #1b1c1c;
  }
  .lead {
    margin: 0;
    color: #6f7a72;
    font-size: 0.95rem;
    line-height: 1.5;
  }

  /* Cards */
  .card {
    background: #fff;
    border-radius: 1.25rem;
    box-shadow: 0 8px 32px rgba(27, 28, 28, 0.03);
  }
  .card-pad {
    padding: 1.75rem 1.5rem;
  }
  @media (min-width: 640px) {
    .card-pad {
      padding: 2rem;
    }
  }

  /* State cards */
  .state-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 0.75rem;
    padding-top: 2.5rem;
    padding-bottom: 2.5rem;
    border: 1px solid rgba(190, 201, 193, 0.35);
  }
  .state-title {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 800;
    color: #1b1c1c;
  }
  .state-sub {
    margin: 0;
    font-size: 0.9rem;
    color: #6f7a72;
    line-height: 1.5;
    max-width: 28rem;
  }
  .state-ico {
    font-size: 3rem;
  }
  .success-state {
    border-color: rgba(0, 103, 71, 0.25);
    background: linear-gradient(135deg, rgba(0, 77, 52, 0.03) 0%, #fff 100%);
  }
  .success-ico { color: #006747; }
  .err-ico { color: #93000a; }
  .expired-ico { color: #775a19; }

  .success-actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    justify-content: center;
    margin-top: 0.5rem;
  }
  .receipt-hint {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.78rem;
    color: #6f7a72;
    margin: 0.25rem 0 0;
  }
  .receipt-hint .material-symbols-outlined {
    font-size: 1rem;
    color: #006747;
  }

  /* Loading spinner */
  .spinner {
    width: 2.5rem;
    height: 2.5rem;
    border: 3px solid rgba(0, 103, 71, 0.15);
    border-top-color: #006747;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Booking bar */
  .booking-bar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 1rem 1.35rem;
    background: #fff;
    border: 1px solid rgba(190, 201, 193, 0.35);
    border-radius: 1rem;
  }
  .bbar-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  .bbar-ico {
    font-size: 1.5rem;
    color: #004d34;
  }
  .bbar-label {
    margin: 0 0 0.15rem;
    font-size: 0.6rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #6f7a72;
  }
  .bbar-code {
    margin: 0;
    font-weight: 800;
    font-size: 1.05rem;
    color: #004d34;
    font-family: ui-monospace, monospace;
    letter-spacing: 0.04em;
  }
  .status-chip {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.45rem 0.9rem;
    border-radius: 999px;
    background: rgba(119, 90, 25, 0.1);
    border: 1px solid rgba(119, 90, 25, 0.2);
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #775a19;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .status-chip.partial {
    background: rgba(0, 103, 71, 0.08);
    border-color: rgba(0, 103, 71, 0.2);
    color: #006747;
  }
  .pulse {
    width: 7px;
    height: 7px;
    border-radius: 999px;
    background: currentColor;
    animation: pulse 1.4s ease-in-out infinite;
  }
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.3; }
  }

  /* Partial payment alert */
  .partial-alert {
    display: flex;
    gap: 0.85rem;
    align-items: flex-start;
    padding: 1rem 1.25rem;
    border-radius: 0.75rem;
    background: rgba(0, 103, 71, 0.07);
    border-left: 3px solid #006747;
    color: #004d34;
  }
  .partial-alert .material-symbols-outlined {
    font-size: 1.5rem;
    flex-shrink: 0;
    color: #006747;
  }
  .pa-title {
    margin: 0 0 0.25rem;
    font-weight: 700;
    font-size: 0.92rem;
  }
  .pa-sub {
    margin: 0;
    font-size: 0.82rem;
    color: #3f6754;
  }

  /* Polling row */
  .polling-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.65rem 1rem;
    border-radius: 0.5rem;
    background: #f5f3f3;
  }
  .poll-dot {
    width: 6px;
    height: 6px;
    border-radius: 999px;
    background: #6f7a72;
    animation: pulse 1.8s ease-in-out infinite;
    flex-shrink: 0;
  }
  .poll-text {
    margin: 0;
    font-size: 0.78rem;
    color: #6f7a72;
    font-style: italic;
  }

  /* Help row */
  .help-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
  }

  /* Buttons */
  .btn-primary {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.35rem;
    padding: 0.9rem 1.75rem;
    border: none;
    border-radius: 999px;
    font-weight: 700;
    font-size: 1rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
    color: #fff;
    background: linear-gradient(90deg, #004d34, #006747);
    box-shadow: 0 12px 28px rgba(0, 77, 52, 0.22);
    cursor: pointer;
    text-decoration: none;
    transition: transform 0.15s ease, opacity 0.15s ease;
  }
  .btn-primary:hover:not(:disabled) {
    transform: scale(1.02);
  }
  .btn-outline {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.75rem 1.25rem;
    border-radius: 999px;
    border: 2px solid #004d34;
    background: transparent;
    color: #004d34;
    font-weight: 700;
    font-size: 0.88rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
    cursor: pointer;
    text-decoration: none;
    transition: background 0.15s ease;
  }
  .btn-outline:hover {
    background: rgba(0, 77, 52, 0.06);
  }
  .link-muted {
    background: none;
    border: none;
    padding: 0;
    font-size: 0.82rem;
    font-weight: 600;
    color: #6f7a72;
    cursor: pointer;
    text-decoration: none;
    font-family: inherit;
  }
  .link-muted:hover {
    color: #775a19;
  }

  /* Aside summary */
  .summary-card {
    border: 1px solid rgba(190, 201, 193, 0.35);
  }
  .sum-row {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 0.5rem;
    padding: 0.5rem 0;
    border-bottom: 1px dashed rgba(190, 201, 193, 0.4);
    font-size: 0.88rem;
  }
  .sum-row:last-of-type {
    border-bottom: none;
  }
  .total-row {
    padding-top: 0.75rem;
    padding-bottom: 0.75rem;
  }
  .sum-label {
    color: #6f7a72;
    font-size: 0.82rem;
  }
  .sum-val {
    font-weight: 700;
    font-size: 0.88rem;
    color: #1b1c1c;
    text-align: right;
  }
  .mono {
    font-family: ui-monospace, monospace;
    font-size: 0.82rem;
    letter-spacing: 0.04em;
    color: #004d34;
  }
  .sum-amt {
    font-weight: 800;
    font-size: 1.1rem;
    color: #775a19;
    font-family: 'Plus Jakarta Sans', sans-serif;
    text-align: right;
  }
  .sum-status {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-top: 1rem;
    padding: 0.55rem 0.85rem;
    border-radius: 0.5rem;
    font-size: 0.78rem;
    font-weight: 700;
  }
  .sum-status .material-symbols-outlined {
    font-size: 1rem;
  }
  .paid-badge {
    background: rgba(0, 103, 71, 0.1);
    color: #006747;
  }
  .partial-badge {
    background: rgba(0, 103, 71, 0.06);
    color: #006747;
  }
  .waiting-badge {
    background: rgba(119, 90, 25, 0.1);
    color: #775a19;
  }
  .expired-badge {
    background: rgba(147, 0, 10, 0.08);
    color: #93000a;
  }
  .sum-note {
    display: flex;
    gap: 0.5rem;
    align-items: flex-start;
    margin-top: 1.25rem;
    padding-top: 1rem;
    border-top: 1px solid rgba(190, 201, 193, 0.35);
    font-size: 0.75rem;
    color: #6f7a72;
    line-height: 1.45;
  }
  .sum-note .material-symbols-outlined {
    font-size: 0.95rem;
    color: #006747;
    flex-shrink: 0;
    margin-top: 0.05rem;
  }
  .sum-note p {
    margin: 0;
  }
  .tips-card {
    border: 1px solid rgba(190, 201, 193, 0.35);
  }
  .tips {
    margin: 0;
    padding-left: 1.1rem;
    font-size: 0.82rem;
    color: #3f4943;
    line-height: 1.55;
  }
  .tips li {
    margin-bottom: 0.45rem;
  }
  .tips li:last-child {
    margin-bottom: 0;
  }
</style>
