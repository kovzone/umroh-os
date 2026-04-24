<script lang="ts">
  const CATEGORIES = [
    {
      title: 'Pengadaan',
      icon: 'shopping_cart',
      color: '#2563eb',
      items: [
        { label: 'Permintaan Pembelian (PR)', href: '/console/logistics/purchase-requests', desc: 'PR, tiered approvals, budget sync' },
        { label: 'Penerimaan Barang (GRN)', href: '/console/logistics/grn', desc: 'GRN parsial, QC, reversal' },
        { label: 'Seleksi Vendor Otomatis', href: '/console/logistics/auto-vendor', desc: 'Auto-select vendor terbaik' },
      ]
    },
    {
      title: 'Inventori',
      icon: 'inventory_2',
      color: '#7c3aed',
      items: [
        { label: 'Inventori & Gudang', href: '/console/logistics/inventory', desc: 'Stock alerts, reorder, multi-gudang, barcode' },
        { label: 'Stok Opname Digital', href: '/console/logistics/stocktake', desc: 'Hitung fisik stok & laporan variance' },
      ]
    },
    {
      title: 'Fulfillment & Pengiriman',
      icon: 'local_shipping',
      color: '#059669',
      items: [
        { label: 'Fulfillment & Pengiriman', href: '/console/logistics/fulfillment', desc: 'Sinkronisasi ukuran, courier tracking, retur & tukar' },
      ]
    },
  ];
</script>

<main class="page-shell">
  <header class="topbar">
    <nav class="breadcrumb">
      <a href="/console" class="bc-link">Console</a>
      <span class="bc-sep">/</span>
      <span>Logistik</span>
    </nav>
  </header>

  <section class="canvas">
    <div class="page-head">
      <h2>Pusat Logistik</h2>
      <p>Kelola pengadaan, inventori, fulfillment, vendor, dan retur barang</p>
    </div>

    {#each CATEGORIES as cat}
      <div class="cat-block">
        <div class="cat-header">
          <span class="material-symbols-outlined cat-icon" style="color:{cat.color}">{cat.icon}</span>
          <h3>{cat.title}</h3>
        </div>
        <div class="card-grid">
          {#each cat.items as item}
            <a href={item.href} class="nav-card">
              <div class="nav-card-label">{item.label}</div>
              <div class="nav-card-desc">{item.desc}</div>
              <span class="material-symbols-outlined nav-arrow">arrow_forward</span>
            </a>
          {/each}
        </div>
      </div>
    {/each}
  </section>
</main>

<style>
  .page-shell { min-height: 100vh; background: #f7f9fb; }
  .topbar { position: sticky; top: 0; z-index: 30; height: 4rem; background: rgb(255 255 255 / 0.9); border-bottom: 1px solid rgb(195 198 215 / 0.45); padding: 0 1.25rem; display: flex; align-items: center; backdrop-filter: blur(8px); }
  .breadcrumb { display: flex; align-items: center; gap: 0.4rem; font-size: 0.82rem; color: #737686; }
  .bc-link { color: #2563eb; text-decoration: none; font-weight: 600; }
  .bc-sep { color: #b0b3c1; }
  .canvas { padding: 1.5rem; max-width: 72rem; }
  .page-head { margin-bottom: 1.5rem; }
  .page-head h2 { margin: 0; font-size: 1.5rem; font-weight: 700; }
  .page-head p { margin: 0.25rem 0 0; font-size: 0.82rem; color: #737686; }
  .cat-block { margin-bottom: 2rem; }
  .cat-header { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 0.85rem; }
  .cat-header h3 { margin: 0; font-size: 0.95rem; font-weight: 700; color: #191c1e; }
  .cat-icon { font-size: 1.25rem; }
  .card-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 0.75rem; }
  .nav-card { display: flex; flex-direction: column; gap: 0.25rem; padding: 1rem 1.1rem; background: #fff; border: 1px solid rgb(195 198 215 / 0.45); border-radius: 0.5rem; text-decoration: none; color: inherit; position: relative; transition: box-shadow 0.15s, border-color 0.15s; }
  .nav-card:hover { border-color: #93c5fd; box-shadow: 0 2px 8px rgb(0 0 0 / 0.07); }
  .nav-card-label { font-size: 0.85rem; font-weight: 700; color: #191c1e; padding-right: 1.5rem; }
  .nav-card-desc { font-size: 0.72rem; color: #737686; }
  .nav-arrow { position: absolute; right: 0.85rem; top: 50%; transform: translateY(-50%); font-size: 1rem; color: #b0b3c1; }
  .nav-card:hover .nav-arrow { color: #2563eb; }
</style>
