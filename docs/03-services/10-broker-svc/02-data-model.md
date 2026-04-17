# broker-svc — Data Model

> ⚠️ **DEFERRED FOR MVP — reserved for F6.** See ADR 0006.

broker-svc owns no business data and no Postgres tables. All workflow state is owned by Temporal in its own datastore.

## What Temporal stores

- Workflow execution history (replayable event log)
- Workflow status (running / completed / failed / terminated)
- Pending activities and timers
- Search attributes (queryable workflow metadata)

## Configuration

Temporal connection details live in `broker-svc/config.json` under a `temporal` section:

```json
{
  "temporal": {
    "host_port": "temporal:7233",
    "namespace": "umrohos",
    "task_queue": "umrohos-default"
  }
}
```

The Temporal server config (`temporal/` directory and docker-compose blocks) is deferred from MVP per ADR 0006. It will be reintroduced when broker-svc and the F6 visa pipeline land.

## Search attributes (planned)

For querying workflows by business identifier:
- `BookingId` (keyword)
- `JamaahId` (keyword)
- `BranchId` (keyword)
- `WorkflowKind` (keyword) — booking_saga / visa_pipeline / refund_flow / payment_reconciliation
