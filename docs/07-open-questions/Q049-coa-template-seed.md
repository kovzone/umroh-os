---
id: Q049
title: Chart of Accounts template seed
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: open
---

# Q049 — Chart of Accounts template seed

## Context

The Chart of Accounts (COA) is the GL's spine — every journal_line references an `account_id` from `chart_of_accounts`. F9's journal engine needs the COA **populated on day one** with all the accounts referenced by W2, W4, W5, W6, W7, W10, W11, W12, W13, W14, W15, W16, W17.

Indonesian accounting has several COA conventions:

- **Standar Akuntansi Nasional (SAN)** — widely-used 4-digit numbering (1 Assets, 2 Liabilities, 3 Equity, 4 Revenue, 5 COGS, 6 Expenses, 7 Other revenue, 8 Other expenses, 9 Tax).
- **EFRS/PSAK-aligned custom** — each company designs its own COA to match its reporting needs.
- **Software-specific templates** (MYOB, Accurate, Jurnal.id, Zahir) — tailored templates per industry.

For a travel agency, the COA has industry-specific accounts: Hutang Jamaah, Hutang Tabungan, Pendapatan Umroh, Beban LA (Land Arrangement), Beban Visa, Hutang Komisi Agen, etc. These need to land on day one so the journal engine has targets.

## The question

1. **Seed template** — start from a canonical Indonesian travel-agency COA, or let the agency build their own from scratch?
2. **Customization** — how much can the agency modify the seeded COA? Are "core" accounts (Hutang Jamaah, Pendapatan Umroh) immutable because F9's hardcoded journaling references them?
3. **Numbering depth** — 4-digit (X.X.X.X) fixed, or variable (1, 1.1, 1.1.1, ...)?
4. **Multi-branch COA** — do sub-accounts per branch exist (e.g. 1.1.1.01 Kas Pusat, 1.1.1.02 Kas Cabang Jakarta, 1.1.1.03 Kas Cabang Surabaya)?
5. **Account ↔ journaling rule binding** — the auto-journal engine (W9) maps events to debit/credit account codes; if an agency renames an account, does the rule still work?
6. **Adding accounts** — who can add new COA entries (finance director only, finance admin, or open)?

## Options considered

- **Option A — Bundled canonical travel-agency COA template; customizable but with system-reserved accounts.** Ship a 60–80 account template tailored for Umrah ERP; mark ~20 accounts as `is_system_seeded` (used by journal engine rules); allow adding children + custom accounts but not renaming system-seeded ones.
  - Pros: day-one functionality; opinionated defaults; journal rules don't break.
  - Cons: harder for agencies with existing COA to migrate.
- **Option B — Agency builds own COA + maps auto-journal rules to their account codes.** Empty start; configuration UI for mapping each event-type to an account.
  - Pros: flexible; matches any existing agency COA.
  - Cons: substantial onboarding effort; risk of misconfigured rules.
- **Option C — Ship multiple templates (compact / detailed / regional); agency picks on setup + customizes.**
  - Pros: middle ground; choice without chaos.
  - Cons: more templates to maintain; agencies may want yet another.

## Recommendation

**Option A — bundled canonical travel-agency COA with ~20 system-seeded accounts + agency-extensible.**

Empty start (Option B) pushes weeks of onboarding friction that the MVP can't absorb. Multiple templates (Option C) is a feature without a current use-case — agencies either want the canonical or want their own; the canonical covers 80%+ by itself. Option A lets us ship fast with opinionated defaults, preserves journal-engine reliability via system-seeded flags, and keeps extensibility open.

Defaults to propose: 4-digit hierarchical COA (X.X.X.X), with examples:

```
1 Asset
  1.1 Aktiva Lancar
    1.1.1 Kas & Bank
      1.1.1.01 Kas Kecil Pusat
      1.1.1.02 Bank BCA Pusat (configurable per agency)
    1.1.2 Piutang
      1.1.2.01 Piutang Jamaah
      1.1.2.02 Piutang Karyawan
  1.2 Aktiva Tetap
    1.2.1 Aset Tetap
    1.2.2 Akumulasi Depresiasi
2 Liabilitas
  2.1 Hutang Lancar
    2.1.01 Hutang Jamaah
    2.1.02 Hutang Tabungan Jamaah
    2.1.03 Hutang Usaha
    2.1.04 Hutang Komisi Agen
    2.1.05 PPh 23 Dipotong
    2.1.06 PPN Keluaran
3 Equity
4 Pendapatan
  4.1 Pendapatan Umroh
  4.9 Pendapatan Lain-Lain (incl. pinalti per Q053)
5 HPP / Direct Cost
  5.1 HPP Umroh (Ticket, Visa, LA, Hotel, Bus) — tagged per job_order
6 Beban Operasional
  6.1 Beban Komisi Agen
  6.2 Beban Depresiasi
  6.3 FX Loss
  6.x Other operational
8 Beban Lain-Lain (non-operating)
9 Pajak Penghasilan Badan
```

System-seeded accounts (immutable code, renameable label): Hutang Jamaah, Hutang Tabungan Jamaah, Hutang Usaha, Hutang Komisi Agen, Pendapatan Umroh, Pendapatan Lain-Lain, Bank (class), FX Gain, FX Loss, PPh 23 Dipotong, PPN Keluaran, Akumulasi Depresiasi, Beban Komisi Agen, Beban Depresiasi. Auto-journal rules reference these by their `is_system_seeded` flag + `code` (not label), so label changes don't break rules. Multi-branch: per-branch Kas + Bank sub-accounts at 1.1.1.XX; agency adds on setup. Adding accounts: finance director role by default; configurable.

Reversibility: switching to a different template later requires reclassifying historical journals — not trivial. Customization within the shipped template is cheap.

## Answer

TBD — awaiting stakeholder input. Deciders: finance director, external accountant, agency owner (existing COA if any).
