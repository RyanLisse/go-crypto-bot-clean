## 5. Deployment and Configuration

### 5.1 Environment Variables

Configure the application using environment variables:

```bash
# AI Provider
AI_PROVIDER=gemini  # or openai, anthropic
AI_PROVIDER_API_KEY=your_api_key_here
AI_PROVIDER_MODEL=gemini-flash  # or gpt-4, claude-3-opus

# Database
DATABASE_URL=libsql://your-turso-db-url
DATABASE_AUTH_TOKEN=your_turso_auth_token

# Security
JWT_SECRET=your_jwt_secret_here
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://your-production-domain.com

# Rate Limiting
MAX_REQUESTS_PER_MINUTE=60
MAX_TOKENS_PER_DAY=100000

# Logging
LOG_LEVEL=info  # debug, info, warn, error
```

### 5.2 Docker Deployment

Deploy the application using Docker:

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app ./cmd/server

# Create a minimal image
FROM alpine:3.18

COPY --from=builder /go/bin/app /app

# Set environment variables
ENV AI_PROVIDER=gemini
ENV AI_PROVIDER_MODEL=gemini-flash
ENV LOG_LEVEL=info

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/app"]
```

### 5.3 Kubernetes Deployment

Deploy the application to Kubernetes:

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: crypto-bot-ai
  labels:
    app: crypto-bot-ai
spec:
  replicas: 3
  selector:
    matchLabels:
      app: crypto-bot-ai
  template:
    metadata:
      labels:
        app: crypto-bot-ai
    spec:
      containers:
      - name: crypto-bot-ai
        image: your-registry/crypto-bot-ai:latest
        ports:
        - containerPort: 8080
        env:
        - name: AI_PROVIDER
          value: "gemini"
        - name: AI_PROVIDER_MODEL
          value: "gemini-flash"
        - name: AI_PROVIDER_API_KEY
          valueFrom:
            secretKeyRef:
              name: ai-secrets
              key: api-key
        - name: DATABASE_URL
          valueFrom:
            configMapKeyRef:
              name: db-config
              key: url
        - name: DATABASE_AUTH_TOKEN
          valueFrom:
            secretKeyRef:
              name: db-secrets
              key: auth-token
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-secrets
              key: jwt-secret
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
```

## 6. Conclusion

This implementation guide provides a comprehensive approach to integrating AI capabilities into the go-crypto-bot-migration project. By following these patterns and best practices, you can create a robust, secure, and performant AI assistant that enhances the trading experience for your users.

Key benefits of this implementation include:

1. **Clean Architecture**: Following the project's dependency injection pattern for maintainable code
2. **Security**: Proper API key management, input validation, and rate limiting
3. **Performance**: Optimized prompts, caching, and streaming responses
4. **User Experience**: Rich, interactive interface with data visualization
5. **Reliability**: Comprehensive testing and monitoring

Remember to regularly update your AI models and prompts as new capabilities become available, and continuously monitor usage to optimize costs and performance.
