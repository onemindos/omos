---
title: "OpenClaw Live State"
type: reference
status: active
created: 2026-04-16
modified: 2026-05-10
owner: zeus
tags: [type/reference, domain/SO, domain/03-agents]
ai_summary: "Current production state of the OpenClaw gateway — agents, skills, MCP servers, sync services, and pending work"
---

# OpenClaw Live State

> Last verified: 2026-05-10
> Auto-synced from VPS via onemind-watch → GitHub → onemind-pull

---

## Gateway

### Nodes

| Node | Hostname | IP / Tailscale | Notes |
|------|----------|----------------|-------|
| onemindos-cloud-1 | omos-core | root@64.227.8.23 · https://omos-core.tail717ba0.ts.net | Primary OpenClaw node |
| Titan | omos-titan | — | Workhorse VM — Services TBD, no OpenClaw node yet |

### Gateway Config

| Field | Value |
|-------|-------|
| OpenClaw version | v2026.5.2 |
| Port | 18789 (loopback + Tailscale serve) |
| Config file | /root/.openclaw/openclaw.json |
| Service | systemd USER service (openclaw-gateway.service) |
| chatCompletions API | Enabled — /v1/chat/completions |
| Model (all agents) | github-copilot/claude-sonnet-4.6 |
| Gateway token | 1Password SO vault → "OpenClaw Gateway (onemindos-cloud)" |
| Workspace repos | [openclaw-core-workspaces](https://github.com/OneMind-OS/openclaw-core-workspaces) · [openclaw-sentinel-workspaces](https://github.com/OneMind-OS/openclaw-sentinel-workspaces) |

---

## Agents (9 Active)

| Agent | ID | Domain | Workspace |
|-------|----|--------|-----------|
| Legacy ⚡ | main | SO (all) | ~/.openclaw/workspace/ |
| Spartan 💪 | spartan | HP 27-28 | ~/.openclaw/workspace-spartan/ |
| Oracle 🔮 | oracle | HP 25-26, 29-34 | ~/.openclaw/workspace-oracle/ |
| Eden 🌱 | eden | LE (environment) | ~/.openclaw/workspace-eden/ |
| Haven 🏡 | haven | LE 50-74 | ~/.openclaw/workspace-haven/ |
| Empire 💰 | empire | GE 75-99 | ~/.openclaw/workspace-empire/ |
| Guardian 🛡️ | guardian | SO 11 | ~/.openclaw/workspace-guardian/ |
| Grid ⚙️ | grid | SO 01, 09 | ~/.openclaw/workspace-grid/ |
| Trinity 🔺 | trinity | SO 03 | ~/.openclaw/workspace-trinity/ |

Each workspace contains: `SOUL.md`, `AGENTS.md`, `USER.md`, `IDENTITY.md`, `MEMORY.md`

---

## Skills — 18 Ready

```
Ready (18):
  diffs, obsidian-vault-maintainer, wiki-maintainer, prose
  tavily, 1password, clawhub, coding-agent
  gh-issues, github, healthcheck, node-connect
  session-logs, skill-creator, taskflow, taskflow-inbox-triage
  tmux, weather

Missing requirements (intentionally skipped):
  sag       — ElevenLabs TTS (not needed)
  gog       — Google Workspace CLI (not needed)
  obsidian  — Mac-only obsidian-cli (low priority, needs Mac node)
  blogwatcher / openai-whisper / summarize — not priority
```

### Key Skills

| Skill | What It Does | Requirements |
|-------|-------------|--------------|
| `github` | GitHub ops via gh CLI | gh v2.89.0 (installed) |
| `gh-issues` | Fetch issues, spawn agents to fix, open PRs | gh CLI |
| `healthcheck` | Security hardening + risk assessment | none |
| `session-logs` | Search past conversation history | ripgrep 14.1.0 (installed) |
| `coding-agent` | Delegate to Claude Code | claude v2.1.110 (installed) |
| `1password` | 1Password CLI (op) operations | op CLI (installed) |
| `clawhub` | Search/install skills from clawhub.com | clawhub v0.9.0 (installed) |
| `obsidian-vault-maintainer` | Maintain Obsidian-compatible wiki | none |
| `taskflow` | Internal OpenClaw task flows | none |
| `tmux` | Remote-control tmux sessions | tmux |
| `tavily` | Web search and content extraction | none |
| `coding-agent` | Spawn Claude Code for dev tasks | claude CLI |

### Using clawhub

```bash
# Search for a skill
clawhub search <query>

# Inspect before installing
clawhub inspect <slug>

# Install a skill
clawhub install <slug>

# Update installed skills
clawhub update

# List installed
clawhub list
```

---

## MCP Servers (3 Active)

| Name | Command | Purpose |
|------|---------|---------|
| `context7` | `npx -y @upstash/context7-mcp@latest` | Up-to-date library docs for coding tasks |
| `github` | `npx -y @modelcontextprotocol/server-github` | Full GitHub API (authed as Zeus-Delacruz) |
| `codex-vault` | `npx -y @modelcontextprotocol/server-filesystem /root/onemind-codex` | Read/write vault files |

---

## Vault Sync Architecture

```
Mac (Obsidian / git commit / git push)
    ↓ push
GitHub (OneMind-OS/onemind-codex) ← SOURCE OF TRUTH
    ↓ onemind-pull.timer (every 10s)
VPS /root/onemind-codex/
    ↓ onemind-watch.service (inotifywait, 3s debounce)
GitHub (auto git add/commit/push on any file change)
```

### Services

| Service | Type | Status | Script |
|---------|------|--------|--------|
| `onemind-watch.service` | systemd simple | active ✓ | /usr/local/bin/onemind-watch.sh |
| `onemind-pull.service` | systemd oneshot | (runs via timer) | /usr/local/bin/onemind-pull.sh |
| `onemind-pull.timer` | systemd timer | active ✓ | every 10s |

### Logs

```bash
# Watch sync in real-time
tail -f /var/log/onemind-watch.log
tail -f /var/log/onemind-pull.log

# Service status
systemctl status onemind-watch onemind-pull.timer
```

### Manual operations on VPS

```bash
# Write a file then commit manually
cd /root/onemind-codex
echo "test" >> some-file.md
git add -A && git commit -m "agent: description" && git push

# Force pull from GitHub
/usr/local/bin/onemind-pull.sh
```

---

## Config Sync

Changes to `~/.openclaw/` on each VM are auto-committed to the workspace repos via `openclaw-config-sync.service`.

| VM | Service | Repo |
|----|---------|------|
| omos-core | `openclaw-config-sync.service` | `OneMind-OS/openclaw-core-workspaces` |
| omos-sentinel | `openclaw-config-sync.service` | `OneMind-OS/openclaw-sentinel-workspaces` |

```bash
# Check sync service status
systemctl status openclaw-config-sync
journalctl -u openclaw-config-sync -f

# Manual sync trigger (touch any file in .openclaw/)
touch ~/.openclaw/openclaw.json
```

---

## Installed Tools (VPS)

| Tool | Version | Purpose |
|------|---------|---------|
| openclaw | v2026.4.14 | AI gateway |
| node.js | v20+ | Runtime |
| gh | v2.89.0 | GitHub CLI (authed as Zeus-Delacruz) |
| claude | v2.1.110 | Claude Code CLI |
| rg (ripgrep) | 14.1.0 | Fast search (session-logs skill) |
| inotifywait | — | Filesystem watcher (onemind-watch) |
| clawhub | v0.9.0 | ClawHub CLI (skill marketplace) |
| op | 2.x | 1Password CLI |

---

## Channels (Pending Configuration)

| Channel | Status | What's Needed |
|---------|--------|---------------|
| Telegram | Plugin enabled, no token | Create bot via @BotFather → set token in openclaw.json |
| Discord | Plugin enabled, no token | Get token from SO vault "discord bot token (Sentinel)" → set in openclaw.json |

### Telegram setup

```bash
# 1. Create bot via @BotFather on Telegram
# 2. Get bot token (format: 123456789:ABCdef...)
# 3. Find your Telegram user ID (message @userinfobot)
# 4. Set in config:
openclaw config set channels.telegram.accounts.default.token "BOT_TOKEN"
openclaw config set channels.telegram.dmPolicy allowlist
openclaw config set channels.telegram.allowFrom '["YOUR_TELEGRAM_USER_ID"]'
openclaw gateway restart
```

### Discord setup

```bash
# 1. Get token from 1Password SO vault → "discord bot token (Sentinel)"
# 2. Get your Discord user ID (Developer Mode → right-click → Copy User ID)
# 3. Set in config:
openclaw config set channels.discord.accounts.default.token "BOT_TOKEN"
openclaw config set channels.discord.dmPolicy allowlist
openclaw config set channels.discord.allowFrom '["YOUR_DISCORD_USER_ID"]'
openclaw gateway restart
```

---

## Pending Work

### High Priority

- [x] **Telegram** — ✅ wired (token in openclaw.json)
- [x] **Discord** — ✅ wired (token in openclaw.json)
- [ ] **Mac node** — install openclaw node on Mac for screen/shell access from VPS Legacy

### Lower Priority

- [ ] **Obsidian skill** — install `obsidian-cli` via brew on Mac node (Mac-only skill)
- [ ] **HACS integration** — OpenClaw HACS plugin on Beelink Home Assistant
- [ ] **Dreaming** — enable memory dreaming (runs at 3 AM, promotes session content → MEMORY.md)
- [ ] **Agent routing rules** — bind specific agents to Telegram/Discord groups
- [ ] **Obsidian Git plugin on Mac** — auto-commit vault edits to GitHub (completes bidirectional loop)

### Dropped (Not Needed)

- ~~sag (ElevenLabs TTS)~~ — not needed
- ~~gog (Google Workspace)~~ — not needed

---

## Config Editing Safety Protocol

```bash
# ALWAYS backup before editing openclaw.json
cp ~/.openclaw/openclaw.json ~/.openclaw/openclaw.json.bak

# Use CLI where available (recommended)
openclaw config set <key> <value>
openclaw plugins enable <name>
openclaw skills <subcommand>

# Edit JSON directly only as fallback
nano ~/.openclaw/openclaw.json

# Restart to apply changes
systemctl --user restart openclaw-gateway

# Verify
openclaw gateway status
```

---

## Quick Reference

```bash
# Chat with agents
openclaw agent --message "..."                          # Legacy (default)
openclaw agent --agent empire --message "..."           # specific agent
openclaw tui                                            # interactive TUI

# Skills
openclaw skills check                                   # 18 ready / 6 missing
openclaw skills list                                    # full table
clawhub search <query>                                  # search marketplace
clawhub install <slug>                                  # install skill

# Memory
openclaw memory promote --apply                         # promote sessions → MEMORY.md

# Gateway
openclaw gateway status
systemctl --user status openclaw-gateway
journalctl --user -u openclaw-gateway -f               # live logs

# Vault sync
systemctl status onemind-watch onemind-pull.timer
tail -f /var/log/onemind-watch.log
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
