# ================================
# Vue.js build stage
# ================================
FROM node:20 AS vue

# Build-time arguments
ARG API_URL
ARG LOG_SOCKET_URL
ARG LANGUAGES

WORKDIR /app

# Copy package.json first to leverage Docker cache
COPY package*.json ./

RUN npm install --legacy-peer-deps

# Copy all source files
COPY . .

# Build Vue app with environment variables
RUN NODE_ENV=production \
    VITE_API_URL="${API_URL}" \
    VITE_LOG_SOCKET_URL="${LOG_SOCKET_URL}" \
    VITE_I18N_LANGUAGES="${LANGUAGES}" \
    npm run build

# ================================
# Nginx production stage
# ================================
FROM nginx:alpine

# Install bash, gettext (for envsubst), and openssl
RUN apk add --no-cache bash gettext openssl \
    && rm -rf /etc/nginx/conf.d/default.conf /usr/share/nginx/html \
    && mkdir -p /var/www/site /etc/nginx/certs

# Copy custom Nginx startup script
COPY configs/nginx.conf  /etc/nginx/conf.d/site.conf

COPY configs/entrypoint.sh /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh

COPY --from=vue /app/dist /var/www/site

VOLUME ["/etc/nginx/certs"]

EXPOSE 80 443

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]
