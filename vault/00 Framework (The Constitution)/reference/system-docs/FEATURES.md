# OpenClaw Features — Live Status

> **Source of truth** for all active features, skills, plugins, channels, and configs running on the VPS.
> Updated: 2026-03-17 | OpenClaw v2026.3.13 | Build 61d171a
> VPS: `104.248.4.212` | User: `onemind` | Config: `/home/onemind/.openclaw/openclaw.json`

---

## Gateway Dashboard

| Item | Value |
|------|-------|
| WebSocket URL | `wss://gateway.onemindos.dev` |
| HTTP URL | `https://gateway.onemindos.dev` |
| Local | `http://127.0.0.1:18789` |
| Token | In 1Password → "OpenClaw Gateway" → `token` field |
| Tailscale IP (VPS) | `100.77.52.93` |
| Caddy proxy | `gateway.onemindos.dev → localhost:18789` |

**Dashboard Tips:**
- Use `wss://` not `ws://` — Caddy handles TLS
- If you get "rate_limited" (too many wrong tokens), restart the gateway: `systemctl restart openclaw`
- The token is ~64 chars hex — paste it exactly, no spaces
- `flush_interval -1` is active in Caddy — WebSocket should stay stable

---

## Agents (9 Total)

| Agent | ID | Model | Heartbeat | Discord Channel |
|-------|----|-------|-----------|----------------|
| ⚡ Legacy | `legacy` | `azure-openai/gpt-5.2` | 1h | `#legacy-chat` |
| ⚔️ Spartan | `spartan` | `azure-openai/gpt-4o` | 2h | `#spartan-chat` |
| 🔮 Oracle | `oracle` | `azure-openai/gpt-4o` | 3h | `#oracle-chat` |
| 🌳 Eden | `eden` | `azure-openai/gpt-4o` | 3h | `#eden-chat` |
| 👑 Empire | `empire` | `azure-openai/gpt-5.2` | 2h | `#empire-chat` |
| 🌿 Haven | `haven` | `azure-openai/gpt-4o` | 3h | `#haven-chat` |
| 🛡️ Guardian | `guardian` | `azure-openai/gpt-4o` | 2h | `#guardian-chat` |
| ⚙️ Grid | `grid` | `azure-openai/gpt-4o` | 3h | `#grid-chat` |
| 🔺 Trinity | `trinity` | `azure-openai/gpt-4o` | 3h | `#trinity-chat` |

**Subagent model (all agents):** `azure-openai/gpt-4o-mini`

---

## Model Providers (5 Active)

### Azure OpenAI ✅
> Base URL from `AZURE_OPENAI_BASE_URL`, Key from `AZURE_OPENAI_API_KEY`

| Alias | Model |
|-------|-------|
| `gpt5` | `gpt-5.2` |
| `gpt5-codex` | `gpt-5.2-codex` |
| `codex-mini` | `gpt-5.1-codex-mini` |
| `gpt5-mini` | `gpt-5-mini` |
| `gpt5-nano` | `gpt-5-nano` |
| `gpt4o` | `gpt-4o` |
| `gpt4o-mini` | `gpt-4o-mini` |
| `o1` | `o1` |
| `o3` | `o3` |
| `o3-mini` | `o3-mini` |
| `o4-mini` | `o4-mini` |
| `gpt41` | `gpt-4.1` |
| `gpt41-mini` | `gpt-4.1-mini` |
| `gpt41-nano` | `gpt-4.1-nano` |
| `gpt45` | `gpt-4.5` |

### Anthropic ✅ (Claude 4.5 only)
> Key from `ANTHROPIC_API_KEY`

| Alias | Model |
|-------|-------|
| `opus45` | `claude-opus-4-5-20250514` |
| `sonnet45` | `claude-sonnet-4-5-20250514` |
| `haiku45` | `claude-haiku-4-5-20250514` |

### OpenAI ✅
> Key from `OPENAI_API_KEY`

| Alias | Model |
|-------|-------|
| _(direct)_ | `gpt-4o`, `gpt-4o-mini`, `o3-mini`, `gpt-4.1`, `gpt-4.1-mini`, `gpt-4.1-nano`, `o4-mini` |

### Amazon Bedrock ✅
> AWS creds from 1Password "AWS Bedrock" item → `username`/`password`, region `us-east-1`

| Model ID |
|----------|
| `anthropic.claude-opus-4-5-v1:0` |
| `anthropic.claude-sonnet-4-5-v1:0` |
| `anthropic.claude-haiku-4-5-v1:0` |
| `anthropic.claude-3-5-sonnet-20241022-v2:0` |
| `anthropic.claude-3-5-haiku-20241022-v1:0` |
| `anthropic.claude-3-opus-20240229-v1:0` |
| `amazon.nova-pro-v1:0` |
| `amazon.nova-lite-v1:0` |
| `amazon.nova-micro-v1:0` |
| `meta.llama3-3-70b-instruct-v1:0` |
| `meta.llama4-scout-17b-instruct-v1:0` |

### Gradient (DigitalOcean) ✅
> Creds from `GRADIENT_BASE_URL` + `GRADIENT_API_KEY`

| Alias | Model |
|-------|-------|
| `llama` | `meta-llama/Llama-3.3-70B-Instruct` |
| `mistral` | `mistralai/Mistral-7B-Instruct-v0.3` |
| `qwen` | `Qwen/Qwen2.5-72B-Instruct` |

---

## Channels

| Channel | Status | Notes |
|---------|--------|-------|
| **Discord** | ✅ Active | Bot connected, 9 agent channels |
| **BlueBubbles (iMessage)** | ✅ Active | Connected to MacBook |
| **WhatsApp** | ❌ Not configured | Plugin disabled, would need `openclaw channels login` |
| **Telegram** | ❌ Not configured | Plugin disabled |
| **Signal** | ❌ Not configured | Plugin disabled |
| **Voice Call** | ❌ Not configured | Plugin disabled, needs Twilio creds |

---

## Plugins (Loaded)

| Plugin | ID | Version | Status | Purpose |
|--------|----|---------|--------|---------|
| ClawVault | `clawvault` | 3.5.0 | ✅ Loaded | Structured memory system, proactive recall — **memory slot** |
| Lossless Context | `lossless-claw` | 0.3.0 | ✅ Loaded | DAG-based conversation summarization — **context engine slot** |
| Lobster | `lobster` | bundled | ✅ Loaded | Named automation routines |
| Diffs | `diffs` | bundled | ✅ Loaded | Shareable diff tool |
| BlueBubbles | `bluebubbles` | bundled | ✅ Loaded | iMessage channel |
| LLM Task | `llm-task` | bundled | ✅ Loaded | Background LLM subtasks |
| OpenProse | `open-prose` | bundled | ✅ Loaded | `/prose` slash command |
| Diagnostics OTEL | `diagnostics-otel` | bundled | ✅ Loaded | OpenTelemetry traces |
| ACPX | `acpx` | bundled | ✅ Loaded | ACP router extension |
| Memory Core | `memory-core` | bundled | — | Disabled (clawvault is memory slot) |
| Memory LanceDB | `memory-lancedb` | bundled | — | Disabled (clawvault is memory slot) |

**Plugin Slots:**
- `memory` → `clawvault`
- `contextEngine` → `lossless-claw` ← activated 2026-03-17

**Security:** `plugins.allow: ["clawvault", "lossless-claw"]` ← set 2026-03-17

---

## Skills (28 Ready / 68 Total)

### Ready — Bundled (`openclaw-bundled`)
| Skill | Description |
|-------|-------------|
| `1password` | 1Password CLI (`op`) — install and use |
| `blucli` | BluOS CLI for music playback control |
| `bluebubbles` | Send/manage iMessages via BlueBubbles |
| `camsnap` | Capture frames/clips from RTSP/ONVIF cameras |
| `discord` | Discord ops via message tool |
| `healthcheck` | Host security hardening and risk config |
| `node-connect` | Diagnose OpenClaw node connection/pairing |
| `skill-creator` | Create/edit/audit AgentSkills |
| `tmux` | Remote-control tmux sessions |

### Ready — Extra (`openclaw-extra`)
| Skill | Description |
|-------|-------------|
| `acp-router` | Route LLM requests across Pi, Claude Code, Codex |
| `diffs` | Produce shareable diffs |
| `prose` | OpenProse VM skill pack + `/prose` command |

### Ready — Workspace
| Skill | Description |
|-------|-------------|
| `github` | GitHub API with managed OAuth |
| `obsidian` | Obsidian vault operations (plain Markdown) |
| `weather` | Current weather & forecasts (no API key) |
| `Agent Browser` | Rust-based headless browser automation |
| `api-gateway` | 100+ API connectors (Google, Microsoft, etc.) |
| `cloudflare-api` | Cloudflare DNS & zone management |
| `defuddle` | Extract clean markdown from web pages |
| `json-canvas` | Create/edit JSON Canvas `.canvas` files |
| `nats-manager` | NATS server monitoring + JetStream management |
| `obsidian-bases` | Create/edit Obsidian Bases `.base` files |
| `obsidian-cli` | Obsidian CLI vault interaction |
| `obsidian-markdown` | Obsidian flavored markdown (wikilinks, callouts) |
| `proactive-agent` | Transform agents into proactive actors |
| `rss-digest` | Fetch + digest RSS/Substack feeds |
| `skill-forge` | Autonomous OpenClaw skill builder |
| `todoist` | Manage Todoist tasks and projects |

### Not Ready (40 skills) — missing bins or platform gates
> Many skills need specific binaries on the Linux VPS (e.g., macOS-only tools, or `uv`, `brew`, `ruby`, etc.).
> Run `openclaw skills list` on VPS for full view.

---

## Lobster Routines (9)

| Routine | File |
|---------|------|
| Morning Briefing | `morning-brief.lob` |
| Habit Check | `habit-checker.lob` |
| Health Check | `health-check.lob` |
| Tech Digest | `tech-digest.lob` |
| Second Brain | `second-brain.lob` |
| Earnings Tracker | `earnings-tracker.lob` |
| Git Sync | `git-sync.lob` |
| Project State | `project-state.lob` |
| Update OpenClaw | `update-openclaw.lob` |

Run with: `lobster run <name>` (from Legacy) or `openclaw run <name>` (CLI)

---

## Voice / TTS

| Feature | Status |
|---------|--------|
| TTS Provider | OpenAI (via `OPENAI_API_KEY` / Azure fallback) |
| Voice Auto-Join | ❌ Disabled (empty array — on-demand only) |
| Voice Call Plugin | ❌ Not installed (`@openclaw/voice-call` needs Twilio creds) |
| LiveKit | ✅ Running in Docker (`LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`, `LIVEKIT_URL`) |

---

## Memory System

| Layer | Plugin | Status |
|-------|--------|--------|
| **Primary** | ClawVault 3.5.0 | ✅ Active (memory slot) |
| **Context Engine** | Lossless Claw 0.3.0 | ✅ Active (context engine slot) — set 2026-03-17 |
| Fallback | memory-core | Disabled (overridden by ClawVault) |
| DB | `/home/onemind/.openclaw/lcm.db` | SQLite, threshold 0.75 |
| Codex memory search | `/home/onemind/onemind-os/codex` | ✅ Indexed, hybrid BM25+vector, MMR |

---

## Infrastructure Services (Docker)

| Service | Port | Status |
|---------|------|--------|
| NATS JetStream | 4222 / 8222 | ✅ Running |
| LiveKit | 7880 | ✅ Running |
| n8n | 5678 | ✅ Running |
| OmniDB | varies | Running |
| legacy-bridge (NATS↔OpenClaw) | internal | Service |

---

## Systemd Services

| Service | Status |
|---------|--------|
| `openclaw.service` | ✅ Active (via `op run`) |
| `caddy.service` | ✅ Active (WebSocket fix live) |
| `onemind-watch.service` | ⏸️ **Disabled** (paused 2026-03-17) |
| `onemind-pull.timer` | ✅ Active (pulls GitHub → VPS every 1 min) |
| `openclaw-reload.service` | ✅ Active |

> **Note:** `onemind-watch` was disabled because it was hitting detached HEAD errors. Re-enable with `systemctl enable --now onemind-watch` when repo is on a clean branch.

---

## Secrets (1Password)

All secrets in vault `aaeyris5z44pdqspk277u2oppi` ("00-24 SO (Sovereign Operations)").

| Secret | 1Password Item |
|--------|----------------|
| Azure OpenAI API Key | "Azure OpenAI" |
| Anthropic API Key | "Anthropic" |
| OpenAI API Key | "OpenAI" |
| AWS Bedrock Key | "AWS Bedrock" (username/password fields) |
| Discord Bot Token | "Discord Bot" |
| BlueBubbles Password | "BlueBubbles" |
| Gateway Token | "OpenClaw Gateway" |
| Hooks Token | "OpenClaw Hooks" |
| LiveKit credentials | "LiveKit" |

---

## Known Issues / TODO

| # | Issue | Priority | Fix |
|---|-------|----------|-----|
| 1 | `onemind-watch` disabled | Medium | Fix detached-HEAD state on VPS repo, then re-enable |
| 2 | WhatsApp not connected | Low | Run `openclaw channels login` to pair WhatsApp if needed |
| 3 | 40/68 skills not ready | Low | Most need macOS-only binaries — acceptable on Linux VPS |
| 4 | `clawvault` async warning on load | Cosmetic | Waiting for OpenClaw runtime update to support async plugin registration |
| 5 | `lcm_grep` `no such table: conversations` | Low | LCM DB may need migration — check with `openclaw doctor` after restart |
| 6 | MacBook node probe timeout | Low | Node offline or sleep — reconnects automatically when awake |
| 7 | NATS JetStream fallbacks | Low | Subjects `comms.>`, `home.>` etc. falling back — NATS stream config |

---

## Quick Reference Commands

```bash
# VPS access
ssh root@104.248.4.212

# OpenClaw (run as onemind)
sudo -u onemind openclaw status
sudo -u onemind openclaw skills list
sudo -u onemind openclaw plugins list
sudo -u onemind openclaw doctor

# Services
systemctl status openclaw caddy onemind-pull onemind-watch
systemctl restart openclaw   # also clears rate-limit lockout
journalctl -u openclaw -f

# Gateway dashboard
# URL: wss://gateway.onemindos.dev
# Token: op run --env-file=/opt/onemind/.env.tpl -- bash -c 'echo $OPENCLAW_GATEWAY_TOKEN'

# Update OpenClaw
sudo -u onemind openclaw update
```


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
