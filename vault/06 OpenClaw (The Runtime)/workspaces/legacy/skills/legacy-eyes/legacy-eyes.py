#!/usr/bin/env python3
"""
legacy-eyes.py — 1fps Vision Loop for Legacy/OpenClaw
Captures webcam or RTSP frames, sends to Claude vision API,
writes digest to AMBIENT_STATE.md in the OpenClaw workspace.

Run as a background daemon alongside the OpenClaw gateway:
  python3 legacy-eyes.py &

Requirements:
  pip install anthropic opencv-python
  ANTHROPIC_API_KEY must be set in environment
"""

import anthropic
import base64
import cv2
import json
import os
import time
from datetime import datetime
from pathlib import Path

# ── Config ────────────────────────────────────────────────────────────────────

WORKSPACE_DIR   = os.environ.get("OPENCLAW_WORKSPACE", os.path.expanduser("~/.openclaw/workspace"))
AMBIENT_STATE   = os.path.join(WORKSPACE_DIR, "AMBIENT_STATE.md")
VISION_SOURCE   = os.environ.get("LEGACY_VISION_SOURCE", "0")  # "0" = webcam, or RTSP URL
CAPTURE_INTERVAL = 1.0        # seconds between captures (1fps)
DIGEST_INTERVAL  = 10         # run vision API every N frames (every 10 seconds)
MAX_NATS_EVENTS  = 20         # max NATS events to buffer in state file

client = anthropic.Anthropic()

# ── Vision digest prompt ───────────────────────────────────────────────────────

VISION_PROMPT = """You are Legacy's eyes. Describe what you see in 2-3 sentences max.
Focus on:
- What Zeus (the user) is doing
- Any notable changes in environment
- Anything actionable or worth flagging

Be concise. Factual. Skip the obvious. No preamble.
If nothing significant is visible, reply: SCENE_UNCHANGED"""

# ── State management ──────────────────────────────────────────────────────────

def read_ambient_state() -> dict:
    """Read current ambient state or return empty template."""
    if not os.path.exists(AMBIENT_STATE):
        return {
            "last_updated": None,
            "vision_digest": "No vision data yet.",
            "scene_unchanged_count": 0,
            "nats_events": [],
            "session_id": None,
            "location": "unknown"
        }
    try:
        with open(AMBIENT_STATE, "r") as f:
            content = f.read()
        start = content.find("```json") + 7
        end   = content.find("```", start)
        if start > 7 and end > start:
            return json.loads(content[start:end].strip())
    except Exception:
        pass
    return {
        "last_updated": None,
        "vision_digest": "No vision data yet.",
        "scene_unchanged_count": 0,
        "nats_events": [],
        "session_id": None,
        "location": "unknown"
    }


def write_ambient_state(state: dict):
    """Write ambient state as readable markdown + embedded JSON to workspace."""
    os.makedirs(WORKSPACE_DIR, exist_ok=True)
    ts = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    state["last_updated"] = ts

    content = f"""# Ambient State — Legacy Eyes
*Updated: {ts}*

## Vision Digest
{state.get("vision_digest", "No data.")}

## Active Session
- Session ID: `{state.get("session_id", "none")}`
- Location: `{state.get("location", "unknown")}`
- Scene unchanged for: {state.get("scene_unchanged_count", 0)} consecutive ticks

## Recent NATS Events (last {MAX_NATS_EVENTS})
"""
    events = state.get("nats_events", [])
    if events:
        for ev in events[-MAX_NATS_EVENTS:]:
            content += f"- `{ev.get('subject','?')}` — {ev.get('summary','?')} @ {ev.get('ts','?')}\n"
    else:
        content += "_No events yet._\n"

    content += f"""
---
```json
{json.dumps(state, indent=2)}
```
"""
    with open(AMBIENT_STATE, "w") as f:
        f.write(content)


# ── Frame capture ─────────────────────────────────────────────────────────────

def capture_frame(source) -> bytes | None:
    """Capture a single frame from webcam or RTSP source."""
    cap = cv2.VideoCapture(int(source) if source.isdigit() else source)
    ret, frame = cap.read()
    cap.release()
    if not ret:
        return None
    _, buf = cv2.imencode(".jpg", frame, [cv2.IMWRITE_JPEG_QUALITY, 75])
    return buf.tobytes()


# ── Vision API call ───────────────────────────────────────────────────────────

def get_vision_digest(frame_bytes: bytes) -> str:
    """Send frame to Claude vision and get a scene digest."""
    b64 = base64.standard_b64encode(frame_bytes).decode("utf-8")
    try:
        msg = client.messages.create(
            model="claude-haiku-4-5-20251001",
            max_tokens=150,
            messages=[{
                "role": "user",
                "content": [
                    {
                        "type": "image",
                        "source": {
                            "type": "base64",
                            "media_type": "image/jpeg",
                            "data": b64
                        }
                    },
                    {
                        "type": "text",
                        "text": VISION_PROMPT
                    }
                ]
            }]
        )
        return msg.content[0].text.strip()
    except Exception as e:
        return f"Vision error: {e}"


# ── Main loop ─────────────────────────────────────────────────────────────────

def run():
    print("[legacy-eyes] Starting 1fps vision loop...")
    print(f"[legacy-eyes] Source: {VISION_SOURCE}")
    print(f"[legacy-eyes] Workspace: {WORKSPACE_DIR}")
    print(f"[legacy-eyes] Vision API call every {DIGEST_INTERVAL}s")

    state = read_ambient_state()
    frame_count = 0

    while True:
        loop_start = time.time()

        frame_bytes = capture_frame(VISION_SOURCE)

        if frame_bytes and (frame_count % DIGEST_INTERVAL == 0):
            digest = get_vision_digest(frame_bytes)

            if digest == "SCENE_UNCHANGED":
                state["scene_unchanged_count"] = state.get("scene_unchanged_count", 0) + 1
            else:
                state["scene_unchanged_count"] = 0
                state["vision_digest"] = digest
                print(f"[legacy-eyes] {digest[:100]}")

            write_ambient_state(state)

        elif not frame_bytes:
            print(f"[legacy-eyes] Warning: failed to capture frame from {VISION_SOURCE}")

        frame_count += 1
        elapsed = time.time() - loop_start
        time.sleep(max(0, CAPTURE_INTERVAL - elapsed))


if __name__ == "__main__":
    run()
