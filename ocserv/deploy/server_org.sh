#!/bin/bash

set -euo pipefail

OCD_LOGFILE="/var/log/ocserv.log"
touch "$OCD_LOGFILE"
chmod 644 "$OCD_LOGFILE"

# --- Trap signals to stop child processes gracefully ---
cleanup() {
  echo "[INFO] Caught signal. Stopping processes..." | tee -a "$OCD_LOGFILE"

  kill -TERM "$SERVICE_PID" 2>/dev/null || true
  echo "[INFO] Ocserv service API stopped successfully..." | tee -a "$OCD_LOGFILE"

  if [[ -n "${OCSERV_PID:-}" ]]; then
    kill -TERM -- "-$OCSERV_PID" 2>/dev/null || true
    echo "[INFO] Ocserv (PGID: $OCSERV_PID) stopped successfully..." | tee -a "$OCD_LOGFILE"
  else
    echo "[WARN] Ocserv PID not found; skipping stop." | tee -a "$OCD_LOGFILE"
  fi

  wait "$SERVICE_PID" "$OCSERV_WRAPPER_PID"
  echo "[INFO] All services stopped. Exiting." | tee -a "$OCD_LOGFILE"
  exit 0
}
trap cleanup SIGINT SIGTERM

# --- ASCII banner for handler service ---
cat <<'EOF'

  ___                                _   _                 _ _                ____                  _
 / _ \  ___ ___  ___ _ ____   __    | | | | __ _ _ __   __| | | ___ _ __     / ___|  ___ _ ____   _(_) ___ ___
| | | |/ __/ __|/ _ \ '__\ \ / /    | |_| |/ _` | '_ \ / _` | |/ _ \ '__|    \___ \ / _ \ '__\ \ / / |/ __/ _ \
| |_| | (__\__ \  __/ |   \ V /     |  _  | (_| | | | | (_| | |  __/ |        ___) |  __/ |   \ V /| | (_|  __/
 \___/ \___|___/\___|_|    \_/      |_| |_|\__,_|_| |_|\__,_|_|\___|_|       |____/ \___|_|    \_/ |_|\___\___|


   🧩 ocserv handler service (Go-based management)
   📄 Logs: stdout (Docker logging)
   🕓 Started at: $(date)

──────────────────────────────────────────────────────

EOF

echo "[INFO] Starting ocserv_service..."
/ocserv_service &
SERVICE_PID=$!

echo "[INFO] ocserv_service PID: $SERVICE_PID"
echo "[INFO] Waiting for ocserv_service to initialize..."
sleep 2

# --- ASCII banner for ocserv VPN server ---
cat <<'EOF'

  ___                                ____
 / _ \  ___ ___  ___ _ ____   __    / ___|  ___ _ ____   _____ _ __
| | | |/ __/ __|/ _ \ '__\ \ / /    \___ \ / _ \ '__\ \ / / _ \ '__|
| |_| | (__\__ \  __/ |   \ V /      ___) |  __/ |   \ V /  __/ |
 \___/ \___|___/\___|_|    \_/      |____/ \___|_|    \_/ \___|_|


   🛡️  OpenConnect VPN Server (ocserv)
   📄 Logs: /logs/ocserv.log + stdout
   ⚙️  Config: /etc/ocserv/ocserv.conf
   🕓 Started at: $(date)

──────────────────────────────────────────────────────

EOF

echo "[INFO] Starting ocserv..."
# Run in a new session to isolate PGID and redirect output directly
#setsid /usr/sbin/ocserv --debug=9999 --foreground --config=/etc/ocserv/ocserv.conf >> "$OCD_LOGFILE" 2>&1 &

#setsid /usr/sbin/ocserv --debug=9999 --foreground --config=/etc/ocserv/ocserv.conf \
#  2>&1 | while IFS= read -r line; do
#    printf '[%s] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$line"
#  done >> "$OCD_LOGFILE" &

setsid /usr/sbin/ocserv --debug=9999 --foreground --config=/etc/ocserv/ocserv.conf \
  2>&1 | while IFS= read -r line; do
    printf '[%s] %s\n' "$(date +%s)" "$line"
  done >> "$OCD_LOGFILE" &


OCSERV_WRAPPER_PID=$!
OCSERV_PID=$(ps -o pgid= "$OCSERV_WRAPPER_PID" | tr -d ' ')

echo "[INFO] ocserv wrapper PID: $OCSERV_WRAPPER_PID, PGID: $OCSERV_PID"

# --- Wait for both services ---
wait "$SERVICE_PID" "$OCSERV_WRAPPER_PID"
