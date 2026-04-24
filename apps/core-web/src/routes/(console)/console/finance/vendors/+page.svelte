<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let vendors = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  // slide-over state
  let showForm = $state(false);
  let editingId = $state<string | null>(null);
  let formLoading = $state(false);
  let formError = $state('');
  let formSuccess = $state('');

  // form fields
  let fName = $state('');
  let fCategory = $state('hotel');
  let fBankAccount = $state('');
  let fTaxId = $state('');
  let fContactEmail = $state('');

  async function loadVendors() {
    loading = true; error = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/vendors`);
      const body = res.ok ? await res.json() : null;
      vendors = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { error = 'Gagal memuat data vendor.'; vendors = []; }
    loading = false;
  }

  $effect(() => { loadVendors(); });

  function openAdd() {
    editingId = null; fName = ''; fCategory = 'hotel'; fBankAccount = ''; fTaxId = ''; fContactEmail = '';
    formError = ''; formSuccess = ''; showForm = true;
  }

  function openEdit(v: any) {
    editingId = v.id; fName = v.name ?? ''; fCategory = v.category ?? 'hotel';
    fBankAccount = v.bank_account ?? ''; fTaxId = v.tax_id ?? ''; fContactEmail = v.contact_email ?? '';
    formError = ''; formSuccess = ''; showForm = true;
  }

  async function submitForm(e: SubmitEvent) {
    e.preventDefault(); formLoading = true; formError = ''; formSuccess = '';
    const body = { name: fName, category: fCategory, bank_account: fBankAccount, tax_id: fTaxId, contact_email: fContactEmail };
    try {
      const res = await fetch(`${GATEWAY}/v1/finance/vendors${editingId ? '/' + editingId : ''}`, {
        method: editingId ? 'PUT' : 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      });
      if (!res.ok) throw new Error(`${res.status}`);
      formSuccess = editingId ? 'Vendor berhasil diperbarui.' : 'Vendor berhasil ditambahkan.';
      await loadVendors();
      setTimeout(() => { showForm = false; formSuccess = ''; }, 1200);
    } catch (err: any) { formError = `Gagal menyimpan: ${err.message}`; }
    formLoading = false;
  }

  async function deleteVendor(id: string) {
    if (!confirm('Hapus vendor ini?')) return;
    try {
      await fetch(`${GATEWAY}/v1/finance/vendors/${id}`, { method: 'DELETE' });
      await loadVendors();
    } catch { alert('Gagal menghapus vendor.'); }
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Master Vendor</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Master Vendor</h2>
        <p>Kelola data vendor & supplier</p>
      </div>
      <button class="btn-primary" onclick={openAdd}>
        <span class="material-symbols-outlined">add</span>
        Tambah Vendor
      </button>
    </div>

    {#if error}
      <div class="alert-err">{error}</div>
    {/if}

    <div class="panel">
      {#if loading}
        <div class="loading-row">
          <span class="material-symbols-outlined spin">progress_activity</span> Memuat...
        </div>
      {:else if vendors.length === 0}
        <div class="empty">Belum ada vendor. Klik "Tambah Vendor" untuk memulai.</div>
      {:else}
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Nama</th><th>Kategori</th><th>Email</th><th>NPWP</th><th>Rekening Bank</th><th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each vendors as v}
                <tr>
                  <td class="fw">{v.name}</td>
                  <td><span class="chip">{v.category}</span></td>
                  <td>{v.contact_email ?? '-'}</td>
                  <td class="mono">{v.tax_id ?? '-'}</td>
                  <td>{v.bank_account ?? '-'}</td>
                  <td>
                    <div class="row-actions">
                      <button class="btn-sm" onclick={() => openEdit(v)}>Edit</button>
                      <button class="btn-sm danger" onclick={() => deleteVendor(v.id)}>Hapus</button>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  </section>
</main>

<!-- Slide-over form -->
{#if showForm}
  <div class="overlay" role="presentation" onclick={() => showForm = false}></div>
  <aside class="slideover">
    <div class="slideover-header">
      <h3>{editingId ? 'Edit Vendor' : 'Tambah Vendor'}</h3>
      <button class="close-btn" onclick={() => showForm = false}>
        <span class="material-symbols-outlined">close</span>
      </button>
    </div>
    <form class="slideover-body" onsubmit={submitForm}>
      {#if formError}<div class="alert-err">{formError}</div>{/if}
      {#if formSuccess}<div class="alert-ok">{formSuccess}</div>{/if}

      <div class="field">
        <label>Nama Vendor *</label>
        <input type="text" bind:value={fName} required placeholder="PT Makkah Tour" />
      </div>
      <div class="field">
        <label>Kategori *</label>
        <select bind:value={fCategory}>
          <option value="hotel">Hotel</option>
          <option value="airline">Maskapai</option>
          <option value="transport">Transportasi</option>
          <option value="catering">Katering</option>
          <option value="other">Lainnya</option>
        </select>
      </div>
      <div class="field">
        <label>Rekening Bank</label>
        <input type="text" bind:value={fBankAccount} placeholder="BCA 1234567890 a.n. PT Makkah" />
      </div>
      <div class="field">
        <label>NPWP</label>
        <input type="text" bind:value={fTaxId} placeholder="01.234.567.8-901.000" />
      </div>
      <div class="field">
        <label>Email Kontak</label>
        <input type="email" bind:value={fContactEmail} placeholder="vendor@example.com" />
      </div>
      <div class="form-actions">
        <button type="button" class="btn-ghost" onclick={() => showForm = false}>Batal</button>
        <button type="submit" class="btn-primary" disabled={formLoading}>
          {formLoading ? 'Menyimpan...' : 'Simpan'}
        </button>
      </div>
    </form>
  </aside>
{/if}

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar {
    position: sticky; top: 0; z-index: 30; height: 4rem;
    background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px);
  }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 1.25rem; gap: 1rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #15803d; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .panel { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.6rem 0.85rem; font-size: 0.78rem; text-align: left; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  th { background: #f2f4f6; font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  tbody tr:hover { background: #f7f9fb; }
  .fw { font-weight: 600; }
  .mono { font-family: monospace; font-size: 0.74rem; }
  .chip { display: inline-flex; padding: 0.1rem 0.45rem; background: #e0f2fe; color: #075985; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .row-actions { display: flex; gap: 0.35rem; }
  .btn-sm { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.2rem; padding: 0.25rem 0.55rem; font-size: 0.72rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #191c1e; }
  .btn-sm:hover { background: #f2f4f6; }
  .btn-sm.danger { color: #dc2626; border-color: #fca5a5; }
  .btn-sm.danger:hover { background: #fef2f2; }
  .empty { text-align: center; color: #b0b3c1; padding: 3rem; font-size: 0.82rem; }
  .loading-row { display: flex; align-items: center; gap: 0.5rem; padding: 1.5rem; color: #737686; font-size: 0.82rem; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .btn-primary .material-symbols-outlined { font-size: 1rem; }
  .btn-ghost { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.85rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; color: #191c1e; }
  .btn-ghost:hover { background: #f2f4f6; }
  .overlay { position: fixed; inset: 0; background: rgb(0 0 0 / 0.3); z-index: 40; }
  .slideover {
    position: fixed; top: 0; right: 0; bottom: 0; width: 380px;
    background: #fff; z-index: 50; box-shadow: -4px 0 24px rgb(0 0 0 / 0.12);
    display: flex; flex-direction: column;
  }
  .slideover-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 1rem 1.25rem; border-bottom: 1px solid rgb(195 198 215 / 0.45);
  }
  .slideover-header h3 { margin: 0; font-size: 1rem; font-weight: 700; }
  .close-btn { border: none; background: transparent; cursor: pointer; color: #737686; display: grid; place-items: center; }
  .close-btn .material-symbols-outlined { font-size: 1.2rem; }
  .slideover-body { flex: 1; overflow-y: auto; padding: 1.25rem; display: flex; flex-direction: column; gap: 1rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select {
    border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem;
    padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none;
  }
  .field input:focus, .field select:focus { border-color: #2563eb; }
  .form-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 0.5rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
