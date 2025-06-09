#!/bin/bash
set -e

go test ./...

# Ensure the code builds for all target platforms.
targets=("linux/amd64" "linux/arm64" "darwin/arm64" "windows/amd64")
for t in "${targets[@]}"; do
    IFS=/ read -r os arch <<<"$t"
    echo "building for $os/$arch"
    env GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o /tmp/todo_test_$os_$arch .
done
