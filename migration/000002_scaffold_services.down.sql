-- Teardown for 000002_scaffold_services.

DROP INDEX IF EXISTS public.diagnostics_service_created_at_idx;
DROP TABLE IF EXISTS public.diagnostics;
