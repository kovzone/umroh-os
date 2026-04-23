CREATE TABLE logistics.shipments (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id          UUID NOT NULL REFERENCES logistics.fulfillment_tasks(id),
    tracking_number  TEXT NOT NULL UNIQUE,
    carrier          TEXT NOT NULL DEFAULT 'manual',
    status           TEXT NOT NULL DEFAULT 'shipped' CHECK (status IN ('shipped','in_transit','delivered','failed')),
    shipped_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    delivered_at     TIMESTAMPTZ,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT shipments_task_id_unique UNIQUE (task_id)
);

CREATE TABLE logistics.pickup_tokens (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id          UUID NOT NULL REFERENCES logistics.fulfillment_tasks(id),
    token            TEXT NOT NULL UNIQUE,
    used             BOOLEAN NOT NULL DEFAULT false,
    used_at          TIMESTAMPTZ,
    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON logistics.shipments (task_id);
CREATE INDEX ON logistics.pickup_tokens (token);
CREATE INDEX ON logistics.pickup_tokens (task_id);
