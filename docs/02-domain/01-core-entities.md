# Core Entities

The principal domain entities, where they live, and their key relationships. Detailed schemas live in each service's `02-data-model.md`.

## Identity & Access (iam-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `user` | A login account (staff, agent, jamaah) | belongs to `branch`; has many `user_role` |
| `role` | A named bundle of permissions (e.g. "ops_admin") | has many `permission` via `role_permission` |
| `permission` | A single capability (resource + action) | many-to-many with `role` |
| `branch` | A company branch / office | has many `user`, scopes data visibility |
| `audit_log` | Immutable record of CRUD actions | references `user`, `branch` |

## Catalog (catalog-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `package` | A sellable Umrah/Hajj/Badal package | has many `package_departure`, references `hotel`, `airline`, `muthawwif` |
| `package_departure` | Specific departure date with seat inventory | belongs to `package` |
| `hotel` | Master hotel record | many `package` reference it |
| `airline` | Master airline record | many `package` reference it |
| `muthawwif` | Master tour-leader record | many `package` reference it |
| `itinerary_template` | Reusable itinerary structure | linked to `package` |

## Pilgrim (jamaah-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `jamaah` | A pilgrim (calon, active, or alumni) | belongs to a `family_unit` |
| `family_unit` | A group of related jamaah | has many `jamaah`, used for mahram resolution |
| `mahram_relation` | Edge in the family graph | links two `jamaah` with a relationship type |
| `document` | A scanned document (KTP, passport, vaccine) | belongs to `jamaah`, has one `ocr_result` |
| `ocr_result` | Structured data extracted from a document | belongs to `document` |

## Booking (booking-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `booking` | A reservation for one or more jamaah on one departure | references `package_departure` (catalog) |
| `booking_item` | One jamaah on one booking | references `jamaah` (jamaah-svc) |
| `room_allocation` | Which jamaah share which hotel room | belongs to `booking` |
| `bus_allocation` | Which jamaah sit in which bus | belongs to `booking` |

## Payment (payment-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `invoice` | A bill against a booking | references `booking_id` |
| `virtual_account` | A VA issued by Midtrans/Xendit | belongs to `invoice` |
| `payment_event` | A gateway webhook record | references `invoice` |
| `refund` | A refund record | references `invoice` |

## Visa (visa-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `visa_application` | A visa application for one jamaah | references `jamaah_id`, `booking_id` |
| `visa_status_history` | Status transition log | belongs to `visa_application` |
| `e_visa` | Issued visa document | belongs to `visa_application` |
| `tasreh` | Permission letter (e.g. Raudhah) | belongs to `visa_application` |

## Operations (ops-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `verification_task` | A document waiting for verification | references `document_id` |
| `manifest` | Generated manifest for a departure | references `package_departure_id` |
| `luggage_tag` | A scannable tag for a piece of luggage | references `jamaah_id`, `booking_id` |
| `handling_event` | Airport scan / event | references `jamaah_id` |

## Logistics (logistics-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `stock_item` | A SKU in a warehouse | belongs to `warehouse` |
| `warehouse` | A physical warehouse | has many `stock_item` |
| `purchase_order` | A PO to a vendor | references `vendor_id` |
| `grn` | Goods received note | belongs to `purchase_order` |
| `kit_definition` | Bill of materials for a pilgrim kit | has many `stock_item` references |
| `shipment` | A shipped package to a jamaah | references `booking_id` |

## Finance (finance-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `chart_of_accounts` | COA tree | self-referencing parent/child |
| `journal_entry` | Header of a double-entry record | has many `journal_line` |
| `journal_line` | Debit or credit row | belongs to `journal_entry`, references `chart_of_accounts` |
| `ar_balance` | Receivable ledger | references `booking_id` |
| `ap_balance` | Payable ledger | references `vendor_id` |
| `tax_record` | PPh / PPN record | references `journal_entry` |
| `fx_rate` | Exchange rate snapshot | dated record |
| `job_order_cost` | Cost center per departure | references `package_departure_id` |

## CRM (crm-svc)

| Entity | Description | Key relationships |
|---|---|---|
| `lead` | A prospective customer | optionally references `agent_id`, `campaign_id` |
| `campaign` | A marketing campaign | has many `lead` |
| `agent` | An agent (reseller) account | references `user_id`, has parent `agent_id` for hierarchy |
| `commission_ledger_entry` | A commission line | references `agent_id`, `booking_id` |
| `broadcast` | A WhatsApp/email blast | targets a segment |
| `alumni_thread` | Community discussion thread | belongs to `alumni` user |

## Cross-context entity references

These IDs cross service boundaries via gRPC, never via foreign key. The owning service is the source of truth.

| ID | Owned by | Used in |
|---|---|---|
| `user_id` | iam-svc | every service (audit, attribution) |
| `branch_id` | iam-svc | every service that scopes data |
| `package_id`, `package_departure_id` | catalog-svc | booking, ops, finance |
| `jamaah_id` | jamaah-svc | booking, visa, ops, logistics |
| `booking_id` | booking-svc | payment, visa, logistics, ops, finance, crm |
| `invoice_id`, `payment_event_id` | payment-svc | finance |
| `agent_id` | crm-svc | booking (attribution), iam (login) |
