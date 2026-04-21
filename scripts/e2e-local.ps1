#Requires -Version 5.1
<#
.SYNOPSIS
  Bring up the dev Docker stack (if needed), wait for Postgres, run DB migrations, then run Playwright e2e.

.DESCRIPTION
  Windows-friendly alternative to `make dev-bootstrap` + `make e2e-install` + `make e2e` when `make` is not on PATH.
  Requires: Docker Desktop running, Node.js 20+, and `migrate` on PATH (golang-migrate) for schema apply — same as Makefile.

.PARAMETER SkipUp
  If set, skip `docker compose up -d` (stack already running).

.PARAMETER SkipMigrate
  If set, skip `migrate up` (schema already applied).

.EXAMPLE
  .\scripts\e2e-local.ps1
#>
param(
  [switch] $SkipUp,
  [switch] $SkipMigrate
)

$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RepoRoot

Write-Host "== UmrohOS e2e (local) ==" -ForegroundColor Cyan

docker info 2>&1 | Out-Null
if ($LASTEXITCODE -ne 0) {
  Write-Error "Docker is not reachable. Start Docker Desktop and wait until it is ready, then retry."
}

$composeFile = Join-Path $RepoRoot "docker-compose.dev.yml"
if (-not $SkipUp) {
  Write-Host ">> docker compose up -d" -ForegroundColor Yellow
  docker compose -f $composeFile up -d
  if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}

Write-Host ">> Waiting for Postgres (pg_isready)..." -ForegroundColor Yellow
$ready = $false
for ($i = 1; $i -le 30; $i++) {
  docker compose -f $composeFile exec -T postgres pg_isready -U postgres -d umrohos_dev 2>$null | Out-Null
  if ($LASTEXITCODE -eq 0) { $ready = $true; break }
  Start-Sleep -Seconds 2
  Write-Host "   attempt $i/30 ..."
}
if (-not $ready) {
  Write-Error "Postgres did not become ready in time. Check: docker compose -f docker-compose.dev.yml ps"
}

if (-not $SkipMigrate) {
  $migrate = Get-Command migrate -ErrorAction SilentlyContinue
  if (-not $migrate) {
    Write-Warning "golang-migrate (`migrate`) not found on PATH. Install from https://github.com/golang-migrate/migrate — or run migrations from Git Bash: make migrate-up"
  }
  else {
    Write-Host ">> migrate up" -ForegroundColor Yellow
    $dbUrl = "postgres://postgres:changeme@localhost:5432/umrohos_dev?sslmode=disable"
    & migrate -source "file://$RepoRoot/migration" -database $dbUrl up
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
  }
}

$e2eDir = Join-Path $RepoRoot "tests\e2e"
Set-Location $e2eDir
Write-Host ">> npm install (e2e)" -ForegroundColor Yellow
npm install
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host ">> playwright install chromium" -ForegroundColor Yellow
npx playwright install chromium
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host ">> playwright test" -ForegroundColor Yellow
npm test
exit $LASTEXITCODE
