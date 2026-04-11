# logistics-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `logistics.po_approved` | PO approved | po_id, vendor_id, amount | finance-svc (AP entry) |
| `logistics.grn_recorded` | Goods received | grn_id, po_id | finance-svc |
| `logistics.shipment_dispatched` | Shipment created | shipment_id, booking_id | crm-svc (notify customer) |
| `logistics.shipment_delivered` | Courier confirms delivery | shipment_id | crm-svc |
| `logistics.stock_low` | Item below reorder level | sku, warehouse_id, quantity | (alerts) |

## Events consumed

| Event | From | Action |
|---|---|---|
| `booking.paid_in_full` | booking-svc | dispatch kit (via gRPC) |
| `booking.cancelled` | booking-svc | release reserved stock |
