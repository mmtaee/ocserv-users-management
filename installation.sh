#!/bin/bash

if [[ $(id -u) != "0" ]]; then
    echo -e "\e[0;31m"Error: You must be root to run this install script."\e[0m"
    exit 1
fi

if [ ! -x /usr/bin/dialog ]; then
    apt-get install dialog -y
fi

mode_menu() {
    TITLE="Installation"
    MENU="Choose an installation mode: "
    OPTIONS=(
        1 "Ocserv & Panel (Dockerized)"
        2 "Ocserv & Panel (Local)"
        3 "Panel          (Local)"
        4 "Exit"
    )
    CHOICE=$(dialog --clear \
        --title "$TITLE" \
        --menu "$MENU" \
        15 60 10 \
        "${OPTIONS[@]}" \
        2>&1 >/dev/tty)
    clear
    case $CHOICE in
    1)
        MODE=1 # "multi_docker"
        ;;
    2)
        MODE=2 # "multi_local"
        ;;
    3)
        MODE=3 # "panel_local"
        ;;
    4)
        clear
        echo "bye ..."
        ;;
    esac
    if [ "${MODE}" == false ]; then
        clear
        exit 1

    else
        MODE=${MODE}
        clear
    fi
}

mode_menu

if [ ! "$(echo "1 2 3" | grep -w $MODE)" ]; then
    clear
    exit 1
fi

if [ "$(echo "2" | grep -w $MODE)" ] && [ "$(grep '^VERSION' /etc/os-release | grep "Focal Fossa" | wc -l)" -eq "0" ]; then
    echo "This script is only stable with Ubuntu 20.04(Focal Fossa)"
    exit 1
fi

if [ "$(echo "1" | grep -w $MODE)" ] && [ ! -x "$(command -v docker)" ]; then
    echo "Docker service is not installed"
    exit 1
fi

# set domain and ip
PRIVATE_IP=$(hostname -I | cut -d' ' -f1)
DOMAIN=$(dialog --nocancel --inputbox "Enter your domain or IP or leave it blank" 10 50 "${PRIVATE_IP}" 2>&1 >/dev/tty)
clear
if [ -z "$DOMAIN" ]; then
    echo get public ip from myip.opendns.com
    echo dig started ...
    IP=$(dig +short myip.opendns.com @resolver1.opendns.com)
    if [ "$?" != "0" ]; then
        IP=${PRIVATE_IP}
    fi
    HOST=$(dialog --nocancel --inputbox "Enter your host domain or ip" 10 50 "${IP}" 2>&1 >/dev/tty)
    if [ -z "$DOMAIN" ]; then
        DOMAIN=${IP}
    fi
    clear
else
    HOST="${DOMAIN}"
fi

# get ocserv vars
if [ "$MODE" != "3" ]; then
    PORT=$(dialog --nocancel --inputbox "Enter your port" 10 50 "20443" 2>&1 >/dev/tty)
    clear
    if [ -z "$PORT" ]; then
        PORT=20443
    fi

    CN=$(dialog --nocancel --inputbox "Enter your company name" 10 50 "End-way-Cisco-VPN" 2>&1 >/dev/tty)
    clear
    if [ -z "$CN" ]; then
        CN=End-way-Cisco-VPN
    fi

    ORG=$(dialog --nocancel --inputbox "Enter your organization name" 10 50 "End-way" 2>&1 >/dev/tty)
    clear
    if [ -z "$ORG" ]; then
        ORG=End-way
    fi

    EXPIRE=$(dialog --nocancel --inputbox "Enter ssl expire days" 10 50 "3650" 2>&1 >/dev/tty)
    clear
    if [ -z "$EXPIRE" ]; then
        EXPIRE=3650
    fi

    OC_NET=$(dialog --nocancel --inputbox "Enter ocserv ipv4-network" 10 50 "172.16.24.0/24" 2>&1 >/dev/tty)
    clear
    if [ -z "$OC_NET" ]; then
        OC_NET='172.16.24.0/24'
    fi

    OCSERV_VARS="CN=${CN} ORG=${ORG} EXPIRE=${EXPIRE} OC_NET=${OC_NET}"
    VARS="DOMAIN=${DOMAIN} HOST=${HOST} PORT=${PORT}"
    export CN=${CN} ORG=${ORG} EXPIRE=${EXPIRE} OC_NET=${OC_NET}
    export DOMAIN=${DOMAIN} HOST=${HOST} PORT=${PORT}
fi

chmod +x ./configs/ocserv.sh
chmod +x ./configs/panel.sh

# multi docker
if [ "$MODE" == "1" ]; then
    COMMAND="DOCKER_SCAN_SUGGEST=false ${OCSERV_VARS} ${VARS}"
    BUILD="${COMMAND} docker compose up --build"
    eval ${BUILD}
fi

# multi local
if [ "$MODE" == "2" ]; then
    ./configs/ocserv.sh
    ./configs/panel.sh
fi

# panel only local
if [ "$MODE" == "3" ]; then
    ./configs/panel.sh
fi
