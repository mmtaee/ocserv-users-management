#!/bin/bash

LOG_FILE=/var/log/ocserv.log

pidfile=/run/ocserv.pid
# ocserv service
printf "\e[33m########### ocserv service starting ... ###########\e[0m"
printf "\n"
/usr/sbin/ocserv --debug=2 --foreground --config=/etc/ocserv/ocserv.conf --pid-file=${pidfile} 2>&1 | tee ${LOG_FILE} &
# /usr/sbin/ocserv -c /etc/ocserv/ocserv.conf -f 2>&1 > /tmp/log.txt &

# django service
printf "\e[33m########### backend service starting ... ###########\e[0m"
printf "\n"
python3 /app/manage.py migrate
python3 /app/manage.py runserver 0.0.0.0:8000 &

# user stats service
printf "\e[33m########### user stats service starting ... ###########\e[0m"
printf "\n"
LOG_FILE=${LOG_FILE} python3 /app/user_stats.py &

wait -n
exit $?
