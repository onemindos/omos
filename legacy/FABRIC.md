# NATS Fabric — The Nervous System

> Last updated: 2026-08-12
> Author: Legacy (self-authored)

The NATS bus is the nervous system of OneMind. Every actor — human, AI, sensor, device — communicates through it.

## Cluster

| Property | Value |
|----------|-------|
| Deployment | 3-node StatefulSet (`nats-0`, `nats-1`, `nats-2`) |
| Namespace | `nats` |
| Cluster | do-nyc3-onemind-hub |
| Domain | `onemind` (JetStream) |
| Client port | 4222 (internal) |
| Leaf port | 7422 (TLS, LoadBalancer) |
| WebSocket | 8080 |
| MQTT | 1883 |
| Monitoring | 8222 |
| LoadBalancer IP | 100.117.169.54 (omos-nats on tailnet) |

## Accounts

| Account | Purpose |
|---------|---------|
| `ONEMIND` | Primary — all agents, services, leaves |
| `SYS` | System monitoring |

**Note:** No `OMOS` account on the hub. OMOS is Mac-local only.

## Leaf Nodes

| Leaf | Location | Credentials | Account |
|------|----------|-------------|---------|
| hermes-legacy-leaf | DOKS sidecar | leaf_tak | ONEMIND |
| omos-ha-leaf | HA Yellow | leaf_tak | ONEMIND |
| omos-mbp (local) | Zeus's Mac | local | OMOS (local) |

### Leaf Configuration Template (Edge devices)
```
port: 4222
leafnodes {
  remotes [
    {
      url: "tls://leaf_tak:<password>@100.117.169.54:7422"
      tls { insecure: true }
    }
  ]
}
```

**Pitfalls:**
- Use `100.117.169.54:7422` (LB IP) — NOT node IPs
- `tls { insecure: true }` goes INSIDE the remote block
- Node IPs cause TLS handshake failure

## Subject Schema

### Agent Communication
```
agents.prompt.hermes.zeus.legacy        ← Mother (me)
agents.prompt.hermes.zeus.legacy-edge   ← Daughter (edge)
agents.prompt.hermes.zeus.forge         ← Forge (coding)
agents.prompt.hermes.zeus.<servant>     ← Future servants
```

### TAK / Situational Awareness
```
tak.cot.out            → CoT from TAK Server (tak-platform plugin)
tak.cot.inject         → Inject CoT into TAK Server (TLS :8089)
tak.cot.>              → Full CoT event stream
```

**TAK platform plugin** (`/opt/data/plugins/platforms/tak-platform/`): Hermes-native bridge.
TAK GeoChat → NATS `tak.cot.out` → Legacy reasons → reply CoT injected directly
to TAK Server via mTLS TCP :8089.

### Voice (Reserved — ATAK-native target)
```
voice.transcript       → STT output (TranscriptMsg JSON)
voice.say              → TTS response
voice.wake             → Wake-word detections (optional mode only)
voice.partial          → Streaming partials (early-trigger)
cmd.<domain>.<action>  → Structured commands (req/reply)
cmd.<domain>.abort     → Safety reflex (no LLM)
evt.<domain>.<event>   → Agent results/telemetry
req.sa.location        → Operator position resolution
req.fleet.drones.available → Fan-out asset discovery
cot.in                 → Render on TAK map
```

omos-voice (Mumble/LiveKit) archived 2026-08-12. Target: ATAK-native PTT plugin
pipe audio directly to STT → NATS. Contract preserved in SO Wiki.

### Home Assistant (via ha-bridge)
```
cmd.ha.service.call    → Call HA service
cmd.ha.state.get       → Get entity state
cmd.ha.states.list     → List all entities (185+)
events.ha.>            → HA state change events
```

### RADS (Homestead Digital Twin)
```
rad.greenwich.land.field.<n>.soil       → Soil readings
rad.greenwich.land.field.<n>.irrigate   → Irrigation commands
rad.greenwich.land.decisions            → Agent decision log
rad.greenwich.robots.<id>.telemetry     → MAVLink telemetry
rad.greenwich.robots.<id>.cmd           → Drone commands
rad.greenwich.legacy.override           → Legacy override
```

### System
```
legacy.health.alive    → Heartbeat
legacy.*               → Legacy's own namespace
perception.*           → Vision/sensor context feeds
```

## JetStream Streams

| Stream | Subject Filter | Purpose | TTL |
|--------|---------------|---------|-----|
| RADS_LAND | rad.*.land.> | Soil/agriculture | 10 years |
| RADS_ROBOTS | rad.*.robots.> | Drone telemetry | 10 years |
| AGENT_EVENTS | agents.event.> | Agent lifecycle | 30 days |
| AGENT_PROMPTS | agents.prompt.> | Prompts to agents | 7 days |

## KV Buckets

| Bucket | Purpose |
|--------|---------|
| `connections` | Agent connections |
| `schedules` | Scheduled tasks |
| `triggers` | Event triggers |
| `config` | System config |
| `agent-state` | Agent state tracking |
| `zeus-context` | Zeus context |
| `rules` | Automation rules |

## NATS CLI Usage (from Legacy pod)

```bash
nats pub legacy.health.alive "$(date -u +%s)"
nats req cmd.ha.state.get '{"entity_id":"switch.plug_mini_ef5a"}' --timeout 5s
nats sub "events.ha.>"
nats kv get zeus-context <key>
nats stream info RADS_LAND
```

---

*The bus owns truth. All state flows through NATS. Never bypass the fabric.*
