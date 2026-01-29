# Makefile for go-inspector

# Variables
COVERAGE_FILE=coverage.out
EXAMPLES_DIR=examples

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

# Build library and examples
.PHONY: build
build:
	go build ./...
	go build ./$(EXAMPLES_DIR)/...

# Clean up artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(COVERAGE_FILE)
	find $(EXAMPLES_DIR) -type f -name 'main' -delete

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download
