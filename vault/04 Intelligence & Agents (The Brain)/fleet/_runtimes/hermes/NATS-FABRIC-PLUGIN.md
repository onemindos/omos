---
domain: 01
name: NATS Fabric Plugin — Legacy Meta Agent
agent: Legacy
status: active
version: 1.2.0
created: 2026-05-22
updated: 2026-05-22
tags: [nats, hermes, plugin, fabric, context, agents]
---

# NATS Fabric Plugin — Legacy Meta Agent (v1.2)

## Purpose

Gives Hermes/Legacy two capabilities via a single native Hermes plugin:

1. **Caller SDK** — Legacy can discover, prompt, and fan-out to any agent on the NATS fabric
2. **Context Awareness** — Legacy sees Zeus's live location/status automatically, without being asked

## Plugin Location

```
~/.hermes/plugins/nats-fabric/
├── plugin.yaml       ← manifest (name: nats-fabric, version: 1.2.0)
└── __init__.py       ← Python module
```

---

## Tools (7 total)

### Agent Orchestration

| Tool | What it does | Params |
|------|-------------|--------|
| `fabric_discover` | List all live agents on the fabric | — |
| `fabric_prompt` | Send a prompt to one agent, get response | `agent`, `prompt`, `timeout` |
| `fabric_fan_out` | Prompt multiple agents in parallel | `agents[]`, `prompt`, `timeout` |
| `fabric_status` | Health/version check of all agents | — |

### Context (zeus-context KV)

| Tool | What it does | Params |
|------|-------------|--------|
| `context_get` | Read current location, focus_mode, status | — |
| `context_set` | Write a context value (notifies all watchers) | `key`, `value` |
| `events_replay` | Pull last N events from AGENT_EVENTS JetStream | `limit`, `subject` |

**Valid `context_set` keys:** `location`, `focus_mode`, `status`

```bash
# Example values
location:   home-office | farm | gym | driving | coffee-shop
focus_mode: deep-work | active | rest | unavailable
status:     online | busy | offline
```

---

## How Context Injection Works

The plugin watches the `zeus-context` KV bucket in the background and automatically
enriches Legacy's view of each conversation when Zeus's context changes.

```
Zeus changes KV (nats CLI or context_set tool)
    ↓
_kv_watch_loop() detects change → sets _context_changed = True
    ↓
Zeus sends any message to Legacy
    ↓
pre_llm_call hook fires (before every LLM call, all paths)
    ↓
Sees _context_changed = True → returns {"context": "[Zeus context update: ...]"}
    ↓
Hermes appends to user message (never touches system prompt — preserves Bedrock cache)
    ↓
Legacy knows Zeus is at the gym without being told
```

**The injection only fires when context has changed** since the last injection. Not on every message.

---

## Hooks (3 total)

| Hook | Type | When it fires | What it does |
|------|------|---------------|-------------|
| `on_session_start` | sync | New session created | Schedules NATS connection warm-up |
| `on_session_end` | sync | Session ends | Schedules NATS drain |
| `pre_llm_call` | async | Before every LLM call (all paths) | Starts KV watch on first call; injects context when changed |

**Important:** `invoke_hook()` in Hermes is synchronous. All hooks MUST be plain `def` (not `async def`).
Async work is scheduled via `asyncio.get_running_loop().create_task()` from within sync hooks.

**`pre_llm_call` vs `pre_gateway_dispatch`:**
- `pre_gateway_dispatch` — fires only in the gateway process (Telegram/Discord messages)
- `pre_llm_call` — fires inside the conversation loop for ALL paths (Telegram, API server, CLI)
- Use `pre_llm_call` for context injection. Use `pre_gateway_dispatch` only for gateway-specific flow control.

---

## KV Watch — Persistence Model

The KV watch background task (`_kv_watch_loop`) is an asyncio task that reconnects on error.
Its persistence depends on which process it runs in:

| Path | Process model | Watch persists? |
|------|--------------|-----------------|
| Telegram message | Gateway main process | ✅ Yes — starts on first message, runs forever |
| API `/v1/chat/completions` | Subprocess (spawned per request) | ❌ No — subprocess exits after response |

For API calls, `context_get` still works — it does a direct KV fetch on demand (no background task needed).
The auto-injection only works for Telegram sessions, which is the primary use case.

---

## Schema Format

Hermes's `registry.get_definitions()` wraps tool schemas as `{"type": "function", "function": schema}`.
**Do NOT include the `type/function` wrapper in the schema you pass to `ctx.register_tool()`.**

```python
# CORRECT
ctx.register_tool(name="context_set", schema={
    "name": "context_set",
    "description": "...",
    "parameters": {"type": "object", "properties": {...}, "required": [...]}
})

# WRONG — produces double-wrap, LLM sees empty parameters
ctx.register_tool(name="context_set", schema={
    "type": "function",
    "function": {"name": "context_set", "parameters": {...}}
})
```

---

## Agent Tokens

| Token | Agent | Subject | Notes |
|-------|-------|---------|-------|
| `cc` | Claude Code | `agents.prompt.cc.zeus.cc-main` | Deep coding, file ops |
| `oc` | OpenClaw | `agents.prompt.oc.zeus.legacy` | Gateway, device control |
| `pi` | Pi Agent | `agents.prompt.pi.onemind.pi-main` | Fast coding (interactive-only) |
| `hermes` | Hermes | `agents.prompt.hermes.zeus.hermes-main` | Self-reference |

---

## Dependencies

- **nats-py** v2.14.0 — in `~/.hermes/hermes-agent/venv/`
- **NATS server** running at `nats://127.0.0.1:4222`
- **Agent host shims** at `~/.nats/agents/` (cc-host.mjs, hermes-host.mjs)

---

## CLI Commands

```bash
# Test context tools directly
nats kv put zeus-context location gym          # change a value
nats kv get zeus-context location              # verify

# Test via Legacy (API path)
curl -s http://127.0.0.1:8642/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"hermes-agent","messages":[{"role":"user","content":"context_get"}],"stream":false}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['choices'][0]['message']['content'])"

# Plugin management
hermes plugins list
hermes plugins enable nats-fabric
hermes gateway restart     # reload after changes to __init__.py

# Watch logs
tail -f ~/.hermes/logs/agent.log | grep nats_fabric
```

---

## NATS MCP Server (separate layer)

A Node.js MCP server gives Legacy (and all other agents) direct NATS access without the plugin:

| Item | Value |
|------|-------|
| Code | `~/.nats/mcp-server/server.js` |
| Configured | `~/.hermes/config.yaml → mcp_servers.nats` |
| Tools | `nats_request`, `nats_publish`, `nats_kv_get`, `nats_kv_put`, `nats_kv_keys`, `nats_stream_read` |

The **plugin** adds: agent orchestration, auto-injection, smart context tools.
The **MCP server** adds: raw fabric access for any agent (including Claude Code, Copilot).

---

## Known Issues / Gotchas

- Pi Agent is interactive-only — fabric_prompt times out (needs stdin)
- Plugin handlers must accept `**kwargs` (Hermes injects `task_id`)
- First NATS connect after gateway restart may add ~1s to first tool call
- KV watch auto-injection only works for Telegram sessions (not API test calls)

---

## Related

- [[NATS-INFRASTRUCTURE]] — full infra (server config, JetStream, KV buckets, MCP server)
- `~/.nats/agents/cc-host.mjs` — Claude Code host shim
- `~/.nats/agents/hermes-host.mjs` — Hermes host shim
- `~/.openclaw/openclaw.json → channels.nats` — OpenClaw side config
