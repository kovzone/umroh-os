<script lang="ts">
  type CommSection = 'overview' | 'history' | 'payout' | 'override' | 'roas' | 'alumni';

  const SECTIONS: Array<{ id: CommSection; label: string; icon: string }> = [
    { id: 'overview', label: 'Saldo & Ringkasan', icon: 'account_balance_wallet' },
    { id: 'history', label: 'Riwayat Transaksi', icon: 'history' },
    { id: 'payout', label: 'Permintaan Payout', icon: 'payments' },
    { id: 'override', label: 'Override Komisi', icon: 'tune' },
    { id: 'roas', label: 'Kalkulator ROAS', icon: 'calculate' },
    { id: 'alumni', label: 'Referral Alumni', icon: 'loyalty' }
  ];

  let activeSection = $state<CommSection>('overview');

  // --- Overview ---
  const balanceCards = $state([
    { label: 'Saldo Tersedia', value: 'Rp 8.450.000', icon: 'account_balance_wallet', color: '#059669', bg: '#d1fae5' },
    { label: 'Komisi Pending', value: 'Rp 4.200.000', icon: 'pending', color: '#d97706', bg: '#fef3c7' },
    { label: 'Total Lifetime', value: 'Rp 847.000.000', icon: 'trending_up', color: '#2563eb', bg: '#dbeafe' },
    { label: 'Agen Aktif Berkontribusi', value: '38', icon: 'storefront', color: '#7c3aed', bg: '#ede9fe' }
  ]);

  // --- Transaction History ---
  const transactions = $state([
    { id: 't1', ref: 'BK-2026-0401', agent: 'PT Mitra Barokah', type: 'Komisi Booking', amount: 1_625_000, status: 'paid', date: '2026-04-20' },
    { id: 't2', ref: 'BK-2026-0388', agent: 'CV Cahaya Umroh', type: 'Komisi Booking', amount: 1_100_000, status: 'paid', date: '2026-04-19' },
    { id: 't3', ref: 'OVR-2026-001', agent: 'Karima Travel Group', type: 'Override Tier Gold', amount: 500_000, status: 'pending', date: '2026-04-21' },
    { id: 't4', ref: 'BK-2026-0395', agent: 'Ustadz Fahmi', type: 'Komisi Booking', amount: 825_000, status: 'paid', date: '2026-04-20' },
    { id: 't5', ref: 'PO-2026-0011', agent: 'PT Mitra Barokah', type: 'Payout', amount: -5_000_000, status: 'disbursed', date: '2026-04-18' }
  ]);

  // --- Payout ---
  let payoutAmount = $state('');
  let payoutBank = $state('BCA');
  let payoutAccount = $state('');
  let payoutNote = $state('');

  const payoutHistory = $state([
    { id: 'po1', agent: 'PT Mitra Barokah', amount: 5_000_000, bank: 'BCA - 1234567890', status: 'disbursed', date: '2026-04-18' },
    { id: 'po2', agent: 'CV Cahaya Umroh', amount: 2_500_000, bank: 'Mandiri - 9876543210', status: 'pending', date: '2026-04-22' }
  ]);

  function submitPayout() {
    if (!payoutAmount || !payoutAccount) return;
    payoutHistory.unshift({
      id: 'po' + Date.now(),
      agent: 'Agen Saat Ini',
      amount: parseFloat(payoutAmount),
      bank: payoutBank + ' - ' + payoutAccount,
      status: 'pending',
      date: new Date().toISOString().slice(0, 10)
    });
    payoutAmount = '';
    payoutAccount = '';
    payoutNote = '';
  }

  // --- Override Commission ---
  const overrideRows = $state([
    { id: 'ov1', agent: 'Karima Travel Group', tier: 'Gold', override_pct: 2.5, base_pct: 5, total_sales: 320_000_000, earned: 8_000_000 },
    { id: 'ov2', agent: 'PT Mitra Barokah', tier: 'Silver', override_pct: 1.5, base_pct: 5, total_sales: 215_000_000, earned: 3_225_000 },
    { id: 'ov3', agent: 'Naufal Berkah Tour', tier: 'Bronze', override_pct: 0.5, base_pct: 5, total_sales: 65_000_000, earned: 325_000 }
  ]);

  // --- ROAS Calculator ---
  let roasAdSpend = $state(3_000_000);
  let roasRevenue = $state(48_000_000);
  const roasValue = $derived(roasRevenue > 0 && roasAdSpend > 0 ? (roasRevenue / roasAdSpend).toFixed(2) : '0.00');
  const roasCpa = $derived(roasAdSpend > 0 ? fmtRp(roasAdSpend / 12) : '—');

  // --- Alumni Referral ---
  const alumniReferrals = $state([
    { id: 'ar1', alumni: 'Ibu Siti Aminah', referred: 'Pak Hendra', package: 'Gold 2026', commission: 500_000, status: 'paid', date: '2026-04-15' },
    { id: 'ar2', alumni: 'Pak Dodi Santoso', referred: 'Ibu Wulan', package: 'Silver 9H', commission: 330_000, status: 'pending', date: '2026-04-20' },
    { id: 'ar3', alumni: 'Ibu Rina Hadiati', referred: 'Pak Wahyu', package: 'Gold 2026', commission: 500_000, status: 'paid', date: '2026-04-12' }
  ]);

  function fmtRp(v: number): string {
    if (v >= 1_000_000) return 'Rp ' + (v / 1_000_000).toFixed(1) + ' jt';
    return 'Rp ' + v.toLocaleString('id-ID');
  }

  function fmtDate(iso: string): string {
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">hub</span>
      <span class="sep">/</span>
      <a href="/console/crm" class="bc-link">CRM</a>
      <span class="sep">/</span>
      <a href="/console/crm/agency" class="bc-link">Portal Agen</a>
      <span class="sep">/</span>
      <span class="topbar-current">Komisi & Wallet</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Komisi & Wallet</h2>
        <p>Saldo komisi, riwayat transaksi, payout, override, ROAS, dan referral alumni</p>
      </div>
    </div>

    <!-- Section Tabs -->
    <div class="section-tabs">
      {#each SECTIONS as sec}
        <button
          class="sec-tab"
          class:active={activeSection === sec.id}
          onclick={() => { activeSection = sec.id; }}
        >
          <span class="material-symbols-outlined">{sec.icon}</span>
          {sec.label}
        </button>
      {/each}
    </div>

    <!-- Overview -->
    {#if activeSection === 'overview'}
      <div class="kpi-grid">
        {#each balanceCards as card}
          <div class="kpi-card">
            <div class="kpi-icon" style="background:{card.bg};color:{card.color}">
              <span class="material-symbols-outlined">{card.icon}</span>
            </div>
            <div>
              <div class="kpi-value">{card.value}</div>
              <div class="kpi-label">{card.label}</div>
            </div>
          </div>
        {/each}
      </div>

    <!-- History -->
    {:else if activeSection === 'history'}
      <div class="section-block">
        <div class="section-title">
          <span class="material-symbols-outlined">history</span>
          Riwayat Transaksi Komisi
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Referensi</th>
                <th>Agen</th>
                <th>Tipe</th>
                <th>Jumlah</th>
                <th>Status</th>
                <th>Tanggal</th>
              </tr>
            </thead>
            <tbody>
              {#each transactions as tx (tx.id)}
                <tr>
                  <td><span class="mono-code">{tx.ref}</span></td>
                  <td class="font-semibold">{tx.agent}</td>
                  <td class="text-muted">{tx.type}</td>
                  <td class:amount-negative={tx.amount < 0} class:amount-positive={tx.amount > 0} class="amount-cell">
                    {tx.amount < 0 ? '-' : '+'}{fmtRp(Math.abs(tx.amount))}
                  </td>
                  <td><span class="tx-status tx-status--{tx.status}">{tx.status === 'paid' ? 'Dibayar' : tx.status === 'pending' ? 'Pending' : 'Dicairkan'}</span></td>
                  <td class="text-muted">{fmtDate(tx.date)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="table-footer">Total {transactions.length} transaksi</div>
      </div>

    <!-- Payout -->
    {:else if activeSection === 'payout'}
      <div class="two-col">
        <div class="section-block">
          <div class="section-title">
            <span class="material-symbols-outlined">payments</span>
            Ajukan Payout
          </div>
          <div class="form-body">
            <div class="field-row">
              <label class="field-label">Jumlah Payout (Rp)</label>
              <input type="number" class="field-input" bind:value={payoutAmount} placeholder="5000000" />
            </div>
            <div class="field-row">
              <label class="field-label">Bank Tujuan</label>
              <select class="field-input" bind:value={payoutBank}>
                <option>BCA</option><option>Mandiri</option><option>BNI</option><option>BRI</option>
              </select>
            </div>
            <div class="field-row">
              <label class="field-label">Nomor Rekening</label>
              <input type="text" class="field-input" bind:value={payoutAccount} placeholder="1234567890" />
            </div>
            <div class="field-row">
              <label class="field-label">Catatan (opsional)</label>
              <input type="text" class="field-input" bind:value={payoutNote} placeholder="Payout bulan April 2026" />
            </div>
            <button class="primary-btn" onclick={submitPayout}>
              <span class="material-symbols-outlined">send</span>
              Ajukan Payout
            </button>
          </div>
        </div>
        <div class="section-block">
          <div class="section-title">
            <span class="material-symbols-outlined">history</span>
            Riwayat Payout
          </div>
          <div class="table-wrap">
            <table>
              <thead>
                <tr><th>Agen</th><th>Jumlah</th><th>Bank</th><th>Status</th><th>Tanggal</th></tr>
              </thead>
              <tbody>
                {#each payoutHistory as po (po.id)}
                  <tr>
                    <td class="font-semibold">{po.agent}</td>
                    <td class="amount-negative">{fmtRp(po.amount)}</td>
                    <td class="text-muted">{po.bank}</td>
                    <td><span class="tx-status tx-status--{po.status}">{po.status === 'disbursed' ? 'Dicairkan' : 'Pending'}</span></td>
                    <td class="text-muted">{fmtDate(po.date)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>

    <!-- Override -->
    {:else if activeSection === 'override'}
      <div class="section-block">
        <div class="section-title">
          <span class="material-symbols-outlined">tune</span>
          Override Komisi per Tier
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Agen</th>
                <th>Tier</th>
                <th>Komisi Dasar</th>
                <th>Override</th>
                <th>Total Penjualan</th>
                <th>Override Diperoleh</th>
              </tr>
            </thead>
            <tbody>
              {#each overrideRows as row (row.id)}
                <tr>
                  <td class="font-semibold">{row.agent}</td>
                  <td><span class="tier-badge tier-{row.tier.toLowerCase()}">{row.tier}</span></td>
                  <td>{row.base_pct}%</td>
                  <td class="text-green">+{row.override_pct}%</td>
                  <td>{fmtRp(row.total_sales)}</td>
                  <td class="amount-positive">{fmtRp(row.earned)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>

    <!-- ROAS Calculator -->
    {:else if activeSection === 'roas'}
      <div class="section-block">
        <div class="section-title">
          <span class="material-symbols-outlined">calculate</span>
          Kalkulator ROAS
        </div>
        <div class="roas-layout">
          <div class="roas-form">
            <div class="field-row">
              <label class="field-label">Biaya Iklan (Rp)</label>
              <input type="number" class="field-input" bind:value={roasAdSpend} />
            </div>
            <div class="field-row">
              <label class="field-label">Pendapatan dari Iklan (Rp)</label>
              <input type="number" class="field-input" bind:value={roasRevenue} />
            </div>
          </div>
          <div class="roas-result">
            <div class="roas-stat">
              <div class="roas-num">{roasValue}x</div>
              <div class="roas-label">ROAS</div>
            </div>
            <div class="roas-stat">
              <div class="roas-num">{roasCpa}</div>
              <div class="roas-label">Est. Biaya per Booking (12 booking)</div>
            </div>
            <div class="roas-rating">
              {#if parseFloat(roasValue) >= 4}
                <span class="roas-good">Sangat Baik (≥ 4x)</span>
              {:else if parseFloat(roasValue) >= 2}
                <span class="roas-ok">Cukup (2–4x)</span>
              {:else}
                <span class="roas-bad">Perlu Ditingkatkan (&lt; 2x)</span>
              {/if}
            </div>
          </div>
        </div>
      </div>

    <!-- Alumni Referral -->
    {:else if activeSection === 'alumni'}
      <div class="section-block">
        <div class="section-title">
          <span class="material-symbols-outlined">loyalty</span>
          Alumni Referral Tracker
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Alumni (Referrer)</th>
                <th>Direferensikan</th>
                <th>Paket</th>
                <th>Komisi Referral</th>
                <th>Status</th>
                <th>Tanggal</th>
              </tr>
            </thead>
            <tbody>
              {#each alumniReferrals as ref (ref.id)}
                <tr>
                  <td class="font-semibold">{ref.alumni}</td>
                  <td>{ref.referred}</td>
                  <td class="text-muted">{ref.package}</td>
                  <td class="amount-positive">{fmtRp(ref.commission)}</td>
                  <td><span class="tx-status tx-status--{ref.status}">{ref.status === 'paid' ? 'Dibayar' : 'Pending'}</span></td>
                  <td class="text-muted">{fmtDate(ref.date)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="table-footer">Total {alumniReferrals.length} referral alumni</div>
      </div>
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; justify-content: space-between; gap: 1rem; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.88rem; color: #434655; }
  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .sep { color: #b0b3c1; }
  .bc-link { color: #004ac6; text-decoration: none; font-weight: 500; }
  .bc-link:hover { text-decoration: underline; }
  .topbar-current { font-weight: 600; color: #191c1e; }
  .top-actions { display: flex; align-items: center; gap: 0.35rem; }
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .section-tabs { display: flex; gap: 0.25rem; flex-wrap: wrap; margin-bottom: 1.25rem; }
  .sec-tab { display: inline-flex; align-items: center; gap: 0.35rem; padding: 0.45rem 0.85rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; background: #fff; font-size: 0.78rem; color: #434655; cursor: pointer; }
  .sec-tab .material-symbols-outlined { font-size: 1rem; }
  .sec-tab:hover { background: #f2f4f6; }
  .sec-tab.active { border-color: #2563eb; color: #004ac6; background: #eff6ff; font-weight: 700; }

  .kpi-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 1rem; margin-bottom: 1rem; }
  .kpi-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; padding: 1rem; display: flex; align-items: center; gap: 0.85rem; }
  .kpi-icon { width: 2.8rem; height: 2.8rem; border-radius: 0.35rem; display: grid; place-items: center; flex-shrink: 0; }
  .kpi-icon .material-symbols-outlined { font-size: 1.4rem; }
  .kpi-value { font-size: 1.25rem; font-weight: 700; color: #191c1e; line-height: 1.1; }
  .kpi-label { font-size: 0.68rem; color: #434655; margin-top: 0.15rem; }

  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; overflow: hidden; margin-bottom: 1rem; }
  .section-title { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; font-weight: 700; color: #191c1e; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .section-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }

  .primary-btn { display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.5rem 1rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.3rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .text-muted { color: #737686; font-size: 0.75rem; }
  .font-semibold { font-weight: 600; color: #191c1e; }
  .mono-code { font-family: monospace; font-size: 0.72rem; color: #004ac6; font-weight: 600; }
  .amount-cell { font-weight: 700; }
  .amount-positive { color: #059669; font-weight: 700; }
  .amount-negative { color: #dc2626; font-weight: 700; }
  .text-green { color: #059669; font-weight: 700; }
  .table-footer { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; font-size: 0.68rem; color: #434655; }

  .tx-status { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .tx-status--paid { background: #d1fae5; color: #065f46; }
  .tx-status--pending { background: #fef3c7; color: #92400e; }
  .tx-status--disbursed { background: #dbeafe; color: #1e40af; }

  .tier-badge { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .tier-gold { background: #fef3c7; color: #92400e; }
  .tier-silver { background: #f1f5f9; color: #334155; }
  .tier-bronze { background: #fde8d8; color: #7c2d12; }

  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  @media (max-width: 900px) { .two-col { grid-template-columns: 1fr; } }

  .form-body { padding: 1rem; display: flex; flex-direction: column; gap: 0.65rem; }
  .field-row { display: flex; flex-direction: column; gap: 0.2rem; }
  .field-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field-input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; }

  /* ROAS */
  .roas-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .roas-layout { grid-template-columns: 1fr; } }
  .roas-form { display: flex; flex-direction: column; gap: 0.65rem; }
  .roas-result { display: flex; flex-direction: column; gap: 1rem; justify-content: center; }
  .roas-stat { text-align: center; }
  .roas-num { font-size: 2.5rem; font-weight: 700; color: #2563eb; line-height: 1.1; }
  .roas-label { font-size: 0.72rem; color: #434655; margin-top: 0.2rem; }
  .roas-rating { text-align: center; }
  .roas-good { display: inline-flex; padding: 0.35rem 0.75rem; background: #d1fae5; color: #065f46; border-radius: 999px; font-size: 0.78rem; font-weight: 700; }
  .roas-ok { display: inline-flex; padding: 0.35rem 0.75rem; background: #fef3c7; color: #92400e; border-radius: 999px; font-size: 0.78rem; font-weight: 700; }
  .roas-bad { display: inline-flex; padding: 0.35rem 0.75rem; background: #fee2e2; color: #991b1b; border-radius: 999px; font-size: 0.78rem; font-weight: 700; }
</style>
