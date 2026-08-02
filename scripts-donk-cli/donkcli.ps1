# Launch DONK as the `donkcli` command.
# Usage:
#   .\scripts\donkcli.ps1                  # open TUI (skip splash)
#   .\scripts\donkcli.ps1 -Theme donk-dark
#   .\scripts\donkcli.ps1 doctor
#   .\scripts\donkcli.ps1 demo
#   .\scripts\donkcli.ps1 --help
#   .\scripts\donkcli.ps1 -BuildOnly       # build only, print path

param(
    [string]$Theme = "",
    [switch]$BuildOnly,
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$Passthrough
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root
$env:CARGO_TARGET_DIR = Join-Path $Root "target"

Write-Host "==> Building donkcli…" -ForegroundColor Cyan
cargo build -p donk-cli --bin donkcli
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

$donkcli = Join-Path $Root "target\debug\donkcli.exe"
if (-not (Test-Path $donkcli)) {
    Write-Error "Missing binary: $donkcli"
}

if ($BuildOnly) {
    Write-Host "Built: $donkcli"
    Write-Host "Install on PATH: .\scripts\install-donkcli.ps1"
    exit 0
}

$argsList = [System.Collections.Generic.List[string]]::new()
if ($Theme) {
    $argsList.Add("--theme")
    $argsList.Add($Theme)
}
if ($Passthrough -and $Passthrough.Count -gt 0) {
    foreach ($a in $Passthrough) { $argsList.Add($a) }
} else {
    $argsList.Add("--skip-splash")
}

Write-Host "==> donkcli $($argsList -join ' ')" -ForegroundColor Cyan
& $donkcli @argsList
exit $LASTEXITCODE
