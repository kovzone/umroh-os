<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

  // Card expand state
  let expanded = $state<string | null>(null);
  function toggle(id: string) { expanded = expanded === id ? null : id; }

  // ---- Luggage Counter ----
  let lg_departureId = $state('');
  let lg_barcode = $state('');
  let lg_loading = $state(false);
  let lg_error = $state('');
  let lg_result = $state<any>(null);

  async function scanLuggage(e: Event) {
    e.preventDefault();
    lg_loading = true; lg_error = ''; lg_result = null;
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/luggage-scan`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ departure_id: lg_departureId, barcode: lg_barcode }),
      });
      if (!res.ok) throw new Error(`Gagal scan koper (${res.status})`);
      lg_result = await res.json();
      lg_barcode = '';
    } catch (err) {
      lg_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    lg_loading = false;
  }

  // ---- Broadcast ----
  let bc_departureId = $state('');
  let bc_type = $state('departure');
  let bc_message = $state('');
  let bc_loading = $state(false);
  let bc_error = $state('');
  let bc_success = $state('');

  async function sendBroadcast(e: Event) {
    e.preventDefault();
    bc_loading = true; bc_error = ''; bc_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/broadcast-schedule`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ departure_id: bc_departureId, type: bc_type, message: bc_message }),
      });
      if (!res.ok) throw new Error(`Gagal kirim broadcast (${res.status})`);
      bc_success = 'Broadcast berhasil dikirim.';
      bc_message = '';
    } catch (err) {
      bc_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    bc_loading = false;
  }

  // ---- Tasreh & Raudhah ----
  let ts_pilgrimId = $state('');
  let ts_type = $state('tasreh');
  let ts_loading = $state(false);
  let ts_error = $state('');
  let ts_success = $state('');

  async function issueTasreh(e: Event) {
    e.preventDefault();
    ts_loading = true; ts_error = ''; ts_success = '';
    try {
      const endpoint = ts_type === 'tasreh' ? 'tasreh' : 'raudhah-entry';
      const res = await fetch(`${GATEWAY}/v1/ops/${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pilgrim_id: ts_pilgrimId }),
      });
      if (!res.ok) throw new Error(`Gagal (${res.status})`);
      ts_success = ts_type === 'tasreh' ? 'Tasreh berhasil diterbitkan.' : 'Masuk Raudhah berhasil dicatat.';
      ts_pilgrimId = '';
    } catch (err) {
      ts_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ts_loading = false;
  }

  // ---- Audio Devices ----
  let au_deviceId = $state('');
  let au_deviceName = $state('');
  let au_groupId = $state('');
  let au_loading = $state(false);
  let au_error = $state('');
  let au_success = $state('');

  async function registerAudio(e: Event) {
    e.preventDefault();
    au_loading = true; au_error = ''; au_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/audio-devices`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ device_id: au_deviceId, device_name: au_deviceName, group_id: au_groupId }),
      });
      if (!res.ok) throw new Error(`Gagal daftarkan perangkat (${res.status})`);
      au_success = 'Perangkat audio berhasil didaftarkan.';
      au_deviceId = ''; au_deviceName = ''; au_groupId = '';
    } catch (err) {
      au_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    au_loading = false;
  }

  // ---- Zamzam Distribution ----
  let zz_departureId = $state('');
  let zz_pilgrimId = $state('');
  let zz_liters = $state('');
  let zz_loading = $state(false);
  let zz_error = $state('');
  let zz_success = $state('');

  async function recordZamzam(e: Event) {
    e.preventDefault();
    zz_loading = true; zz_error = ''; zz_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/zamzam-distribution`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          departure_id: zz_departureId,
          pilgrim_id: zz_pilgrimId,
          liters: parseFloat(zz_liters),
        }),
      });
      if (!res.ok) throw new Error(`Gagal catat distribusi (${res.status})`);
      zz_success = 'Distribusi zamzam berhasil dicatat.';
      zz_pilgrimId = ''; zz_liters = '';
    } catch (err) {
      zz_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    zz_loading = false;
  }

  // ---- Room Check-in ----
  let ci_pilgrimId = $state('');
  let ci_roomNumber = $state('');
  let ci_hotelName = $state('');
  let ci_loading = $state(false);
  let ci_error = $state('');
  let ci_success = $state('');

  async function recordCheckin(e: Event) {
    e.preventDefault();
    ci_loading = true; ci_error = ''; ci_success = '';
    try {
      const res = await fetch(`${GATEWAY}/v1/ops/room-checkin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          pilgrim_id: ci_pilgrimId,
          room_number: ci_roomNumber,
          hotel_name: ci_hotelName,
        }),
      });
      if (!res.ok) throw new Error(`Gagal catat check-in (${res.status})`);
      ci_success = 'Check-in kamar berhasil dicatat.';
      ci_pilgrimId = ''; ci_roomNumber = '';
    } catch (err) {
      ci_error = err instanceof Error ? err.message : 'Terjadi kesalahan';
    }
    ci_loading = false;
  }

  const CARDS = [
    { id: 'luggage', icon: 'luggage', label: 'Luggage Counter', desc: 'Scan koper & hitung per keberangkatan' },
    { id: 'broadcast', icon: 'campaign', label: 'Broadcast Jadwal', desc: 'Kirim broadcast jadwal berangkat/tiba' },
    { id: 'tasreh', icon: 'mosque', label: 'Tasreh & Raudhah', desc: 'Terbitkan tasreh & catat masuk Raudhah' },
    { id: 'audio', icon: 'headphones', label: 'Perangkat Audio', desc: 'Daftarkan dan kelola audio device' },
    { id: 'zamzam', icon: 'water_drop', label: 'Distribusi Zamzam', desc: 'Catat distribusi air zamzam' },
    { id: 'checkin', icon: 'hotel', label: 'Check-in Kamar', desc: 'Rekam check-in kamar jamaah' },
  ];
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/ops" class="bc-link">Ops</a>
      <span class="bc-sep">/</span>
      <span>Operasional Lapangan</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Operasional Lapangan</h2>
      <p>BL-OPS-036..042 — Hub operasional lapangan: koper, broadcast, tasreh, audio, zamzam, check-in</p>
    </div>

    <div class="hub-grid">
      {#each CARDS as card}
        <div class="hub-card" class:open={expanded === card.id}>
          <button class="hub-card-header" onclick={() => toggle(card.id)}>
            <span class="material-symbols-outlined hub-icon">{card.icon}</span>
            <div class="hub-card-text">
              <div class="hub-card-label">{card.label}</div>
              <div class="hub-card-desc">{card.desc}</div>
            </div>
            <span class="material-symbols-outlined chevron">
              {expanded === card.id ? 'expand_less' : 'expand_more'}
            </span>
          </button>

          {#if expanded === card.id}
            <div class="hub-form-area">

              {#if card.id === 'luggage'}
                <form class="mini-form" onsubmit={scanLuggage}>
                  <div class="field">
                    <label>ID Keberangkatan</label>
                    <input type="text" placeholder="dep-001" bind:value={lg_departureId} required />
                  </div>
                  <div class="field">
                    <label>Barcode Koper</label>
                    <input type="text" placeholder="LGG-00001" bind:value={lg_barcode} required />
                  </div>
                  <button type="submit" class="btn-primary" disabled={lg_loading}>
                    {#if lg_loading}<span class="spinner"></span>{/if}Scan
                  </button>
                  {#if lg_error}<div class="alert-err">{lg_error}</div>{/if}
                  {#if lg_result}
                    <div class="alert-ok">Koper <strong>{lg_barcode || lg_result.barcode}</strong> — Total: {lg_result.total_count ?? lg_result.count ?? '?'} koper</div>
                  {/if}
                </form>
              {/if}

              {#if card.id === 'broadcast'}
                <form class="mini-form" onsubmit={sendBroadcast}>
                  <div class="field">
                    <label>ID Keberangkatan</label>
                    <input type="text" placeholder="dep-001" bind:value={bc_departureId} required />
                  </div>
                  <div class="field">
                    <label>Jenis</label>
                    <select bind:value={bc_type}>
                      <option value="departure">Keberangkatan</option>
                      <option value="arrival">Kepulangan</option>
                      <option value="reminder">Pengingat</option>
                    </select>
                  </div>
                  <div class="field">
                    <label>Pesan</label>
                    <textarea rows="3" placeholder="Pesan broadcast..." bind:value={bc_message} required></textarea>
                  </div>
                  <button type="submit" class="btn-primary" disabled={bc_loading}>
                    {#if bc_loading}<span class="spinner"></span>{/if}Kirim Broadcast
                  </button>
                  {#if bc_error}<div class="alert-err">{bc_error}</div>{/if}
                  {#if bc_success}<div class="alert-ok">{bc_success}</div>{/if}
                </form>
              {/if}

              {#if card.id === 'tasreh'}
                <form class="mini-form" onsubmit={issueTasreh}>
                  <div class="field">
                    <label>ID Jamaah</label>
                    <input type="text" placeholder="pilg-001" bind:value={ts_pilgrimId} required />
                  </div>
                  <div class="field">
                    <label>Jenis</label>
                    <select bind:value={ts_type}>
                      <option value="tasreh">Terbitkan Tasreh</option>
                      <option value="raudhah">Catat Masuk Raudhah</option>
                    </select>
                  </div>
                  <button type="submit" class="btn-primary" disabled={ts_loading}>
                    {#if ts_loading}<span class="spinner"></span>{/if}Simpan
                  </button>
                  {#if ts_error}<div class="alert-err">{ts_error}</div>{/if}
                  {#if ts_success}<div class="alert-ok">{ts_success}</div>{/if}
                </form>
              {/if}

              {#if card.id === 'audio'}
                <form class="mini-form" onsubmit={registerAudio}>
                  <div class="field">
                    <label>ID Perangkat</label>
                    <input type="text" placeholder="AUD-001" bind:value={au_deviceId} required />
                  </div>
                  <div class="field">
                    <label>Nama Perangkat</label>
                    <input type="text" placeholder="Clip-on Pembimbing 1" bind:value={au_deviceName} required />
                  </div>
                  <div class="field">
                    <label>Grup/Rombongan</label>
                    <input type="text" placeholder="GRP-A" bind:value={au_groupId} />
                  </div>
                  <button type="submit" class="btn-primary" disabled={au_loading}>
                    {#if au_loading}<span class="spinner"></span>{/if}Daftarkan
                  </button>
                  {#if au_error}<div class="alert-err">{au_error}</div>{/if}
                  {#if au_success}<div class="alert-ok">{au_success}</div>{/if}
                </form>
              {/if}

              {#if card.id === 'zamzam'}
                <form class="mini-form" onsubmit={recordZamzam}>
                  <div class="field">
                    <label>ID Keberangkatan</label>
                    <input type="text" placeholder="dep-001" bind:value={zz_departureId} required />
                  </div>
                  <div class="field">
                    <label>ID Jamaah</label>
                    <input type="text" placeholder="pilg-001" bind:value={zz_pilgrimId} required />
                  </div>
                  <div class="field">
                    <label>Volume (Liter)</label>
                    <input type="number" min="0" step="0.5" placeholder="5" bind:value={zz_liters} required />
                  </div>
                  <button type="submit" class="btn-primary" disabled={zz_loading}>
                    {#if zz_loading}<span class="spinner"></span>{/if}Catat
                  </button>
                  {#if zz_error}<div class="alert-err">{zz_error}</div>{/if}
                  {#if zz_success}<div class="alert-ok">{zz_success}</div>{/if}
                </form>
              {/if}

              {#if card.id === 'checkin'}
                <form class="mini-form" onsubmit={recordCheckin}>
                  <div class="field">
                    <label>ID Jamaah</label>
                    <input type="text" placeholder="pilg-001" bind:value={ci_pilgrimId} required />
                  </div>
                  <div class="field">
                    <label>Nama Hotel</label>
                    <input type="text" placeholder="Hotel Grand Makkah" bind:value={ci_hotelName} required />
                  </div>
                  <div class="field">
                    <label>Nomor Kamar</label>
                    <input type="text" placeholder="401" bind:value={ci_roomNumber} required />
                  </div>
                  <button type="submit" class="btn-primary" disabled={ci_loading}>
                    {#if ci_loading}<span class="spinner"></span>{/if}Catat Check-in
                  </button>
                  {#if ci_error}<div class="alert-err">{ci_error}</div>{/if}
                  {#if ci_success}<div class="alert-ok">{ci_success}</div>{/if}
                </form>
              {/if}

            </div>
          {/if}
        </div>
      {/each}
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
  .page-head { margin-bottom: 1.5rem; }
  .page-head h2 { margin: 0; font-size: 1.4rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.78rem; color: #737686; }
  .hub-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1rem; }
  .hub-card { background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; overflow: hidden; transition: box-shadow 0.15s; }
  .hub-card:hover { box-shadow: 0 2px 8px rgb(0 0 0 / 0.07); }
  .hub-card.open { border-color: #93c5fd; }
  .hub-card-header { display: flex; align-items: center; gap: 0.85rem; padding: 1rem 1.1rem; background: transparent; border: 0; width: 100%; text-align: left; cursor: pointer; font-family: inherit; }
  .hub-card-header:hover { background: #f7f9fb; }
  .hub-icon { font-size: 1.5rem; color: #004ac6; flex-shrink: 0; }
  .hub-card-text { flex: 1; }
  .hub-card-label { font-size: 0.88rem; font-weight: 700; color: #191c1e; }
  .hub-card-desc { font-size: 0.72rem; color: #737686; margin-top: 0.1rem; }
  .chevron { font-size: 1.1rem; color: #737686; flex-shrink: 0; }
  .hub-form-area { padding: 1rem 1.1rem; border-top: 1px solid rgb(195 198 215 / 0.35); background: #fafbfc; }
  .mini-form { display: flex; flex-direction: column; gap: 0.65rem; }
  .field { display: flex; flex-direction: column; gap: 0.25rem; }
  .field label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select, .field textarea { border: 1px solid rgb(195 198 215 / 0.55); border-radius: 0.25rem; padding: 0.42rem 0.6rem; font-size: 0.82rem; color: #191c1e; background: #fff; font-family: inherit; resize: vertical; }
  .btn-primary { display: inline-flex; align-items: center; gap: 0.35rem; background: linear-gradient(90deg,#004ac6,#2563eb); color: #fff; border: 0; border-radius: 0.25rem; padding: 0.48rem 1rem; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; align-self: flex-start; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .spinner { width: 0.8rem; height: 0.8rem; border: 2px solid rgb(255 255 255 / 0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.5rem 0.75rem; font-size: 0.78rem; }
  .alert-ok { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; border-radius: 0.25rem; padding: 0.5rem 0.75rem; font-size: 0.78rem; }
</style>
