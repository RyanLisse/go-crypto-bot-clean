FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
