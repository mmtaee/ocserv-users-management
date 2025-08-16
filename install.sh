#!/bin/bash

set -e

# ===============================
# Configuration
# ===============================
GO_REQUIRED=1.24
NODE_REQUIRED=20
OSCERV_DIR_NAME="$(pwd)/ocserv_users_management"
GIT_BRANCH="devel-v3"
BACKEND_DIR="/var/www/backend"
UI_DIR="/var/www/ui"
HOST=
OCSERV_PORT=20443
SSL_CN="End-way-Cisco-VPN"
SSL_ORG="End-way"
SSL_EXPIRE=3650
SSL_OC_NET="172.16.24.0/24"
OCSERV_DNS="8.8.8.8"
ETH=

# ===============================
# Colorful echo module
# ===============================
print_message() {
    local type="$1"
    local message="$2"

    local RED="\e[31m"
    local GREEN="\e[32m"
    local YELLOW="\e[33m"
    local BLUE="\e[34m"
    local MAGENTA="\e[35m"
    local RESET="\e[0m"

    case "$type" in
        info)
            echo -e "${BLUE}[INFO]${RESET} $message"
            ;;
        success)
            echo -e "${GREEN}[SUCCESS]${RESET} $message"
            ;;
        warn)
            echo -e "${YELLOW}[WARN]${RESET} $message"
            ;;
        error)
            echo -e "${RED}[ERROR]${RESET} $message"
            ;;
        highlight)
            echo -e "${MAGENTA}$message${RESET}"
            ;;
        *)
            echo "$message"
            ;;
    esac
}

# ===============================
# Ensure root or sudo access
# ===============================
ensure_root() {
    if [ "$(id -u)" -ne 0 ]; then
        if groups "$USER" | grep -q '\bsudo\b'; then
            print_message info "You are in the sudoers group. Using sudo..."
            sudo -v || { print_message error "Failed to get sudo privileges."; exit 1; }
        else
            print_message error "You are not root and not in the sudoers group."
            print_message info "Please enter the root password to continue."
            su -c "bash $0" || { print_message error "Failed to switch to root."; exit 1; }
            exit 0
        fi
    fi
}

# ===============================
# Check OS, Architecture, Go, Node
# ===============================
check_dependencies() {
    # --- Check OS ---
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_OK=false

        if [[ "$ID" == "ubuntu" ]]; then
            OS_VERSION=$(echo "$VERSION_ID" | awk -F. '{printf "%d.%02d", $1, $2}')
            if (( $(echo "$OS_VERSION >= 24.04" | bc -l) )); then OS_OK=true; fi
        elif [[ "$ID" == "debian" ]]; then
            OS_VERSION=$(echo "$VERSION_ID" | awk -F. '{print $1}')
            if (( OS_VERSION >= 12 )); then OS_OK=true; fi
        fi

        if [ "$OS_OK" = true ]; then
            print_message info "OS detected: $NAME $VERSION_ID"
        else
            print_message error "Unsupported OS version: $NAME $VERSION_ID"
            exit 1
        fi
    else
        print_message error "Cannot detect OS."
        exit 1
    fi

    # --- Check architecture ---
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)
            print_message info "Architecture: 64-bit"
            ;;
        i386|i686)
            print_message info "Architecture: 32-bit"
            ;;
        *)
            print_message warn "Unknown architecture: $ARCH"
            ;;
    esac

    # --- Check Go ---
    GO_INSTALLED=false
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        GO_MAJOR_MINOR=$(echo "$GO_VERSION" | awk -F. '{print $1"."$2}')
        if (( $(echo "$GO_MAJOR_MINOR >= ${GO_REQUIRED}" | bc -l) )); then
            GO_INSTALLED=true
            print_message success "Go version $GO_VERSION is installed and >= ${GO_REQUIRED}"
        else
            print_message warn "Go version $GO_VERSION is installed but < ${GO_REQUIRED}"
        fi
    else
        print_message warn "Go is not installed"
    fi

    # --- Check Node.js ---
    NODE_INSTALLED=false
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node -v | sed 's/v//')
        NODE_MAJOR=$(echo "$NODE_VERSION" | awk -F. '{print $1}')
        if (( NODE_MAJOR >= NODE_REQUIRED )); then
            NODE_INSTALLED=true
            print_message success "Node.js version $NODE_VERSION is installed and >= ${NODE_REQUIRED}"
        else
            print_message warn "Node.js version $NODE_VERSION is installed but < ${NODE_REQUIRED}"
        fi
    else
        print_message warn "Node.js is not installed"
    fi
}


# ===============================
# Get System Environment From Client
# ===============================
get_environment() {
    get_ip_count=0
    print_message info "🌟 Getting required environment..."

    # ===============================
    # Recursive function to get VPS host or detect public IP
    # ===============================
    get_ip() {
        ((get_ip_count = get_ip_count + 1))

        # Prompt user for host
        read -rp "🌐 Enter your VPS host or IP (leave empty to auto-detect public IP): " HOST

        if [ -z "$HOST" ]; then
            print_message info "🔍 No host provided. Attempting to detect public IP..."
            HOST=$(curl -s --max-time 5 https://api.ipify.org)

            if [ -z "$HOST" ]; then
                print_message warn "⚠️ Failed to detect public IP automatically."

                # Check retry limit (max 3 attempts)
                if (( get_ip_count >= 3 )); then
                    print_message error "❌ Failed to detect public IP after 3 attempts. Exiting."
                    exit 1
                fi

                # Recursive call to prompt user again
                get_ip
            else
                print_message highlight "✅ Detected public IP: $HOST"
            fi
        else
            print_message highlight "✅ Using provided host: $HOST"
        fi
    }

    printf "\n"
    get_ip


    # ===============================
    # Get Interface of system
    # ===============================
    GetInterface() {
        printf "\n"

        # Get all interfaces except lo
        interface_list=$(ip -o link show | awk -F': ' '{print $2}' | grep -v lo)

        if [[ -z "$interface_list" ]]; then
            print_message error "No network interfaces found!"
            return 1
        fi

        # Number the interfaces
        numbered_interfaces=$(echo "$interface_list" | awk '{print NR, $0}')

        print_message highlight "Numbered interfaces:"
        IFS=$'\n'
        for line in $numbered_interfaces; do
            number=$(echo "$line" | awk '{print $1}')
            interface=$(echo "$line" | awk '{$1=""; print $0}' | sed 's/^[[:space:]]*//')
           print_message highlight "$(printf "%4d: %s" "$number" "$interface")"
        done
        unset IFS

        read -p "Enter the number corresponding to the desired network interface: " interface_number

        total_interfaces=$(echo "$numbered_interfaces" | wc -l)
        if [[ "$interface_number" =~ ^[0-9]+$ ]] && [[ "$interface_number" -ge 1 ]] && [[ "$interface_number" -le "$total_interfaces" ]]; then
            ETH=$(echo "$numbered_interfaces" | grep "^$interface_number " | awk '{$1=""; print $0}' | sed 's/^[[:space:]]*//')
        print_message highlight "Selected interface: $ETH"
        else
            print_message warn "Invalid selection. Please enter a valid number. ${interface_number} is out of range!"
            GetInterface
        fi
    }

    GetInterface
    printf "\n"

    # ocserv port
    read -rp "Enter your ocserv port or leave blank to use ${OCSERV_PORT}: " port
    [[ -n "$port" ]] && OCSERV_PORT=$port
    print_message highlight "✅ Using port: ${OCSERV_PORT}"
    printf "\n"

    # Company Name
    read -rp "Enter your company name or leave blank to use '${SSL_CN}': " cn
    [[ -n "$cn" ]] && SSL_CN=$cn
    print_message highlight "✅ Using company name: ${SSL_CN}"
    printf "\n"

    # Organization Name
    read -rp "Enter your organization name or leave blank to use '${SSL_ORG}': " org
    [[ -n "$org" ]] && SSL_ORG=$org
    print_message highlight "✅ Using organization name: ${SSL_ORG}"
    printf "\n"

    # SSL Expiration Days
    read -rp "Enter SSL expire days or leave blank to use ${SSL_EXPIRE} days: " expire
    [[ -n "$expire" ]] && SSL_EXPIRE=$expire
    print_message highlight "✅ Using SSL expiration days: ${SSL_EXPIRE}"
    printf "\n"

    # ocserv IPv4 Network
    read -rp "Enter ocserv IPv4 network or leave blank to use ${SSL_OC_NET}: " oc_net
    [[ -n "$oc_net" ]] && SSL_OC_NET=$oc_net
    print_message highlight "✅ Using ocserv IPv4 network: ${SSL_OC_NET}"
    printf "\n"

    # ocserv DNS
    read -rp "🌐 Enter your DNS server or leave blank to use default (${OCSERV_DNS}): " dns
    if [[ -n "$dns" ]]; then
        OCSERV_DNS="$dns"
    fi
    print_message highlight "✅ Using ocserv DNS: ${OCSERV_DNS}"
    printf "\n"
}

# ===============================
# Create env file
# ===============================
set_environment() {
    DATABASES="ocserv.db"
    SECRET_KEY=$(openssl rand -hex 32)
    JWT_SECRET=$(openssl rand -hex 32)
    ALLOW_ORIGINS="http://${HOST},https://${HOST}"

    ENV_FILE="${BACKEND_DIR}/.env"

    if [ ! -f "$ENV_FILE" ]; then
        print_message info "Creating environment file at $ENV_FILE ..."
        cat > "$ENV_FILE" <<EOL
HOST=$HOST
DATABASES=$DATABASES
SECRET_KEY=$SECRET_KEY
JWT_SECRET=$JWT_SECRET
ALLOW_ORIGINS=$ALLOW_ORIGINS
EOL
        print_message success "✅ Environment file created successfully."
    else
        print_message warning "⚠️ Environment file already exists at $ENV_FILE. Skipping creation."
    fi
}


# ===============================
# Install Go
# ===============================
install_go() {
    print_message info "Installing Go $GO_REQUIRED..."
    wget https://go.dev/dl/go${GO_REQUIRED}.linux-amd64.tar.gz -O /tmp/go.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
    go version
    print_message success "Go installed successfully."
}

# ===============================
# Install Node.js
# ===============================
install_node() {
    print_message info "Installing Node.js $NODE_REQUIRED..."
    curl -fsSL https://deb.nodesource.com/setup_${NODE_REQUIRED}.x | bash -
    apt-get install -y nodejs
    print_message success "Node.js installed successfully."
}

# ===============================
# Set Go environment
# ===============================
go_env() {
    case "$ID" in
        ubuntu|debian) GOOS=linux ;;
        darwin) GOOS=darwin ;;
        windows) GOOS=windows ;;
        *) print_message error "Unsupported OS for Go build: $ID"; exit 1 ;;
    esac

    case "$ARCH" in
        x86_64) GOARCH=amd64 ;;
        i386|i686) GOARCH=386 ;;
        arm64|aarch64) GOARCH=arm64 ;;
        *) print_message error "Unsupported architecture for Go build: $ARCH"; exit 1 ;;
    esac

    export GIN_MODE=release
    export CGO_ENABLED=1
    export GOOS GOARCH

    print_message highlight "📦 Go build environment:"
    print_message highlight "   GOOS=$GOOS"
    print_message highlight "   GOARCH=$GOARCH"
    print_message highlight "   GIN_MODE=$GIN_MODE"
    print_message highlight "   CGO_ENABLED=$CGO_ENABLED"
}

# ===============================
# Clone or update repository
# ===============================
update_repo() {
    if [ -d "${OSCERV_DIR_NAME}" ]; then
        print_message info "${OSCERV_DIR_NAME} exists. Updating..."
        cd "${OSCERV_DIR_NAME}" || { print_message error "Failed to enter ${OSCERV_DIR_NAME}"; exit 1; }
        git pull origin "${GIT_BRANCH}" || { print_message error "Git pull failed!"; rm -rf "${OSCERV_DIR_NAME}"; exit 1; }
    else
        print_message info "Cloning ${GIT_BRANCH} branch into ${OSCERV_DIR_NAME}..."
        git clone --single-branch -b "${GIT_BRANCH}" \
            https://github.com/mmtaee/ocserv-users-management.git "${OSCERV_DIR_NAME}" || { print_message error "Git clone failed!"; exit 1; }
        cd "${OSCERV_DIR_NAME}" || { print_message error "Failed to enter ${OSCERV_DIR_NAME}"; exit 1; }
    fi
}

# ===============================
# Build Go service
# ===============================
build_go_service() {
    local dir="$1"
    local output="$2"

    cd "$dir" || { print_message error "Failed to enter $dir"; exit 1; }

    print_message info "⬇️  Downloading Go module dependencies for $(basename "$output")..."
    go mod download || { print_message error "Failed to download Go modules for $(basename "$output")"; exit 1; }
    print_message success "Go modules for $(basename "$output") downloaded successfully."

    print_message info "Building $(basename "$output")..."
    go build -v -o "$output" main.go || { print_message error "Go build failed for $(basename "$output")"; exit 1; }
    print_message success "$(basename "$output") built successfully."
}

# ===============================
# Build UI
# ===============================
build_ui() {
    cd "${OSCERV_DIR_NAME}/ui" || { print_message error "Failed to enter UI directory"; exit 1; }
    print_message info "Installing UI dependencies..."
    npm install || { print_message error "npm install failed"; exit 1; }

    print_message info "Building UI..."
    npm run build || { print_message error "UI build failed"; exit 1; }

    print_message info "Copying UI files to ${UI_DIR}..."
    cp -r "${OSCERV_DIR_NAME}/ui/dist/"* "${UI_DIR}/" || { print_message error "Failed to copy UI files"; exit 1; }
    print_message success "UI built and copied successfully"
}

# ===============================
# Install Ocserv Service
# ===============================
install_ocserv() {
    print_message info "⬇️ Installing ocserv and required packages..."

    if apt install -y --no-install-recommends \
      ocserv \
      ca-certificates \
      procps \
      gnutls-bin \
      iptables \
      openssl \
      less \
      dnsutils \
      jq; then
        print_message success "✅ ocserv and dependencies installed successfully."
    else
        print_message error "❌ Failed to install ocserv and dependencies."
        exit 1
    fi
}

# ===============================
# Update IP Tables Configs
# ===============================
UpdateIpTables() {
    # Ensure interface is set
    if [ -z "$ETH" ]; then
        GetInterface || { echo "Failed to get interface"; return 1; }
    fi

    # Example iptables rules (adjust as needed)
    iptables -F
    iptables -A INPUT -i lo -j ACCEPT
    iptables -A INPUT -i "$ETH" -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
    iptables -A INPUT -i "$ETH" -p tcp --dport 22 -j ACCEPT
    iptables -A INPUT -j DROP

    # Save rules depending on system
    if command -v netfilter-persistent >/dev/null 2>&1; then
        # Debian 12 / Ubuntu 24+
        netfilter-persistent save
    elif [ -f /etc/iptables/rules.v4 ]; then
        # Older systems
        iptables-save > /etc/iptables/rules.v4
    else
        echo "Warning: Cannot find persistent iptables location. Rules not saved."
    fi
}


# ===============================
# Install nginx
# ===============================
install_nginx() {
    if apt install -y nginx; then
        print_message success "✅ nginx installed successfully."
    else
        print_message error "❌ Failed to install nginx."
        exit 1
    fi
}

# ===============================
# Systemd Services Configurations
# ===============================
systemd_service (){
  print_message info "Configuration system service ..."
  # TODO: copy systemd services
  # 1- panel_api
  # 2- ocserv_service
  # 3- log_service
  # 4- ui
}

# ===============================
# Check Systemd Services
# ===============================
check_systemd (){
  print_message info "Checking Systemd Services ..."
  # TODO: check systemd service
  # 1- panel_api
  # 2- ocserv_service
  # 3- log_service
  # 4- ui
  # 5- ocserv service
}

# ===============================
# Main
# ===============================
main() {
    ensure_root
    check_dependencies

    # Create directories
    mkdir -p "${BACKEND_DIR}" || { print_message error "Failed to create ${BACKEND_DIR}"; exit 1; }
    mkdir -p "${UI_DIR}" || { print_message error "Failed to create ${UI_DIR}"; exit 1; }

    get_environment

    set_environment

    # Install missing dependencies
    [ "$GO_INSTALLED" = false ] && install_go
    [ "$NODE_INSTALLED" = false ] && install_node

    go_env

    # Update repository
    update_repo

    # Build backend services
    build_go_service "${OSCERV_DIR_NAME}/backend" "${BACKEND_DIR}/panel_api"
    build_go_service "${OSCERV_DIR_NAME}/ocserv/service" "${BACKEND_DIR}/ocserv_service"
    build_go_service "${OSCERV_DIR_NAME}/services" "${BACKEND_DIR}/log_service"

    # Build UI
    build_ui

    print_message success "🎉 All services built successfully!"

    # Install Ocserv
    install_ocserv

    # Update IP Tables
    UpdateIpTables

    # Install Nginx
    install_nginx

    # Systemd Services
    systemd_service

    # Check Systemd Services
    check_systemd
}

# Run main
main

exit 0
