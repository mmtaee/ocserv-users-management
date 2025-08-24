
```bash
sudo docker build -t ocserv:latest . && \ 
sudo docker run -it --rm \
    --privileged \
    --env-file .env \
    --cap-add=NET_ADMIN \
    --device /dev/net/tun \
    -v $(pwd)/volumes/ocserv:/etc/ocserv \
    -v $(pwd)/volumes/db:/app/db \
    -v $(pwd)/volumes/db_tmp:/root/ocserv_db \
    -p 443:443/tcp \
    -p 443:443/udp \
    -p 8080:8080 \
    --name ocserv \
    ocserv



ocpasswd -c /etc/ocserv/ocpassd masoud

sudo openconnect -u masoud 127.0.0.1
```


