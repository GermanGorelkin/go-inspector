# Makefile for go-inspector

# Variables
BINARY_NAME=cli
COVERAGE_FILE=coverage.out

# Default target
.PHONY: all
all: fmt lint test build

# Format code
.PHONY: fmt
fmt:
	gofmt -w .

# Lint code
.PHONY: lint
lint:
	go vet ./...
	gofmt -l .

# Run tests
.PHONY: test
test:
	go test ./... -v

# Run tests with coverage
.PHONY: coverage
coverage:
	go test ./... -coverprofile=$(COVERAGE_FILE)
	go tool cover -func=$(COVERAGE_FILE)

# Build library and CLI
.PHONY: build
build:
	go build ./...
	go build -o $(BINARY_NAME) cmd/cli/main.go

# Clean up artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download
