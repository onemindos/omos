---
title: "Skill Pack: Init"
type: skill
status: active
created: 2026-03-27
tags: [type/skill, domain/codex]
---

# Init Pack

> Bootstrap a new agent session by loading identity, preferences, active goals, and available skills.

## CODEX Stage: Utility (System Initialization)

## What It Does
Init is the cold-start skill. When an agent begins a new session with no prior context, Init loads the essential system files — identity, preferences, and active goals — so the agent knows who it is, who it serves, and what matters right now. It is the first skill any agent should run.

## When to Use
- Agent starts a fresh session with no loaded context
- Agent needs to reset and reload its baseline understanding
- New agent is being onboarded to the OneMind system for the first time

## Files
| File | Purpose |
|------|---------|
| `SKILL.md` | Core skill instructions |
| `INSTALL.md` | Self-installation for agents |
| `VERIFY.md` | Post-install verification |
