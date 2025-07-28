#!/bin/bash

SITE_DIR="/var/www/site"
CURRENT_DIR=$(pwd)

if [[ $(id -u) != "0" ]]; then
    echo -e "\e[0;31m"Error: You must be root to run this install script."\e[0m"
    exit 1
fi
apt install -y python3 python3-pip python3-venv python3-dev build-essential \
    nginx cron curl gcc g++ make openssl nodejs ca-certificates curl gnupg
if [ "$?" = "0" ]; then
    echo -e "\e[0;32m"Panel dependencies installation was successful."\e[0m"
else
    echo -e "\e[0;31m"Panel dependencies installation was failed."\e[0m"
    exit 1
fi

# back-end
echo -e "\e[0;32m"Back-end Installing ..."\e[0m"
rm -rf /var/www/html
#rm -rf ${SITE_DIR}
mkdir -p ${SITE_DIR}
cp -r ${CURRENT_DIR}/back-end ${SITE_DIR}/back-end
rm -rf /lib/systemd/system/backend.service
rm -rf /lib/systemd/system/user_stats.service
cp ./configs/backend.service /lib/systemd/system
cp ./configs/user_stats.service /lib/systemd/system
cp ./configs/uwsgi.ini ${SITE_DIR}/back-end
python3 -m venv ${SITE_DIR}/back-end/venv
source ${SITE_DIR}/back-end/venv/bin/activate
pip install -U wheel setuptools
pip install -r ${SITE_DIR}/back-end/requirements.txt
pip install uwsgi==2.0.24
SECRET_KEY=$(openssl rand -base64 '64')
echo "DEBUG=False" >${SITE_DIR}/back-end/.env
echo "SECRET_KEY=${SECRET_KEY}" >>${SITE_DIR}/back-end/.env
echo "CORS_ALLOWED=http://${HOST},https://${HOST}" >>${SITE_DIR}/back-end/.env
mkdir -p ${SITE_DIR}/back-end/db
${SITE_DIR}/back-end/manage.py migrate
deactivate

echo 'www-data ALL=(ALL) NOPASSWD: \
    /usr/bin/rm /etc/ocserv/*, \
    /usr/bin/mkdir /etc/ocserv/*, \
    /usr/bin/touch /etc/ocserv/*, \
    /usr/bin/cat /etc/ocserv/*, \
    /usr/bin/sed /etc/ocserv/*, \
    /usr/bin/tee /etc/ocserv/*, \
    /usr/bin/ocpasswd, \
    /usr/bin/occtl, \
    /usr/bin/systemctl restart ocserv.service, \
    /usr/bin/systemctl status ocserv.service' | sudo tee -a /etc/sudoers >/dev/null

crontab -l | echo "59 23 * * * ${SITE_DIR}/back-end/venv/bin/python3 ${SITE_DIR}/back-end/manage.py user_management" | crontab -

# front-end
echo -e "\e[0;32m"Front-End Installing ..."\e[0m"

CURRENT_VERSION=$(node -v 2>/dev/null | sed 's/^v//')
REQUIRED_VERSION="20.0.0"

version_lt() {
  [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" != "$2" ]
}

if [ -z "$CURRENT_VERSION" ] || version_lt "$CURRENT_VERSION" "$REQUIRED_VERSION"; then
  echo -e "\e[0;33mInstalling or upgrading Node.js to version 20.x...\e[0m"
  curl -fsSL https://deb.nodesource.com/setup_20.x | sudo bash -
  sudo apt-get install -y nodejs
else
  echo -e "\e[0;32mNode.js is installed: v$CURRENT_VERSION (>= 20.0.0)\e[0m"
fi

echo -e "\e[0;32m"Buliding web production..."\e[0m"

cd "${CURRENT_DIR}"/front-end/ || exit
npm install
NODE_ENV=production npm run build

mkdir -p ${SITE_DIR}/front-end
cp -r "${CURRENT_DIR}"/front-end/dist/* ${SITE_DIR}/front-end

# nginx
echo -e "\e[0;32m"Nginx Configurations ..."\e[0m"
rm -rf /etc/nginx/sites-enabled/default

if [[  -n "${DOMAIN}"  ]];then
cat <<\EOT >/etc/nginx/conf.d/site.conf
server {
    listen 80;
    server_name ${DOMAIN} ;
    server_tokens off;

    return 302 https://$server_name$request_uri;
}
server {
    listen 443 ssl http2;
    server_name ${DOMAIN} ;
    server_tokens off;

    ssl_certificate         /etc/nginx/certs/cert.pem;
    ssl_certificate_key    /etc/nginx/certs/cert.key;
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
}
EOT
else
cat <<\EOT >/etc/nginx/conf.d/site.conf
server {
    listen 80;
    server_tokens off;
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
}
EOT
fi



chown -R www-data. /etc/nginx/conf.d/site.conf
chown -R www-data. ${SITE_DIR}

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
OCSERV_STAT=$(systemctl is-active user_stats.service)
if [ "$OCSERV_STAT" = "active" ]; then
    echo -e "\e[0;32m"user_stats.service Is Started."\e[0m"
else
    echo -e "\e[0;31m"user_stats.service Is Not Running."\e[0m"
    exit 1
fi
