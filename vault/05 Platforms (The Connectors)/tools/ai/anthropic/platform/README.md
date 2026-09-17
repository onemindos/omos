---
type: platform
domain: SO-05
status: active
updated: 2026-05-14
owner: zeus
tags: [platform/anthropic, product/platform, domain/so]
ai_summary: "Anthropic platform admin — billing, org settings, team seats, usage monitoring, and API key management."
---

# Anthropic Platform — Admin & Billing

Org-level settings, billing, seats, and usage for all Anthropic products.

## Access

- Console: `console.anthropic.com`
- Claude.ai admin: `claude.ai` → Settings
- Login: 1Password → SO vault → "claude.ai" / "Anthropic Console"

## Current Plan

| Product | Plan | Renewal | Cost |
| ------- | ---- | ------- | ---- |
| Claude.ai | *(Pro / Team — confirm)* | *(monthly/annual)* | *(confirm)* |
| Anthropic API | Pay-as-you-go | Monthly | Usage-based |
| Claude Code | Included with plan | — | — |

## Team Seats

| Member | Role | Seat Type |
| ------ | ---- | --------- |
| Zeus | Admin / Owner | Full |
| *(add as invited)* | — | — |

## Usage Monitoring

| Where | What |
| ----- | ---- |
| `console.anthropic.com` → Usage | API token consumption by model |
| `console.anthropic.com` → Billing | Monthly invoices, cost breakdown |
| `claude.ai` → Settings → Usage | Claude.ai conversation usage |

## API Key Management

- Keys stored in 1Password → SO vault → "Anthropic API Key"
- Rotate every 90 days or on team member departure
- Use separate keys for: production agents, dev/test, Claude Code CLI
- Never commit keys to git — always reference via 1Password

## Alerts

Set spending alerts at `console.anthropic.com` → Billing → Alerts.
Recommended: alert at 50% and 80% of monthly budget.
