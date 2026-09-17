---
title: "Agent Skills Architecture"
domain: SO
subdomain: "01 Intelligence (AI, Memory & Tools)"
type: reference
status: active
created: 2026-04-18
updated: 2026-04-18
owner: zeus
tags: [domain/so, type/reference, status/active]
ai_summary: "Architecture reference for the OpenClaw agent skills system — how skills are structured, loaded, and invoked across all 9 agents."
---
# Agent Skills Architecture
**Maintained by:** ⚡ Legacy
**Last updated:** 2026-04-18

---

## Overview

This document defines how skills are structured, scoped, and assigned across all 9 One Mind OS agents running on OpenClaw. It is the canonical reference for tool access decisions.

---

## Skill Tiers

OpenClaw has three distinct skill tiers. Each agent's effective toolset is the union of all three.

### Tier 1 — Built-in System Skills
- **Location:** `/usr/lib/node_modules/openclaw/skills/`
- **Controlled by:** `~/.openclaw/openclaw.json` → `skills.entries`
- **Scope:** Global — enabled state applies to ALL agents simultaneously
- **Examples:** `coding-agent`, `github`, `tmux`, `weather`, `obsidian`, `healthcheck`

Toggle a skill on/off in `openclaw.json` to affect every agent at once.

### Tier 2 — Workspace Skills
- **Location:** `~/.openclaw/workspace[-agent]/skills/`
- **Controlled by:** presence in that specific agent's workspace
- **Scope:** Agent-specific — only the agent whose workspace contains it can see it
- **Examples:** `todoist-api-rest` (symlinked to all), `1password` (Legacy only), `postiz` (Legacy + Empire)

Use workspace skills for:
- Custom-built integrations
- Skills restricted to specific agents
- Agent-specific configurations

### Tier 3 — Extension Skills
- **Location:** `/usr/lib/node_modules/openclaw/dist/extensions/*/skills/`
- **Controlled by:** parent plugin enabled in `openclaw.json`
- **Scope:** Global when plugin is on
- **Examples:** `tavily`, `diffs`, `wiki-maintainer`, `obsidian-vault-maintainer`

These activate automatically when their parent plugin is enabled. No separate toggle needed.

---

## How Skills Stack Per Agent

```
Agent's effective toolset =
  Tier 1 (globally enabled built-ins)
  + Tier 2 (skills in their workspace/skills/ folder)
  + Tier 3 (extension skills from enabled plugins)
```

---

## Currently Enabled System Skills (Global — All Agents)

| Skill | Purpose |
|-------|---------|
| `coding-agent` | Spawn Claude Code / Codex for coding tasks |
| `github` | GitHub ops via `gh` CLI |
| `gh-issues` | Fetch issues, spawn fix agents, monitor PRs |
| `healthcheck` | Security audits, VPS hardening, risk checks |
| `obsidian` | Obsidian vault integration (Mac node) |
| `session-logs` | Search past conversation history |
| `weather` | Weather and forecasts |
| `tmux` | Remote-control tmux sessions |
| `clawhub` | Install/update skills from clawhub.com |
| `node-connect` | Diagnose node connection issues |
| `blogwatcher` | Monitor blogs/RSS |
| `openai-whisper` | Audio transcription |
| `sag` | TTS (ElevenLabs) |
| `gog` | Google Workspace access |

---

## Workspace Skills — Agent-Specific Assignments

### todoist-api-rest (All 9 agents — symlinked)
Full Todoist integration: REST API v1 + Sync API v1 + Official MCP server.
Master copy: `/root/.openclaw/workspace/skills/todoist-api-rest/`
All agent workspaces symlink to master — updates propagate instantly.

### 1password (Legacy ONLY)
Credential access to all 5 vaults (SO/HP/LE/GE/Legacy AI).
**No other agent has this skill. No exceptions.**
Location: `/root/.openclaw/workspace/skills/1password/`

### postiz + social-media-scheduler + social-monitor (Legacy + Empire)
Social media scheduling, content calendar, engagement monitoring.
Location: `/root/.openclaw/workspace/skills/` and `/root/.openclaw/workspace-empire/skills/`

### hyperframes (Legacy + Empire)
HTML→MP4 video rendering on VPS.

---

## Agent Tools Matrix

| Agent | Todoist | 1Password | GitHub | Coding Agent | tmux | Obsidian | openshell |
|-------|---------|-----------|--------|--------------|------|----------|-----------|
| ⚡ Legacy | ✅ | ✅ ONLY | ✅ Full | ✅ | ✅ | ✅ | ✅ Full |
| 💪 Spartan | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ Light |
| 🔮 Oracle | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ Light |
| 🌱 Eden | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ Light |
| 🏡 Haven | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ Light |
| 💰 Empire | ✅ | ❌ | ✅ Full | ✅ | ❌ | ✅ | ✅ |
| 🛡️ Guardian | ✅ | ❌ | ✅ Read | ✅ Audit | ✅ | ❌ | ✅ Full |
| ⚙️ Grid | ✅ | ❌ | ✅ Full | ✅ | ✅ | ❌ | ✅ Full |
| 🔺 Trinity | ✅ | ❌ | ✅ Read | ✅ | ✅ | ✅ Read | ✅ Full |

---

## Todoist Domain Routing

Each agent owns specific Todoist projects aligned to their domain.

| Agent | Domain | Todoist Projects |
|-------|--------|-----------------|
| ⚡ Legacy | SO (00-24) + All | 00-06, cross-domain |
| 💪 Spartan | HP Body | 27 Body, 28 Mastery, 33 Joy |
| 🔮 Oracle | HP Mind/Money | 25 Identity, 26 Mind, 29 Systems, 30 Finance, 31 Sovereignty, 32 Relationships, 34 Spirituality |
| 🌱 Eden | LE Environment | 50 Home, 58 Community |
| 🏡 Haven | LE Family | 50 Home, 51 Marriage, 52 Offspring, 53 Parents, 54 Heritage, 55 Health, 56 Estate, 58 Community |
| 💰 Empire | GE Business | 75 Brand, 76 Academy, 77 Products, 78 Marketing, 79 Biz Ops, 80 Treasury, 81 Strategy, 82 Partner, 83 Philanthropy |
| 🛡️ Guardian | SO Security | 04 Security |
| ⚙️ Grid | SO Infrastructure | 02 Infrastructure, 03 Command |
| 🔺 Trinity | SO Coordination | 03 Command, 00 Framework |

---

## Credential Access Policy

**Rule: Legacy is the sole gateway to 1Password. No exceptions.**

```
1Password (all 5 vaults)
        ↓
  ⚡ Legacy ONLY
        ↓
  Provisions tokens/keys into workspace skills
        ↓
  All agents use skills directly (no vault access)
```

If an agent needs a credential:
1. Agent requests from Legacy
2. Legacy fetches from 1Password
3. Legacy provisions it into a skill file or environment variable
4. Agent uses the skill — never touches 1Password directly

---

## CLI Access Policy

| Access Level | Agents | Tools |
|-------------|--------|-------|
| **Full CLI** | Legacy, Grid, Guardian, Trinity | openshell + tmux + full shell |
| **Moderate CLI** | Empire | openshell (deploys, builds, repos) |
| **Light CLI** | Spartan, Oracle, Haven, Eden | openshell (scripts, API calls only) |

---

## Adding New Skills — Process

1. **New system skill needed for all agents** → enable in `openclaw.json` `skills.entries`
2. **New skill for one agent only** → add to that agent's `workspace-[agent]/skills/`
3. **New skill for multiple agents** → create in Legacy workspace, symlink to others
4. **New credential-dependent skill** → Legacy provisions token, creates skill, deploys symlinks
5. **Document here** → update this file + agent's `TOOLS.md`

---

## File Locations

| File | Purpose |
|------|---------|
| `~/.openclaw/openclaw.json` | Global skill enable/disable + MCP servers |
| `~/.openclaw/workspace/skills/` | Legacy workspace skills (master copies) |
| `~/.openclaw/workspace-[agent]/skills/` | Agent-specific skill overrides |
| `~/.openclaw/workspace[-agent]/TOOLS.md` | Agent tool reference sheet |
| `/usr/lib/node_modules/openclaw/skills/` | System skill library (read-only) |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
