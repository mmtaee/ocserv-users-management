FROM golang:1.25.0 AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY  . .

RUN go build -v -o stream_log main.go

FROM debian:trixie-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/stream_log .

VOLUME ["/app/db"]

EXPOSE 8080

CMD ["./stream_log"]
