#!/bin/bash

set -e

echo "[INFO] Starting services..."

# Check if webhook exists and is executable
if [ -x "/app/webhook" ]; then
  echo "[INFO] Starting webhook..."
  /app/webhook &
else
  echo "[WARN] /webhook not found or not executable, skipping..."
fi

echo "[INFO] Starting ocserv..."
setsid /usr/sbin/ocserv --foreground --config=/etc/ocserv/ocserv.conf &

# Wait for any process to exit
wait -n

# Exit with the status of the first process that exits
exit $?