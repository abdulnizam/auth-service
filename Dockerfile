# ---------- Stage 1: Builder ----------
FROM golang:1.21-alpine AS builder

# Install git (required for go mod if using private modules)
RUN apk add --no-cache git

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o auth-service ./cmd

# ---------- Stage 2: Final image ----------
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from builder
COPY --from=builder /app/auth-service .

# Copy .env (optional for testing; exclude in prod)
# COPY .env .

# Expose app port
EXPOSE 8080

# Run the binary
CMD ["./auth-service"]