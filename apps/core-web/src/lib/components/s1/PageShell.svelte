<script lang="ts">
  import type { Snippet } from 'svelte';

  // Shared chrome for S1 B2C catalog + booking shells (S1-L-01). Keeps spacing
  // and back-navigation consistent until a design-system package exists (Q009).

  let {
    title,
    subtitle,
    backHref,
    backLabel = 'Back',
    children
  }: {
    title: string;
    subtitle?: string;
    backHref?: string;
    backLabel?: string;
    children: Snippet;
  } = $props();
</script>

<svelte:head>
  <title>{title} — UmrohOS</title>
</svelte:head>

<div class="shell">
  {#if backHref}
    <p class="back-row">
      <a class="back" href={backHref}>{backLabel}</a>
    </p>
  {/if}
  <header class="head">
    <h1>{title}</h1>
    {#if subtitle}
      <p class="sub">{subtitle}</p>
    {/if}
  </header>
  <div class="body">
    {@render children()}
  </div>
</div>

<style>
  .shell {
    max-width: 48rem;
  }
  .back-row {
    margin: 0 0 var(--space-3) 0;
  }
  .back {
    color: var(--color-accent, #4493f8);
    text-decoration: none;
    font-size: 0.875rem;
  }
  .back:hover {
    text-decoration: underline;
  }
  .head h1 {
    margin: 0 0 var(--space-2) 0;
    font-size: 1.75rem;
    line-height: 1.2;
    font-weight: 600;
  }
  .sub {
    margin: 0 0 var(--space-4) 0;
    color: var(--color-muted);
    font-size: 0.95rem;
    line-height: 1.5;
  }
  .body {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }
</style>
