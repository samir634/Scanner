# Check if API key is provided
param(
    [Parameter(Mandatory=$true)]
    [string]$ApiKey
)

# Create a temporary file with the API key
$content = Get-Content -Path "main.go" -Raw
$content = $content.Replace("API_KEY_PLACEHOLDER", $ApiKey)
$content | Set-Content -Path "main.go.tmp"

# Build the binary
Write-Host "Building scanner binary..."
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o scanner main.go.tmp

# Clean up temporary file
Remove-Item -Path "main.go.tmp" -Force

# Check if build was successful
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful! Binary created as 'scanner'"
    Write-Host "You can now copy the 'scanner' binary to your server"
} else {
    Write-Host "Build failed!"
    exit 1
} 