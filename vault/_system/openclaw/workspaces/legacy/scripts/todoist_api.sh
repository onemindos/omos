#!/bin/bash
# Todoist REST API wrapper
# Usage: todoist_api.sh <endpoint> <method> [data_json]
ENDPOINT=$1
METHOD=$2
DATA=$3
TOKEN="${TODOIST_API_KEY}"
if [ -z "$DATA" ]; then
  curl -s -X "$METHOD" "https://api.todoist.com/rest/v2/$ENDPOINT" \
    -H "Authorization: Bearer $TOKEN"
else
  curl -s -X "$METHOD" "https://api.todoist.com/rest/v2/$ENDPOINT" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$DATA"
fi
