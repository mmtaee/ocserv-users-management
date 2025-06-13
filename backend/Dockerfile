# Stage 1: Build the Go binary
FROM golang:1.24.2 AS builder
LABEL authors="masoud"

# Set environment variables for release build with CGO enabled for SQLite
ENV GIN_MODE=release
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary output named "panel_api"
RUN go build -o panel_api main.go


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
COPY --from=builder /app/panel_api .

# Expose the port your app listens on (change if needed)
EXPOSE 8080

# Run the binary
CMD ["./panel_api"]
