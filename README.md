
# OpenConnect VPN Server (ocserv) with Web Panel

This guide provides a simple and efficient way to set up and manage an OpenConnect VPN server (ocserv) with a powerful web panel for managing users and user groups. This solution offers an easy-to-deploy, scalable, and secure VPN setup with minimal configuration.

---
<br />

<p>
<img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/mmtaee/ocserv-users-management"><img alt="GitHub contributors" src="https://img.shields.io/github/contributors/mmtaee/ocserv-users-management"> 
</p>
 

## Key Features

1. **User Management**:
   - Create, update, edit, remove, block, and disconnect users.
   - Set traffic usage limits (e.g., GB or monthly usage).
  
2. **Group Management**:
   - Create, update, and remove user groups.
  
3. **Command Line Tools**:
   - `occtl` command-line utility for various server operations.
  
4. **Statistics**:
   - View statistics on user traffic (RX and TX).
  
5. **Usage Calculation**:
   - Track data usage per user.

---

## Installation Methods

You can install the solution using one of the following methods:

### 1. **Using the `install.sh` Script**

```bash
chmod +x install.sh
./install.sh
```

### 2. **Installing the Panel Without Script**

```bash
chmod +x ./configs/panel.sh
HOST=http://YOUR_DOMAIN_OR_IP ./configs/panel.sh
```

### 3. **Docker Host Setup**

1. Create the `prod.env` file:

```bash
touch prod.env
```

2. Add the following configuration to the `prod.env` file:

```bash
cat << EOF >> prod.env
ORG=End-way
EXPIRE=3650
CN=End-way-Cisco-VPN
OC_NET=172.16.24.0/24

# Replace with your domain or IP
CORS_ALLOWED=http://HOST_IP_OR_DOMAIN,https://HOST_IP_OR_DOMAIN
HOST=HOST_IP_OR_DOMAIN
DOMAIN=
PORT=20443
EOF
```

3. Run the Docker Compose command:

```bash
docker compose up -d --build
```

### 4. **Frontend Development Mode**

```bash
docker compose -f docker-compose.dev.yml up --build
```

---

## Creating an Admin User

- **Docker Mode (In Container)**:

```bash
python3 /app/manage.py createadmin -u USERNAME -p PASSWORD
```

- **System Mode**:

```bash
/var/www/site/back-end/venv/bin/python3 /var/www/site/back-end/manage.py createadmin -u USERNAME -p PASSWORD
```

---

## Admin Panel Configuration

1. **Launch Web Browser**.
2. **Navigate to** `http://YOUR-DOMAIN-OR-IP` in the browser.
3. **Complete the administrative setup**.

---

## Migrating Accounts from the Old Panel to New Panel

### Migration Commands

1. **For users with free traffic**:
   
   ```bash
   --free-traffic
   ```

2. **Path to the old SQLite database**:
   - **For OS**:

   ```bash
   mv /tmp/db.sqlite3 /tmp/db-old.sqlite3
   /var/www/site/back-end/venv/bin/python3 manage.py migrate_to_new --old-path /tmp/db-old.sqlite3
   ```

   - **For Docker Host**:
   
   ```bash
   mv db.sqlite3 db-old.sqlite3
   cp db-old.sqlite3 volumes/db
   python3 /app/manage.py migrate_to_new --old-path /app/db/db-old.sqlite3
   ```

---

## Developer Mode

1. **Create a `dev.env` file**:

```bash
touch dev.env
```

2. **Add the following configuration to `dev.env`**:

```bash
cat << EOF >> dev.env
DEBUG=True
ORG=End-way
EXPIRE=3650
CN=End-way-Cisco-VPN
OC_NET=172.16.24.0/24

# Change to your domain or IP
CORS_ALLOWED=http://127.0.0.1:9000

# Change to your domain or IP
HOST=127.0.0.1
DOMAIN=
PORT=20443
EOF
```

3. **Run the Backend Service**:

```bash
docker compose -f docker-compose.dev.yml up -d --build
```

4. **Run the Frontend Service**:

```bash
cd front-end
npm install && npm run serve
```

5. **Swagger API Documentation**:

Navigate to `http://127.0.0.1:8000/doc/` to access the Swagger documentation.

---

## Additional Notes

- The OpenConnect VPN server (ocserv) is configured with best practices for security.
- The web panel is designed to be easy to use for both admins and end users.
- If you encounter any issues, please refer to the documentation or contact support.

---

By following the above steps, you can easily set up and manage your OpenConnect VPN server and provide users with secure, scalable VPN access.


## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=mmtaee/ocserv-users-management&type=Date)](https://www.star-history.com/#mmtaee/ocserv-users-management&Date)
 