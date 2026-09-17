# Homestead Edge Architecture — 6-Building Intelligence Grid

> **Status**: Architecture Design (Pre-Implementation)
> **Owner**: LE Agent (Legacy Evolution — Homestead Steward)
> **Created**: 2026-03-03
> **Domain**: LE (50-74) + HP (25-49 environment) + GE (75-99 yield/ROI)

---

## Overview

Six buildings on Zeus's homestead, each wired as an intelligent edge node connected to one central OpenClaw Gateway. Three technology layers per building — all using native OpenClaw primitives (zero custom daemons).

| # | Building | Primary Function | Special Integration |
|---|----------|-----------------|---------------------|
| 1 | **Main House** | Living quarters, central hub | HA integration, full sensor suite |
| 2 | **Office** | Command center, development | Power monitoring, network core |
| 3 | **Dojo** | Training, meditation, wellness | Sound/motion, session tracking |
| 4 | **Greenhouse** | Automated growing (FarmBot) | Soil/humidity sensors, FarmBot API |
| 5 | **Guest House** | Guest accommodation | Smart lock, occupancy detection |
| 6 | **Barn** | Livestock, workshop, storage | Climate control, livestock monitoring |

---

## Master Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ZEUS'S HOMESTEAD — 6 BUILDINGS                        │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │              CENTRAL GATEWAY (Main House or VPS)                 │    │
│  │                                                                 │    │
│  │   OpenClaw Gateway :18789                                       │    │
│  │   ├── legacy  (orchestrator — "give me all 6 statuses")         │    │
│  │   ├── hp      (tracks environment health → PI Score)            │    │
│  │   ├── le      (homestead steward — owns the buildings)          │    │
│  │   └── ge      (ROI on greenhouse yield, barn operations)        │    │
│  │                                                                 │    │
│  │   NATS JetStream :4222  (WireClaw sensors + event archive)      │    │
│  │   OmniDB :5432          (time-series + embeddings)              │    │
│  │   Home Assistant :8123  (Zigbee/Z-Wave/WiFi device control)     │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│            ┌────────┬────────┼────────┬─────────┬──────────┐           │
│            ▼        ▼        ▼        ▼         ▼          ▼           │
│         ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐        │
│         │ MAIN │ │OFFICE│ │ DOJO │ │GREEN │ │GUEST │ │ BARN │        │
│         │HOUSE │ │      │ │      │ │HOUSE │ │HOUSE │ │      │        │
│         └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘        │
│            │        │        │        │         │        │             │
│  Layer 1:  Pi/Mini  Pi/Mini  Pi/Mini  Pi/Mini   Pi/Mini  Pi/Mini      │
│  (Node)    headless headless headless headless  headless headless     │
│            node     node     node     node      node     node        │
│                                                                       │
│  Layer 2:  WireClaw WireClaw WireClaw WireClaw  WireClaw WireClaw    │
│  (Sensors) motion   motion   motion   temp      motion   motion      │
│            temp     temp     humidity humidity  temp     temp         │
│            door     power    sound    soil      door     livestock    │
│                                                                       │
│  Layer 3:  HeyClaw  HeyClaw  HeyClaw  ───       HeyClaw  HeyClaw     │
│  (Voice)   kitchen  desk     mat area           bedroom  workbench   │
│            bedroom                                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Three Layers Per Building

### Layer 1 — Headless OpenClaw Node (Brain)

Each building runs a **Raspberry Pi 4/5 or mini PC** as a native OpenClaw headless node. No custom code — this is a built-in OpenClaw feature.

**What each node exposes natively:**

| Command | Purpose |
|---------|---------|
| `system.run` | Execute any command (check HA, read sensors, run scripts) |
| `system.notify` | Push notifications to that building's display |
| `system.which` | Verify what tools are installed on the node |
| `system.execApprovals` | Manage per-node command allowlists |

**Setup per building:**

```bash
# On each Pi (e.g., the Barn):
export OPENCLAW_GATEWAY_TOKEN="<your-token>"
openclaw node install \
  --host gateway.tailnet \
  --port 18789 \
  --display-name "Barn"

openclaw node restart
```

**Approve on gateway:**

```bash
openclaw devices list
openclaw devices approve <requestId>
```

**Example — run a command on the Barn node:**

```json
{
  "action": "run",
  "node": "barn",
  "command": ["sensors", "--json"],
  "commandTimeoutMs": 5000
}
```

**Connectivity:** Tailscale mesh VPN connects all buildings. Each Pi is on the tailnet. Gateway runs at `gateway.tailnet:18789`. Node connects via WebSocket:

```
Pi (barn) ── Tailscale ──→ Gateway :18789 (WebSocket)
```

If running the gateway on the local network, SSH tunnels work too:

```bash
# Terminal A: forward local 18790 -> gateway
ssh -N -L 18790:127.0.0.1:18789 user@gateway-host

# Terminal B: connect node through tunnel
OPENCLAW_GATEWAY_TOKEN="<token>" openclaw node run --host 127.0.0.1 --port 18790 --display-name "Barn"
```

---

### Layer 2 — WireClaw ESP32 Sensors (Nerves)

$5 each, web-flashed at wireclaw.io (no toolchain needed). Each ESP32 publishes sensor data to NATS subjects at **30ms latency**.

**NATS Subject Convention:**

```
wireclaw.{building}.{room}.{sensor_type}
```

**Per-Building Sensor Map:**

| Building | Sensors | NATS Subjects |
|----------|---------|---------------|
| Main House | motion, temp, door (kitchen, bedroom, living) | `wireclaw.main-house.kitchen.motion`, `.temp`, `.door` |
| Office | motion, temp, power, network | `wireclaw.office.desk.motion`, `.temp`, `.power` |
| Dojo | motion, temp, humidity, sound | `wireclaw.dojo.main.motion`, `.temp`, `.sound` |
| Greenhouse | temp, humidity, soil x6 zones | `wireclaw.greenhouse.zone-a.soil`, `.temp`, `.humidity` |
| Guest House | motion, temp, door | `wireclaw.guest-house.bedroom.motion`, `.door` |
| Barn | temp x3 stalls, motion, air quality | `wireclaw.barn.stall-1.temp`, `.main.motion` |

**Local Rule Engine (runs ON the ESP32 — works offline):**

```
IF motion AND time > 10pm        → lights on AND alert Legacy
IF soil_moisture < 30%            → trigger FarmBot irrigation zone
IF barn_temp > 95°F               → alert + turn on fans
IF guest-house.door == "open"     → check occupancy → notify if empty
AND time > 11pm
IF dojo.sound > 70dB              → log training session start
AND time BETWEEN 5am-8am
```

Legacy discovers all WireClaw devices automatically via NATS — the nexus-bridge skill subscribes to `wireclaw.>` subjects.

---

### Layer 3 — HeyClaw Voice Pucks (Ears)

$10 each (ESP32 + speaker + mic). Wake word: **"Hey Legacy"**. Each connects via WiFi to the OpenClaw Gateway WebSocket.

**Per-Building Voice Placement:**

| Building | Location | Example Commands |
|----------|----------|-----------------|
| Main House | Kitchen, Bedroom | "Hey Legacy, what's the greenhouse soil moisture?" |
| Office | Desk | "Hey Legacy, is anyone in the guest house?" |
| Dojo | Mat area | "Hey Legacy, set a 30-minute timer for training" |
| Greenhouse | — (no voice) | — (noisy fans, use phone/watch instead) |
| Guest House | Bedroom | "Hey Legacy, what's the WiFi password?" |
| Barn | Workbench | "Hey Legacy, what's the temperature in stall 3?" |

**Flow:** Voice → ESP32 mic → OpenClaw Gateway → Legacy processes → TTS response → plays on puck speaker.

---

### Layer 4 (Greenhouse Special) — FarmBot Integration

The greenhouse has unique plumbing beyond the standard 3 layers:

```
FarmBot API ←→ Home Assistant ←→ OpenClaw Gateway
                    ↑
              MQTT :1883
                    ↑
         WireClaw soil/humidity/temp sensors (6 zones)
         MimiClaw smart sensors (pre-analyze soil trends)
```

**FarmBot provides:** Automated planting sequences, watering schedules, weeding, soil height mapping.

**HA ties it together:** Environmental sensors trigger FarmBot actions. OpenClaw (via the `le` agent) can:
- "What's growing right now?"
- "Irrigation status for zone B?"
- "Plant the spring lettuce sequence"
- "Greenhouse yield this month?"

---

## MimiClaw — Edge Intelligence Upgrade

MimiClaw is an optional **brain upgrade** for select WireClaw sensors. ESP32-S3, pure C, 512KB RAM, $5. Adds local reasoning via a ReAct agent loop running on bare metal firmware.

**Standard WireClaw vs MimiClaw:**

```
Standard WireClaw:   temp=95°F → publish raw to NATS → Legacy decides

MimiClaw upgrade:    temp=95°F → local reasoning: "unusual for March,
                     rising trend, was 85°F an hour ago, livestock
                     comfort threshold breached" → publish ENRICHED
                     event to NATS → Legacy gets pre-analyzed data
```

**Best use cases:**
- **Greenhouse**: Reason about soil + temp + humidity together locally — "Zone A trending dry, recommend irrigation in ~2 hours"
- **Barn**: Livestock comfort thresholds — "Stall 2 heat stress likely by 3 PM based on trend"

---

## Multi-Agent Routing — Who Handles What

The 4-agent hive each has a role in homestead operations:

```
  "Hey Legacy, give me a homestead status report"
                    │
                    ▼
              ┌─────────┐
              │ legacy   │ ← Orchestrator, receives request
              └────┬─────┘
                   │ sessions_send or direct tool use
        ┌──────────┼──────────────────┐
        ▼          ▼                  ▼
  ┌──────────┐ ┌──────────┐    ┌──────────┐
  │    le    │ │    hp    │    │    ge    │
  │ Steward  │ │  Coach   │    │ Advisor  │
  └──────────┘ └──────────┘    └──────────┘
  "6 buildings   "Barn temp     "Greenhouse
   all online,    affecting      yield: 40lbs
   1 advisory"   sleep quality"  this month,
                                 $2.10/lb ROI"
```

| Agent | Homestead Role | Example Questions |
|-------|---------------|-------------------|
| **legacy** | Orchestrator — routes, compiles, delivers | "Status report", "Turn on barn fans" |
| **le** | Homestead steward — owns buildings, maintenance, 50-year view | "Guest house ready for visitors?", "Maintenance schedule" |
| **hp** | Environment + wellness — tracks how buildings affect Zeus | "Is my sleep environment optimal?", "Air quality trend?" |
| **ge** | ROI + operations — greenhouse economics, barn productivity | "Greenhouse cost vs yield?", "Energy spend per building?" |

---

## Status Report Flow — Step by Step

When Zeus says **"Hey Legacy, give me a homestead status report"**:

### Step 1: Legacy Receives (any channel)

Legacy recognizes this as LE domain. Two routing options:

**Option A — Legacy handles directly (simple queries):**

```
nodes → status                         (all 6 online/offline)
nodes → run --node each -- "sensors + HA"  (per-building data)
nexus_query → wireclaw.*.*.temp        (all temp readings)
```

**Option B — Legacy delegates to LE (complex analysis):**

```
sessions_send({
  sessionKey: "le",
  message: "Run a full homestead status report — all 6 buildings"
})
```

### Step 2: Data Collection (parallel via node.invoke)

```bash
# Each runs on respective building's Pi via system.run
nodes → run --node main-house   -- "sensors --json && ha-cli state summary"
nodes → run --node office       -- "uptime && sensors --json && df -h /"
nodes → run --node dojo         -- "uptime && sensors --json"
nodes → run --node greenhouse   -- "ha-cli state farmbot && sensors --json"
nodes → run --node guest-house  -- "uptime && ha-cli state lock.front_door"
nodes → run --node barn         -- "uptime && sensors --json && ha-cli state climate.barn"
```

Plus WireClaw sensor data via NATS (nexus-bridge skill):

```
nexus_query({ subject: "wireclaw.*.*.temp", last: 1 })
nexus_query({ subject: "wireclaw.greenhouse.*.soil", last: 1 })
```

### Step 3: Compile and Deliver

Legacy synthesizes into a formatted report, delivered on whatever channel was used:

```
🏠 HOMESTEAD STATUS — March 3, 2026 @ 2:15 PM

1. MAIN HOUSE          ✅ Online (uptime 47d)
   Temp: 72°F | Humidity: 45% | Motion: kitchen (2m ago)
   All doors: locked | Energy: 4.2 kWh today

2. OFFICE               ✅ Online (uptime 47d)
   Temp: 70°F | Disk: 62% used | No motion (empty)

3. DOJO                  ✅ Online (uptime 12d)
   Temp: 68°F | Humidity: 52% | Sound: 32dB (quiet)
   Last session: yesterday 6:30 AM

4. GREENHOUSE            ✅ Online (uptime 47d)
   Zone A: Soil 38% 🟡 (irrigating) | Temp: 82°F | Humidity: 76%
   Zone B: Soil 55% ✅ | Temp: 80°F
   FarmBot: Running spring lettuce sequence (Step 4/7)

5. GUEST HOUSE           ✅ Online (uptime 47d)
   Temp: 71°F | Unoccupied | Front door: locked
   Last guest checkout: Feb 28

6. BARN                  ⚠️ Online (uptime 3d — rebooted)
   Temp: 74°F (⬆ warming) | Stall 1: 72°F | Stall 2: 73°F
   Fans: ON | Feed system: nominal
   NOTE: Node rebooted 3 days ago — check power stability

📊 Summary: 6/6 online | 1 advisory (barn power) | 1 active irrigation
```

---

## Complete Wiring Diagram

```
                        ┌─────────────────────────┐
     Zeus's phone ───── │   OpenClaw Gateway       │ ──── Zeus's Mac (node)
     (WhatsApp/         │   :18789 (Tailscale)     │      (Canvas, Voice Wake)
      Discord/          │                          │
      iMessage)         │  legacy ← orchestrator   │
                        │  hp     ← environment    │
                        │  le     ← steward        │
                        │  ge     ← economics      │
                        └──────────┬──────────────┘
                                   │ WebSocket + NATS
                    ┌──────────────┼──────────────────┐
                    │              │                   │
           ┌────────▼──┐  ┌───────▼────┐  ┌──────────▼───┐
           │ 6 Headless │  │ 20+ Wire  │  │ 5 HeyClaw    │
           │ Pi Nodes   │  │ Claw ESP32│  │ Voice Pucks  │
           │ (system.   │  │ (NATS at  │  │ ("Hey Legacy │
           │  run)      │  │  30ms)    │  │   ...")      │
           └────────────┘  └───────────┘  └──────────────┘
                    │              │                   │
           exec on each    sensor data          voice commands
           building's Pi   per room/zone        per building
```

---

## Hardware Bill of Materials

| Item | Qty | Unit Cost | Total | Purpose |
|------|-----|-----------|-------|---------|
| Raspberry Pi 5 (4GB) | 6 | ~$60 | ~$360 | Headless node per building |
| WireClaw ESP32 | 20 | ~$5 | ~$100 | 3-4 sensors per building |
| HeyClaw ESP32 + speaker | 5 | ~$10 | ~$50 | Voice in 5 buildings |
| MimiClaw ESP32-S3 | 2 | ~$5 | ~$10 | Smart sensors (greenhouse + barn) |
| PoE switch (5-port) | 6 | ~$30 | ~$180 | Power + network per building |
| microSD cards (32GB) | 6 | ~$8 | ~$48 | Pi storage |
| Ethernet cable (per building) | 6 | ~$15 | ~$90 | Pi → PoE switch |
| USB-C power supply (Pi) | 6 | ~$12 | ~$72 | Backup if no PoE hat |
| **TOTAL** | | | **~$910** | **Full 6-building intelligence grid** |

---

## Installation Sequence Per Building

```bash
# ── Step 1: Flash Pi with Raspberry Pi OS Lite ──
# Use Raspberry Pi Imager → Raspberry Pi OS Lite (64-bit)
# Set hostname: barn.local (or main-house.local, etc.)
# Enable SSH, configure WiFi as backup

# ── Step 2: Install Tailscale on Pi ──
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up --authkey=<your-tailscale-key>

# ── Step 3: Install OpenClaw ──
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
sudo apt install -y nodejs
npm install -g openclaw@latest

# ── Step 4: Register as headless node ──
export OPENCLAW_GATEWAY_TOKEN="<your-gateway-token>"
openclaw node install \
  --host gateway.tailnet \
  --port 18789 \
  --display-name "Barn"

# ── Step 5: Start the node service ──
openclaw node restart

# ── Step 6: Approve on gateway ──
# (run on gateway machine)
openclaw devices list
openclaw devices approve <requestId>
openclaw nodes status

# ── Step 7: Configure exec approvals ──
openclaw approvals allowlist add --node barn "/usr/bin/sensors"
openclaw approvals allowlist add --node barn "/usr/local/bin/ha-cli"
openclaw approvals allowlist add --node barn "/usr/bin/uptime"

# ── Step 8: Flash WireClaw sensors ──
# → Go to wireclaw.io, select ESP32, web flash
# → Configure WiFi SSID + password
# → Set NATS URL: nats://gateway.tailnet:4222
# → Set device ID: barn-stall-1, barn-main, etc.

# ── Step 9: Flash HeyClaw voice puck ──
# → Same ESP32 hardware, HeyClaw firmware
# → Configure WiFi + Gateway URL: ws://gateway.tailnet:18789
# → Set wake word: "Hey Legacy"

# ── Step 10: Verify ──
openclaw nodes status          # Should show building online
openclaw nodes describe --node barn
openclaw nodes run --node barn -- echo "Hello from the barn"
```

---

## NATS Subject Map (All Buildings)

```
# ── Main House ──
wireclaw.main-house.kitchen.motion
wireclaw.main-house.kitchen.temp
wireclaw.main-house.bedroom.temp
wireclaw.main-house.bedroom.humidity
wireclaw.main-house.living.motion
wireclaw.main-house.living.door
wireclaw.main-house.front.door

# ── Office ──
wireclaw.office.desk.motion
wireclaw.office.desk.temp
wireclaw.office.server.temp
wireclaw.office.main.power

# ── Dojo ──
wireclaw.dojo.main.motion
wireclaw.dojo.main.temp
wireclaw.dojo.main.humidity
wireclaw.dojo.main.sound

# ── Greenhouse ──
wireclaw.greenhouse.zone-a.soil
wireclaw.greenhouse.zone-a.temp
wireclaw.greenhouse.zone-a.humidity
wireclaw.greenhouse.zone-b.soil
wireclaw.greenhouse.zone-b.temp
wireclaw.greenhouse.zone-b.humidity
wireclaw.greenhouse.zone-c.soil
wireclaw.greenhouse.zone-c.temp
wireclaw.greenhouse.ambient.temp
wireclaw.greenhouse.ambient.humidity

# ── Guest House ──
wireclaw.guest-house.bedroom.motion
wireclaw.guest-house.bedroom.temp
wireclaw.guest-house.main.door
wireclaw.guest-house.front.door

# ── Barn ──
wireclaw.barn.stall-1.temp
wireclaw.barn.stall-2.temp
wireclaw.barn.stall-3.temp
wireclaw.barn.main.motion
wireclaw.barn.main.temp
wireclaw.barn.main.air-quality
wireclaw.barn.workshop.motion
```

---

## Proactive Alerts (via WireClaw Rules + Legacy HEARTBEAT)

| Trigger | Action | Agent |
|---------|--------|-------|
| `barn.stall-*.temp > 90°F` | Alert Zeus + turn on fans via HA | le |
| `greenhouse.zone-*.soil < 25%` | Trigger FarmBot irrigation | le |
| `guest-house.door == open AND time > 11pm AND unoccupied` | Security alert | legacy |
| `dojo.sound > 70dB AND time 5am-8am` | Log training session | hp |
| `office.power > 500W sustained 2h` | Check for runaway process | legacy |
| `any-building node offline > 5min` | Infrastructure alert | legacy |
| `main-house.bedroom.temp > 76°F at night` | "Lower AC, affects sleep quality" | hp |

---

## Related Docs

- [MULTI-AGENT.md](MULTI-AGENT.md) — 4-agent hive configuration
- [HYBRID-MIND-ARCHITECTURE.md](HYBRID-MIND-ARCHITECTURE.md) — Central Cortex + edge design
- [NODES.md](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/openclaw/NODES.md) — OpenClaw node reference
- [ecosystem/ECOSYSTEM-42.md](ecosystem/ECOSYSTEM-42.md) — WireClaw (#24), HeyClaw (#25), MimiClaw (#26) details
- [ecosystem/DEPLOYMENT-SEQUENCE.md](ecosystem/DEPLOYMENT-SEQUENCE.md) — Tier 3C IoT deployment plan
