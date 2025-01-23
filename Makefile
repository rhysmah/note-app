BINARY_NAME=note-app
GOARCH=amd64
VERSION?=0.1

BINARY_UNIX=$(BINARY_NAME)-$(VERSION)-linux
BINARY_OSX=$(BINARY_NAME)-$(VERSION)-darwin
BINARY_WIN=$(BINARY_NAME)-$(VERSION)-windows.exe

.DEFAULT_GOAL := build

clean:
	go clean
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_OSX)
	rm -f $(BINARY_WIN)

format: 
	go fmt ./...

vet: format
	go vet ./...

build: vet
	GOARCH=$(GOARCH) GOOS=darwin go build -o $(BINARY_OSX) main.go
	GOARCH=$(GOARCH) GOOS=linux go build -o $(BINARY_UNIX) main.go
	GOARCH=$(GOARCH) GOOS=windows go build -o $(BINARY_WIN) main.go

.PHONY: clean format vet build