---
id: ui-02-20260111053925
title: "Agent Naming Convention & Quick Reference"
type: reference
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Agent Naming Convention & Quick Reference"
---

# Agent Naming Convention & Quick Reference

**Last Updated:** 2025-11-08 (Updated - Communications absorbed into Legacy)

---

## Domain Prefix System

All agent files use a domain prefix for instant identification:

| Prefix | Domain | Color | Description |
|--------|--------|-------|-------------|
| **LEGACY** | Master Orchestrator | Red | Strategic coordination, cross-domain synthesis, **all external communications** |
| **HP** | Holistic Performance | Blue | Personal optimization, health, finance, intelligence |
| **GE** | Generational Entrepreneurship | Green | Business operations, capital, executive support |
| **LE** | Legacy Evolution | Purple | Family, homestead, wealth, multi-location operations |
| **IT** | Innovative Technology | Orange | Infrastructure, robotics, technical systems |

---

## Complete Agent List (Organized by Domain)

### Master Orchestrator
**00 - LEGACY - Master Orchestrator**
- Strategic coordination across all domains
- Cross-domain synthesis and priority management
- User interface (voice, chat, Obsidian)
- **External communications (email, scheduling, meetings)**

---

### Holistic Performance Domain (Blue) - 3 Agents

**01 - HP - Intelligence Agent**
- Information processing from 100+ sources
- Daily briefings and RAG preparation
- Trend detection and research

**02 - HP - Performance Agent**
- Health and fitness monitoring
- Biometric tracking (sleep, HRV, activity)
- Energy state management

**03 - HP - Finance Agent**
- Personal portfolio management
- Market intelligence and analysis
- Financial planning and tax optimization

---

### Generational Entrepreneurship Domain (Green) - 3 Agents

**04 - GE - Business Operations Agent**
- Day-to-day operations across portfolio companies
- CRM, project management, customer relationships
- Sales pipeline and operational reporting

**05 - GE - Founder Assistant Agent**
- Executive support and meeting preparation
- Research, analysis, and communication drafting
- Strategic decision support

**06 - GE - Capital Systems Agent**
- Business finance and financial modeling
- Fundraising and investor relations
- Compliance and reporting

---

### Legacy Evolution Domain (Purple) - 3 Agents

**07 - LE - Home Command Agent**
- Multi-location automation (office, home, farm, warehouse)
- Security monitoring and climate control
- Home Assistant orchestration

**08 - LE - Homestead Operations Agent**
- Farm and land management
- Livestock, crops, equipment
- Weather monitoring and seasonal planning

**09 - LE - Family Wealth Agent**
- Generational wealth planning
- Estate management and succession
- Family financial education

---

### Innovative Technology Domain (Orange) - 2 Agents

**10 - IT - Infrastructure Agent**
- System monitoring and DevOps
- AWS VMs, containers, databases
- Security scanning and backup management

**11 - IT - Robotics & Automation Agent**
- Drone operations and missions
- Ground robot coordination (future)
- Automation workflows (n8n, Node-RED)

---

## Quick Reference Matrix

| # | Prefix | Agent | Primary Function | Autonomy Level |
|---|--------|-------|------------------|----------------|
| 00 | LEGACY | Master Orchestrator | Coordination + Communications | Mixed |
| 01 | HP | Intelligence | Information processing | High |
| 02 | HP | Performance | Health tracking | High |
| 03 | HP | Finance | Portfolio management | Medium |
| 04 | GE | Business Operations | Daily operations | Medium |
| 05 | GE | Founder Assistant | Executive support | Low |
| 06 | GE | Capital Systems | Business finance | Low |
| 07 | LE | Home Command | Automation | High |
| 08 | LE | Homestead Operations | Farm management | Medium |
| 09 | LE | Family Wealth | Wealth planning | Low |
| 10 | IT | Infrastructure | System monitoring | High |
| 11 | IT | Robotics | Drone operations | Medium |

**Total: 1 Master Orchestrator + 11 Specialized Agents = 12 Agents**

---

## Architecture Summary

**Master Orchestrator:**
- **LEGACY** (1 agent) - Handles strategic coordination AND all external communications

**Four Pillars:**
- **HP - Holistic Performance** (3 agents) - Intelligence, Performance, Finance
- **GE - Generational Entrepreneurship** (3 agents) - Business Ops, Founder Assist, Capital
- **LE - Legacy Evolution** (3 agents) - Home Command, Homestead, Family Wealth
- **IT - Innovative Technology** (2 agents) - Infrastructure, Robotics

**Perfect symmetry:** Everything maps to either Legacy or one of the four pillars.

---

## Autonomy Level Key

- **High:** Operates 80%+ autonomously, minimal human intervention
- **Medium:** 50-80% autonomous, regular check-ins and approvals
- **Low:** Suggest-and-wait model, most actions require approval

---

## Why Communications Is in Legacy

**Architectural Reasoning:**
- Legacy is the master interface between Zeus and the external world
- Email, scheduling, and communications are extensions of Legacy's orchestration role
- Keeps architecture clean: everything is either orchestration or one of the four pillars
- Prevents "orphan" domains that don't map to the Legacy Codex structure

**Legacy handles:**
- Voice interface (LiveKit)
- Email management (triage, draft, send)
- Calendar and scheduling
- External relationship management
- All agent coordination

---

## File Location

All agent profiles located in:  
`Legacy AI/Infrastructure - Agents/`

---

*Created: 2025-11-08*  
*Updated: 2025-11-08 (Communications absorbed into Legacy)*
