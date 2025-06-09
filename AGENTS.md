# Development Guidelines

This repository contains a Go command line application. To keep the codebase consistent and reliable, follow these rules when contributing:

1. **Test first.** When adding a feature or fixing a bug, start by writing a failing test in the `cmd` package. Then implement the code that makes the test pass.
2. **Format code.** Run `gofmt -w` on all modified Go files before committing.
3. **Run tests.** Use `./test.sh` to execute `go test ./...` and ensure all tests succeed.
4. **No binaries.** Do not commit compiled executables or build artifacts.
