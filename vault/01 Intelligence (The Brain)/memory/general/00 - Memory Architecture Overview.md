---
id: ui-02-20260111053925
title: "Legacy AI - Memory Block Architecture"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Legacy AI - Memory Block Architecture"
---

# Legacy AI - Memory Block Architecture

## Overview
This folder contains the complete memory architecture for the Legacy AI multi-agent system. Memory is organized in a 3-tier hierarchy designed for Letta's shared memory block system.

## Memory Hierarchy

```
┌─────────────────────────────────────────────────────┐
│         CORE SHARED MEMORY (All Agents)             │
│  - Legacy Codex                                     │
│  - User Profile                                     │
│  - System State                                     │
│  - Cross-Domain Context                             │
└─────────────────────────────────────────────────────┘
                        ↓
        ┌───────────────┼───────────────┐
        ↓               ↓               ↓
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│  DOMAIN       │ │  DOMAIN       │ │  DOMAIN       │
│  MEMORY       │ │  MEMORY       │ │  MEMORY       │
│  (Pillar)     │ │  (Pillar)     │ │  (Pillar)     │
└───────────────┘ └───────────────┘ └───────────────┘
        ↓               ↓               ↓
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│  AGENT        │ │  AGENT        │ │  AGENT        │
│  MEMORY       │ │  MEMORY       │ │  MEMORY       │
│  (Specific)   │ │  (Specific)   │ │  (Specific)   │
└───────────────┘ └───────────────┘ └───────────────┘
```

## File Organization

### Core Shared Memory
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/01 - Legacy Codex]] - Foundational principles and decision frameworks
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/02 - User Profile]] - Zeus's patterns, preferences, and current context
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/03 - System State]] - Real-time status of all agents and infrastructure
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/04 - Cross-Domain Context]] - How domains interact and dependencies

### Domain Memory
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/05 - Holistic Performance Domain]] - Performance, health, cognitive optimization
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/06 - Generational Empire Domain]] - Business operations and portfolio
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/07 - Legacy Evolution Domain]] - Family, homestead, multi-location operations
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/08 - Innovative Technology Domain]] - Infrastructure, robotics, tech stack

### Agent-Specific Memory (Examples)
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/09 - Intelligence Agent Memory]] - Information processing operations
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/10 - Infrastructure Agent Memory]] - System monitoring and maintenance

## Implementation Notes

**Format Choice:** Markdown over JSON
- LLMs understand natural language better than structured data
- More human-readable and editable
- More token-efficient
- Flexible structure without rigid schemas

**Character Limits:** Letta blocks support ~20,000 characters
- Core blocks: 3,000-8,000 chars each
- Domain blocks: 4,000-5,000 chars each
- Agent blocks: 5,000-6,000 chars each

**Access Control:** See individual files for read/write permissions per agent

**Update Protocols:**
- Core Shared: Real-time to daily depending on block
- Domain Memory: Updated by domain agents as needed
- Agent-Specific: Full autonomy for each agent's own memory

## Next Steps
1. Implement Core Memory Blocks in Letta
2. Create agent roster with memory assignments
3. Build unified task engine
4. Design unified RAG system
5. Connect tool registry
6. Build monitoring dashboard

---
*Created: 2025-11-08*
*Status: Architecture defined, ready for implementation*
