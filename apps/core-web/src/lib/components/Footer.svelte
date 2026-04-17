<script lang="ts">
  // App-wide footer. Service-status link lives here per the scaffold design:
  // the landing page stays focused on product framing, while operators still
  // have a one-click route to the backend health dashboard.
</script>

<footer class="footer">
  <div class="inner">
    <nav>
      <a href="/system/status" data-testid="footer-status-link">Service status</a>
    </nav>
    <p class="meta">UmrohOS &middot; core-web &middot; scaffold v0.1</p>
  </div>
</footer>

<style>
  .footer {
    border-top: 1px solid var(--color-border);
    background: var(--color-surface);
    padding: var(--space-3) var(--space-4);
  }
  .inner {
    max-width: 72rem;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    flex-wrap: wrap;
  }
  nav a {
    color: var(--color-text);
    text-decoration: none;
    border-bottom: 1px dashed var(--color-border);
    padding-bottom: 1px;
  }
  nav a:hover { border-bottom-color: var(--color-text); }
  .meta {
    color: var(--color-muted);
    font-family: var(--font-mono);
    font-size: 0.8rem;
    margin: 0;
  }
</style>
