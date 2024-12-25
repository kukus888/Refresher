$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o refresher.exe .