FROM golang:1.24.2 AS builder

ENV GIN_MODE=release
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o panel_api main.go

FROM debian:trixie-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/panel_api .

EXPOSE 8080

CMD ["./panel_api"]
