@echo off

set GOROOT=c:\Tools\go-1.17.x
set GOPATH=c:\Lib\Golang
set PATH=%GOROOT%\bin;%PATH%

go mod init github.com/playrix/it-prod/examples/walk/tableview