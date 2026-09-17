# Lobster Workflows — OneMind OS

> Automated workflows that run on schedule via OpenClaw's Lobster runtime.
> These are the "backbone" — things that happen predictably without being asked.

## What is Lobster?

[Lobster](https://github.com/openclaw/lobster) is OpenClaw's typed workflow runtime.
It turns skills and tools into composable pipelines with:
- Cron scheduling
- Resumable approvals (human-in-the-loop)
- Error handling and retries
- Step-by-step execution with typed inputs/outputs

## Workflows

| File | Schedule | Purpose |
| --- | --- | --- |
| `morning-brief.lob` | 7:00 AM ET daily | Daily briefing — news, tasks, system health, AI action recs → Discord |
| `health-check.lob` | Every 5 min | System health monitoring — gateway, NATS, OmniDB |
| `git-sync.lob` | Every 15 min | Workspace backup to git + change summary |
| `habit-checker.lob` | Morning + evening | HP daily habit check-ins via Discord |
| `second-brain.lob` | 11:00 PM nightly | Memory capture — decisions, ideas, follow-ups |
| `earnings-tracker.lob` | Monday 8 AM | Track AI/tech earnings, surface GE opportunities |
| `tech-digest.lob` | Thursday 9 AM | Tech + AI digest from HN, Reddit, RSS |
| `project-state.lob` | Every heartbeat | Keep workspace/state/projects.yaml current |
| `update-openclaw.lob` | On demand | Update gateway to latest npm release on VPS, remind Zeus to rebuild apps |

## How It Works

Lobster is enabled in `openclaw.json`:
```json
"lobster": { "enabled": true }
```

The `lob-1.0.0` skill is installed in `skills/`. Workflows in this folder
are auto-discovered by the Lobster runtime when the gateway starts.

## Adding New Workflows

Create a `.lob` file in this directory. It will be picked up on next gateway restart
or hot-reload. Follow the pattern in the existing files.

## Available Skills (30 total)

All skills in `../skills/` are available to lobster workflows:

| Skill | Purpose |
| --- | --- |
| `nexus-bridge` | NATS dispatch (cross-runtime) |
| `weather-1.0.0` | Weather lookups |
| `github-api-1.0.3` | GitHub operations |
| `obsidian-1.0.0` | Vault file operations |
| `browser-1.0.0` | Web browsing |
| `api-gateway-1.0.48` | HTTP API calls |
| `proactive-agent-3.1.0` | Proactive checks |
| `home-assistant-1.0.0` | Smart home control |
| `todoist` / `todoist-api` | Task management |
| `zoho-cli` / `zoho-mail-1.0.0` | CRM + calendar + mail |
| `resend-1.0.0` | Email sending |
| `twilio-1.0.0` | SMS |
| `spotify-1.0.0` | Music control |
| `rss-digest-1.0.0` | RSS feed parsing |
| `vision-1.0.0` | Computer vision |
| `desktop-control-1.0.0` | Desktop automation |
| `transit-1.0.0` | Transit lookups |
| `uber-1.0.0` | Uber rides |
| `onepassword-1.0.0` | 1Password credential queries |
| `clickup-1.0.0` | ClickUp tasks |
| `hitl-escalate-1.0.0` | Human-in-the-loop escalation |
| `lob-1.0.0` | Lobster runtime itself |
| `recipe-scraper-1.0.0` | Recipe extraction |
| `homestead-control` | Homestead devices |
| `skill-creator-0.1.0` | Generate new skills |
| `find-skills-0.1.0` | Skill discovery |
| `agent-browser-0.2.0` | Agent-driven browsing |
| `credit-repair-1.0.0` | Credit tools |

Plus built-in tools: `web_search`, `web_fetch`, `shell`, `file_read`, `file_write`.
