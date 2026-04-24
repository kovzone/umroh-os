#!/usr/bin/env bash
# ============================================================
# UmrohOS UAT Full Run — 2026-04-24
# Jalankan dari: umroh-os/tests/e2e/
# Usage: bash run-uat-full.sh
# Output: tests/uat/runs/2026-04-24/raw-results.txt
# ============================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Load env
if [[ -f .env.prod ]]; then
  set -a; source .env.prod; set +a
  echo "✅ Loaded .env.prod"
else
  echo "❌ .env.prod not found — abort"
  exit 1
fi

# Output dir
OUTDIR="$(dirname "$SCRIPT_DIR")/uat/runs/$(date +%Y-%m-%d)"
mkdir -p "$OUTDIR"
RAW="$OUTDIR/raw-results.txt"

echo "===== UmrohOS UAT Run: $(date) =====" | tee "$RAW"
echo "Server: $GATEWAY_SVC_URL" | tee -a "$RAW"
echo "" | tee -a "$RAW"

# --------------- API Project ---------------
echo "===== PROJECT: api =====" | tee -a "$RAW"
npx playwright test \
  --project=api \
  --reporter=list \
  2>&1 | tee -a "$RAW" || true

echo "" | tee -a "$RAW"

# --------------- UAT-API Project ---------------
echo "===== PROJECT: uat-api =====" | tee -a "$RAW"
npx playwright test \
  --project=uat-api \
  --reporter=list \
  2>&1 | tee -a "$RAW" || true

echo "" | tee -a "$RAW"

# --------------- Browser Project ---------------
echo "===== PROJECT: browser =====" | tee -a "$RAW"
npx playwright test \
  --project=browser \
  --reporter=list \
  2>&1 | tee -a "$RAW" || true

echo "" | tee -a "$RAW"
echo "===== DONE — results saved to: $RAW =====" | tee -a "$RAW"
echo ""
echo "📄 Hasil lengkap ada di: tests/uat/runs/$(date +%Y-%m-%d)/raw-results.txt"
