<script lang="ts">
  import type { PageData, AppUser, UserStatus } from './+page.server';

  let { data }: { data: PageData } = $props();

  let users = $state<AppUser[]>([]);

  $effect(() => {
    users = data.users ?? [];
  });

  // ---- filters ----
  let statusFilter = $state<UserStatus | 'all'>('all');
  let searchQuery = $state('');

  const filteredUsers = $derived(
    users.filter((u) => {
      const matchStatus = statusFilter === 'all' || u.status === statusFilter;
      const matchSearch =
        !searchQuery ||
        u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        u.email.toLowerCase().includes(searchQuery.toLowerCase());
      return matchStatus && matchSearch;
    })
  );

  // ---- add user modal ----
  let addModalOpen = $state(false);
  let addForm = $state({ name: '', email: '', password: '', role: 'cs' });
  let addSaving = $state(false);
  let addError = $state('');

  function openAddModal() {
    addForm = { name: '', email: '', password: '', role: 'cs' };
    addError = '';
    addModalOpen = true;
  }

  function closeAddModal() {
    addModalOpen = false;
  }

  async function submitAdd() {
    if (!addForm.name.trim() || !addForm.email.trim() || !addForm.password) {
      addError = 'Semua field wajib diisi.';
      return;
    }
    addSaving = true;
    addError = '';
    try {
      await new Promise((r) => setTimeout(r, 500));
      const newUser: AppUser = {
        id: `u${Date.now()}`,
        name: addForm.name.trim(),
        email: addForm.email.trim(),
        roles: [addForm.role],
        status: 'pending',
        last_login: null,
        created_at: new Date().toISOString()
      };
      users = [...users, newUser];
      closeAddModal();
    } catch {
      addError = 'Gagal menambahkan pengguna.';
    } finally {
      addSaving = false;
    }
  }

  // ---- detail modal ----
  let detailModalOpen = $state(false);
  let detailUser = $state<AppUser | null>(null);
  let detailSaving = $state(false);
  let detailError = $state('');
  let detailEditMode = $state(false);
  let detailForm = $state({ name: '', email: '', roles: '' });

  function openDetail(user: AppUser) {
    detailUser = user;
    detailForm = { name: user.name, email: user.email, roles: user.roles.join(', ') };
    detailEditMode = false;
    detailError = '';
    detailModalOpen = true;
  }

  function closeDetail() {
    detailModalOpen = false;
    detailUser = null;
  }

  async function saveDetail() {
    if (!detailUser) return;
    detailSaving = true;
    detailError = '';
    try {
      await new Promise((r) => setTimeout(r, 400));
      const updated: AppUser = {
        ...detailUser,
        name: detailForm.name.trim(),
        email: detailForm.email.trim(),
        roles: detailForm.roles.split(',').map((r) => r.trim()).filter(Boolean)
      };
      users = users.map((u) => (u.id === updated.id ? updated : u));
      detailUser = updated;
      detailEditMode = false;
    } catch {
      detailError = 'Gagal menyimpan perubahan.';
    } finally {
      detailSaving = false;
    }
  }

  async function toggleSuspend() {
    if (!detailUser) return;
    detailSaving = true;
    try {
      await new Promise((r) => setTimeout(r, 400));
      const newStatus: UserStatus = detailUser.status === 'suspended' ? 'active' : 'suspended';
      const updated = { ...detailUser, status: newStatus };
      users = users.map((u) => (u.id === updated.id ? updated : u));
      detailUser = updated;
    } finally {
      detailSaving = false;
    }
  }

  async function resetPassword() {
    if (!detailUser) return;
    detailSaving = true;
    try {
      await new Promise((r) => setTimeout(r, 400));
      alert(`Link reset password dikirim ke ${detailUser.email}`);
    } finally {
      detailSaving = false;
    }
  }

  // ---- helpers ----
  function relativeTime(iso: string | null): string {
    if (!iso) return 'Belum pernah login';
    const diff = Date.now() - new Date(iso).getTime();
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);
    if (hours < 1) return 'Baru saja';
    if (hours < 24) return `${hours} jam lalu`;
    if (days === 1) return 'Kemarin';
    return `${days} hari lalu`;
  }

  const STATUS_LABELS: Record<UserStatus, string> = {
    active: 'Aktif',
    suspended: 'Ditangguhkan',
    pending: 'Pending'
  };
</script>

<main class="page-shell">
  <header class="topbar">
    <div class="breadcrumb">
      <span class="material-symbols-outlined bc-icon">manage_accounts</span>
      <span class="bc-text">Manajemen Pengguna</span>
    </div>
    <div class="top-actions">
      <button class="icon-btn">
        <span class="material-symbols-outlined">notifications</span>
        <span class="dot"></span>
      </button>
      <button class="avatar" aria-label="Profile">AD</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Pengguna</h2>
        <p>Kelola akun pengguna konsol — status, role, dan akses</p>
      </div>
      <button type="button" class="primary-btn" onclick={openAddModal}>
        <span class="material-symbols-outlined">person_add</span>
        Tambah Pengguna
      </button>
    </div>

    {#if data.error}
      <div class="error-banner">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- Filters -->
    <div class="filters-row">
      <div class="filter-group">
        <label class="field-label" for="f-status">Status</label>
        <select id="f-status" class="filter-select" bind:value={statusFilter}>
          <option value="all">Semua Status</option>
          <option value="active">Aktif</option>
          <option value="suspended">Ditangguhkan</option>
          <option value="pending">Pending</option>
        </select>
      </div>
      <div class="filter-group search-group">
        <label class="field-label" for="f-search">Cari</label>
        <div class="search-inline">
          <span class="material-symbols-outlined">search</span>
          <input id="f-search" type="text" placeholder="Nama atau email..." bind:value={searchQuery} />
        </div>
      </div>
      {#if statusFilter !== 'all' || searchQuery}
        <button type="button" class="clear-btn" onclick={() => { statusFilter = 'all'; searchQuery = ''; }}>
          <span class="material-symbols-outlined">close</span>
          Reset
        </button>
      {/if}
    </div>

    <!-- Table -->
    <div class="panel">
      <div class="table-wrap">
        {#if filteredUsers.length === 0}
          <div class="empty-state">
            <span class="material-symbols-outlined">manage_accounts</span>
            <p>Tidak ada pengguna yang sesuai filter.</p>
          </div>
        {:else}
          <table>
            <thead>
              <tr>
                <th>Nama & Email</th>
                <th>Role</th>
                <th>Status</th>
                <th>Terakhir Login</th>
                <th class="align-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each filteredUsers as user (user.id)}
                <tr class="clickable-row" onclick={() => openDetail(user)}>
                  <td>
                    <div class="user-cell">
                      <div class="user-avatar">{user.name.charAt(0)}</div>
                      <div>
                        <div class="user-name">{user.name}</div>
                        <div class="user-email">{user.email}</div>
                      </div>
                    </div>
                  </td>
                  <td>
                    <div class="roles-cell">
                      {#each user.roles as role}
                        <span class="role-badge">{role}</span>
                      {/each}
                    </div>
                  </td>
                  <td>
                    <span class="status-badge status-badge--{user.status}">
                      {STATUS_LABELS[user.status]}
                    </span>
                  </td>
                  <td class="time-cell">{relativeTime(user.last_login)}</td>
                  <td class="actions-cell" onclick={(e) => e.stopPropagation()}>
                    <button class="action-btn" onclick={() => openDetail(user)}>
                      <span class="material-symbols-outlined">open_in_new</span>
                      Detail
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>
      <div class="table-footer">Menampilkan {filteredUsers.length} dari {users.length} pengguna</div>
    </div>
  </section>
</main>

<!-- Add User Modal -->
{#if addModalOpen}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div class="modal-backdrop" onclick={closeAddModal}></div>
  <div class="modal" role="dialog" aria-modal="true">
    <div class="modal-header">
      <h3>Tambah Pengguna Baru</h3>
      <button class="modal-close" onclick={closeAddModal}><span class="material-symbols-outlined">close</span></button>
    </div>
    <div class="modal-body">
      <div class="field-row">
        <label class="field-label" for="add-name">Nama Lengkap</label>
        <input id="add-name" type="text" class="field-input" bind:value={addForm.name} placeholder="Nama lengkap" />
      </div>
      <div class="field-row">
        <label class="field-label" for="add-email">Email</label>
        <input id="add-email" type="email" class="field-input" bind:value={addForm.email} placeholder="email@umrohos.id" />
      </div>
      <div class="field-row">
        <label class="field-label" for="add-pass">Password</label>
        <input id="add-pass" type="password" class="field-input" bind:value={addForm.password} placeholder="••••••••" />
      </div>
      <div class="field-row">
        <label class="field-label" for="add-role">Role</label>
        <select id="add-role" class="field-input" bind:value={addForm.role}>
          <option value="admin">Admin</option>
          <option value="finance">Finance</option>
          <option value="cs">CS</option>
          <option value="ops">Ops</option>
        </select>
      </div>
      {#if addError}<p class="modal-error">{addError}</p>{/if}
    </div>
    <div class="modal-footer">
      <button class="ghost-btn" onclick={closeAddModal} disabled={addSaving}>Batal</button>
      <button class="primary-btn" onclick={submitAdd} disabled={addSaving}>
        {#if addSaving}<span class="material-symbols-outlined spin">progress_activity</span>Menyimpan...
        {:else}<span class="material-symbols-outlined">save</span>Tambah Pengguna
        {/if}
      </button>
    </div>
  </div>
{/if}

<!-- Detail Modal -->
{#if detailModalOpen && detailUser}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div class="modal-backdrop" onclick={closeDetail}></div>
  <div class="modal modal--wide" role="dialog" aria-modal="true">
    <div class="modal-header">
      <h3>Detail Pengguna</h3>
      <button class="modal-close" onclick={closeDetail}><span class="material-symbols-outlined">close</span></button>
    </div>
    <div class="modal-body">
      <div class="user-detail-head">
        <div class="user-avatar large">{detailUser.name.charAt(0)}</div>
        <div>
          <div class="user-name">{detailUser.name}</div>
          <div class="user-email">{detailUser.email}</div>
          <span class="status-badge status-badge--{detailUser.status} mt-2">{STATUS_LABELS[detailUser.status]}</span>
        </div>
      </div>

      {#if detailEditMode}
        <div class="field-row">
          <label class="field-label" for="det-name">Nama</label>
          <input id="det-name" type="text" class="field-input" bind:value={detailForm.name} />
        </div>
        <div class="field-row">
          <label class="field-label" for="det-email">Email</label>
          <input id="det-email" type="email" class="field-input" bind:value={detailForm.email} />
        </div>
        <div class="field-row">
          <label class="field-label" for="det-roles">Roles (pisah koma)</label>
          <input id="det-roles" type="text" class="field-input" bind:value={detailForm.roles} placeholder="admin, cs, finance" />
        </div>
      {:else}
        <div class="detail-info-grid">
          <div class="di-row">
            <span class="di-label">Roles</span>
            <div class="roles-cell">
              {#each detailUser.roles as role}
                <span class="role-badge">{role}</span>
              {/each}
            </div>
          </div>
          <div class="di-row">
            <span class="di-label">Terakhir Login</span>
            <span>{relativeTime(detailUser.last_login)}</span>
          </div>
          <div class="di-row">
            <span class="di-label">Dibuat</span>
            <span>{new Date(detailUser.created_at).toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' })}</span>
          </div>
        </div>
      {/if}

      {#if detailError}<p class="modal-error">{detailError}</p>{/if}
    </div>
    <div class="modal-footer modal-footer--space">
      <div class="footer-left">
        <button class="ghost-btn" onclick={toggleSuspend} disabled={detailSaving}>
          <span class="material-symbols-outlined">{detailUser.status === 'suspended' ? 'lock_open' : 'block'}</span>
          {detailUser.status === 'suspended' ? 'Aktifkan' : 'Tangguhkan'}
        </button>
        <button class="ghost-btn" onclick={resetPassword} disabled={detailSaving}>
          <span class="material-symbols-outlined">key</span>
          Reset Password
        </button>
      </div>
      <div class="footer-right">
        {#if detailEditMode}
          <button class="ghost-btn" onclick={() => { detailEditMode = false; }} disabled={detailSaving}>Batal</button>
          <button class="primary-btn" onclick={saveDetail} disabled={detailSaving}>
            {#if detailSaving}<span class="material-symbols-outlined spin">progress_activity</span>Menyimpan...
            {:else}<span class="material-symbols-outlined">save</span>Simpan
            {/if}
          </button>
        {:else}
          <button class="ghost-btn" onclick={closeDetail}>Tutup</button>
          <button class="primary-btn" onclick={() => { detailEditMode = true; }}>
            <span class="material-symbols-outlined">edit</span>
            Edit
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }

  .topbar {
    position: sticky; top: 0; z-index: 30; height: 4rem;
    background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center;
    justify-content: space-between; gap: 1rem; backdrop-filter: blur(8px);
  }

  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.88rem; font-weight: 600; color: #191c1e; }
  .bc-icon { font-size: 1.05rem; color: #004ac6; }
  .top-actions { display: flex; align-items: center; gap: 0.5rem; }
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; position: relative; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .dot { position: absolute; width: 0.46rem; height: 0.46rem; border-radius: 999px; background: #ba1a1a; right: 0.4rem; top: 0.35rem; border: 2px solid #fff; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas { padding: 1.5rem; max-width: 96rem; }

  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .error-banner { display: flex; align-items: center; gap: 0.5rem; background: #ffdad6; color: #93000a; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.82rem; margin-bottom: 1rem; }

  .filters-row { display: flex; align-items: flex-end; gap: 0.75rem; flex-wrap: wrap; margin-bottom: 1rem; }
  .filter-group { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .filter-select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; min-width: 10rem; outline: none; }

  .search-group .search-inline { display: flex; align-items: center; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; background: #fff; padding: 0 0.6rem; }
  .search-group .search-inline span { font-size: 0.9rem; color: #737686; }
  .search-group .search-inline input { border: 0; background: transparent; padding: 0.42rem 0.4rem; min-width: 14rem; font-size: 0.82rem; outline: none; font-family: inherit; }

  .clear-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.42rem 0.6rem; border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; background: #fff; font-size: 0.78rem; color: #434655; cursor: pointer; align-self: flex-end; }
  .clear-btn:hover { background: #f2f4f6; }
  .clear-btn .material-symbols-outlined { font-size: 0.9rem; }

  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.25rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.62rem 0.85rem; text-align: left; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); white-space: nowrap; }
  th { text-transform: uppercase; font-size: 0.62rem; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; }
  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .align-right { text-align: right; }

  .clickable-row { cursor: pointer; }

  .user-cell { display: flex; align-items: center; gap: 0.65rem; }
  .user-avatar {
    width: 2rem; height: 2rem; border-radius: 999px;
    background: #b4c5ff; color: #00174b; font-size: 0.75rem;
    font-weight: 700; display: grid; place-items: center; flex-shrink: 0;
  }
  .user-avatar.large { width: 3rem; height: 3rem; font-size: 1.1rem; }
  .user-name { font-weight: 700; font-size: 0.82rem; color: #191c1e; }
  .user-email { font-size: 0.72rem; color: #434655; }

  .roles-cell { display: flex; gap: 0.25rem; flex-wrap: wrap; }
  .role-badge { display: inline-flex; padding: 0.1rem 0.35rem; border-radius: 0.2rem; font-size: 0.62rem; font-weight: 700; background: #ede9fe; color: #4c1d95; }

  .status-badge { display: inline-flex; padding: 0.15rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 700; }
  .status-badge--active { background: #d1fae5; color: #065f46; }
  .status-badge--suspended { background: #fee2e2; color: #991b1b; }
  .status-badge--pending { background: #fef9c3; color: #7d5f00; }
  .mt-2 { margin-top: 0.35rem; }

  .time-cell { font-size: 0.72rem; color: #434655; }

  .actions-cell { text-align: right; }
  .action-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.3rem 0.55rem; border-radius: 0.2rem; border: 1px solid rgb(195 198 215 / 0.55); background: #fff; font-size: 0.72rem; font-weight: 600; color: #191c1e; cursor: pointer; }
  .action-btn:hover { background: #f2f4f6; }
  .action-btn .material-symbols-outlined { font-size: 0.85rem; }

  .table-footer { padding: 0.55rem 0.85rem; font-size: 0.68rem; color: #434655; border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb; }

  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.6rem; padding: 3rem 1rem; color: #b0b3c1; }
  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.82rem; }

  /* modal */
  .modal-backdrop { position: fixed; inset: 0; background: rgb(0 0 0 / 0.35); z-index: 50; }
  .modal {
    position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%);
    z-index: 51; width: min(480px, calc(100vw - 2rem)); background: #fff;
    border-radius: 0.4rem; border: 1px solid rgb(195 198 215 / 0.55);
    box-shadow: 0 8px 24px rgb(0 0 0 / 0.12); display: flex; flex-direction: column;
  }
  .modal--wide { width: min(560px, calc(100vw - 2rem)); }
  .modal-header { display: flex; align-items: center; justify-content: space-between; padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); background: #f2f4f6; }
  .modal-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .modal-close { border: 0; background: transparent; cursor: pointer; color: #434655; display: grid; place-items: center; border-radius: 0.2rem; padding: 0.2rem; }
  .modal-close:hover { background: #e6e8ea; }
  .modal-body { padding: 1rem; display: flex; flex-direction: column; gap: 0.75rem; }
  .modal-error { margin: 0; font-size: 0.75rem; color: #ba1a1a; background: #ffdad6; padding: 0.45rem 0.65rem; border-radius: 0.2rem; }
  .modal-footer { padding: 0.75rem 1rem; border-top: 1px solid rgb(195 198 215 / 0.45); display: flex; justify-content: flex-end; gap: 0.55rem; }
  .modal-footer--space { justify-content: space-between; }
  .footer-left, .footer-right { display: flex; gap: 0.55rem; }

  .user-detail-head { display: flex; align-items: center; gap: 0.85rem; padding: 0.6rem; background: #f2f4f6; border-radius: 0.25rem; }
  .detail-info-grid { display: flex; flex-direction: column; gap: 0.5rem; }
  .di-row { display: flex; align-items: center; gap: 0.5rem; font-size: 0.82rem; }
  .di-label { width: 7rem; font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; flex-shrink: 0; }

  .field-row { display: flex; flex-direction: column; gap: 0.25rem; }
  .field-input { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.5rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; outline: none; }

  .ghost-btn, .primary-btn {
    border-radius: 0.25rem; padding: 0.5rem 0.85rem; font-size: 0.8rem;
    font-weight: 600; cursor: pointer; border: 1px solid rgb(195 198 215 / 0.55);
    display: inline-flex; align-items: center; gap: 0.35rem; font-family: inherit;
  }
  .ghost-btn { background: #fff; color: #191c1e; }
  .ghost-btn:disabled, .primary-btn:disabled { opacity: 0.6; cursor: not-allowed; }
  .primary-btn { border-color: #2563eb; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; }
  .primary-btn .material-symbols-outlined, .ghost-btn .material-symbols-outlined { font-size: 1rem; }

  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; }
</style>
