# Glossary — Indonesian / Domain Terms

The PRD uses Indonesian terminology throughout. This glossary maps each term to its English meaning and the canonical code identifier we use in the codebase.

> **Convention:** code identifiers stay in English (snake_case for DB columns, camelCase for Go fields). The Indonesian term is the human label, not the code symbol. Where the Indonesian term is so domain-specific it has no clean English equivalent (e.g. *mahram*, *tasreh*, *muthawwif*), keep the Indonesian word in code too.

## Religious / pilgrimage terms

| Indonesian | English | Code identifier | Notes |
|---|---|---|---|
| Umrah / Umroh | Lesser Islamic pilgrimage to Mecca | `umrah` | Can be done year-round. Smaller than Hajj. |
| Haji / Hajj | Annual Islamic pilgrimage | `hajj` | Time-bound; quota-controlled by Saudi government. |
| Badal | Pilgrimage performed on behalf of another (deceased / incapable) | `badal` | A package variant, not a separate service. |
| Jamaah | Pilgrim | `jamaah` | The core domain entity. Both pre-departure and active. |
| Calon jamaah | Prospective pilgrim | `prospective_jamaah` | Has not yet booked. |
| Alumni | Returned pilgrim | `alumni` | Post-pilgrimage status. |
| Mahram | Male guardian relationship required for women under 45 | `mahram` | Validated via family graph; legal/religious requirement. |
| Muthawwif | Pilgrimage tour guide | `muthawwif` | Saudi-licensed tour leader. |
| Tasreh | Permission letter / entry permit | `tasreh` | E.g. permits for entering Raudhah. |
| Raudhah | The chamber between the Prophet's grave and pulpit in Medina | `raudhah` | Restricted-access area; permits managed digitally. |
| Mecca / Makkah | The holy city | `mecca` | |
| Medina / Madinah | The holy city | `medina` | |
| Nusuk | Saudi Arabia's official pilgrimage app | `nusuk` | Integration target for Raudhah Shield feature. |
| ZISWAF | Zakat, Infaq, Shadaqah, Wakaf — charitable giving categories | `ziswaf` | Alumni hub feature. |
| Zakat | Mandatory Islamic almsgiving | `zakat` | |

## Business / operational terms

| Indonesian | English | Code identifier | Notes |
|---|---|---|---|
| Paket | Package (travel package) | `package` | The sellable product: itinerary + hotel + flight + extras. |
| Agen | Agent (reseller / sub-agent) | `agent` | B2B network member. |
| Cabang | Branch | `branch` | Office location. Data scope boundary. |
| Perwakilan | Representative / super-agent | `representative` | Hierarchical level above plain `agen`. |
| Maskapai | Airline | `airline` | |
| Hotel | Hotel | `hotel` | |
| Vendor | Vendor / supplier | `vendor` | Hotels, airlines, couriers, visa providers. |
| Manifest / Manifes | Passenger manifest | `manifest` | Submitted to immigration/airline. |
| Manasik | Pre-departure orientation / training | `manasik` | Educational session for jamaah. |
| Gudang | Warehouse | `warehouse` | |
| Stok | Stock / inventory | `stock` | |
| DP (Down Payment) | Deposit | `down_payment` | |
| Lunas | Fully paid | `paid_in_full` | Common payment status. |
| Cicilan | Installment | `installment` | |
| Refund | Refund | `refund` | |
| Komisi | Commission | `commission` | Agent earnings. |
| Overriding | Override commission (super-agent's cut from sub-agents) | `override_commission` | |
| Jurnal | Journal entry (accounting) | `journal_entry` | PSAK double-entry. |
| Buku besar | General ledger | `general_ledger` | |
| Neraca | Balance sheet | `balance_sheet` | |
| Laba rugi | Profit and loss | `profit_loss` | |
| Job order | Job-order costing | `job_order` | Each departure date = one cost center. |
| Virtual Account (VA) | Virtual bank account for payment | `virtual_account` | Issued by Midtrans/Xendit. |
| QRIS | Indonesian unified QR payment standard | `qris` | |
| PPh / PPN | Indonesian income tax / VAT | `pph` / `ppn` | PSAK-compliant tax handling. |
| PSAK | Indonesian Financial Accounting Standards | `psak` | Compliance requirement for finance module. |

## Government / regulatory

| Indonesian | English | Notes |
|---|---|---|
| Kemenag | Ministry of Religious Affairs (Indonesia) | Issues PPIU/PIHK travel permits. |
| PPIU | Permit for Umrah travel operators | Required to operate. |
| PIHK | Permit for Hajj travel operators | Required to operate. |
| SISKOPATUH | Kemenag's Umrah computerization system | Data submission target. |
| KTP | Indonesian national ID card | OCR target. |
| Imigrasi | Immigration authority | Manifest submission. |
| MOFA / Sajil | Saudi Ministry of Foreign Affairs visa system | Integration target. |
| GDS | Global Distribution System (Amadeus, Galileo) | Direct ticket issuance. |

## To be expanded

This list is incomplete. Every session that re-reads the PRD should add any new term encountered. Target: 60+ entries by the end of Phase 1.
