<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  const faqs = [
    { q: 'Apakah boleh wanita umroh tanpa mahram?', a: 'Berdasarkan peraturan Arab Saudi tahun 2021, wanita berusia 45 tahun ke atas diperbolehkan berangkat umroh tanpa mahram dalam rombongan resmi yang telah terdaftar. Wanita di bawah 45 tahun tetap harus didampingi mahram.' },
    { q: 'Bagaimana hukum melewati miqat tanpa berihram?', a: 'Melewati miqat tanpa berihram bagi yang berniat umroh atau haji adalah haram dan berdosa. Harus kembali ke miqat untuk berihram. Jika tidak bisa kembali, wajib membayar dam (denda berupa menyembelih kambing).' },
    { q: 'Apakah sa\'i harus langsung setelah thawaf?', a: 'Tidak harus langsung, boleh ada jeda waktu. Namun, sunnah-nya dilakukan segera setelah thawaf. Sa\'i dimulai dari Shafa dan berakhir di Marwah.' },
    { q: 'Bolehkah thawaf menggunakan kursi roda?', a: 'Ya, boleh. Orang yang tidak mampu berjalan boleh thawaf menggunakan kursi roda atau ditandu, bahkan boleh menggunakan lantai atas dengan fasilitas yang ada di Masjidil Haram.' },
    { q: 'Apakah anak-anak bisa melaksanakan umroh?', a: 'Ya, anak-anak boleh melaksanakan umroh dan sah hukumnya sebagai amalan sunnah. Namun umroh anak tidak menggantikan kewajiban umroh setelah baligh. Orang tua/wali yang mengihramkan anak.' },
  ];

  interface RecentAnswer {
    id: number;
    q: string;
    a: string;
    date: string;
    answered: boolean;
  }

  const recentAnswers: RecentAnswer[] = [
    { id: 1, q: 'Bagaimana hukum memakai masker saat ihram?', a: 'Memakai masker saat berihram diperbolehkan karena kebutuhan kesehatan, tidak termasuk larangan ihram. Larangan menutup wajah berlaku khusus untuk wanita yang menutup wajah dengan cadar/niqab.', date: '20 Jan 2025', answered: true },
    { id: 2, q: 'Apakah boleh menggunakan HP sambil thawaf?', a: 'Boleh menggunakan HP untuk membaca doa atau berdzikir, namun dianjurkan untuk fokus beribadah dan tidak terdistraksi. Jangan sampai mengganggu jamaah lain.', date: '19 Jan 2025', answered: true },
  ];

  let openFaq = $state<number | null>(null);
  let question = $state('');
  let submitting = $state(false);
  let submitted = $state(false);

  function submitQuestion(e: Event) {
    e.preventDefault();
    if (!question.trim()) return;
    submitting = true;
    setTimeout(() => {
      submitting = false;
      submitted = true;
      question = '';
      setTimeout(() => { submitted = false; }, 3000);
    }, 1200);
  }
</script>

<svelte:head>
  <title>Tanya Jawab Agama — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="tanya-root">
    <div class="shell">
      <a href="/ibadah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Panduan Ibadah
      </a>
      <div class="page-header">
        <h1>Tanya Jawab Agama</h1>
        <p>Pertanyaan dan jawaban seputar ibadah umroh & haji</p>
      </div>

      <!-- FAQ accordion -->
      <h2 class="section-title">Pertanyaan Sering Ditanyakan</h2>
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

      <!-- Submit question form -->
      <div class="submit-section">
        <h2 class="section-title">Kirim Pertanyaan</h2>
        <div class="submit-card">
          {#if submitted}
            <div class="success-banner">
              <span class="material-symbols-outlined">check_circle</span>
              Pertanyaan Anda telah terkirim. Tim ustaz kami akan menjawab dalam 1×24 jam.
            </div>
          {/if}
          <form onsubmit={submitQuestion} class="submit-form">
            <label class="field-label" for="q-input">Pertanyaan Anda</label>
            <textarea
              id="q-input"
              rows="4"
              placeholder="Tulis pertanyaan seputar ibadah umroh atau haji..."
              bind:value={question}
              required
            ></textarea>
            <div class="field-hint">Pertanyaan akan dijawab oleh tim pembimbing berpengalaman. Jawaban tersedia dalam 1×24 jam.</div>
            <button type="submit" class="submit-btn" disabled={submitting || !question.trim()}>
              {#if submitting}
                <span class="material-symbols-outlined spin">progress_activity</span>
                Mengirim...
              {:else}
                <span class="material-symbols-outlined">send</span>
                Kirim Pertanyaan
              {/if}
            </button>
          </form>
        </div>
      </div>

      <!-- Recent answers -->
      <h2 class="section-title">Pertanyaan yang Baru Dijawab</h2>
      <div class="answers-list">
        {#each recentAnswers as item (item.id)}
          <div class="answer-card">
            <div class="answer-header">
              <div class="answer-q">
                <span class="q-icon">Q</span>
                {item.q}
              </div>
              <span class="answer-date">{item.date}</span>
            </div>
            {#if item.answered}
              <div class="answer-body">
                <span class="a-icon">A</span>
                <p>{item.a}</p>
              </div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .tanya-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 64rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 2rem; }
  .page-header h1 { margin: 0; font-size: 1.9rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .page-header p { margin: 0.4rem 0 0; color: #6b7280; }
  .section-title { margin: 0 0 1.2rem; font-size: 1.1rem; font-weight: 700; color: #1b1c1c; font-family: 'Plus Jakarta Sans', sans-serif; }
  .faq-list { display: flex; flex-direction: column; gap: 0.65rem; margin-bottom: 2.5rem; }
  .faq-item { background: #fff; border-radius: 1rem; border: 1px solid rgba(190,201,193,0.2); overflow: hidden; }
  .faq-q { display: flex; align-items: center; justify-content: space-between; gap: 1rem; padding: 1.1rem 1.3rem; background: none; border: none; cursor: pointer; font-size: 0.9rem; font-weight: 600; color: #1b1c1c; text-align: left; width: 100%; font-family: inherit; }
  .faq-q:hover { background: #fbf9f8; }
  .faq-arrow { color: #9ca3af; font-size: 1.2rem; transition: transform 0.2s; flex-shrink: 0; }
  .faq-arrow.open { transform: rotate(180deg); color: #006747; }
  .faq-a { padding: 0 1.3rem 1.1rem; font-size: 0.88rem; color: #57534e; line-height: 1.7; }
  /* Submit */
  .submit-section { margin-bottom: 2.5rem; }
  .submit-card { background: #fff; border-radius: 1.5rem; border: 1px solid rgba(190,201,193,0.2); padding: 1.75rem; }
  .success-banner { display: flex; align-items: center; gap: 0.5rem; background: rgba(0,103,71,0.08); color: #006747; border-radius: 0.75rem; padding: 0.75rem 1rem; font-size: 0.85rem; font-weight: 600; margin-bottom: 1.25rem; }
  .success-banner .material-symbols-outlined { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
  .submit-form { display: flex; flex-direction: column; gap: 0.85rem; }
  .field-label { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.07em; color: #6b7280; }
  .submit-form textarea {
    border: 1px solid rgba(190,201,193,0.4);
    border-radius: 0.75rem;
    padding: 0.75rem 1rem;
    font-size: 0.88rem;
    color: #1b1c1c;
    background: #fbf9f8;
    font-family: inherit;
    resize: vertical;
  }
  .field-hint { font-size: 0.78rem; color: #9ca3af; }
  .submit-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: linear-gradient(135deg, #004d34, #006747);
    color: #fff;
    border: none;
    border-radius: 999px;
    padding: 0.8rem 1.8rem;
    font-size: 0.9rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    align-self: flex-start;
  }
  .submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
  /* Answers */
  .answers-list { display: flex; flex-direction: column; gap: 1rem; }
  .answer-card { background: #fff; border-radius: 1.2rem; border: 1px solid rgba(190,201,193,0.2); padding: 1.25rem 1.5rem; }
  .answer-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; margin-bottom: 0.75rem; }
  .answer-q { display: flex; align-items: flex-start; gap: 0.75rem; font-size: 0.9rem; font-weight: 600; color: #1b1c1c; flex: 1; }
  .q-icon, .a-icon {
    width: 1.6rem;
    height: 1.6rem;
    border-radius: 0.4rem;
    display: grid;
    place-items: center;
    font-size: 0.72rem;
    font-weight: 800;
    flex-shrink: 0;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .q-icon { background: rgba(0,103,71,0.08); color: #006747; }
  .a-icon { background: rgba(119,90,25,0.08); color: #775a19; }
  .answer-date { font-size: 0.72rem; color: #9ca3af; flex-shrink: 0; }
  .answer-body { display: flex; gap: 0.75rem; align-items: flex-start; }
  .answer-body p { margin: 0; font-size: 0.85rem; color: #57534e; line-height: 1.65; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .spin { animation: spin 0.8s linear infinite; }
</style>
