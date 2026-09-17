---
title: "OBS Studio Automation Research"
type: research
domain: SO/03 Infrastructure
topic: video, streaming, automation, mcp
status: reference
created: 2026-05-24
updated: 2026-05-26
tags: [type/research, domain/so, topic/video, topic/obs, topic/automation]
---

# OBS Studio Automation Research

**OBS Studio** = Open Broadcaster Software. Free, open-source video recording and live streaming.

**Status:** No native MCP server exists (May 2026). Automation available via WebSocket API.

---

## What OBS Can Do

- Multi-source compositing (cameras, windows, images, text overlays)
- Scene management (switch layouts mid-stream)
- Streaming (YouTube, Twitch, custom RTMP)
- Recording (MP4, MKV, MOV — local)
- Virtual camera (use OBS as webcam in Zoom/Teams)
- Filters (chroma key, color correction, audio processing)

---

## OneMind Use Cases

| Scenario | Trigger | Action |
|----------|---------|--------|
| Tutorial recording | Hermes command | Start screen + mic capture |
| Live stream | Cron job (scheduled) | Stream to YouTube Live |
| Daily progress timelapse | Cron 11pm | Record 8h screen, 16x speed |
| Drone footage | NATS event | Capture RTSP streams from drones |
| Presentation | Hermes command | Activate scene: screen + webcam overlay |

---

## Automation Methods

### Method 1 — OBS WebSocket API (Recommended)

Built into OBS since v28. No plugin needed.

**Enable:** OBS → Tools → WebSocket Server Settings → Enable

```bash
# Install Python client
pip install obs-websocket-py
```

```python
from obs_websocket import obsws, requests

ws = obsws("localhost", 4455, "your-password")
ws.connect()

ws.call(requests.StartRecord())
# ... recording happens ...
ws.call(requests.StopRecord())

ws.disconnect()
```

**Capabilities:** Start/stop recording, start/stop streaming, switch scenes, show/hide sources, adjust audio, take screenshots, control filters.

---

### Method 2 — AppleScript Fallback (Simple)

```bash
# Start OBS
open -a "OBS"

# Trigger recording via hotkey (must be configured in OBS first)
osascript -e 'tell application "System Events" to keystroke "r" using {command down, shift down}'
```

Less reliable than WebSocket — use only if WebSocket setup not worth the overhead.

---

## Proposed MCP Server Design

If OBS control becomes frequent enough to warrant MCP integration, here's the tool schema:

```
obs_start_recording   — params: scene (optional), output_path (optional)
obs_stop_recording    — returns: file_path
obs_start_streaming   — params: service (youtube/twitch/custom)
obs_stop_streaming
obs_switch_scene      — params: scene_name
obs_toggle_source     — params: source_name, visible (bool)
obs_get_status        — returns: recording, streaming, current_scene, output_path
obs_create_scene      — params: name, sources[]
```

Architecture:
```
Hermes MCP Client
  ↓
OBS MCP Server (Node.js or Python)
  ↓
OBS WebSocket API (localhost:4455)
  ↓
OBS Studio (Mac)
```

Note: Hermes already has an `obs-nats-integration` skill — check that before building a custom MCP server.

---

## Decision Tree

```
Need automated video recording/streaming?
├── YES → Install OBS + configure WebSocket
│         Need Hermes to trigger it? → Write obs-control skill (Python)
│         Need cross-agent access?   → Build MCP server (optional)
└── NO  → Skip. Use macOS built-in (Cmd+Shift+5) for ad-hoc recording.
```

**Build MCP server when:**
- Regular content creation (tutorials, live streams)
- Automated daily recordings
- Multi-camera drone productions

**Skip when:**
- No video production planned
- Manual OBS control is sufficient

---

## Installation

```bash
# Install OBS
brew install --cask obs

# WebSocket is built-in since OBS 28 — just enable in settings
# Test connection
python3 -c "from obs_websocket import obsws; print('client ready')"
```

---

## Action Items

- [ ] Determine if video automation is needed (content creation decision)
- [ ] If yes: install OBS, enable WebSocket, write obs-control Hermes skill
- [ ] Check existing `obs-nats-integration` skill before building anything new
- [ ] MCP server only if WebSocket skill proves insufficient

---

## Related

- Hermes skill: `obs-nats-integration`
- `SO/08 Research/lpr-alpr-research.md` — camera feeds for OBS sources
- `GE/75 Zeus Delacruz (Brand & Empire)/Content/` — content calendar
