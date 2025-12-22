.PHONY: fmt test

all: fmt test

fmt:
	go fmt ./...

test:
	go test -v -cover ./...

bench:
	go test -v ./_examples/pool_test.go -bench=. -benchtime=1s