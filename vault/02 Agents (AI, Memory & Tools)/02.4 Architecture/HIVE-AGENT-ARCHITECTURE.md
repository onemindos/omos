---
id: ui-02.4-hive-20260305
title: "Hive Agent Architecture — A2A + SubAgent Design"
type: architecture
status: active
created: 2026-03-05
modified: 2026-03-05
owner: zeus
tags: [type/architecture, domain/ui, topic/agents, topic/a2a, topic/subagents]
ai_summary: "Live architecture of the 4-agent OpenClaw Hive: Legacy + HP + LE + GE. Covers A2A messaging, SubAgent scaling, MCP integration, and channel routing."
---

# Hive Agent Architecture — OneMind OS

## The Equation

```
HUMAN (Zeus) + AI (LEGACY) = ONE MIND
```

---

## Layer 1 — The Hive (4 Sovereign Agents)

```
ZEUS (Human)
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│              OpenClaw Gateway (:18789)                   │
│                   VPS: 134.122.8.135                    │
│              Tailscale: 100.104.61.29                   │
└────────────────────────┬────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
   ┌─────────┐      ┌─────────┐     ┌─────────┐     ┌─────────┐
   │ LEGACY  │◄────►│   HP    │◄───►│   LE    │◄───►│   GE    │
   │ (UI)    │      │ (Self)  │     │(Family) │     │(Wealth) │
   │gpt-4o   │      │ sonnet  │     │ sonnet  │     │ opus    │
   └─────────┘      └─────────┘     └─────────┘     └─────────┘
   00-24 UI         25-49 HP        50-74 LE        75-99 GE
   Orchestrator     PI Score        LI Score        WI Score
```

**A2A (Agent-to-Agent) messaging** is live — all 4 agents can delegate tasks to each other directly.

---

## Layer 2 — Agent Profiles

| Agent | ID | Model | Role | Channels |
|-------|----|-------|------|---------|
| **Legacy** | `legacy` | `azure-openai/gpt-4o` | Sovereign orchestrator, default agent | Discord, BlueBubbles |
| **HP** | `hp` | `anthropic/claude-sonnet-4-5` | Personal coach, body/mind/money | Session picker |
| **LE** | `le` | `anthropic/claude-sonnet-4-5` | Family steward, estate, kids | Session picker |
| **GE** | `ge` | `anthropic/claude-opus-4-6` | Business advisor, revenue, ventures | Session picker |

---

## Layer 3 — SubAgent Scaling (Current Config)

```json
"subagents": {
  "maxConcurrent": 8,
  "maxSpawnDepth": 3,
  "maxChildrenPerAgent": 5,
  "model": "azure-openai/gpt-4o-mini"
}
```

### What This Means

| Setting | Value | Implication |
|---------|-------|-------------|
| `maxConcurrent` | **8** | Up to 8 subagents running in parallel at any time |
| `maxSpawnDepth` | **3** | Legacy → Child → Grandchild → Great-grandchild (3 levels deep) |
| `maxChildrenPerAgent` | **5** | Each agent can spawn up to 5 direct children |
| `model` | `gpt-4o-mini` | Subagents use faster/cheaper model by default |

### Max Theoretical Task Tree

```
Legacy (root)
├── Subagent-1 (depth 1)
│   ├── Subagent-1a (depth 2)
│   │   └── Subagent-1a-i (depth 3) ← MAX DEPTH
│   ├── Subagent-1b (depth 2)
│   └── Subagent-1c (depth 2)
├── Subagent-2 (depth 1)
│   └── ...
└── Subagent-5 (depth 1) ← MAX CHILDREN
```

Max simultaneous: **8 concurrent** across the entire tree.

### Scaling Recommendations

- **Current (8/3/5)**: Good for a single-user assistant. Plenty of headroom.
- **Increase `maxConcurrent`** to 12-16 if complex parallel research tasks are queued (VPS has 4GB RAM, 4 cores — can handle it)
- **Keep `maxSpawnDepth` at 3** — deeper trees create debugging nightmares
- **`gpt-4o-mini` as subagent model is correct** — fast, cheap, sufficient for leaf tasks

---

## Layer 4 — MCP Integration (via mcporter)

```
Legacy
  └── exec: mcporter call <server.tool>
              │
              ├── zoho-mail       (19 tools, SSE)    → Zoho Mail API
              ├── composio-todoist (76 tools, HTTP)  → Todoist API
              ├── github          (26 tools, stdio)  → GitHub API
              └── filesystem      (14 tools, stdio)  → VPS filesystem
```

### MCP Servers (Live as of 2026-03-05)

| Server | Transport | Tools | Auth |
|--------|-----------|-------|------|
| `zoho-mail` | SSE | 19 | API key in URL |
| `composio-todoist` | Streamable HTTP | 76 | `x-api-key` header |
| `github` | stdio (npx) | 26 | `GITHUB_TOKEN` env |
| `filesystem` | stdio (npx) | 14 | None (restricted paths) |

---

## Layer 5 — Channel Routing (Bindings)

```
Discord messages  ──► Legacy agent
BlueBubbles (iMessage) ──► Legacy agent

HP / LE / GE ──► Accessible via:
  - OpenClaw Studio (web UI)
  - openclaw agent --agent <id>
  - Additional channel bindings (WhatsApp, Telegram)
```

To route a separate WhatsApp number to GE agent:
```json
{
  "agentId": "ge",
  "match": { "channel": "whatsapp", "peer": { "kind": "direct", "id": "+1XXXXXXXXXX" } }
}
```

---

## Layer 6 — Memory Architecture

### Current State

| Plugin | Status | Backend | Capability |
|--------|--------|---------|------------|
| `memory-core` | **loaded** | Flat files (`memory/YYYY-MM-DD.md`) | Basic daily logs, manual search |
| `memory-lancedb` | **disabled** | LanceDB vector DB | Semantic auto-recall, auto-capture |

### Recommendation → Enable memory-lancedb

`memory-lancedb` adds:
- **Auto-recall**: Automatically injects relevant memories into agent context
- **Semantic search**: Vector similarity vs keyword-only in memory-core
- **Auto-capture**: Passages flagged as important are stored without manual prompting
- **Compaction-aware**: Works alongside the existing `memoryFlush` compaction

Both plugins can coexist. Enable:
```bash
openclaw plugins enable memory-lancedb
# Then restart:
systemctl restart openclaw
```

---

## Layer 7 — A2A Messaging Config

```json
"tools": {
  "agentToAgent": {
    "enabled": true,
    "allow": ["legacy", "hp", "le", "ge"]
  }
}
```

### A2A Use Cases

| From | To | When |
|------|----|------|
| Legacy | HP | "What's Zeus's PI Score this week?" |
| Legacy | LE | "Draft a family update for the kids' schedule" |
| Legacy | GE | "Analyze the Circle.so revenue this month" |
| HP | Legacy | "Flag this health metric to Zeus immediately" |
| GE | Legacy | "Send Zeus a deal summary via Discord" |

### SubAgent vs A2A — When to Use Which

| Use | Mechanism |
|-----|-----------|
| Parallel research subtasks | SubAgent (spawn + parallel) |
| Domain-specific expertise | A2A (route to domain agent) |
| Long-running background work | SubAgent with isolated session |
| Cross-domain data synthesis | A2A → Legacy aggregates |

---

## Exec Policy

```json
"tools": {
  "exec": {
    "host": "gateway",
    "security": "full",
    "applyPatch": { "enabled": true }
  }
}
```

`security: "full"` — Legacy can run shell commands without per-command approval. This enables the mcporter skill and all other CLI-based skills. Change to `"allowlist"` for tighter control if needed.

---

## Infrastructure

```
VPS: DigitalOcean (134.122.8.135)
OS: Ubuntu 22.04
RAM: 4GB (362MB used by gateway in idle)
User: onemind
Service: systemctl openclaw (auto-start, auto-restart)
Gateway port: 18789 (loopback only)
Tailscale IP: 100.104.61.29
Config: /home/onemind/.openclaw/openclaw.json
Logs: journalctl -u openclaw
```

---

## Quick Reference Commands

```bash
# Check all agents
ssh root@134.122.8.135 'su - onemind -c "openclaw agents list"'

# Send message to a specific agent
ssh root@134.122.8.135 'su - onemind -c "set -a && source /opt/onemind/.env && set +a && openclaw agent --agent hp --session-id <id> --message \"...\" --json"'

# Check MCP servers
ssh root@134.122.8.135 'su - onemind -c "set -a && source /opt/onemind/.env && set +a && mcporter list"'

# Gateway logs
ssh root@134.122.8.135 'journalctl -u openclaw --since "1 hour ago" --no-pager'

# Restart gateway
ssh root@134.122.8.135 'systemctl restart openclaw'
```
