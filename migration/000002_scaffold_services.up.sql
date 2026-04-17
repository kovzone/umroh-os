-- Scaffolding step: create the shared infrastructure table used by every
-- service's /system/diagnostics/db-tx probe.
--
-- The diagnostics endpoint is a proof-of-connectivity reference — it writes a
-- row inside a WithTx callback and reads it back, confirming the service can
-- reach Postgres through the same transaction helper real queries will use.
-- All services share this single table (under the default `public` schema).
-- Per-service schemas (iam, catalog, booking, ...) are created later when
-- each service's real domain tables land (F1.2, F2.2, ...) per ADR 0007.

CREATE TABLE IF NOT EXISTS public.diagnostics (
    id         BIGSERIAL PRIMARY KEY,
    service    TEXT NOT NULL,
    message    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS diagnostics_service_created_at_idx
    ON public.diagnostics (service, created_at DESC);
