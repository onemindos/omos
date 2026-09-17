# OneMind OS - Integration Architecture

> How external services connect to OneMind OS and flow events to Legacy

---

## Overview

OneMind OS uses a **hybrid integration approach**:
- **Composio** for OAuth authentication (cloud services)
- **Direct webhooks** for ALL events (not limited by Composio triggers)
- **Direct API** for local services (Home Assistant, Frigate)

This gives you the best of both worlds: easy OAuth + unlimited real-time events.

---

## Visual Architecture

### High-Level Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        INTEGRATION ARCHITECTURE                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              YOUR UI                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  Settings → Integrations                                               │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐         │  │
│  │  │ Todoist │ │ GitHub  │ │  Home   │ │ Frigate │ │  +Add   │         │  │
│  │  │   ✅    │ │   ✅    │ │  Asst   │ │  ⚪     │ │  New    │         │  │
│  │  │         │ │         │ │   ✅    │ │         │ │         │         │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘         │  │
│  │                                                                        │  │
│  │  Click "Connect" → OAuth popup (Composio) → Done!                     │  │
│  │  Or: Enter URL + Token → Test → Connected                             │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           COMPOSIO (OAuth Only)                              │
│                                                                              │
│  • Handles OAuth for 250+ cloud services                                    │
│  • We use it ONLY for authentication, NOT for triggers                      │
│  • After OAuth, we configure webhooks to point to OUR server                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        EVENT INGESTION (ALL Events)                          │
└─────────────────────────────────────────────────────────────────────────────┘

  Cloud Services                  Local Services               IoT Devices
  (Todoist, GitHub)              (Home Assistant)             (Sensors, Cameras)
        │                              │                            │
        │ webhook                      │ webhook/WS                 │ MQTT
        ▼                              ▼                            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Redpanda Connect (Port 4196)                              │
│                    /webhooks                                                 │
│                                                                              │
│  • Receives ALL webhooks from all services                                  │
│  • Auto-detects service from headers/payload                                │
│  • Transforms to standard OneMind event format                              │
│  • Routes to NATS JetStream                                                 │
│                                                                              │
│  Headers: X-Webhook-Source: todoist                                         │
│           X-Hub-Signature-256: sha256=...                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         NATS JetStream                                       │
│                                                                              │
│  Stream: ONEMIND (7-day retention, 100k messages)                           │
│                                                                              │
│  Subjects:                                                                   │
│  • events.todoist.task.created                                              │
│  • events.todoist.task.completed                                            │
│  • events.github.pr.opened                                                  │
│  • events.github.issue.commented                                            │
│  • events.homeassistant.state.changed                                       │
│  • events.frigate.detection.person                                          │
│  • events.{service}.{event_type}                                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Consciousness Processor                                   │
│                                                                              │
│  Filters events based on:                                                   │
│  • Awareness level (DORMANT → AWARE → PRESENT → OMNIPRESENT)               │
│  • Time of day                                                              │
│  • Location                                                                 │
│  • Focus mode                                                               │
│  • Priority                                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Legacy's Inbox                                       │
│                                                                              │
│  Legacy receives filtered events and can:                                   │
│  • Respond to the user                                                      │
│  • Take actions via tools (Composio or direct API)                          │
│  • Update context and memory                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Integration Types

### Type 1: Cloud Services (via Composio OAuth)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     CLOUD SERVICES (Composio OAuth)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Services: Todoist, GitHub, Notion, Slack, Gmail, Spotify, etc.             │
│                                                                              │
│  AUTHENTICATION (Composio):                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │  1. User clicks "Connect" in UI                                         ││
│  │  2. UI calls POST /integrations/{id}/connect/oauth                      ││
│  │  3. Backend asks Composio for OAuth URL                                 ││
│  │  4. UI opens popup with OAuth URL                                       ││
│  │  5. User logs in to service (e.g., Todoist)                             ││
│  │  6. Composio stores the OAuth token                                     ││
│  │  7. Backend marks integration as connected                              ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
│  EVENTS (Direct Webhooks - NOT Composio triggers):                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │  1. User configures webhook in service (Todoist, GitHub, etc.)          ││
│  │  2. Webhook URL: https://your-domain.com/webhooks                       ││
│  │  3. Service sends ALL events to YOUR webhook                            ││
│  │  4. Redpanda Connect transforms and routes to NATS                      ││
│  │  5. You get EVERY event, not just what Composio exposes                 ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
│  ACTIONS (via Composio):                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │  Legacy: "Complete task 12345"                                          ││
│  │  → Tool calls Composio                                                  ││
│  │  → Composio uses stored OAuth token                                     ││
│  │  → Todoist API: POST /tasks/12345/close                                 ││
│  │  → Task completed!                                                      ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Type 2: Local Services (Direct API)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     LOCAL SERVICES (Direct API)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Services: Home Assistant, Frigate, local APIs                              │
│                                                                              │
│  AUTHENTICATION (Direct):                                                    │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │  1. User enters URL + API token in UI                                   ││
│  │  2. Backend tests connection                                            ││
│  │  3. If healthy, mark as connected                                       ││
│  │  4. Store URL/token in environment or secrets                           ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
│  EVENTS (Webhook or MQTT):                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │  Home Assistant:                                                        ││
│  │    - Create automation that POSTs to your webhook                       ││
│  │    - Or: Subscribe to WebSocket state changes                           ││
│  │                                                                          ││
│  │  Frigate:                                                               ││
│  │    - MQTT events: frigate/events                                        ││
│  │    - EMQX bridges MQTT → NATS                                           ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
│  ACTIONS (Direct API):                                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │  Legacy: "Turn on living room lights"                                   ││
│  │  → Tool calls Home Assistant directly                                   ││
│  │  → POST http://homeassistant:8123/api/services/light/turn_on            ││
│  │  → Headers: Authorization: Bearer {token}                               ││
│  │  → Lights on!                                                           ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Type 3: IoT/MQTT Devices

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     IoT DEVICES (MQTT)                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Devices: Sensors, cameras (via Frigate), ESP32, Zigbee/Z-Wave              │
│                                                                              │
│  EVENT FLOW:                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │                                                                          ││
│  │   Sensor/Device                                                         ││
│  │       │                                                                  ││
│  │       │ MQTT publish                                                     ││
│  │       ▼                                                                  ││
│  │   EMQX Broker (Port 1883)                                               ││
│  │       │                                                                  ││
│  │       │ Bridge                                                           ││
│  │       ▼                                                                  ││
│  │   NATS JetStream                                                        ││
│  │       │                                                                  ││
│  │       ▼                                                                  ││
│  │   Consciousness Processor                                               ││
│  │       │                                                                  ││
│  │       ▼                                                                  ││
│  │   Legacy's Inbox                                                        ││
│  │                                                                          ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Event Schema

All events follow a standard format:

```json
{
  "id": "evt_abc123",
  "type": "task.created",
  "source": "todoist",
  "timestamp": 1705123456.789,
  "priority": 2,
  "data": {
    // Original payload from the service
    "task_id": "12345",
    "content": "Review PR #42",
    "project_id": "67890"
  }
}
```

### NATS Subject Routing

Events are routed to subjects based on source and type:

| Subject Pattern | Example |
|-----------------|---------|
| `events.{source}.{type}` | `events.todoist.task.created` |
| `events.github.pr.*` | `events.github.pr.opened` |
| `events.homeassistant.state.*` | `events.homeassistant.state.changed` |

---

## Webhook Setup Guide

### For Cloud Services (After OAuth)

1. Go to the service's developer settings
2. Add webhook URL: `https://your-domain.com/webhooks`
3. Select **ALL events** (not just specific ones)
4. Copy webhook secret to `.env`

### Todoist

```bash
# 1. Go to: https://developer.todoist.com/appconsole.html
# 2. Select your app
# 3. Under Webhooks, add URL: https://your-domain.com/webhooks
# 4. Select ALL event types
# 5. Save webhook secret:
TODOIST_WEBHOOK_SECRET=your_secret_here
```

### GitHub

```bash
# 1. Go to: Repo Settings > Webhooks > Add webhook
# 2. Payload URL: https://your-domain.com/webhooks
# 3. Content type: application/json
# 4. Secret: Generate a random string
# 5. Events: Select "Send me everything"
GITHUB_WEBHOOK_SECRET=your_secret_here
```

### Home Assistant

```yaml
# In Home Assistant automations.yaml
automation:
  - alias: "Forward all events to OneMind"
    trigger:
      - platform: event
        event_type: state_changed
    action:
      - service: rest_command.onemind_webhook
        data:
          payload: "{{ trigger.event | to_json }}"

rest_command:
  onemind_webhook:
    url: "https://your-domain.com/webhooks"
    method: POST
    headers:
      X-Webhook-Source: "homeassistant"
    content_type: "application/json"
    payload: "{{ payload }}"
```

---

## API Endpoints

### List Integrations

```http
GET /integrations
```

Returns all integrations grouped by category with status.

### Get Integration Details

```http
GET /integrations/{id}
```

Returns detailed info including webhook setup instructions.

### Connect (OAuth)

```http
POST /integrations/{id}/connect/oauth
```

Returns redirect URL for OAuth popup.

### Connect (Direct)

```http
POST /integrations/{id}/connect/direct
Content-Type: application/json

{
  "url": "http://homeassistant:8123",
  "token": "your_api_token"
}
```

### Disconnect

```http
POST /integrations/{id}/disconnect
```

### Health Check

```http
GET /integrations/{id}/health
```

### Summary

```http
GET /integrations/status/summary
```

Quick status for dashboard widgets.

---

## File Structure

```
agents/legacy/integrations/
├── __init__.py         # Module exports
├── types.py            # Type definitions (Integration, State, etc.)
├── registry.py         # Pre-defined integrations (Todoist, GitHub, etc.)
├── manager.py          # Connection management (OAuth, health checks)
└── api.py              # FastAPI router (mount in main.py)

ui/src/lib/integrations/
├── index.ts            # Module exports
├── types.ts            # TypeScript types
└── hooks.ts            # React hooks (useIntegrations, useIntegrationActions)

ui/src/app/settings/integrations/
└── page.tsx            # Integrations settings page
```

---

## Adding a New Integration

1. **Define in registry.py:**

```python
register(Integration(
    id="newservice",
    name="New Service",
    description="Description here",
    icon="icon-name",
    category=IntegrationCategory.PRODUCTIVITY,
    connection=ConnectionConfig(
        type=IntegrationType.COMPOSIO,  # or DIRECT_API, MQTT
        composio_app_name="newservice",
    ),
    webhook=WebhookConfig(
        endpoint_path="/webhooks",
        setup_instructions="...",
    ),
    event_types=["event.type1", "event.type2"],
    actions=["action1", "action2"],
))
```

2. **Update Redpanda Connect (optional):**

If the service has unique payload patterns, add detection logic to `connect.yaml`.

3. **The UI auto-discovers it!**

No frontend changes needed - the new integration appears automatically.

---

## Summary

| Component | Purpose |
|-----------|---------|
| **Composio** | OAuth authentication for cloud services |
| **Direct Webhooks** | Receive ALL events (not limited by triggers) |
| **Redpanda Connect** | Transform webhooks → NATS events |
| **EMQX** | Bridge MQTT → NATS for IoT |
| **NATS JetStream** | Event backbone (pub/sub) |
| **Consciousness Processor** | Filter by awareness level |
| **Integration Registry** | Single source of truth for all integrations |
| **Settings UI** | Click-to-connect experience |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
