#!/bin/bash

API_URL="http://localhost:8080/api/ocserv/groups"
AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYyMDU0NjIsImlhdCI6MTc1MzQ0MDY2MiwiaXNBZG1pbiI6dHJ1ZSwianRpIjoiMDFLMTBLQUQ3Q1kyQzFSSlFGWUtYNVo2UzUiLCJzdWIiOiIwMUsxMEtBQU03RjlXNFRBTkdDQ0g2SEVEUiJ9.QtYwIqA5Z9fIZ0N_Dyc_LQok5z41PJv5dKq3eAzNBwQ"

for i in {1..100}; do
  GROUP_NAME="test${i}"
  echo "Creating group: $GROUP_NAME"

  curl -s -o /dev/null -w "Status: %{http_code}\n" -X POST "$API_URL" \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"${GROUP_NAME}\", \"config\": {}}"
done
