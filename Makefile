.PHONY: help test build clean fmt vet lint coverage examples

# Default target
help:
	@echo "Available targets:"
	@echo "  make test      - Run all tests"
	@echo "  make build     - Build all packages"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make fmt       - Format code with gofmt"
	@echo "  make vet       - Run go vet"
	@echo "  make lint      - Run all linters"
	@echo "  make coverage  - Generate test coverage report"
	@echo "  make examples  - Build all examples"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Build all packages
build:
	@echo "Building packages..."
	@go build ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@go clean ./...
	@rm -f coverage.out coverage.html
	@find examples -type f -name "main" -delete

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run linters
lint: fmt vet
	@echo "Running linters..."
	@go mod tidy
	@go mod verify

# Generate test coverage
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build all examples
examples:
	@echo "Building examples..."
	@cd examples/simple && go build -o simple
	@cd examples/manager && go build -o manager
	@cd examples/custom && go build -o custom
	@echo "Examples built successfully"

# Run all checks before commit
check: lint test
	@echo "All checks passed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

