.PHONY: build test install clean release dev-deps

# Variables
BINARY_NAME=searchctl
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/chronicblondiee/searchctl/internal/version.Version=$(VERSION) -X github.com/chronicblondiee/searchctl/internal/version.Commit=$(COMMIT) -X github.com/chronicblondiee/searchctl/internal/version.Date=$(DATE)"

build:
	mkdir -p bin
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

test:
	go test ./...

install:
	go install $(LDFLAGS) .

clean:
	rm -rf bin/

release:
	goreleaser release --rm-dist

dev-deps:
	go install github.com/goreleaser/goreleaser@latest

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

.DEFAULT_GOAL := build
