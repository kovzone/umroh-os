<script lang="ts">
  const DEMO_PACKAGES = [
    { id: 'pkg1', name: 'Umroh Ramadhan Premium', departure_date: '2026-03-01', duration_days: 15 },
    { id: 'pkg2', name: 'Umroh Reguler April', departure_date: '2026-04-10', duration_days: 12 },
    { id: 'pkg3', name: 'Umroh Plus Turki', departure_date: '2026-05-05', duration_days: 18 },
    { id: 'pkg4', name: 'Umroh Hemat Juni', departure_date: '2026-06-15', duration_days: 10 },
    { id: 'pkg5', name: 'Umroh VIP Juli', departure_date: '2026-07-01', duration_days: 14 }
  ];

  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' });
  }
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">hub</span>
      <span class="sep">/</span>
      <a href="/console/crm" class="crumb-link">CRM Tools</a>
      <span class="sep">/</span>
      <span class="topbar-current">Itinerary Viewer</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn" title="Notifikasi">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar" aria-label="Profile">AD</button>
    </div>
  </header>

  <section class="canvas-page">
    <div class="page-head">
      <div>
        <h2>Itinerary Viewer</h2>
        <p>BL-CRM-015 — Pilih paket untuk melihat itinerary lengkap yang dapat dibagikan</p>
      </div>
    </div>

    <div class="pkg-list">
      {#each DEMO_PACKAGES as pkg}
        <a href="/console/crm/itinerary/{pkg.id}" class="pkg-card">
          <div class="pkg-icon">
            <span class="material-symbols-outlined">route</span>
          </div>
          <div class="pkg-body">
            <div class="pkg-name">{pkg.name}</div>
            <div class="pkg-meta">
              <span class="material-symbols-outlined">calendar_today</span>
              {formatDate(pkg.departure_date)}
              <span class="dot">·</span>
              <span class="material-symbols-outlined">schedule</span>
              {pkg.duration_days} hari
            </div>
          </div>
          <span class="material-symbols-outlined arrow">arrow_forward</span>
        </a>
      {/each}
    </div>
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }

  .topbar {
    position: sticky; top: 0; z-index: 30;
    height: 4rem;
    background: rgb(255 255 255 / 0.9);
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem;
    display: flex; align-items: center; justify-content: space-between; gap: 1rem;
    backdrop-filter: blur(8px);
  }

  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.88rem; color: #434655; }
  .breadcrumb-icon { font-size: 1.1rem; color: #004ac6; }
  .sep { color: #b0b3c1; }
  .crumb-link { color: #434655; text-decoration: none; }
  .crumb-link:hover { color: #004ac6; }
  .topbar-current { font-weight: 600; color: #191c1e; }
  .top-actions { display: flex; align-items: center; gap: 0.35rem; }
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas-page { padding: 1.5rem; max-width: 56rem; }
  .page-head { margin-bottom: 1.5rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .pkg-list { display: flex; flex-direction: column; gap: 0.65rem; }

  .pkg-card {
    display: flex; align-items: center; gap: 1rem;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.35rem;
    padding: 1rem 1.1rem;
    text-decoration: none;
    color: inherit;
    transition: border-color 0.1s, box-shadow 0.1s;
  }

  .pkg-card:hover { border-color: #2563eb; box-shadow: 0 2px 8px rgb(37 99 235 / 0.08); }

  .pkg-icon {
    width: 2.5rem; height: 2.5rem;
    background: #dbeafe; color: #1e3a8a;
    border-radius: 0.35rem;
    display: grid; place-items: center; flex-shrink: 0;
  }
  .pkg-icon .material-symbols-outlined { font-size: 1.3rem; }

  .pkg-body { flex: 1; }
  .pkg-name { font-size: 0.9rem; font-weight: 700; color: #191c1e; margin-bottom: 0.25rem; }
  .pkg-meta {
    display: flex; align-items: center; gap: 0.3rem;
    font-size: 0.75rem; color: #434655;
  }
  .pkg-meta .material-symbols-outlined { font-size: 0.8rem; }
  .dot { color: #b0b3c1; }

  .arrow { font-size: 1.1rem; color: #b0b3c1; }
  .pkg-card:hover .arrow { color: #2563eb; }
</style>
