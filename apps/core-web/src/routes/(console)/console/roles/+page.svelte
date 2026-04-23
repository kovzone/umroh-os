<script lang="ts">
  import type { Role, Permission, PermAction, PageData } from './+page.server';

  let { data }: { data: PageData } = $props();

  // ---- local state — synced via $effect ----
  let roles = $state<Role[]>([]);
  let allResources = $state<string[]>([]);

  $effect(() => {
    roles = data.roles ?? [];
    allResources = data.allResources ?? [];
  });

  // ---- selected role ----
  let selectedRoleId = $state<string | null>(null);

  const selectedRole = $derived(
    selectedRoleId ? roles.find((r) => r.id === selectedRoleId) ?? null : null
  );

  // ---- edit state for selected role permissions ----
  // Map<resource, Set<action>>
  let editPerms = $state<Map<string, Set<PermAction>>>(new Map());

  $effect(() => {
    if (selectedRole) {
      const map = new Map<string, Set<PermAction>>();
      for (const p of selectedRole.permissions) {
        map.set(p.resource, new Set(p.actions));
      }
      editPerms = map;
    } else {
      editPerms = new Map();
    }
  });

  // ---- new role form ----
  let showNewForm = $state(false);
  let newRoleName = $state('');
  let newRoleDesc = $state('');
  let newRoleError = $state('');

  function submitNewRole() {
    newRoleError = '';
    if (!newRoleName.trim()) {
      newRoleError = 'Nama role tidak boleh kosong.';
      return;
    }
    const newRole: Role = {
      id: `r${Date.now()}`,
      name: newRoleName.trim().toLowerCase().replace(/\s+/g, '_'),
      description: newRoleDesc.trim(),
      permissions: []
    };
    roles = [...roles, newRole];
    selectedRoleId = newRole.id;
    showNewForm = false;
    newRoleName = '';
    newRoleDesc = '';
  }

  function cancelNewRole() {
    showNewForm = false;
    newRoleName = '';
    newRoleDesc = '';
    newRoleError = '';
  }

  // ---- permission toggle ----
  const ALL_ACTIONS: PermAction[] = ['read', 'write', 'delete', 'export'];

  const ACTION_LABELS: Record<PermAction, string> = {
    read: 'Baca',
    write: 'Tulis',
    delete: 'Hapus',
    export: 'Ekspor'
  };

  function isChecked(resource: string, action: PermAction): boolean {
    return editPerms.get(resource)?.has(action) ?? false;
  }

  function togglePerm(resource: string, action: PermAction) {
    const next = new Map(editPerms);
    const actions = new Set(next.get(resource) ?? []);
    if (actions.has(action)) {
      actions.delete(action);
    } else {
      actions.add(action);
    }
    next.set(resource, actions);
    editPerms = next;
  }

  // ---- save permissions ----
  let saving = $state(false);
  let saveMsg = $state('');

  function savePermissions() {
    if (!selectedRole) return;
    saving = true;
    saveMsg = '';
    // Build updated permissions array
    const newPerms: Permission[] = [];
    for (const resource of allResources) {
      const actions = [...(editPerms.get(resource) ?? [])] as PermAction[];
      if (actions.length > 0) {
        newPerms.push({ resource, actions });
      }
    }
    // Update local state (optimistic — no API in mock mode)
    roles = roles.map((r) =>
      r.id === selectedRole.id ? { ...r, permissions: newPerms } : r
    );
    saving = false;
    saveMsg = 'Perubahan berhasil disimpan.';
    setTimeout(() => { saveMsg = ''; }, 2500);
  }

  // ---- resource label helper ----
  const RESOURCE_LABELS: Record<string, string> = {
    packages: 'Paket',
    departures: 'Keberangkatan',
    bookings: 'Booking',
    jamaah: 'Jamaah',
    payments: 'Pembayaran',
    visa_docs: 'Dokumen Visa',
    ops: 'Operasional',
    finance: 'Keuangan',
    leads: 'Leads',
    catalog: 'Katalog',
    users: 'Pengguna',
    roles: 'Roles',
    audit_log: 'Audit Log'
  };

  function countPerms(role: Role): number {
    return role.permissions.reduce((sum, p) => sum + p.actions.length, 0);
  }
</script>

<main class="page-shell">
  <!-- Topbar -->
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <a href="/console/dashboard" class="back-link">
        <span class="material-symbols-outlined">chevron_left</span>
        Dashboard
      </a>
      <span class="breadcrumb-sep">/</span>
      <span class="topbar-current">Roles &amp; Permissions</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn" title="Notifikasi">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar" aria-label="Profile">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Roles &amp; Permissions</h2>
        <p>Kelola hak akses pengguna berdasarkan role</p>
      </div>
    </div>

    {#if data.error}
      <div class="error-banner" role="alert">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- ================================================================
         Two-panel layout
    ================================================================= -->
    <div class="two-panel">

      <!-- ---- Left panel — Role list ---- -->
      <aside class="panel-left">
        <div class="panel-left-header">
          <span class="panel-title">Daftar Roles</span>
          <button
            type="button"
            class="new-role-btn"
            onclick={() => { showNewForm = true; selectedRoleId = null; }}
          >
            <span class="material-symbols-outlined">add</span>
            Buat Role
          </button>
        </div>

        <!-- New role inline form -->
        {#if showNewForm}
          <div class="new-role-form">
            <div class="form-field">
              <label for="new-role-name">Nama Role</label>
              <input
                id="new-role-name"
                type="text"
                placeholder="mis. supervisor"
                bind:value={newRoleName}
              />
            </div>
            <div class="form-field">
              <label for="new-role-desc">Deskripsi</label>
              <input
                id="new-role-desc"
                type="text"
                placeholder="Deskripsi singkat..."
                bind:value={newRoleDesc}
              />
            </div>
            {#if newRoleError}
              <p class="form-error">{newRoleError}</p>
            {/if}
            <div class="form-actions">
              <button type="button" class="ghost-btn" onclick={cancelNewRole}>Batal</button>
              <button type="button" class="primary-btn" onclick={submitNewRole}>
                <span class="material-symbols-outlined">save</span>
                Simpan
              </button>
            </div>
          </div>
        {/if}

        <!-- Role list -->
        {#if roles.length === 0}
          <div class="empty-inline">
            <span class="material-symbols-outlined">group</span>
            <p>Belum ada role.</p>
          </div>
        {:else}
          <ul class="role-list">
            {#each roles as role (role.id)}
              <li>
                <button
                  type="button"
                  class="role-card"
                  class:selected={selectedRoleId === role.id}
                  onclick={() => { selectedRoleId = role.id; }}
                >
                  <div class="role-card-top">
                    <span class="role-name">{role.name}</span>
                    <span class="perm-count">{countPerms(role)} hak</span>
                  </div>
                  {#if role.description}
                    <p class="role-desc">{role.description}</p>
                  {/if}
                </button>
              </li>
            {/each}
          </ul>
        {/if}
      </aside>

      <!-- ---- Right panel — Role detail ---- -->
      <section class="panel-right">
        {#if !selectedRole}
          <div class="placeholder">
            <span class="material-symbols-outlined">manage_accounts</span>
            <p>Pilih role di kiri untuk melihat detail dan mengatur permissions</p>
          </div>
        {:else}
          <div class="detail-header">
            <div>
              <h3 class="detail-role-name">{selectedRole.name}</h3>
              {#if selectedRole.description}
                <p class="detail-role-desc">{selectedRole.description}</p>
              {/if}
            </div>
            <div class="detail-actions">
              {#if saveMsg}
                <span class="save-success">
                  <span class="material-symbols-outlined">check_circle</span>
                  {saveMsg}
                </span>
              {/if}
              <button
                type="button"
                class="primary-btn"
                onclick={savePermissions}
                disabled={saving}
              >
                {#if saving}
                  <span class="material-symbols-outlined spin">progress_activity</span>
                  Menyimpan...
                {:else}
                  <span class="material-symbols-outlined">save</span>
                  Simpan Perubahan
                {/if}
              </button>
            </div>
          </div>

          <!-- Permissions table -->
          <div class="panel perms-panel">
            <div class="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Resource</th>
                    {#each ALL_ACTIONS as action}
                      <th class="align-center">{ACTION_LABELS[action]}</th>
                    {/each}
                  </tr>
                </thead>
                <tbody>
                  {#each allResources as resource (resource)}
                    <tr>
                      <td>
                        <span class="resource-name">
                          {RESOURCE_LABELS[resource] ?? resource}
                        </span>
                        <span class="resource-key">{resource}</span>
                      </td>
                      {#each ALL_ACTIONS as action}
                        <td class="align-center">
                          <label class="checkbox-wrap" title="{ACTION_LABELS[action]} {resource}">
                            <input
                              type="checkbox"
                              checked={isChecked(resource, action)}
                              onchange={() => togglePerm(resource, action)}
                            />
                            <span class="checkbox-custom"></span>
                          </label>
                        </td>
                      {/each}
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
            <div class="table-footer">
              {allResources.length} resources &mdash;
              {[...editPerms.values()].reduce((sum, s) => sum + s.size, 0)} hak aktif
            </div>
          </div>
        {/if}
      </section>
    </div>
  </section>
</main>

<style>
  .page-shell {
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

  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: #434655;
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.1rem;
    font-size: 0.82rem;
    color: #434655;
    text-decoration: none;
    font-weight: 500;
  }

  .back-link:hover { color: #004ac6; }
  .back-link .material-symbols-outlined { font-size: 1rem; }

  .breadcrumb-sep {
    color: #b0b3c1;
    font-size: 0.78rem;
  }

  .topbar-current {
    font-size: 0.88rem;
    font-weight: 600;
    color: #191c1e;
  }

  .top-actions {
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
    display: grid;
    place-items: center;
  }

  .icon-btn:hover { background: #eceef0; }

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
    margin-bottom: 1.25rem;
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
    margin-bottom: 1.25rem;
  }

  /* ---- two panel layout ---- */
  .two-panel {
    display: flex;
    gap: 1.25rem;
    align-items: flex-start;
  }

  /* ---- left panel ---- */
  .panel-left {
    width: 35%;
    min-width: 240px;
    flex-shrink: 0;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    overflow: hidden;
  }

  .panel-left-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.85rem 1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    background: #f2f4f6;
  }

  .panel-title {
    font-size: 0.85rem;
    font-weight: 700;
    color: #191c1e;
  }

  .new-role-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.35rem 0.65rem;
    border-radius: 0.25rem;
    border: 1px solid #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    font-size: 0.72rem;
    font-weight: 600;
    color: #fff;
    cursor: pointer;
    font-family: inherit;
  }

  .new-role-btn:hover { opacity: 0.88; }
  .new-role-btn .material-symbols-outlined { font-size: 0.85rem; }

  /* ---- new role form ---- */
  .new-role-form {
    padding: 0.85rem 1rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    background: #f7f9fb;
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
  }

  .form-field {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
  }

  .form-field label {
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }

  .form-field input {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    border-radius: 0.25rem;
    padding: 0.4rem 0.6rem;
    font-size: 0.82rem;
    color: #191c1e;
    outline: none;
    font-family: inherit;
  }

  .form-error {
    margin: 0;
    font-size: 0.72rem;
    color: #ba1a1a;
  }

  .form-actions {
    display: flex;
    gap: 0.4rem;
    justify-content: flex-end;
  }

  /* ---- role list ---- */
  .role-list {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .role-card {
    width: 100%;
    text-align: left;
    padding: 0.85rem 1rem;
    border: none;
    background: transparent;
    cursor: pointer;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
    transition: background 0.1s;
    font-family: inherit;
  }

  .role-list li:last-child .role-card { border-bottom: 0; }

  .role-card:hover { background: #f7f9fb; }

  .role-card.selected {
    background: rgb(37 99 235 / 0.06);
    border-left: 3px solid #004ac6;
  }

  .role-card-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .role-name {
    font-size: 0.85rem;
    font-weight: 700;
    color: #191c1e;
    font-family: 'IBM Plex Mono', 'Courier New', monospace;
  }

  .perm-count {
    font-size: 0.62rem;
    color: #737686;
    flex-shrink: 0;
  }

  .role-desc {
    margin: 0.2rem 0 0;
    font-size: 0.72rem;
    color: #434655;
    line-height: 1.4;
  }

  .empty-inline {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding: 2.5rem 1rem;
    color: #b0b3c1;
  }

  .empty-inline .material-symbols-outlined { font-size: 2rem; }
  .empty-inline p { margin: 0; font-size: 0.78rem; }

  /* ---- right panel ---- */
  .panel-right {
    flex: 1;
    min-width: 0;
  }

  .placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    padding: 4rem 2rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    color: #b0b3c1;
    text-align: center;
  }

  .placeholder .material-symbols-outlined { font-size: 3rem; }
  .placeholder p { margin: 0; font-size: 0.85rem; max-width: 22rem; line-height: 1.5; }

  .detail-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
  }

  .detail-role-name {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 800;
    color: #191c1e;
    font-family: 'IBM Plex Mono', 'Courier New', monospace;
  }

  .detail-role-desc {
    margin: 0.25rem 0 0;
    font-size: 0.82rem;
    color: #434655;
  }

  .detail-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-shrink: 0;
  }

  .save-success {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.75rem;
    color: #065f46;
    font-weight: 600;
  }

  .save-success .material-symbols-outlined { font-size: 0.9rem; }

  /* ---- panel (perms table) ---- */
  .panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .perms-panel { overflow: visible; }

  .table-wrap { overflow-x: auto; }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    padding: 0.55rem 0.85rem;
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

  .align-center { text-align: center; }

  /* ---- resource cells ---- */
  .resource-name {
    display: block;
    font-weight: 700;
    font-size: 0.78rem;
    color: #191c1e;
  }

  .resource-key {
    display: block;
    font-size: 0.62rem;
    color: #737686;
    font-family: 'IBM Plex Mono', 'Courier New', monospace;
    margin-top: 0.05rem;
  }

  /* ---- checkbox ---- */
  .checkbox-wrap {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    position: relative;
    width: 1.25rem;
    height: 1.25rem;
  }

  .checkbox-wrap input[type='checkbox'] {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
    pointer-events: none;
  }

  .checkbox-custom {
    width: 1.1rem;
    height: 1.1rem;
    border: 2px solid rgb(195 198 215 / 0.7);
    border-radius: 0.2rem;
    background: #fff;
    display: grid;
    place-items: center;
    transition: background 0.1s, border-color 0.1s;
  }

  .checkbox-wrap input:checked + .checkbox-custom {
    background: #004ac6;
    border-color: #004ac6;
  }

  .checkbox-wrap input:checked + .checkbox-custom::after {
    content: '';
    display: block;
    width: 0.35rem;
    height: 0.6rem;
    border: 2px solid #fff;
    border-top: none;
    border-left: none;
    transform: rotate(45deg) translateY(-0.08rem);
  }

  .checkbox-wrap:hover .checkbox-custom {
    border-color: #2563eb;
  }

  /* ---- table footer ---- */
  .table-footer {
    padding: 0.55rem 0.85rem;
    font-size: 0.68rem;
    color: #434655;
    border-top: 1px solid rgb(195 198 215 / 0.35);
    background: #f7f9fb;
  }

  /* ---- buttons ---- */
  .ghost-btn,
  .primary-btn {
    border-radius: 0.25rem;
    padding: 0.48rem 0.85rem;
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

  .ghost-btn:hover { background: #f2f4f6; }

  .primary-btn {
    border-color: #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
  }

  .primary-btn:disabled { opacity: 0.6; cursor: not-allowed; }
  .primary-btn .material-symbols-outlined,
  .ghost-btn .material-symbols-outlined { font-size: 0.95rem; }

  /* ---- spin animation ---- */
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  .spin { animation: spin 0.8s linear infinite; font-size: 0.95rem; }
</style>
