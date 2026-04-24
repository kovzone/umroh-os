<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  type Severity = 'rendah' | 'sedang' | 'tinggi';
  type Category = 'akomodasi' | 'transportasi' | 'makanan' | 'kesehatan' | 'ibadah' | 'lainnya';

  interface Report {
    id: string;
    date: string;
    category: Category;
    description: string;
    severity: Severity;
    status: 'open' | 'in_progress' | 'resolved';
  }

  const myReports = $state<Report[]>([
    {
      id: 'RPT-001',
      date: '22 Jan 2025',
      category: 'akomodasi',
      description: 'AC kamar tidak berfungsi dengan baik sejak kemarin malam.',
      severity: 'sedang',
      status: 'resolved',
    },
    {
      id: 'RPT-002',
      date: '20 Jan 2025',
      category: 'makanan',
      description: 'Menu makan malam tidak sesuai yang dijanjikan (tidak ada lauk daging).',
      severity: 'rendah',
      status: 'in_progress',
    },
  ]);

  let category = $state<Category>('akomodasi');
  let description = $state('');
  let severity = $state<Severity>('sedang');
  let photoName = $state('');
  let submitting = $state(false);
  let submitSuccess = $state(false);

  function handlePhoto(e: Event) {
    const input = e.target as HTMLInputElement;
    photoName = input.files?.[0]?.name ?? '';
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    if (!description.trim()) return;
    submitting = true;
    setTimeout(() => {
      myReports.unshift({
        id: `RPT-00${myReports.length + 1}`,
        date: 'Hari ini',
        category,
        description,
        severity,
        status: 'open',
      });
      description = '';
      photoName = '';
      submitting = false;
      submitSuccess = true;
      setTimeout(() => { submitSuccess = false; }, 3000);
    }, 1500);
  }

  const catLabel: Record<Category, string> = {
    akomodasi: 'Akomodasi',
    transportasi: 'Transportasi',
    makanan: 'Makanan',
    kesehatan: 'Kesehatan',
    ibadah: 'Ibadah',
    lainnya: 'Lainnya',
  };

  const sevLabel: Record<Severity, string> = { rendah: 'Rendah', sedang: 'Sedang', tinggi: 'Tinggi' };
  const sevColor: Record<Severity, string> = { rendah: '#006747', sedang: '#775a19', tinggi: '#ba1a1a' };
  const statusLabel: Record<string, string> = { open: 'Terbuka', in_progress: 'Diproses', resolved: 'Selesai' };
  const statusColor: Record<string, string> = { open: '#ba1a1a', in_progress: '#775a19', resolved: '#006747' };
</script>

<svelte:head>
  <title>Laporan Kendala — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="laporan-root">
    <div class="shell">
      <a href="/jemaah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Portal Jamaah
      </a>
      <div class="page-header">
        <h1>Laporan Kendala</h1>
        <p>Sampaikan masalah perjalanan agar segera ditangani tim kami</p>
      </div>

      <div class="two-col">
        <!-- Report form -->
        <div>
          <div class="form-card">
            <h2>Buat Laporan Baru</h2>
            {#if submitSuccess}
              <div class="success-banner">
                <span class="material-symbols-outlined">check_circle</span>
                Laporan berhasil dikirim. Tim kami akan menindaklanjuti segera.
              </div>
            {/if}
            <form onsubmit={handleSubmit} class="report-form">
              <div class="field">
                <label for="cat">Kategori Masalah</label>
                <select id="cat" bind:value={category}>
                  {#each Object.entries(catLabel) as [val, lbl]}
                    <option value={val}>{lbl}</option>
                  {/each}
                </select>
              </div>

              <div class="field">
                <label for="desc">Deskripsi Masalah</label>
                <textarea id="desc" rows="4" placeholder="Ceritakan masalah yang Anda alami secara detail..." bind:value={description} required></textarea>
              </div>

              <div class="field">
                <label>Tingkat Keparahan</label>
                <div class="severity-options">
                  {#each ['rendah', 'sedang', 'tinggi'] as s}
                    <label class="sev-opt" class:selected={severity === s} style="--sev-color: {sevColor[s as Severity]}">
                      <input type="radio" name="severity" value={s} bind:group={severity} />
                      {sevLabel[s as Severity]}
                    </label>
                  {/each}
                </div>
              </div>

              <div class="field">
                <label>Foto Pendukung (opsional)</label>
                <label class="photo-upload">
                  <span class="material-symbols-outlined">photo_camera</span>
                  {photoName || 'Pilih foto...'}
                  <input type="file" accept="image/*" onchange={handlePhoto} />
                </label>
              </div>

              <button type="submit" class="submit-btn" disabled={submitting || !description.trim()}>
                {#if submitting}
                  <span class="material-symbols-outlined spin">progress_activity</span>
                  Mengirim...
                {:else}
                  <span class="material-symbols-outlined">send</span>
                  Kirim Laporan
                {/if}
              </button>
            </form>
          </div>
        </div>

        <!-- History -->
        <div>
          <h2 class="hist-title">Riwayat Laporan Saya</h2>
          {#if myReports.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">inbox</span>
              <p>Belum ada laporan yang dibuat</p>
            </div>
          {:else}
            <div class="reports-list">
              {#each myReports as r (r.id)}
                <div class="report-item">
                  <div class="report-top">
                    <span class="report-id">{r.id}</span>
                    <span class="report-date">{r.date}</span>
                  </div>
                  <div class="report-badges">
                    <span class="cat-badge">{catLabel[r.category]}</span>
                    <span class="sev-badge" style="background: color-mix(in srgb, {sevColor[r.severity]} 12%, transparent); color: {sevColor[r.severity]}">{sevLabel[r.severity]}</span>
                    <span class="status-badge" style="background: color-mix(in srgb, {statusColor[r.status]} 12%, transparent); color: {statusColor[r.status]}">{statusLabel[r.status]}</span>
                  </div>
                  <p class="report-desc">{r.description}</p>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .laporan-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 80rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; align-items: start; }
  .form-card { background: #fff; border-radius: 1.5rem; padding: 1.75rem; border: 1px solid rgba(190,201,193,0.2); }
  .form-card h2 { margin: 0 0 1.25rem; font-size: 1.05rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .success-banner { display: flex; align-items: center; gap: 0.5rem; background: rgba(0,103,71,0.08); color: #006747; border-radius: 0.75rem; padding: 0.75rem 1rem; font-size: 0.85rem; font-weight: 600; margin-bottom: 1.25rem; }
  .success-banner .material-symbols-outlined { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .report-form { display: flex; flex-direction: column; gap: 1rem; }
  .field { display: flex; flex-direction: column; gap: 0.35rem; }
  .field label { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #6b7280; }
  .field select, .field textarea {
    border: 1px solid rgba(190,201,193,0.4);
    border-radius: 0.75rem;
    padding: 0.65rem 0.85rem;
    font-size: 0.88rem;
    color: #1b1c1c;
    background: #fbf9f8;
    font-family: inherit;
    resize: vertical;
  }
  .severity-options { display: flex; gap: 0.5rem; }
  .sev-opt {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.6rem;
    border-radius: 0.75rem;
    border: 1.5px solid rgba(190,201,193,0.3);
    cursor: pointer;
    font-size: 0.82rem;
    font-weight: 600;
    transition: all 0.15s;
  }
  .sev-opt.selected { background: color-mix(in srgb, var(--sev-color) 12%, transparent); border-color: var(--sev-color); color: var(--sev-color); }
  .sev-opt input { display: none; }
  .photo-upload {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    border: 1.5px dashed rgba(190,201,193,0.5);
    border-radius: 0.75rem;
    padding: 0.65rem 1rem;
    cursor: pointer;
    font-size: 0.85rem;
    color: #9ca3af;
  }
  .photo-upload input { display: none; }
  .photo-upload .material-symbols-outlined { font-size: 1.1rem; color: #006747; }
  .submit-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: linear-gradient(135deg, #004d34, #006747);
    color: #fff;
    border: none;
    border-radius: 999px;
    padding: 0.8rem 1.8rem;
    font-size: 0.9rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    align-self: flex-start;
  }
  .submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
  /* History */
  .hist-title { margin: 0 0 1.2rem; font-size: 1.05rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .reports-list { display: flex; flex-direction: column; gap: 0.85rem; }
  .report-item { background: #fff; border-radius: 1.2rem; padding: 1rem 1.2rem; border: 1px solid rgba(190,201,193,0.2); }
  .report-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.5rem; }
  .report-id { font-family: 'IBM Plex Mono', monospace; font-size: 0.75rem; font-weight: 700; color: #006747; }
  .report-date { font-size: 0.72rem; color: #9ca3af; }
  .report-badges { display: flex; gap: 0.4rem; flex-wrap: wrap; margin-bottom: 0.5rem; }
  .cat-badge { font-size: 0.7rem; font-weight: 700; background: rgba(190,201,193,0.2); color: #6b7280; border-radius: 999px; padding: 0.2rem 0.6rem; }
  .sev-badge, .status-badge { font-size: 0.7rem; font-weight: 700; border-radius: 999px; padding: 0.2rem 0.6rem; }
  .report-desc { margin: 0; font-size: 0.82rem; color: #57534e; line-height: 1.5; }
  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.5rem; padding: 3rem 1rem; color: #9ca3af; }
  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.85rem; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; }
  @media (max-width: 760px) { .two-col { grid-template-columns: 1fr; } }
</style>
