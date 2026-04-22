<script lang="ts">
  import { page } from '$app/state';

  const slugParts = $derived(page.params.slug?.split('/') ?? []);
  const title = $derived(
    slugParts.length > 0
      ? slugParts
          .map((part) => part.replace(/-/g, ' '))
          .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
          .join(' / ')
      : 'Placeholder'
  );
</script>

<main class="placeholder-page">
  <header>
    <h2>{title}</h2>
    <p>Halaman ini masih placeholder dan akan diimplementasikan pada task berikutnya.</p>
  </header>

  <section class="placeholder-card">
    <p>Route: <code>{page.url.pathname}</code></p>
  </section>
</main>

<style>
  .placeholder-page {
    padding: 1.5rem;
    max-width: 96rem;
  }

  header h2 {
    margin: 0;
    font-size: 1.45rem;
    line-height: 1.2;
    color: #191c1e;
  }

  header p {
    margin: 0.35rem 0 0;
    font-size: 0.86rem;
    color: #434655;
  }

  .placeholder-card {
    margin-top: 0.9rem;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #fff;
    border-radius: 0.25rem;
    padding: 0.8rem;
    color: #434655;
    font-size: 0.84rem;
  }

  code {
    color: #004ac6;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 0.8rem;
  }
</style>

