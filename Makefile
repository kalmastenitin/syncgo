include e2e/e2e.mk

GOBIN := $(shell pwd)/bin
MODULE  := github.com/syncgo/syncgo

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X main.version=$(VERSION) \
           -X main.commit=$(COMMIT)  \
           -X main.date=$(DATE)

.PHONY: tools generate build test run demo lint

tools:
	mkdir -p bin
	GOBIN=$(GOBIN) go install tool

generate: tools
	go generate ./...

build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin ./...

lint: tools
	$(GOBIN)/golangci-lint run ./...

test:
	go test ./... -v -short

run: build
	./bin/syncgo $(ARGS)

demo: build prepare
	./bin/syncgo --config ./e2e/config.yaml
