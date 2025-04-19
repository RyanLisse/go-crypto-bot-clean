#!/bin/bash

# Default port is 8080, but can be overridden with PORT environment variable
PORT=${PORT:-8080}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to kill process on port
kill_process_on_port() {
    local port=$1
    echo "Checking if port $port is in use..."

    # Check if lsof is available
    if command_exists lsof; then
        if lsof -i :$port > /dev/null 2>&1; then
            echo "Port $port is in use. Killing process..."
            lsof -i :$port -t | xargs kill -9
            echo "Process killed."
            # Wait a moment to ensure the port is released
            sleep 1
        else
            echo "Port $port is not in use."
        fi
    # If lsof is not available, try netstat on Linux
    elif command_exists netstat; then
        if netstat -tuln | grep -q ":$port "; then
            echo "Port $port is in use. Killing process..."
            # Get PID using netstat and awk
            pid=$(netstat -tuln | grep ":$port " | awk '{print $7}' | cut -d'/' -f1)
            if [ -n "$pid" ]; then
                kill -9 $pid
                echo "Process killed."
                # Wait a moment to ensure the port is released
                sleep 1
            else
                echo "Could not find PID for process using port $port."
            fi
        else
            echo "Port $port is not in use."
        fi
    else
        echo "Warning: Neither lsof nor netstat is available. Cannot check if port $port is in use."
    fi
}

# Kill any process using the specified port
kill_process_on_port $PORT

# Start the server
echo "Starting server on port $PORT..."
PORT=$PORT ./server
