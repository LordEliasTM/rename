@echo off

if exist build\* rmdir /s /q build

set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s" -o build\rename-win64.exe
set GOOS=linux
go build -ldflags "-s" -o build\rename-linux64
set GOOS=darwin
go build -ldflags "-s" -o build\rename-darwin64
