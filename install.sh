#!/bin/bash

if [[ $(id -u) != "0" ]]; then
    echo -e "\e[0;31m"Error: You must be root to run this install script."\e[0m"
    exit 1
fi

color_echo() {
    local color=$1
    local text=$2
    echo -e "\e[${color}m${text}\e[0m"
}

GetModes() {
    printf "\n"
    modes=(
        'ocserv + user panel --- (docker)'
        'ocserv + user panel --- (systemd) [only in ubuntu 20.04]'
        'panel  ---------------- (systemd)'
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

    if [ "$(echo "2 3" | grep -w $mode)" ] && [ "$(grep '^VERSION' /etc/os-release | grep "Focal Fossa" | wc -l)" -eq "0" ]; then
        color_echo "31;1" "This script is only stable with Ubuntu 20.04(Focal Fossa)"
        GetModes
    fi

    if [ "$(echo "1" | grep -w $mode)" ] && [ ! -x "$(command -v docker)" ]; then
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
        VALIDATE="^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$"
        if [[ "${domain}" =~ "${VALIDATE}" ]]; then
            DOMAIN=${domain}
        fi
    fi
    CORS_ALLOWED=""
    if [[ -n "${DOMAIN}" ]]; then
        color_echo "36;1" "Your Domain: ${DOMAIN}"
        CORS_ALLOWED="http://${DOMAIN},https://${DOMAIN}"
        HOST=${DOMAIN}
    elif [[ -n "${PUBLIC_IP}" ]]; then
        color_echo "36;1" "Your Public IP: ${PUBLIC_IP}"
        CORS_ALLOWED="http://${PUBLIC_IP},https://${PUBLIC_IP}"
        HOST=${PUBLIC_IP}
    else
        HOST=${HOST_IP}
    fi
    CORS_ALLOWED="${CORS_ALLOWED},http://${HOST_IP},https://${HOST_IP}"
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

GetDNS() {
    printf "\n"
    DNS="8.8.8.8"
    read -p "Enter Your DNS or let it blank to use ${DNS}: " dns
    if [[ -n "${dns}" ]]; then
        DNS=${dns}
    fi
    color_echo "36;1" "DNS: ${DNS}"
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

UpdateDockerProdEnv() {
    file_path="./prod.env"
    declare -A variables=(
        ["ORG"]="${ORG}"
        ["EXPIRE"]="${EXPIRE}"
        ["CN"]="${CN}"
        ["OC_NET"]="${OC_NET}"
        ["CORS_ALLOWED"]="${CORS_ALLOWED}"
        ["HOST"]="${HOST}"
        ["DOMAIN"]="${DOMAIN}"
        ["PORT"]="${PORT}"
    )
    if [ -e "$file_path" ]; then
        truncate -s 0 ${file_path}
    else
        touch "$file_path"
    fi
    for key in "${!variables[@]}"; do
        echo "$key=${variables[$key]}" >>"$file_path"
    done
}

GetModes

GetIP

GetDomain

GetPort

GetDNS

GetOcservVars

if [[ $mode != '1' ]]; then
    chmod +x ./configs/ocserv.sh
    chmod +x ./configs/panel.sh
    GetInterface
fi

color_echo "35;1" "Installing ${selected_mode} .............."

export CN ORG EXPIRE OC_NET DOMAIN HOST PORT DNS

if [[ $mode == '1' ]]; then
    UpdateDockerProdEnv
    DOCKER_VARS="CN=${CN} ORG=${ORG} EXPIRE=${EXPIRE} OC_NET=${OC_NET} DOMAIN=${DOMAIN} HOST=${HOST} PORT=${PORT} DNS=${DNS}"
    BUILD="DOCKER_SCAN_SUGGEST=false ${DOCKER_VARS} docker compose up -d --build"
    eval ${BUILD}
elif [[ $mode == '2' ]]; then
    ./configs/ocserv.sh
    ./configs/panel.sh
    UpdateIpTables
elif [[ $mode == '3' ]]; then
    ./configs/panel.sh
fi
