<script lang="ts">
  // Capability tile used by the landing page's "What's inside" grid.
  // Groundwork for future dashboard/summary cards (admin home, executive
  // dashboard, etc.) — same shape, swap the data.

  export type CapabilityStatus = 'planned' | 'in-progress' | 'live';

  interface Props {
    title: string;
    description: string;
    status: CapabilityStatus;
  }

  let { title, description, status }: Props = $props();
</script>

<article class={['card', status]} data-testid="capability-card" data-status={status}>
  <header>
    <h3>{title}</h3>
    <span class={['badge', status]}>{status}</span>
  </header>
  <p>{description}</p>
</article>

<style>
  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-top: 3px solid var(--color-muted);
    border-radius: var(--radius-card);
    padding: var(--space-3);
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .card.live        { border-top-color: var(--color-ok); }
  .card.in-progress { border-top-color: var(--color-pending); }
  .card.planned     { border-top-color: var(--color-muted); }

  header { display: flex; justify-content: space-between; align-items: baseline; gap: var(--space-2); }
  h3 { margin: 0; font-size: 1.05rem; }

  .badge {
    text-transform: uppercase;
    font-size: 0.7rem;
    font-family: var(--font-mono);
    padding: 0 var(--space-2);
    border-radius: 999px;
    background: var(--color-bg);
    color: var(--color-muted);
  }
  .badge.live        { color: var(--color-ok); }
  .badge.in-progress { color: var(--color-pending); }

  p { color: var(--color-muted); margin: 0; font-size: 0.875rem; line-height: 1.5; }
</style>
