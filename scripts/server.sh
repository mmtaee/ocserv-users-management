#!/bin/bash
set -e

# Forward signals to child processes
trap "echo '[INFO] Caught SIGTERM, stopping...'; kill -TERM $OCSERV_PID $API_PID 2>/dev/null" SIGTERM SIGINT

echo "[INFO] Starting services..."
/api serve -d &
API_PID=$!

echo "[INFO] Starting ocserv..."
/usr/sbin/ocserv --foreground --debug=999 --config=/etc/ocserv/ocserv.conf &
OCSERV_PID=$!

# Wait for any process to exit
wait -n

# If one process exits, kill the other
kill -TERM $OCSERV_PID $API_PID 2>/dev/null || true

# Wait again to let everything clean up
wait