#!/bin/bash

printf "\e[33m########### nginx starting ###########\e[0m \n"
VALIDATE="^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$"

printf "\e[33m########### ${DOMAIN} ###########\e[0m \n"

if [[ "$DOMAIN" =~ $VALIDATE ]]; then
    echo "Valid $DOMAIN name."
    cat <<\EOT >/tmp/site.conf.template
upstream api_backend {
    server ocserv_and_backend:8000;
}

server {
    listen 80;
    server_name ${DOMAIN} ;
    return 302 https://$server_name$request_uri;
}
server {
    listen 443 ssl http2;
    server_name ${DOMAIN} ;

    ssl_certificate         /etc/nginx/certs/cert.pem;
    ssl_certificate_key    /etc/nginx/certs/cert.key;

    location / {
        root /var/www/site;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    location ~ ^/(api) {
        proxy_pass http://api_backend;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Host $host;
    }
}
EOT
    envsubst '$DOMAIN' </tmp/site.conf.template >/etc/nginx/conf.d/site.conf
else
    cat <<\EOT >/etc/nginx/conf.d/site.conf
upstream api_backend {
    server ocserv_and_backend:8000;
}
server {
    listen 80;
    location / {
        root /var/www/site;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    location ~ ^/(api) {
        proxy_pass http://api_backend;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Host $host;
    }
}
EOT
fi

exec "$@"
