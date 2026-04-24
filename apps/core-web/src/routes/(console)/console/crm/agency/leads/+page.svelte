<script lang="ts">
  type LeadsSection = 'leads' | 'drip' | 'segments' | 'broadcast' | 'quote' | 'stale';

  const SECTIONS: Array<{ id: LeadsSection; label: string; icon: string }> = [
    { id: 'leads', label: 'Leads & Follow-up', icon: 'person_search' },
    { id: 'drip', label: 'Drip Sequence', icon: 'schedule_send' },
    { id: 'segments', label: 'Segmen', icon: 'group_work' },
    { id: 'broadcast', label: 'Broadcast', icon: 'send' },
    { id: 'quote', label: 'Quote Generator', icon: 'request_quote' },
    { id: 'stale', label: 'Stale Prospect', icon: 'radar' }
  ];

  let activeSection = $state<LeadsSection>('leads');

  // --- Leads ---
  type LeadStatus = 'new' | 'contacted' | 'qualified' | 'proposal' | 'closed_won' | 'closed_lost';

  const LEAD_STATUS_LABELS: Record<LeadStatus, string> = {
    new: 'Baru',
    contacted: 'Dihubungi',
    qualified: 'Kualifed',
    proposal: 'Proposal',
    closed_won: 'Menang',
    closed_lost: 'Gagal'
  };

  const leads = $state([
    { id: 'l1', name: 'Bapak Hendra Wijaya', source: 'Instagram', agent: 'PT Mitra Barokah', status: 'qualified' as LeadStatus, follow_up: '2026-04-25', last_contact: '2026-04-22' },
    { id: 'l2', name: 'Ibu Siti Rahmawati', source: 'WhatsApp Blast', agent: 'CV Cahaya Umroh', status: 'proposal' as LeadStatus, follow_up: '2026-04-26', last_contact: '2026-04-23' },
    { id: 'l3', name: 'Pak Dodi Santoso', source: 'Referral', agent: 'Karima Travel', status: 'new' as LeadStatus, follow_up: '2026-04-24', last_contact: '2026-04-21' },
    { id: 'l4', name: 'Ibu Aisyah Lubis', source: 'Google Ads', agent: 'PT Mitra Barokah', status: 'contacted' as LeadStatus, follow_up: '2026-04-27', last_contact: '2026-04-20' },
    { id: 'l5', name: 'Pak Rudi Hermawan', source: 'Facebook', agent: 'Ustadz Fahmi', status: 'closed_won' as LeadStatus, follow_up: '', last_contact: '2026-04-18' },
    { id: 'l6', name: 'Ibu Dewi Pratiwi', source: 'Website Replica', agent: 'CV Cahaya Umroh', status: 'qualified' as LeadStatus, follow_up: '2026-04-28', last_contact: '2026-04-22' }
  ]);

  const isOverdue = (date: string): boolean => !!date && new Date(date) < new Date();

  // --- Drip Sequence ---
  const dripSequences = $state([
    { id: 'd1', name: 'Welcome Sequence', steps: 5, active_contacts: 34, trigger: 'Registrasi Baru' },
    { id: 'd2', name: 'Nurture — Belum Beli 30H', steps: 8, active_contacts: 18, trigger: '30 hari tanpa transaksi' },
    { id: 'd3', name: 'Re-engagement 90H', steps: 3, active_contacts: 7, trigger: '90 hari tidak aktif' },
    { id: 'd4', name: 'Post-Purchase Care', steps: 4, active_contacts: 12, trigger: 'Booking lunas' }
  ]);

  // --- Segments ---
  let newSegName = $state('');
  let newSegCriteria = $state('');
  const segments = $state([
    { id: 'sg1', name: 'Prospek Ramadan', count: 142, criteria: 'Source: Instagram OR Facebook' },
    { id: 'sg2', name: 'Alumni 2024', count: 88, criteria: 'Pernah berangkat 2024' },
    { id: 'sg3', name: 'Hot Lead — Budget > 30jt', count: 31, criteria: 'Estimasi budget ≥ 30.000.000' }
  ]);

  function addSegment() {
    if (!newSegName.trim()) return;
    segments.push({ id: 'sg' + Date.now(), name: newSegName, count: 0, criteria: newSegCriteria });
    newSegName = '';
    newSegCriteria = '';
  }

  // --- Broadcast ---
  let broadcastMsg = $state('');
  let broadcastTarget = $state('all');
  const broadcasts = $state([
    { id: 'b1', title: 'Promo Ramadan 2026', target: 'Semua Lead', sent: 432, opened: 298, date: '2026-04-10' },
    { id: 'b2', title: 'Reminder Follow-up Hot Lead', target: 'Hot Lead', sent: 31, opened: 28, date: '2026-04-18' },
    { id: 'b3', title: 'Info Jadwal Keberangkatan Mei', target: 'Qualified + Proposal', sent: 85, opened: 61, date: '2026-04-20' }
  ]);

  function sendBroadcast() {
    if (!broadcastMsg.trim()) return;
    broadcasts.unshift({ id: 'b' + Date.now(), title: broadcastMsg, target: broadcastTarget === 'all' ? 'Semua Lead' : broadcastTarget, sent: 0, opened: 0, date: new Date().toISOString().slice(0, 10) });
    broadcastMsg = '';
  }

  // --- Quote Generator ---
  let quotePackage = $state('Paket Gold Ramadan 2026');
  let quoteQty = $state(2);
  let quoteDiscount = $state(5);
  const PACKAGE_PRICES: Record<string, number> = {
    'Paket Gold Ramadan 2026': 32_500_000,
    'Paket Silver 9H 2026': 22_000_000,
    'Paket Bronze 9H Budget': 16_500_000
  };
  const quoteTotal = $derived(
    (() => {
      const base = PACKAGE_PRICES[quotePackage] ?? 25_000_000;
      return base * quoteQty * (1 - quoteDiscount / 100);
    })()
  );

  // --- Stale Prospects ---
  const staleProspects = $state([
    { id: 'sp1', name: 'Pak Arif Budiman', last_contact: '2026-01-15', days_silent: 99, agent: 'CV Cahaya Umroh', status: 'contacted' as LeadStatus },
    { id: 'sp2', name: 'Ibu Nani Suryani', last_contact: '2026-02-03', days_silent: 80, agent: 'PT Mitra Barokah', status: 'qualified' as LeadStatus },
    { id: 'sp3', name: 'Pak Wahyu Setiawan', last_contact: '2026-01-28', days_silent: 86, agent: 'Karima Travel', status: 'proposal' as LeadStatus }
  ]);

  function fmtDate(iso: string): string {
    if (!iso) return '—';
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }

  function fmtRp(v: number): string {
    return 'Rp ' + (v / 1_000_000).toFixed(1) + ' jt';
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
      <span class="topbar-current">CRM Leads</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>CRM Leads Agen</h2>
        <p>Kelola leads, drip sequence, segmen, broadcast, dan quote generator</p>
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

    <!-- Leads & Follow-up -->
    {#if activeSection === 'leads'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title"><span class="material-symbols-outlined">person_search</span>Leads & Follow-up</div>
          <button class="primary-btn"><span class="material-symbols-outlined">add</span>Tambah Lead</button>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Nama Lead</th>
                <th>Sumber</th>
                <th>Agen</th>
                <th>Status</th>
                <th>Follow-up</th>
                <th>Kontak Terakhir</th>
              </tr>
            </thead>
            <tbody>
              {#each leads as lead (lead.id)}
                <tr>
                  <td><span class="font-semibold">{lead.name}</span></td>
                  <td class="text-muted">{lead.source}</td>
                  <td class="text-muted">{lead.agent}</td>
                  <td><span class="lead-status lead-status--{lead.status}">{LEAD_STATUS_LABELS[lead.status]}</span></td>
                  <td>
                    {#if lead.follow_up}
                      <span class="follow-up" class:overdue={isOverdue(lead.follow_up)}>{fmtDate(lead.follow_up)}</span>
                    {:else}
                      <span class="text-muted">—</span>
                    {/if}
                  </td>
                  <td class="text-muted">{fmtDate(lead.last_contact)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="table-footer">Total {leads.length} leads</div>
      </div>

    <!-- Drip Sequence -->
    {:else if activeSection === 'drip'}
      <div class="section-block">
        <div class="section-header">
          <div class="section-title"><span class="material-symbols-outlined">schedule_send</span>Drip Sequence</div>
          <button class="primary-btn"><span class="material-symbols-outlined">add</span>Buat Sequence</button>
        </div>
        <div class="drip-list">
          {#each dripSequences as drip (drip.id)}
            <div class="drip-row">
              <div class="drip-icon"><span class="material-symbols-outlined">schedule_send</span></div>
              <div class="drip-info">
                <div class="drip-name">{drip.name}</div>
                <div class="drip-meta">{drip.steps} langkah &bull; Trigger: {drip.trigger}</div>
              </div>
              <div class="drip-stat">
                <div class="drip-stat-num">{drip.active_contacts}</div>
                <div class="drip-stat-label">Kontak Aktif</div>
              </div>
              <button class="action-btn"><span class="material-symbols-outlined">edit</span>Edit</button>
            </div>
          {/each}
        </div>
      </div>

    <!-- Segments -->
    {:else if activeSection === 'segments'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">group_work</span>Segmen Leads</div>
        <div class="seg-layout">
          <div class="seg-form">
            <div class="sub-title">Buat Segmen Baru</div>
            <div class="field-row">
              <label class="field-label">Nama Segmen</label>
              <input type="text" class="field-input" bind:value={newSegName} placeholder="Contoh: Hot Lead Q2" />
            </div>
            <div class="field-row">
              <label class="field-label">Kriteria</label>
              <input type="text" class="field-input" bind:value={newSegCriteria} placeholder="Contoh: Status = Qualified AND Source = Instagram" />
            </div>
            <button class="primary-btn" onclick={addSegment}>
              <span class="material-symbols-outlined">add</span>
              Simpan Segmen
            </button>
          </div>
          <div class="seg-list">
            <div class="sub-title">Segmen Aktif</div>
            {#each segments as seg (seg.id)}
              <div class="seg-card">
                <div class="seg-card-top">
                  <span class="seg-name">{seg.name}</span>
                  <span class="seg-count">{seg.count} leads</span>
                </div>
                <div class="seg-criteria">{seg.criteria}</div>
              </div>
            {/each}
          </div>
        </div>
      </div>

    <!-- Broadcast -->
    {:else if activeSection === 'broadcast'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">send</span>Broadcast Center</div>
        <div class="broadcast-layout">
          <div class="compose-form">
            <div class="sub-title">Compose Broadcast</div>
            <div class="field-row">
              <label class="field-label">Target Segmen</label>
              <select class="field-input" bind:value={broadcastTarget}>
                <option value="all">Semua Lead</option>
                {#each segments as seg}
                  <option value={seg.name}>{seg.name}</option>
                {/each}
              </select>
            </div>
            <div class="field-row">
              <label class="field-label">Pesan</label>
              <textarea class="field-textarea" rows="4" bind:value={broadcastMsg} placeholder="Tulis pesan broadcast di sini..."></textarea>
            </div>
            <button class="primary-btn" onclick={sendBroadcast}>
              <span class="material-symbols-outlined">send</span>
              Kirim Broadcast
            </button>
          </div>
          <div class="broadcast-history">
            <div class="sub-title">Broadcast Terakhir</div>
            <div class="table-wrap">
              <table>
                <thead>
                  <tr><th>Judul</th><th>Target</th><th>Dikirim</th><th>Dibuka</th><th>Tanggal</th></tr>
                </thead>
                <tbody>
                  {#each broadcasts as b (b.id)}
                    <tr>
                      <td class="font-semibold">{b.title}</td>
                      <td class="text-muted">{b.target}</td>
                      <td>{b.sent}</td>
                      <td>{b.opened}</td>
                      <td class="text-muted">{fmtDate(b.date)}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

    <!-- Quote Generator -->
    {:else if activeSection === 'quote'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">request_quote</span>Quote Generator</div>
        <div class="quote-layout">
          <div class="quote-form">
            <div class="field-row">
              <label class="field-label">Paket Umroh</label>
              <select class="field-input" bind:value={quotePackage}>
                {#each Object.keys(PACKAGE_PRICES) as pkg}
                  <option value={pkg}>{pkg}</option>
                {/each}
              </select>
            </div>
            <div class="field-row">
              <label class="field-label">Jumlah Jamaah</label>
              <input type="number" class="field-input" bind:value={quoteQty} min="1" max="50" />
            </div>
            <div class="field-row">
              <label class="field-label">Diskon Agen (%)</label>
              <input type="number" class="field-input" bind:value={quoteDiscount} min="0" max="20" />
            </div>
            <div class="field-row">
              <label class="field-label">Nama Prospek</label>
              <input type="text" class="field-input" placeholder="Bapak / Ibu ..." />
            </div>
          </div>
          <div class="quote-preview">
            <div class="sub-title">Preview Quotation</div>
            <div class="quote-card">
              <div class="quote-package">{quotePackage}</div>
              <div class="quote-line">
                <span>Harga per Jamaah</span>
                <span>{fmtRp(PACKAGE_PRICES[quotePackage] ?? 0)}</span>
              </div>
              <div class="quote-line">
                <span>Jumlah Jamaah</span>
                <span>{quoteQty} orang</span>
              </div>
              <div class="quote-line">
                <span>Diskon Agen</span>
                <span class="text-green">-{quoteDiscount}%</span>
              </div>
              <div class="quote-divider"></div>
              <div class="quote-total-line">
                <span>Total</span>
                <span class="quote-total">{fmtRp(quoteTotal)}</span>
              </div>
            </div>
            <button class="primary-btn" style="margin-top:0.75rem">
              <span class="material-symbols-outlined">download</span>
              Unduh PDF Quote
            </button>
          </div>
        </div>
      </div>

    <!-- Stale Prospects -->
    {:else if activeSection === 'stale'}
      <div class="section-block">
        <div class="section-title"><span class="material-symbols-outlined">radar</span>Stale Prospect Radar</div>
        <div class="stale-info">
          <span class="material-symbols-outlined">info</span>
          Prospek yang tidak ada kontak lebih dari 60 hari
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Nama</th>
                <th>Agen</th>
                <th>Status</th>
                <th>Kontak Terakhir</th>
                <th>Hari Diam</th>
                <th class="align-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each staleProspects as sp (sp.id)}
                <tr>
                  <td class="font-semibold">{sp.name}</td>
                  <td class="text-muted">{sp.agent}</td>
                  <td><span class="lead-status lead-status--{sp.status}">{LEAD_STATUS_LABELS[sp.status]}</span></td>
                  <td class="text-muted">{fmtDate(sp.last_contact)}</td>
                  <td><span class="days-badge">{sp.days_silent} hari</span></td>
                  <td class="align-right">
                    <button class="action-btn action-btn--re">
                      <span class="material-symbols-outlined">replay</span>
                      Re-engage
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="table-footer">Total {staleProspects.length} stale prospect</div>
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
  .section-header { display: flex; align-items: center; justify-content: space-between; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .section-title { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; font-weight: 700; color: #191c1e; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .section-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }
  .section-header .section-title { padding: 0; border: 0; }

  .primary-btn { display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.5rem 1rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.3rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .align-right { text-align: right; }
  .text-muted { color: #737686; font-size: 0.75rem; }
  .font-semibold { font-weight: 600; color: #191c1e; }
  .table-footer { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; font-size: 0.68rem; color: #434655; }

  .lead-status { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .lead-status--new { background: #f1f5f9; color: #334155; }
  .lead-status--contacted { background: #dbeafe; color: #1e40af; }
  .lead-status--qualified { background: #fef3c7; color: #92400e; }
  .lead-status--proposal { background: #ede9fe; color: #4c1d95; }
  .lead-status--closed_won { background: #d1fae5; color: #065f46; }
  .lead-status--closed_lost { background: #fee2e2; color: #991b1b; }

  .follow-up { font-size: 0.75rem; color: #191c1e; }
  .follow-up.overdue { color: #dc2626; font-weight: 700; }

  /* Drip */
  .drip-list { display: flex; flex-direction: column; }
  .drip-row { display: flex; align-items: center; gap: 0.85rem; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .drip-row:last-child { border-bottom: 0; }
  .drip-icon { width: 2.4rem; height: 2.4rem; background: #dbeafe; border-radius: 0.35rem; display: grid; place-items: center; flex-shrink: 0; }
  .drip-icon .material-symbols-outlined { font-size: 1.2rem; color: #1e40af; }
  .drip-info { flex: 1; }
  .drip-name { font-weight: 700; font-size: 0.82rem; color: #191c1e; }
  .drip-meta { font-size: 0.68rem; color: #737686; margin-top: 0.15rem; }
  .drip-stat { text-align: right; min-width: 80px; }
  .drip-stat-num { font-size: 1.1rem; font-weight: 700; color: #191c1e; }
  .drip-stat-label { font-size: 0.62rem; color: #737686; }
  .action-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.3rem 0.6rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.72rem; font-weight: 600; color: #191c1e; cursor: pointer; }
  .action-btn:hover { background: #f2f4f6; }
  .action-btn .material-symbols-outlined { font-size: 0.9rem; }
  .action-btn--re { border-color: #2563eb; color: #004ac6; }
  .action-btn--re:hover { background: #eff6ff; }

  /* Segments */
  .seg-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .seg-layout { grid-template-columns: 1fr; } }
  .seg-form { display: flex; flex-direction: column; gap: 0.65rem; }
  .seg-list { display: flex; flex-direction: column; gap: 0.6rem; }
  .sub-title { font-size: 0.78rem; font-weight: 700; color: #191c1e; margin-bottom: 0.5rem; }
  .field-row { display: flex; flex-direction: column; gap: 0.2rem; }
  .field-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field-input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; }
  .seg-card { border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.3rem; padding: 0.65rem 0.8rem; }
  .seg-card-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.25rem; }
  .seg-name { font-weight: 700; font-size: 0.82rem; color: #191c1e; }
  .seg-count { font-size: 0.68rem; color: #437655; font-weight: 700; background: #d1fae5; color: #065f46; padding: 0.1rem 0.35rem; border-radius: 0.2rem; }
  .seg-criteria { font-size: 0.68rem; color: #737686; }

  /* Broadcast */
  .broadcast-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .broadcast-layout { grid-template-columns: 1fr; } }
  .compose-form { display: flex; flex-direction: column; gap: 0.65rem; }
  .broadcast-history { display: flex; flex-direction: column; gap: 0.5rem; }
  .field-textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; resize: vertical; }

  /* Quote */
  .quote-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; padding: 1rem; }
  @media (max-width: 700px) { .quote-layout { grid-template-columns: 1fr; } }
  .quote-form { display: flex; flex-direction: column; gap: 0.65rem; }
  .quote-preview { display: flex; flex-direction: column; }
  .quote-card { border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.35rem; padding: 1rem; display: flex; flex-direction: column; gap: 0.5rem; margin-top: 0.5rem; }
  .quote-package { font-weight: 700; font-size: 0.9rem; color: #191c1e; padding-bottom: 0.5rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .quote-line { display: flex; justify-content: space-between; font-size: 0.8rem; color: #434655; }
  .text-green { color: #059669; font-weight: 700; }
  .quote-divider { height: 1px; background: rgb(195 198 215 / 0.45); }
  .quote-total-line { display: flex; justify-content: space-between; font-size: 0.85rem; font-weight: 700; color: #191c1e; }
  .quote-total { color: #2563eb; font-size: 1.1rem; }

  /* Stale */
  .stale-info { display: flex; align-items: center; gap: 0.4rem; background: #fef3c7; padding: 0.6rem 1rem; font-size: 0.78rem; color: #92400e; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .stale-info .material-symbols-outlined { font-size: 1rem; }
  .days-badge { background: #fee2e2; color: #991b1b; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
</style>
