<script lang="ts">
  type TeamSection = 'downline' | 'tier' | 'zakat' | 'charity';

  const SECTIONS: Array<{ id: TeamSection; label: string; icon: string }> = [
    { id: 'downline', label: 'Struktur Downline', icon: 'account_tree' },
    { id: 'tier', label: 'Tier Leveling', icon: 'military_tech' },
    { id: 'zakat', label: 'Kalkulator Zakat', icon: 'calculate' },
    { id: 'charity', label: 'Infaq & Amal Pagi', icon: 'volunteer_activism' }
  ];

  let activeSection = $state<TeamSection>('downline');

  // --- Downline ---
  const downlineTree = $state([
    { id: 'a1', name: 'PT Mitra Barokah Tours', tier: 'Gold', sponsor: '—', level: 1, sales_ytd: 320_000_000, downline_count: 5 },
    { id: 'a2', name: 'CV Cahaya Umroh', tier: 'Silver', sponsor: 'PT Mitra Barokah', level: 2, sales_ytd: 215_000_000, downline_count: 3 },
    { id: 'a3', name: 'Karima Travel Group', tier: 'Gold', sponsor: '—', level: 1, sales_ytd: 290_000_000, downline_count: 4 },
    { id: 'a4', name: 'Ustadz Fahmi Network', tier: 'Silver', sponsor: 'Karima Travel', level: 2, sales_ytd: 145_000_000, downline_count: 2 },
    { id: 'a5', name: 'Naufal Berkah Tour', tier: 'Bronze', sponsor: 'CV Cahaya Umroh', level: 3, sales_ytd: 65_000_000, downline_count: 1 },
    { id: 'a6', name: 'Rina Hadiati Agency', tier: 'Bronze', sponsor: 'Ustadz Fahmi', level: 3, sales_ytd: 42_000_000, downline_count: 0 },
    { id: 'a7', name: 'PT Andalus Wisata', tier: 'Bronze', sponsor: 'Karima Travel', level: 2, sales_ytd: 28_000_000, downline_count: 0 },
    { id: 'a8', name: 'Ikhsan Santri Agency', tier: 'Bronze', sponsor: 'Naufal Berkah Tour', level: 4, sales_ytd: 18_000_000, downline_count: 0 }
  ]);

  let downlineSearch = $state('');
  const filteredDownline = $derived(
    !downlineSearch.trim()
      ? downlineTree
      : downlineTree.filter(d => d.name.toLowerCase().includes(downlineSearch.toLowerCase()))
  );

  // --- Tier Leveling ---
  const tierRules = $state([
    { tier: 'Bronze', min_sales: 0, max_sales: 99_999_999, commission_pct: 5, override_pct: 0.5, requirements: 'Terdaftar & Aktif' },
    { tier: 'Silver', min_sales: 100_000_000, max_sales: 299_999_999, commission_pct: 5, override_pct: 1.5, requirements: 'Omzet ≥ Rp 100 jt / kuartal' },
    { tier: 'Gold', min_sales: 300_000_000, max_sales: 999_999_999, commission_pct: 5, override_pct: 2.5, requirements: 'Omzet ≥ Rp 300 jt / kuartal' },
    { tier: 'Platinum', min_sales: 1_000_000_000, max_sales: Infinity, commission_pct: 5, override_pct: 3.5, requirements: 'Omzet ≥ Rp 1 M / kuartal' }
  ]);

  const agentTierStatus = $state([
    { name: 'PT Mitra Barokah Tours', current_tier: 'Gold', sales_q: 320_000_000, next_tier: 'Platinum', needed: 680_000_000 },
    { name: 'Karima Travel Group', current_tier: 'Gold', sales_q: 290_000_000, next_tier: 'Platinum', needed: 710_000_000 },
    { name: 'CV Cahaya Umroh', current_tier: 'Silver', sales_q: 215_000_000, next_tier: 'Gold', needed: 85_000_000 },
    { name: 'Ustadz Fahmi Network', current_tier: 'Silver', sales_q: 145_000_000, next_tier: 'Gold', needed: 155_000_000 }
  ]);

  // --- Zakat Calculator ---
  let zakatIncome = $state(0);
  let zakatPeriod = $state<'monthly' | 'yearly'>('yearly');
  const NISAB_GOLD = 85 * 950_000; // ~85 gram gold x harga estimasi

  const zakatCalc = $derived((() => {
    const annual = zakatPeriod === 'monthly' ? zakatIncome * 12 : zakatIncome;
    const isWajib = annual >= NISAB_GOLD;
    const amount = isWajib ? annual * 0.025 : 0;
    return { annual, isWajib, amount };
  })());

  // --- Charity / Infaq ---
  let infaqAmount = $state(10_000);
  let infaqNote = $state('');

  const charityHistory = $state([
    { id: 'ch1', agent: 'PT Mitra Barokah Tours', amount: 500_000, type: 'Infaq Pagi', note: 'Sedekah Senin', date: '2026-04-21' },
    { id: 'ch2', agent: 'CV Cahaya Umroh', amount: 250_000, type: 'Infaq Pagi', note: 'Pagi Selasa', date: '2026-04-22' },
    { id: 'ch3', agent: 'Karima Travel Group', amount: 1_000_000, type: 'Zakat Maal', note: 'Zakat April 2026', date: '2026-04-20' },
    { id: 'ch4', agent: 'Ustadz Fahmi', amount: 150_000, type: 'Infaq Pagi', note: '', date: '2026-04-23' }
  ]);

  function submitInfaq() {
    if (!infaqAmount) return;
    charityHistory.unshift({
      id: 'ch' + Date.now(),
      agent: 'Agen Saat Ini',
      amount: infaqAmount,
      type: 'Infaq Pagi',
      note: infaqNote,
      date: new Date().toISOString().slice(0, 10)
    });
    infaqNote = '';
  }

  function fmtRp(v: number): string {
    if (v >= 1_000_000) return 'Rp ' + (v / 1_000_000).toFixed(1) + ' jt';
    if (v >= 1_000) return 'Rp ' + (v / 1_000).toFixed(0) + ' rb';
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
      <span class="topbar-current">Tim & Downline</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Tim & Downline</h2>
        <p>Struktur downline, tier leveling, kalkulator zakat, dan amal pagi</p>
      </div>
    </div>

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

    <!-- Downline -->
    {#if activeSection === 'downline'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title"><span class="material-symbols-outlined">account_tree</span>Struktur Downline</div>
          <div class="search-wrap">
            <span class="material-symbols-outlined search-icon">search</span>
            <input type="text" placeholder="Cari agen..." bind:value={downlineSearch} />
          </div>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Nama Agen</th>
                <th>Tier</th>
                <th>Sponsor / Upline</th>
                <th>Level</th>
                <th>Sales YTD</th>
                <th>Downline Langsung</th>
              </tr>
            </thead>
            <tbody>
              {#each filteredDownline as agent (agent.id)}
                <tr>
                  <td>
                    <div class="agent-cell" style="padding-left:{(agent.level - 1) * 1.2}rem">
                      {#if agent.level > 1}
                        <span class="material-symbols-outlined indent-icon">subdirectory_arrow_right</span>
                      {/if}
                      <span class="font-semibold">{agent.name}</span>
                    </div>
                  </td>
                  <td><span class="tier-badge tier-{agent.tier.toLowerCase()}">{agent.tier}</span></td>
                  <td class="text-muted">{agent.sponsor}</td>
                  <td><span class="level-chip">L{agent.level}</span></td>
                  <td>{fmtRp(agent.sales_ytd)}</td>
                  <td>{agent.downline_count} agen</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="table-footer">Total {downlineTree.length} agen dalam jaringan</div>
      </div>

    <!-- Tier Leveling -->
    {:else if activeSection === 'tier'}
      <div class="two-col">
        <div class="section-block">
          <div class="section-title"><span class="material-symbols-outlined">military_tech</span>Aturan Tier</div>
          <div class="tier-table">
            {#each tierRules as rule}
              <div class="tier-row">
                <div class="tier-name-cell">
                  <span class="tier-badge tier-{rule.tier.toLowerCase()}">{rule.tier}</span>
                </div>
                <div class="tier-details">
                  <div class="tier-req">{rule.requirements}</div>
                  <div class="tier-comm">Komisi {rule.commission_pct}% + Override {rule.override_pct}%</div>
                </div>
              </div>
            {/each}
          </div>
        </div>
        <div class="section-block">
          <div class="section-title"><span class="material-symbols-outlined">trending_up</span>Status Tier Agen</div>
          <div class="tier-status-list">
            {#each agentTierStatus as ag}
              <div class="tier-status-card">
                <div class="tier-status-top">
                  <span class="font-semibold">{ag.name}</span>
                  <span class="tier-badge tier-{ag.current_tier.toLowerCase()}">{ag.current_tier}</span>
                </div>
                <div class="tier-progress-wrap">
                  <div class="tier-progress-bar">
                    <div
                      class="tier-progress-fill"
                      style="width:{Math.min(ag.sales_q / (ag.sales_q + ag.needed) * 100, 100)}%"
                    ></div>
                  </div>
                </div>
                <div class="tier-status-bottom">
                  <span class="tier-sales">{fmtRp(ag.sales_q)} dicapai</span>
                  <span class="tier-needed text-muted">butuh {fmtRp(ag.needed)} lagi → {ag.next_tier}</span>
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>

    <!-- Zakat Calculator -->
    {:else if activeSection === 'zakat'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">calculate</span>Kalkulator Zakat Maal Agen</div>
        <div class="zakat-layout">
          <div class="zakat-form">
            <div class="info-banner">
              <span class="material-symbols-outlined">info</span>
              Nisab zakat maal ≈ 85 gram emas (~Rp {(NISAB_GOLD / 1_000_000).toFixed(0)} jt)
            </div>
            <div class="field-row">
              <label class="field-label">Periode</label>
              <div class="toggle-group">
                <button
                  class="toggle-btn"
                  class:active={zakatPeriod === 'monthly'}
                  onclick={() => { zakatPeriod = 'monthly'; }}
                >Bulanan</button>
                <button
                  class="toggle-btn"
                  class:active={zakatPeriod === 'yearly'}
                  onclick={() => { zakatPeriod = 'yearly'; }}
                >Tahunan</button>
              </div>
            </div>
            <div class="field-row">
              <label class="field-label">
                Penghasilan {zakatPeriod === 'monthly' ? 'Bulanan' : 'Tahunan'} (Rp)
              </label>
              <input type="number" class="field-input" bind:value={zakatIncome} placeholder="0" />
            </div>
          </div>
          <div class="zakat-result">
            <div class="zakat-result-card" class:wajib={zakatCalc.isWajib} class:not-wajib={!zakatCalc.isWajib}>
              <div class="zakat-result-icon">
                <span class="material-symbols-outlined">
                  {zakatCalc.isWajib ? 'check_circle' : 'cancel'}
                </span>
              </div>
              <div class="zakat-result-main">
                <div class="zakat-status-text">
                  {zakatCalc.isWajib ? 'Wajib Zakat' : 'Belum Wajib Zakat'}
                </div>
                <div class="zakat-annual">Penghasilan tahunan: {fmtRp(zakatCalc.annual)}</div>
              </div>
              {#if zakatCalc.isWajib}
                <div class="zakat-amount-wrap">
                  <div class="zakat-amount">{fmtRp(zakatCalc.amount)}</div>
                  <div class="zakat-amount-label">Zakat yang harus dibayar (2.5%)</div>
                </div>
              {/if}
            </div>
            <div class="zakat-note">
              Catatan: Kalkulator ini hanya estimasi. Konsultasikan dengan amil atau ulama untuk penghitungan yang lebih akurat.
            </div>
          </div>
        </div>
      </div>

    <!-- Charity / Infaq -->
    {:else if activeSection === 'charity'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">volunteer_activism</span>Infaq & Amal Pagi</div>
        <div class="charity-layout">
          <div class="charity-form">
            <div class="morning-quote">
              <span class="material-symbols-outlined">format_quote</span>
              <p>"Sedekah tidak mengurangi harta." — HR. Muslim</p>
            </div>
            <div class="field-row">
              <label class="field-label">Nominal Infaq (Rp)</label>
              <div class="quick-amounts">
                {#each [10_000, 25_000, 50_000, 100_000] as amt}
                  <button
                    class="quick-amt-btn"
                    class:active={infaqAmount === amt}
                    onclick={() => { infaqAmount = amt; }}
                  >{fmtRp(amt)}</button>
                {/each}
              </div>
              <input type="number" class="field-input" bind:value={infaqAmount} />
            </div>
            <div class="field-row">
              <label class="field-label">Catatan (opsional)</label>
              <input type="text" class="field-input" bind:value={infaqNote} placeholder="Niat atau keterangan..." />
            </div>
            <button class="charity-btn" onclick={submitInfaq}>
              <span class="material-symbols-outlined">volunteer_activism</span>
              Catat Infaq Hari Ini
            </button>
          </div>
          <div class="charity-history">
            <div class="sub-title">Riwayat Infaq & Amal</div>
            <div class="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Agen</th>
                    <th>Nominal</th>
                    <th>Jenis</th>
                    <th>Catatan</th>
                    <th>Tanggal</th>
                  </tr>
                </thead>
                <tbody>
                  {#each charityHistory as ch (ch.id)}
                    <tr>
                      <td class="font-semibold">{ch.agent}</td>
                      <td class="amount-positive">{fmtRp(ch.amount)}</td>
                      <td><span class="charity-type">{ch.type}</span></td>
                      <td class="text-muted">{ch.note || '—'}</td>
                      <td class="text-muted">{fmtDate(ch.date)}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
            <div class="charity-total">
              Total infaq & amal terkumpul:
              <strong>{fmtRp(charityHistory.reduce((s, c) => s + c.amount, 0))}</strong>
            </div>
          </div>
        </div>
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

  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; overflow: hidden; margin-bottom: 1rem; }
  .section-header { display: flex; align-items: center; justify-content: space-between; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); flex-wrap: wrap; gap: 0.5rem; }
  .section-title { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; font-weight: 700; color: #191c1e; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .section-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }
  .section-header .section-title { padding: 0; border: 0; }

  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .text-muted { color: #737686; font-size: 0.75rem; }
  .font-semibold { font-weight: 600; color: #191c1e; }
  .table-footer { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; font-size: 0.68rem; color: #434655; }

  /* Search */
  .search-wrap { position: relative; }
  .search-icon { position: absolute; left: 0.65rem; top: 50%; transform: translateY(-50%); font-size: 1rem; color: #737686; }
  .search-wrap input { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.4rem 0.7rem 0.4rem 2.1rem; font-size: 0.82rem; color: #191c1e; min-width: 14rem; }

  /* Downline */
  .agent-cell { display: flex; align-items: center; gap: 0.35rem; }
  .indent-icon { font-size: 0.9rem; color: #b0b3c1; flex-shrink: 0; }
  .tier-badge { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .tier-gold { background: #fef3c7; color: #92400e; }
  .tier-silver { background: #f1f5f9; color: #334155; }
  .tier-bronze { background: #fde8d8; color: #7c2d12; }
  .tier-platinum { background: #cffafe; color: #0e7490; }
  .level-chip { display: inline-flex; padding: 0.12rem 0.4rem; background: #eff6ff; color: #1d4ed8; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }

  /* Tier */
  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  @media (max-width: 900px) { .two-col { grid-template-columns: 1fr; } }
  .tier-table { display: flex; flex-direction: column; }
  .tier-row { display: flex; align-items: center; gap: 0.85rem; padding: 0.75rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .tier-row:last-child { border-bottom: 0; }
  .tier-name-cell { min-width: 70px; }
  .tier-details { flex: 1; }
  .tier-req { font-size: 0.78rem; font-weight: 600; color: #191c1e; }
  .tier-comm { font-size: 0.68rem; color: #737686; margin-top: 0.1rem; }
  .tier-status-list { display: flex; flex-direction: column; gap: 0; }
  .tier-status-card { padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .tier-status-card:last-child { border-bottom: 0; }
  .tier-status-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.4rem; }
  .tier-progress-wrap { margin-bottom: 0.3rem; }
  .tier-progress-bar { height: 6px; background: #e2e8f0; border-radius: 999px; }
  .tier-progress-fill { height: 100%; background: linear-gradient(90deg, #2563eb, #059669); border-radius: 999px; transition: width 0.3s; }
  .tier-status-bottom { display: flex; justify-content: space-between; font-size: 0.68rem; }
  .tier-sales { color: #059669; font-weight: 700; }
  .tier-needed { color: #737686; }

  /* Zakat */
  .zakat-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .zakat-layout { grid-template-columns: 1fr; } }
  .zakat-form { display: flex; flex-direction: column; gap: 0.75rem; }
  .info-banner { display: flex; align-items: center; gap: 0.4rem; background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 0.3rem; padding: 0.55rem 0.75rem; font-size: 0.75rem; color: #1e40af; }
  .info-banner .material-symbols-outlined { font-size: 1rem; flex-shrink: 0; }
  .field-row { display: flex; flex-direction: column; gap: 0.2rem; }
  .field-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field-input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; }
  .toggle-group { display: flex; gap: 0.25rem; }
  .toggle-btn { padding: 0.35rem 0.75rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; background: #fff; font-size: 0.75rem; cursor: pointer; color: #434655; }
  .toggle-btn.active { border-color: #2563eb; background: #eff6ff; color: #004ac6; font-weight: 700; }

  .zakat-result { display: flex; flex-direction: column; gap: 0.75rem; }
  .zakat-result-card { border: 2px solid #e2e8f0; border-radius: 0.5rem; padding: 1.25rem; display: flex; flex-direction: column; align-items: center; gap: 0.5rem; text-align: center; }
  .zakat-result-card.wajib { border-color: #059669; background: #f0fdf4; }
  .zakat-result-card.not-wajib { border-color: #e2e8f0; background: #f8fafc; }
  .zakat-result-icon .material-symbols-outlined { font-size: 2.5rem; }
  .wajib .zakat-result-icon .material-symbols-outlined { color: #059669; }
  .not-wajib .zakat-result-icon .material-symbols-outlined { color: #94a3b8; }
  .zakat-status-text { font-size: 1rem; font-weight: 700; color: #191c1e; }
  .zakat-annual { font-size: 0.75rem; color: #434655; }
  .zakat-amount-wrap { margin-top: 0.35rem; }
  .zakat-amount { font-size: 1.5rem; font-weight: 700; color: #059669; }
  .zakat-amount-label { font-size: 0.68rem; color: #434655; margin-top: 0.15rem; }
  .zakat-note { font-size: 0.68rem; color: #737686; background: #f7f9fb; border: 1px solid rgb(195 198 215 / 0.35); border-radius: 0.3rem; padding: 0.55rem 0.7rem; line-height: 1.4; }

  /* Charity */
  .charity-layout { display: grid; grid-template-columns: 1fr 1.5fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .charity-layout { grid-template-columns: 1fr; } }
  .charity-form { display: flex; flex-direction: column; gap: 0.75rem; }
  .morning-quote { background: #fef9c3; border: 1px solid #fde68a; border-radius: 0.35rem; padding: 0.75rem; display: flex; gap: 0.4rem; align-items: flex-start; }
  .morning-quote .material-symbols-outlined { font-size: 1.2rem; color: #92400e; flex-shrink: 0; margin-top: 0.1rem; }
  .morning-quote p { margin: 0; font-size: 0.78rem; color: #78350f; font-style: italic; line-height: 1.4; }
  .quick-amounts { display: flex; gap: 0.3rem; flex-wrap: wrap; margin-bottom: 0.3rem; }
  .quick-amt-btn { padding: 0.3rem 0.6rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.72rem; cursor: pointer; color: #434655; }
  .quick-amt-btn:hover { background: #f2f4f6; }
  .quick-amt-btn.active { border-color: #059669; background: #d1fae5; color: #065f46; font-weight: 700; }
  .charity-btn { display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.55rem 1rem; background: linear-gradient(90deg, #065f46, #059669); color: #fff; border: none; border-radius: 0.3rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; }
  .charity-btn .material-symbols-outlined { font-size: 1rem; }
  .charity-history { display: flex; flex-direction: column; gap: 0.5rem; }
  .sub-title { font-size: 0.78rem; font-weight: 700; color: #191c1e; }
  .charity-type { display: inline-flex; padding: 0.12rem 0.4rem; background: #d1fae5; color: #065f46; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .amount-positive { color: #059669; font-weight: 700; }
  .charity-total { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; font-size: 0.75rem; color: #434655; }
  .charity-total strong { color: #059669; }
</style>
