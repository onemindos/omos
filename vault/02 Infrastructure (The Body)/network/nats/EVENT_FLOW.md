# OneMind OS: Universal Event Flow Architecture

## Overview

All events in OneMind OS flow through **NATS JetStream**, which acts as the central nervous system. The **OpenMind Bridge** heartbeat loop processes these events at 1 Hz.

---

## Event Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    EVENT LIFECYCLE                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. EVENT INGESTION                                          │
│     ═══════════════                                          │
│                                                              │
│     External Source → EventGateway → NATS Publish           │
│                                                              │
│     Examples:                                                │
│     • GitHub webhook → events.webhook.github                │
│     • ESP32 sensor  → events.sensor.greenhouse              │
│     • Agent action  → events.agent.researcher               │
│     • Heartbeat     → cortex.heartbeat                      │
│                                                              │
│  ────────────────────────────────────────────────────────── │
│                                                              │
│  2. NATS JETSTREAM (Message Broker)                          │
│     ═══════════════════════════════                          │
│                                                              │
│     ┌──────────────────────────────────────┐                │
│     │  JetStream Streams                   │                │
│     │  ─────────────────                   │                │
│     │  • EVENTS (events.>)                 │                │
│     │  • ACTIONS (action.>)                │                │
│     │  • CORTEX (cortex.>)                 │                │
│     │  • SENSORS (sensor.>)                │                │
│     │                                       │                │
│     │  Retention: 7-30 days                │                │
│     │  Storage: File-based persistence     │                │
│     └──────────────────────────────────────┘                │
│                                                              │
│  ────────────────────────────────────────────────────────── │
│                                                              │
│  3. MULTIPLE CONSUMERS (Subscribers)                         │
│     ════════════════════════════════                         │
│                                                              │
│     A. EventArchiver                                         │
│        ├─ Subscribes: events.>, action.>                    │
│        ├─ Action: Write to TimescaleDB                      │
│        └─ Purpose: Long-term storage (90 days)              │
│                                                              │
│     B. OpenMind Bridge (Heartbeat Loop)                      │
│        ├─ Subscribes: events.>, cortex.>                    │
│        ├─ Action: Process in consciousness loop             │
│        └─ Purpose: Real-time awareness & response           │
│                                                              │
│     C. Specialized Agents                                    │
│        ├─ Subscribes: events.agent.{id}                     │
│        ├─ Action: React to specific events                  │
│        └─ Purpose: Automated workflows                      │
│                                                              │
│     D. WebSocket Manager                                     │
│        ├─ Subscribes: events.>, action.>                    │
│        ├─ Action: Broadcast to frontend                     │
│        └─ Purpose: Real-time dashboard updates              │
│                                                              │
│  ────────────────────────────────────────────────────────── │
│                                                              │
│  4. LONG-TERM STORAGE                                        │
│     ══════════════════                                       │
│                                                              │
│     TimescaleDB (Hypertables)                                │
│     ├─ Compressed after 7 days                              │
│     ├─ Retention: 90 days (auto-delete)                     │
│     ├─ Continuous aggregates for analytics                  │
│     └─ Queries: "What did I do yesterday?"                  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## OpenMind Heartbeat Loop (1 Hz)

The bridge runs continuously at 1 Hz and processes events:

### Pseudocode

```python
# backend/openmind/bridge.py

async def heartbeat_loop():
    """
    OpenMind consciousness loop running at 1 Hz.

    This is the "awareness engine" that processes all events
    and decides what actions to take.
    """

    while True:
        cycle_start = time.time()

        # ═══════════════════════════════════════════════════
        # PHASE 1: SENSE (Gather Recent Events)
        # ═══════════════════════════════════════════════════

        # Pull recent events from NATS (last 1 second)
        events = await sense_digital_inputs()

        # Sources:
        # - events.webhook.*     (GitHub PR, ClickUp task)
        # - events.sensor.*      (ESP32 heartbeat, greenhouse temp)
        # - events.agent.*       (Agent completed task)
        # - events.system.*      (Backend health check)

        # ═══════════════════════════════════════════════════
        # PHASE 2: FUSE (Combine with Context)
        # ═══════════════════════════════════════════════════

        # Merge events with:
        # - Recent memories (from PostgreSQL)
        # - Current state (what we're working on)
        # - User preferences
        # - Time/location context

        context = await fuse_context(
            events=events,
            memories=await get_recent_memories(),
            state=get_current_state(),
            user_context=get_user_context()
        )

        # ═══════════════════════════════════════════════════
        # PHASE 3: CORTEX (LLM Decision Making)
        # ═══════════════════════════════════════════════════

        # Send fused context to LLM (Agno Agent)
        # LLM decides:
        # - What events are important?
        # - What actions should we take?
        # - What should we remember?
        # - Should we notify the user?

        plan = await cortex_agent.run(
            prompt=f"""
            Recent events:
            {context.events}

            Current state:
            {context.state}

            What should we do?
            """,
            tools=[
                send_notification,
                create_task,
                update_calendar,
                control_lights,
                # ... all available tools
            ]
        )

        # ═══════════════════════════════════════════════════
        # PHASE 4: ACTION (Execute Plan)
        # ═══════════════════════════════════════════════════

        # Execute actions from the plan
        for action in plan.actions:
            result = await execute_action(action)

            # Publish action result back to NATS
            await nats.publish(
                subject=f"action.{action.type}",
                data={
                    "action": action.name,
                    "result": result,
                    "timestamp": datetime.utcnow()
                }
            )

            # EventArchiver automatically saves this to TimescaleDB

        # ═══════════════════════════════════════════════════
        # PHASE 5: REMEMBER (Update State)
        # ═══════════════════════════════════════════════════

        # Update memories in PostgreSQL
        if plan.should_remember:
            await save_memory(plan.memory)

        # Update current state
        update_current_state(plan.new_state)

        # Publish heartbeat to NATS (for monitoring)
        await nats.publish(
            subject="cortex.heartbeat",
            data={
                "cycle": cycle_count,
                "events_processed": len(events),
                "actions_taken": len(plan.actions),
                "timestamp": datetime.utcnow()
            }
        )

        # ═══════════════════════════════════════════════════
        # Sleep to maintain 1 Hz
        # ═══════════════════════════════════════════════════

        cycle_duration = time.time() - cycle_start
        sleep_time = max(0, 1.0 - cycle_duration)
        await asyncio.sleep(sleep_time)
```

---

## Event Subject Naming Convention

All NATS subjects follow this pattern:

```
{category}.{source}.{type}[.{detail}]

Examples:
- events.webhook.github.pr.opened
- events.webhook.clickup.task.created
- events.sensor.greenhouse.temp
- events.sensor.office.presence
- events.agent.researcher.completed
- events.system.backend.health
- cortex.heartbeat
- action.notification.sent
- action.light.turned_on
```

---

## Event Archival to TimescaleDB

The **EventArchiver** service runs independently:

```python
# backend/openmind/event_archiver.py

class EventArchiver:
    """
    Archives all NATS events to TimescaleDB.

    Runs continuously, subscribing to events.> and action.>
    """

    async def start(self):
        # Subscribe to ALL events
        await self.nats.subscribe(
            subject="events.>",
            callback=self.archive_event
        )

        await self.nats.subscribe(
            subject="action.>",
            callback=self.archive_event
        )

        await self.nats.subscribe(
            subject="cortex.>",
            callback=self.archive_event
        )

    async def archive_event(self, msg):
        """
        Archive a single event to TimescaleDB.
        """
        # Parse event
        event = json.loads(msg.data)

        # Insert into hypertable
        await self.timescale.execute("""
            INSERT INTO events (
                timestamp,
                subject,
                source,
                type,
                data
            ) VALUES ($1, $2, $3, $4, $5)
        """,
            event['timestamp'],
            msg.subject,
            event['source'],
            event['type'],
            event
        )

        # Acknowledge NATS message
        await msg.ack()
```

---

## Query Patterns

### Real-Time (NATS JetStream)

```python
# Get events from last 5 minutes
async def get_recent_events():
    """Query NATS for recent events (fast)."""
    events = []

    async with nats.subscribe("events.>") as sub:
        async for msg in sub.messages:
            if msg.timestamp > (now() - timedelta(minutes=5)):
                events.append(msg)

    return events
```

### Historical (TimescaleDB)

```sql
-- What did I do yesterday?
SELECT
    time_bucket('1 hour', timestamp) AS hour,
    type,
    COUNT(*) as count
FROM events
WHERE timestamp > NOW() - INTERVAL '1 day'
  AND source = 'agent'
GROUP BY hour, type
ORDER BY hour DESC;

-- Sensor readings over time
SELECT
    time_bucket('15 minutes', timestamp) AS bucket,
    AVG((data->>'value')::float) AS avg_temp
FROM events
WHERE source = 'sensor'
  AND type = 'greenhouse.temp'
  AND timestamp > NOW() - INTERVAL '7 days'
GROUP BY bucket
ORDER BY bucket;
```

---

## Event Gateway Implementation

The **EventGateway** normalizes all incoming events:

```python
# backend/events/gateway.py

class EventGateway:
    """
    Universal event ingestion endpoint.

    Receives events from any source and publishes to NATS.
    """

    async def ingest_webhook(self, source: str, payload: dict):
        """
        Webhook endpoint for external services.

        POST /webhooks/{source}
        """
        # Normalize event
        event = await self.normalize(source, payload)

        # Publish to NATS
        subject = f"events.webhook.{source}.{event.type}"
        await self.nats.publish(subject, event.to_json())

        # Return immediately (processing happens async)
        return {"status": "received", "event_id": event.id}

    async def ingest_sensor(self, location: str, sensor_type: str, value: float):
        """
        Sensor data ingestion.

        Called by ESP32, MQTT bridge, Home Assistant, etc.
        """
        event = {
            "timestamp": datetime.utcnow(),
            "source": "sensor",
            "type": f"{location}.{sensor_type}",
            "value": value
        }

        subject = f"events.sensor.{location}.{sensor_type}"
        await self.nats.publish(subject, json.dumps(event))
```

---

## WebSocket Real-Time Updates

Frontend receives events in real-time:

```javascript
// frontend/lib/services/websocket_service.dart

class WebSocketService {
  void connect() {
    ws = WebSocket('ws://localhost:7777/ws/events');

    ws.onMessage((event) {
      // Event from NATS, forwarded via WebSocket
      final data = jsonDecode(event);

      switch (data['subject']) {
        case 'events.webhook.github':
          showNotification('New GitHub PR');
          break;

        case 'events.sensor.greenhouse.temp':
          updateDashboard('greenhouse_temp', data['value']);
          break;

        case 'cortex.heartbeat':
          updateHeartbeatIndicator();
          break;
      }
    });
  }
}
```

---

## Summary

**Universal Event Flow:**
1. **Source** → EventGateway/MQTT/Agent → Publishes to NATS
2. **NATS** → Distributes to all subscribers (archiver, bridge, agents, websocket)
3. **OpenMind Bridge** → Processes events in 1 Hz loop (SENSE → FUSE → CORTEX → ACTION)
4. **EventArchiver** → Stores in TimescaleDB for long-term analytics
5. **Agents** → React to specific events
6. **WebSocket** → Broadcasts to frontend for real-time UI

**Key Points:**
- NATS is the **universal event bus** (nervous system)
- TimescaleDB is an **add-on** for long-term archival (not a replacement)
- Heartbeat loop **processes** events from NATS in real-time
- All services subscribe to NATS subjects they care about
- Events are immutable (publish once, consumed by many)

---

_Last updated: 2026-02-11_


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
