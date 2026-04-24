<script lang="ts">
  const GATEWAY = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';
  function formatIDR(n: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  }

  let ledgerType = $state('Keduanya');
  let thresholdDays = $state('30');
  let customDays = $state('');
  let alerts = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  async function load(e?: SubmitEvent) {
    if (e) e.preventDefault();
    loading = true; error = '';
    try {
      const days = thresholdDays === 'custom' ? customDays : thresholdDays;
      const params = new URLSearchParams();
      if (ledgerType !== 'Keduanya') params.set('ledger_type', ledgerType);
      if (days) params.set('threshold_days', days);
      const res = await fetch(`${GATEWAY}/v1/finance/aging-alerts?${params.toString()}`);
      const body = res.ok ? await res.json() : null;
      alerts = Array.isArray(body?.data) ? body.data : (Array.isArray(body) ? body : []);
    } catch { error = 'Gagal memuat data.'; alerts = []; }
    loading = false;
  }

  $effect(() => { load(); });

  function severity(days: number): 'red' | 'orange' | 'yellow' {
    if (days > 90) return 'red';
    if (days > 60) return 'orange';
    return 'yellow';
  }

  function severityLabel(days: number): string {
    if (days > 90) return '>90 hari';
    if (days > 60) return '60-90 hari';
    return '30-60 hari';
  }

  const grouped = $derived(() => {
    const red = alerts.filter((a: any) => a.overdue_days > 90);
    const orange = alerts.filter((a: any) => a.overdue_days > 60 && a.overdue_days <= 90);
    const yellow = alerts.filter((a: any) => a.overdue_days <= 60);
    return { red, orange, yellow };
  });

  const totalOverdue = $derived(alerts.reduce((sum: number, a: any) => sum + (a.amount ?? 0), 0));
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console/finance" class="bc-link">Keuangan</a>
      <span class="bc-sep">/</span>
      <span>Peringatan Jatuh Tempo</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Peringatan Piutang/Utang Jatuh Tempo</h2>
        <p>Monitor AR & AP yang mendekati atau melewati jatuh tempo</p>
      </div>
      {#if alerts.length > 0}
        <div class="total-overdue">
          Total Tertunggak: <strong>{formatIDR(totalOverdue)}</strong>
        </div>
      {/if}
    </div>

    {#if error}<div class="alert-err">{error}</div>{/if}

    <form class="filter-bar" onsubmit={load}>
      <div class="field">
        <label>Tipe Buku Besar</label>
        <select bind:value={ledgerType}>
          <option value="Keduanya">Keduanya</option>
          <option value="AR">AR (Piutang)</option>
          <option value="AP">AP (Utang)</option>
        </select>
      </div>
      <div class="field">
        <label>Threshold Hari</label>
        <select bind:value={thresholdDays}>
          <option value="30">30 hari</option>
          <option value="60">60 hari</option>
          <option value="90">90 hari</option>
          <option value="custom">Custom</option>
        </select>
      </div>
      {#if thresholdDays === 'custom'}
        <div class="field">
          <label>Hari</label>
          <input type="number" bind:value={customDays} min="1" placeholder="45" />
        </div>
      {/if}
      <button type="submit" class="btn-primary" disabled={loading} style="align-self:flex-end">
        {loading ? 'Memuat...' : 'Terapkan'}
      </button>
    </form>

    {#if loading}
      <div class="loading-center"><span class="material-symbols-outlined spin">progress_activity</span></div>
    {:else if alerts.length === 0}
      <div class="empty-state">
        <span class="material-symbols-outlined">check_circle</span>
        <p>Tidak ada piutang/utang yang mendekati jatuh tempo.</p>
      </div>
    {:else}
      {@const g = grouped()}

      {#if g.red.length > 0}
        <div class="severity-group red">
          <div class="severity-header">
            <span class="material-symbols-outlined">error</span>
            <h3>Kritis — Lebih dari 90 hari ({g.red.length} item)</h3>
          </div>
          <div class="alert-cards">
            {#each g.red as a}
              <div class="alert-card red">
                <div class="ac-top">
                  <span class="ac-id">{a.entity_id ?? a.id}</span>
                  <span class="ac-days">{a.overdue_days} hari</span>
                </div>
                <div class="ac-name">{a.entity_name ?? a.counterparty ?? '-'}</div>
                <div class="ac-amount">{formatIDR(a.amount ?? 0)}</div>
                <div class="ac-due">Jatuh tempo: {a.due_date ?? '-'}</div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if g.orange.length > 0}
        <div class="severity-group orange">
          <div class="severity-header">
            <span class="material-symbols-outlined">warning</span>
            <h3>Perhatian — 60–90 hari ({g.orange.length} item)</h3>
          </div>
          <div class="alert-cards">
            {#each g.orange as a}
              <div class="alert-card orange">
                <div class="ac-top">
                  <span class="ac-id">{a.entity_id ?? a.id}</span>
                  <span class="ac-days">{a.overdue_days} hari</span>
                </div>
                <div class="ac-name">{a.entity_name ?? a.counterparty ?? '-'}</div>
                <div class="ac-amount">{formatIDR(a.amount ?? 0)}</div>
                <div class="ac-due">Jatuh tempo: {a.due_date ?? '-'}</div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if g.yellow.length > 0}
        <div class="severity-group yellow">
          <div class="severity-header">
            <span class="material-symbols-outlined">info</span>
            <h3>Perlu Tindakan — 30–60 hari ({g.yellow.length} item)</h3>
          </div>
          <div class="alert-cards">
            {#each g.yellow as a}
              <div class="alert-card yellow">
                <div class="ac-top">
                  <span class="ac-id">{a.entity_id ?? a.id}</span>
                  <span class="ac-days">{a.overdue_days} hari</span>
                </div>
                <div class="ac-name">{a.entity_name ?? a.counterparty ?? '-'}</div>
                <div class="ac-amount">{formatIDR(a.amount ?? 0)}</div>
                <div class="ac-due">Jatuh tempo: {a.due_date ?? '-'}</div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 96rem; }
  .page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 1.25rem; gap: 1rem; flex-wrap: wrap; }
  .page-head h2 { margin: 0; font-size: 1.4rem; }
  .page-head p { margin: 0.2rem 0 0; font-size: 0.78rem; color: #737686; }
  .total-overdue { background: #fef2f2; border: 1px solid #fecaca; border-radius: 0.35rem; padding: 0.65rem 1rem; font-size: 0.82rem; color: #dc2626; }
  .alert-err { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; border-radius: 0.25rem; padding: 0.65rem 0.85rem; font-size: 0.8rem; margin-bottom: 1rem; }
  .filter-bar { display: flex; gap: 0.75rem; align-items: flex-start; flex-wrap: wrap; margin-bottom: 1.25rem; }
  .field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field label { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #434655; }
  .field input, .field select { border: 1px solid rgb(195 198 215 / 0.55); background: #fff; border-radius: 0.25rem; padding: 0.48rem 0.65rem; font-size: 0.82rem; color: #191c1e; font-family: inherit; outline: none; }
  .btn-primary { display: inline-flex; align-items: center; background: linear-gradient(90deg, #004ac6, #2563eb); color: #fff; border: none; border-radius: 0.25rem; padding: 0.5rem 0.9rem; font-size: 0.8rem; font-weight: 600; cursor: pointer; font-family: inherit; }
  .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
  .loading-center { display: flex; justify-content: center; padding: 3rem; }
  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; padding: 4rem; color: #059669; }
  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.9rem; font-weight: 600; }
  .severity-group { margin-bottom: 1.5rem; }
  .severity-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.75rem; }
  .severity-header h3 { margin: 0; font-size: 0.9rem; font-weight: 700; }
  .severity-group.red .severity-header { color: #dc2626; }
  .severity-group.orange .severity-header { color: #d97706; }
  .severity-group.yellow .severity-header { color: #a16207; }
  .severity-group .material-symbols-outlined { font-size: 1.1rem; }
  .alert-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 0.65rem; }
  .alert-card { background: #fff; border: 1px solid; border-radius: 0.5rem; padding: 0.85rem 1rem; }
  .alert-card.red { border-color: #fca5a5; background: #fef2f2; }
  .alert-card.orange { border-color: #fdba74; background: #fff7ed; }
  .alert-card.yellow { border-color: #fde047; background: #fefce8; }
  .ac-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.35rem; }
  .ac-id { font-family: monospace; font-size: 0.7rem; color: #737686; }
  .ac-days { font-size: 0.68rem; font-weight: 700; }
  .alert-card.red .ac-days { color: #dc2626; }
  .alert-card.orange .ac-days { color: #d97706; }
  .alert-card.yellow .ac-days { color: #a16207; }
  .ac-name { font-size: 0.8rem; font-weight: 600; color: #191c1e; margin-bottom: 0.2rem; }
  .ac-amount { font-size: 1rem; font-weight: 800; font-variant-numeric: tabular-nums; color: #191c1e; }
  .ac-due { font-size: 0.68rem; color: #737686; margin-top: 0.2rem; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; font-size: 1.4rem; color: #737686; }
</style>
