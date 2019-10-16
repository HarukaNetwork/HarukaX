#!/bin/bash

# Build linux targets
echo "Building for linux"

echo "Building x86_64"
env GOOS=linux GOARCH=amd64 go build -o build/ginko_linux_x86_64

echo "Building x86"
env GOOS=linux GOARCH=386 go build -o build/ginko_linux_x86

echo "Building arm"
env GOOS=linux GOARCH=arm go build -o build/ginko_linux_arm

echo "Building arm64"
env GOOS=linux GOARCH=arm64 go build -o build/ginko_linux_arm64

# Build windows targets
echo
echo "Building for windows"

echo "Building x86_64"
env GOOS=windows GOARCH=amd64 go build -o build/ginko_win_x86_64.exe

echo "Building x86"
env GOOS=windows GOARCH=386 go build -o build/ginko_win_x86.exe

# Build mac targets
echo
echo "Building for MacOS"

echo "Building x86_64"
env GOOS=darwin GOARCH=amd64 go build -o build/ginko_darwin_x86_64

echo "Building x86"
env GOOS=darwin GOARCH=386 go build -o build/ginko_darwin_x86