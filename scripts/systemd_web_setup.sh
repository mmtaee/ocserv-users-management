#!/bin/bash

set -euo pipefail

## ===============================
## Deployment Script: Copy Built Binaries
## Description: Copies pre-built binaries to /opt/ocserv_user_management
##              and creates systemd units for each service
## ===============================
#
## Ensure running as root
#if [[ $(id -u) -ne 0 ]]; then
#    echo "❌ You must run this script as root (sudo)."
#    exit 1
#fi
#
#bin_dir="/opt/ocserv_user_management"
#mkdir -p "$bin_dir"
#
## Detect OS and architecture
#OS=$(uname | tr '[:upper:]' '[:lower:]')
#ARCH=$(uname -m)
#
## Map uname architecture to Go architecture
#case "$ARCH" in
#    x86_64) ARCH="amd64" ;;
#    i386|i686) ARCH="386" ;;
#    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
#esac
#
#echo "ℹ️ Detected OS: $OS, ARCH: $ARCH"
#
## ===============================
## Configuration: Services
## Key = service name, Value = path to project
## ===============================
#declare -A SERVICES=(
#    ["api"]="./api"
#    ["stream_log"]="./stream_log"
#)
#
#
#
#for service in "${!SERVICES[@]}"; do
#    systemctl stop "$service"
#done
#
## ===============================
## Copy binaries to deployment directory
## ===============================
#for service in "${!SERVICES[@]}"; do
#    project_dir="${SERVICES[$service]}"
#    src="${project_dir}/bin/${service}-${ARCH}"
#    dest="${bin_dir}/${service}"
#
#    if [[ ! -f "$src" ]]; then
#        echo "❌ Binary not found for $service: $src"
#        exit 1
#    fi
#
#    echo "📦 Installing $service -> $dest"
#    cp "$src" "$dest"
#    chmod +x "$dest"
#done
#
#echo "✅ All binaries installed successfully in $bin_dir"
#
## ===============================
## Create systemd units
## ===============================
#ENV_FILE="${bin_dir}/ocserv_user_management.env"
#
#cp .env $ENV_FILE || exit 1
#
#for service in "${!SERVICES[@]}"; do
#    unit_file="/etc/systemd/system/${service}.service"
#    binary="${bin_dir}/${service}"
#
#    # Determine port
#    if [[ "$service" == "api" ]]; then
#        ARGS="serve --host 127.0.0.1 --port 8080"
#    elif [[ "$service" == "stream_log" ]]; then
#        ARGS="-h 127.0.0.1 -p 8081 --systemd"
#    fi
#
#    echo "⚙️ Creating systemd unit for $service -> $unit_file"
#
#    cat > "$unit_file" <<EOF
#[Unit]
#Description=$service service
#After=network.target
#
#[Service]
#Type=simple
#EnvironmentFile=${ENV_FILE}
#ExecStart=${binary} ${ARGS}
#Restart=on-failure
#RestartSec=5s
#User=root
#WorkingDirectory=${bin_dir}
#StandardOutput=journal
#StandardError=journal
#
#[Install]
#WantedBy=multi-user.target
#EOF
#done
#
## Reload systemd, enable and start services
#systemctl daemon-reload
#
#for service in "${!SERVICES[@]}"; do
#    systemctl enable "$service"
#    systemctl restart "$service"
#    echo "✅ Started $service service"
#done
#
#echo "🎉 All binaries deployed and systemd services created successfully."

apt update

apt install -y nginx

REQUIRED_VERSION="23.11.1"

# Get currently installed Node.js version (strip leading "v")
CURRENT_VERSION=$(node -v 2>/dev/null | sed 's/^v//')

# Compare two versions: returns 0 if $1 < $2
version_lt() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" != "$2" ]
}

if [[ -z "$CURRENT_VERSION" ]] || version_lt "$CURRENT_VERSION" "$REQUIRED_VERSION"; then
    echo "⚠️ Node.js not found or outdated (current: ${CURRENT_VERSION:-none})."
    echo "Installing or upgrading Node.js to version 23.x..."

    # Add NodeSource setup and install Node.js
    curl -fsSL https://deb.nodesource.com/setup_23.x | bash -
    apt-get install -y nodejs

   echo "✅ Node.js installed: $(node -v)"
else
    echo "✅ Node.js is already installed: v$CURRENT_VERSION (>= ${REQUIRED_VERSION})"
fi

cd ./web

echo "📦 Installing npm dependencies..."

npm install --legacy-peer-deps || { echo "❌ npm install failed"; exit 1; }

export VITE_I18N_LANGUAGES="${LANGUAGES}"

export NODE_ENV=production

echo "🏗 Building frontend..."

npm run build || { echo "❌ Frontend build failed"; exit 1; }

