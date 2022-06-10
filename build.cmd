set VERSION=1.0.0
set GOPATH="c:\gopath"
go build -o todo.exe -ldflags="-X 'github.com/cklukas/todo/cmd.AppVersion=%VERSION%'"