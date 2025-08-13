
```bash
sudo docker build -t ocserv:latest . && \ 
sudo docker run -it --rm \
    --cap-add=NET_ADMIN \
    --device /dev/net/tun \
    -v $(pwd)/volumes/ocserv:/etc/ocserv \
    -v $(pwd)/volumes/logs:/logs \
    -p 443:443/tcp \
    -p 443:443/udp \
    -p 8081:8081 \
    --name ocserv \
    ocserv


curl --location 'http://127.0.0.1:8080/api/users/?user=null' \
     --header 'Content-Type: application/json' \
     --data '{
        "username": "masoud",
        "password": "1234",
        "group": "defaults"
     }'


ocpasswd -c /etc/ocserv/ocpassd masoud

sudo openconnect -u masoud 127.0.0.1
```


