dev:
	go run cmd/pastebin/main.go
build:
	go build -o bin/pastebin cmd/pastebin/main.go