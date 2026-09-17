---
title: "Flue — Agent Harness Framework"
type: research-candidate
status: evaluating
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: agent-framework
tags: [research, ai, agents, harness, typescript, local, candidate]
ai_summary: "TypeScript agent harness framework — gives any model context, tools, skills, sandbox, and durable execution. Candidate for local on-machine autonomous agents in OneMind OS."
---

# Flue — Agent Harness Framework

> **Status:** Candidate — evaluating for local on-machine agent deployment in OneMind OS.
> When adopted, move to `SO/01 The Fabric/nodes/digital/services-registry/ai/agents/`.

## Links

- **Website:** https://flueframework.com/
- **Docs:** https://flueframework.com/guide
- **Reference:** https://flueframework.com/reference
- **CLI:** https://flueframework.com/cli
- **SDK:** https://flueframework.com/sdk
- **Ecosystem:** https://flueframework.com/ecosystem
- **GitHub:** https://github.com/flueframework (via site)

---

## What It Is

Flue is a **TypeScript agent harness framework** for building autonomous AI agents that run locally or in the cloud. It gives any model (Claude, GPT, etc.) the full context and environment it needs for truly autonomous work.

---

## Key Features

| Feature | Description |
|---|---|
| **Agents** | Context-persistent autonomous agents across conversations and events |
| **Workflows** | Structured automations — code guides agent reasoning from input to result |
| **Sandboxes** | Secure environment for agents to use tools, modify files, run code |
| **Durable Execution** | Survives failures and restarts — agents resume where they left off |
| **Subagents** | Delegate work to specialized sub-agents |
| **Tools** | Typed actions for calling APIs, querying data, making controlled changes |
| **Skills** | Reusable expertise packaged as `SKILL.md` files — same pattern as Hermes |
| **MCP Servers** | Connect to the full MCP ecosystem |
| **Observability** | OpenTelemetry, Braintrust, Sentry |
| **Chat** | Slack, Teams, Discord, GitHub integration |

---

## Deployment Targets

- Node.js
- Cloudflare Workers
- GitHub Actions
- GitLab CI/CD
- Daytona
- Render

---

## Why It's Interesting for OneMind OS

- Uses `SKILL.md` files — same pattern as Hermes skills, familiar architecture
- Designed for **local on-machine autonomous agents** — aligns with OneMind's sovereignty values
- Secure sandboxed execution — agents can't escape their scope
- Durable execution means agents survive reboots and failures
- Could power always-on local agents (monitoring, automation, fleet ops) without relying on cloud

## Use Case Hypothesis

> Deploy Flue agents locally on `omos-mbp` and fleet nodes for persistent, always-on automation tasks — e.g. vault watcher, NATS event handler, fleet health monitor.

---

## Decision Criteria

- [ ] Compare against existing Hermes cron jobs — is there overlap?
- [ ] Evaluate sandbox security model
- [ ] Check TypeScript skill compatibility with Hermes SKILL.md format
- [ ] Test local deployment on `omos-mbp`
- [ ] Determine which OneMind tasks benefit from always-on vs on-demand agents
