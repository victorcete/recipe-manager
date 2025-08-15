.PHONY: build test fmt clean

# Build the project
build:
	go build -o bin/server ./cmd/server

# Run tests
test:
	go test ./...

# Format the code
fmt:
	go fmt ./...

# Clean binaries and artifacts
clean:
	rm -rf bin/
	go clean