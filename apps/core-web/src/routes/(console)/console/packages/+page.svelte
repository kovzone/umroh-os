<script lang="ts">
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  const packages = $derived(data.packages ?? []);
  const error = $derived(data.error ?? null);

  function formatPrice(amount: number, currency: string): string {
    const symbol = currency === 'IDR' ? 'Rp' : currency;
    return `${symbol} ${new Intl.NumberFormat('id-ID').format(amount)}`;
  }

  function statusLabel(status: string | undefined): string {
    if (status === 'active') return 'Aktif';
    if (status === 'draft') return 'Draft';
    return status ?? '-';
  }

  function statusClass(status: string | undefined): string {
    if (status === 'active') return 'pill-active';
    if (status === 'draft') return 'pill-draft';
    return 'pill-draft';
  }

  function kindLabel(kind: string): string {
    if (kind === 'umroh') return 'Umroh';
    if (kind === 'hajj') return 'Haji';
    if (kind === 'ziarah') return 'Ziarah';
    return kind;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <div class="search-wrap">
      <span class="material-symbols-outlined search-icon">search</span>
      <input type="text" placeholder="Cari paket..." />
    </div>
    <div class="top-actions">
      <div class="icon-actions">
        <button class="icon-btn"><span class="material-symbols-outlined">notifications</span></button>
        <button class="avatar" aria-label="Profile">A</button>
      </div>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Katalog Paket</h2>
        <p>Kelola paket umroh, haji, dan ziarah yang ditawarkan</p>
      </div>
      <div class="head-actions">
        <a href="/console/packages/new" class="primary-btn">
          <span class="material-symbols-outlined">add</span>
          <span>Tambah Package</span>
        </a>
      </div>
    </div>

    {#if error}
      <div class="alert-error" role="alert">
        <span class="material-symbols-outlined">error</span>
        <p>{error}</p>
      </div>
    {/if}

    <div class="panel">
      {#if packages.length === 0 && !error}
        <div class="empty-state">
          <span class="material-symbols-outlined empty-icon">inventory_2</span>
          <p class="empty-title">Belum ada paket</p>
          <p class="empty-hint">Mulai dengan menambahkan paket pertama Anda.</p>
          <a href="/console/packages/new" class="primary-btn" style="margin-top:1rem">
            <span class="material-symbols-outlined">add</span>
            <span>Tambah Package</span>
          </a>
        </div>
      {:else}
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Nama Paket</th>
                <th>Kategori</th>
                <th>Harga Mulai</th>
                <th>Status</th>
                <th>Keberangkatan</th>
                <th class="align-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each packages as pkg (pkg.id)}
                <tr>
                  <td>
                    <div class="pkg-name-cell">
                      <span class="pkg-name">{pkg.name}</span>
                      <span class="pkg-id">ID: {pkg.id}</span>
                    </div>
                  </td>
                  <td>
                    <span class="pill pill-kind">{kindLabel(pkg.kind)}</span>
                  </td>
                  <td class="price-cell">
                    {formatPrice(pkg.starting_price.list_amount, pkg.starting_price.list_currency)}
                  </td>
                  <td>
                    <span class="pill {statusClass(pkg.status)}">{statusLabel(pkg.status)}</span>
                  </td>
                  <td>
                    {#if pkg.next_departure}
                      <div class="dep-cell">
                        <span>{pkg.next_departure.departure_date}</span>
                        <span class="dep-seats">{pkg.next_departure.remaining_seats} kursi tersisa</span>
                      </div>
                    {:else}
                      <span class="muted">Belum terjadwal</span>
                    {/if}
                  </td>
                  <td class="actions-cell">
                    <a href="/console/packages/{pkg.id}/edit" class="action-btn ghost-btn">
                      <span class="material-symbols-outlined">edit</span>
                      Edit
                    </a>
                    <a href="/console/packages/{pkg.id}/departures" class="action-btn primary-btn-sm">
                      <span class="material-symbols-outlined">flight_takeoff</span>
                      Departures
                    </a>
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

<style>
  .page-shell {
    min-height: 100vh;
    background: #f7f9fb;
    font-family: 'IBM Plex Sans', ui-sans-serif, system-ui, -apple-system, sans-serif;
  }

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
    max-width: 28rem;
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
    font-family: inherit;
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
    display: grid;
    place-items: center;
  }

  .icon-btn:hover {
    background: #eceef0;
  }

  .avatar {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #b4c5ff;
    color: #00174b;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
  }

  .canvas {
    padding: 1.5rem;
    max-width: 100%;
  }

  .page-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1.25rem;
  }

  .page-head h2 {
    margin: 0;
    font-size: 1.5rem;
    line-height: 1.2;
    color: #191c1e;
  }

  .page-head p {
    margin: 0.3rem 0 0;
    font-size: 0.82rem;
    color: #434655;
  }

  .head-actions {
    display: flex;
    gap: 0.6rem;
    align-items: center;
  }

  .primary-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    border: 1px solid #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
    border-radius: 0.25rem;
    padding: 0.55rem 0.9rem;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    text-decoration: none;
    font-family: inherit;
  }

  .primary-btn:hover {
    background: linear-gradient(90deg, #003a9e, #1d4ed8);
  }

  .primary-btn .material-symbols-outlined {
    font-size: 1rem;
  }

  .alert-error {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.85rem 1rem;
    background: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 0.25rem;
    margin-bottom: 1rem;
    color: #dc2626;
    font-size: 0.84rem;
  }

  .alert-error .material-symbols-outlined {
    font-size: 1.1rem;
    flex-shrink: 0;
  }

  .alert-error p {
    margin: 0;
  }

  .panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    text-align: center;
  }

  .empty-icon {
    font-size: 3rem;
    color: #c3c6d7;
    margin-bottom: 0.75rem;
  }

  .empty-title {
    margin: 0 0 0.25rem;
    font-size: 1rem;
    font-weight: 700;
    color: #191c1e;
  }

  .empty-hint {
    margin: 0;
    font-size: 0.82rem;
    color: #434655;
  }

  .table-wrap {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 0.7rem 1rem;
    text-align: left;
    font-size: 0.82rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    white-space: nowrap;
  }

  th {
    background: #e6e8ea;
    text-transform: uppercase;
    font-size: 0.62rem;
    letter-spacing: 0.08em;
    color: #434655;
    font-weight: 700;
  }

  tbody tr:hover {
    background: #f7f9fb;
  }

  tbody tr:last-child td {
    border-bottom: 0;
  }

  .align-right {
    text-align: right;
  }

  .pkg-name-cell {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
    min-width: 12rem;
  }

  .pkg-name {
    font-weight: 600;
    color: #191c1e;
    white-space: normal;
    line-height: 1.3;
  }

  .pkg-id {
    font-size: 0.65rem;
    color: #737686;
    font-family: 'IBM Plex Mono', monospace;
  }

  .pill {
    display: inline-flex;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.04em;
  }

  .pill-active {
    background: #dcfce7;
    color: #166534;
  }

  .pill-draft {
    background: #e0e3e5;
    color: #434655;
  }

  .pill-kind {
    background: rgb(37 99 235 / 0.1);
    color: #004ac6;
  }

  .price-cell {
    font-weight: 600;
    color: #191c1e;
  }

  .dep-cell {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
  }

  .dep-seats {
    font-size: 0.68rem;
    color: #737686;
  }

  .muted {
    color: #737686;
    font-size: 0.78rem;
    font-style: italic;
  }

  .actions-cell {
    text-align: right;
    display: flex;
    gap: 0.4rem;
    justify-content: flex-end;
    align-items: center;
  }

  .action-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    border-radius: 0.2rem;
    font-size: 0.72rem;
    padding: 0.3rem 0.55rem;
    cursor: pointer;
    text-decoration: none;
    font-weight: 600;
    font-family: inherit;
    border: 1px solid transparent;
  }

  .action-btn .material-symbols-outlined {
    font-size: 0.9rem;
  }

  .ghost-btn {
    background: #e6e8ea;
    color: #434655;
    border-color: rgb(195 198 215 / 0.55);
  }

  .ghost-btn:hover {
    background: #d6d9db;
  }

  .primary-btn-sm {
    background: #004ac6;
    color: #fff;
    border-color: #004ac6;
  }

  .primary-btn-sm:hover {
    background: #003a9e;
  }
</style>
