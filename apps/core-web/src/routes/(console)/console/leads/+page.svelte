<script lang="ts">
  import { updateLead } from '$lib/features/s4-crm/crm-api';
  import type { CSUser, Lead, LeadSource, LeadStatus, PackageOption, UpdateLeadRequest } from '$lib/features/s4-crm/types';

  // ---- page data ----
  interface PageData {
    leads: Lead[];
    total: number;
    csUsers: CSUser[];
    packages: PackageOption[];
    error: string | null;
  }

  let { data }: { data: PageData } = $props();

  // ---- local state ----
  // $state() initializers do NOT reference data prop directly — synced via $effect
  let leads = $state<Lead[]>([]);
  let statusFilter = $state<LeadStatus | 'all'>('all');
  let csFilter = $state<string>('all');
  let searchQuery = $state('');

  // Sync data.leads into local state on mount/data change
  $effect(() => {
    leads = data.leads ?? [];
  });

  // ---- modal state ----
  let modalOpen = $state(false);
  let modalLead = $state<Lead | null>(null);
  let modalStatus = $state<LeadStatus>('new');
  let modalNotes = $state('');
  let modalSaving = $state(false);
  let modalError = $state('');

  // ---- constants ----
  const LEAD_STATUSES: LeadStatus[] = ['new', 'contacted', 'qualified', 'converted', 'lost'];

  const STATUS_LABELS: Record<LeadStatus, string> = {
    new: 'Baru',
    contacted: 'Dihubungi',
    qualified: 'Qualified',
    converted: 'Konversi',
    lost: 'Tidak Jadi'
  };

  const SOURCE_LABELS: Record<LeadSource, string> = {
    direct: 'Direct',
    whatsapp: 'WhatsApp',
    instagram: 'Instagram',
    facebook: 'Facebook',
    tiktok: 'TikTok',
    organic: 'Organik',
    agent: 'Agen',
    referral: 'Referral',
    landing_page: 'Landing Page',
    other: 'Lainnya'
  };

  // ---- derived filtered list ----
  const filteredLeads = $derived(
    leads.filter((l) => {
      const matchStatus = statusFilter === 'all' || l.status === statusFilter;
      const matchCs = csFilter === 'all' || l.assigned_cs_id === csFilter;
      const matchSearch =
        !searchQuery ||
        l.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        l.phone.includes(searchQuery);
      return matchStatus && matchCs && matchSearch;
    })
  );

  // ---- summary counts ----
  const summaryCounts = $derived(
    LEAD_STATUSES.reduce(
      (acc, s) => {
        acc[s] = leads.filter((l) => l.status === s).length;
        return acc;
      },
      {} as Record<LeadStatus, number>
    )
  );

  // ---- helpers ----
  function relativeTime(iso: string): string {
    const diff = Date.now() - new Date(iso).getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);
    if (minutes < 1) return 'Baru saja';
    if (minutes < 60) return `${minutes} menit lalu`;
    if (hours < 24) return `${hours} jam lalu`;
    if (days === 1) return 'Kemarin';
    if (days < 7) return `${days} hari lalu`;
    return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }

  function openModal(lead: Lead) {
    modalLead = lead;
    modalStatus = lead.status;
    modalNotes = '';
    modalError = '';
    modalOpen = true;
  }

  function closeModal() {
    modalOpen = false;
    modalLead = null;
    modalError = '';
  }

  async function saveUpdate() {
    if (!modalLead) return;
    modalSaving = true;
    modalError = '';
    try {
      const req: UpdateLeadRequest = {
        status: modalStatus,
        notes: modalNotes.trim() || undefined
      };
      const updated = await updateLead(modalLead.id, req);
      leads = leads.map((l) => (l.id === updated.id ? updated : l));
      closeModal();
    } catch (err) {
      modalError = err instanceof Error ? err.message : 'Gagal menyimpan perubahan.';
    } finally {
      modalSaving = false;
    }
  }
</script>

<main class="leads-shell">
  <header class="topbar">
    <div class="search-wrap">
      <span class="material-symbols-outlined search-icon">search</span>
      <input
        type="text"
        placeholder="Cari nama atau nomor HP..."
        bind:value={searchQuery}
      />
    </div>
    <div class="top-actions">
      <div class="icon-actions">
        <button class="icon-btn notif">
          <span class="material-symbols-outlined">notifications</span>
          <span class="dot"></span>
        </button>
        <button class="avatar" aria-label="Profile">CS</button>
      </div>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Leads</h2>
        <p>Daftar prospek yang masuk — kelola dan update status per lead</p>
      </div>
    </div>

    <!-- Summary chips -->
    <div class="summary-bar">
      {#each LEAD_STATUSES as s}
        <button
          type="button"
          class="summary-chip summary-chip--{s}"
          class:active={statusFilter === s}
          onclick={() => { statusFilter = statusFilter === s ? 'all' : s; }}
        >
          <span class="count">{summaryCounts[s]}</span>
          <span class="label">{STATUS_LABELS[s]}</span>
        </button>
      {/each}
      <span class="total-label">Total: {leads.length} lead</span>
    </div>

    <!-- Filters -->
    <div class="filters-row">
      <div class="filter-group">
        <label for="filter-status">Status</label>
        <select id="filter-status" bind:value={statusFilter}>
          <option value="all">Semua Status</option>
          {#each LEAD_STATUSES as s}
            <option value={s}>{STATUS_LABELS[s]}</option>
          {/each}
        </select>
      </div>

      <div class="filter-group">
        <label for="filter-cs">CS Assigned</label>
        <select id="filter-cs" bind:value={csFilter}>
          <option value="all">Semua CS</option>
          {#each (data.csUsers ?? []) as cs (cs.id)}
            <option value={cs.id}>{cs.name}</option>
          {/each}
        </select>
      </div>

      <div class="filter-group search-filter">
        <label for="search-input">Cari</label>
        <div class="search-inline">
          <span class="material-symbols-outlined">search</span>
          <input
            id="search-input"
            type="text"
            placeholder="Nama / Nomor HP"
            bind:value={searchQuery}
          />
        </div>
      </div>

      {#if statusFilter !== 'all' || csFilter !== 'all' || searchQuery}
        <button
          type="button"
          class="clear-btn"
          onclick={() => { statusFilter = 'all'; csFilter = 'all'; searchQuery = ''; }}
        >
          <span class="material-symbols-outlined">close</span>
          Reset Filter
        </button>
      {/if}
    </div>

    {#if data.error}
      <div class="error-banner">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- Table -->
    <div class="panel">
      <div class="table-wrap">
        {#if filteredLeads.length === 0}
          <div class="empty-state">
            <span class="material-symbols-outlined">person_add</span>
            <p>Tidak ada lead yang sesuai filter.</p>
          </div>
        {:else}
          <table>
            <thead>
              <tr>
                <th>Nama &amp; Nomor HP</th>
                <th>Sumber</th>
                <th>Status</th>
                <th>CS Assigned</th>
                <th>Waktu Masuk</th>
                <th class="align-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each filteredLeads as lead (lead.id)}
                <tr>
                  <!-- Nama + Phone -->
                  <td>
                    <div class="name-cell">
                      <span class="lead-name">{lead.name}</span>
                      <span class="lead-phone">
                        <span class="material-symbols-outlined">phone</span>
                        {lead.phone}
                      </span>
                      {#if lead.email}
                        <span class="lead-email">
                          <span class="material-symbols-outlined">mail</span>
                          {lead.email}
                        </span>
                      {/if}
                    </div>
                  </td>

                  <!-- Source badge -->
                  <td>
                    <span class="source-badge source-badge--{lead.source}">
                      {SOURCE_LABELS[lead.source] ?? lead.source}
                    </span>
                    {#if lead.utm_campaign}
                      <span class="utm-label">{lead.utm_campaign}</span>
                    {/if}
                  </td>

                  <!-- Status badge -->
                  <td>
                    <span class="status-badge status-badge--{lead.status}">
                      {STATUS_LABELS[lead.status]}
                    </span>
                  </td>

                  <!-- CS assigned -->
                  <td>
                    {#if lead.assigned_cs_name}
                      <span class="cs-name">
                        <span class="cs-avatar">{lead.assigned_cs_name.charAt(0)}</span>
                        {lead.assigned_cs_name}
                      </span>
                    {:else}
                      <span class="unassigned">Belum assign</span>
                    {/if}
                  </td>

                  <!-- Relative time -->
                  <td>
                    <span class="time-cell" title={new Date(lead.created_at).toLocaleString('id-ID')}>
                      {relativeTime(lead.created_at)}
                    </span>
                  </td>

                  <!-- Actions -->
                  <td class="actions-cell">
                    <button
                      type="button"
                      class="update-btn"
                      onclick={() => openModal(lead)}
                    >
                      <span class="material-symbols-outlined">edit</span>
                      Update Status
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>
      <div class="table-footer">
        Menampilkan {filteredLeads.length} dari {leads.length} lead
      </div>
    </div>
  </section>
</main>

<!-- Update Status Modal -->
{#if modalOpen && modalLead}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div class="modal-backdrop" onclick={closeModal}></div>
  <div class="modal" role="dialog" aria-labelledby="modal-title" aria-modal="true">
    <div class="modal-header">
      <h3 id="modal-title">Update Status Lead</h3>
      <button type="button" class="modal-close" onclick={closeModal} aria-label="Tutup">
        <span class="material-symbols-outlined">close</span>
      </button>
    </div>

    <div class="modal-body">
      <!-- Lead info summary -->
      <div class="modal-info">
        <div class="modal-lead-avatar">{modalLead.name.charAt(0)}</div>
        <div class="modal-lead-detail">
          <span class="modal-lead-name">{modalLead.name}</span>
          <span class="modal-lead-phone">{modalLead.phone}</span>
        </div>
        <span class="status-badge status-badge--{modalLead.status}">
          {STATUS_LABELS[modalLead.status]}
        </span>
      </div>

      <label class="field-label" for="modal-status">Status Baru</label>
      <select id="modal-status" class="modal-select" bind:value={modalStatus}>
        {#each LEAD_STATUSES as s}
          <option value={s}>{STATUS_LABELS[s]}</option>
        {/each}
      </select>

      <label class="field-label" for="modal-notes">
        Catatan
        <span class="optional">(opsional)</span>
      </label>
      <textarea
        id="modal-notes"
        class="modal-textarea"
        rows="3"
        placeholder="Tambahkan catatan update status..."
        bind:value={modalNotes}
      ></textarea>

      {#if modalError}
        <p class="modal-error">{modalError}</p>
      {/if}
    </div>

    <div class="modal-footer">
      <button type="button" class="ghost-btn" onclick={closeModal} disabled={modalSaving}>
        Batal
      </button>
      <button
        type="button"
        class="primary-btn"
        onclick={saveUpdate}
        disabled={modalSaving}
      >
        {#if modalSaving}
          <span class="material-symbols-outlined spin">progress_activity</span>
          Menyimpan...
        {:else}
          <span class="material-symbols-outlined">save</span>
          Simpan
        {/if}
      </button>
    </div>
  </div>
{/if}

<style>
  .leads-shell {
    min-height: 100vh;
    background: #f7f9fb;
  }

  /* ---- topbar ---- */
  .topbar {
    position: sticky;
    top: 0;
    z-index: 30;
    height: 4rem;
    background: rgb(255 255 255 / 0.9);
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    backdrop-filter: blur(8px);
  }

  .search-wrap {
    flex: 1;
    max-width: 32rem;
    position: relative;
  }

  .search-icon {
    position: absolute;
    left: 0.65rem;
    top: 50%;
    transform: translateY(-50%);
    font-size: 1.05rem;
    color: #737686;
  }

  .search-wrap input {
    width: 100%;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #f2f4f6;
    border-radius: 0.25rem;
    padding: 0.48rem 0.7rem 0.48rem 2.1rem;
    font-size: 0.85rem;
    color: #191c1e;
    outline: none;
  }

  .top-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .icon-actions {
    display: flex;
    align-items: center;
    gap: 0.35rem;
  }

  .icon-btn {
    border: 0;
    background: transparent;
    color: #434655;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    cursor: pointer;
    position: relative;
    display: grid;
    place-items: center;
  }

  .icon-btn:hover { background: #eceef0; }

  .dot {
    position: absolute;
    width: 0.46rem;
    height: 0.46rem;
    border-radius: 999px;
  }

  .notif .dot {
    right: 0.4rem;
    top: 0.35rem;
    background: #ba1a1a;
    border: 2px solid #fff;
  }

  .avatar {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #b4c5ff;
    color: #00174b;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    font-weight: 700;
    font-size: 0.65rem;
    cursor: pointer;
  }

  /* ---- canvas ---- */
  .canvas {
    padding: 1.5rem;
    max-width: 96rem;
  }

  .page-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
  }

  .page-head h2 {
    margin: 0;
    font-size: 1.5rem;
    line-height: 1.2;
  }

  .page-head p {
    margin: 0.3rem 0 0;
    font-size: 0.82rem;
    color: #434655;
  }

  /* ---- summary bar ---- */
  .summary-bar {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .summary-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.35rem 0.65rem;
    border-radius: 0.25rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    cursor: pointer;
    font-size: 0.72rem;
    font-weight: 600;
    background: #fff;
    color: #434655;
    transition: background 0.1s;
  }

  .summary-chip:hover { background: #f2f4f6; }
  .summary-chip.active { border-color: #2563eb; color: #004ac6; background: rgb(37 99 235 / 0.07); }

  .summary-chip .count { font-size: 0.9rem; font-weight: 700; }
  .summary-chip .label { text-transform: capitalize; }

  .total-label {
    margin-left: auto;
    font-size: 0.72rem;
    color: #434655;
  }

  /* ---- filters ---- */
  .filters-row {
    display: flex;
    align-items: flex-end;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .filter-group {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
  }

  .filter-group label {
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }

  .filter-group select,
  .filter-group input {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    border-radius: 0.25rem;
    padding: 0.42rem 0.6rem;
    font-size: 0.82rem;
    color: #191c1e;
    min-width: 10rem;
    outline: none;
  }

  .search-filter .search-inline {
    display: flex;
    align-items: center;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    background: #fff;
    padding: 0 0.6rem;
  }

  .search-filter .search-inline span { font-size: 0.9rem; color: #737686; flex-shrink: 0; }

  .search-filter .search-inline input {
    border: 0;
    background: transparent;
    padding: 0.42rem 0.4rem;
    min-width: 14rem;
    font-size: 0.82rem;
    outline: none;
  }

  .clear-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.42rem 0.6rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    background: #fff;
    font-size: 0.78rem;
    color: #434655;
    cursor: pointer;
    align-self: flex-end;
  }

  .clear-btn:hover { background: #f2f4f6; }
  .clear-btn .material-symbols-outlined { font-size: 0.9rem; }

  /* ---- error banner ---- */
  .error-banner {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: #ffdad6;
    color: #93000a;
    border-radius: 0.25rem;
    padding: 0.65rem 0.85rem;
    font-size: 0.82rem;
    margin-bottom: 1rem;
  }

  /* ---- table panel ---- */
  .panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .table-wrap { overflow-x: auto; }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    padding: 0.62rem 0.85rem;
    text-align: left;
    font-size: 0.78rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    white-space: nowrap;
  }

  th {
    text-transform: uppercase;
    font-size: 0.62rem;
    letter-spacing: 0.08em;
    color: #434655;
    background: #f2f4f6;
  }

  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }

  .align-right { text-align: right; }

  /* ---- name cell ---- */
  .name-cell {
    display: flex;
    flex-direction: column;
    gap: 0.08rem;
  }

  .lead-name {
    font-weight: 700;
    font-size: 0.82rem;
    color: #191c1e;
  }

  .lead-phone,
  .lead-email {
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
    font-size: 0.72rem;
    color: #434655;
  }

  .lead-phone .material-symbols-outlined,
  .lead-email .material-symbols-outlined {
    font-size: 0.78rem;
  }

  /* ---- source badge ---- */
  .source-badge {
    display: inline-flex;
    padding: 0.12rem 0.4rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
    background: #e6e8ea;
    color: #434655;
  }

  .source-badge--whatsapp { background: #d1fae5; color: #065f46; }
  .source-badge--instagram { background: #fde8ff; color: #6b00a3; }
  .source-badge--facebook { background: #dbeafe; color: #1e3a8a; }
  .source-badge--tiktok { background: #f1f5f9; color: #0f172a; }
  .source-badge--organic { background: #fef9c3; color: #7d5f00; }
  .source-badge--agent { background: #ede9fe; color: #4c1d95; }
  .source-badge--referral { background: #fee2e2; color: #991b1b; }
  .source-badge--landing_page { background: #e0f2fe; color: #075985; }

  .utm-label {
    display: block;
    font-size: 0.62rem;
    color: #737686;
    margin-top: 0.15rem;
  }

  /* ---- status badge ---- */
  .status-badge {
    display: inline-flex;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
  }

  .status-badge--new { background: #dbeafe; color: #1e3a8a; }
  .status-badge--contacted { background: #fef9c3; color: #7d5f00; }
  .status-badge--qualified { background: #ede9fe; color: #4c1d95; }
  .status-badge--converted { background: #d1fae5; color: #065f46; }
  .status-badge--lost { background: #fee2e2; color: #991b1b; }

  /* ---- CS name ---- */
  .cs-name {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.75rem;
    color: #191c1e;
  }

  .cs-avatar {
    width: 1.5rem;
    height: 1.5rem;
    border-radius: 999px;
    background: #b4c5ff;
    color: #00174b;
    font-size: 0.62rem;
    font-weight: 700;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }

  .unassigned {
    font-size: 0.72rem;
    color: #b0b3c1;
    font-style: italic;
  }

  /* ---- time cell ---- */
  .time-cell {
    font-size: 0.72rem;
    color: #434655;
  }

  /* ---- actions ---- */
  .actions-cell { text-align: right; }

  .update-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.3rem 0.6rem;
    border-radius: 0.2rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    font-size: 0.72rem;
    font-weight: 600;
    color: #191c1e;
    cursor: pointer;
  }

  .update-btn:hover { background: #f2f4f6; }
  .update-btn .material-symbols-outlined { font-size: 0.9rem; }

  .table-footer {
    padding: 0.55rem 0.85rem;
    font-size: 0.68rem;
    color: #434655;
    border-top: 1px solid rgb(195 198 215 / 0.35);
    background: #f7f9fb;
  }

  /* ---- empty state ---- */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.6rem;
    padding: 3rem 1rem;
    color: #b0b3c1;
  }

  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.82rem; }

  /* ---- modal ---- */
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgb(0 0 0 / 0.35);
    z-index: 50;
  }

  .modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 51;
    width: min(480px, calc(100vw - 2rem));
    background: #fff;
    border-radius: 0.4rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    box-shadow: 0 8px 24px rgb(0 0 0 / 0.12);
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.85rem 1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    background: #f2f4f6;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 700;
  }

  .modal-close {
    border: 0;
    background: transparent;
    cursor: pointer;
    color: #434655;
    display: grid;
    place-items: center;
    border-radius: 0.2rem;
    padding: 0.2rem;
  }

  .modal-close:hover { background: #e6e8ea; }

  .modal-body {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .modal-info {
    display: flex;
    align-items: center;
    gap: 0.7rem;
    padding: 0.6rem;
    background: #f2f4f6;
    border-radius: 0.25rem;
  }

  .modal-lead-avatar {
    width: 2rem;
    height: 2rem;
    border-radius: 999px;
    background: #b4c5ff;
    color: #00174b;
    font-size: 0.78rem;
    font-weight: 700;
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }

  .modal-lead-detail {
    display: flex;
    flex-direction: column;
    gap: 0.08rem;
    flex: 1;
  }

  .modal-lead-name {
    font-size: 0.82rem;
    font-weight: 700;
    color: #191c1e;
  }

  .modal-lead-phone {
    font-size: 0.72rem;
    color: #434655;
  }

  .field-label {
    font-size: 0.68rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }

  .optional {
    font-weight: 400;
    text-transform: none;
    letter-spacing: 0;
  }

  .modal-select,
  .modal-textarea {
    border: 1px solid rgb(195 198 215 / 0.55);
    border-radius: 0.25rem;
    padding: 0.5rem 0.65rem;
    font-size: 0.82rem;
    color: #191c1e;
    background: #fff;
    width: 100%;
    font-family: inherit;
    outline: none;
  }

  .modal-textarea { resize: vertical; min-height: 5rem; }

  .modal-error {
    margin: 0;
    font-size: 0.75rem;
    color: #ba1a1a;
    background: #ffdad6;
    padding: 0.45rem 0.65rem;
    border-radius: 0.2rem;
  }

  .modal-footer {
    padding: 0.75rem 1rem;
    border-top: 1px solid rgb(195 198 215 / 0.45);
    display: flex;
    justify-content: flex-end;
    gap: 0.55rem;
  }

  .ghost-btn,
  .primary-btn {
    border-radius: 0.25rem;
    padding: 0.5rem 0.85rem;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    border: 1px solid rgb(195 198 215 / 0.55);
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    font-family: inherit;
  }

  .ghost-btn {
    background: #fff;
    color: #191c1e;
  }

  .ghost-btn:disabled,
  .primary-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .primary-btn {
    border-color: #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
  }

  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
