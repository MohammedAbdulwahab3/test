#!/bin/sh
# Start Alloy in background if configured
if [ ! -z "$PROMETHEUS_URL" ]; then
    echo "Starting Grafana Alloy..."
    /usr/bin/alloy run --server.http.listen-addr=0.0.0.0:12345 /etc/alloy/config.alloy &
fi

# Start the main application
echo "Starting Backend Server..."
./server
