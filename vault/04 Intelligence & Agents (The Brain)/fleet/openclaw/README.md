---
domain: 02
name: OpenClaw Platform
agent: Grid
status: active
updated: 2026-05-11
---

# OpenClaw Platform

OpenClaw is the AI agent gateway running 6 One Mind agents across omos-core and omos-sentinel. It handles messaging channels, agent identity, memory, skill execution, cron jobs, and hooks.

This folder mirrors the structure of [docs.openclaw.ai](https://docs.openclaw.ai) so every topic has a predictable home.

## Our Setup

| Item | Value |
|---|---|
| Platform version | OpenClaw 2026.5.x |
| Primary node | omos-core (`100.123.126.34` via Tailscale) |
| Failover node | omos-sentinel |
| AI provider | GitHub Copilot → `claude-sonnet-4.6` |
| Active agents | 6 (5 super-agents on omos-core + Sentinel on omos-sentinel) |
| Channels | Telegram (5 topics) + Discord (6 agent channels + operational) |
| Fallback model | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` |

## Agent Architecture

### Super-Agents (Sonnet 4.6 — full domain orchestrators)

| Agent | ID | Discord Channel | Telegram Topic | Owns |
|---|---|---|---|---|
| Legacy | `main` | `#legacy-chat` | `292` | Orchestrator — Zeus's primary interface |
| Zenith | `zenith` | `#zenith-chat` | `318` | HP 25-49 — body, mind, money, spirit, intelligence |
| Unity | `unity` | `#unity-chat` | `319` | LE 50-74 — family, home, land, estate, generations |
| Eternal | `eternal` | `#eternal-chat` | `320` | GE 75-99 — empire, brand, capital, ventures, wealth |
| Sovereignty | `sovereignty` | `#sovereignty-chat` | `321` | SO 00-24 — infra, security, automation, dev, research |
| Sentinel | `sentinel` | `#sentinel-chat` | DM only | External watchdog + failover |

### Skills (loaded on demand by super-agents)

Sub-agents were retired — super-agents now load domain skills directly as needed:

| Super-Agent | Skills Available |
|---|---|
| Zenith | `spartan`, `oracle`, `identity`, `librarian`, `pulse` |
| Unity | `haven`, `helix`, `eden`, `nomad` |
| Eternal | `empire`, `empire-ops`, `creed`, `trinity`, `apex`, `finance`, `open-prose` |
| Sovereignty | `grid`, `guardian`, `guardian-ops`, `overwatch`, `eagleeye`, `coding-agent`, `github`, `diffs` |

## Quick Links

| Need to... | Go to |
|---|---|
| Troubleshoot Mac ↔ VPS sync | [gateway/SYNC-ARCHITECTURE.md](gateway/SYNC-ARCHITECTURE.md) — full sync reference |
| Check workspace backup health | [workspaces/backup-repo.md](workspaces%20(Agent%20Workspaces)/backup-repo.md) |
| See all agents + skills | [agents/README.md](agents/README.md) |
| View cron schedule | [cron/CRON.md](cron%20(Scheduled%20Jobs)/CRON.md) |
| Add a new MCP server | [mcp/README.md](mcp%20(MCP%20Servers)/README.md) |

## Sections

| Folder | What It Covers | Official Docs |
|---|---|---|
| [start/](start/) | Setup, onboarding, guides, lore | [docs.openclaw.ai/start](https://docs.openclaw.ai/start/getting-started.md) |
| [gateway/](gateway/) | Config, health, logging, architecture, **sync troubleshooting** | [docs.openclaw.ai/gateway](https://docs.openclaw.ai/gateway/index.md) |
| [channels/](channels/) | Telegram, Discord, and all messaging channels | [docs.openclaw.ai/channels](https://docs.openclaw.ai/channels/index.md) |
| [nodes/](nodes/) | iOS/Android apps, voice, camera, media | [docs.openclaw.ai/nodes](https://docs.openclaw.ai/nodes/index.md) |
| [platforms/](platforms/) | macOS, iOS, onemind-app specifics | [docs.openclaw.ai/platforms](https://docs.openclaw.ai/platforms/index.md) |
| [concepts/](concepts/) | Memory, soul, sessions, architecture deep-dives | [docs.openclaw.ai/concepts](https://docs.openclaw.ai/concepts/architecture.md) |
| [automation/](automation/) | Cron jobs, hooks, taskflow, standing orders | [docs.openclaw.ai/automation](https://docs.openclaw.ai/automation/index.md) |
| [tools/](tools/) | Skills, use cases, browser, search, exec | [docs.openclaw.ai/tools](https://docs.openclaw.ai/tools/index.md) |
| [plugins/](plugins/) | Plugin SDK, manifests, plugin catalog | [docs.openclaw.ai/plugins](https://docs.openclaw.ai/plugins/architecture.md) |
| [clawhub/](clawhub/) | Publishing and installing skills via ClawHub | [docs.openclaw.ai/clawhub](https://docs.openclaw.ai/clawhub/index.md) |
| [providers/](providers/) | AI model providers, costs, API configs | [docs.openclaw.ai/providers](https://docs.openclaw.ai/providers/index.md) |
| [install/](install/) | VPS deployment, Docker, scripts, systemd | [docs.openclaw.ai/install](https://docs.openclaw.ai/install/index.md) |
| [cli/](cli/) | Full CLI reference | [docs.openclaw.ai/cli](https://docs.openclaw.ai/cli/index.md) |
| [reference/](reference/) | Master reference, ecosystem inventory, templates | [docs.openclaw.ai/reference](https://docs.openclaw.ai/reference/credits.md) |
| [security/](security/) | Threat model, CVEs, hardening | [docs.openclaw.ai/gateway/security](https://docs.openclaw.ai/gateway/security/index.md) |
| [help/](help/) | Troubleshooting, runbooks, live state, FAQ | [docs.openclaw.ai/help](https://docs.openclaw.ai/help/index.md) |
| [web/](web/) | Dashboard, WebChat, TUI, API, observability | [docs.openclaw.ai/web](https://docs.openclaw.ai/web/index.md) |
| [costs/](costs/) | Hosting provider cost comparison | — |
| [workspaces/](workspaces/) | Per-agent workspace files (identity, soul, memory) | — |

## Related

- [[AGENTS.md]] — Universal agent instructions
- [[CLAUDE.md]] — Claude Code instructions for this vault
- `SO/03 Infrastructure` — Parent domain
