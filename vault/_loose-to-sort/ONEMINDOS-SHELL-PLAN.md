# OneMindOS Shell — Architecture Plan

## Overview

Flutter native app (macOS, iOS, Android, Web) as the primary UI shell for OneMindOS, backed by SwarmClaw running headless as the API engine and integration layer.

```
┌─────────────────────────────────────────────────────┐
│  Flutter App (OneMindOS Shell)                      │
│  Platforms: macOS • iOS • Android • Web             │
│  Features: Voice • Custom UI • Native Performance   │
└──────────────────────┬──────────────────────────────┘
                       │ REST API (localhost:3456/api)
┌──────────────────────▼──────────────────────────────┐
│  SwarmClaw (headless LaunchAgent)                   │
│  Engine: OpenClaw RPC • Providers • Connectors      │
│  Data: Sessions • Memory • Fleet Topology • Agents  │
│  DB: SQLite (swarmclaw.db + memory.db)              │
└─────────────────────────────────────────────────────┘
```

---

## Why This Architecture

- Flutter gives native voice, mobile, and cross-platform without WebView compromises
- SwarmClaw already implements the hard integration work (OpenClaw WebSocket RPC, Ed25519 auth, 23 providers, 13 connectors)
- Decoupled: Flutter owns UX, SwarmClaw owns integrations — either can be replaced independently
- No vendor lock-in to SwarmClaw's frontend opinions
- Mobile-ready from day one

---

## SwarmClaw Backend (Already Running)

- Deployed as: macOS LaunchAgent (`com.onemind.swarmclaw`)
- Binary: `~/.nvm/versions/node/v24.15.0/bin/swarmclaw`
- Plist: `~/Library/LaunchAgents/com.onemind.swarmclaw.plist`
- Web UI (unused, just the API): `http://localhost:3456`
- Data: `~/.swarmclaw/data/`
- Logs: `~/Library/Logs/swarmclaw.log`, `~/Library/Logs/swarmclaw.error.log`
- Auto-starts on login, restarts on crash (KeepAlive: true)

---

## SwarmClaw API Endpoints (Flutter consumes these)

### Fleet / Infrastructure

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/gateways/fleet` | Fleet topology — all connected OpenClaw gateways, node counts, session counts, health |
| GET | `/api/gateways` | List gateway profiles |
| POST | `/api/gateways` | Add a new gateway profile (name + WebSocket URL) |
| DELETE | `/api/gateways/:id` | Remove a gateway profile |

### Agents

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/agents` | All registered agents (Legacy, Zenith, Unity, Eternal, Sovereignty, Sentinel) |
| GET | `/api/agents/:id` | Single agent details |
| POST | `/api/agents` | Create/register new agent |
| PATCH | `/api/agents/:id` | Update agent config |
| DELETE | `/api/agents/:id` | Remove agent |

### Sessions (Central Nervous System)

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/sessions` | All sessions across all agents and connectors |
| GET | `/api/sessions/:id` | Session detail with messages |
| POST | `/api/sessions` | Create new session |
| POST | `/api/sessions/:id/message` | Send message to agent in session |
| DELETE | `/api/sessions/:id` | End/archive session |

### Memory

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/memory` | Agent memories (semantic + lexical search) |
| GET | `/api/memory?agentId=X` | Memories for specific agent |
| POST | `/api/memory` | Create memory entry |
| DELETE | `/api/memory/:id` | Remove memory |

### Activity Feed

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/activity` | Activity log (entity-based, filterable by type/action/actor) |

### Connectors

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/connectors` | All connector configs and health status |
| POST | `/api/connectors` | Create connector |
| PATCH | `/api/connectors/:id` | Update connector config |
| POST | `/api/connectors/:id/start` | Start connector |
| POST | `/api/connectors/:id/stop` | Stop connector |

### Tasks / Missions

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/tasks` | Task board |
| POST | `/api/tasks` | Create task |
| PATCH | `/api/tasks/:id` | Update task |
| GET | `/api/missions` | Autonomous missions |
| POST | `/api/missions` | Create mission |

### Schedules

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/schedules` | Scheduled jobs |
| POST | `/api/schedules` | Create schedule |
| PATCH | `/api/schedules/:id` | Update schedule |
| DELETE | `/api/schedules/:id` | Remove schedule |

### Knowledge

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/knowledge` | Knowledge sources (files, URLs, manual) |
| POST | `/api/knowledge` | Add knowledge source |
| DELETE | `/api/knowledge/:id` | Remove source |

### Skills

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/skills` | Registered skills |
| POST | `/api/skills` | Add skill |
| DELETE | `/api/skills/:id` | Remove skill |

### MCP Servers

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/mcp-servers` | MCP server configs |
| POST | `/api/mcp-servers` | Add MCP server |
| DELETE | `/api/mcp-servers/:id` | Remove MCP server |

### External Agents

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/external-agents` | Registered external agent runtimes |
| POST | `/api/external-agents` | Register external agent |
| POST | `/api/external-agents/:id/heartbeat` | Heartbeat from external agent |

### Webhooks

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/webhooks` | Webhook configs |
| POST | `/api/webhooks` | Create webhook |
| DELETE | `/api/webhooks/:id` | Remove webhook |

### Projects

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/projects` | Workspace projects |
| POST | `/api/projects` | Create project |

### Chatrooms

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/chatrooms` | Multi-agent group conversations |
| POST | `/api/chatrooms` | Create chatroom |
| POST | `/api/chatrooms/:id/message` | Send message to chatroom |

### Runs

| Method | Endpoint | Returns |
|--------|----------|---------|
| GET | `/api/runs` | Unified execution records (missions + tasks) |

---

## Flutter App — Feature Map

### Core Views

| View | Data Source | Features |
|------|------------|----------|
| Dashboard | `/api/gateways/fleet` + `/api/activity` | Fleet health, live feed, agent status cards |
| Agents | `/api/agents` | Agent profiles (Legacy, Zenith, Unity, Eternal, Sovereignty, Sentinel), config, assign connectors |
| Sessions | `/api/sessions` | Central nervous system — ALL sessions from all agents, all sources |
| Chat | `/api/sessions/:id/message` | Direct conversation with any agent, voice input/output |
| Fleet | `/api/gateways/fleet` | Multi-machine topology, gateway health, drain/cordon controls |
| Memory | `/api/memory` | Agent memories, search, create, link |
| Tasks | `/api/tasks` + `/api/missions` | Task board, mission tracking |
| Connectors | `/api/connectors` | Connector status, start/stop, health monitoring |
| Knowledge | `/api/knowledge` | Knowledge base management |
| Settings | Various | Gateway profiles, connector configs, themes |

### Native Features (Flutter-only, not in SwarmClaw web UI)

| Feature | Implementation |
|---------|---------------|
| Voice Input | `speech_to_text` package, platform Speech APIs |
| Voice Output | `flutter_tts` or macOS `NSSpeechSynthesizer` via platform channel |
| Push Notifications | Firebase Cloud Messaging (mobile), macOS UserNotifications |
| Biometric Auth | `local_auth` package (Face ID, Touch ID) |
| Widgets | iOS WidgetKit, macOS Today widget — agent status at a glance |
| Shortcuts | Siri Shortcuts integration for quick commands |
| Background Sync | Background fetch for session polling on mobile |
| Offline Cache | Local SQLite cache of recent sessions/messages |

---

## Session Routing — Making Everything Appear

### Agents Started Through SwarmClaw (automatic)
Sessions initiated via SwarmClaw's API are captured natively — no extra work.

### External OpenClaw Sessions
OpenClaw connector (built-in) — configure one per gateway with WebSocket URL. Polls `chat.history` every 2.5s + real-time WebSocket events.

### External Claude Code Sessions
Option A: FileQueue connector watching `~/.claude/` for JSONL session files
Option B: Custom extension that parses Claude Code sessions and inserts into SwarmClaw

### External Hermes/Letta Sessions
Custom connector polling Letta REST API at `localhost:8283` → `GET /v1/agents/{id}/messages`

### External Codex Sessions
FileQueue connector watching `~/.codex-sessions/` for thread output

### Any Other CLI Agent
Register via External Agents API + FileQueue ingress pattern

---

## Database (SwarmClaw-managed)

| File | Purpose |
|------|---------|
| `~/.swarmclaw/data/swarmclaw.db` | Main app DB (56 tables, JSON-blob pattern) |
| `~/.swarmclaw/data/memory.db` | Agent memory with FTS5 full-text search |

Schema: `(id TEXT PRIMARY KEY, data TEXT)` — full objects serialized as JSON.
Performance: WAL mode, 64MB cache, 256MB mmap, 15s busy timeout.

---

## Future: Forking SwarmClaw → OneMindOS Engine

When ready to fully own the backend:

| Change | What |
|--------|------|
| Package name | `@onemind/onemindos-engine` |
| Binary | `onemindos-engine` |
| Data dir | `~/.onemind/data/` (or merge with existing vault) |
| Env vars | `ONEMIND_HOME` |
| Branding | Remove SwarmClaw references, update `public/branding/` |
| LaunchAgent | `com.onemind.engine` |
| Theme | Set `--neutral-tint` to OneMindOS brand color |

---

## Phased Rollout

### Phase 1 — Foundation (Current)
- [x] SwarmClaw installed as headless LaunchAgent
- [x] Running at localhost:3456
- [ ] Add OpenClaw gateway profiles (local + Tailscale remotes)
- [ ] Configure OpenClaw connector for session sync
- [ ] Register 6 named agents

### Phase 2 — Flutter Shell (Build)
- [ ] Scaffold Flutter project (macOS first, then iOS/Android)
- [ ] SwarmClaw API client (Dart, generated from endpoint list above)
- [ ] Dashboard view (fleet health + activity feed)
- [ ] Agents view (6 named agents)
- [ ] Sessions view (central nervous system)
- [ ] Chat view with voice input

### Phase 3 — Full Integration
- [ ] Session-watcher extension (Claude Code + Letta + Codex external sessions)
- [ ] Voice output (TTS for agent responses)
- [ ] Push notifications (mobile)
- [ ] Custom connectors for project-specific integrations
- [ ] Biometric auth

### Phase 4 — Own the Stack
- [ ] Fork SwarmClaw → OneMindOS Engine
- [ ] Rebrand and customize data model
- [ ] Merge data directory with existing `/onemind/` vault
- [ ] Custom API extensions beyond SwarmClaw's scope
- [ ] Deploy engine on VPS for remote access (beyond Tailscale)
