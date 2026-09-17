---
id: so-03-dashboards
title: "Gateway Config Dashboards"
type: reference
domain: SO
status: active
created: 2026-05-10
modified: 2026-05-10
owner: zeus
tags: [domain/so, folder/03, infra/openclaw, infra/discord, infra/telegram, ops/monitoring]
---

# Gateway Config Dashboards

HTML snapshots of the live omos-core gateway config. Generated on demand or by cron. Open any `.html` file directly in a browser.

---

## Files

| File | Date | What It Shows |
|------|------|---------------|
| `gateway-config-YYYY-MM-DD.html` | Snapshot date | Full agent roster, channels, skills, cron jobs |

To regenerate: ask any agent `run /om-dashboard gateway` or trigger the `om-dashboard` cron job manually.

---

## Dashboard Sections

### Stats Bar (top)
Quick counts: agents, Discord bindings, Telegram topics, cron jobs. Red numbers = problems requiring action.

### Channel Cards
| Card | What to Check |
|------|--------------|
| Discord | `threadBindings: enabled` — if missing, threads won't spawn isolated sessions |
| Telegram | `streaming: block` — if `partial`, event loop will saturate under load |
| Runtime | Model and concurrency limits |

### Agent Roster Table

| Column | What It Means |
|--------|--------------|
| Discord | Green badge = channel bound. `—` = agent has no Discord channel routing |
| Telegram | Green badge = topic ID assigned. `—` = messages won't route to this agent |
| Skills | All skills the agent can invoke. Base set = 7. Legacy = 25 |
| Cron | Schedule badge. **Orange = broken delivery** (will fail-closed) |

### Cron Jobs Table

| Status | Meaning | Fix |
|--------|---------|-----|
| `✓ ok` | Delivery channel explicitly set | Nothing needed |
| `⚠ will fail` | `channel: "last"` with multiple channels active | `openclaw cron edit <id> --channel discord --to <channel-id>` |
| `no cron` (in roster) | Agent has no scheduled report | Add via `openclaw cron add` or ask Grid |

---

## Common Actions

```bash
# SSH to omos-core
ssh -i ~/.ssh/id_ed25519 ubuntu@100.123.126.34

# Check live channel status
sudo openclaw channels status --probe

# List all cron jobs
sudo -u ubuntu openclaw cron list

# Run a cron job immediately (test)
sudo -u ubuntu openclaw cron run <job-id>

# Fix a broken cron delivery channel
sudo -u ubuntu openclaw cron edit <job-id> --channel discord --to <discord-channel-id>

# Check agent identities loaded
sudo openclaw agents list

# Reload gateway config (hot reload)
pid=$(ps aux | grep 'openclaw.*index.js' | grep -v grep | awk '{print $2}' | head -1)
kill -HUP $pid
```

---

## Related
- [[MESSAGING]] — Full Discord channel ID and Telegram topic ID map
- [[omos-core]] — Primary node docs
- [[TAILSCALE-NETWORK]] — Network topology
