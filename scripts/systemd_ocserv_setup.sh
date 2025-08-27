#!/bin/bash
# =========================================
# Script: ocserv_setup.sh
# Description: Install and configure Ocserv VPN server
# Requirements: Run as root
# =========================================

# Exit on error
set -e

# ===============================
# Function: print_message
# Description: Print colorized messages
# ===============================
print_message() {
    local type="$1"
    local msg="$2"
    local RED="\e[31m"
    local GREEN="\e[32m"
    local YELLOW="\e[33m"
    local BLUE="\e[36m"
    local RESET="\e[0m"

    case "$type" in
        info) echo -e "${BLUE}[INFO]${RESET} $msg" ;;
        success) echo -e "${GREEN}[SUCCESS]${RESET} $msg" ;;
        warn) echo -e "${YELLOW}[WARN]${RESET} $msg" ;;
        error) echo -e "${RED}[ERROR]${RESET} $msg" ;;
        *) echo "$msg" ;;
    esac
}

# ===============================
# Ensure script is run as root
# ===============================
if [[ $(id -u) != 0 ]]; then
    print_message error "You must be root to run this install script."
    exit 1
fi

print_message info "Installing Ocserv VPN server..."

# ===============================
# Install dependencies
# ===============================
apt-get update
DEBIAN_FRONTEND=noninteractive apt install -y --no-install-recommends \
    ocserv ca-certificates procps gnutls-bin iptables openssl \
    less dnsutils jq iptables-persistent

print_message success "Dependencies installed successfully."

# ===============================
# Configure ocserv
# ===============================
print_message info "Writing ocserv configuration..."
cat <<EOT >/etc/ocserv/ocserv.conf
auth="plain[passwd=/etc/ocserv/ocpasswd]"
run-as-user=root
run-as-group=root
socket-file=/var/run/ocserv-socket
isolate-workers=true
max-clients=1024
keepalive=32400
dpd=90
mobile-dpd=1800
switch-to-tcp-timeout=5
try-mtu-discovery=true
server-cert=/etc/ocserv/certs/cert.pem
server-key=/etc/ocserv/certs/cert.key
tls-priorities="NORMAL:%SERVER_PRECEDENCE:%COMPAT:-RSA:-VERS-SSL3.0:-ARCFOUR-128"
auth-timeout=40
min-reauth-time=300
max-ban-score=50
ban-reset-time=300
cookie-timeout=86400
deny-roaming=false
rekey-time=172800
rekey-method=ssl
use-occtl=true
pid-file=/var/run/ocserv.pid
device=vpns
predictable-ips=true
tunnel-all-dns=true
dns=${OCSERV_DNS}
ping-leases=false
mtu=1420
cisco-client-compat=true
dtls-legacy=true
tcp-port=${OCSERV_PORT}
udp-port=${OCSERV_PORT}
max-same-clients=2
ipv4-network=${OC_NET}
config-per-group=/etc/ocserv/groups/
config-per-user=/etc/ocserv/users/
log-level=3
rate-limit-ms=100
EOT

# Create directories
mkdir -p /etc/ocserv/defaults /etc/ocserv/groups /etc/ocserv/users/

# shellcheck disable=SC2188
> /etc/ocserv/defaults/group.conf

# ===============================
# Generate certificates if not exists
# ===============================
if [ ! -f /etc/ocserv/certs/cert.pem ]; then
    print_message info "Generating self-signed certificates..."
    mkdir -p /etc/ocserv/certs
    cd /etc/ocserv/certs || exit
    touch /etc/ocserv/ocpasswd
    servercert="cert.pem"
    serverkey="key.pem"

    # CA
    certtool --generate-privkey --outfile ca-key.pem
    cat <<EOF >ca.tmpl
cn = "${SSL_CN}"
organization = "${SSL_ORG}"
serial = 1
expiration_days = ${SSL_EXPIRE}
ca
signing_key
cert_signing_key
crl_signing_key
EOF
    certtool --generate-self-signed --load-privkey ca-key.pem \
        --template ca.tmpl --outfile ca-cert.pem

    # Server cert
    certtool --generate-privkey --outfile ${serverkey}
    cat <<EOF >server.tmpl
cn = "${SSL_CN}"
organization = "${SSL_ORG}"
serial = 2
expiration_days = ${SSL_EXPIRE}
signing_key
encryption_key
tls_www_server
EOF
    certtool --generate-certificate --load-privkey ${serverkey} \
        --load-ca-certificate ca-cert.pem --load-ca-privkey ca-key.pem \
        --template server.tmpl --outfile ${servercert} >>/tmp/cert.txt 2>&1

    # Optional: store pin info
    grep -r 'pin-sha256' /tmp/cert.txt | tr -d '[:space:]' >> /etc/ocserv/public_key_pin
    echo "Docker Host ip: $(hostname -i)" >> /etc/ocserv/public_key_pin

    rm -rf /tmp/cert.txt
    cp "${servercert}" /etc/ocserv/certs/cert.pem
    cp "${serverkey}" /etc/ocserv/certs/cert.key
fi

# ===============================
# Enable IP forwarding
# ===============================
print_message info "Enabling IP forwarding..."
sysctl -w net.ipv4.ip_forward=1
echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
sysctl -p

# Create TUN device if missing
mkdir -p /dev/net
if [ ! -c /dev/net/tun ]; then
    mknod /dev/net/tun c 10 200
fi

# ===============================
# Configure iptables for VPN
# ===============================
print_message info "Configuring iptables rules..."
# Open VPN port for TCP and UDP
iptables -I INPUT -p tcp --dport "${OCSERV_PORT}" -j ACCEPT
iptables -I INPUT -p udp --dport "${OCSERV_PORT}" -j ACCEPT

# Allow forwarding to/from VPN subnet
iptables -I FORWARD -s "${OC_NET}" -j ACCEPT
iptables -I FORWARD -d "${OC_NET}" -j ACCEPT

# Enable NAT for VPN subnet via selected interface
iptables -t nat -A POSTROUTING -s "${OC_NET}" -o "${ETH}" -j MASQUERADE

# Save rules persistently
sh -c "iptables-save > /etc/iptables/rules.v4"
sh -c "ip6tables-save > /etc/iptables/rules.v6"

# ===============================
# Enable and start ocserv service
# ===============================
print_message info "Enabling ocserv service to start on boot..."
systemctl daemon-reload
systemctl enable --now ocserv.service

if systemctl is-active --quiet ocserv; then
    print_message success "Ocserv is running."
else
    print_message error "Ocserv failed to start."
    exit 1
fi
