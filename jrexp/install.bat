@echo off

set GOARCH=386

set gocmd=%GOPATH%\bin\windows_386\go1.10.exe

echo %gocmd%

@echo on
%gocmd% get github.com/lestrrat-go/file-rotatelogs
%gocmd% get github.com/sirupsen/logrus
