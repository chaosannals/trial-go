@echo off
set CGO_ENABLED=0
set GOOS=js
set GOARCH=wasm
go build -o demo/demo.wasm main.go
