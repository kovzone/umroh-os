<script lang="ts">
  import { enhance } from '$app/forms';
  import type { SubmitFunction } from '@sveltejs/kit';
  import type { ActionData } from './$types';

  let { form }: { form: ActionData } = $props();
  let showPassword = $state(false);
  let submitting = $state(false);

  const formEnhance: SubmitFunction = () => {
    submitting = true;
    return async ({ update }) => {
      await update();
      submitting = false;
    };
  };
</script>

<svelte:head>
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@400;500;600;700&display=swap"
  />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
  />
</svelte:head>

<div class="shell">
  <main class="main">
    <div class="left-panel">
      <div class="panel-content">
        <div class="logo-row">
          <span class="material-symbols-outlined logo-icon">account_balance</span>
          <h1>UmrohOS Console</h1>
        </div>

        <div class="panel-copy">
          <h2>Internal Operations Console</h2>
          <p>Secure gateway for authorized personnel. Authenticate to manage the sovereign ledger.</p>

          <ul class="feature-list">
            <li>
              <span class="material-symbols-outlined">history</span>
              <div>
                <p class="feature-title">Audit trail</p>
                <p class="feature-desc">Immutable record of all operational transactions.</p>
              </div>
            </li>
            <li>
              <span class="material-symbols-outlined">admin_panel_settings</span>
              <div>
                <p class="feature-title">Role-based access</p>
                <p class="feature-desc">Strict authorization protocols per user clearance.</p>
              </div>
            </li>
            <li>
              <span class="material-symbols-outlined">sync_alt</span>
              <div>
                <p class="feature-title">Real-time sync</p>
                <p class="feature-desc">Immediate propagation across the distributed network.</p>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <p class="restricted">Restricted Access Environment</p>
    </div>

    <div class="right-panel">
      <div class="form-wrap">
        <div class="mobile-logo-row">
          <span class="material-symbols-outlined logo-icon">account_balance</span>
          <h1>UmrohOS Console</h1>
        </div>

        <header class="header">
          <h2>Console Login</h2>
          <p>Enter your credentials to access the operational system.</p>
        </header>

        <form method="POST" class="form-card" use:enhance={formEnhance}>
          {#if form?.error}
            <div class="error" role="alert">
              <span class="material-symbols-outlined">error</span>
              <div>
                <p class="error-title">Authentication Failed</p>
                <p class="error-msg">{form.error}</p>
              </div>
            </div>
          {/if}

          <label class="field">
            <span>Email</span>
            <input
              name="email"
              type="email"
              required
              autocomplete="email"
              placeholder="operator@umrohos.com"
              value={form?.values?.email ?? ''}
            />
          </label>

          <label class="field">
            <div class="password-label-row">
              <span>Password</span>
              <a href="/console/login#">Lupa Password?</a>
            </div>
            <div class="password-wrap">
              <input
                name="password"
                class="password-input"
                type={showPassword ? 'text' : 'password'}
                required
                autocomplete="current-password"
                placeholder="••••••••"
              />
              <button
                type="button"
                class="toggle"
                aria-label={showPassword ? 'Hide password' : 'Show password'}
                onclick={() => (showPassword = !showPassword)}
              >
                <span class="material-symbols-outlined">
                  {showPassword ? 'visibility' : 'visibility_off'}
                </span>
              </button>
            </div>
          </label>

          <label class="remember">
            <input name="remember-me" type="checkbox" />
            <span>Keep me logged in for this session</span>
          </label>

          <button class="submit" type="submit" disabled={submitting}>
            {#if submitting}
              <svg class="spinner" viewBox="0 0 24 24" aria-hidden="true">
                <circle class="spinner-track" cx="12" cy="12" r="10"></circle>
                <path
                  class="spinner-fill"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              <span>Authenticating...</span>
            {:else}
              <span>Sign In</span>
            {/if}
          </button>
        </form>

        <div class="secure-line">
          <span class="material-symbols-outlined">lock</span>
          <p>Encrypted Operational Line</p>
        </div>
      </div>
    </div>
  </main>

  <footer class="footer">
    <div class="footer-left">
      <span class="brand">UmrohOS</span>
      <span class="divider">|</span>
      <span class="copy">© 2024 Sovereign Ledger Systems</span>
    </div>
    <nav class="footer-nav">
      <a href="/console/login#">Butuh bantuan IT?</a>
      <a href="/console/login#status" class="status-link"><span></span>Status sistem</a>
    </nav>
  </footer>
</div>

<style>
  :global(.material-symbols-outlined) {
    font-family: 'Material Symbols Outlined', sans-serif;
    font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }

  .shell {
    min-height: 100vh;
    min-height: 100dvh;
    display: flex;
    flex-direction: column;
    background: #f8fafc;
    color: #0f172a;
    font-family:
      'IBM Plex Sans',
      ui-sans-serif,
      system-ui,
      -apple-system,
      sans-serif;
  }

  .main {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .left-panel {
    display: none;
  }

  .right-panel {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1.5rem;
  }

  .form-wrap {
    width: 100%;
    max-width: 30rem;
  }

  .mobile-logo-row {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.7rem;
    margin-bottom: 2.25rem;
  }

  .logo-icon {
    color: #1d4ed8;
    font-size: 2rem;
    font-variation-settings: 'FILL' 1, 'wght' 500, 'GRAD' 0, 'opsz' 24;
  }

  .mobile-logo-row h1 {
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.02em;
    font-size: 1.1rem;
    font-weight: 700;
  }

  .header {
    margin-bottom: 1.5rem;
    text-align: center;
  }

  .header h2 {
    margin: 0 0 0.25rem;
    font-size: 1.6rem;
    line-height: 1.2;
    color: #0f172a;
  }

  .header p {
    margin: 0;
    font-size: 0.875rem;
    color: #64748b;
  }

  .form-card {
    background: #fff;
    border: 1px solid #e2e8f0;
    padding: 2rem;
    display: grid;
    gap: 1.15rem;
  }

  .error {
    background: #fef2f2;
    border: 1px solid #fecaca;
    display: flex;
    gap: 0.6rem;
    padding: 0.8rem;
    align-items: flex-start;
  }

  .error .material-symbols-outlined {
    color: #dc2626;
    font-size: 1.1rem;
    margin-top: 0.1rem;
  }

  .error-title {
    margin: 0;
    font-size: 0.88rem;
    font-weight: 700;
    color: #b91c1c;
  }

  .error-msg {
    margin: 0.15rem 0 0;
    font-size: 0.86rem;
    color: #dc2626;
  }

  .field {
    display: grid;
    gap: 0.45rem;
  }

  .field span {
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.12em;
    color: #475569;
  }

  .field input {
    border: 1px solid #cbd5e1;
    width: 100%;
    padding: 0.7rem 0.85rem;
    font: inherit;
    font-size: 0.9rem;
    color: #0f172a;
    background: #fff;
    border-radius: 0;
    transition: border-color 120ms ease;
  }

  .field input:focus {
    outline: none;
    border-color: #1d4ed8;
    box-shadow: 0 0 0 1px #1d4ed8;
  }

  .field input::placeholder {
    color: #94a3b8;
  }

  .password-label-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .password-label-row a {
    font-size: 0.72rem;
    color: #1d4ed8;
    font-weight: 600;
    text-decoration: none;
  }

  .password-label-row a:hover {
    text-decoration: underline;
  }

  .password-wrap {
    position: relative;
  }

  .toggle {
    position: absolute;
    right: 0.55rem;
    top: 50%;
    transform: translateY(-50%);
    border: none;
    background: transparent;
    color: #64748b;
    padding: 0.1rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .toggle:hover {
    color: #475569;
  }

  .toggle .material-symbols-outlined {
    font-size: 1.2rem;
  }

  .password-input {
    padding-right: 2.6rem;
  }

  .remember {
    display: flex;
    align-items: center;
    gap: 0.55rem;
    color: #475569;
    font-size: 0.86rem;
  }

  .remember input {
    width: 1rem;
    height: 1rem;
    margin: 0;
    accent-color: #1d4ed8;
  }

  .submit {
    margin-top: 0.2rem;
    border: 1px solid #1d4ed8;
    background: #1d4ed8;
    color: #fff;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    padding: 0.78rem 1rem;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.6rem;
    min-height: 2.85rem;
    transition: background-color 140ms ease;
  }

  .submit:hover:not(:disabled) {
    background: #1e40af;
    border-color: #1e40af;
  }

  .submit:disabled {
    opacity: 0.92;
    cursor: not-allowed;
  }

  .spinner {
    width: 1rem;
    height: 1rem;
    animation: spin 0.8s linear infinite;
  }

  .spinner-track {
    fill: none;
    stroke: currentColor;
    stroke-width: 4;
    opacity: 0.25;
  }

  .spinner-fill {
    fill: currentColor;
    opacity: 0.75;
  }

  .secure-line {
    margin-top: 1.5rem;
    display: flex;
    gap: 0.35rem;
    justify-content: center;
    align-items: center;
    color: #94a3b8;
  }

  .secure-line .material-symbols-outlined {
    font-size: 0.95rem;
  }

  .secure-line p {
    margin: 0;
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.14em;
  }

  .footer {
    border-top: 1px solid #e2e8f0;
    background: #fff;
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.5rem;
  }

  .footer-left {
    display: flex;
    align-items: center;
    gap: 0.65rem;
    flex-wrap: wrap;
    justify-content: center;
  }

  .brand {
    font-size: 0.73rem;
    font-weight: 700;
    text-transform: uppercase;
  }

  .divider {
    color: #cbd5e1;
  }

  .copy {
    font-size: 0.62rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #64748b;
    font-weight: 500;
  }

  .footer-nav {
    display: flex;
    gap: 1.2rem;
    align-items: center;
    flex-wrap: wrap;
    justify-content: center;
  }

  .footer-nav a {
    text-decoration: none;
    text-transform: uppercase;
    letter-spacing: 0.12em;
    font-size: 0.62rem;
    font-weight: 700;
    color: #64748b;
  }

  .footer-nav a:hover {
    color: #1d4ed8;
  }

  .status-link {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
  }

  .status-link span {
    width: 0.45rem;
    height: 0.45rem;
    border-radius: 999px;
    background: #10b981;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  @media (min-width: 900px) {
    .main {
      flex-direction: row;
    }

    .left-panel {
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      background: #0f172a;
      color: #cbd5e1;
      width: 50%;
      border-right: 1px solid #1e293b;
      padding: 3rem 2.6rem 2.2rem;
    }

    .panel-content {
      display: flex;
      flex-direction: column;
      gap: 2.6rem;
    }

    .logo-row {
      display: flex;
      align-items: center;
      gap: 0.65rem;
    }

    .logo-row h1 {
      margin: 0;
      font-size: 1.15rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.02em;
      color: #fff;
    }

    .panel-copy h2 {
      margin: 0 0 0.6rem;
      color: #fff;
      font-size: 2rem;
      line-height: 1.2;
    }

    .panel-copy > p {
      margin: 0 0 2rem;
      color: #94a3b8;
      max-width: 31rem;
      line-height: 1.55;
    }

    .feature-list {
      margin: 0;
      padding: 0;
      list-style: none;
      display: grid;
      gap: 1.2rem;
    }

    .feature-list li {
      display: flex;
      gap: 0.8rem;
      align-items: flex-start;
    }

    .feature-list li .material-symbols-outlined {
      color: #3b82f6;
      background: #1e293b;
      border: 1px solid #334155;
      padding: 0.42rem;
      font-size: 1.15rem;
      flex-shrink: 0;
    }

    .feature-title {
      margin: 0;
      color: #f1f5f9;
      font-weight: 600;
      font-size: 0.95rem;
    }

    .feature-desc {
      margin: 0.12rem 0 0;
      color: #94a3b8;
      font-size: 0.85rem;
      line-height: 1.45;
    }

    .restricted {
      margin: 0;
      display: inline-flex;
      align-items: center;
      width: fit-content;
      background: #1e293b;
      border: 1px solid #334155;
      color: #94a3b8;
      padding: 0.45rem 0.65rem;
      font-size: 0.58rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.2em;
    }

    .right-panel {
      width: 50%;
      padding: 2.25rem 2rem;
    }

    .mobile-logo-row {
      display: none;
    }

    .header {
      text-align: left;
      margin-bottom: 1.15rem;
    }

    .footer {
      flex-direction: row;
      padding: 1rem 2rem;
    }

    .footer-left,
    .footer-nav {
      justify-content: flex-start;
    }
  }

  @media (min-width: 1280px) {
    .left-panel {
      width: 40%;
    }

    .right-panel {
      width: 60%;
    }

    .left-panel {
      padding: 3.7rem 3.5rem 2.4rem;
    }

    .right-panel {
      padding: 2rem 3rem;
    }

    .form-wrap {
      max-width: 32rem;
    }

    .header h2 {
      font-size: 2rem;
    }
  }
</style>

