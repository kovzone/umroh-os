<script lang="ts">
  let {
    ctaHref,
    ctaLabel = 'Daftar Sekarang',
    ctaTestId = undefined,
    packagesLinkActive = true
  }: {
    ctaHref: string;
    ctaLabel?: string;
    ctaTestId?: string;
    /** When true, "Paket Umroh" shows active styling (catalog + package detail). */
    packagesLinkActive?: boolean;
  } = $props();
</script>

<nav class="top-nav" aria-label="Utama">
  <div class="shell nav-inner">
    <a class="brand" href="/">UmrohOS</a>
    <div class="nav-links">
      <a href="/packages" class:active={packagesLinkActive}>Paket Umroh</a>
      <a href="/packages">Jadwal</a>
      <a href="/#proses-booking">Manasik</a>
      <a href="/">Tentang Kami</a>
    </div>
    <a class="nav-cta" href={ctaHref} data-testid={ctaTestId}>{ctaLabel}</a>
  </div>
</nav>

<style>
  .top-nav {
    position: fixed;
    top: 0;
    width: 100%;
    z-index: 50;
    background: rgba(251, 249, 248, 0.85);
    backdrop-filter: blur(12px);
    box-shadow: 0 16px 38px -12px rgba(27, 28, 28, 0.06);
  }
  .nav-inner {
    height: 5.2rem;
    display: grid;
    grid-template-columns: auto 1fr auto;
    align-items: center;
    gap: 2rem;
  }
  .brand {
    text-decoration: none;
    color: #006747;
    font-weight: 700;
    font-size: 1.65rem;
    letter-spacing: -0.02em;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .nav-links {
    justify-self: center;
    display: flex;
    align-items: center;
    gap: 2rem;
  }
  .nav-links a {
    color: #57534e;
    text-decoration: none;
    font-family: 'Plus Jakarta Sans', sans-serif;
    font-weight: 500;
  }
  .nav-links a.active {
    color: #006747;
    font-weight: 700;
    border-bottom: 2px solid #775a19;
    padding-bottom: 0.2rem;
  }
  .nav-cta {
    border: none;
    text-decoration: none;
    border-radius: 999px;
    background: linear-gradient(90deg, #004d34, #006747);
    color: #fff;
    font-weight: 700;
    padding: 0.7rem 1.45rem;
    font-size: 0.9rem;
    white-space: nowrap;
  }
  @media (max-width: 1100px) {
    .nav-links {
      display: none;
    }
  }
  @media (max-width: 760px) {
    .nav-inner {
      grid-template-columns: 1fr auto;
    }
  }
</style>
