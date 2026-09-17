---
name: om-dashboard
description: Generate live HTML status dashboards for any area of the One Mind system — gateway config, agent roster, cron jobs, channel health, vault structure, or pillar status. Use on demand or via cron. Output is a self-contained HTML file saved to the vault.
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# One Mind Dashboard Generator

Generate a live HTML status dashboard for any area of the One Mind system. Dashboards are dark-themed, self-contained HTML files saved to the vault at `00-24 SO/03 Command/command (Command Center)/Core/dashboards/`.

## When to Use This Skill

Invoke when:
- User asks for a status overview of any system area
- A cron job fires requesting a dashboard refresh
- After significant config changes (agents, channels, cron, vault structure)
- User types `/om-dashboard` followed by a dashboard type

## Dashboard Types

| Type | Command | What It Shows |
|------|---------|--------------|
| Gateway | `/om-dashboard gateway` | Agents, channels, skills, cron jobs, channel health |
| Agents | `/om-dashboard agents` | All 19 agents — identity, workspace, model, skills, routing rules |
| Cron | `/om-dashboard cron` | All scheduled jobs — schedule, delivery channel, next run, status |
| Vault | `/om-dashboard vault` | Pillar structure, domain count, recent changes, inbox count |
| Channels | `/om-dashboard channels` | Discord channels + Telegram topics, binding status, threading config |
| Pillar | `/om-dashboard pillar <SO\|HP\|LE\|GE>` | Domains, agents, active tasks, health for one pillar |
| Full | `/om-dashboard full` | All of the above combined |

If no type is specified, default to `gateway`.

## How to Generate

### Step 1 — Collect live data

Run these commands on omos-core (SSH: `ssh -i ~/.ssh/id_ed25519 ubuntu@100.123.126.34`):

```bash
# Full config
sudo cat /home/ubuntu/.openclaw/openclaw.json

# Agent list with routing
sudo -u ubuntu openclaw agents list --json

# Cron jobs
sudo -u ubuntu openclaw cron list --json

# Channel status
sudo -u ubuntu openclaw channels status --probe

# Agent bindings
sudo -u ubuntu openclaw agents bindings --json

# Vault inbox count
ls "/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/99 Inbox (Queue)/" | wc -l

# Recent vault commits (last 24h)
git -C /home/ubuntu/onemind log --since="24 hours ago" --oneline
```

### Step 2 — Build the HTML

Generate a self-contained HTML file using the design system below. No external dependencies — all CSS inline.

### Step 3 — Save to vault

Save as:
```
/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/command (Command Center)/Core/dashboards/<type>-<YYYY-MM-DD>.html
```

The `onemind-watch.service` will auto-commit and push to GitHub.

### Step 4 — Announce

Post a one-line summary to the agent's Discord channel:
```
Dashboard saved: <type>-<date>.html — <key stat 1>, <key stat 2>, <N warnings>
```

---

## Design System

All dashboards use this consistent dark theme. Copy exactly.

```css
:root {
  --bg: #0f1117;
  --surface: #1a1d27;
  --surface2: #252836;
  --border: #2e3347;
  --text: #e2e8f0;
  --text-muted: #64748b;
  --accent: #f59e0b;
  --green: #4ade80;
  --green-bg: #052e16;
  --green-border: #14532d;
  --orange: #fb923c;
  --orange-bg: #422006;
  --orange-border: #7c2d12;
  --red: #f87171;
  --muted-bg: #1c1917;
  --muted-border: #292524;
}
```

### Component Classes

```css
/* Layout */
.header        /* top bar with title and subtitle */
.stats         /* flex row of stat cards */
.stat          /* individual stat card: .stat-val + .stat-lbl */
.section       /* padded content section */
.sec-title     /* section heading — uppercase, accent color */
.info-grid     /* responsive card grid */
.card          /* info card with .irow rows */

/* Data display */
table          /* full-width bordered table */
th             /* muted uppercase header */
td             /* standard cell */
.badge         /* inline status pill */
.badge-ok      /* green */
.badge-warn    /* orange */
.badge-missing /* muted gray */
.skill-tag     /* small inline skill label */
.model-tag     /* blue inline model label */
.warn-box      /* orange warning banner */
```

### Pillar Colors
```
SO → #3b82f6 (blue)
HP → #ef4444 (red)
LE → #22c55e (green)
GE → #eab308 (yellow)
```

### Agent Colors
```
main/Legacy     → #f59e0b   spartan    → #ef4444   oracle     → #8b5cf6
genesis         → #f97316   pulse      → #06b6d4   librarian  → #6366f1
eden            → #22c55e   haven      → #10b981   helix      → #84cc16
nomad           → #a3e635   empire     → #eab308   apex       → #f59e0b
creed           → #ec4899   trinity    → #14b8a6   guardian   → #3b82f6
grid            → #64748b   overwatch  → #94a3b8   eagleeye   → #0ea5e9
sentinel        → #475569
```

---

## Status Logic

Apply these rules when setting badge states:

| Condition | Badge |
|-----------|-------|
| Cron `delivery.channel == "last"` | `badge-warn` — will fail-closed |
| Agent has no Discord binding | `badge-missing` |
| Agent has no Telegram topic | `badge-missing` |
| Agent has no cron job | `badge-missing` |
| Channel `enabled: false` | `badge-warn` |
| Gateway `eventLoopUtilization > 0.9` | `badge-warn` |
| Everything nominal | `badge-ok` |

Always show a warning count in the stats bar. Zero warnings = green. Any warnings = orange count.

---

## Cron Schedule (auto-refresh)

This skill is triggered daily by a cron job on omos-core. The job runs as `overwatch` at 09:00 ET and regenerates the `gateway` dashboard. Overwatch announces the result in `#overwatch-chat`.

To add this cron job:
```bash
sudo -u ubuntu openclaw cron add \
  --agent overwatch \
  --name "om-dashboard-daily" \
  --cron "0 9 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --announce \
  --channel discord \
  --to 1502954950764793870 \
  --message "Run /om-dashboard gateway and save the output to the vault. Post a summary here when done."
```

---

## Model Note

This skill is designed to run efficiently on `claude-haiku-4.5` for routine cron-triggered refreshes. Use `claude-sonnet-4.6` only when generating a full multi-section dashboard or when asked explicitly by Zeus.

Set the cron job model override:
```bash
sudo -u ubuntu openclaw cron edit <job-id> --model amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0
```

---

## Output Quality Rules

- **Never hallucinate data** — only report what live commands return
- **Always include generation timestamp** in the HTML footer
- **Always include a warnings count** — 0 warnings is a valid and good state
- **File name must include date** — `gateway-2026-05-10.html` not `gateway-latest.html`
- **Announce to Discord after save** — one line, not a full report
