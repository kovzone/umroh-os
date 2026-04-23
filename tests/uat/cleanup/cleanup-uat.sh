#!/usr/bin/env bash
# cleanup-uat.sh
# Menghapus semua test data [UAT] langsung dari database.
# Jalankan dari root project: bash tests/uat/cleanup/cleanup-uat.sh
#
# Syarat: psql harus terinstall di mesin yang menjalankan script ini.
# Alternatif: uncomment blok Docker di bawah jika psql tidak tersedia.

set -euo pipefail

PG_URL="postgres://postgres:IDL4Ssfdo9ettSaFfleZp4M+3vKA8wX2@216.176.238.161:5432/umrohos?sslmode=disable"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_FILE="$SCRIPT_DIR/cleanup-uat.sql"

echo "🧹 UAT Cleanup dimulai..."
echo "   DB: 216.176.238.161:5432/umrohos"
echo "   SQL: $SQL_FILE"
echo ""

# Jalankan cleanup SQL
psql "$PG_URL" -f "$SQL_FILE"

echo ""
echo "✅ UAT cleanup selesai."
