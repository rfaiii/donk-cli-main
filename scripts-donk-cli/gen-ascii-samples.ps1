# Regenerate DONK ASCII samples from vendored FIGlet fonts.
# Requires `figlet` on PATH (scoop install figlet / apt install figlet / brew install figlet).

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$FontDir = Join-Path $Root "ref\ascii\fonts"
$OutDir = Join-Path $Root "ref\ascii\samples"
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

if (-not (Get-Command figlet -ErrorAction SilentlyContinue)) {
    Write-Host "figlet not found on PATH."
    Write-Host "Install then re-run, or paste from https://patorjk.com/software/taag/"
    Write-Host "Fonts are ready under ref/ascii/fonts/ — see docs/ASCII-FONTS.md"
    exit 0
}

$Text = if ($args.Count -gt 0) { $args[0] } else { "DONK-CLI" }
Get-ChildItem $FontDir -Filter "*.flf" | ForEach-Object {
    $safe = $_.BaseName -replace '[^\w\-]', '_'
    $out = Join-Path $OutDir "$safe-$Text.txt"
    Write-Host "figlet -f $($_.Name) $Text"
    & figlet -f $_.FullName $Text | Set-Content -Encoding utf8 $out
}
Write-Host "Wrote samples to $OutDir"
