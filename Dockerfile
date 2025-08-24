FROM golang:1.25.0 AS builder

ENV GIN_MODE=release
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN mkdir /common

COPY ./common /common

WORKDIR /app

COPY ./api/go.mod ./api/go.sum ./

RUN go mod download

COPY ./api .

RUN go build -o api main.go

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

RUN mkdir /app /logs

COPY --from=builder /app/api /api

COPY scripts/entrypoint.sh /entrypoint.sh

COPY scripts/server.sh /server.sh

RUN chmod +x /entrypoint.sh /server.sh /api

EXPOSE 443/tcp 443/udp 8080

VOLUME ["/etc/ocserv", "/app"]

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/server.sh"]
