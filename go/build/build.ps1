param(
    [string]$appName = "shbfsmreg",
    [string]$rootPath = "C:\0_Priv\github.com\SHBF-SM-Registrering\go",
    [string]$srcPath = "C:\0_Priv\github.com\SHBF-SM-Registrering\go\src",
    [string]$version = "1.0.0"
)

$startPath = (Get-Location).path

$env:GOOS="windows"
$env:GOARCH="amd64"

Set-location "$($srcPath)"

try {
    go build -ldflags="-s -w -X main.version=$($version)" -o "$($rootPath)\build\bin\$($appName).exe" .\cmd
} catch {
    Write-Error "Failed to build the application: $_"
    exit 1
}

#Linux
$env:GOOS="linux"
$env:GOARCH="amd64"
try {
    go build -ldflags="-s -w -X main.version=$($version)" -o "$($rootPath)\build\bin\$($appName)_linux.exe" .\cmd
} catch {
    Write-Error "Failed to build the application: $_"
    exit 1
}

#Darwin (macos)
$env:GOOS="darwin"
$env:GOARCH="amd64"
try {
    go build -ldflags="-s -w -X main.version=$($version)" -o "$($rootPath)\build\bin\$($appName)_darwin.exe" .\cmd
} catch {
    Write-Error "Failed to build the application: $_"
    exit 1
}

$env:GOOS="windows"
Write-Host "Build completed successfully"
Set-location $startPath 