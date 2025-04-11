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
RUN CGO_ENABLED=1 go build -o minimal-api ./cmd/minimal-api

# Create a minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/minimal-api /app/minimal-api

# Copy configuration files
COPY --from=builder /app/configs /app/configs

# Create directories for data persistence
RUN mkdir -p /app/data/logs

# Set environment variables
ENV CONFIG_PATH=/app/configs
ENV CONFIG_FILE=config.minimal.yaml
ENV LOG_PATH=/app/data/logs

# Expose port
EXPOSE 8080

# Run the minimal API
CMD ["/app/minimal-api"]
