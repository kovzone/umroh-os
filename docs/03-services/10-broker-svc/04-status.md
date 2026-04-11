# broker-svc — Status

## Implementation checklist

- [ ] Scaffolded from `baseline/go-backend-template/broker-svc`
- [ ] Wired into `docker-compose.dev.yml` with Temporal dependency
- [ ] Adapters for every other service (iam, catalog, booking, jamaah, payment, visa, ops, logistics, finance, crm)
- [ ] Booking saga workflow + activities
- [ ] Payment reconciliation workflow + activities
- [ ] Visa pipeline workflow + activities
- [ ] Refund flow workflow + activities
- [ ] Compensation handlers for every workflow
- [ ] Search attributes registered with Temporal
- [ ] gRPC API for starting workflows
- [ ] Workflow tests using Temporal's test suite
- [ ] Integration tests with real Temporal
- [ ] Verified by reviewer (walk a saga in Temporal Web UI)

## Current status

**Not started.** broker-svc should be the second or third service scaffolded, after iam-svc and one or two of the data-owning services it needs to orchestrate (e.g. catalog + booking).
