<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let category = $state('');
  let departureId = $state('');
  let requiredBy = $state('');
  let maxBudget = $state('');
  let loading = $state(false);
  let error = $state('');
  let result = $state<any>(null);

  const CATEGORIES = [
    'Kain Ihram',
    'Batik Seragam',
    'Tas Koper',
    'Perlengkapan Haji',
    'Konsumsi',
    'Transportasi',
    'Akomodasi',
    'Perlengkapan Medis',
    'Souvenir',
    'Lainnya',
  ];

  async function submitAutoVendor(e: Event) {
    e.preventDefault();
    loading = true; error = ''; result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/logistics/auto-vendor`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          category,
          departure_id: departureId,
          required_by: requiredBy,
          max_budget: parseFloat(maxBudget),
        }),
      });
      if (!res.ok) throw new Error(`Gagal seleksi vendor (${res.status})`);
      result = await res.json();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    loading = false;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/logistics" class="bc-link">Logistik</a>
      <span class="bc-sep">/</span>
      <span>Seleksi Vendor Otomatis</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Seleksi Vendor Otomatis</h2>
      <p>BL-LOG-016 — Rekomendasikan vendor terbaik berdasarkan kategori dan anggaran</p>
    </div>

    <div class="layout">
      <div class="form-col">
        <div class="section-block">
          <h3 class="section-title">Parameter Seleksi</h3>
          <form class="form-grid" onsubmit={submitAutoVendor}>
            <div class="field">
              <label for="av-cat">Kategori</label>
              <select id="av-cat" bind:value={category} required>
                <option value="" disabled>Pilih kategori...</option>
                {#each CATEGORIES as cat}
                  <option>{cat}</option>
                {/each}
              </select>
            </div>
            <div class="field">
              <label for="av-dep">ID Keberangkatan</label>
              <input id="av-dep" type="text" placeholder="dep-001" bind:value={departureId} required />
            </div>
            <div class="field">
              <label for="av-date">Dibutuhkan Sebelum</label>
              <input id="av-date" type="date" bind:value={requiredBy} required />
            </div>
            <div class="field">
              <label for="av-budget">Anggaran Maksimal (Rp)</label>
              <input id="av-budget" type="number" min="0" placeholder="50000000" bind:value={maxBudget} required />
            </div>
            <div class="field field-wide field-actions">
              <button type="submit" class="btn-primary" disabled={loading}>
                {#if loading}<span class="spinner"></span>{/if}
                Cari Vendor Terbaik
              </button>
            </div>
          </form>
          {#if error}<div class="alert-err">{error}</div>{/if}
        </div>
      </div>

      <div class="result-col">
        {#if result}
          <div class="section-block vendor-card">
            <div class="vendor-header">
              <span class="material-symbols-outlined vendor-icon">storefront</span>
              <div>
                <div class="vendor-label">Vendor Direkomendasikan</div>
                <div class="vendor-name">{result.vendor_name ?? result.name ?? 'N/A'}</div>
              </div>
              <span class="chip chip-green ml-auto">Rekomendasi</span>
            </div>

            <div class="vendor-body">
              <div class="vendor-row">
                <span>Kategori</span>
                <strong>{result.category ?? category}</strong>
              </div>
              {#if result.vendor_id}
                <div class="vendor-row">
                  <span>Vendor ID</span>
                  <strong class="mono">{result.vendor_id}</strong>
                </div>
              {/if}
              <div class="vendor-row">
                <span>Estimasi Biaya</span>
                <strong class="text-blue">{formatIDR(result.estimated_cost ?? result.cost ?? 0)}</strong>
              </div>
              {#if result.rating}
                <div class="vendor-row">
                  <span>Rating</span>
                  <strong>{result.rating} / 5.0</strong>
                </div>
              {/if}
              {#if result.delivery_days}
                <div class="vendor-row">
                  <span>Estimasi Pengiriman</span>
                  <strong>{result.delivery_days} hari</strong>
                </div>
              {/if}
            </div>

            {#if result.reason}
              <div class="vendor-reason">
                <div class="reason-label">Alasan Rekomendasi</div>
                <p>{result.reason}</p>
              </div>
            {/if}

            {#if result.alternatives?.length > 0}
              <div class="vendor-alt">
                <div class="reason-label">Alternatif</div>
                {#each result.alternatives as alt}
                  <div class="alt-row">
                    <span>{alt.name ?? alt.vendor_name}</span>
                    <span>{formatIDR(alt.estimated_cost ?? alt.cost ?? 0)}</span>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {:else if !loading}
          <div class="empty-card">
            <span class="material-symbols-outlined empty-icon">storefront</span>
            <p>Isi form dan klik "Cari Vendor Terbaik" untuk melihat rekomendasi</p>
          </div>
        {:else}
          <div class="loading-card">
            <span class="spinner-lg"></span>
            <p>Menganalisis vendor...</p>
          </div>
        {/if}
      </div>
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 72rem; }
  .page-head { margin-bottom: 1.25rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .layout { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; align-items: start; }
  @media (max-width: 900px) { .layout { grid-template-columns: 1fr; } }
  .section-block { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; padding: 1.25rem; }
  .section-title { margin: 0 0 1rem; font-size: 0.9rem; font-weight: 700; }
  .form-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 0.75rem; align-items: end; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-wide { grid-column: span 2; }
  .field label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.45rem 0.65rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; }
  .field-actions { align-self: flex-end; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.5rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.85rem; height: 0.85rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  .spinner-lg { width: 2rem; height: 2rem; border: 3px solid rgb(37 99 235 / 0.2); border-top-color: #2563eb; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.6rem 0.85rem; font-size: 0.8rem; margin-top: 0.75rem; }
  .vendor-card { border-color: #93c5fd; }
  .vendor-header { display: flex; align-items: center; gap: 0.85rem; margin-bottom: 1rem; padding-bottom: 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .vendor-icon { font-size: 2rem; color: #2563eb; }
  .vendor-label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #737686; margin-bottom: 0.2rem; }
  .vendor-name { font-size: 1.1rem; font-weight: 700; color: #191c1e; }
  .chip { display: inline-flex; padding: 0.12rem 0.45rem; border-radius: 0.2rem; font-size: 0.65rem; font-weight: 600; }
  .chip-green { background: #dcfce7; color: #166534; }
  .ml-auto { margin-left: auto; }
  .vendor-body { display: flex; flex-direction: column; gap: 0.5rem; margin-bottom: 1rem; }
  .vendor-row { display: flex; justify-content: space-between; align-items: center; font-size: 0.82rem; padding: 0.35rem 0; border-bottom: 1px solid rgb(195 198 215 / 0.2); }
  .vendor-row span { color: #737686; }
  .mono { font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem; }
  .text-blue { color: #2563eb; }
  .vendor-reason { background: #f8fafc; border-radius: 0.35rem; padding: 0.75rem; margin-bottom: 0.75rem; }
  .reason-label { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #737686; margin-bottom: 0.35rem; }
  .vendor-reason p { margin: 0; font-size: 0.8rem; color: #434655; line-height: 1.5; }
  .vendor-alt { }
  .alt-row { display: flex; justify-content: space-between; padding: 0.35rem 0; font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.2); color: #737686; }
  .empty-card, .loading-card { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 0.75rem; padding: 3rem 1rem; background: #fff; border: 1px dashed rgb(195 198 215 / 0.55); border-radius: 0.5rem; color: #b0b3c1; text-align: center; }
  .empty-icon { font-size: 2.5rem; }
  .empty-card p, .loading-card p { margin: 0; font-size: 0.82rem; max-width: 240px; }
</style>
