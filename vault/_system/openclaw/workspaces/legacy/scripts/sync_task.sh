#!/bin/bash
# Sync agent task state to Todoist
# Usage: sync_task.sh <content> <status> [task_id] [description]
CONTENT=$1
STATUS=$2
TASK_ID=$3
DESCRIPTION=$4
DIR="$(dirname "$0")"
ESC_CONTENT=$(printf '%s' "$CONTENT" | python3 -c "import sys,json; print(json.dumps(sys.stdin.read()))")
ESC_DESC=$(printf '%s' "$DESCRIPTION" | python3 -c "import sys,json; print(json.dumps(sys.stdin.read()))")
PAYLOAD="{"content": $ESC_CONTENT, "description": $ESC_DESC}"
if [ -n "$TASK_ID" ]; then
  "$DIR/todoist_api.sh" "tasks/$TASK_ID" POST "$PAYLOAD"
else
  "$DIR/todoist_api.sh" "tasks" POST "$PAYLOAD"
fi
