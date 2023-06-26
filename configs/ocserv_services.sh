#!/bin/bash

printf "\e[33m########### ocserv services starting ###########\e[0m"
# ocserv service
pidfile=/run/ocserv.pid
printf "\e[33m########### ocserv service starting ... ###########\e[0m"
printf "\n"
/usr/sbin/ocserv --debug=2 --foreground --config=/etc/ocserv/ocserv.conf --pid-file=${pidfile}
# /usr/sbin/ocserv -c /etc/ocserv/ocserv.conf -f &

wait -n
exit $?
