<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  type DocStatus = 'uploaded' | 'verified' | 'pending' | 'rejected' | 'required';

  interface DocItem {
    id: string;
    name: string;
    description: string;
    icon: string;
    required: boolean;
    status: DocStatus;
    uploadedAt?: string;
    fileName?: string;
    rejectReason?: string;
  }

  let documents = $state<DocItem[]>([
    { id: 'passport', name: 'Paspor', description: 'Foto halaman biodata paspor (min 6 bulan masa berlaku)', icon: 'badge', required: true, status: 'verified', uploadedAt: '1 Des 2024', fileName: 'paspor_bambang.jpg' },
    { id: 'kk', name: 'Kartu Keluarga', description: 'Scan KK yang masih berlaku', icon: 'family_restroom', required: true, status: 'uploaded', uploadedAt: '1 Des 2024', fileName: 'kk_bambang.pdf' },
    { id: 'ktp', name: 'KTP', description: 'Foto KTP tampak depan, jelas dan tidak buram', icon: 'contact_page', required: true, status: 'pending' },
    { id: 'buku_nikah', name: 'Buku Nikah', description: 'Untuk jamaah yang berangkat dengan pasangan', icon: 'favorite', required: false, status: 'required' },
    { id: 'vaksin', name: 'Sertifikat Vaksinasi', description: 'Sertifikat vaksin meningitis (wajib untuk visa Saudi)', icon: 'vaccines', required: true, status: 'rejected', uploadedAt: '28 Nov 2024', fileName: 'vaksin_bambang.jpg', rejectReason: 'Kualitas foto tidak jelas. Harap upload ulang dengan foto yang lebih terang.' },
    { id: 'foto', name: 'Foto 4×6 Background Putih', description: 'Foto terbaru 6 bulan terakhir, background putih, tampak depan', icon: 'photo_camera', required: true, status: 'required' }
  ]);

  let draggingOver = $state<string | null>(null);
  let uploading = $state<string | null>(null);

  const summary = $derived({
    total: documents.length,
    verified: documents.filter(d => d.status === 'verified').length,
    uploaded: documents.filter(d => d.status === 'uploaded').length,
    pending: documents.filter(d => d.status === 'pending').length,
    rejected: documents.filter(d => d.status === 'rejected').length,
    required: documents.filter(d => d.status === 'required').length
  });

  const overallProgress = $derived(
    documents.length > 0 ? ((summary.verified + summary.uploaded) / documents.length) * 100 : 0
  );

  function statusLabel(status: DocStatus): string {
    const labels: Record<DocStatus, string> = {
      uploaded: 'Diunggah', verified: 'Terverifikasi', pending: 'Menunggu Review',
      rejected: 'Ditolak', required: 'Perlu Diunggah'
    };
    return labels[status];
  }

  function handleDrop(e: DragEvent, docId: string) {
    e.preventDefault();
    draggingOver = null;
    const file = e.dataTransfer?.files[0];
    if (file) simulateUpload(docId, file.name);
  }

  function handleFileInput(e: Event, docId: string) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (file) simulateUpload(docId, file.name);
  }

  function simulateUpload(docId: string, fileName: string) {
    uploading = docId;
    setTimeout(() => {
      documents = documents.map(d =>
        d.id === docId
          ? { ...d, status: 'pending', fileName, uploadedAt: 'Baru saja' }
          : d
      );
      uploading = null;
    }, 1500);
  }

  function removeDoc(docId: string) {
    documents = documents.map(d =>
      d.id === docId
        ? { ...d, status: 'required', fileName: undefined, uploadedAt: undefined }
        : d
    );
  }
</script>

<svelte:head>
  <title>Dokumen Saya — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" ctaLabel="Lihat Paket" packagesLinkActive={false}>
  <div class="docs-root">
    <div class="shell">
      <div class="page-header">
        <div>
          <h1>Dokumen Perjalanan</h1>
          <p>Upload dan pantau kelengkapan dokumen umrah Anda</p>
        </div>
        <a class="btn-wa" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
          <span class="material-symbols-outlined">chat</span>
          Bantuan Dokumen
        </a>
      </div>

      <!-- Progress summary -->
      <div class="progress-card">
        <div class="progress-stats">
          {#each [
            { label: 'Terverifikasi', count: summary.verified, color: '#006747' },
            { label: 'Menunggu', count: summary.uploaded + summary.pending, color: '#775a19' },
            { label: 'Belum Diunggah', count: summary.required, color: '#9ca3af' },
            { label: 'Ditolak', count: summary.rejected, color: '#ba1a1a' }
          ] as stat (stat.label)}
            <div class="progress-stat">
              <strong style="color: {stat.color}">{stat.count}</strong>
              <span>{stat.label}</span>
            </div>
          {/each}
        </div>
        <div class="progress-bar-section">
          <div class="progress-header">
            <span>Kelengkapan Dokumen</span>
            <span>{overallProgress.toFixed(0)}%</span>
          </div>
          <div class="progress-bar-wrap">
            <div class="progress-bar" style="width: {overallProgress}%"></div>
          </div>
        </div>
      </div>

      <!-- Documents list -->
      <div class="docs-list">
        {#each documents as doc (doc.id)}
          <div class="doc-card" class:doc-rejected={doc.status === 'rejected'}>
            <div class="doc-icon" class:icon-verified={doc.status === 'verified'} class:icon-rejected={doc.status === 'rejected'}>
              <span class="material-symbols-outlined">{doc.icon}</span>
            </div>
            <div class="doc-info">
              <div class="doc-top">
                <div>
                  <div class="doc-name-row">
                    <h3>{doc.name}</h3>
                    {#if doc.required}
                      <span class="required-badge">Wajib</span>
                    {/if}
                  </div>
                  <p class="doc-desc">{doc.description}</p>
                  {#if doc.fileName}
                    <p class="doc-file">
                      <span class="material-symbols-outlined">attach_file</span>
                      {doc.fileName}
                      {#if doc.uploadedAt} · {doc.uploadedAt}{/if}
                    </p>
                  {/if}
                  {#if doc.rejectReason}
                    <p class="reject-reason">
                      <span class="material-symbols-outlined">error</span>
                      {doc.rejectReason}
                    </p>
                  {/if}
                </div>
                <div class="doc-actions">
                  <span class="doc-status" class:status-verified={doc.status === 'verified'} class:status-pending={doc.status === 'uploaded' || doc.status === 'pending'} class:status-rejected={doc.status === 'rejected'} class:status-required={doc.status === 'required'}>
                    {statusLabel(doc.status)}
                  </span>
                  {#if doc.status !== 'verified'}
                    {#if uploading === doc.id}
                      <div class="uploading-indicator">
                        <span class="material-symbols-outlined spin">refresh</span>
                        Mengupload...
                      </div>
                    {:else}
                      <label class="upload-btn">
                        <span class="material-symbols-outlined">upload</span>
                        {doc.fileName ? 'Ganti File' : 'Upload'}
                        <input type="file" accept=".jpg,.jpeg,.png,.pdf" onchange={(e) => handleFileInput(e, doc.id)} style="display:none" />
                      </label>
                      {#if doc.fileName}
                        <button class="remove-btn" type="button" onclick={() => removeDoc(doc.id)}>
                          <span class="material-symbols-outlined">delete</span>
                        </button>
                      {/if}
                    {/if}
                  {/if}
                </div>
              </div>

              {#if doc.status === 'required' || doc.status === 'rejected'}
                <!-- Drop zone -->
                <div
                  class="drop-zone"
                  class:drag-over={draggingOver === doc.id}
                  role="button"
                  tabindex="0"
                  ondragover={(e) => { e.preventDefault(); draggingOver = doc.id; }}
                  ondragleave={() => draggingOver = null}
                  ondrop={(e) => handleDrop(e, doc.id)}
                >
                  <span class="material-symbols-outlined">cloud_upload</span>
                  <p>Seret & lepas file di sini, atau <label class="inline-upload-label">klik untuk memilih <input type="file" accept=".jpg,.jpeg,.png,.pdf" onchange={(e) => handleFileInput(e, doc.id)} style="display:none" /></label></p>
                  <p class="drop-hint">Format: JPG, PNG, atau PDF · Maks. 5 MB</p>
                </div>
              {/if}
            </div>
          </div>
        {/each}
      </div>

      <div class="info-box">
        <span class="material-symbols-outlined">security</span>
        <p>Semua dokumen yang Anda unggah dienkripsi dan hanya dapat diakses oleh tim UmrohOS untuk keperluan pemrosesan visa dan dokumen perjalanan.</p>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .docs-root {
    padding-top: calc(5.2rem + 2rem);
    padding-bottom: 5rem;
    background: #fbf9f8;
    min-height: 100vh;
  }
  .shell {
    max-width: 72rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  .page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }
  .page-header h1 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .btn-wa {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    background: #006747;
    color: #fff;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.7rem 1.4rem;
    font-size: 0.88rem;
    flex-shrink: 0;
  }
  .btn-wa .material-symbols-outlined { font-size: 1rem; }
  /* Progress card */
  .progress-card {
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.8rem 2rem;
    border: 1px solid rgba(190,201,193,0.3);
    margin-bottom: 2rem;
  }
  .progress-stats {
    display: flex;
    gap: 2.5rem;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
  }
  .progress-stat strong {
    display: block;
    font-size: 1.8rem;
    font-weight: 800;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .progress-stat span { font-size: 0.8rem; color: #9ca3af; }
  .progress-header { display: flex; justify-content: space-between; font-size: 0.82rem; color: #6b7280; margin-bottom: 0.5rem; }
  .progress-bar-wrap { height: 10px; background: #e4e2e2; border-radius: 999px; overflow: hidden; }
  .progress-bar { height: 100%; background: linear-gradient(90deg, #004d34, #22c55e); border-radius: 999px; transition: width 0.5s ease; }
  /* Docs list */
  .docs-list { display: grid; gap: 1.2rem; margin-bottom: 2rem; }
  .doc-card {
    background: #fff;
    border-radius: 1.5rem;
    padding: 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    display: flex;
    gap: 1rem;
    align-items: flex-start;
  }
  .doc-card.doc-rejected { border-color: rgba(186,26,26,0.2); background: rgba(186,26,26,0.02); }
  .doc-icon {
    width: 3.2rem;
    height: 3.2rem;
    border-radius: 1rem;
    background: rgba(190,201,193,0.2);
    display: grid;
    place-items: center;
    flex-shrink: 0;
    color: #9ca3af;
  }
  .doc-icon.icon-verified { background: rgba(0,103,71,0.1); color: #006747; }
  .doc-icon.icon-rejected { background: rgba(186,26,26,0.1); color: #ba1a1a; }
  .doc-icon .material-symbols-outlined { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .doc-info { flex: 1; }
  .doc-top { display: flex; justify-content: space-between; gap: 1rem; align-items: flex-start; flex-wrap: wrap; margin-bottom: 0.5rem; }
  .doc-name-row { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.25rem; }
  .doc-name-row h3 { margin: 0; font-size: 0.95rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .required-badge {
    font-size: 0.7rem;
    font-weight: 700;
    background: rgba(234,88,12,0.1);
    color: #c2410c;
    border-radius: 999px;
    padding: 0.2rem 0.6rem;
  }
  .doc-desc { margin: 0; font-size: 0.82rem; color: #9ca3af; }
  .doc-file { margin: 0.4rem 0 0; display: flex; align-items: center; gap: 0.25rem; font-size: 0.8rem; color: #6b7280; }
  .doc-file .material-symbols-outlined { font-size: 0.9rem; }
  .reject-reason { margin: 0.5rem 0 0; display: flex; gap: 0.4rem; align-items: flex-start; font-size: 0.82rem; color: #ba1a1a; line-height: 1.4; }
  .reject-reason .material-symbols-outlined { font-size: 1rem; flex-shrink: 0; }
  .doc-actions { display: flex; align-items: center; gap: 0.5rem; flex-shrink: 0; flex-wrap: wrap; }
  .doc-status {
    font-size: 0.76rem;
    font-weight: 700;
    padding: 0.3rem 0.7rem;
    border-radius: 999px;
    background: rgba(190,201,193,0.2);
    color: #6b7280;
  }
  .doc-status.status-verified { background: rgba(0,103,71,0.1); color: #006747; }
  .doc-status.status-pending { background: rgba(119,90,25,0.1); color: #775a19; }
  .doc-status.status-rejected { background: rgba(186,26,26,0.1); color: #ba1a1a; }
  .doc-status.status-required { background: rgba(234,88,12,0.08); color: #c2410c; }
  .upload-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.82rem;
    font-weight: 700;
    color: #006747;
    background: rgba(0,103,71,0.08);
    border-radius: 999px;
    padding: 0.4rem 0.9rem;
    cursor: pointer;
    transition: background 0.15s;
  }
  .upload-btn:hover { background: rgba(0,103,71,0.15); }
  .upload-btn .material-symbols-outlined { font-size: 1rem; }
  .remove-btn {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    border: none;
    background: rgba(186,26,26,0.08);
    color: #ba1a1a;
    display: grid;
    place-items: center;
    cursor: pointer;
  }
  .remove-btn .material-symbols-outlined { font-size: 1rem; }
  .uploading-indicator { display: flex; align-items: center; gap: 0.35rem; font-size: 0.82rem; color: #775a19; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .spin { animation: spin 1s linear infinite; }
  /* Drop zone */
  .drop-zone {
    margin-top: 0.75rem;
    border: 2px dashed rgba(190,201,193,0.5);
    border-radius: 1rem;
    padding: 1.5rem;
    text-align: center;
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
  }
  .drop-zone.drag-over { border-color: #006747; background: rgba(0,103,71,0.04); }
  .drop-zone .material-symbols-outlined { font-size: 2rem; color: #d1d5db; display: block; margin-bottom: 0.5rem; }
  .drop-zone p { margin: 0; font-size: 0.85rem; color: #9ca3af; }
  .drop-zone p:first-of-type { color: #57534e; }
  .drop-hint { margin-top: 0.35rem !important; font-size: 0.76rem !important; }
  .inline-upload-label { color: #006747; font-weight: 700; cursor: pointer; text-decoration: underline; }
  /* Info box */
  .info-box {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    background: rgba(0,103,71,0.06);
    border-radius: 1rem;
    padding: 1.1rem 1.4rem;
  }
  .info-box .material-symbols-outlined {
    color: #006747;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
    flex-shrink: 0;
  }
  .info-box p { margin: 0; font-size: 0.85rem; color: #57534e; line-height: 1.6; }
</style>
