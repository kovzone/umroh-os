<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title="Detail keberangkatan"
  subtitle={`${data.package.name} · ${data.departure.departureDate} s.d. ${data.departure.returnDate}`}
  backHref={`/packages/${data.package.id}`}
  backLabel="Package"
>
  <section class="s1-card" data-testid="s1-departure-detail">
    <dl class="meta">
      <div>
        <dt>Total kursi</dt>
        <dd>{data.departure.totalSeats}</dd>
      </div>
      <div>
        <dt>Sisa kursi</dt>
        <dd>{data.departure.remainingSeats}</dd>
      </div>
      <div>
        <dt>Status</dt>
        <dd class="status">{data.departure.status}</dd>
      </div>
    </dl>

    <h2>Harga per tipe kamar</h2>
    <ul class="prices">
      {#each data.departure.pricing as price (price.roomType)}
        <li>
          <span class="room">{price.roomType}</span>
          <strong>{price.amountLabel}</strong>
        </li>
      {/each}
    </ul>

    <div class="s1-actions">
      <a class="s1-btn s1-btn--primary" href={`/booking/${data.package.id}`} data-testid="s1-book-from-departure"
        >Book this departure</a
      >
    </div>
  </section>
</PageShell>

<style>
  .meta {
    margin: 0 0 var(--space-4) 0;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(9rem, 1fr));
    gap: var(--space-3);
  }
  .meta div {
    border: 1px solid var(--color-border);
    border-radius: calc(var(--radius-card) - 2px);
    padding: var(--space-3);
  }
  .meta dt {
    font-size: 0.75rem;
    color: var(--color-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  .meta dd {
    margin: var(--space-1) 0 0 0;
    font-size: 1.1rem;
    font-weight: 600;
  }
  .status {
    text-transform: capitalize;
  }
  .prices {
    list-style: none;
    margin: var(--space-3) 0 var(--space-4) 0;
    padding: 0;
    display: grid;
    gap: var(--space-2);
  }
  .prices li {
    border: 1px solid var(--color-border);
    border-radius: calc(var(--radius-card) - 2px);
    padding: var(--space-3);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .room {
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-muted);
    font-size: 0.78rem;
  }
</style>
