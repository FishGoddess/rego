.PHONY: fmt test

all: fmt test

fmt:
	go fmt ./...

test:
	go mod tidy
	go test -v -cover ./...