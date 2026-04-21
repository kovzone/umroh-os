<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title="New booking (draft)"
  subtitle={`Package ${data.packageId} — B2C self-serve channel (stub)`}
  backHref={`/packages/${data.packageId}`}
  backLabel="Back to package"
>
  <section class="s1-card" data-testid="s1-booking-draft-shell">
    <p class="s1-muted">{data.stubNote}</p>

    <form class="form" onsubmit={(e) => e.preventDefault()}>
      <div class="s1-field">
        <label for="lead-name">Lead pilgrim name</label>
        <input
          id="lead-name"
          class="s1-input"
          type="text"
          name="leadName"
          disabled
          placeholder="Wires to draft payload in S1-L-04"
        />
      </div>
      <div class="s1-field">
        <label for="lead-phone">Contact phone</label>
        <input
          id="lead-phone"
          class="s1-input"
          type="tel"
          name="leadPhone"
          disabled
          placeholder="+62 …"
        />
      </div>
      <div class="s1-actions">
        <button type="button" class="s1-btn s1-btn--primary" disabled data-testid="s1-save-draft-stub"
          >Save draft</button
        >
      </div>
    </form>
  </section>
</PageShell>

<style>
  .form {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
    margin-top: var(--space-3);
  }
</style>
