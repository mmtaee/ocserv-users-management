#!/bin/bash

OCSERV_LOG_FILE=/var/log/ocserv.log
pidfile=/run/ocserv.pid

# ocserv service
printf "\e[33m########### ocserv service starting ... ###########\e[0m"
printf "\n"
/usr/sbin/ocserv --debug=2 --foreground --config=/etc/ocserv/ocserv.conf --pid-file=${pidfile} 2>&1 | tee ${OCSERV_LOG_FILE} &
# /usr/sbin/ocserv -c /etc/ocserv/ocserv.conf -f 2>&1 > /tmp/log.txt &

# django service
printf "\e[33m########### backend service starting ... ###########\e[0m"
printf "\n"
mkdir -p /app/db
# pip install pymemcache
python3 /app/manage.py migrate
OCSERV_LOG_FILE=${OCSERV_LOG_FILE} DOCKERIZED=True python3 /app/manage.py runserver 0.0.0.0:8000 &
# cd /app && uvicorn ocserv.asgi:application --host 0.0.0.0 --port 8000 --workers 1 --ws websockets --lifespan off&
# gunicorn ocserv.asgi:application -w 1 -k uvicorn.workers.UvicornWorker -b 0.0.0.0:8000 --chdir /app


# user stats service
printf "\e[33m########### user stats service starting ... ###########\e[0m"
printf "\n"
OCSERV_LOG_FILE=${OCSERV_LOG_FILE} python3 /app/user_stats.py &


# log monitor
# printf "\e[33m########### log monitor service starting ... ###########\e[0m"
# printf "\n"
# sleep 10
# OCSERV_LOG_FILE=${OCSERV_LOG_FILE} WS_SERVER=ws://127.0.0.1:8000 WS_TOKEN=${WS_TOKEN} /monitor/log_monitor &

wait -n
exit $?


# run only uvicorn:
# uvicorn MAIN_DIR.asgi:application --reload --debug --ws websockets
# uvicorn backend.asgi:application --reload --debug --ws websockets --reload --host 0.0.0.0 --port 8004
# run with gunicorn together:
# gunicorn MAIN_DIR.asgi:application -w 4 -k uvicorn.workers.UvicornWorker
# pip install websockets
# pip install uvicorn
# pip install uvloop
# pip install httptools

