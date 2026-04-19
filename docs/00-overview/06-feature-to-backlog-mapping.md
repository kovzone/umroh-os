# Feature -> Backlog Mapping (Priority-Owned by Detail Fitur)

Dokumen ini adalah layer operasional antara:

- `docs/06-features/*` (detail fitur, acceptance, edge case), dan
- task code slice di [`05-slice-engineering-checklist-and-task-codes.md`](./05-slice-engineering-checklist-and-task-codes.md).

Urutan fase **0 → 6** mengikuti [`04-delivery-plan-2p-sequence-first.md`](./04-delivery-plan-2p-sequence-first.md) (konteks tim & slice); gate dan definisi kode tugas per slice ada di **05**.

Prinsip utama:

1. `Priority` (`MH-MVP`, `MH-V1`, `SH`, `CH`) **milik detail fitur/backlog item**.
2. `Phase` + `Exec seq` + `Slice` + `Task Code` menunjukkan **urutan integrasi** dan **paket kerja** (bukan deadline kalender).
3. Satu detail fitur boleh pecah ke beberapa backlog item agar ukuran PR tetap kecil.

Catatan kelengkapan:

- Fase **0–5** berisi backlog **yang sudah dipetakan** ke slice/task code utama.
- Fase **6** memuat **seluruh 202 nomor modul** dari `docs/Modul UmrohOS - MosCoW.csv` sebagai baris backlog `BL-*` (depth expansion), selain umbrella **6.A–6.E**. Urutan kerja mengikuti **`Exec seq`** hingga **854** setelah **6.E** (lihat blockquote Fase 6).
- **Indeks `No` CSV → subsection Fase 6**: `#1–24`→**6.L**; `#25–70` & `#199–202`→**6.H**; `#71–86`→**6.G**; `#87–128`→**6.K**; `#129–150`→**6.J**; `#151–163`→**6.M**; `#164–176`→**6.N**; `#177–178` & `#187–188`→**6.F**; `#179–186` & `#189`→**6.I**; `#190–198`→**6.O**.
- Untuk **nomor modul** yang punya baris di `docs/Modul UmrohOS - MosCoW.csv`, kolom `MoSCoW` di CSV dipakai sebagai **prioritas backlog** (`Must Have` → `MH-V1`, `Should Have` → `SH`, `Could Have` → `CH`) bila belum ada keputusan berbeda di `docs/06-features/*`.
- **`Slice` / `Task Code` Fase 6** diselaraskan ke [**05**](./05-slice-engineering-checklist-and-task-codes.md) (§ **Fase 6 — Depth backlog**); kolom di tabel Fase 6 memuat kode tugas per domain (bukan satu baris = satu PR).

---

## Format backlog ID

Gunakan pola:

`BL-{DOMAIN}-{NNN}`

Contoh:

- `BL-IAM-001`
- `BL-CAT-004`
- `BL-BOOK-007`
- `BL-PAY-003`
- `BL-B2C-001`
- `BL-JMJ-001`

Status kerja yang disarankan:

`todo -> in_progress -> in_review -> done`

### `Exec seq` (urutan numerik)

Untuk **Fase 0–5**, angka biasanya mengikuti pola longgar **Fase × 100 + urutan kecil**:

- **Fase 0** → `000–099` (contoh: `005`)
- **Fase 1** → `100–199` (contoh: `110`)
- **Fase 2** → `200–299`
- dst.

**Fase 6 (depth):** dipakai rentang **`600+`** (hingga `854`) sebagai jalur dedikasi setelah core; subsection **6.G–6.O** tidak wajib mengikuti rumus `Phase × 100` per baris — urutan per baris mengikuti nilai **`Exec seq`** di tabel.

Angka lebih kecil = lebih dulu. Banyak baris boleh **angka sama** kalau paralel aman lintas domain.

---

## Fase 0 — Engineering bootstrap (S0)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 0 | S0-J-01 | Template kontrak + README `docs/contracts/` | MH-MVP | 005 | S0 | S0-J-01 | BL-ENG-001 | J | todo | — | Template `slice-Sx.md` + README ter-merge |
| 0 | S0-J-02 | Branch strategy + merge ownership | MH-MVP | 006 | S0 | S0-J-02 | BL-ENG-002 | J | todo | BL-ENG-001 | Aturan merge tertulis |
| 0 | S0-J-03 | DoR/DoD per PR | MH-MVP | 007 | S0 | S0-J-03 | BL-ENG-003 | J | todo | BL-ENG-001 | DoR/DoD singkat tersedia |
| 0 | S0-L-01 | Matriks route UI publik vs internal (S1) | MH-MVP | 008 | S0 | S0-L-01 | BL-FE-PLN-001 | L | todo | BL-ENG-001 | Tabel role vs URL awal |
| 0 | S0-E-01 | Service ownership S1–S2 (PR owner) | MH-MVP | 009 | S0 | S0-E-01 | BL-ENG-004 | E | todo | BL-ENG-001 | Tabel service vs owner |

---

## Fase 1 — Discover + draft booking (S1)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 1 | F1-W1 | Login internal + refresh token flow | MH-MVP | 110 | S1 | S1-E-04 | BL-IAM-001 | E | todo | S1-J-01..S1-J-04 | Login sukses, refresh jalan, unauthorized -> 401 |
| 1 | F1-W3 | `CheckPermission` middleware untuk route internal | MH-MVP | 111 | S1 | S1-E-04 | BL-IAM-002 | E | todo | S1-J-01..S1-J-04 | Route finance ditolak untuk role non-finance |
| 1 | F1-W5 | Suspend/revoke session dasar | MH-MVP | 112 | S1 | S1-E-04 | BL-IAM-003 | E | todo | S1-J-01..S1-J-04 | User suspended tidak bisa akses ulang |
| 1 | F1-AC | Audit log untuk state-changing call | MH-MVP | 113 | S1 | S1-E-04 | BL-IAM-004 | E | todo | S1-J-01..S1-J-04 | Aksi create/update booking tercatat di audit |
| 1 | F2-W2 | Read model paket + departure aktif | MH-MVP | 120 | S1 | S1-E-02 | BL-CAT-001 | E | todo | S1-J-01..S1-J-04 | `list/detail` tampil hanya paket/departure valid |
| 1 | F2-W3 | Data departure memuat seat cap + status | MH-MVP | 121 | S1 | S1-E-02 | BL-CAT-002 | E | todo | S1-J-01..S1-J-04 | Detail departure expose sisa seat konsisten |
| 1 | F2-W6 | Atomic `ReserveSeats` + `ReleaseSeats` | MH-MVP | 122 | S1 | S1-E-03 | BL-CAT-003 | E | todo | S1-J-01..S1-J-04 | Last-seat race aman, tanpa oversell |
| 1 | F2-AC | Endpoint katalog publik untuk B2C | MH-MVP | 123 | S1 | S1-E-02 | BL-CAT-004 | E | todo | S1-J-01..S1-J-04 | B2C bisa browse paket tanpa error kontrak |
| 1 | F4-W1 | Create draft booking dari B2C flow | MH-MVP | 130 | S1 | S1-E-03 | BL-BOOK-001 | E | todo | S1-J-01..S1-J-04 | Booking `draft` tersimpan dengan field minimum |
| 1 | F4-W2 | Stamp channel attribution (`b2c_self`/`b2b_agent`) | MH-MVP | 131 | S1 | S1-E-03 | BL-BOOK-002 | E | todo | S1-J-01..S1-J-04 | Booking menyimpan `channel` + `agent_id` bila ada |
| 1 | F4-W4 | Status machine sampai `pending_payment`/`expired` | MH-MVP | 132 | S1 | S1-E-03 | BL-BOOK-003 | E | todo | S1-J-01..S1-J-04 | Transisi status valid tanpa loncat state |
| 1 | F4-W8 | Submit validation dasar (package aktif, seat cukup) | MH-MVP | 133 | S1 | S1-E-03 | BL-BOOK-004 | E | todo | S1-J-01..S1-J-04 | Hard fail jika seat habis / departure invalid |
| 1 | F4-W8 | Gate dokumen minimum (Q006) saat submit | MH-MVP | 134 | S1 | S1-E-03 | BL-BOOK-005 | E | todo | S1-J-01..S1-J-04 | Missing doc -> error jelas per jamaah/doc kind |
| 1 | F4-W6 | Integrasi `ValidateMahram` sebagai warning | MH-MVP | 135 | S1 | S1-E-03 | BL-BOOK-006 | E | todo | S1-J-01..S1-J-04 | Hasil mahram tersimpan, tidak blok submit |
| 1 | F4-UI | UI katalog -> detail -> form booking | MH-MVP | 140 | S1 | S1-L-02 | BL-FE-BOOK-001 | L | todo | S1-J-01..S1-J-04 | User bisa mencapai form booking dari katalog |
| 1 | F4-UI | Integrasi FE ke API katalog | MH-MVP | 141 | S1 | S1-L-03 | BL-FE-BOOK-002 | L | todo | S1-J-01..S1-J-04 | FE memanggil list/detail sesuai kontrak |
| 1 | F4-UI | Integrasi FE create draft booking | MH-MVP | 142 | S1 | S1-L-04 | BL-FE-BOOK-003 | L | todo | S1-J-01..S1-J-04 | Submit form membuat draft booking sukses |

---

## Fase 2 — Get paid (S2)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 2 | F5-W1 | Create invoice + issue VA + TTL | MH-MVP | 210 | S2 | S2-E-02 | BL-PAY-001 | E | todo | S2-J-01..S2-J-04 | VA terbit, `expires_at` tersimpan sesuai config |
| 2 | F5-W1 | Gateway selection + fallback rule (Q013) | MH-MVP | 211 | S2 | S2-E-02 | BL-PAY-002 | E | todo | S2-J-01..S2-J-04 | Failover hanya pada timeout/5xx |
| 2 | F5-W2 | Webhook signature verification | MH-MVP | 220 | S2 | S2-E-02 | BL-PAY-003 | E | todo | S2-J-01..S2-J-04 | Signature salah -> 401, tidak update bisnis |
| 2 | F5-W2 | Webhook idempotency (`gateway_txn_id`) | MH-MVP | 221 | S2 | S2-E-02 | BL-PAY-004 | E | todo | S2-J-01..S2-J-04 | Replay webhook -> no-op aman |
| 2 | F5-W2 | Update `paid_amount` + signal `MarkBookingPaid` | MH-MVP | 222 | S2 | S2-E-02 | BL-PAY-005 | E | todo | S2-J-01..S2-J-04 | Booking status sinkron partial/lunas |
| 2 | F5-W5 | Reconciliation cron untuk miss webhook | MH-MVP | 223 | S2 | S2-E-03 | BL-PAY-006 | E | todo | S2-J-01..S2-J-04 | Webhook drop terecovery pada siklus berikutnya |
| 2 | F5-W8 | Refund dasar terhubung cancellation booking | MH-MVP | 224 | S2 | S2-E-02 | BL-PAY-007 | E | todo | S2-J-01..S2-J-04 | Cancel -> refund flow tercatat dan idempotent |
| 2 | F5-W9 | FX snapshot + rounding rule Q001 | MH-MVP | 225 | S2 | S2-E-02 | BL-PAY-008 | E | todo | S2-J-01..S2-J-04 | Snapshot immutable setelah payment pertama |
| 2 | F5-UI | Halaman checkout menampilkan VA + countdown | MH-MVP | 230 | S2 | S2-L-02 | BL-FE-PAY-001 | L | todo | S2-J-01..S2-J-04 | User lihat akun VA, nominal, expiry real-time |
| 2 | F5-UI | Wiring FE booking -> payment call | MH-MVP | 231 | S2 | S2-L-03 | BL-FE-PAY-002 | L | todo | S2-J-01..S2-J-04 | Dari draft bisa lanjut issue payment |
| 2 | F10-W14 | Payment link generator (CS closing) | MH-V1 | 235 | S2 | S2-E-02 | BL-PAY-020 | E | todo | S2-J-01..S2-J-04 | CS bisa issue link/VA untuk booking existing |
| 2 | F5-E2E | End-to-end uji stub -> gateway nyata | MH-MVP | 240 | S2 | S2-J-05 | BL-QA-001 | J | todo | S2-J-01..S2-J-04 | Skenario draft->bayar->paid lulus E2E |

---

## Fase 3 — Fulfillment + finance minimum pasca bayar (S3)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 3 | F3-W1 | Upload dokumen jamaah (KTP/paspor/foto) | MH-V1 | 310 | S3 | S3-E-02 | BL-DOC-001 | E | todo | S3-J-01..S3-J-03 | Dokumen tersimpan + status lifecycle jalan |
| 3 | F3-W4 | OCR paspor + fallback manual review | MH-V1 | 311 | S3 | S3-E-02 | BL-DOC-002 | E | todo | S3-J-01..S3-J-03 | OCR hasil tersimpan, low confidence masuk review |
| 3 | F3-W3 | Verification queue approve/reject + audit | MH-V1 | 312 | S3 | S3-E-02 | BL-DOC-003 | E | todo | S3-J-01..S3-J-03 | Approve/reject tercatat audit + notif reject terkirim |
| 3 | F7-W5 | Manifest generator (PDF/XLSX) | MH-V1 | 320 | S3 | S3-E-02 | BL-OPS-001 | E | todo | S3-J-01..S3-J-03 | Manifest bisa di-generate per departure |
| 3 | F7-W2 | Smart grouping room allocation run + commit | MH-V1 | 321 | S3 | S3-E-02 | BL-OPS-002 | E | todo | S3-J-01..S3-J-03 | Grouping run menghasilkan alokasi valid |
| 3 | F7-W6 | ID card + luggage tag issuance QR signed | MH-V1 | 322 | S3 | S3-E-02 | BL-OPS-003 | E | todo | S3-J-01..S3-J-03 | QR valid/terverifikasi, tamper ditolak |
| 3 | F8-W10 | Trigger fulfillment hanya `paid_in_full` | MH-V1 | 330 | S3 | S3-E-02 | BL-LOG-001 | E | todo | S3-J-01..S3-J-03 | Queue fulfillment tidak memuat booking non-lunas |
| 3 | F8-W11 | Shipment + tracking number + WA notify | MH-V1 | 331 | S3 | S3-E-02 | BL-LOG-002 | E | todo | S3-J-01..S3-J-03 | Pengiriman menghasilkan resi + notifikasi |
| 3 | F8-W12 | Self-pickup QR single-use | MH-V1 | 332 | S3 | S3-E-02 | BL-LOG-003 | E | todo | S3-J-01..S3-J-03 | QR pickup valid sekali pakai + expiry |
| 3 | F9-W2 | Posting jurnal penerimaan pembayaran (deferred rev) | MH-V1 | 340 | S3 | S3-E-03 | BL-FIN-001 | E | todo | S3-J-01..S3-J-03 | Dr Bank / Cr Hutang Jamaah otomatis |
| 3 | F9-W4 | Auto-AP dari GRN (sinkron) | MH-V1 | 341 | S3 | S3-E-03 | BL-FIN-002 | E | todo | S3-J-01..S3-J-03 | GRN gagal jika posting AP gagal |
| 3 | F9-W9 | Engine jurnal double-entry + idempotency source | MH-V1 | 342 | S3 | S3-E-03 | BL-FIN-003 | E | todo | S3-J-01..S3-J-03 | Tidak ada jurnal unbalanced / duplicate source |
| 3 | F9-W10 | Revenue recognition trigger di event departure | MH-V1 | 343 | S3 | S3-E-03 | BL-FIN-004 | E | todo | S3-J-01..S3-J-03 | Pendapatan diakui saat trigger Q043, bukan saat bayar |
| 3 | F7-UI | Ops board status fulfillment+manifest ringkas | MH-V1 | 350 | S3 | S3-L-02 | BL-FE-OPS-001 | L | todo | S3-J-01..S3-J-03 | UI menampilkan status ops utama per booking |

---

## Fase 4 — Growth loop dasar (S4)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 4 | F10-W3 | Lead tracker + tagging (CRM dasar) | MH-V1 | 410 | S4 | S4-E-02 | BL-CRM-001 | E | todo | S4-J-01..S4-J-02 | Lead tersimpan + bisa di-list/filter/tag |
| 4 | F10-W8 | Attribution + UTM reconciliation (dasar) | MH-V1 | 411 | S4 | S4-E-02 | BL-CRM-002 | E | todo | S4-J-01..S4-J-02 | UTM tersimpan + konsisten ke booking |
| 4 | F10-W4 | CS round-robin + SLA dasar | MH-V1 | 412 | S4 | S4-E-02 | BL-CRM-003 | E | todo | S4-J-01..S4-J-02 | Distribusi lead + SLA minimal jalan |
| 4 | F10-UI | Lead capture form (publik/internal) | MH-V1 | 420 | S4 | S4-L-01 | BL-FE-CRM-001 | L | todo | S4-J-01..S4-J-02 | Form submit sukses + validasi dasar |

---

## Fase 5 — Hardening & readiness (S5)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 5 | F9-W15 | Laporan dasar Neraca + Laba Rugi + Arus Kas | MH-V1 | 510 | S5 | S5-E-01 | BL-FIN-005 | E | todo | S5-J-01..S5-J-02 | 3 laporan utama bisa ditarik periode berjalan |
| 5 | F9-W18 | Finance audit trail + anti-delete jurnal | MH-V1 | 511 | S5 | S5-E-01 | BL-FIN-006 | E | todo | S5-J-01..S5-J-02 | Delete jurnal ditolak, koreksi via counter-entry |
| 5 | F9-UI | Finance view ringkas status jurnal/payment | MH-V1 | 520 | S5 | S5-L-01 | BL-FE-FIN-001 | L | todo | S5-J-01..S5-J-02 | Finance bisa trace booking -> invoice -> jurnal |
| 5 | QA | UAT checklist inti + regression permission/audit | MH-V1 | 530 | S5 | S5-L-01 | BL-QA-002 | L | todo | S5-J-01..S5-J-02 | Daftar skenario lolos + bukti |
| 5 | QA | UAT checklist payment/finance/logistics | MH-V1 | 531 | S5 | S5-E-01 | BL-QA-003 | E | todo | S5-J-01..S5-J-02 | Daftar skenario lolos + bukti |

---

## Fase 6 — Depth expansion (setelah core stabil)

> Bagian ini menyusul **setelah** Fase 1–5 lolos gate integrasi.  
> `Exec seq` sengaja **600+** supaya tidak “loncat” sebelum fondasi jalan.  
> Setelah **6.E** (hingga `Exec seq` **652**), rentang berikutnya berurutan: **6.G** **653–668** → **6.F** **669–672** → **6.H** **673–722** → **6.I** **723–731** → **6.J** **732–753** → **6.K** **754–795** → **6.L** **796–819** → **6.M** **820–832** → **6.N** **833–845** → **6.O** **846–854** (seluruh baris CSV **#1–#202** tercakup).

### 6.A — Marketing/CRM dalam (F10)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F10-W1 | Agent onboarding + E-KYC + e-MoU | MH-V1 | 610 | S4 | S4-L-02 | BL-CRM-010 | L | todo | S4-J-01..S4-J-02 | Alur register→approve→aktif role agent |
| 6 | F10-W2 | Replica site + share UTM + tracking | MH-V1 | 611 | S4 | S4-L-02 | BL-CRM-011 | L | todo | S1-J-01..S1-J-04 | Replica site render katalog + tracking lead |
| 6 | F10-W9 | Dompet komisi (saldo + status dasar) | MH-V1 | 612 | S4 | S4-E-03 | BL-CRM-012 | L | todo | S2-J-01..S2-J-04 | Saldo komisi konsisten dengan event bayar |

Pecahan modul F10 resmi dari CSV ada di **6.H** (`BL-CRM-017`–`BL-CRM-066`). Baris **6.A** (`BL-CRM-010`–`012`) tetap sebagai **paket integrasi** ringkas yang menjembatani beberapa modul sekaligus.

### 6.B — Finance dalam (F9)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F9-W5 | AP disbursement ladder minimum | MH-V1 | 620 | S3 | S3-E-07 | BL-FIN-010 | E | todo | S3-J-01..S3-J-03 | Batch AP + approval + audit |
| 6 | F9-W17 | AR/AP aging alert dasar | MH-V1 | 621 | S3 | S3-E-07 | BL-FIN-011 | E | todo | S3-J-01..S3-J-03 | Aging bucket tampil + alert rule dasar |

### 6.C — Warehouse/procurement dalam (F8)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F8-W1 | PR + budget gate | MH-V1 | 630 | S3 | S3-E-05 | BL-LOG-010 | E | todo | S3-J-01..S3-J-03 | Over-budget PR ditolak |
| 6 | F8-W4 | GRN + QC + auto-AP sync | MH-V1 | 631 | S3 | S3-E-05 | BL-LOG-011 | E | todo | S3-J-01..S3-J-03 | GRN rollback saat finance gagal |
| 6 | F8-W7 | Kit assembly atomic | MH-V1 | 632 | S3 | S3-E-05 | BL-LOG-012 | E | todo | S3-J-01..S3-J-03 | Assembly all-or-nothing |

### 6.D — Ops lapangan dalam (F7)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F7-W7 | ALL system scan + idempotency | MH-V1 | 640 | S3 | S3-E-04 | BL-OPS-010 | E | todo | S3-J-01..S3-J-03 | Event scan idempotent |
| 6 | F7-W10 | Bus boarding scan + roster | MH-V1 | 641 | S3 | S3-E-04 | BL-OPS-011 | E | todo | S3-J-01..S3-J-03 | Boarding roster konsisten |

### 6.E — Visa pipeline Must (#97) (F6)

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F6-W1 | Readiness auto-transition `WAITING_DOCS → READY` | MH-V1 | 650 | S3 | S3-E-06 | BL-VISA-001 | E | todo | S3-J-01..S3-J-03 | Transisi status tercatat + idempotent |
| 6 | F6-W2 | Bulk submit visa all-or-nothing | MH-V1 | 651 | S3 | S3-E-06 | BL-VISA-002 | E | todo | BL-VISA-001 | Bulk submit atomic sesuai spec |
| 6 | F6-W3 | Poll status provider + history | MH-V1 | 652 | S3 | S3-E-06 | BL-VISA-003 | E | todo | BL-VISA-002 | History poll tersimpan |

### 6.G — Master Product modul CSV (#71–#86) (F2)

Sumber prioritas per baris: `docs/Modul UmrohOS - MosCoW.csv` (`No` 71–86).

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F2-CSV-71 | Master data hotel (Database Hotel) | MH-V1 | 653 | S1 | S1-E-05 | BL-CAT-005 | E | todo | S1-J-01..S1-J-04 | CRUD + referensi hotel bisa dipakai komposisi produk |
| 6 | F2-CSV-72 | Master data pembimbing/muthawwif (Database Pembimbing) | MH-V1 | 654 | S1 | S1-E-05 | BL-CAT-006 | E | todo | S1-J-01..S1-J-04 | CRUD + assignment ke produk/departure konsisten |
| 6 | F2-CSV-73 | Master transportasi & maskapai | MH-V1 | 655 | S1 | S1-E-05 | BL-CAT-007 | E | todo | S1-J-01..S1-J-04 | Referensi carrier/rute/mode tersedia untuk komposisi produk |
| 6 | F2-CSV-74 | Varian produk perjalanan (template + constraint) | MH-V1 | 656 | S1 | S1-E-05 | BL-CAT-008 | E | todo | S1-J-01..S1-J-04 | Varian tidak merusak read model publik + validasi publish |
| 6 | F2-CSV-75 | Produk finansial & retail (addon non-core) | SH | 657 | S1 | S1-E-05 | BL-CAT-009 | E | todo | S1-J-01..S1-J-04 | Addon terpisah dari paket inti + pricing rule jelas |
| 6 | F2-CSV-76 | Input massal cerdas (import sheet + validasi) | SH | 658 | S1 | S1-E-05 | BL-CAT-010 | E | todo | S1-J-01..S1-J-04 | Import partial failure aman + laporan error per baris |
| 6 | F2-CSV-77 | Pembaruan massal (bulk edit dengan guard) | SH | 659 | S1 | S1-E-05 | BL-CAT-011 | E | todo | S1-J-01..S1-J-04 | Bulk update preview + audit + rollback policy |
| 6 | F2-CSV-78 | Generator flyer dinamis | SH | 660 | S4 | S4-L-02 | BL-CRM-013 | L | todo | S4-J-01..S4-J-02 | Flyer render dari template + data paket aktual |
| 6 | F2-CSV-79 | Omni-flyer (multi-format/channel) | CH | 661 | S4 | S4-L-02 | BL-CRM-014 | L | todo | S4-J-01..S4-J-02 | Satu sumber konten -> beberapa varian output |
| 6 | F2-CSV-80 | Itinerary interaktif (shareable) | SH | 662 | S4 | S4-L-02 | BL-CRM-015 | L | todo | S4-J-01..S4-J-02 | Itinerary konsisten dengan master itinerary + deep-link |
| 6 | F2-CSV-81 | Otomasi copywriting | CH | 663 | S4 | S4-L-02 | BL-CRM-016 | L | todo | S4-J-01..S4-J-02 | Output reviewable + tidak publish otomatis tanpa gate |
| 6 | F2-CSV-82 | Sinkronisasi satu pintu (katalog -> saluran agen) | MH-V1 | 664 | S4 | S4-E-04 | BL-CAT-012 | E | todo | S4-J-01..S4-J-02 | Perubahan master terpropagasi idempotent per agen |
| 6 | F2-CSV-83 | Auto-update agen (versi katalog + diff) | MH-V1 | 665 | S4 | S4-E-04 | BL-CAT-013 | E | todo | S4-J-01..S4-J-02 | Agen punya snapshot versi + mekanisme upgrade aman |
| 6 | F2-CSV-84 | Pelacakan kursi lintas saluran (agen/B2C) | MH-V1 | 666 | S1 | S1-E-05 | BL-BOOK-007 | E | todo | S1-J-01..S1-J-04 | Seat state tidak double-sell antar saluran |
| 6 | F2-CSV-85 | Tampilan ganda dashboard (role/context switch) | SH | 667 | S5 | S5-L-02 | BL-DASH-005 | L | todo | S5-J-01..S5-J-02 | Dua mode tampilan konsisten permission + filter |
| 6 | F2-CSV-86 | Ceklis kesiapan vendor per departure | MH-V1 | 668 | S3 | S3-E-04 | BL-OPS-020 | E | todo | S3-J-01..S3-J-03 | Checklist item + status + bukti attachment minimum |

### 6.F — Dashboard Must modules (#177–#178, #187–#188) (F11)

Modul Dashboard **Should** selebihnya (**#179–#186**, **#189**) ada di **6.I** (`BL-DASH-006`–`014`). Modul **#85** *tampilan ganda* terkait dashboard ada di **6.G** (`BL-DASH-005`).

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F11-CSV-177 | Eksekusi Vendor (widget readiness) | MH-V1 | 669 | S5 | S5-L-02 | BL-DASH-001 | L | todo | BL-OPS-020 | Widget checklist konsisten dengan event checklist vendor |
| 6 | F11-CSV-178 | Ketersediaan Kursi (widget) | MH-V1 | 670 | S5 | S5-L-02 | BL-DASH-002 | L | todo | S1-J-01..S1-J-04 | Widget seat inventory real-time |
| 6 | F11-CSV-187 | Arus Kas Instan (widget) | MH-V1 | 671 | S5 | S5-L-02 | BL-DASH-003 | L | todo | S5-J-01..S5-J-02 | Ringkasan kas tampil konsisten dengan F9 |
| 6 | F11-CSV-188 | Laporan Keuangan Eksekutif (widget) | MH-V1 | 672 | S5 | S5-L-02 | BL-DASH-004 | L | todo | S5-J-01..S5-J-02 | Ringkasan P&L/neraca eksekutif |

### 6.H — Marketing/CRM & Alumni/ZISWAF modul CSV (#25–#70, #199–#202) (F10)

Sumber prioritas per baris: `docs/Modul UmrohOS - MosCoW.csv` (kolom `MoSCoW`). Modul **#199–#202** tercatat di CSV bawah *Fitur Pelengkap & Daily App*; pemetaan domain mengikuti cakupan F10 di `docs/06-features/10-marketing-crm-agents.md` (alumni referral + ZISWAF slice).

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F10-CSV-25 | Pendaftaran Keagenan Mandiri | MH-V1 | 673 | S4 | S4-L-02 | BL-CRM-017 | L | todo | S4-J-01..S4-J-02 | Form register mitra + status pipeline |
| 6 | F10-CSV-26 | E-KYC & Verifikasi | MH-V1 | 674 | S4 | S4-L-02 | BL-CRM-018 | L | todo | S4-J-01..S4-J-02 | Upload KYC + hasil verifikasi tercatat |
| 6 | F10-CSV-27 | E-Signature | SH | 675 | S4 | S4-L-02 | BL-CRM-019 | L | todo | S4-J-01..S4-J-02 | MoU digital tersimpan + audit trail |
| 6 | F10-CSV-28 | Website Replika | MH-V1 | 676 | S4 | S4-L-02 | BL-CRM-020 | L | todo | S4-J-01..S4-J-02 | Replica render katalog + identitas agen |
| 6 | F10-CSV-29 | One-Click Social Sharing | SH | 677 | S4 | S4-L-02 | BL-CRM-021 | L | todo | S4-J-01..S4-J-02 | Share one-click + UTM terbawa |
| 6 | F10-CSV-30 | Kartu Nama Digital | CH | 678 | S4 | S4-L-02 | BL-CRM-022 | L | todo | S4-J-01..S4-J-02 | Kartu nama digital generate + share |
| 6 | F10-CSV-31 | Bank Konten | SH | 679 | S4 | S4-L-02 | BL-CRM-023 | L | todo | S4-J-01..S4-J-02 | Bank konten browse/search + hak akses |
| 6 | F10-CSV-32 | Auto-Watermark Flyer & Share | SH | 680 | S4 | S4-L-02 | BL-CRM-024 | L | todo | S4-J-01..S4-J-02 | Flyer watermark WA agen + export |
| 6 | F10-CSV-33 | Galeri Dokumentasi per Program | CH | 681 | S4 | S4-L-02 | BL-CRM-025 | L | todo | S4-J-01..S4-J-02 | Galeri per program + kontrol publikasi |
| 6 | F10-CSV-34 | Integrasi Tracking Code Mandiri | CH | 682 | S4 | S4-L-02 | BL-CRM-026 | L | todo | S4-J-01..S4-J-02 | Tracking code agen terpasang + validasi |
| 6 | F10-CSV-35 | Leads Tracker | MH-V1 | 683 | S4 | S4-L-02 | BL-CRM-027 | L | todo | S4-J-01..S4-J-02 | Lead masuk + sumber/atribusi dasar |
| 6 | F10-CSV-36 | Reminder & Follow-up Tagging | SH | 684 | S4 | S4-L-02 | BL-CRM-028 | L | todo | S4-J-01..S4-J-02 | Reminder + tag follow-up tersimpan |
| 6 | F10-CSV-37 | Saldo & Status Komisi | MH-V1 | 685 | S4 | S4-L-02 | BL-CRM-029 | L | todo | S2-J-01..S2-J-04 | Saldo komisi konsisten per agen |
| 6 | F10-CSV-38 | Notifikasi Real-Time | SH | 686 | S4 | S4-L-02 | BL-CRM-030 | L | todo | S2-J-01..S2-J-04 | Notifikasi event komisi real-time |
| 6 | F10-CSV-39 | Riwayat Transaksi & Pencairan | MH-V1 | 687 | S4 | S4-L-02 | BL-CRM-031 | L | todo | S2-J-01..S2-J-04 | Riwayat + ajukan pencairan + status |
| 6 | F10-CSV-40 | LMS | SH | 688 | S4 | S4-L-02 | BL-CRM-032 | L | todo | S4-J-01..S4-J-02 | Kurs/modul academy dasar |
| 6 | F10-CSV-41 | Kuis & Lencana | CH | 689 | S4 | S4-L-02 | BL-CRM-033 | L | todo | S4-J-01..S4-J-02 | Kuis + lencana tercatat |
| 6 | F10-CSV-42 | Script Jualan | SH | 690 | S4 | S4-L-02 | BL-CRM-034 | L | todo | S4-J-01..S4-J-02 | Script jualan searchable per tier |
| 6 | F10-CSV-43 | Papan Peringkat | CH | 691 | S4 | S4-L-02 | BL-CRM-035 | L | todo | S4-J-01..S4-J-02 | Leaderboard per aturan yang ditetapkan |
| 6 | F10-CSV-44 | Push Notification | SH | 692 | S4 | S4-L-02 | BL-CRM-036 | L | todo | S4-J-01..S4-J-02 | Push academy terjadwal |
| 6 | F10-CSV-45 | Dashboard Super-View | MH-V1 | 693 | S4 | S4-L-02 | BL-CRM-037 | L | todo | S4-J-01..S4-J-02 | Super-view agregat downline |
| 6 | F10-CSV-46 | Leveling Otomatis | SH | 694 | S4 | S4-L-02 | BL-CRM-038 | L | todo | S4-J-01..S4-J-02 | Level tier otomatis sesuai aturan |
| 6 | F10-CSV-47 | Overriding Commission | MH-V1 | 695 | S4 | S4-L-02 | BL-CRM-039 | L | todo | S4-J-01..S4-J-02 | Override komisi dihitung deterministik |
| 6 | F10-CSV-48 | Ads Manager Lite | SH | 696 | S4 | S4-E-03 | BL-CRM-040 | E | todo | S4-J-01..S4-J-02 | Iklan spend + sync status kampanye |
| 6 | F10-CSV-49 | Manajemen UTM & Atribusi Kampanye | SH | 697 | S4 | S4-E-03 | BL-CRM-041 | E | todo | S4-J-01..S4-J-02 | UTM builder + konsistensi ke lead/booking |
| 6 | F10-CSV-50 | Landing Page Builder & A/B Testing | SH | 698 | S4 | S4-E-03 | BL-CRM-042 | E | todo | S4-J-01..S4-J-02 | LP publish + varian A/B + metrik |
| 6 | F10-CSV-51 | Content Planner & Calendar | CH | 699 | S4 | S4-E-03 | BL-CRM-043 | E | todo | S4-J-01..S4-J-02 | Kalender konten + assignment |
| 6 | F10-CSV-52 | Content Publisher & Scheduler | CH | 700 | S4 | S4-E-03 | BL-CRM-044 | E | todo | S4-J-01..S4-J-02 | Jadwal publish multi-channel |
| 6 | F10-CSV-53 | Omni-Channel Distribution | CH | 701 | S4 | S4-E-03 | BL-CRM-045 | E | todo | S4-J-01..S4-J-02 | Distribusi konten omni saluran |
| 6 | F10-CSV-54 | Social Media & Content Analytics | CH | 702 | S4 | S4-E-03 | BL-CRM-046 | E | todo | S4-J-01..S4-J-02 | Analitik konten + export dasar |
| 6 | F10-CSV-55 | Bot Filter & Auto-Classification | SH | 703 | S4 | S4-E-03 | BL-CRM-047 | E | todo | S4-J-01..S4-J-02 | Bot filter + klasifikasi lead |
| 6 | F10-CSV-56 | Pesan Berantai | SH | 704 | S4 | S4-E-03 | BL-CRM-048 | E | todo | S4-J-01..S4-J-02 | Pesan berantai template + limit |
| 6 | F10-CSV-57 | Pemicu Momen | CH | 705 | S4 | S4-E-03 | BL-CRM-049 | E | todo | S4-J-01..S4-J-02 | Pemicu momen otomatis |
| 6 | F10-CSV-58 | Segmentasi Database Cerdas | SH | 706 | S4 | S4-E-03 | BL-CRM-050 | E | todo | S4-J-01..S4-J-02 | Segmentasi query tersimpan + preview count |
| 6 | F10-CSV-59 | Pusat Siaran Massal | SH | 707 | S4 | S4-E-03 | BL-CRM-051 | E | todo | S4-J-01..S4-J-02 | Siaran massal consent + rate limit |
| 6 | F10-CSV-60 | Distribusi Adil | MH-V1 | 708 | S4 | S4-E-03 | BL-CRM-052 | E | todo | S4-J-01..S4-J-02 | Distribusi lead adil + audit |
| 6 | F10-CSV-61 | Pemicu Kecepatan Respons | SH | 709 | S4 | S4-E-03 | BL-CRM-053 | E | todo | S4-J-01..S4-J-02 | SLA trigger + eskalasi |
| 6 | F10-CSV-62 | Rekam Jejak Leads & Tagging | MH-V1 | 710 | S4 | S4-E-03 | BL-CRM-054 | E | todo | S4-J-01..S4-J-02 | Timeline lead + tagging multi |
| 6 | F10-CSV-63 | Generator Penawaran Harga | SH | 711 | S4 | S4-E-03 | BL-CRM-055 | E | todo | S4-J-01..S4-J-02 | Penawaran PDF/link + nomor referensi |
| 6 | F10-CSV-64 | Pembuat Link Pembayaran | MH-V1 | 712 | S4 | S4-E-03 | BL-CRM-056 | E | todo | BL-PAY-020 | Link bayar terbit untuk booking existing (CS closing) |
| 6 | F10-CSV-65 | e-Approval Diskon | SH | 713 | S4 | S4-E-03 | BL-CRM-057 | E | todo | S4-J-01..S4-J-02 | Alur diskon multi-level approval |
| 6 | F10-CSV-66 | Loyalitas & Referral Alumni | CH | 714 | S4 | S4-E-03 | BL-CRM-058 | E | todo | S4-J-01..S4-J-02 | Referral code + reward tracking |
| 6 | F10-CSV-67 | Dashboard Kinerja CS | SH | 715 | S4 | S4-E-03 | BL-CRM-059 | E | todo | S4-J-01..S4-J-02 | Dashboard metrik CS konsisten SLA |
| 6 | F10-CSV-68 | Kalkulator ROAS | SH | 716 | S4 | S4-E-03 | BL-CRM-060 | E | todo | S4-J-01..S4-J-02 | ROAS input spend + revenue attrib |
| 6 | F10-CSV-69 | Sinkronisasi Retargeting | CH | 717 | S4 | S4-E-03 | BL-CRM-061 | E | todo | S4-J-01..S4-J-02 | Sinkron audience retargeting |
| 6 | F10-CSV-70 | Radar Prospek Lama | CH | 718 | S4 | S4-E-03 | BL-CRM-062 | E | todo | S4-J-01..S4-J-02 | Radar prospek dorman + tugas |
| 6 | F10-CSV-199 | Pusat Referral Alumni | SH | 719 | S4 | S4-E-03 | BL-CRM-063 | E | todo | S4-J-01..S4-J-02 | Referral alumni + atribusi booking |
| 6 | F10-CSV-200 | Tabungan Niat Kembali | CH | 720 | S4 | S4-E-03 | BL-CRM-064 | E | todo | S4-J-01..S4-J-02 | Tabungan niat kembali alur dasar |
| 6 | F10-CSV-201 | Kalkulator Zakat | CH | 721 | S4 | S4-E-03 | BL-CRM-065 | E | todo | S4-J-01..S4-J-02 | Kalkulator zakat input + hasil |
| 6 | F10-CSV-202 | Sedekah & Infaq Pagi | CH | 722 | S4 | S4-E-03 | BL-CRM-066 | E | todo | S4-J-01..S4-J-02 | Sedekah/infaq pagi flow + bukti |

### 6.I — Dashboard modul CSV — sisa Should (#179–#186, #189) (F11)

Sumber prioritas per baris: `docs/Modul UmrohOS - MosCoW.csv` (`No` **179–186**, **189**). Modul Must **#177, #178, #187, #188** sudah diwakili **6.F** (`BL-DASH-001`–`004`); cakupan penuh F11: `docs/06-features/11-dashboards.md`.

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F11-CSV-179 | Pantauan Anggaran Iklan (widget) | SH | 723 | S4 | S4-L-03 | BL-DASH-006 | L | todo | S4-J-01..S4-J-02 | Spend vs closings + CPL/CPA konsisten sumber F10 |
| 6 | F11-CSV-180 | Papan Kinerja CS (widget) | SH | 724 | S4 | S4-L-03 | BL-DASH-007 | L | todo | S4-J-01..S4-J-02 | Metrik tim CS + SLA + leaderboard dari CRM |
| 6 | F11-CSV-181 | Radar Bus (live) | SH | 725 | S3 | S3-L-05 | BL-DASH-008 | L | todo | S3-J-01..S3-J-03 | Peta/status armada dari feed GPS/boarding F7 |
| 6 | F11-CSV-182 | Status Raudhah (aggregate) | SH | 726 | S3 | S3-L-05 | BL-DASH-009 | L | todo | S3-J-01..S3-J-03 | % masuk Raudhah per departure + drill-down |
| 6 | F11-CSV-183 | Pelacakan Koper (aggregate) | SH | 727 | S3 | S3-L-05 | BL-DASH-010 | L | todo | S3-J-01..S3-J-03 | Agregat posisi koper per departure dari scan F7 |
| 6 | F11-CSV-184 | Laporan Insiden (feed) | SH | 728 | S3 | S3-L-05 | BL-DASH-011 | L | todo | S3-J-01..S3-J-03 | Feed insiden + filter severitas + notifikasi HQ |
| 6 | F11-CSV-185 | Kesehatan Gudang (widget) | SH | 729 | S3 | S3-L-05 | BL-DASH-012 | L | todo | S3-J-01..S3-J-03 | Nilai stok + chart kritikal vs reorder (read F8) |
| 6 | F11-CSV-186 | Pantauan Eksekusi Logistik (widget) | SH | 730 | S3 | S3-L-05 | BL-DASH-013 | L | todo | S3-J-01..S3-J-03 | Paid-unshipped aging + backlog GRN/PO ringkas |
| 6 | F11-CSV-189 | Likuiditas — AR/AP aging ringkas | SH | 731 | S5 | S5-L-02 | BL-DASH-014 | L | todo | BL-FIN-011 | Bucket aging + alert konsisten dengan sumber F9 |

Modul CSV **#126–#128** (*Executive Dashboard* operasional di kategori *Operational & Handling*) dipecah lagi di **6.K** sebagai `F11-CSV-126`–`128` dengan `BL-DASH-015`–`017` — harmonisasi dengan **BL-DASH-012** / **#185** (*Kesehatan Gudang*) dan **BL-DASH-013** / **#186** (*Pantauan Eksekusi Logistik*) di tabel atas.

### 6.J — Finance modul CSV (#129–#150) (F9)

Sumber prioritas: `docs/Modul UmrohOS - MosCoW.csv`. Cakupan domain: `docs/06-features/09-finance-and-accounting.md`. Baris **6.B** (`BL-FIN-010`, `BL-FIN-011`) tetap umbrella; tabel ini memperinci modul **Finance** di CSV.

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F9-CSV-129 | Penagihan Otomatis | MH-V1 | 732 | S3 | S3-E-07 | BL-FIN-020 | E | todo | S3-J-01..S3-J-03 | Penagihan terjadwal + status piutang konsisten |
| 6 | F9-CSV-130 | Integrasi Bank | MH-V1 | 733 | S3 | S3-E-07 | BL-FIN-021 | E | todo | S3-J-01..S3-J-03 | Mutasi bank terhubung + rekonsiliasi dasar |
| 6 | F9-CSV-131 | Buku Pembantu Piutang | MH-V1 | 734 | S3 | S3-E-07 | BL-FIN-022 | E | todo | S3-J-01..S3-J-03 | Buku pembantu piutang per jamaah/booking |
| 6 | F9-CSV-132 | Kwitansi Digital | MH-V1 | 735 | S3 | S3-E-07 | BL-FIN-023 | E | todo | S3-J-01..S3-J-03 | Kwitansi digital + nomor urut + audit |
| 6 | F9-CSV-133 | Pembayaran Manual & Cek Muka [Tambahan] | SH | 736 | S3 | S3-E-07 | BL-FIN-024 | E | todo | S3-J-01..S3-J-03 | Pencatatan pembayaran manual/cek muka + bukti |
| 6 | F9-CSV-134 | Database Vendor | MH-V1 | 737 | S3 | S3-E-07 | BL-FIN-025 | E | todo | S3-J-01..S3-J-03 | Master vendor + relasi ke AP |
| 6 | F9-CSV-135 | Buku Pembantu Hutang | MH-V1 | 738 | S3 | S3-E-07 | BL-FIN-026 | E | todo | S3-J-01..S3-J-03 | Buku pembantu hutang per vendor |
| 6 | F9-CSV-136 | Otorisasi Pembayaran | MH-V1 | 739 | S3 | S3-E-07 | BL-FIN-027 | E | todo | S3-J-01..S3-J-03 | Otorisasi multi-level sebelum bayar |
| 6 | F9-CSV-137 | Kas Kecil & Bon Sementara | SH | 740 | S3 | S3-E-07 | BL-FIN-028 | E | todo | S3-J-01..S3-J-03 | Kas kecil + bon sementara + tutup periode |
| 6 | F9-CSV-138 | Pembiayaan Berbasis Proyek | MH-V1 | 741 | S3 | S3-E-07 | BL-FIN-029 | E | todo | S3-J-01..S3-J-03 | Cost object per proyek/departure + alokasi |
| 6 | F9-CSV-139 | Laba Rugi Keberangkatan | MH-V1 | 742 | S3 | S3-E-07 | BL-FIN-030 | E | todo | S3-J-01..S3-J-03 | P&L per departure vs anggaran awal |
| 6 | F9-CSV-140 | Analisis Anggaran vs Aktual | SH | 743 | S3 | S3-E-07 | BL-FIN-031 | E | todo | S3-J-01..S3-J-03 | Variance anggaran vs aktual + drill-down |
| 6 | F9-CSV-141 | Jurnal Otomatis | MH-V1 | 744 | S3 | S3-E-07 | BL-FIN-032 | E | todo | S3-J-01..S3-J-03 | Jurnal otomatis idempotent per sumber event |
| 6 | F9-CSV-142 | Pengakuan Pendapatan | MH-V1 | 745 | S3 | S3-E-07 | BL-FIN-033 | E | todo | S3-J-01..S3-J-03 | Pengakuan pendapatan sesuai policy Q043 |
| 6 | F9-CSV-143 | Multi-Currency | SH | 746 | S3 | S3-E-07 | BL-FIN-034 | E | todo | S3-J-01..S3-J-03 | Multi-mata uang + rate snapshot konsisten Q001 |
| 6 | F9-CSV-144 | Manajemen Aset Tetap | SH | 747 | S3 | S3-E-07 | BL-FIN-035 | E | todo | S3-J-01..S3-J-03 | Kartu aset + penyusutan dasar |
| 6 | F9-CSV-145 | Integrasi Pajak | SH | 748 | S3 | S3-E-07 | BL-FIN-036 | E | todo | S3-J-01..S3-J-03 | Pajak terintegrasi (PPN/PPh) sesuai Q046/Q047 |
| 6 | F9-CSV-146 | Pencairan Komisi Agen | MH-V1 | 749 | S3 | S3-E-07 | BL-FIN-037 | E | todo | S3-J-01..S3-J-03 | Alur pencairan komisi agen + audit |
| 6 | F9-CSV-147 | Laporan Keuangan Real-Time | MH-V1 | 750 | S3 | S3-E-07 | BL-FIN-038 | E | todo | S5-J-01..S5-J-02 | Laporan real-time (neraca/P&L ringkas) |
| 6 | F9-CSV-148 | Dashboard Arus Cash | MH-V1 | 751 | S3 | S3-E-07 | BL-FIN-039 | E | todo | S5-J-01..S5-J-02 | Dashboard arus kas + saldo bank/petty |
| 6 | F9-CSV-149 | Peringatan Umur Hutang/Piutang | SH | 752 | S3 | S3-E-07 | BL-FIN-040 | E | todo | S5-J-01..S5-J-02 | Alert umur piutang/hutang + bucket |
| 6 | F9-CSV-150 | Jejak Audit & Hak Akses | MH-V1 | 753 | S3 | S3-E-07 | BL-FIN-041 | E | todo | S5-J-01..S5-J-02 | Audit trail finance + RBAC anti-delete |

### 6.K — Operational & Handling modul CSV (#87–#128) (F7/F8)

Sumber prioritas: `docs/Modul UmrohOS - MosCoW.csv`. **#87–#108** dipetakan ke **F7** (`docs/06-features/07-operations-handling.md` + dokumen jamaah terkait); **#109–#125** ke **F8** (`docs/06-features/08-warehouse-and-fulfillment.md`); **#126–#128** adalah widget *Executive Dashboard* operasional (`BL-DASH-015`–`017`; harmonisasi dengan widget serupa di **6.I** bila ada overlap nama).

Baris **6.C**/**6.D**/**6.E** tetap umbrella; tabel ini memperinci modul CSV **Operational & Handling**.

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F7-CSV-87 | Penyimpanan Kolektif | MH-V1 | 754 | S3 | S3-E-04 | BL-OPS-021 | E | todo | S3-J-01..S3-J-03 | Vault kolektif per departure + ACL |
| 6 | F7-CSV-88 | OCR Paspor & Mahram Logic | MH-V1 | 755 | S3 | S3-E-04 | BL-OPS-022 | E | todo | S3-J-01..S3-J-03 | OCR paspor + flag mahram sesuai aturan |
| 6 | F7-CSV-89 | Progress Tracker & Expiry Alert | SH | 756 | S3 | S3-E-04 | BL-OPS-023 | E | todo | S3-J-01..S3-J-03 | Progress dokumen + expiry alert |
| 6 | F7-CSV-90 | Generator Surat Resmi | MH-V1 | 757 | S3 | S3-E-04 | BL-OPS-024 | E | todo | S3-J-01..S3-J-03 | Surat resmi generate dari template |
| 6 | F7-CSV-91 | Manifest Imigrasi | MH-V1 | 758 | S3 | S3-E-04 | BL-OPS-025 | E | todo | S3-J-01..S3-J-03 | Manifest imigrasi format + versioning |
| 6 | F7-CSV-92 | Algoritma Penempatan Kamar | MH-V1 | 759 | S3 | S3-E-04 | BL-OPS-026 | E | todo | S3-J-01..S3-J-03 | Algoritma rooming + constraint valid |
| 6 | F7-CSV-93 | Pengatur Transportasi | SH | 760 | S3 | S3-E-04 | BL-OPS-027 | E | todo | S3-J-01..S3-J-03 | Penempatan seat/transportasi grup |
| 6 | F7-CSV-94 | Handling Manifest | SH | 761 | S3 | S3-E-04 | BL-OPS-028 | E | todo | S3-J-01..S3-J-03 | Handling manifest + delta publish |
| 6 | F7-CSV-95 | ID Card & Staff Assignment | MH-V1 | 762 | S3 | S3-E-04 | BL-OPS-029 | E | todo | S3-J-01..S3-J-03 | ID card + assignment staff tercetak |
| 6 | F7-CSV-96 | Log Fisik Paspor | SH | 763 | S3 | S3-E-04 | BL-OPS-030 | E | todo | S3-J-01..S3-J-03 | Log fisik paspor + handover audit |
| 6 | F7-CSV-97 | Visa Progress Tracker | MH-V1 | 764 | S3 | S3-E-04 | BL-OPS-031 | E | todo | BL-VISA-003 | Status visa per jamaah + SLA |
| 6 | F7-CSV-98 | E-Visa Repository | SH | 765 | S3 | S3-E-04 | BL-OPS-032 | E | todo | BL-VISA-003 | Repository e-visa + metadata |
| 6 | F7-CSV-99 | Integrasi API Eksternal Lanjutan | CH | 766 | S3 | S3-E-04 | BL-OPS-033 | E | todo | BL-VISA-003 | Konektor provider lanjutan (optional) |
| 6 | F7-CSV-100 | Administrasi Refund & Pinalti | SH | 767 | S3 | S3-E-04 | BL-OPS-034 | E | todo | S3-J-01..S3-J-03 | Admin refund + pinalti + approval |
| 6 | F7-CSV-101 | ALL System | MH-V1 | 768 | S3 | S3-E-04 | BL-OPS-035 | E | todo | S3-J-01..S3-J-03 | ALL system scan + idempotensi |
| 6 | F7-CSV-102 | Penghitung Koper | SH | 769 | S3 | S3-E-04 | BL-OPS-036 | E | todo | S3-J-01..S3-J-03 | Penghitung koper + event scan |
| 6 | F7-CSV-103 | Broadcast Keberangkatan & Kedatangan | SH | 770 | S3 | S3-E-04 | BL-OPS-037 | E | todo | S3-J-01..S3-J-03 | Broadcast jadwal berangkat/datang |
| 6 | F7-CSV-104 | Smart Bus Boarding | MH-V1 | 771 | S3 | S3-E-04 | BL-OPS-038 | E | todo | S3-J-01..S3-J-03 | Boarding bus + roster |
| 6 | F7-CSV-105 | Raudhah Shield & Tasreh Digital | SH | 772 | S3 | S3-E-04 | BL-OPS-039 | E | todo | S3-J-01..S3-J-03 | Raudhah shield + tasreh digital |
| 6 | F7-CSV-106 | Manajemen Perangkat Audio | CH | 773 | S3 | S3-E-04 | BL-OPS-040 | E | todo | S3-J-01..S3-J-03 | Inventaris perangkat audio lapangan |
| 6 | F7-CSV-107 | Distribusi Zamzam | SH | 774 | S3 | S3-E-04 | BL-OPS-041 | E | todo | S3-J-01..S3-J-03 | Distribusi zamzam + bukti serah |
| 6 | F7-CSV-108 | Check-In Kamar Cepat [Tambahan] | CH | 775 | S3 | S3-E-04 | BL-OPS-042 | E | todo | S3-J-01..S3-J-03 | Check-in kamar cepat (add-on) |
| 6 | F8-CSV-109 | Permintaan Pembelian | MH-V1 | 776 | S3 | S3-E-05 | BL-LOG-013 | E | todo | S3-J-01..S3-J-03 | PR + approval + link anggaran |
| 6 | F8-CSV-110 | Sinkronisasi Anggaran | MH-V1 | 777 | S3 | S3-E-05 | BL-LOG-014 | E | todo | S3-J-01..S3-J-03 | Sinkron anggaran PR vs aktual |
| 6 | F8-CSV-111 | Persetujuan Berjenjang | SH | 778 | S3 | S3-E-05 | BL-LOG-015 | E | todo | S3-J-01..S3-J-03 | Approval berjenjang PR/PO |
| 6 | F8-CSV-112 | Otomasi Vendor | SH | 779 | S3 | S3-E-05 | BL-LOG-016 | E | todo | S3-J-01..S3-J-03 | Otomasi pilihan vendor (rule) |
| 6 | F8-CSV-113 | Goods Receipt | MH-V1 | 780 | S3 | S3-E-05 | BL-LOG-017 | E | todo | S3-J-01..S3-J-03 | GRN + partial + reversal policy |
| 6 | F8-CSV-114 | Quality Control | SH | 781 | S3 | S3-E-05 | BL-LOG-018 | E | todo | S3-J-01..S3-J-03 | QC inbound + status reject |
| 6 | F8-CSV-115 | Pemicu Hutang Otomatis | MH-V1 | 782 | S3 | S3-E-05 | BL-LOG-019 | E | todo | S3-J-01..S3-J-03 | Posting hutang otomatis dari GRN |
| 6 | F8-CSV-116 | Pelabelan Barcode/SKU | SH | 783 | S3 | S3-E-05 | BL-LOG-020 | E | todo | S3-J-01..S3-J-03 | Label barcode/SKU konsisten master |
| 6 | F8-CSV-117 | Multi-Warehouse | SH | 784 | S3 | S3-E-05 | BL-LOG-021 | E | todo | S3-J-01..S3-J-03 | Multi-warehouse + transfer |
| 6 | F8-CSV-118 | Peringatan Stok Kritis | SH | 785 | S3 | S3-E-05 | BL-LOG-022 | E | todo | S3-J-01..S3-J-03 | Alert stok di bawah reorder |
| 6 | F8-CSV-119 | Perakitan Paket | SH | 786 | S3 | S3-E-05 | BL-LOG-023 | E | todo | S3-J-01..S3-J-03 | Perakitan kit all-or-nothing |
| 6 | F8-CSV-120 | Stock Opname Digital | SH | 787 | S3 | S3-E-05 | BL-LOG-024 | E | todo | S3-J-01..S3-J-03 | Stock opname digital + variance |
| 6 | F8-CSV-121 | Sinkronisasi Ukuran | MH-V1 | 788 | S3 | S3-E-05 | BL-LOG-025 | E | todo | S3-J-01..S3-J-03 | Ukuran/preset fulfillment sinkron katalog |
| 6 | F8-CSV-122 | Pemicu Pengiriman | MH-V1 | 789 | S3 | S3-E-05 | BL-LOG-026 | E | todo | S3-J-01..S3-J-03 | Pemicu kirim setelah paid-in-full |
| 6 | F8-CSV-123 | Integrasi Ekspedisi | SH | 790 | S3 | S3-E-05 | BL-LOG-027 | E | todo | S3-J-01..S3-J-03 | Integrasi ekspedisi + resi |
| 6 | F8-CSV-124 | Pengambilan Mandiri | SH | 791 | S3 | S3-E-05 | BL-LOG-028 | E | todo | S3-J-01..S3-J-03 | Pickup mandiri + QR single-use |
| 6 | F8-CSV-125 | Retur & Penukaran Barang [Tambahan] | CH | 792 | S3 | S3-E-05 | BL-LOG-029 | E | todo | S3-J-01..S3-J-03 | Retur/penukaran (add-on) |
| 6 | F11-CSV-126 | Inventory Health | SH | 793 | S5 | S5-L-02 | BL-DASH-015 | L | todo | S5-J-01..S5-J-02 | Widget inventory health (read F8) |
| 6 | F11-CSV-127 | Fulfillment & PO Monitor | SH | 794 | S5 | S5-L-02 | BL-DASH-016 | L | todo | S5-J-01..S5-J-02 | Monitor PO + fulfillment backlog |
| 6 | F11-CSV-128 | Laporan Kerusakan | CH | 795 | S5 | S5-L-02 | BL-DASH-017 | L | todo | S5-J-01..S5-J-02 | Laporan kerusakan barang (add-on) |


### 6.L — B2C Front-End modul CSV (#1–#24)

Sumber prioritas: `docs/Modul UmrohOS - MosCoW.csv`. Cakupan: situs publik B2C — selaraskan dengan `docs/06-features/02-catalog-and-master-data.md`, `docs/06-features/04-booking-and-allocation.md`, `docs/06-features/05-payment-and-reconciliation.md`.

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | B2C-CSV-01 | Modul Homepage Dinamis | MH-V1 | 796 | S1 | S1-L-05 | BL-B2C-001 | L | todo | S1-J-01..S1-J-04 | Homepage dinamis + performa mobile |
| 6 | B2C-CSV-02 | Modul Validasi Legalitas & About Us | MH-V1 | 797 | S1 | S1-L-05 | BL-B2C-002 | L | todo | S1-J-01..S1-J-04 | Legalitas & About Us terverifikasi tampil |
| 6 | B2C-CSV-03 | Modul Galeri & Social Proof | SH | 798 | S1 | S1-L-05 | BL-B2C-003 | L | todo | S1-J-01..S1-J-04 | Galeri + bukti sosial (testimoni) |
| 6 | B2C-CSV-04 | Modul Blog & Artikel | CH | 799 | S1 | S1-L-05 | BL-B2C-004 | L | todo | S1-J-01..S1-J-04 | Blog/list artikel + SEO dasar |
| 6 | B2C-CSV-05 | Modul Brand Identity & White-labeling | SH | 800 | S1 | S1-L-05 | BL-B2C-005 | L | todo | S1-J-01..S1-J-04 | White-label tema + asset brand |
| 6 | B2C-CSV-06 | Modul Menu Builder & Navigation | MH-V1 | 801 | S1 | S1-L-05 | BL-B2C-006 | L | todo | S1-J-01..S1-J-04 | Menu navigasi terkonfigurasi + ACL |
| 6 | B2C-CSV-07 | Modul Smart Search & Filter | MH-V1 | 802 | S1 | S1-L-05 | BL-B2C-007 | L | todo | S1-J-01..S1-J-04 | Search/filter paket konsisten katalog |
| 6 | B2C-CSV-08 | Modul Detail Produk & Itinerary Interaktif | MH-V1 | 803 | S1 | S1-L-05 | BL-B2C-008 | L | todo | S1-J-01..S1-J-04 | Detail produk + itinerary interaktif |
| 6 | B2C-CSV-09 | Modul Ketersediaan Real-Time | MH-V1 | 804 | S1 | S1-L-05 | BL-B2C-009 | L | todo | S1-J-01..S1-J-04 | Seat/ketersediaan real-time dari katalog |
| 6 | B2C-CSV-10 | Modul Kalkulator Simulasi Tabungan | CH | 805 | S1 | S1-L-05 | BL-B2C-010 | L | todo | S1-J-01..S1-J-04 | Simulasi tabungan (add-on) + disclaimer |
| 6 | B2C-CSV-11 | Info Esensial & Seat Tracker | MH-V1 | 806 | S1 | S1-L-05 | BL-B2C-011 | L | todo | S1-J-01..S1-J-04 | Info esensial + seat tracker publik |
| 6 | B2C-CSV-12 | Spesifikasi Akomodasi Cerdas | SH | 807 | S1 | S1-L-05 | BL-B2C-012 | L | todo | S1-J-01..S1-J-04 | Spesifikasi akomodasi ringkas per paket |
| 6 | B2C-CSV-13 | Profil Pembimbing | SH | 808 | S1 | S1-L-05 | BL-B2C-013 | L | todo | S1-J-01..S1-J-04 | Profil pembimbing tampil per paket |
| 6 | B2C-CSV-14 | Micro-Web Itinerary | SH | 809 | S1 | S1-L-05 | BL-B2C-014 | L | todo | S1-J-01..S1-J-04 | Micro-web itinerary shareable |
| 6 | B2C-CSV-15 | Call-to-Action | MH-V1 | 810 | S1 | S1-L-05 | BL-B2C-015 | L | todo | S1-J-01..S1-J-04 | CTA ke WA/booking konsisten tracking |
| 6 | B2C-CSV-16 | Self-Booking Engine | MH-V1 | 811 | S1 | S1-L-05 | BL-B2C-016 | L | todo | S1-J-01..S1-J-04 | Alur self-booking end-to-end B2C |
| 6 | B2C-CSV-17 | Guest Data Form | MH-V1 | 812 | S1 | S1-L-05 | BL-B2C-017 | L | todo | S1-J-01..S1-J-04 | Form data tamu + validasi field |
| 6 | B2C-CSV-18 | Payment Gateway B2C | MH-V1 | 813 | S2 | S2-L-04 | BL-B2C-018 | L | todo | S2-J-01..S2-J-04 | Checkout B2C terhubung payment slice |
| 6 | B2C-CSV-19 | Riwayat Keuangan | SH | 814 | S1 | S1-L-05 | BL-B2C-019 | L | todo | S1-J-01..S1-J-04 | Riwayat transaksi pembayaran jamaah |
| 6 | B2C-CSV-20 | Self-Upload Dokumen | SH | 815 | S3 | S3-L-03 | BL-B2C-020 | L | todo | BL-DOC-001 | Upload dokumen self-service + status |
| 6 | B2C-CSV-21 | Informasi Logistik & Perlengkapan | CH | 816 | S1 | S1-L-05 | BL-B2C-021 | L | todo | S1-J-01..S1-J-04 | Info logistik/perlengkapan (add-on) |
| 6 | B2C-CSV-22 | Papan Informasi Keberangkatan | SH | 817 | S1 | S1-L-05 | BL-B2C-022 | L | todo | S1-J-01..S1-J-04 | Papan info keberangkatan (read) |
| 6 | B2C-CSV-23 | Knowledge Base | CH | 818 | S1 | S1-L-05 | BL-B2C-023 | L | todo | S1-J-01..S1-J-04 | Knowledge base browse/search |
| 6 | B2C-CSV-24 | Floating Chat | SH | 819 | S1 | S1-L-05 | BL-B2C-024 | L | todo | S1-J-01..S1-J-04 | Floating chat + routing channel |

### 6.M — Admin & Security modul CSV (#151–#163) (F1)

Sumber prioritas: `docs/Modul UmrohOS - MosCoW.csv`. Domain: `docs/06-features/01-identity-and-access.md` (RBAC, audit, konfigurasi).

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F1-CSV-151 | Pembuatan Peran Jabatan | MH-V1 | 820 | S1 | S1-E-06 | BL-IAM-005 | E | todo | S1-J-01..S1-J-04 | CRUD peran jabatan + mapping IAM |
| 6 | F1-CSV-152 | Pemetaan Izin Spesifik | MH-V1 | 821 | S1 | S1-E-06 | BL-IAM-006 | E | todo | S1-J-01..S1-J-04 | Izin granular per route/aksi |
| 6 | F1-CSV-153 | Hierarki Visibilitas Data | MH-V1 | 822 | S1 | S1-E-06 | BL-IAM-007 | E | todo | S1-J-01..S1-J-04 | Hierarki scope data (global/cabang) |
| 6 | F1-CSV-154 | Pendaftaran Akun Staf | MH-V1 | 823 | S1 | S1-E-06 | BL-IAM-008 | E | todo | S1-J-01..S1-J-04 | Onboarding akun staf + invite |
| 6 | F1-CSV-155 | Kontrol Status Pengguna | MH-V1 | 824 | S1 | S1-E-06 | BL-IAM-009 | E | todo | S1-J-01..S1-J-04 | Suspend/aktif user + alasan audit |
| 6 | F1-CSV-156 | Keamanan Akun & Sandi | MH-V1 | 825 | S1 | S1-E-06 | BL-IAM-010 | E | todo | S1-J-01..S1-J-04 | Kebijakan sandi + MFA opsional |
| 6 | F1-CSV-157 | Log Aktivitas Terpusat | MH-V1 | 826 | S1 | S1-E-06 | BL-IAM-011 | E | todo | S1-J-01..S1-J-04 | Log aktivitas terpusat searchable |
| 6 | F1-CSV-158 | Peringatan Anomali | SH | 827 | S1 | S1-E-06 | BL-IAM-012 | E | todo | S1-J-01..S1-J-04 | Alert anomali login/aksi (SH) |
| 6 | F1-CSV-159 | Riwayat Sesi Pengguna | SH | 828 | S1 | S1-E-06 | BL-IAM-013 | E | todo | S1-J-01..S1-J-04 | Riwayat sesi + revoke |
| 6 | F1-CSV-160 | Konfigurasi Integrasi API | MH-V1 | 829 | S1 | S1-E-06 | BL-IAM-014 | E | todo | S1-J-01..S1-J-04 | Konfigurasi kunci API + rotasi |
| 6 | F1-CSV-161 | Manajemen Template Komunikasi | SH | 830 | S1 | S1-E-06 | BL-IAM-015 | E | todo | S1-J-01..S1-J-04 | Template WA/email terpusat |
| 6 | F1-CSV-162 | Konfigurasi Global Variabel | MH-V1 | 831 | S1 | S1-E-06 | BL-IAM-016 | E | todo | S1-J-01..S1-J-04 | Global config key-value + audit |
| 6 | F1-CSV-163 | Pencadangan Database | MH-V1 | 832 | S1 | S1-E-06 | BL-IAM-017 | E | todo | S1-J-01..S1-J-04 | Jadwal/restore backup DB (prosedur) |

### 6.N — Jamaah Journey modul CSV (#164–#176) (F12)

Sumber prioritas: `docs/Modul UmrohOS - MosCoW.csv`. Domain: `docs/06-features/12-alumni-and-daily-app.md` (pengalaman jamaah in-trip / post-booking).

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F12-CSV-164 | Jadwal & Informasi Langsung | MH-V1 | 833 | S3 | S3-L-04 | BL-JMJ-001 | L | todo | S3-J-01..S3-J-03 | Jadwal live + push info per trip |
| 6 | F12-CSV-165 | Pengingat Pintar | SH | 834 | S3 | S3-L-04 | BL-JMJ-002 | L | todo | S3-J-01..S3-J-03 | Pengingat pintar (itinerary) |
| 6 | F12-CSV-166 | Panduan Ibadah Digital | SH | 835 | S3 | S3-L-04 | BL-JMJ-003 | L | todo | S3-J-01..S3-J-03 | Panduan ibadah offline-friendly |
| 6 | F12-CSV-167 | Dompet Dokumen Digital | MH-V1 | 836 | S3 | S3-L-04 | BL-JMJ-004 | L | todo | S3-J-01..S3-J-03 | Dompet dokumen digital per jamaah |
| 6 | F12-CSV-168 | Tombol Darurat | SH | 837 | S3 | S3-L-04 | BL-JMJ-005 | L | todo | S3-J-01..S3-J-03 | Tombol darurat + eskalasi |
| 6 | F12-CSV-169 | E-Certificate Generator | CH | 838 | S3 | S3-L-04 | BL-JMJ-006 | L | todo | S3-J-01..S3-J-03 | E-certificate generate (add-on) |
| 6 | F12-CSV-170 | Absensi Naik Bus | MH-V1 | 839 | S3 | S3-L-04 | BL-JMJ-007 | L | todo | S3-J-01..S3-J-03 | Absensi naik bus + scan |
| 6 | F12-CSV-171 | Kontrol Waktu Bus | SH | 840 | S3 | S3-L-04 | BL-JMJ-008 | L | todo | S3-J-01..S3-J-03 | Kontrol waktu bus + SLA |
| 6 | F12-CSV-172 | Pelacakan Koper Bandara | SH | 841 | S3 | S3-L-04 | BL-JMJ-009 | L | todo | S3-J-01..S3-J-03 | Lacak koper bandara (status) |
| 6 | F12-CSV-173 | Manajemen Alat Komunikasi | CH | 842 | S3 | S3-L-04 | BL-JMJ-010 | L | todo | S3-J-01..S3-J-03 | Manajemen alat komunikasi (add-on) |
| 6 | F12-CSV-174 | Pembagian Zamzam | SH | 843 | S3 | S3-L-04 | BL-JMJ-011 | L | todo | S3-J-01..S3-J-03 | Distribusi zamzam + bukti |
| 6 | F12-CSV-175 | Pelaporan Isu Harian | SH | 844 | S3 | S3-L-04 | BL-JMJ-012 | L | todo | S3-J-01..S3-J-03 | Pelaporan isu harian |
| 6 | F12-CSV-176 | Edukasi Mitigasi | CH | 845 | S3 | S3-L-04 | BL-JMJ-013 | L | todo | S3-J-01..S3-J-03 | Edukasi mitigasi (add-on) |

### 6.O — Fitur Pelengkap & Daily App modul CSV (#190–#198) (F12)

Sumber prioritas: `docs/Modul UmrohOS - MosCoW.csv`. Modul **#199–#202** (referral/ZISWAF) sudah di **6.H**. Baris berikut: konten & komunitas harian — `docs/06-features/12-alumni-and-daily-app.md`.

| Phase | Ref detail fitur | Ringkasan detail | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by Gate | Acceptance ringkas |
|------:|---|---|:---:|---:|---|---|---|---|---|---|---|
| 6 | F12-CSV-190 | Jadwal Shalat & Adzan | SH | 846 | S5 | S5-L-03 | BL-PLG-001 | L | todo | S5-J-01..S5-J-02 | Jadwal shalat + notifikasi adzan |
| 6 | F12-CSV-191 | Kompas Arah Kiblat | SH | 847 | S5 | S5-L-03 | BL-PLG-002 | L | todo | S5-J-01..S5-J-02 | Kompas kiblat akurat per lokasi |
| 6 | F12-CSV-192 | Al-Quran Digital | CH | 848 | S5 | S5-L-03 | BL-PLG-003 | L | todo | S5-J-01..S5-J-02 | Al-Quran digital (add-on) |
| 6 | F12-CSV-193 | Dzikir & Kumpulan Doa | CH | 849 | S5 | S5-L-03 | BL-PLG-004 | L | todo | S5-J-01..S5-J-02 | Dzikir & doa harian |
| 6 | F12-CSV-194 | Ensiklopedia Manasik | SH | 850 | S5 | S5-L-03 | BL-PLG-005 | L | todo | S5-J-01..S5-J-02 | Ensiklopedia manasik searchable |
| 6 | F12-CSV-195 | Artikel & Kajian Routine | CH | 851 | S5 | S5-L-03 | BL-PLG-006 | L | todo | S5-J-01..S5-J-02 | Artikel/kajian rutin |
| 6 | F12-CSV-196 | Tanya Jawab Agama | CH | 852 | S5 | S5-L-03 | BL-PLG-007 | L | todo | S5-J-01..S5-J-02 | Tanya jawab agama (moderated) |
| 6 | F12-CSV-197 | Forum Grup Angkatan | CH | 853 | S5 | S5-L-03 | BL-PLG-008 | L | todo | S5-J-01..S5-J-02 | Forum grup angkatan |
| 6 | F12-CSV-198 | Papan Informasi Reuni | CH | 854 | S5 | S5-L-03 | BL-PLG-009 | L | todo | S5-J-01..S5-J-02 | Papan reuni & pengumuman |

---

## Aturan pakai harian

1. Saat buka pekerjaan, sebut minimal: `Backlog ID + Task Code` (jika sudah dipetakan ke slice).
2. Jika satu backlog item terasa terlalu besar untuk satu PR, pecah jadi item baru (`BL-...-next`) sebelum coding lanjut.
3. Perubahan priority harus dilakukan di tabel ini dulu, baru turun ke board ticket.
4. Kolom `Blocked by Gate` berisi gate kontrak yang harus selesai dulu; jika belum selesai, status item tetap `todo`.
5. Baris Fase 6 memakai **`Slice` / `Task Code` dari paket depth** di [`05-slice-engineering-checklist-and-task-codes.md`](./05-slice-engineering-checklist-and-task-codes.md) (§ **Fase 6 — Depth backlog**). Satu kode tugas boleh menampung banyak `BL-*`; pecah subtugas jika PR membengkak.

Contoh perintah:

> Kerjakan `BL-PAY-004` pada `S2-E-02`. Jangan ubah kontrak di luar `slice-S2.md`.

---

## Referensi

- [`04-delivery-plan-2p-sequence-first.md`](./04-delivery-plan-2p-sequence-first.md) — urutan fase & slice
- [`05-slice-engineering-checklist-and-task-codes.md`](./05-slice-engineering-checklist-and-task-codes.md) — gate & kode tugas
- `docs/06-features/01-identity-and-access.md`
- `docs/06-features/02-catalog-and-master-data.md`
- `docs/06-features/03-pilgrim-and-documents.md`
- `docs/06-features/04-booking-and-allocation.md`
- `docs/06-features/05-payment-and-reconciliation.md`
- `docs/06-features/06-visa-pipeline.md`
- `docs/06-features/07-operations-handling.md`
- `docs/06-features/08-warehouse-and-fulfillment.md`
- `docs/06-features/09-finance-and-accounting.md`
- `docs/06-features/10-marketing-crm-agents.md`
- `docs/06-features/11-dashboards.md`
- `docs/06-features/12-alumni-and-daily-app.md`
- `docs/Modul UmrohOS - MosCoW.csv`
