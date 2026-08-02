# Launch DONK inside Alacritty using the bundled host config.
# Usage:
#   .\scripts\run-alacritty-host.ps1
#   .\scripts\run-alacritty-host.ps1 -SkipSplash
#   $env:DONK_BIN = "C:\path\to\donk.exe"; .\scripts\run-alacritty-host.ps1

param(
    [switch]$SkipSplash,
    [string]$Theme = ""
)

$ErrorActionPreference = "Stop"

$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
$ConfigTemplate = Join-Path $Root "resources\host\alacritty\donk.toml"
$RuntimeDir = Join-Path $env:TEMP "donk-alacritty-host"
$RuntimeConfig = Join-Path $RuntimeDir "donk.toml"

if (-not (Test-Path $ConfigTemplate)) {
    Write-Error "Missing host config: $ConfigTemplate"
}

$donkBin = $env:DONK_BIN
if (-not $donkBin) {
    $candidate = Join-Path $Root "target\release\donk.exe"
    if (Test-Path $candidate) {
        $donkBin = $candidate
    } else {
        $candidate = Join-Path $Root "target\debug\donk.exe"
        if (Test-Path $candidate) {
            $donkBin = $candidate
        }
    }
}
if (-not $donkBin) {
    $cmd = Get-Command donk -ErrorAction SilentlyContinue
    if ($cmd) { $donkBin = $cmd.Source }
}
if (-not $donkBin) {
    Write-Host "Building donk (debug)..."
    Push-Location $Root
    cargo build -p donk-terminal --bin donk
    Pop-Location
    $donkBin = Join-Path $Root "target\debug\donk.exe"
}
if (-not (Test-Path $donkBin)) {
    Write-Error "Could not find donk binary. Build with: cargo build -p donk-terminal --bin donk"
}

$alacritty = $env:ALACRITTY_BIN
if (-not $alacritty) {
    $cmd = Get-Command alacritty -ErrorAction SilentlyContinue
    if ($cmd) { $alacritty = $cmd.Source }
}
if (-not $alacritty) {
    Write-Error "Alacritty not found on PATH. Install it or set ALACRITTY_BIN."
}

New-Item -ItemType Directory -Force -Path $RuntimeDir | Out-Null

# Rewrite shell.program to absolute path (TOML-escaped backslashes).
$escaped = ($donkBin -replace '\\', '/')
$argsList = @()
if ($SkipSplash) { $argsList += '"--skip-splash"' }
if ($Theme) { $argsList += "`"--theme`"", "`"$Theme`"" }
$argsToml = if ($argsList.Count -gt 0) { "[" + ($argsList -join ", ") + "]" } else { "[]" }

$content = Get-Content -Raw $ConfigTemplate
$content = $content -replace 'program = "donk"', "program = `"$escaped`""
$content = $content -replace 'args = \[\]', "args = $argsToml"
Set-Content -Path $RuntimeConfig -Value $content -Encoding UTF8

Write-Host "DONK host"
Write-Host "  alacritty: $alacritty"
Write-Host "  donk:      $donkBin"
Write-Host "  config:    $RuntimeConfig"

& $alacritty --config-file $RuntimeConfig
