---
type: audit
title: OpenClaw Optimization — May 2026
date: 2026-05-11
author: Legacy (Claude Code)
status: complete
---

# OpenClaw Optimization Audit — 2026-05-11

## What Was Found

### Critical Gaps

1. **Super-agents not registered** — Zenith, Unity, Eternal, Sovereignty had workspace files and plugin skills but were **completely absent from `openclaw.json`**. They could not receive messages, be routed to, or run cron jobs.

2. **All domain agents had zero domain skills** — Every non-main agent had the same 5 generic skills (`onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync`). None loaded their domain skill. Spartan didn't know it was Spartan. Grid didn't know it was Grid.

3. **Super-agents missing MEMORY.md** — Zenith, Unity, Eternal, Sovereignty had no MEMORY.md in their workspaces — no persistent context storage.

4. **Watch service not monitoring new workspaces** — `inotifywait` glob expanded at process start time. Zenith, Unity, Eternal, Sovereignty workspaces were invisible to the backup watcher.

5. **No Discord/Telegram bindings for super-agents** — The 4 new super-agent channels had no routing config. Messages to `#zenith-chat` or Telegram topic 318 went nowhere.

### Optimization Gaps

6. **`web-readability` plugin disabled** — Agents couldn't parse URLs into clean text.

7. **`tavily` plugin disabled** — Despite being in the skills allowlist and loaded as a skill, the plugin itself was off. Tavily search was broken for all agents.

8. **No model fallbacks** — `defaults.model.fallbacks` was `[]`. If GitHub Copilot went down, all agents failed hard.

---

## What Was Fixed

### `openclaw.json` changes

| Change | Details |
|---|---|
| Added 4 super-agents | `zenith`, `unity`, `eternal`, `sovereignty` with Sonnet 4.6 model + domain skills |
| Injected domain skills | 17 domain agents now load their correct skill on startup |
| Enabled `web-readability` | Agents can parse webpage content |
| Enabled `tavily` plugin | Tavily search now functional (was skill-only, not plugin) |
| Added model fallbacks | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` as fallback for all |
| Added 4 Discord bindings | zenith→`#zenith-chat`, unity→`#unity-chat`, eternal→`#eternal-chat`, sovereignty→`#sovereignty-chat` |
| Added 4 Discord channel configs | Channels 1503272334990118952–338274390076 now enabled with autoThread |
| Added 4 Telegram topic bindings | Topics 318–321 now route to respective super-agents |

### Super-agent skill lists

| Agent | Identity Skill | Domain Skills |
|---|---|---|
| `zenith` | `agent-zenith` | `spartan`, `oracle`, `identity`, `librarian`, `pulse` |
| `unity` | `agent-unity` | `haven`, `helix`, `eden`, `nomad` |
| `eternal` | `agent-eternal` | `empire`, `empire-ops`, `creed`, `trinity`, `apex`, `finance` |
| `sovereignty` | `agent-sovereignty` | `grid`, `guardian`, `guardian-ops`, `overwatch`, `eagleeye`, `coding-agent`, `github`, `diffs` |

All super-agents also get: `onemind-maintainer`, `taskflow`, `taskflow-inbox-triage`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync`

### Other actions

- Created `MEMORY.md` in all 4 super-agent workspaces
- Reloaded gateway via `kill -HUP <pid>`
- Restarted `watch.sh` — new inotifywait process now monitors all 15 workspaces including the 4 new ones

---

## Current State (Post-Optimization)

```
Total agents registered: 23
Super-agents (Sonnet 4.6): main, zenith, unity, eternal, sovereignty, sentinel
Domain agents (Haiku 4.5): 17
Discord channels active: 26
Telegram topics active: 22
Cron jobs: 23 (4 new super-agent reports added)
Model fallback: aws-bedrock Sonnet 4.6
Plugins active: web-readability, tavily, exa, open-prose, discord, telegram, active-memory, memory-core, memory-wiki, document-extract, browser, deepgram, azure-speech, duckduckgo, llm-task, openshell, phone-control, skill-workshop, tokenjuice, webhooks, device-pair, file-transfer, amazon-bedrock, github-copilot, microsoft-foundry
Default search provider: tavily
Discord streaming: progressive (real-time output)
```

---

## Wave 2 — MCP & Skills Optimization (same day)

| Change | Details |
|---|---|
| Default search provider | Set to `tavily` (was unset — duckduckgo was fallback) |
| Voice-call skill enabled | `voice-call` skill entry now `enabled: true` |
| `om-dashboard` + `onemind-status` + `onemind-report` | Added to Legacy + all 4 super-agents |
| `obsidian-vault-maintainer` | Added to Legacy for direct Obsidian vault ops |
| `exa` plugin enabled | Added to sovereignty, zenith, librarian, trinity, pulse |
| `open-prose` plugin enabled | Added to creed, empire, eternal for writing tasks |
| Discord streaming → `progressive` | Real-time output instead of waiting for full response |
| `groupChat.visibleReplies` → `all` | Full tool results visible, not just summaries |
| Cron jobs for super-agents | 4 new jobs: sovereignty/zenith/eternal daily at 23:05-23:20 ET, unity weekly Sunday 10:30 ET |
| Super-agent cron cadence | Super-agents synthesize after sub-agents finish — build a pillar-level view from sub-reports |

### Super-Agent Final Skill Counts

| Agent | Skills |
|---|---|
| Legacy (`main`) | 29 |
| Zenith | 18 (identity + HP stack + exa + dashboard) |
| Unity | 16 (identity + LE stack + dashboard) |
| Eternal | 19 (identity + GE stack + open-prose + dashboard) |
| Sovereignty | 21 (identity + SO stack + exa + dev tools + dashboard) |

---

## Remaining Opportunities

| Item | Priority | Notes |
|---|---|---|
| Sentinel node config alignment | High | Sentinel on omos-sentinel needs same super-agent registration audit |
| `openrouter` as second fallback | Low | Multi-cloud fallback chain for maximum resilience |
| Exa API key verification | Low | Plugin enabled — needs key configured in credentials if not already set |
