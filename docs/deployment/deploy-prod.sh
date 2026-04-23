#!/usr/bin/env bash
# Run on the production host from the repo root (e.g. /home/infra/umrohos).
# Usage: ./docs/deployment/deploy-prod.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT}"

if [[ ! -f .env.prod ]]; then
  echo "error: .env.prod not found in ${ROOT} — copy from env.prod.sample and fill secrets." >&2
  exit 1
fi

# Load .env.prod for the wait loop (pg_isready user/db).
set -a
# shellcheck source=/dev/null
source .env.prod
set +a

COMPOSE=(docker compose -f docker-compose.prod.yml --env-file .env.prod)

echo "→ building images (if needed)…"
"${COMPOSE[@]}" build

echo "→ starting postgres + otel…"
"${COMPOSE[@]}" up -d postgres otel-collector

echo "→ waiting for postgres…"
for _ in {1..30}; do
  if "${COMPOSE[@]}" exec -T postgres pg_isready -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-umrohos}" 2>/dev/null; then
    break
  fi
  sleep 2
done

echo "→ running DB migrations…"
"${COMPOSE[@]}" --profile tooling run --rm migrate

echo "→ starting application stack…"
"${COMPOSE[@]}" up -d

echo "→ done."
"${COMPOSE[@]}" ps
