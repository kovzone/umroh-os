-- Reverting 000007 semantically returns the system to "no audit writes have
-- happened yet": only under that precondition is 000003's blanket reject
-- trigger safe alongside subsequent seed rollbacks (000004..000006 delete
-- fixture users, which cascade ON DELETE SET NULL on iam.audit_logs.user_id).
--
-- TRUNCATE bypasses the per-row BEFORE UPDATE / BEFORE DELETE triggers by
-- design, so it runs while the narrowed-or-strict trigger is in place. After
-- TRUNCATE, restoring the strict body leaves the system in the exact state
-- migration 000003 left it (empty append-only table, blanket reject trigger).

TRUNCATE iam.audit_logs;

CREATE OR REPLACE FUNCTION iam.audit_logs_reject_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'iam.audit_logs is append-only; % is not permitted',
        TG_OP
        USING ERRCODE = 'insufficient_privilege';
END;
$$;
