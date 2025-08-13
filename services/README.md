
```bash
go build main.go && sudo --preserve-env=DOCKERIZED DOCKERIZED=true ./main

# Docker
sudo docker build -t ocserv:service .

sudo docker run -it --rm -p "8082:8082" --env-file .env \
     -v "/home/masoud/Projects/ocserv-users-management/backend/ocserv.db:/app/db/ocserv.db" \
     -v "/var/run/docker.sock:/var/run/docker.sock" ocserv:service

```




