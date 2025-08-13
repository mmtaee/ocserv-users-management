FROM golang:1.24.2 AS builder
LABEL authors="masoud"

ENV GIN_MODE=release
ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY  . .

RUN go build -v -o ocserv_service main.go

# Stage 2: Final runtime image (Debian-slim)
FROM debian:bookworm-slim

LABEL authors="masoud"

# Install SQLite runtime libraries and certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/ocserv_service .

# Expose the port your app listens on (change if needed)
EXPOSE 8082

# Run the binary
CMD ["./ocserv_service"]
