# Hermes Architecture Plan — Legacy Agent System

> **Created:** 2026-05-19
> **Purpose:** Definitive architecture for Zeus's single-agent system running all physical + digital infrastructure
> **Primary Agent:** Legacy (Hermes)
> **Interface Layer:** OpenClaw (voice/display/nodes)

---

## Question 1: Why One Process Matters (TypeScript vs Python)

### The Core Difference

**Hermes (Python) — Single async event loop:**
```
Voice in → STT → Agent reasons → import nats → publish command → drone moves
                                → import frame_sdk → take photo → vision analysis
                                → import mediapipe → detect gesture → robot command
                                → import mavsdk → direct drone control
                                → import meshtastic → LoRa broadcast
```
All in ONE `asyncio` event loop. No IPC. No serialization. No bridge processes. Shared memory between all plugins.

**OpenClaw (TypeScript) — Requires bridges for every hardware SDK:**
```
Voice in → STT → Agent reasons → spawn python subprocess → wait for stdout → parse JSON
                                → MCP server (separate process) → IPC round-trip → result
                                → exec('python3 script.py') → parse output
                                → another MCP server → another IPC round-trip
```

### Why This Can't Be Done in TypeScript

| SDK | Python | TypeScript Equivalent |
|---|---|---|
| `frame-sdk` (Brilliant Labs) | `pip install frame-sdk` → direct import | Does not exist |
| `mediapipe` | `pip install mediapipe` → direct import | `@mediapipe/tasks-vision` (browser-only, limited) |
| `mavsdk` | `pip install mavsdk` → direct import | No maintained package |
| `meshtastic` | `pip install meshtastic` → direct import | Does not exist |
| `ultralytics` (YOLO) | `pip install ultralytics` → direct import | Does not exist |
| `nats` | `pip install nats-py` → async native | `nats` npm ✅ (this one works) |
| `opencv` | `pip install opencv-python` → direct import | `opencv.js` (limited WASM port) |
| `dmxpy` | `pip install dmxpy` → direct import | `dmx` npm (exists but unmaintained) |

**Result:** In TypeScript, you'd maintain 6+ separate Python processes (MCP servers or scripts) alongside your agent. Each tool call crosses a process boundary. State can't be shared. Errors are harder to trace. The "one process" advantage is real and compounds with every hardware SDK you add.

### What This Means Practically

**Hermes (one process, shared state):**
```python
# Frame takes photo → MediaPipe detects gesture → NATS sends command
# All in-memory, zero serialization, <10ms between steps

photo = await frame.camera.take_photo()  # BLE read, ~200ms
gesture = hands.process(cv2.imdecode(photo))  # CPU, ~50ms  
await nc.publish("fleet.cmd.robot", json.dumps({"action": gesture}))  # TCP, ~5ms
# Total: ~255ms
```

**OpenClaw (multi-process, serialized):**
```typescript
// Frame MCP server → MediaPipe MCP server → NATS publish
// Each crosses process boundary, JSON serialization, ~100ms overhead per hop

const photo = await mcp.call("frame-server", "take_photo", {});  // IPC ~100ms + BLE ~200ms
const gesture = await mcp.call("mediapipe-server", "detect", {image: photo});  // IPC ~100ms + CPU ~50ms
await nats.publish("fleet.cmd.robot", JSON.stringify({action: gesture}));  // TCP ~5ms
// Total: ~555ms (2x slower, 3 processes to maintain)
```

**For real-time drone/robot control, that 2x difference matters.**

---

## Question 2: Agent Architecture — Single Supercharged Agent

### Decision: ONE Agent (Legacy) + Delegation, NOT Multi-Agent

**Why NOT port OpenClaw's multi-agent system:**
- OpenClaw has 12 agents because it can't change persona mid-conversation (rigid routing)
- Hermes has profiles (same agent, different toolset/persona per context)
- Multi-agent = you have to know which agent to talk to
- You explicitly said: "I don't want to talk to many different agents"

**Why ONE agent with delegation:**
- You talk to Legacy. Period.
- Legacy decides when to spawn background workers (delegation)
- Legacy has ALL toolsets available, switches context naturally
- Delegation handles parallel work WITHOUT you managing it

### The Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│  LEGACY (Primary Hermes Agent)                                       │
│  Model: Claude Sonnet 4 (fast) / Claude Opus 4 (complex)            │
│  Toolsets: ALL (farm, drone, robot, home, voice, code, web, files)  │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  PROFILES (Context Switches — Same Agent, Different Mode)     │   │
│  │                                                                │   │
│  │  farm-ops    → farm toolset + weather + irrigation context     │   │
│  │  security    → camera + drone + alert toolset                  │   │
│  │  dev         → code + terminal + git toolset                   │   │
│  │  home        → HomeAssistant + DMX + IoT toolset               │   │
│  │  fleet       → NATS + drone + robot toolset                    │   │
│  │  research    → web + file + analysis toolset                   │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  PLUGINS (Always Active — Hardware Access Layer)               │   │
│  │                                                                │   │
│  │  frame-bridge     → Brilliant Labs glasses (camera/mic/display)│   │
│  │  nats-fleet       → All edge device communication              │   │
│  │  meshtastic-mesh  → Off-grid LoRa backup                      │   │
│  │  dmx-lights       → Lighting automation                        │   │
│  │  drone-direct     → MAVSDK simple ops                          │   │
│  │  vision-detect    → YOLO object detection                      │   │
│  │  mediapipe-vision → Gesture/pose detection                     │   │
│  │  gesture-to-robot → Combo pipeline (Frame→Gesture→Command)     │   │
│  │  livekit-voice    → Real-time voice rooms (future)             │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  MCP SERVERS (External Capabilities)                           │   │
│  │                                                                │   │
│  │  dimos-mcp        → DimOS drone runtime (complex flight ops)   │   │
│  │  homeassistant    → IoT device control                         │   │
│  │  filesystem       → Extended file operations                   │   │
│  │  github           → Repo management                            │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  DELEGATION (Background Workers — Agent Spawns as Needed)      │   │
│  │                                                                │   │
│  │  "Research best spray schedule for wheat" → leaf worker        │   │
│  │  "Deploy new version of farm-dashboard"   → leaf worker        │   │
│  │  "Analyze all drone footage from today"   → leaf worker        │   │
│  │  "Run morning report + backup + sync"     → 3 parallel leaves  │   │
│  │                                                                │   │
│  │  Config:                                                       │   │
│  │    max_concurrent_children: 5                                  │   │
│  │    max_spawn_depth: 2 (orchestrator can spawn orchestrators)   │   │
│  │    child_timeout_seconds: 900                                  │   │
│  │    inherit_mcp_toolsets: true                                  │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  CRON (Scheduled Automations — Run Without You)                │   │
│  │                                                                │   │
│  │  06:00  morning-report (RPC script: weather+irrigation+status) │   │
│  │  06:30  drone-patrol-north (DimOS MCP: patrol waypoints)       │   │
│  │  12:00  livestock-count (drone photo → YOLO count → alert)     │   │
│  │  17:00  evening-summary (delegate: compile day's events)       │   │
│  │  22:00  security-sweep (drone + cameras + fence check)         │   │
│  │  02:00  system-backup (RPC script: backup all services)        │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  RPC SCRIPTS (Lobster Replacements — Zero-Context Pipelines)   │   │
│  │                                                                │   │
│  │  morning-farm-report.py   → weather+telemetry+irrigation       │   │
│  │  drone-patrol.py          → takeoff→waypoints→photos→land      │   │
│  │  backup-all.py            → db dump+compress+upload+notify     │   │
│  │  deploy-service.py        → git pull+build+restart+healthcheck │   │
│  │  livestock-audit.py       → drone photo→YOLO→count→compare     │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### How Concurrent Sessions Work

Hermes gateway handles this natively:

```
Telegram message from you     ─┐
Discord VC audio from you      ─┤── ALL route to same Legacy agent
WhatsApp message from you      ─┤   (single conversation context)
Cron trigger (6:00 AM)         ─┤
Webhook from HomeAssistant     ─┘

The agent processes ONE thing at a time in its main loop.
Delegation handles parallelism (background workers).
Cron tasks get their own isolated sessions (don't pollute your chat).
```

**No limitations on what it can do** because:
- All Python SDKs available via plugins
- All external services available via MCP
- Delegation for parallel work
- RPC scripts for multi-step pipelines
- Cron for scheduled operations
- Gateway for multi-platform input

### Config Changes to Supercharge

```yaml
# ~/.hermes/config.yaml additions

delegation:
  max_concurrent_children: 5      # up from 3
  max_spawn_depth: 2              # allow nested orchestration
  child_timeout_seconds: 900      # 15 min per child task
  orchestrator_enabled: true
  inherit_mcp_toolsets: true

code_execution:
  mode: project                   # use your real venv + tools
  timeout: 600                    # 10 min for complex scripts
  max_tool_calls: 100             # up from 50

model:
  default: claude-sonnet-4-20250514    # fast for most things
  thinking: claude-opus-4-20250514     # complex reasoning when needed
```

---

## Question 3: Hosting — Best Setup for Your Stack

### Decision: Local Primary + Modal for Heavy Compute

**Why NOT purely serverless (Modal/Daytona):**
- Your agent needs BLE access to Frame glasses (local only)
- Your agent needs LAN access to Meshtastic device (local only)  
- Your agent needs low-latency NATS to edge devices (local = lowest latency)
- Your agent needs USB DMX (local only)
- Serverless = cold start delays (2-5s) → kills real-time voice/drone response

**Why NOT purely cloud VPS:**
- You already have an M2 Max Mac (powerful enough)
- Tailscale gives you global access anyway
- Local = zero egress costs for video/audio processing

### The Hybrid Architecture

```
┌─────────────────────────────────────────────────────────┐
│  LOCAL (Mac Mini / M2 Max — Primary Runtime)            │
│                                                          │
│  Hermes Agent (LaunchAgent, always running)              │
│  ├── Gateway: Telegram + Discord + WhatsApp + Signal     │
│  ├── Plugins: frame-bridge, nats-fleet, meshtastic,      │
│  │            dmx-lights, drone-direct, vision-detect     │
│  ├── MCP: DimOS, HomeAssistant                           │
│  ├── NATS client (fleet telemetry subscriber)            │
│  ├── STT: local faster-whisper (GPU accelerated on M2)   │
│  └── Hardware: BLE, USB, LAN access                      │
│                                                          │
│  Why local:                                              │
│  - Hardware access (BLE, USB, LAN)                       │
│  - Zero cold start (always warm)                         │
│  - M2 Max = plenty of compute for ML inference           │
│  - Free (no per-minute costs)                            │
│  - Tailscale = globally accessible anyway                │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  MODAL (Serverless GPU — Heavy Compute Tasks)            │
│                                                          │
│  Used via delegation or RPC scripts for:                 │
│  ├── Large YOLO inference (batch processing farm photos) │
│  ├── Video analysis (drone footage processing)           │
│  ├── Training/fine-tuning (custom models)                │
│  ├── Batch STT (transcribe hours of audio)               │
│  └── Heavy data processing (satellite imagery, etc.)     │
│                                                          │
│  Config:                                                 │
│  terminal_backends:                                      │
│    modal:                                                │
│      gpu: A10G                                           │
│      timeout: 900                                        │
│      image: "python:3.12-cuda"                           │
│                                                          │
│  Why Modal:                                              │
│  - Pay per second (hibernates when idle → $0)            │
│  - GPU access for ML inference                           │
│  - Auto-scales for batch jobs                            │
│  - Agent delegates heavy work here, keeps local snappy   │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  VPS (onemindos-cloud-1, 100.98.25.10 — Infrastructure) │
│                                                          │
│  ├── NATS server (fleet coordination hub)                │
│  ├── Backup storage                                      │
│  ├── Public webhook endpoints                            │
│  └── Fallback agent (if Mac is offline)                  │
│                                                          │
│  Why VPS:                                                │
│  - NATS needs a fixed address for all devices            │
│  - Webhooks need public endpoint                         │
│  - Redundancy if your Mac dies                           │
└─────────────────────────────────────────────────────────┘
```

### Hermes Terminal Backend Config

```yaml
# ~/.hermes/config.yaml
terminal:
  default: local          # primary = local Mac
  backends:
    local:
      shell: zsh
    ssh:
      host: 100.98.25.10  # VPS via Tailscale
      user: zeus
    modal:
      gpu: A10G
      timeout: 900
      packages: ["ultralytics", "opencv-python", "torch"]
```

When the agent needs GPU, it uses `terminal(command="...", backend="modal")`. Otherwise everything runs local.

---

## Question 4: One Agent to Rule Them All + OpenClaw's Role

### The "Talk to One Agent" Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         YOU (Zeus)                                │
│                                                                   │
│   Voice: "Legacy, launch the security sweep"                     │
│   Text:  "What's the irrigation status?"                         │
│   Gesture: *points at barn* (Frame → gesture → command)          │
└──────────────────────────────┬────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│  OPENCLAW (Interface Layer — Voice + Display + Nodes)            │
│                                                                   │
│  What it does:                                                    │
│  ├── Talk mode (always-listening, wake word "Legacy")             │
│  ├── Mac node (desktop dashboard, Canvas rendering)               │
│  ├── iOS node (mobile interface, A2UI)                            │
│  ├── Voice (Deepgram STT → text → forward to Hermes)             │
│  └── Display (render farm dashboard, drone feeds, alerts)         │
│                                                                   │
│  What it does NOT do:                                             │
│  ├── ❌ Reason about tasks                                       │
│  ├── ❌ Call hardware SDKs                                        │
│  ├── ❌ Manage fleet                                              │
│  └── ❌ Run cron jobs                                             │
│                                                                   │
│  It's a THIN RELAY:                                               │
│  Voice → transcribe → forward to Hermes → get response → TTS     │
│  Canvas request → forward to Hermes → get HTML → render on node   │
└──────────────────────────────┬────────────────────────────────────┘
                               │ Webhook / NATS / direct API
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│  HERMES / LEGACY (Brain — All Decision-Making)                   │
│                                                                   │
│  Receives:                                                        │
│  ├── Voice commands (from OpenClaw relay OR Telegram direct)      │
│  ├── Text commands (Telegram, Discord, WhatsApp)                  │
│  ├── Cron triggers (scheduled tasks)                              │
│  ├── Webhook events (HomeAssistant, GitHub, external)             │
│  ├── NATS alerts (fleet telemetry, device events)                 │
│  └── Frame input (photos, audio, tap gestures)                    │
│                                                                   │
│  Decides & Acts:                                                  │
│  ├── Direct tool calls (fast, simple operations)                  │
│  ├── RPC scripts (multi-step pipelines, zero-context)             │
│  ├── Delegation (heavy background work → spawns workers)          │
│  └── Response back to you (via whichever channel you used)        │
│                                                                   │
│  Background (runs without you asking):                            │
│  ├── Cron jobs firing on schedule                                 │
│  ├── NATS subscriptions processing alerts                         │
│  ├── Self-improvement (skill capture after complex tasks)         │
│  ├── Curator (7-day skill lifecycle maintenance)                  │
│  └── Delegated workers completing tasks asynchronously            │
└─────────────────────────────────────────────────────────────────┘
```

### How OpenClaw Fits (Thin Relay, Not a Brain)

**Keep OpenClaw running as a process.** But strip its role to INTERFACE ONLY:

```
OpenClaw's new job:
1. Listen for wake word ("Legacy")
2. Transcribe what you say
3. Forward text to Hermes via webhook/NATS
4. Receive Hermes response
5. TTS the response back to you
6. Render dashboards on Canvas when Hermes says to

That's it. No reasoning. No tool calling. No task management.
```

**Implementation — OpenClaw as relay:**

```typescript
// OpenClaw plugin: hermes-relay
// Intercepts all messages and forwards to Hermes, returns Hermes response

export default {
  name: 'hermes-relay',
  hooks: {
    async onMessage(message, context) {
      // Forward to Hermes gateway via HTTP
      const response = await fetch('http://localhost:4200/api/message', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          text: message.text,
          platform: 'openclaw-relay',
          user: 'zeus',
          channel: context.channel
        })
      });
      
      const result = await response.json();
      
      // Display result via OpenClaw's native rendering
      if (result.canvas_html) {
        await context.canvas.present(result.canvas_html);
      }
      
      return result.text;  // TTS this back
    }
  }
};
```

**Or via NATS (cleaner):**

```typescript
// OpenClaw subscribes to "legacy.response.openclaw"
// OpenClaw publishes to "legacy.input.voice"

nc.subscribe('legacy.response.openclaw', (msg) => {
  const response = JSON.parse(msg.data);
  if (response.tts) speak(response.text);
  if (response.canvas) canvas.present(response.canvas);
});

// When user speaks:
nc.publish('legacy.input.voice', JSON.stringify({
  text: transcription,
  source: 'openclaw-talk-mode',
  timestamp: Date.now()
}));
```

### Why Not Just Kill OpenClaw?

Because Hermes CANNOT do these things:
- **Always-listening wake word detection** (OpenClaw's Talk mode with Deepgram streaming)
- **Native iOS/Mac canvas rendering** (A2UI, WebView on device)
- **Desktop node integration** (macOS menu bar, notifications)
- **Phone camera/screen access** (phone-control plugin)

These are DISPLAY/INPUT capabilities. Hermes is headless by design. OpenClaw gives it a face.

### The "No Agent Spawn" Concern

You said "we don't have agent spawn like OpenClaw." But Hermes HAS delegation:

| OpenClaw | Hermes Equivalent |
|---|---|
| 12 named agents (zeus, coder, researcher...) | 1 agent + profiles (switch context, not identity) |
| Agent-to-agent messaging | `delegate_task` (parent-child, not peer-to-peer) |
| Subagent spawning (max 12) | `delegate_task(tasks=[...])` batch mode (max 5 concurrent) |
| Each agent has own session | Each delegated task has isolated session |
| You choose which agent to talk to | Legacy decides internally whether to delegate |

**The key difference:** In OpenClaw, YOU route to the right agent. In Hermes, LEGACY routes internally. You never see or manage the workers.

### What Happens When You Say "Run the Security Sweep"

```
You: "Legacy, run the security sweep"

Legacy's internal reasoning:
  "Security sweep requires: drone patrol, camera check, fence sensors, compile report.
   This is a multi-step operation. I'll delegate the heavy parts and coordinate."

Legacy executes:
  1. delegate_task(
       tasks=[
         {goal: "Run drone patrol over property perimeter, capture photos every 100m",
          toolsets: ["drone-direct", "nats-fleet", "vision-detect"]},
         {goal: "Check all security camera feeds for anomalies in last hour",
          toolsets: ["homeassistant", "vision-detect"]},
         {goal: "Query fence sensor status from NATS telemetry",
          toolsets: ["nats-fleet"]}
       ]
     )
  2. Workers run in parallel (you don't see them)
  3. Results return to Legacy
  4. Legacy compiles: "Security sweep complete. All clear. Drone spotted 2 deer 
     near north fence. No anomalies on cameras. All 12 fence sensors reporting normal."
  5. Sends to you via Telegram + renders map on OpenClaw Canvas

You never talked to "drone-agent" or "security-agent" or "camera-agent."
You talked to Legacy. Legacy handled it.
```

---

## RPC Scripts — Lobster Workflow Replacements

### What "Zero-Context-Cost" Means

Normal agent conversation:
```
Turn 1: LLM sees full history (10KB) + decides to search    = 10KB processed
Turn 2: LLM sees full history (12KB) + decides to extract   = 12KB processed
Turn 3: LLM sees full history (15KB) + decides to save      = 15KB processed
Turn 4: LLM sees full history (16KB) + summarizes           = 16KB processed
Total: 53KB of context processed across 4 turns
```

RPC script:
```
Turn 1: LLM sees history (10KB) + writes script             = 10KB processed
Turn 2: LLM sees history (11KB) + reads script result       = 11KB processed
Total: 21KB of context processed across 2 turns (60% savings)
```

The 3 tool calls inside the script cost ZERO context because they never enter the LLM conversation. They run as plain Python function calls over a socket.

### How RPC Scripts Work (Mechanical Detail)

```
┌─────────────────────────────────────────────────────────┐
│  Hermes Agent (Parent Process)                           │
│                                                          │
│  1. Agent decides to execute a script                    │
│  2. Generates hermes_tools.py (stub module):             │
│     def web_search(query): → UDS call → parent          │
│     def terminal(command): → UDS call → parent          │
│     def read_file(path):   → UDS call → parent          │
│                                                          │
│  3. Opens Unix Domain Socket (UDS)                       │
│  4. Spawns child Python process with:                    │
│     - HERMES_RPC_SOCKET=/tmp/hermes_rpc_abc123.sock      │
│     - hermes_tools.py in PYTHONPATH                      │
│     - NO API keys, NO tokens (scrubbed env)              │
│                                                          │
│  5. Child runs script:                                   │
│     from hermes_tools import web_search, terminal        │
│     result = web_search("weather today")  # → UDS → parent dispatches → returns │
│     terminal("echo done")                 # → UDS → parent dispatches → returns │
│                                                          │
│  6. Parent collects stdout + return value                │
│  7. Returns to agent as single tool result               │
└─────────────────────────────────────────────────────────┘
```

### For Remote Execution (Modal/Docker/SSH):

```
Same concept but file-based RPC:
- Parent ships script + hermes_tools.py to remote
- Script writes req_000001 (JSON request) to HERMES_RPC_DIR
- Parent polls for request files, dispatches locally, writes res_000001
- Script reads response file, continues
- Works across any network boundary
```

### Example: Morning Farm Report (Was Lobster Workflow)

```python
#!/usr/bin/env python3
"""Morning farm report — runs via cron at 06:00 daily.
Calls 5 tools via RPC, costs 0 LLM context tokens."""

from hermes_tools import web_search, terminal, read_file, write_file
from datetime import date

# 1. Weather forecast
weather = web_search(query=f"weather forecast {date.today()} Scottsdale AZ agriculture")
weather_summary = weather[0]["snippet"] if weather else "Weather unavailable"

# 2. Irrigation system status
irrigation = terminal(command="curl -s http://192.168.1.100/api/status | jq '.zones'")

# 3. Drone fleet battery levels
fleet = terminal(command="""python3 -c "
import asyncio, json, nats
async def check():
    nc = await nats.connect('nats://100.98.25.10:4222')
    resp = await nc.request('fleet.status.all', b'battery', timeout=3)
    print(resp.data.decode())
    await nc.close()
asyncio.run(check())
" """)

# 4. Yesterday's activity log
activity = read_file(path="/var/logs/hermes/yesterday-summary.txt")

# 5. Compile and save report
report = f"""# Morning Farm Report — {date.today()}

## Weather
{weather_summary}

## Irrigation
{irrigation}

## Fleet Status
{fleet}

## Yesterday's Activity
{activity[:500]}

---
Generated automatically by Legacy at 06:00
"""

write_file(path=f"/tmp/farm-report-{date.today()}.md", content=report)
print(f"Report saved. Summary: {weather_summary[:100]}")
```

### Example: Security Sweep (Was Lobster Workflow)

```python
#!/usr/bin/env python3
"""Security sweep — drone patrol + camera check + sensor audit."""

from hermes_tools import terminal, web_extract
import json

# 1. Launch drone patrol
terminal(command="""python3 -c "
import asyncio, nats, json
async def patrol():
    nc = await nats.connect('nats://100.98.25.10:4222')
    waypoints = [
        {'lat': 33.4500, 'lon': -112.0700, 'alt': 30},
        {'lat': 33.4510, 'lon': -112.0700, 'alt': 30},
        {'lat': 33.4510, 'lon': -112.0690, 'alt': 30},
        {'lat': 33.4500, 'lon': -112.0690, 'alt': 30},
    ]
    await nc.publish('fleet.cmd.drone-1', json.dumps({
        'action': 'patrol',
        'waypoints': waypoints,
        'photo_interval': 100
    }).encode())
    await nc.close()
asyncio.run(patrol())
" """)

# 2. Check security cameras via HomeAssistant
cameras = terminal(command="""curl -s -H "Authorization: Bearer $HA_TOKEN" \
    http://192.168.1.50:8123/api/states | \
    jq '[.[] | select(.entity_id | startswith("camera.")) | {id: .entity_id, state: .state}]'""")

# 3. Check fence sensors
sensors = terminal(command="""python3 -c "
import asyncio, nats, json
async def check():
    nc = await nats.connect('nats://100.98.25.10:4222')
    resp = await nc.request('fleet.status.fence-sensors', b'status', timeout=5)
    print(resp.data.decode())
    await nc.close()
asyncio.run(check())
" """)

# 4. Compile results
print(json.dumps({
    "drone_patrol": "launched",
    "cameras": cameras,
    "fence_sensors": sensors,
    "status": "sweep initiated"
}))
```

### Example: Deploy Service (Was Lobster Workflow)

```python
#!/usr/bin/env python3
"""Deploy a service — git pull, build, restart, healthcheck."""

from hermes_tools import terminal, write_file
import sys

service = sys.argv[1] if len(sys.argv) > 1 else "farm-dashboard"
repo_path = f"/opt/services/{service}"

# 1. Git pull
pull = terminal(command=f"cd {repo_path} && git pull origin main")
if "Already up to date" in pull:
    print(f"{service}: already up to date, skipping deploy")
    sys.exit(0)

# 2. Install deps + build
terminal(command=f"cd {repo_path} && npm install && npm run build")

# 3. Restart service
terminal(command=f"sudo systemctl restart {service}")

# 4. Wait and healthcheck
import time
time.sleep(3)
health = terminal(command=f"curl -sf http://localhost:3000/health || echo 'UNHEALTHY'")

if "UNHEALTHY" in health:
    # Rollback
    terminal(command=f"cd {repo_path} && git checkout HEAD~1 && npm run build")
    terminal(command=f"sudo systemctl restart {service}")
    print(f"DEPLOY FAILED: {service} — rolled back")
else:
    print(f"DEPLOY SUCCESS: {service} — healthy")
```

---

## Real-Time Voice + Vision Pipeline

### Current State (Works Today)

```
┌─────────────────────────────────────────────────────────────┐
│  OPTION A: Telegram Voice (6-9s latency)                     │
│                                                               │
│  You → hold mic → record → send .ogg → Hermes transcribes   │
│  → agent reasons → executes tools → TTS → voice bubble back  │
│                                                               │
│  Good for: non-urgent commands, queries, reports             │
│  Bad for: real-time drone control, emergency stops           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  OPTION B: Discord Voice Channel (3-4s latency)              │
│                                                               │
│  You → join VC on phone → speak freely (hands-free)          │
│  → 1.5s silence detection → STT → agent → TTS → plays back  │
│                                                               │
│  Good for: hands-free outdoor work, continuous conversation  │
│  Bad for: sub-second response needed                         │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  OPTION C: OpenClaw Talk Mode → Hermes Relay (<3s)           │
│                                                               │
│  OpenClaw always-listening → wake word "Legacy"              │
│  → Deepgram streaming STT (real-time) → text                │
│  → Forward to Hermes via NATS/webhook                       │
│  → Hermes executes → response back to OpenClaw              │
│  → OpenClaw TTS → speaker                                    │
│                                                               │
│  Good for: at your desk, ambient control                     │
│  Best latency because Deepgram streams (no record+send)      │
└─────────────────────────────────────────────────────────────┘
```

### Future State (Needs Plugin Development)

```
┌─────────────────────────────────────────────────────────────┐
│  OPTION D: LiveKit Voice Room (500ms-1s latency) [FUTURE]    │
│                                                               │
│  Phone app → LiveKit room → Hermes plugin subscribes         │
│  → Real-time audio stream → continuous STT                   │
│  → Agent processes → TTS published back to room              │
│  → You hear response instantly                               │
│                                                               │
│  Also enables:                                                │
│  → Drone video streams into same room                        │
│  → Robot cameras publishing to room                          │
│  → Multi-device audio (Frame glasses mic → room → agent)     │
│                                                               │
│  This is the ULTIMATE setup for outdoor ranch work:          │
│  Frame glasses mic → LiveKit → Hermes → drone command        │
│  Drone camera → LiveKit → Frame glasses display              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  OPTION E: OpenAI Realtime API (300-500ms) [FUTURE]          │
│                                                               │
│  WebSocket → OpenAI Realtime API → bidirectional audio       │
│  → Function calling happens mid-stream                       │
│  → Response starts before you finish speaking                │
│                                                               │
│  Lowest possible latency but:                                │
│  → Locked to OpenAI models                                   │
│  → More expensive (per-minute billing)                       │
│  → Less control over tool execution                          │
└─────────────────────────────────────────────────────────────┘
```

### Vision Pipeline (Frame Glasses → Agent → Action)

```
Brilliant Frame Glasses (on your face)
    │ BLE 5.3 → phone hub app
    │
    ▼
Hermes frame-bridge plugin
    │ await frame.camera.take_photo()
    │ await frame.microphone.record_audio()
    │
    ▼
Agent Processing
    │ Vision model: "I see a broken fence post at the north boundary"
    │ Decision: "Alert Zeus + schedule repair + mark location"
    │
    ▼
Actions
    ├── Send alert to Telegram: "Broken fence post detected at 33.451, -112.069"
    ├── Add to task list: "Repair fence post — north boundary"
    ├── Publish GPS marker: fleet.marker.fence-repair
    └── Show on Frame display: "⚠️ Fence damage logged"
```

---

## Complete System Diagram

```
═══════════════════════════════════════════════════════════════════
                        ZEUS (Human)
═══════════════════════════════════════════════════════════════════
        │              │              │              │
   Voice (Talk)   Telegram      Frame Glasses    Discord VC
        │              │              │              │
        ▼              ▼              ▼              ▼
┌──────────────┐  ┌─────────────────────────────────────────┐
│  OPENCLAW    │  │  HERMES / LEGACY (Brain)                 │
│  (Interface) │  │                                          │
│              │──│→ Gateway: Telegram, Discord, WhatsApp,   │
│  Talk mode   │  │           Signal, HomeAssistant          │
│  Mac node    │  │                                          │
│  iOS node    │  │  Plugins: frame-bridge, nats-fleet,      │
│  Canvas/A2UI │  │           mediapipe, drone-direct,       │
│              │  │           vision-detect, dmx-lights,     │
│  Relay only: │  │           meshtastic-mesh                │
│  voice→text  │  │                                          │
│  text→display│  │  MCP: DimOS, HomeAssistant, GitHub       │
│              │  │                                          │
│  Renders:    │  │  Delegation: up to 5 parallel workers    │
│  dashboards  │  │  RPC Scripts: zero-context pipelines     │
│  drone feeds │  │  Cron: 6+ scheduled daily automations    │
│  alerts      │  │                                          │
└──────────────┘  └────────────────┬────────────────────────┘
                                   │
                          NATS (over Tailscale)
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
        ▼                          ▼                          ▼
┌──────────────┐         ┌──────────────┐          ┌──────────────┐
│  DimOS       │         │  HomeAssistant│          │  Edge Devices│
│  (Jetson)    │         │  (RPi/NUC)   │          │              │
│              │         │              │          │  Meshtastic  │
│  Drone-1     │         │  Zigbee      │          │  Fence sensors│
│  Drone-2     │         │  Z-Wave      │          │  Weather stn │
│  Unitree G1  │         │  MQTT        │          │  Irrigation  │
│              │         │  ESPHome     │          │              │
└──────────────┘         └──────────────┘          └──────────────┘
```

---

## Summary: What Goes Where

| Component | Role | Why |
|---|---|---|
| **Hermes (local Mac)** | Brain — all reasoning, tool calling, delegation | Python SDKs, hardware access, zero cold start |
| **OpenClaw (local Mac)** | Face — voice input, display output, device nodes | Wake word, Canvas, A2UI, always-listening |
| **Modal (serverless)** | Muscle — heavy compute when needed | GPU, batch ML, pay-per-second |
| **VPS (Tailscale)** | Backbone — NATS hub, webhooks, redundancy | Fixed address for fleet, public endpoints |
| **DimOS (edge)** | Hands — real-time robotics control | 10-400 Hz loops, obstacle avoidance |
| **HomeAssistant (local)** | Nervous system — IoT protocol translation | Zigbee, Z-Wave, MQTT, ESPHome |
| **NATS (VPS)** | Blood vessels — message delivery everywhere | Pub/sub, request/reply, fleet coordination |
| **Tailscale (mesh)** | Skeleton — connects everything securely | Zero-config VPN, all devices reachable |

---

## n8n Comparison: Why RPC Scripts + Delegation Replace It

| n8n Workflow | Hermes Equivalent |
|---|---|
| Trigger: Webhook | Gateway platform event / `pre_gateway_dispatch` hook |
| Trigger: Cron | Hermes cron (`hermes cron add ...`) |
| Trigger: MQTT message | NATS subscription in plugin (fires on event) |
| Node: HTTP Request | `web_search` / `web_extract` in RPC script |
| Node: Execute Command | `terminal(command="...")` in RPC script |
| Node: IF condition | Python `if/else` in script |
| Node: Split in Batches | `delegate_task(tasks=[...])` batch mode |
| Node: Set variable | Python variable assignment |
| Node: Send Telegram | Agent responds to Telegram channel directly |
| Error handling: Retry | Agent reasons about error, adapts approach |
| Workflow versioning | Git (scripts are files in your repo) |
| Visual editor | None needed — agent WRITES the scripts for you |

**What Hermes does that n8n can't:**
- Agent can DECIDE mid-workflow to change the plan
- Agent can WRITE new workflows dynamically from natural language
- Agent can reason about failures and try alternative approaches
- Zero-context-cost means complex pipelines don't burn money
- Direct hardware access (no separate integration servers)

**What n8n does that you'd miss:**
- Visual monitoring dashboard (see workflow execution in real-time)
  - Hermes alternative: gateway SSE shows live agent activity
- OAuth credential management UI
  - Hermes alternative: `requires_env` in plugin.yaml + env vars
- Community workflow templates
  - Hermes alternative: skills marketplace (skills.sh)
