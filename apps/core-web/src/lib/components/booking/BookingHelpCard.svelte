<script lang="ts">
  let { whatsappHref }: { whatsappHref: string } = $props();
</script>

<div class="help">
  <div class="help-inner">
    <p class="t1">Butuh bantuan?</p>
    <p class="t2">Konsultan kami siap membantu proses pemesanan Anda.</p>
    <a class="btn" href={whatsappHref} target="_blank" rel="noopener noreferrer">Hubungi CS</a>
  </div>
  <span class="material-symbols-outlined deco" aria-hidden="true">support_agent</span>
</div>

<style>
  .help {
    position: relative;
    overflow: hidden;
    border-radius: 1rem;
    padding: 1.5rem 1.5rem 1.75rem;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
    color: #fff;
  }
  .help-inner {
    position: relative;
    z-index: 1;
  }
  .t1 {
    margin: 0 0 0.25rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
    font-weight: 700;
    font-size: 1.05rem;
  }
  .t2 {
    margin: 0 0 1rem;
    font-size: 0.75rem;
    line-height: 1.45;
    opacity: 0.88;
    max-width: 18rem;
  }
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem 1rem;
    border-radius: 999px;
    font-size: 0.72rem;
    font-weight: 700;
    font-family: 'Plus Jakarta Sans', sans-serif;
    text-decoration: none;
    color: #004d34;
    background: rgba(255, 255, 255, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.35);
    transition: background 0.15s ease, transform 0.15s ease;
  }
  .btn:hover {
    background: #a0f4ca;
  }
  .deco {
    position: absolute;
    right: -0.5rem;
    bottom: -0.5rem;
    font-size: 6rem;
    opacity: 0.08;
    pointer-events: none;
    transition: transform 0.25s ease;
  }
  .help:hover .deco {
    transform: rotate(8deg) scale(1.02);
  }
</style>
