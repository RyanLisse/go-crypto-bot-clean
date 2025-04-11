FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache netcat-openbsd

# Set working directory
WORKDIR /app

# Create a simple health check script
RUN echo '#!/bin/sh' > health.sh && \
    echo 'while true; do' >> health.sh && \
    echo '  echo -e "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nOK" | nc -l -p 8080' >> health.sh && \
    echo 'done' >> health.sh && \
    chmod +x health.sh

# Expose port
EXPOSE 8080

# Run the health check script
CMD ["/app/health.sh"]
