# Slice Engineering — Checklist & Kode Pekerjaan (Lutfi + Elda)

Dokumen ini melengkapi [`04-delivery-plan-2p-sequence-first.md`](./04-delivery-plan-2p-sequence-first.md) (urutan fase & slice). Daftar baris backlog `BL-*` + `Exec seq` ada di [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md).

**Tujuan:** menuntaskan **dependency engineering** *sebelum* coding slice dimulai, supaya paralel kerja minim bentrok dan minim rework.

---

## Format kode pekerjaan (untuk perintah ke AI / ke rekan)

```
S{slice}-{owner}-{seq}
```

| Segmen | Arti | Nilai |
|--------|------|--------|
| `S{slice}` | Slice vertikal | `S0` … `S5` (lihat bawah) |
| `{owner}` | Pemilik eksekusi | `L` = Lutfi, `E` = Elda, `J` = Joint (keduanya wajib hadir / merge bersama) |
| `{seq}` | Urut 2 digit | `01`, `02`, … |

**Contoh perintah ke AI:**

> Kerjakan **S1-L-04** dan **S1-E-03**. Patuhi kontrak di `docs/contracts/slice-S1.md` (jika sudah ada). Jangan ubah kontrak tanpa ACC.

**Saran tambahan (opsional tapi bagus):**

- Satu file kontrak per slice: `docs/contracts/slice-Sx.md` (Markdown + contoh JSON) — *folder `docs/contracts/` boleh dibuat saat S0-J-01 selesai.*
- Issue tracker (GitHub/Linear): judul issue = kode tugas, supaya 1:1 dengan checklist ini.

---

## Aturan prioritas vs slice (yang benar)

Gunakan label berikut untuk semua **feature/backlog item**:

- `MH-MVP`: wajib selesai untuk MVP tercepat
- `MH-V1`: wajib selesai untuk rilis v1 (setelah MVP)
- `SH`: penting, dikerjakan setelah semua `MH-MVP` + `MH-V1` aman
- `CH`: opsional, dikerjakan jika kapasitas tersedia

Prinsip:

- **Priority itu milik feature/backlog item**, bukan milik slice.
- **Slice itu urutan delivery** (wadah eksekusi end-to-end).
- **Task code** adalah paket kerja teknis yang mengeksekusi beberapa backlog item.

Konsekuensi:

- Satu task code boleh berisi backlog item dengan label prioritas berbeda, tapi disarankan dipisah (`-a`, `-b`) agar eksekusi fokus.
- Penandaan `MH-MVP/MH-V1/SH/CH` dilakukan di dokumen mapping backlog (feature-level), bukan dipatok default per slice.

---

## Definisi slice

| Kode | Nama slice | User journey (ringkas) |
|------|--------------|-------------------------|
| **S0** | Engineering bootstrap | Konvensi repo, CI, format kontrak, ownership merge |
| **S1** | Discover + draft booking | B2C: katalog → detail → booking form → **draft** (tanpa login jamaah); staff/internal pakai auth sesuai F1 |
| **S2** | Get paid | Draft → invoice/VA → webhook → `pending_payment` / partial / lunas |
| **S3** | Fulfill + keuangan minimum | Lunas → tugas fulfillment → posting jurnal / keuangan dasar |
| **S4** | Growth loop | Lead + attribution + hook komisi/CRM read |
| **S5** | Hardening | UAT, perf, security checklist, freeze |

---

## Aturan “slice boleh mulai coding”

Sebuah slice `Sx` (untuk `x ≥ 1`) status **READY TO BUILD** hanya jika:

1. Semua baris di tabel **Engineering freeze** untuk slice itu = **centang** (tidak ada gate yang masih terbuka).
2. Artefak kontrak untuk slice itu sudah ada (minimal Markdown) dan **sudah di-review** oleh owner non-eksekutor (Lutfi review Elda, sebaliknya).
3. Tidak ada item **[GATE]** terbuka di slice sebelumnya (kecuali secara eksplisit di-*waive* dengan catatan tanggal).

---

# S0 — Engineering bootstrap (sebelum S1)

**Tujuan:** satu bahasa untuk kontrak, merge, dan kualitas gate.

## Checklist — Joint **[GATE]**

| Kode | Owner | Pekerjaan | Output / bukti selesai |
|------|-------|-----------|-------------------------|
| S0-J-01 | J | Buat folder & template kontrak `docs/contracts/README.md` + `slice-Sx.md` template | Folder + template ter-merge |
| S0-J-02 | J | Sepakat **branch strategy** (mis. `main` + short-lived `feat/*`) + aturan “siapa merge” | 1 paragraf di README kontrak atau di wiki internal |
| S0-J-03 | J | Definisi **Definition of Ready (DoR)** & **Definition of Done (DoD)** per PR | Tabel singkat di README kontrak |
| S0-L-01 | L | Daftar **role + route UI** yang akan disentuh S1 (publik vs internal) | Tabel role vs URL |
| S0-E-01 | E | Daftar **service** yang disentuh S1–S2 + ownership file (siapa PR owner) | Tabel service vs owner |

---

# S1 — Discover + draft booking

## Engineering freeze (WAJIB sebelum coding fitur S1) **[GATE]**

| Kode | Owner | Pekerjaan | Output / bukti selesai |
|------|-------|-----------|-------------------------|
| S1-J-01 | J | **Kontrak API publik** katalog: list paket, detail paket, detail departure + sisa seat (read) | `docs/contracts/slice-S1.md` § Catalog |
| S1-J-02 | J | **Kontrak API** `POST /v1/bookings` (draft): field wajib, error shape, idempotency key (jika ada) | `docs/contracts/slice-S1.md` § Booking |
| S1-J-03 | J | **Kontrak** `ReserveSeats` / `ReleaseSeats` (gRPC atau REST internal): parameter, failure codes, kompensasi | `docs/contracts/slice-S1.md` § Inventory |
| S1-J-04 | J | **State** booking yang dipakai S1: minimal `draft` (dokumen lengkap **belum** wajib di S1 jika MVP kamu memang baru KTP+passport di gate berikutnya — tetap tulis eksplisit) | 1 paragraf keputusan + referensi Q006 |
| S1-E-01 | E | Review **konkurensi seat** + transaksi DB (statement atomic) pada kontrak S1-J-03 | Comment “approved” di dokumen kontrak atau PR review |
| S1-L-01 | L | Wireframe / daftar screen S1 (URL + komponen utama) | Link figma *atau* bullet di kontrak |

## Checklist implementasi (setelah freeze)

| Kode | Owner | Pekerjaan | Depends on |
|------|-------|-----------|------------|
| S1-L-02 | L | UI katalog + detail + alur ke form booking | S1-J-01 |
| S1-L-03 | L | Integrasi client → API katalog | S1-J-01 |
| S1-L-04 | L | Integrasi client → create draft booking | S1-J-02 |
| S1-E-02 | E | `catalog-svc` read endpoints sesuai kontrak | S1-J-01 |
| S1-E-03 | E | `booking-svc` draft + orkestrasi reserve seat sesuai kontrak | S1-J-02, S1-J-03 |
| S1-E-04 | E | Middleware auth **internal** untuk route admin/CS yang dipakai tes (jika S1 butuh) | S0-L-01 |

---

# S2 — Get paid (VA + webhook)

## Engineering freeze **[GATE]**

| Kode | Owner | Pekerjaan | Output / bukti selesai |
|------|-------|-----------|-------------------------|
| S2-J-01 | J | Kontrak `POST` invoice + `POST` VA issue: `amount_total`, `currency`, `fx_snapshot`, TTL | `docs/contracts/slice-S2.md` |
| S2-J-02 | J | Kontrak webhook: header signature, body minimal, dedupe key, response codes | `slice-S2.md` § Webhook |
| S2-J-03 | J | Kontrak callback ke booking: status transition + idempotensi | `slice-S2.md` § Booking integration |
| S2-J-04 | J | **Stub** `payment-svc` (response tetap sesuai kontrak) atau toggle `MOCK_GATEWAY` | Stub merge / env flag dokumentasi |
| S2-E-01 | E | Tabel DB invoice/events sesuai `docs/03-services/04-payment-svc/02-data-model.md` | Migrasi direview |
| S2-L-01 | L | UI checkout: menampilkan VA/QR + status polling strategy | Deskripsi di kontrak atau komentar UI |

## Checklist implementasi

| Kode | Owner | Pekerjaan | Depends on |
|------|-------|-----------|------------|
| S2-E-02 | E | Implement `payment-svc` invoice + VA + webhook | S2-J-01–J-04 |
| S2-E-03 | E | Reconcile cron minimal | S2-J-02 |
| S2-L-02 | L | Halaman checkout + error UX | S2-J-01 |
| S2-L-03 | L | Wiring booking flow → panggil payment | S2-J-03 |
| S2-L-04 | L | Checkout B2C mendalam (VA/QR + error UX lanjutan; `BL-B2C-018`) | S2-L-02 |
| S2-J-05 | J | Uji end-to-end: stub lalu gateway nyata | S2-E-02, S2-L-03 |

---

# S3 — Fulfillment + keuangan minimum

## Engineering freeze **[GATE]**

| Kode | Owner | Pekerjaan | Output |
|------|-------|-----------|--------|
| S3-J-01 | J | Event `payment.received` / `booking.paid_in_full` → payload untuk logistics + finance | `slice-S3.md` |
| S3-J-02 | J | Kontrak tugas fulfillment minimal (status, assignee) | `slice-S3.md` |
| S3-J-03 | J | Kontrak jurnal minimal (akun placeholder + amount rules) | `slice-S3.md` |
| S3-E-01 | E | Review beban posting vs refund | Komentar di kontrak |

## Checklist implementasi

| Kode | Owner | Pekerjaan | Depends on |
|------|-------|-----------|------------|
| S3-E-02 | E | `logistics-svc` trigger + status | S3-J-02 |
| S3-E-03 | E | `finance-svc` posting dasar | S3-J-03 |
| S3-L-02 | L | UI status “perlengkapan” di portal (read-only OK) | S3-J-02 |

---

# S4 — Growth loop (CRM)

## Engineering freeze **[GATE]**

| Kode | Owner | Pekerjaan | Output |
|------|-------|-----------|--------|
| S4-J-01 | J | Skema lead + snapshot UTM + atribusi (Q019/Q057) | `slice-S4.md` |
| S4-J-02 | J | Event dari booking → CRM (nama event + payload) | `slice-S4.md` |

## Checklist implementasi

| Kode | Owner | Pekerjaan | Depends on |
|------|-------|-----------|------------|
| S4-L-01 | L | Lead list + capture form | S4-J-01 |
| S4-E-02 | E | Endpoint penyimpanan lead + konsumsi event | S4-J-02 |

---

# S5 — Hardening

## Engineering freeze **[GATE]**

| Kode | Owner | Pekerjaan | Output |
|------|-------|-----------|--------|
| S5-J-01 | J | Daftar skenario UAT wajib (dari gate MVP) | Checklist di `slice-S5.md` |
| S5-J-02 | J | Matriks severity bug + SLA fix | Tabel |

## Checklist implementasi

| Kode | Owner | Pekerjaan | Depends on |
|------|-------|-----------|------------|
| S5-L-01 | L | UAT journey B2C/agen | S5-J-01 |
| S5-E-01 | E | UAT payment/finance/logistics | S5-J-01 |

---

# Fase 6 — Depth backlog (setelah core S1–S5 stabil)

Paket teknis berikut memetakan **`BL-*` Fase 6** di `docs/00-overview/06-feature-to-backlog-mapping.md` ke **Slice + Task Code**. Satu kode boleh menampung banyak baris backlog; pecah menjadi subtugas (`S4-E-03a`, dll.) jika ukuran PR membengkak.

## Checklist implementasi — depth per domain

| Kode | Owner | Pekerjaan | Depends on | Domain backlog (ringkas; detail di mapping Fase 6) |
|------|-------|-----------|------------|------------------------------------------------------|
| **S1-E-05** | E | **Katalog & master data dalam** — hotel/pembimbing/transport, varian/addon, import & bulk edit, **seat lintas saluran** | S1-E-02, S1-J-01 | `BL-CAT-005`–`011`, `BL-BOOK-007` (`BL-CAT-012`–`013` → **S4-E-04**) |
| **S1-E-06** | E | **IAM & admin platform** — RBAC granular, staf, keamanan sesi/MFA, log terpusat, API keys, template komunikasi, global config, prosedur backup | S1-E-04, S1-J-01 | `BL-IAM-005`–`017` |
| **S1-L-05** | L | **Situs B2C dalam** — homepage sampai self-booking, guest form, riwayat, info logistik, KB, chat, dll. (**bukan** VA checkout & upload dokumen) | S1-L-02, S1-J-01 | `BL-B2C-001`–`017`, `019`, `021`–`024` |
| **S2-L-04** | L | **Checkout B2C dalam** — UX gateway + wiring ke invoice/VA | S2-L-02, S2-J-01 | `BL-B2C-018` |
| **S3-E-04** | E | **Operasi lapangan dalam** — dokumen kolektif, manifest, rooming/transport, visa UI data, ALL, bus, raudhah, zamzam, refund admin, checklist vendor, … | S3-E-02, S3-J-02 | `BL-OPS-010`, `011`, `020`, `BL-OPS-021`–`042` |
| **S3-E-05** | E | **Gudang & procurement dalam** — PR/PO/GRN/QC, stok multi-gudang, assembly, kirim/retur | S3-E-02, S3-J-02 | `BL-LOG-010`–`029` |
| **S3-E-06** | E | **Visa pipeline dalam** — readiness, bulk submit, poll provider + history | S3-J-02, `BL-VISA-001` | `BL-VISA-001`–`003` |
| **S3-E-07** | E | **Finance modul dalam** — penagihan, bank, subledger, AP ladder, pajak, rev-rec lanjutan, komisi, laporan & audit | S3-E-03, S3-J-03 | `BL-FIN-010`, `011`, `BL-FIN-020`–`041` |
| **S3-L-03** | L | **Self-upload dokumen B2C** + status baca | S3-J-02, `BL-DOC-001` | `BL-B2C-020` |
| **S3-L-04** | L | **Jamaah journey** — in-trip / post-booking (jadwal, dompet dokumen, bus, zamzam, …) | S3-L-02, S3-J-02 | `BL-JMJ-001`–`013` |
| **S3-L-05** | L | **Widget dashboard operasional** — bus live, raudhah, koper, insiden, kesehatan gudang, pantauan logistik | S3-L-02, S3-J-02 | `BL-DASH-008`–`013` |
| **S4-E-03** | E | **CRM & growth API/back-office dalam** — kampanye, otomasi, routing lead, wallet/pencairan, ads/UTM backend, diskon approval, alumni/ZISWAF data | S4-E-02, S4-J-01 | `BL-CRM-012`, `BL-CRM-040`–`066` (+ baris Owner **E** di **6.H**) |
| **S4-E-04** | E | **Sinkron katalog → saluran agen** — snapshot versi, diff, push idempotent | S4-J-02, S1-E-05 | `BL-CAT-012`, `013` |
| **S4-L-02** | L | **Portal agen & aset pemasaran** — onboarding, replica, konten, flyer/itinerary, academy UI, super-view, … | S4-L-01, S4-J-01 | `BL-CRM-010`, `011`, `BL-CRM-013`–`039` (+ baris Owner **L** di **6.H**) |
| **S4-L-03** | L | **Widget dashboard iklan & CS** — spend vs closings, kinerja CS | S4-L-01, S4-J-01 | `BL-DASH-006`, `007` |
| **S5-L-02** | L | **Widget dashboard eksekutif** — vendor readiness, seat, kas/P&L ringkas, likuiditas, inventory/PO exec, kerusakan | S5-L-01, S5-J-01 | `BL-DASH-001`–`005`, `014`–`017` |
| **S5-L-03** | L | **Daily app & pelengkap** — sholat, kiblat, konten manasik/komunitas | S5-J-01 | `BL-PLG-001`–`009` |

---

## Template mapping backlog (feature-level)

Gunakan template ini untuk memetakan detail fitur ke task code:

| Feature ID | Ringkasan fitur | Priority | Slice | Task Code | Backlog ID | Owner | Acceptance ringkas |
|------------|------------------|----------|-------|-----------|------------|-------|--------------------|
| F4-BOOK-001 | Create draft booking | MH-MVP | S1 | S1-E-03 | BL-BOOK-001 | E | Draft tersimpan + seat reserve atomik |
| F5-PAY-001 | Issue VA | MH-MVP | S2 | S2-E-02 | BL-PAY-001 | E | VA issued + TTL tersimpan |
| F10-CRM-001 | Lead capture basic | SH | S4 | S4-L-01 | BL-CRM-001 | L | Lead tersimpan dengan UTM |

Aturan isi:

1. `Priority` wajib diisi dari level feature/backlog.
2. `Slice` dan `Task Code` menunjukkan kapan dan oleh siapa dikerjakan.
3. Satu `Feature ID` boleh pecah ke beberapa `Backlog ID` jika terlalu besar untuk 1 PR.

---

## Masukan singkat (biar sistem kode tugasnya hidup)

1. **Jangan mulai slice tanpa merge kontrak** — itu satu PR kecil yang bisa direview dalam <1 jam.
2. Saat minta AI mengerjakan kode, lampirkan: **kode tugas + isi kontrak slice + service owner** — supaya AI tidak mengubah boundary sembarangan.
3. Kalau kontrak berubah: **bump versi** (`slice-S2-v2.md` atau seksi *Changelog* di file yang sama) supaya histori jelas.

---

## Referensi

- Urutan fase & slice (sequence-first): [`04-delivery-plan-2p-sequence-first.md`](./04-delivery-plan-2p-sequence-first.md)
- Mapping `BL-*` + `Exec seq` (termasuk Fase 6 → kode tugas depth): [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md)
- Model data pembayaran: `docs/03-services/04-payment-svc/02-data-model.md`
- Alur booking (produk): `docs/06-features/04-booking-and-allocation.md`
