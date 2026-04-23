#!/usr/bin/env bash
# run-uat.sh — UAT runner untuk UmrohOS
#
# Usage:
#   bash tests/uat/run-uat.sh            # Jalankan semua 3 agent paralel
#   bash tests/uat/run-uat.sh agent1     # Hanya Agent 1 (Auth + Catalog)
#   bash tests/uat/run-uat.sh agent2     # Hanya Agent 2 (Booking + Payment)
#   bash tests/uat/run-uat.sh agent3     # Hanya Agent 3 (CRM + Ops)
#   bash tests/uat/run-uat.sh cleanup    # Jalankan cleanup saja

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
E2E_DIR="$ROOT_DIR/tests/e2e"
RUN_DATE=$(date +%Y-%m-%d_%H-%M)
REPORT_DIR="$SCRIPT_DIR/runs/$RUN_DATE"

# Warna output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║       UmrohOS UAT Runner — $RUN_DATE       ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"
echo ""

mkdir -p "$REPORT_DIR"

# ─── Cleanup mode ─────────────────────────────────────────────────────────────
if [[ "${1:-}" == "cleanup" ]]; then
  echo -e "${YELLOW}🧹 Menjalankan UAT cleanup...${NC}"
  bash "$SCRIPT_DIR/cleanup/cleanup-uat.sh"
  exit 0
fi

# ─── Cek dependencies ─────────────────────────────────────────────────────────
if ! command -v node &> /dev/null; then
  echo -e "${RED}❌ Node.js tidak ditemukan. Install dulu.${NC}"
  exit 1
fi

if [[ ! -f "$E2E_DIR/.env.prod" ]]; then
  echo -e "${RED}❌ File $E2E_DIR/.env.prod tidak ditemukan.${NC}"
  echo "   Buat dulu dengan copy dari .env.prod.example atau lihat tests/uat/CONFIG.md"
  exit 1
fi

# ─── Install playwright browsers jika belum ada ───────────────────────────────
cd "$E2E_DIR"
if [[ ! -d "node_modules" ]]; then
  echo -e "${YELLOW}📦 Installing dependencies...${NC}"
  npm install
fi

echo -e "${YELLOW}🌐 Target: http://216.176.238.161 (gateway :4000, web :80)${NC}"
echo ""

# ─── Fungsi run per agent ─────────────────────────────────────────────────────
run_agent() {
  local agent_num=$1
  local spec_file=$2
  local scope=$3
  local log_file="$REPORT_DIR/agent-${agent_num}-output.log"

  echo -e "${BLUE}▶ Agent ${agent_num} (${scope}) starting...${NC}"

  # Run API tests
  echo "=== API Tests ===" >> "$log_file"
  set +e
  NODE_ENV=production \
    $(cat "$E2E_DIR/.env.prod" | grep -v '^#' | grep -v '^$' | xargs -I{} echo "export {};" | head -20) \
    npx dotenv -e "$E2E_DIR/.env.prod" -- \
    npx playwright test "tests/${spec_file}" \
    --project=api \
    --reporter=list \
    2>&1 | tee -a "$log_file"

  echo "" >> "$log_file"
  echo "=== Browser Tests ===" >> "$log_file"

  # Run browser tests (chromium harus terinstall)
  npx dotenv -e "$E2E_DIR/.env.prod" -- \
    npx playwright test "tests/${spec_file}" \
    --project=browser \
    --reporter=list \
    2>&1 | tee -a "$log_file"
  set -e

  echo -e "${GREEN}✓ Agent ${agent_num} selesai. Log: $log_file${NC}"
}

# ─── Main run logic ───────────────────────────────────────────────────────────
cd "$E2E_DIR"

case "${1:-all}" in
  agent1)
    run_agent 1 "04-uat-s1-auth-catalog.spec.ts" "Auth+Catalog"
    ;;
  agent2)
    run_agent 2 "05-uat-s2-booking-payment.spec.ts" "Booking+Payment"
    ;;
  agent3)
    run_agent 3 "06-uat-s3s4-crm-ops.spec.ts" "CRM+Ops"
    ;;
  all)
    echo -e "${YELLOW}🚀 Menjalankan 3 agent secara paralel...${NC}"
    echo ""

    # Jalankan paralel dengan background processes
    run_agent 1 "04-uat-s1-auth-catalog.spec.ts" "Auth+Catalog" &
    PID1=$!
    run_agent 2 "05-uat-s2-booking-payment.spec.ts" "Booking+Payment" &
    PID2=$!
    run_agent 3 "06-uat-s3s4-crm-ops.spec.ts" "CRM+Ops" &
    PID3=$!

    # Tunggu semua selesai
    wait $PID1 && echo -e "${GREEN}✓ Agent 1 done${NC}" || echo -e "${RED}✗ Agent 1 ada failures${NC}"
    wait $PID2 && echo -e "${GREEN}✓ Agent 2 done${NC}" || echo -e "${RED}✗ Agent 2 ada failures${NC}"
    wait $PID3 && echo -e "${GREEN}✓ Agent 3 done${NC}" || echo -e "${RED}✗ Agent 3 ada failures${NC}"
    ;;
  *)
    echo "Usage: bash run-uat.sh [all|agent1|agent2|agent3|cleanup]"
    exit 1
    ;;
esac

echo ""
echo -e "${BLUE}══════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Logs tersimpan di: $REPORT_DIR${NC}"
echo -e "${BLUE}  Issues tracker: tests/uat/ISSUES.md${NC}"
echo -e "${BLUE}══════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}💡 Untuk cleanup test data: bash tests/uat/run-uat.sh cleanup${NC}"
