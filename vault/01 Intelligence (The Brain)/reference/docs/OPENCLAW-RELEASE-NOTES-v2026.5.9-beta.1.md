---
title: OpenClaw v2026.5.9-beta.1 + v2026.5.10-beta.2 — Release Notes & Architecture Analysis
domain: 01
tags: [openclaw, release-notes, architecture, sqlite-migration, agents]
status: active
updated: 2026-05-10
---

# OpenClaw v2026.5.9-beta.1 — Release Notes & Architecture Analysis

Released: 2026-05-09 (pre-release)
Reviewed: 2026-05-10

---

## 1. What's New — Priority Items for This Stack

### Model & Provider Updates

| Change | Impact |
|--------|--------|
| **GitHub Copilot catalog refreshes live** from `/models` API | Per-account entitlements + accurate context windows at runtime. Static catalog (now includes `gpt-5.5`) is the fallback. |
| **Bedrock `serviceTier` support** | Set `agents.defaults.params.serviceTier` to `default`, `flex`, `priority`, or `reserved` for Bedrock models. |
| **Bedrock credential refresh** | Long-running Gateway processes now pick up renewed AWS profile credentials without restart. Was a known workaround. |
| **Google/Gemini 3 Pro → 3.1 auto-migration** | Retired `gemini-3-pro-preview` normalized to `gemini-3.1-pro-preview` everywhere. |
| **Mistral Medium 3.5** | Added to catalog with reasoning support. |

### Agent Behavior

- **Agents now inject their own model identity into system prompts** — agents can accurately answer "what model are you running?"
- **`/think default` and `/fast default`** — clears session thinking overrides, returns to configured defaults
- **Compaction bug fixed** — consecutive compaction cycles now run correctly instead of alternating (long-standing bug)
- **Sub-agent spawn fix** — `model: "default"` now resolves correctly
- **Session repair** — null-role JSONL corruption cleaned up on startup automatically

### Channels

**Telegram:**
- Draft previews continue in a new message when they hit Telegram's 4096-char limit
- Polling restart watermark persisted properly
- Reply chain hydration fixed — agents now see full nested reply context

**Discord:**
- Progress draft previews default to ON — tool activity shows in one edited message
- Full STT/TTS + realtime voice modes added to voice channels
- `agent-proxy` is now the default voice mode

**iMessage:**
- **Offline catchup replay** added — replays messages received while Gateway was down
  ```json
  channels.imessage.catchup.enabled: true
  // tunables: maxAgeMinutes, perRunLimit, firstRunLookbackMinutes
  ```
- **BlueBubbles channel removed** (BREAKING — see below)

**Slack:**
- Thread routing improved — channel turns in implicit-conversation channels now route to thread sessions correctly

### Performance

- **Gateway plugin metadata reuse** — no longer rescanned before every provider call (major reduction in turn overhead)
- **Session list fast-path** — `openclaw sessions` and Control UI no longer do heavyweight model resolution per row
- **Control UI decoupled** — chat/config refreshes no longer block on slow history/schema work

### Security Fixes

| Fix | Severity |
|-----|----------|
| **Memory Wiki** (GHSA-72fw-cqh5-f324) — sandboxed subagents could read parent/sibling session transcripts via `wiki_search`/`wiki_get` | High |
| **Exec `argPattern`** was silently not enforced on Linux/macOS — now fixed | Medium |
| **`trusted-proxy` mode** now fails closed if identity checks fail (was accepting local password fallback) | Medium |
| **`openclaw.json` permissions** restored to `0600` after recovery | Medium |
| **Root execution refused by default** — `openclaw` now refuses to run as root without `OPENCLAW_ALLOW_ROOT=1` | Medium |

### Operations & CLI

- **`openclaw gateway restart --safe --skip-deferral`** — bypass the safe-restart deferral gate when a pinned task blocks restart
- **`gateway.tailscale.preserveFunnel: true`** — skip re-applying `tailscale serve` on startup if an external Funnel already covers the port (relevant to this setup)
- **`openclaw cron list --agent <id>`** — filter cron jobs by specific agent
- **`openclaw sessions cleanup --fix-dm-scope`** — clean up stale DM session rows
- **Node 22 minimum raised to `22.16+`**, Node 24 recommended

---

## 2. Breaking Change

**BlueBubbles channel removed.**
Any agent config referencing `channels.bluebubbles` must migrate to `channels.imessage` using the `imsg` bridge on a signed-in Mac or SSH wrapper.

---

## 3. Action Items for This Stack

### Immediate

- [ ] **Root user** — if OpenClaw runs as root on `onemindos-core`, add `OPENCLAW_ALLOW_ROOT=1` to the systemd unit env, or migrate to non-root user
- [ ] **Tailscale Funnel preserve** — add `gateway.tailscale.preserveFunnel: true` to `openclaw.json` to stop OpenClaw re-applying serve on every restart
- [ ] **Run `openclaw update`** on both `onemindos-core` and `omos-sentinel` nodes

### Before the SQLite Migration Lands (see Section 4)

- [ ] **Run `openclaw backup create`** on both nodes to capture a clean JSONL snapshot before migration auto-runs
- [ ] Verify backup location and test restore on a non-critical agent

### Optional Enhancements

- [ ] **Bedrock `serviceTier`** — set `agents.defaults.params.serviceTier: "priority"` for faster Bedrock response lanes
- [ ] **iMessage catchup** — enable `channels.imessage.catchup.enabled: true` to catch messages missed during downtime/restarts

---

## 4. The SQLite Runtime State Migration (Coming Next)

> PR #78595 — "Refactor runtime state into SQLite" by @steipete  
> Status: Active development on `main`. 2k+ files changed, 298 commits. Not in beta.1.

### The Problem It Solves

Currently OpenClaw state is scattered across:

```
~/.openclaw/openclaw.json          ← config (single file, can get clobbered)
~/.openclaw/agents/*/sessions/*.jsonl ← one file per conversation (thousands of files)
sidecar SQLite                     ← memory search only
lock files                         ← can get stuck and block restarts
pruning state                      ← the alternating compaction bug lived here
```

Every gateway crash, network hiccup, or mid-write interruption can leave these files out of sync. The session repair code, truncation-repair paths, and the null-role JSONL cleanup in this beta are all band-aids on this architecture.

### What Changes

```
Current:  scattered JSON + JSONL + sidecar SQLite + lock files
Future:   single ~/.openclaw/openclaw.db (SQLite, WAL mode)
```

- ACID transactions — no partial writes, no corrupt JSONL mid-turn
- No lock files — SQLite handles concurrent access natively
- The entire class of session repair / truncation-repair bugs stops existing
- One-file backup: `cp ~/.openclaw/openclaw.db /backup/` — done

### What It Unlocks

#### Direct Database Queries (No Gateway Required)

```sql
-- All agent activity this week
SELECT agent_id, COUNT(*) as turns, SUM(tokens) as total_tokens
FROM sessions JOIN turns USING (session_id)
WHERE created_at > datetime('now', '-7 days')
GROUP BY agent_id;

-- Find every time Spartan discussed a specific topic
SELECT agent_id, content, created_at FROM turns
WHERE agent_id = 'spartan' AND content LIKE '%deadlift%';

-- Cross-pillar intelligence: which domains got the most agent attention?
SELECT agent_id, AVG(tokens) as avg_tokens_per_turn FROM turns
GROUP BY agent_id ORDER BY avg_tokens_per_turn DESC;
```

SQLite WAL mode allows reads while the Gateway is writing — dashboards can query live without blocking agent turns.

#### Custom Dashboard Capabilities (Post-Migration)

```
Current:  Dashboard → WebSocket → Gateway (must be running) → JSONL files
Future:   Dashboard → SQLite DB (read-only, always available)
         + Dashboard → WebSocket → Gateway (for live streaming turns)
```

**What this enables for the One Mind stack:**

| Feature | How |
|---------|-----|
| **Live agent activity feed** | Poll SQLite WAL every 1s — no WebSocket needed |
| **Agent health monitor** | Query `last_active` per agent; detect silent/stuck agents without relying on Gateway health checks |
| **Unified pillar view** | All 9 agents' history in one query — SO/HP/LE/GE activity in one dashboard panel |
| **Cross-agent intelligence** | "What decisions did all agents make about GE ventures this week?" |
| **Offline analytics** | Pull DB locally, run in DBeaver, Datasette, or Python notebooks |
| **Vault auto-sync** | Script reads SQLite → writes weekly summary notes to `_codex/memory/` automatically |

#### Vault Integration Example (Future)

```python
import sqlite3
from datetime import datetime, timedelta

conn = sqlite3.connect("/root/.openclaw/openclaw.db")
week_ago = (datetime.now() - timedelta(days=7)).isoformat()

# Get all agent summaries for the weekly vault note
rows = conn.execute("""
    SELECT agent_id, content FROM turns
    WHERE role = 'assistant' AND created_at > ?
    ORDER BY agent_id, created_at
""", (week_ago,)).fetchall()

# Write to vault
with open("/root/onemind-codex/_codex/memory/work/weekly-agent-summary.md", "w") as f:
    f.write(f"# Weekly Agent Summary — {datetime.now().strftime('%Y-W%W')}\n\n")
    for agent_id, content in rows:
        f.write(f"## {agent_id}\n{content[:500]}...\n\n")
```

### Timeline Estimate

| Milestone | Estimated |
|-----------|-----------|
| PR lands in `main` | 2-4 weeks |
| Available in beta | v2026.6.x |
| Stable release | v2026.6.x or v2026.7.x |
| Auto-migration on update | Included (JSONL → SQLite) |

**Critical:** Run `openclaw backup create` on both nodes **before** updating when this lands. The auto-migration is one-way.

---

## 5. Codex Harness as Default for OpenAI Models

> Announced by @pash in OpenClaw Discord, May 2026.

OpenAI models in OpenClaw are moving to the **Codex harness by default**.

- The old `openai-codex/*` provider path is removed
- OpenAI/GPT models now run through OpenAI's Codex app-server natively
- Legacy `openai-codex/*` config keys auto-migrate to `openai/*` via `openclaw doctor --fix`
- Enables native guardian mode vs YOLO mode with Codex, GPT-5.5 out of the box

**Impact on this stack:** Low — primary model is `github-copilot/claude-sonnet-4.6` with Bedrock fallback. No OpenAI-primary agents configured. If OpenAI models are added later, routing has changed.

---

## 6. Reference Links

- Release beta.1: https://github.com/openclaw/openclaw/releases/tag/v2026.5.9-beta.1
- Release beta.2: https://github.com/openclaw/openclaw/releases/tag/v2026.5.10-beta.2
- SQLite PR: https://github.com/openclaw/openclaw/pull/78595
- Security advisory: GHSA-72fw-cqh5-f324 (Memory Wiki session visibility)

---

## 7. v2026.5.10-beta.2 — What's New (Released 2026-05-10)

Dropped the same day as this review. Smaller release focused on stability, model routing fixes, and Discord voice polish.

### Changes

**Voice & Realtime**
- `talk.realtime.instructions` config — append custom voice style instructions while preserving OpenClaw's built-in agent-consult guidance
- Discord voice realtime diagnostics added — speaker turns, playback resets, barge-in detection, audio cutoff analysis
- Discord voice now defaults to pure-JS `opusscript` decoder (avoids slow native addon compiles); opt-in native `@discordjs/opus` available for performance lanes
- OpenAI realtime voice now correctly handles Codex-compatible legacy audio/transcript event aliases

**Skills & Plugins**
- Opt-in private skill archive upload install path added, gated by `skills.install.allowUploadedArchives` — lets trusted Gateway clients install zip-backed skills when explicitly enabled by operator
- Codex dynamic-tools profile removed — Codex app-server now always owns workspace/edit/patch/exec/process/plan tools

**CLI & Onboarding**
- Setup, onboarding, configure, and channel commands now explain the next useful command in terminal flows instead of terse labels
- `config set` and `config patch` now correctly persist values that equal runtime defaults (was silently dropping them)
- Strict JSON parse failures now show a valid example and plain-string escape hatch

**Models**
- `provider/*` wildcard entries now supported in `agents.defaults.models` — `/model`, `/models`, and model pickers can show dynamically discovered models for a provider without exact allowlists
- Anthropic: `claude-haiku-4-5` added to API-key agent allowlist defaults
- Anthropic: Claude Opus/Sonnet 4 now correctly reports 1M session context (was showing 200k from stale local config)
- xAI/Grok: `/think low|medium|high` now exposed for reasoning-capable Grok models
- OpenRouter: `403 API key budget limit exceeded` now treated as billing error — model fallback advances instead of retrying
- Kimi K2.5 added to bundled OpenRouter catalog
- `moonshotai/...` and `moonshot-ai/...` refs now accepted as aliases for canonical `moonshot/...`
- Ollama: DeepSeek V4 cloud models keep thinking capability even when `/api/show` omits it — `/think high` no longer rejected

**GitHub Copilot Fix (Important for This Stack)**
- Copilot API tokens now minted with `vscode-chat` integration identity — fixes image-capable Copilot models inheriting the wrong `copilot-language-server` scope. Stale cached tokens auto-refreshed. Fixes #79946.

### Fixes Relevant to This Stack

| Fix | Why It Matters |
|-----|----------------|
| **Gateway config reread after SIGUSR1 restart** — was reusing stale startup snapshot, dropping config written after boot | Affects every gateway restart on omos-core |
| **LLM idle watchdog now applies during stream setup** — silent pre-stream model hangs no longer wait for full agent timeout | Reduces silent stuck turns |
| **Cron isolated self-cleanup** — can now inspect own job run history | Agent cron health checks work correctly |
| **ACPX runtime startup probe** — gateway `ready` signal no longer fires before ACP runtime is usable | Prevents race on startup |
| **Telegram: no-response DM turns stay quiet** — was rewriting them into visible silent-reply chatter | Cleaner Telegram behavior |
| **Telegram duplicate reply fix** — canonical final answer delivered instead of accumulated pre-tool text blocks | Fixes #79621 and #79986 |
| **Telegram draft preview debounce** — edited drafts no longer duplicate earlier text when providers emit incremental output | Cleaner streaming |
| **Agents: blank model-emitted tool calls stopped before dispatch** — fixes Kimi/NVIDIA blank-name retry loops | Prevents stuck turns |
| **ACP error chain preserved** — actual upstream failure shown in Gateway logs/Telegram replies instead of opaque `Internal error` | Easier debugging |
| **Security audit test isolation** — test suite was silently truncating real `~/.openclaw/exec-approvals.json` on the host running tests | Important if running tests locally |

### Dependencies Updated

`@agentclientprotocol/claude-agent-acp 0.33.1`, Codex ACP `0.14.0`, Baileys `7.0.0-rc10`, Google GenAI `2.0.1`, OpenAI `6.37.0`, AWS SDK `3.1045.0`, Kysely `0.29.0`

### Action Items from beta.2

- [ ] **Verify Copilot image models work correctly** after the `vscode-chat` token identity fix (affects image-capable Copilot routes)
- [ ] **Gateway restart behavior improved** — SIGUSR1 restarts now reread config from disk; no action needed but confirms restart stability
- [ ] **No breaking changes** in beta.2
