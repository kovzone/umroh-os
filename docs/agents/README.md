# UmrohOS — Agent System

Dokumen ini menjelaskan sistem multi-agent yang digunakan untuk mengerjakan project UmrohOS.

## Cara Kerja

Lutfi (owner) berbicara dengan **Claude Utama** (orchestrator). Claude Utama membaca konteks,
mendelegasikan pekerjaan ke agent yang tepat berdasarkan task code, lalu melaporkan hasilnya.

```
Lutfi  →  Claude Utama  →  [Agent yang tepat]  →  Repository
                       ←  hasil / laporan      ←
```

## Cara Mengaktifkan Agent

Cukup sebut task code-nya ke Claude Utama:

```
"Kerjakan S1-E-03"           → Elda Agent
"Kerjakan S1-L-02"           → Lutfi Agent
"Cek apakah S1-E-03 sudah sesuai acceptance criteria"  → QA Agent
"Provision server production" → DevSecOps Agent
"Klarifikasi scope S2-J-01"  → PM Agent
```

Atau minta secara natural:
```
"Lanjutkan semua S1 backend tasks yang masih todo"
"Deploy ke server setelah S1 selesai"
"Pastikan semua endpoint sudah aman sebelum merge"
```

## Agent Index

| File | Agent | Domain | Task Codes |
|------|-------|--------|------------|
| [01-pm-agent.md](./01-pm-agent.md) | PM Agent | Product, Contracts, Open Questions | `S*-J-*` (contract docs) |
| [02-lutfi-agent.md](./02-lutfi-agent.md) | Lutfi Agent (L) | Svelte 5, B2C, CRM frontend | `S*-L-*` |
| [03-elda-agent.md](./03-elda-agent.md) | Elda Agent (E) | Go, gRPC, PostgreSQL, backend services | `S*-E-*` |
| [04-qa-security-agent.md](./04-qa-security-agent.md) | QA + Security Agent | Playwright, k6, OWASP | Per-slice acceptance |
| [05-devsecops-agent.md](./05-devsecops-agent.md) | DevSecOps Agent | Docker, Nginx, Ubuntu, CI/CD | Infrastructure tasks |

## Aturan Umum (Berlaku untuk Semua Agent)

1. **Repository adalah source of truth** — selalu baca `docs/contracts/`, feature specs, dan ADR sebelum mulai coding.
2. **Jangan ubah contract** (`docs/contracts/slice-Sx.md`) tanpa persetujuan eksplisit PM Agent + review dari kedua developer (L dan E).
3. **Satu task = satu PR** — jaga PR size kecil, sesuai `CONTRIBUTING.md`.
4. **Definition of Done** harus terpenuhi sebelum task dianggap selesai (lihat `docs/04-delivery-plan-2p-sequence-first.md`).
5. **Joint tasks (J)** — PM Agent menulis contract draft, Lutfi Agent dan Elda Agent masing-masing review dari perspektif domain mereka.

## Kapan Agent Boleh Bertanya ke Lutfi

Hanya untuk:
- Kredensial eksternal (payment gateway production, MOFA/Sajil, dll)
- Keputusan bisnis yang belum ada di docs dan tidak bisa diasumsikan (open questions yang benar-benar butuh owner decision)
- Trade-off arsitektur baru yang belum ada ADR-nya
- Konfirmasi sebelum tindakan **destruktif** (hapus data, reset production, dll)
