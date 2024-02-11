#!/bin/bash
OCSERV_LOG_FILE=/var/log/ocserv.log
BACKEND=0.0.0.0:8000

pidfile=/run/ocserv.pid

if [ ! -f "${OCSERV_LOG_FILE}" ]; then
    touch "${OCSERV_LOG_FILE}"
    chmod +x "${OCSERV_LOG_FILE}"
fi

# ocserv service
printf "\e[33m########### ocserv service starting ... ###########\e[0m"
printf "\n"
cd ${LOG_FILE}
/usr/sbin/ocserv --debug=2 --foreground --pid-file=${pidfile} --config=/etc/ocserv/ocserv.conf >${OCSERV_LOG_FILE} 2>&1 &

# django service
printf "\e[33m########### backend service starting ... ###########\e[0m"
printf "\n"
OCSERV_LOG_FILE=${OCSERV_LOG_FILE} DOCKERIZED=True python3 /app/manage.py runserver ${BACKEND} &

# user stats service
printf "\e[33m########### user stats service starting ... ###########\e[0m"
printf "\n"
#OCSERV_LOG_FILE=${OCSERV_LOG_FILE} python3 /app/user_stats.py &
OCSERV_LOG_FILE=${OCSERV_LOG_FILE} python3 /app/manage.py user_stats &

wait -n
exit $?
