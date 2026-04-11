# payment-svc — Status

## Implementation checklist

- [ ] Scaffolded from baseline template
- [ ] Wired into `docker-compose.dev.yml`
- [ ] Initial DDL written (`_init/payment_db/`)
- [ ] sqlc queries
- [ ] OpenAPI spec
- [ ] Invoice CRUD
- [ ] Midtrans adapter (issue VA, verify webhook, refund)
- [ ] Xendit adapter (same surface)
- [ ] Webhook handlers with signature verification
- [ ] Idempotent payment_event ingestion
- [ ] Reconciliation cron job
- [ ] gRPC methods for saga
- [ ] Unit tests
- [ ] Integration tests with mocked gateway
- [ ] Verified by reviewer

## Current status

**Not started.**
