---
id: root-00-20260113225954
title: "LEGACY HA REFERENCE"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "LEGACY HA REFERENCE"
---
# Legacy AI - Home Assistant Reference

> **Created:** 2026-01-06
> **Purpose:** Complete Home Assistant capabilities reference and API documentation for Legacy-AI integration

---

## Connection Details

### API Access (For Claude/Legacy-AI)

```
NABU CASA URL:    https://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa
LOCAL URL:        http://homeassistant.local:8123
API ENDPOINT:     {URL}/api/
WEBSOCKET:        wss://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa/api/websocket

LONG-LIVED TOKEN: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJiMzY3Zjc1MzQ4MWI0YWZlODBlOGFlODk0ZjE5YmQ2NCIsImlhdCI6MTc2NzczODk4MiwiZXhwIjoyMDgzMDk4OTgyfQ.C__KGAQFBHRUxB4j-Mg7uu4NnbQv2WntovyZLNzKXeE

TOKEN EXPIRES:    2036-01-06 (10 year token)
```

### Current Instance Info

| Property | Value |
|----------|-------|
| **Version** | 2025.12.5 |
| **Location** | Greenwich, CT (41.01, -73.64) |
| **Time Zone** | America/New_York |
| **State** | Running |
| **Config Source** | Storage |
| **Unit System** | Imperial (US) |

---

## Current Entity Inventory

### Summary (32 entities)

| Domain | Count | Purpose |
|--------|-------|---------|
| sensor | 13 | Sun, backup, weather, energy |
| update | 7 | Core, Supervisor, OS, add-ons |
| binary_sensor | 2 | Remote UI, router WAN |
| tts | 2 | HA Cloud, Google Translate |
| weather | 1 | Met.no forecast |
| sun | 1 | Sun position |
| person | 1 | zeuslegacy |
| zone | 1 | Home |
| todo | 1 | Shopping list |
| conversation | 1 | HA voice assistant |
| stt | 1 | HA Cloud speech-to-text |
| event | 1 | Backup events |

### Key Entities

```yaml
# TTS (Text-to-Speech)
tts.home_assistant_cloud       # Cloud-based TTS (included in subscription)
tts.google_translate_en_com    # Free Google TTS

# STT (Speech-to-Text)
stt.home_assistant_cloud       # Cloud-based STT (included in subscription)

# Conversation
conversation.home_assistant    # Voice assistant processing

# Remote Access
binary_sensor.remote_ui        # HA Cloud connection status (currently: on)

# Network
binary_sensor.verizon_router_wan_status  # Internet connection (currently: on)

# Weather
weather.forecast_home          # Met.no weather (currently: partlycloudy)

# Updates
update.home_assistant_core_update
update.home_assistant_supervisor_update
update.home_assistant_operating_system_update
update.music_assistant_update
update.let_s_encrypt_update
```

---

## Installed Add-ons

| Add-on | Purpose | Status |
|--------|---------|--------|
| **Music Assistant** | Multi-room audio, streaming services | Installed |
| **Let's Encrypt** | SSL certificates | Installed |

### Recommended Add-ons to Install

| Add-on | Purpose | Priority |
|--------|---------|----------|
| **SSH & Web Terminal** | Remote CLI access for Claude | HIGH |
| **File Editor** | Edit config files in browser | Medium |
| **Samba Share** | Network file access | Medium |
| **ESPHome** | Flash ATOM Echo and other ESP devices | HIGH |
| **HACS** | Community integrations store | HIGH |
| **Frigate** | NVR with AI object detection | Future |
| **Node-RED** | Advanced automation flows | Optional |

---

## Backup System

### Automatic Backups (HA Cloud)

HA Cloud includes automatic daily backups stored in the cloud.

**Configuration:**
```yaml
# Settings → System → Backups → Configure automatic backups
Schedule: Daily at 3:00 AM
Retention: 3 backups
Location: Home Assistant Cloud
Contents: Full (config + add-ons + media)
```

**Entities:**
```yaml
sensor.backup_backup_manager_state              # idle, creating, uploading
sensor.backup_next_scheduled_automatic_backup   # Next backup time
sensor.backup_last_successful_automatic_backup  # Last success
sensor.backup_last_attempted_automatic_backup   # Last attempt
event.backup_automatic_backup                   # Backup events
```

### Manual Backups via API

```python
# Create full backup
POST /api/services/backup/create_automatic
{
    "include_addons": true,
    "include_database": true,
    "include_folders": ["config", "media", "share"]
}

# List backups
GET /api/backup/backups

# Download specific backup
GET /api/backup/download/{backup_id}

# Restore from backup
POST /api/hassio/backups/{backup_id}/restore
```

### Backup Locations

| Location | Method | Pros | Cons |
|----------|--------|------|------|
| **HA Cloud** | Automatic | Off-site, encrypted, easy restore | Requires subscription |
| **Local USB** | Manual mount | Fast, no internet needed | Can fail with hardware |
| **NAS/SMB** | Network share | Large storage, versioning | Network dependency |
| **Google Drive** | Add-on | Free 15GB, off-site | Add-on required |
| **Dropbox** | Add-on | Easy setup | Limited free tier |

### Backup Strategy Recommendation

```
Daily:   HA Cloud automatic (keep 7)
Weekly:  Full backup to NAS (keep 4)
Monthly: Full backup to Google Drive (keep 12)
Before updates: Manual snapshot
```

---

## Update System

### Update Channels

| Component | Current | Update Method |
|-----------|---------|---------------|
| **Core** | 2025.12.5 | Settings → System → Updates |
| **Supervisor** | Latest | Auto-updates |
| **OS** | 16.x | Settings → System → Updates |
| **Add-ons** | Various | Add-on store |

### Update via API

```python
# Check for updates
GET /api/states/update.home_assistant_core_update
# Returns: {"state": "off"} = no update, {"state": "on"} = update available

# Get update details
GET /api/states/update.home_assistant_core_update
# attributes.latest_version, attributes.installed_version

# Install update
POST /api/services/update/install
{
    "entity_id": "update.home_assistant_core_update"
}

# Skip update
POST /api/services/update/skip
{
    "entity_id": "update.home_assistant_core_update"
}
```

### Update Notifications

Create automation to notify Legacy when updates are available:

```yaml
automation:
  - alias: "Notify Legacy of HA Updates"
    trigger:
      - platform: state
        entity_id: update.home_assistant_core_update
        to: "on"
    action:
      - service: notify.notify
        data:
          message: "Home Assistant {{ state_attr('update.home_assistant_core_update', 'latest_version') }} is available"
```

### Update Best Practices

1. **Read release notes** before updating
2. **Create backup** before major updates
3. **Update add-ons first**, then Core
4. **Wait 2-3 days** after release for bug reports
5. **Test automations** after updates

---

## Network Configuration

### Current Network

| Component | Status | Details |
|-----------|--------|---------|
| **WAN** | Connected | Verizon router |
| **Remote Access** | HA Cloud | Nabu Casa tunnel |
| **Local Access** | mDNS | homeassistant.local |

### Network Architecture

```
Internet
    │
    ▼
┌─────────────┐
│ Verizon     │
│ Router      │
└──────┬──────┘
       │
       ├── Beelink (HA) ─── homeassistant.local:8123
       │
       ├── Zeus's Mac
       │
       ├── Yi Cameras (future)
       │
       ├── SwitchBot Hub (future)
       │
       └── ATOM Echo (future)
```

### Remote Access Options

| Method | Setup | Security | Speed |
|--------|-------|----------|-------|
| **HA Cloud (current)** | Subscription | Encrypted tunnel | Good |
| **Tailscale** | Free add-on | WireGuard VPN | Excellent |
| **Cloudflare Tunnel** | Free | Zero-trust | Good |
| **Port Forward + DuckDNS** | Free | Needs SSL | Direct |
| **Nginx Proxy Manager** | Self-hosted | Full control | Direct |

### Tailscale Setup (Recommended Backup)

Even with HA Cloud, Tailscale provides direct VPN access:

```yaml
# Install Tailscale add-on
# Settings → Add-ons → Search "Tailscale"

# Benefits:
# - Direct connection (faster than cloud relay)
# - Works even if HA Cloud is down
# - Access other devices on same Tailnet
# - SSH access to HA OS
```

### Network Diagnostics via API

```python
# Check network status
GET /api/hassio/network/info

# Response includes:
# - interfaces (eth0, wlan0)
# - IP addresses
# - DNS servers
# - Gateway
```

---

## Available Services (API Actions)

### Core Services

```python
# Turn any device on/off
POST /api/services/homeassistant/turn_on
POST /api/services/homeassistant/turn_off
{"entity_id": "switch.office_plug"}

# Toggle device
POST /api/services/homeassistant/toggle
{"entity_id": "light.bedroom"}

# Restart Home Assistant
POST /api/services/homeassistant/restart

# Reload configuration
POST /api/services/homeassistant/reload_core_config
```

### TTS Services (Legacy Speaking)

```python
# HA Cloud TTS (best quality, requires subscription)
POST /api/services/tts/cloud_say
{
    "entity_id": "media_player.living_room",  # Target speaker
    "message": "Mr. Delacruz, you have a meeting in 10 minutes",
    "language": "en-US",
    "options": {
        "voice": "JennyNeural"  # Azure voice
    }
}

# Google Translate TTS (free, lower quality)
POST /api/services/tts/google_translate_say
{
    "entity_id": "media_player.office",
    "message": "Hello from Legacy"
}

# Generic TTS (uses default engine)
POST /api/services/tts/speak
{
    "entity_id": "tts.home_assistant_cloud",
    "media_player_entity_id": "media_player.office",
    "message": "Testing voice output"
}
```

### Notification Services

```python
# Push notification (mobile app)
POST /api/services/notify/mobile_app_zeus_iphone
{
    "message": "Legacy needs your attention",
    "title": "Legacy AI",
    "data": {
        "push": {"sound": "default"},
        "actions": [
            {"action": "RESPOND", "title": "Respond"},
            {"action": "DISMISS", "title": "Dismiss"}
        ]
    }
}

# Persistent notification (HA dashboard)
POST /api/services/persistent_notification/create
{
    "message": "System update available",
    "title": "Legacy Alert",
    "notification_id": "legacy_alert_001"
}
```

### Automation Services

```python
# Trigger automation
POST /api/services/automation/trigger
{"entity_id": "automation.morning_routine"}

# Turn automation on/off
POST /api/services/automation/turn_on
POST /api/services/automation/turn_off
{"entity_id": "automation.motion_lights"}

# Reload automations
POST /api/services/automation/reload
```

### Script Services

```python
# Run a script
POST /api/services/script/turn_on
{"entity_id": "script.goodnight_routine"}

# With variables
POST /api/services/script/goodnight_routine
{
    "variables": {
        "brightness": 50,
        "delay": 30
    }
}
```

### Conversation Services

```python
# Process voice/text command through HA assistant
POST /api/services/conversation/process
{
    "text": "Turn on the living room lights",
    "language": "en"
}

# Response:
{
    "speech": {
        "plain": {"speech": "Turned on the living room lights"}
    },
    "response_type": "action_done"
}
```

### Add-on Control (Hassio)

```python
# Start add-on
POST /api/services/hassio/addon_start
{"addon": "core_ssh"}

# Stop add-on
POST /api/services/hassio/addon_stop
{"addon": "core_ssh"}

# Restart add-on
POST /api/services/hassio/addon_restart
{"addon": "a]_music_assistant"}

# Note: Requires special permissions on the token
```

---

## REST API Reference

### Authentication

All requests require the Authorization header:
```
Authorization: Bearer {LONG_LIVED_ACCESS_TOKEN}
Content-Type: application/json
```

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/` | GET | API status check |
| `/api/config` | GET | HA configuration |
| `/api/states` | GET | All entity states |
| `/api/states/{entity_id}` | GET | Single entity state |
| `/api/states/{entity_id}` | POST | Update entity state |
| `/api/services` | GET | List all services |
| `/api/services/{domain}/{service}` | POST | Call a service |
| `/api/events` | GET | List event types |
| `/api/events/{event_type}` | POST | Fire an event |
| `/api/history/period/{timestamp}` | GET | Historical data |
| `/api/logbook/{timestamp}` | GET | Logbook entries |
| `/api/error_log` | GET | Error log |
| `/api/camera_proxy/{entity_id}` | GET | Camera snapshot |
| `/api/template` | POST | Render a template |

### Example API Calls

```python
import urllib.request
import json

BASE_URL = "https://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa"
TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

def ha_request(endpoint, method="GET", data=None):
    url = f"{BASE_URL}{endpoint}"
    req = urllib.request.Request(url, method=method)
    req.add_header("Authorization", f"Bearer {TOKEN}")
    req.add_header("Content-Type", "application/json")

    if data:
        req.data = json.dumps(data).encode()

    with urllib.request.urlopen(req) as response:
        return json.loads(response.read().decode())

# Get all states
states = ha_request("/api/states")

# Get specific entity
weather = ha_request("/api/states/weather.forecast_home")

# Call a service
ha_request("/api/services/tts/speak", "POST", {
    "entity_id": "tts.home_assistant_cloud",
    "media_player_entity_id": "media_player.office",
    "message": "Hello from Legacy"
})

# Render a template
result = ha_request("/api/template", "POST", {
    "template": "The weather is {{ states('weather.forecast_home') }}"
})
```

---

## WebSocket API

For real-time updates, use the WebSocket API:

```javascript
const ws = new WebSocket('wss://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa/api/websocket');

ws.onopen = () => {
    // Authenticate
    ws.send(JSON.stringify({
        type: 'auth',
        access_token: 'YOUR_TOKEN'
    }));
};

ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);

    if (msg.type === 'auth_ok') {
        // Subscribe to state changes
        ws.send(JSON.stringify({
            id: 1,
            type: 'subscribe_events',
            event_type: 'state_changed'
        }));
    }

    if (msg.type === 'event') {
        console.log('State changed:', msg.event.data);
    }
};
```

### WebSocket Commands

| Command | Description |
|---------|-------------|
| `subscribe_events` | Subscribe to event types |
| `unsubscribe_events` | Unsubscribe |
| `call_service` | Call a service |
| `get_states` | Get all states |
| `get_config` | Get configuration |
| `get_services` | Get available services |
| `get_panels` | Get dashboard panels |

---

## Voice Assistant Integration

### Wyoming Protocol (Local Voice)

Home Assistant uses Wyoming protocol for voice:

```
Audio In → openWakeWord → Whisper (STT) → Conversation → Piper (TTS) → Audio Out
                │              │                │              │
           Wake word      Speech to       Process         Text to
           detection        text          command          speech
```

### ATOM Echo Setup

1. Install ESPHome add-on
2. Create new device, select "M5Stack ATOM Echo"
3. Flash firmware via USB
4. Configure in HA as voice satellite

```yaml
# ESPHome configuration for ATOM Echo
esphome:
  name: atom-echo-office
  friendly_name: Office Voice

esp32:
  board: m5stack-atom
  framework:
    type: esp-idf

microphone:
  - platform: i2s_audio
    id: echo_mic
    i2s_din_pin: GPIO23
    adc_type: external
    pdm: true

speaker:
  - platform: i2s_audio
    id: echo_speaker
    i2s_dout_pin: GPIO21
    dac_type: external
    mode: mono

voice_assistant:
  microphone: echo_mic
  speaker: echo_speaker
  use_wake_word: true

  on_wake_word_detected:
    - light.turn_on:
        id: led
        effect: pulse
```

### Custom Wake Word

```yaml
# Create custom wake word "Legacy"
wake_word:
  - platform: openwakeword
    models:
      - model: legacy  # Custom trained model
        sensitivity: 0.5
```

---

## Automation Capabilities

### Trigger Types

| Trigger | Example | Use Case |
|---------|---------|----------|
| **State** | Device turns on | React to changes |
| **Time** | Every day at 7 AM | Scheduled tasks |
| **Sun** | Sunset/sunrise | Daylight-based |
| **Zone** | Zeus enters home | Location-based |
| **Device** | Button pressed | Physical input |
| **Webhook** | HTTP request | External integration |
| **Template** | Complex condition | Advanced logic |
| **Calendar** | Event starts | Schedule-based |
| **Conversation** | Voice command | Voice control |

### Action Types

| Action | Example | Use Case |
|--------|---------|----------|
| **Service** | Turn on light | Control devices |
| **Delay** | Wait 5 minutes | Timing |
| **Wait** | Wait for door close | Conditional |
| **Choose** | If/then/else | Logic |
| **Repeat** | Loop 5 times | Iteration |
| **Parallel** | Run simultaneously | Efficiency |
| **Script** | Run script | Reusability |

### Example Automations for Legacy

```yaml
# Morning briefing
automation:
  - alias: "Legacy Morning Briefing"
    trigger:
      - platform: time
        at: "07:00:00"
    condition:
      - condition: state
        entity_id: person.zeuslegacy
        state: "home"
    action:
      - service: tts.cloud_say
        data:
          entity_id: media_player.bedroom
          message: >
            Good morning Mr. Delacruz.
            It's {{ now().strftime('%A, %B %d') }}.
            The weather is {{ states('weather.forecast_home') }}
            with a high of {{ state_attr('weather.forecast_home', 'temperature') }} degrees.

# Legacy presence detection
automation:
  - alias: "Notify Legacy of Arrival"
    trigger:
      - platform: state
        entity_id: person.zeuslegacy
        to: "home"
    action:
      - service: rest_command.notify_legacy
        data:
          message: "Zeus has arrived home"

# Emergency alert
automation:
  - alias: "Smoke Detected Alert"
    trigger:
      - platform: state
        entity_id: binary_sensor.smoke_detector
        to: "on"
    action:
      - service: notify.mobile_app_zeus
        data:
          message: "EMERGENCY: Smoke detected!"
          data:
            push:
              sound:
                name: default
                critical: 1
                volume: 1.0
```

---

## Security Best Practices

### Token Management

| Token Type | Use Case | Expiry | Permissions |
|------------|----------|--------|-------------|
| **Long-Lived** | API access | 10 years | Full user permissions |
| **Short-Lived** | Temporary access | 30 min | Limited |
| **Refresh** | OAuth flows | Configurable | Varies |

**Current Token Info:**
- Created: 2026-01-06
- Expires: 2036-01-06
- Rotation recommended: Annually

### Security Checklist

- [x] HA Cloud (encrypted tunnel, no port forwarding)
- [ ] Two-factor authentication enabled
- [ ] Regular token rotation
- [ ] IP allowlist for sensitive operations
- [ ] Separate user accounts per integration
- [ ] Audit log monitoring

### Network Security

```yaml
# Recommended firewall rules (if using port forwarding)
# DO NOT port forward if using HA Cloud

# Block external access except HA Cloud tunnel
iptables -A INPUT -p tcp --dport 8123 -s 127.0.0.1 -j ACCEPT
iptables -A INPUT -p tcp --dport 8123 -j DROP
```

---

## Integration Roadmap

### Phase 1: Current (Basic)

- [x] HA Cloud remote access
- [x] REST API connection
- [x] Weather integration
- [x] Update monitoring

### Phase 2: Devices (Next)

- [ ] ATOM Echo voice satellite
- [ ] Yi cameras via RTSP
- [ ] SwitchBot integration
- [ ] Mobile app for Zeus

### Phase 3: Intelligence (Future)

- [ ] LLM Vision for presence detection
- [ ] Frigate NVR for object detection
- [ ] Custom wake word "Legacy"
- [ ] Proactive TTS announcements

### Phase 4: Multi-Location

- [ ] NYC Office HA instance
- [ ] Remote HA integration
- [ ] Cross-location automations
- [ ] Unified Legacy-AI control

---

## Troubleshooting

### Common Issues

| Issue | Diagnosis | Solution |
|-------|-----------|----------|
| API returns 401 | Token invalid/expired | Regenerate token |
| API returns 404 | Wrong endpoint | Check URL spelling |
| Remote UI offline | HA Cloud issue | Check `binary_sensor.remote_ui` |
| Slow API response | Network latency | Use local URL when on same network |
| Service fails | Missing entity | Verify entity_id exists |

### Debug Commands

```python
# Check API is responding
GET /api/

# Check configuration valid
GET /api/config

# Get error log
GET /api/error_log

# Check specific entity
GET /api/states/{entity_id}

# Test service (dry run via automation trace)
POST /api/services/automation/trigger
{"entity_id": "automation.test", "skip_condition": true}
```

### Log Access

```bash
# Via SSH (when add-on installed)
ha core logs
ha supervisor logs
ha host logs

# Via API
GET /api/error_log
```

---

## Quick Reference

### Most Used API Calls

```python
# Get weather
GET /api/states/weather.forecast_home

# Get all entities
GET /api/states

# Speak through HA
POST /api/services/tts/cloud_say
{"entity_id": "media_player.X", "message": "Hello"}

# Send notification
POST /api/services/notify/mobile_app_zeus
{"message": "Alert", "title": "Legacy"}

# Turn device on/off
POST /api/services/homeassistant/turn_on
{"entity_id": "switch.X"}

# Get history (last 24h)
GET /api/history/period?filter_entity_id=sensor.X

# Process voice command
POST /api/services/conversation/process
{"text": "What's the weather?"}
```

### Entity ID Patterns

```
{domain}.{object_id}

Domains:
- light.*, switch.*, cover.*   → Controllable
- sensor.*, binary_sensor.*    → Read-only
- automation.*, script.*       → Executable
- person.*, zone.*             → Location
- weather.*, sun.*             → Environment
- media_player.*               → Audio/video
- camera.*                     → Video feeds
- climate.*                    → HVAC
- lock.*, alarm_control_panel.* → Security
```

---

## Related Documentation

- [LEGACY-PLAN.md](./LEGACY-PLAN.md) - Master Legacy AI plan
- [LEGACY-COMMS.md](./LEGACY-COMMS.md) - Communications system
- [LEGACY-LOCATIONS.md](./LEGACY-LOCATIONS.md) - Multi-location architecture
- [Home Assistant REST API](https://developers.home-assistant.io/docs/api/rest/)
- [Home Assistant WebSocket API](https://developers.home-assistant.io/docs/api/websocket/)
