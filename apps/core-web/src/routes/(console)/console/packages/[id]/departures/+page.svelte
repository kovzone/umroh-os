<script lang="ts">
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  const pkg = $derived(data.package);
  const departures = $derived(data.departures ?? []);

  function statusLabel(status: string): string {
    if (status === 'open') return 'Buka';
    if (status === 'closed') return 'Tutup';
    if (status === 'draft') return 'Draft';
    return status;
  }

  function statusClass(status: string): string {
    if (status === 'open') return 'pill-open';
    if (status === 'closed') return 'pill-closed';
    return 'pill-draft';
  }

  function formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return dateStr;
    return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }

  function seatFill(remaining: number, total: number): string {
    if (!total) return '-';
    const filled = total - remaining;
    const pct = Math.round((filled / total) * 100);
    return `${filled}/${total} (${pct}% terisi)`;
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <a href="/console/packages" class="back-link">
      <span class="material-symbols-outlined">arrow_back</span>
      Katalog Paket
    </a>
    <span class="topbar-sep">/</span>
    <span class="topbar-pkg-name">{pkg?.name ?? 'Loading...'}</span>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Keberangkatan</h2>
        <p>
          Paket: <strong>{pkg?.name}</strong>
          &nbsp;·&nbsp;
          <a href="/console/packages/{pkg?.id}/edit" class="edit-link">
            <span class="material-symbols-outlined">edit</span> Edit paket
          </a>
        </p>
      </div>
      <a href="/console/packages/{pkg?.id}/departures/new" class="primary-btn">
        <span class="material-symbols-outlined">add</span>
        Tambah Departure
      </a>
    </div>

    <div class="panel">
      {#if departures.length === 0}
        <div class="empty-state">
          <span class="material-symbols-outlined empty-icon">flight_takeoff</span>
          <p class="empty-title">Belum ada keberangkatan</p>
          <p class="empty-hint">Tambahkan jadwal keberangkatan untuk paket ini.</p>
          <a href="/console/packages/{pkg?.id}/departures/new" class="primary-btn" style="margin-top:1rem">
            <span class="material-symbols-outlined">add</span>
            Tambah Departure
          </a>
        </div>
      {:else}
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Tgl Keberangkatan</th>
                <th>Tgl Kembali</th>
                <th>Kapasitas</th>
                <th>Sisa Kursi</th>
                <th>Status</th>
                <th class="align-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each departures as dep (dep.id)}
                <tr>
                  <td>
                    <div class="date-cell">
                      <span class="date-main">{formatDate(dep.departure_date)}</span>
                      <span class="dep-id">ID: {dep.id}</span>
                    </div>
                  </td>
                  <td>{formatDate(dep.return_date)}</td>
                  <td>{dep.total_seats} kursi</td>
                  <td>
                    <div class="seat-cell">
                      <span class:low-seats={dep.remaining_seats <= 5}>{dep.remaining_seats}</span>
                      <span class="seat-hint">{seatFill(dep.remaining_seats, dep.total_seats)}</span>
                    </div>
                  </td>
                  <td>
                    <span class="pill {statusClass(dep.status)}">{statusLabel(dep.status)}</span>
                  </td>
                  <td class="actions-cell">
                    <a
                      href="/console/packages/{pkg?.id}/departures/{dep.id}/edit"
                      class="action-btn ghost-btn"
                    >
                      <span class="material-symbols-outlined">edit</span>
                      Edit
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
    gap: 0.6rem;
    backdrop-filter: blur(8px);
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    color: #004ac6;
    text-decoration: none;
    font-size: 0.84rem;
    font-weight: 500;
    white-space: nowrap;
  }

  .back-link:hover { text-decoration: underline; }
  .back-link .material-symbols-outlined { font-size: 1rem; }

  .topbar-sep { color: #c3c6d7; }

  .topbar-pkg-name {
    font-size: 0.84rem;
    color: #434655;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .canvas { padding: 1.5rem; max-width: 100%; }

  .page-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1.25rem;
  }

  .page-head h2 {
    margin: 0;
    font-size: 1.5rem;
    color: #191c1e;
  }

  .page-head p {
    margin: 0.3rem 0 0;
    font-size: 0.82rem;
    color: #434655;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    flex-wrap: wrap;
  }

  .edit-link {
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
    color: #004ac6;
    text-decoration: none;
    font-size: 0.78rem;
    font-weight: 500;
  }

  .edit-link:hover { text-decoration: underline; }
  .edit-link .material-symbols-outlined { font-size: 0.9rem; }

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
    white-space: nowrap;
    flex-shrink: 0;
  }

  .primary-btn:hover { background: linear-gradient(90deg, #003a9e, #1d4ed8); }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

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
    padding: 4rem 2rem;
    text-align: center;
  }

  .empty-icon { font-size: 3rem; color: #c3c6d7; margin-bottom: 0.75rem; }
  .empty-title { margin: 0 0 0.25rem; font-size: 1rem; font-weight: 700; color: #191c1e; }
  .empty-hint { margin: 0; font-size: 0.82rem; color: #434655; }

  .table-wrap { overflow-x: auto; }

  table { width: 100%; border-collapse: collapse; }

  th, td {
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

  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .align-right { text-align: right; }

  .date-cell { display: flex; flex-direction: column; gap: 0.1rem; }
  .date-main { font-weight: 600; color: #191c1e; }
  .dep-id { font-size: 0.65rem; color: #737686; font-family: 'IBM Plex Mono', monospace; }

  .seat-cell { display: flex; flex-direction: column; gap: 0.1rem; }
  .low-seats { color: #dc2626; font-weight: 700; }
  .seat-hint { font-size: 0.68rem; color: #737686; }

  .pill {
    display: inline-flex;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.04em;
  }

  .pill-open { background: #dcfce7; color: #166534; }
  .pill-closed { background: #fef2f2; color: #dc2626; }
  .pill-draft { background: #e0e3e5; color: #434655; }

  .actions-cell {
    text-align: right;
    display: flex;
    gap: 0.4rem;
    justify-content: flex-end;
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

  .action-btn .material-symbols-outlined { font-size: 0.9rem; }

  .ghost-btn {
    background: #e6e8ea;
    color: #434655;
    border-color: rgb(195 198 215 / 0.55);
  }

  .ghost-btn:hover { background: #d6d9db; }
</style>
