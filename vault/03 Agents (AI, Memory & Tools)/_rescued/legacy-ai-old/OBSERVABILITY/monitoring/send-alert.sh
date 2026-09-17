#!/bin/bash
################################################################################
# LEGACY AI ALERT SYSTEM
# Purpose: Send alerts via Slack, Discord, or Email
# Usage: ./scripts/monitoring/send-alert.sh <level> <message> [details]
# Levels: critical, warning, info
################################################################################

ALERT_LEVEL=$1
MESSAGE=$2
DETAILS=$3

if [ -z "$ALERT_LEVEL" ] || [ -z "$MESSAGE" ]; then
    echo "Usage: $0 <level> <message> [details]"
    echo "Levels: critical, warning, info"
    exit 1
fi

# Configuration (can be overridden by environment variables)
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL:-}"
DISCORD_WEBHOOK_URL="${DISCORD_WEBHOOK_URL:-}"
EMAIL_TO="${ALERT_EMAIL_TO:-}"
HOSTNAME=$(hostname)
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Color and emoji based on alert level
case "$ALERT_LEVEL" in
    critical)
        COLOR="danger"
        DISCORD_COLOR=15158332  # Red
        EMOJI="🚨"
        ;;
    warning)
        COLOR="warning"
        DISCORD_COLOR=16776960  # Yellow
        EMOJI="⚠️"
        ;;
    info)
        COLOR="good"
        DISCORD_COLOR=3447003   # Blue
        EMOJI="ℹ️"
        ;;
    *)
        COLOR="good"
        DISCORD_COLOR=3447003
        EMOJI="ℹ️"
        ;;
esac

################################################################################
# SEND TO SLACK
################################################################################
if [ -n "$SLACK_WEBHOOK_URL" ]; then
    SLACK_PAYLOAD=$(cat <<EOF
{
  "text": "$EMOJI $MESSAGE",
  "attachments": [
    {
      "color": "$COLOR",
      "fields": [
        {
          "title": "Host",
          "value": "$HOSTNAME",
          "short": true
        },
        {
          "title": "Level",
          "value": "${ALERT_LEVEL^^}",
          "short": true
        },
        {
          "title": "Timestamp",
          "value": "$TIMESTAMP",
          "short": false
        }
        $([ -n "$DETAILS" ] && echo ",{\"title\": \"Details\", \"value\": \"$DETAILS\", \"short\": false}")
      ]
    }
  ]
}
EOF
)
    
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$SLACK_WEBHOOK_URL" \
        -H 'Content-Type: application/json' \
        -d "$SLACK_PAYLOAD")
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo "✅ Alert sent to Slack"
    else
        echo "❌ Failed to send to Slack (HTTP $HTTP_CODE)"
    fi
fi

################################################################################
# SEND TO DISCORD
################################################################################
if [ -n "$DISCORD_WEBHOOK_URL" ]; then
    DISCORD_PAYLOAD=$(cat <<EOF
{
  "content": "$EMOJI **${ALERT_LEVEL^^}**: $MESSAGE",
  "embeds": [
    {
      "title": "Legacy AI Infrastructure Alert",
      "color": $DISCORD_COLOR,
      "fields": [
        {
          "name": "Host",
          "value": "$HOSTNAME",
          "inline": true
        },
        {
          "name": "Level",
          "value": "${ALERT_LEVEL^^}",
          "inline": true
        },
        {
          "name": "Timestamp",
          "value": "$TIMESTAMP",
          "inline": false
        }
        $([ -n "$DETAILS" ] && echo ",{\"name\": \"Details\", \"value\": \"$DETAILS\", \"inline\": false}")
      ]
    }
  ]
}
EOF
)
    
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$DISCORD_WEBHOOK_URL" \
        -H 'Content-Type: application/json' \
        -d "$DISCORD_PAYLOAD")
    
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "204" ]; then
        echo "✅ Alert sent to Discord"
    else
        echo "❌ Failed to send to Discord (HTTP $HTTP_CODE)"
    fi
fi

################################################################################
# SEND EMAIL
################################################################################
if [ -n "$EMAIL_TO" ]; then
    # Check if mail command is available
    if command -v mail >/dev/null 2>&1; then
        EMAIL_SUBJECT="[$ALERT_LEVEL_UPPER] Legacy AI Alert: $MESSAGE"
        EMAIL_BODY=$(cat <<EOF
Legacy AI Infrastructure Alert

Level: ${ALERT_LEVEL^^}
Host: $HOSTNAME
Timestamp: $TIMESTAMP

Message:
$MESSAGE

$([ -n "$DETAILS" ] && echo -e "Details:\n$DETAILS")

--
Automated alert from Legacy AI monitoring system
EOF
)
        
        echo "$EMAIL_BODY" | mail -s "$EMAIL_SUBJECT" "$EMAIL_TO"
        
        if [ $? -eq 0 ]; then
            echo "✅ Alert sent via email to $EMAIL_TO"
        else
            echo "❌ Failed to send email"
        fi
    else
        echo "⚠️  Mail command not available - skipping email alert"
    fi
fi

################################################################################
# LOG TO FILE
################################################################################
LOG_DIR="/var/log/legacy-ai-alerts"
mkdir -p "$LOG_DIR" 2>/dev/null || LOG_DIR="/tmp/legacy-ai-alerts"
mkdir -p "$LOG_DIR"

LOG_FILE="$LOG_DIR/alerts-$(date +%Y%m).log"

echo "[$(date +%Y-%m-%d\ %H:%M:%S)] [$ALERT_LEVEL_UPPER] $MESSAGE" >> "$LOG_FILE"
if [ -n "$DETAILS" ]; then
    echo "  Details: $DETAILS" >> "$LOG_FILE"
fi

echo "📝 Alert logged to: $LOG_FILE"

################################################################################
# SUMMARY
################################################################################
if [ -z "$SLACK_WEBHOOK_URL" ] && [ -z "$DISCORD_WEBHOOK_URL" ] && [ -z "$EMAIL_TO" ]; then
    echo ""
    echo "⚠️  No alert channels configured!"
    echo "To enable alerts, set one or more of:"
    echo "  - export SLACK_WEBHOOK_URL='https://hooks.slack.com/...'"
    echo "  - export DISCORD_WEBHOOK_URL='https://discord.com/api/webhooks/...'"
    echo "  - export ALERT_EMAIL_TO='your@email.com'"
    echo ""
    echo "Alert message: $EMOJI $MESSAGE"
    [ -n "$DETAILS" ] && echo "Details: $DETAILS"
fi

exit 0
