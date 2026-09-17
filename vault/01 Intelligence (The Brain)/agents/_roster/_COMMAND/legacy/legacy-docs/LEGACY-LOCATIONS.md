---
id: root-00-20260113225954
title: "LEGACY LOCATIONS"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "LEGACY LOCATIONS"
---
# Legacy AI - Multi-Location Infrastructure Plan

> **Created:** 2025-01-06
> **Purpose:** Architecture for Legacy AI to manage multiple physical locations through Home Assistant

---

## Overview

Legacy-AI will serve as a unified intelligence layer across all of Zeus's properties, with Home Assistant providing the physical control layer at each location.

---

## Locations

| Location | Code | Network Type | HA Instance | Role |
|----------|------|--------------|-------------|------|
| **Greenwich Home** | `GRW` | Residential | HA Master | Primary hub, aggregates all |
| **NYC Office** | `NYC` | Commercial | HA Satellite | Work environment |
| **Mexico Farm** | `MEX` | Rural/Remote | HA Satellite | Vacation + agriculture |
| **Upstate Homestead** | `UPS` | Rural compound | HA Satellite | Multi-building complex |

---

## Architecture

### Federated Model

```
                              Legacy-AI (Letta Cloud)
                                      │
                                      │ Single MCP/API Connection
                                      ▼
                          ┌───────────────────────┐
                          │   MASTER HA (Greenwich)│
                          │   home.legacy.domain   │
                          │   HA Cloud: $6.50/mo   │
                          └───────────┬───────────┘
                                      │
                    Remote Home Assistant Integration
                                      │
          ┌───────────────┬───────────┴───────────┬───────────────┐
          ▼               ▼                       ▼               ▼
   ┌─────────────┐ ┌─────────────┐       ┌─────────────┐ ┌─────────────┐
   │  NYC Office │ │ Mexico Farm │       │  Homestead  │ │  (Future)   │
   │  HA Cloud   │ │  HA Cloud   │       │  HA Cloud   │ │             │
   │  $6.50/mo   │ │  $6.50/mo   │       │  $6.50/mo   │ │             │
   └─────────────┘ └─────────────┘       └──────┬──────┘ └─────────────┘
                                                │
                                    Single Network (UniFi)
                                                │
                          ┌──────────┬──────────┼──────────┬──────────┐
                          ▼          ▼          ▼          ▼          ▼
                     Main House  Greenhouse   Barn    3D Print   FarmBot
                                              House     Room
```

### Benefits of This Model

1. **Single point of access** - Legacy connects to ONE API endpoint
2. **Unified entity namespace** - All devices prefixed by location
3. **Centralized automations** - Cross-location rules in one place
4. **Resilient** - Each location works independently if master offline
5. **Scalable** - Add new locations without changing Legacy's config

---

## Location Details

### 1. Greenwich Home (MASTER)

**Role:** Primary residence, master aggregator hub

**Hardware:**
- Beelink Mini PC (i3, 8GB RAM, 256GB SSD)
- ATOM Echo (voice satellite)
- Yi Cameras (presence detection)
- SwitchBot devices (plugs, sensors)

**Network:** Standard residential WiFi/Ethernet

**Zones:**
- Living areas
- Office/workspace
- Bedrooms
- Outdoor/garage

**Entity Prefix:** `grw_` or none (as master, can be default)

---

### 2. NYC Office

**Role:** Work environment automation

**Hardware:**
- Mini PC (Intel N100 or similar, ~$150)
- Smart plugs for desk equipment
- Temperature/humidity sensors
- Camera for presence detection

**Network:** Office network (may need IT coordination)

**Zones:**
- Desk area
- Meeting room (if applicable)
- Common areas

**Entity Prefix:** `nyc_`

**Automations:**
- Arrival/departure detection
- Desk equipment power management
- Climate comfort
- "I'm at the office" status for Legacy

---

### 3. Mexico Farm

**Role:** Vacation property + agricultural monitoring

**Hardware:**
- Rugged mini PC (heat-resistant)
- Outdoor sensors (weather, soil moisture)
- Pool equipment controls
- Security cameras
- Gate/access controls

**Network:** Rural internet (may need Starlink backup)

**Zones:**
- Main house
- Pool area
- Garden/agriculture
- Perimeter/security

**Entity Prefix:** `mex_`

**Automations:**
- Remote monitoring when away
- Pool maintenance automation
- Security alerts
- Weather-based irrigation
- "Property check" routines

---

### 4. Upstate Homestead (Compound)

**Role:** Multi-building rural property with agriculture

**Hardware:**
- Central mini PC in main building
- UniFi network equipment (coverage across buildings)
- Per-building sensors and controls

**Network:** Single UniFi network spanning all buildings

**Buildings:**

#### Main House
- Standard smart home devices
- HVAC controls
- Security system

#### Greenhouse
- Climate sensors (temp, humidity, CO2)
- Automated ventilation
- Irrigation controls
- Grow lights
- FarmBot integration (if applicable)

#### Barn House
- Temperature monitoring
- Lighting controls
- Security cameras
- Equipment monitoring

#### 3D Printing Room
- OctoPrint integration (print monitoring)
- Cameras for print surveillance
- Temperature/ventilation
- Fire/smoke detection
- Power monitoring per printer

#### FarmBot Area
- FarmBot API integration
- Soil sensors
- Weather station
- Water tank levels

**Entity Prefix:** `ups_` with zone suffixes:
- `ups_main_*`
- `ups_greenhouse_*`
- `ups_barn_*`
- `ups_3dprint_*`
- `ups_farmbot_*`

---

## Network Architecture

### Per-Location Network

| Location | Recommended Setup | Internet |
|----------|-------------------|----------|
| Greenwich | Standard router/mesh | Cable/Fiber |
| NYC Office | Existing office network | Commercial |
| Mexico | UniFi + Starlink backup | Local ISP + Starlink |
| Homestead | UniFi (full deployment) | Rural fiber/Starlink |

### Homestead Network Detail

```
Internet (Fiber/Starlink)
         │
         ▼
   ┌───────────┐
   │  UDM Pro  │  ← Main router/controller
   │  (Main)   │
   └─────┬─────┘
         │
    Fiber/Cat6 backbone
         │
    ┌────┴────┬────────────┬────────────┐
    ▼         ▼            ▼            ▼
┌───────┐ ┌───────┐   ┌───────┐   ┌───────┐
│ Main  │ │Green- │   │ Barn  │   │ 3D    │
│ House │ │ house │   │ House │   │ Print │
│ AP    │ │ AP    │   │ AP    │   │ AP    │
└───────┘ └───────┘   └───────┘   └───────┘

All on same VLAN = One HA sees everything
```

---

## Legacy-AI Integration

### Why API-Based Control (Not HA Voice Assistant)

Legacy-AI controls Home Assistant via **REST API**, not as HA's internal voice assistant. This is intentional:

| Approach | How It Works | Pros | Cons |
|----------|--------------|------|------|
| **HA Voice Assistant** | ATOM Echo → HA STT → LLM → HA TTS | Local wake word | Tied to one location, streaming issues with Letta |
| **Legacy via API** | Any Interface → Letta → HA REST API | Control from anywhere, unified brain | No local wake word (use button or phone) |

**Our Choice: API-Based Control**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           VOICE FROM ANYWHERE                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐                                                           │
│   │ Vapi Phone  │─────┐                                                     │
│   │ (203)760-   │     │                                                     │
│   └─────────────┘     │                                                     │
│                       │         ┌─────────────┐                             │
│   ┌─────────────┐     │         │             │                             │
│   │ LiveKit Web │─────┼────────→│  LEGACY-AI  │                             │
│   └─────────────┘     │         │   (Letta)   │                             │
│                       │         │             │                             │
│   ┌─────────────┐     │         │  Memory:    │                             │
│   │ Discord     │─────┘         │  - Locations│                             │
│   └─────────────┘               │  - Devices  │                             │
│                                 │  - Prefs    │                             │
│                                 └──────┬──────┘                             │
│                                        │                                     │
│                                        │ REST API                            │
│                                        ▼                                     │
│   ┌────────────────────────────────────────────────────────────────────┐   │
│   │                    FEDERATED HOME ASSISTANT                         │   │
│   │                                                                     │   │
│   │   ┌─────────────┐                                                   │   │
│   │   │  GREENWICH  │ ← Master (all entities aggregated here)          │   │
│   │   │   (Master)  │                                                   │   │
│   │   └──────┬──────┘                                                   │   │
│   │          │ Remote HA Integration                                    │   │
│   │    ┌─────┼─────┬─────────────┐                                      │   │
│   │    ▼     ▼     ▼             ▼                                      │   │
│   │  ┌───┐ ┌───┐ ┌───┐       ┌───────┐                                  │   │
│   │  │NYC│ │MEX│ │UPS│       │Future │                                  │   │
│   │  └───┘ └───┘ └───┘       └───────┘                                  │   │
│   │                                                                     │   │
│   └────────────────────────────────────────────────────────────────────┘   │
│                                        │                                     │
│                                        ▼                                     │
│   ┌────────────────────────────────────────────────────────────────────┐   │
│   │                         PHYSICAL DEVICES                            │   │
│   │                                                                     │   │
│   │  grw_office_light    nyc_desk_plug    mex_pool_pump                │   │
│   │  grw_speaker_echo    nyc_sensor_temp  mex_gate_main                │   │
│   │  ups_greenhouse_*    ups_3dprint_*    ups_farmbot_*                │   │
│   └────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Single Connection Architecture

```python
# Legacy-AI only needs ONE configuration
MASTER_HA = {
    "url": "https://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",  # Long-lived token
}

# All locations accessible through master
# Entity naming convention:
# - grw_living_room_light
# - nyc_desk_plug
# - mex_pool_pump
# - ups_greenhouse_temp
# - ups_3dprint_printer1_status
```

---

## Legacy Memory: Location Registry

Legacy stores location information in Letta's **core memory** for instant access:

```python
# Stored in Legacy-AI's core memory block
LOCATIONS = {
    "grw": {
        "name": "Greenwich Home",
        "code": "GRW",
        "role": "primary_residence",
        "ha_url": "https://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa",
        "zones": ["office", "living_room", "bedroom", "kitchen", "garage", "outdoor"],
        "speakers": ["media_player.grw_speaker_echo"],
        "presence_sensor": "binary_sensor.grw_presence_zeus",
        "is_master": True
    },
    "nyc": {
        "name": "NYC Office",
        "code": "NYC",
        "role": "work",
        "zones": ["main", "desk", "meeting"],
        "speakers": ["media_player.nyc_speaker_desk"],
        "presence_sensor": "binary_sensor.nyc_presence_zeus",
        "is_master": False
    },
    "mex": {
        "name": "Mexico Farm",
        "code": "MEX",
        "role": "vacation",
        "zones": ["main_house", "pool", "garden", "security"],
        "speakers": ["media_player.mex_speaker_main"],
        "presence_sensor": "binary_sensor.mex_presence_zeus",
        "is_master": False
    },
    "ups": {
        "name": "Upstate Homestead",
        "code": "UPS",
        "role": "compound",
        "zones": ["main", "greenhouse", "barn", "3dprint", "farmbot"],
        "speakers": ["media_player.ups_main_speaker", "media_player.ups_barn_speaker"],
        "presence_sensor": "binary_sensor.ups_presence_zeus",
        "is_master": False
    }
}

# Device type mappings for natural language
DEVICE_TYPES = {
    "light": ["light"],
    "plug": ["switch"],
    "sensor": ["sensor", "binary_sensor"],
    "camera": ["camera"],
    "speaker": ["media_player"],
    "climate": ["climate"],
    "lock": ["lock"],
    "cover": ["cover"]  # blinds, garage doors
}

# Location-specific device hints (for disambiguation)
LOCATION_HINTS = {
    "greenhouse": "ups",
    "farmbot": "ups",
    "3d print": "ups",
    "pool": "mex",
    "farm": "mex",
    "office": "nyc",  # unless "home office" specified
    "desk": "nyc"
}
```

---

## Legacy Location Awareness

### How Legacy Knows Where Zeus Is

```python
def get_zeus_location() -> dict:
    """
    Determine Zeus's current location using multiple methods.

    Returns:
        {
            "location": "grw",
            "location_name": "Greenwich Home",
            "confidence": 0.95,
            "method": "presence_sensor",
            "last_updated": "2026-01-07T10:30:00Z"
        }
    """
    # Priority 1: HA presence detection at each location
    for loc_code, loc_info in LOCATIONS.items():
        presence = get_ha_state(loc_info["presence_sensor"])
        if presence == "on":
            return {
                "location": loc_code,
                "location_name": loc_info["name"],
                "confidence": 0.95,
                "method": "presence_sensor"
            }

    # Priority 2: Check calendar for scheduled location
    calendar_location = check_calendar_location()
    if calendar_location:
        return {
            "location": calendar_location,
            "confidence": 0.7,
            "method": "calendar"
        }

    # Priority 3: Last known location from memory
    last_known = get_from_memory("zeus_last_location")
    if last_known:
        return {
            "location": last_known,
            "confidence": 0.5,
            "method": "last_known"
        }

    return {"location": "away", "confidence": 0.3, "method": "default"}
```

### Smart Command Routing

```python
def parse_location_from_command(command: str) -> str:
    """
    Extract location from natural language command.

    Examples:
        "Turn on the Greenwich lights" → "grw"
        "What's the pool temperature?" → "mex" (pool = Mexico)
        "Turn on the lights" → current location
    """
    command_lower = command.lower()

    # Check for explicit location mention
    for loc_code, loc_info in LOCATIONS.items():
        if loc_info["name"].lower() in command_lower:
            return loc_code
        if loc_code in command_lower:
            return loc_code

    # Check for location-specific device hints
    for hint, loc_code in LOCATION_HINTS.items():
        if hint in command_lower:
            return loc_code

    # Default to current location
    return get_zeus_location()["location"]

def route_command(command: str, explicit_location: str = None):
    """
    Route a command to the appropriate location.

    Examples:
        "Turn on the living room light"
            → If at GRW: grw_living_room_light
            → If at NYC: error (no living room at office)

        "Turn on the Greenwich porch light"
            → Always routes to grw_outdoor_light_porch
    """
    location = explicit_location or parse_location_from_command(command)

    # Parse device from command
    device = parse_device_from_command(command)

    # Construct entity_id
    entity_id = f"{location}_{device['zone']}_{device['type']}_{device['name']}"

    return execute_ha_command(entity_id, device["action"])
```

---

## Voice Control From Anywhere

### The Key Advantage

Since Legacy uses Vapi (phone) and controls HA via API, you can control **any location from anywhere**:

```
┌─────────────────────────────────────────────────────────────┐
│                  EXAMPLE: You're in Paris                   │
│                                                             │
│  You: "Hey Legacy, turn on the Greenwich porch light        │
│        and check if the Mexico pool pump is running"        │
│                                                             │
│  Legacy:                                                    │
│    1. Parses: Greenwich → grw, porch light                  │
│    2. Calls: POST /api/services/light/turn_on               │
│       Body: {"entity_id": "light.grw_outdoor_light_porch"}  │
│                                                             │
│    3. Parses: Mexico → mex, pool pump                       │
│    4. Calls: GET /api/states/switch.mex_pool_pump_main      │
│                                                             │
│    5. Responds: "I've turned on the Greenwich porch light.  │
│       The Mexico pool pump is currently running."           │
└─────────────────────────────────────────────────────────────┘
```

### Communication Routing

```python
def speak_to_zeus(message: str, urgency: str = "normal"):
    """
    Route communication to Zeus at the correct location.

    Args:
        message: What Legacy wants to say
        urgency: "normal", "important", or "critical"
    """
    location_info = get_zeus_location()
    location = location_info["location"]

    if location != "away":
        # Zeus is at a known location - speak through local speaker
        speaker = LOCATIONS[location]["speakers"][0]
        call_ha_service("tts.speak", {
            "entity_id": "tts.home_assistant_cloud",
            "media_player_entity_id": speaker,
            "message": message
        })
        return f"Spoke to Zeus at {LOCATIONS[location]['name']}"

    elif urgency == "critical":
        # Zeus is away and it's critical - call their phone
        initiate_vapi_call(ZEUS_PHONE, message)
        return "Called Zeus's phone"

    else:
        # Zeus is away, not critical - send notification
        send_push_notification(message)
        return "Sent push notification"
```

---

## Letta Tools for Multi-Location HA Control

### Core Tools (Already Created)

```python
# Tool: get_ha_states
def get_ha_states(location: str = None, device_type: str = None) -> dict:
    """
    Get device states from Home Assistant.

    Args:
        location: Filter by location code (grw, nyc, mex, ups).
                  If None, returns all locations.
        device_type: Filter by type (light, switch, sensor, etc.)

    Returns:
        Dictionary of entity states

    Examples:
        get_ha_states() → All devices, all locations
        get_ha_states(location="grw") → Greenwich devices only
        get_ha_states(device_type="light") → All lights
        get_ha_states(location="mex", device_type="sensor") → Mexico sensors
    """

# Tool: control_ha_device
def control_ha_device(entity_id: str, action: str, **params) -> str:
    """
    Control any Home Assistant device.

    Args:
        entity_id: The entity to control (e.g., "light.grw_office_light_desk")
        action: Action to perform (on, off, toggle, set_temperature, etc.)
        **params: Additional parameters (brightness, temperature, etc.)

    Examples:
        control_ha_device("light.grw_office_light_desk", "on")
        control_ha_device("light.grw_office_light_desk", "on", brightness=50)
        control_ha_device("climate.nyc_thermostat", "set_temperature", temperature=72)
    """

# Tool: speak_in_room
def speak_in_room(message: str, location: str = None, room: str = None) -> str:
    """
    Make Legacy speak through a Home Assistant speaker.

    Args:
        message: What to say
        location: Location code (auto-detected if not specified)
        room: Specific room/speaker (uses default if not specified)

    Examples:
        speak_in_room("Dinner is ready") → Speaks at current location
        speak_in_room("Package arrived", location="grw") → Greenwich speaker
    """
```

### Extended Tools (To Be Added)

```python
# Tool: trigger_automation
def trigger_automation(automation_id: str) -> str:
    """
    Trigger a Home Assistant automation.

    Args:
        automation_id: The automation entity ID (e.g., "automation.grw_morning_routine")

    Examples:
        trigger_automation("automation.grw_morning_routine")
        trigger_automation("automation.mex_security_check")
    """

# Tool: run_script
def run_script(script_id: str, **variables) -> str:
    """
    Run a Home Assistant script.

    Args:
        script_id: The script entity ID (e.g., "script.grw_movie_mode")
        **variables: Variables to pass to the script

    Examples:
        run_script("script.grw_movie_mode")
        run_script("script.ups_start_print", printer="printer1", file="benchy.gcode")
    """

# Tool: activate_scene
def activate_scene(scene_id: str) -> str:
    """
    Activate a Home Assistant scene.

    Args:
        scene_id: The scene entity ID (e.g., "scene.grw_dinner_time")

    Examples:
        activate_scene("scene.grw_dinner_time")
        activate_scene("scene.nyc_focus_mode")
    """

# Tool: get_location_summary
def get_location_summary(location: str) -> dict:
    """
    Get a summary of a location's status.

    Returns:
        {
            "location": "Greenwich Home",
            "temperature": 72,
            "humidity": 45,
            "lights_on": 3,
            "security": "armed_home",
            "occupancy": "home",
            "last_motion": "5 minutes ago"
        }
    """

# Tool: where_is_zeus
def where_is_zeus() -> dict:
    """
    Determine Zeus's current location.

    Returns:
        {
            "location": "grw",
            "location_name": "Greenwich Home",
            "confidence": 0.95,
            "method": "presence_sensor"
        }
    """
```

---

## Natural Language Examples

| You Say | Legacy Understands | Action |
|---------|-------------------|--------|
| "Turn on the lights" | Current location → default zone → lights | `light.{loc}_*` turn on |
| "Turn on Greenwich lights" | GRW → all lights | `light.grw_*` turn on |
| "What's the pool temperature?" | Pool = MEX → sensor | Query `sensor.mex_pool_temp` |
| "Start the movie mode" | Current location → scene | `scene.{loc}_movie_mode` |
| "Run the morning routine" | Current location → automation | `automation.{loc}_morning_routine` |
| "Is anyone at the office?" | NYC → presence | Query `binary_sensor.nyc_presence_*` |
| "How's the greenhouse?" | Greenhouse = UPS | Query `sensor.ups_greenhouse_*` |
| "Start a 3D print" | 3D print = UPS | Trigger OctoPrint integration |
| "Set away mode everywhere" | All locations | Trigger security automations |

---

## Implementation Phases

### Phase 1: Greenwich (NOW)
- [x] Acquire Beelink hardware
- [ ] Install HA OS via Ubuntu method
- [ ] Subscribe to HA Cloud
- [ ] Set up SSH access for Claude
- [ ] Add local devices (ATOM Echo, Yi cameras, SwitchBot)
- [ ] Create Long-Lived Access Token
- [ ] Connect Legacy-AI via REST API
- [ ] Test basic automations

### Phase 2: Foundation
- [ ] Establish entity naming convention
- [ ] Create backup system
- [ ] Document all integrations
- [ ] Set up monitoring/alerts

### Phase 3: NYC Office
- [ ] Acquire mini PC for office
- [ ] Install HA OS
- [ ] Subscribe to HA Cloud (#2)
- [ ] Connect to Master via Remote HA integration
- [ ] Add office devices
- [ ] Create office-specific automations

### Phase 4: Mexico Farm
- [ ] Site survey for network needs
- [ ] Install Starlink (if needed for backup)
- [ ] Acquire and deploy hardware
- [ ] Subscribe to HA Cloud (#3)
- [ ] Connect to Master
- [ ] Focus on security + remote monitoring

### Phase 5: Upstate Homestead
- [ ] Network infrastructure (UniFi deployment)
- [ ] Central HA installation
- [ ] Subscribe to HA Cloud (#4)
- [ ] Connect buildings to network
- [ ] Per-building device deployment
- [ ] FarmBot integration
- [ ] OctoPrint integration
- [ ] Agricultural automations

---

## Cost Summary

### Monthly Recurring

| Item | Cost |
|------|------|
| HA Cloud - Greenwich | $6.50 |
| HA Cloud - NYC | $6.50 |
| HA Cloud - Mexico | $6.50 |
| HA Cloud - Homestead | $6.50 |
| **Total HA Cloud** | **$26.00/mo** |

### One-Time Hardware (Estimated)

| Location | Hardware | Est. Cost |
|----------|----------|-----------|
| Greenwich | Beelink (owned) + devices | ~$200 |
| NYC Office | Mini PC + sensors | ~$250 |
| Mexico | Mini PC + outdoor sensors | ~$400 |
| Homestead | Mini PC + UniFi network | ~$1,500 |
| **Total Hardware** | | **~$2,350** |

---

## Integrations Per Location

### All Locations
- Home Assistant Cloud (remote access)
- Legacy-AI (via master)
- Mobile app notifications

### Greenwich
- SwitchBot
- Yi Cameras / Frigate
- ATOM Echo (voice)
- LLM Vision

### NYC Office
- Calendar sync
- Presence detection
- Equipment power management

### Mexico
- Weather integration
- Security cameras
- Pool automation
- Starlink monitoring

### Homestead
- FarmBot API
- OctoPrint
- Weather station
- Agricultural sensors
- UniFi network monitoring

---

## Security Considerations

1. **All HA instances behind HA Cloud** - No direct port exposure
2. **Long-Lived Tokens** - Rotate periodically
3. **User accounts per location** - Audit trail
4. **VPN backup** - Tailscale for emergency access
5. **Local fallback** - Each location works independently

---

## Backup Strategy

### Per-Instance Backups
- HA Cloud automatic backups (included)
- Local backup to USB/NAS
- Critical config in git repo

### Master Backup
- Full configuration export weekly
- Automation definitions in version control
- Credentials in secure vault (1Password/Bitwarden)

---

## Related Documentation

- [LEGACY-PLAN.md](LEGACY-PLAN.md) - Master Legacy AI plan
- [LEGACY-COMMS.md](LEGACY-COMMS.md) - Communications system
- [LEGACY-AI-SYSTEM.md](LEGACY-AI-SYSTEM.md) - Core AI documentation

---

## Quick Reference

### URLs (After Setup)

| Location | Local URL | Remote URL |
|----------|-----------|------------|
| Greenwich | `http://homeassistant.local:8123` | `https://[id].ui.nabu.casa` |
| NYC | `http://ha-nyc.local:8123` | `https://[id].ui.nabu.casa` |
| Mexico | `http://ha-mex.local:8123` | `https://[id].ui.nabu.casa` |
| Homestead | `http://ha-ups.local:8123` | `https://[id].ui.nabu.casa` |

### Entity Naming Convention

```
{location}_{zone}_{device_type}_{name}

Examples:
- grw_office_light_desk
- nyc_main_plug_monitor
- mex_pool_pump_main
- ups_greenhouse_sensor_temp
- ups_3dprint_camera_printer1
```
