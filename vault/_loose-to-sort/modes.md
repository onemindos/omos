---
type: operating-modes-index
category: intelligence
updated: 2026-05-26
tags:
  - modes
  - operating-modes
  - context
  - intelligence
  - legacy
vault_path: "SO/01 Intelligence (The Brain)/modes/"
---

# OneMind OS — Operating Modes

**Operating Modes** = Contextual configurations that change how Legacy and OneMind OS behave based on Zeus's activity, location, security posture, or urgency level.

Modes are not personality presets. They change **what the system does** — automations, alert thresholds, data sources, agent behavior — not how it talks.

Full mode definitions live at: `00-24 SO (Sovereign Operations)/01 Intelligence (The Brain)/modes/`

---

## Mode Categories

### Activity Modes
*What Zeus is doing*

| Mode | Trigger | System Behavior |
|------|---------|----------------|
| **Big Daddy** | Family lead role | Family-oriented automations, protective posture |
| **Daddy** | Active parenting | Kid-safe context, schedule-aware |
| **Farmer** | Physical outdoor work | Low interruption, weather awareness |
| **Flying** | In flight | Offline-tolerant, minimal alerts |
| **Hunter** | Hunting / outdoor pursuit | Location privacy, outdoor context |
| **Learning** | Study / deep research | Focus mode, distraction suppression |
| **PR Live** | Public relations / media | Comms-ready, brand posture active |
| **Training Wheels** | Onboarding / delegation | Verbose output, extra confirmations |

### Location Modes
*Where Zeus is*

| Mode | Location | System Behavior |
|------|----------|----------------|
| **Home** | Greenwich CT residence | Full home automation, speaker alerts enabled |
| **Office** | Work location | Professional posture, office integrations |
| **NYC** | New York City | Transit-aware, no home speaker, Telegram primary |
| **Greenwich** | Greenwich area (not home) | Local context, nearby POIs |

### Security Modes
*System threat posture*

| Mode | Threat Level | System Behavior |
|------|-------------|----------------|
| **Incognito** | Privacy-first | Minimal data logging, reduced tracking |
| **Overwatch** | Heightened monitoring | All sensors active, tighter alert thresholds |
| **Locked Down** | Maximum security | Non-essential services suspended, strict access |

### Urgent Modes
*Active situation response*

| Mode | Situation | System Behavior |
|------|-----------|----------------|
| **Avalanche** | Overwhelm / backlog crisis | Triage-only, queue management, no new intake |
| **Combat** | Active high-stakes execution | Maximum focus, all non-critical alerts suppressed |
| **Danger** | Physical/security threat | Emergency protocols, all channels active |
| **Omni** | Full system activation | Every agent active, maximum awareness |
| **Zeus** | Zeus-tracking active (FIND ZEUS MODE) | See below |

---

## FIND ZEUS MODE

The centerpiece of the modes system. A persistent intelligence layer — Legacy is always aware of Zeus's location, context, and environment.

### States

```
PASSIVE  (default) — Legacy silently monitors. No output.
ACTIVE   (triggered) — Legacy uses presence as context for all responses.
ALERT    (urgent signal) — Legacy interrupts Zeus via awareness-filtered channel.
FIND     (requested) — "Where is Zeus?" → Legacy answers in <3 seconds.
EMERGENCY — All channels fire simultaneously.
```

### Data Sources

| Source | NATS Subject | Data |
|--------|-------------|------|
| Home Assistant | `zeus.presence.zone` | Zone resolution (home/nyc/office) |
| iPhone GPS | `zeus.presence.device` | Device state, battery |
| Bouncie GPS | `zeus.presence.vehicle` | Lat/lon, speed, heading |
| Legacy (computed) | `zeus.presence.context` | Enriched context string |
| Mode state | `zeus.presence.mode` | Current mode (passive/active/alert) |

### Output Channel Routing

| Channel | When Used | Min Urgency |
|---------|-----------|-------------|
| Telegram | Default always-on | Any |
| Home Speaker (TTS) | Home zone only | High+ |
| Twilio Call | Driving / Emergency | High+ |
| iOS Push | Away from Mac | Medium+ |
| macOS Alert | Active session | Medium |

### Triggering

```bash
# Via Telegram to Legacy:
"find zeus mode on"  → ACTIVE mode
"quiet mode"         → raise threshold to critical+
"emergency mode"     → all channels open
"find me"            → immediate state report

# Via NATS:
nats pub zeus.presence.mode '{"mode": "active", "reason": "user_requested"}'
```

### Implementation Status

| Phase | Task | Status |
|-------|------|--------|
| 1 | HA presence sensors (zones, person entity, device trackers) | ⏳ Pending |
| 2 | NATS pub/sub for zone changes (HA → NATS bridge) | ⏳ Pending |
| 3 | Bouncie → NATS vehicle publisher | ✅ Skill defined |
| 4 | Legacy subscribes to `zeus.presence.>` | ⏳ Pending |
| 5 | Awareness filter in Legacy context | ⏳ Pending |
| 6 | Output channel routing (Telegram, Twilio, speaker) | ⏳ Pending |

Full architecture doc: `SO/01 Intelligence (The Brain)/modes/urgent/zeus-mode/FIND-ZEUS-MODE.md`

---

## Mode Activation

Modes are activated via Legacy (Telegram, voice, CLI) or triggered automatically by sensors/automations.

```bash
# Activate a mode
"activate combat mode"
"switch to home mode"
"engage overwatch"

# Check current mode
"what mode am I in?"
"current operating mode"
```

Mode state is published to NATS: `zeus.presence.mode`

---

## Related Files

- **Full Mode Definitions:** `00-24 SO/01 Intelligence (The Brain)/modes/`
- **FIND ZEUS Architecture:** `SO/01 Intelligence (The Brain)/modes/urgent/zeus-mode/FIND-ZEUS-MODE.md`
- **NATS Infrastructure:** `SO/02 Infrastructure (The Body)/NATS-INFRASTRUCTURE.md`
- **Home Assistant Bridge:** `SO/02 Infrastructure (The Body)/services/home-assistant/HA-NATS-Bridge.md`
- **Agent Personalities:** `SO/06 Agents (The Fleet)/hermes/hermes-personalities.md`

---

**Total Modes:** 21 (8 activity + 4 location + 3 security + 5 urgent + 1 omnipresent)
**Status:** System defined — FIND ZEUS MODE implementation in progress ★
