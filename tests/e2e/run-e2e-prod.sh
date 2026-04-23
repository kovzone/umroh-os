#!/usr/bin/env bash
# run-e2e-prod.sh — Jalankan SEMUA e2e tests terhadap server 216.176.238.161
#
# Mencakup e2e lama (01-03) + UAT baru (04-06) sekaligus.
#
# Usage (dari root project):
#   bash tests/e2e/run-e2e-prod.sh            # Semua tests
#   bash tests/e2e/run-e2e-prod.sh api        # Hanya API tests (01-02*)
#   bash tests/e2e/run-e2e-prod.sh browser    # Hanya browser tests (03-06)
#   bash tests/e2e/run-e2e-prod.sh uat        # Hanya UAT tests (04-06)
#   bash tests/e2e/run-e2e-prod.sh <pattern>  # Filter nama file, e.g. "auth" atau "04-uat"
#
# Syarat:
#   - tests/e2e/.env.prod harus ada (sudah dibuat)
#   - npx playwright install chromium (kalau belum pernah)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$SCRIPT_DIR/.env.prod"
RUN_DATE=$(date +%Y-%m-%d_%H-%M)
REPORT_DIR="$ROOT_DIR/tests/uat/runs/$RUN_DATE"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${BLUE}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   UmrohOS E2E — Production Run — $RUN_DATE  ║${NC}"
echo -e "${BLUE}╠══════════════════════════════════════════════════╣${NC}"
echo -e "${BLUE}║  Target  : http://216.176.238.161                ║${NC}"
echo -e "${BLUE}║  Gateway : http://216.176.238.161:4000           ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════╝${NC}"
echo ""

# ─── Cek kebutuhan ──────────────────────────────────────────────────────────
if [[ ! -f "$ENV_FILE" ]]; then
  echo -e "${RED}❌ $ENV_FILE tidak ditemukan.${NC}"
  exit 1
fi

if ! command -v node &>/dev/null; then
  echo -e "${RED}❌ Node.js tidak ditemukan.${NC}"
  exit 1
fi

cd "$SCRIPT_DIR"

if [[ ! -d "node_modules" ]]; then
  echo -e "${YELLOW}📦 Installing dependencies...${NC}"
  npm install
fi

# Install chromium jika belum ada (untuk browser tests)
if ! npx playwright install --dry-run chromium &>/dev/null 2>&1; then
  echo -e "${YELLOW}🌐 Installing Playwright Chromium...${NC}"
  npx playwright install chromium
fi

mkdir -p "$REPORT_DIR"

# ─── Tentukan apa yang dijalankan ─────────────────────────────────────────────
MODE="${1:-all}"

run_with_env() {
  local args=("$@")
  # Load .env.prod sebagai env vars, lalu jalankan playwright
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
  npx playwright test "${args[@]}"
}

echo -e "${CYAN}Mode: $MODE${NC}"
echo ""

# ─── Jalankan berdasarkan mode ────────────────────────────────────────────────
case "$MODE" in

  api)
    echo -e "${YELLOW}▶ API tests saja (project=api, spec 01-02*)...${NC}"
    run_with_env --project=api \
      --reporter=list \
      --output="$REPORT_DIR" \
      2>&1 | tee "$REPORT_DIR/api-run.log"
    ;;

  browser)
    echo -e "${YELLOW}▶ Browser tests saja (project=browser, spec 03-06)...${NC}"
    run_with_env --project=browser \
      --reporter=list \
      --output="$REPORT_DIR" \
      2>&1 | tee "$REPORT_DIR/browser-run.log"
    ;;

  uat)
    echo -e "${YELLOW}▶ UAT tests saja (04-06)...${NC}"
    run_with_env \
      "tests/04-uat-s1-auth-catalog.spec.ts" \
      "tests/05-uat-s2-booking-payment.spec.ts" \
      "tests/06-uat-s3s4-crm-ops.spec.ts" \
      --reporter=list \
      --output="$REPORT_DIR" \
      2>&1 | tee "$REPORT_DIR/uat-run.log"
    ;;

  all)
    echo -e "${YELLOW}▶ Semua tests (API + Browser)...${NC}"
    echo ""
    echo -e "${CYAN}[1/2] API tests (01-02*, project=api)${NC}"
    run_with_env --project=api \
      --reporter=list \
      --output="$REPORT_DIR" \
      2>&1 | tee "$REPORT_DIR/api-run.log" || true

    echo ""
    echo -e "${CYAN}[2/2] Browser tests (03-06, project=browser)${NC}"
    run_with_env --project=browser \
      --reporter=list \
      --output="$REPORT_DIR" \
      2>&1 | tee "$REPORT_DIR/browser-run.log" || true
    ;;

  *)
    # Pattern match — e.g. "auth", "catalog", "04-uat"
    echo -e "${YELLOW}▶ Tests matching pattern: *${MODE}*${NC}"
    run_with_env \
      --grep "$MODE" \
      --reporter=list \
      --output="$REPORT_DIR" \
      2>&1 | tee "$REPORT_DIR/filtered-run.log"
    ;;

esac

# ─── Summary ─────────────────────────────────────────────────────────────────
echo ""
echo -e "${BLUE}══════════════════════════════════════════════════${NC}"

# Hitung pass/fail dari log
PASS=$(grep -c "✓\|passed" "$REPORT_DIR"/*.log 2>/dev/null | awk -F: '{sum += $2} END {print sum+0}')
FAIL=$(grep -c "✗\|failed\|×" "$REPORT_DIR"/*.log 2>/dev/null | awk -F: '{sum += $2} END {print sum+0}')

echo -e "  Log dir  : ${CYAN}$REPORT_DIR${NC}"
echo -e "  Issues   : ${CYAN}tests/uat/ISSUES.md${NC}"
echo ""

if [[ "$FAIL" -gt 0 ]]; then
  echo -e "  ${RED}⚠️  Ada test failures. Lihat log di atas, tulis ke ISSUES.md.${NC}"
else
  echo -e "  ${GREEN}✅ Semua test lulus!${NC}"
fi

echo -e "${BLUE}══════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}💡 Cleanup test data: bash tests/uat/cleanup/cleanup-uat.sh${NC}"
