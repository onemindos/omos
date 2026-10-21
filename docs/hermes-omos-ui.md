# Hermes ↔ omos-ui Integration Reference

This document describes what omos-ui expects from Hermes (and any agent) on the NATS bus. It is the contract between the agent runtime and the browser UI. Keep this in sync with `platform/hermes/agent/AGENTS.yaml`.

---

## How omos-ui Connects

The browser SPA connects directly to the in-cluster NATS hub via WebSocket (`wss://ws.onemindos.dev`, port 443 → NATS port 8080 internally). There is no bridge server. The browser is a NATS client with the `onemind` user credentials.

**Critical:** `$SRV.*` subjects are server-local and do NOT propagate through NATS leaf nodes. Hermes connects to the hub via a leaf node sidecar (`localhost:4222` → cluster `:7422`). This means the browser cannot discover Hermes via `$SRV.PING.agents` or `$SRV.INFO.agents`. Discovery uses JetStream instead (see below).

---

## Agent Discovery — What Hermes Must Publish

omos-ui discovers agents in two phases on page load.

### Phase 1: JetStream AGENTS Stream (primary — works for leaf-connected agents)

The browser reads the `AGENTS` JetStream stream using `DeliverPolicy.LastPerSubject` filtered to `agents.hb.>`. This returns the most recent heartbeat per unique subject immediately, with no wait.

**Hermes MUST publish a heartbeat to the AGENTS stream every 30 seconds:**

```
Subject: agents.hb.hermes.zeus.legacy
Subject: agents.hb.hermes.zeus.legacy-edge  (edge body)
```

**Heartbeat payload (JSON):**

```json
{
  "agent": "hermes",
  "owner": "zeus",
  "session": "legacy",
  "instance_id": "hermes-zeus-legacy",
  "name": "Legacy",
  "description": "The authoritative primary — sovereign state, final decisions",
  "version": "v2",
  "metadata": {
    "role": "primary"
  }
}
```

Fields:
- `agent`, `owner`, `session` — required. Must match the subject pattern.
- `instance_id` — required. Unique stable identifier. omos-ui keys agent cards on this. Without it, UI synthesizes `{agent}-{owner}-{session}` as fallback, but stable `instance_id` is preferred.
- `name` — display name shown on the agent card.
- `description` — one-line description shown in the agent card subtitle.
- `metadata.role` — if `"controller"`, omos-ui renders the ControllerPanel (piexec/ccexec) instead of the chat panel.

### Phase 2: `$SRV.INFO.agents` (fallback — hub-direct agents only)

If Hermes is ever reconfigured to connect directly to the hub (not via leaf), standard NATS micro service discovery works and phase 2 will surface it. No action needed — Hermes already registers as a NATS micro service named `agents`. This is a bonus fallback, not the primary path.

### Live heartbeat watch

After discovery, omos-ui subscribes to `agents.hb.*.*.*` live. New or updated heartbeats update agent card state in real time. The heartbeat LED on each agent card flashes on each received message.

---

## Prompt Endpoint — Wire Protocol (Synadia Agent Protocol v0.3)

Hermes registers a NATS micro service endpoint that handles incoming prompts.

**Subject:** `agents.prompt.hermes.zeus.legacy`

The browser sends prompts as NATS `requestMany` calls. The agent streams response chunks back to the reply subject until it sends an empty-body terminator message.

### Inbound envelope (browser → agent)

```json
{
  "type": "query",
  "data": {
    "id": "<uuid>",
    "reply_subject": "_INBOX.abc123",
    "prompt": "What is the current SA picture?",
    "attachments": []
  }
}
```

- `id` — correlation ID for this prompt. The agent echoes it in all response frames.
- `reply_subject` — where to publish response chunks. The browser's `requestMany` reply inbox.
- `prompt` — the user's message text.
- `attachments` — array of `{ filename, base64 }` objects (optional).

### Response frames (agent → reply subject)

**Text chunk:**
```json
{ "type": "response", "data": "partial text..." }
```

**Status / ack:**
```json
{ "type": "status", "data": "ack" }
```

**Permission query (tool gate):**
```json
{
  "type": "query",
  "data": {
    "id": "<prompt-id>",
    "query_id": "<uuid>",
    "reply_subject": "_INBOX.xyz456",
    "question": "Can I execute: kubectl get pods -n hermes?"
  }
}
```
The browser renders a permission prompt and responds by publishing to `reply_subject`.

**Error:**
```json
{ "type": "error", "data": "error message" }
```

**Stream terminator (§6.5):**
An empty-body message (`Uint8Array(0)`) published to the reply subject signals end-of-stream. `requestMany` with `RequestStrategy.SentinelMsg` stops consuming on receipt.

---

## AGENTS JetStream Stream — Configuration

The AGENTS stream must exist in the cluster with these settings:

```
Name:            AGENTS
Subjects:        agents.hb.>
Retention:       Limits
Storage:         File
MaxAge:          90s          (3 missed 30s heartbeats = agent considered offline)
MaxMsgsPer:      1            (one message per unique subject — LastPerSubject semantics)
Replicas:        3
```

This is provisioned via `nats-migration/` — run `./provision.sh` to create/update. The `MaxMsgsPer: 1` combined with `DeliverPolicy.LastPerSubject` in the browser consumer is what gives zero-delay discovery: the stream holds exactly one heartbeat per agent, always the latest.

---

## Subject Reference (agent-relevant subjects)

| Subject | Direction | Purpose |
|---------|-----------|---------|
| `agents.hb.hermes.zeus.legacy` | Hermes → AGENTS stream | Heartbeat / liveness |
| `agents.hb.hermes.zeus.legacy-edge` | Edge → AGENTS stream | Edge body heartbeat |
| `agents.prompt.hermes.zeus.legacy` | Browser → Hermes | Incoming prompts |
| `agents.prompt.hermes.zeus.legacy-edge` | Browser → Edge | Edge-directed prompts |
| `agents.hb.*.*.*` | Hermes subscribe | Live heartbeat watch (UI) |
| `svc.tak.sa` | Browser → tak-agent | Get SA picture (NATS RPC) |
| `svc.tak.chat` | Browser → tak-agent | Send GeoChat to TAK operators |
| `svc.tak.inject_cot` | Browser → tak-agent | Inject CoT marker onto map |
| `svc.tak.mission.*` | Browser → tak-agent | Mission create/list/add |

---

## Agent Identity (from AGENTS.yaml)

| Field | Legacy (Mother) | Legacy (Edge) |
|-------|----------------|---------------|
| `agent` | `hermes` | `hermes` |
| `owner` | `zeus` | `zeus` |
| `session` | `legacy` | `legacy-edge` |
| `instance_id` | `hermes-zeus-legacy` | `hermes-zeus-legacy-edge` |
| `bus_subject` | `agents.prompt.hermes.zeus.legacy` | `agents.prompt.hermes.zeus.legacy-edge` |
| `role` | `primary` | `edge-body` |
| `location` | DOKS / hermes namespace | Zeus's MacBook Pro |
| `nats` | leaf sidecar (`localhost:4222`) | local leaf (`localhost:4222`) |

---

## Adding a New Agent (servant)

1. Add an entry to `platform/hermes/agent/AGENTS.yaml` under `planned:` (or `agents:` once it runs)
2. Create `platform/hermes/agent/agents/{session}/SOUL.md` from `platform/hermes/agent/agents/_template/SOUL.md`
3. Ensure the agent publishes heartbeats to `agents.hb.hermes.zeus.{session}` every ≤30s with `instance_id` set
4. Ensure the agent listens on `agents.prompt.hermes.zeus.{session}` and speaks Synadia Agent Protocol v0.3
5. The agent will appear in omos-ui automatically on next heartbeat — no UI changes needed
