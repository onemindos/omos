---
id: openclaw-channels
title: "OpenClaw Channels — Discord & Telegram"
type: reference
status: active
created: 2026-05-04
modified: 2026-05-11
owner: zeus
tags: [type/reference, status/active, domain/SO-02, agent/all]
---

# OpenClaw Channels — Discord & Telegram

Last updated: 2026-05-11. Full audit complete. Both channels live. Telegram security fixed. Talk mode confirmed on Azure Realtime. 23 agents active (19 old skill-agents + 4 new super-agents).

---

## Channel Architecture

```
Discord Server: OneMind OS (guild 1453411552181883027)
│
├── [onemindos bot — omos-core] — 5 active agent channels
│   ├── legacy-chat       → Legacy (main) — orchestrator, Zeus's primary interface
│   ├── zenith-chat       → Zenith — HP super-agent (body, mind, money, spirit)
│   ├── unity-chat        → Unity — LE super-agent (family, home, land, estate)
│   ├── eternal-chat      → Eternal — GE super-agent (empire, brand, capital)
│   └── sovereignty-chat  → Sovereignty — SO super-agent (infra, security, automation)
│
├── [onemindos bot — omos-core] — 3 utility channels (outbound only, inbound disabled)
│   ├── notifications      → System alerts / webhook targets
│   ├── agent-heartbeat    → Agent health reports
│   └── morning-brief      → Legacy daily briefing
│
├── COMMAND & CONTROL (category) — operational control
│   ├── approvals          → Exec approval requests (execApprovals.target)
│   ├── inbox              → Incoming captures queued for triage
│   └── directives         → Standing orders and mission-level directives
│
├── REPORTING (category) — agent output
│   ├── daily-reports      → Agent daily summaries
│   ├── weekly-reports     → Agent weekly summaries
│   └── briefings          → Morning brief and strategic briefings
│
├── SYSTEM HEALTH (category) — monitoring
│   ├── system-status      → Node health, uptime, connectivity
│   ├── errors             → Agent errors, exceptions, failure alerts
│   └── cron-runs          → Cron job execution logs
│
├── WORK PIPELINE (category) — execution tracking
│   ├── active-work        → Tasks currently in progress
│   └── completed          → Finished tasks and deliverables log
│
└── [sentinel bot — omos-sentinel] — 1 channel
    └── sentinel-chat     → Sentinel (failover/watchdog)

Telegram Group: @OneMindOS (group -1003905828404) — Forum mode
│
├── [onemindos_bot — omos-core] — 5 active topics
│   ├── topic 292  → Legacy (main)
│   ├── topic 318  → Zenith
│   ├── topic 319  → Unity
│   ├── topic 320  → Eternal
│   └── topic 321  → Sovereignty
│
└── [omos_sentinel_bot — omos-sentinel] — DM only
    └── (no group topics — Sentinel responds to DMs only when core unreachable)

```

---

## omos-core — Discord Channel Bindings

**Node:** `omos-core` (Tailscale `100.123.126.34`)
**Bot token:** 1Password → `op://aaeyris5z44pdqspk277u2oppi/qwomupt2iajig4twcv6ijdkyj4/token`
**Guild ID:** `1453411552181883027`
**Authorized user:** Zeus (`1024395579742228570`)

### Active Agent Channels

| Agent | Channel Name | Channel ID |
|-------|-------------|------------|
| `main` (Legacy) | legacy-chat | `1478604243509448705` |
| `zenith` | zenith-chat | `1503272334990118952` |
| `unity` | unity-chat | `1503272335921123369` |
| `eternal` | eternal-chat | `1503272336978087986` |
| `sovereignty` | sovereignty-chat | `1503272338274390076` |


**Utility channels (outbound only — `enabled: false` for inbound):**

| Channel Name | Channel ID | Purpose |
|-------------|------------|---------|
| notifications | `1478604241013833859` | System alerts / webhook targets |
| agent-heartbeat | `1478604235494133801` | Agent health reports |
| morning-brief | `1478604238287536148` | Legacy daily 7AM briefing output |

**Per-channel settings:** `autoThread: true`, `autoArchiveDuration: 10080` (7 days), `includeThreadStarter: true`, `requireMention: false`

---

## omos-core — Operational Channel Groups

### COMMAND & CONTROL (`1503277359586082847`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| approvals | `1503277446475157625` | Exec approval requests — `execApprovals.target` |
| inbox | `1503277448924500061` | Incoming captures queued for triage |
| directives | `1503277451826954300` | Standing orders and mission-level directives |

### REPORTING (`1503277390057439292`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| daily-reports | `1503277454318637157` | Agent daily summaries |
| weekly-reports | `1503277456939814913` | Agent weekly summaries |
| briefings | `1503277459565445201` | Morning brief and strategic briefings |

### SYSTEM HEALTH (`1503277392603644167`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| system-status | `1503277462686011403` | Node health, uptime, connectivity |
| errors | `1503277465529745438` | Agent errors, exceptions, failure alerts |
| cron-runs | `1503277468340060230` | Cron job execution logs |

### WORK PIPELINE (`1503277395220758559`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| active-work | `1503277470982475876` | Tasks currently in progress |
| completed | `1503277473440333978` | Finished tasks and deliverables log |

> **Note:** `execApprovals.target` in openclaw.json is wired to `channel:1503277446475157625` (`#approvals`). Approval buttons appear in this channel rather than DMs.

---

## omos-core — Telegram Topic Bindings

**Bot:** `@onemindos_bot`
**Bot token:** 1Password → `op://aaeyris5z44pdqspk277u2oppi/c2hsruwzaum6z4d7qjjyduadpq/credential`
**Group ID:** `-1003905828404`
**Authorized user:** Zeus (`577573558`)

### Active Topics

| Agent | Topic Name | Topic ID |
|-------|-----------|----------|
| `main` (Legacy) | ⚡ Legacy | `292` |
| `zenith` | 🏔 Zenith | `318` |
| `unity` | 🌿 Unity | `319` |
| `eternal` | 👑 Eternal | `320` |
| `sovereignty` | 🌐 Sovereignty | `321` |


**Group-level settings:** `requireMention: false`, `groupPolicy: allowlist`

---

## omos-sentinel — Channel Config

**Node:** `omos-sentinel` (Tailscale `100.66.235.38`)
**Discord bot token:** 1Password → `op://aaeyris5z44pdqspk277u2oppi/czrullgbzaz3dty2pki6s4b4oq/credential`
**Telegram bot:** `@omos_sentinel_bot` — token in 1Password → `op://aaeyris5z44pdqspk277u2oppi/lmpukgoo7w4bb5katnt24vqobe/credential`

**Discord:** Sentinel responds in `#sentinel-chat` only (`1482355265092190319`)
**Telegram:** DM only — no group topics. Activates when omos-core is unreachable.

---

## System Status (2026-05-11)

| Channel | Connected | Auth | Issues |
|---------|-----------|------|--------|
| Discord `@OneMind OS` | ✅ | Token configured | 17 stale channel IDs in config; limited intents |
| Telegram `@onemindos_bot` | ✅ | Token configured | execApprovals sub-client pairing issue |

**Talk mode:** Azure OpenAI Realtime `gpt-realtime-1.5` via `onemind-openai.openai.azure.com` (eastus2) — active on omos-core browser Talk.
**Discord voice:** `voice.model=gpt-realtime-1.5` configured — `/vc join` pending test.
**Telegram voice:** Inbound transcription active (Deepgram). Outbound TTS pending activation.

---

## Active Config Settings (2026-05-11)

### Discord
| Key | Value | Notes |
|-----|-------|-------|
| `streaming.mode` | `progress` | Live tool progress shown while agent works |
| `streaming.progress.toolProgress` | `true` | Tool calls visible in progress display |
| `commands.native` | `auto` | Slash commands auto-registered |
| `replyToMode` | `first` | Native Discord quote-reply on first message |
| `ackReaction` | `⚡` | Reaction while processing |
| `autoPresence.enabled` | `true` | Bot status reflects node health |
| `autoPresence.healthyText` | `One Mind OS — Online` | |
| `execApprovals.enabled` | `true` | Button-based approvals via DM |
| `execApprovals.cleanupAfterResolve` | `true` | Prompt removed after approve/deny |
| `threadBindings.enabled` | `true` | Threads bind to spawned sessions |
| `threadBindings.idleHours` | `24` | Sessions expire after 24h idle |
| `historyLimit` | `20` | Guild channel message history |
| `groupPolicy` | `allowlist` | |
| `allowFrom` | `["1024395579742228570"]` | Zeus Discord user ID only |

### Telegram
| Key | Value | Notes |
|-----|-------|-------|
| `streaming.mode` | `partial` | Live edit preview as agent types |
| `streaming.preview.toolProgress` | `true` | Tool calls shown in preview |
| `allowFrom` | `["577573558"]` | Zeus Telegram user ID only |
| `groupAllowFrom` | `["577573558"]` | Group sender lock — was missing, now fixed |
| `capabilities.inlineButtons` | `all` | Approve/reject buttons active everywhere |
| `ackReaction` | `⚡` | Reaction while processing |
| `replyToMode` | `first` | Quote-reply threading |
| `execApprovals.enabled` | `true` | DM-based approval flow |
| `execApprovals.approvers` | `["577573558"]` | Zeus only |
| `dmPolicy` | `allowlist` | |
| `groupPolicy` | `allowlist` | (no separate key — enforced via `groupAllowFrom`) |
| `historyLimit` | `50` | Group topic context window |
| `commands.native` | `auto` | Bot command menu auto-registered (66 commands, budget-capped) |

### System
| Key | Value |
|-----|-------|
| `commands.ownerAllowFrom` | `["discord:1024395579742228570", "telegram:577573558"]` |
| `session.dmScope` | `per-channel-peer` |

---

## How Threads Work (Discord)

```
User sends message → e.g. legacy-chat
    → OpenClaw creates a thread from that message
    → Legacy responds inside the thread
    → All follow-up messages stay in thread
    → Channel root stays clean (thread summaries only)
    → Thread auto-archives after 7 days idle
    → Session expires after 24h idle
```

Each thread = one isolated session. Subagents spawned via ACP also bind to threads.

---

## Testing

**Telegram — send test to topic:**
```bash
curl -s "https://api.telegram.org/bot<TOKEN>/sendMessage" \
  -d "chat_id=-1003905828404&message_thread_id=292&text=Test"
```

**Discord — send test to channel:**
```bash
curl -X POST "https://discord.com/api/v10/channels/1478604243509448705/messages" \
  -H "Authorization: Bot <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"content":"Test"}'
```

**Check channel status:**
```bash
ssh ubuntu@100.123.126.34 "openclaw channels status"
ssh ubuntu@100.123.126.34 "openclaw channels capabilities --channel discord --target channel:1478604243509448705"
```

---

## Related

- [[discord (Discord Channel)/discord]] — Discord channel config detail
- [[telegram (Telegram Channel)/telegram]] — Telegram channel config detail
- [[../../../integrations (App Connections)/messaging (Messaging)/MESSAGING|MESSAGING]] — Bot registry and token locations
- [[../../../omos-core|omos-core]] — Primary node
- [[../../../omos-sentinel|omos-sentinel]] — Rescue/failover node
