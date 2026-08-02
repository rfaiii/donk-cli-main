# Install `donkcli` onto PATH for this user (Windows).
# Copies target\release\donkcli.exe (or debug) into %LOCALAPPDATA%\donk\bin
# and prepends that folder to the user PATH if needed.
#
# Usage:
#   .\scripts\install-donkcli.ps1
#   .\scripts\install-donkcli.ps1 -Debug
# Then open a new terminal and run:  donkcli

param(
    [switch]$Debug
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root
$env:CARGO_TARGET_DIR = Join-Path $Root "target"

$profile = if ($Debug) { "debug" } else { "release" }
Write-Host "==> Building donkcli ($profile)…" -ForegroundColor Cyan
if ($Debug) {
    cargo build -p donk-cli --bin donkcli
} else {
    cargo build -p donk-cli --bin donkcli --release
}
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

$src = Join-Path $Root "target\$profile\donkcli.exe"
if (-not (Test-Path $src)) {
    Write-Error "Missing $src"
}

$binDir = Join-Path $env:LOCALAPPDATA "donk\bin"
New-Item -ItemType Directory -Force -Path $binDir | Out-Null
$dest = Join-Path $binDir "donkcli.exe"
Copy-Item -Force $src $dest
Write-Host "Installed: $dest" -ForegroundColor Green

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not ($userPath -split ';' | Where-Object { $_ -eq $binDir })) {
    $newPath = if ([string]::IsNullOrEmpty($userPath)) { $binDir } else { "$binDir;$userPath" }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = "$binDir;$env:Path"
    Write-Host "Added to user PATH: $binDir" -ForegroundColor Green
    Write-Host "Open a new terminal (or refresh PATH) then run:  donkcli" -ForegroundColor Cyan
} else {
    $env:Path = "$binDir;$env:Path"
    Write-Host "PATH already includes $binDir" -ForegroundColor DarkGray
    Write-Host "Run:  donkcli" -ForegroundColor Cyan
}

& $dest --version
