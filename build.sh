#!/bin/bash
VERSION="1.${GITHUB_RUN_NUMBER}.${GITHUB_RUN_ID}"
echo "Build command line tool with version $VERSION"
rm -f ./todo ./todo.exe
CGO_ENABLED=0 go build -o todo -ldflags="-X 'github.com/cklukas/todo/cmd.AppVersion=$VERSION'" && \
GOOS=windows GOARCH=amd64 go build -o todo.exe -ldflags="-X 'github.com/cklukas/todo/cmd.AppVersion=$VERSION'" && \
chmod +x ./todo && \
./todo version && \
echo "[OK] build todo and todo.exe"
