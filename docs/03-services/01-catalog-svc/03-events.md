# catalog-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `catalog.package.created` | New package | package_id, kind | crm-svc (for marketing) |
| `catalog.package.activated` | Package status → active | package_id | crm-svc |
| `catalog.departure.created` | New departure date | package_departure_id, package_id, date | crm-svc |
| `catalog.seats.reserved` | Saga reservation | package_departure_id, count | (none yet) |
| `catalog.seats.released` | Compensation or cancel | package_departure_id, count | (none yet) |
| `catalog.departure.closed` | Sold out or manual close | package_departure_id | crm-svc |

> Event delivery TBD. For Phase 1, these are notional. The seat reservation/release calls happen synchronously via gRPC from the orchestrating service (booking-svc submit saga; payment-svc refund saga) per ADR 0006.

## Events consumed

| Event | From | Action |
|---|---|---|
| (none yet) | | catalog-svc is upstream of bookings |
