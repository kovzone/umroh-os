<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title={data.package.name}
  subtitle={data.package.description}
  backHref="/packages"
  backLabel="All packages"
>
  <section class="s1-card" data-testid="s1-package-detail">
    <img src={data.package.coverPhotoUrl} alt={data.package.name} class="cover" loading="lazy" />
    <ul class="highlights">
      {#each data.package.highlights as item (item)}
        <li>{item}</li>
      {/each}
    </ul>

    <h2>Pilihan keberangkatan</h2>
    <ul class="departures">
      {#each data.package.departures as departure (departure.id)}
        <li>
          <div>
            <p class="date">
              {departure.departureDate} s.d. {departure.returnDate}
            </p>
            <p class="s1-muted">Status: {departure.status} · Sisa {departure.remainingSeats} kursi</p>
          </div>
          <a class="s1-btn s1-btn--ghost" href={`/packages/${data.package.id}/departures/${departure.id}`}>
            Lihat detail
          </a>
        </li>
      {/each}
    </ul>

    <div class="s1-actions">
      <a
        class="s1-btn s1-btn--primary"
        href={`/booking/${data.package.id}`}
        data-testid="s1-start-booking">Start draft booking</a
      >
    </div>
  </section>
</PageShell>

<style>
  .cover {
    width: 100%;
    max-height: 15rem;
    object-fit: cover;
    border-radius: calc(var(--radius-card) - 2px);
    border: 1px solid var(--color-border);
    margin-bottom: var(--space-3);
  }
  .highlights {
    margin: 0 0 var(--space-4) 0;
    padding-left: 1.1rem;
    color: var(--color-text);
  }
  .departures {
    list-style: none;
    margin: var(--space-3) 0 var(--space-4) 0;
    padding: 0;
    display: grid;
    gap: var(--space-3);
  }
  .departures li {
    display: flex;
    justify-content: space-between;
    gap: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: calc(var(--radius-card) - 2px);
    padding: var(--space-3);
  }
  .date {
    margin: 0 0 var(--space-1) 0;
    font-weight: 600;
  }
</style>
