.PHONY: build
build:
	go build -v ./cmd/MyEcho

.PHONY: test
test:
	go test -v -race -timeout 30s ./ ...

.PHONY: start
start:
	go run ./cmd/MyEcho

.DEFAULT_GOAL := build