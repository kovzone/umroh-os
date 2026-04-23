<script lang="ts">
  import type { PageData, Hotel, Airline, Muthawwif, Addon } from './+page.server';

  let { data }: { data: PageData } = $props();

  // ---- local state synced via $effect ----
  let hotels = $state<Hotel[]>([]);
  let airlines = $state<Airline[]>([]);
  let muthawwifs = $state<Muthawwif[]>([]);
  let addons = $state<Addon[]>([]);

  $effect(() => {
    hotels = data.hotels ?? [];
    airlines = data.airlines ?? [];
    muthawwifs = data.muthawwifs ?? [];
    addons = data.addons ?? [];
  });

  // ---- active tab ----
  type Tab = 'hotel' | 'maskapai' | 'muthawwif' | 'addon';
  let activeTab = $state<Tab>('hotel');

  // ---- modal state ----
  type ModalMode = 'create' | 'edit';
  let modalOpen = $state(false);
  let modalMode = $state<ModalMode>('create');
  let modalSaving = $state(false);
  let modalError = $state('');

  // hotel form
  let hotelForm = $state({ id: '', name: '', city: '', stars: 3, distance_m: 500 });
  // airline form
  let airlineForm = $state({ code: '', name: '', type: 'airline' as 'airline' | 'rail' | 'bus' });
  // muthawwif form
  let muthawwifForm = $state({ id: '', name: '' });
  // addon form
  let addonForm = $state({ id: '', name: '', price_idr: 0 });

  function openCreate() {
    modalMode = 'create';
    modalError = '';
    if (activeTab === 'hotel') hotelForm = { id: '', name: '', city: '', stars: 3, distance_m: 500 };
    if (activeTab === 'maskapai') airlineForm = { code: '', name: '', type: 'airline' };
    if (activeTab === 'muthawwif') muthawwifForm = { id: '', name: '' };
    if (activeTab === 'addon') addonForm = { id: '', name: '', price_idr: 0 };
    modalOpen = true;
  }

  function openEdit(item: Hotel | Airline | Muthawwif | Addon) {
    modalMode = 'edit';
    modalError = '';
    if (activeTab === 'hotel') {
      const h = item as Hotel;
      hotelForm = { id: h.id, name: h.name, city: h.city, stars: h.stars, distance_m: h.distance_m };
    }
    if (activeTab === 'maskapai') {
      const a = item as Airline;
      airlineForm = { code: a.code, name: a.name, type: a.type };
    }
    if (activeTab === 'muthawwif') {
      const m = item as Muthawwif;
      muthawwifForm = { id: m.id, name: m.name };
    }
    if (activeTab === 'addon') {
      const d = item as Addon;
      addonForm = { id: d.id, name: d.name, price_idr: d.price_idr };
    }
    modalOpen = true;
  }

  function closeModal() {
    modalOpen = false;
    modalError = '';
  }

  async function saveModal() {
    modalSaving = true;
    modalError = '';
    try {
      // Mock save — in production POST/PUT to API
      await new Promise((r) => setTimeout(r, 400));
      if (activeTab === 'hotel') {
        const h: Hotel = { ...hotelForm, id: hotelForm.id || `h${Date.now()}` };
        if (modalMode === 'create') {
          hotels = [...hotels, h];
        } else {
          hotels = hotels.map((x) => (x.id === h.id ? h : x));
        }
      }
      if (activeTab === 'maskapai') {
        const a: Airline = { ...airlineForm };
        if (modalMode === 'create') {
          airlines = [...airlines, a];
        } else {
          airlines = airlines.map((x) => (x.code === a.code ? a : x));
        }
      }
      if (activeTab === 'muthawwif') {
        const m: Muthawwif = { ...muthawwifForm, photo_url: null, id: muthawwifForm.id || `m${Date.now()}` };
        if (modalMode === 'create') {
          muthawwifs = [...muthawwifs, m];
        } else {
          muthawwifs = muthawwifs.map((x) => (x.id === m.id ? m : x));
        }
      }
      if (activeTab === 'addon') {
        const d: Addon = { ...addonForm, id: addonForm.id || `a${Date.now()}` };
        if (modalMode === 'create') {
          addons = [...addons, d];
        } else {
          addons = addons.map((x) => (x.id === d.id ? d : x));
        }
      }
      closeModal();
    } catch {
      modalError = 'Gagal menyimpan data.';
    } finally {
      modalSaving = false;
    }
  }

  function deleteItem(id: string) {
    if (!confirm('Hapus data ini?')) return;
    if (activeTab === 'hotel') hotels = hotels.filter((x) => x.id !== id);
    if (activeTab === 'maskapai') airlines = airlines.filter((x) => x.code !== id);
    if (activeTab === 'muthawwif') muthawwifs = muthawwifs.filter((x) => x.id !== id);
    if (activeTab === 'addon') addons = addons.filter((x) => x.id !== id);
  }

  function formatIDR(n: number): string {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  const TAB_LABELS: Record<Tab, string> = {
    hotel: 'Hotel',
    maskapai: 'Maskapai',
    muthawwif: 'Muthawwif',
    addon: 'Addon'
  };

  const AIRLINE_TYPE_LABELS: Record<string, string> = {
    airline: 'Maskapai',
    rail: 'Kereta',
    bus: 'Bus'
  };

  const modalTitle = $derived(
    modalMode === 'create'
      ? `Tambah ${TAB_LABELS[activeTab]}`
      : `Edit ${TAB_LABELS[activeTab]}`
  );
</script>

<main class="page-shell">
  <header class="topbar">
    <div class="breadcrumb">
      <span class="material-symbols-outlined bc-icon">tune</span>
      <span class="bc-text">Master Data</span>
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
        <h2>Master Data Hub</h2>
        <p>Kelola data referensi: hotel, maskapai, muthawwif, dan addon paket umroh</p>
      </div>
      <button type="button" class="primary-btn" onclick={openCreate}>
        <span class="material-symbols-outlined">add</span>
        Tambah {TAB_LABELS[activeTab]}
      </button>
    </div>

    {#if data.error}
      <div class="error-banner">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- Tabs -->
    <div class="tabs">
      {#each (['hotel', 'maskapai', 'muthawwif', 'addon'] as Tab[]) as tab}
        <button
          type="button"
          class="tab-btn"
          class:active={activeTab === tab}
          onclick={() => { activeTab = tab; }}
        >
          {TAB_LABELS[tab]}
        </button>
      {/each}
    </div>

    <!-- Hotel Tab -->
    {#if activeTab === 'hotel'}
      <div class="panel">
        <div class="table-wrap">
          {#if hotels.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">hotel</span>
              <p>Belum ada data hotel. Tambahkan hotel baru.</p>
            </div>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Nama Hotel</th>
                  <th>Kota</th>
                  <th>Bintang</th>
                  <th>Jarak Masjid</th>
                  <th class="align-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {#each hotels as hotel (hotel.id)}
                  <tr>
                    <td><code class="id-code">{hotel.id}</code></td>
                    <td><strong>{hotel.name}</strong></td>
                    <td>{hotel.city}</td>
                    <td>
                      <span class="stars">
                        {#each { length: hotel.stars } as _}★{/each}
                      </span>
                    </td>
                    <td>{hotel.distance_m} m</td>
                    <td class="actions-cell">
                      <button class="action-btn" onclick={() => openEdit(hotel)}>
                        <span class="material-symbols-outlined">edit</span>
                      </button>
                      <button class="action-btn danger" onclick={() => deleteItem(hotel.id)}>
                        <span class="material-symbols-outlined">delete</span>
                      </button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
        <div class="table-footer">{hotels.length} hotel terdaftar</div>
      </div>
    {/if}

    <!-- Maskapai Tab -->
    {#if activeTab === 'maskapai'}
      <div class="panel">
        <div class="table-wrap">
          {#if airlines.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">flight</span>
              <p>Belum ada data maskapai. Tambahkan maskapai baru.</p>
            </div>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Kode</th>
                  <th>Nama Maskapai</th>
                  <th>Tipe</th>
                  <th class="align-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {#each airlines as airline (airline.code)}
                  <tr>
                    <td><code class="id-code">{airline.code}</code></td>
                    <td><strong>{airline.name}</strong></td>
                    <td>
                      <span class="type-badge type-badge--{airline.type}">
                        {AIRLINE_TYPE_LABELS[airline.type]}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button class="action-btn" onclick={() => openEdit(airline)}>
                        <span class="material-symbols-outlined">edit</span>
                      </button>
                      <button class="action-btn danger" onclick={() => deleteItem(airline.code)}>
                        <span class="material-symbols-outlined">delete</span>
                      </button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
        <div class="table-footer">{airlines.length} maskapai terdaftar</div>
      </div>
    {/if}

    <!-- Muthawwif Tab -->
    {#if activeTab === 'muthawwif'}
      <div class="panel">
        <div class="table-wrap">
          {#if muthawwifs.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">person</span>
              <p>Belum ada data muthawwif. Tambahkan muthawwif baru.</p>
            </div>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Muthawwif</th>
                  <th>Foto</th>
                  <th class="align-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {#each muthawwifs as m (m.id)}
                  <tr>
                    <td><strong>{m.name}</strong></td>
                    <td>
                      {#if m.photo_url}
                        <img src={m.photo_url} alt={m.name} class="avatar-img" />
                      {:else}
                        <div class="avatar-initial">{m.name.charAt(0)}</div>
                      {/if}
                    </td>
                    <td class="actions-cell">
                      <button class="action-btn" onclick={() => openEdit(m)}>
                        <span class="material-symbols-outlined">edit</span>
                      </button>
                      <button class="action-btn danger" onclick={() => deleteItem(m.id)}>
                        <span class="material-symbols-outlined">delete</span>
                      </button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
        <div class="table-footer">{muthawwifs.length} muthawwif terdaftar</div>
      </div>
    {/if}

    <!-- Addon Tab -->
    {#if activeTab === 'addon'}
      <div class="panel">
        <div class="table-wrap">
          {#if addons.length === 0}
            <div class="empty-state">
              <span class="material-symbols-outlined">add_box</span>
              <p>Belum ada data addon. Tambahkan addon baru.</p>
            </div>
          {:else}
            <table>
              <thead>
                <tr>
                  <th>Nama Addon</th>
                  <th>Harga (IDR)</th>
                  <th class="align-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {#each addons as addon (addon.id)}
                  <tr>
                    <td><strong>{addon.name}</strong></td>
                    <td class="price-cell">{formatIDR(addon.price_idr)}</td>
                    <td class="actions-cell">
                      <button class="action-btn" onclick={() => openEdit(addon)}>
                        <span class="material-symbols-outlined">edit</span>
                      </button>
                      <button class="action-btn danger" onclick={() => deleteItem(addon.id)}>
                        <span class="material-symbols-outlined">delete</span>
                      </button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/if}
        </div>
        <div class="table-footer">{addons.length} addon terdaftar</div>
      </div>
    {/if}
  </section>
</main>

<!-- Modal -->
{#if modalOpen}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div class="modal-backdrop" onclick={closeModal}></div>
  <div class="modal" role="dialog" aria-labelledby="modal-title" aria-modal="true">
    <div class="modal-header">
      <h3 id="modal-title">{modalTitle}</h3>
      <button type="button" class="modal-close" onclick={closeModal} aria-label="Tutup">
        <span class="material-symbols-outlined">close</span>
      </button>
    </div>

    <div class="modal-body">
      {#if activeTab === 'hotel'}
        <div class="field-row">
          <label class="field-label" for="f-hotel-name">Nama Hotel</label>
          <input id="f-hotel-name" type="text" class="field-input" bind:value={hotelForm.name} placeholder="Nama hotel" />
        </div>
        <div class="field-row">
          <label class="field-label" for="f-hotel-city">Kota</label>
          <select id="f-hotel-city" class="field-input" bind:value={hotelForm.city}>
            <option value="Makkah">Makkah</option>
            <option value="Madinah">Madinah</option>
          </select>
        </div>
        <div class="field-row">
          <label class="field-label" for="f-hotel-stars">Bintang (1-5)</label>
          <input id="f-hotel-stars" type="number" min="1" max="5" class="field-input" bind:value={hotelForm.stars} />
        </div>
        <div class="field-row">
          <label class="field-label" for="f-hotel-dist">Jarak ke Masjid (m)</label>
          <input id="f-hotel-dist" type="number" min="0" class="field-input" bind:value={hotelForm.distance_m} />
        </div>
      {/if}

      {#if activeTab === 'maskapai'}
        <div class="field-row">
          <label class="field-label" for="f-airline-code">Kode IATA</label>
          <input id="f-airline-code" type="text" class="field-input" bind:value={airlineForm.code} placeholder="GA" maxlength="3" />
        </div>
        <div class="field-row">
          <label class="field-label" for="f-airline-name">Nama Maskapai</label>
          <input id="f-airline-name" type="text" class="field-input" bind:value={airlineForm.name} placeholder="Garuda Indonesia" />
        </div>
        <div class="field-row">
          <label class="field-label" for="f-airline-type">Tipe</label>
          <select id="f-airline-type" class="field-input" bind:value={airlineForm.type}>
            <option value="airline">Maskapai</option>
            <option value="rail">Kereta</option>
            <option value="bus">Bus</option>
          </select>
        </div>
      {/if}

      {#if activeTab === 'muthawwif'}
        <div class="field-row">
          <label class="field-label" for="f-muth-name">Nama Muthawwif</label>
          <input id="f-muth-name" type="text" class="field-input" bind:value={muthawwifForm.name} placeholder="Ust. Nama Lengkap" />
        </div>
      {/if}

      {#if activeTab === 'addon'}
        <div class="field-row">
          <label class="field-label" for="f-addon-name">Nama Addon</label>
          <input id="f-addon-name" type="text" class="field-input" bind:value={addonForm.name} placeholder="Nama layanan tambahan" />
        </div>
        <div class="field-row">
          <label class="field-label" for="f-addon-price">Harga (IDR)</label>
          <input id="f-addon-price" type="number" min="0" step="5000" class="field-input" bind:value={addonForm.price_idr} />
        </div>
      {/if}

      {#if modalError}
        <p class="modal-error">{modalError}</p>
      {/if}
    </div>

    <div class="modal-footer">
      <button type="button" class="ghost-btn" onclick={closeModal} disabled={modalSaving}>Batal</button>
      <button type="button" class="primary-btn" onclick={saveModal} disabled={modalSaving}>
        {#if modalSaving}
          <span class="material-symbols-outlined spin">progress_activity</span>
          Menyimpan...
        {:else}
          <span class="material-symbols-outlined">save</span>
          Simpan
        {/if}
      </button>
    </div>
  </div>
{/if}

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }

  .topbar {
    position: sticky; top: 0; z-index: 30;
    height: 4rem; background: rgb(255 255 255 / 0.9);
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem; display: flex; align-items: center;
    justify-content: space-between; gap: 1rem; backdrop-filter: blur(8px);
  }

  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.88rem; font-weight: 600; color: #191c1e; }
  .bc-icon { font-size: 1.05rem; color: #004ac6; }

  .top-actions { display: flex; align-items: center; gap: 0.5rem; }

  .icon-btn {
    border: 0; background: transparent; color: #434655;
    width: 2rem; height: 2rem; border-radius: 0.25rem;
    cursor: pointer; position: relative; display: grid; place-items: center;
  }
  .icon-btn:hover { background: #eceef0; }
  .dot { position: absolute; width: 0.46rem; height: 0.46rem; border-radius: 999px; background: #ba1a1a; right: 0.4rem; top: 0.35rem; border: 2px solid #fff; }

  .avatar {
    border: 1px solid rgb(195 198 215 / 0.55); background: #b4c5ff;
    color: #00174b; width: 2rem; height: 2rem; border-radius: 0.25rem;
    font-weight: 700; font-size: 0.65rem; cursor: pointer;
  }

  .canvas { padding: 1.5rem; max-width: 96rem; }

  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; }
  .page-head p { margin: 0.3rem 0 0; font-size: 0.82rem; color: #434655; }

  .error-banner {
    display: flex; align-items: center; gap: 0.5rem;
    background: #ffdad6; color: #93000a; border-radius: 0.25rem;
    padding: 0.65rem 0.85rem; font-size: 0.82rem; margin-bottom: 1rem;
  }

  .tabs {
    display: flex; gap: 0; margin-bottom: 1rem;
    border-bottom: 2px solid rgb(195 198 215 / 0.45);
  }

  .tab-btn {
    padding: 0.55rem 1.1rem; border: 0; background: transparent;
    font-size: 0.82rem; font-weight: 600; color: #434655;
    cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -2px;
  }
  .tab-btn:hover { color: #191c1e; }
  .tab-btn.active { color: #004ac6; border-bottom-color: #004ac6; }

  .panel {
    background: #fff; border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem; overflow: hidden;
  }

  .table-wrap { overflow-x: auto; }

  table { width: 100%; border-collapse: collapse; }

  th, td {
    padding: 0.62rem 0.85rem; text-align: left;
    font-size: 0.78rem; border-bottom: 1px solid rgb(195 198 215 / 0.45);
    white-space: nowrap;
  }

  th {
    text-transform: uppercase; font-size: 0.62rem;
    letter-spacing: 0.08em; color: #434655; background: #f2f4f6;
  }

  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }
  .align-right { text-align: right; }

  .id-code {
    font-family: 'IBM Plex Mono', monospace; font-size: 0.72rem;
    background: #f2f4f6; padding: 0.1rem 0.3rem; border-radius: 0.15rem;
  }

  .stars { color: #d97706; letter-spacing: -1px; }

  .price-cell { font-weight: 600; color: #004ac6; }

  .type-badge {
    display: inline-flex; padding: 0.12rem 0.4rem; border-radius: 0.2rem;
    font-size: 0.65rem; font-weight: 700; background: #e6e8ea; color: #434655;
  }
  .type-badge--airline { background: #dbeafe; color: #1e3a8a; }
  .type-badge--rail { background: #d1fae5; color: #065f46; }
  .type-badge--bus { background: #fef9c3; color: #7d5f00; }

  .avatar-img { width: 2rem; height: 2rem; border-radius: 999px; object-fit: cover; }
  .avatar-initial {
    width: 2rem; height: 2rem; border-radius: 999px; background: #b4c5ff;
    color: #00174b; font-size: 0.72rem; font-weight: 700;
    display: grid; place-items: center;
  }

  .actions-cell { text-align: right; }
  .action-btn {
    border: 1px solid rgb(195 198 215 / 0.55); background: #fff;
    border-radius: 0.2rem; cursor: pointer; padding: 0.25rem 0.35rem;
    display: inline-flex; align-items: center; color: #434655;
    margin-left: 0.25rem;
  }
  .action-btn:hover { background: #f2f4f6; }
  .action-btn.danger:hover { background: #ffdad6; color: #ba1a1a; border-color: #fca5a5; }
  .action-btn .material-symbols-outlined { font-size: 0.95rem; }

  .table-footer {
    padding: 0.55rem 0.85rem; font-size: 0.68rem; color: #434655;
    border-top: 1px solid rgb(195 198 215 / 0.35); background: #f7f9fb;
  }

  .empty-state {
    display: flex; flex-direction: column; align-items: center;
    gap: 0.6rem; padding: 3rem 1rem; color: #b0b3c1;
  }
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

  .modal-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 0.85rem 1rem; border-bottom: 1px solid rgb(195 198 215 / 0.45); background: #f2f4f6;
  }
  .modal-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .modal-close { border: 0; background: transparent; cursor: pointer; color: #434655; display: grid; place-items: center; border-radius: 0.2rem; padding: 0.2rem; }
  .modal-close:hover { background: #e6e8ea; }

  .modal-body { padding: 1rem; display: flex; flex-direction: column; gap: 0.75rem; }

  .field-row { display: flex; flex-direction: column; gap: 0.25rem; }
  .field-label { font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field-input {
    border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem;
    padding: 0.5rem 0.65rem; font-size: 0.82rem; color: #191c1e;
    background: #fff; font-family: inherit; outline: none;
  }

  .modal-error { margin: 0; font-size: 0.75rem; color: #ba1a1a; background: #ffdad6; padding: 0.45rem 0.65rem; border-radius: 0.2rem; }

  .modal-footer { padding: 0.75rem 1rem; border-top: 1px solid rgb(195 198 215 / 0.45); display: flex; justify-content: flex-end; gap: 0.55rem; }

  .ghost-btn, .primary-btn {
    border-radius: 0.25rem; padding: 0.5rem 0.85rem; font-size: 0.8rem;
    font-weight: 600; cursor: pointer; border: 1px solid rgb(195 198 215 / 0.55);
    display: inline-flex; align-items: center; gap: 0.35rem; font-family: inherit;
  }
  .ghost-btn { background: #fff; color: #191c1e; }
  .ghost-btn:disabled, .primary-btn:disabled { opacity: 0.6; cursor: not-allowed; }
  .primary-btn {
    border-color: #2563eb; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff;
  }
  .primary-btn .material-symbols-outlined { font-size: 1rem; }

  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1rem; }
</style>
