# logistics-svc — Data Model

## Tables (planned)

### `warehouses`
| col | type |
|---|---|
| id | uuid pk |
| name | text |
| address | text |
| branch_id | uuid |

### `stock_items`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| sku | text unique | |
| name | text | |
| unit | text | piece / set / box |
| warehouse_id | uuid fk | |
| quantity | int | |
| reorder_level | int | |
| created_at, updated_at | timestamptz | |

### `purchase_orders`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| code | text unique | |
| vendor_id | uuid | |
| status | po_status enum | |
| total_amount | numeric(15,2) | |
| currency | text | |
| created_by, approved_by | uuid | |
| created_at, updated_at | timestamptz | |

### `po_lines`
| col | type |
|---|---|
| id | uuid pk |
| purchase_order_id | uuid fk |
| stock_item_id | uuid fk |
| quantity | int |
| unit_price | numeric(15,2) |

### `goods_received_notes`
| col | type |
|---|---|
| id | uuid pk |
| purchase_order_id | uuid fk |
| received_by | uuid |
| qc_passed | boolean |
| created_at | timestamptz |

### `kit_definitions`
| col | type |
|---|---|
| id | uuid pk |
| name | text |
| description | text |
| created_at | timestamptz |

### `kit_components`
| col | type |
|---|---|
| kit_definition_id | uuid fk |
| stock_item_id | uuid fk |
| quantity | int |
| pk(kit_definition_id, stock_item_id) |  |

### `shipments`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| booking_id | uuid | |
| recipient_name | text | |
| recipient_address | text | |
| courier | text | |
| tracking_number | text null | |
| status | shipment_status enum | |
| created_at, updated_at | timestamptz | |

## Enums

```sql
CREATE TYPE po_status AS ENUM ('draft', 'submitted', 'approved', 'received', 'cancelled');
CREATE TYPE shipment_status AS ENUM ('pending', 'picked', 'shipped', 'delivered', 'failed', 'returned');
```
