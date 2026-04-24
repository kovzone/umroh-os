<script lang="ts">
  import type { Snippet } from 'svelte';
  import FloatingChat from './FloatingChat.svelte';
  import MarketingFonts from './MarketingFonts.svelte';
  import MarketingFooter from './MarketingFooter.svelte';
  import MarketingTopNav from './MarketingTopNav.svelte';

  let {
    ctaHref,
    ctaLabel = 'Daftar Sekarang',
    ctaTestId = undefined,
    packagesLinkActive = true,
    children
  }: {
    ctaHref: string;
    ctaLabel?: string;
    ctaTestId?: string;
    packagesLinkActive?: boolean;
    children: Snippet;
  } = $props();
</script>

<MarketingFonts />

<div class="m-page">
  <MarketingTopNav {ctaHref} {ctaLabel} {ctaTestId} {packagesLinkActive} />
  <!-- Reserves space for fixed .top-nav (see MarketingTopNav .nav-inner height) -->
  <div class="top-nav-spacer" aria-hidden="true"></div>
  <div class="m-page-body">
    {@render children()}
  </div>
  <div class="m-page-footer">
    <MarketingFooter />
  </div>
  <FloatingChat />
</div>

<style>
  .m-page {
    font-family: 'Manrope', sans-serif;
    background: #fbf9f8;
    color: #1b1c1c;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    min-height: 100dvh;
  }
  .top-nav-spacer {
    height: 5.2rem;
    flex-shrink: 0;
  }
  .m-page-body {
    flex: 1;
    width: 100%;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }
  .m-page-footer {
    flex-shrink: 0;
    margin-top: auto;
  }
  .m-page :global(h1),
  .m-page :global(h2),
  .m-page :global(h3),
  .m-page :global(h4) {
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .m-page :global(.material-symbols-outlined) {
    font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
    font-size: 1.2rem;
    line-height: 1;
  }
  .m-page :global(.material-symbols-outlined.fill) {
    font-variation-settings: 'FILL' 1, 'wght' 500, 'GRAD' 0, 'opsz' 24;
  }
  .m-page :global(.shell) {
    max-width: 80rem;
    margin: 0 auto;
    padding: 0 2rem;
  }
  @media (max-width: 760px) {
    .m-page :global(.shell) {
      padding: 0 1rem;
    }
  }
</style>
