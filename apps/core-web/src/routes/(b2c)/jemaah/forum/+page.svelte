<script lang="ts">
  import { MarketingPageLayout } from '$lib/components/marketing';

  interface Post {
    id: number;
    author: string;
    initials: string;
    avatarColor: string;
    message: string;
    timestamp: string;
  }

  let posts = $state<Post[]>([
    { id: 1, author: 'Ustaz Ahmad Fauzi', initials: 'AF', avatarColor: '#006747', message: 'Assalamu\'alaikum warahmatullahi wabarakatuh. Alhamdulillah kita semua sudah tiba dengan selamat. Ingat shalat subuh berjamaah di masjid ya!', timestamp: '05:45' },
    { id: 2, author: 'Bambang Suryanto', initials: 'BS', avatarColor: '#1565c0', message: 'Wa\'alaikumsalam. Siap pak ustaz! Apakah setelah subuh langsung thawaf atau ada briefing dulu?', timestamp: '05:50' },
    { id: 3, author: 'Siti Rahayu', initials: 'SR', avatarColor: '#775a19', message: 'Ada yang tahu klinik kesehatan jamaah Indonesia di mana ya? Suami saya agak kurang enak badan.', timestamp: '07:12' },
    { id: 4, author: 'Hendra Wijaya', initials: 'HW', avatarColor: '#c62828', message: 'Bu Siti, klinik ada di lantai B1 hotel, dekat lift. Buka 24 jam. Semoga cepat sembuh pak suaminya!', timestamp: '07:15' },
    { id: 5, author: 'Dewi Lestari', initials: 'DL', avatarColor: '#6a1c6a', message: 'Teman-teman jangan lupa bawa kartu tanda jamaah saat ke masjid ya. Kemarin ada yang ketinggalan dan susah masuk.', timestamp: '08:30' },
  ]);

  let newMessage = $state('');
  let posting = $state(false);

  const myName = 'Bambang Suryanto';
  const myInitials = 'BS';
  const myColor = '#1565c0';
  const groupName = 'Rombongan Umroh Januari 2025 — Grup A';
  const memberCount = 38;

  function submitPost() {
    if (!newMessage.trim()) return;
    posting = true;
    setTimeout(() => {
      const now = new Date();
      posts.push({
        id: Date.now(),
        author: myName,
        initials: myInitials,
        avatarColor: myColor,
        message: newMessage.trim(),
        timestamp: `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`,
      });
      newMessage = '';
      posting = false;
    }, 500);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) submitPost();
  }
</script>

<svelte:head>
  <title>Forum Rombongan — UmrohOS</title>
</svelte:head>

<MarketingPageLayout ctaHref="/packages" packagesLinkActive={false}>
  <div class="forum-root">
    <div class="shell">
      <a href="/jemaah" class="back-link">
        <span class="material-symbols-outlined">arrow_back</span>
        Portal Jamaah
      </a>
      <div class="page-header">
        <div class="group-info">
          <div class="group-icon">
            <span class="material-symbols-outlined">groups</span>
          </div>
          <div>
            <h1>{groupName}</h1>
            <div class="member-count">
              <span class="material-symbols-outlined">person</span>
              {memberCount} anggota rombongan
            </div>
          </div>
        </div>
      </div>

      <!-- Posts feed -->
      <div class="feed-card">
        <div class="posts-feed">
          {#each posts as post (post.id)}
            <div class="post-item" class:mine={post.author === myName}>
              {#if post.author !== myName}
                <div class="avatar" style="background: {post.avatarColor}">{post.initials}</div>
              {/if}
              <div class="post-bubble-wrap">
                {#if post.author !== myName}
                  <div class="post-author">{post.author}</div>
                {/if}
                <div class="post-bubble" class:mine={post.author === myName}>
                  {post.message}
                </div>
                <div class="post-time" class:mine={post.author === myName}>{post.timestamp}</div>
              </div>
            </div>
          {/each}
        </div>

        <!-- Composer -->
        <div class="composer">
          <div class="avatar composer-avatar" style="background: {myColor}">{myInitials}</div>
          <div class="composer-input-wrap">
            <textarea
              rows="2"
              placeholder="Tulis pesan... (Ctrl+Enter untuk kirim)"
              bind:value={newMessage}
              onkeydown={handleKeydown}
            ></textarea>
          </div>
          <button class="send-btn" onclick={submitPost} disabled={posting || !newMessage.trim()}>
            <span class="material-symbols-outlined">send</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</MarketingPageLayout>

<style>
  .forum-root { padding-top: calc(5.2rem + 2rem); padding-bottom: 5rem; background: #fbf9f8; min-height: 100vh; }
  .shell { max-width: 64rem; margin: 0 auto; padding: 0 1.5rem; }
  .back-link { display: inline-flex; align-items: center; gap: 0.35rem; color: #006747; font-weight: 600; font-size: 0.85rem; text-decoration: none; margin-bottom: 0.75rem; }
  .back-link .material-symbols-outlined { font-size: 1rem; }
  .page-header { margin-bottom: 2rem; }
  .group-info { display: flex; align-items: center; gap: 1rem; }
  .group-icon { width: 3.5rem; height: 3.5rem; border-radius: 1rem; background: rgba(0,103,71,0.1); display: grid; place-items: center; color: #006747; flex-shrink: 0; }
  .group-icon .material-symbols-outlined { font-size: 2rem; font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 48; }
  .group-info h1 { margin: 0; font-size: 1.25rem; font-weight: 800; color: #004d34; font-family: 'Plus Jakarta Sans', sans-serif; }
  .member-count { display: flex; align-items: center; gap: 0.3rem; font-size: 0.82rem; color: #9ca3af; margin-top: 0.25rem; }
  .member-count .material-symbols-outlined { font-size: 0.9rem; }
  .feed-card { background: #fff; border-radius: 1.5rem; border: 1px solid rgba(190,201,193,0.2); overflow: hidden; }
  .posts-feed { padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; min-height: 400px; }
  .post-item { display: flex; gap: 0.75rem; align-items: flex-end; }
  .post-item.mine { flex-direction: row-reverse; }
  .avatar {
    width: 2.2rem;
    height: 2.2rem;
    border-radius: 50%;
    display: grid;
    place-items: center;
    font-size: 0.75rem;
    font-weight: 700;
    color: #fff;
    flex-shrink: 0;
  }
  .post-bubble-wrap { max-width: 70%; display: flex; flex-direction: column; gap: 0.15rem; }
  .post-author { font-size: 0.72rem; font-weight: 700; color: #9ca3af; margin-bottom: 0.1rem; }
  .post-bubble {
    background: #f2f4f6;
    border-radius: 1rem 1rem 1rem 0.25rem;
    padding: 0.75rem 1rem;
    font-size: 0.88rem;
    color: #1b1c1c;
    line-height: 1.5;
  }
  .post-bubble.mine { background: rgba(0,103,71,0.1); border-radius: 1rem 1rem 0.25rem 1rem; color: #004d34; }
  .post-time { font-size: 0.68rem; color: #9ca3af; }
  .post-time.mine { text-align: right; }
  .composer {
    display: flex;
    align-items: flex-end;
    gap: 0.75rem;
    padding: 1rem 1.5rem;
    border-top: 1px solid rgba(190,201,193,0.2);
    background: #fbf9f8;
  }
  .composer-avatar { flex-shrink: 0; }
  .composer-input-wrap { flex: 1; }
  .composer textarea {
    width: 100%;
    border: 1px solid rgba(190,201,193,0.3);
    border-radius: 0.75rem;
    padding: 0.65rem 0.85rem;
    font-size: 0.88rem;
    font-family: inherit;
    background: #fff;
    resize: none;
    color: #1b1c1c;
  }
  .send-btn {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    background: #006747;
    color: #fff;
    border: none;
    cursor: pointer;
    display: grid;
    place-items: center;
    flex-shrink: 0;
    transition: background 0.15s;
  }
  .send-btn:hover { background: #004d34; }
  .send-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .send-btn .material-symbols-outlined { font-size: 1.1rem; }
</style>
