<script lang="ts">
  // App-wide header. Lives inside +layout.svelte so every route inherits it.
  // The Sign-in button is intentionally disabled until F1.5 wires real iam-svc auth;
  // the tooltip flags that explicitly so future devs don't chase a bug.
</script>

<header class="topbar">
  <div class="brand">
    <a href="/" class="brand-link">
      <strong>UmrohOS</strong>
    </a>
    <span class="muted">core-web</span>
  </div>

  <nav>
    <button
      type="button"
      disabled
      title="Sign-in lands with F1.5 — iam-svc auth"
      class="btn-signin"
      data-testid="signin-button"
    >
      Sign in
    </button>
  </nav>
</header>

<style>
  .topbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
    background: var(--color-surface);
  }
  .brand { display: flex; align-items: baseline; gap: var(--space-3); }
  .brand-link { color: var(--color-text); text-decoration: none; }
  .muted { color: var(--color-muted); font-family: var(--font-mono); font-size: 0.875rem; }

  .btn-signin {
    background: transparent;
    color: var(--color-muted);
    border: 1px solid var(--color-border);
    padding: 0.4rem 1rem;
    border-radius: var(--radius-card);
    cursor: not-allowed;
    font-family: var(--font-sans);
    font-size: 0.875rem;
  }
  .btn-signin:disabled { opacity: 0.5; }
</style>
