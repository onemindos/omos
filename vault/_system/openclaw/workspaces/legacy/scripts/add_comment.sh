#!/bin/bash
# Add progress comment to Todoist task
# Usage: add_comment.sh <task_id> <comment_text>
TASK_ID=$1
TEXT=$2
DIR="$(dirname "$0")"
ESC=$(printf '%s' "$TEXT" | python3 -c "import sys,json; print(json.dumps(sys.stdin.read()))")
"$DIR/todoist_api.sh" "comments" POST "{"task_id": "$TASK_ID", "content": $ESC}"
