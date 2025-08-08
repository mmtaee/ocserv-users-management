#!/bin/bash

set -euo pipefail

# --- Trap signals to stop child processes gracefully ---

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

setsid /usr/sbin/ocserv --debug=9999 --foreground --config=/etc/ocserv/ocserv.conf &



wait -n
exit $?