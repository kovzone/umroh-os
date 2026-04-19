# Delivery Plan 2 Person (Sequence-First)

## Context

Dokumen ini memecah pekerjaan **2 orang** dengan prinsip **sequence-first**:

- **Tanpa patokan kalender** sebagai komitmen delivery (durasi fleksibel).
- **Urutan fase** di bawah ini adalah komitmen engineering: apa dulu, apa belakangan, gate mana yang harus lolos sebelum lanjut.

Penamaan file ini **tidak** menyiratkan durasi minggu tetap (bukan “8w”). Tautan lama ke `04-delivery-plan-2p-8w.md` diganti ke nama file ini.

## Dokumen overview terkait (04 / 05 / 06)

| File | Fungsi |
|------|--------|
| **Ini (04)** | **Urutan fase & slice** — konteks 2 orang, prioritas MH-* / SH / CH, ownership, *tanpa* komitmen kalender. |
| [`05-slice-engineering-checklist-and-task-codes.md`](./05-slice-engineering-checklist-and-task-codes.md) | **Gate + kode tugas** — checklist joint/L/E sebelum coding, format `Sx-J-01`, **§ Fase 6 — Depth backlog** (`S1-E-05`, …). |
| [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md) | **Backlog operasional** — baris `BL-*`, `Exec seq`, kolom Slice/Task Code yang mengarah ke **05**; isi produk tetap mengacu `docs/06-features/*`. |

Ringkas: **04** = peta *apa dulu*; **05** = SOP & nomor tugas per slice; **06** = daftar tiket integrasi.

## Label Prioritas (Wajib Pakai)

Semua backlog dan detail fitur di repo ini memakai label:

- **MH-MVP**: wajib selesai untuk fase fungsional paling awal (biasanya inti alur bisnis)
- **MH-V1**: wajib selesai untuk kelengkapan v1 (setelah inti stabil)
- **SH**: penting, dikerjakan setelah seluruh `MH-MVP` + `MH-V1` aman
- **CH**: opsional, dikerjakan jika kapasitas tersedia

Aturan urutan eksekusi: **MH-MVP -> MH-V1 -> SH -> CH**.

## Cara Kerja Shared

"Shared" bukan berarti coding bersamaan di task yang sama.

Setiap task wajib punya:

- **A (Accountable):** penentu final
- **R (Responsible):** eksekutor utama
- **C (Consulted):** reviewer wajib

Aturan:

1. Satu task hanya 1 orang **R**
2. Reviewer tidak mengubah scope tanpa sinkron 15 menit harian
3. Jika kontrak API berubah, update dulu dokumen service API sebelum lanjut coding

## Ownership Utama

### Lutfi (Marketing owner)

- F10 Marketing/CRM/Agent (utama)
- F4 Booking flow side (channel, attribution, UX flow)
- B2C/B2B funnel UI
- F11 dashboard sales/marketing (fase 2)

### Elda (Finance/Inventory owner)

- F5 Payment core (invoice/VA/webhook/reconcile)
- F9 Finance core (journal, AR/AP, report dasar)
- F8 Logistics/Fulfillment core
- F11 dashboard finance/inventory (fase 2)

### Shared Foundation (Lutfi + Elda)

- F1 auth/role/audit minimum
- F2 package/departure/seat minimum
- Kontrak event lintas service (booking/payment/finance/logistics/crm)

## Vertical Slices (Urutan End-to-End)

Urutan terbaik adalah slice yang cepat menghasilkan nilai bisnis dan menurunkan risiko integrasi.

### Slice 1 - Discover & Book Draft

Alur: login -> lihat paket -> pilih departure -> buat draft booking.

### Slice 2 - Get Paid

Alur: draft booking -> issue VA -> webhook masuk -> status paid/partial.

### Slice 3 - Fulfillment + Accounting Minimum

Alur: paid -> trigger fulfillment -> posting jurnal dasar -> status terlihat.

### Slice 4 - Growth Loop

Alur: attribution -> lead tracking -> komisi basic -> dashboard sales awal.

### Slice 5 - Hardening & Go-Live Readiness

Alur: reliability, audit, permission, UAT, bugfix prioritas.

**Checklist dependency engineering per slice + kode pekerjaan (Lutfi/Elda/Joint):** [`05-slice-engineering-checklist-and-task-codes.md`](./05-slice-engineering-checklist-and-task-codes.md).
**Mapping detail fitur -> backlog -> task code:** [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md).

## Urutan Eksekusi (Sequence-First)

### Fase 0 — Engineering bootstrap (S0)

**Tujuan:** kontrak, konvensi merge, definisi DoR/DoD, template `docs/contracts/*`.

**Joint gate:** selesaikan checklist `S0-*` di `05-slice-engineering-checklist-and-task-codes.md`.

### Fase 1 — Discover + draft booking (S1)

**Tujuan:** B2C bisa browse katalog → detail → form → **draft booking**; internal auth minimum untuk route yang dipakai staff bila diperlukan.

**Joint gate:** selesaikan engineering freeze `S1-J-*` sebelum implementasi besar dimulai.

### Fase 2 — Get paid (S2)

**Tujuan:** draft → invoice/VA → webhook → status pembayaran sinkron ke booking.

**Joint gate:** selesaikan engineering freeze `S2-J-*` + uji end-to-end `S2-J-05`.

### Fase 3 — Fulfillment + finance minimum pasca bayar (S3)

**Tujuan:** `paid_in_full` memicu tugas fulfillment/logistik minimal + posting keuangan minimum yang bisa diaudit.

**Joint gate:** selesaikan engineering freeze `S3-J-*`.

### Fase 4 — Growth loop dasar (S4)

**Tujuan:** lead + attribution dasar + event minimal ke CRM read-model.

**Joint gate:** selesaikan engineering freeze `S4-J-*`.

### Fase 5 — Hardening & readiness (S5)

**Tujuan:** UAT, reliability, permission regression, dokumentasi operasional minimum.

**Joint gate:** selesaikan engineering freeze `S5-J-*`.

### Fase 6 — Depth expansion (setelah core stabil)

Bagian ini **sengaja tidak dipatok waktu**. Daftar berikut adalah **prioritas domain produk** (bukan salinan urutan baris di mapping). Urutan konkret antar-`BL-*` dan subsection **6.G–6.O** memakai kolom **`Exec seq`** di [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md).

Urutan domain yang disarankan:

1. **F10** dalam (onboarding agent, commission view, reporting dasar, SLA lead)
2. **F9** dalam (AP layak operasional, aging AR/AP, period close checklist minimum)
3. **F8** dalam (warehouse/QC/reorder critical flow, dispatch tracking minimum)
4. **F7** dalam (manifest/grouping yang paling kritis untuk operasi)
5. **F6** starter (visa tracker/read model dasar) setelah data dokumen + booking stabil
6. **F11** dashboard sesuai kebutuhan operasional (setelah metric backend tersedia)

#### Joint Final Gate (untuk “v1 layak operasi”)

- Alur inti Fase 1–3 tetap stabil setelah fitur depth masuk
- Dokumentasi operasional minimum tersedia
- Bug triage tersisa hanya P2/P3 yang sudah diprioritaskan

## RACI Ringkas untuk Shared Foundation

### F1 minimal (login + role dasar)

- **A:** Lutfi
- **R:** Elda
- **C:** Lutfi

### F2 minimal (paket + departure + seat basic)

- **A:** Lutfi
- **R:** Lutfi
- **C:** Elda

### F4 minimal (create draft booking)

- **A:** Lutfi
- **R:** Lutfi
- **C:** Elda

## Definisi Done per Task

Satu task dianggap selesai jika memenuhi semua:

1. Endpoint/UI jalan sesuai acceptance
2. Event/status antar service sinkron
3. Audit log tercatat untuk aksi state-changing
4. Permission role tidak bocor
5. Minimal 1 test skenario happy path + 1 edge case lolos

## Ritual sinkron (disarankan)

### Sinkron singkat harian

- Kemarin selesai apa
- Hari ini kerjakan apa
- Blocker apa
- Kontrak API/event apa yang berubah

### Sinkron non-harian

- Demo fase/slice yang sedang dikerjakan
- Triage bug (P0/P1/P2)
- Freeze scope fase berikutnya (scope guardrail per fase, bukan per kalender)

## Scope Guardrails

Untuk menjaga **kualitas integrasi** tanpa patokan tanggal:

- Jangan ambil fitur `CH` sebelum seluruh `MH-MVP` pada alur inti (Fase 1–3) stabil
- Visa pipeline lengkap (F6 mendalam), advanced ops lapangan, dan fitur alumni/daily lengkap tetap **setelah** core + depth dasar stabil (lihat Fase 6)
- Jika ada scope baru, wajib trade-off: keluarkan item lain dari fase berjalan atau turunkan prioritas (`SH/CH`)
