FROM golang:1.25.0 AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY  . .

RUN go build -v -o webhook main.go

FROM debian:trixie-slim

WORKDIR /app

COPY --from=builder /app/webhook .

VOLUME ["/app"]

EXPOSE 8080

CMD ["./webhook"]
