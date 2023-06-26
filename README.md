# ocserv-users-management ocserv installation (docker or local) and user management web-panel (Ubuntu 20.04)
```
installation:
    >>> chmod +x installation.sh 
    >>> ./installation.sh   
```
<center><img src="dialog.png"></center>
    
```
install ocserv:
    >>> chmod +x ./configs/ocserv.sh 
    >>> CN=COMPANY_NAME ORG=ORG_NAME EXPIRE=DAYS OC_NET=OCSERV_IPV4_NETWORK ./configs/ocserv.sh
```
```
install user panel:
    >>> chmod +x ./configs/panel.sh 
    >>> HOST=http://YOUR_DOMAIN_OR_IP ./configs/panel.sh
```
```
docker build:
    >>> DOCKER_SCAN_SUGGEST=false \
        DOMAIN=YOUR_DOMAIN \
        HOST=http://YOUR_DOMAIN_OR_IP \
        docker compose up -d --build
```

```
frontend development:
docker compose -f docker-compose.dev.yml up --build
```

---

