.PHONY: build clean dev docker-build docker-run docker-run-bg docker-stop fmt lint test test-coverage test-coverage-detail test-coverage-html test-verbose

# Build the project
build:
	go build -o bin/server ./cmd/server

# Clean binaries and artifacts
clean:
	rm -rf bin/ coverage.out coverage.html
	go clean

# Start development server
dev: build
	./bin/server

# Build Docker image
docker-build:
	docker build -t learn-go:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 learn-go:latest

# Run Docker container in background
docker-run-bg:
	docker run -d -p 8080:8080 --name learn-go-container learn-go:latest

# Stop and remove Docker container
docker-stop:
	docker stop learn-go-container || true
	docker rm learn-go-container || true

# Format the code
fmt:
	go fmt ./...

# Lint the code
lint:
	golangci-lint run

# Run tests
test:
	go test ./...

# Run tests with coverage report
test-coverage:
	go test -cover ./...

# Run tests with detailed coverage breakdown by function
test-coverage-detail:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm coverage.out

# Run tests with detailed coverage report and generate HTML
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose:
	go test -v ./...
