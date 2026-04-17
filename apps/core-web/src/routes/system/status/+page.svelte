<script lang="ts">
  import ServiceStatusCard from '$lib/components/ServiceStatusCard.svelte';
  import { ServiceStatus, BACKEND_SERVICES } from '$lib/state/service-status.svelte';

  // Single instance per page mount; the class manages its own polling lifecycle
  // via createSubscriber, so it self-cleans when the page unmounts.
  const status = new ServiceStatus(BACKEND_SERVICES);
</script>

<svelte:head>
  <title>UmrohOS — Service Status</title>
</svelte:head>

<section>
  <h1>Service status</h1>
  <p class="muted" data-testid="last-poll">
    Last poll: {status.lastPolledAt ? status.lastPolledAt.toLocaleTimeString() : 'pending…'}
  </p>

  <div class="grid" data-testid="status-grid">
    {#each status.services as svc (svc.name)}
      <ServiceStatusCard {svc} />
    {/each}
  </div>
</section>

<style>
  h1 { margin: 0 0 var(--space-2) 0; }
  .muted { color: var(--color-muted); font-size: 0.875rem; margin-bottom: var(--space-4); }
  .grid {
    display: grid;
    gap: var(--space-3);
    grid-template-columns: repeat(auto-fill, minmax(14rem, 1fr));
  }
</style>
