<script lang="ts">
  interface ReplicaSite {
    id: string;
    agent_name: string;
    subdomain: string;
    clicks: number;
    leads: number;
    last_updated: string;
    active: boolean;
  }

  let sites = $state<ReplicaSite[]>([
    { id: 's1', agent_name: 'PT Mitra Barokah Tours', subdomain: 'mitrabarokah', clicks: 1_240, leads: 38, last_updated: '2026-04-20', active: true },
    { id: 's2', agent_name: 'CV Cahaya Umroh', subdomain: 'cahayaumroh', clicks: 875, leads: 22, last_updated: '2026-04-18', active: true },
    { id: 's3', agent_name: 'Karima Travel Group', subdomain: 'karimatravel', clicks: 2_150, leads: 61, last_updated: '2026-04-21', active: true },
    { id: 's4', agent_name: 'Ustadz Fahmi Network', subdomain: 'fahminetwork', clicks: 340, leads: 8, last_updated: '2026-04-15', active: true },
    { id: 's5', agent_name: 'PT Andalus Wisata', subdomain: 'andaluswisata', clicks: 0, leads: 0, last_updated: '2026-04-12', active: false },
    { id: 's6', agent_name: 'Naufal Berkah Tour', subdomain: 'naufalberkahtour', clicks: 110, leads: 3, last_updated: '2026-04-19', active: true }
  ]);

  const utmStats = $derived({
    total_clicks: sites.reduce((s, r) => s + r.clicks, 0),
    total_leads: sites.reduce((s, r) => s + r.leads, 0),
    active_sites: sites.filter(r => r.active).length,
    avg_cvr: (() => {
      const active = sites.filter(r => r.active && r.clicks > 0);
      if (!active.length) return 0;
      return (active.reduce((s, r) => s + r.leads / r.clicks, 0) / active.length * 100).toFixed(1);
    })()
  });

  let editModal = $state(false);
  let editSite = $state<ReplicaSite | null>(null);
  let editSubdomain = $state('');

  function openEdit(site: ReplicaSite) {
    editSite = site;
    editSubdomain = site.subdomain;
    editModal = true;
  }

  function saveEdit() {
    if (!editSite) return;
    sites = sites.map(s => s.id === editSite!.id ? { ...s, subdomain: editSubdomain } : s);
    editModal = false;
    editSite = null;
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
      <span class="topbar-current">Replica Site</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
      <button class="avatar">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Replica Site Manager</h2>
        <p>Kelola landing page personal setiap agen dan pantau performa UTM</p>
      </div>
    </div>

    <!-- UTM Stats -->
    <div class="stats-row">
      <div class="stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:#2563eb">mouse</span>
        <div>
          <div class="stat-value">{utmStats.total_clicks.toLocaleString('id-ID')}</div>
          <div class="stat-label">Total Klik</div>
        </div>
      </div>
      <div class="stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:#059669">person_add</span>
        <div>
          <div class="stat-value">{utmStats.total_leads}</div>
          <div class="stat-label">Total Leads dari Replica</div>
        </div>
      </div>
      <div class="stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:#7c3aed">language</span>
        <div>
          <div class="stat-value">{utmStats.active_sites}</div>
          <div class="stat-label">Site Aktif</div>
        </div>
      </div>
      <div class="stat-card">
        <span class="material-symbols-outlined stat-icon" style="color:#d97706">trending_up</span>
        <div>
          <div class="stat-value">{utmStats.avg_cvr}%</div>
          <div class="stat-label">Avg. Konversi</div>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="section-block">
      <div class="section-title">
        <span class="material-symbols-outlined">language</span>
        Daftar Replica Site Agen
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Nama Agen</th>
              <th>URL Replica Site</th>
              <th>Status</th>
              <th>Klik</th>
              <th>Leads</th>
              <th>CVR</th>
              <th>Update Terakhir</th>
              <th class="align-right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each sites as site (site.id)}
              <tr>
                <td><span class="agent-name">{site.agent_name}</span></td>
                <td>
                  <a
                    href="https://{site.subdomain}.umrohos.id"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="url-link"
                  >{site.subdomain}.umrohos.id</a>
                </td>
                <td>
                  <span class="site-status" class:active={site.active} class:inactive={!site.active}>
                    {site.active ? 'Aktif' : 'Nonaktif'}
                  </span>
                </td>
                <td>{site.clicks.toLocaleString('id-ID')}</td>
                <td>{site.leads}</td>
                <td>
                  {site.clicks > 0 ? (site.leads / site.clicks * 100).toFixed(1) + '%' : '—'}
                </td>
                <td class="text-muted">{fmtDate(site.last_updated)}</td>
                <td class="align-right actions-cell">
                  <a
                    href="https://{site.subdomain}.umrohos.id"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="action-btn"
                  >
                    <span class="material-symbols-outlined">open_in_new</span>
                    Preview
                  </a>
                  <button class="action-btn action-btn--edit" onclick={() => openEdit(site)}>
                    <span class="material-symbols-outlined">edit</span>
                    Edit
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <div class="table-footer">Total {sites.length} site terdaftar</div>
    </div>
  </section>
</main>

<!-- Edit Modal -->
{#if editModal && editSite}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div class="modal-backdrop" onclick={() => { editModal = false; }}></div>
  <div class="modal" role="dialog" aria-labelledby="edit-modal-title" aria-modal="true">
    <div class="modal-header">
      <h3 id="edit-modal-title">Edit Replica Site</h3>
      <button class="modal-close" onclick={() => { editModal = false; }}>
        <span class="material-symbols-outlined">close</span>
      </button>
    </div>
    <div class="modal-body">
      <div class="field-group">
        <label class="field-label" for="subdomain-input">Subdomain</label>
        <div class="subdomain-wrap">
          <input id="subdomain-input" type="text" bind:value={editSubdomain} placeholder="namaagen" />
          <span class="subdomain-suffix">.umrohos.id</span>
        </div>
      </div>
    </div>
    <div class="modal-footer">
      <button class="ghost-btn" onclick={() => { editModal = false; }}>Batal</button>
      <button class="primary-btn" onclick={saveEdit}>Simpan</button>
    </div>
  </div>
{/if}

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

  .stats-row { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 0.85rem; margin-bottom: 1.25rem; }
  .stat-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; padding: 0.9rem 1rem; display: flex; align-items: center; gap: 0.75rem; }
  .stat-icon { font-size: 1.6rem; flex-shrink: 0; }
  .stat-value { font-size: 1.3rem; font-weight: 700; color: #191c1e; line-height: 1.1; }
  .stat-label { font-size: 0.68rem; color: #434655; margin-top: 0.1rem; }

  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.4rem; overflow: hidden; }
  .section-title { display: flex; align-items: center; gap: 0.5rem; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); font-size: 0.82rem; font-weight: 700; color: #191c1e; }
  .section-title .material-symbols-outlined { font-size: 1rem; color: #2563eb; }

  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); white-space: nowrap; }
  th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; font-weight: 700; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .align-right { text-align: right; }
  .text-muted { color: #737686; font-size: 0.75rem; }
  .agent-name { font-weight: 600; color: #191c1e; }

  .url-link { color: #2563eb; text-decoration: none; font-size: 0.75rem; }
  .url-link:hover { text-decoration: underline; }

  .site-status { display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; }
  .site-status.active { background: #d1fae5; color: #065f46; }
  .site-status.inactive { background: #f1f5f9; color: #64748b; }

  .actions-cell { display: flex; justify-content: flex-end; gap: 0.35rem; align-items: center; }
  .action-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.3rem 0.6rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.2rem; background: #fff; font-size: 0.72rem; font-weight: 600; color: #191c1e; cursor: pointer; text-decoration: none; }
  .action-btn:hover { background: #f2f4f6; }
  .action-btn .material-symbols-outlined { font-size: 0.9rem; }
  .action-btn--edit:hover { border-color: #2563eb; color: #2563eb; }

  .table-footer { padding: 0.55rem 0.85rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; font-size: 0.68rem; color: #434655; }

  /* Modal */
  .modal-backdrop { position: fixed; inset: 0; background: rgb(0 0 0 / 0.35); z-index: 50; }
  .modal { position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%); z-index: 51; width: min(440px, calc(100vw - 2rem)); background: #fff; border-radius: 0.4rem; border: 1px solid rgb(195 198 215 / 0.55); box-shadow: 0 8px 24px rgb(0 0 0 / 0.12); }
  .modal-header { display: flex; align-items: center; justify-content: space-between; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); background: #f2f4f6; }
  .modal-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .modal-close { border: 0; background: transparent; cursor: pointer; color: #434655; display: grid; place-items: center; border-radius: 0.2rem; padding: 0.2rem; }
  .modal-close:hover { background: #e6e8ea; }
  .modal-body { padding: 1rem; }
  .field-group { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .subdomain-wrap { display: flex; align-items: center; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; overflow: hidden; }
  .subdomain-wrap input { flex: 1; border: 0; padding: 0.5rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; outline: none; }
  .subdomain-suffix { padding: 0.5rem 0.65rem; background: #f2f4f6; font-size: 0.78rem; color: #434655; border-left: 1px solid rgb(195 198 215 / 0.55); white-space: nowrap; }
  .modal-footer { padding: 0.75rem 1rem; border-top: 1px solid rgb(195 198 215 / 0.45); display: flex; justify-content: flex-end; gap: 0.55rem; }
  .ghost-btn, .primary-btn { border-radius: 0.25rem; padding: 0.5rem 0.85rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; border: 1px solid rgb(195 198 215 / 0.55); }
  .ghost-btn { background: #fff; color: #191c1e; }
  .primary-btn { border-color: #2563eb; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; }
</style>
