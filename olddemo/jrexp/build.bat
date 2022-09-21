@echo off

set GOARCH=386

set gocmd=%GOPATH%\bin\windows_386\go1.10.exe

%gocmd% build
