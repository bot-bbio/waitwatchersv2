# Windows PowerShell Build Script for WaitWatchersV2
Write-Host "=== Compiling WaitWatchersV2 ===" -ForegroundColor Cyan

# Save existing environments
$oldGoos = $env:GOOS
$oldGoarch = $env:GOARCH

# Build WASM
Write-Host "[BUILD] Compiling WebAssembly client (main.wasm)..." -ForegroundColor Yellow
$env:GOOS = "js"
$env:GOARCH = "wasm"
go build -o main.wasm cmd/wasm/main.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] main.wasm compiled successfully." -ForegroundColor Green
} else {
    Write-Host "[ERROR] Failed to compile main.wasm" -ForegroundColor Red
    $env:GOOS = $oldGoos
    $env:GOARCH = $oldGoarch
    exit 1
}

# Restore environment
$env:GOOS = $oldGoos
$env:GOARCH = $oldGoarch

# Build CLI
Write-Host "[BUILD] Compiling CLI tool (mach.exe)..." -ForegroundColor Yellow
go build -o mach.exe main.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] mach.exe compiled successfully." -ForegroundColor Green
} else {
    Write-Host "[ERROR] Failed to compile mach.exe" -ForegroundColor Red
    exit 1
}

Write-Host "=== Build process completed successfully! ===" -ForegroundColor Green
