#!/bin/sh
set -e

CERT_DIR="/etc/nginx/certs"
CERT_KEY="$CERT_DIR/cert.key"
CERT_PEM="$CERT_DIR/cert.pem"

mkdir -p "$CERT_DIR"

# Generate self-signed cert if missing
if [ ! -f "$CERT_KEY" ] || [ ! -f "$CERT_PEM" ]; then
  echo "Generating self-signed SSL certificate..."
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$CERT_KEY" \
    -out "$CERT_PEM" \
    -subj "/C=US/ST=CA/L=SanFrancisco/O=MyCompany/OU=IT/CN=localhost"
fi

# Run nginx in foreground
exec "$@"
