#!/usr/bin/env python3
"""
legacy-session-start.py — Fire a session start event into NATS
Run this when Zeus opens his command center (Discord, webchat, etc.)
Legacy receives it on the next heartbeat tick and wakes fully aware.

Usage:
  python3 legacy-session-start.py --location home
  python3 legacy-session-start.py --location office --omcs IT.31.P01

Requirements:
  pip install nats-py
"""

import asyncio
import json
import uuid
import argparse
import subprocess
from datetime import datetime
import nats
import os

NATS_URL = os.environ.get("NATS_URL", "nats://localhost:4222")


def get_active_apps() -> list:
    """Detect what's running (macOS only, graceful fallback)."""
    try:
        out = subprocess.check_output(
            ["osascript", "-e",
             'tell application "System Events" to get name of (processes where background only is false)'],
            timeout=3, text=True
        )
        return [x.strip() for x in out.split(",")][:10]
    except Exception:
        return []


async def fire_session_start(location: str, omcs_coord: str):
    nc = await nats.connect(NATS_URL)

    session_id = str(uuid.uuid4())[:8]
    ts = datetime.now().isoformat()

    payload = {
        "session_id": session_id,
        "location": location,
        "omcs_coordinate": omcs_coord,
        "started_at": ts,
        "active_apps": get_active_apps()
    }

    subject = "legacy.session.start"
    await nc.publish(subject, json.dumps(payload).encode())
    print(f"[session-start] Fired: {subject}")
    print(f"[session-start] Session ID: {session_id}")
    print(f"[session-start] Location: {location}")
    print(f"[session-start] OMCS: {omcs_coord}")

    await nc.drain()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--location", default="home", help="home|farm|office|transit")
    parser.add_argument("--omcs", default="IT.00.P00", help="OMCS coordinate for this session")
    args = parser.parse_args()
    asyncio.run(fire_session_start(args.location, args.omcs))
