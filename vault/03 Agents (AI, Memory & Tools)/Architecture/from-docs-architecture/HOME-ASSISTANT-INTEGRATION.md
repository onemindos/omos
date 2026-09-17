# Home Assistant Integration

> Full Home Assistant API integration for OneMind OS
> Status: COMPLETE (Backend + Flutter UI)
> Last Updated: 2026-01-25

## Overview

OneMind OS integrates with Home Assistant to provide:
- **Legacy Control** - Natural language control of all devices
- **Real-Time Events** - WebSocket subscription to state changes
- **Smart Alerts** - AI-aware notifications for motion, doors, temperature
- **Dashboard Summary** - Quick status card for Flutter UI
- **Deep Links** - Tap notifications to open HA Companion app
- **Flutter UI** - Full device control screen with filters, scenes, quick actions

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           OneMind OS                                     │
│                                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐ │
│  │   Flutter    │  │    Legacy    │  │    NATS      │  │ Notification│ │
│  │     UI       │  │    Agent     │  │   Events     │  │   Router    │ │
│  └──────────────┘  └──────────────┘  └──────────────┘  └─────────────┘ │
│         │                 │                 ▲                 ▲         │
│         │                 │                 │                 │         │
│         ▼                 ▼                 │                 │         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    core/homeassistant/                           │   │
│  │                                                                  │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │   │
│  │  │  client  │  │   api    │  │  events  │  │  models  │        │   │
│  │  │ REST+WS  │  │ FastAPI  │  │  Bridge  │  │ Pydantic │        │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                               │                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    core/tools/home_assistant.py                  │   │
│  │                                                                  │   │
│  │  home_assistant() | get_home_status() | control_device()        │   │
│  │  list_devices()   | activate_scene()                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                               │                                         │
└───────────────────────────────│─────────────────────────────────────────┘
                                │
                                ▼
                        ┌──────────────┐
                        │     HAOS     │
                        │  (Your HA)   │
                        │              │
                        │  REST API    │
                        │  WebSocket   │
                        └──────────────┘
```

## File Structure

### Backend Files

```
core/homeassistant/
├── __init__.py           # Module exports (~50 lines)
├── models.py             # Pydantic models (~450 lines)
├── client.py             # REST + WebSocket client (~750 lines)
├── events.py             # NATS event bridge (~350 lines)
└── api.py                # FastAPI endpoints (~400 lines)

core/tools/
└── home_assistant.py     # Legacy's HA tools (~710 lines)
```

### Flutter UI Files

```
ui/lib/features/home_assistant/
├── models/
│   ├── home_assistant.dart           # Freezed models (~148 lines)
│   ├── home_assistant.freezed.dart   # Generated
│   └── home_assistant.g.dart         # Generated JSON
├── providers/
│   └── home_assistant_provider.dart  # Riverpod state (~340 lines)
└── screens/
    └── home_assistant_screen.dart    # Main UI screen (~976 lines)
```

### Modified Files

| File | Changes |
|------|---------|
| `core/orchestrator/main.py` | Router mount + startup init |
| `core/tools/__init__.py` | Added tool exports |
| `ui/lib/core/router/app_router.dart` | Added `/home-assistant` route |
| `ui/lib/shared/widgets/app_drawer.dart` | Added nav item in AGENTS section |
| `ui/lib/core/api/agno_client.dart` | Added generic HTTP methods |

**Total Backend: ~2,700 lines**
**Total Flutter UI: ~1,500 lines**

## Environment Variables

```bash
# Required - Add to your .env on server
HOME_ASSISTANT_URL=http://your-home-assistant:8123
HOME_ASSISTANT_TOKEN=your_long_lived_access_token

# Get token from: HA → Profile → Long-Lived Access Tokens → Create Token
```

## Backend Components

### 1. Models (`core/homeassistant/models.py`)

Pydantic models for all HA entities:

| Model | Purpose |
|-------|---------|
| `HAConnectionStatus` | Connection state, version, errors |
| `HAHomeSummary` | Dashboard summary (lights, locks, temp) |
| `HAEntityState` | Current state of any entity |
| `HAEntity` | Entity from registry with metadata |
| `HAArea` | Room/area definition |
| `HADevice` | Physical device info |
| `HAFloor` | Floor definition |
| `HAService` | Available service definition |
| `HAServiceDomain` | Service domain with all services |
| `HAEvent` | Generic HA event |
| `HAStateChangedEvent` | State change event |
| `EntityDomain` | Enum of all entity domains |

### 2. Client (`core/homeassistant/client.py`)

Full REST + WebSocket client:

```python
from core.homeassistant.client import get_ha_client, create_ha_client

# Create client on startup
client = await create_ha_client(websocket=True)

# Get states
states = await client.get_states()

# Control devices
await client.turn_on("light.living_room")
await client.turn_off("light.living_room")
await client.toggle("switch.fan")
await client.lock("lock.front_door")
await client.unlock("lock.front_door")

# Set values
await client.call_service(
    domain="light",
    service="turn_on",
    entity_id="light.office",
    service_data={"brightness_pct": 50}
)

# Get home summary
summary = await client.get_home_summary()
# Returns: lights_on, locks_locked, temperature, alerts, etc.

# Get entities with filters
entities = await client.get_entities(domain="light", area_id="living_room")

# Sync registries
await client.sync_registries()
```

**WebSocket Features:**
- Auth flow: `auth_required` → `auth` → `auth_ok`
- Subscribe to state changes
- Get entity/device/area/floor registries
- Real-time event callbacks

### 3. Event Bridge (`core/homeassistant/events.py`)

Bridges HA events to NATS with alert detection:

```python
from core.homeassistant.events import create_event_bridge

# Create and start bridge
bridge = await create_event_bridge(
    nats_client=nats,
    notification_router=router,
    ha_client=client
)

# Bridge automatically:
# - Subscribes to HA state_changed events
# - Publishes to NATS subjects
# - Detects alert conditions
# - Sends notifications with cooldowns
```

**Event Filters:**
- Domains: light, switch, lock, cover, binary_sensor, sensor, climate, etc.
- State changes only (ignores attribute-only changes)
- Ignores unavailable transitions

**Alert Detection:**

| Alert Type | Trigger | Priority |
|------------|---------|----------|
| Motion | binary_sensor.motion → on | Normal |
| Door | binary_sensor.door → on | Normal |
| Lock | lock → unlocked | High |
| Temperature High | sensor.temperature ≥ 85°F | High |
| Temperature Low | sensor.temperature ≤ 35°F | High |
| Humidity High | sensor.humidity ≥ 80% | Normal |
| Alarm | alarm_control_panel → triggered | Critical |

**Cooldowns:** 15 minutes per alert type + entity (prevents spam)

### 4. API Endpoints (`core/homeassistant/api.py`)

FastAPI router at `/api/ha`:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/ha/status` | Connection status |
| `GET` | `/api/ha/summary` | Home summary for dashboard |
| `GET` | `/api/ha/entities` | List entities (filterable) |
| `GET` | `/api/ha/entities/{id}` | Get entity details |
| `POST` | `/api/ha/entities/{id}/call` | Call service on entity |
| `GET` | `/api/ha/areas` | List areas |
| `GET` | `/api/ha/areas/{id}/entities` | Entities in area |
| `GET` | `/api/ha/devices` | List devices |
| `GET` | `/api/ha/services` | List available services |
| `POST` | `/api/ha/services/call` | Call any service |
| `GET` | `/api/ha/scenes` | List scenes |
| `POST` | `/api/ha/scenes/{id}/activate` | Activate a scene |
| `GET` | `/api/ha/history/{id}` | Entity history |
| `POST` | `/api/ha/sync` | Force sync registries |
| `GET` | `/api/ha/bridge/stats` | Event bridge statistics |
| `POST` | `/api/ha/quick/lights-off` | Turn off all lights |
| `POST` | `/api/ha/quick/lock-all` | Lock all locks |
| `POST` | `/api/ha/quick/arm-away` | Arm alarm (away) |

### 5. Legacy's Tools (`core/tools/home_assistant.py`)

Five tools for natural language control:

#### Main Tool: `home_assistant`

```python
# Get home status
home_assistant(action="status")
# Returns: lights on, locks, temperature, alerts

# Control devices
home_assistant(action="turn_on", entity_id="light.barn")
home_assistant(action="turn_off", entity_id="light.barn")
home_assistant(action="toggle", entity_id="switch.fan")

# Set values
home_assistant(action="set_brightness", entity_id="light.office", value=50)
home_assistant(action="set_temperature", entity_id="climate.living_room", value=72)

# Locks
home_assistant(action="lock", entity_id="lock.front_door")
home_assistant(action="unlock", entity_id="lock.front_door")

# Scenes
home_assistant(action="scene", entity_id="scene.good_night")

# List devices
home_assistant(action="list", domain="light")
home_assistant(action="list", area="barn")
```

#### Convenience Tools

| Tool | Description | Example |
|------|-------------|---------|
| `get_home_status()` | Get summary of home state | "What's happening at home?" |
| `control_device(entity_id, action)` | Control any device | "Turn off the barn lights" |
| `list_devices(device_type, area)` | List devices | "What lights are in the barn?" |
| `activate_scene(scene_name)` | Activate a scene | "Activate good night mode" |

## Flutter UI Components

### 1. Models (`home_assistant.dart`)

Freezed models matching backend:

```dart
@freezed
sealed class HAConnectionStatus with _$HAConnectionStatus {
  const factory HAConnectionStatus({
    @Default(false) bool connected,
    String? haVersion,
    String? locationName,
    DateTime? lastSync,
    String? lastError,
  }) = _HAConnectionStatus;
}

@freezed
sealed class HAHomeSummary with _$HAHomeSummary {
  const factory HAHomeSummary({
    @Default('Home') String locationName,
    @Default(0) int lightsOn,
    @Default(0) int lightsTotal,
    @Default(0) int locksLocked,
    @Default(0) int locksTotal,
    double? temperature,
    String? temperatureUnit,
    double? humidity,
    @Default([]) List<String> alerts,
    @Default([]) List<String> activeScenes,
    int? personsHome,
  }) = _HAHomeSummary;
}

// Also: HAEntity, HAArea, HAScene, HADevice, HAAlert, HAEntityListResponse
```

### 2. Provider (`home_assistant_provider.dart`)

Riverpod state management:

```dart
// State class
class HomeAssistantState {
  final HAConnectionStatus status;
  final HAHomeSummary? summary;
  final List<HAEntity> entities;
  final List<HAArea> areas;
  final List<HAScene> scenes;
  final List<HAAlert> alerts;
  final bool isLoading;
  final String? error;
  final String? selectedAreaId;
  final String? selectedDomain;

  // Computed properties
  List<HAEntity> get filteredEntities;
  Map<String?, List<HAEntity>> get entitiesByArea;
  Set<String> get availableDomains;
}

// Provider
final homeAssistantProvider = StateNotifierProvider<
  HomeAssistantNotifier,
  HomeAssistantState
>((ref) {
  final apiClient = ref.watch(agnoClientProvider);
  return HomeAssistantNotifier(apiClient);
});

// Notifier methods
class HomeAssistantNotifier {
  Future<void> loadAll();           // Load all data
  Future<void> fetchStatus();       // Connection status
  Future<void> fetchSummary();      // Home summary
  Future<void> fetchEntities();     // All entities
  Future<void> fetchAreas();        // All areas
  Future<void> fetchScenes();       // All scenes

  Future<bool> turnOn(entityId);    // Turn on
  Future<bool> turnOff(entityId);   // Turn off
  Future<bool> toggle(entityId);    // Toggle
  Future<bool> lock(entityId);      // Lock
  Future<bool> unlock(entityId);    // Unlock
  Future<bool> activateScene(id);   // Activate scene

  Future<bool> allLightsOff();      // Quick action
  Future<bool> lockAll();           // Quick action
  Future<bool> setBrightness(id, value);

  void setAreaFilter(areaId);       // Filter by area
  void setDomainFilter(domain);     // Filter by type
  void clearFilters();              // Clear filters

  Future<void> syncRegistries();    // Force sync
}
```

### 3. Screen (`home_assistant_screen.dart`)

Full UI with these components:

| Component | Purpose |
|-----------|---------|
| `_ConnectionStatusCard` | Shows connected/disconnected, HA version |
| `_HomeSummaryCard` | Lights on/off, locks, temperature, alerts |
| `_QuickActionsCard` | "All Lights Off", "Lock All" buttons |
| `_ScenesCard` | Tap-to-activate scene chips |
| `_FiltersCard` | Filter by area and device type dropdowns |
| `_EntitiesSection` | List of entities grouped by area |
| `_EntityCard` | Individual entity with domain-specific controls |
| `_ErrorCard` | Error display with dismiss |

**Entity Controls by Domain:**

| Domain | Control |
|--------|---------|
| `light`, `switch`, `fan` | Toggle switch |
| `lock` | Lock/Unlock button |
| `cover` | Up/Stop/Down buttons |
| `sensor`, `binary_sensor` | State display |

**Features:**
- Pull-to-refresh
- Sync button in app bar
- Deep link to HA Companion app
- Haptic feedback on actions
- Loading states
- Error handling

### 4. Navigation

**Route:** `/home-assistant`

**Drawer Location:** AGENTS section
```
AGENTS
├── Agents
├── Teams
├── Workflows
├── Approvals
├── Guardrails
├── Robotics
├── Drones
└── Home Assistant  ← HERE
```

## NATS Event Subjects

The event bridge publishes to these NATS subjects:

| Subject | Description |
|---------|-------------|
| `events.ha.state_changed.{domain}.{entity}` | State changes |
| `events.ha.alert.motion` | Motion detected |
| `events.ha.alert.door` | Door opened |
| `events.ha.alert.lock` | Lock unlocked |
| `events.ha.alert.temperature_high` | High temperature |
| `events.ha.alert.temperature_low` | Low temperature |
| `events.ha.alert.humidity_high` | High humidity |
| `events.ha.alert.alarm` | Alarm triggered |

**Event Payload:**
```json
{
  "entity_id": "binary_sensor.front_door",
  "domain": "binary_sensor",
  "old_state": "off",
  "new_state": "on",
  "attributes": {},
  "timestamp": "2026-01-25T12:00:00Z"
}
```

## Alert Thresholds

Default thresholds (configurable in `events.py`):

| Alert Type | Threshold |
|------------|-----------|
| Temperature High | 85°F |
| Temperature Low | 35°F |
| Humidity High | 80% |
| Humidity Low | 20% |
| Cooldown Period | 15 minutes |

## Integration with Other Systems

### Notification Router

Alerts integrate with the custom notification system:

```python
await notification_router.notify(
    user_id="zeus",
    title="Home Alert: Motion",
    body="Motion detected in barn",
    priority=NotificationPriority.NORMAL,
    source="homeassistant",
    category="alert",
    action_url="homeassistant://navigate/entity/binary_sensor.barn_motion",
)
```

### Consciousness Hooks

HA events can trigger consciousness hooks:

```python
# Subscribe to HA events in consciousness system
async def on_ha_motion(subject, data, reply):
    # Process motion event
    # Maybe alert if in certain awareness mode
    pass

await nats.subscribe("events.ha.alert.motion", cb=on_ha_motion)
```

### Deep Links

Notifications include deep links to HA Companion:

```dart
// Open HA Companion to specific entity
launchUrl('homeassistant://navigate/entity/{entity_id}');

// Open HA Companion to area/room
launchUrl('homeassistant://navigate/lovelace/{dashboard}');
```

## Usage Examples

### Natural Language Control

```
User: "Turn off all the barn lights and lock up"

Legacy: I'll turn off the barn lights and lock the barn door.

[Calls home_assistant(action="list", area="barn") to find devices]
[Calls home_assistant(action="turn_off", entity_id="light.barn_main")]
[Calls home_assistant(action="turn_off", entity_id="light.barn_work")]
[Calls home_assistant(action="lock", entity_id="lock.barn_door")]

Done. I've turned off 2 lights in the barn and locked the barn door.
```

### Motion Alert Flow

```
1. Motion sensor triggers in HA
2. HA sends state_changed event via WebSocket
3. OneMind's event bridge detects motion
4. Checks awareness mode and time of day
5. Sends notification via router
6. User receives push: "Motion detected: Barn Camera"
7. User taps notification
8. HA Companion opens to barn camera
```

### Flutter UI Flow

```
1. User opens Home Assistant screen
2. Provider calls loadAll() on init
3. Fetches status, summary, entities, areas, scenes
4. Displays connection status card
5. Shows home summary (lights, locks, temp)
6. User taps "All Lights Off"
7. Provider calls allLightsOff()
8. Backend calls HA service
9. Provider refreshes summary and entities
10. UI updates to show all lights off
```

## Deployment

### 1. Set Environment Variables

```bash
# On server, add to .env
HOME_ASSISTANT_URL=http://your-ha:8123
HOME_ASSISTANT_TOKEN=your_token
```

### 2. Restart Backend

```bash
docker restart onemind-agno
```

### 3. Verify Connection

```bash
# Test status endpoint
curl http://100.102.21.44:8081/api/ha/status

# Expected response:
{
  "connected": true,
  "ha_version": "2024.1.0",
  "location_name": "Home",
  "last_sync": "2026-01-25T12:00:00Z"
}

# Test summary endpoint
curl http://100.102.21.44:8081/api/ha/summary

# Expected response:
{
  "location_name": "Home",
  "lights_on": 3,
  "lights_total": 12,
  "locks_locked": 2,
  "locks_total": 2,
  "temperature": 72.0,
  "temperature_unit": "F",
  "alerts": []
}
```

### 4. Access Flutter UI

Navigate to: **http://100.102.21.44:3000** → Drawer → AGENTS → Home Assistant

Or direct: **http://100.102.21.44:3000/#/home-assistant**

## Testing

### Backend Testing

```bash
# Test connection
curl http://localhost:8081/api/ha/status

# Get home summary
curl http://localhost:8081/api/ha/summary

# List lights
curl "http://localhost:8081/api/ha/entities?domain=light"

# List areas
curl http://localhost:8081/api/ha/areas

# Turn on a light
curl -X POST http://localhost:8081/api/ha/entities/light.barn/call \
  -H "Content-Type: application/json" \
  -d '{"action": "turn_on"}'

# Activate scene
curl -X POST http://localhost:8081/api/ha/scenes/good_night/activate

# Quick actions
curl -X POST http://localhost:8081/api/ha/quick/lights-off
curl -X POST http://localhost:8081/api/ha/quick/lock-all

# Force sync
curl -X POST http://localhost:8081/api/ha/sync

# Get bridge stats
curl http://localhost:8081/api/ha/bridge/stats
```

### Flutter UI Testing

1. Navigate to `/home-assistant`
2. Verify connection status shows "Connected"
3. Check home summary displays correctly
4. Test quick actions (All Lights Off, Lock All)
5. Test scene activation
6. Test entity toggles
7. Test area/domain filters
8. Test pull-to-refresh
9. Test sync button
10. Test deep link to HA Companion

## What This DOESN'T Do

We intentionally don't rebuild HA's UI:

- ❌ Full device browser with sliders
- ❌ Room-by-room detail views
- ❌ Automation builder
- ❌ History graphs
- ❌ Scene editor
- ❌ Integration setup
- ❌ Add-on management

**HA Companion app does all this better.** We focus on:

- ✅ AI control via Legacy
- ✅ Smart alerts with awareness
- ✅ Dashboard summary
- ✅ Quick actions
- ✅ Deep links to HA app

## Future Development

### Phase 1: Current (COMPLETE)
- [x] REST API client
- [x] WebSocket client
- [x] Event bridge to NATS
- [x] Alert detection
- [x] Legacy tools
- [x] FastAPI endpoints
- [x] Flutter UI screen
- [x] Riverpod state management

### Phase 2: Enhancements
- [ ] Real-time WebSocket updates in Flutter UI
- [ ] Entity history charts
- [ ] Automation triggers from UI
- [ ] Camera stream integration
- [ ] Voice control via Legacy
- [ ] Widget for mobile home screen

### Phase 3: Advanced
- [ ] Multi-location support (multiple HA instances)
- [ ] Custom alert rules configuration
- [ ] Energy monitoring dashboard
- [ ] Integration with presence detection
- [ ] Automation suggestions from Legacy

## Cost

| Component | Cost |
|-----------|------|
| Integration | $0 |
| HA Companion App | $0 |
| Alerts (via OneMind notifications) | See notification costs |

**Total: $0** (uses existing notification infrastructure)

## Troubleshooting

### Connection Failed

```
Error: "Client not initialized"
```

**Solution:** Check environment variables are set and restart backend.

### WebSocket Disconnects

```
Error: "WebSocket connection closed"
```

**Solution:** The client auto-reconnects. If persistent, check HA is accessible.

### No Entities Found

```
Empty entity list
```

**Solution:** Run `POST /api/ha/sync` to force registry refresh.

### Alerts Not Sending

```
No notifications received
```

**Solution:** Check cooldown (15 min), verify notification router is configured.

### Flutter UI Not Loading

```
Error fetching data
```

**Solution:**
1. Check backend is running: `curl http://100.102.21.44:8081/api/ha/status`
2. Check API URL in Flutter config
3. Check CORS settings if running locally

## References

- [Home Assistant REST API](https://developers.home-assistant.io/docs/api/rest)
- [Home Assistant WebSocket API](https://developers.home-assistant.io/docs/api/websocket)
- [HA Companion App Deep Links](https://companion.home-assistant.io/docs/integrations/url-handler)
- [OneMind Notification System](./NOTIFICATION-SYSTEM.md)
- [OneMind Consciousness Processor](../legacy/CONSCIOUSNESS-PROCESSOR.md)
