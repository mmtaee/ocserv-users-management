FROM golang:1.24.2 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY service/go.mod service/go.sum ./

RUN go mod download

COPY service/ .

RUN go build -v -o ocserv_service main.go

FROM debian:trixie-slim

RUN apt update && \
    apt install -y --no-install-recommends \
    ocserv \
    ca-certificates \
    procps \
    gnutls-bin \
    iptables \
    openssl \
    less \
    dnsutils \
    jq \
    && apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir -p /app /logs

COPY deploy/entrypoint.sh /entrypoint.sh

COPY deploy/server.sh /server.sh

COPY --from=builder /app/ocserv_service /ocserv_service

RUN chmod +x /entrypoint.sh /server.sh /ocserv_service

EXPOSE 443/tcp 443/udp 8080

VOLUME ["/etc/ocserv"]

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/server.sh"]
