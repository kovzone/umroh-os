<script lang="ts">
  import { onDestroy } from 'svelte';

  interface Props {
    /** ISO 8601 expiry timestamp, e.g. "2026-04-24T12:00:00Z" */
    expiresAt: string;
    /** Called when countdown reaches zero */
    onExpired?: () => void;
  }

  let { expiresAt, onExpired }: Props = $props();

  function computeSecondsLeft(expires: string): number {
    return Math.max(0, Math.floor((new Date(expires).getTime() - Date.now()) / 1000));
  }

  let secondsLeft = $state(computeSecondsLeft(expiresAt));

  const hours = $derived(Math.floor(secondsLeft / 3600));
  const minutes = $derived(Math.floor((secondsLeft % 3600) / 60));
  const seconds = $derived(secondsLeft % 60);

  const isUrgent = $derived(secondsLeft > 0 && secondsLeft <= 600); // <= 10 minutes
  const isExpired = $derived(secondsLeft === 0);

  function pad(n: number): string {
    return n.toString().padStart(2, '0');
  }

  let interval: ReturnType<typeof setInterval> | undefined;

  $effect(() => {
    // Re-sync when expiresAt prop changes
    secondsLeft = computeSecondsLeft(expiresAt);

    clearInterval(interval);

    if (secondsLeft > 0) {
      interval = setInterval(() => {
        secondsLeft = computeSecondsLeft(expiresAt);
        if (secondsLeft === 0) {
          clearInterval(interval);
          onExpired?.();
        }
      }, 1000);
    } else {
      onExpired?.();
    }

    return () => {
      clearInterval(interval);
    };
  });

  onDestroy(() => {
    clearInterval(interval);
  });
</script>

<div class="countdown" class:urgent={isUrgent} class:expired={isExpired} aria-live="polite" aria-label="Batas waktu pembayaran">
  {#if isExpired}
    <span class="material-symbols-outlined ic">timer_off</span>
    <span class="label expired-text">VA telah kadaluarsa</span>
  {:else}
    <span class="material-symbols-outlined ic">schedule</span>
    <div class="segments">
      <div class="seg">
        <span class="digit">{pad(hours)}</span>
        <span class="unit">jam</span>
      </div>
      <span class="sep">:</span>
      <div class="seg">
        <span class="digit">{pad(minutes)}</span>
        <span class="unit">mnt</span>
      </div>
      <span class="sep">:</span>
      <div class="seg">
        <span class="digit">{pad(seconds)}</span>
        <span class="unit">dtk</span>
      </div>
    </div>
  {/if}
</div>

<style>
  .countdown {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.55rem 1rem;
    border-radius: 0.75rem;
    background: rgba(254, 212, 136, 0.25);
    border: 1px solid rgba(119, 90, 25, 0.2);
    color: #775a19;
    font-family: 'Plus Jakarta Sans', sans-serif;
    font-size: 0.88rem;
    font-weight: 700;
    transition: background 0.3s ease, border-color 0.3s ease;
  }
  .countdown.urgent {
    background: rgba(254, 212, 136, 0.55);
    border-color: rgba(119, 90, 25, 0.45);
    animation: urgentPulse 1s ease-in-out infinite;
  }
  .countdown.expired {
    background: rgba(147, 0, 10, 0.08);
    border-color: rgba(147, 0, 10, 0.2);
    color: #93000a;
  }
  @keyframes urgentPulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.75; }
  }
  .ic {
    font-size: 1.1rem;
  }
  .label {
    font-size: 0.82rem;
  }
  .expired-text {
    font-size: 0.82rem;
  }
  .segments {
    display: flex;
    align-items: center;
    gap: 0.2rem;
  }
  .seg {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 2.5rem;
  }
  .digit {
    font-size: 1.25rem;
    font-weight: 800;
    line-height: 1;
    font-variant-numeric: tabular-nums;
    letter-spacing: 0.02em;
  }
  .unit {
    font-size: 0.52rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    opacity: 0.65;
  }
  .sep {
    font-size: 1.1rem;
    font-weight: 800;
    margin-bottom: 0.4rem;
    opacity: 0.6;
  }
</style>
