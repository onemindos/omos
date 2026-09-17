# Device Hub - Mobile Sensor Platform

> Transform your iPhone into an always-aware sensor hub for OneMind AI

The Device Hub is a comprehensive mobile sensor platform that streams real-time device data to OneMind, enabling context-aware AI assistance. Your agents always know where you are, what you're doing, and can proactively help based on your situation.

---

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
   - [Location Tracking](#1-location-tracking)
   - [Geofencing](#2-geofencing)
   - [Background Tracking](#3-background-tracking)
   - [Step Counting](#4-step-counting)
   - [Battery Monitoring](#5-battery-monitoring)
   - [Network Awareness](#6-network-awareness)
   - [Motion Sensors](#7-motion-sensors)
3. [Agent Integration](#agent-integration)
4. [Architecture](#architecture)
5. [API Reference](#api-reference)
6. [Data Models](#data-models)
7. [Event Types](#event-types)
8. [Privacy & Settings](#privacy--settings)
9. [Platform Configuration](#platform-configuration)
10. [Troubleshooting](#troubleshooting)

---

## Overview

### What is Device Hub?

Device Hub turns your mobile device into a continuous sensor stream for OneMind. Instead of manually telling your AI assistant where you are or what you're doing, it knows automatically.

### Key Capabilities

| Capability | Description | Battery Impact |
|------------|-------------|----------------|
| **GPS Location** | Real-time position tracking | ~3%/hour (balanced) |
| **Geofencing** | Zone enter/exit detection | <1%/hour |
| **Activity Recognition** | Walking, driving, stationary | <1%/hour |
| **Step Counting** | Daily steps with history | Negligible |
| **Battery State** | Level, charging, low power | Negligible |
| **Network State** | WiFi name, cellular type | Negligible |
| **Motion Sensors** | Accelerometer, gyroscope | ~1%/hour when active |

### Why It Matters

**Before Device Hub:**
```
You: "Turn off the lights"
Agent: "Which lights?"
You: "The ones at home"
Agent: "Are you at home?"
You: "No, I just left"
```

**After Device Hub:**
```
You: "Turn off the lights"
Agent: "I see you just left home. Turning off all lights at Main House."
```

---

## Features

### 1. Location Tracking

Real-time GPS positioning with configurable accuracy and update frequency.

#### How It Works

The app uses iOS CoreLocation (or Android Location Services) to track your position. Three tracking modes balance accuracy vs battery:

| Mode | Update Interval | Accuracy | Battery | Use Case |
|------|-----------------|----------|---------|----------|
| **Aggressive** | 5 seconds | High (10m) | ~10%/hour | Navigation, real-time tracking |
| **Balanced** | 30 seconds | Medium (50m) | ~3%/hour | **Default** - daily use |
| **Passive** | Significant changes | Low (500m) | <1%/hour | Background awareness |

#### What You See

**In the App:**
```
┌─────────────────────────────────────┐
│  📍 Current Location                │
│  ─────────────────────────────────  │
│  Latitude:   41.0186° N             │
│  Longitude:  73.6416° W             │
│  Altitude:   42m                    │
│  Accuracy:   ±15m                   │
│  Speed:      0 km/h (stationary)    │
│  Heading:    N/A                    │
│  Updated:    2 seconds ago          │
│                                     │
│  Activity: 🧍 Stationary            │
│  Zone:     🏠 Home                  │
└─────────────────────────────────────┘
```

**In the Backend:**
```json
{
  "latitude": 41.0186,
  "longitude": -73.6416,
  "altitude": 42.0,
  "accuracy": 15.0,
  "speed": 0.0,
  "heading": null,
  "activity": "stationary",
  "timestamp": "2025-02-02T14:30:00Z",
  "device_id": "zeus-iphone"
}
```

#### Location History

All location data is stored for 1 year, enabling:
- Route visualization on maps
- Pattern detection (daily commute, frequent locations)
- Historical queries ("Where was I last Tuesday?")

---

### 2. Geofencing

Create virtual boundaries that trigger events when you enter or exit.

#### How It Works

Geofences are circular regions defined by a center point and radius. The system monitors these zones efficiently using iOS region monitoring (not continuous GPS), making it very battery-efficient.

#### Built-in Zones

| Zone | Trigger | Description |
|------|---------|-------------|
| **Home** | Enter/Exit | Primary residence - triggers "arrived home" / "left home" |
| **Work** | Enter/Exit | Workplace - triggers "arrived at work" / "left work" |

#### Custom Zones

Create unlimited custom zones:

```
Zone: "Gym"
Center: 41.0234, -73.6501
Radius: 100m
Triggers: Enter, Exit, Dwell (15 min)
```

#### What You See

**When You Leave Home:**
```
┌─────────────────────────────────────┐
│  🏠 → 🚗                            │
│  ─────────────────────────────────  │
│  You left Home                      │
│  Time: 8:15 AM                      │
│                                     │
│  [Actions triggered:]               │
│  ✓ Lights turned off                │
│  ✓ Thermostat set to Away           │
│  ✓ Security cameras activated       │
└─────────────────────────────────────┘
```

**When You Arrive at Work:**
```
┌─────────────────────────────────────┐
│  🚗 → 🏢                            │
│  ─────────────────────────────────  │
│  You arrived at Work                │
│  Time: 9:02 AM                      │
│  Commute: 47 minutes                │
│                                     │
│  [Calendar today:]                  │
│  • 10:00 AM - Team standup          │
│  • 2:00 PM - Client call            │
└─────────────────────────────────────┘
```

#### Geofence Events

| Event | Description | NATS Topic |
|-------|-------------|------------|
| `entered` | Crossed into zone | `presence.geofence.entered` |
| `exited` | Left the zone | `presence.geofence.exited` |
| `dwell` | Stayed in zone for X minutes | `presence.geofence.dwell` |

#### Dwell Detection

Dwell triggers fire when you've been in a zone for a specified duration. Useful for:
- "You've been at the grocery store for 30 minutes" → Send shopping list
- "You've been at the gym for 1 hour" → Log workout

---

### 3. Background Tracking

Location updates continue even when the app is closed or the phone is locked.

#### How It Works

**iOS Implementation:**
1. **Significant Location Changes** - iOS wakes app for major position changes (~500m)
2. **Background App Refresh** - Periodic updates when system allows
3. **Region Monitoring** - Geofence triggers wake the app
4. **Silent Push** - Backend can request location update

**Android Implementation:**
1. **Foreground Service** - Persistent notification shows tracking is active
2. **WorkManager** - Periodic background jobs for batch uploads

#### Battery Optimization

Background tracking uses "significant location changes" by default:
- Only triggers when you move ~500m
- Uses cell tower + WiFi triangulation (not GPS)
- Battery impact: <1%/hour

For higher accuracy, the app can request temporary GPS access:
- When entering/exiting geofences
- When activity changes (started driving)
- On-demand from agent request

#### What You See

**iOS Control Center:**
```
📍 OneMind is using your location
   [Always] [While Using App] [Never]
```

**Notification (when in background):**
```
┌─────────────────────────────────────┐
│ 📍 OneMind                          │
│ Background location active          │
│ Last update: 2 min ago • 🏠 Home    │
└─────────────────────────────────────┘
```

#### Background Events

The backend receives location updates even when you're not using the app:

```python
# Example: Location received at 3:42 AM (phone charging, screen off)
{
  "type": "device.location.updated",
  "source": "background",
  "data": {
    "latitude": 41.0186,
    "longitude": -73.6416,
    "accuracy": 65.0,  # Lower accuracy (WiFi-based)
    "trigger": "significant_location_change"
  }
}
```

---

### 4. Step Counting

Daily step tracking with history, goals, and milestones.

#### How It Works

Uses the device's built-in pedometer (M-series coprocessor on iPhone), which:
- Runs continuously with zero battery impact
- Counts steps even when app is not running
- Provides historical data for past 7 days

#### What You See

**Daily Step Widget:**
```
┌─────────────────────────────────────┐
│  👟 Today's Steps                   │
│  ─────────────────────────────────  │
│                                     │
│     ████████████░░░░░  8,432        │
│                        ────         │
│                       /10,000       │
│                                     │
│  Distance: 6.2 km                   │
│  Calories: 342 kcal                 │
│  Active time: 1h 23m                │
│                                     │
│  [📈 View History]                  │
└─────────────────────────────────────┘
```

**Step History:**
```
┌─────────────────────────────────────┐
│  This Week                          │
│  ─────────────────────────────────  │
│  Mon │ ████████████████   12,341    │
│  Tue │ ██████████         8,102     │
│  Wed │ ███████████████    11,892    │
│  Thu │ █████████          7,234     │
│  Fri │ ████████           6,891     │
│  Sat │ ████████████████   12,003    │
│  Sun │ ████████░░░░░░░░   8,432     │
│  ─────────────────────────────────  │
│  Weekly Total: 66,895               │
│  Daily Average: 9,556               │
└─────────────────────────────────────┘
```

#### Step Milestones

Events fire when you hit step milestones:

| Milestone | Event | Agent Response |
|-----------|-------|----------------|
| 1,000 | `device.steps.milestone` | "Good start to the day!" |
| 5,000 | `device.steps.milestone` | "Halfway to your goal" |
| 10,000 | `device.steps.goal` | "Congratulations! Daily goal reached!" |
| 15,000 | `device.steps.milestone` | "Exceptional day!" |

#### Step Data Model

```json
{
  "date": "2025-02-02",
  "steps": 8432,
  "distance_meters": 6200,
  "calories": 342,
  "active_minutes": 83,
  "floors_climbed": 12,
  "goal": 10000,
  "goal_reached": false
}
```

---

### 5. Battery Monitoring

Real-time battery level and charging state tracking.

#### How It Works

Monitors battery events through the system API:
- Level changes (reported in 1% increments)
- Charging state changes (plugged in/unplugged)
- Low power mode activation

#### What You See

**Battery Widget:**
```
┌─────────────────────────────────────┐
│  🔋 Battery                         │
│  ─────────────────────────────────  │
│                                     │
│     ████████████████░░░░  78%       │
│                                     │
│  Status:     Not charging           │
│  Low Power:  Off                    │
│  Health:     98%                    │
│                                     │
│  Estimate:   6h 42m remaining       │
└─────────────────────────────────────┘
```

#### Battery Events

| Event | Trigger | NATS Topic |
|-------|---------|------------|
| Level update | Every 1% change | `device.battery.level` |
| Low battery | Below 20% | `device.battery.low` |
| Critical battery | Below 5% | `device.battery.critical` |
| Charging started | Plugged in | `device.battery.charging` |
| Fully charged | Reached 100% | `device.battery.full` |

#### Agent Responses

**Low Battery:**
```
Agent: "Your battery is at 15%. You're 30 minutes from home.
        Should I enable Low Power Mode and pause non-essential syncs?"
```

**Charging Detection:**
```
Agent: "I see you're charging at work. Perfect time to run
        that large file sync you postponed."
```

---

### 6. Network Awareness

Track WiFi and cellular connectivity for presence detection and smart syncing.

#### How It Works

Monitors network state:
- Connected WiFi network name (SSID) and identifier (BSSID)
- WiFi signal strength
- Cellular connection type (LTE, 5G, etc.)
- Online/offline status

#### What You See

**Network Widget:**
```
┌─────────────────────────────────────┐
│  📶 Network                         │
│  ─────────────────────────────────  │
│                                     │
│  Type:    WiFi                      │
│  Network: HomeNetwork_5G            │
│  Signal:  ████░ Strong              │
│  Speed:   ~200 Mbps                 │
│                                     │
│  Cellular: 5G (backup)              │
└─────────────────────────────────────┘
```

#### WiFi-Based Presence

WiFi BSSID can determine location more precisely than GPS indoors:

| WiFi Network | Location | Confidence |
|--------------|----------|------------|
| `HomeNetwork_5G` | Home - Living Room | 95% |
| `HomeNetwork_2.4G` | Home - Basement | 90% |
| `CorpWiFi-Floor3` | Work - 3rd Floor | 95% |
| `Starbucks_Guest` | Coffee Shop | 80% |

#### Network Events

| Event | Description | NATS Topic |
|-------|-------------|------------|
| WiFi changed | Connected to new network | `device.network.wifi.changed` |
| Went offline | Lost all connectivity | `device.network.disconnected` |
| Came online | Regained connectivity | `device.network.connected` |
| Cellular changed | LTE → 5G, etc. | `device.network.cellular.changed` |

#### Smart Sync Behavior

The Device Hub batches uploads based on network:

| Network Type | Behavior |
|--------------|----------|
| **WiFi** | Immediate upload, full history sync |
| **5G** | Immediate upload, skip heavy data |
| **LTE** | Batch every 5 minutes |
| **Offline** | Cache locally, sync when online |

---

### 7. Motion Sensors

Accelerometer and gyroscope data for advanced motion detection.

#### How It Works

Accesses raw sensor data:
- **Accelerometer**: Measures device acceleration (gravity + movement)
- **Gyroscope**: Measures rotation rate
- **Device Motion**: Combined/processed sensor data

#### Use Cases

| Feature | Sensors Used | Description |
|---------|--------------|-------------|
| **Fall Detection** | Accelerometer | Detect sudden impacts |
| **Shake to Action** | Accelerometer | Trigger actions by shaking |
| **Orientation** | Gyroscope | Portrait/landscape/face-down |
| **Driving Detection** | Both | Detect vehicle movement patterns |

#### Fall Detection

If enabled, monitors for fall patterns:

```
┌─────────────────────────────────────┐
│  ⚠️ Fall Detected                   │
│  ─────────────────────────────────  │
│                                     │
│  A possible fall was detected       │
│  at 2:34 PM                         │
│                                     │
│  Location: 41.0186, -73.6416        │
│            (Home - Kitchen)         │
│                                     │
│  [I'm OK]  [Need Help]  [Dismiss]   │
│                                     │
│  Auto-alerting emergency contact    │
│  in 30 seconds...                   │
└─────────────────────────────────────┘
```

#### Motion Events

| Event | Trigger | NATS Topic |
|-------|---------|------------|
| Fall detected | High-G impact + no movement | `device.motion.fall` |
| Shake detected | Rapid back-and-forth | `device.motion.shake` |
| Orientation changed | Device rotated | `device.motion.orientation` |

---

## Agent Integration

The Device Hub enables powerful context-aware agent interactions.

### Accessing Device Context

Agents can query current device state via the `get_device_context` tool:

```python
# Agent tool call
result = await get_device_context()

# Returns:
{
  "location": {
    "latitude": 41.0186,
    "longitude": -73.6416,
    "accuracy": 15.0,
    "zone": "home",
    "zone_name": "Main House"
  },
  "activity": "stationary",
  "battery": {
    "level": 78,
    "charging": false,
    "low_power_mode": false
  },
  "network": {
    "type": "wifi",
    "wifi_name": "HomeNetwork_5G",
    "connected": true
  },
  "steps_today": 8432,
  "last_updated": "2025-02-02T14:30:00Z"
}
```

### Subscribing to Events

Agents can react to device events via NATS subscriptions:

```python
# Subscribe to all device events
await bus.subscribe("device.*", on_device_event)

# Subscribe to presence events only
await bus.subscribe("presence.*", on_presence_event)

# Subscribe to specific events
await bus.subscribe("presence.home.left", on_left_home)
await bus.subscribe("device.battery.low", on_low_battery)
```

### Example Agent Scenarios

#### Scenario 1: Smart Commute

```
Event: presence.home.left (8:15 AM)

Agent Actions:
1. Check calendar for first meeting
2. Check traffic to work
3. Send notification: "Good morning! Traffic looks light,
   you should arrive by 8:52 AM. Your first meeting is at 10:00 AM."
4. Start commute playlist on CarPlay
```

#### Scenario 2: Proactive Home Control

```
Event: presence.home.arrived (6:30 PM)

Context Check:
- Outside temperature: 35°F
- Time: Evening
- User pattern: Usually home for the evening

Agent Actions:
1. Turn on living room lights (60% brightness)
2. Set thermostat to 72°F
3. Send notification: "Welcome home! I've warmed up the house."
```

#### Scenario 3: Location-Based Reminder

```
Event: presence.geofence.entered (zone: "Grocery Store")

Agent Memory Check:
- Shopping list has 5 items
- Last visited grocery store: 4 days ago

Agent Actions:
1. Send notification with shopping list
2. Ask: "You're at the grocery store. Here's your list.
         Anything else you need?"
```

#### Scenario 4: Battery-Aware Behavior

```
Event: device.battery.low (15%)

Context Check:
- Location: 12 miles from home
- Driving: No (stationary)
- Time: 3:00 PM

Agent Actions:
1. Reduce background sync frequency
2. Disable non-essential notifications
3. Send: "Battery at 15%. I've reduced background activity.
         Nearest charger: Starbucks (0.3 miles)"
```

#### Scenario 5: Activity-Based Assistance

```
Event: device.activity.changed (stationary → driving)

Agent Actions:
1. Switch to driving mode
2. Enable hands-free interaction
3. Queue important messages for audio read-out
4. Pause non-urgent notifications
```

### Natural Language Queries

Users can ask location-aware questions:

| Question | Device Hub Context | Agent Response |
|----------|-------------------|----------------|
| "Where am I?" | GPS coordinates + zone | "You're at the office, 3rd floor" |
| "How far from home?" | Current location + home geofence | "You're 12.3 miles from home, about 25 minutes" |
| "How many steps today?" | Pedometer data | "You've walked 8,432 steps today, 84% of your goal" |
| "What's my battery?" | Battery state | "Your phone is at 78%, not charging" |
| "Am I on WiFi?" | Network state | "Yes, connected to HomeNetwork_5G" |

---

## Architecture

### System Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         FLUTTER APP                              │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    SENSOR SERVICES                        │   │
│  │                                                           │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│  │  │ Location    │  │ Motion      │  │ DeviceState │      │   │
│  │  │ Service     │  │ Service     │  │ Service     │      │   │
│  │  │             │  │             │  │             │      │   │
│  │  │ • GPS       │  │ • Steps     │  │ • Battery   │      │   │
│  │  │ • Geofence  │  │ • Accel     │  │ • Network   │      │   │
│  │  │ • Activity  │  │ • Gyro      │  │ • Screen    │      │   │
│  │  │ • Background│  │ • Fall det. │  │             │      │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘      │   │
│  │         │                │                │              │   │
│  │         └────────────────┼────────────────┘              │   │
│  │                          │                               │   │
│  └──────────────────────────┼───────────────────────────────┘   │
│                             │                                   │
│                             ▼                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                  DEVICE HUB SERVICE                       │  │
│  │                                                           │  │
│  │  • Orchestrates all sensor services                       │  │
│  │  • Batches events for efficient upload                    │  │
│  │  • Manages app lifecycle (foreground/background)          │  │
│  │  • Handles offline caching                                │  │
│  │  • WebSocket/REST streaming to backend                    │  │
│  │                                                           │  │
│  └───────────────────────────┬──────────────────────────────┘  │
│                              │                                  │
└──────────────────────────────┼──────────────────────────────────┘
                               │
                               │ WebSocket / REST API
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                          BACKEND                                  │
│                                                                  │
│  ┌────────────────┐     ┌────────────────┐                      │
│  │                │     │                │                      │
│  │  /api/devices  │────▶│   NATS Bus     │                      │
│  │                │     │                │                      │
│  │  • location    │     │  device.*      │                      │
│  │  • sensors     │     │  presence.*    │                      │
│  │  • state       │     │                │                      │
│  │  • geofences   │     └───────┬────────┘                      │
│  │  • config      │             │                               │
│  │                │             │                               │
│  └───────┬────────┘             │                               │
│          │                      │                               │
│          ▼                      ▼                               │
│  ┌────────────────┐     ┌────────────────┐                      │
│  │                │     │                │                      │
│  │  PostgreSQL    │     │    Agents      │                      │
│  │                │     │                │                      │
│  │  • locations   │     │  • Subscribe   │                      │
│  │  • geofences   │     │  • React       │                      │
│  │  • steps       │     │  • Query       │                      │
│  │  • device_state│     │  • Tools       │                      │
│  │                │     │                │                      │
│  └────────────────┘     └────────────────┘                      │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### Data Flow

1. **Sensor Collection** (Device)
   - Sensor services collect raw data from device APIs
   - Each service emits events via StreamControllers
   - Events include timestamp and metadata

2. **Batching & Upload** (Device)
   - DeviceHubService aggregates events from all services
   - Batches events based on network type and config
   - Uploads via WebSocket (real-time) or REST (batch)

3. **Processing** (Backend)
   - API receives and validates sensor data
   - Stores in PostgreSQL for history
   - Publishes to NATS for real-time distribution

4. **Agent Consumption** (Backend)
   - Agents subscribe to relevant NATS topics
   - Query API for historical data
   - Use tools to access current context

### Service Structure

```
frontend/lib/platform/services/device_hub/
├── location_service.dart        # GPS, geofencing, activity
├── motion_service.dart          # Accelerometer, gyroscope, steps
├── device_state_service.dart    # Battery, network, screen
├── device_hub_service.dart      # Orchestrator
└── models/
    ├── location_models.dart     # Position, Geofence, Activity
    ├── motion_models.dart       # SensorEvent, StepCount
    └── device_state_models.dart # BatteryState, NetworkState
```

### Provider Structure

```dart
// Singleton services
locationServiceProvider        // LocationService instance
motionServiceProvider          // MotionService instance
deviceStateServiceProvider     // DeviceStateService instance
deviceHubServiceProvider       // DeviceHubService orchestrator

// Stream providers (reactive)
currentPositionProvider        // Stream<Position>
currentActivityProvider        // Stream<ActivityType>
geofenceEventsProvider         // Stream<GeofenceEvent>
stepCountProvider              // Stream<int>
batteryStateProvider           // Stream<BatteryState>
networkStateProvider           // Stream<NetworkState>

// State notifiers
deviceHubConfigNotifierProvider    // DeviceHubConfig state
geofenceListNotifierProvider       // List<Geofence> state
```

---

## API Reference

### Location Endpoints

#### Report Location
```http
POST /api/devices/location
Content-Type: application/json

{
  "latitude": 41.0186,
  "longitude": -73.6416,
  "altitude": 42.0,
  "accuracy": 15.0,
  "speed": 0.0,
  "heading": null,
  "activity": "stationary",
  "timestamp": "2025-02-02T14:30:00Z",
  "device_id": "zeus-iphone"
}

Response: 201 Created
{
  "id": "loc_abc123",
  "status": "recorded"
}
```

#### Get Current Location
```http
GET /api/devices/location/current

Response: 200 OK
{
  "latitude": 41.0186,
  "longitude": -73.6416,
  "altitude": 42.0,
  "accuracy": 15.0,
  "activity": "stationary",
  "zone": "home",
  "zone_name": "Main House",
  "timestamp": "2025-02-02T14:30:00Z"
}
```

#### Get Location History
```http
GET /api/devices/location/history?start=2025-02-01T00:00:00Z&end=2025-02-02T23:59:59Z&limit=1000

Response: 200 OK
{
  "locations": [
    {
      "id": "loc_abc123",
      "latitude": 41.0186,
      "longitude": -73.6416,
      "activity": "stationary",
      "timestamp": "2025-02-02T14:30:00Z"
    },
    ...
  ],
  "count": 847,
  "has_more": false
}
```

### Geofence Endpoints

#### List Geofences
```http
GET /api/devices/geofences

Response: 200 OK
{
  "geofences": [
    {
      "id": "home",
      "name": "Home",
      "latitude": 41.0186,
      "longitude": -73.6416,
      "radius": 100,
      "transition_types": ["enter", "exit"],
      "enabled": true
    },
    {
      "id": "work",
      "name": "Office",
      "latitude": 41.0523,
      "longitude": -73.5891,
      "radius": 150,
      "transition_types": ["enter", "exit", "dwell"],
      "dwell_time": 900,
      "enabled": true
    }
  ]
}
```

#### Create Geofence
```http
POST /api/devices/geofences
Content-Type: application/json

{
  "name": "Gym",
  "latitude": 41.0234,
  "longitude": -73.6501,
  "radius": 100,
  "transition_types": ["enter", "exit"],
  "metadata": {
    "category": "fitness"
  }
}

Response: 201 Created
{
  "id": "geofence_xyz789",
  "name": "Gym",
  ...
}
```

#### Report Geofence Event
```http
POST /api/devices/geofences/event
Content-Type: application/json

{
  "geofence_id": "home",
  "transition": "exit",
  "latitude": 41.0192,
  "longitude": -73.6420,
  "timestamp": "2025-02-02T08:15:00Z"
}

Response: 201 Created
{
  "status": "recorded",
  "event_id": "evt_abc123"
}
```

### Sensor Endpoints

#### Report Sensor Batch
```http
POST /api/devices/sensors/batch
Content-Type: application/json

{
  "device_id": "zeus-iphone",
  "readings": [
    {
      "sensor_type": "pedometer",
      "timestamp": "2025-02-02T14:30:00Z",
      "values": {
        "steps": 8432,
        "distance": 6200,
        "calories": 342
      }
    },
    {
      "sensor_type": "accelerometer",
      "timestamp": "2025-02-02T14:30:01Z",
      "values": {
        "x": 0.02,
        "y": 0.01,
        "z": 9.81
      }
    }
  ]
}

Response: 201 Created
{
  "status": "recorded",
  "count": 2
}
```

#### Get Today's Steps
```http
GET /api/devices/sensors/steps/today

Response: 200 OK
{
  "date": "2025-02-02",
  "steps": 8432,
  "distance_meters": 6200,
  "calories": 342,
  "goal": 10000,
  "goal_percentage": 84.32
}
```

### Device State Endpoints

#### Report Device State
```http
POST /api/devices/state
Content-Type: application/json

{
  "battery": {
    "level": 78,
    "is_charging": false,
    "is_low_power_mode": false
  },
  "network": {
    "type": "wifi",
    "wifi_name": "HomeNetwork_5G",
    "wifi_bssid": "AA:BB:CC:DD:EE:FF",
    "is_connected": true
  },
  "screen_on": true,
  "timestamp": "2025-02-02T14:30:00Z"
}

Response: 201 Created
```

#### Get Current Device State
```http
GET /api/devices/state/current

Response: 200 OK
{
  "battery": {
    "level": 78,
    "is_charging": false,
    "is_low_power_mode": false
  },
  "network": {
    "type": "wifi",
    "wifi_name": "HomeNetwork_5G",
    "is_connected": true
  },
  "screen_on": true,
  "last_updated": "2025-02-02T14:30:00Z"
}
```

### Configuration Endpoints

#### Get Device Config
```http
GET /api/devices/config

Response: 200 OK
{
  "location_mode": "balanced",
  "track_location": true,
  "track_motion": true,
  "track_battery": true,
  "track_network": true,
  "track_environment": false,
  "batch_interval_seconds": 30,
  "upload_on_wifi_only": true,
  "location_history_retention_days": 365
}
```

#### Update Device Config
```http
PATCH /api/devices/config
Content-Type: application/json

{
  "location_mode": "aggressive",
  "batch_interval_seconds": 10
}

Response: 200 OK
{
  "status": "updated",
  "config": { ... }
}
```

### WebSocket Endpoint

#### Real-time Streaming
```http
WS /api/devices/ws

# Client sends:
{
  "type": "location",
  "data": {
    "latitude": 41.0186,
    "longitude": -73.6416,
    ...
  }
}

# Server responds:
{
  "type": "ack",
  "id": "loc_abc123"
}

# Server pushes:
{
  "type": "config_update",
  "data": {
    "location_mode": "passive"
  }
}
```

---

## Data Models

### Location Models

```dart
/// GPS position with metadata
@freezed
class Position with _$Position {
  const factory Position({
    required double latitude,
    required double longitude,
    double? altitude,
    double? accuracy,
    double? speed,
    double? heading,
    ActivityType? activity,
    required DateTime timestamp,
    String? deviceId,
  }) = _Position;
}

/// Activity recognition types
enum ActivityType {
  stationary,
  walking,
  running,
  cycling,
  driving,
  unknown,
}

/// Geofence zone definition
@freezed
class Geofence with _$Geofence {
  const factory Geofence({
    required String id,
    required String name,
    required double latitude,
    required double longitude,
    required double radius,
    @Default([GeofenceTransition.enter, GeofenceTransition.exit])
    List<GeofenceTransition> transitions,
    int? dwellTime,
    @Default(true) bool enabled,
    @Default({}) Map<String, dynamic> metadata,
  }) = _Geofence;
}

/// Geofence transition types
enum GeofenceTransition {
  enter,
  exit,
  dwell,
}

/// Geofence event
@freezed
class GeofenceEvent with _$GeofenceEvent {
  const factory GeofenceEvent({
    required String geofenceId,
    required String geofenceName,
    required GeofenceTransition transition,
    required double latitude,
    required double longitude,
    required DateTime timestamp,
  }) = _GeofenceEvent;
}
```

### Device State Models

```dart
/// Battery state
@freezed
class BatteryState with _$BatteryState {
  const factory BatteryState({
    required int level,
    required BatteryChargeState chargeState,
    @Default(false) bool isLowPowerMode,
  }) = _BatteryState;
}

enum BatteryChargeState {
  unknown,
  charging,
  discharging,
  full,
  notCharging,
}

/// Network state
@freezed
class NetworkState with _$NetworkState {
  const factory NetworkState({
    required NetworkType type,
    String? wifiName,
    String? wifiBssid,
    int? wifiSignalStrength,
    String? cellularType,
    @Default(true) bool isConnected,
  }) = _NetworkState;
}

enum NetworkType {
  wifi,
  cellular,
  ethernet,
  none,
}
```

### Motion Models

```dart
/// Step count data
@freezed
class StepCount with _$StepCount {
  const factory StepCount({
    required String date,
    required int steps,
    double? distanceMeters,
    double? calories,
    int? activeMinutes,
    int? floorsClimbed,
    @Default(10000) int goal,
  }) = _StepCount;
}

/// Accelerometer event
@freezed
class AccelerometerEvent with _$AccelerometerEvent {
  const factory AccelerometerEvent({
    required double x,
    required double y,
    required double z,
    required DateTime timestamp,
  }) = _AccelerometerEvent;
}

/// Fall detection event
@freezed
class FallEvent with _$FallEvent {
  const factory FallEvent({
    required DateTime timestamp,
    required double latitude,
    required double longitude,
    required double impactMagnitude,
    @Default(false) bool confirmed,
  }) = _FallEvent;
}
```

### Configuration Models

```dart
/// Device Hub configuration
@freezed
class DeviceHubConfig with _$DeviceHubConfig {
  const factory DeviceHubConfig({
    @Default(LocationTrackingMode.balanced) LocationTrackingMode locationMode,
    @Default(true) bool trackLocation,
    @Default(true) bool trackMotion,
    @Default(true) bool trackBattery,
    @Default(true) bool trackNetwork,
    @Default(false) bool trackEnvironment,
    @Default(Duration(seconds: 30)) Duration batchInterval,
    @Default(true) bool uploadOnWifiOnly,
    @Default(365) int historyRetentionDays,
  }) = _DeviceHubConfig;
}

enum LocationTrackingMode {
  aggressive,  // 5s updates
  balanced,    // 30s updates
  passive,     // Significant changes only
  off,
}
```

---

## Event Types

### NATS Event Topics

#### Device Events

| Topic | Description | Payload |
|-------|-------------|---------|
| `device.location.updated` | New location recorded | Position data |
| `device.location.stale` | No location for 5+ min | Last known position |
| `device.activity.changed` | Activity type changed | Old → new activity |
| `device.battery.level` | Battery level changed | Level percentage |
| `device.battery.low` | Battery below 20% | Level, estimate |
| `device.battery.critical` | Battery below 5% | Level, estimate |
| `device.battery.charging` | Started charging | Level, charger type |
| `device.battery.full` | Fully charged | - |
| `device.network.connected` | Came online | Network type |
| `device.network.disconnected` | Went offline | Last network |
| `device.network.wifi.changed` | WiFi changed | Old → new SSID |
| `device.steps.milestone` | Step milestone reached | Steps, milestone |
| `device.steps.goal` | Daily goal reached | Steps, goal |
| `device.motion.fall` | Fall detected | Position, impact |
| `device.motion.shake` | Device shaken | Duration |

#### Presence Events

| Topic | Description | Payload |
|-------|-------------|---------|
| `presence.geofence.entered` | Entered a zone | Geofence, position |
| `presence.geofence.exited` | Left a zone | Geofence, position |
| `presence.geofence.dwell` | Dwelled in zone | Geofence, duration |
| `presence.home.arrived` | Arrived at home | Position, time |
| `presence.home.left` | Left home | Position, time |
| `presence.work.arrived` | Arrived at work | Position, time |
| `presence.work.left` | Left work | Position, time |

### Event Payload Examples

#### Location Updated
```json
{
  "type": "device.location.updated",
  "timestamp": "2025-02-02T14:30:00Z",
  "source": "device_hub",
  "data": {
    "latitude": 41.0186,
    "longitude": -73.6416,
    "accuracy": 15.0,
    "activity": "stationary",
    "zone": "home"
  }
}
```

#### Geofence Entered
```json
{
  "type": "presence.geofence.entered",
  "timestamp": "2025-02-02T18:30:00Z",
  "source": "device_hub",
  "data": {
    "geofence_id": "home",
    "geofence_name": "Home",
    "latitude": 41.0188,
    "longitude": -73.6418,
    "previous_zone": null
  }
}
```

#### Battery Low
```json
{
  "type": "device.battery.low",
  "timestamp": "2025-02-02T16:45:00Z",
  "source": "device_hub",
  "priority": 2,
  "data": {
    "level": 18,
    "estimate_minutes": 42,
    "is_charging": false
  }
}
```

---

## Privacy & Settings

### Privacy Controls

| Setting | Default | Description |
|---------|---------|-------------|
| **Track Location** | On | Enable GPS tracking |
| **Track Motion** | On | Enable step/motion sensors |
| **Track Battery** | On | Share battery state |
| **Track Network** | On | Share WiFi/cellular info |
| **Background Tracking** | On | Track when app closed |
| **Upload on WiFi Only** | On | Batch uploads for cellular |
| **History Retention** | 1 year | How long to keep data |

### Data Storage

| Data Type | Storage Location | Retention |
|-----------|------------------|-----------|
| Location history | PostgreSQL (backend) | 1 year |
| Geofences | PostgreSQL (backend) | Until deleted |
| Daily steps | PostgreSQL (backend) | 1 year |
| Device state | PostgreSQL (backend) | 30 days |
| Pending uploads | SQLite (device) | Until uploaded |

### Deleting Data

```http
# Delete all location history
DELETE /api/devices/location/history

# Delete specific date range
DELETE /api/devices/location/history?start=2025-01-01&end=2025-01-31

# Delete all device data
DELETE /api/devices/data
```

### Export Data

```http
# Export all data as JSON
GET /api/devices/export?format=json

# Export location history as GPX
GET /api/devices/location/export?format=gpx
```

---

## Platform Configuration

### iOS Configuration

#### Info.plist
```xml
<!-- Location permissions -->
<key>NSLocationAlwaysAndWhenInUseUsageDescription</key>
<string>OneMind needs location access to provide context-aware assistance and automation. Your location is used to trigger smart home actions, provide relevant suggestions, and enable location-based reminders.</string>

<key>NSLocationWhenInUseUsageDescription</key>
<string>OneMind uses your location to provide relevant suggestions and show nearby information.</string>

<!-- Background modes -->
<key>UIBackgroundModes</key>
<array>
  <string>location</string>
  <string>fetch</string>
  <string>processing</string>
</array>

<!-- Motion sensors -->
<key>NSMotionUsageDescription</key>
<string>OneMind uses motion data to track your steps and detect falls for safety.</string>
```

#### Capabilities (Xcode)
- [x] Background Modes
  - [x] Location updates
  - [x] Background fetch
  - [x] Background processing
- [x] Push Notifications (for silent push location requests)

### Android Configuration

#### AndroidManifest.xml
```xml
<!-- Location permissions -->
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_BACKGROUND_LOCATION" />

<!-- Foreground service for background tracking -->
<uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_LOCATION" />

<!-- Activity recognition -->
<uses-permission android:name="android.permission.ACTIVITY_RECOGNITION" />

<!-- Wake lock for reliable background operation -->
<uses-permission android:name="android.permission.WAKE_LOCK" />

<!-- Network state -->
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
<uses-permission android:name="android.permission.ACCESS_WIFI_STATE" />

<!-- Foreground service declaration -->
<service
    android:name=".DeviceHubForegroundService"
    android:foregroundServiceType="location"
    android:exported="false" />
```

---

## Troubleshooting

### Location Issues

#### Location not updating
1. Check location permissions in Settings
2. Ensure "Always" permission is granted for background
3. Check if Low Power Mode is blocking location
4. Verify internet connectivity for upload

#### Poor accuracy
1. Move to area with clear sky view
2. Check if WiFi is enabled (helps accuracy)
3. Switch to "Aggressive" mode temporarily
4. Restart the app

#### High battery drain
1. Switch to "Passive" mode
2. Enable "Upload on WiFi only"
3. Reduce geofence count
4. Check for location-using apps in background

### Geofence Issues

#### Not triggering
1. Check geofence is enabled
2. Verify radius is appropriate (minimum 100m recommended)
3. Ensure background location is permitted
4. Wait up to 5 minutes after entering/exiting

#### False triggers
1. Increase geofence radius
2. Add dwell time requirement
3. Check for GPS drift near zone edge

### Step Counting Issues

#### Steps not counting
1. Ensure motion permission is granted
2. Check if phone is being carried (not in bag with other items)
3. Restart the app
4. Check device pedometer is working (Health app)

#### Inaccurate count
1. Calibrate by walking known distance
2. Check step length setting
3. Note: Cycling/driving should not count as steps

### Sync Issues

#### Data not uploading
1. Check internet connectivity
2. Verify backend is reachable
3. Check "Upload on WiFi only" setting
4. Look for pending uploads in app settings

#### Missing historical data
1. Check retention period settings
2. Verify data wasn't deleted
3. Check database storage space

---

## Verification Checklist

Use this checklist to verify the Device Hub is working correctly:

### 1. Location Tracking
- [ ] Open app, see current coordinates displayed
- [ ] Position updates every 30 seconds (balanced mode)
- [ ] Activity shows correctly (stationary/walking/driving)
- [ ] Location appears in backend database

### 2. Geofencing
- [ ] Create "Home" geofence via settings
- [ ] Leave home area → receive "left home" notification
- [ ] Backend receives `presence.home.left` event
- [ ] Return home → receive "arrived home" notification
- [ ] Backend receives `presence.home.arrived` event

### 3. Background Tracking
- [ ] Close app completely
- [ ] Move to new location (500m+)
- [ ] Check backend received location update
- [ ] Re-open app, see correct current location

### 4. Step Counting
- [ ] View today's step count
- [ ] Walk around for 100 steps
- [ ] Verify count increments
- [ ] Check step history shows past days

### 5. Battery Monitoring
- [ ] View current battery level in app
- [ ] Backend receives battery state
- [ ] Unplug charger → backend receives state change
- [ ] At <20% → receive low battery notification

### 6. Network Awareness
- [ ] View current network in app
- [ ] Switch WiFi networks → backend receives update
- [ ] Go offline → verify cached data syncs when online

### 7. Agent Integration
- [ ] Ask agent "Where am I?"
- [ ] Agent responds with current location/zone
- [ ] Ask agent "What's my battery?"
- [ ] Agent responds with battery level
- [ ] Ask agent "How many steps today?"
- [ ] Agent responds with step count

---

## Implementation Status

### Backend Components (COMPLETE)

| Component | File | Status |
|-----------|------|--------|
| **Location Tools** | `backend/agno/tools/location.py` | ✅ Complete |
| **Device API** | `backend/agno/api/devices.py` | ✅ Complete |
| **Context Provider** | `backend/agno/context/device_context.py` | ✅ Complete |
| **Tool Registry** | `backend/agno/tools/registry.py` | ✅ Updated |
| **Legacy Agent** | `backend/agno/agents/presets.py` | ✅ Has location tools |

### Location Tools

Legacy and other agents have access to these tools:

```python
# Get user's current GPS location
get_user_location()
# Returns: {"latitude": 41.0186, "longitude": -73.6416, "zone": "Home", "activity": "stationary"}

# Get user's physical activity
get_user_activity()
# Returns: {"activity": "walking", "speed_kmh": 4.2}

# Check if user is at a specific zone
check_user_at_location("Home")
# Returns: {"is_at_location": true, "current_zone": "Home"}
```

### API Endpoints

```bash
# Report location from mobile app
POST /api/devices/location
{
  "latitude": 41.0186,
  "longitude": -73.6416,
  "activity": "stationary",
  "zone_name": "Home"
}

# Get current location
GET /api/devices/location

# Report device state (battery, network)
POST /api/devices/state

# Get full device context
GET /api/devices/context
```

### Context Injection

Legacy automatically knows your zone and activity via context injection:

```python
from backend.agno.context import get_device_context_string

context = get_device_context_string()
# Returns: "User is at Home, stationary."
```

### Frontend Integration (TODO)

The Flutter app needs to send location updates to the backend:

```dart
// In location_service.dart - add this to _handlePositionStream
Future<void> _sendLocationToBackend(Position position) async {
  final response = await http.post(
    Uri.parse('${Environment.apiUrl}/api/devices/location'),
    headers: {'Content-Type': 'application/json'},
    body: jsonEncode({
      'latitude': position.latitude,
      'longitude': position.longitude,
      'altitude': position.altitude,
      'accuracy': position.accuracy,
      'speed': position.speed,
      'heading': position.heading,
      'activity': _currentActivity?.name,
      'zone_name': _currentZone?.name,
      'device_id': 'zeus-iphone',
    }),
  );
}
```

---

## Changelog

### v1.1.0 (2025-02-03)
- Added backend location tools (`get_user_location`, `get_user_activity`, `check_user_at_location`)
- Added Device Hub API (`/api/devices/location`, `/api/devices/state`, `/api/devices/context`)
- Added context injection for automatic zone/activity awareness
- Legacy agent now has location tools in preset

### v1.0.0 (Initial Release)
- Core location tracking with three modes
- Geofencing with enter/exit/dwell
- Background location on iOS and Android
- Step counting with daily history
- Battery and network monitoring
- Agent integration via tools and events
- 1-year location history retention
- Privacy controls and data export

---

*Device Hub is part of OneMind OS - Your Personal AI Command Center*
