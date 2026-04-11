# broker-svc — API

## gRPC methods (planned)

`BrokerService` — exposes a thin gRPC surface to start workflows from outside Temporal:

- `StartBookingSaga(StartBookingSagaRequest) → StartBookingSagaResponse` — called by booking-svc on `Submit`
- `StartVisaPipeline(StartVisaPipelineRequest) → StartVisaPipelineResponse` — called by jamaah-svc when documents are ready
- `StartRefundFlow(StartRefundFlowRequest) → StartRefundFlowResponse` — called by booking-svc on `Cancel`
- `SignalPaymentReceived(SignalPaymentReceivedRequest) → SignalPaymentReceivedResponse` — called by payment-svc on webhook

## REST endpoints

None. broker-svc is internal-only and not exposed via the gateway. Workflow visibility happens via Temporal Web UI.

## Activities (internal)

Activities are Go functions registered with Temporal, not exposed externally. They live in `broker-svc/internal/activities/<svc>_activities.go` and wrap gRPC calls to the corresponding service via the adapter pattern.
