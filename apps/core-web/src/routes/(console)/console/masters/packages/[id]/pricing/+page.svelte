<script lang="ts">
  import type { PageData, DeparturePricing, RoomPricing } from './+page.server';

  let { data }: { data: PageData } = $props();

  let departures = $state<DeparturePricing[]>([]);
  let packageName = $state('');

  $effect(() => {
    departures = data.departures ?? [];
    packageName = data.package_name ?? '';
  });

  // expanded departures
  let expandedDeps = $state<Set<string>>(new Set());

  function toggleDep(id: string) {
    const next = new Set(expandedDeps);
    if (next.has(id)) { next.delete(id); } else { next.add(id); }
    expandedDeps = next;
  }

  // editing state: Map<departure_id, Map<room_type, price>>
  let editPrices = $state<Record<string, Record<string, number>>>({});
  let savingDeps = $state<Set<string>>(new Set());
  let savedDeps = $state<Set<string>>(new Set());

  $effect(() => {
    const prices: Record<string, Record<string, number>> = {};
    for (const dep of departures) {
      prices[dep.departure_id] = {};
      for (const p of dep.pricing) {
        prices[dep.departure_id][p.room_type] = p.price_idr;
      }
    }
    editPrices = prices;
  });

  function updatePrice(depId: string, roomType: string, value: number) {
    editPrices = {
      ...editPrices,
      [depId]: { ...editPrices[depId], [roomType]: value }
    };
    // clear saved indicator
    const next = new Set(savedDeps);
    next.delete(depId);
    savedDeps = next;
  }

  async function savePricing(depId: string) {
    const saving = new Set(savingDeps);
    saving.add(depId);
    savingDeps = saving;
    try {
      // Mock save
      await new Promise((r) => setTimeout(r, 500));
      // Update local state
      departures = departures.map((dep) => {
        if (dep.departure_id !== depId) return dep;
        return {
          ...dep,
          pricing: dep.pricing.map((p) => ({
            ...p,
            price_idr: editPrices[depId]?.[p.room_type] ?? p.price_idr
          }))
        };
      });
      const saved = new Set(savedDeps);
      saved.add(depId);
      savedDeps = saved;
    } finally {
      const s = new Set(savingDeps);
      s.delete(depId);
      savingDeps = s;
    }
  }

  function formatIDR(n: number): string {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  function formatDate(d: string): string {
    return new Date(d).toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' });
  }

  const ROOM_LABELS: Record<string, string> = { double: 'Double', triple: 'Triple', quad: 'Quad' };
</script>

<main class="page-shell">
  <header class="topbar">
    <div class="breadcrumb">
      <a href="/console/masters" class="bc-link">Master Data</a>
      <span class="material-symbols-outlined bc-sep">chevron_right</span>
      <span class="bc-text">{packageName || 'Paket'}</span>
      <span class="material-symbols-outlined bc-sep">chevron_right</span>
      <span class="bc-text">Pricing Keberangkatan</span>
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
        <h2>{packageName}</h2>
        <p>Kelola harga per tipe kamar untuk setiap keberangkatan</p>
      </div>
    </div>

    {#if data.error}
      <div class="error-banner">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    {#if departures.length === 0}
      <div class="empty-full">
        <span class="material-symbols-outlined">flight_takeoff</span>
        <p>Belum ada keberangkatan untuk paket ini.</p>
      </div>
    {:else}
      <div class="dep-list">
        {#each departures as dep (dep.departure_id)}
          {@const isExpanded = expandedDeps.has(dep.departure_id)}
          {@const isSaving = savingDeps.has(dep.departure_id)}
          {@const isSaved = savedDeps.has(dep.departure_id)}
          <div class="dep-card">
            <button
              type="button"
              class="dep-header"
              onclick={() => toggleDep(dep.departure_id)}
              aria-expanded={isExpanded}
            >
              <div class="dep-info">
                <span class="material-symbols-outlined dep-icon">flight_takeoff</span>
                <div>
                  <div class="dep-date">{formatDate(dep.departure_date)}</div>
                  <div class="dep-airline">{dep.airline}</div>
                </div>
              </div>
              <span class="material-symbols-outlined chevron" class:rotated={isExpanded}>expand_more</span>
            </button>

            {#if isExpanded}
              <div class="dep-body">
                <table class="pricing-table">
                  <thead>
                    <tr>
                      <th>Tipe Kamar</th>
                      <th>Harga (IDR)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {#each dep.pricing as p (p.room_type)}
                      <tr>
                        <td>
                          <span class="room-badge room-badge--{p.room_type}">{ROOM_LABELS[p.room_type]}</span>
                        </td>
                        <td>
                          <input
                            type="number"
                            class="price-input"
                            value={editPrices[dep.departure_id]?.[p.room_type] ?? p.price_idr}
                            min="0"
                            step="500000"
                            oninput={(e) => updatePrice(dep.departure_id, p.room_type, Number((e.target as HTMLInputElement).value))}
                          />
                          <span class="price-preview">{formatIDR(editPrices[dep.departure_id]?.[p.room_type] ?? p.price_idr)}</span>
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </table>

                <div class="dep-footer">
                  {#if isSaved}
                    <span class="saved-msg">
                      <span class="material-symbols-outlined">check_circle</span>
                      Harga berhasil disimpan
                    </span>
                  {/if}
                  <button
                    type="button"
                    class="primary-btn"
                    onclick={() => savePricing(dep.departure_id)}
                    disabled={isSaving}
                  >
                    {#if isSaving}
                      <span class="material-symbols-outlined spin">progress_activity</span>
                      Menyimpan...
                    {:else}
                      <span class="material-symbols-outlined">save</span>
                      Simpan Harga
                    {/if}
                  </button>
                </div>
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }

  .topbar {
    position: sticky; top: 0; z-index: 30; height: 4rem;
    background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center;
    justify-content: space-between; gap: 1rem; backdrop-filter: blur(8px);
  }

  .breadcrumb { display: flex; align-items: center; gap: 0.3rem; font-size: 0.82rem; }
  .bc-link { color: #004ac6; text-decoration: none; font-weight: 600; }
  .bc-link:hover { text-decoration: underline; }
  .bc-sep { font-size: 1rem; color: #b0b3c1; }
  .bc-text { font-weight: 600; color: #191c1e; }

  .top-actions { display: flex; align-items: center; gap: 0.5rem; }
  .icon-btn { border: 0; background: transparent; color: #434655; width: 2rem; height: 2rem; border-radius: 0.25rem; cursor: pointer; position: relative; display: grid; place-items: center; }
  .icon-btn:hover { background: #eceef0; }
  .dot { position: absolute; width: 0.46rem; height: 0.46rem; border-radius: 999px; background: #ba1a1a; right: 0.4rem; top: 0.35rem; border: 2px solid #fff; }
  .avatar { border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff; color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem; font-weight: 700; font-size: 0.65rem; cursor: pointer; }

  .canvas { padding: 1.5rem; max-width: 72rem; }

  .page-head { margin-bottom: 1.5rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .error-banner { display: flex; align-items: center; gap: 0.5rem; background: #ffdad6; color: #93000a; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.82rem; margin-bottom: 1rem; }

  .empty-full { display: flex; flex-direction: column; align-items: center; gap: 0.6rem; padding: 4rem 1rem; color: #b0b3c1; }
  .empty-full .material-symbols-outlined { font-size: 3rem; }
  .empty-full p { margin: 0; font-size: 0.9rem; }

  .dep-list { display: flex; flex-direction: column; gap: 0.75rem; }

  .dep-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.25rem; overflow: hidden; }

  .dep-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 0.85rem 1rem; background: transparent; border: 0;
    cursor: pointer; width: 100%; text-align: left;
  }
  .dep-header:hover { background: #f7f9fb; }

  .dep-info { display: flex; align-items: center; gap: 0.75rem; }
  .dep-icon { font-size: 1.2rem; color: #004ac6; }
  .dep-date { font-weight: 700; font-size: 0.88rem; color: #191c1e; }
  .dep-airline { font-size: 0.72rem; color: #434655; margin-top: 0.1rem; }

  .chevron { font-size: 1.2rem; color: #434655; transition: transform 0.2s; }
  .chevron.rotated { transform: rotate(180deg); }

  .dep-body { padding: 0 1rem 1rem; border-top: 1px solid rgb(195 198 215 / 0.35); }

  .pricing-table { width: 100%; border-collapse: collapse; margin-top: 0.75rem; }
  .pricing-table th, .pricing-table td { padding: 0.55rem 0.7rem; text-align: left; font-size: 0.82rem; border-bottom: 1px solid rgb(195 198 215 / 0.35); }
  .pricing-table th { font-size: 0.62rem; text-transform: uppercase; letter-spacing: 0.08em; color: #434655; background: #f2f4f6; }
  .pricing-table tbody tr:last-child td { border-bottom: 0; }

  .room-badge { display: inline-flex; padding: 0.15rem 0.45rem; border-radius: 0.2rem; font-size: 0.72rem; font-weight: 700; }
  .room-badge--double { background: #dbeafe; color: #1e3a8a; }
  .room-badge--triple { background: #ede9fe; color: #4c1d95; }
  .room-badge--quad { background: #d1fae5; color: #065f46; }

  .price-input {
    border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem;
    padding: 0.4rem 0.6rem; font-size: 0.82rem; width: 10rem;
    font-family: inherit; outline: none;
  }
  .price-input:focus { border-color: #2563eb; }

  .price-preview { margin-left: 0.5rem; font-size: 0.72rem; color: #434655; font-style: italic; }

  .dep-footer { display: flex; align-items: center; justify-content: flex-end; gap: 0.75rem; margin-top: 0.75rem; }

  .saved-msg { display: inline-flex; align-items: center; gap: 0.3rem; font-size: 0.75rem; color: #065f46; }
  .saved-msg .material-symbols-outlined { font-size: 1rem; }

  .primary-btn {
    border-radius: 0.25rem; padding: 0.5rem 0.85rem; font-size: 0.8rem;
    font-weight: 600; cursor: pointer; border: 1px solid #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff;
    display: inline-flex; align-items: center; gap: 0.35rem; font-family: inherit;
  }
  .primary-btn:disabled { opacity: 0.6; cursor: not-allowed; }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; }
</style>
