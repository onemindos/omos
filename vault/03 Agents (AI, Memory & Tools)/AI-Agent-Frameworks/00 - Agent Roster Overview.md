---
id: ui-02-20260111053925
title: "Legacy AI - Agent Roster Overview (UPDATED)"
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
ai_summary: "Legacy AI - Agent Roster Overview (UPDATED)"
---

# Legacy AI - Agent Roster Overview (UPDATED)

**Purpose:** Complete documentation of all agents in the Legacy AI multi-agent system  
**Status:** Architecture finalized - Communications absorbed into Legacy  
**Last Updated:** 2025-11-08

---

## System Architecture

```
                    ┌─────────────────────────────┐
                    │                             │
                    │   LEGACY                    │
                    │   (Master Orchestrator)     │
                    │   + Communications          │
                    │                             │
                    └─────────────┬───────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
┌───────▼──────┐  ┌───────────────▼──────────┐  ┌──────────▼─────────┐
│ HOLISTIC     │  │ GENERATIONAL             │  │ LEGACY             │
│ PERFORMANCE  │  │ ENTREPRENEURSHIP         │  │ EVOLUTION          │
│              │  │                          │  │                    │
│ 3 Agents     │  │ 3 Agents                 │  │ 3 Agents           │
└──────────────┘  └──────────────────────────┘  └────────────────────┘
                                  │
                         ┌────────▼────────┐
                         │ INNOVATIVE      │
                         │ TECHNOLOGY      │
                         │                 │
                         │ 2 Agents        │
                         └─────────────────┘
```

**Total:** 1 Master Orchestrator + 11 Specialized Agents = **12 Agents**

---

## Agent Roster Summary

### Master Orchestrator (1)
**00 - LEGACY - Master Orchestrator**
- Strategic coordination across all domains
- Cross-domain synthesis and priority management
- User interface (voice via LiveKit, chat, Obsidian)
- **External communications management (email, scheduling, meetings, relationship management)**
- System health oversight

### Holistic Performance Domain (3 Agents)
**01 - HP - Intelligence Agent** - Information processing, distillation, RAG preparation  
**02 - HP - Performance Agent** - Health, fitness, cognitive performance optimization  
**03 - HP - Finance Agent** - Personal finance, portfolio tracking, market intelligence

### Generational Endeavor Domain (3 Agents)
**04 - GE - Business Operations Agent** - Day-to-day operations, CRM, project management  
**05 - GE - Founder Assistant Agent** - Strategic support, meeting prep, executive communications  
**06 - GE - Capital Systems Agent** - Business finance, fundraising, investor relations

### Legacy Evolution Domain (3 Agents)
**07 - LE - Home Command Agent** - Multi-location automation, family operations, household management  
**08 - LE - Homestead Operations Agent** - Farm/land management, agriculture, livestock, equipment  
**09 - LE - Family Wealth Agent** - Generational wealth planning, estate management, financial education

### Innovative Technology Domain (2 Agents)
**10 - IT - Infrastructure Agent** - IT infrastructure monitoring, DevOps, security operations  
**11 - IT - Robotics & Automation Agent** - Robotics systems, drone operations, automation workflows

---

## Design Principles

### Clean Architecture: Legacy + Four Pillars

**Everything maps to either the Master Orchestrator or one of the four Legacy Codex pillars.**

**Why Communications Is in Legacy:**
- Legacy is the master interface between Zeus and everything else
- Email, scheduling, meetings are extensions of orchestration
- Keeps architecture clean - no orphan domains outside the four pillars
- Voice + email + scheduling = unified communication interface

**No "External Interface" domain needed** - it's all orchestration.

---

### Why 11 Specialized Agents (Not More or Less)

**Your operational scope requires this level of specialization:**
- Multiple businesses across portfolio companies
- Personal/family operations across multiple properties
- Technical infrastructure (self-hosted, multi-location)
- 100+ daily information sources
- Robotics/automation systems
- Financial operations and investments

**11 specialized agents provide:**
- Deep domain expertise with persistent context
- Parallel processing (agents work simultaneously)
- Graceful degradation (one agent failing doesn't break others)
- Clear ownership and accountability
- Easier debugging and optimization

**Condensed from 25+ potential agents:**
- "Office Manager" + "Home Assistant" → **Home Command Agent**
- "Eagle Eye" + "Overwatch" + "Hacker" → **Infrastructure Agent** (security operations mode)
- "Zoho Manager" + "Circle Manager" → **Business Operations Agent**
- "Communications" → **Absorbed into Legacy**

---

## Memory Integration

Each agent has access to:
1. **Core Shared Memory** (all agents)
   - [[01 - Legacy Codex]]
   - [[02 - User Profile]]
   - [[03 - System State]]
   - [[04 - Cross-Domain Context]]

2. **Domain Memory** (pillar-specific)
   - HP agents share [[05 - Holistic Performance Domain]]
   - GE agents share [[06 - Generational Endeavor Domain]]
   - LE agents share [[07 - Legacy Evolution Domain]]
   - Grid agents share [[08 - Innovative Technology Domain]]

3. **Agent-Specific Memory** (private to each agent)
   - Skills and capabilities
   - Job context and operational data
   - Learning patterns and optimization notes
   - Operational history and statistics

See: [[Memory Block Architecture]] folder for full details

---

## Autonomous Operation Design

**Proactive Triggers:**
1. **Time-Based:** Daily briefings, weekly reviews, monthly reports
2. **Event-Based:** Alert thresholds, system state changes, incoming emails
3. **Pattern-Based:** Detected trends, anomalies, opportunities
4. **Context-Based:** Meeting prep (1hr before), follow-ups (24hr after)

**Self-Learning Mechanisms:**
1. **Feedback Loops:** Your responses inform agent behavior
2. **Pattern Recognition:** What you engage with vs. ignore
3. **Optimization:** Refine timing, content, priorities based on usage
4. **Tool Discovery:** Agents learn new capabilities over time

**Autonomy Levels:**
- **Fully Autonomous:** Health monitoring, system checks, data collection, email triage
- **Suggest & Execute:** Daily briefs, routine task creation, email drafts
- **Suggest & Wait:** Strategic decisions, financial moves, external communications (send after approval)
- **Human Required:** Major changes, risk decisions, family matters

See: [[01 - Legacy Codex]] for decision escalation framework

---

## Agent Coordination Patterns

### Direct Reporting to Legacy
All agents report directly to Legacy (no middle managers). Legacy routes requests and coordinates cross-agent workflows.

### Common Coordination Patterns

**Pattern 1: Intelligence → Finance**
- Intelligence Agent processes market news
- Finance Agent receives relevant intel for portfolio decisions
- Coordination: Intelligence tags financial content, Finance subscribes to tags

**Pattern 2: Performance → All Agents**
- Performance Agent monitors your energy state
- All agents adjust task urgency based on cognitive capacity
- Coordination: Performance publishes state, others subscribe

**Pattern 3: Legacy → Founder Assistant (Meeting Prep)**
- Legacy receives meeting request via email
- Founder Assistant preps context and talking points
- Legacy coordinates delivery 1 hour before meeting
- Coordination: Sequential handoff via task system

**Pattern 4: Homestead → Home Command**
- Homestead Operations needs irrigation control
- Home Command executes physical automation
- Coordination: Homestead requests, Home Command executes

**Pattern 5: Infrastructure → All Agents**
- Infrastructure detects system degradation
- All agents receive capacity constraints
- Coordination: Infrastructure broadcasts alerts, agents adapt

See: [[04 - Cross-Domain Context]] for integration flows

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1-2)
1. Implement Legacy (master orchestrator with communications)
2. Set up core shared memory blocks
3. Create Infrastructure Agent (must monitor system health)
4. Verify base system is stable

### Phase 2: Data Pipeline (Week 3)
1. Implement Intelligence Agent (information processing)
2. Connect to RSS feeds and sources
3. Test RAG document preparation
4. Verify daily brief delivery

### Phase 3: Core Operations (Week 4-5)
1. Test Legacy's communications functions (email triage, scheduling)
2. Implement Performance Agent (health tracking)
3. Implement Finance Agent (portfolio monitoring)
4. Test cross-agent coordination

### Phase 4: Business Layer (Week 6-7)
1. Implement Business Operations Agent
2. Implement Founder Assistant Agent
3. Implement Capital Systems Agent
4. Connect to business systems (Zoho, Circle, Shopify)

### Phase 5: Home & Family (Week 8)
1. Implement Home Command Agent
2. Implement Homestead Operations Agent
3. Implement Family Wealth Agent
4. Connect to Home Assistant instances

### Phase 6: Advanced Automation (Week 9-10)
1. Implement Robotics & Automation Agent
2. Connect to drone systems
3. Build advanced automation workflows
4. Test end-to-end system

---

## Quality Metrics

**How we'll know it's working:**

**Week 1:** Infrastructure Agent monitoring all systems, daily health reports  
**Week 3:** Morning intelligence brief delivered automatically at 6am  
**Week 4:** Email triage and draft responses working  
**Week 5:** All agents operational, coordinating autonomously  
**Week 8:** Proactive behaviors working (meeting prep, task creation, alerts)  
**Week 10:** System running 80%+ autonomously, only strategic escalations

**Success Criteria:**
- 95%+ uptime for critical agents
- <5 minute response time for queries
- 80%+ tasks handled autonomously
- Zero missed daily briefings
- Inbox zero achievable (good email filtering)
- <24 hour response time for important emails
- Proactive alerts catching issues before they become problems

---

## Document Index

### Agent Profiles (11 Specialized Agents)
**Master Orchestrator:**
- [[00 - LEGACY - Master Orchestrator]] - Includes communications management

**Holistic Performance (3):**
- [[01 - HP - Intelligence Agent]]
- [[02 - HP - Performance Agent]]
- [[03 - HP - Finance Agent]]

**Generational Endeavor (3):**
- [[04 - GE - Business Operations Agent]]
- [[05 - GE - Founder Assistant Agent]]
- [[06 - GE - Capital Systems Agent]]

**Legacy Evolution (3):**
- [[07 - LE - Home Command Agent]]
- [[08 - LE - Homestead Operations Agent]]
- [[09 - LE - Family Wealth Agent]]

**Innovative Technology (2):**
- [[10 - IT - Infrastructure Agent]]
- [[11 - IT - Robotics & Automation Agent]]

### Supporting Documentation
- [[Agent Naming Convention - FINAL]] - Naming system and quick reference
- [[Memory Block Architecture]] - Memory system design (separate folder)
- [[Task Engine Architecture]] - Task orchestration (to be created)
- [[RAG System Architecture]] - Knowledge retrieval (to be created)
- [[Tool Registry]] - API integrations (to be created)
- [[Dashboard Design]] - Monitoring interface (to be created)

---

## Next Steps

1. ✅ **Agent Roster Finalized** - Clean architecture with Legacy + Four Pillars
2. ✅ **Memory Architecture Documented** - Foundation complete
3. ⏭️ **Build Task Engine** - How agents coordinate work
4. ⏭️ **Deploy Agents** - One by one, following roadmap
5. ⏭️ **Test & Iterate** - Refine based on real-world usage

---

*Created: 2025-11-08*  
*Updated: 2025-11-08 (Communications absorbed into Legacy - Final architecture)*  
*Status: Architecture finalized, ready for implementation*
