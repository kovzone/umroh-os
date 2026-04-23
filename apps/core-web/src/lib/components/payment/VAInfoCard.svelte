<script lang="ts">
  import CountdownTimer from './CountdownTimer.svelte';
  import type { VirtualAccount } from '$lib/features/s2-payment/types';
  import { formatIdrAmountLabel } from '$lib/utils/format-idr';

  interface Props {
    va: VirtualAccount;
    /** Called when VA countdown reaches zero */
    onExpired?: () => void;
  }

  let { va, onExpired }: Props = $props();

  let copied = $state(false);
  let copyTimeout: ReturnType<typeof setTimeout> | undefined;

  const bankLabel = $derived(formatBankLabel(va.bank_code));
  const amountLabel = $derived(formatIdrAmountLabel(va.amount_total));

  function formatBankLabel(code: string): string {
    const map: Record<string, string> = {
      BCA: 'BCA (Bank Central Asia)',
      BNI: 'BNI (Bank Negara Indonesia)',
      BRI: 'BRI (Bank Rakyat Indonesia)',
      MANDIRI: 'Bank Mandiri',
      PERMATA: 'Bank Permata',
      CIMB: 'CIMB Niaga',
      BTN: 'BTN (Bank Tabungan Negara)'
    };
    return map[code.toUpperCase()] ?? code;
  }

  function formatAccountNumber(raw: string): string {
    // Group digits in blocks of 4 for readability: "8000 1234 5678 90"
    return raw.replace(/\D/g, '').replace(/(.{4})/g, '$1 ').trim();
  }

  async function copyAccountNumber() {
    const plain = va.account_number.replace(/\s/g, '');
    try {
      await navigator.clipboard.writeText(plain);
      copied = true;
      clearTimeout(copyTimeout);
      copyTimeout = setTimeout(() => {
        copied = false;
      }, 2000);
    } catch {
      // Clipboard API not available; silently fail
    }
  }
</script>

<div class="va-card">
  <div class="va-header">
    <div class="bank-badge">
      <span class="material-symbols-outlined">account_balance</span>
      <span>{bankLabel}</span>
    </div>
    <CountdownTimer expiresAt={va.expires_at} {onExpired} />
  </div>

  <div class="va-body">
    <div class="va-amount-row">
      <p class="label">Total yang harus dibayar</p>
      <p class="amount">{amountLabel}</p>
    </div>

    <div class="va-number-row">
      <p class="label">Nomor Virtual Account</p>
      <div class="va-number-wrap">
        <p class="va-number">{formatAccountNumber(va.account_number)}</p>
        <button
          type="button"
          class="copy-btn"
          class:copied
          onclick={() => void copyAccountNumber()}
          aria-label={copied ? 'Disalin!' : 'Salin nomor VA'}
          title={copied ? 'Disalin!' : 'Salin nomor VA'}
        >
          <span class="material-symbols-outlined">
            {copied ? 'check' : 'content_copy'}
          </span>
          <span class="copy-label">{copied ? 'Disalin!' : 'Salin'}</span>
        </button>
      </div>
    </div>
  </div>

  <div class="va-instructions">
    <h4 class="instr-title">
      <span class="material-symbols-outlined">info</span>
      Cara pembayaran via {va.bank_code}
    </h4>
    <ol class="steps">
      <li>Buka aplikasi mobile banking atau ATM <strong>{va.bank_code}</strong>.</li>
      <li>Pilih menu <strong>Transfer</strong> atau <strong>Virtual Account</strong>.</li>
      <li>Masukkan nomor VA: <strong>{formatAccountNumber(va.account_number)}</strong></li>
      <li>Pastikan nominal yang tertera: <strong>{amountLabel}</strong>.</li>
      <li>Konfirmasi dan selesaikan transaksi. Pembayaran terverifikasi otomatis.</li>
    </ol>
  </div>
</div>

<style>
  .va-card {
    background: #fff;
    border: 1px solid rgba(0, 103, 71, 0.2);
    border-radius: 1.25rem;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 77, 52, 0.06);
  }
  .va-header {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 1.25rem 1.5rem;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
    color: #fff;
  }
  .bank-badge {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    font-size: 0.95rem;
    font-weight: 700;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .bank-badge .material-symbols-outlined {
    font-size: 1.2rem;
  }
  /* Override countdown colors on dark header */
  .va-header :global(.countdown) {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.25);
    color: #fff;
  }
  .va-header :global(.countdown.urgent) {
    background: rgba(254, 212, 136, 0.35);
    border-color: rgba(254, 212, 136, 0.55);
    color: #ffd98a;
  }
  .va-header :global(.countdown.expired) {
    background: rgba(255, 100, 100, 0.2);
    border-color: rgba(255, 100, 100, 0.4);
    color: #ffb3b3;
  }

  .va-body {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    border-bottom: 1px solid rgba(190, 201, 193, 0.3);
  }
  .label {
    margin: 0 0 0.3rem;
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    color: #6f7a72;
  }
  .amount {
    margin: 0;
    font-size: 2rem;
    font-weight: 800;
    color: #775a19;
    font-family: 'Plus Jakarta Sans', sans-serif;
    letter-spacing: -0.02em;
  }
  .va-number-wrap {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  .va-number {
    margin: 0;
    font-size: 1.65rem;
    font-weight: 800;
    font-family: ui-monospace, 'Cascadia Code', 'Source Code Pro', monospace;
    letter-spacing: 0.06em;
    color: #004d34;
  }
  .copy-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.45rem 1rem;
    border-radius: 999px;
    border: 2px solid #004d34;
    background: transparent;
    color: #004d34;
    font-weight: 700;
    font-size: 0.82rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
    cursor: pointer;
    transition: background 0.15s ease, color 0.15s ease, border-color 0.15s ease;
  }
  .copy-btn:hover {
    background: rgba(0, 77, 52, 0.07);
  }
  .copy-btn.copied {
    background: #004d34;
    color: #fff;
    border-color: #004d34;
  }
  .copy-btn .material-symbols-outlined {
    font-size: 1rem;
  }
  .copy-label {
    font-size: 0.78rem;
  }

  .va-instructions {
    padding: 1.25rem 1.5rem;
    background: #fbfaf9;
  }
  .instr-title {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin: 0 0 0.85rem;
    font-size: 0.88rem;
    font-weight: 700;
    color: #3f4943;
  }
  .instr-title .material-symbols-outlined {
    font-size: 1rem;
    color: #006747;
  }
  .steps {
    margin: 0;
    padding-left: 1.15rem;
    font-size: 0.85rem;
    line-height: 1.6;
    color: #3f4943;
  }
  .steps li {
    margin-bottom: 0.4rem;
  }
  .steps li:last-child {
    margin-bottom: 0;
  }
</style>
