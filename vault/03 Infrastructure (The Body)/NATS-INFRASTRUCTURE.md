---
domain: 02
name: NATS Infrastructure
status: active
updated: 2026-05-22
owner: zeus
tags: [nats, messaging, infrastructure, agents, event-driven]
---

# NATS Infrastructure — One Mind

Complete documentation of the NATS messaging infrastructure for One Mind.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                    NATS MESSAGE BUS (localhost:4222)                  │
├─────────────────────────────────────────────────────────────────────┤
│  JetStream: 5 Streams (30-day retention)                            │
│  KV Stores: 6 buckets (agent-state, zeus-context, config, etc.)     │
│  WebSocket: ws://localhost:8080                                      │
│  HTTP Monitor: http://localhost:8222                                 │
│  Leafnodes: port 7422 (for VPS cluster)                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐           │
│  │ OpenClaw │  │  Claude   │  │  Hermes  │  │    PI    │           │
│  │  (NATS   │  │   Code   │  │  (NATS   │  │  (NATS   │           │
│  │ Channel) │  │ (Plugin)  │  │ Gateway) │  │ Channel) │           │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘           │
│       │              │              │              │                 │
│       └──────────────┴──────────────┴──────────────┘                │
│              Synadia Agent Protocol v0.3                             │
│       subjects: agents.{verb}.{type}.{owner}.{session}              │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Server Configuration

| Field | Value |
|-------|-------|
| Binary | `/opt/homebrew/bin/nats-server` |
| Version | 2.14.1 |
| Config | `~/.nats/config/onemind.conf` |
| Data | `~/.nats/data/` |
| Log | `~/.nats/nats-server.log` |
| LaunchAgent | `~/Library/LaunchAgents/io.nats.server.plist` |
| Server Name | `onemind-mac` |
| Client Port | 4222 |
| HTTP Monitor | 8222 |
| WebSocket | 8080 |
| Leafnodes | 7422 |
| Memory Store | 1 GB |
| File Store | 10 GB |

### Config File (`~/.nats/config/onemind.conf`)

```conf
server_name: onemind-mac
host: 0.0.0.0
port: 4222
http: 8222

jetstream {
  store_dir: /Users/onemind/.nats/data
  max_mem: 1G
  max_file: 10G
}

leafnodes {
  port: 7422
}

websocket {
  port: 8080
  no_tls: true
}

log_file: /Users/onemind/.nats/nats-server.log
```

### LaunchAgent (auto-start on boot)

```xml
<!-- ~/Library/LaunchAgents/io.nats.server.plist -->
<plist version="1.0">
<dict>
  <key>Label</key><string>io.nats.server</string>
  <key>ProgramArguments</key>
  <array>
    <string>/opt/homebrew/bin/nats-server</string>
    <string>-c</string>
    <string>/Users/onemind/.nats/config/onemind.conf</string>
  </array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
</dict>
</plist>
```

### CLI Context

```bash
nats context ls              # list contexts
nats context select onemind-local  # switch context
```

Context file: `~/.config/nats/context/onemind-local.json`

---

## JetStream Streams

All streams have **30-day retention** and are file-backed.

| Stream | Subjects | Purpose |
|--------|----------|---------|
| `AGENT_EVENTS` | `agent.>` | Agent lifecycle, tool calls, completions |
| `SYSTEM` | `system.>` | System health, alerts, monitoring |
| `SAAS_EVENTS` | `saas.>` | SaaS webhook events (Composio, Todoist, etc.) |
| `COMMS` | `comms.>` | Cross-agent communication |
| `TASKS` | `task.>` | Task creation, status updates, completions |

### Manage Streams

```bash
nats stream ls                    # list all streams
nats stream info AGENT_EVENTS     # details on a stream
nats stream view AGENT_EVENTS     # view recent messages
nats pub agent.legacy.status '{"status":"active"}'  # publish
```

---

## KV Stores

| Bucket | History | Purpose | Example Keys |
|--------|---------|---------|--------------|
| `agent-state` | 5 | Agent current state | `legacy.status`, `hermes.status` |
| `zeus-context` | 10 | Zeus's current context | `location`, `focus_mode`, `status` |
| `config` | 3 | System configuration | `nats.version`, `setup.date` |
| `triggers` | 5 | Event triggers | automation triggers |
| `connections` | 5 | Active connections | service connections |
| `schedules` | 10 | Scheduled tasks | cron-like schedules |

### Manage KV

```bash
nats kv ls                          # list all buckets
nats kv get zeus-context location   # get a value
nats kv put zeus-context status "focused"  # set a value
nats kv history zeus-context location      # view history
nats kv watch zeus-context          # watch for changes (real-time)
```

---

## Installed SDKs

### TypeScript (global npm)

| Package | Version | Purpose |
|---------|---------|---------|
| `@nats-io/transport-node` | 3.4.0 | Core NATS connection (TCP) |
| `@nats-io/nats-core` | latest | Core types + WebSocket |
| `@nats-io/jetstream` | latest | JetStream client |
| `@nats-io/kv` | latest | KV store client |
| `@nats-io/obj` | latest | Object store client |
| `@nats-io/services` | latest | Micro services framework |
| `@synadia-ai/agents` | 0.5.1 | Synadia caller SDK |
| `@synadia-ai/agent-service` | 0.5.2 | Synadia host SDK |
| `@synadia-ai/nats-channel` | latest | OpenClaw NATS channel |
| `@synadia-ai/nats-pi-channel` | latest | PI NATS channel |

### Python (system)

| Package | Version | Purpose |
|---------|---------|---------|
| `nats-py` | 2.14.0 | Core NATS client |
| `nats-core` | 0.2.0 | Modular core |
| `nats-jetstream` | 0.3.0 | Modular JetStream |
| `nats-key-value` | 0.1.0 | Modular KV |
| `synadia-ai-agents` | 0.7.1 | Synadia caller SDK |
| `synadia-ai-agent-service` | 0.4.1 | Synadia host SDK |

---

## Monitoring & Dashboard

### Netdata (Real-time Monitoring Dashboard)

| Field | Value |
|-------|-------|
| Local Dashboard | http://localhost:19999 |
| Cloud Dashboard | https://app.netdata.cloud |
| Install | `brew install netdata` |
| Service | `brew services start netdata` |
| Config | `/opt/homebrew/etc/netdata/` |
| NATS Collector | `/opt/homebrew/etc/netdata/go.d/nats.conf` |
| Version | 2.10.3 |

Features: Per-second metrics, ML anomaly detection, NATS-native collector (JetStream, connections, accounts, HTTP endpoints), 2693+ charts, mobile app access via Netdata Cloud.

**NATS Metrics Collected:**
- JetStream: streams, consumers, storage, memory, API requests/errors
- Per-account: traffic, messages, connections, subscriptions, leaf nodes, slow consumers
- HTTP endpoints: all NATS monitoring API routes
- Process: CPU, memory, threads, uptime

```bash
# Manage service
brew services start netdata
brew services stop netdata
brew services restart netdata

# Access dashboards
open http://localhost:19999              # local dashboard
open https://app.netdata.cloud          # cloud (phone + computer)
```

**Note:** The `macos` plugin is disabled in `/opt/homebrew/etc/netdata/netdata.conf` due to a stack buffer overflow bug on macOS 26.5 Tahoe. Process monitoring still works via apps.plugin.

### nats-top (Terminal)

```bash
nats-top               # real-time TUI (like htop for NATS)
```

### HTTP Monitoring

```bash
curl http://localhost:8222/varz    # server stats
curl http://localhost:8222/jsz     # jetstream stats
curl http://localhost:8222/connz   # connections
curl http://localhost:8222/subsz   # subscriptions
```

---

## Agent Channel Plugins (Synadia Agent Protocol v0.3)

All agents register as NATS micro services and are discoverable via:
```bash
nats req '$SRV.INFO.agents' '' --replies=0 --timeout=2s
nats micro ls
nats sub 'agents.hb.>'   # watch heartbeats
```

### OpenClaw NATS Channel

**Status:** Installed and configured.

| Field | Value |
|-------|-------|
| Package | `@synadia-ai/nats-channel` |
| Installed via | `openclaw plugins install @synadia-ai/nats-channel` |
| Subject | `agents.prompt.oc.zeus.legacy` |
| Heartbeat | `agents.hb.oc.zeus.legacy` (every 5s) |
| Status | `agents.status.oc.zeus.legacy` |
| Config key | `channels.nats` in `openclaw.json` |
| Context | `onemind-local` |

**Configuration in `~/.openclaw/openclaw.json`:**
```json
{
  "channels": {
    "nats": {
      "accounts": {
        "default": {
          "agentName": "legacy",
          "owner": "zeus",
          "context": "onemind-local"
        }
      }
    }
  }
}
```

**Activate:** Restart OpenClaw gateway → `openclaw gateway restart`

**Verify:**
```bash
nats req '$SRV.INFO.agents' '' --replies=0 --timeout=2s
nats sub 'agents.hb.oc.zeus.legacy'
```

**Talk to it:**
```bash
nats req agents.prompt.oc.zeus.legacy "Hello!" \
  --wait-for-empty --reply-timeout 30s --timeout 60s
```

---

### Claude Code NATS Plugin

**Status:** Ready to install (requires Bun).

| Field | Value |
|-------|-------|
| Package | `claude-channel-nats` (Claude plugin marketplace) |
| Subject | `agents.prompt.cc.zeus.<session>` |
| Heartbeat | `agents.hb.cc.zeus.<session>` (every 5s) |
| Prerequisites | Bun (`~/.bun/bin/bun` — installed) |
| Config dir | `~/.claude/channels/nats/` |

**Install:**
```bash
# 1. Add the Synadia marketplace (inside a claude session)
/plugin marketplace add synadia-ai/synadia-agents

# 2. Install the plugin
/plugin install nats-channel@synadia-plugins

# 3. Launch with the channel
claude --dangerously-load-development-channels plugin:nats-channel@synadia-plugins
```

**Configure to use local server:**
```bash
# Inside Claude Code session:
/nats-channel:configure onemind-local
/nats-channel:configure session onemind
```

This writes to `~/.claude/channels/nats/config.json`:
```json
{
  "context": "onemind-local",
  "sessionName": "onemind",
  "permissions": { "mode": "terminal" }
}
```

**Verify from outside:**
```bash
nats req '$SRV.INFO.agents' '' --replies=0 --timeout=2s
nats req agents.prompt.cc.zeus.onemind "Hello Claude" \
  --replies=0 --reply-timeout=30s --timeout=90s
```

---

### Hermes NATS Gateway

**Status:** Not yet installed (requires Hermes v0.14.0+).

| Field | Value |
|-------|-------|
| Repository | `synadia-ai/hermes-nats-gateway` |
| Subject | `agents.prompt.hermes.zeus.default` |
| Heartbeat | `agents.hb.hermes.zeus.default` (every 30s) |
| Plugin dir | `~/.hermes/plugins/nats-platform/` |

**Install:**
```bash
# 1. Install Hermes (if not already)
curl -fsSL https://hermes-agent.nousresearch.com/install.sh | bash -s -- --skip-setup

# 2. Install the NATS gateway plugin
hermes plugins install synadia-ai/hermes-nats-gateway

# 3. Install NATS SDKs into Hermes venv (REQUIRED)
bash ~/.hermes/plugins/nats-platform/scripts/install-sdks.sh

# 4. Configure NATS
hermes setup gateway   # pick "NATS", select onemind-local context

# 5. Or manual config in ~/.hermes/.env:
# NATS_CONTEXT=onemind-local
# HERMES_NATS_OWNER=zeus
# HERMES_NATS_SESSION_NAME=default

# 6. Run
hermes gateway run
```

**Verify:**
```bash
nats micro ls
nats req agents.status.hermes.zeus.default ''
nats sub 'agents.hb.hermes.zeus.default'
```

---

### PI Agent NATS Channel

**Status:** npm package installed globally.

| Field | Value |
|-------|-------|
| Package | `@synadia-ai/nats-pi-channel` |
| Subject | `agents.prompt.pi.zeus.<session>` |

Refer to PI Agent docs for activation.

---

## Synadia Agent Protocol v0.3 — Wire Format

### Subject Pattern
```
agents.{verb}.{type}.{owner}.{session}
```

| Verb | Purpose |
|------|---------|
| `prompt` | Send a prompt to an agent |
| `hb` | Heartbeat (published every 5-30s) |
| `status` | On-demand liveness check |

| Type | Agent |
|------|-------|
| `oc` | OpenClaw |
| `cc` | Claude Code |
| `hermes` | Hermes |
| `pi` | PI Agent |

### Request Format

Plain text:
```bash
nats req agents.prompt.oc.zeus.legacy "Hello!" \
  --wait-for-empty --reply-timeout 30s --timeout 60s
```

JSON envelope:
```json
{"prompt": "Hello!", "attachments": [{"filename": "doc.pdf", "content": "<base64>"}]}
```

### Response Chunks (streamed on reply subject)

```json
{"type": "status", "data": "ack"}         // immediate acknowledgment
{"type": "response", "data": "Hello..."}  // content chunk
{"type": "response", "data": " world"}    // more content
{"type": "query", "data": {...}}          // permission request (mid-stream)
                                           // empty body = end of stream (terminator)
```

### Discovery

```bash
# Find all agents
nats req '$SRV.INFO.agents' '' --replies=0 --timeout=2s

# Ping all agents
nats req '$SRV.PING.agents' '' --replies=0 --timeout=2s

# List micro services
nats micro ls

# Watch all heartbeats
nats sub 'agents.hb.>'
```

---

## Composio → NATS Webhook Bridge

### Architecture (Verified Path)

```
Composio Action Trigger
    → Composio Webhook (POST to callback URL)
        → Bridge Service (lightweight HTTP → NATS publisher)
            → NATS subject: saas.{app}.{event}
                → JetStream SAAS_EVENTS stream (30-day retention)
                    → Consumers (agents, automations)
```

### How Composio Webhooks Work

Composio supports **trigger subscriptions** — when an event occurs in a connected app (new Todoist task, GitHub push, Strava activity), Composio sends a POST to your callback URL.

**Connected accounts (verified):**
- Todoist, GitHub, Strava, Dropbox, Instagram, Facebook, YouTube, Gemini

### Composio CLI Commands

```bash
# List available triggers for an app
composio triggers --app todoist

# Subscribe to a trigger
composio triggers subscribe TODOIST_NEW_TASK --callback https://your-bridge/webhook

# List active subscriptions
composio triggers list

# Check connected accounts
composio connections list
```

### Bridge Service Design

A minimal Node.js/Python service that:
1. Receives POST from Composio webhook
2. Validates the payload
3. Publishes to NATS on `saas.{app}.{event_type}`

**Example (Node.js):**
```javascript
import { connect } from "@nats-io/transport-node";
import { createServer } from "http";

const nc = await connect({ servers: "nats://127.0.0.1:4222" });

createServer((req, res) => {
  let body = "";
  req.on("data", (chunk) => body += chunk);
  req.on("end", () => {
    const event = JSON.parse(body);
    const subject = `saas.${event.app}.${event.trigger}`;
    nc.publish(subject, JSON.stringify(event));
    res.writeHead(200);
    res.end("ok");
  });
}).listen(9222);
```

**Subject naming:**
- `saas.todoist.new_task` — new Todoist task
- `saas.github.push` — GitHub push event
- `saas.strava.activity_created` — new Strava activity

### Expose via Tailscale Funnel (no public IP needed)

```bash
tailscale funnel --bg 9222
# → https://onemindos-cloud-1.tail717ba0.ts.net → localhost:9222
```

---

## Useful Commands Cheatsheet

```bash
# Server
launchctl kickstart -k gui/$(id -u)/io.nats.server  # restart server
nats server info                                      # server details

# Publish / Subscribe
nats pub agent.legacy.heartbeat '{"ts":"now"}'
nats sub 'agent.>'                                    # watch all agent subjects
nats sub 'saas.>'                                     # watch SaaS events

# Streams
nats stream ls
nats stream info AGENT_EVENTS
nats stream view AGENT_EVENTS --last 10
nats stream purge AGENT_EVENTS                        # ⚠️ deletes all messages

# KV
nats kv ls
nats kv get agent-state legacy.status
nats kv put agent-state legacy.status '{"state":"active","ts":"2026-05-22"}'
nats kv watch agent-state                             # real-time watch

# Monitoring
nats-top                                              # terminal dashboard
open http://localhost:19999                            # Netdata dashboard (local)
open https://app.netdata.cloud                        # Netdata Cloud (phone/computer)
curl -s http://localhost:8222/jsz | python3 -m json.tool  # JetStream stats

# Agent Discovery (Synadia Protocol)
nats req '$SRV.INFO.agents' '' --replies=0 --timeout=2s
nats micro ls
nats sub 'agents.hb.>'
```

---

## SDK Usage Examples

### TypeScript — Discover & Prompt an Agent

```typescript
import { connect } from "@nats-io/transport-node";
import { Agents } from "@synadia-ai/agents";

const nc = await connect({ servers: "nats://localhost:4222" });
const agents = new Agents({ nc });

// Discover all agents on this NATS
const found = await agents.discover();
console.log(`Found ${found.length} agents`);

// Prompt the first one
const [agent] = found;
for await (const msg of await agent!.prompt("what can you do?")) {
  if (msg.type === "response") process.stdout.write(msg.text);
}

await agents.close();
await nc.close();
```

### Python — Discover & Prompt

```python
import asyncio
import nats
from synadia_ai.agents import Agents, ResponseChunk

async def main():
    nc = await nats.connect("nats://127.0.0.1:4222")
    agents = Agents(nc=nc)

    found = await agents.discover()
    print(f"Found {len(found)} agents")

    if found:
        async for msg in found[0].prompt("hello"):
            if isinstance(msg, ResponseChunk):
                print(msg.text, end="", flush=True)

    await agents.close()
    await nc.close()

asyncio.run(main())
```

### TypeScript — Host a Custom Agent

```typescript
import { connect } from "@nats-io/transport-node";
import { AgentService } from "@synadia-ai/agent-service";

const nc = await connect({ servers: "nats://localhost:4222" });

const service = new AgentService({
  nc,
  agent: "custom",
  owner: "zeus",
  name: "router",
  description: "OneMind event router",
});

service.onPrompt(async (envelope, response) => {
  // Process the prompt and respond
  await response.send(`Received: ${envelope.prompt}`);
});

await service.start();
console.log(`Listening on ${service.subject.prompt}`);
```

---

## Leafnode Cluster (Active — May 2026)

Multi-machine NATS fabric connected via Tailscale leafnodes to the Mac hub.

```
┌─────────────────────────────────────────────────────────────┐
│              Mac Hub (omos-mbp / 100.98.178.23)              │
│   nats-server v2.14.1 — ports 4222/7422/8080/8222           │
│   Leafnode listener: 0.0.0.0:7422                           │
├─────────────────────────────────────────────────────────────┤
│         ↕ leafnode                    ↕ leafnode (planned)   │
│  ┌─────────────────────┐      ┌──────────────────────┐     │
│  │ Dell (omos-dell)    │      │ HA (omos-ha)         │     │
│  │ 100.66.84.71        │      │ 100.64.25.12         │     │
│  │ nats-server v2.14.1 │      │ (not yet deployed)   │     │
│  │ systemd: enabled    │      │                      │     │
│  └─────────────────────┘      └──────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### Dell Leafnode (`omos-dell`)

| Field | Value |
|-------|-------|
| Binary | `/usr/local/bin/nats-server` v2.14.1 |
| Config | `/etc/nats/nats-server.conf` |
| Service | `nats-server.service` (systemd, enabled) |
| Server name | `onemind-dell` |
| Local port | `127.0.0.1:4222` |
| Connects to | `nats-leaf://100.98.178.23:7422` |
| Log | `/var/log/nats/nats-server.log` |
| Deployed | 2026-05-22 |

---

## HA↔NATS Bridge (Active — May 2026)

Connects Home Assistant WebSocket API to the NATS fabric. Runs as a macOS LaunchAgent on the Mac.

| Field | Value |
|-------|-------|
| Code | `~/.nats/services/ha-bridge/bridge.js` |
| LaunchAgent | `io.onemind.ha-bridge` |
| HA Host | `http://192.168.1.180:8123` |
| HA Version | 2026.4.4 |
| Entities cached | 183 |
| Token | Long-lived access token (in `.env`) |

### NATS Subjects

| Subject | Direction | Purpose |
|---------|-----------|---------|
| `cmd.ha.service.call` | NATS → HA | Call any HA service (lights, switches, scripts) |
| `cmd.ha.state.get` | NATS → HA | Get current state of an entity |
| `cmd.ha.states.list` | NATS → HA | List/filter entities by domain |
| `events.ha.state` | HA → NATS | All state change events |
| `events.ha.{domain}` | HA → NATS | Domain-specific events (e.g. events.ha.light) |

### Usage Examples

```bash
# Get weather
nats req cmd.ha.state.get '{"entity_id":"weather.forecast_home"}'

# List all sensors
nats req cmd.ha.states.list '{"domain":"sensor"}'

# Call a service
nats req cmd.ha.service.call '{"domain":"switch","service":"turn_on","entity_id":"switch.office_fan"}'

# Subscribe to real-time events
nats sub "events.ha.>"
```

---

## Universal NATS MCP Server (Active — May 2026)

Gives every AI agent (Hermes, Claude Code, Copilot, Gemini) full NATS access via MCP protocol.

| Field | Value |
|-------|-------|
| Code | `~/.nats/mcp-server/server.js` |
| Package | `@onemind/nats-mcp-server` v1.0.0 |
| Configured in | `~/.hermes/config.yaml`, `.claude/mcp.json` |

### Tools

| Tool | Purpose |
|------|---------|
| `nats_request` | req/reply on any subject (cmd.ha.*, cmd.system.*, etc.) |
| `nats_publish` | fire-and-forget publish (events, notifications, triggers) |
| `nats_kv_get` | read key from KV bucket (zeus-context, agent-state) |
| `nats_kv_put` | write key to KV bucket |
| `nats_kv_keys` | list all keys in a bucket (discovery) |
| `nats_stream_read` | read last N messages from a JetStream stream (history) |

### Architecture

```
Agent (Hermes/Claude/Copilot) 
  → MCP tool call (nats_request, nats_kv_get, etc.)
    → NATS MCP Server (stdio)
      → NATS bus (localhost:4222)
        → HA Bridge (cmd.ha.*) 
        → Dell services (cmd.dell.* — future)
        → Any future microservice
```

One universal tool layer — new NATS microservices are immediately accessible without tool changes.

---

## Future Plans

- [ ] HA NATS leafnode (NATS on Home Assistant Alpine box, connecting to Mac:7422)
- [ ] VPS Atlas leafnode (3.234.227.222 — Tailscale + NATS)
- [x] ~~VPS leafnode cluster~~ → Dell leafnode active (2026-05-22)
- [ ] Composio webhook bridge service (running on VPS or Mac)
- [ ] n8n NATS node (publish/subscribe from n8n workflows)
- [x] ~~Agent-to-agent communication over NATS~~ → nats-fabric plugin active
- [ ] JetStream consumers for durable event processing
- [ ] NATS Auth (accounts + NKey) for multi-tenant security

---

## Core Concepts — How NATS Works

NATS is built around three primitives that map directly to how OneMind OS thinks:

| NATS Term | OneMind Analogy | Role |
|-----------|-----------------|------|
| **Publisher** | Sensor | Emits a signal — detects something, fires an event |
| **Subject** | Pathway | The named channel the signal travels through (e.g. `cmd.ha.light.on`) |
| **Subscriber** | Reaction | Listens on a subject and acts when a message arrives |

### Flow
```
Publisher (Sensor)
  → fires message on Subject (Pathway)
    → Subscriber(s) (Reaction) receive and act
```

### In OneMind Context
- **Publisher** = IoT device, agent, sensor, Home Assistant, script — anything that detects or generates an event
- **Subject** = the routing address (e.g. `saas.execute`, `cmd.ha.*`, `agent.legacy.inbox`) — hierarchical, wildcard-capable
- **Subscriber** = NATS saas worker, HA bridge, agent listener, automation — whatever reacts to that signal

### Subject Wildcards
- `*` — single token wildcard: `cmd.ha.*` matches `cmd.ha.light` but not `cmd.ha.light.on`
- `>` — multi-token wildcard: `cmd.ha.>` matches everything under `cmd.ha.`

### Why This Matters for OneMind
NATS is the nervous system. Every action in the OS follows this pattern — a sensor fires, travels a pathway, triggers a reaction. No direct API calls, no polling. Pure event-driven.

---

## Related

- [[03 Infrastructure (The Body)/README|Infrastructure README]]
- [[01 Intelligence (The Brain)/README|Intelligence README]]
- NATS docs: https://docs.nats.io
- Netdata: https://www.netdata.cloud (monitoring dashboard)
- Synadia Agents: https://github.com/synadia-ai/synadia-agents
- Agent Protocol Spec: https://github.com/synadia-ai/synadia-agent-sdk-docs
