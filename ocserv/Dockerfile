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

RUN chmod +x /entrypoint.sh /server.sh

EXPOSE 443/tcp 443/udp 8080

VOLUME ["/etc/ocserv", "/app"]

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/server.sh"]
