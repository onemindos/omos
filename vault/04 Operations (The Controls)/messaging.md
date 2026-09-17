---
id: so-02-messaging
title: "Messaging — Telegram"
type: reference
domain: SO
status: active
created: 2026-05-10
modified: 2026-05-12
owner: zeus
tags: [domain/so, folder/02, infra/messaging, infra/telegram]
---

# Messaging — Telegram

> Primary communication layer. Discord deactivated 2026-05-12 — system transitioned to autonomous mode.

---

## Bots

| Bot | Username | Node | Token Location |
|-----|----------|------|----------------|
| OneMindOS | @onemindos_bot | omos-core | `op://aaeyris5z44pdqspk277u2oppi/c2hsruwzaum6z4d7qjjyduadpq/credential` |
| Sentinel | @omos_sentinel_bot | omos-sentinel | `~/.openclaw/openclaw.json` |

---

## Telegram

**Group:** OneMindOS (`-1003905828404`) — Forum mode

### Agent Topics (omos-core / @onemindos_bot)

| Topic | ID | Agent |
|-------|----|-------|
| General | 1 | Legacy (main) |
| ⚡ Legacy | 292 | Legacy (main) |
| Zenith | 318 | Zenith |
| Unity | 319 | Unity |
| Eternal | 320 | Eternal |
| Sovereignty | 321 | Sovereignty |

### Operational Topics (cron delivery)

| Topic | ID | Receives |
|-------|----|----------|
| Approvals | 322 | HITL approval requests |
| 📥 Inbox | 323 | Captured inputs |
| 📋 Directives | 324 | Standing orders |
| 📊 Daily Reports | 325 | Agent daily reports (all 5 agents) |
| 📆 Weekly Reports | 326 | Weekly scorecards + reviews |
| 📣 Briefings | 327 | Legacy status briefing |
| 🖥️ System Status | 328 | Heartbeats + infra checks |
| 🚨 Errors | 329 | Agent errors |
| ⏱️ Cron Runs | 330 | Cron execution logs |
| 🔄 Active Work | 331 | Tasks in progress |
| Completed | 332 | Finished tasks |

### Sentinel (omos-sentinel / @omos_sentinel_bot)

| Topic | ID | Agent |
|-------|----|-------|
| Sentinel | 256 | Sentinel |

DM only for alerts when omos-core is unreachable.

---

## Exec Approvals

Approval requests DM Zeus directly (`target: dm`, approver: `577573558`).

---

## Discord

Deactivated 2026-05-12. Plugin disabled, adapter `enabled: false`.

Six agent chat channel configs remain in openclaw.json (disabled) as a rewire reference:

| Channel | ID | Agent |
|---------|----|-------|
| #legacy-chat | 1478604243509448705 | Legacy |
| #sentinel-chat | 1482355265092190319 | Sentinel |
| #zenith-chat | 1503272334990118952 | Zenith |
| #unity-chat | 1503272335921123369 | Unity |
| #eternal-chat | 1503272336978087986 | Eternal |
| #sovereignty-chat | 1503272338274390076 | Sovereignty |

To re-enable: set `channels.discord.enabled: true` and `plugins.entries.discord.enabled: true`.

---

## Config Locations

| Node | Config File |
|------|------------|
| omos-core | `~/.openclaw/openclaw.json` |
| omos-sentinel | `~/.openclaw/openclaw.json` |

```bash
# Reload after config change (no restart needed for most settings)
kill -HUP $(pgrep -f 'openclaw.*gateway')

# List Telegram topics
openclaw config get channels.telegram.groups

# List all bindings
openclaw config get bindings
```

---

## Related

- [[omos-core]] — Primary node
- [[omos-sentinel]] — Rescue/failover node
- [[TAILSCALE-NETWORK]] — Network topology
