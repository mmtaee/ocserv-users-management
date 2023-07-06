#!/bin/bash

LOG_FILE=/var/log/ocserv.log
BACKEND=0.0.0.0:8000
WS_TOKEN=1234

pidfile=/run/ocserv.pid
# ocserv service
printf "\e[33m########### ocserv service starting ... ###########\e[0m"
printf "\n"
/usr/sbin/ocserv --debug=2 --foreground --config=/etc/ocserv/ocserv.conf --pid-file=${pidfile} 2>&1 | tee ${LOG_FILE} &
# /usr/sbin/ocserv -c /etc/ocserv/ocserv.conf -f 2>&1 > /tmp/log.txt &

# log monitor
LOG_FILE=${LOG_FILE} WS_SERVER=ws://${BACKEND} WS_TOKEN=${WS_TOKEN} /log_service/log_monitor &

# django service
printf "\e[33m########### backend service starting ... ###########\e[0m"
printf "\n"
python3 /app/manage.py runserver ${BACKEND} &

# user stats service
printf "\e[33m########### user stats service starting ... ###########\e[0m"
printf "\n"
LOG_FILE=${LOG_FILE} python3 /app/user_stats.py &

wait -n
exit $?
