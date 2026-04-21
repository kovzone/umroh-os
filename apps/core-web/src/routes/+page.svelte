<script lang="ts">
  import CapabilityCard from '$lib/components/CapabilityCard.svelte';

  // Capabilities shown on the landing page — one tile per high-level module
  // group from docs/00-overview/03-module-list.md. Keep these to 4–6 items;
  // the full 202-module catalogue lives in the PRD, not on the landing page.
  const capabilities: Array<{ title: string; description: string; status: 'planned' | 'in-progress' | 'live' }> = [
    {
      title: 'Booking & packages',
      description:
        'B2C and B2B booking across package catalogues, with per-departure seat reservation and family-aware allocation.',
      status: 'planned'
    },
    {
      title: 'Jamaah, documents & visa',
      description:
        'Pilgrim profiles, mahram graph, OCR-assisted passport reading, and the MOFA / Sajil / Nusuk visa pipeline.',
      status: 'planned'
    },
    {
      title: 'Finance (PSAK)',
      description:
        'Double-entry accounting aligned with Indonesian PSAK standards: AR / AP, tax (PPh, PPN), multi-currency FX, and per-departure cost rollup.',
      status: 'planned'
    },
    {
      title: 'Field operations',
      description:
        'Saudi-side manifests, luggage-tag QR, airport handling, Raudhah Shield monitoring, and real-time incident tracking.',
      status: 'planned'
    }
  ];
</script>

<svelte:head>
  <title>UmrohOS — ERP for Umrah and Hajj travel agencies</title>
</svelte:head>

<section class="hero" data-testid="landing-hero">
  <h1>UmrohOS</h1>
  <p class="tagline">End-to-end ERP for Umrah and Hajj travel agencies.</p>
  <p class="description">
    One platform for B2C package booking, B2B agent networks, marketing and CRM, document
    and visa processing, warehouse logistics, PSAK-compliant finance, and jamaah field
    operations in Saudi Arabia.
  </p>
  <p class="cta-row">
    <a class="cta" href="/packages" data-testid="browse-packages-cta">Browse packages</a>
    <span class="cta-hint">S1 catalog shell — API in S1-L-03</span>
  </p>
</section>

<section class="capabilities">
  <h2>What&rsquo;s inside</h2>
  <div class="grid" data-testid="capability-grid">
    {#each capabilities as c (c.title)}
      <CapabilityCard title={c.title} description={c.description} status={c.status} />
    {/each}
  </div>
</section>

<style>
  .hero {
    padding: var(--space-4) 0;
    border-bottom: 1px solid var(--color-border);
    margin-bottom: var(--space-4);
  }
  .hero h1 { font-size: 2.5rem; margin: 0 0 var(--space-2) 0; line-height: 1.1; }
  .tagline { font-size: 1.125rem; margin: 0 0 var(--space-3) 0; }
  .description {
    color: var(--color-muted);
    max-width: 44rem;
    margin: 0;
    line-height: 1.6;
  }
  .cta-row {
    margin: var(--space-4) 0 0 0;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-3);
  }
  .cta {
    display: inline-flex;
    align-items: center;
    padding: 0.55rem 1.1rem;
    border-radius: var(--radius-card);
    background: var(--color-accent, #4493f8);
    color: #fff;
    font-weight: 600;
    font-size: 0.9375rem;
    text-decoration: none;
  }
  .cta:hover {
    filter: brightness(1.08);
  }
  .cta-hint {
    color: var(--color-muted);
    font-size: 0.8rem;
    font-family: var(--font-mono);
  }

  .capabilities h2 {
    font-size: 1.35rem;
    margin: 0 0 var(--space-3) 0;
  }
  .grid {
    display: grid;
    gap: var(--space-3);
    grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
  }
</style>
