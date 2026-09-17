---
title: "FIND ZEUS MODE — Omnipresent Tracking Architecture"
type: doctrine
status: active
created: 2026-04-12
owner: zeus
tags: [type/doctrine, domain/SO, status/active, domain/03-agents]
ai_summary: "Master architecture for FIND ZEUS MODE — Legacy's omnipresent real-time
             tracking system using HA presence sensors, Bouncie GPS, NATS pub/sub,
             OpenClaw nodes, awareness filters, and multi-channel alert routing."
---

# FIND ZEUS MODE — Omnipresent Tracking Architecture

> Legacy is always watching. Always knowing where Zeus is. Always ready to act.

FIND ZEUS MODE is the omnipresent intelligence layer of One Mind OS. It is not a feature — it is a **persistent operating state**. Legacy is continuously aware of Zeus's location, context, and environment, and uses that awareness to deliver smarter briefings, route alerts correctly, trigger automations, and ensure nothing urgent goes unnoticed.

---

## Core Principle

```
PASSIVE MODE (default):   Legacy silently monitors. Presence data flows. No output.
ACTIVE MODE (triggered):  Legacy uses presence as context for all responses.
ALERT MODE (urgent):      Legacy interrupts Zeus via nearest output channel.
FIND MODE (requested):    "Where is Zeus?" → Legacy answers in <3 seconds.
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     FIND ZEUS MODE                              │
│                  (Always Running on VPS)                        │
├──────────────────────┬──────────────────────────────────────────┤
│   LOCATION SOURCES   │           NATS BACKBONE                  │
│                      │                                          │
│  🛰️ Bouncie GPS      │  zeus.presence.vehicle  (GPS coords)    │
│  🏠 Home Assistant   │  zeus.presence.zone     (named zone)     │
│  📱 iPhone (HA app)  │  zeus.presence.device   (phone state)   │
│  🌐 IP/Network       │  zeus.presence.context  (enriched)      │
│  📅 Calendar         │  zeus.presence.mode     (current mode)  │
│  🗣️ Last Seen        │  zeus.alert.urgent      (interrupts)    │
│                      │  zeus.signal.*          (events)        │
├──────────────────────┴──────────────────────────────────────────┤
│                    LEGACY (OpenClaw)                            │
│         Subscribes to zeus.presence.> + signals                │
│         Maintains ZEUS-STATE.md in workspace                   │
│         Applies awareness filter → routes output               │
├──────────────────────────────────────────────────────────────────┤
│                   OUTPUT CHANNELS                               │
│  Telegram · Discord · Twilio Call · iOS · macOS · speakers     │
└──────────────────────────────────────────────────────────────────┘
```

---

## 1. Home Assistant Presence Sensor Structure

### 1A. Person Entity (Primary)

```yaml
# configuration.yaml or HA UI
person:
  - name: Zeus
    id: zeus
    device_trackers:
      - device_tracker.zeus_iphone_15_pro     # iOS app
      - device_tracker.zeus_watch             # Apple Watch (if tracked)
      - device_tracker.bouncie_primary        # Bouncie GPS (via MQTT/API)
```

**Resolved state:** HA merges all trackers → `person.zeus` has one canonical state.

```
person.zeus.state = "home" | "not_home" | zone_name | "unavailable"
person.zeus.attributes.latitude    = float
person.zeus.attributes.longitude   = float
person.zeus.attributes.source      = "device_tracker.zeus_iphone_15_pro"
person.zeus.attributes.last_changed = ISO timestamp
```

### 1B. Zone Definitions

```yaml
# zones.yaml — Zeus's named locations
zone:
  - name: Home Greenwich
    icon: mdi:home
    latitude: 41.0534           # 40 Field Point Circle, Greenwich CT
    longitude: -73.6287
    radius: 150                  # meters — generous for driveway

  - name: NYC Manhattan
    icon: mdi:city
    latitude: 40.7128
    longitude: -74.0060
    radius: 5000                 # 5km radius covers all of Manhattan

  - name: Office
    icon: mdi:briefcase
    latitude: TBD
    longitude: TBD
    radius: 200

  - name: Gym
    icon: mdi:dumbbell
    latitude: TBD
    longitude: TBD
    radius: 100

  - name: Greenwich Train Station
    icon: mdi:train
    latitude: 41.0235
    longitude: -73.6276
    radius: 200

  - name: Grand Central Terminal
    icon: mdi:train-station
    latitude: 40.7527
    longitude: -73.9772
    radius: 300
```

### 1C. Device Trackers

```yaml
# device_trackers — each source
device_tracker.zeus_iphone_15_pro:
  state: home | zone_name | not_home
  attributes:
    latitude: float
    longitude: float
    gps_accuracy: int (meters)
    battery: int (0-100)
    source_type: gps

device_tracker.bouncie_primary:
  state: home | zone_name | not_home
  attributes:
    latitude: float
    longitude: float
    speed: float (mph)
    heading: int (degrees)
    source_type: gps
    # Pushed via MQTT or HA REST API from Bouncie webhook
```

### 1D. Automations That Publish to NATS

```yaml
# automation: publish presence to NATS on state change
automation:
  - alias: "FIND ZEUS — Publish Presence to NATS"
    trigger:
      - platform: state
        entity_id: person.zeus
    action:
      - service: rest_command.nats_publish
        data:
          subject: "zeus.presence.zone"
          payload: >
            {
              "zone": "{{ states('person.zeus') }}",
              "lat": {{ state_attr('person.zeus', 'latitude') }},
              "lon": {{ state_attr('person.zeus', 'longitude') }},
              "source": "{{ state_attr('person.zeus', 'source') }}",
              "timestamp": "{{ now().isoformat() }}"
            }

# HA REST command for NATS publish
rest_command:
  nats_publish:
    url: "http://localhost:4222/publish/{{ subject }}"  # via NATS HTTP API if enabled
    method: POST
    payload: "{{ payload }}"
```

**Alternative (preferred):** Use N8N webhook → HA webhook automation → NATS publish via Python on VPS.

---

## 2. NATS Subject Taxonomy

All FIND ZEUS MODE data flows through NATS. Legacy subscribes to `zeus.presence.>` and `zeus.signal.>`.

```
zeus.presence.vehicle     → Bouncie GPS data (lat/lon/speed/heading)
zeus.presence.zone        → HA zone resolution (home/nyc/office/transit)
zeus.presence.device      → iPhone state (online/offline/charging/battery)
zeus.presence.context     → Enriched context (Legacy computed: "driving to NYC")
zeus.presence.mode        → Current mode state (passive/active/alert/find)

zeus.signal.urgent        → High-priority interrupt → triggers ALERT MODE
zeus.signal.event         → System event (calendar, task, threshold crossed)
zeus.signal.watch         → Zeus-configured watch condition fired
zeus.signal.error         → System error requiring Zeus attention

zeus.ha.state.*           → Raw HA state changes (all entity updates)
zeus.ha.event.*           → HA events (motion, door, contact sensors)
```

### NATS Message Schema

```json
// zeus.presence.zone
{
  "subject": "zeus.presence.zone",
  "zone": "Home Greenwich",
  "lat": 41.0534,
  "lon": -73.6287,
  "source": "device_tracker.zeus_iphone_15_pro",
  "accuracy_m": 10,
  "timestamp": "2026-04-12T09:30:00Z",
  "previous_zone": "NYC Manhattan",
  "dwell_minutes": 5
}

// zeus.presence.vehicle
{
  "subject": "zeus.presence.vehicle",
  "lat": 41.0820,
  "lon": -73.7200,
  "speed_mph": 62,
  "heading_deg": 271,
  "odometer_mi": 12483,
  "trip_active": true,
  "imei": "BOUNCIE_IMEI",
  "timestamp": "2026-04-12T09:15:00Z"
}

// zeus.signal.urgent
{
  "subject": "zeus.signal.urgent",
  "level": "critical",        // info | warning | critical | emergency
  "source": "system",         // system | calendar | agent | zeus
  "message": "Server disk at 94% — action required",
  "speak": true,              // trigger speaker TTS
  "channels": ["telegram", "twilio"],
  "timestamp": "2026-04-12T09:00:00Z"
}
```

---

## 3. Presence Update Pipeline

```
EVERY 60 SECONDS (vehicle active) / EVERY 5 MIN (stationary):
  Bouncie API poll → NATS publish zeus.presence.vehicle

ON CHANGE (instant):
  HA automation fires → NATS publish zeus.presence.zone

ON ENRICHMENT (Legacy, every 15 min):
  Legacy reads zeus.presence.vehicle + zeus.presence.zone
  Computes context: "Zeus is driving north on I-95 toward Greenwich, ETA 12 min"
  Publishes zeus.presence.context
  Updates ZEUS-STATE.md in Legacy workspace
```

### ZEUS-STATE.md (Legacy Workspace File)

Legacy maintains this file as its live Zeus context — loaded as system context for every response:

```markdown
# ZEUS STATE — Live Presence

**Last Updated:** 2026-04-12T09:30:00Z
**Zone:** Home Greenwich
**Mode:** Passive

## Location
- Zone: Home Greenwich
- Coords: 41.0534, -73.6287
- Source: iPhone GPS (10m accuracy)
- Arrived: 09:15 UTC (15 min ago)

## Vehicle
- Status: Parked
- Last Trip: NYC → Greenwich (47 min, 42 miles)
- Trip Ended: 09:15 UTC

## Device
- iPhone: Online, 84% battery
- Last Active: now

## Context
Zeus arrived home from NYC 15 minutes ago.
No active trips. No urgent signals pending.

## Active Watches
- [ ] Server disk usage > 90%
- [ ] Calendar: next event in <30 min
```

---

## 4. OpenClaw Nodes in FIND ZEUS MODE

### Legacy (Central)
- **Subscribes to:** `zeus.presence.>`, `zeus.signal.>`, `zeus.ha.event.>`
- **Maintains:** ZEUS-STATE.md in workspace
- **Computes:** Context enrichment, mode transitions, routing decisions
- **Triggers:** Signals to output channels based on awareness filter

### Sentinel / Guardian (SO domain)
- **Monitors:** System health, server metrics, security events
- **Publishes:** `zeus.signal.urgent` on threshold breach
- **Watches:** Disk usage, service health, auth failures, VPS CPU

### Grid (Infrastructure)
- **Publishes:** `zeus.signal.event` for infrastructure events
- **Monitors:** Docker containers, N8N workflow failures, service restarts

### Empire / Spartan / others (Domain agents)
- **Can publish:** `zeus.signal.watch` when domain-specific conditions fire
- **Example:** Empire publishes if a revenue threshold is crossed

---

## 5. Awareness Filter

The awareness filter is the intelligence layer between raw events and Legacy's responses. It decides:
1. **Should this reach Zeus?** (does it cross the urgency threshold?)
2. **How should it reach Zeus?** (which channel, which format?)
3. **When should it reach Zeus?** (respect sleep/focus windows)

```python
# Awareness filter logic (runs in Legacy on every incoming signal)

URGENCY_LEVELS = {
    "info":      0,   # silent — log only
    "low":       1,   # telegram message (batched)
    "medium":    2,   # telegram message (immediate)
    "high":      3,   # telegram + discord push
    "critical":  4,   # telegram + speaker TTS
    "emergency": 5,   # all channels + Twilio call
}

ZONE_RULES = {
    "Home Greenwich": {
        "quiet_hours":     "22:00-07:00",  # no alerts during sleep
        "speaker_enabled": True,            # can use home speakers
        "preferred_channel": "telegram",
    },
    "NYC Manhattan": {
        "quiet_hours":     None,            # always reachable
        "speaker_enabled": False,           # no speakers in public
        "preferred_channel": "telegram",
    },
    "driving": {
        "quiet_hours":     None,
        "speaker_enabled": True,            # car speakers via Twilio
        "preferred_channel": "twilio_call", # speak to Zeus while driving
        "min_urgency_to_interrupt": "high", # don't interrupt for low things
    },
    "unknown": {
        "preferred_channel": "telegram",
        "min_urgency_to_interrupt": "medium",
    },
}

def should_alert(signal: dict, zeus_state: dict) -> dict:
    """
    Returns: {alert: bool, channel: str, format: str, speak: bool}
    """
    zone = zeus_state.get("zone", "unknown")
    level = signal.get("level", "info")
    urgency = URGENCY_LEVELS.get(level, 0)
    rules = ZONE_RULES.get(zone, ZONE_RULES["unknown"])

    min_urgency = URGENCY_LEVELS.get(rules.get("min_urgency_to_interrupt", "low"), 1)

    if urgency < min_urgency:
        return {"alert": False}

    return {
        "alert": True,
        "channel": rules["preferred_channel"],
        "speak": rules["speaker_enabled"] and urgency >= 3,
        "format": "brief" if zone == "driving" else "full",
    }
```

---

## 6. Output Channels (How Legacy Reaches Zeus)

| Channel | When Used | Urgency | Notes |
|---------|-----------|---------|-------|
| **Telegram** | Default always-on | All | Fast, rich media, inline menus |
| **Discord** | Team/secondary | Medium+ | Webhooks to personal server |
| **Twilio Call** | Driving / Emergency | High+ | Speaks message audibly via TTS |
| **Home Speaker (TTS)** | Home, Critical | High+ | HA TTS via media_player |
| **iOS OpenClaw App** | iPhone push | Medium+ | Native push notification |
| **macOS OpenClaw App** | Desktop alert | Medium | macOS notification |
| **OneMind Web** | Active session | Any | In-app toast/banner |

### Speaker TTS (Home Assistant)

```yaml
# HA service call for speaker alert
service: tts.google_translate_say
# or: service: tts.cloud_say (Nabu Casa)
target:
  entity_id: media_player.home_speaker  # configure your speaker entity
data:
  message: "{{ message }}"
  language: "en-US"
```

**From Legacy (via HA API):**
```bash
curl -s -X POST "$HA_URL/api/services/tts/google_translate_say" \
  -H "Authorization: Bearer $HA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "media_player.home_speaker", "message": "Zeus, your server disk is at 94 percent."}'
```

---

## 7. Mode States

```
PASSIVE (default)
  ↓ Zeus requests "find me" or agent needs location context
ACTIVE (location aware)
  ↓ Urgent signal crosses threshold
ALERT (interrupt mode)
  ↓ Critical signal
EMERGENCY (all channels)

PASSIVE → context-aware responses, no output
ACTIVE  → Legacy uses location in every response
ALERT   → Legacy interrupts via awareness-filtered channels
EMERGENCY → All channels fire simultaneously
```

### Triggering Modes

```bash
# Zeus can set mode via Telegram to Legacy:
"find zeus mode on"     → Legacy enters ACTIVE mode
"quiet mode"            → awareness filter raises threshold to critical+
"emergency mode"        → all channels open
"find me"               → Legacy reports current state immediately

# NATS publish mode change:
nats pub zeus.presence.mode '{"mode": "active", "reason": "user_requested"}'
```

---

## 8. Implementation Phases

| Phase | What | Status |
|-------|------|--------|
| **1** | HA presence sensor structure (zones, person entity, device trackers) | ⏳ Configure HA |
| **2** | NATS pub/sub for zone changes (HA → NATS bridge) | ⏳ N8N or Python |
| **3** | Bouncie → NATS vehicle publisher (bouncie skill) | ✅ Skill defined |
| **4** | Legacy subscribes to `zeus.presence.>` + maintains ZEUS-STATE.md | ⏳ Workspace update |
| **5** | Awareness filter in Legacy context | ⏳ Document + deploy |
| **6** | Signals skill (urgent alerts, speaker TTS) | ✅ Skill defined (see signals/) |
| **7** | Output channel routing (Telegram, Twilio, speaker) | ⏳ Configure |
| **8** | ZEUS-STATE.md auto-update via presence-pulse service | ⏳ Build service |

---

## 9. presence-pulse Service (VPS)

A lightweight VPS service that polls Bouncie + HA and publishes to NATS:

```python
# /home/onemind/onemind-os/_system/openclaw/workspaces/legacy/skills/presence-pulse/presence-pulse.py
# Runs every 60s (vehicle active) / 300s (stationary)

import os, time, json, urllib.request
import nats  # nats.py library

NATS_URL = "nats://localhost:4222"
HA_URL = os.environ["HOME_ASSISTANT_URL"]
HA_TOKEN = os.environ["HOME_ASSISTANT_TOKEN"]
BOUNCIE_TOKEN = os.environ["BOUNCIE_ACCESS_TOKEN"]

async def poll_and_publish():
    nc = await nats.connect(NATS_URL)

    while True:
        # Poll HA person.zeus
        ha_state = get_ha_state("person.zeus")

        # Poll Bouncie
        bouncie_state = get_bouncie_location()

        # Publish both
        await nc.publish("zeus.presence.zone", json.dumps(ha_state).encode())
        if bouncie_state:
            await nc.publish("zeus.presence.vehicle", json.dumps(bouncie_state).encode())

        # Sleep: 60s if driving, 300s if home
        speed = bouncie_state.get("speed_mph", 0) if bouncie_state else 0
        await asyncio.sleep(60 if speed > 5 else 300)
```

---

## Related Files

| File | Purpose |
|------|---------|
| `.claude/skills/bouncie/SKILL.md` | Bouncie GPS API integration |
| `.claude/skills/transit/SKILL.md` | Transit routing with location awareness |
| `.claude/skills/home-assistant/SKILL.md` | HA API integration |
| `.claude/skills/signals/SKILL.md` | Alert + speaker notification system |
| `SO/03 Agents/OPENCLAW-INTERFACES.md` | All OpenClaw access interfaces |
| `_system/openclaw/workspaces/legacy/skills/presence-pulse/` | VPS polling service |
| `_system/openclaw/workspaces/legacy/ZEUS-STATE.md` | Live Zeus presence context |

---

*Part of SO — Sovereign Operations. Doctrine file — do not delete.*
