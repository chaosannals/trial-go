@echo off

:: 配置 GO 交叉编译 wasm 的环境变量
set GOOS=js
set GOARCH=wasm

:: 编译
go build -o public/demo.wasm main.go
