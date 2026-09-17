---
id: ui-15-20260111053926
title: "One Mind Skills Library"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/15]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "One Mind Skills Library"
---

# One Mind Skills Library

> Modular capabilities for Hive Mind agents

---

## What Are Skills?

Skills are markdown documents that expand agent capabilities without requiring new agents. Each skill provides:

- **Domain knowledge** — Specialized information the agent needs
- **Procedures** — Step-by-step instructions for tasks
- **Context** — Background information for decision-making
- **Examples** — Sample inputs and expected outputs

Skills are loaded into Letta agent archival memory and can be added, updated, or removed without redeploying agents.

---

## Skills by Agent

### LEGACY (Mother AI)

| Skill | File | Purpose |
|-------|------|---------|
| System Coordination | `LEGACY/system-coordination.md` | Cross-agent communication |
| Approval Management | `LEGACY/approval-management.md` | HITL queue and rules |
| Morning Brief | `LEGACY/morning-brief.md` | Daily briefing generation |
| Evening Review | `LEGACY/evening-review.md` | Daily review protocol |
| Inbox Processing | `LEGACY/inbox-processing.md` | Triage incoming items |

### HP Agent (Holistic Performance)

| Skill | File | Purpose |
|-------|------|---------|
| Health Monitoring | `HP/health-monitoring.md` | Vitals, symptoms, appointments |
| Fitness Programming | `HP/fitness-programming.md` | Workouts, tracking, recovery |
| Nutrition Tracking | `HP/nutrition-tracking.md` | Meals, supplements, diet |
| Sleep Optimization | `HP/sleep-optimization.md` | Sleep tracking, protocols |
| Habit Management | `HP/habit-management.md` | Habits, streaks, stacking |
| Energy Management | `HP/energy-management.md` | Energy levels, peak times |
| Learning Coach | `HP/learning-coach.md` | Reading, courses, skills |

### LE Agent (Legacy Evolution)

| Skill | File | Purpose |
|-------|------|---------|
| Family Coordination | `LE/family-coordination.md` | Events, communication |
| Home Management | `LE/home-management.md` | Household, maintenance |
| Homestead Operations | `LE/homestead-operations.md` | Property, agriculture |
| Legacy Planning | `LE/legacy-planning.md` | Heritage, traditions |
| Wealth Stewardship | `LE/wealth-stewardship.md` | Family finances, estate |
| Children's Education | `LE/childrens-education.md` | Learning, development |

### GE Agent (Generational Endeavor)

| Skill | File | Purpose |
|-------|------|---------|
| Business Operations | `GE/business-operations.md` | Day-to-day ops, projects |
| Revenue Tracking | `GE/revenue-tracking.md` | MRR, pipeline, metrics |
| Product Management | `GE/product-management.md` | Development, launches |
| Marketing Strategy | `GE/marketing-strategy.md` | Content, campaigns |
| Customer Success | `GE/customer-success.md` | CRM, community |
| Compliance Management | `GE/compliance-management.md` | Contracts, legal |

### Grid Agent (Innovative Technology)

| Skill | File | Purpose |
|-------|------|---------|
| Infrastructure Monitoring | `IT/infrastructure-monitoring.md` | System health, uptime |
| Security Operations | `IT/security-operations.md` | Threats, access control |
| AI Agent Management | `IT/ai-agent-management.md` | Letta agents, memory |
| Automation Engineering | `IT/automation-engineering.md` | N8N, scripts, cron |
| Knowledge Architecture | `IT/knowledge-architecture.md` | Vault organization |
| DevOps | `IT/devops.md` | Deployments, CI/CD |

---

## Skill File Format

Each skill follows this template:

```markdown
# Skill: [Name]

## Purpose
What this skill enables the agent to do.

## Knowledge
Domain-specific information, facts, rules.

## Procedures
Step-by-step instructions for common tasks.

## Context
Background information for decision-making.

## Examples
Sample inputs and expected agent responses.

## Integration
- Tools: Which MCP tools this skill uses
- Memory: Which memory blocks to reference
- Escalation: When to escalate to LEGACY
```

---

## Loading Skills into Letta

Skills are loaded via the Letta API:

```python
from letta_client import Letta

client = Letta(token=LETTA_API_KEY)

# Read skill file
with open("Skills/HP/habit-management.md") as f:
    skill_content = f.read()

# Add to agent archival memory
client.agents.archival_memory.insert(
    agent_id="hp-agent-id",
    content=skill_content
)
```

See `00-Infrastructure/scripts/load_skills.py` for the full implementation.

---

## Adding New Skills

1. Create markdown file in appropriate agent folder
2. Follow the skill template format
3. Run the skill loader script
4. Update this overview document

---

#skills #legacy-ai #hive-mind
