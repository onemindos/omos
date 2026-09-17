---
id: so-02-20260111053925
title: "Hive Mind Architecture"
type: concept
domain: SO
domain_range: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-03-26
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Hive Mind Architecture — 9-agent system powered by OpenClaw"
---

# Hive Mind Architecture

> 9-Agent System for One Mind — powered by OpenClaw

---

## Overview

The Legacy AI Hive Mind consists of **9 agents** that manage all aspects of Zeus Delacruz's life and work. Each agent is an OpenClaw-powered stateful AI with persistent memory, shared context, and modular skills.

---

## Architecture Diagram

```
Zeus (human)
  ↕ Telegram / iMessage / Voice / iOS App / Discord
Legacy ⚡ (Sovereign Orchestrator)
  ├── Spartan ⚔️   ← Body & Physical Mastery (HP 27-28)
  ├── Oracle 🔮    ← Mind, Money & Meaning (HP 25-26, 29-34)
  ├── Eden 🌳      ← Sustainability & Agriculture (LE)
  ├── Empire 👑    ← Business & Wealth (GE 75-84)
  ├── Haven 🌿     ← Home, Family & Estate (LE 50-58)
  ├── Guardian 🛡️  ← Security & Compliance (SO 11 + HP 31 + LE 57)
  ├── Grid ⚙️      ← Infrastructure Ops (SO 00-12)
  └── Trinity 🔺   ← Community & Partnerships (GE 75 + Community)

  + NATS → Agno Workers (stateless):
    researcher, coder, analyst, web_automator, db_manager
```

---

## The 9 Agents

| Agent | Role | Domain Range |
|-------|------|-------------|
| **Legacy** ⚡ | Sovereign Orchestrator | All domains (00-99) |
| **Spartan** ⚔️ | Body & Physical Mastery | HP 27-28 |
| **Oracle** 🔮 | Mind, Money & Meaning | HP 25-26, 29-34 |
| **Eden** 🌳 | Sustainability & Agriculture | LE (farming) |
| **Empire** 👑 | Business & Wealth | GE 75-84 |
| **Haven** 🌿 | Home, Family & Estate | LE 50-58 |
| **Guardian** 🛡️ | Security & Compliance | SO 11 + HP 31 + LE 57 |
| **Grid** ⚙️ | Infrastructure Ops | SO 00-12 |
| **Trinity** 🔺 | Community & Partnerships | GE 75 + Community |

---

## Shared Memory Architecture

All agents share access to core memory blocks while maintaining domain-specific context:

### Core Memory (All Agents)

| Block | Purpose | Update Frequency |
|-------|---------|------------------|
| **Zeus Profile** | Identity, values, preferences | Rarely |
| **Current Context** | Today's priorities, active focus | Daily |
| **System State** | Agent status, health metrics | Real-time |
| **Cross-Domain Queue** | Inter-agent messages | As needed |

### Domain Memory (Per Agent)

| Agent | Memory Block | Contents |
|-------|--------------|----------|
| Spartan | HP Body Context | Fitness goals, workouts, body metrics |
| Oracle | HP Mind Context | Learning goals, finances, growth plans |
| Eden | LE Farm Context | Crops, seasons, sustainability plans |
| Empire | GE Biz Context | Business metrics, active deals, revenue |
| Haven | LE Home Context | Family state, home status, estate plans |
| Guardian | Security Context | Credentials, compliance, threat status |
| Grid | Infra Context | Infrastructure status, services, uptime |
| Trinity | Community Context | Partnerships, community, outreach |

---

## Skills System

Each agent has modular **skills** attached as archival memory. Skills expand agent capabilities without requiring new agents.

### How Skills Work

1. Skills are markdown documents describing specialized knowledge
2. Skills are loaded into agent archival memory via OpenClaw API
3. Agents reference skills when handling relevant requests
4. Skills can be added, updated, or removed without redeploying agents

### Skill Format

```markdown
# Skill: [Skill Name]

## Purpose
What this skill enables the agent to do.

## Knowledge
Domain-specific information, procedures, rules.

## Actions
What the agent should do when this skill is invoked.

## Examples
Sample inputs and expected outputs.
```

---

## Agent Routing

Legacy routes all requests to the appropriate domain agent:

| If request involves... | Route to |
|------------------------|----------|
| Fitness, body, physical training | Spartan |
| Mind, learning, money, personal growth | Oracle |
| Agriculture, sustainability, food growing | Eden |
| Business, revenue, products, marketing | Empire |
| Home, family, estate, legacy planning | Haven |
| Security, compliance, credentials | Guardian |
| Tech, automation, infrastructure, AI | Grid |
| Community, partnerships, outreach | Trinity |
| Cross-domain or unclear | Legacy handles directly |

---

## Autonomy Levels

| Level | Description | Examples |
|-------|-------------|----------|
| **L3** | Full autonomy | Reading files, routing requests |
| **L2** | Act then notify | Creating tasks, updating metrics |
| **L1** | Request approval | Financial decisions, external comms |
| **L0** | Human only | Legal matters, health emergencies |

---

## Implementation

### OpenClaw Configuration

- **Platform**: OpenClaw (self-hosted)
- **Gateway**: OpenClaw 2026.3.7
- **Model**: azure-openai/gpt-4o (primary), gpt-5.2 (available)
- **Memory**: Shared blocks + agent-specific archival
- **Tools**: MCP servers (Todoist, GitHub, Obsidian, Home Assistant)

### Agent Creation

Agents are created via OpenClaw workspace configuration:

```bash
python 00-Infrastructure/scripts/openclaw_agent_setup.py --all
```

---

## Agent Specifications

- [[AGENT-ASSETS]] — Full 9-agent roster with workspace details
- [[00-24 SO (Sovereign Operations)/03 Intelligence (AI, Memory & Tools)/Agents/Legacy-ai/_Agent]] — Legacy: Sovereign Orchestrator

---

## Skills Library

Skills are defined in: [[01.03 AI communications - (Skills W Tools)]]

---

## Change Log

| Date | Change |
|------|--------|
| 2026-03-26 | Updated to 9-agent architecture, Letta→OpenClaw |
| 2024-12-23 | Consolidated to 5-agent architecture |
| 2024-11-08 | Previous 12-agent design (archived) |

---

#agent #legacy-ai #architecture #hive-mind
