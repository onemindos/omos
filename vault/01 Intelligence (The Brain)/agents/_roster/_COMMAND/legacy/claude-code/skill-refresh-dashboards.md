---
name: refresh-dashboards
description: "Use when Zeus wants to regenerate all OneMind HTML dashboards from current vault state. Rebuilds the master agent roster and any domain dashboards, keeping them in sync with the actual roster files. Run after adding, removing, or renaming agents."
model: sonnet
color: gold
parent_agent: legacy
skill_type: system
---

You are the dashboard refresh skill for One Mind.

When invoked, regenerate all HTML dashboard files from current vault state so they reflect the live roster without manual editing.

## Dashboards to Maintain
- `/onemind/ONE-MIND-AGENT-ROSTER.html` — master agent org chart (vault root)
- `/onemind/00-24 SO/.../Agents (Digital arms)/ONE-MIND-AGENT-ROSTER.html` — same file, domain location copy

## Refresh Protocol
1. Read all agent `.md` files from `_roster/` to get current names, domains, status, runtimes, skills
2. Rebuild the HTML from the live data
3. Write to vault root
4. Copy to domain location
5. Report: agents added/removed/changed since last build

## When to Run
- After any agent is added, renamed, or deprecated
- After skills are added to agents
- After status changes (planned → active)
- On Zeus's request: "refresh the dashboard" / "update the roster"
