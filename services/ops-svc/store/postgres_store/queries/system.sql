-- name: ReadyCheck :one
SELECT 1 AS ok;

-- name: InsertDbTxDiagnostic :one
INSERT INTO public.diagnostics (service, message)
VALUES ($1, $2)
RETURNING id, service, message, created_at;

-- name: GetDbTxDiagnostic :one
SELECT id, service, message, created_at
FROM public.diagnostics
WHERE id = $1;
