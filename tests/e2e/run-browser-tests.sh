#!/usr/bin/env bash
# ============================================================
# run-browser-tests.sh — Jalankan browser E2E tests (B2C + Console)
#
# Tests ini menjalankan Playwright terhadap browser Chrome nyata
# pada URL http://216.176.238.161 (production server).
#
# Prerequisites:
#   1. cd ke folder tests/e2e
#   2. npm install (jika belum)
#   3. npx playwright install chromium (jika belum)
#   4. Copy .env.prod ke .env (atau set env var manual)
#
# Cara jalankan:
#   bash tests/e2e/run-browser-tests.sh
#
# Atau dari root repo:
#   bash tests/e2e/run-browser-tests.sh
#
# Options (env var):
#   CORE_WEB_URL   — URL frontend (default: http://216.176.238.161)
#   GATEWAY_SVC_URL — URL gateway API (default: http://216.176.238.161)
#   TEST_FILTER    — Filter test by name, e.g. TEST_FILTER="B2C-01"
#   HEADED         — Set ke "1" untuk jalankan dengan browser visible
# ============================================================

set -e

# Pindah ke direktori tests/e2e (script bisa dijalankan dari mana saja)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
E2E_DIR="$SCRIPT_DIR"

# Jika dipanggil dari root repo
if [[ ! -f "$E2E_DIR/package.json" ]]; then
  E2E_DIR="$SCRIPT_DIR/tests/e2e"
fi

cd "$E2E_DIR"

echo ""
echo "=================================================="
echo "  UmrohOS Browser E2E Tests"
echo "  Target: ${CORE_WEB_URL:-http://216.176.238.161}"
echo "=================================================="
echo ""

# Load .env.prod jika ada dan .env belum ada
if [[ ! -f ".env" && -f ".env.prod" ]]; then
  echo "📋 Menggunakan .env.prod sebagai environment..."
  cp .env.prod .env
elif [[ -f ".env" ]]; then
  echo "📋 Menggunakan .env yang sudah ada..."
else
  echo "⚠️  Tidak ada .env — menggunakan default values dari playwright.config.ts"
fi

# Set default env vars jika belum di-set
export CORE_WEB_URL="${CORE_WEB_URL:-http://216.176.238.161}"
export GATEWAY_SVC_URL="${GATEWAY_SVC_URL:-http://216.176.238.161}"

# Build playwright command
PLAYWRIGHT_CMD="npx playwright test 03b-browser-b2c --project=browser"

# Filter test berdasarkan nama (opsional)
if [[ -n "$TEST_FILTER" ]]; then
  PLAYWRIGHT_CMD="$PLAYWRIGHT_CMD --grep \"$TEST_FILTER\""
  echo "🔍 Filter: $TEST_FILTER"
fi

# Headed mode (browser visible)
if [[ "$HEADED" == "1" ]]; then
  PLAYWRIGHT_CMD="$PLAYWRIGHT_CMD --headed"
  echo "👁️  Headed mode: browser akan terlihat"
fi

echo ""
echo "🚀 Menjalankan: $PLAYWRIGHT_CMD"
echo ""
echo "--------------------------------------------------"
echo ""

# Jalankan tests
eval "$PLAYWRIGHT_CMD"

EXIT_CODE=$?

echo ""
echo "--------------------------------------------------"
echo ""

if [[ $EXIT_CODE -eq 0 ]]; then
  echo "✅ Semua browser tests PASSED!"
else
  echo "❌ Ada test yang FAILED (exit code: $EXIT_CODE)"
  echo ""
  echo "💡 Tips debug:"
  echo "   - Jalankan ulang dengan HEADED=1 untuk lihat browser:"
  echo "     HEADED=1 bash $0"
  echo "   - Jalankan 1 test group saja, contoh:"
  echo "     TEST_FILTER='B2C-01' bash $0"
  echo "   - Lihat HTML report:"
  echo "     npx playwright show-report"
fi

echo ""
echo "📊 Lihat detail report:"
echo "   npx playwright show-report"
echo ""

exit $EXIT_CODE
