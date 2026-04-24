<script lang="ts">
  let open = $state(false);
  let message = $state('');

  const quickMessages = [
    'Saya ingin mengetahui harga paket umrah terbaru',
    'Bagaimana cara mendaftar dan membayar cicilan?',
    'Apa saja dokumen yang diperlukan untuk umrah?',
    'Apakah ada jadwal keberangkatan bulan ini?'
  ];

  function sendMessage(text: string = message) {
    if (!text.trim()) return;
    const encoded = encodeURIComponent(text.trim());
    window.open(`https://wa.me/6281200000000?text=${encoded}`, '_blank');
    message = '';
    open = false;
  }

  function toggleOpen() {
    open = !open;
  }
</script>

<div class="floating-chat-wrap" class:is-open={open}>
  <!-- Bubble -->
  {#if open}
    <div class="chat-bubble" role="dialog" aria-label="Chat dengan kami">
      <div class="bubble-header">
        <div class="agent-info">
          <div class="agent-avatar">
            <span class="material-symbols-outlined">support_agent</span>
          </div>
          <div>
            <strong>Tim UmrohOS</strong>
            <p>Biasanya membalas dalam beberapa menit</p>
          </div>
        </div>
        <button class="close-btn" type="button" onclick={toggleOpen} aria-label="Tutup chat">
          <span class="material-symbols-outlined">close</span>
        </button>
      </div>

      <div class="bubble-body">
        <div class="greeting-msg">
          <div class="bot-bubble">
            <p>Assalamu'alaikum! 👋 Ada yang bisa kami bantu?</p>
          </div>
        </div>

        <div class="quick-options">
          <p class="quick-label">Pertanyaan populer:</p>
          {#each quickMessages as qm (qm)}
            <button type="button" class="quick-btn" onclick={() => sendMessage(qm)}>
              {qm}
            </button>
          {/each}
        </div>
      </div>

      <div class="bubble-footer">
        <div class="input-row">
          <input
            type="text"
            placeholder="Ketik pesan Anda..."
            bind:value={message}
            onkeydown={(e) => e.key === 'Enter' && sendMessage()}
          />
          <button
            class="send-btn"
            type="button"
            onclick={() => sendMessage()}
            disabled={!message.trim()}
          >
            <span class="material-symbols-outlined">send</span>
          </button>
        </div>
        <p class="wa-note">
          <span class="material-symbols-outlined">open_in_new</span>
          Pesan dikirim via WhatsApp
        </p>
      </div>
    </div>
  {/if}

  <!-- Trigger Button -->
  <button class="chat-trigger" type="button" onclick={toggleOpen} aria-label="Buka chat">
    {#if !open}
      <span class="material-symbols-outlined trigger-icon">chat</span>
      <span class="trigger-label">Tanya Admin</span>
    {:else}
      <span class="material-symbols-outlined trigger-icon">close</span>
    {/if}
    <span class="notification-dot"></span>
  </button>
</div>

<style>
  .floating-chat-wrap {
    position: fixed;
    bottom: 1.8rem;
    right: 1.8rem;
    z-index: 100;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 0.75rem;
  }
  /* Trigger */
  .chat-trigger {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    border: none;
    background: linear-gradient(135deg, #004d34, #006747);
    color: #fff;
    border-radius: 999px;
    padding: 0.85rem 1.4rem;
    font-size: 0.95rem;
    font-weight: 700;
    cursor: pointer;
    box-shadow: 0 8px 24px rgba(0, 103, 71, 0.4);
    position: relative;
    transition: transform 0.2s, box-shadow 0.2s;
    font-family: 'Plus Jakarta Sans', sans-serif;
  }
  .chat-trigger:hover {
    transform: scale(1.04);
    box-shadow: 0 12px 30px rgba(0, 103, 71, 0.5);
  }
  .trigger-icon {
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
    font-size: 1.3rem;
  }
  .notification-dot {
    position: absolute;
    top: 0.3rem;
    right: 0.3rem;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: #fed488;
    border: 2px solid #004d34;
    animation: pulse 2s infinite;
  }
  @keyframes pulse {
    0%, 100% { transform: scale(1); opacity: 1; }
    50% { transform: scale(1.3); opacity: 0.7; }
  }
  .is-open .notification-dot { display: none; }
  /* Chat bubble */
  .chat-bubble {
    width: 22rem;
    background: #fff;
    border-radius: 1.5rem;
    box-shadow: 0 20px 48px rgba(0, 0, 0, 0.18);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.2s ease;
  }
  @keyframes slideUp {
    from { opacity: 0; transform: translateY(12px); }
    to { opacity: 1; transform: translateY(0); }
  }
  /* Header */
  .bubble-header {
    background: linear-gradient(135deg, #004d34, #006747);
    color: #fff;
    padding: 1.1rem 1.3rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .agent-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  .agent-avatar {
    width: 2.8rem;
    height: 2.8rem;
    border-radius: 50%;
    background: rgba(255,255,255,0.2);
    display: grid;
    place-items: center;
    flex-shrink: 0;
  }
  .agent-avatar .material-symbols-outlined {
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
    font-size: 1.4rem;
  }
  .agent-info strong { display: block; font-size: 0.95rem; font-family: 'Plus Jakarta Sans', sans-serif; }
  .agent-info p { margin: 0.15rem 0 0; font-size: 0.76rem; opacity: 0.8; }
  .close-btn {
    border: none;
    background: rgba(255,255,255,0.15);
    color: #fff;
    border-radius: 50%;
    width: 2rem;
    height: 2rem;
    display: grid;
    place-items: center;
    cursor: pointer;
    transition: background 0.15s;
  }
  .close-btn:hover { background: rgba(255,255,255,0.25); }
  .close-btn .material-symbols-outlined { font-size: 1rem; }
  /* Body */
  .bubble-body {
    padding: 1.2rem;
    flex: 1;
    max-height: 20rem;
    overflow-y: auto;
  }
  .greeting-msg { margin-bottom: 1rem; }
  .bot-bubble {
    background: #f0f9f4;
    border-radius: 1rem 1rem 1rem 0.3rem;
    padding: 0.8rem 1rem;
    display: inline-block;
    max-width: 85%;
  }
  .bot-bubble p {
    margin: 0;
    font-size: 0.88rem;
    color: #1b1c1c;
    line-height: 1.55;
  }
  .quick-label {
    margin: 0 0 0.6rem;
    font-size: 0.74rem;
    font-weight: 700;
    color: #9ca3af;
    text-transform: uppercase;
    letter-spacing: 0.07em;
  }
  .quick-options { display: grid; gap: 0.4rem; }
  .quick-btn {
    text-align: left;
    border: 1px solid rgba(190,201,193,0.4);
    border-radius: 0.75rem;
    background: #fff;
    padding: 0.6rem 0.8rem;
    font-size: 0.82rem;
    color: #1b1c1c;
    cursor: pointer;
    transition: background 0.15s, border-color 0.15s;
    line-height: 1.35;
  }
  .quick-btn:hover { background: #f0f9f4; border-color: #006747; color: #004d34; }
  /* Footer */
  .bubble-footer {
    padding: 0.9rem 1.2rem;
    border-top: 1px solid rgba(190,201,193,0.2);
    background: #fafaf9;
  }
  .input-row {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }
  .input-row input {
    flex: 1;
    border: 1px solid rgba(190,201,193,0.4);
    border-radius: 999px;
    padding: 0.6rem 1rem;
    font-size: 0.85rem;
    outline: none;
    color: #1b1c1c;
    transition: border-color 0.15s;
  }
  .input-row input:focus { border-color: #006747; }
  .send-btn {
    width: 2.4rem;
    height: 2.4rem;
    border-radius: 50%;
    border: none;
    background: #006747;
    color: #fff;
    display: grid;
    place-items: center;
    cursor: pointer;
    flex-shrink: 0;
    transition: opacity 0.15s;
  }
  .send-btn:disabled { opacity: 0.4; cursor: not-allowed; }
  .send-btn .material-symbols-outlined {
    font-size: 1.1rem;
    font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
  .wa-note {
    margin: 0.5rem 0 0;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.72rem;
    color: #9ca3af;
    justify-content: center;
  }
  .wa-note .material-symbols-outlined { font-size: 0.8rem; }
  @media (max-width: 480px) {
    .floating-chat-wrap {
      bottom: 1rem;
      right: 1rem;
    }
    .chat-bubble { width: calc(100vw - 2rem); }
    .trigger-label { display: none; }
    .chat-trigger { padding: 0.9rem; }
  }
</style>
