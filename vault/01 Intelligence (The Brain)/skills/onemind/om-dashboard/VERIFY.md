# Verification: One Mind Dashboard Skill

Run these checks after install or after any SKILL.md update.

## Skill File Integrity
- [ ] `/home/ubuntu/.openclaw/plugin-skills/om-dashboard/SKILL.md` exists on omos-core
- [ ] Vault copy matches VPS copy (no drift)
- [ ] All 19 agents have `om-dashboard` in their skills array in `openclaw.json`

## On-Demand Invocation
- [ ] Send `run /om-dashboard gateway` to any agent channel in Discord
- [ ] Agent responds with a summary and saves HTML to vault dashboards folder
- [ ] HTML file appears at `00-24 SO/03 Command/.../Core/dashboards/gateway-YYYY-MM-DD.html`
- [ ] `onemind-watch.service` commits and pushes the file within 30s

## Dashboard Types
- [ ] `/om-dashboard gateway` — agent roster, channels, skills, cron
- [ ] `/om-dashboard cron` — all scheduled jobs, delivery status, next run
- [ ] `/om-dashboard channels` — Discord + Telegram binding health
- [ ] `/om-dashboard agents` — identity, workspace, routing rules per agent
- [ ] `/om-dashboard vault` — pillar structure, inbox count, recent changes
- [ ] `/om-dashboard full` — all sections combined

## Model Efficiency
- [ ] Cron-triggered runs use `claude-haiku-4.5` (fast, cheap)
- [ ] On-demand runs may use `claude-sonnet-4.6` (default)

## Output Quality
- [ ] HTML is self-contained (no external dependencies)
- [ ] Dark theme applied consistently
- [ ] Stats bar shows warning count
- [ ] Badge colors match status logic (green/orange/gray)
- [ ] Generation timestamp in HTML footer
- [ ] File name includes date (`gateway-2026-05-10.html`)

## Cron Auto-Refresh (Overwatch)
- [ ] Overwatch has `om-dashboard-daily` cron at 09:00 ET
- [ ] Job uses haiku model
- [ ] Delivery goes to `#overwatch-chat` (Discord `1502954950764793870`)
