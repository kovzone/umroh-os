<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title="Paket Umrah"
  subtitle="Pilih paket, cek keberangkatan, lalu lanjut ke form draft booking. Data saat ini masih mock (siap diganti ke API katalog)."
>
  <section class="s1-card" data-testid="s1-package-catalog">
    <h2 class="section-title">Pilihan paket tersedia</h2>
    <ul class="s1-grid list">
      {#each data.packages as pkg (pkg.id)}
        <li class="s1-card tile">
          <img src={pkg.coverPhotoUrl} alt={pkg.name} class="cover" loading="lazy" />
          <p class="badge">{pkg.kind}</p>
          <h3>{pkg.name}</h3>
          <p class="s1-muted">{pkg.blurb}</p>
          <dl class="meta">
            <div>
              <dt>Harga mulai</dt>
              <dd>{pkg.startingPriceLabel}</dd>
            </div>
            <div>
              <dt>Keberangkatan terdekat</dt>
              <dd>{pkg.nextDepartureLabel}</dd>
            </div>
            <div>
              <dt>Sisa kursi</dt>
              <dd>{pkg.remainingSeats} kursi</dd>
            </div>
          </dl>
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
  .section-title {
    margin: 0 0 var(--space-3) 0;
  }
  .list {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .cover {
    width: 100%;
    height: 9rem;
    object-fit: cover;
    border-radius: calc(var(--radius-card) - 2px);
    margin-bottom: var(--space-3);
    border: 1px solid var(--color-border);
  }
  .badge {
    margin: 0 0 var(--space-2) 0;
    color: var(--color-accent);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .tile h3 {
    margin: 0 0 var(--space-2) 0;
    font-size: 1rem;
    font-weight: 600;
  }
  .meta {
    margin: var(--space-3) 0;
    display: grid;
    gap: var(--space-2);
  }
  .meta div {
    display: grid;
    gap: 0;
  }
  .meta dt {
    font-size: 0.78rem;
    color: var(--color-muted);
  }
  .meta dd {
    margin: 0;
    font-size: 0.92rem;
    font-weight: 500;
  }
</style>
