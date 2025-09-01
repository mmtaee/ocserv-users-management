#!/bin/bash

# ==============================================================
# Deployment Script: ocserv_user_management + Web + Nginx
# ==============================================================
#
# Description:
#  - Builds Go services (api, stream_log) and installs them into /opt/ocserv_user_management
#  - Creates and enables systemd service units
#  - Builds frontend (npm + Vite) and deploys to /var/www/site
#  - Installs and configures Nginx with self-signed SSL and reverse proxy
#
# Usage:
#   sudo ./deploy.sh
#
# Prerequisites:
#  - Debian 12
#  - Git repo with `api`, `stream_log`, and `web/` directories
#  - Environment variables for SSL:
#       SSL_EXPIRE, SSL_C, SSL_ST, SSL_L, SSL_ORG, SSL_OU, SSL_CN
# ==============================================================

set -euo pipefail
trap 'echo "❌ Deployment failed at line $LINENO."; exit 1' ERR

echo "ℹ️ Starting deployment..."

sudo apt install -y gcc

# -----------------------
# Deployment directories
# -----------------------
BIN_DIR="/opt/ocserv_user_management"
sudo mkdir -p "$BIN_DIR"
echo "ℹ️ Using deployment directory: $BIN_DIR"

# -----------------------
# Detect OS and ARCH
# -----------------------
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    i386|i686) ARCH="386" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "ℹ️ Detected OS: $OS, ARCH: $ARCH"

# -----------------------
# Services configuration
# -----------------------
declare -A SERVICES=(
    ["api"]="./api"
    ["stream_log"]="./stream_log"
)

# -----------------------
# Build Go binaries
# -----------------------
for service in "${!SERVICES[@]}"; do
    project_dir="${SERVICES[$service]}"
    dest="${BIN_DIR}/${service}"

    echo "⚙️ Building $service from $project_dir ..."

    (
        cd "$project_dir" || { echo "❌ Missing project dir: $project_dir"; exit 1; }

        # Build binary in project dir
        CGO_ENABLED=1 GOOS=linux GOARCH="${ARCH}" \
            go build -ldflags="-s -w" -o "$service"

        # Move binary to BIN_DIR
        sudo mv "$service" "$dest"
    )

    echo "📦 Build $service completed"
    sudo chmod +x "$dest"
done
echo "✅ All binaries built and deployed into $BIN_DIR"

# -----------------------
# Stop existing services
# -----------------------
echo "ℹ️ Stopping existing services (if any)..."
for service in "${!SERVICES[@]}"; do
    sudo systemctl stop "$service" 2>/dev/null || true
done

# -----------------------
# Environment file
# -----------------------
ENV_FILE="${BIN_DIR}/ocserv_user_management.env"
if [[ -f ".env" ]]; then
    sudo cp .env "$ENV_FILE"
    echo "ℹ️ Copied environment file to $ENV_FILE"
else
    echo "⚠️ .env file not found, skipping environment copy"
fi

# -----------------------
# Create systemd units
# -----------------------
for service in "${!SERVICES[@]}"; do
    unit_file="/etc/systemd/system/${service}.service"
    binary="${BIN_DIR}/${service}"

    # Service-specific args
    case "$service" in
        api) ARGS="serve --host 127.0.0.1 --port 8080" ;;
        stream_log) ARGS="-h 127.0.0.1 -p 8081 --systemd" ;;
        *) ARGS="" ;;
    esac

    echo "⚙️ Creating systemd unit for $service -> $unit_file"
    sudo tee "$unit_file" > /dev/null <<EOF
[Unit]
Description=$service service
After=network.target

[Service]
Type=simple
EnvironmentFile=${ENV_FILE}
ExecStart=${binary} ${ARGS}
Restart=on-failure
RestartSec=5s
User=root
WorkingDirectory=${BIN_DIR}
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
done

echo "ℹ️ Reloading systemd and starting services..."
sudo systemctl daemon-reload

for service in "${!SERVICES[@]}"; do
    sudo systemctl enable --now "$service"
    echo "✅ Started $service service"
done

echo "🎉 Go services deployed successfully."

# -----------------------
# Node.js check & install
# -----------------------
echo "Checking Node.js version..."
REQUIRED_NODE_VERSION="23.11.1"
if ! command -v node >/dev/null 2>&1; then
    CURRENT_NODE_VERSION=""
else
    CURRENT_NODE_VERSION=$(node -v 2>/dev/null | sed 's/^v//')
fi

# Function to compare versions
version_lt() { dpkg --compare-versions "$1" lt "$2"; }

if [[ -z "$CURRENT_NODE_VERSION" ]] || version_lt "$CURRENT_NODE_VERSION" "$REQUIRED_NODE_VERSION"; then
    echo "⚠️ Node.js not found or outdated (current: ${CURRENT_NODE_VERSION:-none}). Installing Node.js 23.x..."
    curl -fsSL https://deb.nodesource.com/setup_23.x | sudo bash -
    sudo apt install -y nodejs
    echo "✅ Node.js installed: $(node -v)"
else
    echo "✅ Node.js is already installed: v$CURRENT_NODE_VERSION"
fi

# -----------------------
# Build frontend
# -----------------------
cd ./web
echo "📦 Installing npm dependencies..."
npm install --legacy-peer-deps || { echo "❌ npm install failed"; exit 1; }

export NODE_ENV=production
export VITE_I18N_LANGUAGES="${LANGUAGES:-en}"

echo "🏗 Building frontend..."
npm run build || { echo "❌ Frontend build failed"; exit 1; }

[ -d dist ] || { echo "❌ dist folder not found after build"; exit 1; }

# -----------------------
# Install and configure Nginx
# -----------------------
sudo apt update
sudo apt install -y nginx

sudo rm -rf /etc/nginx/sites-enabled/default

CERT_DIR="/etc/nginx/certs"
sudo mkdir -p "$CERT_DIR"
CERT_KEY="${CERT_DIR}/cert.key"
CERT_PEM="${CERT_DIR}/cert.pem"

if [[ ! -f "$CERT_KEY" || ! -f "$CERT_PEM" ]]; then
    echo "🔐 Generating self-signed SSL certificate..."
    sudo openssl req -x509 -nodes -days "${SSL_EXPIRE:-365}" -newkey rsa:2048 \
        -keyout "$CERT_KEY" -out "$CERT_PEM" \
        -subj "/C=${SSL_C:-US}/ST=${SSL_ST:-State}/L=${SSL_L:-City}/O=${SSL_ORG:-Org}/OU=${SSL_OU:-Unit}/CN=${SSL_CN:-localhost}"
fi

# -----------------------
# Write Nginx config
# -----------------------
sudo tee /etc/nginx/conf.d/site.conf > /dev/null <<'EOF'
upstream api_backend {
    server 127.0.0.1:8080;
}
upstream stream_log_backend {
    server 127.0.0.1:8081;
}
server {
    listen 3000 ssl http2;
    server_name _;

    ssl_certificate     /etc/nginx/certs/cert.pem;
    ssl_certificate_key /etc/nginx/certs/cert.key;

    location / {
        root /var/www/site;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location ~ ^/(api) {
        proxy_pass http://api_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /ws/ {
        proxy_pass http://stream_log_backend/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 86400s;
        proxy_send_timeout 86400s;
    }
}
EOF

# -----------------------
# Deploy frontend
# -----------------------
sudo mkdir -p /var/www/site
sudo cp -r dist/* /var/www/site
sudo chown -R www-data:www-data /var/www/site

# -----------------------
# Test & start Nginx
# -----------------------
echo "ℹ️ Testing Nginx configuration..."
if ! sudo nginx -t; then
    echo "❌ Nginx configuration test failed."
    exit 1
fi

sudo systemctl daemon-reload
sudo systemctl enable --now nginx.service

if sudo systemctl is-active --quiet nginx; then
    echo "✅ Nginx is running."
else
    echo "❌ Nginx failed to start."
    exit 1
fi

echo "🎉 Deployment completed successfully!"
