<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  let sosPressed = $state(false);
  let sosCountdown = $state(0);
  let sosSent = $state(false);

  function pressSOS() {
    if (sosSent) return;
    sosPressed = true;
    sosCountdown = 5;
    const interval = setInterval(() => {
      sosCountdown--;
      if (sosCountdown <= 0) {
        clearInterval(interval);
        sosPressed = false;
        sosSent = true;
      }
    }, 1000);
  }

  function cancelSOS() {
    sosPressed = false;
    sosCountdown = 0;
  }

  function resetSOS() {
    sosSent = false;
  }

  const contacts = [
    { name: 'Pembimbing Rombongan', role: 'Ustaz Ahmad Fauzi', phone: '+62 812-0000-0001', icon: 'person' },
    { name: 'Maktab / PPIH Saudi', role: 'Kantor Urusan Haji Indonesia', phone: '+966 12 555-0001', icon: 'domain' },
    { name: 'KBRI Riyadh', role: 'Kontak Darurat Kedutaan', phone: '+966 11 488-0000', icon: 'account_balance' },
    { name: 'Rumah Sakit Arab Saudi', role: 'Darurat Medis 24 Jam', phone: '911', icon: 'local_hospital' },
    { name: 'Kantor UmrohOS', role: 'Hotline 24/7 Indonesia', phone: '+62 800-000-0000', icon: 'support_agent' },
  ];

  const faqs = [
    {
      q: 'Saya terpisah dari rombongan. Apa yang harus dilakukan?',
      a: 'Tetap tenang. Pergi ke titik kumpul yang telah ditentukan (Lobi Hotel Grand Makkah). Hubungi pembimbing rombongan. Jika tidak bisa, hubungi KBRI Riyadh atau tekan tombol SOS di atas untuk memberi tahu tim kami.'
    },
    {
      q: 'Saya kehilangan paspor / dokumen penting.',
      a: 'Segera laporkan ke pembimbing rombongan. Hubungi KBRI Riyadh untuk proses dokumen darurat. Jangan panik — tim kami memiliki salinan digital semua dokumen Anda.'
    },
    {
      q: 'Saya atau anggota rombongan sakit / membutuhkan pertolongan medis.',
      a: 'Hubungi pembimbing langsung. Untuk darurat medis segera, hubungi 911 (Saudi Arabia). Klinik kesehatan jemaah Indonesia tersedia di sekitar hotel.'
    },
    {
      q: 'Bagaimana jika saya ketinggalan bus?',
      a: 'Tetap di lokasi keberangkatan terakhir Anda. Hubungi pembimbing rombongan. Gunakan transportasi resmi (taksi berlisensi) jika perlu, dan simpan struk.'
    },
  ];

  let openFaq = $state<number | null>(null);

  const location = 'Hotel Grand Makkah, Al-Aziziyah, Makkah';
</script>

<svelte:head>
  <title>Darurat — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="darurat-root">
    <div class="shell">
      <div class="page-header">
        <a href="/jemaah" class="back-link">
          <span class="material-symbols-outlined">arrow_back</span>
          Portal Jamaah
        </a>
        <h1>Darurat & SOS</h1>
        <p>Bantuan segera untuk situasi darurat perjalanan</p>
      </div>

      <!-- SOS Button -->
      <div class="sos-section">
        {#if sosSent}
          <div class="sos-sent">
            <span class="material-symbols-outlined sos-sent-icon">check_circle</span>
            <div class="sos-sent-title">SOS Berhasil Dikirim!</div>
            <div class="sos-sent-msg">Tim kami telah menerima sinyal darurat Anda. Pembimbing rombongan akan segera menghubungi Anda.</div>
            <button class="sos-reset-btn" onclick={resetSOS}>Reset</button>
          </div>
        {:else if sosPressed}
          <div class="sos-countdown-wrapper">
            <div class="sos-btn pressing">
              <span class="sos-counter">{sosCountdown}</span>
              <span class="sos-label-sm">Mengirim SOS...</span>
            </div>
            <button class="cancel-sos" onclick={cancelSOS}>Batalkan</button>
          </div>
        {:else}
          <button class="sos-btn" onclick={pressSOS}>
            <span class="sos-icon material-symbols-outlined">emergency</span>
            <span class="sos-text">SOS DARURAT</span>
            <span class="sos-hint">Tahan 5 detik untuk kirim</span>
          </button>
        {/if}

        <div class="location-info">
          <span class="material-symbols-outlined">location_on</span>
          <div>
            <div class="loc-label">Lokasi Saya (Perkiraan)</div>
            <div class="loc-val">{location}</div>
          </div>
        </div>
      </div>

      <!-- Quick contacts -->
      <h2 class="section-title">Kontak Darurat</h2>
      <div class="contacts-list">
        {#each contacts as c}
          <a href="tel:{c.phone}" class="contact-card">
            <div class="contact-icon">
              <span class="material-symbols-outlined">{c.icon}</span>
            </div>
            <div class="contact-info">
              <div class="contact-name">{c.name}</div>
              <div class="contact-role">{c.role}</div>
              <div class="contact-phone">{c.phone}</div>
            </div>
            <span class="material-symbols-outlined call-icon">call</span>
          </a>
        {/each}
      </div>

      <!-- FAQ accordion -->
      <h2 class="section-title">Pertanyaan Darurat Umum</h2>
      <div class="faq-list">
        {#each faqs as item, idx}
          <div class="faq-item">
            <button class="faq-q" onclick={() => openFaq = openFaq === idx ? null : idx}>
              <span>{item.q}</span>
              <span class="material-symbols-outlined faq-arrow" class:open={openFaq === idx}>expand_more</span>
            </button>
            {#if openFaq === idx}
              <div class="faq-a">{item.a}</div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .darurat-root {
    padding-top: calc(5.2rem + 2rem);
    padding-bottom: 5rem;
    background: #fbf9f8;
    min-height: 100vh;
  }
  .shell { max-width: 64rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #ba1a1a; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  /* SOS section */
  .sos-section {
    background: #fff;
    border-radius: 1.5rem;
    padding: 2.5rem;
    border: 1px solid rgba(190,201,193,0.2);
    margin-bottom: 2.5rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
  }
  .sos-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    width: 180px;
    height: 180px;
    border-radius: 50%;
    background: linear-gradient(135deg, #ba1a1a, #dc2626);
    border: 6px solid rgba(186,26,26,0.15);
    color: #fff;
    cursor: pointer;
    transition: transform 0.15s, box-shadow 0.15s;
    box-shadow: 0 8px 32px rgba(186,26,26,0.3);
  }
  .sos-btn:hover { transform: scale(1.04); box-shadow: 0 12px 40px rgba(186,26,26,0.4); }
  .sos-btn.pressing { animation: pulse 0.8s ease-in-out infinite; background: linear-gradient(135deg, #7f1d1d, #ba1a1a); }
  @keyframes pulse { 0%, 100% { box-shadow: 0 8px 32px rgba(186,26,26,0.3); } 50% { box-shadow: 0 8px 48px rgba(186,26,26,0.6); } }
  .sos-icon { font-size: 2.5rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 48; }
  .sos-text { font-size: 1rem; font-weight: 800; letter-spacing: 0.05em; font-family: 'Plus Jakarta Sans', sans-serif; }
  .sos-hint { font-size: 0.68rem; opacity: 0.8; }
  .sos-counter { font-size: 3rem; font-weight: 800; font-family: 'Plus Jakarta Sans', sans-serif; }
  .sos-label-sm { font-size: 0.78rem; opacity: 0.9; }
  .sos-countdown-wrapper { display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .cancel-sos { background: none; border: 1.5px solid #ba1a1a; color: #ba1a1a; font-weight: 700; border-radius: 999px; padding: 0.5rem 1.5rem; cursor: pointer; font-size: 0.85rem; }
  .sos-sent { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; text-align: center; }
  .sos-sent-icon { font-size: 3.5rem; color: #006747; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 48; }
  .sos-sent-title { font-size: 1.3rem; font-weight: 800; color: #006747; font-family: 'Plus Jakarta Sans', sans-serif; }
  .sos-sent-msg { font-size: 0.9rem; color: #57534e; max-width: 32rem; line-height: 1.6; }
  .sos-reset-btn { background: none; border: 1.5px solid #d1d5db; color: #9ca3af; border-radius: 999px; padding: 0.4rem 1.2rem; cursor: pointer; font-size: 0.8rem; margin-top: 0.5rem; }
  .location-info { display: flex; align-items: flex-start; gap: 0.75rem; background: rgba(190,201,193,0.1); border-radius: 1rem; padding: 1rem 1.25rem; width: 100%; }
  .location-info .material-symbols-outlined { color: #006747; font-size: 1.2rem; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .loc-label { font-size: 0.72rem; color: #9ca3af; text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 0.2rem; }
  .loc-val { font-size: 0.88rem; color: #1b1c1c; font-weight: 600; }
  /* Contacts */
  .section-title { margin: 0 0 1.2rem; font-size: 1.1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .contacts-list { display: grid; gap: 0.85rem; margin-bottom: 2.5rem; }
  .contact-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem 1.2rem;
    background: #fff;
    border-radius: 1.2rem;
    border: 1px solid rgba(190,201,193,0.2);
    text-decoration: none;
    color: inherit;
    transition: box-shadow 0.15s;
  }
  .contact-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.07); }
  .contact-icon { width: 3rem; height: 3rem; border-radius: 0.85rem; background: rgba(186,26,26,0.1); display: grid; place-items: center; flex-shrink: 0; color: #ba1a1a; }
  .contact-icon .material-symbols-outlined { font-size: 1.3rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .contact-info { flex: 1; }
  .contact-name { font-size: 0.92rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .contact-role { font-size: 0.78rem; color: #9ca3af; margin-top: 0.1rem; }
  .contact-phone { font-size: 0.85rem; color: #006747; font-weight: 600; margin-top: 0.2rem; }
  .call-icon { color: #006747; font-size: 1.4rem; flex-shrink: 0; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  /* FAQ */
  .faq-list { display: flex; flex-direction: column; gap: 0.75rem; }
  .faq-item { background: #fff; border-radius: 1rem; border: 1px solid rgba(190,201,193,0.2); overflow: hidden; }
  .faq-q { display: flex; align-items: center; justify-content: space-between; gap: 1rem; padding: 1.1rem 1.3rem; background: none; border: none; cursor: pointer; font-size: 0.9rem; font-weight: 600; color: #1b1c1c; text-align: left; width: 100%; font-family: inherit; }
  .faq-arrow { color: #9ca3af; font-size: 1.2rem; transition: transform 0.2s; flex-shrink: 0; }
  .faq-arrow.open { transform: rotate(180deg); color: #006747; }
  .faq-a { padding: 0 1.3rem 1.1rem; font-size: 0.85rem; color: #57534e; line-height: 1.7; }
</style>
