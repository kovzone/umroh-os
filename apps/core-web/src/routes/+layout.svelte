<script lang="ts">
  import '../app.css';
  import Header from '$lib/components/Header.svelte';
  import Footer from '$lib/components/Footer.svelte';

  let { children } = $props();
</script>

<div class="app">
  <Header />
  <main class="container">
    {@render children()}
  </main>
  <Footer />
</div>

<style>
  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }
  .container {
    max-width: 72rem;
    width: 100%;
    margin: 0 auto;
    padding: var(--space-4);
    flex: 1;
  }
</style>
