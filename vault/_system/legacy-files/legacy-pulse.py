#!/usr/bin/env python3
"""
legacy-pulse.py — NATS Event Bridge for Legacy/OpenClaw
Subscribes to relevant NATS subjects and writes events
into AMBIENT_STATE.md so Legacy reads them on each heartbeat.

Run alongside OpenClaw gateway and legacy-eyes:
  python3 legacy-pulse.py &

Requirements:
  pip install nats-py
  NATS_URL defaults to nats://localhost:4222
"""

import asyncio
import json
import os
from datetime import datetime
from pathlib import Path

import nats

# ── Config ────────────────────────────────────────────────────────────────────

NATS_URL       = os.environ.get("NATS_URL", "nats://localhost:4222")
WORKSPACE_DIR  = os.path.expanduser("~/.openclaw/workspace")
AMBIENT_STATE  = os.path.join(WORKSPACE_DIR, "AMBIENT_STATE.md")
MAX_EVENTS     = 20

# Subjects Legacy cares about — in priority order
SUBSCRIPTIONS = [
    "system.alerts.>",     # P0 — always surface
    "legacy.session.>",    # session lifecycle
    "omcs.hp.>",           # health/wellness signals
    "omcs.ge.>",           # business/revenue signals
    "omcs.it.>",           # infra/tech signals
    "omcs.le.>",           # personal growth signals
    "it.droplet.>",        # DigitalOcean droplet events
    "location.update",     # Zeus location changes
]

# ── State helpers ─────────────────────────────────────────────────────────────

def read_state() -> dict:
    if not os.path.exists(AMBIENT_STATE):
        return {"nats_events": [], "vision_digest": "No data.", "last_updated": None}
    try:
        with open(AMBIENT_STATE, "r") as f:
            content = f.read()
        start = content.find("```json") + 7
        end   = content.find("```", start)
        if start > 7 and end > start:
            return json.loads(content[start:end].strip())
    except Exception:
        pass
    return {"nats_events": [], "vision_digest": "No data.", "last_updated": None}


def write_state(state: dict):
    os.makedirs(WORKSPACE_DIR, exist_ok=True)
    ts = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    state["last_updated"] = ts

    vision  = state.get("vision_digest", "No data.")
    session = state.get("session_id", "none")
    loc     = state.get("location", "unknown")
    unchanged = state.get("scene_unchanged_count", 0)
    events  = state.get("nats_events", [])[-MAX_EVENTS:]

    lines = [
        "# Ambient State — Legacy Pulse",
        f"*Updated: {ts}*",
        "",
        "## Vision Digest",
        vision,
        "",
        "## Active Session",
        f"- Session ID: `{session}`",
        f"- Location: `{loc}`",
        f"- Scene unchanged for: {unchanged} consecutive ticks",
        "",
        f"## Recent NATS Events (last {MAX_EVENTS})",
    ]

    if events:
        for ev in events:
            lines.append(f"- `{ev.get('subject','?')}` — {ev.get('summary','?')} @ {ev.get('ts','?')}")
    else:
        lines.append("_No events yet._")

    lines += ["", "---", "```json", json.dumps(state, indent=2), "```", ""]

    with open(AMBIENT_STATE, "w") as f:
        f.write("\n".join(lines))


def summarize_payload(subject: str, data: dict) -> str:
    """Generate a human-readable 1-line summary of a NATS event."""
    if subject.startswith("system.alerts"):
        return f"⚠️ {data.get('message', str(data))}"
    if subject.startswith("legacy.session.start"):
        return f"Session started @ {data.get('location','unknown')}"
    if subject.startswith("legacy.session.end"):
        return f"Session ended — {data.get('duration','?')} min"
    if subject.startswith("omcs.hp"):
        return f"HP signal: {data.get('metric','?')} = {data.get('value','?')}"
    if subject.startswith("omcs.ge"):
        return f"GE event: {data.get('event','?')} — {data.get('detail','?')}"
    if subject.startswith("it.droplet"):
        return f"Droplet: {data.get('status','?')} CPU={data.get('cpu_pct','?')}%"
    if subject == "location.update":
        return f"Zeus → {data.get('location','?')}"
    return str(data)[:80]


# ── NATS message handler ──────────────────────────────────────────────────────

async def message_handler(msg):
    subject = msg.subject
    try:
        data = json.loads(msg.data.decode())
    except Exception:
        data = {"raw": msg.data.decode()}

    ts = datetime.now().strftime("%H:%M:%S")
    summary = summarize_payload(subject, data)

    print(f"[legacy-pulse] {subject} → {summary}")

    state = read_state()
    events = state.get("nats_events", [])
    events.append({
        "subject": subject,
        "summary": summary,
        "ts": ts,
        "raw": data
    })
    # Keep only last MAX_EVENTS
    state["nats_events"] = events[-MAX_EVENTS:]

    # Update session state if it's a session event
    if subject.startswith("legacy.session.start"):
        state["session_id"] = data.get("session_id", "unknown")
        state["location"]   = data.get("location", "unknown")

    write_state(state)


# ── Main ──────────────────────────────────────────────────────────────────────

async def run():
    print(f"[legacy-pulse] Connecting to NATS at {NATS_URL}...")

    nc = await nats.connect(NATS_URL)
    print("[legacy-pulse] Connected.")

    for subject in SUBSCRIPTIONS:
        await nc.subscribe(subject, cb=message_handler)
        print(f"[legacy-pulse] Subscribed: {subject}")

    print("[legacy-pulse] Listening for OMCS events...")

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        print("[legacy-pulse] Shutting down.")
        await nc.drain()


if __name__ == "__main__":
    asyncio.run(run())
