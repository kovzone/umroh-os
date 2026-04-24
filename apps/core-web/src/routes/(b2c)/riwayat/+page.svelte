<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  // In production this would come from the authenticated user's booking service
  const transactions = [
    {
      id: 'TRX-2024-001',
      bookingCode: 'BKG-20241205-001',
      packageName: 'Paket Gold 9 Hari',
      departureDate: '15 Februari 2025',
      amount: 34200000,
      type: 'DP',
      method: 'Transfer Bank (BCA)',
      date: '5 Desember 2024',
      status: 'confirmed'
    },
    {
      id: 'TRX-2024-002',
      bookingCode: 'BKG-20241205-001',
      packageName: 'Paket Gold 9 Hari',
      departureDate: '15 Februari 2025',
      amount: 17100000,
      type: 'Cicilan 1',
      method: 'Transfer Bank (Mandiri)',
      date: '5 Januari 2025',
      status: 'confirmed'
    },
    {
      id: 'TRX-2024-003',
      bookingCode: 'BKG-20241205-001',
      packageName: 'Paket Gold 9 Hari',
      departureDate: '15 Februari 2025',
      amount: 17100000,
      type: 'Cicilan 2',
      method: 'Virtual Account BNI',
      date: '5 Februari 2025',
      status: 'pending'
    }
  ];

  const totalPaid = $derived(
    transactions.filter(t => t.status === 'confirmed').reduce((acc, t) => acc + t.amount, 0)
  );
  const totalAmount = $derived(transactions.reduce((acc, t) => acc + t.amount, 0));
  const progressPercent = $derived(totalAmount > 0 ? (totalPaid / totalAmount) * 100 : 0);

  function formatRp(val: number): string {
    return 'Rp ' + new Intl.NumberFormat('id-ID').format(val);
  }

  function statusLabel(status: string): string {
    if (status === 'confirmed') return 'Dikonfirmasi';
    if (status === 'pending') return 'Menunggu';
    if (status === 'failed') return 'Gagal';
    return status;
  }
</script>

<svelte:head>
  <title>Riwayat Pembayaran — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" ctaLabel="Lihat Paket" packagesLinkActive={false}>
  <div class="history-root">
    <div class="shell">
      <div class="page-header">
        <div>
          <h1>Riwayat Pembayaran</h1>
          <p>Pantau status pembayaran dan cicilan perjalanan umrah Anda</p>
        </div>
        <a class="btn-wa" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
          <span class="material-symbols-outlined">chat</span>
          Hubungi Admin
        </a>
      </div>

      <!-- Summary Card -->
      <div class="summary-card">
        <div class="summary-info">
          <div class="summary-item">
            <span class="summary-label">Paket</span>
            <strong>Paket Gold 9 Hari — 15 Februari 2025</strong>
          </div>
          <div class="summary-item">
            <span class="summary-label">Total Biaya</span>
            <strong class="total">{formatRp(totalAmount)}</strong>
          </div>
          <div class="summary-item">
            <span class="summary-label">Sudah Dibayar</span>
            <strong class="paid">{formatRp(totalPaid)}</strong>
          </div>
          <div class="summary-item">
            <span class="summary-label">Sisa Tagihan</span>
            <strong class="remaining">{formatRp(totalAmount - totalPaid)}</strong>
          </div>
        </div>
        <div class="progress-section">
          <div class="progress-header">
            <span>Progres Pembayaran</span>
            <span>{progressPercent.toFixed(0)}%</span>
          </div>
          <div class="progress-bar-wrap">
            <div class="progress-bar" style="width: {progressPercent}%"></div>
          </div>
        </div>
      </div>

      <!-- Transactions -->
      <div class="transactions-section">
        <h2>Detail Transaksi</h2>
        <div class="transactions-list">
          {#each transactions as txn (txn.id)}
            <div class="txn-card">
              <div class="txn-icon" class:icon-confirmed={txn.status === 'confirmed'} class:icon-pending={txn.status === 'pending'}>
                <span class="material-symbols-outlined">
                  {txn.status === 'confirmed' ? 'check_circle' : txn.status === 'pending' ? 'schedule' : 'cancel'}
                </span>
              </div>
              <div class="txn-main">
                <div class="txn-top">
                  <div>
                    <h3 class="txn-type">{txn.type}</h3>
                    <p class="txn-method">{txn.method}</p>
                    <p class="txn-date">{txn.date}</p>
                  </div>
                  <div class="txn-right">
                    <strong class="txn-amount">{formatRp(txn.amount)}</strong>
                    <span class="txn-status" class:status-confirmed={txn.status === 'confirmed'} class:status-pending={txn.status === 'pending'}>
                      {statusLabel(txn.status)}
                    </span>
                  </div>
                </div>
                <p class="txn-ref">Ref: {txn.id} · Booking: {txn.bookingCode}</p>
              </div>
            </div>
          {/each}
        </div>
      </div>

      <!-- Info Box -->
      <div class="info-box">
        <span class="material-symbols-outlined info-icon">info</span>
        <div>
          <p>Konfirmasi pembayaran biasanya diproses dalam 1×24 jam kerja. Jika pembayaran Anda belum terkonfirmasi setelah 24 jam, hubungi tim kami dengan menyertakan bukti transfer.</p>
          <a href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">Hubungi Admin →</a>
        </div>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .history-root {
    padding-top: calc(5.2rem + 2rem);
    padding-bottom: 5rem;
    background: #fbf9f8;
    min-height: 100vh;
  }
  .shell {
    max-width: 72rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  .page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }
  .page-header h1 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .page-header p {
    margin: 0.4rem 0 0;
    color: #6b7280;
  }
  .btn-wa {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    background: #006747;
    color: #fff;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.7rem 1.4rem;
    font-size: 0.88rem;
    flex-shrink: 0;
  }
  .btn-wa .material-symbols-outlined { font-size: 1rem; }
  /* Summary */
  .summary-card {
    background: #fff;
    border-radius: 1.5rem;
    padding: 2rem;
    border: 1px solid rgba(190,201,193,0.3);
    margin-bottom: 2rem;
    box-shadow: 0 4px 16px rgba(0,0,0,0.04);
  }
  .summary-info {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr 1fr;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
  }
  .summary-label {
    display: block;
    font-size: 0.76rem;
    font-weight: 700;
    color: #9ca3af;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    margin-bottom: 0.4rem;
  }
  .summary-item strong {
    font-size: 0.95rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .summary-item .total { color: #004d34; font-size: 1.1rem; }
  .summary-item .paid { color: #006747; }
  .summary-item .remaining { color: #c2410c; }
  .progress-header {
    display: flex;
    justify-content: space-between;
    font-size: 0.82rem;
    color: #6b7280;
    margin-bottom: 0.5rem;
  }
  .progress-bar-wrap {
    height: 10px;
    background: #e4e2e2;
    border-radius: 999px;
    overflow: hidden;
  }
  .progress-bar {
    height: 100%;
    background: linear-gradient(90deg, #004d34, #22c55e);
    border-radius: 999px;
    transition: width 0.5s ease;
  }
  /* Transactions */
  .transactions-section h2 {
    margin: 0 0 1.2rem;
    font-size: 1.15rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .transactions-list {
    display: grid;
    gap: 1rem;
    margin-bottom: 2rem;
  }
  .txn-card {
    display: flex;
    gap: 1rem;
    background: #fff;
    border-radius: 1.25rem;
    padding: 1.2rem 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    align-items: flex-start;
  }
  .txn-icon {
    width: 3rem;
    height: 3rem;
    border-radius: 50%;
    background: rgba(190,201,193,0.2);
    display: grid;
    place-items: center;
    flex-shrink: 0;
    color: #9ca3af;
  }
  .txn-icon.icon-confirmed { background: rgba(0,103,71,0.1); color: #006747; }
  .txn-icon.icon-pending { background: rgba(234,88,12,0.1); color: #c2410c; }
  .txn-icon .material-symbols-outlined {
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .txn-main { flex: 1; }
  .txn-top {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
  }
  .txn-type {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .txn-method, .txn-date {
    margin: 0.2rem 0 0;
    font-size: 0.82rem;
    color: #9ca3af;
  }
  .txn-right { text-align: right; flex-shrink: 0; }
  .txn-amount {
    display: block;
    font-size: 1.05rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .txn-status {
    display: inline-block;
    margin-top: 0.35rem;
    font-size: 0.76rem;
    font-weight: 700;
    padding: 0.25rem 0.65rem;
    border-radius: 999px;
    background: rgba(190,201,193,0.2);
    color: #6b7280;
  }
  .txn-status.status-confirmed { background: rgba(0,103,71,0.1); color: #006747; }
  .txn-status.status-pending { background: rgba(234,88,12,0.1); color: #c2410c; }
  .txn-ref {
    margin: 0.6rem 0 0;
    font-size: 0.76rem;
    color: #d1d5db;
    font-family: monospace;
  }
  /* Info box */
  .info-box {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
    background: rgba(0,103,71,0.06);
    border-left: 3px solid #006747;
    border-radius: 0 1rem 1rem 0;
    padding: 1.2rem 1.5rem;
  }
  .info-icon {
    color: #006747;
    font-size: 1.3rem;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
    flex-shrink: 0;
    margin-top: 0.1rem;
  }
  .info-box p {
    margin: 0 0 0.5rem;
    font-size: 0.88rem;
    color: #57534e;
    line-height: 1.65;
  }
  .info-box a {
    font-size: 0.88rem;
    font-weight: 700;
    color: #006747;
    text-decoration: none;
  }
  @media (max-width: 760px) {
    .summary-info { grid-template-columns: 1fr 1fr; }
  }
  @media (max-width: 480px) {
    .summary-info { grid-template-columns: 1fr; }
  }
</style>
