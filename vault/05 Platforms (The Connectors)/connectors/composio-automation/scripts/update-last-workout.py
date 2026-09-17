#!/usr/bin/env python3
"""
update-last-workout.py
Fetches the most recent Strava workout-type activity and writes it to
/Users/zeuslegacy/Documents/onemind-codex/LAST-WORKOUT.md

Run:
  composio execute STRAVA_LIST_ATHLETE_ACTIVITIES -d '{"per_page": 10}' | python3 update-last-workout.py
  -- OR --
  python3 update-last-workout.py   (will call composio directly)
"""

import json
import subprocess
import sys
from datetime import datetime, timezone

CODEX_ROOT = "/Users/zeuslegacy/Documents/onemind-codex"
OUTPUT_FILE = f"{CODEX_ROOT}/LAST-WORKOUT.md"

WORKOUT_TYPES = {"Workout", "WeightTraining", "Crossfit", "Yoga", "Pilates",
                 "Rowing", "StairStepper", "Elliptical", "RockClimbing",
                 "Swim", "Run", "Ride", "Walk", "Hike", "VirtualRun",
                 "VirtualRide", "EBikeRide"}

def fetch_activities():
    """Fetch recent activities via composio CLI, or read from stdin."""
    if not sys.stdin.isatty():
        raw = sys.stdin.read()
    else:
        result = subprocess.run(
            ["composio", "execute", "STRAVA_LIST_ATHLETE_ACTIVITIES", "-d", '{"per_page": 20}'],
            capture_output=True, text=True, timeout=30
        )
        if result.returncode != 0:
            raise RuntimeError(f"composio failed: {result.stderr}")
        raw = result.stdout
    return json.loads(raw)


def find_last_workout(data):
    """Return the most recent activity of any workout type."""
    activities = data.get("data", {}).get("details", [])
    for act in activities:
        if act.get("type") in WORKOUT_TYPES or act.get("sport_type") in WORKOUT_TYPES:
            return act
    # Fallback: return the most recent activity regardless of type
    return activities[0] if activities else None


def format_duration(seconds):
    h = seconds // 3600
    m = (seconds % 3600) // 60
    s = seconds % 60
    if h:
        return f"{h}h {m}m"
    return f"{m}m {s}s"


def format_date(iso_str):
    """Convert ISO datetime string to readable format."""
    try:
        dt = datetime.fromisoformat(iso_str.replace("Z", "+00:00"))
        return dt.strftime("%B %d, %Y at %I:%M %p UTC")
    except Exception:
        return iso_str


def build_markdown(act):
    name = act.get("name", "Unknown")
    sport = act.get("sport_type") or act.get("type", "Unknown")
    date = format_date(act.get("start_date_local") or act.get("start_date", ""))
    duration = format_duration(act.get("elapsed_time", 0))
    moving_time = format_duration(act.get("moving_time", 0))
    avg_hr = act.get("average_heartrate")
    max_hr = act.get("max_heartrate")
    distance_m = act.get("distance", 0)
    distance = f"{distance_m / 1000:.2f} km" if distance_m > 100 else "—"
    calories = act.get("kilojoules")
    calories_str = f"{calories:.0f} kJ" if calories else "—"
    trainer = act.get("trainer", False)
    gear = act.get("gear_id", "—")
    device = act.get("device_name", "—")
    avg_watts = act.get("average_watts", "—")
    max_watts = act.get("max_watts", "—")
    elevation = act.get("total_elevation_gain", 0)
    activity_id = act.get("id", "")
    strava_url = f"https://www.strava.com/activities/{activity_id}" if activity_id else "—"
    updated_at = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M UTC")

    hr_analysis = ""
    if avg_hr and max_hr:
        if max_hr > 185:
            hr_note = "🔴 Peak effort moment detected (max HR > 185)"
        elif max_hr > 170:
            hr_note = "🟠 High intensity peak (max HR 170–185)"
        else:
            hr_note = "🟢 Moderate effort throughout"

        if avg_hr > 160:
            zone = "Zone 4–5 (Threshold/VO2 Max)"
        elif avg_hr > 140:
            zone = "Zone 3–4 (Aerobic/Tempo)"
        elif avg_hr > 120:
            zone = "Zone 2–3 (Aerobic Base)"
        else:
            zone = "Zone 1–2 (Recovery/Easy)"

        hr_analysis = f"""
## Heart Rate Analysis

| Metric | Value |
|--------|-------|
| Average HR | {avg_hr:.1f} bpm |
| Max HR | {max_hr:.0f} bpm |
| Estimated Zone | {zone} |

{hr_note}
"""

    power_section = ""
    if avg_watts and avg_watts != "—":
        power_section = f"""
## Power Data

| Metric | Value |
|--------|-------|
| Avg Power | {avg_watts:.0f} W |
| Max Power | {max_watts} W |
"""

    return f"""---
title: "Last Workout — {name}"
type: metric
status: active
updated: {updated_at}
tags: [type/metric, domain/hp, topic/fitness, source/strava]
ai_summary: "Last Strava workout: {name} on {date[:12]}. Duration {duration}. Avg HR {avg_hr:.1f if avg_hr else '—'} bpm."
---

# Last Workout: {name}

> Auto-generated from Strava. Last updated: {updated_at}

## Overview

| Field | Value |
|-------|-------|
| Activity | {name} |
| Type | {sport} |
| Date | {date} |
| Total Time | {duration} |
| Moving Time | {moving_time} |
| Distance | {distance} |
| Indoor/Trainer | {"Yes" if trainer else "No"} |
| Device | {device} |
| Calories | {calories_str} |
| Elevation Gain | {elevation:.0f} m |
| Strava Link | [{activity_id}]({strava_url}) |
{hr_analysis}{power_section}
## Raw JSON (for agents)

```json
{json.dumps(act, indent=2)}
```
"""


def main():
    try:
        data = fetch_activities()
        act = find_last_workout(data)
        if not act:
            print("No activities found.", file=sys.stderr)
            sys.exit(1)

        md = build_markdown(act)
        with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
            f.write(md)

        print(f"✓ Updated {OUTPUT_FILE}")
        print(f"  Activity: {act.get('name')} ({act.get('sport_type') or act.get('type')})")
        print(f"  Date: {act.get('start_date_local', '')[:10]}")
        if act.get("average_heartrate"):
            print(f"  Avg HR: {act['average_heartrate']:.1f} bpm  |  Max HR: {act.get('max_heartrate', '—')} bpm")

    except Exception as e:
        print(f"✗ Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
