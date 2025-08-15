.PHONY: build clean run dev format

# Build the server binary
build:
	go build -o bin/server cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run the server directly
run:
	go run cmd/server/main.go

# Build and run for development
dev: build
	./bin/server

# Format and lint code
format:
	gofmt -s -w .
	go vet ./...
