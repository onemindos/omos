---
id: ui-02-20260111053925
title: "Hive Mind Architecture"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Hive Mind Architecture"
---

# Hive Mind Architecture

> 5-Agent System for One Mind

---

## Overview

The Legacy AI Hive Mind consists of **5 core agents** that manage all aspects of Zeus Delacruz's life and work. Each agent is a Letta-powered stateful AI with persistent memory, shared context, and modular skills.

---

## Architecture Diagram

```
                    ┌─────────────────────────────┐
                    │       LEGACY (Mother AI)    │
                    │    Central Coordinator      │
                    │  + Shared Memory Blocks     │
                    └─────────────┬───────────────┘
                                  │
        ┌─────────────┬───────────┼───────────┬─────────────┐
        │             │           │           │             │
        ▼             ▼           ▼           ▼             ▼
   ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
   │   HP    │  │   LE    │  │   GE    │  │   IT    │
   │ Agent   │  │ Agent   │  │ Agent   │  │ Agent   │
   └─────────┘  └─────────┘  └─────────┘  └─────────┘1
```

---

## The 5 Agents

| Agent | Role | Domain | Pillar |
|-------|------|--------|--------|
| **LEGACY** | Mother AI / Orchestrator | Cross-domain coordination | Command |
| **HP Agent** | Performance Conductor | Personal optimization | 10-19 |
| **LE Agent** | Legacy Conductor | Family & heritage | 20-29 |
| **GE Agent** | Business Conductor | Entrepreneurship | 30-39 |
| **IT Agent** | Technology Conductor | Infrastructure | 40-49 |

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

### Pillar Memory (Domain Agents)

| Agent | Memory Block | Contents |
|-------|--------------|----------|
| HP | HP Context | Health goals, habits, metrics |
| LE | LE Context | Family state, home status, legacy plans |
| GE | GE Context | Business metrics, active deals, revenue |
| IT | IT Context | Infrastructure status, projects, security |

---

## Skills System

Each agent has modular **skills** attached as archival memory. Skills expand agent capabilities without requiring new agents.

### How Skills Work

1. Skills are markdown documents describing specialized knowledge
2. Skills are loaded into agent archival memory via Letta API
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

LEGACY (Mother AI) routes all requests to the appropriate pillar agent:

| If request involves... | Route to |
|------------------------|----------|
| Health, fitness, habits, personal growth | HP Agent |
| Family, home, homestead, legacy | LE Agent |
| Business, revenue, products, marketing | GE Agent |
| Tech, automation, infrastructure, AI | IT Agent |
| Cross-domain or unclear | LEGACY handles directly |

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

### Letta Configuration

- **Platform**: Letta Cloud (or self-hosted)
- **Model**: Claude Sonnet 4.5 (primary), Haiku (fast tasks)
- **Memory**: Shared blocks + agent-specific archival
- **Tools**: MCP servers (Todoist, GitHub, Obsidian)

### Agent Creation

Agents are created via Python script:

```bash
python 00-Infrastructure/scripts/letta_agent_setup.py --all
```

See [[00-Infrastructure/scripts/letta_agent_setup.py]] for full implementation.

---

## Agent Specifications

- [[40-49 IT (Innovative Technology)/41 Legacy-AI (Intelligent Core)/Agents/Mother-AI]] — Central orchestrator
- [[HP-Agent]] — Holistic Performance conductor
- [[LE-Agent]] — Legacy Evolution conductor
- [[GE-Agent]] — Generational Entrepreneurship conductor
- [[IT-Agent]] — Innovative Technology conductor

---

## Skills Library

Skills are defined in: [[01.03 AI communications - (Skills W Tools)]]

---

## Change Log

| Date | Change |
|------|--------|
| 2024-12-23 | Consolidated to 5-agent architecture |
| 2024-11-08 | Previous 12-agent design (archived) |

---

#agent #legacy-ai #architecture #hive-mind
