# Stage 1: Build the Go binary
FROM golang:1.24.2 AS builder
LABEL authors="masoud"

ENV GIN_MODE=release
ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /app

# Only copy what's needed for dependency resolution
COPY service/go.mod service/go.sum ./
RUN go mod download

# Copy the entire service for building
COPY service/ .

# Build statically-linked binary (even with CGO)
RUN go build -v -o ocserv_service main.go

# Stage 2: Final runtime image
FROM ubuntu:24.04

LABEL maintainer="masoud"

# Install only runtime dependencies
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

# Setup directory structure
RUN mkdir -p /app /logs

# Copy entrypoint and server logic
COPY deploy/entrypoint.sh /entrypoint.sh
COPY deploy/server.sh /server.sh

# Copy the compiled binary from builder
COPY --from=builder /app/ocserv_service /ocserv_service

# Secure permissions
RUN chmod +x /entrypoint.sh /server.sh /ocserv_service

EXPOSE 443/tcp 443/udp 8080
VOLUME ["/etc/ocserv", "/logs"]

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/server.sh"]
