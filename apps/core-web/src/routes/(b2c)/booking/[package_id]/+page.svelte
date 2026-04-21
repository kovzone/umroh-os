<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import MarketingPageLayout from '$lib/components/marketing/MarketingPageLayout.svelte';
  import BookingStepper from '$lib/components/booking/BookingStepper.svelte';
  import BookingHelpCard from '$lib/components/booking/BookingHelpCard.svelte';
  import BookingSummaryAside from '$lib/components/booking/BookingSummaryAside.svelte';
  import { createDraftBooking } from '$lib/features/s1-booking/repository';
  import type { DraftBookingError } from '$lib/features/s1-booking/types';
  import type { DepartureDetail, DepartureSummary, RoomType } from '$lib/features/s1-catalog/types';
  import { formatIdrAmountLabel, parseIdrAmountLabel } from '$lib/utils/format-idr';
  import {
    airlineOrDefault,
    departureSeatsLabel,
    formatDepartureDayId
  } from '$lib/utils/format-departure';

  let { data } = $props();

  const pkg = $derived(data.package);
  const wa = $derived(pkg.whatsappHref ?? 'https://wa.me/6281200000000');

  const step = $derived.by(() => {
    const raw = Number(page.url.searchParams.get('step') ?? '1');
    return raw >= 1 && raw <= 4 ? raw : 1;
  });

  const departureId = $derived.by(() => {
    const list = pkg.departures;
    const q = page.url.searchParams.get('departure');
    if (q && list.some((d) => d.id === q)) {
      return q;
    }
    if (data.preferredDepartureId && list.some((d) => d.id === data.preferredDepartureId)) {
      return data.preferredDepartureId;
    }
    return list[0]?.id ?? '';
  });

  const selectedDep = $derived(pkg.departures.find((d) => d.id === departureId));
  const selectedDetail = $derived(data.departureDetailById[departureId] ?? null);

  function maxJamaahForRoom(r: RoomType): number {
    if (r === 'triple') {
      return 3;
    }
    if (r === 'double') {
      return 2;
    }
    return 4;
  }

  function unitPriceIdrFromDetail(detail: DepartureDetail | null, rt: RoomType): number | null {
    if (!detail) {
      return null;
    }
    const row = detail.pricing.find((p) => p.roomType === rt);
    if (row?.listAmountIdr != null) {
      return row.listAmountIdr;
    }
    const parsed = parseIdrAmountLabel(row?.amountLabel);
    if (parsed != null) {
      return parsed;
    }
    const quad = detail.pricing.find((p) => p.roomType === 'quad');
    if (quad?.listAmountIdr != null) {
      return quad.listAmountIdr;
    }
    return parseIdrAmountLabel(quad?.amountLabel);
  }

  function durationLabelFrom(dep: DepartureSummary | undefined): string {
    if (!dep) {
      return '—';
    }
    const a = new Date(dep.departureDate);
    const b = new Date(dep.returnDate);
    const nights = Math.max(0, Math.round((b.getTime() - a.getTime()) / 86400000));
    return `${nights + 1} Hari`;
  }

  const durationLabel = $derived(durationLabelFrom(selectedDep));

  function roomTypeLabel(r: RoomType): string {
    const m: Record<RoomType, string> = {
      quad: 'Quad (4 orang)',
      triple: 'Triple (3 orang)',
      double: 'Double (2 orang)'
    };
    return m[r];
  }

  function maskNik(s: string): string {
    const t = s.replace(/\D/g, '');
    if (t.length < 4) {
      return '—';
    }
    return `****${t.slice(-4)}`;
  }

  function gotoQuery(partial: Record<string, string | undefined>) {
    const u = new URL(page.url);
    for (const [k, v] of Object.entries(partial)) {
      if (v === undefined) {
        u.searchParams.delete(k);
      } else {
        u.searchParams.set(k, v);
      }
    }
    void goto(`${u.pathname}${u.search}`, { noScroll: false });
  }

  function setDeparture(id: string) {
    gotoQuery({ departure: id, step: String(step) });
  }

  function gotoStep(next: number) {
    gotoQuery({ step: String(next), departure: departureId || undefined });
  }

  type JamaahRow = { fullName: string; nik: string; dob: string; passport: string };

  let leadName = $state('');
  let leadEmail = $state('');
  let leadWhatsapp = $state('');
  let leadDomicile = $state('');
  let roomType = $state<RoomType>('quad');
  let jamaahCount = $state(2);
  let jamaahRows = $state<JamaahRow[]>([
    { fullName: '', nik: '', dob: '', passport: '' },
    { fullName: '', nik: '', dob: '', passport: '' }
  ]);

  const unitPriceIdr = $derived(unitPriceIdrFromDetail(selectedDetail, roomType));
  const unitPriceLabel = $derived(
    unitPriceIdr != null ? formatIdrAmountLabel(unitPriceIdr) : pkg.startingPriceLabel
  );
  const billJamaahCount = $derived(jamaahRows.length);
  const totalPriceIdr = $derived(
    unitPriceIdr != null && billJamaahCount > 0 ? unitPriceIdr * billJamaahCount : null
  );
  const totalPriceLabel = $derived(
    totalPriceIdr != null ? formatIdrAmountLabel(totalPriceIdr) : pkg.startingPriceLabel
  );
  const priceBreakdownHint = $derived(
    unitPriceIdr != null && billJamaahCount > 0
      ? `${billJamaahCount} jamaah × ${unitPriceLabel}`
      : ''
  );

  $effect(() => {
    const cap = maxJamaahForRoom(roomType);
    let n = jamaahCount;
    if (!Number.isFinite(n) || n < 1) {
      n = 1;
    }
    if (n > cap) {
      n = cap;
    }
    if (n !== jamaahCount) {
      jamaahCount = n;
      return;
    }
    const rows = jamaahRows;
    if (rows.length < n) {
      jamaahRows = [
        ...rows,
        ...Array.from({ length: n - rows.length }, () => ({
          fullName: '',
          nik: '',
          dob: '',
          passport: ''
        }))
      ];
    } else if (rows.length > n) {
      jamaahRows = rows.slice(0, n);
    }
  });

  let agreePassport = $state(false);
  let agreeVaccine = $state(false);
  let agreeItinerary = $state(false);
  let agreeLegal = $state(false);

  const reviewReady = $derived(
    agreePassport && agreeVaccine && agreeItinerary && agreeLegal
  );

  let paymentMethod = $state<'va' | 'transfer' | 'ewallet'>('va');
  let draftResult = $state<{ bookingId: string; createdAt: string } | null>(null);
  let isSubmitting = $state(false);
  let submitError = $state('');

  function quadPriceForDeparture(id: string): string {
    const d = data.departureDetailById[id];
    return d?.pricing.find((p) => p.roomType === 'quad')?.amountLabel ?? '—';
  }

  function step1Continue() {
    if (!departureId) {
      return;
    }
    gotoStep(2);
  }

  function step2Valid(): boolean {
    if (!leadName.trim() || !leadEmail.trim() || !leadWhatsapp.trim()) {
      return false;
    }
    for (const row of jamaahRows) {
      if (!row.fullName.trim() || !row.nik.trim() || !row.dob) {
        return false;
      }
    }
    return true;
  }

  async function submitDraftAndGoPayment() {
    submitError = '';
    if (!reviewReady || !departureId) {
      return;
    }
    isSubmitting = true;
    try {
      const result = await createDraftBooking({
        channel: 'b2c_self',
        packageId: pkg.id,
        departureId,
        roomType,
        jamaahCount: jamaahRows.length,
        lead: {
          fullName: leadName.trim(),
          whatsapp: leadWhatsapp.trim(),
          email: leadEmail.trim(),
          domicile: leadDomicile.trim() || '—'
        }
      });
      draftResult = { bookingId: result.bookingId, createdAt: result.createdAt };
      gotoStep(4);
    } catch (err) {
      const bookingErr = err as DraftBookingError;
      submitError = `${bookingErr.code}: ${bookingErr.message}`;
    } finally {
      isSubmitting = false;
    }
  }

  const bookingCodeDisplay = $derived(
    draftResult ? `UMR-${draftResult.bookingId.slice(0, 8).toUpperCase()}` : '—'
  );

  function copyText(text: string) {
    void navigator.clipboard?.writeText(text.replace(/\s/g, ''));
  }
</script>

<MarketingPageLayout ctaHref="/packages" ctaLabel="Masuk" packagesLinkActive={true}>
  <div class="bw" data-testid="s1-booking-draft-shell">
    <main class="shell bw-main">
      <a class="back" href={`/packages/${pkg.id}`}>
        <span class="material-symbols-outlined">arrow_back</span>
        Kembali ke detail paket
      </a>

      <BookingStepper current={step} />

      <div class="grid">
        <div class="main-col">
          {#if step === 1}
            <section class="card card-pad">
              <h2 class="h2">Pilih jadwal keberangkatan</h2>
              <div class="dep-list">
                {#each pkg.departures as dep (dep.id)}
                  {@const seats = departureSeatsLabel(dep)}
                  {@const priceLine = quadPriceForDeparture(dep.id)}
                  {@const selected = dep.id === departureId}
                  <button
                    type="button"
                    class="dep-row"
                    class:selected
                    onclick={() => setDeparture(dep.id)}
                  >
                    <div class="dep-left">
                      <div class="radio" class:filled={selected} aria-hidden="true">
                        {#if selected}
                          <span class="dot"></span>
                        {/if}
                      </div>
                      <div>
                        <p class="dep-date">{formatDepartureDayId(dep)}</p>
                        <div class="airline">
                          <span class="material-symbols-outlined ap">flight</span>
                          <span>{airlineOrDefault(dep)}</span>
                        </div>
                      </div>
                    </div>
                    <div class="dep-right">
                      <p class="price">{priceLine}</p>
                      <div class="pill" class:urgent={seats.urgent}>
                        <span class="dot-s"></span>
                        <span>{seats.text}</span>
                      </div>
                    </div>
                  </button>
                {/each}
              </div>
              <div class="actions-row">
                <a class="link-muted" href={`/packages/${pkg.id}`}>Simpan dulu (lanjut nanti)</a>
                <button type="button" class="btn-primary" onclick={() => step1Continue()}>Lanjut</button>
              </div>
            </section>
          {:else if step === 2}
            <header class="page-h">
              <h1 class="h1">Data jamaah</h1>
              <p class="lead">
                Mohon lengkapi data seluruh jamaah sesuai dengan identitas resmi (KTP/Paspor).
              </p>
            </header>
            <section class="card card-pad">
              <div class="sec-head">
                <span class="material-symbols-outlined fill sec-ico">contact_mail</span>
                <h2 class="h3">Kontak pemesan</h2>
              </div>
              <div class="form-grid">
                <div class="field">
                  <label for="ln">Nama lengkap</label>
                  <input id="ln" bind:value={leadName} type="text" placeholder="Contoh: Ahmad Fauzan" />
                </div>
                <div class="field">
                  <label for="em">Email</label>
                  <input id="em" bind:value={leadEmail} type="email" placeholder="nama@email.com" />
                </div>
                <div class="field">
                  <label for="wa">Nomor WhatsApp</label>
                  <input id="wa" bind:value={leadWhatsapp} type="tel" placeholder="+62 812 3456 7890" />
                </div>
                <div class="field">
                  <label for="dom">Alamat singkat (opsional)</label>
                  <input id="dom" bind:value={leadDomicile} type="text" placeholder="Kota atau domisili" />
                </div>
                <div class="field span2">
                  <label for="rtc">Tipe kamar</label>
                  <select id="rtc" bind:value={roomType} class="inp-select">
                    <option value="quad">Quad (4 orang)</option>
                    <option value="triple">Triple (3 orang)</option>
                    <option value="double">Double (2 orang)</option>
                  </select>
                </div>
                <div class="field span2">
                  <label for="jc">Jumlah jamaah</label>
                  <input
                    id="jc"
                    type="number"
                    min="1"
                    max={maxJamaahForRoom(roomType)}
                    bind:value={jamaahCount}
                  />
                </div>
              </div>
            </section>
            <section class="card card-pad">
              <div class="sec-head">
                <span class="material-symbols-outlined fill sec-ico">group</span>
                <h2 class="h3">
                  Data jamaah <span class="muted-inline">({jamaahRows.length} dari {roomType === 'quad' ? 4 : roomType === 'triple' ? 3 : 2} kuota)</span>
                </h2>
              </div>
              {#each jamaahRows as row, i (i)}
                <div class="jamaah-block" class:mb={i < jamaahRows.length - 1}>
                  <div class="jamaah-head">
                    <div class="num-badge" class:primary={i === 0}>{i + 1}</div>
                    <h3 class="jh">Jamaah {i + 1}</h3>
                    {#if i === 0}
                      <span class="tag">Pemesan</span>
                    {/if}
                  </div>
                  <div class="form-grid">
                    <div class="field">
                      <label for="jn{i}">Nama sesuai KTP</label>
                      <input id="jn{i}" bind:value={row.fullName} type="text" />
                    </div>
                    <div class="field">
                      <label for="nik{i}">NIK</label>
                      <input id="nik{i}" bind:value={row.nik} type="text" inputmode="numeric" />
                    </div>
                    <div class="field">
                      <label for="dob{i}">Tanggal lahir</label>
                      <input id="dob{i}" bind:value={row.dob} type="date" />
                    </div>
                    <div class="field">
                      <label for="pp{i}">Paspor (opsional)</label>
                      <input id="pp{i}" bind:value={row.passport} type="text" placeholder="No. paspor jika ada" />
                    </div>
                  </div>
                </div>
              {/each}
            </section>
            <div class="actions-row">
              <a class="link-muted" href={`/packages/${pkg.id}`}>Simpan dulu (lanjut nanti)</a>
              <button
                type="button"
                class="btn-primary"
                disabled={!step2Valid()}
                onclick={() => gotoStep(3)}
              >
                Lanjut ke Review
              </button>
            </div>
          {:else if step === 3}
            <header class="page-h">
              <h1 class="h1">Review pemesanan</h1>
              <p class="lead">Silakan periksa kembali detail keberangkatan dan kelengkapan dokumen Anda.</p>
            </header>
            <section class="card card-pad bordered">
              <div class="sec-head">
                <span class="material-symbols-outlined sec-ico">flight_takeoff</span>
                <h2 class="h3">Ringkasan keberangkatan</h2>
              </div>
              <div class="review-grid">
                <div>
                  <p class="rk">Tanggal</p>
                  <p class="rv">{selectedDep ? formatDepartureDayId(selectedDep) : '—'}</p>
                </div>
                <div>
                  <p class="rk">Maskapai</p>
                  <p class="rv">{selectedDep ? airlineOrDefault(selectedDep) : '—'}</p>
                </div>
                <div>
                  <p class="rk">Durasi</p>
                  <p class="rv">{durationLabel.toLowerCase()}</p>
                </div>
                <div>
                  <p class="rk">Tipe kamar</p>
                  <span class="chip">{roomTypeLabel(roomType)}</span>
                </div>
              </div>
            </section>
            <section class="card card-pad bordered">
              <div class="sec-head space">
                <div class="sec-head">
                  <span class="material-symbols-outlined sec-ico">groups</span>
                  <h2 class="h3">Daftar jamaah</h2>
                </div>
                <button type="button" class="link-btn" onclick={() => gotoStep(2)}>Ubah data</button>
              </div>
              <div class="jlist">
                {#each jamaahRows as row, i (i)}
                  <div class="jrow">
                    <div class="jleft">
                      <div class="num-lg">{i + 1}</div>
                      <div>
                        <p class="jname">{row.fullName.trim() || `Jamaah ${i + 1}`}</p>
                        <p class="jnik">NIK: {maskNik(row.nik)}</p>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            </section>
            <section class="card card-pad bordered">
              <div class="sec-head">
                <span class="material-symbols-outlined sec-ico">fact_check</span>
                <h2 class="h3">Syarat & persetujuan</h2>
              </div>
              <div class="checks">
                <label class="chk-card">
                  <input type="checkbox" bind:checked={agreePassport} />
                  <div>
                    <p class="ct">Dokumen valid min 6 bulan</p>
                    <p class="cs">
                      Saya menjamin bahwa semua paspor jamaah memiliki masa berlaku minimal 6 bulan dari tanggal
                      keberangkatan.
                    </p>
                  </div>
                </label>
                <label class="chk-card">
                  <input type="checkbox" bind:checked={agreeVaccine} />
                  <div>
                    <p class="ct">Vaksin sesuai ketentuan</p>
                    <p class="cs">
                      Seluruh jamaah telah/akan mendapatkan vaksinasi yang dipersyaratkan oleh pemerintah Arab Saudi.
                    </p>
                  </div>
                </label>
                <label class="chk-card">
                  <input type="checkbox" bind:checked={agreeItinerary} />
                  <div>
                    <p class="ct">Setuju dengan itinerary</p>
                    <p class="cs">
                      Saya memahami dan menyetujui jadwal perjalanan serta fasilitas yang telah ditentukan dalam paket
                      ini.
                    </p>
                  </div>
                </label>
              </div>
              <div class="legal-row">
                <label class="legal">
                  <input type="checkbox" bind:checked={agreeLegal} />
                  <span>
                    Saya telah membaca dan menyetujui
                    <a href="/">Syarat & Ketentuan</a>
                    serta
                    <a href="/">Kebijakan Privasi</a>
                  </span>
                </label>
              </div>
              {#if submitError}
                <p class="err" role="alert">{submitError}</p>
              {/if}
            </section>
            <div class="nav-2">
              <button type="button" class="btn-outline" onclick={() => gotoStep(2)}>
                <span class="material-symbols-outlined">arrow_back</span>
                Kembali
              </button>
              <button
                type="button"
                class="btn-primary"
                disabled={!reviewReady || isSubmitting}
                onclick={() => void submitDraftAndGoPayment()}
              >
                {isSubmitting ? 'Menyimpan…' : 'Lanjut ke pembayaran'}
                <span class="material-symbols-outlined">arrow_forward</span>
              </button>
            </div>
            <button type="button" class="link-muted center" onclick={() => gotoStep(2)}>
              Simpan dulu (lanjut nanti)
            </button>
          {:else if !draftResult}
            <section class="card card-pad bordered">
              <h2 class="h3">Pembayaran belum siap</h2>
              <p class="lead" style:margin-top="0.75rem">
                Selesaikan langkah review dan persetujuan terlebih dahulu untuk mendapatkan instruksi pembayaran.
              </p>
              <button type="button" class="btn-primary" style:margin-top="1.25rem" onclick={() => gotoStep(3)}>
                Ke review & syarat
              </button>
            </section>
          {:else}
            <header class="page-h">
              <h1 class="h1">Pembayaran</h1>
              <p class="lead">Selesaikan pembayaran sebelum batas waktu untuk mengamankan reservasi Anda.</p>
            </header>
            <div class="alert-time">
              <span class="material-symbols-outlined ic">schedule</span>
              <div>
                <p class="at-t">Bayar sebelum: 23:59:59</p>
                <p class="at-s">
                  Pesanan Anda akan dibatalkan secara otomatis jika melewati batas waktu (contoh UI).
                </p>
              </div>
            </div>
            <section class="card card-pad bordered">
              <div class="pay-top">
                <div>
                  <p class="rk">Total yang harus dibayar</p>
                  <p class="pay-big">{totalPriceLabel}</p>
                </div>
                <div class="code-box">
                  <div>
                    <p class="rk">Kode booking</p>
                    <p class="code">{bookingCodeDisplay}</p>
                  </div>
                  <button
                    type="button"
                    class="icon-btn"
                    aria-label="Salin kode"
                    onclick={() => copyText(bookingCodeDisplay)}
                  >
                    <span class="material-symbols-outlined">content_copy</span>
                  </button>
                </div>
              </div>
            </section>
            <div>
              <h3 class="h3 mb">Metode pembayaran</h3>
              <div class="pay-grid">
                <label class="pay-tile" class:sel={paymentMethod === 'va'}>
                  <input type="radio" bind:group={paymentMethod} value="va" class="sr" />
                  <span class="pt-t">Virtual Account BCA</span>
                  <span class="pt-s">Verifikasi otomatis 24/7</span>
                </label>
                <label class="pay-tile" class:sel={paymentMethod === 'transfer'}>
                  <input type="radio" bind:group={paymentMethod} value="transfer" class="sr" />
                  <span class="pt-t">Transfer manual</span>
                  <span class="pt-s">Konfirmasi manual 1–2 jam</span>
                </label>
                <label class="pay-tile" class:sel={paymentMethod === 'ewallet'}>
                  <input type="radio" bind:group={paymentMethod} value="ewallet" class="sr" />
                  <span class="pt-t">E-wallet / QRIS</span>
                  <span class="pt-s">Biaya admin mengikuti penyedia</span>
                </label>
              </div>
            </div>
            {#if paymentMethod === 'va'}
              <section class="card card-pad instr">
                <h4 class="h4">Instruksi pembayaran BCA Virtual Account</h4>
                <div class="va-box">
                  <p class="rk">Nomor Virtual Account</p>
                  <p class="va-num">80008 12345 67890</p>
                  <button type="button" class="btn-primary sm" onclick={() => copyText('800081234567890')}>
                    <span class="material-symbols-outlined">content_copy</span>
                    Salin nomor VA
                  </button>
                </div>
                <ol class="steps">
                  <li>Buka aplikasi <strong>m-BCA</strong>, pilih <strong>m-Transfer</strong> lalu <strong>BCA Virtual Account</strong>.</li>
                  <li>Masukkan nomor VA dan pastikan nominal <strong>{totalPriceLabel}</strong>.</li>
                  <li>Ikuti instruksi hingga transaksi berhasil.</li>
                </ol>
                <p class="hint">
                  <span class="material-symbols-outlined">info</span>
                  Pembayaran akan terverifikasi otomatis (mock UI; integrasi payment menyusul).
                </p>
              </section>
            {/if}
            <div class="actions-row bottom">
              <span class="link-muted">Unduh ringkasan (PDF) — segera hadir</span>
              <div class="btns">
                <a class="btn-outline sm" href={wa} target="_blank" rel="noopener noreferrer">Hubungi CS</a>
                <button
                  type="button"
                  class="btn-primary sm"
                  onclick={() => void goto(`/packages/${pkg.id}`)}
                >
                  Saya sudah bayar
                </button>
              </div>
            </div>
          {/if}
        </div>

        <aside class="aside-col">
          {#if step < 4 || (step === 4 && !draftResult)}
            <BookingSummaryAside
              pkg={data.package}
              {durationLabel}
              roomLabel={roomTypeLabel(roomType)}
              totalLabel={totalPriceLabel}
              breakdownHint={priceBreakdownHint}
              departureDayLabel={selectedDep ? formatDepartureDayId(selectedDep) : '—'}
              jamaahCount={jamaahRows.length}
              step={step === 4 && !draftResult ? 3 : step}
            />
            <BookingHelpCard whatsappHref={wa} />
          {:else}
            <div class="status-chip">
              <span class="pulse"></span>
              <span>Menunggu pembayaran</span>
            </div>
            <div class="card card-pad bordered aside-pay">
              <h4 class="h4">Ringkasan paket</h4>
              <div class="pay-thumb">
                <img src={pkg.coverPhotoUrl} alt="" role="presentation" />
              </div>
              <p class="pkg-t">{pkg.name}</p>
              <p class="pkg-sub">{durationLabel} · {selectedDep ? formatDepartureDayId(selectedDep) : ''}</p>
              <div class="line-total">
                <span>Total akhir</span>
                <span class="amt">{totalPriceLabel}</span>
              </div>
            </div>
            <BookingHelpCard whatsappHref={wa} />
          {/if}
        </aside>
      </div>
    </main>
  </div>
</MarketingPageLayout>

<style>
  .bw-main {
    /* Top offset: MarketingPageLayout .top-nav-spacer (5.2rem) + breathing room */
    padding-top: 1.5rem;
    padding-bottom: 4rem;
  }
  .back {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    margin-bottom: 1.75rem;
    font-size: 0.9rem;
    font-weight: 600;
    color: #6f7a72;
    text-decoration: none;
    transition: color 0.15s ease, gap 0.15s ease;
  }
  .back:hover {
    color: #004d34;
    gap: 0.5rem;
  }
  .grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2.25rem;
    align-items: flex-start;
  }
  @media (min-width: 1024px) {
    .grid {
      grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
      gap: 2.5rem;
    }
  }
  .main-col {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .aside-col {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  @media (min-width: 1024px) {
    .aside-col {
      position: sticky;
      top: 5.75rem;
    }
  }
  .card {
    background: #fff;
    border-radius: 1.25rem;
    box-shadow: 0 8px 32px rgba(27, 28, 28, 0.03);
  }
  .card-pad {
    padding: 1.75rem 1.5rem;
  }
  @media (min-width: 640px) {
    .card-pad {
      padding: 2rem 2rem;
    }
  }
  .bordered {
    border: 1px solid rgba(190, 201, 193, 0.35);
  }
  .h1 {
    margin: 0 0 0.35rem;
    font-size: clamp(1.75rem, 3vw, 2.25rem);
    font-weight: 800;
    color: #004d34;
    letter-spacing: -0.02em;
  }
  .h2 {
    margin: 0 0 1.5rem;
    font-size: 1.35rem;
    font-weight: 800;
    color: #004d34;
  }
  .h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 700;
    color: #1b1c1c;
  }
  .h4 {
    margin: 0 0 0.75rem;
    font-size: 1rem;
    font-weight: 700;
    color: #1b1c1c;
  }
  .page-h .lead {
    margin: 0;
    color: #6f7a72;
    font-size: 0.95rem;
    line-height: 1.5;
  }
  .dep-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .dep-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    width: 100%;
    text-align: left;
    padding: 1.25rem 1.35rem;
    border-radius: 0.75rem;
    border: 2px solid transparent;
    background: #fbf9f8;
    cursor: pointer;
    transition:
      border-color 0.15s ease,
      background 0.15s ease;
    font: inherit;
    color: inherit;
  }
  .dep-row:hover {
    background: #f5f3f3;
  }
  .dep-row.selected {
    border-color: rgba(0, 103, 71, 0.35);
    background: #f5f3f3;
  }
  .dep-left {
    display: flex;
    align-items: center;
    gap: 1.1rem;
    min-width: 0;
  }
  .radio {
    width: 1.35rem;
    height: 1.35rem;
    border-radius: 999px;
    border: 2px solid #6f7a72;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .radio.filled {
    border-color: #006747;
    background: #006747;
  }
  .radio .dot {
    width: 0.45rem;
    height: 0.45rem;
    border-radius: 999px;
    background: #fff;
  }
  .dep-date {
    margin: 0 0 0.2rem;
    font-weight: 700;
    font-size: 1.05rem;
  }
  .airline {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.82rem;
    color: #6f7a72;
  }
  .airline .ap {
    font-size: 1rem;
    opacity: 0.65;
  }
  .dep-right {
    text-align: right;
    flex-shrink: 0;
  }
  .price {
    margin: 0 0 0.35rem;
    font-size: 1.15rem;
    font-weight: 800;
    color: #775a19;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .dep-row:not(.selected) .price {
    color: #1b1c1c;
  }
  .pill {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.65rem;
    font-weight: 700;
    padding: 0.2rem 0.65rem;
    border-radius: 999px;
    background: #e4e2e2;
    color: #6f7a72;
  }
  .pill.urgent {
    background: rgba(254, 212, 136, 0.45);
    color: #785a1a;
  }
  .pill .dot-s {
    width: 5px;
    height: 5px;
    border-radius: 999px;
    background: currentColor;
  }
  .actions-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-top: 1.75rem;
  }
  .actions-row.bottom {
    padding-top: 1.5rem;
    border-top: 1px solid rgba(190, 201, 193, 0.35);
  }
  .link-muted {
    background: none;
    border: none;
    padding: 0;
    font-size: 0.82rem;
    font-weight: 600;
    color: #6f7a72;
    cursor: pointer;
    text-decoration: none;
    font-family: inherit;
  }
  .link-muted:hover {
    color: #775a19;
  }
  .link-muted.center {
    display: block;
    margin: 0.75rem auto 0;
    text-align: center;
  }
  .btn-primary {
    border: none;
    cursor: pointer;
    padding: 0.9rem 1.75rem;
    border-radius: 999px;
    font-weight: 700;
    font-size: 1rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
    color: #fff;
    background: linear-gradient(90deg, #004d34, #006747);
    box-shadow: 0 12px 28px rgba(0, 77, 52, 0.22);
    transition: transform 0.15s ease, opacity 0.15s ease;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.35rem;
  }
  .btn-primary:hover:not(:disabled) {
    transform: scale(1.02);
  }
  .btn-primary:disabled {
    opacity: 0.45;
    cursor: not-allowed;
  }
  .btn-primary.sm {
    padding: 0.65rem 1.25rem;
    font-size: 0.85rem;
  }
  .btn-outline {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.85rem 1.25rem;
    border-radius: 0.75rem;
    border: 2px solid #004d34;
    background: transparent;
    color: #004d34;
    font-weight: 700;
    font-family: 'Plus Jakarta Sans', sans-serif;
    cursor: pointer;
    transition: background 0.15s ease;
  }
  .btn-outline:hover {
    background: rgba(0, 77, 52, 0.06);
  }
  .btn-outline.sm {
    padding: 0.6rem 1rem;
    font-size: 0.82rem;
    border-radius: 999px;
    text-decoration: none;
  }
  .sec-head {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1.25rem;
  }
  .sec-head.space {
    justify-content: space-between;
  }
  .sec-ico {
    color: #775a19;
    font-size: 1.35rem;
  }
  .form-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1.1rem;
  }
  @media (min-width: 720px) {
    .form-grid {
      grid-template-columns: 1fr 1fr;
    }
    .span2 {
      grid-column: span 2;
    }
  }
  .field label {
    display: block;
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #6f7a72;
    margin-bottom: 0.35rem;
    margin-left: 0.15rem;
  }
  .field input,
  .field select,
  .inp-select {
    width: 100%;
    box-sizing: border-box;
    border: none;
    border-radius: 0.75rem;
    padding: 0.75rem 1rem;
    background: #e9e8e7;
    font: inherit;
    color: #1b1c1c;
    outline: none;
    transition:
      background 0.15s ease,
      box-shadow 0.15s ease;
  }
  .field input:focus,
  .field select:focus {
    background: #fff;
    box-shadow: 0 0 0 1px #004d34;
  }
  .muted-inline {
    font-weight: 500;
    color: #6f7a72;
    font-size: 0.95rem;
  }
  .jamaah-block.mb {
    margin-bottom: 2rem;
    padding-bottom: 2rem;
    border-bottom: 1px dashed #efeded;
  }
  .jamaah-head {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    margin-bottom: 1rem;
  }
  .num-badge {
    width: 2rem;
    height: 2rem;
    border-radius: 999px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.8rem;
    font-weight: 700;
    background: #efeded;
    color: #6f7a72;
  }
  .num-badge.primary {
    background: #a0f4ca;
    color: #004d34;
  }
  .jh {
    margin: 0;
    font-size: 1.05rem;
    font-weight: 700;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .tag {
    font-size: 0.58rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    padding: 0.2rem 0.55rem;
    border-radius: 999px;
    background: #ffdea5;
    color: #261900;
  }
  .review-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1.25rem;
  }
  @media (min-width: 720px) {
    .review-grid {
      grid-template-columns: repeat(4, 1fr);
    }
  }
  .rk {
    margin: 0 0 0.25rem;
    font-size: 0.62rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #6f7a72;
  }
  .rv {
    margin: 0;
    font-weight: 600;
    font-size: 0.92rem;
  }
  .chip {
    display: inline-flex;
    padding: 0.25rem 0.65rem;
    border-radius: 999px;
    font-size: 0.68rem;
    font-weight: 700;
    background: #ffdea5;
    color: #5d4201;
  }
  .jlist {
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
  }
  .jrow {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.85rem 1rem;
    background: #fbf9f8;
    border-radius: 0.75rem;
  }
  .jleft {
    display: flex;
    align-items: center;
    gap: 0.85rem;
  }
  .num-lg {
    width: 2.25rem;
    height: 2.25rem;
    border-radius: 999px;
    background: rgba(0, 77, 52, 0.08);
    color: #004d34;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.9rem;
  }
  .jname {
    margin: 0 0 0.15rem;
    font-weight: 700;
  }
  .jnik {
    margin: 0;
    font-size: 0.82rem;
    color: #6f7a72;
  }
  .checks {
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
  }
  .chk-card {
    display: flex;
    gap: 0.85rem;
    padding: 1rem;
    border: 1px solid rgba(190, 201, 193, 0.45);
    border-radius: 0.75rem;
    cursor: pointer;
    transition: background 0.15s ease;
  }
  .chk-card:hover {
    background: #fbf9f8;
  }
  .chk-card input {
    margin-top: 0.2rem;
    width: 1.1rem;
    height: 1.1rem;
    accent-color: #004d34;
  }
  .ct {
    margin: 0 0 0.25rem;
    font-weight: 700;
    font-size: 0.92rem;
  }
  .cs {
    margin: 0;
    font-size: 0.82rem;
    color: #6f7a72;
    line-height: 1.45;
  }
  .legal-row {
    margin-top: 1.75rem;
    padding-top: 1.5rem;
    border-top: 1px solid rgba(190, 201, 193, 0.35);
  }
  .legal {
    display: flex;
    align-items: flex-start;
    gap: 0.65rem;
    cursor: pointer;
    font-size: 0.88rem;
  }
  .legal input {
    width: 1.15rem;
    height: 1.15rem;
    margin-top: 0.1rem;
    accent-color: #004d34;
  }
  .legal a {
    color: #775a19;
    font-weight: 700;
  }
  .nav-2 {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    gap: 1rem;
    padding-top: 0.5rem;
  }
  .link-btn {
    border: none;
    background: none;
    padding: 0;
    font-size: 0.82rem;
    font-weight: 700;
    color: #004d34;
    cursor: pointer;
    text-decoration: underline;
    font-family: inherit;
  }
  .err {
    margin: 1rem 0 0;
    padding: 0.65rem 0.85rem;
    border-radius: 0.5rem;
    background: #ffdad6;
    color: #93000a;
    font-size: 0.85rem;
  }
  .alert-time {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
    padding: 1.25rem 1.35rem;
    border-radius: 0.75rem;
    background: rgba(254, 212, 136, 0.35);
    border-left: 4px solid #775a19;
  }
  .alert-time .ic {
    font-size: 2rem;
    color: #775a19;
  }
  .at-t {
    margin: 0 0 0.25rem;
    font-size: 0.72rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: #775a19;
  }
  .at-s {
    margin: 0;
    font-size: 0.82rem;
    color: #5d4201;
    line-height: 1.45;
  }
  .pay-top {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  @media (min-width: 640px) {
    .pay-top {
      flex-direction: row;
      align-items: center;
      justify-content: space-between;
    }
  }
  .pay-big {
    margin: 0.25rem 0 0;
    font-size: 1.85rem;
    font-weight: 800;
    color: #775a19;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .code-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 0.65rem 1rem;
    background: #f5f3f3;
    border-radius: 0.5rem;
    min-width: min(100%, 15rem);
  }
  .code {
    margin: 0.15rem 0 0;
    font-weight: 700;
    color: #004d34;
  }
  .icon-btn {
    border: none;
    background: transparent;
    cursor: pointer;
    color: #775a19;
    border-radius: 999px;
    padding: 0.35rem;
  }
  .icon-btn:hover {
    background: rgba(119, 90, 25, 0.08);
  }
  .mb {
    margin-bottom: 0.65rem;
  }
  .pay-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.75rem;
  }
  @media (min-width: 720px) {
    .pay-grid {
      grid-template-columns: repeat(3, 1fr);
    }
  }
  .pay-tile {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
    padding: 1.1rem 1.15rem;
    border-radius: 0.75rem;
    border: 1px solid rgba(190, 201, 193, 0.45);
    background: #fff;
    cursor: pointer;
    transition:
      border-color 0.15s ease,
      background 0.15s ease;
  }
  .pay-tile.sel {
    border: 2px solid #004d34;
    background: #fff;
  }
  .pay-tile:hover {
    background: #f5f3f3;
  }
  .sr {
    position: absolute;
    top: 0.85rem;
    right: 0.85rem;
    width: 1rem;
    height: 1rem;
    accent-color: #004d34;
    cursor: pointer;
  }
  .pt-t {
    font-weight: 700;
    font-size: 0.92rem;
    color: #1b1c1c;
  }
  .pay-tile.sel .pt-t {
    color: #004d34;
  }
  .pt-s {
    font-size: 0.72rem;
    color: #6f7a72;
  }
  .instr .va-box {
    text-align: center;
    padding: 1.25rem;
    border: 1px solid rgba(190, 201, 193, 0.35);
    border-radius: 0.65rem;
    margin-bottom: 1.25rem;
    background: #fff;
  }
  .va-num {
    margin: 0.35rem 0 0.85rem;
    font-size: 1.35rem;
    font-weight: 800;
    font-family: ui-monospace, monospace;
    letter-spacing: 0.04em;
    color: #004d34;
  }
  .steps {
    margin: 0;
    padding-left: 1.1rem;
    font-size: 0.88rem;
    line-height: 1.55;
    color: #3f4943;
  }
  .steps li {
    margin-bottom: 0.5rem;
  }
  .hint {
    margin: 1.25rem 0 0;
    display: flex;
    gap: 0.35rem;
    font-size: 0.72rem;
    color: #6f7a72;
    font-style: italic;
  }
  .btns {
    display: flex;
    flex-wrap: wrap;
    gap: 0.65rem;
  }
  .status-chip {
    display: flex;
    align-items: center;
    gap: 0.65rem;
    padding: 0.65rem 1rem;
    border-radius: 0.75rem;
    background: rgba(119, 90, 25, 0.1);
    border: 1px solid rgba(119, 90, 25, 0.2);
    font-size: 0.78rem;
    font-weight: 700;
    font-family: 'Plus Jakarta Sans', sans-serif;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #775a19;
  }
  .pulse {
    width: 8px;
    height: 8px;
    border-radius: 999px;
    background: #775a19;
    animation: pulse 1.4s ease-in-out infinite;
  }
  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.35;
    }
  }
  .aside-pay .pay-thumb {
    height: 10rem;
    border-radius: 0.75rem;
    overflow: hidden;
    margin-bottom: 0.85rem;
  }
  .aside-pay .pay-thumb img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .pkg-t {
    margin: 0 0 0.25rem;
    font-weight: 700;
    color: #004d34;
    font-size: 0.95rem;
  }
  .pkg-sub {
    margin: 0 0 0.85rem;
    font-size: 0.75rem;
    color: #6f7a72;
  }
  .line-total {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 0.85rem;
    border-top: 1px dashed rgba(190, 201, 193, 0.5);
    font-weight: 800;
    font-size: 0.92rem;
  }
  .line-total .amt {
    color: #775a19;
    font-size: 1.05rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
</style>
