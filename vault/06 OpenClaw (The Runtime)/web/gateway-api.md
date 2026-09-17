---
title: "Interface — Gateway API (REST + WebSocket)"
type: reference
status: active
created: 2026-04-12
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/api]
ai_summary: "OpenClaw Gateway API — REST endpoints, NATS WebSocket, used by N8N, integrations, and custom code."
---

# Interface: Gateway API (REST + WebSocket)

**Status:** ✅ Active
**Type:** REST + WebSocket
**Primary use:** N8N workflows, integrations, custom code, NATS direct access

---

## Base URLs

```
OpenClaw HTTP:     http://100.77.52.93:443
OpenClaw WS:       ws://100.77.52.93:443
NATS WebSocket:    ws://100.77.52.93:9222
```

All require Tailscale. Auth via `Authorization: Bearer <access-key>` header.

---

## REST Endpoints

```bash
# Send a message to an agent (start or continue a session)
POST /api/sessions/{session_id}/messages
{
  "content": "What's my current task priority?",
  "role": "user"
}

# Create a new session
POST /api/sessions
{
  "agent_id": "legacy",
  "title": "Morning briefing"
}

# List all agents
GET /api/agents

# Get session message history
GET /api/sessions/{session_id}/messages

# Stream a response (Server-Sent Events)
GET /api/sessions/{session_id}/stream
```

---

## NATS WebSocket (Direct)

For code that needs to subscribe to the live NATS feed directly (no agent intermediary):

```javascript
import { connect } from "nats.ws"

const nc = await connect({ servers: "ws://100.77.52.93:9222" })

// Subscribe to Zeus presence updates
const sub = nc.subscribe("zeus.presence.>")
for await (const msg of sub) {
    const data = JSON.parse(new TextDecoder().decode(msg.data))
    console.log(msg.subject, data)
}

// Publish a signal
await nc.publish("zeus.signal.event", JSON.stringify({
    level: "medium",
    message: "N8N workflow completed",
    source: "n8n"
}))
```

---

## Common Use Cases

- **N8N workflows** → POST `/api/sessions/{id}/messages` to trigger agent
- **Cron jobs** → publish to `zeus.signal.event` for timed alerts
- **Custom scripts** → NATS WebSocket to subscribe/publish directly
- **onemind-mobile** → REST API for all agent interactions

---

## Notes

- Full API spec: `onemind-web/src/cli/spec.js` and route files under `src/app/api/`
- Access key stored in `/opt/onemind-web/data/` and in **1Password → SO vault → OpenClaw Access Key**


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
