<script lang="ts">
  type AgentStatus = 'registered' | 'kyc' | 'mou' | 'active';

  interface AgentRow {
    id: string;
    name: string;
    email: string;
    phone: string;
    status: AgentStatus;
    registered_at: string;
    kyc_docs: number;
  }

  const STATUS_LABELS: Record<AgentStatus, string> = {
    registered: 'Terdaftar',
    kyc: 'KYC Diajukan',
    mou: 'MoU Dikirim',
    active: 'Aktif'
  };

  let statusFilter = $state<AgentStatus | 'all'>('all');
  let searchQuery = $state('');

  let agents = $state<AgentRow[]>([
    { id: 'a1', name: 'PT Mitra Barokah Tours', email: 'admin@mitrabarokah.co.id', phone: '0811-2233-4455', status: 'active', registered_at: '2026-01-15', kyc_docs: 5 },
    { id: 'a2', name: 'CV Cahaya Umroh', email: 'info@cahayaumroh.com', phone: '0821-3344-5566', status: 'active', registered_at: '2026-02-03', kyc_docs: 5 },
    { id: 'a3', name: 'Ustadz Fahmi Network', email: 'fahmi@fahminet.id', phone: '0812-9900-1122', status: 'kyc', registered_at: '2026-04-10', kyc_docs: 3 },
    { id: 'a4', name: 'PT Andalus Wisata', email: 'ops@andaluswisata.com', phone: '0855-6677-8899', status: 'mou', registered_at: '2026-04-12', kyc_docs: 5 },
    { id: 'a5', name: 'Rina Hadiati Agency', email: 'rina.hadiati@gmail.com', phone: '0878-1234-5678', status: 'registered', registered_at: '2026-04-20', kyc_docs: 0 },
    { id: 'a6', name: 'Karima Travel Group', email: 'support@karimatravel.id', phone: '0819-9988-7766', status: 'active', registered_at: '2025-11-22', kyc_docs: 5 },
    { id: 'a7', name: 'Naufal Berkah Tour', email: 'naufal@berkah.id', phone: '0831-2200-4400', status: 'kyc', registered_at: '2026-04-18', kyc_docs: 2 },
    { id: 'a8', name: 'PT Baitullah Mandiri', email: 'contact@baitullahmandiri.com', phone: '0812-5544-3322', status: 'registered', registered_at: '2026-04-22', kyc_docs: 0 }
  ]);

  const filteredAgents = $derived(
    agents.filter(a => {
      const matchStatus = statusFilter === 'all' || a.status === statusFilter;
      const q = searchQuery.toLowerCase();
      const matchSearch = !q || a.name.toLowerCase().includes(q) || a.email.toLowerCase().includes(q);
      return matchStatus && matchSearch;
    })
  );

  const counts = $derived({
    all: agents.length,
    registered: agents.filter(a => a.status === 'registered').length,
    kyc: agents.filter(a => a.status === 'kyc').length,
    mou: agents.filter(a => a.status === 'mou').length,
    active: agents.filter(a => a.status === 'active').length
  });

  function approveKyc(id: string) {
    agents = agents.map(a => a.id === id ? { ...a, status: 'mou' as AgentStatus } : a);
  }

  function signMou(id: string) {
    agents = agents.map(a => a.id === id ? { ...a, status: 'active' as AgentStatus } : a);
  }

  function fmtDate(iso: string): string {
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }

  const TABS: Array<{ value: AgentStatus | 'all'; label: string }> = [
    { value: 'all', label: 'Semua' },
    { value: 'registered', label: 'Terdaftar' },
    { value: 'kyc', label: 'KYC' },
    { value: 'mou', label: 'MoU' },
    { value: 'active', label: 'Aktif' }
  ];
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
      <span class="topbar-current">Onboarding</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Pipeline Onboarding Agen</h2>
        <p>Kelola registrasi, proses KYC, dan penandatanganan MoU agen baru</p>
      </div>
      <button class="primary-btn">
        <span class="material-symbols-outlined">person_add</span>
        Tambah Agen
      </button>
    </div>

    <!-- Status Tabs -->
    <div class="tab-row">
      {#each TABS as tab}
        <button
          class="tab-btn"
          class:active={statusFilter === tab.value}
          onclick={() => { statusFilter = tab.value; }}
        >
          {tab.label}
          <span class="tab-count">{tab.value === 'all' ? counts.all : counts[tab.value as AgentStatus]}</span>
        </button>
      {/each}
    </div>

    <!-- Search -->
    <div class="search-row">
      <div class="search-wrap">
        <span class="material-symbols-outlined search-icon">search</span>
        <input type="text" placeholder="Cari nama atau email agen..." bind:value={searchQuery} />
      </div>
    </div>

    <!-- Table -->
    <div class="section-block">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Nama Agen</th>
              <th>Email</th>
              <th>Telepon</th>
              <th>Status</th>
              <th>Dokumen KYC</th>
              <th>Terdaftar</th>
              <th class="align-right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredAgents as agent (agent.id)}
              <tr>
                <td><span class="agent-name">{agent.name}</span></td>
                <td class="text-muted">{agent.email}</td>
                <td class="text-muted">{agent.phone}</td>
                <td>
                  <span class="status-badge status-{agent.status}">{STATUS_LABELS[agent.status]}</span>
                </td>
                <td>
                  <div class="doc-bar">
                    <div class="doc-progress" style="width:{(agent.kyc_docs/5)*100}%"></div>
                  </div>
                  <span class="doc-count">{agent.kyc_docs}/5 dokumen</span>
                </td>
                <td class="text-muted">{fmtDate(agent.registered_at)}</td>
                <td class="align-right actions-cell">
                  {#if agent.status === 'kyc'}
                    <button class="action-btn action-btn--approve" onclick={() => approveKyc(agent.id)}>
                      <span class="material-symbols-outlined">verified</span>
                      Approve KYC
                    </button>
                  {:else if agent.status === 'mou'}
                    <button class="action-btn action-btn--mou" onclick={() => signMou(agent.id)}>
                      <span class="material-symbols-outlined">handshake</span>
                      Sign MoU
                    </button>
                  {:else if agent.status === 'registered'}
                    <button class="action-btn">
                      <span class="material-symbols-outlined">upload_file</span>
                      Minta KYC
                    </button>
                  {:else}
                    <span class="text-muted" style="font-size:0.72rem">—</span>
                  {/if}
                </td>
              </tr>
            {/each}
            {#if filteredAgents.length === 0}
              <tr>
                <td colspan="7" class="empty-cell">
                  <span class="material-symbols-outlined">inbox</span>
                  <p>Tidak ada agen yang sesuai filter.</p>
                </td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
      <div class="table-footer">Menampilkan {filteredAgents.length} dari {agents.length} agen</div>
    </div>
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
  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.25rem; flex-wrap: wrap; gap: 0.75rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; line-height: 1.2; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .primary-btn { display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.5rem 1rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.3rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  .tab-row { display: flex; gap: 0.25rem; margin-bottom: 1rem; flex-wrap: wrap; }
  .tab-btn { display: inline-flex; align-items: center; gap: 0.35rem; padding: 0.4rem 0.8rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; background: #fff; font-size: 0.78rem; color: #434655; cursor: pointer; }
  .tab-btn:hover { background: #f2f4f6; }
  .tab-btn.active { border-color: #2563eb; color: #004ac6; background: #eff6ff; font-weight: 700; }
  .tab-count { font-size: 0.65rem; background: #e2e8f0; color: #334155; padding: 0.1rem 0.35rem; border-radius: 999px; font-weight: 700; }
  .tab-btn.active .tab-count { background: #bfdbfe; color: #1e40af; }

  .search-row { margin-bottom: 1rem; }
  .search-wrap { position: relative; max-width: 28rem; }
  .search-icon { position: absolute; left: 0.65rem; top: 50%; transform: translateY(-50%); font-size: 1rem; color: #737686; }
  .search-wrap input { width: 100%; border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.7rem 0.48rem 2.1rem; font-size: 0.85rem; color: #191c1e; }

  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .align-right { text-align: right; }
  .text-muted { color: #737686; font-size: 0.75rem; }
  .agent-name { font-weight: 600; color: #191c1e; }

  .status-badge { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .status-registered { background: #f1f5f9; color: #334155; }
  .status-kyc { background: #fef3c7; color: #92400e; }
  .status-mou { background: #ede9fe; color: #4c1d95; }
  .status-active { background: #d1fae5; color: #065f46; }

  .doc-bar { height: 4px; background: #e2e8f0; border-radius: 999px; width: 80px; margin-bottom: 0.2rem; }
  .doc-progress { height: 100%; background: #059669; border-radius: 999px; transition: width 0.3s; }
  .doc-count { font-size: 0.62rem; color: #737686; }

  .actions-cell { text-align: right; }
  .action-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.3rem 0.6rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.72rem; font-weight: 600; color: #191c1e; cursor: pointer; }
  .action-btn:hover { background: #f2f4f6; }
  .action-btn .material-symbols-outlined { font-size: 0.9rem; }
  .action-btn--approve { border-color: #059669; color: #065f46; }
  .action-btn--approve:hover { background: #d1fae5; }
  .action-btn--mou { border-color: #7c3aed; color: #4c1d95; }
  .action-btn--mou:hover { background: #ede9fe; }

  .empty-cell { text-align: center; padding: 3rem 1rem; color: #b0b3c1; }
  .empty-cell .material-symbols-outlined { font-size: 2rem; display: block; margin: 0 auto 0.5rem; }
  .empty-cell p { margin: 0; font-size: 0.82rem; }

  .table-footer { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; font-size: 0.68rem; color: #434655; }
</style>
