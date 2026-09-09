#!/bin/bash
# Activate a servant agent's cron jobs
# Usage: ./activate-servant.sh <agent-name|all>

AGENT=$1
PATH="/opt/hermes/.venv/bin:$PATH"

if [ -z "$AGENT" ]; then
    echo "Usage: $0 <agent-name|all>"
    echo "  agent-name: oracle|guardian|grid|spartan|haven|heritage|empire|eden|trinity"
    echo "  all: activate ALL servant agents"
    exit 1
fi

if [ "$AGENT" = "all" ]; then
    for name in oracle guardian grid spartan haven heritage empire eden trinity; do
        echo "=== Activating $name ==="
        hermes cron list | grep -i "$name" | while read -r line; do
            jobid=$(echo "$line" | awk '{print $1}')
            if [ -n "$jobid" ]; then
                hermes cron resume "$jobid"
                echo "  + $jobid"
            fi
        done
    done
else
    echo "=== Activating $AGENT ==="
    hermes cron list | grep -i "$AGENT" | while read -r line; do
        jobid=$(echo "$line" | awk '{print $1}')
        if [ -n "$jobid" ]; then
            hermes cron resume "$jobid"
            echo "  + $jobid"
        fi
    done
fi
echo "Done"
