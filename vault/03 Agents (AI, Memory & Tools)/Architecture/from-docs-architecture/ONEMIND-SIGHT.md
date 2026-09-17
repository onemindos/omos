# OneMind Sight - Platform Bird's Eye View

> **Status:** Implemented and Deployed
> **Last Updated:** 2026-01-25
> **Location:** `ui/lib/features/sight/`

## Overview

**1SIGHT** provides a comprehensive bird's eye view of the entire OneMind OS platform. It focuses on infrastructure visibility that isn't covered by other screens:

- **Platform topology** - Cloud services and edge location connectivity
- **Edge locations** - Remote agent status across Home/Office/Farm/Car
- **Event pipeline** - NATS streams and consciousness processor metrics

## Architecture

```
1SIGHT (Platform Overview)
├── Overview (/sight)          - Platform topology + real-time system gauges
├── Edge (/sight/edge)         - Remote agent locations and heartbeats
└── Pipeline (/sight/pipeline) - NATS streams + consciousness processor stats
```

## Avoiding Redundancy

1SIGHT was designed to complement existing screens without overlap:

| Existing Screen | What It Covers |
|-----------------|----------------|
| `/dashboard` | Live feed, system health services, integrations status, workforce stats |
| `/activity` | Event stream with filtering and unread tracking |
| `/metrics` | Token usage, model costs, tool usage, run counts |
| `/eagle-eye` | Camera/surveillance feeds and sensor alerts |
| `/agents` | Agent list and selection |

**1SIGHT focuses on what's MISSING:**
- Edge location status (RemoteAgents across physical locations)
- NATS event pipeline and consciousness processor metrics
- Platform topology visualization
- Real-time system resource gauges

---

## Screens

### 1. Sight Overview (`/sight`)

**Purpose:** Visual platform topology with real-time system gauges

**Components:**
- **Platform Topology Diagram** - Visual representation of cloud-to-edge architecture
- **System Resource Gauges** - 4 circular gauges (CPU, Memory, Disk, Network)
- **Quick Stats Row** - Edge locations online, NATS msgs/sec, filtered percentage
- **Navigation Links** - Quick access to Edge and Pipeline screens

```
┌─────────────────────────────────────────────────┐
│                 CLOUD (AWS)                      │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐   │
│  │Postgres│ │ Redis  │ │  NATS  │ │ Agno   │   │
│  │   ✅   │ │   ✅   │ │   ✅   │ │   ✅   │   │
│  └────────┘ └────────┘ └────────┘ └────────┘   │
└──────────────────┬──────────────────────────────┘
                   │ NATS Leaf Connections
     ┌─────────────┼─────────────┬─────────────┐
     ▼             ▼             ▼             ▼
┌────────┐   ┌────────┐   ┌────────┐   ┌────────┐
│  HOME  │   │ OFFICE │   │  FARM  │   │  CAR   │
│   ✅   │   │   ✅   │   │   ⚠️   │   │   🔴   │
└────────┘   └────────┘   └────────┘   └────────┘
```

**File:** `ui/lib/features/sight/screens/sight_overview_screen.dart`

---

### 2. Edge Locations (`/sight/edge`)

**Purpose:** Monitor RemoteAgent status across all physical locations

**Components:**
- **Aggregate Stats** - Total/Online/Degraded/Offline counts
- **Location Cards** - Detailed card per location with:
  - Status badge (ONLINE/DEGRADED/OFFLINE)
  - Last heartbeat timestamp
  - Running services (Frigate, Ollama, NATS Leaf)
  - Connected devices
  - Ollama model (if running)
  - Resource gauges (CPU/MEM/DISK)
- **Heartbeat Timeline** - 24-hour visual history per location

**Location Card Example:**
```
┌─────────────────────────────────────┐
│ 🏠 HOME                    ✅ ONLINE │
│─────────────────────────────────────│
│ Last Heartbeat: 2s ago              │
│ Services: Frigate, Ollama, NATS Leaf│
│ Devices: 4 cameras, 2 sensors       │
│ Ollama Model: phi-3-mini            │
│─────────────────────────────────────│
│ CPU: 45%  MEM: 62%  DISK: 38%       │
└─────────────────────────────────────┘
```

**File:** `ui/lib/features/sight/screens/sight_edge_screen.dart`

---

### 3. Event Pipeline (`/sight/pipeline`)

**Purpose:** Deep-dive into NATS streams and consciousness processor

**Components:**
- **Event Flow Diagram** - Visual flow from external events to inbox
- **NATS JetStream Stats** - Messages, consumers, throughput, stream size, lag
- **Consciousness Processor Stats** - Processed/filtered/passed counts, filter breakdown
- **Active Subscriptions** - List of NATS subjects being consumed
- **Latency Stats** - P50/P95/P99 processing latency

**Event Flow:**
```
External Events → NATS → Consciousness Processor → Legacy Inbox
     2,341         ↓              ↓                    141
              Filtered       Time/Location
               (94%)          Filters
```

**Filter Breakdown:**
- Time-based filters: 45% (1,048 events)
- Location filters: 28% (616 events)
- Focus mode: 15% (330 events)
- Noise patterns: 6% (132 events)

**File:** `ui/lib/features/sight/screens/sight_pipeline_screen.dart`

---

## Files Created

### Screens
```
ui/lib/features/sight/screens/
├── sight_overview_screen.dart   (~480 lines)
├── sight_edge_screen.dart       (~520 lines)
└── sight_pipeline_screen.dart   (~600 lines)
```

### Routes Added
**File:** `ui/lib/core/router/app_router.dart`
```dart
// 1SIGHT - Platform Overview & Monitoring
GoRoute(
  path: '/sight',
  name: 'sight',
  builder: (context, state) => const SightOverviewScreen(),
),
GoRoute(
  path: '/sight/edge',
  name: 'sight-edge',
  builder: (context, state) => const SightEdgeScreen(),
),
GoRoute(
  path: '/sight/pipeline',
  name: 'sight-pipeline',
  builder: (context, state) => const SightPipelineScreen(),
),
```

### Navigation Added
**File:** `ui/lib/shared/widgets/app_drawer.dart`
```dart
// 1SIGHT - Platform Overview & Monitoring
_TacticalSection(
  title: '1SIGHT',
  currentRoute: currentRoute,
  items: const [
    _DrawerNavItem('Overview', Icons.hub, '/sight'),
    _DrawerNavItem('Edge', Icons.location_on, '/sight/edge'),
    _DrawerNavItem('Pipeline', Icons.timeline, '/sight/pipeline'),
  ],
),
```

---

## Design System

### Colors
| Purpose | Color | Hex |
|---------|-------|-----|
| Background | Dark | `0xFF0A0A0A` |
| Card | Slightly lighter | `0xFF121212` |
| Primary accent | Red | `0xFFFF0000` |
| Online/Healthy | Green | `0xFF00FF00` |
| Warning | Yellow | `0xFFFFD700` |
| Degraded | Orange | `0xFFFF6B00` |
| Offline/Error | Red | `0xFFFF0000` |
| Text Primary | White | `0xFFFFFFFF` |
| Text Secondary | Light gray | `0xFFB0B0B0` |
| Text Dim | Dark gray | `0xFF6A6A6A` |

### Typography
- **Section headers:** Monospace, 11px, weight 700, letter-spacing 2
- **Stats values:** Monospace, 24px, weight 700
- **Labels:** Monospace, 9-10px, weight 600

### Components
- **Circular gauges** with gradient fill based on value
- **Cards** with subtle colored borders
- **Status badges** with glow effect
- **Timeline bars** for heartbeat history

---

## Backend Endpoints (Future Integration)

Currently using mock data. Future integration with:

| Screen | Endpoints |
|--------|-----------|
| Overview | `GET /health/services`, `GET /api/system/resources` |
| Edge | `GET /api/edge/locations` |
| Pipeline | `GET /consciousness/processor/stats`, `GET /api/nats/stats` |

---

## Data Models (Future)

### EdgeLocation
```dart
class EdgeLocation {
  final String id;
  final String name;
  final IconData icon;
  final String status; // online, offline, degraded
  final DateTime? lastHeartbeat;
  final List<String> services;
  final int deviceCount;
  final String? ollamaModel;
  final double cpuPercent;
  final double memoryPercent;
  final double diskPercent;
}
```

### SystemResources
```dart
class SystemResources {
  final double cpuPercent;
  final double memoryPercent;
  final double diskPercent;
  final double networkMbps;
  final DateTime timestamp;
}
```

### ConsciousnessStats
```dart
class ConsciousnessStats {
  final int eventsTotal;
  final int eventsFiltered;
  final int eventsPassed;
  final double filterPercent;
  final double latencyP50Ms;
  final double latencyP95Ms;
  final double latencyP99Ms;
  final int activeSubscriptions;
}
```

---

## Drawer Placement

```
COMMAND
├── Dashboard
├── Comms
└── Inbox

INTELLIGENCE
├── Awareness
├── Memory
├── Sessions
└── Traces

1INTEL (Personal Analytics)
├── Dashboard
├── Trends
├── Correlations
├── Alerts
└── Log

1SIGHT (Platform Overview)  ← THIS SECTION
├── Overview
├── Edge
└── Pipeline

HARDWARE
├── Eagle Eye
├── Watch
└── Voice
```

---

## Access

| Platform | URL |
|----------|-----|
| Web | http://100.102.21.44:3000 |
| Mobile | Same URL (responsive) |
| Desktop | Same URL in browser |

Navigate via drawer menu → **1SIGHT** section

---

## Future Enhancements

1. **Real-time updates** - WebSocket/SSE for live gauge updates
2. **Backend integration** - Connect to actual endpoints
3. **Providers** - Riverpod providers for state management
4. **Edge actions** - Restart services, ping locations
5. **Pipeline controls** - Adjust filter thresholds
6. **Alerts** - Notifications when edge goes offline
