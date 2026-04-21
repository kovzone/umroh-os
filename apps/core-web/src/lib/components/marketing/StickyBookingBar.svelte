<script lang="ts">
  let {
    priceLabel,
    dateLabel,
    bookingHref,
    whatsappHref,
    canBook,
    bookingTestId = undefined
  }: {
    priceLabel: string;
    dateLabel: string | null;
    bookingHref: string;
    whatsappHref: string;
    canBook: boolean;
    bookingTestId?: string;
  } = $props();
</script>

<div class="sticky-cta" role="region" aria-label="Ringkasan harga dan booking">
  <div class="shell sticky-inner">
    <div class="sticky-meta">
      <div>
        <p class="sticky-label">Perkiraan harga</p>
        <p class="sticky-price">{priceLabel}</p>
      </div>
      {#if dateLabel}
        <div class="sticky-divider"></div>
        <div>
          <p class="sticky-label">Tanggal terpilih</p>
          <p class="sticky-date">{dateLabel}</p>
        </div>
      {/if}
    </div>
    <div class="sticky-actions">
      <a class="btn-outline sm" href={whatsappHref} target="_blank" rel="noreferrer">Chat WhatsApp</a>
      <a
        class="btn-gradient sm {!canBook ? 'is-disabled' : ''}"
        href={bookingHref}
        data-testid={bookingTestId}>Lanjut booking</a>
    </div>
  </div>
</div>

<style>
  .sticky-cta {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 45;
    background: rgba(255, 255, 255, 0.88);
    backdrop-filter: blur(12px);
    border-top: 1px solid rgba(190, 201, 193, 0.35);
    box-shadow: 0 -8px 30px rgb(0, 0, 0, 0.06);
    padding: 0.75rem 0;
  }
  .sticky-inner {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
  }
  .sticky-meta {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    flex-wrap: wrap;
  }
  .sticky-label {
    margin: 0;
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #6f7a72;
  }
  .sticky-price {
    margin: 0;
    font-size: 1.35rem;
    font-weight: 800;
    color: #004d34;
  }
  .sticky-date {
    margin: 0;
    font-weight: 700;
    font-size: 0.95rem;
  }
  .sticky-divider {
    width: 1px;
    height: 2.25rem;
    background: rgba(190, 201, 193, 0.5);
    display: none;
  }
  @media (min-width: 640px) {
    .sticky-divider {
      display: block;
    }
  }
  .sticky-actions {
    display: flex;
    align-items: center;
    gap: 0.65rem;
    flex-wrap: wrap;
  }
  .btn-gradient {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    text-decoration: none;
    color: #fff;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
    font-weight: 700;
    box-shadow: 0 12px 28px rgba(0, 77, 52, 0.22);
  }
  .btn-gradient.is-disabled {
    opacity: 0.45;
    pointer-events: none;
    cursor: not-allowed;
  }
  .btn-outline {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    text-decoration: none;
    color: #1b1c1c;
    border: 2px solid rgba(190, 201, 193, 0.45);
    background: transparent;
    font-weight: 700;
  }
  .btn-outline:hover {
    background: #f5f3f3;
  }
  .btn-outline.sm,
  .btn-gradient.sm {
    padding: 0.65rem 1.1rem;
    font-size: 0.9rem;
    border-radius: 0.85rem;
  }
</style>
