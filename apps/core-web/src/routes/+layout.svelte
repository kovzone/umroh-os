<script lang="ts">
  import '../app.css';
  import { page } from '$app/state';
  import Header from '$lib/components/Header.svelte';
  import Footer from '$lib/components/Footer.svelte';

  let { children } = $props();
  const useCustomMarketingChrome = $derived(
    page.url.pathname === '/' ||
      page.url.pathname === '/packages' ||
      /^\/packages\/[^/]+$/.test(page.url.pathname) ||
      /^\/booking\/[^/]+$/.test(page.url.pathname) ||
      /^\/console(\/.*)?$/.test(page.url.pathname)
  );
</script>

<div class="app" class:home-mode={useCustomMarketingChrome}>
  {#if !useCustomMarketingChrome}
    <Header />
  {/if}
  <main class="container">
    {@render children()}
  </main>
  {#if !useCustomMarketingChrome}
    <Footer />
  {/if}
</div>

<style>
  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }
  .app.home-mode {
    background: #f8f7f4;
    color: #1b1c1c;
  }
  .container {
    max-width: 72rem;
    width: 100%;
    margin: 0 auto;
    padding: var(--space-4);
    flex: 1;
  }
  .app.home-mode .container {
    max-width: none;
    padding: 0;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }
</style>
