# Generates cryptographically random 256-bit secret files in ./.secrets/
# for Docker Compose or local production testing.
$ErrorActionPreference = "Stop"
$secretsDir = Join-Path $PSScriptRoot "..\.secrets"
New-Item -ItemType Directory -Force -Path $secretsDir | Out-Null

$rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()

function New-HexSecret {
    $bytes = New-Object byte[] 32
    $rng.GetBytes($bytes)
    return ($bytes | ForEach-Object { $_.ToString("x2") }) -join ""
}

Set-Content -Path (Join-Path $secretsDir "jwt_secret.txt") -Value (New-HexSecret) -NoNewline
Set-Content -Path (Join-Path $secretsDir "key_encryption_secret.txt") -Value (New-HexSecret) -NoNewline
Write-Host "Generated CSPRNG secret files in $secretsDir" -ForegroundColor Green
