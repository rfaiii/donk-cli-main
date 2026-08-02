# First DONK demo — build, doctor, then launch the TUI.
# Usage:
#   .\scripts\demo.ps1
#   .\scripts\demo.ps1 -DoctorOnly
#   .\scripts\demo.ps1 -Alacritty   # launch inside Alacritty if available

param(
    [switch]$DoctorOnly,
    [switch]$Alacritty,
    [string]$Theme = "donk-dark"
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

# Prefer repo-local target/ (some agent sandboxes redirect CARGO_TARGET_DIR).
$env:CARGO_TARGET_DIR = Join-Path $Root "target"

Write-Host "==> Building donk…" -ForegroundColor Cyan
cargo build -p donk-terminal --bin donk
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

$donk = Join-Path $Root "target\debug\donk.exe"
if (-not (Test-Path $donk)) {
    Write-Error "donk binary missing after build: $donk"
}

Write-Host ""
Write-Host "==> doctor" -ForegroundColor Cyan
& $donk doctor
$doctorCode = $LASTEXITCODE
# AI backend may be missing — still allow UI demo
Write-Host ""

if ($DoctorOnly) {
    exit $doctorCode
}

Write-Host "==> demo guide" -ForegroundColor Cyan
& $donk demo
Write-Host ""

if ($Alacritty) {
    Write-Host "==> Alacritty host" -ForegroundColor Cyan
    & (Join-Path $PSScriptRoot "run-alacritty-host.ps1") -SkipSplash -Theme $Theme
    exit $LASTEXITCODE
}

Write-Host "==> launching TUI (ctrl+c to quit)" -ForegroundColor Cyan
Write-Host "    try: /help  /sys  /files  /read  /animations" -ForegroundColor DarkGray
& $donk --skip-splash --theme $Theme
exit $LASTEXITCODE
