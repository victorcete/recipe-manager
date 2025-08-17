# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .
RUN go build -o bin/mcp-server ./cmd/mcp

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/mcp-server .

# MCP servers communicate via stdio, not HTTP ports
CMD ["./mcp-server"]
