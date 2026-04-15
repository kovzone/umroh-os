# logistics-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/warehouses` | List warehouses |
| `POST` | `/v1/warehouses` | Create warehouse |
| `GET` | `/v1/stock-items` | List stock items (filterable by warehouse) |
| `POST` | `/v1/stock-items` | Add SKU |
| `PATCH` | `/v1/stock-items/{id}` | Adjust stock |
| `GET` | `/v1/purchase-orders` | List POs |
| `POST` | `/v1/purchase-orders` | Create PR |
| `POST` | `/v1/purchase-orders/{id}/approve` | Approve PR |
| `POST` | `/v1/purchase-orders/{id}/grn` | Record goods received |
| `GET` | `/v1/kits` | List kit definitions |
| `POST` | `/v1/kits` | Create kit definition |
| `GET` | `/v1/shipments` | List shipments |
| `POST` | `/v1/shipments` | Create shipment |
| `POST` | `/v1/shipments/{id}/label` | Generate courier label |

## gRPC methods (planned)

`LogisticsService`:
- `DispatchKit(...)` — called by payment-svc directly when booking is paid (per ADR 0006)
- `CheckStock(...)`
- `ReleaseStock(...)` — compensating action

> All endpoints are stubs.
