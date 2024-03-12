# ocserv-users-management

### Web panel to manage ocserv and openconnect users

#### Requirements(Ubuntu 20.04 or Docker host)

## Features:

1- create an account with a limit of gigabytes or monthly usage

2- users: add, edit, update, remove, block and disconnect

3- group: add, edit, update and remove

4- occtl command tools

5- statistics

6- Calculation of users' rx and tx


## Installation :
Choose Your Installation Method:

1- Use install.sh script
```bash
>>> chmod +x install.sh

>>> ./install.sh
```

2- Installing panel without script
```bash
>>> chmod +x ./configs/panel.sh

>>> HOST=http://YOUR_DOMAIN_OR_IP ./configs/panel.sh
```

3- Docker host
```bash
>>> touch prod.env

>>>  cat << EOF >> prod.env
ORG=End-way
EXPIRE=3650
CN=End-way-Cisco-VPN
OC_NET=172.16.24.0/24

# change it to your ip or domain
CORS_ALLOWED=http://HOST_IP_OR_DOMAIN,https://HOST_IP_OR_DOMAIN

# change it to your ip or domain
HOST=HOST_IP_OR_DOMAIN
DOMAIN=
PORT=20443
EOF

>>> DOCKER_SCAN_SUGGEST=false docker-compose up -d --build

```

4- frontend developing
```bash
>>> docker compose -f docker-compose.dev.yml up --build
```

## create extra admin user in terminal
*docker mode 
-- in container
```bash
python3 /app/manage.py createadmin -u USERNAME -p PASSWORD 
```

*systemd
```bash
/var/www/site/back-end/venv/bin/python3 manage.py createadmin -u USERNAME -p PASSWORD 
```


## Admin panel configuration:

1- Launch your web browser.

2- Navigate to http://YOUR-DOMAIN-OR-IP in the address bar.

3- Configure the administrative settings as needed and proceed with the setup.


## Migrate accounts from old panel to new panel:

## commands

1-  --free-traffic: migrate users with free usage traffic

2- --old-path: Path to the old SQLite database


### in os

1- rename /tmp/db.sqlite3 to /tmp/
```bash
>>> mv /tmp/db.sqlite3 /tmp/db-old.sqlite3
```

2- run script to migrate users
```bash
>>> /var/www/site/back-end/venv/bin/python3 manage.py migrate_to_new --old-path /tmp/db-old.sqlite3
```      


### in docker host:

1- rename db.sqlite3 to db-old.sqlite3
```bash
>>> mv db.sqlite3 db-old.sqlite3
```

2- copy db-old.sqlite3 to volumes/db
```bash 
>>> cp db-old.sqlite3 volumes/db
```

3- run command in docker container
```bash
>>> python3 /app/manage.py migrate_to_new --old-path /app/db/db-old.sqlite3
```

## developer mode

1- create dev.env file 
```bash
>>> touch dev.env
```
        
2- copy to dev.env
```bash
>>> cat << EOF >> prod.env
ORG=End-way
EXPIRE=3650
CN=End-way-Cisco-VPN
OC_NET=172.16.24.0/24

# change it to your ip or domain
CORS_ALLOWED=http://127.0.0.1:9000

# change it to your ip or domain
HOST=127.0.0.1
DOMAIN=
PORT=20443
EOF
```
3- run backend service
```bash
>>> docker compose -f docker-compose.dev.yml up -d --build
```
       
4- run frontend service
```bash
>>> cd front-end
>>> npm install && npm run serve
```

5- swagger documents api

Navigate to http://127.0.0.1:8000/doc/ in the address bar.

