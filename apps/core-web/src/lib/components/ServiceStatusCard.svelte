<script lang="ts">
  import type { ServiceState } from '$lib/state/service-status.svelte';

  let { svc }: { svc: ServiceState } = $props();
</script>

<article
  class={['card', svc.status]}
  data-testid="service-card"
  data-service={svc.name}
  data-status={svc.status}
>
  <header>
    <strong>{svc.name}</strong>
    <span class={['badge', svc.status]}>{svc.status}</span>
  </header>
  {#if svc.error}
    <p class="error" title={svc.error}>{svc.error}</p>
  {/if}
</article>

<style>
  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-left: 4px solid var(--color-muted);
    border-radius: var(--radius-card);
    padding: var(--space-3);
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .card.ok      { border-left-color: var(--color-ok); }
  .card.fail    { border-left-color: var(--color-fail); }
  .card.pending { border-left-color: var(--color-pending); }

  header { display: flex; justify-content: space-between; align-items: baseline; gap: var(--space-2); }

  .badge {
    text-transform: uppercase;
    font-size: 0.7rem;
    font-family: var(--font-mono);
    padding: 0 var(--space-2);
    border-radius: 999px;
    background: var(--color-bg);
    color: var(--color-muted);
  }
  .badge.ok      { color: var(--color-ok); }
  .badge.fail    { color: var(--color-fail); }
  .badge.pending { color: var(--color-pending); }

  .error {
    margin: 0;
    color: var(--color-fail);
    font-size: 0.75rem;
    font-family: var(--font-mono);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
