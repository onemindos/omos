#!/usr/bin/env python3
"""
high-hr-recovery-task.py
Checks the last Strava activity's average HR.
If it exceeds the threshold, auto-creates a Todoist recovery task for tomorrow.

Usage:
  python3 high-hr-recovery-task.py
  python3 high-hr-recovery-task.py --threshold 160
  python3 high-hr-recovery-task.py --dry-run

Requires: composio CLI authenticated
"""

import argparse
import json
import subprocess
import sys

DEFAULT_THRESHOLD = 155   # avg HR bpm above which recovery task is created
TODOIST_HP_PROJECT = None  # set to your HP project ID to route correctly, e.g. "6g6P9pqrf5PCR9P5"


def composio(slug, params):
    result = subprocess.run(
        ["composio", "execute", slug, "-d", json.dumps(params)],
        capture_output=True, text=True, timeout=30
    )
    if result.returncode != 0:
        raise RuntimeError(f"composio {slug} failed: {result.stderr}")
    return json.loads(result.stdout)


def main():
    parser = argparse.ArgumentParser(description="Create recovery task if last Strava HR was high")
    parser.add_argument("--threshold", type=int, default=DEFAULT_THRESHOLD,
                        help=f"Avg HR threshold (default: {DEFAULT_THRESHOLD})")
    parser.add_argument("--dry-run", action="store_true",
                        help="Print what would happen without creating the task")
    args = parser.parse_args()

    # Fetch last activity
    print("🔍 Fetching last Strava activity...")
    data = composio("STRAVA_LIST_ATHLETE_ACTIVITIES", {"per_page": 1})
    activities = data.get("data", {}).get("details", [])

    if not activities:
        print("No activities found. Nothing to do.")
        sys.exit(0)

    act = activities[0]
    name = act.get("name", "Unknown")
    sport = act.get("sport_type") or act.get("type", "?")
    avg_hr = act.get("average_heartrate")
    max_hr = act.get("max_heartrate", 0)
    date_str = act.get("start_date_local", "")[:10]
    elapsed = act.get("elapsed_time", 0)
    dur = f"{elapsed // 3600}h {(elapsed % 3600) // 60}m" if elapsed >= 3600 else f"{elapsed // 60}m"

    print(f"  Last activity: {name} ({sport}) on {date_str}")
    print(f"  Duration: {dur}  |  Avg HR: {avg_hr if avg_hr else '—'} bpm  |  Max HR: {max_hr} bpm")

    if not avg_hr:
        print("⚠️  No heart rate data in this activity. Skipping recovery check.")
        sys.exit(0)

    if avg_hr < args.threshold:
        print(f"✅ Avg HR {avg_hr:.1f} < threshold {args.threshold} — no recovery task needed.")
        sys.exit(0)

    # Build recovery task
    intensity = "high intensity" if avg_hr > 165 else "moderate-high intensity"
    task_content = f"🔄 Recovery Day — after {intensity} session ({name})"
    task_description = (
        f"Yesterday: {sport} · {dur} · Avg HR {avg_hr:.1f} bpm · Max HR {max_hr:.0f} bpm\n"
        f"Priority: hydrate, sleep 8h, light walk or stretch only. No hard training."
    )

    task_params = {
        "content": task_content,
        "description": task_description,
        "due_string": "tomorrow",
        "priority": 3,  # P2 (high)
        "labels": ["recovery", "hp"],
    }
    if TODOIST_HP_PROJECT:
        task_params["project_id"] = TODOIST_HP_PROJECT

    print(f"\n⚠️  Avg HR {avg_hr:.1f} >= threshold {args.threshold} — recovery task needed!")
    print(f"   Task: {task_content}")

    if args.dry_run:
        print("\n[DRY RUN] Would create:")
        print(json.dumps(task_params, indent=2))
        sys.exit(0)

    # Create the task
    result = composio("TODOIST_CREATE_TASK", task_params)
    task_id = result.get("data", {}).get("task_id") or result.get("id", "?")
    print(f"\n✅ Created recovery task (ID: {task_id})")
    print(f"   Due: tomorrow  |  Priority: P2")


if __name__ == "__main__":
    main()
