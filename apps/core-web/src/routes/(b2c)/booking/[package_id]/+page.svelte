<script lang="ts">
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();
</script>

<PageShell
  title="Form draft booking"
  subtitle={`${data.package.name} · channel b2c_self`}
  backHref={`/packages/${data.package.id}`}
  backLabel="Back to package"
>
  <div class="layout">
    <section class="s1-card" data-testid="s1-booking-draft-shell">
      <h2>Data ketua rombongan</h2>
      <form class="form" onsubmit={(e) => e.preventDefault()}>
        <div class="s1-field">
          <label for="lead-name">Nama lengkap</label>
          <input id="lead-name" class="s1-input" type="text" name="leadName" placeholder="Budi Santoso" />
        </div>
        <div class="s1-field">
          <label for="lead-whatsapp">WhatsApp</label>
          <input id="lead-whatsapp" class="s1-input" type="tel" name="leadWhatsapp" placeholder="+62812..." />
        </div>
        <div class="s1-field">
          <label for="lead-email">Email</label>
          <input id="lead-email" class="s1-input" type="email" name="leadEmail" placeholder="nama@email.com" />
        </div>
        <div class="s1-field">
          <label for="lead-domicile">Domisili</label>
          <input id="lead-domicile" class="s1-input" type="text" name="leadDomicile" placeholder="Jakarta" />
        </div>
        <div class="s1-field">
          <label for="jamaah-count">Jumlah jamaah</label>
          <input id="jamaah-count" class="s1-input" type="number" min="1" value="1" />
        </div>
        <div class="s1-actions">
          <button type="button" class="s1-btn s1-btn--primary" data-testid="s1-save-draft-stub">Simpan draft</button>
          <span class="s1-muted">Submit API `POST /v1/bookings` akan diaktifkan di S1-L-04.</span>
        </div>
      </form>
    </section>

    <aside class="s1-card summary">
      <h2>Ringkasan paket</h2>
      <p class="package">{data.package.name}</p>
      <p class="s1-muted">{data.package.startingPriceLabel}</p>
      <ul>
        {#each data.package.highlights as item (item)}
          <li>{item}</li>
        {/each}
      </ul>
      <a class="s1-btn s1-btn--ghost" href={`/packages/${data.package.id}`}>Kembali ke detail paket</a>
    </aside>
  </div>
</PageShell>

<style>
  .layout {
    display: grid;
    gap: var(--space-4);
    grid-template-columns: 1.4fr 1fr;
  }
  .summary {
    align-self: start;
  }
  .summary ul {
    margin: var(--space-3) 0;
    padding-left: 1rem;
  }
  .summary li {
    margin-bottom: var(--space-2);
  }
  .package {
    margin: 0 0 var(--space-1) 0;
    font-size: 1rem;
    font-weight: 600;
  }
  .form {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
    margin-top: var(--space-3);
  }
  @media (max-width: 860px) {
    .layout {
      grid-template-columns: 1fr;
    }
  }
</style>
