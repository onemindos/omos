---
type: integration-plan
category: video-automation
priority: deferred
created: 2026-05-24
tags:
  - obs
  - nats
  - video
  - automation
  - real-time
---

# OBS Studio + NATS Integration Plan

**Status:** Deferred until video automation needed  
**Purpose:** Real-time video control via NATS pub/sub → OBS WebSocket API

---

## Overview

**OBS WebSocket API** = One-time action commands (not continuous monitoring)  
**NATS Integration** = Event-driven triggers for automated video production

### Architecture

```
NATS Topic (e.g., video.obs.control)
  ↓
NATS Subscriber (Python/Node.js)
  ↓
OBS WebSocket API (localhost:4455)
  ↓
OBS Studio (Mac M2 Max)
```

---

## Use Cases

### 1. Automated Meeting Recording
```
Meeting starts (Calendar event) → NATS publish
  ↓
NATS subscriber receives event
  ↓
OBS WebSocket → Start recording
  ↓
Meeting ends → NATS publish
  ↓
OBS WebSocket → Stop recording
  ↓
Video saved to ~/Videos/meetings/2026-05-24.mp4
```

### 2. Drone Mission Recording
```
Drone takeoff → NATS publish
  ↓
OBS switches to "Drone Feed" scene
  ↓
OBS starts recording
  ↓
Drone lands → NATS publish
  ↓
OBS stops recording + switches back
```

### 3. Live Stream Automation
```
7pm ET cron → NATS publish
  ↓
OBS starts streaming to YouTube
  ↓
8pm ET → NATS publish
  ↓
OBS stops streaming
```

### 4. Emergency Recording
```
Security alert → NATS publish (priority:high)
  ↓
OBS immediately starts recording all cameras
  ↓
Alert cleared → NATS publish
  ↓
OBS stops + saves evidence
```

---

## NATS Topics (Proposed)

| **Topic** | **Payload** | **Action** |
|-----------|-------------|------------|
| `video.obs.record.start` | `{"scene": "Screen", "path": "/tmp/video.mp4"}` | Start recording |
| `video.obs.record.stop` | `{}` | Stop recording |
| `video.obs.stream.start` | `{"platform": "youtube"}` | Start streaming |
| `video.obs.stream.stop` | `{}` | Stop streaming |
| `video.obs.scene.switch` | `{"scene": "Webcam"}` | Switch scene |
| `video.obs.source.toggle` | `{"source": "Microphone", "visible": false}` | Toggle source |
| `video.obs.screenshot` | `{"path": "/tmp/screenshot.png"}` | Take screenshot |
| `video.obs.status.get` | `{}` | Get current status (recording/streaming) |

---

## Implementation

### Phase 1: OBS Setup
```bash
# 1. Install OBS Studio
brew install --cask obs

# 2. Enable WebSocket in OBS
# Tools → WebSocket Server Settings
# Port: 4455
# Password: <set secure password>

# 3. Install Python library
pip install obs-websocket-py
```

### Phase 2: NATS Subscriber (Python)

**File:** `~/.nats/subscribers/obs-control.py`

```python
#!/usr/bin/env python3
import asyncio
import nats
from obs_websocket import obsws, requests

# OBS connection
obs = obsws("localhost", 4455, "your-password")
obs.connect()

# NATS connection
async def main():
    nc = await nats.connect("nats://localhost:4222")
    
    async def obs_handler(msg):
        subject = msg.subject
        data = msg.data.decode()
        
        print(f"📹 Received: {subject} → {data}")
        
        # Route commands
        if subject == "video.obs.record.start":
            obs.call(requests.StartRecord())
            print("🔴 Recording started")
            
        elif subject == "video.obs.record.stop":
            obs.call(requests.StopRecord())
            print("✅ Recording stopped")
            
        elif subject == "video.obs.stream.start":
            obs.call(requests.StartStreaming())
            print("🔴 Streaming started")
            
        elif subject == "video.obs.stream.stop":
            obs.call(requests.StopStreaming())
            print("✅ Streaming stopped")
            
        elif subject.startswith("video.obs.scene.switch"):
            import json
            payload = json.loads(data)
            scene = payload.get("scene")
            obs.call(requests.SetCurrentScene(scene))
            print(f"🎬 Switched to scene: {scene}")
    
    # Subscribe to all OBS topics
    await nc.subscribe("video.obs.>", cb=obs_handler)
    
    print("✓ OBS NATS subscriber ready")
    await asyncio.Event().wait()

if __name__ == "__main__":
    asyncio.run(main())
```

### Phase 3: Test Commands

```bash
# Publish test messages (requires nats CLI)
nats pub video.obs.record.start '{}'
# OBS should start recording

nats pub video.obs.record.stop '{}'
# OBS should stop recording

nats pub video.obs.scene.switch '{"scene": "Webcam"}'
# OBS should switch to webcam scene
```

### Phase 4: launchd Service (Auto-start on boot)

**File:** `~/Library/LaunchAgents/com.onemind.obs-nats.plist`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.onemind.obs-nats</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/python3</string>
        <string>/Users/onemind/.nats/subscribers/obs-control.py</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/Users/onemind/.nats/logs/obs-nats-error.log</string>
    <key>StandardOutPath</key>
    <string>/Users/onemind/.nats/logs/obs-nats-out.log</string>
</dict>
</plist>
```

Load it:
```bash
launchctl load ~/Library/LaunchAgents/com.onemind.obs-nats.plist
launchctl start com.onemind.obs-nats
```

---

## OBS WebSocket Commands Reference

### Recording
```python
from obs_websocket import obsws, requests

ws = obsws("localhost", 4455, "password")
ws.connect()

# Start recording
ws.call(requests.StartRecord())

# Stop recording
ws.call(requests.StopRecord())

# Pause recording
ws.call(requests.PauseRecord())

# Resume recording
ws.call(requests.ResumeRecord())
```

### Streaming
```python
# Start streaming
ws.call(requests.StartStreaming())

# Stop streaming
ws.call(requests.StopStreaming())
```

### Scenes
```python
# Switch scene
ws.call(requests.SetCurrentScene("Webcam"))

# Get current scene
scene = ws.call(requests.GetCurrentScene())
print(scene.getName())

# List all scenes
scenes = ws.call(requests.GetSceneList())
for scene in scenes.getScenes():
    print(scene['name'])
```

### Sources
```python
# Toggle source visibility
ws.call(requests.SetSourceVisibility("Microphone", False))  # Hide
ws.call(requests.SetSourceVisibility("Microphone", True))   # Show

# Adjust volume
ws.call(requests.SetVolume("Microphone", 0.5))  # 50%
```

### Status
```python
# Get recording status
status = ws.call(requests.GetStreamingStatus())
print(f"Recording: {status.getRecording()}")
print(f"Streaming: {status.getStreaming()}")
```

### Screenshot
```python
# Take screenshot of source
ws.call(requests.TakeSourceScreenshot(
    sourceName="Screen Capture",
    embedPictureFormat="png",
    saveToFilePath="/tmp/screenshot.png"
))
```

---

## Security Considerations

### 1. OBS WebSocket Password
- Store in 1Password: `OBS WebSocket Password`
- Never commit to git
- Use environment variable: `OBS_WEBSOCKET_PASSWORD`

### 2. NATS Authentication
- Secure NATS with credentials
- Limit topic subscriptions (only `video.obs.>`)

### 3. Access Control
- Only OneMind agents can publish to `video.obs.>` topics
- Add NATS subject permissions in `nats-server.conf`

---

## Installation Checklist

When ready to deploy:

- [ ] Install OBS Studio (`brew install --cask obs`)
- [ ] Enable OBS WebSocket (port 4455)
- [ ] Set OBS WebSocket password (store in 1Password)
- [ ] Install `obs-websocket-py` (`pip install obs-websocket-py`)
- [ ] Create `~/.nats/subscribers/obs-control.py`
- [ ] Test manual OBS WebSocket connection
- [ ] Test NATS pub → OBS action
- [ ] Create launchd service for auto-start
- [ ] Add to Homepage dashboard (OBS status widget)
- [ ] Document OBS scenes/sources in vault

---

## Next Steps

### Immediate (When Video Automation Needed)
1. **Define use case:** What videos do we need to automate?
2. **Install OBS + WebSocket**
3. **Create NATS subscriber**
4. **Test end-to-end**

### Future Enhancements
- **AI-driven scene switching:** YOLOv8 detects person → switch to webcam
- **Multi-camera sync:** 3+ drone feeds → OBS multi-view
- **Auto-upload:** Recording done → upload to YouTube/S3
- **Real-time alerts:** Recording failed → NATS publish → notify Zeus

---

## Related Files
- **OBS MCP Evaluation:** `obs-studio-mcp-evaluation.md` (WebSocket API docs)
- **Master Installation Guide:** `installation-guide-master.md` (Phase 4: Optional Tools)
- **NATS Infrastructure:** `00-24 SO/03 Infrastructure (The Body)/NATS-INFRASTRUCTURE.md`

---

**Status:** Ready to implement when video automation becomes priority ★
