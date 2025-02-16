#!/bin/bash

if [[ $(id -u) != "0" ]]; then
    echo -e "\e[0;31m"Error: You must be root to run this install script."\e[0m"
    exit 1
fi
OS_VERSION=$(grep '^VERSION' /etc/os-release | grep "Focal Fossa" | wc -l)
if [ "${OS_VERSION}" -eq "0" ]; then
    echo "This script is only stable with Ubuntu 20.04(Focal Fossa)"
    exit 1
fi
if [ -z "$PORT" ]; then
    PORT=20443
fi
echo -e "\e[0;36m"Installing Ocserv..."\e[0m"
apt-get update
apt-get install -y ocserv gnutls-bin
if [ "$?" = "0" ]; then
    echo -e "\e[0;32m"Ocserv Installation Was Successful."\e[0m"
else
    echo -e "\e[0;31m"Ocserv Installation Is Failed"\e[0m"
    exit 1
fi
if [ ! -f /etc/ocserv/certs/cert.pem ]; then
    mkdir -p /etc/ocserv/certs
    touch /etc/ocserv/ocpasswd
    servercert="cert.pem"
    serverkey="key.pem"
    if [ -z "$CN" ]; then
        CN="End-way-Cisco-VPN"
    fi
    if [ -z "$ORG" ]; then
        ORG="End-way"
    fi
    if [ -z "$EXPIRE" ]; then
        EXPIRE=3650
    fi
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
        --template server.tmpl --outfile ${servercert}
    cp "${servercert}" /etc/ocserv/certs/cert.pem
    cp "${serverkey}" /etc/ocserv/certs/cert.key
fi
if [ ! -f '/etc/ocserv/ocserv.conf' ] || [ $(grep -r "custom config" /etc/ocserv/ocserv.conf | wc -l) == "0" ]; then
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
auth-timeout=240
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
tcp-port=${PORT}
udp-port=${PORT}
max-same-clients=2
ipv4-network=${OC_NET}
config-per-group=/etc/ocserv/groups/
EOT
    mkdir /etc/ocserv/defaults
    >/etc/ocserv/defaults/group.conf
    mkdir /etc/ocserv/groups
fi

# run in script -- UpdateIpTables
# iptables rules
# echo -e "\e[0;32m"Adding iptables rules."\e[0m"
# apt install -y iptables-persistent
# iptables -I INPUT -p tcp --dport ${PORT} -j ACCEPT
# iptables -I INPUT -p udp --dport ${PORT} -j ACCEPT
# iptables -I FORWARD -s ${OC_NET} -j ACCEPT
# iptables -I FORWARD -d ${OC_NET} -j ACCEPT
# iptables -t nat -A POSTROUTING -s ${OC_NET} -o ${ETH} -j MASQUERADE
# #iptables -t nat -A POSTROUTING -j MASQUERADE
# sh -c "iptables-save > /etc/iptables/rules.v4"
# sh -c "ip6tables-save > /etc/iptables/rules.v6"

echo "Enable IP forward."
sysctl -w net.ipv4.ip_forward=1
echo net.ipv4.ip_forward = 1 >>"/etc/sysctl.conf"
systemctl daemon-reload
echo "Enable ocserv service to start during bootup."
systemctl enable ocserv.service
systemctl restart ocserv.service
OCSERV_STATE=$(systemctl is-active ocserv)
if [ "$OCSERV_STATE" = "active" ]; then
    echo -e "\e[0;32m"Ocserv Is Started."\e[0m"
else
    echo -e "\e[0;31m"Ocserv Is Not Running."\e[0m"
    exit 1
fi
