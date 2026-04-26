$ErrorActionPreference = "Stop"

$scriptPath = Join-Path $PSScriptRoot "dev-up.ps1"
$scriptText = Get-Content -Raw -LiteralPath $scriptPath

function Assert-True {
    param(
        [bool]$Condition,
        [string]$Message
    )

    if (-not $Condition) {
        throw $Message
    }
}

Assert-True `
    ($scriptText -match "-EncodedCommand") `
    "dev-up.ps1 should pass child PowerShell commands with -EncodedCommand so quoted window titles survive Start-Process."

Assert-True `
    ($scriptText -match "Wait-BackendReady") `
    "dev-up.ps1 should wait for the backend API before starting Vite."

$outPath = Join-Path $env:TEMP "poker-dev-up-title-test.out.txt"
$errPath = Join-Path $env:TEMP "poker-dev-up-title-test.err.txt"

Remove-Item -LiteralPath $outPath, $errPath -ErrorAction SilentlyContinue

$command = '$Host.UI.RawUI.WindowTitle = "poker backend"; Write-Output "after-title"'
$encodedCommand = [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($command))

$process = Start-Process powershell -ArgumentList @(
    "-NoLogo",
    "-NoProfile",
    "-EncodedCommand",
    $encodedCommand
) -RedirectStandardOutput $outPath -RedirectStandardError $errPath -WindowStyle Hidden -Wait -PassThru

$stdout = if (Test-Path $outPath) { Get-Content -Raw -LiteralPath $outPath } else { "" }
$stderr = if (Test-Path $errPath) { Get-Content -Raw -LiteralPath $errPath } else { "" }

Remove-Item -LiteralPath $outPath, $errPath -ErrorAction SilentlyContinue

Assert-True ($process.ExitCode -eq 0) "Encoded child PowerShell command should exit successfully."
Assert-True ($stdout -match "after-title") "Encoded child PowerShell command should run the command after setting the title."
Assert-True ($stderr -notmatch "The term 'poker' is not recognized") "Encoded child PowerShell command should not treat the quoted title as a command."

Write-Host "dev-up.ps1 tests passed"
