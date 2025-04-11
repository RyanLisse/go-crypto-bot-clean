FROM golang:1.24-alpine

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 go build -o main ./cmd/api

# Create a minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=0 /app/main /app/main

# Copy configuration files
COPY configs/ /app/configs/
COPY deploy/config.production.yaml /app/configs/config.production.yaml

# Create directories for data persistence
RUN mkdir -p /app/data/db /app/data/logs /app/data/backtest

# Set environment variables
ENV CONFIG_PATH=/app/configs
ENV DB_PATH=/app/data/db
ENV LOG_PATH=/app/data/logs
ENV CONFIG_FILE=config.production.yaml

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
