#!/bin/bash

API_URL="http://localhost:8080/api/ocserv/users"
AUTH_TOKEN="eyJhbGciOi..."

for i in $(seq 1 100); do
  username="user$(openssl rand -hex 3)"
  echo "Creating user: $username"

  response=$(curl -s -w "%{http_code}" -o /dev/null -X POST "$API_URL" \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -H "Content-Type: application/json" \
    -d "$(jq -n --arg user "$username" \
          '{group:"defaults", password:"1234", traffic_size:0, traffic_type:"Free", username:$user, config:{}}')")

  if [ "$response" == "201" ]; then
    echo "$username"
  else
    echo "Failed to create $username"
  fi
done

