<script lang="ts">
  import type { FulfillmentStatus } from '$lib/features/s3-ops/types';

  interface Props {
    status: FulfillmentStatus;
  }

  let { status }: Props = $props();

  const labelMap: Record<FulfillmentStatus, string> = {
    queued: 'Queued',
    processing: 'Processing',
    shipped: 'Shipped',
    delivered: 'Delivered',
    cancelled: 'Cancelled'
  };
</script>

<span class="badge badge--{status}">{labelMap[status]}</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.15rem 0.45rem;
    border-radius: 0.2rem;
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    white-space: nowrap;
  }

  /* queued — abu-abu */
  .badge--queued {
    background: #e0e3e5;
    color: #434655;
  }

  /* processing — kuning */
  .badge--processing {
    background: #fff3cd;
    color: #7d4f00;
  }

  /* shipped — biru */
  .badge--shipped {
    background: rgb(37 99 235 / 0.12);
    color: #004ac6;
  }

  /* delivered — hijau */
  .badge--delivered {
    background: #d1fae5;
    color: #065f46;
  }

  /* cancelled — merah */
  .badge--cancelled {
    background: #ffdad6;
    color: #93000a;
  }
</style>
