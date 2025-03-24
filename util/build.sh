#!/bin/bash

# Check if API key is provided
if [ -z "$1" ]; then
    echo "Error: API key is required"
    echo "Usage: ./build.sh <your-api-key>"
    exit 1
fi

API_KEY="$1"

# Create a temporary file with the API key
sed "s/API_KEY_PLACEHOLDER/$API_KEY/" main.go > main.go.tmp

# Build the binary
echo "Building scanner binary..."
GOOS=linux GOARCH=amd64 go build -o scanner main.go.tmp

# Clean up temporary file
rm main.go.tmp

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful! Binary created as 'scanner'"
    echo "You can now copy the 'scanner' binary to your server"
else
    echo "Build failed!"
    exit 1
fi 