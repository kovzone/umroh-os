<footer class="footer">
  <div class="shell footer-wrap">
    <div class="footer-brand">
      <h4>UmrohOS</h4>
      <p>© 2024 UmrohOS. The Serene Path to Holy Land.</p>
    </div>
    <div class="footer-links">
      <a href="/">Kebijakan Privasi</a>
      <a href="/">Syarat & Ketentuan</a>
      <a href="/">Bantuan</a>
      <a href="/">Kontak</a>
    </div>
    <div class="footer-icons">
      <span><span class="material-symbols-outlined">language</span></span>
      <span><span class="material-symbols-outlined">mail</span></span>
    </div>
  </div>
</footer>

<style>
  .footer {
    border-radius: 2rem 2rem 0 0;
    background: #f5f3f3;
    padding: 3rem 0;
    font-size: 0.875rem;
  }
  .footer-wrap {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
  }
  .footer-brand h4 {
    margin: 0;
    color: #006747;
    font-weight: 800;
    font-size: 1.6rem;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .footer-brand p {
    margin: 0.45rem 0 0;
    color: #6b7280;
    font-size: 0.88rem;
  }
  .footer-links {
    display: flex;
    gap: 1.5rem;
    flex-wrap: wrap;
  }
  .footer-links a {
    color: #6b7280;
    text-decoration: none;
    font-size: 0.9rem;
  }
  .footer-icons {
    display: flex;
    gap: 0.6rem;
  }
  .footer-icons span {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 999px;
    background: #efeded;
    color: #004d34;
    display: grid;
    place-items: center;
  }
  .material-symbols-outlined {
    font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
    font-size: 1.2rem;
  }
  @media (max-width: 760px) {
    .footer-wrap {
      flex-direction: column;
      align-items: start;
    }
  }
</style>
