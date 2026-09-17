---
title: "Paperclip"
type: reference
status: active
tags: [platform, ai, agents, orchestration, multi-agent, open-source]
---
# Paperclip

> Open-source multi-agent orchestration platform. If OpenClaw is an employee, Paperclip is the company.

- **GitHub:** https://github.com/paperclipai/paperclip
- **Stars:** 68k+ ⭐
- **Stack:** Node.js server + React UI

## What It Is
Paperclip is the app people use to manage AI agents for work. It orchestrates a team of AI agents to run a business — bring your own agents, assign goals, track work and costs from one dashboard. Looks like a task manager; under the hood: org charts, budgets, governance, goal alignment, agent coordination.

## How It Works
| Step | Example |
|------|---------|
| Define the goal | "Build the #1 AI note-taking app to $1M MRR." |
| Hire the team | CEO, CTO, engineers, designers, marketers — any bot, any provider |
| Approve and run | Review strategy, set budgets, hit go, monitor from dashboard |

## Supported Agents
Works with: OpenClaw, Claude Code, Codex, Cursor, Bash, HTTP
Rule: if it can receive a heartbeat, it's hired.

## Key Features
| Feature | Notes |
|---------|-------|
| Bring Your Own Agent | Any agent, any runtime, one org chart |
| Goal Alignment | Every task traces back to the company mission |
| Heartbeats | Agents wake on schedule, check work, act — flows up/down org chart |
| Cost Control | Monthly budgets per agent — hits limit, stops. No runaway costs |
| Multi-Company | One deployment, many companies, complete data isolation |
| Ticket System | Every conversation traced, every decision explained, full tool-call trace |
| Phone Control | Manage autonomous businesses from mobile |

## Self-Host
```bash
# Docker
docker compose up
```
Docs: https://paperclip.ai/docs

## OneMind OS Angle
Natural control plane for Zeus's multi-agent fleet — coordinates OpenClaw, Claude Code, Codex, Hermes toward a shared mission with budget guardrails and audit trails. Pairs with NATS as the messaging backbone and OpenClaw as the execution layer.

## Cost
Open-source (MIT). Self-hosted. Bring your own LLM API keys.
