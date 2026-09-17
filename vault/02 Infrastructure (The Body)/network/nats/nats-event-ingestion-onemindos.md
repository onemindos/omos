# NATS Event Ingestion — One Mind OS
**How NATS receives events across the entire Legacy Platform stack**

---

## What is NATS?

NATS is a lightweight, high-performance open-source messaging system. It is the **central nervous system** of One Mind OS — every event, command, agent trigger, IoT signal, and SaaS webhook flows through it.

NATS does not speak HTTP natively. It operates on its own protocol over TCP (port 4222). This means external sources need either a **NATS client SDK** or a **bridge layer** to publish events into the system.

---

## The 7 Ways NATS Receives Events

---

### 1. Direct Publish via NATS Client SDK

The fastest and most native path. Any internal service, agent, or script with a NATS client library publishes directly to a subject.

**Available SDKs:** Python, Go, Node.js, Java, Rust, C, .NET, Ruby, and more.

```python
# Python example
import nats

async def publish():
    nc = await nats.connect("nats://localhost:4222")
    await nc.publish("legacy.saas.zoho.contact.created", b'{"id": "123"}')
    await nc.close()
```

```javascript
// Node.js example
const nc = await connect({ servers: "nats://localhost:4222" });
nc.publish("legacy.agent.loop.completed", sc.encode({ status: "done" }));
```

**Best for:**
- Internal One Mind OS services
- OpenClaw agent loop publishing results
- Custom scripts and automation

---

### 2. HTTP → NATS Bridge (Webhook Gateway)

NATS doesn't natively accept HTTP POST requests. A bridge service sits between your SaaS webhooks and NATS — it receives HTTP, validates the payload, then publishes to the correct NATS subject.

**Implementation options:**

**Option A: n8n (Already in your stack)**
- Configure a Webhook Trigger node in n8n
- Point your SaaS app at the n8n webhook URL
- Add a NATS publish step to route the event
- Zero new infrastructure required

**Option B: Lightweight Custom Gateway (Recommended long-term)**
A small Go or Python service:
```python
# FastAPI webhook gateway example
from fastapi import FastAPI, Request
import nats

app = FastAPI()
nc = None

@app.on_event("startup")
async def startup():
    global nc
    nc = await nats.connect("nats://localhost:4222")

@app.post("/webhook/zoho")
async def zoho_webhook(request: Request):
    payload = await request.body()
    await nc.publish("legacy.saas.zoho.events", payload)
    return {"status": "received"}

@app.post("/webhook/clickup")
async def clickup_webhook(request: Request):
    payload = await request.body()
    await nc.publish("legacy.saas.clickup.events", payload)
    return {"status": "received"}
```

**SaaS apps that send webhooks to this gateway:**
- Zoho CRM / Zoho Mail
- ClickUp
- Gmail (via Google Pub/Sub)
- GoDaddy
- Cloudflare
- Vercel
- Any app with webhook support

---

### 3. JetStream HTTP API

JetStream exposes a built-in management API. While primarily for admin operations, it can be used to publish messages directly into streams via HTTP — useful for services that cannot install a NATS client.

```bash
# Publish directly to a JetStream stream via HTTP
curl -X POST http://localhost:8222/js/api/STREAM.MSG.PUT \
  -d '{"subject": "legacy.saas.events", "payload": "base64encodedpayload"}'
```

**Best for:**
- Serverless functions
- Third-party services with limited runtime environments
- Quick testing and debugging

---

### 4. MQTT Bridge (IoT & Smart Home)

NATS speaks **MQTT natively** — no bridge required. Any MQTT-capable device publishes directly to NATS. This is critical for the homestead, farm, and office smart home layer.

**Enable MQTT in NATS server config:**
```
mqtt {
  port: 1883
}
```

**Devices that connect via MQTT:**
- Home Assistant (HA → MQTT → NATS directly)
- Frigate NVR camera events
- Mason Android devices (MASC09T / C210)
- Echo Show 5 (via LineageOS / Home Assistant)
- Environmental sensors (temperature, motion, water)
- Solar/energy monitoring devices
- Any IoT sensor on the homestead/farm/office mesh

```
# MQTT topic → NATS subject mapping
mqtt topic:  home/frigate/camera/motion
nats subject: home.frigate.camera.motion

mqtt topic:  farm/sensors/temperature
nats subject: farm.sensors.temperature
```

**Best for:**
- All smart home and IoT events
- Home Assistant integration
- Edge sensor data from homestead and farm locations

---

### 5. WebSocket

NATS supports native WebSocket connections. Browser-based clients and real-time frontends can publish events directly without HTTP.

**Enable WebSocket in NATS server config:**
```
websocket {
  port: 9222
  no_tls: true  # Use TLS in production
}
```

```javascript
// Browser client connecting to NATS via WebSocket
import { connect, wsconnect } from "@nats-io/transport-browser";

const nc = await wsconnect({ servers: "ws://localhost:9222" });
nc.publish("legacy.community.member.joined", sc.encode({ userId: "abc123" }));
```

**Best for:**
- Legacy Community web portal events
- Real-time dashboard updates
- Browser-based agent interactions

---

### 6. NATS Leaf Nodes (Multi-Location Mesh)

Leaf Nodes are lightweight NATS connections that link remote locations back to your central NATS cluster — forwarding local events up to the hub and receiving commands back down.

**Perfect for the One Mind OS multi-location architecture:**

```
[Home]          [Farm]          [Office]
Leaf Node   →   Leaf Node   →   Leaf Node
    ↓               ↓               ↓
         Central NATS Cluster (DigitalOcean)
                    ↓
              OpenClaw Gateway
```

**Leaf node config example:**
```
leafnodes {
  remotes [
    {
      url: "nats://central.legacyhq.ai:7422"
      credentials: "/etc/nats/node.creds"
    }
  ]
}
```

**Events forwarded from leaf nodes:**
- `home.>` — All home events
- `farm.>` — All farm/homestead events
- `office.>` — All office events

**Secured via Tailscale mesh networking** — each location's leaf node connects through the Tailscale VPN, never exposing NATS directly to the public internet.

---

### 7. NATS Service (Request/Reply)

Beyond events, NATS supports synchronous **Request/Reply** — an agent or service sends a request to a subject and waits for a response. OpenClaw uses this pattern internally.

```python
# Request pattern — agent asking a service for data
response = await nc.request("legacy.omcs.coordinate.lookup", 
                            b'{"coordinate": "HP.71.P01"}', 
                            timeout=2.0)
print(response.data)
```

**Best for:**
- OpenClaw agent tool calls
- OMCS coordinate lookups
- Real-time queries between agents

---

## Subject Hierarchy — One Mind OS

Subjects are dot-separated addresses. Every event in One Mind OS has a structured subject that maps to the OMCS coordinate system and Legacy Codex pillars.

### Structure
```
{domain}.{source}.{entity}.{action}
```

### Full Subject Map

```
# SaaS Integration Events
legacy.saas.zoho.contact.created
legacy.saas.zoho.contact.updated
legacy.saas.zoho.deal.closed
legacy.saas.clickup.task.created
legacy.saas.clickup.task.completed
legacy.saas.gmail.email.received
legacy.saas.gmail.email.sent
legacy.saas.godaddy.domain.renewed
legacy.saas.cloudflare.dns.updated
legacy.saas.vercel.deployment.completed

# Agent Loop Events
legacy.agent.loop.started
legacy.agent.loop.completed
legacy.agent.tool.called
legacy.agent.tool.result
legacy.agent.error

# OMCS Coordinate Events
legacy.omcs.HP.>    # Holistic Performance pillar
legacy.omcs.LE.>    # Legacy Evolution pillar
legacy.omcs.GE.>    # Generational Entrepreneurship pillar
legacy.omcs.IT.>    # Innovative Technology pillar

# Home Events (via MQTT bridge or Leaf Node)
legacy.home.frigate.camera.motion
legacy.home.frigate.camera.person
legacy.home.assistant.automation.triggered
legacy.home.sensor.temperature
legacy.home.sensor.door.opened
legacy.home.energy.solar.update

# Farm / Homestead Events
legacy.farm.sensor.soil.moisture
legacy.farm.sensor.weather
legacy.farm.security.motion

# Office Events
legacy.office.device.mason.online
legacy.office.network.tailscale.connected

# Community / Platform Events
legacy.community.member.joined
legacy.community.course.completed
legacy.community.payment.received
legacy.community.certification.issued

# System Events
legacy.system.openclaw.gateway.health
legacy.system.nats.stream.error
legacy.system.supabase.connection.status
```

### Wildcard Subscriptions

```
legacy.saas.>         → All SaaS events
legacy.home.>         → All home/smart home events
legacy.farm.>         → All farm/homestead events
legacy.agent.>        → All agent loop events
legacy.omcs.HP.>      → All HP pillar events
legacy.>              → Everything (use sparingly)
```

---

## JetStream Streams — Persistence Layer

Raw NATS subjects are ephemeral. JetStream adds persistence, replay, and delivery guarantees. Define streams that capture subjects into durable storage.

### Recommended Stream Configuration

| Stream Name | Subjects | Retention | Use Case |
|---|---|---|---|
| `SAAS_EVENTS` | `legacy.saas.>` | 7 days | All SaaS webhook events |
| `AGENT_EVENTS` | `legacy.agent.>` | 30 days | Agent loop audit trail |
| `OMCS_EVENTS` | `legacy.omcs.>` | 90 days | Life coordinate changes |
| `HOME_EVENTS` | `legacy.home.>` | 3 days | Smart home/IoT events |
| `FARM_EVENTS` | `legacy.farm.>` | 7 days | Homestead sensor data |
| `COMMUNITY_EVENTS` | `legacy.community.>` | 365 days | Platform activity log |
| `SYSTEM_EVENTS` | `legacy.system.>` | 1 day | Health and diagnostics |

---

## Complete One Mind OS Event Flow

```
┌─────────────────────────────────────────────────────────┐
│                    EVENT SOURCES                         │
├──────────────┬──────────────┬──────────────┬────────────┤
│  SaaS Apps   │  IoT/MQTT    │  Browser     │  Internal  │
│  Zoho        │  Home Asst.  │  Community   │  Services  │
│  ClickUp     │  Frigate     │  Portal      │  OpenClaw  │
│  Gmail       │  Sensors     │  Dashboard   │  Agents    │
│  Vercel      │  Mason Devs  │              │            │
└──────┬───────┴──────┬───────┴──────┬───────┴─────┬──────┘
       │              │              │             │
       ▼              ▼              ▼             ▼
┌──────────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐
│ n8n Webhook  │ │  MQTT    │ │WebSocket │ │ NATS Client  │
│   Gateway    │ │  Bridge  │ │ (Port    │ │     SDK      │
│  (HTTP→NATS) │ │(Port1883)│ │  9222)   │ │  (Direct)    │
└──────┬───────┘ └────┬─────┘ └────┬─────┘ └──────┬───────┘
       │              │            │              │
       └──────────────┴────────────┴──────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │      NATS JETSTREAM           │
              │   (Central Cluster on DO)     │
              │                               │
              │  Streams: SAAS / AGENT /      │
              │  OMCS / HOME / COMMUNITY      │
              │                               │
              │  Subject: legacy.saas.zoho.>  │
              │  Subject: legacy.home.>       │
              │  Subject: legacy.agent.>      │
              └────────────────┬──────────────┘
                               │
                    ┌──────────┴──────────┐
                    │                     │
                    ▼                     ▼
          ┌─────────────────┐   ┌─────────────────┐
          │ OpenClaw        │   │  Leaf Nodes      │
          │ Gateway +       │   │  Home / Farm /   │
          │ Agent Loop      │   │  Office          │
          └────────┬────────┘   └─────────────────┘
                   │
         ┌─────────┴──────────┐
         │                    │
         ▼                    ▼
  ┌────────────┐      ┌───────────────┐
  │  Supabase  │      │  External     │
  │ PostgreSQL │      │  APIs / Tools │
  └────────────┘      └───────────────┘
```

---

## Leaf Node Architecture (Multi-Location)

```
Tailscale Mesh Network
─────────────────────────────────────────────────────

[Home Location]          [Farm/Homestead]        [Office]
NATS Leaf Node           NATS Leaf Node          NATS Leaf Node
Port 4222                Port 4222               Port 4222
   │                         │                       │
   │ Home Assistant           │ Soil/Weather           │ Mason Devices
   │ Frigate Cameras          │ Sensors                │ Office Agents
   │ Echo Show 5              │ Security Cams          │
   │                         │                       │
   └─────────────────────────┴───────────────────────┘
                             │
                      Tailscale VPN
                             │
                    ┌────────────────┐
                    │ Central NATS   │
                    │ JetStream      │
                    │ DigitalOcean   │
                    └────────────────┘
```

---

## NATS Server Configuration Reference

### Base Config (One Mind OS)

```
# nats-server.conf
host: 0.0.0.0
port: 4222

# HTTP monitoring
http_port: 8222

# JetStream
jetstream {
  store_dir: /data/nats
  max_memory_store: 1GB
  max_file_store: 10GB
}

# MQTT for IoT / Home Assistant
mqtt {
  port: 1883
}

# WebSocket for browser clients
websocket {
  port: 9222
  no_tls: false
  tls {
    cert_file: "/etc/nats/tls/cert.pem"
    key_file: "/etc/nats/tls/key.pem"
  }
}

# Leaf node hub (accepts connections from locations)
leafnodes {
  port: 7422
}

# Logging
log_file: "/var/log/nats/nats.log"
logtime: true
```

---

## Quick Reference — Which Path to Use

| Event Source | Connection Method | Subject Pattern |
|---|---|---|
| Zoho / ClickUp / Gmail | n8n Webhook → NATS | `legacy.saas.{app}.{event}` |
| Home Assistant | MQTT Native | `legacy.home.{device}.{event}` |
| Frigate NVR | MQTT Native | `legacy.home.frigate.{camera}.{event}` |
| Mason Devices | NATS Client SDK | `legacy.office.device.{id}.{event}` |
| Farm Sensors | MQTT via Leaf Node | `legacy.farm.sensor.{type}` |
| Browser/Community Portal | WebSocket | `legacy.community.{entity}.{event}` |
| OpenClaw Agents | NATS Client SDK (Direct) | `legacy.agent.{action}` |
| OMCS Coordinate Updates | NATS Client SDK (Direct) | `legacy.omcs.{pillar}.{coord}` |
| Cross-location events | Leaf Node → Central Cluster | `legacy.{location}.>` |

---

## Priority Setup Checklist

| Priority | Task |
|---|---|
| **Immediate** | Deploy NATS JetStream on DigitalOcean (bare metal, no Docker) |
| **Immediate** | Enable JetStream and create core streams |
| **Immediate** | Enable MQTT on port 1883 for Home Assistant |
| **Week 1** | Configure n8n webhook gateway for Zoho + ClickUp + Gmail |
| **Week 1** | Set up Tailscale-secured leaf node at home location |
| **Week 2** | Connect Home Assistant MQTT → NATS |
| **Week 2** | Configure Frigate NVR MQTT events |
| **Week 2** | Set up WebSocket for community portal |
| **Month 1** | Deploy leaf nodes at farm and office locations |
| **Month 1** | Build lightweight custom webhook gateway (replace n8n for this) |
| **Ongoing** | Monitor streams via NATS HTTP monitoring (port 8222) |

---

*One Mind OS Stack: OpenClaw → NATS JetStream → Supabase → Tailscale Mesh*
*Legacy Platform Holdings LLC — Confidential*
*OMCS Pillars: HP (Holistic Performance) | LE (Legacy Evolution) | GE (Generational Entrepreneurship) | IT (Innovative Technology)*


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
