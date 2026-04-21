<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title="Package detail"
  subtitle={`Package id: ${data.packageId}`}
  backHref="/packages"
  backLabel="All packages"
>
  <section class="s1-card" data-testid="s1-package-detail">
    <p class="s1-muted">{data.stubNote}</p>
    <p class="s1-muted">
      Error handling for <code>404</code> / invalid id will mirror the gateway error shape in S1-L-03.
    </p>
    <div class="s1-actions">
      <a
        class="s1-btn s1-btn--ghost"
        href={`/packages/${data.packageId}/departures/demo-departure-1`}
        data-testid="s1-departure-sample-link">Sample departure</a
      >
      <a
        class="s1-btn s1-btn--primary"
        href={`/booking/${data.packageId}`}
        data-testid="s1-start-booking">Start draft booking</a
      >
    </div>
  </section>
</PageShell>
