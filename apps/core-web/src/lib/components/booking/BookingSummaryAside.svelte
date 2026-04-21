<script lang="ts">
  import type { PackageDetail } from '$lib/features/s1-catalog/types';

  let {
    pkg,
    durationLabel,
    roomLabel,
    totalLabel,
    departureDayLabel,
    jamaahCount,
    step,
    finePrint,
    breakdownHint
  }: {
    pkg: PackageDetail;
    durationLabel: string;
    roomLabel: string;
    totalLabel: string;
    departureDayLabel: string;
    jamaahCount: number;
    step: number;
    finePrint?: string;
    /** e.g. "2 jamaah × Rp 38.500.000" */
    breakdownHint?: string;
  } = $props();

  const foot = $derived(finePrint ?? pkg.priceFinePrint ?? '');
</script>

<div class="aside-card">
  <div class="accent"></div>
  <h3 class="title">
    <span class="material-symbols-outlined ico">analytics</span>
    Ringkasan paket
  </h3>
  {#if step >= 2}
    <div class="hero-row">
      <div class="thumb">
        <img src={pkg.coverPhotoUrl} alt="" role="presentation" />
      </div>
      <div>
        <p class="pkg-name">{pkg.name}</p>
        <p class="sub">{durationLabel}</p>
      </div>
    </div>
  {/if}
  <div class="rows">
    {#if step === 1}
      <div class="row">
        <span class="k">Nama paket</span>
        <span class="v">{pkg.name}</span>
      </div>
      <div class="row">
        <span class="k">Durasi</span>
        <span class="v">{durationLabel}</span>
      </div>
      <div class="row">
        <span class="k">Kamar</span>
        <span class="v">{roomLabel}</span>
      </div>
    {:else}
      <div class="row">
        <span class="k">Tipe kamar</span>
        <span class="v">{roomLabel}</span>
      </div>
      <div class="row">
        <span class="k">Tanggal</span>
        <span class="v">{departureDayLabel}</span>
      </div>
      <div class="row">
        <span class="k">Jumlah jamaah</span>
        <span class="v">{jamaahCount} jamaah</span>
      </div>
    {/if}
  </div>
  <div class="total-block">
    <div class="total-row">
      <span class="total-lbl">Total biaya</span>
      <span class="total-amt">{totalLabel}</span>
    </div>
    {#if breakdownHint}
      <p class="breakdown">{breakdownHint}</p>
    {/if}
    {#if foot}
      <p class="disclaimer">{foot}</p>
    {/if}
  </div>
</div>

<style>
  .aside-card {
    position: relative;
    background: #fff;
    border-radius: 1.25rem;
    padding: 1.35rem 1.35rem 1.5rem;
    box-shadow: 0 8px 32px rgba(27, 28, 28, 0.04);
    overflow: hidden;
  }
  .accent {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: #775a19;
    border-radius: 0 2px 2px 0;
  }
  .title {
    margin: 0 0 1.1rem;
    padding-left: 0.35rem;
    font-size: 1.05rem;
    font-weight: 700;
    color: #004d34;
    display: flex;
    align-items: center;
    gap: 0.45rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .ico {
    font-size: 1.35rem;
    color: #775a19;
  }
  .hero-row {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
    margin-bottom: 1rem;
    padding: 0.35rem 0 0 0.35rem;
  }
  .thumb {
    width: 4rem;
    height: 4rem;
    border-radius: 0.5rem;
    overflow: hidden;
    flex-shrink: 0;
  }
  .thumb img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .pkg-name {
    margin: 0 0 0.2rem;
    font-weight: 700;
    font-size: 0.92rem;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
    line-height: 1.25;
  }
  .sub {
    margin: 0;
    font-size: 0.72rem;
    color: #6f7a72;
  }
  .rows {
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
    padding-left: 0.35rem;
  }
  .row {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 0.75rem;
    font-size: 0.88rem;
  }
  .k {
    color: #6f7a72;
  }
  .v {
    font-weight: 700;
    text-align: right;
  }
  .total-block {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid #efeded;
    padding-left: 0.35rem;
  }
  .total-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    gap: 0.5rem;
    margin-bottom: 0.35rem;
  }
  .total-lbl {
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #6f7a72;
  }
  .total-amt {
    font-size: 1.35rem;
    font-weight: 800;
    color: #775a19;
    font-family: 'Plus Jakarta Sans', sans-serif;
    letter-spacing: -0.02em;
  }
  .breakdown {
    margin: 0 0 0.5rem;
    font-size: 0.68rem;
    color: #6f7a72;
    line-height: 1.35;
  }
  .disclaimer {
    margin: 0;
    font-size: 0.62rem;
    line-height: 1.45;
    color: #6f7a72;
  }
</style>
