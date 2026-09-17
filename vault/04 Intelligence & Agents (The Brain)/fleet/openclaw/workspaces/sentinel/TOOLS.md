---
type: agent-tools
agent: sentinel
status: active
updated: 2026-05-11
---
# TOOLS.md — Sentinel 👁️

## Skills Loaded

| Skill | Purpose |
|---|---|
| `agent-sentinel` | Identity, operating rules, escalation protocol |
| `onemind-maintainer` | Vault maintenance operations |
| `taskflow` | Task management |
| `tavily` | Web search |
| `summarize` | Content summarization |
| `telegram` | Telegram messaging (DM alerts) |
| `onemind-inbox` | Inbox processing |
| `onemind-sync` | Vault sync operations |

## Infrastructure Access

| Resource | Access |
|---|---|
| omos-sentinel shell | Full (local) |
| omos-core | HTTP health check only — no SSH |
| Vault (`/home/ubuntu/onemind/`) | Read/write via git sync |
| Discord | Post to `#sentinel-chat` and `#agent-heartbeat` |
| Telegram | DM Zeus via `@omos_sentinel_bot` |

## Constraints

- No Docker, no databases — bare process only
- No heavy compute — preserve resources for gateway uptime
- No direct access to omos-core internals — observe from outside only
