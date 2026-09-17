#!/usr/bin/env python3
"""
strava-weekly-summary.py
Pulls the last 7 days of Strava activities and prints a weekly fitness summary.
Optionally creates a Todoist task summarizing the week's training load.

Usage:
  python3 strava-weekly-summary.py
  python3 strava-weekly-summary.py --create-task
  python3 strava-weekly-summary.py --days 14

Requires: composio CLI authenticated
"""

import argparse
import json
import subprocess
import sys
from datetime import datetime, timedelta, timezone
from collections import defaultdict


def composio(slug, params):
    result = subprocess.run(
        ["composio", "execute", slug, "-d", json.dumps(params)],
        capture_output=True, text=True, timeout=30
    )
    if result.returncode != 0:
        raise RuntimeError(f"composio {slug} failed: {result.stderr}")
    return json.loads(result.stdout)


def fmt_dur(seconds):
    h, m = seconds // 3600, (seconds % 3600) // 60
    return f"{h}h {m}m" if h else f"{m}m"


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--days", type=int, default=7, help="Days to look back (default: 7)")
    parser.add_argument("--create-task", action="store_true", help="Create Todoist summary task")
    args = parser.parse_args()

    cutoff = datetime.now(timezone.utc) - timedelta(days=args.days)
    cutoff_ts = int(cutoff.timestamp())

    print(f"🏃 Fetching activities from last {args.days} days...")

    data = composio("STRAVA_LIST_ATHLETE_ACTIVITIES", {
        "after": cutoff_ts,
        "per_page": 50
    })
    activities = data.get("data", {}).get("details", [])

    if not activities:
        print(f"No activities found in the last {args.days} days.")
        return

    # Aggregate
    total_time = sum(a.get("elapsed_time", 0) for a in activities)
    total_dist = sum(a.get("distance", 0) for a in activities)
    hr_readings = [a["average_heartrate"] for a in activities if a.get("average_heartrate")]
    avg_hr_overall = sum(hr_readings) / len(hr_readings) if hr_readings else None
    max_hr_overall = max((a.get("max_heartrate", 0) for a in activities), default=0)

    by_type = defaultdict(int)
    for a in activities:
        by_type[a.get("sport_type") or a.get("type", "Other")] += 1

    # Print summary
    period_str = f"Last {args.days} Days"
    print(f"\n{'='*56}")
    print(f"  📊 STRAVA SUMMARY — {period_str}")
    print(f"  {cutoff.strftime('%b %d')} → {datetime.now().strftime('%b %d, %Y')}")
    print(f"{'='*56}")
    print(f"  Sessions:    {len(activities)}")
    print(f"  Total Time:  {fmt_dur(total_time)}")
    if total_dist > 0:
        print(f"  Distance:    {total_dist/1000:.1f} km")
    if avg_hr_overall:
        print(f"  Avg HR:      {avg_hr_overall:.1f} bpm")
        print(f"  Peak HR:     {max_hr_overall:.0f} bpm")
    print(f"\n  By type:")
    for t, count in sorted(by_type.items(), key=lambda x: -x[1]):
        print(f"    {t}: {count}x")

    print(f"\n  Sessions:")
    for a in activities:
        date = a.get("start_date_local", "")[:10]
        name = a.get("name", "?")[:40]
        sport = a.get("sport_type") or a.get("type", "?")
        dur = fmt_dur(a.get("elapsed_time", 0))
        hr = f"HR {a['average_heartrate']:.0f}" if a.get("average_heartrate") else ""
        print(f"    {date}  {sport}: {name:<40}  {dur}  {hr}")

    print(f"{'='*56}")

    if args.create_task:
        # Build task description
        hr_line = f"Avg HR: {avg_hr_overall:.1f} bpm · Peak: {max_hr_overall:.0f} bpm" if avg_hr_overall else ""
        dist_line = f"Distance: {total_dist/1000:.1f} km" if total_dist > 100 else ""
        types_line = " · ".join(f"{k} ×{v}" for k, v in sorted(by_type.items(), key=lambda x: -x[1]))

        desc_parts = [
            f"{len(activities)} sessions in {args.days} days · {fmt_dur(total_time)} total",
            dist_line,
            hr_line,
            types_line,
        ]
        description = "\n".join(p for p in desc_parts if p)

        result = composio("TODOIST_CREATE_TASK", {
            "content": f"📊 Weekly Training Review ({cutoff.strftime('%b %d')}–{datetime.now().strftime('%b %d')})",
            "description": description,
            "due_string": "today",
            "priority": 2,
            "labels": ["review", "hp"],
        })
        task_id = result.get("data", {}).get("task_id") or result.get("id", "?")
        print(f"\n✅ Created review task in Todoist (ID: {task_id})")


if __name__ == "__main__":
    main()
