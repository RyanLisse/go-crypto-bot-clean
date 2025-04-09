FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o crypto-bot .

# Create a minimal runtime image
FROM alpine:latest

# Add non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install runtime dependencies and security updates
RUN apk add --no-cache ca-certificates tzdata sqlite && \
    apk upgrade --no-cache

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/crypto-bot /app/crypto-bot

# Copy configuration files
COPY configs/ /app/configs/
COPY deploy/config.production.yaml /app/configs/config.production.yaml

# Create directories for data persistence
RUN mkdir -p /app/data/db /app/data/logs /app/data/backtest && \
    chown -R appuser:appgroup /app

# Set environment variables
ENV GIN_MODE=release
ENV CONFIG_PATH=/app/configs
ENV DB_PATH=/app/data/db
ENV LOG_PATH=/app/data/logs
ENV CONFIG_FILE=config.production.yaml

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/crypto-bot"]

# Default command
CMD ["serve"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
