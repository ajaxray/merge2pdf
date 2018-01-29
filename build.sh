env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/merge2pdf_darwin-amd64
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/merge2pdf_linux-amd64
env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o builds/merge2pdf_linux-arm64
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o builds/merge2pdf_windows-386.exe