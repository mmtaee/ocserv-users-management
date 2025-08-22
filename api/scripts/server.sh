#!/bin/bash

set -e

# Ensure required env vars are set
if [[ -z "${ADMIN_USERNAME}" || -z "${ADMIN_PASSWORD}" ]]; then
  echo "[ERROR] ADMIN_USERNAME and ADMIN_PASSWORD must be set"
  exit 1
fi

echo "[INFO] Creating admin user (if not exists)..."
/app/api create-admin -u "${ADMIN_USERNAME}" -p "${ADMIN_PASSWORD}" || true

echo "[INFO] Starting API server..."
exec /app/api serve

