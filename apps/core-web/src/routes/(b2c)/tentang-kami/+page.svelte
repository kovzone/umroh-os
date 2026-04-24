<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  const milestones = [
    { year: '2010', label: 'Berdiri sebagai agen perjalanan umrah pertama berbasis digital di Indonesia' },
    { year: '2014', label: 'Mendapatkan izin resmi PPIU dari Kementerian Agama RI (No. 123/2024)' },
    { year: '2017', label: 'Melayani 5.000 jamaah pertama dengan tingkat kepuasan 97%' },
    { year: '2020', label: 'Meluncurkan platform UmrohOS untuk manajemen perjalanan berbasis data' },
    { year: '2022', label: 'Dipercaya melayani 15.000+ jamaah dari seluruh Indonesia' },
    { year: '2024', label: 'Lebih dari 25.000 jamaah telah berangkat — inovasi terus berlanjut' }
  ];

  const values = [
    { icon: 'verified_user', title: 'Amanah & Transparan', desc: 'Kami menjaga kepercayaan jamaah dengan harga yang jelas, tanpa biaya tersembunyi, dan laporan perjalanan yang dapat dipantau secara real-time.' },
    { icon: 'mosque', title: 'Bimbingan Spiritual', desc: 'Setiap paket dilengkapi program manasik komprehensif dan pendampingan muthawwif berpengalaman dari Kementerian Agama.' },
    { icon: 'support_agent', title: 'Layanan 24/7', desc: 'Tim kami siap melayani kebutuhan jamaah dari proses pendaftaran, pengurusan dokumen, hingga kepulangan ke tanah air.' },
    { icon: 'workspace_premium', title: 'Standar Kualitas Tinggi', desc: 'Bermitra dengan maskapai kelas dunia dan hotel berbintang yang menjamin kenyamanan perjalanan ibadah Anda.' }
  ];

  const team = [
    { name: 'H. Ahmad Fauzi', role: 'CEO & Co-Founder', photo: 'https://images.unsplash.com/photo-1560250097-0b93528c311a?w=400&q=80' },
    { name: 'Siti Rahmawati', role: 'Direktur Operasional', photo: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=400&q=80' },
    { name: 'Ustadz Wahyu Hidayat', role: 'Kepala Bimbingan Ibadah', photo: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&q=80' },
    { name: 'Dewi Santika', role: 'Manajer Pengalaman Jamaah', photo: 'https://images.unsplash.com/photo-1580489944761-15a19d654956?w=400&q=80' }
  ];

  const certifications = [
    { title: 'Izin PPIU', number: 'No. 123/2024', issuer: 'Kementerian Agama RI', icon: 'verified_user' },
    { title: 'Akreditasi A', number: 'Penilaian Tertinggi', issuer: 'Lembaga Akreditasi Haji & Umrah', icon: 'workspace_premium' },
    { title: 'Anggota HIMPUH', number: 'No. H/456/2020', issuer: 'Himpunan Penyelenggara Umrah & Haji', icon: 'groups' },
    { title: 'ISO 9001:2015', number: 'Sertifikat No. 789', issuer: 'Sistem Manajemen Mutu', icon: 'military_tech' }
  ];
</script>

<svelte:head>
  <title>Tentang Kami — UmrohOS</title>
  <meta name="description" content="UmrohOS adalah penyelenggara perjalanan Umrah digital terpercaya dengan izin resmi PPIU Kemenag RI, melayani lebih dari 25.000 jamaah sejak 2010." />
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="about-root">

    <!-- Hero -->
    <section class="hero">
      <div class="shell">
        <div class="hero-grid">
          <div class="hero-copy">
            <p class="kicker">Tentang UmrohOS</p>
            <h1>Mitra Terpercaya Perjalanan Spiritual Anda</h1>
            <p>
              Sejak 2010, UmrohOS telah menjadi jembatan bagi lebih dari 25.000 jamaah Indonesia
              dalam mewujudkan impian ibadah Umrah dan Haji yang nyaman, aman, dan penuh makna.
              Kami hadir dengan teknologi modern namun tetap menjunjung tinggi nilai-nilai keislaman.
            </p>
            <div class="hero-stats">
              <div class="stat">
                <strong>25.000+</strong>
                <span>Jamaah Berangkat</span>
              </div>
              <div class="stat">
                <strong>14 Tahun</strong>
                <span>Pengalaman</span>
              </div>
              <div class="stat">
                <strong>97%</strong>
                <span>Tingkat Kepuasan</span>
              </div>
            </div>
          </div>
          <div class="hero-img">
            <img src="https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&q=80" alt="Masjidil Haram" loading="eager" />
            <div class="license-badge">
              <span class="material-symbols-outlined">verified_user</span>
              <div>
                <strong>Izin Resmi PPIU</strong>
                <p>No. 123/2024 — Kemenag RI</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Values -->
    <section class="values">
      <div class="shell">
        <div class="section-head">
          <h2>Nilai-Nilai Kami</h2>
          <p>Fondasi yang membangun kepercayaan lebih dari satu dekade</p>
        </div>
        <div class="values-grid">
          {#each values as v (v.title)}
            <article class="value-card">
              <div class="value-icon">
                <span class="material-symbols-outlined">{v.icon}</span>
              </div>
              <h3>{v.title}</h3>
              <p>{v.desc}</p>
            </article>
          {/each}
        </div>
      </div>
    </section>

    <!-- Milestones -->
    <section class="milestones">
      <div class="shell">
        <div class="section-head light">
          <h2>Perjalanan Kami</h2>
          <p>Lebih dari satu dekade melayani umat dengan dedikasi penuh</p>
        </div>
        <div class="timeline">
          {#each milestones as m, idx (m.year)}
            <div class="timeline-item" class:right={idx % 2 === 1}>
              <div class="timeline-year">{m.year}</div>
              <div class="timeline-content">
                <p>{m.label}</p>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </section>

    <!-- Team -->
    <section class="team">
      <div class="shell">
        <div class="section-head">
          <h2>Tim Kami</h2>
          <p>Profesional berpengalaman yang berdedikasi untuk perjalanan ibadah Anda</p>
        </div>
        <div class="team-grid">
          {#each team as member (member.name)}
            <article class="team-card">
              <div class="team-photo">
                <img src={member.photo} alt={member.name} loading="lazy" />
              </div>
              <div class="team-info">
                <h3>{member.name}</h3>
                <p>{member.role}</p>
              </div>
            </article>
          {/each}
        </div>
      </div>
    </section>

    <!-- Certifications -->
    <section class="certs">
      <div class="shell">
        <div class="section-head">
          <h2>Legalitas & Sertifikasi</h2>
          <p>Terdaftar dan diakui secara resmi oleh lembaga berwenang</p>
        </div>
        <div class="certs-grid">
          {#each certifications as cert (cert.title)}
            <article class="cert-card">
              <div class="cert-icon">
                <span class="material-symbols-outlined">{cert.icon}</span>
              </div>
              <div>
                <h3>{cert.title}</h3>
                <p class="cert-number">{cert.number}</p>
                <p class="cert-issuer">{cert.issuer}</p>
              </div>
            </article>
          {/each}
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="cta-section">
      <div class="shell cta-inner">
        <div class="cta-copy">
          <h2>Siap Memulai Perjalanan Ibadah?</h2>
          <p>Bergabunglah dengan lebih dari 25.000 jamaah yang telah mempercayakan perjalanan umrah mereka kepada kami.</p>
        </div>
        <div class="cta-btns">
          <a class="btn-primary" href="/packages">Lihat Paket Umrah</a>
          <a class="btn-outline" href="https://wa.me/6281200000000" target="_blank" rel="noreferrer">
            <span class="material-symbols-outlined">chat</span>
            Konsultasi Gratis
          </a>
        </div>
      </div>
    </section>

  </div>
</MarketingPageLayout>

<style>
  .about-root {
    padding-top: 5.2rem;
    background: #fbf9f8;
  }
  .shell {
    max-width: 80rem;
    margin: 0 auto;
    padding: 0 1.5rem;
  }
  /* Hero */
  .hero {
    padding: 5rem 0 4rem;
  }
  .hero-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 4rem;
    align-items: center;
  }
  .kicker {
    display: inline-block;
    margin: 0 0 1rem;
    background: rgba(254, 212, 136, 0.3);
    color: #775a19;
    border-radius: 999px;
    padding: 0.4rem 1rem;
    font-size: 0.76rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .hero-copy h1 {
    margin: 0 0 1rem;
    font-size: clamp(2rem, 3.5vw, 3.2rem);
    line-height: 1.15;
    color: #004d34;
    letter-spacing: -0.02em;
    font-family: 'Plus Jakarta Sans', sans-serif;
    font-weight: 800;
  }
  .hero-copy p {
    margin: 0 0 2rem;
    color: #3f4943;
    font-size: 1.08rem;
    line-height: 1.75;
  }
  .hero-stats {
    display: flex;
    gap: 2.5rem;
    flex-wrap: wrap;
  }
  .stat strong {
    display: block;
    font-size: 2.2rem;
    font-weight: 800;
    color: #004d34;
    line-height: 1;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .stat span {
    font-size: 0.85rem;
    color: #6b7280;
    font-weight: 500;
  }
  .hero-img {
    position: relative;
    border-radius: 2rem;
    overflow: hidden;
    box-shadow: 0 24px 48px rgba(0,0,0,0.15);
  }
  .hero-img img {
    width: 100%;
    height: 32rem;
    object-fit: cover;
    display: block;
  }
  .license-badge {
    position: absolute;
    bottom: 1.5rem;
    left: 1.5rem;
    background: rgba(255,255,255,0.95);
    backdrop-filter: blur(8px);
    border-radius: 1rem;
    padding: 0.9rem 1.2rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    box-shadow: 0 8px 20px rgba(0,0,0,0.12);
  }
  .license-badge .material-symbols-outlined {
    color: #006747;
    font-size: 2rem;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .license-badge strong {
    display: block;
    font-size: 0.92rem;
    font-weight: 700;
    color: #004d34;
  }
  .license-badge p {
    margin: 0.2rem 0 0;
    font-size: 0.78rem;
    color: #6b7280;
  }
  /* Section head */
  .section-head {
    text-align: center;
    margin-bottom: 3.5rem;
  }
  .section-head h2 {
    margin: 0;
    font-size: clamp(1.8rem, 3.5vw, 2.8rem);
    font-weight: 800;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .section-head p {
    margin: 0.8rem auto 0;
    max-width: 38rem;
    color: #6b7280;
    font-size: 1.05rem;
  }
  .section-head.light h2 { color: #fff; }
  .section-head.light p { color: rgba(255,255,255,0.75); }
  /* Values */
  .values {
    padding: 5rem 0;
    background: #fff;
  }
  .values-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 2rem;
  }
  .value-card {
    padding: 2rem;
    border-radius: 1.5rem;
    border: 1px solid rgba(190,201,193,0.3);
    background: #fbf9f8;
    transition: box-shadow 0.2s;
  }
  .value-card:hover {
    box-shadow: 0 12px 28px rgba(0,0,0,0.07);
  }
  .value-icon {
    width: 3.5rem;
    height: 3.5rem;
    border-radius: 1rem;
    background: rgba(0,103,71,0.1);
    display: grid;
    place-items: center;
    margin-bottom: 1.2rem;
    color: #006747;
  }
  .value-icon .material-symbols-outlined {
    font-size: 1.5rem;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .value-card h3 {
    margin: 0 0 0.6rem;
    font-size: 1.15rem;
    font-weight: 700;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .value-card p {
    margin: 0;
    color: #57534e;
    font-size: 0.92rem;
    line-height: 1.65;
  }
  /* Milestones */
  .milestones {
    padding: 5rem 0;
    background: linear-gradient(135deg, #004d34 0%, #006747 100%);
  }
  .timeline {
    display: grid;
    gap: 2rem;
    max-width: 52rem;
    margin: 0 auto;
  }
  .timeline-item {
    display: grid;
    grid-template-columns: 6rem 1fr;
    gap: 1.5rem;
    align-items: start;
  }
  .timeline-year {
    background: rgba(255,255,255,0.15);
    border-radius: 0.75rem;
    padding: 0.7rem 1rem;
    text-align: center;
    font-size: 1.1rem;
    font-weight: 800;
    color: #fed488;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .timeline-content {
    background: rgba(255,255,255,0.1);
    border-radius: 1rem;
    padding: 1.1rem 1.3rem;
  }
  .timeline-content p {
    margin: 0;
    color: rgba(255,255,255,0.9);
    line-height: 1.6;
  }
  /* Team */
  .team {
    padding: 5rem 0;
    background: #fbf9f8;
  }
  .team-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 2rem;
  }
  .team-card {
    border-radius: 1.5rem;
    overflow: hidden;
    border: 1px solid rgba(190,201,193,0.3);
    background: #fff;
    box-shadow: 0 4px 12px rgba(0,0,0,0.05);
    transition: transform 0.2s, box-shadow 0.2s;
  }
  .team-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 16px 32px rgba(0,0,0,0.1);
  }
  .team-photo img {
    width: 100%;
    height: 14rem;
    object-fit: cover;
    display: block;
  }
  .team-info {
    padding: 1.3rem;
  }
  .team-info h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .team-info p {
    margin: 0.35rem 0 0;
    font-size: 0.84rem;
    color: #775a19;
  }
  /* Certs */
  .certs {
    padding: 5rem 0;
    background: #fff;
  }
  .certs-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 2rem;
  }
  .cert-card {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
    padding: 1.5rem;
    border-radius: 1.5rem;
    border: 1px solid rgba(190,201,193,0.3);
    background: #fbf9f8;
  }
  .cert-icon {
    width: 3rem;
    height: 3rem;
    border-radius: 0.75rem;
    background: rgba(119,90,25,0.1);
    display: grid;
    place-items: center;
    flex-shrink: 0;
    color: #775a19;
  }
  .cert-icon .material-symbols-outlined {
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .cert-card h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: #1b1c1c;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .cert-number {
    margin: 0.25rem 0 0.15rem;
    font-size: 0.82rem;
    font-weight: 600;
    color: #006747;
  }
  .cert-issuer {
    margin: 0;
    font-size: 0.78rem;
    color: #6b7280;
    line-height: 1.4;
  }
  /* CTA section */
  .cta-section {
    padding: 5rem 0;
    background: #f5f3f3;
  }
  .cta-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 2rem;
    flex-wrap: wrap;
  }
  .cta-copy h2 {
    margin: 0;
    font-size: clamp(1.6rem, 3vw, 2.4rem);
    font-weight: 800;
    color: #004d34;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .cta-copy p {
    margin: 0.7rem 0 0;
    color: #57534e;
    max-width: 38rem;
  }
  .cta-btns {
    display: flex;
    gap: 0.9rem;
    flex-wrap: wrap;
  }
  .btn-primary {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    background: linear-gradient(90deg, #004d34, #006747);
    color: #fff;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.8rem 1.8rem;
  }
  .btn-outline {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    border: 1.5px solid #006747;
    color: #006747;
    font-weight: 700;
    border-radius: 999px;
    padding: 0.8rem 1.8rem;
    background: transparent;
  }
  .btn-outline .material-symbols-outlined {
    font-size: 1rem;
  }
  @media (max-width: 1100px) {
    .hero-grid, .values-grid, .team-grid, .certs-grid {
      grid-template-columns: 1fr 1fr;
    }
  }
  @media (max-width: 640px) {
    .hero-grid, .values-grid, .team-grid, .certs-grid {
      grid-template-columns: 1fr;
    }
    .hero-img img {
      height: 22rem;
    }
    .cta-inner {
      flex-direction: column;
      align-items: flex-start;
    }
  }
</style>
