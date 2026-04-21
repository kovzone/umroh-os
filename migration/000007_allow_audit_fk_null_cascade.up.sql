-- Narrow the append-only trigger on iam.audit_logs so it permits the
-- legitimate FK cascade path (ON DELETE SET NULL → UPDATE user_id/branch_id
-- to NULL) while still rejecting every tamper UPDATE.
--
-- Surfaced by BL-IAM-004: SuspendUser's in-tx audit emit is the first write
-- path into iam.audit_logs, which in turn is the first time any migration
-- rollback (e.g. 000004 down deleting the admin fixture) has a chance to
-- trip the trigger via FK cascade. Previously migration 000003 blocked ALL
-- UPDATEs, so the rollback cascade would fail with insufficient_privilege
-- and leave the DB dirty.
--
-- Compliance posture (F1-AC + UU PDP): audit rows must outlive the users
-- they reference. When a user is deleted (e.g. right-to-erasure), the FK
-- SET NULL is the correct semantic — the audit trail stays, the actor
-- identity is wiped. Tamper-evidence is preserved because resource /
-- resource_id / action / old_value / new_value / ip / created_at remain
-- frozen; only the FK columns may be nulled.

CREATE OR REPLACE FUNCTION iam.audit_logs_reject_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    user_id_nulled   BOOLEAN := NEW.user_id   IS NULL AND OLD.user_id   IS NOT NULL;
    branch_id_nulled BOOLEAN := NEW.branch_id IS NULL AND OLD.branch_id IS NOT NULL;
    user_id_same     BOOLEAN := NEW.user_id   IS NOT DISTINCT FROM OLD.user_id;
    branch_id_same   BOOLEAN := NEW.branch_id IS NOT DISTINCT FROM OLD.branch_id;
BEGIN
    IF TG_OP = 'UPDATE'
       AND NEW.id          =  OLD.id
       AND NEW.resource    =  OLD.resource
       AND NEW.resource_id =  OLD.resource_id
       AND NEW.action      =  OLD.action
       AND NEW.old_value   IS NOT DISTINCT FROM OLD.old_value
       AND NEW.new_value   IS NOT DISTINCT FROM OLD.new_value
       AND NEW.ip          IS NOT DISTINCT FROM OLD.ip
       AND NEW.created_at  =  OLD.created_at
       AND (user_id_nulled   OR user_id_same)
       AND (branch_id_nulled OR branch_id_same)
       AND (user_id_nulled   OR branch_id_nulled)
    THEN
        RETURN NEW;
    END IF;

    RAISE EXCEPTION 'iam.audit_logs is append-only; % is not permitted',
        TG_OP
        USING ERRCODE = 'insufficient_privilege';
END;
$$;
