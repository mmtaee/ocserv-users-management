#!/bin/bash

GetModes() {
    printf "\n"
    modes=(
        'ocserv + user panel --- (systemd)'
        'ocserv ---------------- (systemd)'
        'panel  ---------------- (systemd)'
        'ocserv + user panel --- (docker) [only in ubuntu 20.04]'
        'oscerv ---------------- (docker)[only in ubuntu 20.04]'
    )
    color_echo "32;1" "Instalation modes:"
    for ((i = 0; i < "${#modes[@]}"; i++)); do
        color_echo "33" "$(printf "%4d: %s" $((i + 1)) "${modes[i]}")"
    done
    read -p "Enter the number corresponding to the desired mode: " mode
    if [[ "$mode" -ge 1 && "$mode" -le "${#modes[@]}" ]]; then
        selected_mode="${modes[$((mode - 1))]}"
    else
        color_echo "31" "\nInvalid selection. Please enter a valid mode number. ${mode} is out of range!\n"
        GetModes
    fi
    color_echo "36;1" "Selected mode: ${selected_mode}"

    if [ "$(echo "4 5" | grep -w $mode)" ] && [ "$(grep '^VERSION' /etc/os-release | grep "Focal Fossa" | wc -l)" -eq "0" ]; then
        color_echo "31;1" "This script is only stable with Ubuntu 20.04(Focal Fossa)"
        GetModes
    fi

    if [ "$(echo "4 5" | grep -w $mode)" ] && [ ! -x "$(command -v docker)" ]; then
        color_echo "31" "Docker service is not installed"
        exit 1
    fi
}

GetInterface() {
    printf "\n"
    interface_list=$(ip -o link show | awk '{print $2}' | grep -v lo | tr -d :)
    numbered_interfaces=$(echo "$interface_list" | awk '{print NR, $0}')
    color_echo "32;1" "Numbered interfaces:"
    IFS=$'\n'
    for line in $numbered_interfaces; do
        number=$(echo "$line" | awk '{print $1}')
        interface=$(echo "$line" | awk '{$1=""; print $0}' | sed 's/^[[:space:]]*//')
        color_echo "33" "$(printf "%4d: %s" "$number" "$interface")"
    done
    unset IFS
    read -p "Enter the number corresponding to the desired network interface: " interface_number
    if [[ "$interface_number" -ge 1 && "$interface_number" -le $(echo "$numbered_interfaces" | wc -l) ]]; then
        ETH=$(echo "$numbered_interfaces" | grep "^$interface_number" | awk '{$1=""; print $0}' | sed 's/^[[:space:]]*//')
    else
        color_echo "31" "\nInvalid selection. Please enter a valid number. ${interface_number} is out of range!\n"
        GetInterface
    fi
    color_echo "36;1" "Selected interface: $ETH"
}

GetIP() {
    printf "\n"
    HOST_IP=$(hostname -I | cut -d' ' -f1)
    color_echo "32;1" "Your host IP: ${HOST_IP}"
    read -p "Enter Your Host IP or let it blank : " host_ip
    if [[ -n "${host_ip}" ]]; then
        HOST_IP=${host_ip}
    fi
    color_echo "36;1" "Host IP: ${HOST_IP}"
}

GetDomain() {
    printf "\n"
    color_echo "35;3" "get public ip from myip.opendns.com"
    color_echo "35;3" "dig started ..."
    PUBLIC_IP=$(dig +short myip.opendns.com @resolver1.opendns.com)
    if [ "$?" != "0" ]; then
        PUBLIC_IP=${HOST_IP}
    fi
    echo "Your Public IP: ${PUBLIC_IP}"
    DOMAIN=
    read -p "Enter Your Domain or let it blank to use ${PUBLIC_IP}: " domain
    if [[ -n "${domain}" ]]; then
        DOMAIN=${domain}
    fi
    if [[ -n "${DOMAIN}" ]]; then
        color_echo "36;1" "Your Domain: ${DOMAIN}"
    else
        color_echo "36;1" "Your Public IP: ${PUBLIC_IP}"
    fi
}

GetPort() {
    printf "\n"
    PORT=20443
    read -p "Enter Your ocserv port or let it blank to use ${PORT}: " port
    if [[ -n "${port}" ]]; then
        PORT=${port}
    fi
    color_echo "36;1" "Ocserv Port: ${PORT}"
}

GetOcservVars() {
    printf "\n"
    CN="End-way-Cisco-VPN"
    read -p "Enter Your company name or let it blank to use ${CN}: " cn
    if [[ -n "${cn}" ]]; then
        CN=${cn}
    fi
    color_echo "36;1" "Your company name: ${CN}\n"

    ORG="End-way"
    read -p "Enter your organization name or let it blank to use ${ORG}: " org
    if [[ -n "${org}" ]]; then
        ORG=${org}
    fi
    color_echo "36;1" "Your organization name: ${ORG}\n"

    EXPIRE=3650
    read -p "Enter ssl expire days or let it blank to use ${EXPIRE} days: " expire
    if [[ -n "${expire}" ]]; then
        EXPIRE=${expire}
    fi
    color_echo "36;1" "Your organization name: ${EXPIRE}\n"

    OC_NET="172.16.24.0/24"
    read -p "Enter ocserv ipv4-network or let it blank to use ${OC_NET}: " oc_net
    if [[ -n "${oc_net}" ]]; then
        OC_NET=${oc_net}
    fi
    color_echo "36;1" "Your ocserv ipv4-network: ${OC_NET}\n"
}

UpdateIpTables() {
    color_echo "35;3" "Installing iptables-persistent ..."
    apt install -y iptables-persistent
    iptables -I INPUT -p tcp --dport ${PORT} -j ACCEPT
    iptables -I INPUT -p udp --dport ${PORT} -j ACCEPT
    iptables -I FORWARD -s ${OC_NET} -j ACCEPT
    iptables -I FORWARD -d ${OC_NET} -j ACCEPT
    iptables -t nat -A POSTROUTING -s ${OC_NET} -o ${ETH} -j MASQUERADE
    #iptables -t nat -A POSTROUTING -j MASQUERADE
    sh -c "iptables-save > /etc/iptables/rules.v4"
    sh -c "ip6tables-save > /etc/iptables/rules.v6"
}
