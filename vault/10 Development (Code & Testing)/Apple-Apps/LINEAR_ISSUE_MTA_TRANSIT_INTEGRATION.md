# MTA & Metro North Real-Time Transit Integration

**Priority:** High (P1)
**Estimate:** 2-3 weeks
**Labels:** integration, transit, real-time, flutter, backend
**Epic:** Personal Intelligence

---

## Overview

Integrate NYC MTA and Metro North real-time train data into OneMind OS for commute tracking, arrival notifications, and intelligent travel assistance. Legacy will know your train schedule, alert you to delays, and suggest optimal departure times.

---

## Why This Matters

- **Commute Intelligence**: Legacy knows when to leave for the train
- **Delay Alerts**: Get notified instantly when your train is delayed
- **Route Optimization**: Suggest alternate routes during service disruptions
- **Calendar Integration**: Factor transit time into meeting scheduling
- **Presence Awareness**: Legacy knows you're "in transit" between locations

---

## MTA API Details

### Registration (FREE)

| Item | Details |
|------|---------|
| **Portal** | https://api.mta.info |
| **Cost** | FREE (no paid tiers) |
| **API Key** | Required (instant approval) |
| **Rate Limits** | Unofficial ~60 req/min (be reasonable) |
| **Data Format** | GTFS-RT (Protocol Buffers) |
| **Documentation** | https://api.mta.info/#/subwayRealTimeFeeds |

### Available Feeds

| Feed | Endpoint | Coverage |
|------|----------|----------|
| **Subway Real-Time** | `/nyct/gtfs-ace`, `/nyct/gtfs-bdfm`, etc. | All subway lines |
| **Metro-North** | `/mnr/gtfs-mnr` | All MNR lines |
| **LIRR** | `/lirr/gtfs-lirr` | All LIRR lines |
| **Bus Time** | `/bus-time/v2/siri/...` | NYC buses |
| **Service Alerts** | `/alerts` | System-wide alerts |

### GTFS-RT Data Types

```protobuf
// Trip Update - Arrival/departure predictions
message TripUpdate {
  TripDescriptor trip = 1;
  VehicleDescriptor vehicle = 2;
  repeated StopTimeUpdate stop_time_update = 3;
  uint64 timestamp = 4;
}

// Vehicle Position - Real-time location
message VehiclePosition {
  TripDescriptor trip = 1;
  Position position = 2;
  VehicleStopStatus current_status = 3;
  string stop_id = 4;
}

// Service Alert - Delays, changes, emergencies
message Alert {
  repeated TimeRange active_period = 1;
  repeated EntitySelector informed_entity = 2;
  Cause cause = 3;
  Effect effect = 4;
  TranslatedString header_text = 5;
  TranslatedString description_text = 6;
}
```

---

## Technical Requirements

### Dependencies

```txt
# Add to core/requirements.txt
gtfs-realtime-bindings>=1.0.0  # Google's official GTFS-RT protobuf parser
protobuf>=4.0.0                 # Protocol buffer support
httpx>=0.24.0                   # Async HTTP client (already have)
```

### Environment Variables

```bash
# Add to .env and docker-compose
MTA_API_KEY=your_api_key_here
MTA_FEED_URL=https://api-endpoint.mta.info  # Base URL

# Optional - User preferences
MTA_HOME_STATION=125            # Grand Central = 125
MTA_WORK_STATION=GCT            # Station codes vary by line
MTA_PREFERRED_LINES=["Metro-North", "4", "5", "6"]
```

### Database Schema

```sql
-- Create transit schema
CREATE SCHEMA IF NOT EXISTS transit;

-- User's saved stations/routes
CREATE TABLE transit.user_stations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) NOT NULL,
    station_id VARCHAR(20) NOT NULL,
    station_name VARCHAR(100) NOT NULL,
    line VARCHAR(50) NOT NULL,           -- "Metro-North", "4/5/6", "A/C/E"
    type VARCHAR(20) NOT NULL,           -- "home", "work", "frequent", "saved"
    lat DECIMAL(10, 7),
    lon DECIMAL(10, 7),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, station_id, type)
);

-- Real-time arrival cache (hot data)
CREATE TABLE transit.arrivals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    station_id VARCHAR(20) NOT NULL,
    line VARCHAR(50) NOT NULL,
    direction VARCHAR(20) NOT NULL,      -- "Northbound", "Southbound"
    destination VARCHAR(100),
    arrival_time TIMESTAMPTZ NOT NULL,
    delay_seconds INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'on_time', -- on_time, delayed, cancelled
    vehicle_id VARCHAR(50),
    fetched_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,     -- TTL for cleanup
    INDEX idx_arrivals_station (station_id),
    INDEX idx_arrivals_expires (expires_at)
);

-- Service alerts (active alerts only)
CREATE TABLE transit.alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id VARCHAR(100) UNIQUE NOT NULL,
    affected_lines TEXT[] NOT NULL,      -- ["Metro-North", "4", "5"]
    cause VARCHAR(50),                   -- CONSTRUCTION, WEATHER, etc.
    effect VARCHAR(50),                  -- SIGNIFICANT_DELAYS, DETOUR, etc.
    header TEXT NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_alerts_lines USING GIN (affected_lines)
);

-- Trip history (for analytics)
CREATE TABLE transit.trip_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) NOT NULL,
    origin_station VARCHAR(20) NOT NULL,
    destination_station VARCHAR(20) NOT NULL,
    line VARCHAR(50) NOT NULL,
    scheduled_departure TIMESTAMPTZ,
    actual_departure TIMESTAMPTZ,
    scheduled_arrival TIMESTAMPTZ,
    actual_arrival TIMESTAMPTZ,
    delay_minutes INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_trip_history_user (user_id, created_at DESC)
);
```

---

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1)

#### Task 1.1: MTA API Registration & Testing
- [ ] Register at https://api.mta.info
- [ ] Obtain API key (instant)
- [ ] Test feed access with curl
- [ ] Document all available endpoints
- [ ] Verify Metro-North feed access

```bash
# Test commands
curl -H "x-api-key: $MTA_API_KEY" "https://api-endpoint.mta.info/mnr/gtfs-mnr" --output mnr.pb
# Decode with protoc or Python
```

#### Task 1.2: Create Transit Module Structure

```
core/transit/
├── __init__.py
├── types.py           # Pydantic models for stations, arrivals, alerts
├── mta_client.py      # MTA API client (GTFS-RT parsing)
├── service.py         # Business logic (fetch, cache, notify)
├── api.py             # REST + WebSocket endpoints
├── scheduler.py       # Background polling (every 30s)
└── stations.py        # Station metadata (names, coords, lines)
```

#### Task 1.3: Implement MTA Client

```python
# core/transit/mta_client.py

import httpx
from google.transit import gtfs_realtime_pb2
from typing import List, Optional
from datetime import datetime
from .types import Arrival, Alert, VehiclePosition

class MTAClient:
    """MTA GTFS-RT API Client"""

    def __init__(self, api_key: str, base_url: str = "https://api-endpoint.mta.info"):
        self.api_key = api_key
        self.base_url = base_url
        self.client = httpx.AsyncClient(
            headers={"x-api-key": api_key},
            timeout=10.0
        )

    async def get_metro_north_arrivals(self, station_id: str) -> List[Arrival]:
        """Get real-time arrivals for Metro-North station"""
        response = await self.client.get(f"{self.base_url}/mnr/gtfs-mnr")
        feed = gtfs_realtime_pb2.FeedMessage()
        feed.ParseFromString(response.content)

        arrivals = []
        for entity in feed.entity:
            if entity.HasField("trip_update"):
                for stop_time in entity.trip_update.stop_time_update:
                    if stop_time.stop_id == station_id:
                        arrivals.append(Arrival(
                            station_id=station_id,
                            line=entity.trip_update.trip.route_id,
                            direction=self._get_direction(entity.trip_update.trip),
                            arrival_time=datetime.fromtimestamp(stop_time.arrival.time),
                            delay_seconds=stop_time.arrival.delay or 0,
                            vehicle_id=entity.trip_update.vehicle.id if entity.trip_update.HasField("vehicle") else None
                        ))

        return sorted(arrivals, key=lambda a: a.arrival_time)

    async def get_subway_arrivals(self, station_id: str, feed_id: str) -> List[Arrival]:
        """Get real-time arrivals for subway station"""
        # Feed IDs: ace, bdfm, g, jz, nqrw, l, 123456, 7, sir
        response = await self.client.get(f"{self.base_url}/nyct/gtfs-{feed_id}")
        # ... similar parsing

    async def get_service_alerts(self, lines: Optional[List[str]] = None) -> List[Alert]:
        """Get active service alerts, optionally filtered by line"""
        response = await self.client.get(f"{self.base_url}/alerts")
        feed = gtfs_realtime_pb2.FeedMessage()
        feed.ParseFromString(response.content)

        alerts = []
        for entity in feed.entity:
            if entity.HasField("alert"):
                alert = Alert(
                    alert_id=entity.id,
                    affected_lines=[ie.route_id for ie in entity.alert.informed_entity],
                    cause=entity.alert.cause,
                    effect=entity.alert.effect,
                    header=entity.alert.header_text.translation[0].text if entity.alert.header_text.translation else "",
                    description=entity.alert.description_text.translation[0].text if entity.alert.description_text.translation else ""
                )

                if lines is None or any(l in alert.affected_lines for l in lines):
                    alerts.append(alert)

        return alerts

    async def close(self):
        await self.client.aclose()
```

#### Task 1.4: Define Pydantic Models

```python
# core/transit/types.py

from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime
from enum import Enum

class TransitLine(str, Enum):
    METRO_NORTH = "Metro-North"
    LIRR = "LIRR"
    SUBWAY_ACE = "A/C/E"
    SUBWAY_BDFM = "B/D/F/M"
    SUBWAY_G = "G"
    SUBWAY_JZ = "J/Z"
    SUBWAY_NQRW = "N/Q/R/W"
    SUBWAY_L = "L"
    SUBWAY_123456 = "1/2/3/4/5/6"
    SUBWAY_7 = "7"
    SUBWAY_SIR = "SIR"

class ArrivalStatus(str, Enum):
    ON_TIME = "on_time"
    DELAYED = "delayed"
    CANCELLED = "cancelled"

class AlertEffect(str, Enum):
    NO_SERVICE = "NO_SERVICE"
    REDUCED_SERVICE = "REDUCED_SERVICE"
    SIGNIFICANT_DELAYS = "SIGNIFICANT_DELAYS"
    DETOUR = "DETOUR"
    ADDITIONAL_SERVICE = "ADDITIONAL_SERVICE"
    MODIFIED_SERVICE = "MODIFIED_SERVICE"
    OTHER_EFFECT = "OTHER_EFFECT"
    UNKNOWN_EFFECT = "UNKNOWN_EFFECT"

class Station(BaseModel):
    id: str
    name: str
    line: str
    lat: Optional[float] = None
    lon: Optional[float] = None

class UserStation(BaseModel):
    id: str
    user_id: str
    station_id: str
    station_name: str
    line: str
    type: str  # home, work, frequent, saved
    lat: Optional[float] = None
    lon: Optional[float] = None

class Arrival(BaseModel):
    station_id: str
    line: str
    direction: str
    destination: Optional[str] = None
    arrival_time: datetime
    delay_seconds: int = 0
    status: ArrivalStatus = ArrivalStatus.ON_TIME
    vehicle_id: Optional[str] = None

    @property
    def minutes_until(self) -> int:
        delta = self.arrival_time - datetime.now()
        return max(0, int(delta.total_seconds() / 60))

    @property
    def is_delayed(self) -> bool:
        return self.delay_seconds > 60  # More than 1 minute

class Alert(BaseModel):
    alert_id: str
    affected_lines: List[str]
    cause: Optional[str] = None
    effect: Optional[AlertEffect] = None
    header: str
    description: Optional[str] = None
    start_time: Optional[datetime] = None
    end_time: Optional[datetime] = None

class TransitUpdate(BaseModel):
    """WebSocket message for real-time updates"""
    type: str  # arrival, alert, status
    data: dict
    timestamp: datetime = datetime.now()
```

---

### Phase 2: API & Real-Time Updates (Week 1-2)

#### Task 2.1: REST API Endpoints

```python
# core/transit/api.py

from fastapi import APIRouter, WebSocket, WebSocketDisconnect, Query
from typing import List, Optional
from .service import TransitService
from .types import Arrival, Alert, UserStation

router = APIRouter(prefix="/api/transit", tags=["transit"])

# Singleton service
_service: Optional[TransitService] = None

def get_service() -> TransitService:
    global _service
    if _service is None:
        _service = TransitService()
    return _service

# --- REST Endpoints ---

@router.get("/arrivals/{station_id}", response_model=List[Arrival])
async def get_arrivals(
    station_id: str,
    line: Optional[str] = None,
    limit: int = Query(default=10, le=50)
):
    """Get real-time arrivals for a station"""
    service = get_service()
    arrivals = await service.get_arrivals(station_id, line=line)
    return arrivals[:limit]

@router.get("/alerts", response_model=List[Alert])
async def get_alerts(
    lines: Optional[str] = Query(default=None, description="Comma-separated line IDs")
):
    """Get active service alerts, optionally filtered by lines"""
    service = get_service()
    line_list = lines.split(",") if lines else None
    return await service.get_alerts(line_list)

@router.get("/stations/search", response_model=List[Station])
async def search_stations(
    q: str = Query(..., min_length=2),
    line: Optional[str] = None
):
    """Search for stations by name"""
    service = get_service()
    return await service.search_stations(q, line=line)

@router.get("/user/stations", response_model=List[UserStation])
async def get_user_stations(user_id: str = Query(default="zeus")):
    """Get user's saved stations"""
    service = get_service()
    return await service.get_user_stations(user_id)

@router.post("/user/stations", response_model=UserStation)
async def add_user_station(station: UserStation):
    """Save a station to user's list"""
    service = get_service()
    return await service.add_user_station(station)

@router.delete("/user/stations/{station_id}")
async def remove_user_station(station_id: str, user_id: str = Query(default="zeus")):
    """Remove a station from user's list"""
    service = get_service()
    await service.remove_user_station(user_id, station_id)
    return {"status": "deleted"}

@router.get("/commute/next")
async def get_next_commute(user_id: str = Query(default="zeus")):
    """Get next train for user's commute (home→work or work→home based on time)"""
    service = get_service()
    return await service.get_next_commute(user_id)

@router.get("/commute/status")
async def get_commute_status(user_id: str = Query(default="zeus")):
    """Get overall commute status (delays, alerts affecting route)"""
    service = get_service()
    return await service.get_commute_status(user_id)

# --- WebSocket for Real-Time Updates ---

class TransitWSManager:
    def __init__(self):
        self.connections: dict[str, list[WebSocket]] = {}  # user_id → connections

    async def connect(self, user_id: str, websocket: WebSocket):
        await websocket.accept()
        if user_id not in self.connections:
            self.connections[user_id] = []
        self.connections[user_id].append(websocket)

    async def disconnect(self, user_id: str, websocket: WebSocket):
        if user_id in self.connections:
            self.connections[user_id].remove(websocket)

    async def broadcast(self, user_id: str, message: str):
        if user_id in self.connections:
            for ws in self.connections[user_id]:
                await ws.send_text(message)

    async def broadcast_all(self, message: str):
        for user_connections in self.connections.values():
            for ws in user_connections:
                await ws.send_text(message)

ws_manager = TransitWSManager()

@router.websocket("/ws")
async def transit_websocket(
    websocket: WebSocket,
    user_id: str = Query(default="zeus")
):
    """WebSocket for real-time transit updates"""
    await ws_manager.connect(user_id, websocket)
    try:
        await websocket.send_json({"type": "connected", "user_id": user_id})
        while True:
            # Keep connection alive, updates sent by scheduler
            data = await websocket.receive_text()
            # Handle subscription changes if needed
    except WebSocketDisconnect:
        await ws_manager.disconnect(user_id, websocket)
```

#### Task 2.2: Background Scheduler (Polling)

```python
# core/transit/scheduler.py

import asyncio
import json
from datetime import datetime
from typing import Optional
from .service import TransitService
from .api import ws_manager

class TransitScheduler:
    """Background scheduler for polling MTA feeds"""

    def __init__(self, service: TransitService, poll_interval: int = 30):
        self.service = service
        self.poll_interval = poll_interval
        self.running = False
        self._task: Optional[asyncio.Task] = None

    async def start(self):
        """Start the background polling loop"""
        self.running = True
        self._task = asyncio.create_task(self._poll_loop())

    async def stop(self):
        """Stop the background polling loop"""
        self.running = False
        if self._task:
            self._task.cancel()

    async def _poll_loop(self):
        """Main polling loop"""
        while self.running:
            try:
                await self._poll_and_notify()
            except Exception as e:
                print(f"Transit poll error: {e}")

            await asyncio.sleep(self.poll_interval)

    async def _poll_and_notify(self):
        """Poll MTA feeds and notify subscribers"""

        # 1. Get all users with saved stations
        users = await self.service.get_all_users_with_stations()

        for user_id in users:
            # 2. Get user's stations
            stations = await self.service.get_user_stations(user_id)

            for station in stations:
                # 3. Fetch arrivals
                arrivals = await self.service.get_arrivals(station.station_id)

                # 4. Check for delays > 5 minutes
                for arrival in arrivals[:3]:  # Next 3 trains
                    if arrival.delay_seconds > 300:  # > 5 min delay
                        await ws_manager.broadcast(user_id, json.dumps({
                            "type": "transit.delay",
                            "data": {
                                "station": station.station_name,
                                "line": arrival.line,
                                "delay_minutes": arrival.delay_seconds // 60,
                                "new_arrival": arrival.arrival_time.isoformat()
                            }
                        }))

            # 5. Check for alerts affecting user's lines
            user_lines = list(set(s.line for s in stations))
            alerts = await self.service.get_alerts(user_lines)

            for alert in alerts:
                await ws_manager.broadcast(user_id, json.dumps({
                    "type": "transit.alert",
                    "data": {
                        "alert_id": alert.alert_id,
                        "header": alert.header,
                        "affected_lines": alert.affected_lines,
                        "effect": alert.effect
                    }
                }))
```

#### Task 2.3: NATS Event Publishing

```python
# Add to core/transit/service.py

async def publish_transit_event(self, event_type: str, data: dict):
    """Publish transit event to NATS"""
    subject = f"events.transit.{event_type}"
    await self.nats_client.publish(
        subject,
        json.dumps({
            "type": event_type,
            "data": data,
            "timestamp": datetime.now().isoformat()
        }).encode()
    )

# Event types:
# events.transit.arrival.soon     - Train arriving in < 5 min
# events.transit.delay.detected   - Delay > 5 min detected
# events.transit.alert.new        - New service alert
# events.transit.commute.leave    - Time to leave for commute
```

---

### Phase 3: Integration Registry (Week 2)

#### Task 3.1: Add MTA to Integration Registry

```python
# Add to core/integrations/registry.py

register(Integration(
    id="mta",
    name="MTA Transit",
    category=IntegrationCategory.PRODUCTIVITY,  # Or create TRANSIT category
    description="NYC MTA real-time train data (Subway, Metro-North, LIRR)",
    icon="train",  # Lucide icon
    color="#0039A6",  # MTA blue
    connection=ConnectionConfig(
        type=IntegrationType.DIRECT_API,
        requires_api_key=True,
        api_key_env_var="MTA_API_KEY",
        base_url="https://api-endpoint.mta.info",
        auth_header="x-api-key",
    ),
    webhook=None,  # No webhooks, we poll
    event_types=[
        "arrival.soon",
        "delay.detected",
        "alert.new",
        "commute.leave"
    ],
    actions=[
        "get_arrivals",
        "get_alerts",
        "search_stations",
        "get_commute_status"
    ],
    deep_links={
        "station": DeepLinkPattern(
            mobile="mtasubwaytime://station/{station_id}",
            web="https://new.mta.info/station/{station_id}"
        )
    },
    capabilities=["real_time", "alerts", "trip_planning"],
    docs_url="https://api.mta.info/",
))
```

---

### Phase 4: Flutter UI (Week 2-3)

#### Task 4.1: Create Transit Feature Directory

```
ui/lib/features/transit/
├── screens/
│   ├── transit_dashboard_screen.dart    # Main transit view
│   ├── station_search_screen.dart       # Search/add stations
│   └── commute_settings_screen.dart     # Configure home/work stations
├── widgets/
│   ├── arrival_card.dart                # Single arrival display
│   ├── arrivals_list.dart               # List of arrivals
│   ├── alert_banner.dart                # Service alert banner
│   ├── station_card.dart                # Saved station card
│   ├── commute_status_card.dart         # Overall commute health
│   └── delay_indicator.dart             # Visual delay indicator
├── providers/
│   ├── transit_provider.dart            # Riverpod state
│   └── transit_websocket_provider.dart  # WebSocket connection
└── models/
    ├── arrival.dart                     # Arrival model (Freezed)
    ├── alert.dart                       # Alert model
    └── station.dart                     # Station model
```

#### Task 4.2: Transit Dashboard Screen

```dart
// ui/lib/features/transit/screens/transit_dashboard_screen.dart

class TransitDashboardScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final arrivals = ref.watch(arrivalsProvider);
    final alerts = ref.watch(transitAlertsProvider);
    final commuteStatus = ref.watch(commuteStatusProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text('Transit'),
        actions: [
          IconButton(
            icon: Icon(Icons.add),
            onPressed: () => context.push('/transit/stations/search'),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () => ref.refresh(arrivalsProvider.future),
        child: ListView(
          children: [
            // Commute Status Card
            if (commuteStatus != null)
              CommuteStatusCard(status: commuteStatus),

            // Active Alerts
            if (alerts.isNotEmpty)
              ...alerts.map((a) => AlertBanner(alert: a)),

            // Saved Stations with Arrivals
            SectionHeader('Your Stations'),
            ...arrivals.entries.map((entry) =>
              StationArrivalsSection(
                station: entry.key,
                arrivals: entry.value,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

#### Task 4.3: Arrival Card Widget

```dart
// ui/lib/features/transit/widgets/arrival_card.dart

class ArrivalCard extends StatelessWidget {
  final Arrival arrival;

  @override
  Widget build(BuildContext context) {
    final isDelayed = arrival.delaySeconds > 60;
    final minutes = arrival.minutesUntil;

    return Card(
      color: isDelayed ? Colors.red.shade50 : null,
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: _getLineColor(arrival.line),
          child: Text(
            arrival.line.substring(0, 1),
            style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
          ),
        ),
        title: Text(arrival.destination ?? arrival.direction),
        subtitle: Text(arrival.line),
        trailing: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Text(
              minutes == 0 ? 'NOW' : '$minutes min',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: minutes <= 2 ? Colors.red : null,
              ),
            ),
            if (isDelayed)
              Text(
                '+${arrival.delaySeconds ~/ 60} min delay',
                style: TextStyle(color: Colors.red, fontSize: 12),
              ),
          ],
        ),
      ),
    );
  }

  Color _getLineColor(String line) {
    // MTA line colors
    final colors = {
      'A': Color(0xFF0039A6), 'C': Color(0xFF0039A6), 'E': Color(0xFF0039A6),
      'B': Color(0xFFFF6319), 'D': Color(0xFFFF6319), 'F': Color(0xFFFF6319), 'M': Color(0xFFFF6319),
      'G': Color(0xFF6CBE45),
      'J': Color(0xFF996633), 'Z': Color(0xFF996633),
      'L': Color(0xFFA7A9AC),
      'N': Color(0xFFFCCC0A), 'Q': Color(0xFFFCCC0A), 'R': Color(0xFFFCCC0A), 'W': Color(0xFFFCCC0A),
      '1': Color(0xFFEE352E), '2': Color(0xFFEE352E), '3': Color(0xFFEE352E),
      '4': Color(0xFF00933C), '5': Color(0xFF00933C), '6': Color(0xFF00933C),
      '7': Color(0xFFB933AD),
      'Metro-North': Color(0xFF0039A6),
      'LIRR': Color(0xFF0039A6),
    };
    return colors[line] ?? Color(0xFF0039A6);
  }
}
```

#### Task 4.4: Add Routes to App Router

```dart
// Add to ui/lib/core/router/app_router.dart

GoRoute(
  path: '/transit',
  builder: (context, state) => TransitDashboardScreen(),
),
GoRoute(
  path: '/transit/stations/search',
  builder: (context, state) => StationSearchScreen(),
),
GoRoute(
  path: '/transit/settings',
  builder: (context, state) => CommuteSettingsScreen(),
),
```

#### Task 4.5: Add to Navigation Drawer

```dart
// Add to COMMAND section in app_drawer.dart

_TacticalSection(
  title: 'COMMAND',
  items: [
    DrawerItem(icon: Icons.dashboard, label: 'Dashboard', route: '/dashboard'),
    DrawerItem(icon: Icons.chat, label: 'Comms', route: '/'),
    DrawerItem(icon: Icons.inbox, label: 'Inbox', route: '/inbox'),
    DrawerItem(icon: Icons.train, label: 'Transit', route: '/transit'),  // NEW
  ],
),
```

---

### Phase 5: Legacy Agent Integration (Week 3)

#### Task 5.1: Transit Tool for Legacy

```python
# core/agent/tools/transit.py

from agno.tools import tool
from typing import Optional, List
from core.transit.service import TransitService

_service: Optional[TransitService] = None

def get_transit_service() -> TransitService:
    global _service
    if _service is None:
        _service = TransitService()
    return _service

@tool(
    name="get_train_arrivals",
    description="Get real-time train arrivals for a station. Returns next trains with arrival times and any delays."
)
async def get_train_arrivals(
    station_name: str,
    line: Optional[str] = None,
    limit: int = 5
) -> str:
    """Get upcoming train arrivals"""
    service = get_transit_service()

    # Search for station if name given
    stations = await service.search_stations(station_name, line=line)
    if not stations:
        return f"No station found matching '{station_name}'"

    station = stations[0]
    arrivals = await service.get_arrivals(station.id, line=line)

    if not arrivals:
        return f"No upcoming arrivals at {station.name}"

    result = f"Next trains at {station.name}:\n"
    for arr in arrivals[:limit]:
        delay_str = f" (+{arr.delay_seconds // 60} min delay)" if arr.delay_seconds > 60 else ""
        result += f"- {arr.line} to {arr.destination}: {arr.minutes_until} min{delay_str}\n"

    return result

@tool(
    name="get_transit_alerts",
    description="Get active service alerts for MTA subway and Metro-North. Optionally filter by specific lines."
)
async def get_transit_alerts(
    lines: Optional[List[str]] = None
) -> str:
    """Get active service alerts"""
    service = get_transit_service()
    alerts = await service.get_alerts(lines)

    if not alerts:
        return "No active service alerts" + (f" for {', '.join(lines)}" if lines else "")

    result = "Active Service Alerts:\n"
    for alert in alerts:
        result += f"\n**{alert.header}**\n"
        result += f"Affects: {', '.join(alert.affected_lines)}\n"
        if alert.description:
            result += f"{alert.description[:200]}...\n"

    return result

@tool(
    name="get_commute_status",
    description="Get the current commute status including next train, delays, and any alerts affecting the route."
)
async def get_commute_status(
    user_id: str = "zeus"
) -> str:
    """Get overall commute status"""
    service = get_transit_service()
    status = await service.get_commute_status(user_id)

    return f"""Commute Status:
Next train: {status['next_train']['line']} in {status['next_train']['minutes']} min
Route health: {status['health']}
Active alerts: {len(status['alerts'])}
Recommended departure: {status['recommended_departure']}"""

@tool(
    name="check_if_should_leave",
    description="Check if it's time to leave for a train based on walking time to station and next departures."
)
async def check_if_should_leave(
    destination: str,
    walk_minutes: int = 10
) -> str:
    """Check if user should leave now to catch train"""
    service = get_transit_service()
    # Logic to determine optimal departure time
    # ... implementation
```

#### Task 5.2: Add Tools to Legacy Agent

```python
# Add to agents/legacy/agent.py

from core.agent.tools.transit import (
    get_train_arrivals,
    get_transit_alerts,
    get_commute_status,
    check_if_should_leave
)

TOOLS = [
    # ... existing tools
    get_train_arrivals,
    get_transit_alerts,
    get_commute_status,
    check_if_should_leave,
]
```

#### Task 5.3: Commute Workflow

```python
# core/workflows/commute_check.py

from agno.workflows import Workflow, Step
from datetime import datetime, time

class CommuteCheckWorkflow(Workflow):
    """Daily commute check workflow - runs at configured times"""

    name = "commute_check"
    description = "Check commute status and alert if delays"

    # Run at 6:30 AM and 4:30 PM on weekdays
    schedule = "30 6,16 * * 1-5"

    steps = [
        Step(
            name="check_status",
            action="get_commute_status",
            inputs={"user_id": "zeus"}
        ),
        Step(
            name="evaluate_delays",
            action="evaluate",
            condition="status.delay_minutes > 5",
            on_true="notify_delay",
            on_false="end"
        ),
        Step(
            name="notify_delay",
            action="send_notification",
            inputs={
                "title": "Commute Delay",
                "body": "{status.alert_message}",
                "priority": "high"
            }
        )
    ]
```

---

### Phase 6: Notifications & Inbox (Week 3)

#### Task 6.1: Transit Notifications to Inbox

```python
# Add to core/transit/service.py

async def create_delay_notification(self, user_id: str, arrival: Arrival):
    """Create inbox item for significant delay"""
    from core.inbox.service import InboxService

    inbox = InboxService()
    await inbox.create({
        "type": "notification",
        "title": f"Train Delay: {arrival.line}",
        "body": f"Your {arrival.line} train to {arrival.destination} is delayed by {arrival.delay_seconds // 60} minutes. New arrival: {arrival.arrival_time.strftime('%I:%M %p')}",
        "source": "mta",
        "priority": "high" if arrival.delay_seconds > 600 else "normal",
        "user_id": user_id,
        "metadata": {
            "station_id": arrival.station_id,
            "line": arrival.line,
            "delay_seconds": arrival.delay_seconds
        }
    })

async def create_leave_now_notification(self, user_id: str, station: str, train_time: datetime):
    """Create inbox item telling user to leave for train"""
    from core.inbox.service import InboxService

    inbox = InboxService()
    await inbox.create({
        "type": "reminder",
        "title": "Time to Leave!",
        "body": f"Leave now to catch the {train_time.strftime('%I:%M %p')} train from {station}",
        "source": "mta",
        "priority": "high",
        "user_id": user_id,
        "assignee_type": "human"
    })
```

---

## Acceptance Criteria

### Must Have (P0)
- [ ] MTA API key registered and working
- [ ] Metro-North real-time arrivals fetching
- [ ] Subway real-time arrivals fetching
- [ ] Service alerts fetching
- [ ] REST API endpoints (`/api/transit/*`)
- [ ] WebSocket real-time updates
- [ ] User can save home/work stations
- [ ] Background polling (30s interval)
- [ ] Flutter Transit Dashboard screen
- [ ] Arrivals display with delay indicators
- [ ] Alert banners in UI

### Should Have (P1)
- [ ] Legacy agent transit tools (4 tools)
- [ ] Inbox notifications for delays > 5 min
- [ ] "Time to leave" notifications
- [ ] Commute status endpoint
- [ ] NATS event publishing
- [ ] Trip history tracking
- [ ] Integration registry entry

### Nice to Have (P2)
- [ ] Push notifications (FCM)
- [ ] Apple Watch complication
- [ ] Calendar integration (factor transit into meetings)
- [ ] Multi-stop trip planning
- [ ] Historical delay analytics
- [ ] Alternate route suggestions

---

## Test Plan

### Unit Tests
```python
# tests/transit/test_mta_client.py
def test_parse_gtfs_rt_feed()
def test_extract_arrivals()
def test_filter_by_station()
def test_calculate_delay()

# tests/transit/test_service.py
def test_get_user_stations()
def test_get_commute_status()
def test_search_stations()
```

### Integration Tests
```python
# tests/transit/test_api.py
def test_get_arrivals_endpoint()
def test_get_alerts_endpoint()
def test_websocket_connection()
def test_save_user_station()
```

### Manual Tests
- [ ] Verify arrivals match MTA app
- [ ] Verify delays display correctly
- [ ] Test WebSocket reconnection
- [ ] Test notification delivery
- [ ] Test on slow network

---

## Files to Create/Modify

### New Files
| File | Purpose |
|------|---------|
| `core/transit/__init__.py` | Module init |
| `core/transit/types.py` | Pydantic models |
| `core/transit/mta_client.py` | MTA API client |
| `core/transit/service.py` | Business logic |
| `core/transit/api.py` | REST + WebSocket |
| `core/transit/scheduler.py` | Background polling |
| `core/transit/stations.py` | Station metadata |
| `core/agent/tools/transit.py` | Legacy tools |
| `core/workflows/commute_check.py` | Commute workflow |
| `ui/lib/features/transit/screens/*.dart` | 3 screens |
| `ui/lib/features/transit/widgets/*.dart` | 6 widgets |
| `ui/lib/features/transit/providers/*.dart` | 2 providers |
| `ui/lib/features/transit/models/*.dart` | 3 models |
| `tests/transit/*.py` | Test files |

### Modified Files
| File | Changes |
|------|---------|
| `core/requirements.txt` | Add gtfs-realtime-bindings |
| `core/integrations/registry.py` | Add MTA integration |
| `core/orchestrator/main.py` | Mount transit router |
| `agents/legacy/agent.py` | Add transit tools |
| `ui/lib/core/router/app_router.dart` | Add transit routes |
| `ui/lib/shared/widgets/app_drawer.dart` | Add Transit nav item |
| `infra/compose/docker-compose.core.yml` | Add MTA env vars |
| `.env.example` | Add MTA_API_KEY |

---

## Resources

- **MTA Developer Portal**: https://api.mta.info
- **GTFS-RT Spec**: https://gtfs.org/realtime/
- **GTFS-RT Python Bindings**: https://github.com/google/gtfs-realtime-bindings
- **MTA Line Colors**: https://new.mta.info/maps
- **Metro-North Stations**: https://new.mta.info/agency/metro-north-railroad

---

## Estimated Effort

| Phase | Tasks | Hours |
|-------|-------|-------|
| Phase 1: Core Infrastructure | 4 tasks | 8h |
| Phase 2: API & Real-Time | 3 tasks | 6h |
| Phase 3: Integration Registry | 1 task | 1h |
| Phase 4: Flutter UI | 5 tasks | 12h |
| Phase 5: Agent Integration | 3 tasks | 6h |
| Phase 6: Notifications | 1 task | 2h |
| Testing & Polish | - | 5h |
| **Total** | **17 tasks** | **40h** |

---

*Generated with Claude Code*
