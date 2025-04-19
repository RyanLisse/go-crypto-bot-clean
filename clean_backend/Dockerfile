FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker cmd/worker/main.go

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/bin/server /app/bin/server
COPY --from=builder /app/bin/worker /app/bin/worker

# Copy the .env.example file
COPY .env.example /app/.env

# Create data directory
RUN mkdir -p /app/data

# Expose the server port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/bin/server"]
