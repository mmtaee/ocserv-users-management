#!/bin/bash

if [[ $(id -u) != "0" ]]; then
    echo -e "\e[0;31m"Error: You must be root to run this install script."\e[0m"
    exit 1
fi
apt install -y python3 python3-pip python3-venv python3-dev build-essential \
    nginx cron curl gcc g++ make openssl
if [ "$?" = "0" ]; then
    echo -e "\e[0;32m"Panel dependencies installation was successful."\e[0m"
else
    echo -e "\e[0;31m"Panel dependencies installation was failed."\e[0m"
    exit 1
fi
# back-end
echo -e "\e[0;32m"Back-end Installing ..."\e[0m"
rm -rf /var/www/html
SITE_DIR="/var/www/site"
CURRENT_DIR=$(pwd)
rm -rf ${SITE_DIR}
mkdir -p ${SITE_DIR}
mkdir -p ${SITE_DIR}/back-end
touch /var/log/socket_passwd
chown -R www-data. /var/log/socket_passwd
cp -r ${CURRENT_DIR}/back-end/* ${SITE_DIR}/back-end
cp ./configs/uwsgi.ini ${SITE_DIR}/back-end
python3 -m venv ${SITE_DIR}/back-end/venv
source ${SITE_DIR}/back-end/venv/bin/activate
pip install -U wheel setuptools
pip install -r ${SITE_DIR}/back-end/requirements.txt
pip install uwsgi
SECRET_KEY=$(openssl rand -base64 '64')
echo "DEBUG=False" >${SITE_DIR}/back-end/.env
echo "SECRET_KEY=${SECRET_KEY}" >>${SITE_DIR}/back-end/.env
echo "CORS_ALLOWED=http://${HOST},https://${HOST}" >>${SITE_DIR}/back-end/.env
rm -rf /lib/systemd/system/backend.service
rm -rf /lib/systemd/system/user_stats.service
cp ./configs/backend.service /lib/systemd/system
cp ./configs/user_stats.service /lib/systemd/system
echo www-data ALL=NOPASSWD: /usr/bin/ocpasswd >>/etc/sudoers
echo www-data ALL=NOPASSWD: /usr/bin/occtl >>/etc/sudoers
echo www-data ALL=NOPASSWD: /usr/bin/systemctl restart ocserv.service >>/etc/sudoers
echo www-data ALL=NOPASSWD: /usr/bin/systemctl status ocserv.service >>/etc/sudoers
crontab -l | {
    cat
    echo "59 23 * * * ${SITE_DIR}/back-end/venv/bin/python ${SITE_DIR}/back-end/manage.py user_management"
} | crontab -
${SITE_DIR}/back-end/manage.py migrate
deactivate
# front-end
echo -e "\e[0;32m"Front-End Installing ..."\e[0m"
curl -sL https://deb.nodesource.com/setup_18.x -o /tmp/nodesource_setup.sh
bash /tmp/nodesource_setup.sh
apt install -y nodejs
cd ${CURRENT_DIR}/front-end/
npm install
NODE_ENV=production npm run build
mkdir -p ${SITE_DIR}/front-end
cp -r ${CURRENT_DIR}/front-end/dist/* ${SITE_DIR}/front-end
chown -R www-data ${SITE_DIR}
# nginx
echo -e "\e[0;32m"Nginx Configurations ..."\e[0m"
rm -rf /etc/nginx/sites-enabled/default
cat <<\EOT >/etc/nginx/conf.d/site.conf
server {
    listen 80;
    location / {
        root /var/www/site/front-end;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    location ~ ^/(api) {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Host $host;
    }
    location /ws {
        rewrite ^/ws(.*)$ $1 break;
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
    }
}
EOT

# monitor service
# cp -r ${CURRENT_DIR}/monitor /opt/monitor
# cp /opt/monitor/monitor.service /lib/systemd/system/monitor.service
# systemctl restart monitor.service
# systemctl enable monitor.service

chown -R www-data. /etc/nginx/conf.d/site.conf
systemctl disable backend.service
systemctl disable user_stats.service
systemctl daemon-reload
systemctl restart nginx.service
systemctl enable nginx.service
systemctl restart backend.service
systemctl enable backend.service
systemctl restart user_stats.service
systemctl enable user_stats.service
NGINX_STATE=$(systemctl is-active nginx)
if [ "$NGINX_STATE" = "active" ]; then
    echo -e "\e[0;32m"Nginx Is Started."\e[0m"
else
    echo -e "\e[0;31m"Nginx Is Not Running."\e[0m"
    exit 1
fi
OCSERV_UWSGI_STATE=$(systemctl is-active backend.service)
if [ "$OCSERV_UWSGI_STATE" = "active" ]; then
    echo -e "\e[0;32m"backend.service Is Started."\e[0m"
else
    echo -e "\e[0;31m"backend.service Is Not Running."\e[0m"
    exit 1
fi
OCSERV_USERSTAT=$(systemctl is-active user_stats.service)
if [ "$OCSERV_USERSTAT" = "active" ]; then
    echo -e "\e[0;32m"user_stats.service Is Started."\e[0m"
else
    echo -e "\e[0;31m"user_stats.service Is Not Running."\e[0m"
    exit 1
fi
