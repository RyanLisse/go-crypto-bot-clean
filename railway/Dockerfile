FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o api ./cmd/api

# Create a minimal runtime image
FROM alpine:latest

# Add non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install runtime dependencies and security updates
RUN apk add --no-cache ca-certificates tzdata && \
    apk upgrade --no-cache

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api .

# Use non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./api"]
