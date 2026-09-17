# AG-UI Protocol Integration

> Real-time bi-directional state synchronization between Legacy and UIs

## What is AG-UI?

AG-UI (Agent-User Interaction Protocol) is an open, event-based protocol that standardizes how AI agents connect to frontend applications. It provides:

- **Real-time streaming**: Token-by-token response streaming
- **Bi-directional state sync**: State flows both ways (agent ↔ UI)
- **Human-in-the-loop**: Frontend tools for user intervention
- **Generative UI**: Agents can create UI components

## Why AG-UI for Legacy?

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Legacy STATE SYNC VIA AG-UI                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  WITHOUT AG-UI (polling):                                                │
│  ┌──────────┐     GET /state      ┌──────────┐                          │
│  │    UI    │ ──────────────────▶ │  Legacy  │                          │
│  │          │ ◀────────────────── │          │                          │
│  └──────────┘    every 5 sec      └──────────┘                          │
│                  ❌ Slow, wasteful                                        │
│                                                                          │
│  WITH AG-UI (events):                                                    │
│  ┌──────────┐     SSE stream      ┌──────────┐                          │
│  │    UI    │ ◀═══════════════════│  Legacy  │                          │
│  │          │ ═══════════════════▶│          │                          │
│  └──────────┘   StateSnapshot     └──────────┘                          │
│                 StateDelta                                               │
│                 ✅ Real-time, efficient                                  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## AG-UI Event Types

### Lifecycle Events

```
┌─────────────────────────────────────────────────────────────────────────┐
│  LIFECYCLE EVENTS - Track agent run progression                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  RunStarted ────▶ StepStarted ────▶ StepFinished ────▶ RunFinished      │
│       │                                                      │           │
│       │              (on error)                              │           │
│       └──────────────────────▶ RunError ◀────────────────────┘           │
│                                                                          │
│  Event Payloads:                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ RunStarted:  { runId, threadId, timestamp }                     │    │
│  │ StepStarted: { runId, stepId, name, input }                     │    │
│  │ StepFinished:{ runId, stepId, output }                          │    │
│  │ RunFinished: { runId, result }                                  │    │
│  │ RunError:    { runId, error, code }                             │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Message Events (Streaming)

```
┌─────────────────────────────────────────────────────────────────────────┐
│  MESSAGE EVENTS - Stream text responses token-by-token                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  TextMessageStart ──▶ TextMessageContent ──▶ TextMessageEnd             │
│        │                     │ │ │                    │                  │
│        │                     ▼ ▼ ▼                    │                  │
│        │              (multiple deltas)               │                  │
│        │                                              │                  │
│        └──────────────────────────────────────────────┘                  │
│                                                                          │
│  Example stream:                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ { type: "TextMessageStart", messageId: "msg_1", role: "agent" } │    │
│  │ { type: "TextMessageContent", delta: "Hello" }                  │    │
│  │ { type: "TextMessageContent", delta: ", Zeus" }                 │    │
│  │ { type: "TextMessageContent", delta: "!" }                      │    │
│  │ { type: "TextMessageEnd" }                                      │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### State Events (The Key Feature!)

```
┌─────────────────────────────────────────────────────────────────────────┐
│  STATE EVENTS - Bi-directional state synchronization                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  SNAPSHOT: Complete state at a point in time                             │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ {                                                               │    │
│  │   type: "StateSnapshot",                                        │    │
│  │   state: {                                                      │    │
│  │     location: { place: "Home", lat: 37.77, lng: -122.41 },     │    │
│  │     biometrics: { hr: 72, stress: 0.3, steps: 4521 },          │    │
│  │     activity: { type: "working", duration: 45 },               │    │
│  │     devices: { watch: true, phone: true }                      │    │
│  │   }                                                             │    │
│  │ }                                                               │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  DELTA: Incremental update (RFC 6902 JSON Patch)                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ {                                                               │    │
│  │   type: "StateDelta",                                           │    │
│  │   delta: [                                                      │    │
│  │     { op: "replace", path: "/biometrics/hr", value: 75 },      │    │
│  │     { op: "replace", path: "/biometrics/stress", value: 0.4 }, │    │
│  │     { op: "replace", path: "/activity/duration", value: 46 }   │    │
│  │   ]                                                             │    │
│  │ }                                                               │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  FLOW:                                                                   │
│  1. On connect: Send StateSnapshot (full state)                          │
│  2. On change: Send StateDelta (only what changed)                       │
│  3. On reconnect: Send StateSnapshot again                               │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Tool Call Events

```
┌─────────────────────────────────────────────────────────────────────────┐
│  TOOL EVENTS - Agent tool execution                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ToolCallStart ──▶ ToolCallArgs ──▶ ToolCallEnd ──▶ ToolCallResult      │
│                                                                          │
│  Example (Legacy checking calendar):                                     │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ { type: "ToolCallStart", toolId: "t_1", name: "check_calendar" }│    │
│  │ { type: "ToolCallArgs", delta: '{"date": "today"}' }            │    │
│  │ { type: "ToolCallEnd" }                                         │    │
│  │ { type: "ToolCallResult", result: { meetings: [...] } }         │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Integration Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Legacy + AG-UI ARCHITECTURE                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                         FRONTEND (UI)                           │    │
│  │  ┌──────────────────────────────────────────────────────────┐  │    │
│  │  │  CopilotKit / Custom React                               │  │    │
│  │  │  - Renders StateSnapshot                                 │  │    │
│  │  │  - Applies StateDelta patches                            │  │    │
│  │  │  - Streams messages in real-time                         │  │    │
│  │  │  - Shows tool executions                                 │  │    │
│  │  └──────────────────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                              ▲                                           │
│                              │ SSE / WebSocket                           │
│                              │ (AG-UI events)                            │
│                              ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                         AGNO AGUI                               │    │
│  │  ┌──────────────────────────────────────────────────────────┐  │    │
│  │  │  POST /agui                                              │  │    │
│  │  │  - Accepts RunAgentInput                                 │  │    │
│  │  │  - Streams AG-UI events                                  │  │    │
│  │  │  - Injects state into context                            │  │    │
│  │  └──────────────────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                              ▲                                           │
│                              │                                           │
│                              ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Legacy AGENT                               │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │    │
│  │  │ StateManager│  │PatternEngine│  │     Claude Opus 4.5     │ │    │
│  │  │  (state.py) │  │(patterns.py)│  │     (via Bedrock)       │ │    │
│  │  └──────┬──────┘  └──────┬──────┘  └─────────────────────────┘ │    │
│  │         │                │                                      │    │
│  │         └────────┬───────┘                                      │    │
│  │                  ▼                                              │    │
│  │         ┌─────────────────┐                                     │    │
│  │         │  AG-UI Bridge   │ ← Converts state to AG-UI events    │    │
│  │         └─────────────────┘                                     │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                              ▲                                           │
│                              │ Events (NATS)                             │
│                              ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      DATA SOURCES                               │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐            │    │
│  │  │ T-Watch │  │  Phone  │  │ Calendar│  │ Glasses │            │    │
│  │  │  (HR)   │  │  (GPS)  │  │ (Events)│  │ (View)  │            │    │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘            │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Implementation

### 1. Agno AGUI Setup

```python
# agents/legacy/main.py
from agno.agent import Agent
from agno.interfaces.agui import AGUI
from agno.db.redis import RedisDb
import os

# Create Legacy agent
agent = Agent(
    name="Legacy",
    instructions=Legacy_INSTRUCTIONS,
    db=RedisDb(db_url=os.environ.get("REDIS_URL")),
)

# Wrap in AG-UI interface
agui = AGUI(agent=agent)

# Get FastAPI router
router = agui.get_router()

# Mount on your app
app.include_router(router, prefix="/api")
```

### 2. AG-UI State Bridge

```python
# agents/legacy/consciousness/agui_bridge.py
"""
Bridge between Legacy StateManager and AG-UI events.

Converts state changes to StateSnapshot/StateDelta events
for real-time UI synchronization.
"""

import json
from dataclasses import dataclass, asdict
from typing import Any, Optional
from enum import Enum
import jsonpatch  # RFC 6902 JSON Patch

from agents.legacy.consciousness.state import get_state_manager, LegacyState


class AGUIEventType(str, Enum):
    """AG-UI event types we emit."""
    STATE_SNAPSHOT = "StateSnapshot"
    STATE_DELTA = "StateDelta"
    ACTIVITY_SNAPSHOT = "ActivitySnapshot"
    ACTIVITY_DELTA = "ActivityDelta"
    CUSTOM = "Custom"


@dataclass
class AGUIEvent:
    """AG-UI protocol event."""
    type: AGUIEventType
    data: dict

    def to_sse(self) -> str:
        """Format as Server-Sent Event."""
        return f"data: {json.dumps({'type': self.type.value, **self.data})}\n\n"


class AGUIStateBridge:
    """
    Bridges Legacy state to AG-UI events.

    Tracks previous state to compute JSON Patch deltas.
    """

    def __init__(self):
        self._last_state: Optional[dict] = None
        self._last_activity: Optional[dict] = None

    def get_state_snapshot(self) -> AGUIEvent:
        """
        Get complete state snapshot for initial sync.

        Called on:
        - Client connect
        - Client reconnect
        - Explicit refresh
        """
        state_manager = get_state_manager()
        if not state_manager:
            return AGUIEvent(
                type=AGUIEventType.STATE_SNAPSHOT,
                data={"state": {}}
            )

        state = state_manager.get_state()
        state_dict = self._state_to_agui_format(state)

        # Store for delta computation
        self._last_state = state_dict.copy()

        return AGUIEvent(
            type=AGUIEventType.STATE_SNAPSHOT,
            data={"state": state_dict}
        )

    def get_state_delta(self) -> Optional[AGUIEvent]:
        """
        Get state delta (JSON Patch) if state changed.

        Called periodically or on state change notification.
        Returns None if no changes.
        """
        state_manager = get_state_manager()
        if not state_manager or self._last_state is None:
            return None

        state = state_manager.get_state()
        current_dict = self._state_to_agui_format(state)

        # Compute JSON Patch (RFC 6902)
        patch = jsonpatch.make_patch(self._last_state, current_dict)

        if not patch.patch:
            return None  # No changes

        # Update last state
        self._last_state = current_dict.copy()

        return AGUIEvent(
            type=AGUIEventType.STATE_DELTA,
            data={"delta": patch.patch}
        )

    def get_activity_snapshot(self) -> AGUIEvent:
        """
        Get activity state (what Legacy is doing).

        Includes:
        - Current thinking/processing status
        - Active tools
        - Pending insights
        """
        from agents.legacy.consciousness.patterns import get_pattern_engine

        pattern_engine = get_pattern_engine()

        activity = {
            "status": "idle",
            "currentStep": None,
            "pendingInsights": [],
            "biometricContext": "",
        }

        if pattern_engine:
            biometric_context = pattern_engine.get_biometric_context()
            activity["biometricContext"] = biometric_context

            # Get pending insights
            # (would need async call in real impl)
            activity["pendingInsights"] = []

        self._last_activity = activity.copy()

        return AGUIEvent(
            type=AGUIEventType.ACTIVITY_SNAPSHOT,
            data={"activity": activity}
        )

    def emit_custom_event(self, name: str, data: dict) -> AGUIEvent:
        """
        Emit custom event for app-specific needs.

        Examples:
        - biometric_alert: High stress detected
        - location_change: User arrived at office
        - insight_generated: New pattern detected
        """
        return AGUIEvent(
            type=AGUIEventType.CUSTOM,
            data={
                "name": name,
                **data
            }
        )

    def _state_to_agui_format(self, state: LegacyState) -> dict:
        """Convert LegacyState to AG-UI compatible format."""
        return {
            "user": {
                "id": state.user_id,
                "name": state.user_name,
            },
            "location": {
                "place": state.location.place_name,
                "city": state.location.city,
                "latitude": state.location.latitude,
                "longitude": state.location.longitude,
                "confidence": state.location.confidence,
            },
            "biometrics": {
                "heartRate": state.heart_rate,
                "hrv": state.hrv,
                "stressScore": state.stress_score,
                "fatigueScore": state.fatigue_score,
                "stepsToday": state.steps_today,
            },
            "activity": {
                "type": state.context.activity.value if state.context.activity else "unknown",
                "task": state.context.current_task,
                "app": state.context.current_app,
                "inMeeting": state.context.in_meeting,
            },
            "devices": {
                device_id: {
                    "name": device.name,
                    "type": device.device_type.value,
                    "active": device.is_active,
                    "battery": device.battery_percent,
                }
                for device_id, device in state.devices.items()
            },
            "awareness": {
                "mode": state.awareness_mode,
                "lastUpdate": state.last_state_save.isoformat() if state.last_state_save else None,
            },
        }


# Singleton
_bridge: Optional[AGUIStateBridge] = None


def get_agui_bridge() -> AGUIStateBridge:
    """Get or create the AG-UI bridge."""
    global _bridge
    if _bridge is None:
        _bridge = AGUIStateBridge()
    return _bridge
```

### 3. Custom AG-UI Events for Legacy

```python
# agents/legacy/tools/agui_events.py
"""
Custom AG-UI events that Legacy can emit from tools.
"""

from dataclasses import dataclass
from typing import Optional
from agno.tools import tool


@dataclass
class BiometricAlertEvent:
    """Emitted when biometric thresholds are exceeded."""
    alert_type: str  # "stress", "fatigue", "spo2_low"
    severity: str    # "warning", "critical"
    message: str
    current_value: float
    threshold: float
    suggested_action: Optional[str] = None


@dataclass
class LocationChangeEvent:
    """Emitted when user location changes significantly."""
    previous_place: Optional[str]
    current_place: str
    geofence_triggered: Optional[str] = None  # "home", "work", "gym"


@dataclass
class InsightGeneratedEvent:
    """Emitted when pattern engine generates an insight."""
    insight_id: str
    message: str
    priority: str  # "low", "medium", "high", "urgent"
    action_suggested: Optional[str] = None
    expires_at: Optional[str] = None


@tool
def emit_biometric_alert(
    alert_type: str,
    severity: str,
    message: str,
    current_value: float,
    threshold: float,
    suggested_action: str = None,
) -> BiometricAlertEvent:
    """
    Emit a biometric alert to the UI.

    Use when biometric readings exceed safe thresholds.

    Args:
        alert_type: Type of alert (stress, fatigue, spo2_low)
        severity: Alert severity (warning, critical)
        message: Human-readable alert message
        current_value: Current biometric value
        threshold: Threshold that was exceeded
        suggested_action: Optional action to suggest
    """
    yield BiometricAlertEvent(
        alert_type=alert_type,
        severity=severity,
        message=message,
        current_value=current_value,
        threshold=threshold,
        suggested_action=suggested_action,
    )
    return {"emitted": True}
```

### 4. Frontend Integration (CopilotKit)

```tsx
// ui/components/LegacyChat.tsx
import { CopilotKit, useCopilotContext } from "@copilotkit/react-core";
import { CopilotSidebar } from "@copilotkit/react-ui";
import { useEffect, useState } from "react";

interface LegacyState {
  location: { place: string; city: string };
  biometrics: {
    heartRate: number;
    stressScore: number;
    fatigueScore: number;
    stepsToday: number;
  };
  activity: { type: string; task: string; inMeeting: boolean };
  devices: Record<string, { name: string; active: boolean; battery: number }>;
}

export function LegacyChat() {
  const [state, setState] = useState<LegacyState | null>(null);

  return (
    <CopilotKit
      runtimeUrl="/api/agui"
      // State is synced automatically via AG-UI events
      onStateSnapshot={(snapshot) => setState(snapshot.state)}
      onStateDelta={(delta) => {
        // Apply JSON Patch to current state
        setState((prev) => applyPatch(prev, delta.delta));
      }}
    >
      <div className="flex h-screen">
        {/* State Display Panel */}
        <div className="w-64 border-r p-4">
          <AwarenessPanel state={state} />
        </div>

        {/* Chat Interface */}
        <CopilotSidebar
          labels={{
            title: "Legacy",
            initial: "How can I help you today?",
          }}
        />
      </div>
    </CopilotKit>
  );
}

function AwarenessPanel({ state }: { state: LegacyState | null }) {
  if (!state) return <div>Connecting...</div>;

  return (
    <div className="space-y-4">
      {/* Location */}
      <div>
        <h3 className="font-semibold">📍 Location</h3>
        <p>{state.location.place || state.location.city || "Unknown"}</p>
      </div>

      {/* Biometrics */}
      <div>
        <h3 className="font-semibold">💓 Biometrics</h3>
        <div className="grid grid-cols-2 gap-2 text-sm">
          <div>HR: {state.biometrics.heartRate} bpm</div>
          <div>Steps: {state.biometrics.stepsToday}</div>
          <div>Stress: {(state.biometrics.stressScore * 100).toFixed(0)}%</div>
          <div>Fatigue: {(state.biometrics.fatigueScore * 100).toFixed(0)}%</div>
        </div>
      </div>

      {/* Activity */}
      <div>
        <h3 className="font-semibold">🎯 Activity</h3>
        <p>{state.activity.type}</p>
        {state.activity.task && (
          <p className="text-sm text-gray-600">{state.activity.task}</p>
        )}
        {state.activity.inMeeting && (
          <span className="text-red-500 text-sm">In Meeting</span>
        )}
      </div>

      {/* Devices */}
      <div>
        <h3 className="font-semibold">📱 Devices</h3>
        {Object.entries(state.devices).map(([id, device]) => (
          <div key={id} className="flex items-center gap-2 text-sm">
            <span className={device.active ? "text-green-500" : "text-gray-400"}>
              ●
            </span>
            <span>{device.name}</span>
            {device.battery && <span>({device.battery}%)</span>}
          </div>
        ))}
      </div>
    </div>
  );
}
```

## Event Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AG-UI EVENT FLOW                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  1. USER CONNECTS                                                        │
│  ─────────────────                                                       │
│  UI ──── POST /agui ────▶ Agno AGUI                                     │
│                             │                                            │
│                             ▼                                            │
│                    ┌─────────────────┐                                   │
│                    │ StateSnapshot   │ ◀── Full state sent               │
│                    │ ActivitySnapshot│                                   │
│                    └─────────────────┘                                   │
│                             │                                            │
│  UI ◀───────────────────────┘                                           │
│                                                                          │
│  2. USER SENDS MESSAGE                                                   │
│  ─────────────────────────                                               │
│  UI ──── "What's my heart rate?" ────▶ Legacy                           │
│                                           │                              │
│                                           ▼                              │
│                                    ┌─────────────┐                       │
│                                    │ RunStarted  │                       │
│                                    │ TextMessage*│ (streaming response)  │
│                                    │ ToolCall*   │ (if checking sensors) │
│                                    │ StateDelta  │ (if state changed)    │
│                                    │ RunFinished │                       │
│                                    └─────────────┘                       │
│                                           │                              │
│  UI ◀─────────────────────────────────────┘                             │
│                                                                          │
│  3. BIOMETRIC ALERT                                                      │
│  ──────────────────                                                      │
│  T-Watch ──▶ NATS ──▶ PatternEngine ──▶ InsightGenerated                │
│                                               │                          │
│                                               ▼                          │
│                                    ┌─────────────────┐                   │
│                                    │ Custom Event:   │                   │
│                                    │ BiometricAlert  │                   │
│                                    │ StateDelta      │                   │
│                                    └─────────────────┘                   │
│                                               │                          │
│  UI ◀─────────────────────────────────────────┘                         │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## State Sync Comparison

| Approach | Latency | Bandwidth | Complexity |
|----------|---------|-----------|------------|
| Polling (GET /state) | 5s+ | High (full state each time) | Low |
| WebSocket (custom) | <100ms | Medium | High |
| **AG-UI (snapshot+delta)** | **<100ms** | **Low (only changes)** | **Medium** |

## Why This is Better

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AG-UI ADVANTAGES                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  1. STANDARDIZED PROTOCOL                                                │
│     - Works with CopilotKit, any AG-UI client                           │
│     - Future-proof (growing ecosystem)                                   │
│     - No custom WebSocket protocol to maintain                           │
│                                                                          │
│  2. EFFICIENT STATE SYNC                                                 │
│     - Snapshot on connect (full state)                                   │
│     - Delta on change (JSON Patch, only what changed)                   │
│     - Automatic reconnect handling                                       │
│                                                                          │
│  3. BUILT INTO AGNO                                                      │
│     - from agno.interfaces.agui import AGUI                             │
│     - One line to expose agent as AG-UI                                 │
│     - Custom events from tools                                           │
│                                                                          │
│  4. BI-DIRECTIONAL                                                       │
│     - UI → Agent: User input, context enrichment                        │
│     - Agent → UI: Responses, state, alerts                              │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Next Steps

1. ✅ **Understand AG-UI protocol** - Done
2. **Implement AGUIStateBridge in state.py**
3. **Add custom events for biometric alerts**
4. **Update UI to use CopilotKit with AG-UI**
5. **Test real-time state sync**
