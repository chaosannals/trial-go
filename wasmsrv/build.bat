set GOOS=js
set GOARCH=wasm
go build -o public/demo.wasm main.go
