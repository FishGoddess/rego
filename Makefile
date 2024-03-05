.PHONY: fmt test

fmt:
	go fmt ./...

test:
	go mod tidy
	go test -cover -count=1 -test.cpu=1 ./...

all: fmt test