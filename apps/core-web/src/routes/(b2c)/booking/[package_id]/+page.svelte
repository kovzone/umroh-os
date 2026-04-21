<script lang="ts">
  import { createDraftBooking } from '$lib/features/s1-booking/repository';
  import type { DraftBookingError } from '$lib/features/s1-booking/types';
  import PageShell from '$lib/components/s1/PageShell.svelte';

  let { data } = $props();

  let leadName = $state('');
  let leadWhatsapp = $state('');
  let leadEmail = $state('');
  let leadDomicile = $state('');
  let jamaahCount = $state(1);
  let roomType = $state<'quad' | 'triple' | 'double'>('quad');
  const departures = $derived(data.package.departures);
  let departureId = $state('');

  let isSubmitting = $state(false);
  let submitError = $state('');
  let submitSuccess = $state<{ bookingId: string; createdAt: string } | null>(null);

  $effect(() => {
    if (!departureId && departures.length > 0) {
      departureId = departures[0].id;
    }
  });

  function isFormValid(): boolean {
    return (
      leadName.trim().length > 0 &&
      leadWhatsapp.trim().length > 0 &&
      leadEmail.trim().length > 0 &&
      leadDomicile.trim().length > 0 &&
      jamaahCount >= 1 &&
      departureId.length > 0
    );
  }

  async function onSubmitDraft() {
    submitError = '';
    submitSuccess = null;

    if (!isFormValid()) {
      submitError = 'Lengkapi semua field wajib sebelum simpan draft.';
      return;
    }

    isSubmitting = true;
    try {
      const result = await createDraftBooking({
        channel: 'b2c_self',
        packageId: data.package.id,
        departureId,
        roomType,
        jamaahCount,
        lead: {
          fullName: leadName,
          whatsapp: leadWhatsapp,
          email: leadEmail,
          domicile: leadDomicile
        }
      });
      submitSuccess = { bookingId: result.bookingId, createdAt: result.createdAt };
    } catch (err) {
      const bookingErr = err as DraftBookingError;
      submitError = `${bookingErr.code}: ${bookingErr.message}`;
    } finally {
      isSubmitting = false;
    }
  }
</script>

<PageShell
  title="Form draft booking"
  subtitle={`${data.package.name} · channel b2c_self`}
  backHref={`/packages/${data.package.id}`}
  backLabel="Back to package"
>
  <div class="layout">
    <section class="s1-card" data-testid="s1-booking-draft-shell">
      <h2>Data ketua rombongan</h2>
      <form
        class="form"
        onsubmit={(e) => {
          e.preventDefault();
          void onSubmitDraft();
        }}
      >
        <div class="s1-field">
          <label for="lead-name">Nama lengkap</label>
          <input
            id="lead-name"
            class="s1-input"
            type="text"
            name="leadName"
            placeholder="Budi Santoso"
            bind:value={leadName}
          />
        </div>
        <div class="s1-field">
          <label for="lead-whatsapp">WhatsApp</label>
          <input
            id="lead-whatsapp"
            class="s1-input"
            type="tel"
            name="leadWhatsapp"
            placeholder="+62812..."
            bind:value={leadWhatsapp}
          />
        </div>
        <div class="s1-field">
          <label for="lead-email">Email</label>
          <input
            id="lead-email"
            class="s1-input"
            type="email"
            name="leadEmail"
            placeholder="nama@email.com"
            bind:value={leadEmail}
          />
        </div>
        <div class="s1-field">
          <label for="lead-domicile">Domisili</label>
          <input
            id="lead-domicile"
            class="s1-input"
            type="text"
            name="leadDomicile"
            placeholder="Jakarta"
            bind:value={leadDomicile}
          />
        </div>
        <div class="s1-field">
          <label for="departure-id">Keberangkatan</label>
          <select id="departure-id" class="s1-input" bind:value={departureId}>
            {#each data.package.departures as departure (departure.id)}
              <option value={departure.id}>
                {departure.departureDate} s.d. {departure.returnDate} · sisa {departure.remainingSeats}
              </option>
            {/each}
          </select>
        </div>
        <div class="s1-field">
          <label for="room-type">Tipe kamar</label>
          <select id="room-type" class="s1-input" bind:value={roomType}>
            <option value="quad">quad</option>
            <option value="triple">triple</option>
            <option value="double">double</option>
          </select>
        </div>
        <div class="s1-field">
          <label for="jamaah-count">Jumlah jamaah</label>
          <input id="jamaah-count" class="s1-input" type="number" min="1" bind:value={jamaahCount} />
        </div>
        {#if submitError}
          <p class="state state-error" role="alert">{submitError}</p>
        {/if}
        {#if submitSuccess}
          <p class="state state-success" role="status">
            Draft berhasil dibuat: <strong>{submitSuccess.bookingId}</strong> ({submitSuccess.createdAt})
          </p>
        {/if}
        <div class="s1-actions">
          <button
            type="submit"
            class="s1-btn s1-btn--primary"
            data-testid="s1-save-draft-stub"
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Menyimpan...' : 'Simpan draft'}
          </button>
          <span class="s1-muted">Draft disimpan via `POST /v1/bookings` (mock-first, API-ready).</span>
        </div>
      </form>
    </section>

    <aside class="s1-card summary">
      <h2>Ringkasan paket</h2>
      <p class="package">{data.package.name}</p>
      <p class="s1-muted">{data.package.startingPriceLabel}</p>
      <ul>
        {#each data.package.highlights as item (item)}
          <li>{item}</li>
        {/each}
      </ul>
      <a class="s1-btn s1-btn--ghost" href={`/packages/${data.package.id}`}>Kembali ke detail paket</a>
    </aside>
  </div>
</PageShell>

<style>
  .layout {
    display: grid;
    gap: var(--space-4);
    grid-template-columns: 1.4fr 1fr;
  }
  .summary {
    align-self: start;
  }
  .summary ul {
    margin: var(--space-3) 0;
    padding-left: 1rem;
  }
  .summary li {
    margin-bottom: var(--space-2);
  }
  .package {
    margin: 0 0 var(--space-1) 0;
    font-size: 1rem;
    font-weight: 600;
  }
  .form {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
    margin-top: var(--space-3);
  }
  .state {
    margin: 0;
    border-radius: calc(var(--radius-card) - 2px);
    padding: var(--space-2) var(--space-3);
    font-size: 0.88rem;
  }
  .state-error {
    border: 1px solid #7f1d1d;
    background: #2a1212;
    color: #fecaca;
  }
  .state-success {
    border: 1px solid #14532d;
    background: #10281a;
    color: #bbf7d0;
  }
  @media (max-width: 860px) {
    .layout {
      grid-template-columns: 1fr;
    }
  }
</style>
