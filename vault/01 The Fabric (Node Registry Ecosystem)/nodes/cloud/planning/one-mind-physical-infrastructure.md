# One Mind OS — Physical Infrastructure Architecture
# Home Assistant + OpenClaw Node + Tailscale Per Location
**Created:** 2026-04-18
**Author:** Legacy ⚡
**Domain:** SO 02 Infrastructure (VPS, Servers & Networking)
**Status:** Architecture defined — implementation pending (office setup)

---

## Overview

One Mind OS physical infrastructure follows a **centralized intelligence, distributed execution** model. One VPS gateway runs all agents. Physical locations (home, office, studio) each get one always-on machine that acts as the local node — connecting the digital agent layer to the physical world.

**Core equation per location:**
```
Tailscale subnet router + OpenClaw node + Home Assistant + Frigate NVR
= Legacy has full eyes, ears, and hands in that space
```

---

## Why This Architecture Is Right

There is no better way. The alternatives:
- **Full gateway per location** — more cost, more complexity, no single brain
- **Cloud-only** — latency, privacy concerns, no local execution when internet is down
- **Device-by-device integration** — hundreds of custom integrations to maintain

**One gateway + many nodes** is:
- Cheaper (one VPS, cheap mini PCs per location)
- Simpler (one brain to maintain)
- Sovereign (local-first, cloud-orchestrated)
- Infinitely scalable (add a location = one machine + one command)

---

## Full Stack Per Location

```
Physical Space (Office / Home / Studio / Any Location)
│
├── Hardware Layer
│   ├── TVs, displays, projectors
│   ├── Cameras (Ring, Reolink, IP cameras)
│   ├── 3D printers (OctoPrint integration)
│   ├── Robots (REST API or local IP)
│   ├── Smart locks, doors, security
│   ├── Lighting, HVAC, environment sensors
│   ├── NAS / local storage
│   ├── Smart plugs / power monitoring
│   └── Any device with a network connection
│
├── Home Assistant (same always-on machine)
│   ├── 3,000+ integrations — one API for all devices
│   ├── Local automations (run even if VPS is down)
│   ├── Exposes everything via REST API + WebSocket
│   ├── Frigate NVR for AI camera processing
│   └── Runs in Docker or Home Assistant OS
│
├── OpenClaw Node (same always-on machine)
│   ├── Connected to VPS gateway via Tailscale
│   ├── Executes agent commands locally (low latency)
│   ├── Calls Home Assistant API on behalf of agents
│   └── Works even if VPS is temporarily unreachable
│
└── Tailscale (same always-on machine)
    ├── Subnet router: advertises entire LAN to Tailscale
    ├── VPS can reach every device on LAN directly
    ├── No Tailscale app needed on individual devices
    └── One command: tailscale up --advertise-routes=192.168.x.0/24
```

---

## How Agents Interact With Physical Spaces

```
Legacy / any agent (VPS)
        ↓
calls Home Assistant REST API or WebSocket
        ↓
Home Assistant controls the device
        ↓
Physical action happens in the space

— OR for direct commands —

Legacy / Grid agent (VPS)
        ↓
exec host=office-node (OpenClaw node)
        ↓
command runs locally on office machine
        ↓
physical action (cast, file op, local process)
```

**Intelligence stays centralized on VPS. Execution is distributed to physical nodes.**

---

## Home Assistant — The 3,000 Integration Layer

Home Assistant is the device abstraction layer. Instead of building custom integrations for every device, HA handles the device layer and Legacy talks to one API.

### What Legacy Gets Through HA

| Category | Integrations / Capabilities |
|----------|----------------------------|
| **Cameras** | Live feed, motion detection, recordings, Frigate AI |
| **Presence detection** | Who's in the building, which room, arrival/departure |
| **Environment** | Temperature, humidity, CO2, air quality, noise levels |
| **Security** | Smart locks, alarm system, motion sensors, door/window sensors |
| **Displays & TVs** | Power on/off, input switching, volume, brightness |
| **Lighting** | Scene control, color, brightness, schedules |
| **3D Printer** | OctoPrint — start/stop/monitor prints remotely |
| **Robots** | Any robot with REST API or local network control |
| **Power** | Smart plugs — know what's on, control it, monitor usage |
| **Voice/Audio** | Speak through smart speakers in any room |
| **HVAC** | Thermostat control, schedules, energy monitoring |
| **Projectors** | WiFi projectors — power, input, content |
| **NAS** | Network storage monitoring and control |
| **Network** | Router stats, device presence, bandwidth |

### Key Integrations for One Mind OS

- **Frigate** — local AI NVR for cameras (see facial recognition below)
- **LLM Vision** — connects cameras to AI vision models (Claude/GPT-4o/Gemini) via HA
- **OctoPrint** — 3D printer control
- **WLED** — addressable LED control (studio/office ambiance)
- **Unifi** — network presence detection (who's on WiFi)
- **MQTT** — universal IoT protocol (robots, sensors, custom devices)
- **Webhook** — Legacy triggers HA automations via HTTP

---

## Facial Recognition — Zeus Detection

### Stack
```
IP Camera (Ring, Reolink, or any RTSP camera)
        ↓
Frigate NVR (runs on same always-on machine)
Local AI model — detects faces, objects, motion in real-time
No cloud, no subscription, entirely on-device
        ↓
Frigate fires event: "Zeus detected at front door"
        ↓
Home Assistant receives Frigate event
        ↓
HA automation triggers → calls Legacy via webhook
        ↓
Legacy responds with action:
  - "Welcome back Zeus" cast to nearest TV/display
  - Unlock door automatically
  - Set lighting to Zeus's preferred scene
  - Start morning brief on displays
  - Notify iPhone if unexpected person detected
```

### What Frigate Can Detect
- Faces (with face recognition models)
- People, vehicles, animals, packages
- Motion zones (only alert for certain areas)
- Loitering (person present for >X seconds)
- Multiple cameras simultaneously
- Stores clips and snapshots locally

### Hardware Requirements for Frigate — NO GPU NEEDED
Frigate runs entirely on **CPU** — no GPU required. It's slower without acceleration but fully functional for 1-3 cameras.

| Hardware | Performance | Cost | Notes |
|----------|-------------|------|-------|
| CPU only | ~5-10 FPS detection | Free | Works fine for 1-3 cameras, start here |
| **Google Coral USB TPU** | 100+ FPS, near zero CPU load | ~$60 | ✅ Recommended upgrade — USB stick, plug and go |
| Intel iGPU (OpenVINO) | Good acceleration | Free | Built into Intel CPUs, no extra hardware |
| NVIDIA GPU | Overkill | $200+ | Not needed for home/office use |

**Plan:** Start with CPU only. Add Coral USB TPU ($60) when you want to scale to more cameras or need faster detection.

---

## Tailscale Subnet Router — Setup

### What It Does
Advertises a local network (`192.168.x.0/24`) to the Tailscale network. After setup, VPS can reach every device on that LAN directly — no Tailscale app needed on individual devices.

### Setup Steps (when ready)
```bash
# On the always-on office machine:
tailscale up --advertise-routes=192.168.1.0/24

# In Tailscale admin console (one-time approval):
# https://login.tailscale.com/admin/machines
# Click the machine → "Edit route settings" → approve the subnet

# On VPS (enable route acceptance):
tailscale up --accept-routes

# Test from VPS:
ping 192.168.1.172   # Should reach OMOS TV directly
```

### After Setup
```bash
# From VPS, scan entire office network:
nmap -sn 192.168.1.0/24           # Find all live devices
nmap -p 8008,8009 192.168.1.0/24  # Find all Cast devices
nmap -p 80,443,8080 192.168.1.0/24 # Find all web interfaces
```

---

## Subnet Router vs Exit Node — Key Distinction

| Feature | Subnet Router | Exit Node |
|---------|--------------|-----------|
| Purpose | Legacy reaches INTO your LAN | Your internet traffic routes THROUGH a machine |
| Direction | Outside → LAN devices | Your device → Internet |
| Use case for One Mind OS | ✅ Yes — agent controls devices | ❌ Not needed for device control |
| Hides your IP? | No | Yes (like a VPN) |
| LAN devices need Tailscale? | No | No |
| Can run both on same machine? | Yes | Yes |

**One Mind OS needs subnet router, not exit node.**

---

## LLM Vision — AI Scene Understanding for Cameras

**What it is:** A Home Assistant integration (`llmvision.org`) that sends camera frames to multimodal AI models for full scene understanding — not just "person detected" but actual intelligent analysis.

**Without LLM Vision:**
```
Camera → Frigate → "person detected" (basic object detection only)
```

**With LLM Vision:**
```
Camera → Frigate event → LLM Vision → Claude/GPT-4o analyzes frame
→ "Zeus is at the front door wearing a grey hoodie, carrying a bag"
→ HA automation → Legacy gets full context and responds intelligently
```

### How It Fits in One Mind OS
```
Camera (RTSP stream)
        ↓
Frigate NVR — motion/object detection trigger
        ↓
LLM Vision (HA integration)
        ↓ sends frame to Claude / GPT-4o / Gemini (your choice)
"Zeus detected at front door" / "Unknown person" / "Package delivered"
        ↓
Home Assistant automation fires
        ↓
Legacy webhook → intelligent response
  - Welcome Zeus → cast greeting to TV
  - Unknown person → alert Zeus's iPhone
  - Package → notify Zeus
  - Security event → Guardian agent handles it
```

### Supported AI Providers (via LLM Vision)
- Claude (Anthropic) — already used by Legacy ✅
- GPT-4o (OpenAI)
- Gemini (Google)
- Any OpenAI-compatible API endpoint
- Local models (Ollama) for zero-cost analysis

### Cost
- Integration: free (install via HACS)
- Per query: ~$0.001/image with Claude — negligible for event-triggered analysis
- With local Ollama model: completely free, runs on-device

### Install
```
HACS → Integrations → search "LLM Vision" → install
Configuration → add provider (Claude recommended)
Set up automations: Frigate event → LLM Vision analyze → action
```

---

## Camera Compatibility — Yi Cameras (YYS.2016 / YHS.2016)

### The Problem
Yi cameras from 2016 use a **proprietary RTMP stream** — not standard RTSP. Frigate requires RTSP. Yi locked this down to push their cloud service.

### Options

**Option 1 — Yi-Hack custom firmware (free)**
- Open source firmware that unlocks RTSP on many Yi models
- Actively maintained: `github.com/roleoroleo/yi-hack`
- Check compatibility list for YYS.2016 / YHS.2016 specifically
- If supported: flashes via SD card, enables RTSP at `rtsp://camera_ip/ch0_0.h264`
- Risk: voids warranty (already old cameras — acceptable)

**Option 2 — Check for hidden RTSP**
- Some Yi models have RTSP already at: `rtsp://camera_ip/ch0_0.h264`
- Try connecting with VLC before flashing anything

**Option 3 — Replace them (recommended long-term)**
- Yi 2016 cameras are nearly 10 years old
- Limited resolution, no modern features
- Reolink 5MP cameras: $35 each, native RTSP, rock solid with Frigate
- Dahua / Hikvision: best quality, recommended by Frigate docs

### Recommended Camera Stack for One Mind OS
| Camera | Price | Why |
|--------|-------|-----|
| **Reolink RLC-510A** | ~$45 | RTSP native, 5MP, PoE, AI detection built-in |
| **Dahua IPC-T549M** | ~$80 | Best quality, rock solid streams, large sensor |
| **Amcrest IP5M-T1179** | ~$60 | Rebranded Dahua, good value |

**Frigate recommends:** Dahua > Hikvision > Amcrest > Reolink (5MP max for Reolink)
**WiFi cameras NOT recommended** — streams are unreliable for Frigate. Use PoE where possible.

---

## Multi-Location Scale

```
Home (192.168.1.0/24)
├── Always-on machine: Mac Mini / NUC
├── Tailscale subnet router ✅
├── OpenClaw node ✅
├── Home Assistant ✅
└── Legacy sees: TV, cameras, locks, lights, everything

Office (192.168.2.0/24)
├── Always-on machine: Mac Mini / NUC
├── Tailscale subnet router ✅
├── OpenClaw node ✅
├── Home Assistant ✅
└── Legacy sees: 3D printer, robots, displays, projector, everything

Studio (192.168.3.0/24)
├── Always-on machine: Raspberry Pi 5
├── Tailscale subnet router ✅
├── OpenClaw node ✅
├── Home Assistant ✅
└── Legacy sees: cameras, lighting, audio, recording gear

          All connect back to:
┌─────────────────────────────┐
│   VPS onemindos-cloud-1     │
│   OpenClaw Gateway          │
│   All 9 agents              │
│   Single brain for all      │
│   physical spaces           │
└─────────────────────────────┘
```

---

## Recommended Hardware Per Location

| Use Case | Hardware | Cost | Notes |
|----------|----------|------|-------|
| **Primary node (office/home)** | Mac Mini M2 | ~$600 | Silent, powerful, runs HA + OpenClaw + Frigate |
| **Secondary/budget node** | Intel NUC | ~$300 | Good for HA + OpenClaw |
| **Minimal node** | Raspberry Pi 5 | ~$80 | Studio/simple spaces |
| **AI camera acceleration** | Google Coral USB TPU | ~$60 | Plugs into node machine |
| **Cameras** | Reolink (RTSP) | ~$30-60 | Local stream, no cloud needed |
| **Smart hub** | Zigbee/Z-Wave USB stick | ~$30 | Connects 100s of cheap sensors |

**Best setup:** Mac Mini M2 (silent, always-on, handles everything) + Coral TPU for Frigate.

---

## Agent Routing for Physical Spaces

| Agent | Physical Role |
|-------|--------------|
| **Legacy ⚡** | Orchestrates across all locations, receives presence events |
| **Grid ⚙️** | Executes on local nodes, monitors hardware health |
| **Guardian 🛡️** | Handles security cameras, alerts, access control |
| **Empire 💰** | Renders content → casts to displays/TVs |
| **Haven 🏡** | Home automation, family presence, comfort settings |

---

## Home Assistant API — How Legacy Calls It

### REST API (simple actions)
```bash
# Turn on office lights
curl -X POST http://192.168.1.HA_IP:8123/api/services/light/turn_on \
  -H "Authorization: Bearer HA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "light.office_main"}'

# Get all entity states
curl http://192.168.1.HA_IP:8123/api/states \
  -H "Authorization: Bearer HA_TOKEN"

# Trigger an automation
curl -X POST http://192.168.1.HA_IP:8123/api/services/automation/trigger \
  -H "Authorization: Bearer HA_TOKEN" \
  -d '{"entity_id": "automation.welcome_zeus"}'
```

### WebSocket (real-time events)
- Subscribe to any event (motion detected, door opened, Zeus arrived)
- HA pushes events to Legacy in real-time
- Legacy reacts immediately without polling

---

## Implementation Roadmap

### Phase 1 — Current (done)
- [x] OMOS TV on Tailscale ✅
- [x] Cast infrastructure (pychromecast + receiver) ✅
- [x] Mac as temporary cast bridge ✅

### Phase 2 — Home Setup
- [ ] Always-on machine at home (Mac Mini or NUC)
- [ ] Tailscale subnet router configured
- [ ] OpenClaw node connected
- [ ] Home Assistant installed + basic devices connected
- [ ] Cast moves from Mac-bridge to direct VPS → TV

### Phase 3 — Cameras + Presence
- [ ] IP cameras installed (Reolink recommended)
- [ ] Frigate NVR configured
- [ ] Google Coral TPU installed
- [ ] Zeus facial recognition model trained
- [ ] Legacy receives presence events

### Phase 4 — Office
- [ ] Office machine (Mac Mini M2)
- [ ] Full HA setup: 3D printer, displays, projector
- [ ] Grid agent assigned to office node
- [ ] Robots connected via MQTT or REST

### Phase 5 — Full Multi-Location
- [ ] All locations on subnet routers
- [ ] Cross-location automations (Zeus leaves home → office prepares)
- [ ] Unified presence across all spaces
- [ ] Legacy has real-time awareness of all physical spaces

---

## Key Principle

**Every physical space Zeus controls becomes an extension of One Mind OS.**
One machine per location. One Home Assistant per location. One OpenClaw node per location.
Legacy orchestrates all of it from a single VPS gateway.
Sovereign. Local-first. Infinitely scalable.

---

*Last updated: 2026-04-18 (v2 — added Frigate GPU clarification, LLM Vision, Yi camera compatibility, camera recommendations) | Legacy ⚡ | SO 02 Infrastructure*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
