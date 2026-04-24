<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';
  import { onMount } from 'svelte';
  import { getCatalogPackages } from '$lib/features/s1-catalog/repository';
  import type { PackageCard } from '$lib/features/s1-catalog/types';

  let packages = $state<PackageCard[]>([]);
  let loading = $state(true);
  let filterMonth = $state('all');
  let filterStatus = $state('all');

  onMount(async () => {
    try {
      packages = await getCatalogPackages();
    } finally {
      loading = false;
    }
  });

  // Generate upcoming departure rows from package data
  const departures = $derived.by(() => {
    const rows: Array<{
      id: string;
      packageName: string;
      packageId: string;
      departureDate: string;
      returnDate: string;
      seats: number;
      status: 'open' | 'full' | 'soon';
    }> = [];

    for (const pkg of packages) {
      if (pkg.nextDepartureLabel && pkg.nextDepartureLabel !== 'Belum ada keberangkatan terjadwal') {
        const parts = pkg.nextDepartureLabel.split(' s.d. ');
        rows.push({
          id: pkg.id,
          packageName: pkg.name,
          packageId: pkg.id,
          departureDate: parts[0] ?? pkg.nextDepartureLabel,
          returnDate: parts[1] ?? '',
          seats: pkg.remainingSeats,
          status: pkg.remainingSeats <= 0 ? 'full' : pkg.remainingSeats <= 5 ? 'soon' : 'open'
        });
      }
    }

    // Add demo rows if empty
    if (rows.length === 0) {
      rows.push(
        { id: 'd1', packageName: 'Paket Gold 9 Hari', packageId: 'demo', departureDate: '2025-02-15', returnDate: '2025-02-23', seats: 12, status: 'open' },
        { id: 'd2', packageName: 'Paket Platinum Ramadhan', packageId: 'demo', departureDate: '2025-03-28', returnDate: '2025-04-09', seats: 4, status: 'soon' },
        { id: 'd3', packageName: 'Paket Plus Turki', packageId: 'demo', departureDate: '2025-04-10', returnDate: '2025-04-25', seats: 18, status: 'open' },
        { id: 'd4', packageName: 'Paket Silver 9 Hari', packageId: 'demo', departureDate: '2025-01-22', returnDate: '2025-01-30', seats: 0, status: 'full' },
        { id: 'd5', packageName: 'Umrah Khusus Idul Adha', packageId: 'demo', departureDate: '2025-06-05', returnDate: '2025-06-15', seats: 25, status: 'open' }
      );
    }

    return rows.filter(r => {
      if (filterStatus === 'open' && r.status === 'full') return false;
      if (filterStatus === 'full' && r.status !== 'full') return false;
      return true;
    });
  });

  function formatDate(dateStr: string): string {
    if (!dateStr || dateStr.includes('-') === false) return dateStr;
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return dateStr;
    return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' });
  }

  function statusLabel(status: string): string {
    if (status === 'full') return 'Penuh';
    if (status === 'soon') return 'Hampir Penuh';
    return 'Masih Ada';
  }
</script>

<svelte:head>
  <title>Jadwal Keberangkatan Umrah — UmrohOS</title>
  <meta name="description" content="Pantau jadwal keberangkatan umrah terbaru, sisa kursi, dan status pendaftaran dari UmrohOS." />
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="board-root">

    <section class="board-hero">
      <div class="shell">
        <p class="kicker">Jadwal Keberangkatan</p>
        <h1>Papan Informasi Keberangkatan</h1>
        <p class="hero-sub">Pantau jadwal keberangkatan terbaru, ketersediaan kursi, dan status pendaftaran secara real-time.</p>
        <div class="hero-stats">
          <div class="stat">
            <span class="material-symbols-outlined">flight_takeoff</span>
            <strong>{departures.filter(d => d.status !== 'full').length}</strong>
            <span>Keberangkatan Tersedia</span>
          </div>
          <div class="stat">
            <span class="material-symbols-outlined">event_seat</span>
            <strong>{departures.reduce((acc, d) => acc + d.seats, 0)}</strong>
            <span>Total Kursi Tersisa</span>
          </div>
          <div class="stat">
            <span class="material-symbols-outlined">update</span>
            <strong>Real-time</strong>
            <span>Data Diperbarui</span>
          </div>
        </div>
      </div>
    </section>

    <section class="board-table-section">
      <div class="shell">
        <!-- Filters -->
        <div class="filters-row">
          <div class="filter-group">
            <label>Status Kursi:</label>
            <div class="filter-chips">
              {#each [['all', 'Semua'], ['open', 'Tersedia'], ['full', 'Penuh']] as [val, label] (val)}
                <button
                  type="button"
                  class="filter-chip"
                  class:active={filterStatus === val}
                  onclick={() => filterStatus = val}
                >{label}</button>
              {/each}
            </div>
          </div>
          <p class="update-note">
            <span class="material-symbols-outlined">info</span>
            Data diperbarui secara otomatis. Ketersediaan kursi dapat berubah sewaktu-waktu.
          </p>
        </div>

        <!-- Table -->
        {#if loading}
          <div class="loading-state">
            <span class="material-symbols-outlined spin">refresh</span>
            <p>Memuat data keberangkatan...</p>
          </div>
        {:else if departures.length === 0}
          <div class="empty-state">
            <span class="material-symbols-outlined">flight_takeoff</span>
            <p>Tidak ada keberangkatan yang cocok dengan filter ini.</p>
          </div>
        {:else}
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Paket Umrah</th>
                  <th>Tanggal Berangkat</th>
                  <th>Tanggal Kembali</th>
                  <th>Sisa Kursi</th>
                  <th>Status</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                {#each departures as dep (dep.id)}
                  <tr class:row-full={dep.status === 'full'}>
                    <td class="pkg-cell">
                      <span class="material-symbols-outlined pkg-icon">mosque</span>
                      <span>{dep.packageName}</span>
                    </td>
                    <td>
                      <span class="material-symbols-outlined date-icon">flight_takeoff</span>
                      {formatDate(dep.departureDate)}
                    </td>
                    <td>
                      <span class="material-symbols-outlined date-icon">flight_land</span>
                      {dep.returnDate ? formatDate(dep.returnDate) : '—'}
                    </td>
                    <td>
                      <span class="seats-badge" class:seats-low={dep.seats <= 5 && dep.seats > 0} class:seats-zero={dep.seats === 0}>
                        {dep.seats === 0 ? 'Penuh' : dep.seats + ' kursi'}
                      </span>
                    </td>
                    <td>
                      <span class="status-chip" class:chip-open={dep.status === 'open'} class:chip-soon={dep.status === 'soon'} class:chip-full={dep.status === 'full'}>
                        {statusLabel(dep.status)}
                      </span>
                    </td>
                    <td>
                      {#if dep.status !== 'full'}
                        <a class="action-btn" href="/packages/{dep.packageId}">Daftar</a>
                      {:else}
                        <span class="action-disabled">Ditutup</span>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        <div class="cta-row">
          <p>Tidak melihat jadwal yang sesuai? Hubungi tim kami untuk informasi keberangkatan terbaru.</p>
          <div class="cta-btns">
            <a class="btn-wa" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
              <span class="material-symbols-outlined">chat</span>
              Tanya via WhatsApp
            </a>
            <a class="btn-packages" href="/packages">Lihat Semua Paket</a>
          </div>
        </div>
      </div>
    </section>

  </div>
</MarketingPageLayout>

<style>
  .board-root {
    padding-top: 5.2rem;
    background: #fbf9f8;
  }
  .shell {
    max-width: 80rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  /* Hero */
  .board-hero {
    padding: 4rem 0 3rem;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
    text-align: center;
    color: #fff;
  }
  .kicker {
    display: inline-block;
    margin: 0 0 1rem;
    background: rgba(254,212,136,0.25);
    color: #fed488;
    border-radius: 999px;
    padding: 0.4rem 1rem;
    font-size: 0.76rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .board-hero h1 {
    margin: 0;
    font-size: clamp(1.8rem, 3.5vw, 3rem);
    font-weight: 800;
    letter-spacing: -0.02em;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .hero-sub {
    margin: 0.8rem auto 2rem;
    max-width: 38rem;
    opacity: 0.85;
    font-size: 1.05rem;
  }
  .hero-stats {
    display: flex;
    justify-content: center;
    gap: 3rem;
    flex-wrap: wrap;
  }
  .stat {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
  }
  .stat .material-symbols-outlined {
    font-size: 1.4rem;
    opacity: 0.7;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .stat strong {
    font-size: 1.6rem;
    font-weight: 800;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .stat span:last-child {
    font-size: 0.8rem;
    opacity: 0.8;
  }
  /* Table section */
  .board-table-section {
    padding: 3rem 0 5rem;
  }
  .filters-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
    margin-bottom: 2rem;
  }
  .filter-group {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  .filter-group label {
    font-size: 0.85rem;
    font-weight: 600;
    color: #57534e;
  }
  .filter-chips {
    display: flex;
    gap: 0.4rem;
  }
  .filter-chip {
    border: 1px solid rgba(190,201,193,0.5);
    border-radius: 999px;
    padding: 0.35rem 0.9rem;
    background: #fff;
    color: #57534e;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
  }
  .filter-chip.active {
    background: #006747;
    color: #fff;
    border-color: transparent;
  }
  .update-note {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.78rem;
    color: #9ca3af;
    margin: 0;
  }
  .update-note .material-symbols-outlined {
    font-size: 0.9rem;
  }
  /* Table */
  .loading-state, .empty-state {
    text-align: center;
    padding: 5rem 0;
    color: #9ca3af;
  }
  .loading-state .material-symbols-outlined, .empty-state .material-symbols-outlined {
    font-size: 3rem;
    display: block;
    margin-bottom: 0.75rem;
  }
  @keyframes spin { to { transform: rotate(360deg); } }
  .spin { animation: spin 1s linear infinite; }
  .table-wrap {
    overflow-x: auto;
    border-radius: 1.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    box-shadow: 0 4px 16px rgba(0,0,0,0.04);
  }
  table {
    width: 100%;
    border-collapse: collapse;
    background: #fff;
    border-radius: 1.5rem;
    overflow: hidden;
  }
  thead tr {
    background: #f8faf9;
    border-bottom: 1px solid rgba(190,201,193,0.25);
  }
  th {
    padding: 0.9rem 1.2rem;
    text-align: left;
    font-size: 0.78rem;
    font-weight: 700;
    color: #6b7280;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    white-space: nowrap;
  }
  tbody tr {
    border-bottom: 1px solid rgba(190,201,193,0.15);
    transition: background 0.15s;
  }
  tbody tr:last-child { border-bottom: none; }
  tbody tr:hover { background: #fafaf9; }
  tbody tr.row-full { opacity: 0.6; }
  td {
    padding: 1.1rem 1.2rem;
    font-size: 0.9rem;
    color: #1b1c1c;
    white-space: nowrap;
  }
  .pkg-cell {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    font-weight: 600;
  }
  .pkg-icon {
    font-size: 1.1rem;
    color: #006747;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .date-icon {
    font-size: 0.95rem;
    color: #9ca3af;
    margin-right: 0.3rem;
    vertical-align: middle;
  }
  .seats-badge {
    display: inline-block;
    font-size: 0.82rem;
    font-weight: 700;
    padding: 0.3rem 0.7rem;
    border-radius: 999px;
    background: rgba(0,103,71,0.08);
    color: #006747;
  }
  .seats-badge.seats-low {
    background: rgba(234,88,12,0.1);
    color: #c2410c;
  }
  .seats-badge.seats-zero {
    background: rgba(186,26,26,0.1);
    color: #ba1a1a;
  }
  .status-chip {
    display: inline-block;
    font-size: 0.78rem;
    font-weight: 700;
    padding: 0.3rem 0.8rem;
    border-radius: 999px;
  }
  .chip-open { background: rgba(0,103,71,0.1); color: #006747; }
  .chip-soon { background: rgba(234,88,12,0.1); color: #c2410c; }
  .chip-full { background: rgba(186,26,26,0.1); color: #ba1a1a; }
  .action-btn {
    display: inline-flex;
    align-items: center;
    text-decoration: none;
    background: #004d34;
    color: #fff;
    font-size: 0.82rem;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.45rem 1.1rem;
    transition: opacity 0.15s;
  }
  .action-btn:hover { opacity: 0.85; }
  .action-disabled {
    font-size: 0.82rem;
    color: #9ca3af;
    font-weight: 600;
  }
  .cta-row {
    margin-top: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1.5rem;
    flex-wrap: wrap;
    padding: 1.8rem 2rem;
    background: #f0f9f4;
    border-radius: 1.5rem;
  }
  .cta-row p {
    margin: 0;
    color: #57534e;
    font-size: 0.95rem;
  }
  .cta-btns {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  .btn-wa {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    background: #006747;
    color: #fff;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.7rem 1.5rem;
    font-size: 0.88rem;
  }
  .btn-wa .material-symbols-outlined { font-size: 1rem; }
  .btn-packages {
    text-decoration: none;
    border: 1.5px solid #006747;
    color: #006747;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.7rem 1.5rem;
    font-size: 0.88rem;
  }
</style>
