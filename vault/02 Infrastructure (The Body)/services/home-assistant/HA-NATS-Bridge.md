---
title: Home Assistant — NATS Bridge Architecture
domain: 02 Infrastructure (The Body)
service: home-assistant
status: live
last-updated: 2026-05-24
tags: [home-assistant, nats, ha-bridge, smart-home, switchbot, onemind-os]
related: [[NATS]], [[NATS-Intent-Mesh]]
---

# Home Assistant — NATS Bridge Architecture

## Overview

Home Assistant (omos-ha) is integrated into OneMind OS via a **persistent NATS bridge** (`ha-bridge.js`) running on the MBP. All HA commands and state reads go through the NATS bus — no direct HTTP, no MCP wrappers, no cloud relays.

Core equation: **`nats_request` is the only tool for HA**. Same pattern as SaaS (`saas.execute`).

---

## Architecture

```
Legacy (Hermes) 
    │
    │  nats_request("cmd.ha.*")
    ▼
NATS Server (localhost:4222)
    │
    │  JetStream subject routing
    ▼
ha-bridge.js (~/.nats/services/ha-bridge/bridge.js)
    │  launchd: com.onemind.ha-bridge
    │  persistent WebSocket → 192.168.1.180:8123
    │  183+ entity states cached in memory
    ▼
Home Assistant (omos-ha, 192.168.1.180 / 100.64.25.12)
```

### Why This Topology

| Option | RTT | Status |
|--------|-----|--------|
| MCP → Nabu Casa cloud | ~348ms | ❌ Removed |
| `homeassistant` toolset (direct HTTP) | ~3.6ms | ❌ Disabled |
| **NATS → ha-bridge → HA WebSocket** | **~1–5ms** | ✅ Active |
| Pure NATS bus overhead | ~0.4ms | — |

**Worker on MBP (not omos-ha):** NATS server is on MBP → sub-ms bus hop. HA WebSocket is persistent from MBP → LAN (`192.168.1.180:8123`, not Tailscale). Entity states cached in-memory → reads never hit HA network. Worker on omos-ha would add a Tailscale hop (3ms each way) with zero benefit.

**WebSocket over REST:** Persistent connection → no TCP setup per call. Real-time state push → cache always fresh without polling. Service calls go through the same connection.

---

## NATS Subjects

| Subject | RTT | Payload | Returns |
|---------|-----|---------|---------|
| `cmd.ha.state.get` | ~1ms | `{"entity_id": "sun.sun"}` | Full entity state object (from cache) |
| `cmd.ha.states.list` | ~3ms | `{}` or `{"domain": "switch"}` | `{"entities": [...], "count": N}` |
| `cmd.ha.service.call` | ~5ms | `{"domain", "service", "entity_id", "data?"}` | `{"ok": true}` |
| `events.ha.*` | publish only | — | State change events (subscribe for reactive patterns) |

---

## Usage (nats_request)

```python
# Read a single state
nats_request(subject="cmd.ha.state.get", payload={"entity_id": "switch.plug_mini_a886"})

# List all entities
nats_request(subject="cmd.ha.states.list", payload={})

# List by domain
nats_request(subject="cmd.ha.states.list", payload={"domain": "switch"})

# Call a service
nats_request(subject="cmd.ha.service.call", payload={
    "domain": "switch",
    "service": "turn_on",
    "entity_id": "switch.plug_mini_a886"
})
# Returns: {"ok": true}
```

---

## ha-bridge Service

| Property | Value |
|----------|-------|
| Script | `~/.nats/services/ha-bridge/bridge.js` |
| launchd label | `com.onemind.ha-bridge` |
| HA WebSocket | `ws://192.168.1.180:8123` (LAN, not Tailscale) |
| Entity cache | 185 entities in memory |

**Ops:**
```bash
# Check status
launchctl list | grep ha-bridge

# Restart
launchctl kickstart -k gui/$(id -u)/com.onemind.ha-bridge

# Logs
tail -f ~/.nats/services/ha-bridge/logs/stdout.log
```

---

## What Was Removed (2026-05-24)

| Component | What | Why |
|-----------|------|-----|
| `homeassistant` toolset | `hermes tools disable homeassistant` | Fights NATS architecture |
| `home-assistant` MCP block | Removed from `~/.hermes/config.yaml` | 348ms cloud overhead, redundant |
| Nabu Casa cloud URL | Was `HASS_HOST` in config | Replaced by LAN IP via Tailscale then bridge |

`~/.hermes/.env` still has `HASS_URL` and `HASS_TOKEN` — kept as reference for `homeassistant_tool.py` (disabled on disk). Do not use for new code.

---

## Device Inventory (as of 2026-05-24)

### Switches (SwitchBot Plug Mini)

| Entity | State | Notes |
|--------|-------|-------|
| `switch.plug_mini_a886` | ON, 0.0W | Active |
| `switch.plug_mini_ef5a` | OFF | Left off after 3x cycle benchmark |
| `switch.plug_mini_6826` | **unavailable** | Hub range issue — power-cycle hub |

### Sensors (SwitchBot Meter Plus)

| Entity | Value |
|--------|-------|
| `sensor.meter_plus_7616_temperature` | 75.2°F |
| `sensor.meter_plus_7616_humidity` | 60% |
| `sensor.meter_plus_ceff_temperature` | 79.5°F |
| `sensor.meter_plus_ceff_humidity` | 52% |

### Zeus iPhone

| Entity | Value |
|--------|-------|
| `sensor.zeus_iphone_battery_level` | 100% |
| `sensor.zeus_iphone_wifi_network` | onemindos |
| `sensor.zeus_iphone_storage_available` | 5.68% ⚠️ critically low |
| `device_tracker.zeus_iphone` | home |

### Unavailable / Unknown Entities

- `binary_sensor.*` — bedroom motion, keypad, lab entry, lab leak, zeus keychain: all `unavailable`
- `alarm_control_panel.alarm_control_panel`: `unavailable`
- HA backup sensors: all `unknown` — scheduled backups not running

---

## HA Instance Info

| Property | Value |
|----------|-------|
| LAN IP | `192.168.1.180:8123` |
| Tailscale IP | `100.64.25.12:8123` |
| HA Version | 2026.4.4 (2026.5.3 available, `auto_update: false`) |
| Supervisor | 2026.05.0 (current) |
| Total entities | 185 |

---

## Benchmark (2026-05-24)

```
Test: 3x ON/OFF cycle on switch.plug_mini_ef5a
ON  avg: 1.6ms
OFF avg: 3.2ms

NATS pure bus overhead:   0.4ms
ha-bridge state.get:      ~1ms (cached, zero HA network)
ha-bridge states.list:    ~3ms (183 entities from cache)
ha-bridge service.call:   ~5ms (WebSocket to HA)
```

---

## Pitfalls

- **ha-bridge must be running** — if `cmd.ha.*` times out, check `launchctl list | grep ha-bridge`
- **Never use direct HTTP for HA** — the bridge is the only approved path
- **Do NOT re-enable `homeassistant` toolset** — it bypasses NATS and fights the architecture
- **Do NOT add `home-assistant` MCP block back to config** — cloud relay was 348ms
- **ha-bridge uses LAN IP `192.168.1.180`** — not Tailscale, not Nabu Casa. Requires MBP and omos-ha on same WiFi (`onemindos`)
- **Leaf node on omos-ha** — config has `leafnodes { port: 7422 }` ready but not yet needed. Bridge on MBP is the right topology until MBP-down resilience becomes a requirement

---

## Real-Time Events

Subscribe to `events.ha.*` for reactive patterns (no polling):

```python
# All HA events
nats_subscribe(subject="events.ha.>")

# Specific entity
nats_subscribe(subject="events.ha.switch.plug_mini_a886")
```

Events are published by ha-bridge whenever HA pushes a WebSocket state change.

---

## Related Notes

- [[NATS]] — central bus architecture
- [[NATS-Intent-Mesh]] — full subject namespace map
- [[Hermes-Overview]] — Legacy/Hermes agent layer

## Related Skills (Hermes)

- `home-assistant` skill — canonical usage patterns and pitfalls
- `nats-server` skill — benchmark data, ha-bridge architecture, full device inventory
