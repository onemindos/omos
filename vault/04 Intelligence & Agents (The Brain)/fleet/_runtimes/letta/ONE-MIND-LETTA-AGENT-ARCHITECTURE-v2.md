---
id: ui-02-20260111053925
title: "🧠 ONE MIND — Complete Letta Agent Architecture v2.0"
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
ai_summary: "🧠 ONE MIND — Complete Letta Agent Architecture v2.0"
---

# 🧠 ONE MIND — Complete Letta Agent Architecture v2.0

> **Version:** 2.0  
> **Created:** January 2025  
> **Author:** Zeus Delacruz + Legacy AI  
> **Status:** Master Architecture Document

---

## Executive Summary

This document defines the complete multi-agent system for One Mind using Letta Cloud as the orchestration layer. Each of the 4 pillars (HP, LE, GE, IT) has equal representation with specialized agents. Command layer provides cross-pillar coordination and human-in-the-loop governance.

---

## 🏗️ System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                               HUMAN INTERFACE LAYER                              │
│                                                                                  │
│     ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   │
│     │ Telegram │   │ Todoist  │   │ Obsidian │   │  Voice   │   │  GitHub  │   │
│     │  (Chat)  │   │ (Mobile) │   │  (Vault) │   │ (Future) │   │  (Sync)  │   │
│     └────┬─────┘   └────┬─────┘   └────┬─────┘   └────┬─────┘   └────┬─────┘   │
└──────────┼──────────────┼──────────────┼──────────────┼──────────────┼──────────┘
           │              │              │              │              │
           └──────────────┴──────────────┼──────────────┴──────────────┘
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            00 - COMMAND LAYER                                    │
│                          (Orchestration & Governance)                            │
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │                         MOTHER AI — COORDINATOR                            │ │
│  │                                                                            │ │
│  │  Role: Cross-domain routing, approval workflows, memory sync              │ │
│  │  Model: claude-sonnet-4-5                                                 │ │
│  │  Memory Blocks: zeus_profile, global_state, approval_queue                │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│                                         │                                        │
│          ┌──────────────────────────────┼──────────────────────────────┐        │
│          │                              │                              │        │
│          ▼                              ▼                              ▼        │
│  ┌───────────────┐            ┌───────────────┐            ┌───────────────┐   │
│  │   PROTOCOL    │            │   APPROVAL    │            │   MEMORY      │   │
│  │    GUARDIAN   │            │    MANAGER    │            │   CURATOR     │   │
│  │               │            │               │            │               │   │
│  │  • SOP        │            │  • Human-in-  │            │  • Archival   │   │
│  │    Execution  │            │    the-loop   │            │    Management │   │
│  │  • Workflow   │            │  • Approvals  │            │  • Sleep-time │   │
│  │    Validation │            │  • Escalation │            │    Compute    │   │
│  └───────────────┘            └───────────────┘            └───────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
                                         │
          ┌──────────────────────────────┼──────────────────────────────┐
          │                              │                              │
          │                              │                              │
          ▼                              ▼                              ▼
┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐
│   10-19 HP PILLAR   │ │   20-29 LE PILLAR   │ │   30-39 GE PILLAR   │
│ Holistic Performance│ │  Legacy Evolution   │ │ Generational Enter. │
│                     │ │                     │ │                     │
│ ┌─────────────────┐ │ │ ┌─────────────────┐ │ │ ┌─────────────────┐ │
│ │  HP-CONDUCTOR   │ │ │ │  LE-CONDUCTOR   │ │ │ │  GE-CONDUCTOR   │ │
│ │  (Pillar Lead)  │ │ │ │  (Pillar Lead)  │ │ │ │  (Pillar Lead)  │ │
│ └─────────────────┘ │ │ └─────────────────┘ │ │ └─────────────────┘ │
│         │           │ │         │           │ │         │           │
│    ┌────┼────┐      │ │    ┌────┼────┐      │ │    ┌────┼────┐      │
│    ▼    ▼    ▼      │ │    ▼    ▼    ▼      │ │    ▼    ▼    ▼      │
│ ┌───┐ ┌───┐ ┌───┐   │ │ ┌───┐ ┌───┐ ┌───┐   │ │ ┌───┐ ┌───┐ ┌───┐   │
│ │HLT│ │FIT│ │HAB│   │ │ │IDN│ │REL│ │GRW│   │ │ │VEN│ │OPS│ │FIN│   │
│ │   │ │   │ │   │   │ │ │   │ │   │ │   │   │ │ │   │ │   │ │   │   │
│ └───┘ └───┘ └───┘   │ │ └───┘ └───┘ └───┘   │ │ └───┘ └───┘ └───┘   │
│  +3    +2    +2     │ │  +2    +2    +3     │ │  +4    +3    +3     │
└─────────────────────┘ └─────────────────────┘ └─────────────────────┘
                                         │
                              ┌──────────┴──────────┐
                              ▼                     ▼
                    ┌─────────────────────┐
                    │   40-49 IT PILLAR   │
                    │ Innovative Technology│
                    │                     │
                    │ ┌─────────────────┐ │
                    │ │  IT-CONDUCTOR   │ │
                    │ │  (Pillar Lead)  │ │
                    │ └─────────────────┘ │
                    │         │           │
                    │    ┌────┼────┐      │
                    │    ▼    ▼    ▼      │
                    │ ┌───┐ ┌───┐ ┌───┐   │
                    │ │INF│ │AUT│ │DEV│   │
                    │ │   │ │   │ │   │   │
                    │ └───┘ └───┘ └───┘   │
                    │  +3    +3    +3     │
                    └─────────────────────┘
```

---

## 📊 Agent Count Summary

| Layer | Agent Count | Description |
|-------|-------------|-------------|
| **Command (00)** | 4 | Mother AI + 3 specialized |
| **HP (10-19)** | 8 | Conductor + 7 specialists |
| **LE (20-29)** | 8 | Conductor + 7 specialists |
| **GE (30-39)** | 11 | Conductor + 10 specialists |
| **IT (40-49)** | 10 | Conductor + 9 specialists |
| **TOTAL** | **41 Agents** | Full cognitive workforce |

---

# 🎛️ COMMAND LAYER (00)

## Mother AI — The Coordinator

**Agent ID:** `cmd-mother-ai`  
**Model:** claude-sonnet-4-5  
**Enable Sleep-time:** ✅ Yes  
**Tags:** `command`, `orchestrator`, `governance`

### Memory Blocks

| Block | Purpose | Size |
|-------|---------|------|
| `zeus_profile` | Identity, preferences, 4 pillars | 2000 tokens |
| `global_state` | Current date, active context, energy level | 500 tokens |
| `approval_queue` | Pending human-in-the-loop items | 1000 tokens |
| `pillar_status` | HP/LE/GE/IT daily summaries | 1500 tokens |

### Tools
- `route_to_pillar` - Determine which pillar handles request
- `request_approval` - Queue for human-in-the-loop
- `broadcast_update` - Notify all pillar conductors
- `sync_memory_blocks` - Coordinate shared memory

### Responsibilities
1. Receive all external messages (Telegram, Todoist, etc.)
2. Route to appropriate pillar conductor
3. Manage human-in-the-loop approvals
4. Cross-pillar coordination
5. Global state management

---

## Protocol Guardian

**Agent ID:** `cmd-protocol-guardian`  
**Model:** claude-haiku-3  
**Tags:** `command`, `protocols`, `validation`

### Memory Blocks
| Block | Purpose |
|-------|---------|
| `active_protocols` | Currently executing protocols |
| `protocol_index` | Quick lookup of all protocols by OMCS coordinate |

### Tools
- `execute_protocol` - Run a protocol by OMCS ID
- `validate_protocol` - Check protocol prerequisites
- `log_protocol_execution` - Record protocol runs

### Responsibilities
1. Execute protocols on command
2. Validate prerequisites before execution
3. Track protocol execution history
4. Suggest protocol improvements

---

## Approval Manager

**Agent ID:** `cmd-approval-manager`  
**Model:** claude-haiku-3  
**Tags:** `command`, `hitl`, `approvals`

### Memory Blocks
| Block | Purpose |
|-------|---------|
| `pending_approvals` | Queue of items awaiting human approval |
| `approval_history` | Past decisions for learning |
| `approval_rules` | Auto-approval thresholds |

### Tools
- `request_human_approval` - Create approval request
- `check_approval_status` - Query pending items
- `auto_approve_check` - Check if meets auto-approval criteria

### Human-in-the-Loop Flow
```
Agent Request → Check Auto-Approve Rules → 
  ├── Auto-Approved → Execute
  └── Requires Human → Queue → Notify (Todoist/Telegram) → 
      └── Human Response → Execute or Reject
```

---

## Memory Curator

**Agent ID:** `cmd-memory-curator`  
**Model:** claude-sonnet-4-5 (sleep-time only)  
**Enable Sleep-time:** ✅ Yes  
**Tags:** `command`, `memory`, `background`

### Memory Blocks
| Block | Purpose |
|-------|---------|
| `memory_health` | Memory block usage stats |
| `consolidation_queue` | Items needing archival |

### Tools
- `consolidate_archival` - Move old data to archival memory
- `reorganize_memory` - Optimize memory block structure
- `generate_insights` - Extract patterns from history

### Responsibilities
1. Background memory optimization (sleep-time)
2. Archival management
3. Cross-agent memory sync
4. Pattern extraction from conversations

---

# 💪 HP — HOLISTIC PERFORMANCE (10-19)

## HP-Conductor

**Agent ID:** `hp-conductor`  
**Model:** claude-sonnet-4-5  
**Tags:** `hp`, `conductor`, `health`

### Memory Blocks (Shared with HP agents)
| Block | Purpose |
|-------|---------|
| `hp_daily_state` | Today's health metrics, energy, mood |
| `hp_goals_active` | Current HP goals (fitness, habits, etc.) |
| `hp_protocols_cache` | Quick-access HP protocols |

### Tools
- `route_to_hp_agent` - Delegate to specialist
- `hp_daily_briefing` - Generate HP status
- `update_hp_metrics` - Record health data

---

## HP Specialist Agents

### 1. Health Monitor
**Agent ID:** `hp-health-monitor`  
**OMCS:** HP.10  
**Tags:** `hp`, `health`, `medical`

**Responsibilities:**
- Track medical appointments
- Manage medication schedules
- Monitor health conditions
- Insurance/provider info

**Tools:**
- `log_health_event` - Record symptoms, appointments
- `check_medication_schedule` - Review meds

---

### 2. Fitness Coach
**Agent ID:** `hp-fitness-coach`  
**OMCS:** HP.20  
**Tags:** `hp`, `fitness`, `training`

**Responsibilities:**
- Workout programming
- Training log tracking
- Recovery monitoring
- PR tracking

**Tools:**
- `log_workout` - Record training session
- `generate_workout` - Create workout plan
- `check_recovery_status` - Assess readiness

---

### 3. Nutrition Guide
**Agent ID:** `hp-nutrition-guide`  
**OMCS:** HP.50  
**Tags:** `hp`, `nutrition`, `diet`

**Responsibilities:**
- Meal planning
- Supplement management
- Food logging
- Dietary tracking

**Tools:**
- `log_meal` - Record food intake
- `generate_meal_plan` - Create weekly meals
- `check_supplement_schedule` - Review stack

---

### 4. Sleep Optimizer
**Agent ID:** `hp-sleep-optimizer`  
**OMCS:** HP.40  
**Tags:** `hp`, `sleep`, `recovery`

**Responsibilities:**
- Sleep tracking analysis
- Sleep environment optimization
- Recovery protocols

**Tools:**
- `analyze_sleep_data` - Review sleep metrics
- `suggest_sleep_optimization` - Improvement recommendations

---

### 5. Mental Wellness
**Agent ID:** `hp-mental-wellness`  
**OMCS:** HP.30  
**Tags:** `hp`, `mental`, `mindfulness`

**Responsibilities:**
- Meditation tracking
- Stress management
- Therapy notes
- Mindfulness prompts

**Tools:**
- `log_meditation` - Record practice
- `stress_check_in` - Assess stress level
- `suggest_mindfulness` - Recommend practice

---

### 6. Habit Tracker
**Agent ID:** `hp-habit-tracker`  
**OMCS:** HP.70  
**Tags:** `hp`, `habits`, `routines`

**Responsibilities:**
- Daily habit tracking
- Routine execution
- Streak monitoring
- Habit formation

**Tools:**
- `log_habit` - Mark habit complete
- `check_streaks` - Review habit streaks
- `suggest_habit_stack` - Recommend habit bundles

---

### 7. Energy Manager
**Agent ID:** `hp-energy-manager`  
**OMCS:** HP.60  
**Tags:** `hp`, `energy`, `focus`

**Responsibilities:**
- Energy tracking throughout day
- Peak performance time identification
- Focus block scheduling

**Tools:**
- `log_energy_level` - Record current energy
- `identify_peak_times` - Analyze best focus hours
- `schedule_focus_block` - Block focused time

---

# 🏠 LE — LEGACY EVOLUTION (20-29)

## LE-Conductor

**Agent ID:** `le-conductor`  
**Model:** claude-sonnet-4-5  
**Tags:** `le`, `conductor`, `legacy`

### Memory Blocks (Shared with LE agents)
| Block | Purpose |
|-------|---------|
| `le_family_state` | Family updates, events, needs |
| `le_growth_tracker` | Personal development progress |
| `le_relationships_cache` | Key relationship context |

---

## LE Specialist Agents

### 1. Identity Guardian
**Agent ID:** `le-identity-guardian`  
**OMCS:** LE.10  
**Tags:** `le`, `identity`, `values`

**Responsibilities:**
- Values clarification
- Life philosophy documentation
- Identity document management
- Personal history

**Tools:**
- `review_values` - Check alignment with values
- `update_identity_doc` - Maintain identity records

---

### 2. Relationship Manager
**Agent ID:** `le-relationship-manager`  
**OMCS:** LE.20  
**Tags:** `le`, `relationships`, `family`

**Responsibilities:**
- Family relationship tracking
- Friend connection maintenance
- Relationship check-ins
- Important dates

**Tools:**
- `log_interaction` - Record relationship touchpoint
- `suggest_connection` - Recommend reaching out
- `check_important_dates` - Review birthdays, anniversaries

---

### 3. Growth Coach
**Agent ID:** `le-growth-coach`  
**OMCS:** LE.30  
**Tags:** `le`, `growth`, `learning`

**Responsibilities:**
- Learning goal tracking
- Book/course management
- Skill development
- Journal prompts

**Tools:**
- `log_learning` - Record learning session
- `suggest_next_read` - Recommend books
- `generate_journal_prompt` - Create reflection prompt

---

### 4. Creativity Director
**Agent ID:** `le-creativity-director`  
**OMCS:** LE.40  
**Tags:** `le`, `creativity`, `expression`

**Responsibilities:**
- Creative project tracking
- Writing management
- Art/music organization
- Creative habit support

**Tools:**
- `log_creative_work` - Record creative session
- `suggest_creative_prompt` - Generate inspiration

---

### 5. Life Admin
**Agent ID:** `le-life-admin`  
**OMCS:** LE.50  
**Tags:** `le`, `admin`, `home`

**Responsibilities:**
- Home management
- Vehicle maintenance
- Insurance tracking
- Personal legal docs

**Tools:**
- `log_home_task` - Record maintenance
- `check_renewals` - Review upcoming renewals
- `manage_documents` - Organize important docs

---

### 6. Experience Curator
**Agent ID:** `le-experience-curator`  
**OMCS:** LE.60  
**Tags:** `le`, `travel`, `experiences`

**Responsibilities:**
- Travel planning
- Event tracking
- Bucket list management
- Memory documentation

**Tools:**
- `plan_trip` - Create travel itinerary
- `log_experience` - Record meaningful experience
- `update_bucket_list` - Manage life goals

---

### 7. Legacy Architect
**Agent ID:** `le-legacy-architect`  
**OMCS:** LE.70  
**Tags:** `le`, `legacy`, `impact`

**Responsibilities:**
- Legacy letter writing
- Memoir organization
- Charitable giving tracking
- Long-term impact planning

**Tools:**
- `draft_legacy_letter` - Create message for future
- `log_giving` - Track charitable contributions
- `review_impact` - Assess legacy progress

---

# 💼 GE — GENERATIONAL ENTREPRENEURSHIP (30-39)

## GE-Conductor

**Agent ID:** `ge-conductor`  
**Model:** claude-sonnet-4-5  
**Tags:** `ge`, `conductor`, `business`

### Memory Blocks (Shared with GE agents)
| Block | Purpose |
|-------|---------|
| `ge_ventures_status` | Active business status |
| `ge_revenue_tracker` | Financial overview |
| `ge_priorities_queue` | Business priorities |

---

## GE Specialist Agents

### 1. Venture Lead
**Agent ID:** `ge-venture-lead`  
**OMCS:** GE.10  
**Tags:** `ge`, `ventures`, `strategy`

**Responsibilities:**
- Business entity management
- Strategic planning
- Venture health monitoring
- New opportunity evaluation

**Tools:**
- `assess_venture_health` - Evaluate business status
- `log_strategic_decision` - Record major decisions

---

### 2. Product Manager
**Agent ID:** `ge-product-manager`  
**OMCS:** GE.20  
**Tags:** `ge`, `products`, `services`

**Responsibilities:**
- Product development tracking
- Service offering management
- Launch planning
- Product roadmap

**Tools:**
- `update_product_status` - Track product progress
- `plan_launch` - Create launch checklist

---

### 3. Marketing Strategist
**Agent ID:** `ge-marketing-strategist`  
**OMCS:** GE.30  
**Tags:** `ge`, `marketing`, `sales`

**Responsibilities:**
- Marketing campaign management
- Content strategy
- Funnel optimization
- Brand management

**Tools:**
- `plan_campaign` - Create marketing campaign
- `analyze_funnel` - Review conversion metrics
- `generate_content_ideas` - Brainstorm content

---

### 4. Community Manager
**Agent ID:** `ge-community-manager`  
**OMCS:** GE.40  
**Tags:** `ge`, `community`, `customers`

**Responsibilities:**
- Customer relationship management
- Community engagement
- Testimonial collection
- Customer success

**Tools:**
- `log_customer_interaction` - Record touchpoint
- `identify_advocates` - Find potential testimonials
- `check_community_health` - Assess engagement

---

### 5. Operations Manager
**Agent ID:** `ge-operations-manager`  
**OMCS:** GE.50  
**Tags:** `ge`, `operations`, `execution`

**Responsibilities:**
- Process management
- Team coordination
- Tool/subscription management
- SOP creation

**Tools:**
- `create_sop` - Generate standard procedure
- `assign_task` - Delegate to team
- `review_processes` - Audit operations

---

### 6. Revenue Accountant
**Agent ID:** `ge-revenue-accountant`  
**OMCS:** GE.60  
**Tags:** `ge`, `finance`, `accounting`

**Responsibilities:**
- Revenue tracking
- Expense management
- Budget forecasting
- Tax preparation

**Tools:**
- `log_revenue` - Record income
- `log_expense` - Record cost
- `generate_financial_report` - Create summary

---

### 7. Investment Analyst
**Agent ID:** `ge-investment-analyst`  
**OMCS:** GE.70  
**Tags:** `ge`, `investments`, `wealth`

**Responsibilities:**
- Investment portfolio tracking
- Real estate management
- Wealth building strategy
- Retirement planning

**Tools:**
- `analyze_portfolio` - Review investments
- `log_investment` - Record new investment
- `project_wealth` - Forecast net worth

---

### 8. Compliance Officer
**Agent ID:** `ge-compliance-officer`  
**OMCS:** GE.80  
**Tags:** `ge`, `legal`, `compliance`

**Responsibilities:**
- Contract management
- IP protection
- Regulatory compliance
- Legal document organization

**Tools:**
- `track_contract` - Manage contracts
- `check_compliance` - Review requirements
- `log_ip` - Record intellectual property

---

### 9. Head of Education (Legacy Acadacom)
**Agent ID:** `ge-head-education`  
**Entity:** Legacy Acadacom  
**Tags:** `ge`, `education`, `legacy-acadacom`

**Responsibilities:**
- Curriculum development
- Course content management
- Student success tracking
- Educational content creation

**Tools:**
- `create_course_module` - Design course content
- `track_student_progress` - Monitor learners

---

### 10. Churn Manager (Legacy Acadacom)
**Agent ID:** `ge-churn-manager`  
**Entity:** Legacy Acadacom  
**Tags:** `ge`, `retention`, `legacy-acadacom`

**Responsibilities:**
- Customer retention
- Churn analysis
- Win-back campaigns
- Satisfaction monitoring

**Tools:**
- `identify_at_risk` - Find churn candidates
- `create_winback_campaign` - Design retention offer
- `analyze_churn_reasons` - Review cancellation causes

---

# 🔧 IT — INNOVATIVE TECHNOLOGY (40-49)

## IT-Conductor

**Agent ID:** `it-conductor`  
**Model:** claude-sonnet-4-5  
**Tags:** `it`, `conductor`, `technology`

### Memory Blocks (Shared with Grid agents)
| Block | Purpose |
|-------|---------|
| `it_infrastructure_status` | System health overview |
| `it_active_projects` | Development projects |
| `it_automation_registry` | Active automations |

---

## IT Specialist Agents

### 1. Infrastructure Admin
**Agent ID:** `it-infrastructure-admin`  
**OMCS:** IT.10  
**Tags:** `it`, `infrastructure`, `devops`

**Responsibilities:**
- Cloud service management
- Server administration
- Domain/DNS management
- Backup systems

**Tools:**
- `check_system_health` - Monitor infrastructure
- `manage_dns` - DNS record management
- `trigger_backup` - Initiate backup

---

### 2. AI Orchestrator
**Agent ID:** `it-ai-orchestrator`  
**OMCS:** IT.20  
**Tags:** `it`, `ai`, `agents`

**Responsibilities:**
- Letta agent management
- LLM configuration
- Training data organization
- Agent performance monitoring

**Tools:**
- `list_agents` - Get all Letta agents
- `modify_agent` - Update agent config
- `analyze_agent_performance` - Review metrics

---

### 3. Automation Engineer
**Agent ID:** `it-automation-engineer`  
**OMCS:** IT.30  
**Tags:** `it`, `automation`, `n8n`

**Responsibilities:**
- N8N workflow management
- Script automation
- Integration maintenance
- Cron job management

**Tools:**
- `create_automation` - Build new workflow
- `check_automation_status` - Monitor automations
- `trigger_workflow` - Execute workflow

---

### 4. Knowledge Architect
**Agent ID:** `it-knowledge-architect`  
**OMCS:** IT.40  
**Tags:** `it`, `obsidian`, `pkm`

**Responsibilities:**
- Vault structure maintenance
- Template management
- Plugin configuration
- Note organization

**Tools:**
- `search_vault` - Semantic search across vault
- `create_note` - Create new note
- `update_template` - Modify templates

---

### 5. Developer Lead
**Agent ID:** `it-developer-lead`  
**OMCS:** IT.50  
**Tags:** `it`, `development`, `code`

**Responsibilities:**
- Code project management
- Repository organization
- Development environment
- Technical debt tracking

**Tools:**
- `create_repo` - Initialize repository
- `review_code` - Code review assistance
- `track_tech_debt` - Log technical debt

---

### 6. Hardware Ops
**Agent ID:** `it-hardware-ops`  
**OMCS:** IT.60  
**Tags:** `it`, `hardware`, `devices`

**Responsibilities:**
- Device inventory
- Smart home management
- Peripheral tracking
- Hardware maintenance

**Tools:**
- `inventory_devices` - List all hardware
- `check_smart_home` - Home Assistant status
- `schedule_maintenance` - Plan hardware upkeep

---

### 7. Security Officer
**Agent ID:** `it-security-officer`  
**OMCS:** IT.70  
**Tags:** `it`, `security`, `cyber`

**Responsibilities:**
- Password/auth management
- Security audits
- Encryption key management
- Threat monitoring

**Tools:**
- `audit_security` - Run security check
- `rotate_credentials` - Update passwords/keys
- `check_threats` - Review security alerts

---

### 8. GitHub Sync Agent
**Agent ID:** `it-github-sync`  
**Tags:** `it`, `github`, `sync`

**Responsibilities:**
- Vault ↔ GitHub synchronization
- Change tracking
- Commit attribution
- Version control

**Tools:**
- `commit_changes` - Push vault changes
- `pull_updates` - Get latest from GitHub
- `check_sync_status` - Verify sync state

---

### 9. MCP Manager
**Agent ID:** `it-mcp-manager`  
**Tags:** `it`, `mcp`, `tools`

**Responsibilities:**
- MCP server management
- Tool registration
- Integration health
- New tool discovery

**Tools:**
- `list_mcp_servers` - Get all MCP connections
- `add_mcp_tool` - Register new MCP tool
- `check_mcp_health` - Verify MCP status

---

# 🔗 SHARED MEMORY ARCHITECTURE

## Cross-Agent Memory Blocks

All pillar conductors share these blocks with Mother AI:

```
┌─────────────────────────────────────────────────────────────────────┐
│                    SHARED MEMORY BLOCKS                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  │
│  │  zeus_profile    │  │  global_state    │  │ cross_pillar_ctx │  │
│  │                  │  │                  │  │                  │  │
│  │  • Identity      │  │  • Current date  │  │  • Active proj   │  │
│  │  • 4 Pillars     │  │  • Energy level  │  │  • Priorities    │  │
│  │  • Preferences   │  │  • Focus block   │  │  • Dependencies  │  │
│  │  • Comm style    │  │  • Location      │  │  • Blocked items │  │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘  │
│                                                                      │
│  Shared by: Mother AI + HP/LE/GE/IT Conductors                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Pillar-Specific Memory Blocks

Each pillar has dedicated blocks shared among its agents:

```
┌──────────────────────────────────────────────────────────────────────┐
│  HP_daily_state   │  LE_family_state  │  GE_business_state  │ IT_system_state │
│                   │                   │                     │                 │
│  • Today metrics  │  • Family needs   │  • Revenue status   │  • Infra health │
│  • Energy level   │  • Relationships  │  • Active deals     │  • Active jobs  │
│  • Habits done    │  • Growth goals   │  • Team status      │  • Automations  │
└──────────────────────────────────────────────────────────────────────┘
```

---

# 🔄 TODOIST INTEGRATION

## Project Structure (Matching Agents)

```
📁 One Mind (Workspace)
├── 📁 00 Command
│   ├── Inbox
│   ├── Approvals
│   └── Weekly Review
│
├── 📁 10-19 HP
│   ├── HP.10 Health
│   ├── HP.20 Fitness
│   ├── HP.30 Mental
│   ├── HP.40 Sleep
│   ├── HP.50 Nutrition
│   ├── HP.60 Energy
│   └── HP.70 Habits
│
├── 📁 20-29 LE
│   ├── LE.10 Identity
│   ├── LE.20 Relationships
│   ├── LE.30 Growth
│   ├── LE.40 Creativity
│   ├── LE.50 Life Admin
│   ├── LE.60 Experiences
│   └── LE.70 Legacy
│
├── 📁 30-39 GE
│   ├── GE.10 Ventures
│   ├── GE.20 Products
│   ├── GE.30 Marketing
│   ├── GE.40 Community
│   ├── GE.50 Operations
│   ├── GE.60 Finance
│   ├── GE.70 Investments
│   └── GE.80 Legal
│
└── 📁 40-49 IT
    ├── IT.10 Infrastructure
    ├── IT.20 AI & Agents
    ├── IT.30 Automation
    ├── IT.40 Knowledge
    ├── IT.50 Development
    ├── IT.60 Hardware
    └── IT.70 Security
```

## Labels (Agent Tags)

| Label | Color | Purpose |
|-------|-------|---------|
| `@mother-ai` | Purple | Mother AI assigned |
| `@hp-conductor` | Green | HP pillar lead |
| `@le-conductor` | Blue | LE pillar lead |
| `@ge-conductor` | Orange | GE pillar lead |
| `@it-conductor` | Red | IT pillar lead |
| `@approval-needed` | Yellow | Human-in-the-loop |
| `@blocked` | Gray | Waiting on dependency |
| `@urgent` | Red | Priority escalation |

---

# ⏰ N8N CRON WORKFLOWS

## Daily Workflows

### 1. Morning Briefing (5:30 AM)
```yaml
Trigger: Cron 5:30 AM
Action: POST to Letta (Mother AI)
Message: "TRIGGER: morning_briefing"
```

### 2. Evening Review (9:00 PM)
```yaml
Trigger: Cron 9:00 PM
Action: POST to Letta (Mother AI)
Message: "TRIGGER: evening_review"
```

### 3. Agent Health Check (Every 4 hours)
```yaml
Trigger: Cron 0 */4 * * *
Action: POST to Letta (IT-AI-Orchestrator)
Message: "TRIGGER: agent_health_check"
```

## Weekly Workflows

### 4. Weekly Review Prep (Sunday 6:00 PM)
```yaml
Trigger: Cron 0 18 * * 0
Action: POST to Letta (Mother AI)
Message: "TRIGGER: weekly_review_prep"
```

### 5. Memory Consolidation (Saturday 3:00 AM)
```yaml
Trigger: Cron 0 3 * * 6
Action: POST to Letta (Memory Curator)
Message: "TRIGGER: memory_consolidation"
```

---

# 🙋 HUMAN-IN-THE-LOOP SYSTEM

## Approval Categories

| Category | Auto-Approve | Requires Human |
|----------|--------------|----------------|
| **Read Operations** | ✅ Always | Never |
| **Minor Updates** | ✅ < $50 impact | ≥ $50 impact |
| **Task Creation** | ✅ Personal tasks | Team tasks |
| **Financial** | ✅ < $100 | ≥ $100 |
| **Publishing** | Never | ✅ Always |
| **Deletion** | Never | ✅ Always |
| **External Comms** | Never | ✅ Always |

## Approval Flow

```
Agent wants to execute action
         │
         ▼
┌─────────────────────────────┐
│  Check auto-approve rules    │
└─────────────────────────────┘
         │
    ┌────┴────┐
    │         │
  Auto     Requires
 Approve   Approval
    │         │
    ▼         ▼
Execute   Queue in Todoist
          (@approval-needed)
               │
               ▼
         Notify Zeus
         (Telegram)
               │
               ▼
         Zeus responds
         (Approve/Reject)
               │
          ┌────┴────┐
          │         │
       Approved   Rejected
          │         │
          ▼         ▼
       Execute   Log reason
                 Archive
```

## Implementation (Letta API)

```python
# Agent requests approval
response = client.agents.messages.create(
    agent_id=agent_id,
    messages=[{
        "role": "user",
        "content": "REQUEST_APPROVAL: Send marketing email to 500 subscribers"
    }]
)

# Get tool_call_id from pending approval
tool_call_id = response.messages[-1].tool_call.tool_call_id

# Human approves (via Todoist webhook → n8n → Letta)
client.agents.messages.create(
    agent_id=agent_id,
    messages=[{
        "type": "approval",
        "approvals": [{
            "type": "approval",
            "approve": True,  # or False to reject
            "tool_call_id": tool_call_id,
        }]
    }]
)
```

---

# 🛠️ TOOL TYPES

## 1. MCP Tools (External Services)

| MCP Server | Tools | Agent Access |
|------------|-------|--------------|
| `todoist-mcp` | create_task, complete_task, get_tasks | All agents |
| `github-mcp` | commit, pull, search_code | Grid agents + Mother AI |
| `fetch-mcp` | fetch_url, extract_content | All agents |
| `obsidian-mcp` | search, read, write | All agents |

## 2. Python Tools (Custom Logic)

```python
# Example: Custom tool for HP-Fitness-Coach
def log_workout(
    self,
    workout_type: str,
    duration_minutes: int,
    exercises: list[str],
    notes: str = ""
) -> str:
    """
    Log a completed workout session to HP tracking.
    
    Args:
        workout_type: Type of workout (strength, cardio, mobility)
        duration_minutes: How long the workout lasted
        exercises: List of exercises performed
        notes: Additional notes about the session
        
    Returns:
        Confirmation with workout ID
    """
    workout_id = f"HP.20.LOG.{datetime.now().strftime('%Y%m%d%H%M')}"
    
    # Save to archival memory
    self.archival_insert(json.dumps({
        "id": workout_id,
        "type": workout_type,
        "duration": duration_minutes,
        "exercises": exercises,
        "notes": notes,
        "timestamp": datetime.now().isoformat()
    }))
    
    # Update HP daily state
    self.memory_replace("hp_daily_state", {
        "last_workout": workout_id,
        "workout_done": True
    })
    
    return f"Workout logged: {workout_id}"
```

## 3. API Tools (Direct Integrations)

For services without MCP:
- Google Calendar API
- Twilio SMS
- OpenAI (for specific models)
- Custom webhooks

---

# 🚀 IMPLEMENTATION PHASES

## Phase 1: Foundation (Week 1)
- [ ] Create Mother AI agent in Letta Cloud
- [ ] Set up shared memory blocks
- [ ] Connect Todoist MCP
- [ ] Connect GitHub MCP
- [ ] Create 4 pillar conductor agents

## Phase 2: HP Pillar (Week 2)
- [ ] Create HP specialist agents (7)
- [ ] Configure HP memory blocks
- [ ] Build HP custom tools
- [ ] Test HP workflows

## Phase 3: LE Pillar (Week 2-3)
- [ ] Create LE specialist agents (7)
- [ ] Configure LE memory blocks
- [ ] Build LE custom tools
- [ ] Test LE workflows

## Phase 4: GE Pillar (Week 3)
- [ ] Create GE specialist agents (10)
- [ ] Configure GE memory blocks
- [ ] Build GE custom tools
- [ ] Test GE workflows

## Phase 5: IT Pillar (Week 3-4)
- [ ] Create IT specialist agents (9)
- [ ] Configure IT memory blocks
- [ ] Build IT custom tools
- [ ] Test IT workflows

## Phase 6: Integration (Week 4)
- [ ] Configure N8N cron workflows
- [ ] Set up human-in-the-loop
- [ ] Test cross-domain routing
- [ ] Enable sleep-time compute
- [ ] Final testing

---

# 📚 Quick Reference

## Agent Naming Convention

```
[layer]-[role]
```

Examples:
- `cmd-mother-ai` — Command layer, Mother AI
- `hp-fitness-coach` — HP pillar, Fitness Coach
- `ge-revenue-accountant` — GE pillar, Revenue Accountant
- `it-automation-engineer` — IT pillar, Automation Engineer

## Memory Block Naming

```
[pillar]_[purpose]
```

Examples:
- `zeus_profile` — Global identity
- `hp_daily_state` — HP daily metrics
- `ge_ventures_status` — GE business health

## Tool Naming

```
[action]_[noun]
```

Examples:
- `log_workout` — Record workout
- `create_task` — Make new task
- `analyze_portfolio` — Review investments

---

# 🔗 MCP Server Configuration

## Letta MCP Server (oculairmedia)

For automated agent management:

```json
{
  "name": "letta-mcp",
  "type": "sse",
  "endpoint": "http://your-server:3001/mcp",
  "env": {
    "LETTA_BASE_URL": "https://api.letta.com",
    "LETTA_PASSWORD": "your-api-key"
  }
}
```

## N8N Letta Node

Install in N8N:
```bash
npm install @letta-ai/n8n-nodes-letta
```

Configuration:
- API Token: Letta Cloud API key
- Base URL: `https://api.letta.com`

---

*One Mind Agent Architecture v2.0 © 2025 Zeus Delacruz*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
