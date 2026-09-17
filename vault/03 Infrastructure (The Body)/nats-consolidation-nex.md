# NATS Service Consolidation & NEX Integration

**Date:** 2026-06-14  
**Author:** Zeus Delacruz  
**Status:** Production — Implemented

---

## What Was Done

Reduced 9 individual NATS LaunchAgents down to 3 consolidated services, and installed the NEX execution engine for future zero-plist service deployment.

### Before (9 LaunchAgents)

```
com.onemind.subscriber.agents      subscriber.py agents     prompt.agent.>
com.onemind.subscriber.comms       subscriber.py comms      events.comms.>
com.onemind.subscriber.health      subscriber.py health     events.health.>
com.onemind.subscriber.home        subscriber.py home       events.zeus.>
com.onemind.subscriber.security    subscriber.py security   events.security.>
io.onemind.ha-bridge               ha-bridge/bridge.js
io.onemind.hermes-bridge           hermes-bridge/bridge.js
io.onemind.session-bridge          session-bridge/bridge.js
io.onemind.task-bridge             task-bridge/bridge.js
```

### After (3 LaunchAgents)

```
com.onemind.subscriber-runner      ~/subscribers/subscriber_runner.py   (all 5 channels)
com.onemind.bridge-runner          ~/.nats/services/bridge-runner/runner.js  (all 4 bridges)
com.onemind.nex-node               ~/go/bin/nex node up                  (NEX execution engine)
```

---

## Subscriber Runner

**File:** `~/subscribers/subscriber_runner.py`  
**Log:** `~/subscribers/logs/runner.log`  
**LaunchAgent:** `~/Library/LaunchAgents/com.onemind.subscriber-runner.plist`

Runs all 5 subscriber channels in a single Python process with one NATS connection. Each channel is an independent async task with automatic restart on crash (3s backoff) — one channel failing does not affect the others.

### Adding a New Channel

Edit the `CHANNELS` dict at the top of `subscriber_runner.py` — no new plist needed:

```python
CHANNELS = {
    "home":        "events.zeus.>",
    "security":    "events.security.>",
    "health":      "events.health.>",
    "comms":       "events.comms.>",
    "agents":      "prompt.agent.>",
    "mynewthing":  "events.mynewthing.>",   # ← add one line
}
```

Then restart the runner:

```bash
launchctl kickstart -k gui/$(id -u)/com.onemind.subscriber-runner
```

### Channel Architecture

```
NATS Bus
   │
   ├── events.zeus.>        → [home]     task supervisor → rules KV → actions
   ├── events.security.>    → [security] task supervisor → rules KV → actions
   ├── events.health.>      → [health]   task supervisor → rules KV → actions
   ├── events.comms.>       → [comms]    task supervisor → rules KV → actions
   └── prompt.agent.>       → [agents]   task supervisor → rules KV → actions
```

Rules are loaded from NATS KV bucket `rules` on every event — live config changes with no restarts.

---

## Bridge Runner

**File:** `~/.nats/services/bridge-runner/runner.js`  
**Log:** `~/.nats/services/bridge-runner/runner.log`  
**LaunchAgent:** `~/Library/LaunchAgents/com.onemind.bridge-runner.plist`

Runs all 4 bridges in a single Node.js process. Uses NATS transport-node v3 API. Each bridge is supervised independently — a bridge crash restarts only that bridge after 5s.

### Bridges Included

| Bridge | What It Does | Subjects |
|--------|-------------|---------|
| `ha` | Home Assistant ↔ NATS (WebSocket) | `cmd.ha.service.call`, `cmd.ha.state.get`, `events.ha.*` |
| `hermes` | Polls Hermes gateway status + sessions | `evt.hermes.{profile}.*`, `evt.hermes.health` |
| `session` | Hermes session control via NATS | `cmd.session.>`, `evt.session.*` |
| `task` | Task engine event forwarding | `cmd.task.>`, `evt.task.*` |

### Adding a New Bridge

Add an async function and register it in the `BRIDGES` object in `runner.js`:

```js
async function myBridge(nc) {
  // nc is the shared NATS connection
  const sub = nc.subscribe('cmd.myservice.>');
  for await (const m of sub) {
    // handle messages
  }
}

const BRIDGES = {
  ha:      haBridge,
  hermes:  hermesBridge,
  session: sessionBridge,
  task:    taskBridge,
  myservice: myBridge,   // ← add one line
};
```

Then restart:

```bash
launchctl kickstart -k gui/$(id -u)/com.onemind.bridge-runner
```

---

## NEX — NATS Execution Engine

**Binary:** `~/go/bin/nex`  
**Built from:** `github.com/synadia-io/nex` (source build — June 14, 2026)  
**Log:** `~/.nats/logs/nex-node.log`  
**LaunchAgent:** `~/Library/LaunchAgents/com.onemind.nex-node.plist`  
**Node name:** `onemind-mac`  
**Nexus:** `onemind`

### What Is NEX?

NEX is the official NATS process execution engine made by **Synadia** — the company that created and maintains NATS.io. Synadia is also the commercial entity behind Synadia Cloud (managed NATS hosting) and is the primary steward of the NATS open-source project, which is a CNCF incubating project.

NEX is not a third-party tool. It is first-party Synadia infrastructure, developed alongside the NATS server itself. Current version: v0.4.1 (pre-1.0, active development, consistent release cadence).

### What NEX Does

NEX is a distributed execution engine that runs workloads (processes, functions) directly over NATS with no additional orchestration layer. Instead of creating a LaunchAgent per service, you deploy workloads via NATS publish commands. NEX manages process lifecycle, restarts, observability, and scoped NATS credentials for each workload.

```
Traditional approach:           NEX approach:
  new service → new plist         new service → nex workload deploy
  restart → launchctl             restart → nex workload stop/deploy
  status → launchctl list         status → nex node ls / nex workload list
  logs → log file                 logs → NATS subjects + log file
```

### NEX Architecture on onemind-mac

```
launchd
  └── com.onemind.nex-node
        └── nex node (onemind-mac, nexus=onemind)
              ├── native agent (go_exec) — runs OS binaries
              └── [future] wasm agent — runs WebAssembly workloads
```

NEX registers itself on NATS under `$NEX.>` subjects. All control plane operations happen over NATS.

### Deploy a Workload via NEX

```bash
# Deploy a native binary
nex workload deploy \
  --name my-worker \
  --nats.servers nats://127.0.0.1:4222 \
  --nats.user onemind \
  --nexus onemind \
  /path/to/binary

# List running workloads
nex workload list \
  --nats.servers nats://127.0.0.1:4222 \
  --nats.user onemind \
  --nexus onemind

# Stop a workload
nex workload stop my-worker \
  --nats.servers nats://127.0.0.1:4222 \
  --nats.user onemind \
  --nexus onemind

# List NEX nodes
nex node ls \
  --nats.servers nats://127.0.0.1:4222 \
  --nats.user onemind \
  --nexus onemind
```

### Workload Types Supported

| Type | Description |
|------|-------------|
| `native` | Any OS binary (Python, Node, Go, Rust, etc.) |
| `function` | Lightweight JS functions run in embedded V8 |
| `wasm` | WebAssembly modules (sandboxed, near-zero overhead) |

### Updating NEX

NEX is built from source (no pre-built macOS arm64 release). To update:

```bash
cd /tmp/nex-src && git pull
cd cmd/nex && go build -o ~/go/bin/nex .
launchctl kickstart -k gui/$(id -u)/com.onemind.nex-node
```

---

## Scaling Model

This is the new pattern for adding any NATS-connected service to onemind:

| Service type | How to add |
|---|---|
| New event subscriber (Python) | Add 1 line to `CHANNELS` in `subscriber_runner.py` |
| New bridge (Node.js) | Add 1 function + 1 line to `BRIDGES` in `bridge-runner/runner.js` |
| New standalone worker (any language) | `nex workload deploy` — no file changes at all |
| New infrastructure service (needs its own plist) | Only for services that cannot run under NEX or the runners |

The target steady state: `io.nats.server`, `io.nats.surveyor`, `com.onemind.subscriber-runner`, `com.onemind.bridge-runner`, and `com.onemind.nex-node` are the only 5 NATS-layer LaunchAgents that should ever exist. Everything else deploys through them.

---

## Quick Reference

```bash
# Check all NATS service health
launchctl list io.nats.server io.nats.surveyor \
  com.onemind.subscriber-runner com.onemind.bridge-runner com.onemind.nex-node

# Tail subscriber logs (all channels)
tail -f ~/subscribers/logs/runner.log

# Tail bridge logs (all bridges)
tail -f ~/.nats/services/bridge-runner/runner.log

# Tail NEX node log
tail -f ~/.nats/logs/nex-node.log

# NATS server monitoring
open http://localhost:8222

# Prometheus metrics (Surveyor)
curl http://localhost:7777/metrics
```
