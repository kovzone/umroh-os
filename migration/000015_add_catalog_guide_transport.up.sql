-- 000015 — Catalog ground transport tables.
--
-- Adds catalog.transports (bus/shuttle/private_car operators) and the
-- catalog.package_transports join table so a package can reference one or
-- more transport operators.
--
-- ID convention: TEXT ULID with 'trn_' prefix, consistent with the
-- existing catalog master ID conventions (htl_, arl_, mtw_, addon_, pkgpr_).

CREATE TABLE catalog.transports (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    type        TEXT NOT NULL,   -- 'bus', 'private_car', 'shuttle'
    capacity    INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT transports_id_prefix CHECK (id LIKE 'trn_%'),
    CONSTRAINT transports_type_valid CHECK (type IN ('bus', 'private_car', 'shuttle')),
    CONSTRAINT transports_capacity_nonneg CHECK (capacity >= 0)
);

CREATE INDEX transports_type_idx ON catalog.transports (type);

CREATE TABLE catalog.package_transports (
    package_id    TEXT NOT NULL REFERENCES catalog.packages(id) ON DELETE CASCADE,
    transport_id  TEXT NOT NULL REFERENCES catalog.transports(id) ON DELETE RESTRICT,
    PRIMARY KEY (package_id, transport_id)
);

CREATE INDEX package_transports_transport_id_idx ON catalog.package_transports (transport_id);
