# PLATFORM-EXTRACTION.md — What We Took From Five Platforms

> **Authority:** Architecture reference for features extracted from Mission Control,
> SwarmClaw, Paperclip, Agent of Empires, and Multica into onemindos-core.
> Maintained by: Legacy (architect). Last updated: 2026-06-13.

---

## Why This Document Exists

Five platforms were evaluated as potential additions to the OneMind stack:
**Mission Control**, **SwarmClaw**, **Paperclip**, **Agent of Empires (AoE)**, and **Multica**.

The goal was never to run all five. The goal was to understand what ideas each one solved
well, then take those ideas and build them the OneMind way — natively on NATS, using the
subject grammar and envelope schema already defined in `NATS-DOCTRINE.md`, feeding
ClickHouse automatically, and visible to the Flutter app through the same bus that powers
everything else.

Think of each platform as a *prototype* — someone else spent months learning what works.
We read the code, extracted the lesson, and built it our way.

---

## The Two Platforms We Kept

Two platforms were kept running as-is because they each fill a role that would take
significant effort to replace:

| Platform | Role | Why keep it |
|----------|------|-------------|
| **Mission Control** | Fleet management dashboard | GitHub Issues sync, 120+ REST API, already reads `~/.hermes/cron`, reagraph fleet visualization, token cost ledger |
| **Agent of Empires** | Coding session manager | Rust binary (zero overhead), tmux-native sessions, git worktrees, Docker sandboxing, REST API on `:7777` that Hermes can call, ACP structured events |

Everything else (SwarmClaw, Paperclip, Multica) is stripped. The value from each
was extracted into onemindos-core services. Details below.

---

## Platform 1: Mission Control

**Tech stack:** TypeScript/Node.js, React, reagraph, SSE event bus.
**What it does:** Visualizes your agent fleet. Reads `~/.hermes/` directly. Bidirectional
GitHub Issues sync. Token cost tracking per agent. 120+ REST endpoints.

### What we extracted

#### 1.1 — GitHub ↔ NATS Sync Pattern

Mission Control's `github-sync-engine.ts` does bidirectional sync between GitHub Issues
and an internal task store. We adapt this as a NATS service: `svc.github.sync`.

**The idea in plain English:**
When a GitHub issue is created, a NATS event fires. When a task on NATS changes status,
GitHub is updated. The two systems stay in sync without either one being the "master."

**NATS subjects:**

```
evt.sys.github.issue.created       → new issue opened on GitHub
evt.sys.github.issue.closed        → issue closed
evt.sys.github.issue.updated       → fields changed
cmd.sys.github.issue.create        → command: create an issue on GitHub
cmd.sys.github.issue.close         → command: close an issue on GitHub
```

**How to apply:** This pairs with the Task Engine (`om-task`). When a task is created in
NATS KV (`tasks` bucket), the github-bridge optionally mirrors it to GitHub if the task
has `labels: ["github"]` set. Incoming GitHub webhooks publish to `evt.sys.github.>` and
`om-task` picks them up.

**Status:** Pattern documented here. Service `services/github-bridge/` not yet built.

---

#### 1.2 — Token Cost Ledger Pattern

Mission Control tracks API token cost per agent over time. We capture this on the NATS
bus so it flows into ClickHouse automatically.

**The idea:** Hermes publishes cost snapshots already via hermes-bridge. The subject is
`evt.hermes.{profile}.cost`. The ClickHouse `onemind.events` table captures it. Build a
Grafana or Flutter panel that queries `SELECT profile, SUM(tokens), SUM(cost_usd) FROM
onemind.events WHERE type = 'evt.hermes.cost' GROUP BY profile`.

**What we already have:** hermes-bridge publishes `evt.hermes.health` every 10s with full
gateway status JSON including active sessions. Cost events can be added to bridge.js when
Hermes exposes a cost API endpoint.

**Status:** hermes-bridge is live. Cost subject is reserved. No action needed today.

---

#### 1.3 — Fleet Visualization Subject Set

Mission Control's reagraph shows nodes, edges, and agent relationships. To power a similar
view in the Flutter app or a future cockpit, we need a canonical "fleet state" subject.

**Subjects reserved:**

```
evt.sys.fleet.snapshot             → periodic full fleet state (all nodes, agents, statuses)
evt.sys.fleet.node.joined          → new node came online
evt.sys.fleet.node.left            → node disconnected
```

**How to apply:** The agent-registry service (see SwarmClaw section below) publishes
`evt.sys.fleet.snapshot` every 30s based on the `agent-registry` KV bucket. The Flutter
app subscribes and renders the fleet view without needing a REST call.

**Status:** Reserved. agent-registry service is the prerequisite.

---

## Platform 2: SwarmClaw

**Tech stack:** TypeScript, JSON-RPC over HTTP, CLI process scanner.
**What it does:** Discovers every running coding agent on the machine by scanning ports and
process names. Exposes them via a JSON-RPC API: `discoverAgents`, `executeTask`,
`getStatus`, `cancelTask`. Any agent can call any other agent through SwarmClaw.

### What we extracted

#### 2.1 — Agent Auto-Discovery (Most Important Feature)

This is the feature that makes a fleet a fleet, not just a list of processes. SwarmClaw
scans running processes for names like `claude`, `codex`, `opencode`, `aider`, `hermes`
and known ports. It builds a live registry of who is reachable and what they can do.

**The idea in plain English:**
When a new coding tool starts, the system automatically knows it's there. No manual config.
No "add this agent to the list." It discovers itself. Every other service can then ask
"who's available right now?" and get a live answer.

**NATS subjects we define:**

```
sys.agent.register                 → agent announces itself on startup
sys.agent.heartbeat                → alive ping every 30s
sys.agent.deregister               → graceful shutdown notification
sys.agent.query                    → any service asks: who's alive?
evt.sys.agent.discovered           → broadcast: new agent joined the fleet
evt.sys.agent.lost                 → broadcast: agent stopped responding
```

**NATS KV bucket: `agent-registry`**

Each key is an agent ID. Value is a JSON record:

```json
{
  "id": "claude-code.coder",
  "type": "coding",
  "name": "Claude Code — coder profile",
  "port": 0,
  "pid": 48291,
  "host": "onemind-mac",
  "caps": ["code.edit", "code.run", "code.test"],
  "status": "idle",
  "last_seen": "2026-06-13T10:00:00.000Z"
}
```

**How it works (three layers):**

1. **Self-registration:** Each harness (Hermes profile, Claude Code session, AoE session)
   publishes `sys.agent.register` on startup with its identity JSON. The registry service
   writes it to KV.

2. **Heartbeat:** Every 30s, each harness publishes `sys.agent.heartbeat` with its ID.
   The registry service updates `last_seen`. Agents silent for >2 minutes are marked lost.

3. **Bootstrap scan:** At system start, the registry service runs SwarmClaw's port + process
   scan logic to discover agents that were already running before the registry started.
   This seeds the KV bucket with any pre-existing agents.

**Where this is already partially built:**
`services/task-engine/om_task.py` already subscribes to `agents.hb.>` (Synadia Agent
Protocol heartbeats) and writes to the `agent-registry` KV bucket. That's the heartbeat
layer for Hermes agents. The bootstrap scan and `sys.agent.register` self-registration are
the remaining pieces.

**Status:** Heartbeat aggregation live in `om_task.py`. Bootstrap scan + registration
protocol not yet built. Service stub: `services/agent-registry/`.

---

#### 2.2 — A2A Task Dispatch Protocol (JSON-RPC → NATS)

SwarmClaw's `executeTask` / `getStatus` / `cancelTask` is agent-to-agent task delegation.
One agent says "do this" to another agent. We adapt this as NATS request/reply.

**The idea in plain English:**
Hermes's ops profile needs to delegate a coding task to Claude Code. Instead of a
hardcoded REST call to a hardcoded port, it publishes a request to NATS. The registry
knows which agent is available. The right agent picks it up and replies when done.

**Subjects:**

```
cmd.sys.agent.{agent_id}.execute   → dispatch a task to a specific agent (req/reply)
cmd.sys.agent.execute              → dispatch to any available agent of a given type
evt.sys.agent.{agent_id}.task_done → agent signals task completion
evt.sys.agent.{agent_id}.task_fail → agent signals task failure
```

**Request payload (cmd.sys.agent.execute):**

```json
{
  "task_id": "uuid",
  "type": "code.edit",
  "input": "Refactor auth.py — extract JWT validation to a standalone function",
  "context": { "repo": "/Users/onemind/onemindos-core", "branch": "feature/auth" },
  "timeout_s": 300,
  "requester": "agent.hermes.ops"
}
```

**How to apply:** This is the mechanism by which Hermes profiles (ops, ceo, etc.) can
spawn work in Claude Code or other coding agents without knowing their port or PID.
The registry tells you who's alive; the dispatch protocol sends the work.

**Status:** Protocol defined here. No service built yet. Prerequisite: agent-registry.

---

## Platform 3: Paperclip

**Tech stack:** TypeScript, NATS-native, clean subject router.
**What it does:** A NATS connection pool manager with a clean subject router. Every
message goes through a typed handler. No REST, no HTTP — pure NATS.

### What we extracted

#### 3.1 — NATS Transport Pool Pattern

hermes-bridge today does one connection per process. Paperclip's `nats-transport.ts`
manages a pool: reconnect on drop, deduplicate subscriptions, retry with backoff,
expose a typed `subscribe(subject, handler)` wrapper.

**The idea in plain English:**
Instead of every service managing its own NATS connection from scratch (with its own
reconnect logic, its own error handling), there is one shared pattern. Every service
uses the same wrapper and gets robust connection behavior for free.

**Where to apply:** This pattern should be the base class for every service in
`~/.nats/services/`. It replaces the raw `connect()` call in hermes-bridge with
something that handles drops, retries, and exposes a clean API.

**The pattern in code:**

```javascript
// services/shared/nats-transport.js
export class NatsTransport {
  async connect(url, name) { /* retry loop */ }
  subscribe(subject, handler) { /* deduplicated */ }
  publish(subject, payload) { /* typed, with envelope */ }
  request(subject, payload, timeoutMs) { /* req/reply */ }
  async close() { /* graceful */ }
}
```

**Status:** Pattern documented. hermes-bridge uses raw connect today. Refactor is
low-priority — it works. Apply the pattern when building new services (agent-registry,
github-bridge, aoe-bridge).

---

#### 3.2 — Typed Envelope on Every Message

Paperclip enforced that every NATS message carries a typed header block. We already have
this in `schemas/envelope/envelope.schema.json`. The discipline: never publish raw JSON
directly — always wrap in the envelope.

**The current gap:** hermes-bridge today publishes raw JSON payloads (not wrapped in the
envelope). This means ClickHouse doesn't get a typed `type` field to route by.

**Fix:** Update hermes-bridge to wrap every publish in the envelope:

```javascript
function pub(subject, payload) {
  const envelope = {
    v: 1,
    id: crypto.randomUUID(),
    ts: new Date().toISOString(),
    type: subject,
    source: `agent.hermes.${PROFILE}`,
    subject: subject,
    data: payload
  };
  nc.publish(subject, encoder.encode(JSON.stringify(envelope)));
}
```

**Status:** Gap identified. Low-urgency but should be done before ClickHouse starts
reading hermes subjects at high volume.

---

## Platform 4: Agent of Empires (AoE)

**Tech stack:** Rust binary, tmux-native sessions, git worktrees, Docker per session,
ACP WebSocket protocol, REST API on `:7777`, Tailscale PWA.
**What it does:** Creates isolated coding sessions. Each session is a tmux window, a git
worktree, and optionally a Docker container. A structured event stream (ACP) reports
what the agent is doing in real time: plan panels, tool-call cards, status cards.

**Decision:** AoE is kept. It is the coding session manager for Hermes to delegate into.

### What we extracted

#### 4.1 — ACP Event Schema → NATS Bridge

AoE's ACP (Agent Client Protocol) is a structured WebSocket event stream per session.
Instead of reading raw terminal output, you get typed events:

| ACP Event Type | What it means |
|----------------|---------------|
| `plan.panel` | The agent has formed a plan. Fields: steps, reasoning. |
| `tool.call` | The agent is calling a tool. Fields: tool_name, args. |
| `tool.result` | Tool call returned. Fields: output, success. |
| `status.update` | Agent status changed (thinking / running / waiting / done). |
| `session.end` | Session completed. Fields: exit_code, summary. |

**The idea:** We build an `aoe-bridge` service that subscribes to AoE's ACP WebSocket and
republishes each event as a NATS message. Now Mission Control, the Flutter app, and
ClickHouse all see live coding session state without any of them needing to speak ACP.

**NATS subjects:**

```
evt.agent.{session_id}.plan        → ACP plan.panel event
evt.agent.{session_id}.tool_call   → ACP tool.call event
evt.agent.{session_id}.tool_result → ACP tool.result event
evt.agent.{session_id}.status      → ACP status.update event
evt.agent.{session_id}.done        → ACP session.end event
```

**How Hermes uses AoE:**

1. Hermes ops profile calls `POST http://localhost:7777/sessions` with the task.
2. AoE creates a tmux session + git worktree + optional Docker container.
3. aoe-bridge subscribes to that session's ACP WebSocket.
4. Every event publishes to `evt.agent.{session_id}.>` on NATS.
5. Hermes subscribes to `evt.agent.{session_id}.>` and watches progress.
6. On `evt.agent.{session_id}.done`, Hermes picks up the result.

**Service stub location:** `services/aoe-bridge/` (not yet built).

**Status:** Architecture defined. Bridge service to be built. AoE is running.

---

#### 4.2 — Git Worktree Per Session Pattern

AoE creates a fresh git worktree for each coding session. This means:
- Multiple coding agents can work on different features simultaneously without conflicts.
- Each session has its own branch. Merge happens explicitly when you approve the work.
- If the session fails, you discard the worktree — main branch is untouched.

**Where to apply in onemindos-core:** When `om-task` dispatches a code task to an agent,
it should include a `worktree` field in the task payload specifying the branch name.
The receiving agent (AoE or Claude Code) creates the worktree, does the work, and reports
back with the branch name so it can be reviewed and merged.

**Status:** Pattern documented. No code change needed today. Apply when building task
dispatch in the agent-registry service.

---

## Platform 5: Multica

**Tech stack:** Go 1.26, PostgreSQL 17, Redis, React Native iOS app.
**What it does:** Multi-agent orchestration with daemon-based task lifecycle.
Agents self-register on start via a Go daemon socket. Squads define hierarchical
agent teams with budget caps and delegation rules.

### What we extracted

#### 5.1 — Daemon Self-Registration Pattern (Most Important Multica Idea)

This is the key feature SwarmClaw partially had but Multica implemented cleanly. In
Multica, when any agent process starts, it connects to the daemon socket and says
"I'm here, this is what I can do." The daemon records it. No polling. No external scan.

**The idea in plain English:**
The agent announces itself. The registry doesn't need to go looking. This is more reliable
than port scanning (SwarmClaw) because agents behind NAT or in Docker containers register
themselves — the registry doesn't need network access to find them.

**How we implement this on NATS:**
The `sys.agent.register` subject IS the daemon socket, except it's NATS. Every process
that wants to be discoverable publishes one registration message on startup. This is
exactly what `om_task.py` already does for Hermes agents via `agents.hb.>`.

The registration message:

```json
{
  "v": 1,
  "id": "uuid",
  "ts": "2026-06-13T10:00:00.000Z",
  "type": "sys.agent.register",
  "source": "agent.hermes.coder",
  "subject": "sys.agent.register",
  "data": {
    "agent_id": "hermes.coder",
    "display_name": "Hermes — coder profile",
    "agent_type": "llm",
    "host": "onemind-mac",
    "port": 8642,
    "caps": ["llm.chat", "code.edit", "mcp.filesystem", "mcp.todoist"],
    "harness": "hermes",
    "profile": "coder",
    "version": "0.16.0"
  }
}
```

**Status:** `om_task.py` already handles Hermes heartbeats. Extend to handle
`sys.agent.register` for non-Hermes agents (Claude Code, AoE sessions, etc.).

---

#### 5.2 — Squads: Hierarchical Delegation With Budget Caps

Multica's Squads define a team of agents with a parent and children. The parent agent
can delegate subtasks to child agents but cannot exceed a token/cost budget.

**The idea in plain English:**
Hermes ops (parent) has a budget of $5 for a task. It can delegate to Claude Code
(child) and Hermes coder (child) but those children cannot collectively spend more
than $5 total. When the budget is exhausted, the parent is notified and decides whether
to stop or request more.

**How to apply in onemindos-core:**
The Task Engine (`om-task`) already has task budget fields (planned). Add a `squad_id`
field to tasks. The squad definition lives in KV bucket `squads`:

```json
{
  "squad_id": "ops-squad",
  "parent": "hermes.ops",
  "members": ["hermes.coder", "claude-code.coder"],
  "budget_usd": 5.0,
  "spent_usd": 0.0,
  "policy": "notify_on_80pct"
}
```

When a task is dispatched to a squad member, `om-task` checks the squad budget,
deducts the cost on completion, and publishes `evt.sys.squad.budget.warning` at 80%.

**Status:** Pattern documented. Not yet implemented in `om-task`. Card: P2-S.

---

## Summary: What Gets Built in onemindos-core

This table maps each extracted feature to its implementation location and status.

| # | Feature | Source | Target | Status |
|---|---------|--------|--------|--------|
| 1 | Agent auto-discovery (port + process scan) | SwarmClaw | `services/agent-registry/` | Not built |
| 2 | Daemon self-registration via NATS | Multica | `services/agent-registry/` + `sys.agent.register` | Partial (`om_task.py` has Hermes HB) |
| 3 | A2A task dispatch (execute / status / cancel) | SwarmClaw | `cmd.sys.agent.{id}.execute` subjects | Not built |
| 4 | ACP event bridge (AoE → NATS) | AoE | `services/aoe-bridge/` | Not built |
| 5 | GitHub ↔ NATS sync | Mission Control | `services/github-bridge/` | Not built |
| 6 | Fleet snapshot subject | Mission Control | `evt.sys.fleet.snapshot` via agent-registry | Not built |
| 7 | Git worktree per task pattern | AoE | Task dispatch protocol (worktree field) | Pattern only |
| 8 | Squads + budget caps | Multica | `om-task` extension, `squads` KV bucket | Not built |
| 9 | NATS transport pool base class | Paperclip | `services/shared/nats-transport.js` | Not built |
| 10 | Envelope wrapping in hermes-bridge | Paperclip | `~/.nats/services/hermes-bridge/bridge.js` | Gap identified |

---

## The agent-registry Service (Build Order #1)

Everything above depends on knowing who is alive. That makes the agent-registry service
the first thing to build. Here is what it does:

```
STARTUP:
  1. Connect to NATS
  2. Run bootstrap scan (ports 7777/8642/8644/9119 + ps for claude/codex/opencode/aider)
  3. Write any found agents to `agent-registry` KV bucket
  4. Subscribe to sys.agent.register → write to KV
  5. Subscribe to agents.hb.> (Hermes heartbeats) → update last_seen in KV
  6. Subscribe to sys.agent.heartbeat → update last_seen in KV

EVERY 30s:
  7. Scan KV for agents with last_seen > 2min → mark as lost
  8. Publish evt.sys.fleet.snapshot with full registry state

ON REGISTRY CHANGE:
  9. Publish evt.sys.agent.discovered or evt.sys.agent.lost
```

Location: `services/agent-registry/agent_registry.py`
(Python, same pattern as `om_task.py` — NATS `$SRV` micro service)

LaunchAgent: `io.onemind.agent-registry.plist` (same pattern as hermes-bridge)

---

## NATS Subject Reference — All New Subjects

These subjects are new. They do not conflict with existing subjects in `NATS-DOCTRINE.md`.

```
# Agent registration + discovery
sys.agent.register                     agent startup self-announcement
sys.agent.heartbeat                    30s alive ping
sys.agent.deregister                   graceful shutdown
sys.agent.query                        who's alive? (req/reply)
evt.sys.agent.discovered               new agent joined
evt.sys.agent.lost                     agent timed out

# Fleet state
evt.sys.fleet.snapshot                 full fleet state, every 30s
evt.sys.fleet.node.joined              new compute node online
evt.sys.fleet.node.left                compute node disconnected

# A2A task dispatch
cmd.sys.agent.{agent_id}.execute       dispatch task to specific agent
cmd.sys.agent.execute                  dispatch to any available agent of a type
evt.sys.agent.{agent_id}.task_done     task completed
evt.sys.agent.{agent_id}.task_fail     task failed

# Coding session events (from aoe-bridge)
evt.agent.{session_id}.plan            AoE plan panel
evt.agent.{session_id}.tool_call       AoE tool call
evt.agent.{session_id}.tool_result     AoE tool result
evt.agent.{session_id}.status          AoE status update
evt.agent.{session_id}.done            AoE session end

# GitHub integration
evt.sys.github.issue.created           new GitHub issue
evt.sys.github.issue.closed            issue closed
evt.sys.github.issue.updated           issue updated
cmd.sys.github.issue.create            create issue on GitHub
cmd.sys.github.issue.close             close issue on GitHub

# Squad budget
evt.sys.squad.budget.warning           squad at 80% budget
evt.sys.squad.budget.exhausted         squad budget hit 0
```

---

## KV Buckets Added

| Bucket | Owner | Contents | TTL |
|--------|-------|----------|-----|
| `agent-registry` | agent-registry service | Live agent records | 5 min (refreshed by heartbeat) |
| `squads` | om-task | Squad definitions + budgets | None (permanent) |

---

## Platforms Killed — No Residual Services Needed

| Platform | Decision | Reason |
|----------|----------|--------|
| SwarmClaw | Killed | A2A protocol extracted. Port scan logic extracted. No running service needed. |
| Paperclip | Killed | Transport pool pattern extracted. Envelope discipline extracted. No running service needed. |
| Multica | Never installed | Registration pattern + Squads pattern extracted. React Native app is iOS-only and redundant with our Flutter app. |
