# Stage 1 - Builder
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the binary with optimizations for size and static linking
# - CGO_ENABLED=0: Disable CGO for static binary
# - GOOS=linux: Target Linux OS
# - -a: Force rebuilding of packages
# - -installsuffix cgo: Add suffix to package directory
# - -ldflags '-w -s': Strip debug info and symbol table
# - -trimpath: Remove file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags '-w -s -extldflags "-static"' \
    -trimpath \
    -o homelab-api \
    ./cmd/api

# Stage 2 - Runtime
FROM alpine:latest

# Install runtime dependencies (ca-certificates for HTTPS, curl for health checks)
RUN apk add --no-cache ca-certificates curl tzdata && \
    # Create non-root user and group
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    # Create app directory
    mkdir -p /app && \
    chown -R appuser:appuser /app

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder --chown=appuser:appuser /build/homelab-api /app/homelab-api

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Add health check using curl to hit the /health endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set environment variable defaults
ENV PORT=8080

# Run the application
CMD ["/app/homelab-api"]
