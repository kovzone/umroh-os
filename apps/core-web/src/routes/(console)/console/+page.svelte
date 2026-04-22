<main class="dashboard-shell">
  <header class="topbar">
    <div class="search-wrap">
      <span class="material-symbols-outlined search-icon">search</span>
      <input type="text" placeholder="Search operational data..." />
    </div>
    <div class="top-actions">
      <button class="branch-btn">
        <span class="material-symbols-outlined">domain</span>
        <span>Jakarta Pusat (HQ)</span>
        <span class="material-symbols-outlined">arrow_drop_down</span>
      </button>
      <div class="icon-actions">
        <button class="icon-btn notif">
          <span class="material-symbols-outlined">notifications</span>
          <span class="dot"></span>
        </button>
        <button class="icon-btn"><span class="material-symbols-outlined">settings</span></button>
        <button class="avatar" aria-label="Profile">A</button>
      </div>
    </div>
  </header>

  <section class="canvas">
    <div class="page-head">
      <div>
        <h2>Dashboard Operasional</h2>
        <p>Real-time overview of core operations</p>
      </div>
      <div class="head-actions">
        <button class="ghost-btn">Download Report</button>
        <button class="primary-btn">
          <span class="material-symbols-outlined">add</span>
          <span>New Incident</span>
        </button>
      </div>
    </div>

    <div class="kpi-grid">
      <article class="kpi-card">
        <div class="kpi-top">
          <span>Booking Hari Ini</span>
          <span class="material-symbols-outlined">confirmation_number</span>
        </div>
        <p class="kpi-number">142</p>
        <p class="kpi-hint good"><span class="material-symbols-outlined">trending_up</span>+12% vs yday</p>
      </article>
      <article class="kpi-card">
        <div class="kpi-top">
          <span>Pending Payment</span>
          <span class="material-symbols-outlined">payments</span>
        </div>
        <p class="kpi-number">Rp 2.4M</p>
        <p class="kpi-hint bad">
          <span class="material-symbols-outlined">warning</span>15 invoices overdue
        </p>
      </article>
      <article class="kpi-card">
        <div class="kpi-top">
          <span>Visa Overdue</span>
          <span class="material-symbols-outlined bad">description</span>
        </div>
        <p class="kpi-number">8</p>
        <p class="kpi-hint">Requires immediate action</p>
      </article>
      <article class="kpi-card">
        <div class="kpi-top">
          <span>Berangkat 7 Hari</span>
          <span class="material-symbols-outlined">flight_takeoff</span>
        </div>
        <p class="kpi-number">450</p>
        <p class="kpi-hint">Across 4 flights</p>
      </article>
      <article class="kpi-card">
        <div class="kpi-top">
          <span>Open Incident</span>
          <span class="material-symbols-outlined bad">report</span>
        </div>
        <p class="kpi-number">3</p>
        <p class="kpi-hint bad">1 High Priority</p>
      </article>
      <article class="kpi-card">
        <div class="kpi-top">
          <span>SLA CS</span>
          <span class="material-symbols-outlined">support_agent</span>
        </div>
        <p class="kpi-number">98.2%</p>
        <p class="kpi-hint good">
          <span class="material-symbols-outlined">check_circle</span>Meeting target
        </p>
      </article>
    </div>

    <div class="main-grid">
      <div class="left-col">
        <section class="panel">
          <div class="panel-head">
            <h3>Antrian Persetujuan</h3>
            <a href="/console/approvals">View All</a>
          </div>
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Type</th>
                  <th>Branch</th>
                  <th>Requester</th>
                  <th>Age</th>
                  <th class="align-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><span class="pill blue">Refund Request</span></td>
                  <td>JKT-Selatan</td>
                  <td>Budi Santoso</td>
                  <td>2h 15m</td>
                  <td class="actions-cell">
                    <button class="small-ghost">Reject</button>
                    <button class="small-primary">Approve</button>
                  </td>
                </tr>
                <tr>
                  <td><span class="pill gray">Discount Override</span></td>
                  <td>SBY-Pusat</td>
                  <td>Siti Aminah</td>
                  <td>4h 30m</td>
                  <td class="actions-cell">
                    <button class="small-ghost">Reject</button>
                    <button class="small-primary">Approve</button>
                  </td>
                </tr>
                <tr>
                  <td><span class="pill blue">Visa Expedite</span></td>
                  <td>BDG-Timur</td>
                  <td>Agus Pratama</td>
                  <td class="bad">1d 2h</td>
                  <td class="actions-cell">
                    <button class="small-ghost">Reject</button>
                    <button class="small-primary">Approve</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section class="panel">
          <div class="panel-head">
            <h3>Prioritas Hari Ini</h3>
          </div>
          <div class="priority-list">
            <article class="priority-item">
              <div>
                <p>Manifest Flight GA-980 Incomplete</p>
                <small>Missing 3 passport scans for departure tomorrow</small>
              </div>
              <span class="chip critical">Critical</span>
            </article>
            <article class="priority-item">
              <div>
                <p>Payment Gateway Sync Error</p>
                <small>BSI Virtual Account webhook failing since 08:00 AM</small>
              </div>
              <span class="chip critical">Critical</span>
            </article>
            <article class="priority-item">
              <div>
                <p>Hotel Allocation Overbooked - Makkah</p>
                <small>Package Safa-10 needs 2 extra double rooms</small>
              </div>
              <span class="chip warning">Warning</span>
            </article>
            <article class="priority-item">
              <div>
                <p>Review New Agent Application</p>
                <small>PT Bina Insani Tour documents ready for review</small>
              </div>
              <span class="chip ok">OK</span>
            </article>
          </div>
        </section>
      </div>

      <div class="right-col">
        <section class="panel">
          <div class="panel-head"><h3>Status Integrasi</h3></div>
          <div class="status-grid">
            <div class="status-item"><span>IAM</span><i class="dot ok"></i></div>
            <div class="status-item"><span>Booking</span><i class="dot ok"></i></div>
            <div class="status-item"><span>Payment</span><i class="dot fail"></i></div>
            <div class="status-item"><span>Visa Siskopatuh</span><i class="dot ok"></i></div>
            <div class="status-item wide"><span>Logistics / Inventory</span><i class="dot ok"></i></div>
          </div>
        </section>

        <section class="panel">
          <div class="panel-head"><h3>Aktivitas Terbaru</h3></div>
          <div class="timeline">
            <article>
              <p>Batch <strong>#PKG-102</strong> visa submitted</p>
              <small>System • 10 mins ago</small>
            </article>
            <article>
              <p>Hotel allocation updated for Makkah</p>
              <small>Admin Ops • 45 mins ago</small>
            </article>
            <article>
              <p>Payment sync failure reported</p>
              <small>Monitor • 2 hrs ago</small>
            </article>
            <article>
              <p>New agent application received</p>
              <small>System • 3 hrs ago</small>
            </article>
          </div>
        </section>
      </div>
    </div>
  </section>
</main>

<style>
  .dashboard-shell {
    min-height: 100vh;
    background: #f7f9fb;
  }

  .topbar {
    position: sticky;
    top: 0;
    z-index: 30;
    height: 4rem;
    background: rgb(255 255 255 / 0.9);
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0 1.25rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    backdrop-filter: blur(8px);
  }

  .search-wrap {
    flex: 1;
    max-width: 32rem;
    position: relative;
  }

  .search-icon {
    position: absolute;
    left: 0.65rem;
    top: 50%;
    transform: translateY(-50%);
    font-size: 1.05rem;
    color: #737686;
  }

  .search-wrap input {
    width: 100%;
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #f2f4f6;
    border-radius: 0.25rem;
    padding: 0.48rem 0.7rem 0.48rem 2.1rem;
    font-size: 0.85rem;
    color: #191c1e;
  }

  .top-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .branch-btn {
    border: 0;
    border-radius: 0.25rem;
    background: #e6e8ea;
    color: #191c1e;
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.42rem 0.6rem;
    font-size: 0.82rem;
    font-weight: 500;
    cursor: pointer;
  }

  .icon-actions {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    padding-left: 0.75rem;
    border-left: 1px solid rgb(195 198 215 / 0.55);
  }

  .icon-btn {
    border: 0;
    background: transparent;
    color: #434655;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    cursor: pointer;
    position: relative;
  }

  .icon-btn:hover {
    background: #eceef0;
  }

  .dot {
    position: absolute;
    width: 0.46rem;
    height: 0.46rem;
    border-radius: 999px;
  }

  .notif .dot {
    right: 0.4rem;
    top: 0.35rem;
    background: #ba1a1a;
    border: 2px solid #fff;
  }

  .avatar {
    border: 1px solid rgb(195 198 215 / 0.55);
    background: #b4c5ff;
    color: #00174b;
    width: 2rem;
    height: 2rem;
    border-radius: 0.25rem;
    font-weight: 700;
    cursor: pointer;
  }

  .canvas {
    padding: 1.5rem;
    max-width: 96rem;
  }

  .page-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .page-head h2 {
    margin: 0;
    font-size: 1.6rem;
    line-height: 1.2;
  }

  .page-head p {
    margin: 0.35rem 0 0;
    font-size: 0.82rem;
    color: #434655;
  }

  .head-actions {
    display: flex;
    gap: 0.6rem;
  }

  .ghost-btn,
  .primary-btn {
    border-radius: 0.25rem;
    padding: 0.55rem 0.8rem;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    border: 1px solid rgb(195 198 215 / 0.55);
  }

  .ghost-btn {
    background: #fff;
    color: #191c1e;
  }

  .primary-btn {
    border-color: #2563eb;
    background: linear-gradient(90deg, #004ac6, #2563eb);
    color: #fff;
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
  }

  .kpi-grid {
    display: grid;
    grid-template-columns: repeat(6, minmax(0, 1fr));
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .kpi-card {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    padding: 0.75rem;
    min-height: 7rem;
  }

  .kpi-top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    font-size: 0.62rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #434655;
    font-weight: 600;
  }

  .kpi-top .material-symbols-outlined {
    color: #004ac6;
    font-size: 1rem;
  }

  .kpi-top .material-symbols-outlined.bad {
    color: #ba1a1a;
  }

  .kpi-number {
    margin: 0.45rem 0 0.18rem;
    font-size: 1.3rem;
    font-weight: 700;
  }

  .kpi-hint {
    margin: 0;
    font-size: 0.68rem;
    color: #434655;
    display: inline-flex;
    align-items: center;
    gap: 0.15rem;
  }

  .kpi-hint .material-symbols-outlined {
    font-size: 0.9rem;
  }

  .kpi-hint.good {
    color: #004ac6;
  }

  .kpi-hint.bad {
    color: #ba1a1a;
  }

  .main-grid {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 1rem;
  }

  .left-col,
  .right-col {
    display: grid;
    gap: 1rem;
    align-content: start;
  }

  .panel {
    background: #fff;
    border: 1px solid rgb(195 198 215 / 0.45);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .panel-head {
    background: #e6e8ea;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    padding: 0.7rem 0.85rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
  }

  .panel-head h3 {
    margin: 0;
    font-size: 0.82rem;
    font-weight: 700;
  }

  .panel-head a {
    font-size: 0.72rem;
    color: #004ac6;
    text-decoration: none;
  }

  .table-wrap {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 0.62rem 0.85rem;
    text-align: left;
    font-size: 0.78rem;
    border-bottom: 1px solid rgb(195 198 215 / 0.45);
    white-space: nowrap;
  }

  th {
    text-transform: uppercase;
    font-size: 0.62rem;
    letter-spacing: 0.08em;
    color: #434655;
  }

  tbody tr:hover {
    background: #f2f4f6;
  }

  tbody tr:last-child td {
    border-bottom: 0;
  }

  .pill {
    display: inline-flex;
    padding: 0.12rem 0.4rem;
    border-radius: 0.2rem;
    font-size: 0.62rem;
    font-weight: 700;
  }

  .pill.blue {
    color: #004ac6;
    background: rgb(37 99 235 / 0.12);
  }

  .pill.gray {
    color: #191c1e;
    background: #e0e3e5;
  }

  .actions-cell {
    text-align: right;
  }

  .small-ghost,
  .small-primary {
    border-radius: 0.2rem;
    font-size: 0.68rem;
    padding: 0.24rem 0.48rem;
    cursor: pointer;
    border: 1px solid transparent;
  }

  .small-ghost {
    background: #e6e8ea;
    color: #434655;
    margin-right: 0.3rem;
  }

  .small-primary {
    background: #004ac6;
    color: #fff;
  }

  td.bad {
    color: #ba1a1a;
    font-weight: 600;
  }

  .priority-list {
    display: grid;
    gap: 0.6rem;
    padding: 0.85rem;
  }

  .priority-item {
    background: #f2f4f6;
    border: 1px solid rgb(195 198 215 / 0.35);
    border-radius: 0.25rem;
    padding: 0.65rem;
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.7rem;
  }

  .priority-item p {
    margin: 0;
    font-size: 0.8rem;
    font-weight: 600;
  }

  .priority-item small {
    font-size: 0.68rem;
    color: #434655;
  }

  .chip {
    border-radius: 0.2rem;
    padding: 0.18rem 0.35rem;
    font-size: 0.58rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    white-space: nowrap;
  }

  .chip.critical {
    background: #ffdad6;
    color: #93000a;
  }

  .chip.warning {
    background: #ffede6;
    color: #7d2d00;
  }

  .chip.ok {
    background: #e0e3e5;
    color: #434655;
  }

  .status-grid {
    padding: 0.85rem;
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0.55rem;
  }

  .status-item {
    background: #f2f4f6;
    border: 1px solid rgb(195 198 215 / 0.35);
    border-radius: 0.25rem;
    padding: 0.6rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 0.68rem;
    color: #434655;
    font-weight: 600;
  }

  .status-item.wide {
    grid-column: 1 / -1;
  }

  .status-item .dot {
    position: static;
    width: 0.55rem;
    height: 0.55rem;
  }

  .status-item .dot.ok {
    background: #10b981;
  }

  .status-item .dot.fail {
    background: #ba1a1a;
  }

  .timeline {
    padding: 0.85rem;
    display: grid;
    gap: 0.75rem;
  }

  .timeline article p {
    margin: 0;
    font-size: 0.78rem;
  }

  .timeline article small {
    color: #434655;
    font-size: 0.68rem;
  }

  @media (max-width: 1300px) {
    .kpi-grid {
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
  }

  @media (max-width: 980px) {
    .topbar {
      height: auto;
      flex-direction: column;
      align-items: stretch;
      padding: 0.8rem;
    }

    .search-wrap {
      max-width: none;
    }

    .top-actions {
      justify-content: space-between;
    }

    .main-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 700px) {
    .canvas {
      padding: 1rem;
    }

    .kpi-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .page-head {
      flex-direction: column;
      align-items: flex-start;
    }
  }
</style>

