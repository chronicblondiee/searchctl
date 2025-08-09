.PHONY: build test test-unit test-integration test-all test-conditions test-config test-delete-confirmation test-lifecycle-policies start-test-env stop-test-env install clean release dev-deps

# Variables
BINARY_NAME=searchctl
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/chronicblondiee/searchctl/internal/version.Version=$(VERSION) -X github.com/chronicblondiee/searchctl/internal/version.Commit=$(COMMIT) -X github.com/chronicblondiee/searchctl/internal/version.Date=$(DATE)"

build:
	mkdir -p bin
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

test: test-unit

test-unit:
	@echo "Running unit tests..."
	go test -v ./pkg/... ./cmd/... ./internal/...

test-integration:
	@echo "Running integration tests..."
	./scripts/integration-test.sh

test-conditions:
	@echo "Running conditions tests..."
	./scripts/test-conditions.sh

test-config:
	@echo "Running config tests..."
	./scripts/test-config.sh

test-delete-confirmation:
	@echo "Running delete confirmation tests..."
	./scripts/test-delete-confirmation.sh

test-lifecycle-policies:
	@echo "Running lifecycle policies tests..."
	./scripts/test-lifecycle-policies.sh

test-all: test-unit test-integration test-conditions test-config test-delete-confirmation test-lifecycle-policies
	@echo "All tests completed!"

start-test-env:
	@echo "Starting test environment..."
	./scripts/start-test-env.sh

stop-test-env:
	@echo "Stopping test environment..."
	./scripts/stop-test-env.sh

test-env: start-test-env test-integration stop-test-env

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./pkg/... ./cmd/... ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-short:
	go test -short ./pkg/... ./cmd/... ./internal/...

install:
	go install $(LDFLAGS) .

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

release:
	goreleaser release --rm-dist

dev-deps:
	go install github.com/goreleaser/goreleaser@latest

test-env-start:
	@echo "Starting test environment..."
	./scripts/start-test-env.sh

test-env-stop:
	@echo "Stopping test environment..."
	./scripts/stop-test-env.sh

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet lint test

.DEFAULT_GOAL := build
