
```bash

sudo docker network create ocserv

sudo docker build -t ocserv:latest . 

sudo docker run -it --rm \
    --privileged \
    --env-file ./.env \
    --cap-add=NET_ADMIN \
    --device /dev/net/tun \
    -v $(pwd)/volumes/ocserv:/etc/ocserv \
    -v $(pwd)/volumes/db:/usr/local/bin/db \
    -v $(pwd)/volumes/db_tmp:/root/ocserv_db \
    -p 443:443/tcp \
    -p 443:443/udp \
    -p 8080:8080 \
    --network ocserv \
    --name ocserv \
    ocserv


ocpasswd -c /etc/ocserv/ocpassd masoud

sudo openconnect -u masoud 127.0.0.1



sudo docker build -t stream_log:latest -f Dockerfile-Stream-Log .

sudo docker run -it --rm \
    --env-file ./.env \
    -v $(pwd)/volumes/db:/usr/local/bin/db \
    -v $(pwd)/volumes/db_tmp:/root/ocserv_db \
    -v "/var/run/docker.sock:/var/run/docker.sock" \
    -p 8081:8080 \
    --name stream_log \
    --network ocserv \
    stream_log


source .env  && docker build --build-arg LANGUAGES="${LANGUAGES}" -t web:latest -f Dockerfile-Web .
  
docker run -it --rm \
  --name web \
  -v $(pwd)/volumes/nginx/certs:/etc/nginx/certs \
  -p 3000:3000 \
  --network ocserv \
  web:latest



```


