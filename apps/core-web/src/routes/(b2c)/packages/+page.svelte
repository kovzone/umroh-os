<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title="Packages"
  subtitle="Browse departures and start a draft booking. API-backed list ships in S1-L-03."
>
  <section class="s1-card" data-testid="s1-package-catalog">
    <h2>Available packages</h2>
    <p class="s1-muted">
      Loading and empty states will use this layout once <code>GET /v1/packages</code> is integrated.
    </p>
    <ul class="s1-grid list">
      {#each data.packages as pkg (pkg.id)}
        <li class="s1-card tile">
          <h3>{pkg.name}</h3>
          <p class="s1-muted">{pkg.blurb}</p>
          <a
            class="s1-btn s1-btn--primary"
            href={`/packages/${pkg.id}`}
            data-testid="package-link-{pkg.id}">View package</a
          >
        </li>
      {/each}
    </ul>
  </section>
</PageShell>

<style>
  .list {
    list-style: none;
    margin: var(--space-3) 0 0 0;
    padding: 0;
  }
  .tile h3 {
    margin: 0 0 var(--space-2) 0;
    font-size: 1rem;
    font-weight: 600;
  }
</style>
