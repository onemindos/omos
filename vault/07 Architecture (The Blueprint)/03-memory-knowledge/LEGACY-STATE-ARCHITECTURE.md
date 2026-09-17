# Legacy State Architecture

> How Legacy maintains persistent awareness across restarts

## The Problem

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    WITHOUT STATE PERSISTENCE                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   User: "I'm heading to the coffee shop"                                │
│   Legacy: "Have a nice time!"                                           │
│                                                                          │
│   [Service restarts]                                                     │
│                                                                          │
│   User: "What did I say I was doing?"                                   │
│   Legacy: "I don't know, we just started talking"  ← AMNESIA!           │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## The Solution: Three-Tier State

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    THREE-TIER STATE ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│                         ┌─────────────┐                                  │
│                         │   Legacy    │                                  │
│                         │   Agent     │                                  │
│                         └──────┬──────┘                                  │
│                                │                                         │
│                    ┌───────────┼───────────┐                            │
│                    ▼           ▼           ▼                            │
│             ┌──────────┐ ┌──────────┐ ┌──────────┐                      │
│             │  Redis   │ │ NATS KV  │ │ Postgres │                      │
│             │   HOT    │ │   WARM   │ │   COLD   │                      │
│             └────┬─────┘ └────┬─────┘ └────┬─────┘                      │
│                  │            │            │                             │
│  Access time:  <10ms       <100ms        <1s                            │
│  TTL:          1 hour      1 day         Forever                         │
│  Updates:      Real-time   Per-session   Periodic                        │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## State Categories

### HOT State (Redis) - Real-time Awareness

```
┌─────────────────────────────────────────────────────────────────────────┐
│  HOT STATE - Changes frequently, <10ms access                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  legacy:state:current                                                    │
│  ├── location                                                            │
│  │   ├── latitude: 37.7749                                              │
│  │   ├── longitude: -122.4194                                           │
│  │   ├── place_name: "Home"                                             │
│  │   └── updated_at: "2025-01-19T10:30:00Z"                             │
│  │                                                                       │
│  ├── biometrics                                                          │
│  │   ├── heart_rate: 72                                                 │
│  │   ├── stress_score: 0.3                                              │
│  │   └── steps_today: 4521                                              │
│  │                                                                       │
│  ├── context                                                             │
│  │   ├── activity: "working"                                            │
│  │   ├── current_task: "Writing code"                                   │
│  │   ├── current_app: "VSCode"                                          │
│  │   └── in_meeting: false                                              │
│  │                                                                       │
│  └── devices                                                             │
│      ├── watch_001: { active: true, battery: 85% }                      │
│      └── phone_001: { active: true, battery: 62% }                      │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### WARM State (NATS KV) - Session Context

```
┌─────────────────────────────────────────────────────────────────────────┐
│  WARM STATE - Session context, <100ms access                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  legacy.session.current                                                  │
│  ├── conversation_summary: "User working on OneMind OS..."              │
│  ├── recent_topics: ["biometrics", "state persistence"]                 │
│  ├── pending_tasks: ["Order T-Watch", "Deploy EMQX"]                    │
│  └── mood: "focused"                                                     │
│                                                                          │
│  legacy.events.today                                                     │
│  ├── [0]: { type: "calendar.meeting", time: "14:00" }                   │
│  ├── [1]: { type: "task.due", title: "Review PR" }                      │
│  └── [2]: { type: "reminder", title: "Stand up!" }                      │
│                                                                          │
│  legacy.insights.pending                                                 │
│  ├── [0]: { message: "You've been sedentary for 90min" }                │
│  └── [1]: { message: "Stress indicators elevated" }                     │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### COLD State (PostgreSQL) - Long-term Memory

```
┌─────────────────────────────────────────────────────────────────────────┐
│  COLD STATE - Long-term memory, <1s access                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Table: legacy_state                                                     │
│  ┌────────────┬──────────────────────────────────────────────────────┐  │
│  │ user_id    │ state_data (JSONB)                                   │  │
│  │ updated_at │                                                      │  │
│  └────────────┴──────────────────────────────────────────────────────┘  │
│                                                                          │
│  Table: user_profile                                                     │
│  ┌────────────────────────────────────────────────────────────────────┐ │
│  │ preferences: { voice: "warm", proactivity: "high" }                │ │
│  │ baselines: { resting_hr: 68, typical_hrv: 45 }                     │ │
│  │ known_places: { home: {...}, work: {...}, gym: {...} }             │ │
│  │ routines: { wake_time: "07:00", sleep_time: "23:00" }              │ │
│  └────────────────────────────────────────────────────────────────────┘ │
│                                                                          │
│  Table: conversations (via pgvector)                                     │
│  ┌────────────────────────────────────────────────────────────────────┐ │
│  │ Semantic search over past conversations                            │ │
│  │ "What did we talk about regarding biometrics?"                     │ │
│  └────────────────────────────────────────────────────────────────────┘ │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Data Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         DATA FLOW                                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  SOURCES                    STATE MANAGER              CONSUMERS         │
│  ───────                    ─────────────              ─────────         │
│                                                                          │
│  ┌─────────┐                                                             │
│  │ T-Watch │──biometrics──┐                                             │
│  └─────────┘              │                                              │
│                           │     ┌──────────────┐                         │
│  ┌─────────┐              │     │              │      ┌──────────┐      │
│  │  Phone  │──location────┼────▶│    State     │─────▶│  Legacy  │      │
│  └─────────┘              │     │   Manager    │      │  Agent   │      │
│                           │     │              │      └──────────┘      │
│  ┌─────────┐              │     └──────┬───────┘                         │
│  │ Calendar│──events──────┘            │                                 │
│  └─────────┘                           │                                 │
│                                        ▼                                 │
│  ┌─────────┐                   ┌───────────────┐      ┌──────────┐      │
│  │ Glasses │──context─────────▶│   Storage     │─────▶│   UI     │      │
│  └─────────┘                   │ Redis/NATS/PG │      │ Dashboard│      │
│                                └───────────────┘      └──────────┘      │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Startup Sequence (Hydration)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    STARTUP HYDRATION SEQUENCE                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  1. Legacy starts                                                        │
│     │                                                                    │
│     ▼                                                                    │
│  2. StateManager.start()                                                 │
│     │                                                                    │
│     ▼                                                                    │
│  3. Try Redis (hot state, most recent)                                   │
│     │                                                                    │
│     ├── Found? ──▶ Use it (freshest data)                               │
│     │                                                                    │
│     └── Not found? ──▶ Try NATS KV                                      │
│                        │                                                 │
│                        ├── Found? ──▶ Use it                            │
│                        │                                                 │
│                        └── Not found? ──▶ Try PostgreSQL                │
│                                           │                              │
│                                           ├── Found? ──▶ Use it         │
│                                           │                              │
│                                           └── Not found? ──▶ Fresh      │
│                                                                          │
│  4. State loaded - Legacy knows:                                         │
│     ├── Where you are                                                    │
│     ├── What you were doing                                              │
│     ├── Your biometric state                                             │
│     └── Recent conversation context                                      │
│                                                                          │
│  5. Start periodic save (every 30s)                                      │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Location Sources

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    LOCATION AWARENESS SOURCES                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Source          Accuracy    Update Rate    Battery Impact               │
│  ──────          ────────    ───────────    ──────────────               │
│                                                                          │
│  📱 Phone GPS    ~3m         Continuous     High                         │
│  ⌚ Watch GPS    ~10m        On request     Medium                       │
│  📡 WiFi         ~50m        When connect   Low                          │
│  🌐 IP           ~1km        Per request    None                         │
│  📅 Calendar     ~varies     Per event      None                         │
│  🏠 Geofence     ~100m       On enter/exit  Low                          │
│                                                                          │
│  FUSION STRATEGY:                                                        │
│  1. Use GPS when available and fresh (<5min)                            │
│  2. Fall back to WiFi if GPS stale                                       │
│  3. Use Calendar if no device location                                   │
│  4. Use IP as last resort                                                │
│                                                                          │
│  GEOFENCES (learned places):                                             │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  Home     │  lat: 37.7749  │  lng: -122.4194  │  radius: 100m   │    │
│  │  Work     │  lat: 37.7849  │  lng: -122.4094  │  radius: 200m   │    │
│  │  Gym      │  lat: 37.7649  │  lng: -122.4294  │  radius: 50m    │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Context Injection

Legacy's system prompt gets context injected:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CONTEXT INJECTION INTO PROMPT                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  System Prompt:                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  You are Legacy, Zeus's AI companion.                           │    │
│  │                                                                 │    │
│  │  ## Current Awareness                                           │    │
│  │  📍 Location: Home                                              │    │
│  │  🎯 Activity: Working (45min)                                   │    │
│  │  💻 Task: Writing code                                          │    │
│  │  ⏰ Next: Team standup in 15min                                 │    │
│  │  💓 HR:68 | Stress:30% | Steps:4521                             │    │
│  │  📱 Devices: T-Watch, iPhone                                    │    │
│  │                                                                 │    │
│  │  Use this awareness to be proactive and contextual.             │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  Result: Legacy can say things like:                                     │
│  - "You have standup in 15 minutes, want me to wrap up?"               │
│  - "You've been coding for 45 minutes, your stress is a bit high"      │
│  - "Since you're at home, want me to start your focus playlist?"       │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Integration Points

### 1. Watch → State Manager

```python
# In biometric event handler
async def handle_biometric_event(event):
    state_manager = get_state_manager()
    if state_manager:
        state_manager.update_biometrics(
            heart_rate=event.data.get("hr"),
            hrv=event.data.get("hrv"),
            stress_score=biometric_state.stress_score,
            steps=event.data.get("steps"),
        )
```

### 2. Phone → State Manager

```python
# In location event handler
async def handle_location_event(event):
    state_manager = get_state_manager()
    if state_manager:
        state_manager.update_location(
            latitude=event.data.get("lat"),
            longitude=event.data.get("lng"),
            place_name=event.data.get("place"),
            source=LocationSource.GPS,
        )
```

### 3. Legacy Agent → Context

```python
# In agent context provider
def get_dynamic_context() -> str:
    parts = []

    # State-based context
    state_manager = get_state_manager()
    if state_manager:
        parts.append(state_manager.get_context_summary())

    # Pattern-based context
    pattern_engine = get_pattern_engine()
    if pattern_engine:
        parts.append(pattern_engine.get_biometric_context())

    return "\n".join(parts)
```

## Database Schema

```sql
-- State persistence table
CREATE TABLE legacy_state (
    user_id VARCHAR(64) PRIMARY KEY,
    state_data JSONB NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User profile (long-term)
CREATE TABLE user_profile (
    user_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255),
    preferences JSONB DEFAULT '{}',
    baselines JSONB DEFAULT '{}',
    known_places JSONB DEFAULT '{}',
    routines JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for fast lookups
CREATE INDEX idx_legacy_state_updated ON legacy_state(updated_at DESC);
```

## Agno Framework Integration

Agno has built-in session state management that complements our StateManager.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AGNO + Legacy STATE LAYERS                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                     AGNO SESSION STATE                          │    │
│  │  (Conversation memory, tool results, response history)          │    │
│  │  Storage: RedisDb → PostgreSQL (automatic)                      │    │
│  └───────────────────────────┬─────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                   Legacy STATE MANAGER                          │    │
│  │  (Location, biometrics, devices, user context)                  │    │
│  │  Storage: Redis (hot) → NATS KV (warm) → PostgreSQL (cold)      │    │
│  └───────────────────────────┬─────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                   PATTERN ENGINE                                │    │
│  │  (Biometric patterns, stress/fatigue/recovery scores)           │    │
│  │  Storage: In-memory + NATS KV                                   │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Agno RedisDb Configuration

```python
from agno.agent import Agent
from agno.db.redis import RedisDb
from agno.db.postgres import PostgresDb

# Initialize Redis for session storage (hot state)
redis_db = RedisDb(
    db_url="redis://redis:6379/0",
    db_prefix="agno",              # Key prefix
    expire=3600,                    # 1 hour TTL
)

# Initialize PostgreSQL for long-term memory (cold state)
postgres_db = PostgresDb(
    db_url="postgresql://onemind:password@postgres:5432/onemind"
)

# Create agent with dual storage
agent = Agent(
    name="Legacy",
    db=redis_db,           # Hot state (sessions)
    memory_db=postgres_db, # Cold state (long-term memory)
    session_id="user_session",
)
```

### Agno Redis Parameters

| Parameter | Type | Default | Purpose |
|-----------|------|---------|---------|
| `db_url` | str | — | `redis://host:port/db` or `rediss://` for TLS |
| `db_prefix` | str | `"agno"` | Prefix for all Redis keys |
| `expire` | int | — | TTL in seconds (None = no expiry) |
| `session_table` | str | — | Custom table for sessions |
| `memory_table` | str | — | Custom table for memory |
| `metrics_table` | str | — | Custom table for metrics |

### Session State Variables

Agno supports session state that persists across interactions:

```python
from agno.agent import Agent

# Define session state variables
agent = Agent(
    name="Legacy",
    session_state={
        "user_location": None,
        "current_activity": None,
        "stress_level": 0.0,
        "last_check_in": None,
    },
    # State available in instructions via {session_state}
    instructions="""
    You are Legacy.

    ## Current Awareness
    Location: {session_state.user_location}
    Activity: {session_state.current_activity}
    Stress: {session_state.stress_level}
    """,
)

# Update state during run
agent.session_state["user_location"] = "Home"
agent.session_state["stress_level"] = 0.35
```

### Bridging Agno + StateManager

```python
from agno.agent import Agent
from agno.db.redis import RedisDb
from agents.legacy.consciousness import create_state_manager, get_state_manager

async def setup_legacy():
    """Initialize Legacy with unified state management."""

    # 1. Create our StateManager (location, biometrics, devices)
    state_manager = await create_state_manager(
        redis_client=redis_client,
        nats_kv=nats_kv,
        pg_pool=pg_pool,
    )

    # 2. Create Agno agent with Redis session storage
    redis_db = RedisDb(db_url="redis://redis:6379/0")

    agent = Agent(
        name="Legacy",
        db=redis_db,
        session_state={
            "awareness_context": "",  # Populated dynamically
        },
        # Dynamic context injection
        additional_context=get_awareness_context,
    )

    return agent

def get_awareness_context() -> str:
    """Inject real-time awareness into Legacy's context."""
    state_manager = get_state_manager()
    pattern_engine = get_pattern_engine()

    parts = []

    # Location + activity from StateManager
    if state_manager:
        parts.append(state_manager.get_context_summary())

    # Biometric insights from PatternEngine
    if pattern_engine:
        parts.append(pattern_engine.get_biometric_context())

    return "\n".join(parts)
```

## Docker Compose Services

```yaml
# docker-compose/cloud/docker-compose.yml

services:
  redis:
    image: redis:7-alpine
    container_name: onemind-redis
    command: redis-server --appendonly yes --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    ports:
      - "127.0.0.1:6379:6379"
    networks:
      - onemind
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  agno:
    environment:
      DATABASE_URL: postgresql://onemind:${POSTGRES_PASSWORD}@postgres:5432/onemind
      REDIS_URL: redis://redis:6379/0
      NATS_URL: nats://nats:4222
    depends_on:
      redis:
        condition: service_healthy
```

## Complete Integration Example

```python
"""
agents/legacy/main.py - Legacy Agent with persistent state
"""
import os
import asyncio
from agno.agent import Agent
from agno.db.redis import RedisDb
from agno.db.postgres import PostgresDb

from agents.legacy.consciousness import (
    create_state_manager,
    create_pattern_engine,
    get_state_manager,
    get_pattern_engine,
)

async def create_legacy_agent():
    """Create Legacy with full state persistence."""

    # Initialize storage backends
    redis_db = RedisDb(
        db_url=os.environ.get("REDIS_URL", "redis://localhost:6379/0"),
        db_prefix="legacy",
        expire=3600,  # 1 hour session TTL
    )

    postgres_db = PostgresDb(
        db_url=os.environ.get("DATABASE_URL")
    )

    # Initialize our custom state management
    state_manager = await create_state_manager(
        redis_client=redis_db.redis_client,
        pg_pool=postgres_db.pool,
    )

    # Create the agent
    agent = Agent(
        name="Legacy",
        db=redis_db,
        memory_db=postgres_db,
        instructions=Legacy_INSTRUCTIONS,
        # Dynamic context provider
        additional_context=lambda: state_manager.get_context_summary(),
    )

    return agent


Legacy_INSTRUCTIONS = """
You are Legacy, Zeus's AI companion and personal assistant.

## Current Awareness
{additional_context}

## Guidelines
- Be proactive about the user's wellbeing based on biometric data
- Reference location and activity context naturally
- Remember conversation history and user preferences
- Alert on important patterns (stress, fatigue, sedentary time)
"""
```

## AG-UI Protocol Integration

For real-time bi-directional state synchronization with the UI, we use the AG-UI protocol.

See: [AG-UI-INTEGRATION.md](./AG-UI-INTEGRATION.md)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AG-UI STATE SYNC                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  StateManager                          UI (CopilotKit)                   │
│  ────────────                          ───────────────                   │
│                                                                          │
│  state changes ──▶ StateSnapshot ─────▶ Initial render                  │
│                   (full state)                                           │
│                                                                          │
│  state changes ──▶ StateDelta ────────▶ Patch existing                  │
│                   (JSON Patch)          (only what changed)              │
│                                                                          │
│  The StateManager now has AG-UI methods:                                 │
│  ├── to_agui_state()      → Convert to AG-UI format                     │
│  ├── get_agui_snapshot()  → Full StateSnapshot event                    │
│  ├── get_agui_delta()     → StateDelta with JSON Patch                  │
│  └── get_agui_activity()  → ActivitySnapshot event                      │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Quick Usage

```python
from agents.legacy.consciousness import get_state_manager

state_manager = get_state_manager()

# On client connect - send full state
snapshot = state_manager.get_agui_snapshot()
# → {"type": "StateSnapshot", "state": {...}}

# On state change - send delta
delta = state_manager.get_agui_delta(previous_state)
# → {"type": "StateDelta", "delta": [{"op": "replace", "path": "/biometrics/hr", "value": 75}]}
```

## State Streaming API

The following endpoints provide real-time state access and updates:

### GET /api/state
Returns current state snapshot.

```json
{
  "user": {"id": "zeus", "name": "Zeus"},
  "location": {"place": "Home", "city": "NYC", "latitude": 40.7128, ...},
  "biometrics": {"heartRate": 72, "hrv": 45, "stressScore": 0.2, ...},
  "activity": {"type": "working", "task": "coding", ...},
  "devices": {...},
  "awareness": {"mode": "active", "lastUpdate": "..."}
}
```

### GET /api/state/stream
Server-Sent Events (SSE) stream for real-time state updates.

```javascript
// Frontend usage
const eventSource = new EventSource('/api/state/stream');

eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);

    if (data.type === 'StateSnapshot') {
        // Full state on connect
        setState(data.state);
    } else if (data.type === 'StateDelta') {
        // JSON Patch for incremental updates
        applyPatch(state, data.delta);
    } else if (data.type === 'Heartbeat') {
        // Connection keepalive every 30s
    }
};
```

### POST /api/state/location
Update user location (called by phone app, Home Assistant, etc.)

```bash
curl -X POST "http://localhost:8080/api/state/location?\
latitude=40.7128&longitude=-74.0060&place_name=NYC&source=gps"
```

### POST /api/state/biometrics
Update biometric data (called by watch app, wearables).

```bash
curl -X POST "http://localhost:8080/api/state/biometrics?\
heart_rate=72&hrv=45&steps=5000&stress_score=0.2"
```

## Next Steps

1. ✅ **Add Redis to docker-compose** - Done
2. ✅ **Add AG-UI methods to StateManager** - Done
3. ✅ **Initialize StateManager in agent startup** - Done
4. ✅ **Add state streaming API endpoints** - Done
5. **Wire up NATS event handlers to update state**
6. **Inject context into Legacy's prompt**
7. **Integrate with Agno UI frontend**
8. **Add location reporting from phone app**


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
