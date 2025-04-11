FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the minimal API
RUN go build -o minimal-api ./cmd/minimal-api

# Create a minimal runtime image
FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache ca-certificates tzdata

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
