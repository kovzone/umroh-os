<script lang="ts">
  const steps = [
    { n: 1, label: 'Pilih keberangkatan' },
    { n: 2, label: 'Data jamaah' },
    { n: 3, label: 'Review & syarat' },
    { n: 4, label: 'Pembayaran' }
  ] as const;

  let { current }: { current: number } = $props();

  const safeStep = $derived(current >= 1 && current <= 4 ? current : 1);
  const progressPct = $derived(((safeStep - 1) / 3) * 100);
</script>

<div class="stepper-root">
  <div class="stepper-row">
    {#each steps as s (s.n)}
      <div class="step-col" class:active={safeStep === s.n} class:done={safeStep > s.n}>
        <div class="circle">
          {#if safeStep > s.n}
            <span class="material-symbols-outlined fill">check</span>
          {:else}
            <span class="num">{s.n}</span>
          {/if}
        </div>
        <span class="lbl">{s.label}</span>
      </div>
    {/each}
  </div>
  <div class="track">
    <div class="track-bg"></div>
    <div class="track-fill" style:width="{progressPct}%"></div>
  </div>
</div>

<style>
  .stepper-root {
    position: relative;
    margin-bottom: 2.5rem;
  }
  .stepper-row {
    position: relative;
    z-index: 2;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    max-width: 48rem;
    margin: 0 auto;
    gap: 0.5rem;
  }
  .step-col {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.65rem;
    text-align: center;
    min-width: 0;
  }
  .circle {
    width: 2.75rem;
    height: 2.75rem;
    border-radius: 999px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #efeded;
    color: #6f7a72;
    font-weight: 700;
    font-size: 0.95rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
    box-shadow: none;
  }
  .step-col.done .circle {
    background: #006747;
    color: #fff;
  }
  .step-col.active .circle {
    background: #006747;
    color: #fff;
    box-shadow: 0 0 0 4px rgba(0, 103, 71, 0.12);
  }
  .step-col.active:not(.done) .circle {
    background: linear-gradient(135deg, #006747, #004d34);
  }
  .num {
    line-height: 1;
  }
  .lbl {
    font-size: 0.68rem;
    font-weight: 600;
    color: #6f7a72;
    line-height: 1.2;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  @media (min-width: 640px) {
    .lbl {
      font-size: 0.78rem;
    }
  }
  .step-col.active .lbl {
    color: #004d34;
    font-weight: 700;
  }
  .step-col.done .lbl {
    color: #3f4943;
  }
  .track {
    position: absolute;
    top: 1.35rem;
    left: 8%;
    right: 8%;
    height: 3px;
    z-index: 1;
    pointer-events: none;
  }
  .track-bg {
    position: absolute;
    inset: 0;
    background: #efeded;
    border-radius: 999px;
  }
  .track-fill {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    border-radius: 999px;
    background: linear-gradient(90deg, #775a19, #e9c176);
    transition: width 0.35s ease;
  }
</style>
