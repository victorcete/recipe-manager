.PHONY: build clean dev fmt lint test test-coverage test-coverage-detail test-coverage-html test-verbose

build:
	go build -o bin/mcp-server ./cmd/mcp

clean:
	rm -rf bin/ coverage.out coverage.html
	go clean

dev: build
	./bin/mcp-server

fmt:
	go fmt ./...

lint:
	golangci-lint run

test:
	go test ./...

test-coverage:
	go test -cover ./...

test-coverage-detail:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm coverage.out

test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-verbose:
	go test -v ./...
