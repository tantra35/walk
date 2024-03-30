@echo off

set GOROOT=c:\Tools\go-1.17.x
set GOPATH=c:\Lib\Golang
set PATH=%GOROOT%\bin;%PATH%

rem go build -ldflags="-H windowsgui"
go build
