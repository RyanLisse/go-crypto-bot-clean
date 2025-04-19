#!/bin/bash

# This script kills any process running on the specified port
# Usage: ./kill_port.sh <port>

PORT=${1:-8080}  # Default to port 8080 if not specified

echo "Looking for processes using port $PORT..."

# Find the process ID using the port
PID=$(lsof -i :$PORT -t)

if [ -z "$PID" ]; then
  echo "No process found using port $PORT"
  exit 0
fi

echo "Found process(es) using port $PORT: $PID"
echo "Killing process(es)..."

# Kill the process(es)
for pid in $PID; do
  echo "Killing process $pid"
  kill -9 $pid
done

# Verify the port is now free
sleep 1
PID=$(lsof -i :$PORT -t)
if [ -z "$PID" ]; then
  echo "Port $PORT is now free"
else
  echo "Warning: Port $PORT is still in use by process(es): $PID"
fi
