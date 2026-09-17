---
title: "Agent — Sentinel"
type: agent-profile
agent: sentinel
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 👁️ Sentinel — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` on omos-sentinel
**Workspace on VM:** `~/.openclaw/workspace/` (omos-sentinel)

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `sentinel` |
| **Name** | Sentinel |
| **Model** | `github-copilot/claude-haiku-4.5` |
| **Fallback** | `amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0` |
| **Pillar** | External watchdog — independent of all pillars |
| **Type** | Failover agent — monitors omos-core from outside |
| **Node** | omos-sentinel (`100.66.235.38` via Tailscale) |
| **Workspace** | `/home/ubuntu/.openclaw/workspace/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/sentinel (Sentinel Workspace)/` |

## Role

External watchdog running on omos-sentinel — a physically separate Lightsail instance from omos-core. Monitors Core from outside. If Core is unreachable, Sentinel alerts Zeus via Telegram DM. Haiku 4.5 by design — lightweight, fast, cheap, always-on.

## Skills

| Skill | Purpose |
|---|---|
| `agent-sentinel` | Identity + operating rules |
| `onemind-maintainer` | Vault maintenance |
| `taskflow` | Task management |
| `tavily` | Web search |
| `summarize` | Content summarization |
| `telegram` | Telegram messaging |
| `onemind-inbox` | Inbox processing |
| `onemind-sync` | Vault sync operations |

## Channels

| Channel | ID | Notes |
|---|---|---|
| Discord | `#sentinel-chat` (`1482355265092190319`) | Status reports |
| Telegram | DM via `@omos_sentinel_bot` | Emergency alerts when Core is down |

## Cron Jobs

| Job | Schedule | Time (ET) | Channel |
|---|---|---|---|
| `sentinel-heartbeat` | `0 */6 * * *` | Every 6 hours | `#agent-heartbeat` |
| `sentinel-daily-report` | `55 22 * * *` | 10:55 PM daily | `#sentinel-chat` |

## Node Details

| Property | Value |
|---|---|
| Instance | AWS Lightsail |
| Hostname | `omos-sentinel` |
| Tailscale IP | `100.66.235.38` |
| SSH | `ssh omos-sentinel` |
| OS | Ubuntu |
| Gateway service | `openclaw-gateway.service` (systemd --user) |
| Config | `/home/ubuntu/.openclaw/openclaw.json` |
| Vault path | `/home/ubuntu/onemind/` (git-synced) |

## Why Sentinel Is Separate

Sentinel must survive when everything else fails. If omos-core crashes, loses network, or has a corrupted config — Sentinel is on a different machine, different process, different config. It can still reach Zeus via Telegram DM to say "Core is down."

## Related

- [Workspace files](../workspaces%20(Agent%20Workspaces)/sentinel%20(Sentinel%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [Config backup](../install%20(Install%20&%20Deploy)/backups/omos-sentinel-openclaw-2026-05-11.json) — Current live config
- [Setup script](../install%20(Install%20&%20Deploy)/scripts%20(Deploy%20Scripts)/sentinel-setup.sh) — Full node deployment
