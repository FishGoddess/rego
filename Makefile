.PHONY: fmt test bench

all: fmt test

fmt:
	go fmt ./...

test:
	go mod tidy
	go test -v -cover ./...

bench:
	go test -v ./_examples/performance_test.go -bench=. -benchtime=1s