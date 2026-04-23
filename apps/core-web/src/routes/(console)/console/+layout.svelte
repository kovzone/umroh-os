<script lang="ts">
  import { page } from '$app/state';
  import { onMount } from 'svelte';

  let { children, data } = $props();
  let sidebarCollapsed = $state(false);
  const isLoginPage = $derived(data?.pathname === '/console/login');

  const navItems = [
    { label: 'Dashboard', icon: 'dashboard', href: '/console', badge: '' },
    { label: 'Katalog Paket', icon: 'inventory_2', href: '/console/packages', badge: '' },
    { label: 'Keberangkatan', icon: 'flight_takeoff', href: '/console/departures', badge: '' },
    { label: 'Booking', icon: 'confirmation_number', href: '/console/bookings', badge: '12' },
    { label: 'Jamaah', icon: 'groups', href: '/console/jamaah', badge: '' },
    { label: 'Pembayaran', icon: 'payments', href: '/console/payments', badge: '' },
    { label: 'Visa & Dokumen', icon: 'description', href: '/console/visa-docs', badge: '5' },
    { label: 'Operasional', icon: 'settings_applications', href: '/console/operations', badge: '' },
    { label: 'Ops Board', icon: 'local_shipping', href: '/console/ops', badge: '' },
    { label: 'Logistik', icon: 'conveyor_belt', href: '/console/logistics', badge: '' },
    { label: 'Finance', icon: 'account_balance', href: '/console/finance', badge: '3' },
    { label: 'CRM', icon: 'hub', href: '/console/crm', badge: '' },
    { label: 'Leads', icon: 'person_add', href: '/console/leads', badge: '' },
    { label: 'IAM & Akses', icon: 'admin_panel_settings', href: '/console/iam', badge: '' },
    { label: 'Audit Log', icon: 'history_edu', href: '/console/audit-log', badge: '' },
    { label: 'Pengaturan Sistem', icon: 'settings', href: '/console/system-settings', badge: '' },
    { label: 'Bantuan IT', icon: 'contact_support', href: '/console/help', badge: '' }
  ];

  function isActive(href: string): boolean {
    if (href === '/console') {
      return page.url.pathname === '/console';
    }
    return page.url.pathname.startsWith(href);
  }

  function toggleSidebar(): void {
    sidebarCollapsed = !sidebarCollapsed;
    localStorage.setItem('console-sidebar-collapsed', String(sidebarCollapsed));
  }

  onMount(() => {
    const saved = localStorage.getItem('console-sidebar-collapsed');
    sidebarCollapsed = saved === 'true';
  });
</script>

<svelte:head>
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@400;500;600;700&display=swap"
  />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
  />
</svelte:head>

{#if isLoginPage}
  <div class="login-wrap">
    {@render children()}
  </div>
{:else}
  <div class="console-shell">
    <aside class="sidebar" class:collapsed={sidebarCollapsed}>
      <div class="brand-row">
        <div class="logo">O</div>
        {#if !sidebarCollapsed}
          <div class="brand-copy">
            <h1>Internal Ops</h1>
            <p>Precision Ledger v2.4</p>
          </div>
        {/if}
        <button type="button" class="collapse-btn" onclick={toggleSidebar} aria-label="Toggle sidebar">
          <span class="material-symbols-outlined">
            {sidebarCollapsed ? 'keyboard_double_arrow_right' : 'keyboard_double_arrow_left'}
          </span>
        </button>
      </div>

      <nav class="menu">
        {#each navItems as item, i}
          {#if i === 7 || i === 11}
            <div class="divider"></div>
          {/if}
          <a href={item.href} class="menu-item" class:active={isActive(item.href)} title={item.label}>
            <span class="material-symbols-outlined">{item.icon}</span>
            {#if !sidebarCollapsed}
              <span class="lbl">{item.label}</span>
              {#if item.badge}
                <span class="badge">{item.badge}</span>
              {/if}
            {:else if item.badge}
              <span class="badge-dot"></span>
            {/if}
          </a>
        {/each}
      </nav>
    </aside>

    <div class="content-wrap" class:collapsed={sidebarCollapsed}>
      {@render children()}
    </div>
  </div>
{/if}

<style>
  :global(.material-symbols-outlined) {
    font-family: 'Material Symbols Outlined', sans-serif;
    font-variation-settings: 'FILL' 0, 'wght' 450, 'GRAD' 0, 'opsz' 24;
  }

  .console-shell {
    min-height: 100vh;
    min-height: 100dvh;
    font-family:
      'IBM Plex Sans',
      ui-sans-serif,
      system-ui,
      -apple-system,
      sans-serif;
    background: #f7f9fb;
    color: #191c1e;
  }

  .login-wrap {
    min-height: 100vh;
    min-height: 100dvh;
    background: #f7f9fb;
  }

  .sidebar {
    position: fixed;
    inset: 0 auto 0 0;
    width: 16rem;
    background: #f2f4f6;
    border-right: 1px solid rgb(195 198 215 / 0.45);
    display: flex;
    flex-direction: column;
    z-index: 40;
    transition: width 150ms ease;
  }

  .sidebar.collapsed {
    width: 4.5rem;
  }

  .brand-row {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.9rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
  }

  .logo {
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    background: #2563eb;
    color: #fff;
    display: grid;
    place-items: center;
    font-weight: 700;
    flex-shrink: 0;
  }

  .brand-copy h1 {
    margin: 0;
    font-size: 1.05rem;
    line-height: 1.1;
  }

  .brand-copy p {
    margin: 0.1rem 0 0;
    font-size: 0.68rem;
    color: #434655;
  }

  .collapse-btn {
    margin-left: auto;
    border: 0;
    background: transparent;
    color: #434655;
    cursor: pointer;
    padding: 0.2rem;
    border-radius: 0.25rem;
  }

  .collapse-btn:hover {
    background: #e6e8ea;
  }

  .menu {
    padding: 0.65rem;
    overflow-y: auto;
    display: grid;
    gap: 0.25rem;
    flex: 1;
  }

  .menu-item {
    display: flex;
    align-items: center;
    gap: 0.65rem;
    padding: 0.55rem 0.6rem;
    border-radius: 0.25rem;
    color: #434655;
    text-decoration: none;
    position: relative;
  }

  .menu-item .material-symbols-outlined {
    font-size: 1.08rem;
  }

  .menu-item:hover {
    background: #e6e8ea;
  }

  .menu-item.active {
    background: #fff;
    color: #004ac6;
    border: 1px solid rgb(195 198 215 / 0.5);
  }

  .lbl {
    font-size: 0.82rem;
    font-weight: 500;
    white-space: nowrap;
  }

  .badge {
    margin-left: auto;
    font-size: 0.64rem;
    font-weight: 700;
    background: #ba1a1a;
    color: #fff;
    border-radius: 0.2rem;
    padding: 0.05rem 0.32rem;
  }

  .badge-dot {
    width: 0.36rem;
    height: 0.36rem;
    border-radius: 999px;
    background: #ba1a1a;
    position: absolute;
    right: 0.6rem;
    top: 0.45rem;
  }

  .divider {
    height: 1px;
    background: rgb(195 198 215 / 0.45);
    margin: 0.55rem 0.15rem;
  }

  .content-wrap {
    margin-left: 16rem;
    min-height: 100vh;
    min-height: 100dvh;
    transition: margin-left 150ms ease;
    background: #f7f9fb;
  }

  .content-wrap.collapsed {
    margin-left: 4.5rem;
  }
</style>

