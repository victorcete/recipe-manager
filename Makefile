.PHONY: build test test-verbose test-coverage fmt clean dev

# Build the project
build:
	go build -o bin/server ./cmd/server

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage report
test-coverage:
	go test -cover ./...

# Run tests with detailed coverage report and generate HTML
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format the code
fmt:
	go fmt ./...

# Clean binaries and artifacts
clean:
	rm -rf bin/ coverage.out coverage.html
	go clean

# Start development server
dev: build
	./bin/server