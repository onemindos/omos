#!/usr/bin/env bash
# morning-briefing.sh
# Pulls Strava (last activity) + Todoist (today's tasks) in parallel
# and prints a combined morning brief to stdout.
#
# Usage:
#   bash morning-briefing.sh
#   bash morning-briefing.sh --quiet   # JSON only, no formatting
#
# Requires: composio CLI authenticated, python3, jq (optional)

set -euo pipefail

QUIET=0
[[ "${1:-}" == "--quiet" ]] && QUIET=1

STRAVA_FILE="/tmp/composio_strava_$$.json"
TODOIST_FILE="/tmp/composio_todoist_$$.json"

cleanup() { rm -f "$STRAVA_FILE" "$TODOIST_FILE"; }
trap cleanup EXIT

# ── Parallel fetch ──────────────────────────────────────────────────────────
if [[ $QUIET -eq 0 ]]; then
    echo "⏳ Fetching Strava + Todoist in parallel..."
fi

composio execute STRAVA_LIST_ATHLETE_ACTIVITIES -d '{"per_page": 1}' > "$STRAVA_FILE" &
PID_STRAVA=$!

# Get today's tasks using Todoist filter syntax
composio execute TODOIST_FILTER_TASKS -d '{"filter": "today | overdue"}' > "$TODOIST_FILE" &
PID_TODOIST=$!

wait $PID_STRAVA
wait $PID_TODOIST

if [[ $QUIET -eq 1 ]]; then
    echo "=== STRAVA ==="
    cat "$STRAVA_FILE"
    echo ""
    echo "=== TODOIST ==="
    cat "$TODOIST_FILE"
    exit 0
fi

# ── Parse and format ─────────────────────────────────────────────────────────
TODAY=$(date "+%A, %B %d, %Y")
echo ""
echo "╔══════════════════════════════════════════════════════════╗"
echo "  🌅 ONE MIND MORNING BRIEFING — $TODAY"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""

# Strava section
echo "──── 💪 LAST ACTIVITY (Strava) ────────────────────────────"
python3 << PYEOF
import json, sys

with open("$STRAVA_FILE") as f:
    raw = json.load(f)

activities = raw.get("data", {}).get("details", [])
if not activities:
    print("  No recent activities found.")
    sys.exit(0)

act = activities[0]
name = act.get("name", "Unknown")
sport = act.get("sport_type") or act.get("type", "?")
date = act.get("start_date_local", "")[:10]
elapsed = act.get("elapsed_time", 0)
h, m = elapsed // 3600, (elapsed % 3600) // 60
dur = f"{h}h {m}m" if h else f"{m}m"
avg_hr = act.get("average_heartrate")
max_hr = act.get("max_heartrate")
dist_m = act.get("distance", 0)
dist = f"{dist_m/1000:.2f} km" if dist_m > 100 else "—"

print(f"  📅 {date}  |  {sport}: {name}")
print(f"  ⏱  Duration: {dur}  |  Distance: {dist}")
if avg_hr:
    print(f"  ❤️  Avg HR: {avg_hr:.1f} bpm  |  Max HR: {max_hr:.0f} bpm")
    if max_hr and max_hr > 180:
        print("  ⚠️  High intensity session — consider recovery today")
PYEOF

echo ""

# Todoist section
echo "──── ✅ TODAY'S TASKS (Todoist) ────────────────────────────"
python3 << PYEOF
import json, sys

with open("$TODOIST_FILE") as f:
    raw = json.load(f)

# Handle both array and {results: [...]} shapes
if isinstance(raw, dict):
    tasks = raw.get("data", raw.get("results", []))
    if isinstance(tasks, dict):
        tasks = tasks.get("tasks", tasks.get("items", []))
else:
    tasks = raw

if not tasks:
    print("  ✨ No tasks due today — clear schedule!")
    sys.exit(0)

# Sort: priority desc (p1 first = value 4), then by content
PRIORITY_LABEL = {4: "🔴 P1", 3: "🟠 P2", 2: "🟡 P3", 1: "⚪ P4"}
tasks_sorted = sorted(tasks, key=lambda t: -t.get("priority", 1))

overdue = [t for t in tasks_sorted if t.get("is_overdue") or (t.get("due") or {}).get("is_overdue")]
today_tasks = [t for t in tasks_sorted if t not in overdue]

if overdue:
    print(f"  ⏰ OVERDUE ({len(overdue)}):")
    for t in overdue[:5]:
        p = PRIORITY_LABEL.get(t.get("priority", 1), "")
        print(f"     {p} {t['content'][:60]}")

if today_tasks:
    print(f"  📋 TODAY ({len(today_tasks)}):")
    for t in today_tasks[:10]:
        p = PRIORITY_LABEL.get(t.get("priority", 1), "")
        print(f"     {p} {t['content'][:60]}")

if len(tasks) > 15:
    print(f"  ... and {len(tasks) - 15} more tasks")
PYEOF

echo ""
echo "──── 🔗 Quick Links ────────────────────────────────────────"
echo "  Strava:  https://www.strava.com/dashboard"
echo "  Todoist: https://app.todoist.com/app/today"
echo ""
