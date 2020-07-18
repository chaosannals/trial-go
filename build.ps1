
if (!(Test-Path .\bin)) {
    mkdir .\bin
}
go build -o bin/trial.exe main.go