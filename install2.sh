#!/bin/bash

color_echo() {
    local color=$1
    local text=$2
    echo -e "\e[${color}m${text}\e[0m"
}

source scripts.sh

# GetOcservVars

# GetPort

# GetDomain

# GetIP

# GetModes

# UpdateIpTables

# if [[ $mode == '1' ]]; then
#     GetInterface
# elif [[ $mode == '2' ]]; then
#     GetInterface
# elif [[ $mode == '3' ]]; then
#     GetInterface
# elif [[ $mode == '4' ]]; then
#     echo "4 docker"
# else
#     echo "5 docker"
# fi



# TODO: update script to remove ip tables and add to docker file entrypoint