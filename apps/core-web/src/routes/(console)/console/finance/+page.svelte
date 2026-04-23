<script lang="ts">
  import type { AccountSummary, JournalEntry, PageData } from './+page.server';

  let { data }: { data: PageData } = $props();

  // ---- local state — synced from data via $effect (Svelte 5 rule) ----
  let summary = $state<AccountSummary[]>([]);
  let journals = $state<JournalEntry[]>([]);
  let nextCursor = $state<string | null>(null);

  $effect(() => {
    summary = data.summary ?? [];
    journals = data.journals ?? [];
    nextCursor = data.next_cursor ?? null;
  });

  // ---- journal filter state ----
  let filterFrom = $state('');
  let filterTo = $state('');
  let filterLoading = $state(false);
  let filterError = $state('');

  // ---- expand/collapse journal entries ----
  let expandedIds = $state<Set<string>>(new Set());

  function toggleExpand(id: string) {
    const next = new Set(expandedIds);
    if (next.has(id)) {
      next.delete(id);
    } else {
      next.add(id);
    }
    expandedIds = next;
  }

  // ---- "Load More" pagination state ----
  let loadingMore = $state(false);
  let loadMoreError = $state('');

  // ---- helpers ----
  function formatIDR(amount: number): string {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(amount);
  }

  function formatDate(iso: string): string {
    return new Date(iso).toLocaleDateString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function shortKey(key: string): string {
    return key.length > 40 ? key.slice(0, 37) + '…' : key;
  }

  // ---- filter submit — reload page with query params ----
  async function applyFilter(e: SubmitEvent) {
    e.preventDefault();
    filterLoading = true;
    filterError = '';
    try {
      const params = new URLSearchParams();
      if (filterFrom) params.set('from', filterFrom);
      if (filterTo) params.set('to', filterTo);
      // Navigate with new params — SvelteKit will re-run load
      window.location.href = `/console/finance?${params.toString()}`;
    } catch {
      filterError = 'Gagal menerapkan filter.';
      filterLoading = false;
    }
  }

  function clearFilter() {
    filterFrom = '';
    filterTo = '';
    window.location.href = '/console/finance';
  }

  // ---- load more journals ----
  // Navigates to the same page with cursor param so SvelteKit re-runs load().
  // This is the correct pattern for SSR-backed pagination — no direct API call needed.
  async function loadMore() {
    if (!nextCursor) return;
    loadingMore = true;
    loadMoreError = '';
    try {
      const params = new URLSearchParams(window.location.search);
      params.set('cursor', nextCursor);
      if (filterFrom) params.set('from', filterFrom);
      if (filterTo) params.set('to', filterTo);
      params.set('limit', '50');
      window.location.href = `/console/finance?${params.toString()}`;
    } catch (err) {
      loadMoreError = err instanceof Error ? err.message : 'Gagal memuat lebih banyak jurnal.';
      loadingMore = false;
    }
  }
</script>

<main class="finance-shell">
  <!-- Topbar / breadcrumb -->
  <header class="topbar">
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <span class="material-symbols-outlined breadcrumb-icon">account_balance</span>
      <span class="breadcrumb-text">Keuangan</span>
    </nav>
    <div class="top-actions">
      <button class="icon-btn" title="Notifikasi">
        <span class="material-symbols-outlined">notifications</span>
      </button>
      <button class="avatar" aria-label="Profile">FN</button>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Laporan Keuangan</h2>
        <p>Ringkasan saldo akun dan riwayat jurnal transaksi</p>
      </div>
    </div>

    {#if data.error}
      <div class="error-banner" role="alert">
        <span class="material-symbols-outlined">error</span>
        {data.error}
      </div>
    {/if}

    <!-- ================================================================
         Section 1 — Ringkasan Akun
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">account_balance_wallet</span>
        <h3>Ringkasan Akun</h3>
      </div>

      {#if summary.length === 0}
        <div class="empty-state">
          <span class="material-symbols-outlined">account_balance</span>
          <p>Belum ada data ringkasan akun.</p>
        </div>
      {:else}
        <div class="panel">
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Kode Akun</th>
                  <th class="align-right">Total Debit</th>
                  <th class="align-right">Total Kredit</th>
                  <th class="align-right">Saldo</th>
                </tr>
              </thead>
              <tbody>
                {#each summary as acct (acct.account_code)}
                  <tr>
                    <td>
                      <span class="account-code">{acct.account_code}</span>
                    </td>
                    <td class="align-right amount-debit">{formatIDR(acct.debit_total)}</td>
                    <td class="align-right amount-credit">{formatIDR(acct.credit_total)}</td>
                    <td class="align-right">
                      <span class="balance-cell" class:balance-pos={acct.net >= 0} class:balance-neg={acct.net < 0}>
                        {formatIDR(Math.abs(acct.net))}
                        {#if acct.net < 0}
                          <span class="balance-sign">&minus;</span>
                        {/if}
                      </span>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          <div class="table-footer">
            {summary.length} akun
          </div>
        </div>
      {/if}
    </section>

    <!-- ================================================================
         Section 2 — Riwayat Jurnal
    ================================================================= -->
    <section class="section-block">
      <div class="section-header">
        <span class="material-symbols-outlined section-icon">receipt_long</span>
        <h3>Riwayat Jurnal</h3>
      </div>

      <!-- Date range filter -->
      <form class="filter-bar" onsubmit={applyFilter}>
        <div class="filter-group">
          <label for="filter-from">Dari Tanggal</label>
          <input
            id="filter-from"
            type="date"
            bind:value={filterFrom}
          />
        </div>
        <div class="filter-group">
          <label for="filter-to">Sampai Tanggal</label>
          <input
            id="filter-to"
            type="date"
            bind:value={filterTo}
          />
        </div>
        <div class="filter-actions">
          <button type="submit" class="primary-btn" disabled={filterLoading}>
            {#if filterLoading}
              <span class="material-symbols-outlined spin">progress_activity</span>
              Memuat...
            {:else}
              <span class="material-symbols-outlined">filter_alt</span>
              Terapkan
            {/if}
          </button>
          {#if filterFrom || filterTo}
            <button type="button" class="ghost-btn" onclick={clearFilter}>
              <span class="material-symbols-outlined">close</span>
              Reset
            </button>
          {/if}
        </div>
        {#if filterError}
          <p class="inline-error">{filterError}</p>
        {/if}
      </form>

      {#if journals.length === 0}
        <div class="empty-state">
          <span class="material-symbols-outlined">receipt_long</span>
          <p>Belum ada jurnal.</p>
        </div>
      {:else}
        <div class="panel journal-panel">
          {#each journals as entry (entry.id)}
            {@const isExpanded = expandedIds.has(entry.id)}
            <div class="journal-entry" class:expanded={isExpanded}>
              <!-- Entry header row (clickable) -->
              <!-- svelte-ignore a11y_interactive_supports_focus -->
              <div
                class="journal-row"
                role="button"
                aria-expanded={isExpanded}
                onclick={() => toggleExpand(entry.id)}
                onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') toggleExpand(entry.id); }}
                tabindex="0"
              >
                <span class="material-symbols-outlined expand-icon">
                  {isExpanded ? 'expand_less' : 'expand_more'}
                </span>
                <div class="journal-meta">
                  <span class="journal-date">{formatDate(entry.posted_at)}</span>
                  <span class="journal-desc">{entry.description}</span>
                </div>
                <div class="journal-key-wrap">
                  <span class="journal-key" title={entry.idempotency_key}>
                    {shortKey(entry.idempotency_key)}
                  </span>
                </div>
                <span class="journal-lines-count">
                  {entry.lines.length} baris
                </span>
              </div>

              <!-- Expanded lines detail -->
              {#if isExpanded}
                <div class="journal-detail">
                  <table class="lines-table">
                    <thead>
                      <tr>
                        <th>Kode Akun</th>
                        <th class="align-right">Debit</th>
                        <th class="align-right">Kredit</th>
                      </tr>
                    </thead>
                    <tbody>
                      {#each entry.lines as line, i (i)}
                        <tr>
                          <td><span class="account-code">{line.account_code}</span></td>
                          <td class="align-right">
                            {#if line.debit > 0}
                              <span class="amount-debit">{formatIDR(line.debit)}</span>
                            {:else}
                              <span class="zero">—</span>
                            {/if}
                          </td>
                          <td class="align-right">
                            {#if line.credit > 0}
                              <span class="amount-credit">{formatIDR(line.credit)}</span>
                            {:else}
                              <span class="zero">—</span>
                            {/if}
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                </div>
              {/if}
            </div>
          {/each}

          <!-- Pagination -->
          <div class="journal-footer">
            <span class="journal-count">{journals.length} entri ditampilkan</span>
            {#if nextCursor}
              <button
                type="button"
                class="load-more-btn"
                onclick={loadMore}
                disabled={loadingMore}
              >
                {#if loadingMore}
                  <span class="material-symbols-outlined spin">progress_activity</span>
                  Memuat...
                {:else}
                  <span class="material-symbols-outlined">expand_more</span>
                  Muat Lebih Banyak
                {/if}
              </button>
            {/if}
            {#if loadMoreError}
              <p class="inline-error">{loadMoreError}</p>
            {/if}
          </div>
        </div>
      {/if}
    </section>
  </section>
</main>

<style>
  .finance-shell {
    min-height: 100vh;
    background: #f7f9fb;
  }

  /* ---- topbar ---- */
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

  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #434655;
  }

  .breadcrumb-icon {
    font-size: 1.1rem;
    color: #004ac6;
  }

  .breadcrumb-text {
    font-size: 0.88rem;
    font-weight: 600;
    color: #191c1e;
  }

  .top-actions {
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

  .icon-btn:hover { background: #eceef0; }

  .avatar {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #b4c5ff;
    color: #00174b;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    font-weight: 700;
    font-size: 0.65rem;
    cursor: pointer;
  }

  /* ---- canvas ---- */
  .canvas {
    padding: 1.5rem;
    max-width: 96rem;
  }

  .page-head {
    margin-bottom: 1.25rem;
  }

  .page-head h2 {
    margin: 0;
    font-size: 1.5rem;
    line-height: 1.2;
  }

  .page-head p {
    margin: 0.3rem 0 0;
    font-size: 0.82rem;
    color: #434655;
  }

  /* ---- error banner ---- */
  .error-banner {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: #ffdad6;
    color: #93000a;
    border-radius: 0.25rem;
    padding: 0.65rem 0.85rem;
    font-size: 0.82rem;
    margin-bottom: 1.25rem;
  }

  /* ---- section blocks ---- */
  .section-block {
    margin-bottom: 2rem;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.85rem;
  }

  .section-header h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: #191c1e;
  }

  .section-icon {
    font-size: 1.1rem;
    color: #004ac6;
  }

  /* ---- panel ---- */
  .panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .table-wrap { overflow-x: auto; }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    padding: 0.62rem 0.85rem;
    text-align: left;
    font-size: 0.78rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    white-space: nowrap;
  }

  th {
    text-transform: uppercase;
    font-size: 0.62rem;
    letter-spacing: 0.08em;
    color: #434655;
    background: #f2f4f6;
  }

  tbody tr:hover { background: #f7f9fb; }
  tbody tr:last-child td { border-bottom: 0; }

  .align-right { text-align: right; }

  /* ---- account summary cells ---- */
  .account-code {
    display: inline-flex;
    padding: 0.1rem 0.4rem;
    background: #eceef0;
    border-radius: 0.2rem;
    font-size: 0.7rem;
    font-weight: 700;
    color: #191c1e;
    font-variant-numeric: tabular-nums;
  }

  .account-name {
    font-weight: 500;
  }

  .amount-debit {
    color: #065f46;
    font-variant-numeric: tabular-nums;
  }

  .amount-credit {
    color: #1e3a8a;
    font-variant-numeric: tabular-nums;
  }

  .balance-cell {
    display: inline-flex;
    align-items: center;
    gap: 0.1rem;
    font-weight: 700;
    font-variant-numeric: tabular-nums;
  }

  .balance-pos { color: #065f46; }
  .balance-neg { color: #ba1a1a; }

  .balance-sign {
    font-weight: 700;
    font-size: 0.9em;
  }

  .zero { color: #b0b3c1; }

  .table-footer {
    padding: 0.55rem 0.85rem;
    font-size: 0.68rem;
    color: #434655;
    border-top: 1px solid rgb(195 198 215 / 0.35);
    background: #f7f9fb;
  }

  /* ---- empty state ---- */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.6rem;
    padding: 3rem 1rem;
    color: #b0b3c1;
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
  }

  .empty-state .material-symbols-outlined { font-size: 2.5rem; }
  .empty-state p { margin: 0; font-size: 0.82rem; }

  /* ---- filter bar ---- */
  .filter-bar {
    display: flex;
    align-items: flex-end;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-bottom: 1rem;
  }

  .filter-group {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
  }

  .filter-group label {
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #434655;
  }

  .filter-group input {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    border-radius: 0.25rem;
    padding: 0.42rem 0.6rem;
    font-size: 0.82rem;
    color: #191c1e;
    outline: none;
    font-family: inherit;
  }

  .filter-actions {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    align-self: flex-end;
  }

  .inline-error {
    margin: 0;
    font-size: 0.75rem;
    color: #ba1a1a;
    align-self: flex-end;
  }

  /* ---- journal panel ---- */
  .journal-panel {
    overflow: visible;
  }

  .journal-entry {
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
  }

  .journal-entry:last-of-type {
    border-bottom: 0;
  }

  .journal-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 0.85rem;
    cursor: pointer;
    user-select: none;
    transition: background 0.1s;
  }

  .journal-row:hover { background: #f7f9fb; }
  .journal-entry.expanded > .journal-row { background: #f2f4f6; }

  .expand-icon {
    font-size: 1.05rem;
    color: #737686;
    flex-shrink: 0;
  }

  .journal-meta {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
    min-width: 0;
    flex: 1;
  }

  .journal-date {
    font-size: 0.68rem;
    color: #737686;
    white-space: nowrap;
  }

  .journal-desc {
    font-size: 0.82rem;
    font-weight: 500;
    color: #191c1e;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .journal-key-wrap {
    flex-shrink: 0;
  }

  .journal-key {
    display: inline-flex;
    padding: 0.1rem 0.45rem;
    border-radius: 0.2rem;
    background: #e0f2fe;
    color: #075985;
    font-size: 0.65rem;
    font-weight: 600;
    font-family: 'IBM Plex Mono', 'Courier New', monospace;
    white-space: nowrap;
  }

  .journal-lines-count {
    font-size: 0.68rem;
    color: #737686;
    white-space: nowrap;
    flex-shrink: 0;
  }

  /* ---- expanded lines table ---- */
  .journal-detail {
    padding: 0 0.85rem 0.85rem;
    background: #f7f9fb;
    border-top: 1px solid rgb(195 198 215 / 0.35);
  }

  .lines-table {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.2rem;
    overflow: hidden;
    margin-top: 0.6rem;
    background: #fff;
  }

  .lines-table th {
    font-size: 0.6rem;
    padding: 0.45rem 0.7rem;
  }

  .lines-table td {
    padding: 0.45rem 0.7rem;
    font-size: 0.76rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.35);
  }

  .lines-table tbody tr:last-child td {
    border-bottom: 0;
  }

  /* ---- journal footer ---- */
  .journal-footer {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.65rem 0.85rem;
    border-top: 1px solid rgb(195 198 215 / 0.35);
    background: #f7f9fb;
    flex-wrap: wrap;
  }

  .journal-count {
    font-size: 0.68rem;
    color: #434655;
    flex: 1;
  }

  /* ---- buttons ---- */
  .ghost-btn,
  .primary-btn {
    border-radius: 0.25rem;
    padding: 0.48rem 0.85rem;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    border: 1px solid rgb(195 198 215 / 0.55);
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    font-family: inherit;
  }

  .ghost-btn {
    background: #fff;
    color: #191c1e;
  }

  .ghost-btn:hover { background: #f2f4f6; }

  .primary-btn {
    border-color: #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
  }

  .primary-btn:disabled,
  .ghost-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .primary-btn .material-symbols-outlined,
  .ghost-btn .material-symbols-outlined {
    font-size: 0.95rem;
  }

  .load-more-btn {
    border-radius: 0.25rem;
    padding: 0.4rem 0.75rem;
    font-size: 0.78rem;
    font-weight: 600;
    cursor: pointer;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    color: #004ac6;
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-family: inherit;
  }

  .load-more-btn:hover { background: #f2f4f6; }
  .load-more-btn:disabled { opacity: 0.6; cursor: not-allowed; }
  .load-more-btn .material-symbols-outlined { font-size: 0.95rem; }

  /* ---- spin animation ---- */
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  .spin { animation: spin 0.8s linear infinite; font-size: 0.95rem; }
</style>
