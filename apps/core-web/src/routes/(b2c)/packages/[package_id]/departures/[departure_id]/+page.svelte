<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title="Departure"
  subtitle={`Departure ${data.departureId} — package ${data.packageId}`}
  backHref={`/packages/${data.packageId}`}
  backLabel="Package"
>
  <section class="s1-card" data-testid="s1-departure-detail">
    <p class="s1-muted">{data.stubNote}</p>
    <p class="s1-muted">
      Remaining seats and date labels will render here from the catalog read model.
    </p>
    <div class="s1-actions">
      <a class="s1-btn s1-btn--primary" href={`/booking/${data.packageId}`} data-testid="s1-book-from-departure"
        >Book this departure</a
      >
    </div>
  </section>
</PageShell>
