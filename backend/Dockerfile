FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the minimal API with CGO enabled for SQLite support
# Use build flags to reduce binary size
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o minimal-api ./cmd/minimal-api

# Create a minimal runtime image
FROM alpine:latest

# Install runtime dependencies and security updates
RUN apk update && \
    apk add --no-cache ca-certificates tzdata sqlite curl && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Create directories for data persistence with proper permissions
RUN mkdir -p /app/data/logs /app/configs && \
    chown -R appuser:appgroup /app

# Copy the binary from the builder stage
COPY --from=builder --chown=appuser:appgroup /app/minimal-api /app/minimal-api

# Copy configuration files
COPY --from=builder --chown=appuser:appgroup /app/configs /app/configs

# Set environment variables
ENV CONFIG_PATH=/app/configs
ENV CONFIG_FILE=config.minimal.yaml
ENV LOG_PATH=/app/data/logs

# Expose port
EXPOSE 8080

# Switch to non-root user
USER appuser

# Configure health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the minimal API
CMD ["/app/minimal-api"]
