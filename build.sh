#!/bin/bash
VERSION="1.${GITHUB_RUN_NUMBER}.${GITHUB_RUN_ID}"
echo "Build command line tool with version $VERSION"
rm -f ./todo ./todo.exe ./todo_linux ./todo_linux_amd64 ./todo_linux_arm64 ./todo_mac_arm64
go mod tidy
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o todo_linux_amd64 -ldflags="-X 'github.com/cklukas/todo/cmd.AppVersion=$VERSION'" && \
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o todo_linux_arm64 -ldflags="-X 'github.com/cklukas/todo/cmd.AppVersion=$VERSION'" && \
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o todo_mac_arm64 -ldflags="-X 'github.com/cklukas/todo/cmd.AppVersion=$VERSION'" && \
GOOS=windows GOARCH=amd64 go build -o todo.exe -ldflags="-X 'github.com/cklukas/todo/cmd.AppVersion=$VERSION'" && \
chmod +x ./todo_linux_amd64 && \
chmod +x ./todo_linux_arm64 && \
chmod +x ./todo_mac_arm64 && \
(
    (./todo_linux_amd64 version && echo "(on Linux AMD64)") || 
    (./todo_mac_arm64 version && echo "(on Mac ARM64)") ||
    (./todo_linux_arm64 version && echo "(on Linux ARM64)")) && \
echo "[OK] build completed" || echo "[ERROR] build not OK"
