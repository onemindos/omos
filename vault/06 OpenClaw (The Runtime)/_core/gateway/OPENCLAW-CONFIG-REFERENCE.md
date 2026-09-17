---
title: OpenClaw Config Reference
type: reference
status: active
created: 2026-05-10
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO, domain/02-infra]
---

# OpenClaw Config Reference

> Schema reference for `openclaw.json`. Live configs live in the private workspace repos — never in this vault.
> This document maps every key section, what it controls, and where secrets are stored in 1Password.

---

## Config File Locations

| VM | Live config path | Backup |
|----|-----------------|--------|
| omos-core | `/home/ubuntu/.openclaw/openclaw.json` | `.openclaw/openclaw.json.bak` (last write) |
| omos-sentinel | `/home/ubuntu/.openclaw/openclaw.json` | `.openclaw/openclaw.json.bak` (last write) |

Edit via CLI: `openclaw config patch --file patch.json5` or `openclaw config set <key> <value>`

---

## Top-Level Structure

```json
{
  "agents": { ... },
  "gateway": { ... },
  "plugins": { ... },
  "talk": { ... },
  "channels": { ... },
  "session": { ... },
  "tools": { ... },
  "auth": { ... },
  "mcp": { ... },
  "hooks": { ... }
}
```

---

## `agents`

### `agents.defaults`

| Field | omos-core | omos-sentinel | Purpose |
|-------|-----------|---------------|---------|
| `model.primary` | `github-copilot/claude-sonnet-4.6` | `github-copilot/claude-sonnet-4.6` | Default model for all agents |
| `model.fallbacks` | `["amazon-bedrock/us.anthropic.claude-sonnet-4-6"]` | same | Fallback if primary fails |
| `thinkingDefault` | `adaptive` | `adaptive` | Global default — model decides when reasoning helps |
| `sandbox.mode` | `off` | `off` | No Docker isolation — exec runs directly on host |
| `maxConcurrent` | **`6`** | **`2`** | Max simultaneous agent runs (core has room; sentinel is RAM-constrained) |

### `agents.defaults.subagents` — calibrated for 23-agent fleet

Subagents are short-lived workers spawned by an agent during a task. Haiku 4.5 (fast + cheap) with thinking enabled so they can reason through multi-step work.

| Field | omos-core | omos-sentinel | Purpose |
|-------|-----------|---------------|---------|
| `model` | `github-copilot/claude-haiku-4.5` | same | Fast workers — Haiku is the right call |
| `thinking` | **`low`** | **`off`** | Core workers get basic reasoning; sentinel workers are simple pings |
| `maxConcurrent` | **`12`** | **`4`** | 6 agents may dispatch simultaneously on core; sentinel capped tight |
| `maxSpawnDepth` | **`2`** | `1` (default) | Legacy → Eternal → subagent is valid now with a2a live |
| `maxChildrenPerAgent` | **`6`** | `5` (default) | Slight bump — Legacy benefits during parallel research |
| `runTimeoutSeconds` | **`300`** | **`300`** | 5-min hard kill — prevents stuck workers holding slots forever |
| `archiveAfterMinutes` | **`60`** | **`30`** | Auto-clean completed sessions (sentinel tighter — 3.8GB RAM) |
| `announceTimeoutMs` | **`5000`** | unset | Fail-fast if subagent doesn't announce within 5s |
| `requireAgentId` | unset | unset | Not enforced — spawns can be anonymous for now |

### `tools.subagents` — tool-level safety

| Field | Value | Purpose |
|-------|-------|---------|
| `tools.deny` | `["elevated"]` | Subagents can never escalate to host exec — only parent agents gated by execApprovals |

**Design rationale:** Core has 31GB RAM and runs 6 super-agents that can each dispatch parallel workers. 12 concurrent subagents with depth-2 nesting covers the realistic ceiling for morning-brief, multi-agent flows, and agent-to-agent dispatches without overloading. Sentinel on 3.8GB RAM gets tight caps — its subagents are simple health pings and status checks, no deep reasoning needed.

### `agents.list` — per-agent overrides (omos-core)

| Agent | Model | thinkingDefault | Notes |
|-------|-------|-----------------|-------|
| `main` (Legacy) | inherited (Sonnet 4.6) | `high` | Orchestrator — deep thinking on routing/strategy |
| `eternal` | inherited (Sonnet 4.6) | `high` | Empire/capital — strategic reasoning required |
| `sovereignty` | inherited (Sonnet 4.6) | `medium` | Infra — correctness over speed |
| `zenith` | inherited (Sonnet 4.6) | `low` | Conversational HP — thinking rarely helps |
| `unity` | inherited (Sonnet 4.6) | `low` | Conversational LE — same rationale |
| `sentinel` | `github-copilot/claude-sonnet-4.6` | `adaptive` | Upgraded from Haiku — failover recovery needs full model |

### `agents.list` — per-agent overrides (omos-sentinel)

| Agent | Model | thinkingDefault | Notes |
|-------|-------|-----------------|-------|
| `main` (Sentinel) | inherited (Sonnet 4.6) | `adaptive` | Watchdog + failover — adaptive thinking for recovery decisions |

---

## `tools`

| Field | Live value | Purpose |
|-------|-----------|---------|
| `tools.loopDetection.enabled` | `true` | Detect and break agent reasoning loops |
| `tools.alsoAllow` | `["lobster"]` | Adds Lobster workflow shell on top of default tool profile |
| `tools.agentToAgent.enabled` | `true` | Legacy can invoke domain agents as a tool call |
| `tools.agentToAgent.allow` | `["zenith","unity","eternal","sovereignty","sentinel"]` | Explicit target allowlist — no unconstrained chaining |
| `tools.subagents.tools.deny` | `["elevated"]` | Subagents blocked from host-level exec escalation |

**`agentToAgent`** lets Legacy dispatch domain work mid-session. Example: a cross-domain request comes into Legacy, it calls `eternal` for the capital angle and `zenith` for the health angle simultaneously, synthesises the result.

**`subagents.tools.deny`** ensures spawned workers can never escalate to elevated/host execution. Only parent agents — gated by execApprovals — can run privileged commands.

---

## `session`

| Field | omos-core | omos-sentinel | Purpose |
|-------|-----------|---------------|---------|
| `dmScope` | `per-channel-peer` | — | Session isolation per channel+user pair |
| `agentToAgent.maxPingPongTurns` | `3` | — | Caps agent-to-agent reply chains — prevents runaway loops |
| `threadBindings.defaultSpawnContext` | `fork` | — | Spawned work sessions inherit parent conversation context |
| `threadBindings.spawnSessions` | `true` | — | Thread-bound work session spawning enabled |
| `threadBindings.idleHours` | `24` | — | Thread sessions unfocus after 24h idle |
| `maintenance.mode` | `enforce` | `enforce` | Actively prunes old sessions |
| `maintenance.pruneAfter` | `30d` | `14d` | Sentinel tighter — 3.8GB RAM node |
| `maintenance.maxEntries` | `500` | `200` | Hard cap on session store entries |

---

## `channels`

### execApprovals

Both channels on both nodes now require approval before any agent exec. This gates Sovereignty and Sentinel from running host commands without Zeus confirming.

| Node | Channel | Target | Approver |
|------|---------|--------|---------|
| omos-core | Discord | `#approvals` channel | Zeus (`1024395579742228570`) |
| omos-core | Telegram | DM | Zeus (`577573558`) |
| omos-sentinel | Discord | `#approvals` channel | Zeus (`1024395579742228570`) |
| omos-sentinel | Telegram | DM | Zeus (`577573558`) |

`cleanupAfterResolve: true` on Discord — approval messages auto-delete after action.

### Telegram

| Field | Secret location |
|-------|----------------|
| `channels.telegram.botToken` | `op://aaeyris5z44pdqspk277u2oppi/c2hsruwzaum6z4d7qjjyduadpq/credential` |
| `channels.telegram.dmPolicy` | `allowlist` — no secret |
| `channels.telegram.allowFrom` | `["577573558"]` (Zeus) — no secret |

### Discord

| Field | Secret location |
|-------|----------------|
| `channels.discord.accounts.default.token` | `op://aaeyris5z44pdqspk277u2oppi/qwomupt2iajig4twcv6ijdkyj4/token` |
| `channels.discord.groupPolicy` | `allowlist` — no secret |
| `channels.discord.allowFrom` | `["1024395579742228570"]` (Zeus) — no secret |

---

## `auth`

| Profile | Purpose | Secret location |
|---------|---------|----------------|
| GitHub Copilot | Primary model provider (both nodes) | `op://aaeyris5z44pdqspk277u2oppi/GitHub Copilot (OpenClaw)/token` |
| AWS Bedrock | Fallback model provider (both nodes) | `op://aaeyris5z44pdqspk277u2oppi/AWS Bedrock (OpenClaw)/accessKeyId` + `/secretAccessKey` |

---

## `plugins`

70+ plugin entries in Core. Enable/disable without editing JSON:

```bash
openclaw plugins enable <name>
openclaw plugins disable <name>
openclaw plugins list
```

Key enabled: `amazon-bedrock`, `github-copilot`, `discord`, `telegram`, `tavily`, `memory-core`, `active-memory`, `openshell`, `webhooks`, `file-transfer`, `thread-ownership`, `llm-task`, `tokenjuice`, `document-extract`, `browser`, `deepgram`, `azure-speech`, `microsoft-foundry`, `skill-workshop`

---

## `talk`

Voice config. See `nodes/TALK-MODE-CONFIG.md` for full architecture.

| Node | Provider | Model | Voice | Status |
|------|----------|-------|-------|--------|
| omos-core | Azure OpenAI Realtime (`openai` key) | `gpt-realtime-1.5` | `coral` | Config correct — deployment pending in Azure portal |
| omos-sentinel | Not configured | — | — | Not active |

Azure key source: `AZURE_OPENAI_API_KEY` injected via `/home/ubuntu/.openclaw/gateway.systemd.env`.

---

## `mcp`

| Server | URL | Auth |
|--------|-----|------|
| `circle-mcp` | `https://app.circle.so/api/mcp` | Bearer token in config |

---

## `hooks`

| Hook | Node | Source | Purpose |
|------|------|--------|---------|
| `onemind-session-logger` | omos-core | `/home/ubuntu/onemind/_system/openclaw/hooks/onemind-session-logger` | Logs all agent sessions |

---

## Config Sync Architecture

Service: `openclaw-config-sync.service` — running on both VMs, enabled at boot. Verified active 2026-05-10.

```
You change openclaw.json (via CLI or direct edit)
    ↓ inotifywait detects write event (instant)
    ↓ 3s debounce — waits for burst of changes to settle
    ↓ rsync copies changed files → ~/openclaw-workspaces/
    ↓ git add -A && git commit && git push
    ↓ GitHub commit appears
Total time: ~15 seconds after your change
```

### What triggers a commit

| Trigger | What causes it |
|---------|---------------|
| `openclaw.json` written | Any `openclaw config set ...`, plugin enable/disable, or direct JSON edit |
| `agents/` changed | Agent profiles updated — model swaps, skill changes, routing rules |
| `identity/` changed | Agent identity files modified |

### What does NOT trigger a commit (excluded intentionally)

Logs, sessions, task queues, canvas files, discord/telegram message state, exec approvals, update checks — runtime noise, not config.

### Check sync status

```bash
# Is it running?
systemctl status openclaw-config-sync

# Watch it in real time
journalctl -u openclaw-config-sync -f

# See recent commits (proves it's working)
su - ubuntu -c 'cd ~/openclaw-workspaces && git log --oneline -5'

# Force a test commit (touch config → wait 15s → check log)
su - ubuntu -c 'touch ~/.openclaw/openclaw.json && sleep 15 && cd ~/openclaw-workspaces && git log --oneline -3'
```

---

## Editing Config Safely

```bash
# Preferred: use CLI where available
openclaw config set <key> <value>
openclaw plugins enable <name>
openclaw skills install <slug>

# Direct JSON edit (fallback)
cp ~/.openclaw/openclaw.json ~/.openclaw/openclaw.json.bak
nano ~/.openclaw/openclaw.json
systemctl restart openclaw

# Verify gateway is healthy
openclaw gateway status
journalctl -u openclaw -f
```

---

## Disaster Recovery

Any VM can be fully restored in ~5 minutes:

```bash
# On a fresh Debian/Ubuntu VPS:
export OP_SERVICE_ACCOUNT_TOKEN="<from 1Password SO vault>"
curl -sL https://raw.githubusercontent.com/OneMind-OS/openclaw-core-workspaces/main/bootstrap.sh | bash
# (replace -core- with -sentinel- for Sentinel VM)
```

The `bootstrap.sh` script:
1. Creates `ubuntu` user (already exists on Lightsail)
2. Installs node, op, gh, openclaw
3. Clones workspace repo → deploys `openclaw.json` + workspace files
4. Installs + enables all systemd services
5. Starts openclaw gateway

---

## Related

- [[OPENCLAW-LIVE-STATE]] — current production state (versions, agents, skills)
- [[OPENCLAW-ARCHITECTURE]] — architecture overview
- [[OPENCLAW-CONCURRENCY]] — concurrency model
- Workspace repos: `OneMind-OS/openclaw-core-workspaces` · `OneMind-OS/openclaw-sentinel-workspaces`
