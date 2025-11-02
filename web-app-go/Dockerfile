# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for swag generation
RUN apk add --no-cache git

# Copy go.mod and go.sum first (caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
RUN swag init -g cmd/server/main.go -o docs

# Build the Go binary
RUN go build -o /text-editor ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy built binary and docs
COPY --from=builder /text-editor /app/text-editor
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/public /app/public

# Expose port
EXPOSE 5007

# Create files directory
RUN mkdir -p /app/files

# Run the app
CMD ["/app/text-editor"]
