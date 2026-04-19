# Documentation language (`docs/`)

## Rule

All **committed prose** under `docs/` MUST be written in **English** (US or UK; stay consistent within a file).

## Scope

- Applies to: Markdown and other text specs in `docs/` (architecture, conventions, feature specs, open-question records, service folders, overview, ADRs).
- **Exception — cited source material:** verbatim quotes, legal text, or product copy taken from stakeholder documents may stay in another language if clearly marked as a **quote** or **source excerpt**. Wrap in a blockquote or label with `Source (Indonesian):` when needed.
- **Exception — MoSCoW CSV:** `docs/Modul UmrohOS - MosCoW.csv` may retain Indonesian module titles as the inventory source; English summaries belong in mapping/spec tables when those rows are edited.

## PRD

The main PRD export (`docs/UmrohOS - Product Requirements Document.docx.md`) may remain largely Indonesian until a deliberate English rewrite. Feature specs in `docs/06-features/` are the **operating** English layer for implementation decisions (per `AGENTS.md` authority order).

## New and updated docs

When you add or change documentation:

1. Write new content in English.
2. If you touch an older file that is still non-English, translate at least the sections you changed (prefer translating the whole file when practical).

## Commits

Do **not** commit a documentation change set until **every file in that commit** that lives under `docs/` is fully English for the parts you are changing (no mixed-language commits for the same migration). Prefer one coherent PR per doc batch.

## Related

- Team workflow: `CONTRIBUTING.md`
- Agent onboarding: `AGENTS.md`
