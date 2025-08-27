#!/bin/bash
set -e

# ===============================
# Configuration: Services
# ===============================
declare -A SERVICES=(
    ["api"]="./api/main.go"
    ["stream_log"]="./stream_log/main.go"
)

# ===============================
# Function: build_binary
# Parameters:
#   $1 - service name
#   $2 - path to main.go
#   $3 - GOOS
#   $4 - GOARCH
# ===============================
build_binary() {
    local service="$1"
    local main_file="$2"
    local goos="$3"
    local goarch="$4"

    # Output directory for this service

    # shellcheck disable=SC2155
    local bin_dir="$(dirname "$main_file")/bin"
    mkdir -p "${bin_dir}"

    # Output binary path
    local output="${bin_dir}/${service}-${goarch}"

    echo "🛠 Building ${service} for ${goos}/${goarch}..."
    GOOS="${goos}" GOARCH="${goarch}" go build -ldflags="-s -w" -o "${output}" "${main_file}"
    echo "✅ Built: ${output}"
}

# ===============================
# Build all services for multiple architectures
# ===============================
for service in "${!SERVICES[@]}"; do
    main_file="${SERVICES[$service]}"

    build_binary "$service" "$main_file" "linux" "amd64"
    build_binary "$service" "$main_file" "linux" "386"

    # Optional: add more architectures
    # build_binary "$service" "$main_file" "darwin" "amd64"
    # build_binary "$service" "$main_file" "windows" "amd64.exe"
done

echo "🎉 All builds completed successfully!"
