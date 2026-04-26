param(
    [string]$BackendHost = "127.0.0.1",
    [int]$BackendPort = 8080,
    [string]$FrontendHost = "127.0.0.1",
    [int]$FrontendPort = 4173,
    [int]$BackendReadyTimeoutSeconds = 60,
    [switch]$SkipInstall
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$webRoot = Join-Path $repoRoot "web"
$repoPattern = [regex]::Escape($repoRoot)

function Stop-RepoProcesses {
    $targets = Get-CimInstance Win32_Process | Where-Object {
        $_.CommandLine -and
        $_.CommandLine -match $repoPattern -and
        (
            $_.CommandLine -match "cmd[\\/]+pokerd" -or
            $_.CommandLine -match "npm(?:\.cmd)?\s+run\s+dev" -or
            $_.CommandLine -match "vite"
        )
    }

    foreach ($target in ($targets | Sort-Object ProcessId -Unique)) {
        Write-Host "Stopping repo dev process $($target.ProcessId): $($target.Name)"
        Stop-Process -Id $target.ProcessId -Force -ErrorAction SilentlyContinue
    }
}

function Stop-PortListeners {
    param([int[]]$Ports)

    foreach ($port in $Ports) {
        $listeners = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue |
            Select-Object -ExpandProperty OwningProcess -Unique

        foreach ($processId in $listeners) {
            if ($null -eq $processId) {
                continue
            }

            Write-Host "Stopping port listener $processId on :$port"
            Stop-Process -Id $processId -Force -ErrorAction SilentlyContinue
        }
    }
}

function Assert-Command {
    param([string]$Name)

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "missing required command: $Name"
    }
}

function New-EncodedPowerShellCommand {
    param([string]$Command)

    [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($Command))
}

function Wait-BackendReady {
    param(
        [string]$HostName,
        [int]$Port,
        [int]$TimeoutSeconds
    )

    $url = "http://$HostName`:$Port/api/rooms"
    $deadline = (Get-Date).AddSeconds($TimeoutSeconds)

    Write-Host "Waiting for backend API at $url ..."

    do {
        try {
            $response = Invoke-WebRequest -UseBasicParsing -Uri $url -TimeoutSec 2
            if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 500) {
                Write-Host "Backend API is ready."
                return
            }
        }
        catch {
            Start-Sleep -Milliseconds 500
        }
    } while ((Get-Date) -lt $deadline)

    throw "backend API was not ready at $url within $TimeoutSeconds seconds"
}

Assert-Command go
Assert-Command npm

Stop-RepoProcesses
Stop-PortListeners -Ports @($BackendPort, $FrontendPort)
Start-Sleep -Milliseconds 500

if (-not $SkipInstall -and -not (Test-Path (Join-Path $webRoot "node_modules"))) {
    Write-Host "Installing web dependencies..."
    Push-Location $webRoot
    try {
        npm install
    }
    finally {
        Pop-Location
    }
}

$backendCommand = '$Host.UI.RawUI.WindowTitle = "poker backend"; go run ./cmd/pokerd -addr {0}:{1}' -f $BackendHost, $BackendPort
$frontendCommand = '$Host.UI.RawUI.WindowTitle = "poker frontend"; npm run dev -- --host {0} --port {1}' -f $FrontendHost, $FrontendPort

$backendProcess = Start-Process powershell -WorkingDirectory $repoRoot -ArgumentList @(
    "-NoLogo",
    "-NoProfile",
    "-NoExit",
    "-EncodedCommand",
    (New-EncodedPowerShellCommand $backendCommand)
) -PassThru

Wait-BackendReady -HostName $BackendHost -Port $BackendPort -TimeoutSeconds $BackendReadyTimeoutSeconds

$frontendProcess = Start-Process powershell -WorkingDirectory $webRoot -ArgumentList @(
    "-NoLogo",
    "-NoProfile",
    "-NoExit",
    "-EncodedCommand",
    (New-EncodedPowerShellCommand $frontendCommand)
) -PassThru

Write-Host ""
Write-Host "Backend PID:  $($backendProcess.Id)"
Write-Host "Frontend PID: $($frontendProcess.Id)"
Write-Host "Backend URL:  http://$BackendHost`:$BackendPort/"
Write-Host "Frontend URL: http://$FrontendHost`:$FrontendPort/"
Write-Host ""
Write-Host "The script cleared existing repo dev processes and listeners on :$BackendPort / :$FrontendPort before starting."
