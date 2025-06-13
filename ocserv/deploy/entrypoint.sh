#!/bin/bash

if [ -z "$CN" ]; then
    CN="End-way-Cisco-VPN"
fi
if [ -z "$ORG" ]; then
    ORG="End-way"
fi
if [ -z "$EXPIRE" ]; then
    EXPIRE=3650
fi
if [ -z "$OC_NET" ]; then
    OC_NET=172.16.24.0/24
fi
if [ -z "$HOST" ]; then
    HOST=$(dig +short myip.opendns.com @resolver1.opendns.com)
    if [ "$?" != "0" ]; then
        HOST=$(hostname -I | cut -d' ' -f1)
    fi
fi

cat <<EOT >/etc/ocserv/ocserv.conf
# custom config
auth="plain[passwd=/etc/ocserv/ocpasswd]"
run-as-user=root
run-as-group=root
socket-file=ocserv.sock
chroot-dir=/run
isolate-workers=true
max-clients=1024
keepalive=32400
dpd=90
mobile-dpd=1800
switch-to-tcp-timeout=5
try-mtu-discovery=true
server-cert=/etc/ocserv/certs/cert.pem
server-key=/etc/ocserv/certs/cert.key
#tls-priorities="NORMAL:%SERVER_PRECEDENCE:%COMPAT:-VERS-SSL3.0"
tls-priorities="NORMAL:%SERVER_PRECEDENCE:%COMPAT:-VERS-SSL3.0:-VERS-TLS1.0:-VERS-TLS1.1"
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
dns=${DNS}
ping-leases=false
mtu=1420
cisco-client-compat=true
dtls-legacy=true
tcp-port=443
udp-port=443
max-same-clients=2
ipv4-network=${OC_NET}
config-per-group=/etc/ocserv/groups/
config-per-user=/etc/ocserv/users/
log-level=2
EOT

mkdir -p /etc/ocserv/defaults /etc/ocserv/groups /etc/ocserv/users/

>/etc/ocserv/defaults/group.conf

if [ ! -f /etc/ocserv/certs/cert.pem ]; then
    mkdir -p /etc/ocserv/certs
    cd /etc/ocserv/certs
    touch /etc/ocserv/ocpasswd
    servercert="cert.pem"
    serverkey="key.pem"
    certtool --generate-privkey --outfile ca-key.pem
    cat <<_EOF_ >ca.tmpl
cn = "${CN}"
organization = "${ORG}"
serial = 1
expiration_days = ${EXPIRE}
ca
signing_key
cert_signing_key
crl_signing_key
_EOF_
    certtool --generate-self-signed --load-privkey ca-key.pem \
        --template ca.tmpl --outfile ca-cert.pem
    certtool --generate-privkey --outfile ${serverkey}
    cat <<_EOF_ >server.tmpl
cn = "${CN}"
organization = "${ORG}"
serial = 2
expiration_days = ${EXPIRE}
signing_key
encryption_key
tls_www_server
_EOF_
    certtool --generate-certificate --load-privkey ${serverkey} \
        --load-ca-certificate ca-cert.pem --load-ca-privkey ca-key.pem \
        --template server.tmpl --outfile ${servercert} >>/tmp/cert.txt 2>&1
    echo "Server Cert pin: $(grep -r 'pin-sha256' /tmp/cert.txt | tr -d '[:space:]')" >>/etc/ocserv/public_key_pin
    echo "Docker Host ip: $(hostname -i)" >>/etc/ocserv/public_key_pin
    rm -rf /tmp/cert.txt
    cp "${servercert}" /etc/ocserv/certs/cert.pem
    cp "${serverkey}" /etc/ocserv/certs/cert.key
fi

iptables -t nat -A POSTROUTING -j MASQUERADE
sysctl -w net.ipv4.ip_forward=1 # ipv4 ip forward
sysctl -p
mkdir -p /dev/net               #TUN device
mknod /dev/net/tun c 10 200
chmod 600 /dev/net/tun
mkdir -p /var/log/ocserv/

exec "$@"