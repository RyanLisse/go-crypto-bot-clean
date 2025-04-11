FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the minimal API code
COPY cmd/minimal-api/main.go ./cmd/minimal-api/main.go

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

# Expose port
EXPOSE 8080

# Run the minimal API
CMD ["/app/minimal-api"]
