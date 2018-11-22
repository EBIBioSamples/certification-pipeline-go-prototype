build:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bscurate cmd/bscurate/main.go