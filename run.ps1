# Windows PowerShell Run Script for WaitWatchersV2
# Builds the project first, then starts the local server.

# Run build
powershell -ExecutionPolicy Bypass -File .\build.ps1
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Build failed. Aborting execution." -ForegroundColor Red
    exit 1
}

Write-Host "[RUN] Launching local server..." -ForegroundColor Cyan
Write-Host "-> Open http://localhost:8080 in your web browser." -ForegroundColor Green
go run cmd/serve/main.go
