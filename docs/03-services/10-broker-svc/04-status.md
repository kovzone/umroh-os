# broker-svc — Status

## Current status

**Deferred — reserved for F6 (see ADR 0006).** Not scaffolded in MVP. Temporal is not running in the compose stack. MVP short-lived sagas (booking submit, payment refund) are coordinated in-process by the orchestrating service (booking-svc, payment-svc). broker-svc will be scaffolded + Temporal brought back as part of F6 visa pipeline implementation.

## Implementation checklist (for F6)

The checklist below is reserved for when this service is brought online with F6. No items are in progress.

- [ ] Scaffolded from `baseline/go-backend-template/broker-svc`
- [ ] Wired into `docker-compose.dev.yml` with Temporal dependency
- [ ] Adapters for the services it orchestrates at F6 (visa, jamaah, booking — others added as new workflows land)
- [ ] Visa pipeline workflow + activities (the F6 workflow)
- [ ] Compensation handlers
- [ ] Search attributes registered with Temporal
- [ ] gRPC API for starting workflows (`StartVisaPipeline`)
- [ ] Workflow tests using Temporal's test suite
- [ ] Integration tests with real Temporal
- [ ] Verified by reviewer (walk the visa pipeline in Temporal Web UI)

Note: the earlier checklist listed booking saga, payment reconciliation, and refund flow as broker-svc responsibilities. Those workflows moved to in-process coordination in their respective services per ADR 0006 and are no longer broker-svc scope.
